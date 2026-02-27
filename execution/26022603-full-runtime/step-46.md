# Step 46 - Wire auth APIs through worker runtime

## What changed
- Routed auth API operations through worker commands in `convex/client.go`:
  - `SetAuth` now dispatches `workerCommandSetAuth` when worker is active.
  - `SetAuthCallback` now dispatches `workerCommandSetAuthCB` after immediate fetch.
  - retained pre-worker fallback path for pre-connect auth setup.
- Added worker auth handlers:
  - `handleWorkerSetAuth`
  - `handleWorkerSetAuthCallback`
  - both update token/identity state and enqueue authenticate payload when connected.
- Extended worker models in `convex/worker_types.go` for auth command payloads.
- Added auth worker routing tests in `convex/client_test.go`:
  - `TestSetAuthThroughWorkerEnqueuesAuthenticate`
  - `TestSetAuthCallbackThroughWorkerEnqueuesAuthenticate`
- Marked step 46 as completed in `plans/26022603-convex-go-full-runtime-parity-plan/steps.json`.
- Prepended step learning in `plans/26022603-convex-go-full-runtime-parity-plan/learnings.md`.

## Validations run
- `go test ./...`

## Outcomes
- PASS: `go test ./...`
