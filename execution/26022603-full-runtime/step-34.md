# Step 34 - Implement strict transition apply logic

## What changed
- Hardened transition apply path in `convex/client.go`:
  - requires both `startVersion` and `endVersion`
  - validates transition start version against last applied transition
  - triggers reconnect flow on version mismatch
  - tracks last applied transition end version
- Added mismatch regression test in `convex/client_test.go` (`TestTransitionVersionMismatchTriggersReconnect`).
- Marked step 34 as completed in `plans/26022603-convex-go-full-runtime-parity-plan/steps.json`.
- Prepended step learning in `plans/26022603-convex-go-full-runtime-parity-plan/learnings.md`.

## Validations run
- `go test ./...`

## Outcomes
- PASS: `go test ./...`
