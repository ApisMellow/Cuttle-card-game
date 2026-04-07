package card

import "testing"

func TestCardString(t *testing.T) {
	cases := []struct {
		c    Card
		want string
	}{
		{Card{Ace, Spades}, "A♠"},
		{Card{Ten, Diamonds}, "10♦"},
		{Card{King, Clubs}, "K♣"},
		{Card{Two, Hearts}, "2♥"},
	}
	for _, tc := range cases {
		if got := tc.c.String(); got != tc.want {
			t.Errorf("Card{%v,%v}.String() = %q, want %q", tc.c.Rank, tc.c.Suit, got, tc.want)
		}
	}
}

func TestBeats_HigherRank(t *testing.T) {
	if !(Card{Nine, Clubs}).Beats(Card{Four, Spades}) {
		t.Error("9♣ should beat 4♠ (rank dominates suit)")
	}
}

func TestBeats_LowerRank(t *testing.T) {
	if (Card{Four, Spades}).Beats(Card{Nine, Clubs}) {
		t.Error("4♠ should not beat 9♣")
	}
}

func TestBeats_SameRankHigherSuit(t *testing.T) {
	if !(Card{Five, Spades}).Beats(Card{Five, Hearts}) {
		t.Error("5♠ should beat 5♥ (same rank, higher suit)")
	}
}

func TestBeats_SameRankLowerSuit(t *testing.T) {
	if (Card{Five, Clubs}).Beats(Card{Five, Diamonds}) {
		t.Error("5♣ should not beat 5♦")
	}
}

func TestBeats_Identical(t *testing.T) {
	if (Card{Seven, Hearts}).Beats(Card{Seven, Hearts}) {
		t.Error("a card should not beat itself")
	}
}

func TestSuitOrdering(t *testing.T) {
	// Suit order: Clubs < Diamonds < Hearts < Spades
	if !(Clubs < Diamonds && Diamonds < Hearts && Hearts < Spades) {
		t.Error("suit ordering wrong")
	}
}
