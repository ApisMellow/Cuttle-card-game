package engine

import (
	"testing"

	"github.com/ApisMellow/cuttle/card"
)

// Task 19: Nine one-off — bounce opponent field card and freeze.

func TestNine_BouncePlainPoint_ReturnsToOwnerAndFreezes(t *testing.T) {
	nine := card.Card{Rank: card.Nine, Suit: card.Hearts}
	seven := card.Card{Rank: card.Seven, Suit: card.Diamonds}
	s := twoPlayerStart([]card.Card{nine}, []card.Card{{Rank: card.King, Suit: card.Clubs}}, nil)
	s.Players[P2].Points = []PointEntry{{Card: seven, Owner: P2}}

	target := &Target{Owner: P2, Zone: ZonePoints, Index: 0}
	s2, err := Apply(s, Move{Kind: MoveOneOff, Card: nine, HandIndex: 0, Target: target})
	if err != nil {
		t.Fatal(err)
	}
	if len(s2.Players[P2].Points) != 0 {
		t.Errorf("expected P2 points empty, got %+v", s2.Players[P2].Points)
	}
	// The seven should be back in P2's hand.
	foundAt := -1
	for i, c := range s2.Players[P2].Hand {
		if c == seven {
			foundAt = i
			break
		}
	}
	if foundAt < 0 {
		t.Fatalf("expected seven back in P2 hand, got %+v", s2.Players[P2].Hand)
	}
	// That index should be frozen for P2.
	if !s2.Players[P2].FrozenIDs[foundAt] {
		t.Errorf("expected P2 hand index %d to be frozen, got %+v", foundAt, s2.Players[P2].FrozenIDs)
	}
	// The nine itself goes to scrap.
	if len(s2.Scrap) != 1 || s2.Scrap[0] != nine {
		t.Errorf("expected scrap=[nine], got %+v", s2.Scrap)
	}
	// Turn should have advanced to P2.
	if s2.Active != P2 {
		t.Errorf("expected Active=P2 after turn end, got %v", s2.Active)
	}
}

func TestNine_BouncePointUnderJack_ScrapsJacksReturnsToOriginalOwner(t *testing.T) {
	nine := card.Card{Rank: card.Nine, Suit: card.Hearts}
	seven := card.Card{Rank: card.Seven, Suit: card.Diamonds}
	jack := card.Card{Rank: card.Jack, Suit: card.Spades}
	// P2 currently controls a 7 that was originally P1's, stolen with a Jack.
	s := twoPlayerStart([]card.Card{nine}, nil, nil)
	s.Players[P2].Points = []PointEntry{{
		Card:       seven,
		Owner:      P1,
		JackStack:  []card.Card{jack},
		JackOwners: []PlayerID{P2},
	}}

	target := &Target{Owner: P2, Zone: ZonePoints, Index: 0}
	s2, err := Apply(s, Move{Kind: MoveOneOff, Card: nine, HandIndex: 0, Target: target})
	if err != nil {
		t.Fatal(err)
	}
	if len(s2.Players[P2].Points) != 0 {
		t.Errorf("expected P2 points empty, got %+v", s2.Players[P2].Points)
	}
	// Point returns to original owner P1's hand (P1 is the active/played-by player).
	foundAt := -1
	for i, c := range s2.Players[P1].Hand {
		if c == seven {
			foundAt = i
			break
		}
	}
	if foundAt < 0 {
		t.Fatalf("expected seven in P1 hand, got %+v", s2.Players[P1].Hand)
	}
	if !s2.Players[P1].FrozenIDs[foundAt] {
		t.Errorf("expected P1 hand index %d frozen, got %+v", foundAt, s2.Players[P1].FrozenIDs)
	}
	// Scrap: nine + jack.
	if len(s2.Scrap) != 2 {
		t.Errorf("expected scrap size 2 (nine+jack), got %+v", s2.Scrap)
	}
	sawJack, sawNine := false, false
	for _, c := range s2.Scrap {
		if c == jack {
			sawJack = true
		}
		if c == nine {
			sawNine = true
		}
	}
	if !sawJack || !sawNine {
		t.Errorf("expected both nine and jack in scrap, got %+v", s2.Scrap)
	}
}

