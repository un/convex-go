# Step 12 - Implement remaining ClientMessage variants

## What changed
- Extended `internal/protocol/types.go` `ClientMessage` custom JSON encode/decode to strict variant handling for:
  - `ModifyQuerySet`
  - `Mutation`
  - `Action`
  - `Authenticate`
  - `Event`
- Embedded typed `AuthenticationToken` in authenticate client messages (`ClientMessage.Token`).
- Added required-field validation for each variant and unknown variant rejection.
- Updated runtime auth send path (`convex/client.go`) to emit typed token variants.
- Expanded protocol tests (`internal/protocol/codec_test.go`) for roundtrip coverage and variant validation failures.
- Marked step 12 as completed in `plans/26022603-convex-go-full-runtime-parity-plan/steps.json`.
- Prepended step learning in `plans/26022603-convex-go-full-runtime-parity-plan/learnings.md`.

## Validations run
- `go test ./...`

## Outcomes
- PASS: `go test ./...`
