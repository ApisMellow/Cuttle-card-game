package card

import "fmt"

type Suit uint8

const (
	Clubs Suit = iota
	Diamonds
	Hearts
	Spades
)

func (s Suit) String() string {
	return [...]string{"♣", "♦", "♥", "♠"}[s]
}

type Rank uint8

const (
	Ace Rank = 1 + iota
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

func (r Rank) String() string {
	return [...]string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}[r-1]
}

type Card struct {
	Rank Rank
	Suit Suit
}

func (c Card) String() string {
	return fmt.Sprintf("%s%s", c.Rank, c.Suit)
}

// Beats reports whether c can scuttle other:
// strictly higher rank, or equal rank with strictly higher suit.
func (c Card) Beats(other Card) bool {
	if c.Rank != other.Rank {
		return c.Rank > other.Rank
	}
	return c.Suit > other.Suit
}
