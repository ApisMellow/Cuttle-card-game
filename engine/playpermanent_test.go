package engine

import (
	"testing"

	"github.com/ApisMellow/cuttle/card"
)

func TestLegalMoves_PlayPermanent_QueenKingGlasses(t *testing.T) {
	hand := []card.Card{
		{Rank: card.Queen, Suit: card.Hearts},
		{Rank: card.King, Suit: card.Spades},
		{Rank: card.Eight, Suit: card.Clubs},
		{Rank: card.Five, Suit: card.Diamonds}, // not a permanent
	}
	s := twoPlayerStart(hand, nil, nil)
	moves := LegalMoves(s)
	count := 0
	for _, m := range moves {
		if m.Kind == MovePlayPermanent {
			count++
		}
	}
	// 8 should produce a permanent move (glasses); 5 should not. Q and K each one.
	// Plus 8 also produces a PlayPoint move (covered separately).
	if count != 3 {
		t.Errorf("expected 3 PlayPermanent moves (Q,K,8-glasses), got %d", count)
	}
}

func TestApply_PlayPermanent_Queen(t *testing.T) {
	q := card.Card{Rank: card.Queen, Suit: card.Hearts}
	s := twoPlayerStart([]card.Card{q}, nil, nil)
	s2, err := Apply(s, Move{Kind: MovePlayPermanent, Card: q, HandIndex: 0})
	if err != nil {
		t.Fatal(err)
	}
	if len(s2.Players[P1].Permanents) != 1 || s2.Players[P1].Permanents[0] != q {
		t.Errorf("queen not added: %+v", s2.Players[P1].Permanents)
	}
	if s2.Active != P2 {
		t.Error("turn should advance")
	}
}

func TestApply_PlayPermanent_KingTriggersWinViaThreshold(t *testing.T) {
	k := card.Card{Rank: card.King, Suit: card.Spades}
	s := twoPlayerStart([]card.Card{k}, nil, nil)
	s.Players[P1].Points = []PointEntry{
		{Card: card.Card{Rank: card.Seven, Suit: card.Hearts}, Owner: P1},
		{Card: card.Card{Rank: card.Seven, Suit: card.Clubs}, Owner: P1},
	}
	s2, err := Apply(s, Move{Kind: MovePlayPermanent, Card: k, HandIndex: 0})
	if err != nil {
		t.Fatal(err)
	}
	if s2.Phase != PhaseGameOver || s2.Winner == nil || *s2.Winner != P1 {
		t.Errorf("king should drop threshold to 14 and win; phase=%v winner=%v", s2.Phase, s2.Winner)
	}
}

func TestApply_PlayPermanent_GlassesEight(t *testing.T) {
	e := card.Card{Rank: card.Eight, Suit: card.Diamonds}
	s := twoPlayerStart([]card.Card{e}, nil, nil)
	s2, err := Apply(s, Move{Kind: MovePlayPermanent, Card: e, HandIndex: 0})
	if err != nil {
		t.Fatal(err)
	}
	if len(s2.Players[P1].Permanents) != 1 || s2.Players[P1].Permanents[0] != e {
		t.Error("8 should be added as permanent")
	}
}
