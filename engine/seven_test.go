package engine

import (
	"testing"

	"github.com/ApisMellow/cuttle/card"
)

// Task 21: Seven one-off — reveal top 2 deck cards, play one, return one.

func TestSeven_RevealsTwoAndEntersChoosing(t *testing.T) {
	seven := card.Card{Rank: card.Seven, Suit: card.Hearts}
	top1 := card.Card{Rank: card.Five, Suit: card.Clubs}
	top2 := card.Card{Rank: card.King, Suit: card.Spades}
	s := twoPlayerStart([]card.Card{seven}, nil, []card.Card{top1, top2, {Rank: card.Three, Suit: card.Diamonds}})

	s2, err := Apply(s, Move{Kind: MoveOneOff, Card: seven, HandIndex: 0})
	if err != nil {
		t.Fatal(err)
	}
	if s2.Phase != PhaseSevenChoosing {
		t.Fatalf("expected PhaseSevenChoosing, got %v", s2.Phase)
	}
	if s2.Active != P1 {
		t.Errorf("active should remain P1, got %v", s2.Active)
	}
	if s2.Pending == nil || len(s2.Pending.Revealed) != 2 {
		t.Fatalf("expected 2 revealed, got %+v", s2.Pending)
	}
	if s2.Pending.Revealed[0] != top1 || s2.Pending.Revealed[1] != top2 {
		t.Errorf("revealed order wrong: %+v", s2.Pending.Revealed)
	}
	if len(s2.Deck) != 1 {
		t.Errorf("expected deck len 1, got %d", len(s2.Deck))
	}
	// Seven went to scrap.
	if len(s2.Scrap) != 1 || s2.Scrap[0] != seven {
		t.Errorf("expected seven in scrap, got %+v", s2.Scrap)
	}
}

func TestSeven_OneCardDeck_RevealsOne(t *testing.T) {
	seven := card.Card{Rank: card.Seven, Suit: card.Hearts}
	top1 := card.Card{Rank: card.Five, Suit: card.Clubs}
	s := twoPlayerStart([]card.Card{seven}, nil, []card.Card{top1})
	s2, err := Apply(s, Move{Kind: MoveOneOff, Card: seven, HandIndex: 0})
	if err != nil {
		t.Fatal(err)
	}
	if s2.Phase != PhaseSevenChoosing {
		t.Fatalf("expected PhaseSevenChoosing, got %v", s2.Phase)
	}
	if len(s2.Pending.Revealed) != 1 {
		t.Fatalf("expected 1 revealed, got %+v", s2.Pending.Revealed)
	}
	if len(s2.Deck) != 0 {
		t.Errorf("expected empty deck, got %d", len(s2.Deck))
	}
}

func TestSeven_EmptyDeck_NotLegal(t *testing.T) {
	seven := card.Card{Rank: card.Seven, Suit: card.Hearts}
	s := twoPlayerStart([]card.Card{seven}, nil, nil)
	moves := LegalMoves(s)
	for _, m := range moves {
		if m.Kind == MoveOneOff && m.Card.Rank == card.Seven {
			t.Errorf("seven one-off should not be legal with empty deck")
		}
	}
}

func TestSeven_LegalMovesInChoosingWrapInnerMoves(t *testing.T) {
	seven := card.Card{Rank: card.Seven, Suit: card.Hearts}
	five := card.Card{Rank: card.Five, Suit: card.Clubs}
	king := card.Card{Rank: card.King, Suit: card.Spades}
	s := twoPlayerStart([]card.Card{seven}, nil, []card.Card{five, king, {Rank: card.Three, Suit: card.Diamonds}})
	s2, _ := Apply(s, Move{Kind: MoveOneOff, Card: seven, HandIndex: 0})

	moves := LegalMoves(s2)
	if len(moves) == 0 {
		t.Fatal("expected some SevenPick moves")
	}
	var sawFivePoint, sawFiveOneOff, sawKingPerm bool
	for _, m := range moves {
		if m.Kind != MoveSevenPick {
			t.Errorf("expected only MoveSevenPick, got %v", m.Kind)
			continue
		}
		if m.SubMove == nil {
			t.Errorf("missing submove")
			continue
		}
		if m.SubMove.Kind == MoveDraw || m.SubMove.Kind == MovePass {
			t.Errorf("sub-move should not be draw/pass, got %v", m.SubMove.Kind)
		}
		if m.Card == five && m.SubMove.Kind == MovePlayPoint {
			sawFivePoint = true
		}
		if m.Card == five && m.SubMove.Kind == MoveOneOff {
			sawFiveOneOff = true
		}
		if m.Card == king && m.SubMove.Kind == MovePlayPermanent {
			sawKingPerm = true
		}
	}
	if !sawFivePoint || !sawFiveOneOff || !sawKingPerm {
		t.Errorf("missing expected inner moves: fivePoint=%v fiveOneOff=%v kingPerm=%v", sawFivePoint, sawFiveOneOff, sawKingPerm)
	}
}

