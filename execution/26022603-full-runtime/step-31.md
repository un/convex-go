# Step 31 - Add transport integration tests

## What changed
- Added transport lifecycle integration test in `internal/sync/websocket_manager_test.go` (`TestWebSocketManagerIntegrationLifecycle`) covering open/send/read/reconnect flow with a controllable websocket server.
- Expanded transport test suite with decode/close/heartbeat/inactivity coverage already in place.
- Fixed race-detector failure in transport shutdown by removing `responses` channel close from `WebSocketManager.Close` and relying on closed-state gating.
- Marked step 31 as completed in `plans/26022603-convex-go-full-runtime-parity-plan/steps.json`.
- Prepended step learning in `plans/26022603-convex-go-full-runtime-parity-plan/learnings.md`.

## Validations run
- `go test ./...`
- `go test -race ./...`

## Outcomes
- PASS: `go test ./...`
- PASS: `go test -race ./...`
