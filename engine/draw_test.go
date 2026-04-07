package engine

import (
	"testing"

	"github.com/ApisMellow/cuttle/card"
)

func twoPlayerStart(p1Hand, p2Hand, deck []card.Card) GameState {
	return GameState{
		Players: [2]PlayerState{
			{Hand: p1Hand},
			{Hand: p2Hand},
		},
		Deck:   deck,
		Active: P1,
		Phase:  PhaseNormal,
	}
}

func TestLegalMoves_DrawAvailable(t *testing.T) {
	s := twoPlayerStart(nil, nil, []card.Card{{Rank: card.Ace, Suit: card.Spades}})
	moves := LegalMoves(s)
	if !containsKind(moves, MoveDraw) {
		t.Error("draw should be legal when deck non-empty and hand under limit")
	}
}

func TestLegalMoves_NoDrawWhenDeckEmpty(t *testing.T) {
	s := twoPlayerStart(nil, nil, nil)
	if containsKind(LegalMoves(s), MoveDraw) {
		t.Error("draw should not be legal when deck is empty")
	}
}

func TestLegalMoves_NoDrawWhenHandFull(t *testing.T) {
	full := make([]card.Card, 8)
	for i := range full {
		full[i] = card.Card{Rank: card.Two, Suit: card.Suit(i % 4)}
	}
	s := twoPlayerStart(full, nil, []card.Card{{Rank: card.Ace, Suit: card.Spades}})
	if containsKind(LegalMoves(s), MoveDraw) {
		t.Error("draw should not be legal when hand is full")
	}
}

func TestApply_Draw(t *testing.T) {
	top := card.Card{Rank: card.Five, Suit: card.Hearts}
	s := twoPlayerStart(nil, nil, []card.Card{top, {Rank: card.Six, Suit: card.Clubs}})
	s2, err := Apply(s, Move{Kind: MoveDraw})
	if err != nil {
		t.Fatal(err)
	}
	if len(s2.Players[P1].Hand) != 1 || s2.Players[P1].Hand[0] != top {
		t.Errorf("hand wrong: %+v", s2.Players[P1].Hand)
	}
	if len(s2.Deck) != 1 {
		t.Errorf("deck should have 1 left, got %d", len(s2.Deck))
	}
	if s2.Active != P2 {
		t.Error("turn should advance to P2 after draw")
	}
	if s2.PassesInARow != 0 {
		t.Error("draw resets pass counter")
	}
}

func TestApply_PassWhenNoMovesPossible(t *testing.T) {
	s := twoPlayerStart(nil, nil, nil)
	moves := LegalMoves(s)
	if len(moves) != 1 || moves[0].Kind != MovePass {
		t.Fatalf("expected only Pass, got %+v", moves)
	}
	s2, err := Apply(s, moves[0])
	if err != nil {
		t.Fatal(err)
	}
	if s2.PassesInARow != 1 {
		t.Errorf("pass counter should be 1, got %d", s2.PassesInARow)
	}
	if s2.Active != P2 {
		t.Error("turn should advance after pass")
	}
}

func TestApply_ThreeConsecutivePassesEndsGame(t *testing.T) {
	s := twoPlayerStart(nil, nil, nil)
	for i := 0; i < 3; i++ {
		var err error
		s, err = Apply(s, Move{Kind: MovePass})
		if err != nil {
			t.Fatal(err)
		}
	}
	if s.Phase != PhaseGameOver {
		t.Errorf("phase should be GameOver after 3 passes, got %v", s.Phase)
	}
	if s.Winner != nil {
		t.Error("stalemate has no winner")
	}
}

// containsKind is a small helper used across tests.
func containsKind(moves []Move, k MoveKind) bool {
	for _, m := range moves {
		if m.Kind == k {
			return true
		}
	}
	return false
}
