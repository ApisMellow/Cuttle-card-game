package render

import (
	"fmt"
	"strings"

	"github.com/ApisMellow/cuttle/card"
	"github.com/ApisMellow/cuttle/engine"
)

func Render(s engine.GameState) string {
	var b strings.Builder
	fmt.Fprintf(&b, "=== CUTTLE ===  Deck: %d  Scrap: %d\n\n", len(s.Deck), len(s.Scrap))
	renderPlayer(&b, s, engine.P2, "PLAYER 2")
	b.WriteString("─────────────────────────────────────\n")
	renderPlayer(&b, s, engine.P1, "PLAYER 1")
	if s.Phase == engine.PhaseGameOver {
		if s.Winner != nil {
			fmt.Fprintf(&b, "\n*** GAME OVER — winner: P%d ***\n", *s.Winner+1)
		} else {
			b.WriteString("\n*** GAME OVER — stalemate ***\n")
		}
	} else {
		fmt.Fprintf(&b, "\n*** P%d TO MOVE ***\n", s.Active+1)
		if s.Phase != engine.PhaseNormal {
			fmt.Fprintf(&b, "Phase: %v\n", s.Phase)
		}
	}
	return b.String()
}

func renderPlayer(b *strings.Builder, s engine.GameState, p engine.PlayerID, label string) {
	ps := s.Players[p]
	fmt.Fprintf(b, "%s (%d cards): %s\n", label, len(ps.Hand), cardsString(ps.Hand))
	fmt.Fprintf(b, "  Points (%d): %s\n", engine.PointTotal(ps), pointsString(ps.Points))
	fmt.Fprintf(b, "  Permanents:  %s\n", cardsString(ps.Permanents))
	fmt.Fprintf(b, "  Threshold: %d\n", engine.Threshold(engine.KingCount(ps)))
}

func cardsString(cs []card.Card) string {
	if len(cs) == 0 {
		return "(none)"
	}
	parts := make([]string, len(cs))
	for i, c := range cs {
		parts[i] = c.String()
	}
	return strings.Join(parts, " ")
}

func pointsString(pes []engine.PointEntry) string {
	if len(pes) == 0 {
		return "(none)"
	}
	parts := make([]string, len(pes))
	for i, pe := range pes {
		s := pe.Card.String()
		for _, j := range pe.JackStack {
			s += "+" + j.String()
		}
		parts[i] = s
	}
	return strings.Join(parts, "  ")
}
