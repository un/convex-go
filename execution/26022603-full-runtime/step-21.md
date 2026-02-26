# Step 21 - Run fixture conformance gate

## What changed
- Added `TestRustFixtureConformance` in `internal/protocol/codec_test.go`.
- Gate validates:
  - timestamp encode/decode parity against imported Rust vectors
  - client decode/re-encode compatibility for Rust-derived auth vectors
- Conformance gate now runs in CI via existing `go test ./...` workflow execution.
- Marked step 21 as completed in `plans/26022603-convex-go-full-runtime-parity-plan/steps.json`.
- Prepended step learning in `plans/26022603-convex-go-full-runtime-parity-plan/learnings.md`.

## Validations run
- `go test ./...`

## Outcomes
- PASS: `go test ./...`
