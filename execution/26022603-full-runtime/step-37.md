# Step 37 - Implement replay queue rebuild

## What changed
- Updated reconnect replay construction in `convex/client.go`:
  - replayed pending requests are now rebuilt from sorted request IDs (deterministic order).
  - replayed query-set modifications are now rebuilt from sorted query IDs (deterministic order).
  - replay flow remains ordered as `Authenticate` -> rebuilt `ModifyQuerySet` -> rebuilt in-flight requests.
- Added reconnect replay integration coverage in `convex/client_test.go`:
  - `TestReconnectReplayOrderAuthQueriesThenPendingRequests`
  - validates second-connection replay ordering is deterministic and complete.
- Marked step 37 as completed in `plans/26022603-convex-go-full-runtime-parity-plan/steps.json`.
- Prepended step learning in `plans/26022603-convex-go-full-runtime-parity-plan/learnings.md`.

## Validations run
- `go test ./...`

## Outcomes
- PASS: `go test ./...`
