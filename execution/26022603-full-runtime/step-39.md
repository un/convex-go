# Step 39 - Implement worker main loop

## What changed
- Replaced client listener loop with a dedicated worker loop in `convex/client.go`:
  - `workerLoop` now reads normalized transport events and worker commands.
  - transport responses are converted through `workerEventFromProtocolResponse` and routed through `handleWorkerEvent`.
  - worker command handling path added via `handleWorkerCommand`.
- Added worker lifecycle fields to `Client`:
  - `workerStarted`, `workerCommands`, and `workerDone`.
  - `ensureConnected` now boots the worker loop once per client.
- Updated client shutdown path to notify/wait for worker exit.
- Added loop-level tests in `convex/worker_loop_test.go`:
  - `TestWorkerLoopProcessesCommandAndTransportMessage`
  - `TestWorkerCommandUnsupportedKind`
- Marked step 39 as completed in `plans/26022603-convex-go-full-runtime-parity-plan/steps.json`.
- Prepended step learning in `plans/26022603-convex-go-full-runtime-parity-plan/learnings.md`.

## Validations run
- `go test ./...`

## Outcomes
- PASS: `go test ./...`
