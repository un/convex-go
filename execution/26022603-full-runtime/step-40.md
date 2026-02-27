# Step 40 - Implement flush-before-select behavior

## What changed
- Added worker-managed outbound queue in `convex/client.go`:
  - `outboundQueue` state with `enqueueOutbound` wake-up signaling.
  - `flushOutboundBeforeSelect` sends all queued messages before entering normal worker `select`.
  - worker loop now drains outbound queue first on every iteration.
- Added `flushWake` signaling so queued outbound messages wake worker immediately.
- Routed reconnect replay (`replayState`) through outbound queue to preserve flush priority for auth/query/request replay ordering.
- Added `sendFn` injection hook used for deterministic worker flush tests.
- Added invariant test in `convex/worker_loop_test.go`:
  - `TestWorkerFlushesOutboundBeforeSelect`
- Marked step 40 as completed in `plans/26022603-convex-go-full-runtime-parity-plan/steps.json`.
- Prepended step learning in `plans/26022603-convex-go-full-runtime-parity-plan/learnings.md`.

## Validations run
- `go test ./...`

## Outcomes
- PASS: `go test ./...`