func TestSeven_PickPlayPointReturnsOther(t *testing.T) {
	seven := card.Card{Rank: card.Seven, Suit: card.Hearts}
	five := card.Card{Rank: card.Five, Suit: card.Clubs}
	king := card.Card{Rank: card.King, Suit: card.Spades}
	deckRest := card.Card{Rank: card.Three, Suit: card.Diamonds}
	s := twoPlayerStart([]card.Card{seven}, nil, []card.Card{five, king, deckRest})
	s2, _ := Apply(s, Move{Kind: MoveOneOff, Card: seven, HandIndex: 0})

	pick := Move{Kind: MoveSevenPick, Card: five, SubMove: &Move{Kind: MovePlayPoint, Card: five}}
	s3, err := Apply(s2, pick)
	if err != nil {
		t.Fatal(err)
	}
	if s3.Phase != PhaseNormal {
		t.Errorf("expected PhaseNormal, got %v", s3.Phase)
	}
	// Five played as a point card on P1.
	if len(s3.Players[P1].Points) != 1 || s3.Players[P1].Points[0].Card != five {
		t.Errorf("expected five on P1 points, got %+v", s3.Players[P1].Points)
	}
	// King returned to top of deck.
	if len(s3.Deck) != 2 || s3.Deck[0] != king || s3.Deck[1] != deckRest {
		t.Errorf("expected king on top of deck, got %+v", s3.Deck)
	}
	// Turn ended; P2 now active.
	if s3.Active != P2 {
		t.Errorf("expected P2 active after point play, got %v", s3.Active)
	}
	// Pending cleared.
	if s3.Pending != nil {
		t.Errorf("expected nil pending, got %+v", s3.Pending)
	}
}

func TestSeven_PickOneOffCanBeCountered(t *testing.T) {
	seven := card.Card{Rank: card.Seven, Suit: card.Hearts}
	ace := card.Card{Rank: card.Ace, Suit: card.Clubs}
	king := card.Card{Rank: card.King, Suit: card.Spades}
	two := card.Card{Rank: card.Two, Suit: card.Hearts}
	s := twoPlayerStart([]card.Card{seven}, []card.Card{two}, []card.Card{ace, king})
	// Note: opponent has a 2 but the 7 itself has a deck draw — opponent
	// will also be offered to counter the 7. For this test we want the 7
	// to resolve cleanly first, so give P1 no hand 2s; opponent still can
	// counter the 7 before reveal. Bypass by playing to a state that skips
	// the 7's counter window — easiest: start s directly in PhaseSevenChoosing.
	s2 := s
	s2.Phase = PhaseSevenChoosing
	s2.Pending = &PendingOneOff{PlayedBy: P1, Card: seven, Revealed: []card.Card{ace, king}}
	s2.Deck = nil
	s2.Scrap = []card.Card{seven}
	s2.Players[P1].Hand = nil

	pick := Move{Kind: MoveSevenPick, Card: ace, SubMove: &Move{Kind: MoveOneOff, Card: ace}}
	s3, err := Apply(s2, pick)
	if err != nil {
		t.Fatal(err)
	}
	// Because P2 has a 2, the inner Ace should enter PhaseAwaitingCounter.
	if s3.Phase != PhaseAwaitingCounter {
		t.Errorf("expected PhaseAwaitingCounter after picking ace, got %v", s3.Phase)
	}
	if s3.Active != P2 {
		t.Errorf("expected P2 active for counter, got %v", s3.Active)
	}
	// King was returned to top of deck before dispatch.
	if len(s3.Deck) != 1 || s3.Deck[0] != king {
		t.Errorf("expected king on deck, got %+v", s3.Deck)
	}
	if s3.Pending == nil || s3.Pending.Card != ace {
		t.Errorf("expected ace pending, got %+v", s3.Pending)
	}
}

func TestSeven_UnchosenReturnsToDeckTop(t *testing.T) {
	seven := card.Card{Rank: card.Seven, Suit: card.Hearts}
	five := card.Card{Rank: card.Five, Suit: card.Clubs}
	king := card.Card{Rank: card.King, Suit: card.Spades}
	deckRest := card.Card{Rank: card.Three, Suit: card.Diamonds}
	s := twoPlayerStart([]card.Card{seven}, nil, []card.Card{five, king, deckRest})
	s2, _ := Apply(s, Move{Kind: MoveOneOff, Card: seven, HandIndex: 0})

	// Pick the king as a permanent; five should return to deck top.
	pick := Move{Kind: MoveSevenPick, Card: king, SubMove: &Move{Kind: MovePlayPermanent, Card: king}}
	s3, err := Apply(s2, pick)
	if err != nil {
		t.Fatal(err)
	}
	if len(s3.Deck) != 2 || s3.Deck[0] != five || s3.Deck[1] != deckRest {
		t.Errorf("expected five on top of deck, got %+v", s3.Deck)
	}
	if len(s3.Players[P1].Permanents) != 1 || s3.Players[P1].Permanents[0] != king {
		t.Errorf("expected king permanent, got %+v", s3.Players[P1].Permanents)
	}
}
