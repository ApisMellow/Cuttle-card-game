package engine

import (
	"testing"

	"github.com/ApisMellow/cuttle/card"
)

func TestClone_DeepCopiesHands(t *testing.T) {
	s := GameState{
		Players: [2]PlayerState{
			{Hand: []card.Card{{Rank: card.Ace, Suit: card.Spades}}},
			{Hand: []card.Card{{Rank: card.Two, Suit: card.Hearts}}},
		},
		Deck:  []card.Card{{Rank: card.Three, Suit: card.Clubs}},
		Scrap: []card.Card{{Rank: card.Four, Suit: card.Diamonds}},
	}
	c := clone(s)
	c.Players[0].Hand[0] = card.Card{Rank: card.King, Suit: card.Spades}
	c.Deck[0] = card.Card{Rank: card.King, Suit: card.Spades}
	c.Scrap[0] = card.Card{Rank: card.King, Suit: card.Spades}
	if s.Players[0].Hand[0].Rank != card.Ace {
		t.Error("clone shared hand slice")
	}
	if s.Deck[0].Rank != card.Three {
		t.Error("clone shared deck slice")
	}
	if s.Scrap[0].Rank != card.Four {
		t.Error("clone shared scrap slice")
	}
}

func TestClone_DeepCopiesPointsAndJackStacks(t *testing.T) {
	s := GameState{
		Players: [2]PlayerState{
			{Points: []PointEntry{{
				Card:       card.Card{Rank: card.Seven, Suit: card.Hearts},
				JackStack:  []card.Card{{Rank: card.Jack, Suit: card.Spades}},
				JackOwners: []PlayerID{P2},
			}}},
			{},
		},
	}
	c := clone(s)
	c.Players[0].Points[0].JackStack[0] = card.Card{Rank: card.King, Suit: card.Clubs}
	if s.Players[0].Points[0].JackStack[0].Rank != card.Jack {
		t.Error("clone shared jack stack")
	}
}

func TestPointEntry_Controller(t *testing.T) {
	plain := PointEntry{Card: card.Card{Rank: card.Ten, Suit: card.Clubs}, Owner: P2}
	if got := plain.Controller(); got != P2 {
		t.Errorf("no jacks: Controller() = P%d, want P2", got+1)
	}
	stacked := PointEntry{
		Card:       card.Card{Rank: card.Ten, Suit: card.Clubs},
		Owner:      P1,
		JackStack:  []card.Card{{Rank: card.Jack, Suit: card.Clubs}, {Rank: card.Jack, Suit: card.Hearts}},
		JackOwners: []PlayerID{P2, P1},
	}
	if got := stacked.Controller(); got != P1 {
		t.Errorf("stacked jacks: Controller() = P%d, want top jack's owner P1", got+1)
	}
}
