package engine

import (
	"errors"

	"github.com/ApisMellow/cuttle/card"
)

var ErrIllegalMove = errors.New("illegal move")

const HandLimit = 8

// LegalMoves returns every legal move from the current state.
// The REPL displays this list and the user picks one by index.
func LegalMoves(s GameState) []Move {
	if s.Phase == PhaseGameOver {
		return nil
	}
	if s.Phase == PhaseAwaitingCounter {
		return legalCounterMoves(s)
	}
	if s.Phase != PhaseNormal {
		// later phases handled in later tasks
		return nil
	}
	var moves []Move
	active := s.Players[s.Active]
	if len(s.Deck) > 0 && len(active.Hand) < HandLimit {
		moves = append(moves, Move{Kind: MoveDraw})
	}
	for i, c := range active.Hand {
		if active.FrozenIDs[i] {
			continue
		}
		if c.Rank == card.Ace {
			moves = append(moves, Move{Kind: MoveOneOff, Card: c, HandIndex: i})
		}
		if c.Rank >= card.Ace && c.Rank <= card.Ten {
			moves = append(moves, Move{Kind: MovePlayPoint, Card: c, HandIndex: i})
			opp := s.Active.Other()
			for j, pe := range s.Players[opp].Points {
				if c.Beats(pe.Card) {
					moves = append(moves, Move{
						Kind: MoveScuttle, Card: c, HandIndex: i,
						Target: &Target{Owner: opp, Zone: ZonePoints, Index: j},
					})
				}
			}
		}
		if c.Rank == card.Queen || c.Rank == card.King || c.Rank == card.Eight {
			moves = append(moves, Move{Kind: MovePlayPermanent, Card: c, HandIndex: i})
		}
	}
	if len(moves) == 0 {
		moves = append(moves, Move{Kind: MovePass})
	}
	return moves
}

// Apply validates the move (must appear in LegalMoves(s)) and returns the new state.
func Apply(s GameState, m Move) (GameState, error) {
	if s.Phase == PhaseGameOver {
		return s, ErrIllegalMove
	}
	out := clone(s)
	switch m.Kind {
	case MoveDraw:
		if len(out.Deck) == 0 || len(out.Players[out.Active].Hand) >= HandLimit {
			return s, ErrIllegalMove
		}
		top := out.Deck[0]
		out.Deck = out.Deck[1:]
		out.Players[out.Active].Hand = append(out.Players[out.Active].Hand, top)
		out.PassesInARow = 0
		endTurn(&out)
		return out, nil
	case MovePass:
		out.PassesInARow++
		if out.PassesInARow >= 3 {
			out.Phase = PhaseGameOver
			return out, nil
		}
		endTurn(&out)
		return out, nil
	case MovePlayPoint:
		p := &out.Players[out.Active]
		if m.HandIndex < 0 || m.HandIndex >= len(p.Hand) || p.Hand[m.HandIndex] != m.Card {
			return s, ErrIllegalMove
		}
		if m.Card.Rank < card.Ace || m.Card.Rank > card.Ten {
			return s, ErrIllegalMove
		}
		p.Hand = removeAt(p.Hand, m.HandIndex)
		p.Points = append(p.Points, PointEntry{Card: m.Card, Owner: out.Active})
		out.PassesInARow = 0
		if checkWin(&out, out.Active) {
			return out, nil
		}
		endTurn(&out)
		return out, nil
	case MovePlayPermanent:
		p := &out.Players[out.Active]
		if m.HandIndex < 0 || m.HandIndex >= len(p.Hand) || p.Hand[m.HandIndex] != m.Card {
			return s, ErrIllegalMove
		}
		if m.Card.Rank != card.Queen && m.Card.Rank != card.King && m.Card.Rank != card.Eight {
			// Jacks handled in a later task
			return s, ErrIllegalMove
		}
		p.Hand = removeAt(p.Hand, m.HandIndex)
		p.Permanents = append(p.Permanents, m.Card)
		out.PassesInARow = 0
		if checkWin(&out, out.Active) {
			return out, nil
		}
		endTurn(&out)
		return out, nil
	case MoveScuttle:
		p := &out.Players[out.Active]
		if m.HandIndex < 0 || m.HandIndex >= len(p.Hand) || p.Hand[m.HandIndex] != m.Card {
			return s, ErrIllegalMove
		}
		if m.Target == nil || m.Target.Zone != ZonePoints {
			return s, ErrIllegalMove
		}
		opp := &out.Players[m.Target.Owner]
		if m.Target.Index < 0 || m.Target.Index >= len(opp.Points) {
			return s, ErrIllegalMove
		}
		target := opp.Points[m.Target.Index]
		if !m.Card.Beats(target.Card) {
			return s, ErrIllegalMove
		}
		p.Hand = removeAt(p.Hand, m.HandIndex)
		out.Scrap = append(out.Scrap, m.Card, target.Card)
		out.Scrap = append(out.Scrap, target.JackStack...)
		opp.Points = removeAt(opp.Points, m.Target.Index)
		out.PassesInARow = 0
		endTurn(&out)
		return out, nil
	case MoveOneOff:
		p := &out.Players[out.Active]
		if m.HandIndex < 0 || m.HandIndex >= len(p.Hand) || p.Hand[m.HandIndex] != m.Card {
			return s, ErrIllegalMove
		}
		// Task 13/14: only Ace is implemented as a one-off effect.
		if m.Card.Rank != card.Ace {
			return s, ErrIllegalMove
		}
		played := out.Active
		p.Hand = removeAt(p.Hand, m.HandIndex)
		out.PassesInARow = 0
		// If the opponent has a non-frozen 2, enter PhaseAwaitingCounter.
		opp := played.Other()
		if hasLegalCounter(out.Players[opp]) {
			out.Pending = &PendingOneOff{
				PlayedBy: played,
				Card:     m.Card,
			}
			out.Phase = PhaseAwaitingCounter
			out.Active = opp
			return out, nil
		}
		resolveOneOff(&out, m.Card, played)
		return out, nil
	case MoveCounter:
		if s.Phase != PhaseAwaitingCounter || out.Pending == nil {
			return s, ErrIllegalMove
		}
		p := &out.Players[out.Active]
		if m.HandIndex < 0 || m.HandIndex >= len(p.Hand) || p.Hand[m.HandIndex] != m.Card {
			return s, ErrIllegalMove
		}
		if m.Card.Rank != card.Two {
			return s, ErrIllegalMove
		}
		if p.FrozenIDs[m.HandIndex] {
			return s, ErrIllegalMove
		}
		p.Hand = removeAt(p.Hand, m.HandIndex)
		out.Pending.CounterChain = append(out.Pending.CounterChain, m.Card)
		// Flip waiting player; if they can counter, stay in phase; else auto-resolve.
		next := out.Active.Other()
		if hasLegalCounter(out.Players[next]) {
			out.Active = next
			return out, nil
		}
		resolvePending(&out)
		return out, nil
	case MoveDecline:
		if s.Phase != PhaseAwaitingCounter || out.Pending == nil {
			return s, ErrIllegalMove
		}
		resolvePending(&out)
		return out, nil
	}
	return s, ErrIllegalMove
}

