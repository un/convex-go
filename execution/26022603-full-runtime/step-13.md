# Step 13 - Implement ServerMessage Transition and Chunk variants

## What changed
- Added strict `ServerMessage` custom JSON encode/decode handling for:
  - `Transition`
  - `TransitionChunk`
- Enforced required fields and invariants:
  - Transition requires `startVersion`, `endVersion`, and `modifications`.
  - TransitionChunk requires `chunk`, `partNumber`, `totalParts`, `transitionId` and validates part range.
- Added transition/chunk protocol tests in `internal/protocol/codec_test.go` for roundtrip and malformed payload rejection.
- Marked step 13 as completed in `plans/26022603-convex-go-full-runtime-parity-plan/steps.json`.
- Prepended step learning in `plans/26022603-convex-go-full-runtime-parity-plan/learnings.md`.

## Validations run
- `go test ./...`

## Outcomes
- PASS: `go test ./...`
