# Step 27 - Implement heartbeat ping ticker

## What changed
- Added heartbeat loop to transport (`internal/sync/websocket_manager.go`) that sends websocket ping control frames on a fixed interval.
- Integrated heartbeat stop lifecycle with open/reconnect/close/read-error paths.
- Added heartbeat test in `internal/sync/websocket_manager_test.go` (`TestWebSocketHeartbeatSendsPingFrames`).
- Marked step 27 as completed in `plans/26022603-convex-go-full-runtime-parity-plan/steps.json`.
- Prepended step learning in `plans/26022603-convex-go-full-runtime-parity-plan/learnings.md`.

## Validations run
- `go test ./...`

## Outcomes
- PASS: `go test ./...`
