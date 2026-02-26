# Step 14 - Implement remaining ServerMessage variants

## What changed
- Extended `ServerMessage` strict custom JSON handling in `internal/protocol/types.go` for:
  - `MutationResponse`
  - `ActionResponse`
  - `AuthError`
  - `FatalError`
  - `Ping`
- Enforced variant-required fields and response union validation (e.g. required `success`, required error payload for failed responses).
- Rejected unknown server message types at decode time.
- Expanded protocol tests (`internal/protocol/codec_test.go`) for roundtrip and malformed payload rejection across all remaining variants.
- Marked step 14 as completed in `plans/26022603-convex-go-full-runtime-parity-plan/steps.json`.
- Prepended step learning in `plans/26022603-convex-go-full-runtime-parity-plan/learnings.md`.

## Validations run
- `go test ./...`

## Outcomes
- PASS: `go test ./...`
