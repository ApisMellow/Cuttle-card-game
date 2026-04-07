package render

import (
	"fmt"
	"strings"

	"github.com/ApisMellow/cuttle/card"
	"github.com/ApisMellow/cuttle/engine"
)

// oneOffReference is the static side panel showing what each one-off does.
// It's a quick lookup so a player learning the game doesn't have to memorise
// every effect.
var oneOffReference = []string{
	"ONE-OFF EFFECTS",
	"───────────────────────────",
	"A  scrap ALL point cards",
	"2  counter one-off, OR",
	"   scrap a royal/glasses-8",
	"3  take a card from scrap",
	"4  opponent discards 2",
	"5  draw 2 cards",
	"6  scrap ALL royals + 8s",
	"7  reveal top 2, play one",
	"9  bounce a field card",
	"   (frozen next turn)",
	"───────────────────────────",
	"PERMANENTS",
	"8  glasses: see opp hand",
	"J  steal an opp point",
	"Q  protect your other cards",
	"K  lower win threshold",
}

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
	return joinSideBySide(b.String(), strings.Join(oneOffReference, "\n"), 4)
}

// joinSideBySide places right beside left, separated by `gap` spaces.
// The left column is padded to its widest line so the right column is flush.
func joinSideBySide(left, right string, gap int) string {
	leftLines := strings.Split(left, "\n")
	rightLines := strings.Split(right, "\n")
	leftWidth := 0
	for _, l := range leftLines {
		if w := visibleWidth(l); w > leftWidth {
			leftWidth = w
		}
	}
	n := len(leftLines)
	if len(rightLines) > n {
		n = len(rightLines)
	}
	pad := strings.Repeat(" ", gap)
	var b strings.Builder
	for i := 0; i < n; i++ {
		var l, r string
		if i < len(leftLines) {
			l = leftLines[i]
		}
		if i < len(rightLines) {
			r = rightLines[i]
		}
		b.WriteString(l)
		b.WriteString(strings.Repeat(" ", leftWidth-visibleWidth(l)))
		if r != "" {
			b.WriteString(pad)
			b.WriteString(r)
		}
		b.WriteByte('\n')
	}
	return b.String()
}

// visibleWidth counts runes, which is correct for the ASCII + suit-glyph
// content this renderer produces (each suit ♠♥♦♣ is one rune wide in a
// monospace terminal).
func visibleWidth(s string) int {
	n := 0
	for range s {
		n++
	}
	return n
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
