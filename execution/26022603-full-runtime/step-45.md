# Step 45 - Wire watch_all, clone, and lifecycle semantics

## What changed
- Routed watch-all lifecycle through worker in `convex/client.go`:
  - `WatchAll` now dispatches `workerCommandWatchAll` and returns worker-owned watcher channel.
  - watch close path dispatches `workerCommandUnwatch` through worker.
  - added fallback path for pre-worker watch creation to avoid blocking before connection startup.
- Added worker handlers:
  - `handleWorkerWatchAll`
  - `handleWorkerUnwatch`
- Extended worker models in `convex/worker_types.go` with watch/unwatch payload/result types.
- Added lifecycle/coherency tests in `convex/client_test.go`:
  - `TestWatchAllSnapshotCoherentAcrossSubscriptions`
  - `TestCloneCloseLifecycleNoPanic`
- Marked step 45 as completed in `plans/26022603-convex-go-full-runtime-parity-plan/steps.json`.
- Prepended step learning in `plans/26022603-convex-go-full-runtime-parity-plan/learnings.md`.

## Validations run
- `go test ./...`

## Outcomes
- PASS: `go test ./...`
