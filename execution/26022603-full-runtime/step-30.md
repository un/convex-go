# Step 30 - Implement websocket state callback ordering

## What changed
- Added centralized/deduplicated websocket state callback emission in `convex/client.go` via `emitState`.
- Routed lifecycle state transitions through explicit ordering points:
  - `connecting`
  - `connected`
  - `reconnecting`
  - `disconnected`
- Added callback ordering tests in `convex/client_test.go` for:
  - initial connection sequence
  - reconnect cycle sequence
- Marked step 30 as completed in `plans/26022603-convex-go-full-runtime-parity-plan/steps.json`.
- Prepended step learning in `plans/26022603-convex-go-full-runtime-parity-plan/learnings.md`.

## Validations run
- `go test ./...`

## Outcomes
- PASS: `go test ./...`
