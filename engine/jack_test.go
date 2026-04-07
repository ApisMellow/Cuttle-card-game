package engine

import (
	"testing"

	"github.com/ApisMellow/cuttle/card"
)

// Task 22: Jack as permanent (steal) + Queen protection.

func hasJackPlay(moves []Move, owner PlayerID, idx int) bool {
	for _, m := range moves {
		if m.Kind != MovePlayPermanent || m.Card.Rank != card.Jack {
			continue
		}
		if m.JackTarget == nil {
			continue
		}
		if m.JackTarget.Owner == owner && m.JackTarget.Zone == ZonePoints && m.JackTarget.Index == idx {
			return true
		}
	}
	return false
}

func TestJack_LegalMoves_EnumeratesOpponentPoints(t *testing.T) {
	jack := card.Card{Rank: card.Jack, Suit: card.Spades}
	s := twoPlayerStart([]card.Card{jack}, nil, nil)
	s.Players[P2].Hand = []card.Card{{Rank: card.Five, Suit: card.Clubs}}
	s.Players[P2].Points = []PointEntry{
		{Card: card.Card{Rank: card.Seven, Suit: card.Hearts}, Owner: P2},
		{Card: card.Card{Rank: card.Four, Suit: card.Diamonds}, Owner: P2},
	}
	moves := LegalMoves(s)
	if !hasJackPlay(moves, P2, 0) || !hasJackPlay(moves, P2, 1) {
		t.Errorf("expected Jack plays for both P2 points, got %+v", moves)
	}
}

func TestJack_StealPointTransplants(t *testing.T) {
	jack := card.Card{Rank: card.Jack, Suit: card.Spades}
	seven := card.Card{Rank: card.Seven, Suit: card.Hearts}
	s := twoPlayerStart([]card.Card{jack}, nil, nil)
	s.Players[P2].Hand = []card.Card{{Rank: card.Five, Suit: card.Clubs}}
	s.Players[P2].Points = []PointEntry{{Card: seven, Owner: P2}}

	tgt := Target{Owner: P2, Zone: ZonePoints, Index: 0}
	s2, err := Apply(s, Move{Kind: MovePlayPermanent, Card: jack, HandIndex: 0, JackTarget: &tgt})
	if err != nil {
		t.Fatal(err)
	}
	if len(s2.Players[P2].Points) != 0 {
		t.Errorf("point should leave P2, got %+v", s2.Players[P2].Points)
	}
	if len(s2.Players[P1].Points) != 1 {
		t.Fatalf("point should be on P1, got %d", len(s2.Players[P1].Points))
	}
	pe := s2.Players[P1].Points[0]
	if pe.Card != seven || pe.Owner != P2 {
		t.Errorf("bad entry: %+v", pe)
	}
	if len(pe.JackStack) != 1 || pe.JackStack[0] != jack {
		t.Errorf("jack stack wrong: %+v", pe.JackStack)
	}
	if len(pe.JackOwners) != 1 || pe.JackOwners[0] != P1 {
		t.Errorf("jack owners wrong: %+v", pe.JackOwners)
	}
	if pe.Controller() != P1 {
		t.Error("P1 should control")
	}
	if s2.Active != P2 {
		t.Error("turn should advance")
	}
}

