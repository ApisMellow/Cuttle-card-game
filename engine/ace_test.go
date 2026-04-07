package engine

import (
	"testing"

	"github.com/ApisMellow/cuttle/card"
)

// NOTE: Task 13 ignores counter-2 logic. All opponent hands in these tests
// contain ZERO 2s so resolution is immediate. Counters arrive in Task 14.

func TestLegalMoves_AceEnumeratesOneOff(t *testing.T) {
	hand := []card.Card{{Rank: card.Ace, Suit: card.Spades}}
	s := twoPlayerStart(hand, nil, nil)
	// Opponent hand: no 2s.
	s.Players[P2].Hand = []card.Card{{Rank: card.Five, Suit: card.Clubs}}
	found := false
	for _, m := range LegalMoves(s) {
		if m.Kind == MoveOneOff && m.Card.Rank == card.Ace {
			found = true
		}
	}
	if !found {
		t.Error("expected MoveOneOff for Ace in hand")
	}
}

func TestApply_AceOneOff_ClearsAllPoints(t *testing.T) {
	ace := card.Card{Rank: card.Ace, Suit: card.Spades}
	hand := []card.Card{ace}
	s := twoPlayerStart(hand, nil, nil)
	s.Players[P2].Hand = []card.Card{{Rank: card.Five, Suit: card.Clubs}} // no 2s
	s.Players[P1].Points = []PointEntry{
		{Card: card.Card{Rank: card.Seven, Suit: card.Hearts}, Owner: P1},
		{Card: card.Card{Rank: card.Five, Suit: card.Diamonds}, Owner: P1},
	}
	s.Players[P2].Points = []PointEntry{
		{Card: card.Card{Rank: card.Ten, Suit: card.Diamonds}, Owner: P2},
		{
			Card:       card.Card{Rank: card.Four, Suit: card.Spades},
			Owner:      P2,
			JackStack:  []card.Card{{Rank: card.Jack, Suit: card.Hearts}},
			JackOwners: []PlayerID{P1},
		},
	}

	s2, err := Apply(s, Move{Kind: MoveOneOff, Card: ace, HandIndex: 0})
	if err != nil {
		t.Fatal(err)
	}
	if len(s2.Players[P1].Points) != 0 {
		t.Error("P1 points should be cleared")
	}
	if len(s2.Players[P2].Points) != 0 {
		t.Error("P2 points should be cleared")
	}
	// Scrap should contain: ace + 7H + 5D + 10D + 4S + JH = 6 cards.
	if len(s2.Scrap) != 6 {
		t.Errorf("scrap should have 6 cards, got %d", len(s2.Scrap))
	}
	if len(s2.Players[P1].Hand) != 0 {
		t.Error("ace should leave hand")
	}
	if s2.Active != P2 {
		t.Error("turn should advance to P2")
	}
}

func TestApply_AceOneOff_NoPointsOnBoard(t *testing.T) {
	ace := card.Card{Rank: card.Ace, Suit: card.Clubs}
	s := twoPlayerStart([]card.Card{ace}, nil, nil)
	s.Players[P2].Hand = []card.Card{{Rank: card.King, Suit: card.Hearts}} // no 2s
	s2, err := Apply(s, Move{Kind: MoveOneOff, Card: ace, HandIndex: 0})
	if err != nil {
		t.Fatal(err)
	}
	if len(s2.Scrap) != 1 {
		t.Errorf("scrap should have just the ace, got %d", len(s2.Scrap))
	}
	if s2.Active != P2 {
		t.Error("turn should advance")
	}
}
