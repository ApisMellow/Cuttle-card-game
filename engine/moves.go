package engine

import "github.com/ApisMellow/cuttle/card"

type MoveKind uint8

const (
	MoveDraw MoveKind = iota
	MovePlayPoint
	MovePlayPermanent
	MoveScuttle
	MoveOneOff
	MoveCounter
	MoveDecline
	MoveSevenPick
	MoveDiscardPair
	MovePass
)

type Move struct {
	Kind       MoveKind
	Card       card.Card // card from hand being played (or chosen, for 7)
	HandIndex  int       // index in active player's hand of Card
	Target     *Target
	JackTarget *Target // for MovePlayPermanent when Card is a Jack
	ScrapIndex int     // for MoveOneOff with rank 3 (which scrap card to take)
	DiscardA   int     // for MoveDiscardPair: hand indices to discard
	DiscardB   int
	SubMove    *Move // for MoveSevenPick
}
