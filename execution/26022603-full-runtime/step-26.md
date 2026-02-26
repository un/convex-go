# Step 26 - Implement websocket read loop

## What changed
- Updated websocket reader in `internal/sync/websocket_manager.go` to classify failures explicitly:
  - close-frame classification
  - generic read failures
  - unsupported frame-type detection
  - protocol decode failure wrapping
- Added close/decode classification tests in `internal/sync/websocket_manager_test.go`.
- Marked step 26 as completed in `plans/26022603-convex-go-full-runtime-parity-plan/steps.json`.
- Prepended step learning in `plans/26022603-convex-go-full-runtime-parity-plan/learnings.md`.

## Validations run
- `go test ./...`

## Outcomes
- PASS: `go test ./...`
