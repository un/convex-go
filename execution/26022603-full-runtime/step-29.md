# Step 29 - Implement reconnect backoff jitter parity

## What changed
- Updated backoff implementation in `internal/baseclient/backoff.go` to Rust-style full-jitter exponential behavior:
  - exponential growth with overflow/cap handling
  - random jitter in `[0,1]`
  - default RNG switched to real randomness
  - added `SetFailures` for deterministic stateful testing
- Expanded tests in `internal/baseclient/backoff_test.go` for deterministic growth, cap semantics, and jitter parity behavior.
- Marked step 29 as completed in `plans/26022603-convex-go-full-runtime-parity-plan/steps.json`.
- Prepended step learning in `plans/26022603-convex-go-full-runtime-parity-plan/learnings.md`.

## Validations run
- `go test ./...`

## Outcomes
- PASS: `go test ./...`
