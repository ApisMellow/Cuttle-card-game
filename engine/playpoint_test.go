package engine

import (
	"testing"

	"github.com/ApisMellow/cuttle/card"
)

func TestLegalMoves_PlayPoint_AllNumberCards(t *testing.T) {
	hand := []card.Card{
		{Rank: card.Ace, Suit: card.Clubs},
		{Rank: card.Ten, Suit: card.Hearts},
		{Rank: card.Jack, Suit: card.Spades}, // not a point card
	}
	s := twoPlayerStart(hand, nil, []card.Card{{Rank: card.Two, Suit: card.Diamonds}})
	moves := LegalMoves(s)
	pointMoves := 0
	for _, m := range moves {
		if m.Kind == MovePlayPoint {
			pointMoves++
		}
	}
	if pointMoves != 2 {
		t.Errorf("expected 2 PlayPoint moves (A and 10), got %d", pointMoves)
	}
}

func TestApply_PlayPoint_AddsToField(t *testing.T) {
	hand := []card.Card{{Rank: card.Nine, Suit: card.Hearts}}
	s := twoPlayerStart(hand, nil, nil)
	s2, err := Apply(s, Move{Kind: MovePlayPoint, Card: hand[0], HandIndex: 0})
	if err != nil {
		t.Fatal(err)
	}
	if len(s2.Players[P1].Hand) != 0 {
		t.Error("card should leave hand")
	}
	if len(s2.Players[P1].Points) != 1 || s2.Players[P1].Points[0].Card.Rank != card.Nine {
		t.Errorf("point not added: %+v", s2.Players[P1].Points)
	}
	if s2.Players[P1].Points[0].Owner != P1 {
		t.Error("owner should be P1")
	}
	if s2.Active != P2 {
		t.Error("turn should advance")
	}
}

func TestApply_PlayPoint_TriggersWin(t *testing.T) {
	hand := []card.Card{{Rank: card.Nine, Suit: card.Hearts}}
	s := twoPlayerStart(hand, nil, nil)
	s.Players[P1].Points = []PointEntry{
		{Card: card.Card{Rank: card.Seven, Suit: card.Hearts}, Owner: P1},
		{Card: card.Card{Rank: card.Five, Suit: card.Clubs}, Owner: P1},
	}
	s2, err := Apply(s, Move{Kind: MovePlayPoint, Card: hand[0], HandIndex: 0})
	if err != nil {
		t.Fatal(err)
	}
	if s2.Phase != PhaseGameOver {
		t.Errorf("phase should be GameOver, got %v", s2.Phase)
	}
	if s2.Winner == nil || *s2.Winner != P1 {
		t.Errorf("winner should be P1, got %+v", s2.Winner)
	}
}
