# Step 9 - Implement typed StateModification models

## What changed
- Refactored `internal/protocol/types.go` `StateModification` into explicit variants:
  - `QueryUpdated`
  - `QueryFailed`
  - `QueryRemoved`
- Added strict variant constructors/accessors and custom JSON marshal/unmarshal with malformed-variant rejection.
- Updated transition apply usage in `convex/client.go` to consume typed variant accessors.
- Updated integration-style client test fixture building in `convex/client_test.go` to use typed constructors.
- Expanded protocol tests in `internal/protocol/codec_test.go` for variant roundtrip behavior and malformed decode rejection.
- Marked step 9 as completed in `plans/26022603-convex-go-full-runtime-parity-plan/steps.json`.
- Prepended step learning in `plans/26022603-convex-go-full-runtime-parity-plan/learnings.md`.

## Validations run
- `go test ./...`

## Outcomes
- PASS: `go test ./...`
