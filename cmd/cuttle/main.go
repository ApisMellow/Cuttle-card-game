package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/ApisMellow/cuttle/card"
	"github.com/ApisMellow/cuttle/engine"
	"github.com/ApisMellow/cuttle/render"
)

// Replay is the on-disk record of a finished game: the shuffle seed plus the
// sequence of chosen move indices (0-based positions in LegalMoves at each
// step). Together they reproduce the game bit-for-bit.
type Replay struct {
	Seed  int64 `json:"seed"`
	Moves []int `json:"moves"`
}

func main() {
	seed := flag.Int64("seed", 0, "deck shuffle seed (0 = pick one at random and print it)")
	replayPath := flag.String("replay", "", "replay a saved game JSON instead of playing")
	record := flag.String("record", "", "path for the replay JSON written at game end (default cuttle-replay-<seed>.json)")
	flag.Parse()

	if *replayPath != "" {
		if err := replayFromFile(*replayPath, os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "replay: %v\n", err)
			os.Exit(1)
		}
		return
	}

	s := *seed
	if s == 0 {
		s = rand.Int63()
	}
	moves, err := run(newGame(s), s, os.Stdin, os.Stdout)
	if err != nil {
		// A read error (e.g. EOF on ctrl-D) ends the session; nothing to save.
		return
	}
	path := *record
	if path == "" {
		path = fmt.Sprintf("cuttle-replay-%d.json", s)
	}
	if err := saveReplay(path, Replay{Seed: s, Moves: moves}); err != nil {
		fmt.Fprintf(os.Stderr, "saving replay: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("replay saved to %s (replay with: cuttle --replay %s)\n", path, path)
}

// run drives the hot-seat REPL loop over the given state, reading move
// choices from in and rendering to out, until the game ends or in is
// exhausted. It returns the 0-based index of every move applied, in order.
// The seed is displayed on every frame so any game can be reproduced.
func run(state engine.GameState, seed int64, in io.Reader, out io.Writer) ([]int, error) {
	reader := bufio.NewReader(in)
	var chosen []int
	for state.Phase != engine.PhaseGameOver {
		fmt.Fprint(out, "\033[H\033[2J") // clear screen
		fmt.Fprintf(out, "seed: %d\n", seed)
		fmt.Fprintln(out, render.Render(state))
		moves := engine.LegalMoves(state)
		for i, m := range moves {
			fmt.Fprintf(out, "  %2d) %s\n", i+1, m.Describe(state))
		}
		fmt.Fprint(out, "> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return chosen, err
		}
		idx, err := strconv.Atoi(strings.TrimSpace(line))
		if err != nil || idx < 1 || idx > len(moves) {
			fmt.Fprintln(out, "invalid choice")
			continue
		}
		next, err := engine.Apply(state, moves[idx-1])
		if err != nil {
			fmt.Fprintf(out, "engine error: %v\n", err)
			continue
		}
		state = next
		chosen = append(chosen, idx-1)
	}
	fmt.Fprint(out, "\033[H\033[2J")
	fmt.Fprintf(out, "seed: %d\n", seed)
	fmt.Fprintln(out, render.Render(state))
	return chosen, nil
}

// replayGame rebuilds a game from a replay: reshuffle from the seed, then
// apply each recorded move index in order. Every step is rendered to out
// (pass io.Discard to skip rendering). Returns the final state.
func replayGame(rep Replay, out io.Writer) (engine.GameState, error) {
	state := newGame(rep.Seed)
	for step, idx := range rep.Moves {
		moves := engine.LegalMoves(state)
		if idx < 0 || idx >= len(moves) {
			return state, fmt.Errorf("step %d: move index %d out of range (%d legal moves)", step+1, idx, len(moves))
		}
		fmt.Fprintf(out, "--- step %d: %s\n", step+1, moves[idx].Describe(state))
		next, err := engine.Apply(state, moves[idx])
		if err != nil {
			return state, fmt.Errorf("step %d: %w", step+1, err)
		}
		state = next
		fmt.Fprintln(out, render.Render(state))
	}
	return state, nil
}

func replayFromFile(path string, out io.Writer) error {
	rep, err := loadReplay(path)
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "replaying %s (seed %d, %d moves)\n", path, rep.Seed, len(rep.Moves))
	_, err = replayGame(rep, out)
	return err
}

func saveReplay(path string, rep Replay) error {
	data, err := json.MarshalIndent(rep, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func loadReplay(path string) (Replay, error) {
	var rep Replay
	data, err := os.ReadFile(path)
	if err != nil {
		return rep, err
	}
	if err := json.Unmarshal(data, &rep); err != nil {
		return rep, err
	}
	return rep, nil
}

// newGame builds the opening deal from a seeded shuffle: the same seed
// always produces the same game.
func newGame(seed int64) engine.GameState {
	rng := rand.New(rand.NewSource(seed))
	deck := make([]card.Card, 0, 52)
	for s := card.Clubs; s <= card.Spades; s++ {
		for r := card.Ace; r <= card.King; r++ {
			deck = append(deck, card.Card{Rank: r, Suit: s})
		}
	}
	rng.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })
	p1 := append([]card.Card(nil), deck[:5]...)
	p2 := append([]card.Card(nil), deck[5:11]...)
	rest := append([]card.Card(nil), deck[11:]...)
	return engine.GameState{
		Players: [2]engine.PlayerState{
			{Hand: p1},
			{Hand: p2},
		},
		Deck:   rest,
		Active: engine.P1,
		Phase:  engine.PhaseNormal,
	}
}
