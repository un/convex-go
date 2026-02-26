# Step 8 - Implement typed Query and QuerySetModification models

## What changed
- Refactored `internal/protocol/types.go` `QuerySetModification` into explicit Add/Remove variants with constructors:
  - `NewQuerySetAdd(Query)`
  - `NewQuerySetRemove(QueryID)`
- Added strict custom JSON marshal/unmarshal logic enforcing:
  - exactly one variant active
  - required `udfPath` for Add
  - required `queryId` for Remove
  - unknown type rejection
- Updated runtime call sites (`convex/client.go`, `convex/client_test.go`) to use variant constructors/accessors.
- Expanded protocol tests (`internal/protocol/codec_test.go`) for Add/Remove roundtrip and malformed Add rejection.
- Marked step 8 as completed in `plans/26022603-convex-go-full-runtime-parity-plan/steps.json`.
- Prepended step learning in `plans/26022603-convex-go-full-runtime-parity-plan/learnings.md`.

## Validations run
- `go test ./...`

## Outcomes
- PASS: `go test ./...`
