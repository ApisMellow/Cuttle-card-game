package main

import (
	"io"
	"strings"
	"testing"

	"github.com/ApisMellow/cuttle/card"
	"github.com/ApisMellow/cuttle/engine"
)

// endgameState builds a deterministic position where P1 holds only 10♠ with
// 14 points on the table and an empty deck: exactly one legal move (play the
// 10 as a point card), which wins the game 24 > 21.
func endgameState() engine.GameState {
	return engine.GameState{
		Players: [2]engine.PlayerState{
			{
				Hand: []card.Card{{Rank: card.Ten, Suit: card.Spades}},
				Points: []engine.PointEntry{
					{Card: card.Card{Rank: card.Ten, Suit: card.Hearts}, Owner: engine.P1},
					{Card: card.Card{Rank: card.Four, Suit: card.Clubs}, Owner: engine.P1},
				},
			},
			{},
		},
		Active: engine.P1,
		Phase:  engine.PhaseNormal,
	}
}

// TestRun_ScriptedWin drives a full (one-move) game through the REPL loop:
// a garbage line, an out-of-range choice, then the winning move.
func TestRun_ScriptedWin(t *testing.T) {
	var out strings.Builder
	err := run(endgameState(), strings.NewReader("zz\n9\n1\n"), &out)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	got := out.String()
	if n := strings.Count(got, "invalid choice"); n != 2 {
		t.Errorf("got %d 'invalid choice' lines, want 2", n)
	}
	if !strings.Contains(got, "play 10♠ as point card") {
		t.Errorf("menu should offer playing 10♠ as a point card; output:\n%s", got)
	}
	if !strings.Contains(got, "winner: P1") {
		t.Errorf("final render should announce P1 as winner; output:\n%s", got)
	}
}

// TestRun_EOFEndsSession: exhausting input mid-game returns the read error
// instead of looping forever.
func TestRun_EOFEndsSession(t *testing.T) {
	var out strings.Builder
	err := run(endgameState(), strings.NewReader(""), &out)
	if err != io.EOF {
		t.Fatalf("run = %v, want io.EOF", err)
	}
	if !strings.Contains(out.String(), "> ") {
		t.Error("prompt should be printed before reading input")
	}
}

// TestNewGame_Deal: the standard opening deal is 5/6 cards with the
// remaining 41 in the deck.
func TestNewGame_Deal(t *testing.T) {
	s := newGame()
	if len(s.Players[0].Hand) != 5 || len(s.Players[1].Hand) != 6 || len(s.Deck) != 41 {
		t.Errorf("deal = %d/%d hand, %d deck; want 5/6 hand, 41 deck",
			len(s.Players[0].Hand), len(s.Players[1].Hand), len(s.Deck))
	}
	seen := map[card.Card]bool{}
	for _, c := range append(append(append([]card.Card{}, s.Players[0].Hand...), s.Players[1].Hand...), s.Deck...) {
		if seen[c] {
			t.Errorf("card %s dealt twice", c)
		}
		seen[c] = true
	}
	if len(seen) != 52 {
		t.Errorf("%d distinct cards, want 52", len(seen))
	}
}
