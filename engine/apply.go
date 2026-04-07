package engine

import (
	"errors"
)

var ErrIllegalMove = errors.New("illegal move")

// LegalMoves returns every legal move from the current state.
// The REPL displays this list and the user picks one by index.
func LegalMoves(s GameState) []Move {
	if s.Phase == PhaseGameOver {
		return nil
	}
	// Implemented incrementally in later tasks.
	return nil
}

// Apply validates the move (must appear in LegalMoves(s)) and returns the new state.
func Apply(s GameState, m Move) (GameState, error) {
	if s.Phase == PhaseGameOver {
		return s, ErrIllegalMove
	}
	// Implemented incrementally in later tasks.
	return s, ErrIllegalMove
}
