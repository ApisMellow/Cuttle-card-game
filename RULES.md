# Cuttle — Rules (as implemented)

Two-player card game played with a standard 52-card deck.

## Setup
- Standard 52-card deck, shuffled.
- One player is dealt 6 cards, the other is dealt 5 and goes first.
- A shared **scrap pile** sits between the players.
- Each player has their own field with two zones: **point cards** and **permanents**.

## Win Condition
- A player wins when they have **21 or more points** on their side of the field at the end of their own action.
- Wins are only checked on your own turn. The 2 (counter) is the only card playable on the opponent's turn, and a counter cannot add points, so an off-turn win is impossible.
- **Kings** lower your own threshold: 1 King → 14, 2 Kings → 10, 3 Kings → 7, 4 Kings → 5.

## Hand Limit
- Maximum hand size: **8 cards**. You cannot draw past it.

## Turn Structure
On your turn, take exactly one action. If the deck is empty and you cannot otherwise act, you pass. Three consecutive passes = stalemate.

### Legal Actions
1. **Draw** one card from the deck.
2. Play a number card (A–10) as a **point card** on your side.
3. **Scuttle**: play a number card from hand onto an opponent's point card of strictly lower rank, OR equal rank with higher suit. Both cards go to scrap. Suit order (low → high): ♣ < ♦ < ♥ < ♠.
4. Play a **permanent** (Jack, Queen, King, or any 8 played as "glasses").
5. Play a **one-off** (A, 2, 3, 4, 5, 6, 7, 9). Resolve its effect, then scrap it.

## Card Effects

### Point Cards
- **A–10** may be played as point cards. Ace = 1, face value otherwise. 10 is point-only (no one-off).
- **8** may be played as a point card (8 points) OR as a permanent (glasses).

### One-Offs
| Card | Effect |
|---|---|
| **A** | Scrap all point cards on both sides of the field. |
| **2** | Counter an opponent's one-off as it is played (resolves before that one-off's effect; **2s can counter 2s**); OR scrap a target royal or glasses-8. The 2 is the only card that may be played on the opponent's turn. |
| **3** | Take any one card from the scrap pile into your hand. |
| **4** | Opponent discards 2 cards of their choice from hand to scrap. If their hand has fewer than 2 cards, they discard whatever they have. |
| **5** | Draw 2 cards and add them to your hand (respecting the 8-card limit). |
| **6** | Scrap all royals and glasses-8s on both sides. |
| **7** | Reveal the top 2 cards of the deck. Play one immediately (as point, scuttle, permanent, or one-off — normal restrictions apply). Return the other to the top of the deck. If neither revealed card has a legal play, choose one to scrap; the other returns to the top of the deck. |
| **9** | Return an opponent's field card (point or permanent) to their hand. That card cannot be played on their next turn. |

### Permanents
| Card | Effect |
|---|---|
| **8 (glasses)** | While in play, your opponent's hand is visible to you. Any 8 may be played this way. |
| **J** | Played on top of an opponent's point card to **steal** it onto your side (it now counts for you). If the Jack is later scrapped or stolen back, the underlying point card returns to its original owner. A Jack may target a point card already under another Jack (chain-steal). |
| **Q** | Your **other** cards cannot be targeted by your opponent's cards. The Queen does not protect itself, does not block board-wipes (Ace, Six), and does not block scuttling. **A Queen does block Jacks**: a protected point card cannot be stolen. |
| **K** | Lowers your win threshold. Stacks: 1K=14, 2K=10, 3K=7, 4K=5. |

## Notes / Clarifications
- **2 vs. 2**: 2s can counter 2s. Counter chains resolve last-in-first-out.
- **Jacks vs. Queens**: a Queen on your side prevents your point cards from being stolen by Jacks.
- **Jack chain-steal**: a Jack may be played onto a point card that already has a Jack on it; control transfers to the new Jack's owner.
- **Glasses 8**: any 8, regardless of suit.
- **Win check**: only on your own turn, after your action resolves. Counters (2s) cannot grant points, so winning off-turn is not possible.
