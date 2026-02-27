# Step 38 - Define worker command and event model

## What changed
- Added explicit worker command/event model in `convex/worker_types.go`:
  - typed command kinds covering API/lifecycle operations.
  - `workerCommand` with context-aware cancellation and completion semantics.
  - typed event kinds for command dispatch and transport message/error/done routing.
  - conversion helper `workerEventFromProtocolResponse` for normalized transport event handling.
- Added model-level tests in `convex/worker_types_test.go`:
  - `TestWorkerCommandCancellationSemantics`
  - `TestWorkerEventFromProtocolResponse`
- Marked step 38 as completed in `plans/26022603-convex-go-full-runtime-parity-plan/steps.json`.
- Prepended step learning in `plans/26022603-convex-go-full-runtime-parity-plan/learnings.md`.

## Validations run
- `go test ./...`

## Outcomes
- PASS: `go test ./...`
