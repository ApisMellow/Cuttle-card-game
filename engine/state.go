package engine

// INVARIANT: Jack-as-permanent / point ownership.
//
// A PointEntry always lives in the Points slice of its current *controller*
// — the player who, right now, counts that point toward their score. When a
// Jack is played onto an opponent's point, the entire PointEntry is
// TRANSPLANTED from the victim's Points slice into the attacker's Points
// slice, and the Jack is appended to the entry's JackStack (with the
// attacker in JackOwners). `Owner` on the entry always records the ORIGINAL
// owner (the player who played the point card).
//
// When the top Jack of a PointEntry is scrapped (via 2-as-scrap, 6, 9, etc.)
// and that empties the JackStack, the PointEntry is transplanted back to
// `Owner`'s Points slice. A 9 that bounces a point returns the point card to
// `Owner`'s hand (jack(s) scrapped). A 6 that scraps a stack's jacks returns
// the underlying point card to `Owner`'s Points slice.
//
// Consequence: PointTotal(p) is simply the sum of ranks of p.Points — there
// is no cross-player lookup. `Controller()` is retained as a convenience but
// is always equal to the index of the Players slice that currently contains
// the entry.

import "github.com/ApisMellow/cuttle/card"

type PlayerID uint8

const (
	P1 PlayerID = 0
	P2 PlayerID = 1
)

func (p PlayerID) Other() PlayerID {
	return 1 - p
}

// PointEntry is a point card on a player's side, possibly with Jacks layered on top.
// The current controller is whoever owns the top Jack, or Owner if JackStack is empty.
type PointEntry struct {
	Card       card.Card
	Owner      PlayerID    // original owner; receives the point card if all Jacks are scrapped
	JackStack  []card.Card // bottom-up
	JackOwners []PlayerID  // parallel to JackStack
}

func (pe PointEntry) Controller() PlayerID {
	if len(pe.JackOwners) == 0 {
		return pe.Owner
	}
	return pe.JackOwners[len(pe.JackOwners)-1]
}

type PlayerState struct {
	Hand       []card.Card
	Points     []PointEntry // points currently controlled by this player
	Permanents []card.Card  // Queens, Kings, glasses-8s only (NOT Jacks)
	FrozenIDs  map[int]bool // hand indices frozen by an opponent's 9 last turn
}

type Phase uint8

const (
	PhaseNormal Phase = iota
	PhaseAwaitingCounter
	PhaseSevenChoosing
	PhaseAwaitingDiscard
	PhaseGameOver
)

type TargetZone uint8

const (
	ZonePoints TargetZone = iota
	ZonePermanents
)

type Target struct {
	Owner PlayerID
	Zone  TargetZone
	Index int
}

type PendingOneOff struct {
	PlayedBy     PlayerID
	Card         card.Card
	Target       *Target
	ScrapIndex   int         // for 3: which scrap card to take if this resolves
	Revealed     []card.Card // for 7
	CounterChain []card.Card
}

type GameState struct {
	Players      [2]PlayerState
	Deck         []card.Card // index 0 = top
	Scrap        []card.Card // index len-1 = top
	Active       PlayerID
	Phase        Phase
	Pending      *PendingOneOff
	PassesInARow int
	Winner       *PlayerID
}