func TestJack_ChainStealAppendsToStack(t *testing.T) {
	// P1 already controls a P2-owned 7 via a Jack. P2 plays another Jack to re-steal.
	jackP1 := card.Card{Rank: card.Jack, Suit: card.Spades}
	jackP2 := card.Card{Rank: card.Jack, Suit: card.Hearts}
	seven := card.Card{Rank: card.Seven, Suit: card.Diamonds}
	s := twoPlayerStart(nil, []card.Card{jackP2}, nil)
	s.Players[P1].Points = []PointEntry{
		{Card: seven, Owner: P2, JackStack: []card.Card{jackP1}, JackOwners: []PlayerID{P1}},
	}
	s.Active = P2

	tgt := Target{Owner: P1, Zone: ZonePoints, Index: 0}
	s2, err := Apply(s, Move{Kind: MovePlayPermanent, Card: jackP2, HandIndex: 0, JackTarget: &tgt})
	if err != nil {
		t.Fatal(err)
	}
	if len(s2.Players[P1].Points) != 0 {
		t.Error("point should leave P1")
	}
	if len(s2.Players[P2].Points) != 1 {
		t.Fatalf("point should be on P2, got %d", len(s2.Players[P2].Points))
	}
	pe := s2.Players[P2].Points[0]
	if len(pe.JackStack) != 2 {
		t.Errorf("expected 2 jacks, got %d", len(pe.JackStack))
	}
	if pe.JackOwners[0] != P1 || pe.JackOwners[1] != P2 {
		t.Errorf("jack owners wrong: %+v", pe.JackOwners)
	}
	if pe.Owner != P2 {
		t.Error("original owner preserved")
	}
	if pe.Controller() != P2 {
		t.Error("P2 should control")
	}
}

func TestJack_QueenBlocksSteal(t *testing.T) {
	jack := card.Card{Rank: card.Jack, Suit: card.Spades}
	s := twoPlayerStart([]card.Card{jack}, nil, nil)
	s.Players[P2].Hand = []card.Card{{Rank: card.Five, Suit: card.Clubs}}
	s.Players[P2].Points = []PointEntry{{Card: card.Card{Rank: card.Seven, Suit: card.Hearts}, Owner: P2}}
	s.Players[P2].Permanents = []card.Card{{Rank: card.Queen, Suit: card.Diamonds}}
	moves := LegalMoves(s)
	if hasJackPlay(moves, P2, 0) {
		t.Error("Queen should protect P2 points from Jack")
	}
}

func TestJack_ScrapReturnsPointToOriginalOwner(t *testing.T) {
	// Already covered by two_test.go; add a direct version here.
	two := card.Card{Rank: card.Two, Suit: card.Spades}
	jack := card.Card{Rank: card.Jack, Suit: card.Clubs}
	seven := card.Card{Rank: card.Seven, Suit: card.Hearts}
	s := twoPlayerStart(nil, []card.Card{two}, nil)
	s.Players[P1].Points = []PointEntry{
		{Card: seven, Owner: P2, JackStack: []card.Card{jack}, JackOwners: []PlayerID{P1}},
	}
	s.Active = P2

	tgt := Target{Owner: P1, Zone: ZonePoints, Index: 0}
	s2, err := Apply(s, Move{Kind: MoveOneOff, Card: two, HandIndex: 0, Target: &tgt})
	if err != nil {
		t.Fatal(err)
	}
	if len(s2.Players[P2].Points) != 1 || s2.Players[P2].Points[0].Card != seven {
		t.Errorf("point should have returned to P2, got %+v", s2.Players[P2].Points)
	}
}

func TestJack_PlayWinsGame(t *testing.T) {
	jack := card.Card{Rank: card.Jack, Suit: card.Spades}
	s := twoPlayerStart([]card.Card{jack}, nil, nil)
	s.Players[P2].Hand = []card.Card{{Rank: card.Five, Suit: card.Clubs}}
	// P1 currently at 11. Stealing a 10 from P2 → 21 → win.
	s.Players[P1].Points = []PointEntry{
		{Card: card.Card{Rank: card.Seven, Suit: card.Hearts}, Owner: P1},
		{Card: card.Card{Rank: card.Four, Suit: card.Spades}, Owner: P1},
	}
	s.Players[P2].Points = []PointEntry{
		{Card: card.Card{Rank: card.Ten, Suit: card.Diamonds}, Owner: P2},
	}
	tgt := Target{Owner: P2, Zone: ZonePoints, Index: 0}
	s2, err := Apply(s, Move{Kind: MovePlayPermanent, Card: jack, HandIndex: 0, JackTarget: &tgt})
	if err != nil {
		t.Fatal(err)
	}
	if s2.Phase != PhaseGameOver || s2.Winner == nil || *s2.Winner != P1 {
		t.Errorf("expected P1 win, got phase=%v winner=%v", s2.Phase, s2.Winner)
	}
}

