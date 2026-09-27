# AGENTS.md

Guidance for agent sessions working in this repo. `RULES.md` is the
authoritative rulebook; `docs/plans/` holds the original implementation
plan. Module path is `github.com/ApisMellow/cuttle` (GitHub repo
`ApisMellow/cuttle`, releases tagged semver from `v0.1.0`).

## This engine is the rules canon

`cuttle-web` (the web client) compiles this package to WASM and implements zero game
rules itself — every rule question is answered by `LegalMoves`/`Apply`
output. Rule changes happen here, never in the UI. `cuttle-web` consumes
tagged releases of this module; locally it may also use a `go.mod`
`replace` pointing at this sibling checkout, so unreleased changes here
reach it immediately.

**`LegalMoves` (`engine/apply.go:15`) enumeration order is a stability
hazard**, even though the engine gives no ordering guarantee: `cuttle-web`'s
replay scenarios index into that slice, so reordering how moves are
appended silently repoints its scenarios (their `expect` guards catch it).

## Layout

- `card/` — suits, ranks, cards; no external deps
- `engine/` — immutable `GameState`, `LegalMoves`, `Apply`; one test file
  per rank/phase plus `edge_cases_test.go`, `scenarios_test.go`,
  `random_playout_test.go`
- `render/` — pure text formatter for the REPL
- `cmd/cuttle` — hot-seat REPL; every menu option is a `LegalMoves` entry

## Commands

- `go test ./...` — the gate; test-first, test files exist before/alongside implementations.
- `go vet ./...` static checks, `gofmt -l .` formatting (empty output = clean).
- `go run ./cmd/cuttle` to play, or `go build -o cuttle ./cmd/cuttle`.
- CI (`.github/workflows/ci.yml`) runs all three above on every push to `main` and every PR.
- `engine/random_playout_test.go` (`TestRandomPlayouts`) plays 500 seeded random
  games through `LegalMoves`/`Apply`, checking conservation and that every
  offered move applies — the harness that surfaced E-1/E-2 below.

## E-1 / E-2: fixed, verified in current code

Two `PhaseSevenChoosing` defects (`cuttle-web`'s `docs/SPEC.md` §2.10,
`docs/loop-log/engine-issues.md`) made ~1.9% of random games unwinnable.
**E-1** — a stale `FrozenIDs` hand index rejected every offered seven pick.
**E-2** — `LegalMoves` returned empty when every revealed card was an
unplayable Jack. Both fixed as the code stands:

- E-1: `removeFromHand` (`engine/apply.go:784`) remaps `FrozenIDs` keys on
  hand removal instead of leaving stale indices. `TestSeven_FrozenIndexRemappedWhenSevenPlayed`
  (`engine/seven_test.go:209`) reproduces the SPEC §2.10 case exactly and passes.
- E-2: `legalSevenPickMoves` (`engine/apply.go:524`) falls back to dead-end
  scrap moves instead of returning empty. `TestSeven_NoLegalPlayForRevealed_ScrapsChosen`
  (`engine/seven_test.go:242`) uses the exact two-Jacks case; confirms 2 moves, not 0.

No git commands from an agent session here without explicit instruction —
commits/pushes go through a separate git-focused agent.
