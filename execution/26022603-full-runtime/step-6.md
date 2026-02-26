# Step 6 - Implement strict protocol identifiers

## What changed
- Added typed identifier/version conversion helpers in `internal/protocol/identifiers.go`:
  - `QueryIDFromUint64`, `QuerySetVersionFromUint64`, `IdentityVersionFromUint64`, `RequestSequenceNumberFromUint64`
  - helper methods (`Uint32`, `Uint64`) and session ID validation helpers.
- Updated runtime usage in `convex/client.go` to use conversion helpers instead of raw casts at protocol boundaries.
- Updated websocket session ID construction in `internal/sync/websocket_manager.go` to use validated `protocol.MustSessionID`.
- Added tests in `internal/protocol/identifiers_test.go` covering conversion success, overflow rejection, and session ID validation.
- Marked step 6 as completed in `plans/26022603-convex-go-full-runtime-parity-plan/steps.json`.
- Prepended step learning in `plans/26022603-convex-go-full-runtime-parity-plan/learnings.md`.

## Validations run
- `go test ./...`

## Outcomes
- PASS: `go test ./...`
