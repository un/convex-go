# Step 41 - Implement communicate-during-send behavior

## What changed
- Added `sendWhileCommunicating` in `convex/client.go`:
  - outbound send now runs asynchronously.
  - while waiting for send completion, worker continues handling transport events and commands.
  - closed/stream-close conditions are surfaced through worker event handling.
- Updated outbound flush path to use `sendWhileCommunicating`, preserving flush ordering while avoiding inbound starvation.
- Added worker event hook support (`workerEventHook`) used by tests to assert in-send communication behavior.
- Added regression test in `convex/worker_loop_test.go`:
  - `TestWorkerCommunicatesDuringSend`
- Marked step 41 as completed in `plans/26022603-convex-go-full-runtime-parity-plan/steps.json`.
- Prepended step learning in `plans/26022603-convex-go-full-runtime-parity-plan/learnings.md`.

## Validations run
- `go test ./...`

## Outcomes
- PASS: `go test ./...`
