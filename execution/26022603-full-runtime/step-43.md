# Step 43 - Wire public Subscribe and Query through worker

## What changed
- Routed `Subscribe` through worker command dispatch in `convex/client.go`:
  - added `executeWorkerCommand` helper for typed command dispatch/response.
  - `Subscribe` now issues `workerCommandSubscribe` and receives worker-owned subscription state.
  - subscription close path now dispatches `workerCommandUnsubscribe` through worker.
- Implemented worker handlers for query APIs:
  - `handleWorkerSubscribe`
  - `handleWorkerUnsubscribe`
  - uses existing local sync state and query-set message building under worker ownership.
- Added worker payload/result models in `convex/worker_types.go` for subscribe/unsubscribe commands.
- `Query` continues to rely on subscribe-first-value behavior, now via worker-routed subscription path.
- Marked step 43 as completed in `plans/26022603-convex-go-full-runtime-parity-plan/steps.json`.
- Prepended step learning in `plans/26022603-convex-go-full-runtime-parity-plan/learnings.md`.

## Validations run
- `go test ./...`

## Outcomes
- PASS: `go test ./...`