func TestJack_MultiJackChain(t *testing.T) {
	// Stack: P2-owned 10, jacked by P1, re-jacked by P2, re-jacked by P1.
	j1 := card.Card{Rank: card.Jack, Suit: card.Spades}
	j2 := card.Card{Rank: card.Jack, Suit: card.Hearts}
	j3 := card.Card{Rank: card.Jack, Suit: card.Diamonds}
	ten := card.Card{Rank: card.Ten, Suit: card.Clubs}
	s := twoPlayerStart([]card.Card{j3}, nil, nil)
	s.Players[P2].Hand = []card.Card{{Rank: card.Five, Suit: card.Clubs}}
	s.Players[P2].Points = []PointEntry{
		{Card: ten, Owner: P2, JackStack: []card.Card{j1, j2}, JackOwners: []PlayerID{P1, P2}},
	}
	tgt := Target{Owner: P2, Zone: ZonePoints, Index: 0}
	s2, err := Apply(s, Move{Kind: MovePlayPermanent, Card: j3, HandIndex: 0, JackTarget: &tgt})
	if err != nil {
		t.Fatal(err)
	}
	if len(s2.Players[P1].Points) != 1 {
		t.Fatalf("point should be on P1, got %d", len(s2.Players[P1].Points))
	}
	pe := s2.Players[P1].Points[0]
	if len(pe.JackStack) != 3 {
		t.Errorf("expected 3 jacks, got %d", len(pe.JackStack))
	}
	if pe.Owner != P2 || pe.Controller() != P1 {
		t.Errorf("owner/controller wrong: owner=%v controller=%v", pe.Owner, pe.Controller())
	}
}

func TestJack_QueenAllowsTargetingOpposingQueen_For9(t *testing.T) {
	// Nine can target a Queen even if opponent has a Queen (queen doesn't protect itself).
	nine := card.Card{Rank: card.Nine, Suit: card.Hearts}
	queen := card.Card{Rank: card.Queen, Suit: card.Diamonds}
	s := twoPlayerStart([]card.Card{nine}, nil, nil)
	s.Players[P2].Permanents = []card.Card{queen}
	moves := LegalMoves(s)
	found := false
	for _, m := range moves {
		if m.Kind == MoveOneOff && m.Card.Rank == card.Nine && m.Target != nil && m.Target.Zone == ZonePermanents && m.Target.Index == 0 {
			found = true
		}
	}
	if !found {
		t.Error("9 should be allowed to target opposing Queen")
	}
}

func TestNine_QueenProtectsOtherPermanents(t *testing.T) {
	nine := card.Card{Rank: card.Nine, Suit: card.Hearts}
	s := twoPlayerStart([]card.Card{nine}, nil, nil)
	s.Players[P2].Permanents = []card.Card{
		{Rank: card.Queen, Suit: card.Diamonds},
		{Rank: card.King, Suit: card.Spades},
	}
	moves := LegalMoves(s)
	for _, m := range moves {
		if m.Kind != MoveOneOff || m.Card.Rank != card.Nine || m.Target == nil {
			continue
		}
		if m.Target.Zone == ZonePermanents && m.Target.Index == 1 {
			t.Error("9 should not target King when Queen protects")
		}
	}
}

func TestNine_QueenProtectsPoints(t *testing.T) {
	nine := card.Card{Rank: card.Nine, Suit: card.Hearts}
	s := twoPlayerStart([]card.Card{nine}, nil, nil)
	s.Players[P2].Permanents = []card.Card{{Rank: card.Queen, Suit: card.Diamonds}}
	s.Players[P2].Points = []PointEntry{{Card: card.Card{Rank: card.Seven, Suit: card.Hearts}, Owner: P2}}
	moves := LegalMoves(s)
	for _, m := range moves {
		if m.Kind == MoveOneOff && m.Card.Rank == card.Nine && m.Target != nil && m.Target.Zone == ZonePoints {
			t.Error("9 should not target point when Queen protects")
		}
	}
}
