package engine

import (
	"testing"

	"github.com/ApisMellow/cuttle/card"
)

func TestThreshold_NoKings(t *testing.T) {
	if got := Threshold(0); got != 21 {
		t.Errorf("0 kings: got %d, want 21", got)
	}
}

func TestThreshold_OneKing(t *testing.T)    { mustThreshold(t, 1, 14) }
func TestThreshold_TwoKings(t *testing.T)   { mustThreshold(t, 2, 10) }
func TestThreshold_ThreeKings(t *testing.T) { mustThreshold(t, 3, 7) }
func TestThreshold_FourKings(t *testing.T)  { mustThreshold(t, 4, 5) }

func mustThreshold(t *testing.T, kings, want int) {
	t.Helper()
	if got := Threshold(kings); got != want {
		t.Errorf("%d kings: got %d, want %d", kings, got, want)
	}
}

func TestPointTotal_Empty(t *testing.T) {
	ps := PlayerState{}
	if got := PointTotal(ps); got != 0 {
		t.Errorf("empty: got %d, want 0", got)
	}
}

func TestPointTotal_MixedRanks(t *testing.T) {
	ps := PlayerState{
		Points: []PointEntry{
			{Card: card.Card{Rank: card.Ace, Suit: card.Spades}},
			{Card: card.Card{Rank: card.Seven, Suit: card.Hearts}},
			{Card: card.Card{Rank: card.Eight, Suit: card.Clubs}},
			{Card: card.Card{Rank: card.Ten, Suit: card.Diamonds}},
		},
	}
	if got := PointTotal(ps); got != 1+7+8+10 {
		t.Errorf("got %d, want 26", got)
	}
}

func TestKingCount(t *testing.T) {
	ps := PlayerState{
		Permanents: []card.Card{
			{Rank: card.King, Suit: card.Spades},
			{Rank: card.Queen, Suit: card.Hearts},
			{Rank: card.King, Suit: card.Clubs},
			{Rank: card.Eight, Suit: card.Diamonds},
		},
	}
	if got := KingCount(ps); got != 2 {
		t.Errorf("got %d, want 2", got)
	}
}

func TestHasWon(t *testing.T) {
	ps := PlayerState{
		Points: []PointEntry{
			{Card: card.Card{Rank: card.Seven, Suit: card.Hearts}},
			{Card: card.Card{Rank: card.Five, Suit: card.Clubs}},
			{Card: card.Card{Rank: card.Nine, Suit: card.Hearts}},
		},
	}
	if !HasWon(ps) {
		t.Error("21 points, 0 kings: should have won")
	}
}

func TestHasWon_KingReducesThreshold(t *testing.T) {
	ps := PlayerState{
		Points: []PointEntry{
			{Card: card.Card{Rank: card.Seven, Suit: card.Hearts}},
			{Card: card.Card{Rank: card.Seven, Suit: card.Clubs}},
		},
		Permanents: []card.Card{{Rank: card.King, Suit: card.Spades}},
	}
	if !HasWon(ps) {
		t.Error("14 points, 1 king (threshold 14): should have won")
	}
}
