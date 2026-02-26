# Step 28 - Implement inactivity watchdog

## What changed
- Added server inactivity tracking in `internal/sync/websocket_manager.go`:
  - records last server response timestamp
  - heartbeat loop checks inactivity threshold and emits `InactiveServer` failure
  - closes stale connections on inactivity failure
- Added watchdog regression test in `internal/sync/websocket_manager_test.go` (`TestWebSocketInactivityWatchdogTriggersFailure`).
- Marked step 28 as completed in `plans/26022603-convex-go-full-runtime-parity-plan/steps.json`.
- Prepended step learning in `plans/26022603-convex-go-full-runtime-parity-plan/learnings.md`.

## Validations run
- `go test ./...`

## Outcomes
- PASS: `go test ./...`
