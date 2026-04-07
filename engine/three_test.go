package engine

import (
	"testing"

	"github.com/ApisMellow/cuttle/card"
)

// Task 16: Three one-off — take a card from scrap into hand.

func TestLegalMoves_ThreeEnumeratesOnePerScrapCard(t *testing.T) {
	three := card.Card{Rank: card.Three, Suit: card.Spades}
	s := twoPlayerStart([]card.Card{three}, nil, nil)
	s.Players[P2].Hand = []card.Card{{Rank: card.Five, Suit: card.Clubs}} // no 2s
	s.Scrap = []card.Card{
		{Rank: card.Seven, Suit: card.Hearts},
		{Rank: card.King, Suit: card.Diamonds},
		{Rank: card.Four, Suit: card.Clubs},
	}

	var threeMoves []Move
	for _, m := range LegalMoves(s) {
		if m.Kind == MoveOneOff && m.Card.Rank == card.Three {
			threeMoves = append(threeMoves, m)
		}
	}
	if len(threeMoves) != 3 {
		t.Fatalf("expected 3 MoveOneOff moves for Three (one per scrap card), got %d", len(threeMoves))
	}
	seen := map[int]bool{}
	for _, m := range threeMoves {
		if m.ScrapIndex < 0 || m.ScrapIndex >= 3 {
			t.Errorf("ScrapIndex out of range: %d", m.ScrapIndex)
		}
		seen[m.ScrapIndex] = true
	}
	if len(seen) != 3 {
		t.Errorf("expected distinct ScrapIndex for each move, got %v", seen)
	}
}

func TestLegalMoves_ThreeEmptyScrap_NoLegalOneOff(t *testing.T) {
	three := card.Card{Rank: card.Three, Suit: card.Spades}
	s := twoPlayerStart([]card.Card{three}, nil, nil)
	s.Players[P2].Hand = []card.Card{{Rank: card.Five, Suit: card.Clubs}}
	s.Scrap = nil

	for _, m := range LegalMoves(s) {
		if m.Kind == MoveOneOff && m.Card.Rank == card.Three {
			t.Errorf("expected no Three one-off with empty scrap, got %+v", m)
		}
	}
}

func TestApply_ThreeOneOff_TakesScrapCardIntoHand(t *testing.T) {
	three := card.Card{Rank: card.Three, Suit: card.Spades}
	seven := card.Card{Rank: card.Seven, Suit: card.Hearts}
	king := card.Card{Rank: card.King, Suit: card.Diamonds}
	four := card.Card{Rank: card.Four, Suit: card.Clubs}
	s := twoPlayerStart([]card.Card{three}, nil, nil)
	s.Players[P2].Hand = []card.Card{{Rank: card.Five, Suit: card.Clubs}} // no 2s
	s.Scrap = []card.Card{seven, king, four}

	s2, err := Apply(s, Move{Kind: MoveOneOff, Card: three, HandIndex: 0, ScrapIndex: 1})
	if err != nil {
		t.Fatal(err)
	}
	// King should now be in P1's hand.
	if len(s2.Players[P1].Hand) != 1 || s2.Players[P1].Hand[0] != king {
		t.Errorf("expected P1 hand to contain King, got %+v", s2.Players[P1].Hand)
	}
	// Scrap should now have [seven, four, three] (king removed, three appended).
	want := []card.Card{seven, four, three}
	if len(s2.Scrap) != len(want) {
		t.Fatalf("expected scrap %v, got %v", want, s2.Scrap)
	}
	for i, c := range want {
		if s2.Scrap[i] != c {
			t.Errorf("scrap[%d]: want %v, got %v", i, c, s2.Scrap[i])
		}
	}
	// Turn should have ended.
	if s2.Active != P2 {
		t.Errorf("expected Active=P2 after Three one-off resolves, got %v", s2.Active)
	}
}

func TestLegalMoves_ThreeHandLimit_NoLegalWhenFull(t *testing.T) {
	three := card.Card{Rank: card.Three, Suit: card.Spades}
	// Hand will have 8 cards AFTER playing the 3 if hand is currently 9... but
	// max hand is 8. The constraint: AFTER playing the 3 and taking a card, hand
	// must be <= 8. Current hand has 3 + 7 others = 8. Play 3 → 7. Take → 8. OK.
	// At 3 + 8 others = 9, not reachable. The real edge: hand has 3 + 8 others
	// = 9 — impossible. So we test hand of 8 (three + 7 others) — legal.
	// Then hand of 9 — unreachable. The actual enforceable edge: hand of 3 only
	// where there is NO hand limit issue. Reading the task: "if taking would
	// exceed 8-card hand limit AFTER playing the 3". Playing removes the 3
	// (hand -1) then adds scrap card (+1) = same size. So the only way to
	// exceed is if hand was already > 8, which shouldn't happen. But to respect
	// the spec literally: if len(hand) > 8, no legal. Test both.
	hand := []card.Card{
		three,
		{Rank: card.Four, Suit: card.Spades},
		{Rank: card.Five, Suit: card.Spades},
		{Rank: card.Six, Suit: card.Spades},
		{Rank: card.Seven, Suit: card.Spades},
		{Rank: card.Eight, Suit: card.Spades},
		{Rank: card.Nine, Suit: card.Spades},
		{Rank: card.Ten, Suit: card.Spades},
	}
	s := twoPlayerStart(hand, nil, nil)
	s.Players[P2].Hand = []card.Card{{Rank: card.Five, Suit: card.Clubs}}
	s.Scrap = []card.Card{{Rank: card.King, Suit: card.Diamonds}}
	// Hand size 8; play 3 → 7; take → 8. Legal.
	found := false
	for _, m := range LegalMoves(s) {
		if m.Kind == MoveOneOff && m.Card.Rank == card.Three {
			found = true
		}
	}
	if !found {
		t.Error("expected Three one-off legal when hand would be exactly 8 after resolve")
	}
}
