# Edge Case Test Sweep — Brief for Executor Session

You are doing an edge-case test sweep on the Cuttle engine in this directory. The engine is complete (tagged `mvp-engine-v1`) and 95 unit tests pass. A human just finished a successful manual REPL playthrough. We want to harden the engine against rule edge cases before building a replay/evaluation harness on top of it.

## Authoritative references — read these FIRST

1. `RULES.md` — the rules of Cuttle. **Source of truth.** If anything in this brief contradicts `RULES.md`, the rules win.
2. `engine/state.go` and `engine/apply.go` — the engine itself.
3. `engine/*_test.go` — existing tests, so you don't duplicate coverage. In particular: `draw_test.go`, `playpoint_test.go`, `playpermanent_test.go`, `scuttle_test.go`, `ace_test.go`, `counter_test.go`, `six_test.go`, `three_test.go`, `four_test.go`, `five_test.go`, `nine_test.go`, `two_test.go`, `seven_test.go`, `jack_test.go`, `scenarios_test.go`.
4. `docs/plans/2026-04-06-mvp-engine.md` — Tasks 9–22 describe intended per-card semantics.

## Task

Create a new file `engine/edge_cases_test.go` covering the 36 cases below. Use `t.Run` to group cases by category (A through I). Reuse existing helpers (`twoPlayerStart` in `draw_test.go`, `assertDeckCount`/`padDeckTo52` in `scenarios_test.go`). **Do not modify existing test files.**

For each case: write the test → run it → observe pass or fail. If it fails, decide whether the **test** is wrong (fix it against `RULES.md`) or the **engine** is wrong (fix engine in `engine/apply.go` or wherever the bug lives, add a code comment explaining the fix). Each engine bug fix is its own commit.

## Cases (36 total)

### A. Hand-limit interactions (HandLimit = 8)

1. Hand at 7, play a 5 (draw 2): exactly 1 card drawn, hand ends at 8.
2. Hand at 8, play a 5: legal? Draws 0. `RULES.md` says "respecting the 8-card limit" — make a judgment call and document it in a code comment.
3. Hand at 8, play a 3 (take from scrap): would overflow the hand. Should be suppressed from `LegalMoves`. Verify.
4. Hand at 8 in `PhaseSevenChoosing`: the chosen card goes to the field or scrap, not into hand, so this should be legal even at hand=8.
5. After `MoveDraw` brings hand to exactly 8, on this player's next own turn, `MoveDraw` is not in `LegalMoves`.

### B. Deck-edge interactions

6. Deck empty + hand has playable cards: `MoveDraw` absent, `MovePass` absent (other moves exist), player must play.
7. Deck empty + hand has zero playable cards: `MovePass` is the only legal move.
8. 7 played with `len(Deck) == 0`: not legal at all (no reveals possible).
9. 7 played with `len(Deck) == 1`: only one revealed card; menu has exactly one `MoveSevenPick` option (or however the engine models the no-choice case — match the engine's model).
10. 5 with `len(Deck) == 0`: legal? Draws 0. Document.
11. 5 with `len(Deck) == 1`: draws 1.
12. Draw the very last deck card: deck size becomes 0, turn ends normally.

### C. Opponent-hand-size interactions

13. 4 (opponent discards 2) when opponent has 0 cards: legal? Auto-resolves with no discards. Document.
14. 4 when opponent has exactly 1 card: discards 1, transitions back to normal phase.
15. 4 when opponent has 2 cards: standard case, discards both.

### D. Empty-pile interactions

16. 3 with empty scrap: not legal as a one-off.
17. 2-as-scrap when no royals/glasses-8 on field: not legal in `PhaseNormal` (counter mode is unreachable from `PhaseNormal`).
18. 6 played when there are no permanents anywhere: legal? Trivially resolves. Document.
19. Ace played when there are no points anywhere: legal? Trivially resolves. Document.

### E. Counter-chain edge cases

20. Defender has 2 twos, attacker has 1 two: A → 2(d) → 2(a) → 2(d). Chain length 3, original A cancelled. All four cards in scrap.
21. A 2 frozen by an opponent's 9 cannot be played as a counter — it must not appear in `PhaseAwaitingCounter`'s legal moves.
22. Counter chain length 5 (rare but legal): A → 2 → 2 → 2 → 2 → 2. Parity odd → A cancelled.
23. Defender has no 2s and no other relevant resources: auto-resolve, no `PhaseAwaitingCounter` ever entered.

### F. Jack / Queen interactions

24. Jack-on-Jack-on-Jack chain steal. P1 plays J on P2's 7. P2 plays J on the same point. P1 plays J on the same point. `JackStack` length 3, controller correct (= P1).
25. Scrapping the top Jack from a 2-Jack stack (via 2-as-scrap): top Jack scraps, the next Jack down still controls the underlying point.
26. Six wiping a Jack stack: the underlying point returns to its `Owner`, all Jacks in the stack scrap.
27. Nine bouncing a point card under a Jack stack: per the plan, the underlying point returns to its `Owner`'s hand, the Jacks scrap.
28. Queen scrapped (by 2 or 6) while protecting points: those points become targetable (e.g. by Jack steal) on the opponent's next turn. Test by playing 2 to scrap Q, then verifying Jack steal targets appear in the opponent's `LegalMoves`.

### G. Win-check timing

29. Playing a King that drops your threshold below your current points → instant win on your turn.
30. 2-as-scrap targeting an opponent Jack to reclaim your stolen point that pushes you to threshold → instant win on your turn.
31. 7 reveals a card that, when played as a point card, wins the game → game over on the seven-pick resolution.
32. Jack chain steal that pushes the new controller over their threshold → instant win (the steal is the new controller's own action, so this is allowed).

### H. Frozen-card edge cases

33. A card frozen by an opponent's 9 cannot be played as point, permanent, scuttle, one-off, or counter — test all five paths.
34. The frozen marker clears at the start of the affected player's next own turn (not the opponent's). Verify with a careful sequence.

### I. Stalemate

35. 3 consecutive passes ends the game with `Phase == PhaseGameOver` and `Winner == nil`.
36. `MovePass` is not legal when other moves exist (e.g. when draw is available).

## Workflow

- Strict TDD: write test → `go test ./engine/... -run TestEdgeCases -count=1` → observe → fix.
- After completing each category (A through I), commit:
  - Test additions: `test(engine): edge cases — <category>`
  - Engine fixes: `fix(engine): <specific bug> — found by edge case sweep`
- When all 36 cases pass: run `go test ./... -count=1`, `go vet ./...`, and `gofmt -l .`. All must be clean.

## Constraints

- Do **not** modify files outside `engine/`.
- Do **not** modify `RULES.md`.
- Do **not** modify existing test files.
- If you find what you believe is a real rule ambiguity that should change `RULES.md`, **flag it in the summary**; do not edit `RULES.md` yourself.

## Output

When complete, return a summary under 400 words containing:
- Number of edge case tests added
- Bugs found, with one-line descriptions and the commit hash for each fix
- Judgment calls made against ambiguous rule wording (e.g. "5 with hand=8: I made it legal, drawing 0 cards, because RULES.md says 'respecting the limit', not 'forbidden when full'")
- Any cases that look suspicious to a human reviewer but you couldn't conclusively resolve

Begin work immediately. Read the files, write the tests, run them, fix bugs, commit per category, and report when complete.
