package engine

import (
	"testing"

	"github.com/ApisMellow/cuttle/card"
)

func TestLegalMoves_ScuttleEnumeratesValidTargets(t *testing.T) {
	hand := []card.Card{{Rank: card.Nine, Suit: card.Hearts}}
	s := twoPlayerStart(hand, nil, nil)
	s.Players[P2].Points = []PointEntry{
		{Card: card.Card{Rank: card.Four, Suit: card.Spades}, Owner: P2},
		{Card: card.Card{Rank: card.Ten, Suit: card.Diamonds}, Owner: P2},
		{Card: card.Card{Rank: card.Nine, Suit: card.Clubs}, Owner: P2}, // 9♥ beats 9♣ (suit)
	}
	scuttles := 0
	for _, m := range LegalMoves(s) {
		if m.Kind == MoveScuttle {
			scuttles++
		}
	}
	if scuttles != 2 { // 4♠ and 9♣
		t.Errorf("expected 2 scuttle targets, got %d", scuttles)
	}
}

func TestApply_Scuttle_BothCardsToScrap(t *testing.T) {
	hand := []card.Card{{Rank: card.Nine, Suit: card.Hearts}}
	s := twoPlayerStart(hand, nil, nil)
	s.Players[P2].Points = []PointEntry{
		{Card: card.Card{Rank: card.Four, Suit: card.Spades}, Owner: P2},
	}
	target := &Target{Owner: P2, Zone: ZonePoints, Index: 0}
	s2, err := Apply(s, Move{Kind: MoveScuttle, Card: hand[0], HandIndex: 0, Target: target})
	if err != nil {
		t.Fatal(err)
	}
	if len(s2.Players[P2].Points) != 0 {
		t.Error("opponent point should be removed")
	}
	if len(s2.Scrap) != 2 {
		t.Errorf("scrap should have 2 cards, got %d", len(s2.Scrap))
	}
	if len(s2.Players[P1].Hand) != 0 {
		t.Error("scuttle card should leave hand")
	}
}

func TestApply_Scuttle_ScrapsJackStackWithPoint(t *testing.T) {
	hand := []card.Card{{Rank: card.Ten, Suit: card.Spades}}
	s := twoPlayerStart(hand, nil, nil)
	s.Players[P2].Points = []PointEntry{
		{
			Card:       card.Card{Rank: card.Five, Suit: card.Hearts},
			Owner:      P1,
			JackStack:  []card.Card{{Rank: card.Jack, Suit: card.Diamonds}},
			JackOwners: []PlayerID{P2},
		},
	}
	target := &Target{Owner: P2, Zone: ZonePoints, Index: 0}
	s2, err := Apply(s, Move{Kind: MoveScuttle, Card: hand[0], HandIndex: 0, Target: target})
	if err != nil {
		t.Fatal(err)
	}
	if len(s2.Players[P2].Points) != 0 {
		t.Error("point should be gone")
	}
	if len(s2.Scrap) != 3 { // 10♠, 5♥, J♦
		t.Errorf("scrap should have 3 cards, got %d", len(s2.Scrap))
	}
}
