package main

import (
	"bufio"
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

func main() {
	// A read error (e.g. EOF on ctrl-D) simply ends the session.
	_ = run(newGame(), os.Stdin, os.Stdout)
}

// run drives the hot-seat REPL loop over the given state, reading move
// choices from in and rendering to out, until the game ends or in is
// exhausted. Extracted from main so tests can script a session.
func run(state engine.GameState, in io.Reader, out io.Writer) error {
	reader := bufio.NewReader(in)
	for state.Phase != engine.PhaseGameOver {
		fmt.Fprint(out, "\033[H\033[2J") // clear screen
		fmt.Fprintln(out, render.Render(state))
		moves := engine.LegalMoves(state)
		for i, m := range moves {
			fmt.Fprintf(out, "  %2d) %s\n", i+1, m.Describe(state))
		}
		fmt.Fprint(out, "> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
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
	}
	fmt.Fprint(out, "\033[H\033[2J")
	fmt.Fprintln(out, render.Render(state))
	return nil
}

func newGame() engine.GameState {
	deck := make([]card.Card, 0, 52)
	for s := card.Clubs; s <= card.Spades; s++ {
		for r := card.Ace; r <= card.King; r++ {
			deck = append(deck, card.Card{Rank: r, Suit: s})
		}
	}
	rand.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })
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
