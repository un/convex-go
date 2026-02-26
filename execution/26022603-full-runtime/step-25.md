# Step 25 - Implement websocket write queue loop

## What changed
- Refactored `internal/sync/websocket_manager.go` send path to use a dedicated per-connection write queue and writer goroutine.
- `Send` now enqueues encoded payloads with non-blocking backpressure behavior (`websocket write queue full`).
- Added connection lifecycle handling for writer stop channels across open/reconnect/close/error paths.
- Added ordering regression test in `internal/sync/websocket_manager_test.go` (`TestWebSocketWriteQueuePreservesOrdering`).
- Marked step 25 as completed in `plans/26022603-convex-go-full-runtime-parity-plan/steps.json`.
- Prepended step learning in `plans/26022603-convex-go-full-runtime-parity-plan/learnings.md`.

## Validations run
- `go test ./...`

## Outcomes
- PASS: `go test ./...`
