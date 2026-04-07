package engine

import (
	"testing"

	"github.com/ApisMellow/cuttle/card"
)

// Task 15: Six one-off. Scraps all royals (J, Q, K) and glasses-8s on both
// sides. Jacks on point stacks are scrapped; underlying point returns to its
// original Owner first.

func TestLegalMoves_SixEnumeratesOneOff(t *testing.T) {
	six := card.Card{Rank: card.Six, Suit: card.Clubs}
	s := twoPlayerStart([]card.Card{six}, nil, nil)
	s.Players[P2].Hand = []card.Card{{Rank: card.Five, Suit: card.Clubs}} // no 2s
	found := false
	for _, m := range LegalMoves(s) {
		if m.Kind == MoveOneOff && m.Card.Rank == card.Six {
			found = true
		}
	}
	if !found {
		t.Error("expected MoveOneOff for Six in hand")
	}
}

func TestApply_SixOneOff_WipesRoyalsAndGlasses(t *testing.T) {
	six := card.Card{Rank: card.Six, Suit: card.Clubs}
	s := twoPlayerStart([]card.Card{six}, nil, nil)
	s.Players[P2].Hand = []card.Card{{Rank: card.Five, Suit: card.Clubs}} // no 2s
	s.Players[P1].Permanents = []card.Card{
		{Rank: card.Queen, Suit: card.Diamonds},
		{Rank: card.King, Suit: card.Spades},
		{Rank: card.Eight, Suit: card.Clubs}, // glasses
	}
	s.Players[P2].Permanents = []card.Card{
		{Rank: card.Queen, Suit: card.Hearts},
	}
	// Point cards should survive a Six.
	s.Players[P1].Points = []PointEntry{
		{Card: card.Card{Rank: card.Seven, Suit: card.Hearts}, Owner: P1},
	}

	s2, err := Apply(s, Move{Kind: MoveOneOff, Card: six, HandIndex: 0})
	if err != nil {
		t.Fatal(err)
	}
	if len(s2.Players[P1].Permanents) != 0 {
		t.Errorf("P1 permanents should be empty, got %v", s2.Players[P1].Permanents)
	}
	if len(s2.Players[P2].Permanents) != 0 {
		t.Errorf("P2 permanents should be empty, got %v", s2.Players[P2].Permanents)
	}
	// Points untouched.
	if len(s2.Players[P1].Points) != 1 {
		t.Errorf("P1 points should survive, got %d", len(s2.Players[P1].Points))
	}
	// Scrap: six + Q♦ + K♠ + 8♣ + Q♥ = 5 cards.
	if len(s2.Scrap) != 5 {
		t.Errorf("scrap should have 5 cards, got %d", len(s2.Scrap))
	}
	if s2.Active != P2 {
		t.Error("turn should advance to P2")
	}
}

func TestApply_SixOneOff_ScrapsJackReturnsPoint(t *testing.T) {
	six := card.Card{Rank: card.Six, Suit: card.Clubs}
	s := twoPlayerStart([]card.Card{six}, nil, nil)
	s.Players[P2].Hand = []card.Card{{Rank: card.Five, Suit: card.Clubs}} // no 2s
	// A 10 originally owned by P2, stolen by P1 via a Jack. Per the
	// Task-22 invariant the PointEntry transplants to the controller, so
	// we place it on P1's Points slice here.
	stolen := PointEntry{
		Card:       card.Card{Rank: card.Ten, Suit: card.Diamonds},
		Owner:      P2,
		JackStack:  []card.Card{{Rank: card.Jack, Suit: card.Hearts}},
		JackOwners: []PlayerID{P1},
	}
	s.Players[P1].Points = []PointEntry{stolen}

	s2, err := Apply(s, Move{Kind: MoveOneOff, Card: six, HandIndex: 0})
	if err != nil {
		t.Fatal(err)
	}
	if len(s2.Players[P1].Points) != 0 {
		t.Errorf("P1 should no longer control the stolen point, got %d", len(s2.Players[P1].Points))
	}
	if len(s2.Players[P2].Points) != 1 {
		t.Fatalf("P2 should get their 10 back, got %d", len(s2.Players[P2].Points))
	}
	back := s2.Players[P2].Points[0]
	if back.Card.Rank != card.Ten || back.Owner != P2 {
		t.Errorf("returned point wrong: %+v", back)
	}
	if len(back.JackStack) != 0 || len(back.JackOwners) != 0 {
		t.Error("jack stack should be cleared")
	}
	// Scrap: six + JH = 2 cards. Point card is NOT scrapped.
	if len(s2.Scrap) != 2 {
		t.Errorf("scrap should have 2 cards (six + jack), got %d", len(s2.Scrap))
	}
}

func TestApply_SixOneOff_JackReturnWinsForOwner(t *testing.T) {
	// P1 is active and plays a six. P2 has a stolen 10 sitting under a P1
	// Jack. When the Jack is scrapped the 10 returns to P2, but win checks
	// only fire for the active player, so no win.
	six := card.Card{Rank: card.Six, Suit: card.Clubs}
	s := twoPlayerStart([]card.Card{six}, nil, nil)
	s.Players[P2].Hand = []card.Card{{Rank: card.Five, Suit: card.Clubs}}
	// P1 already at 14 via points.
	s.Players[P1].Points = []PointEntry{
		{Card: card.Card{Rank: card.Ten, Suit: card.Hearts}, Owner: P1},
		{Card: card.Card{Rank: card.Four, Suit: card.Spades}, Owner: P1},
	}
	// Stolen point: originally P1's, currently held by P2 via a Jack. When
	// six scraps the Jack, the 7 returns to P1, pushing P1 to 21 → win.
	s.Players[P2].Points = []PointEntry{
		{
			Card:       card.Card{Rank: card.Seven, Suit: card.Hearts},
			Owner:      P1,
			JackStack:  []card.Card{{Rank: card.Jack, Suit: card.Clubs}},
			JackOwners: []PlayerID{P2},
		},
	}
	s2, err := Apply(s, Move{Kind: MoveOneOff, Card: six, HandIndex: 0})
	if err != nil {
		t.Fatal(err)
	}
	if s2.Phase != PhaseGameOver || s2.Winner == nil || *s2.Winner != P1 {
		t.Errorf("expected P1 win, phase=%v winner=%v", s2.Phase, s2.Winner)
	}
}
