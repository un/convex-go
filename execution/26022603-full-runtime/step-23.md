# Step 23 - Implement websocket dial handshake

## What changed
- Hardened websocket dial path in `internal/sync/websocket_manager.go`:
  - explicit pre-dial context cancellation check
  - detailed handshake error formatting with HTTP status/body when available
  - cancellation-aware dial error classification
- Added transport tests in `internal/sync/websocket_manager_test.go` for:
  - handshake failure detail propagation
  - context cancellation behavior
- Marked step 23 as completed in `plans/26022603-convex-go-full-runtime-parity-plan/steps.json`.
- Prepended step learning in `plans/26022603-convex-go-full-runtime-parity-plan/learnings.md`.

## Validations run
- `go test ./...`

## Outcomes
- PASS: `go test ./...`
