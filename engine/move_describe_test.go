package engine

import (
	"strings"
	"testing"

	"github.com/ApisMellow/cuttle/card"
)

func TestDescribe_PlayPoint(t *testing.T) {
	s := twoPlayerStart([]card.Card{{Rank: card.Nine, Suit: card.Hearts}}, nil, nil)
	m := Move{Kind: MovePlayPoint, Card: card.Card{Rank: card.Nine, Suit: card.Hearts}, HandIndex: 0}
	got := m.Describe(s)
	if !strings.Contains(got, "9♥") || !strings.Contains(strings.ToLower(got), "point") {
		t.Errorf("unexpected description: %q", got)
	}
}

func TestDescribe_Draw(t *testing.T) {
	s := twoPlayerStart(nil, nil, []card.Card{{Rank: card.Two, Suit: card.Clubs}})
	m := Move{Kind: MoveDraw}
	if !strings.Contains(strings.ToLower(m.Describe(s)), "draw") {
		t.Error("draw description should contain 'draw'")
	}
}
