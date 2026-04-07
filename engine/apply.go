package engine

import (
	"errors"

	"github.com/ApisMellow/cuttle/card"
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

func clone(s GameState) GameState {
	out := s
	for i := 0; i < 2; i++ {
		out.Players[i] = clonePlayer(s.Players[i])
	}
	out.Deck = append([]card.Card(nil), s.Deck...)
	out.Scrap = append([]card.Card(nil), s.Scrap...)
	if s.Pending != nil {
		p := *s.Pending
		p.Revealed = append([]card.Card(nil), s.Pending.Revealed...)
		p.CounterChain = append([]card.Card(nil), s.Pending.CounterChain...)
		if s.Pending.Target != nil {
			tgt := *s.Pending.Target
			p.Target = &tgt
		}
		out.Pending = &p
	}
	if s.Winner != nil {
		w := *s.Winner
		out.Winner = &w
	}
	return out
}

func clonePlayer(p PlayerState) PlayerState {
	out := p
	out.Hand = append([]card.Card(nil), p.Hand...)
	out.Permanents = append([]card.Card(nil), p.Permanents...)
	out.Points = make([]PointEntry, len(p.Points))
	for i, pe := range p.Points {
		out.Points[i] = pe
		out.Points[i].JackStack = append([]card.Card(nil), pe.JackStack...)
		out.Points[i].JackOwners = append([]PlayerID(nil), pe.JackOwners...)
	}
	if p.FrozenIDs != nil {
		out.FrozenIDs = make(map[int]bool, len(p.FrozenIDs))
		for k, v := range p.FrozenIDs {
			out.FrozenIDs[k] = v
		}
	}
	return out
}