func TestNine_BounceQueenPermanent_ReturnsAndFreezes(t *testing.T) {
	nine := card.Card{Rank: card.Nine, Suit: card.Hearts}
	queen := card.Card{Rank: card.Queen, Suit: card.Diamonds}
	s := twoPlayerStart([]card.Card{nine}, nil, nil)
	s.Players[P2].Permanents = []card.Card{queen}

	target := &Target{Owner: P2, Zone: ZonePermanents, Index: 0}
	s2, err := Apply(s, Move{Kind: MoveOneOff, Card: nine, HandIndex: 0, Target: target})
	if err != nil {
		t.Fatal(err)
	}
	if len(s2.Players[P2].Permanents) != 0 {
		t.Errorf("expected P2 permanents empty, got %+v", s2.Players[P2].Permanents)
	}
	foundAt := -1
	for i, c := range s2.Players[P2].Hand {
		if c == queen {
			foundAt = i
			break
		}
	}
	if foundAt < 0 {
		t.Fatalf("expected queen in P2 hand, got %+v", s2.Players[P2].Hand)
	}
	if !s2.Players[P2].FrozenIDs[foundAt] {
		t.Errorf("expected P2 frozen at index %d, got %+v", foundAt, s2.Players[P2].FrozenIDs)
	}
}

func TestNine_FrozenCardCannotBePlayedNextTurn(t *testing.T) {
	nine := card.Card{Rank: card.Nine, Suit: card.Hearts}
	seven := card.Card{Rank: card.Seven, Suit: card.Diamonds}
	// Give P2 the seven in their hand already along with another playable card.
	other := card.Card{Rank: card.Five, Suit: card.Clubs}
	_ = other
	s := twoPlayerStart([]card.Card{nine}, nil, nil)
	s.Players[P2].Points = []PointEntry{{Card: seven, Owner: P2}}

	target := &Target{Owner: P2, Zone: ZonePoints, Index: 0}
	s2, err := Apply(s, Move{Kind: MoveOneOff, Card: nine, HandIndex: 0, Target: target})
	if err != nil {
		t.Fatal(err)
	}
	// Now P2's turn. The seven should be in hand but frozen → cannot be played.
	if s2.Active != P2 {
		t.Fatalf("expected Active=P2, got %v", s2.Active)
	}
	// Find seven in P2 hand.
	sevenIdx := -1
	for i, c := range s2.Players[P2].Hand {
		if c == seven {
			sevenIdx = i
			break
		}
	}
	if sevenIdx < 0 {
		t.Fatal("seven not in P2 hand")
	}
	moves := LegalMoves(s2)
	for _, m := range moves {
		if m.HandIndex == sevenIdx && (m.Kind == MovePlayPoint || m.Kind == MoveScuttle || m.Kind == MovePlayPermanent || m.Kind == MoveOneOff) {
			if m.Card == seven {
				t.Errorf("frozen seven should not appear in LegalMoves, got %+v", m)
			}
		}
	}
	// After P2 takes any action (e.g., pass/draw), on their NEXT turn the freeze clears.
	// Simulate P2 pass → turn to P1 → P1 pass → back to P2, freeze should clear.
	s3, err := Apply(s2, Move{Kind: MovePass})
	if err != nil {
		t.Fatal(err)
	}
	// Now Active=P1. P1 pass.
	s4, err := Apply(s3, Move{Kind: MovePass})
	if err != nil {
		t.Fatal(err)
	}
	if s4.Active != P2 {
		t.Fatalf("expected Active=P2, got %v", s4.Active)
	}
	if s4.Players[P2].FrozenIDs[sevenIdx] {
		t.Error("freeze should have cleared by P2's next turn")
	}
}
