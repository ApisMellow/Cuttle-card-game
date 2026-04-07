package engine

import "github.com/ApisMellow/cuttle/card"

func Threshold(kings int) int {
	switch {
	case kings <= 0:
		return 21
	case kings == 1:
		return 14
	case kings == 2:
		return 10
	case kings == 3:
		return 5
	default:
		return 0
	}
}

func PointTotal(p PlayerState) int {
	total := 0
	for _, pe := range p.Points {
		total += int(pe.Card.Rank)
	}
	return total
}

func KingCount(p PlayerState) int {
	n := 0
	for _, c := range p.Permanents {
		if c.Rank == card.King {
			n++
		}
	}
	return n
}

func HasWon(p PlayerState) bool {
	return PointTotal(p) >= Threshold(KingCount(p))
}
