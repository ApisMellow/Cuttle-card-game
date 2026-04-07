# Next Phase Roadmap — After MVP Engine Acceptance

**Status:** Drafted 2026-04-06, immediately after the executor session completed `mvp-engine-v1` and all 95 automated tests passed. This roadmap takes effect **only if** the manual REPL playthrough (the one acceptance test the automated suite cannot perform) succeeds and the engine "feels right" to play.

If the playthrough exposes bugs, those become Phase 0 — fix them as failing tests in the appropriate `engine/*_test.go` file before doing anything below.

---

## Why this roadmap exists

Two things came up at the end of the MVP review that need to shape what comes next:

1. **The adversary scenario test is a regression witness, not independent confirmation.** It was hand-encoded from our design conversation, so the engine passing it only proves the engine produces today's expected output on today's expected input. It does not prove the engine plays well in general, and it does not prove a Claude-driven decision process plays well either.

2. **The REPL is genuinely random** (Go ≥1.20 auto-seeds `math/rand` per process), but **nothing about the game is reproducible**. The moment we want to compare two decision policies, or replay a position, or debug a "you should have seen what happened" moment, we have no seed to grab.

Both observations point at the same next move: before we put a Claude-driven player in the loop, we need a deterministic, replayable harness to evaluate it inside.

---

## Phase 1: Determinism + replay infrastructure

**Goal:** Any Cuttle game can be reproduced bit-for-bit from a small piece of input, and any game can be saved as a replay.

### 1a. Seedable shuffle
- Add a `--seed N` flag to `cmd/cuttle`. When provided, build the deck shuffle from a `rand.New(rand.NewSource(N))` instead of the global source.
- When omitted, generate a random seed at startup, **print it on the first screen of the game**, and use it. This way every game is reproducible after the fact, even ones you didn't plan to reproduce.
- Add a unit test that the same seed produces the same initial `GameState`.

### 1b. Move log
- Extend the REPL to record `(state-before, move-chosen)` to an in-memory slice as the game progresses.
- On game end, write the seed + the sequence of move indices (not the move structs — just the integers the player typed) to a JSON file. Index-based replay is robust because it doesn't depend on `Move` struct layout staying stable.
- Add a `--replay path.json` flag: load a replay file, drive the game by feeding the recorded indices into `Apply` one at a time, render each step. Useful for "show me how that game went" and for regression-testing the engine against known-good histories.

### 1c. Deterministic test for replay round-trip
- Generate a game with a fixed seed, play 30 moves picking move index 0 each time, save to JSON, load it back, replay, assert final `GameState` is identical.

**Deliverable:** any game can be reproduced from a seed + a list of integers. This is the substrate everything else in this roadmap stands on.

---

## Phase 2: Engine state projection for external decision-makers

**Goal:** A clean, stable, JSON-serializable view of `GameState` that an external process (a Claude Code subagent, a headless `claude` invocation, a script, anything) can read and reason about without depending on Go internals.

### 2a. Define the projection schema
- A new package `engine/view` (or similar) with a `View` struct: plain fields, JSON tags, no methods. Fields cover everything a player needs to make a legal decision: both hands (open-handed for now), both point arrays with Jack stacks, both permanent arrays, scrap pile contents, deck count (NOT contents), active player, phase, pending one-off, win threshold for each player.
- Include a parallel `MoveView` for each entry in `LegalMoves`: kind, human-readable description, and the integer index used to reference it.
- The schema is documented in `docs/json-schema.md` with a worked example matching the adversary scenario.

### 2b. Pure projection function
- `func Project(s engine.GameState) View` — pure, deterministic, no side effects, no I/O. The same `GameState` always produces the same `View`.
- Round-trip is **not** required (we're not deserializing back into `GameState`). The view is read-only.

### 2c. CLI subcommand
- Add `cuttle export --seed N --moves 1,3,2,4` that runs a partial game from a seed and a move sequence and prints the resulting view as JSON to stdout. This is the integration point that any external decision process will use.
- Add a corresponding `cuttle apply --seed N --moves 1,3,2,4 --pick K` that picks move index K from the resulting state and prints the new view. This lets a decision process say "what happens if I play option 2?" without embedding the engine.

**Deliverable:** an external process can ask "what's the current state?" and "what would happen if I picked move K?" via shell commands, with no Go knowledge.

---

## Phase 3: Headless Claude as a player

**Goal:** Let a Claude Code instance running in headless mode read a `View`, pick a legal move, and play a real game.

### 3a. Decision prompt
- Write a single self-contained prompt template that takes a `View` as JSON and the rules (`RULES.md`) and asks: "Pick the integer index of your move, and explain why in two sentences."
- The prompt forces a specific output format (e.g. `MOVE: 3` on its own line) so a wrapping script can parse it without Claude prose ambiguity.
- Iterate the prompt against the adversary scenario as a sanity check: does Claude pick the winning move? If not, the prompt — or the projection — needs work before continuing.

### 3b. Game-driver script
- A small Go (or shell) program that:
  1. Starts a game with a seed.
  2. On each turn, projects the state to JSON, invokes `claude -p "..."` headless with the prompt + JSON + rules, parses the `MOVE: N` line.
  3. Validates that N is in `LegalMoves` (defense against malformed Claude output — fall back to a deterministic policy if invalid).
  4. Applies the move, logs it, loops.
- Run modes: Claude vs. random, Claude vs. always-draw, Claude vs. Claude (self-play).

### 3c. Evaluation harness
- Run N games (start small: 10) for each match-up, fixed seeds so re-runs are comparable.
- Record win rate, average game length, and the full move logs.
- Manually spot-check 2-3 logs per match-up for "did Claude make any obviously dumb moves?" — this is the qualitative check that automated win rates can't replace.

**Deliverable:** an answer to the original question that started this whole project — *"would a thinking model with the rules in its prompt actually be a competent Cuttle opponent without specialized harness?"*

---

## Phase 4 and beyond (not planning yet)

If Phase 3 shows Claude is a competent player, the natural follow-ups are:
- A nicer display (Bubble Tea, reusing the `Card-Game` ASCII renderer)
- Hidden hands and per-player projections (filter the `View` based on `glasses-8` permanents)
- Variant rule support (configurable thresholds, alternate counter rules)
- Maybe a web UI

If Phase 3 shows Claude is a weak player, the interesting follow-ups are:
- Why? Bad rule comprehension? Bad lookahead? Bad probability intuition?
- Do better prompts close the gap (chain-of-thought, list-and-evaluate)?
- Does giving Claude a small lookahead helper (e.g. "here are the resulting `View`s for each candidate move") improve play substantially?

These don't need plans yet. We'll know which question to ask after Phase 3 runs.

---

## What to bring back here when starting Phase 1

If the manual playthrough was clean, the next session can start from this document directly. Useful things to have on hand:
- Confirmation the playthrough was clean (or a list of bugs to file as Phase 0)
- The `mvp-engine-v1` tag is already in place; Phase 1 work happens on a new branch (`feature/replay-infra` or similar)
- This roadmap stays here as a living document — update it as the picture clarifies