// hasLegalCounter reports whether p has any non-frozen 2 in hand.
func hasLegalCounter(p PlayerState) bool {
	for i, c := range p.Hand {
		if c.Rank == card.Two && !p.FrozenIDs[i] {
			return true
		}
	}
	return false
}

// legalCounterMoves lists decline + one MoveCounter per legal 2.
func legalCounterMoves(s GameState) []Move {
	moves := []Move{{Kind: MoveDecline}}
	p := s.Players[s.Active]
	for i, c := range p.Hand {
		if c.Rank == card.Two && !p.FrozenIDs[i] {
			moves = append(moves, Move{Kind: MoveCounter, Card: c, HandIndex: i})
		}
	}
	return moves
}

// resolvePending finalizes a counter chain: even len → resolve, odd → cancel.
// Either way, the original card and every 2 in the chain go to scrap, then
// control returns to the original player and endTurn runs.
func resolvePending(s *GameState) {
	pend := s.Pending
	s.Pending = nil
	s.Phase = PhaseNormal
	s.Active = pend.PlayedBy
	cancelled := len(pend.CounterChain)%2 == 1
	if cancelled {
		// Card is cancelled; just scrap original + chain.
		s.Scrap = append(s.Scrap, pend.Card)
		s.Scrap = append(s.Scrap, pend.CounterChain...)
		endTurn(s)
		return
	}
	// Resolve the original one-off's effect. The chain 2s go to scrap alongside.
	resolveOneOff(s, pend.Card, pend.PlayedBy)
	// Append the chain 2s to scrap (resolveOneOff already scrapped the original + effect).
	s.Scrap = append(s.Scrap, pend.CounterChain...)
}

// resolveOneOff applies the effect of a one-off card played by `played` and
// scraps the card itself, runs win check, and ends the turn.
func resolveOneOff(s *GameState, c card.Card, played PlayerID) {
	s.Active = played
	s.Scrap = append(s.Scrap, c)
	switch c.Rank {
	case card.Ace:
		for i := 0; i < 2; i++ {
			pl := &s.Players[i]
			for _, pe := range pl.Points {
				s.Scrap = append(s.Scrap, pe.Card)
				s.Scrap = append(s.Scrap, pe.JackStack...)
			}
			pl.Points = nil
		}
	}
	if checkWin(s, played) {
		return
	}
	endTurn(s)
}

func removeAt[T any](xs []T, i int) []T {
	out := make([]T, 0, len(xs)-1)
	out = append(out, xs[:i]...)
	out = append(out, xs[i+1:]...)
	return out
}

// checkWin sets phase/winner if p has reached threshold. Returns true if game ended.
func checkWin(s *GameState, p PlayerID) bool {
	if HasWon(s.Players[p]) {
		s.Phase = PhaseGameOver
		winner := p
		s.Winner = &winner
		return true
	}
	return false
}

// endTurn advances Active and clears the new active player's frozen marks.
func endTurn(s *GameState) {
	s.Active = s.Active.Other()
	s.Players[s.Active].FrozenIDs = nil
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
