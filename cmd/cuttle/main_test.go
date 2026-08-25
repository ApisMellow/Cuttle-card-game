package main

import (
	"io"
	"path/filepath"
	"reflect"
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
// a garbage line, an out-of-range choice, then the winning move. The move
// log must record exactly the one applied move.
func TestRun_ScriptedWin(t *testing.T) {
	var out strings.Builder
	moves, err := run(endgameState(), 7, strings.NewReader("zz\n9\n1\n"), &out)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if !reflect.DeepEqual(moves, []int{0}) {
		t.Errorf("move log = %v, want [0]", moves)
	}
	got := out.String()
	if n := strings.Count(got, "invalid choice"); n != 2 {
		t.Errorf("got %d 'invalid choice' lines, want 2", n)
	}
	if !strings.Contains(got, "seed: 7") {
		t.Error("each frame should display the seed")
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
	_, err := run(endgameState(), 0, strings.NewReader(""), &out)
	if err != io.EOF {
		t.Fatalf("run = %v, want io.EOF", err)
	}
	if !strings.Contains(out.String(), "> ") {
		t.Error("prompt should be printed before reading input")
	}
}

// TestNewGame_Deal: the standard opening deal is 5/6 cards with the
// remaining 41 in the deck, and no card appears twice.
func TestNewGame_Deal(t *testing.T) {
	s := newGame(1)
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

// TestNewGame_SeedDeterminism: the same seed produces the identical initial
// state; a different seed produces a different shuffle.
func TestNewGame_SeedDeterminism(t *testing.T) {
	if !reflect.DeepEqual(newGame(42), newGame(42)) {
		t.Error("same seed should produce identical initial states")
	}
	if reflect.DeepEqual(newGame(42), newGame(43)) {
		t.Error("different seeds should produce different shuffles")
	}
}

// TestReplay_RoundTrip: play a fixed-seed game for up to 30 moves always
// picking the first legal move, save the replay to JSON, load it back,
// replay it, and require the final state to match bit-for-bit.
func TestReplay_RoundTrip(t *testing.T) {
	const seed = 42
	state := newGame(seed)
	var chosen []int
	for i := 0; i < 30 && state.Phase != engine.PhaseGameOver; i++ {
		moves := engine.LegalMoves(state)
		next, err := engine.Apply(state, moves[0])
		if err != nil {
			t.Fatalf("move %d: %v", i+1, err)
		}
		state = next
		chosen = append(chosen, 0)
	}

	path := filepath.Join(t.TempDir(), "game.json")
	if err := saveReplay(path, Replay{Seed: seed, Moves: chosen}); err != nil {
		t.Fatalf("save: %v", err)
	}
	rep, err := loadReplay(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	replayed, err := replayGame(rep, io.Discard)
	if err != nil {
		t.Fatalf("replay: %v", err)
	}
	if !reflect.DeepEqual(state, replayed) {
		t.Errorf("replayed final state differs from original\noriginal: %+v\nreplayed: %+v", state, replayed)
	}
}

// TestReplayGame_BadIndex: a corrupt replay (index out of range) fails with
// an error naming the step instead of panicking.
func TestReplayGame_BadIndex(t *testing.T) {
	_, err := replayGame(Replay{Seed: 1, Moves: []int{9999}}, io.Discard)
	if err == nil || !strings.Contains(err.Error(), "step 1") {
		t.Errorf("err = %v, want out-of-range error naming step 1", err)
	}
}
