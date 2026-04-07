package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/ApisMellow/cuttle/card"
	"github.com/ApisMellow/cuttle/engine"
	"github.com/ApisMellow/cuttle/render"
)

func main() {
	state := newGame()
	reader := bufio.NewReader(os.Stdin)
	for state.Phase != engine.PhaseGameOver {
		fmt.Print("\033[H\033[2J") // clear screen
		fmt.Println(render.Render(state))
		moves := engine.LegalMoves(state)
		for i, m := range moves {
			fmt.Printf("  %2d) %s\n", i+1, m.Describe(state))
		}
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		idx, err := strconv.Atoi(strings.TrimSpace(line))
		if err != nil || idx < 1 || idx > len(moves) {
			fmt.Println("invalid choice")
			continue
		}
		next, err := engine.Apply(state, moves[idx-1])
		if err != nil {
			fmt.Printf("engine error: %v\n", err)
			continue
		}
		state = next
	}
	fmt.Print("\033[H\033[2J")
	fmt.Println(render.Render(state))
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
