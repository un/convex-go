# Step 44 - Wire public Mutation and Action through worker

## What changed
- Routed mutation/action API request creation through worker in `convex/client.go`:
  - `runRequest` now dispatches `workerCommandMutation` / `workerCommandAction`.
  - worker returns typed request handle (`requestID` + response channel).
  - context cancellation now dispatches `workerCommandCancelReq` to clear pending request state deterministically.
- Added worker handlers:
  - `handleWorkerRunRequest`
  - `handleWorkerCancelRequest`
- Extended worker command model in `convex/worker_types.go` with request/cancel payload/result types.
- Added cancellation regression test in `convex/client_test.go`:
  - `TestMutationContextCancellationRemovesPendingRequest`
- Marked step 44 as completed in `plans/26022603-convex-go-full-runtime-parity-plan/steps.json`.
- Prepended step learning in `plans/26022603-convex-go-full-runtime-parity-plan/learnings.md`.

## Validations run
- `go test ./...`

## Outcomes
- PASS: `go test ./...`
