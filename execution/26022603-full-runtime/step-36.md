# Step 36 - Finalize auth callback refresh semantics

## What changed
- Hardened auth callback semantics in `convex/client.go`:
  - `SetAuthCallback` now rejects nil callbacks.
  - `SetAuth` uses a shared locked auth-token update helper.
  - Reconnect flow now treats `authFetcher(true)` failure as retryable reconnect failure (backoff + retry).
  - Reconnect flow applies refreshed auth token to local state without sending `Authenticate` before reconnect completes.
- Added auth semantics tests in `convex/client_test.go`:
  - `TestSetAuthCallbackRequiresFetcher`
  - `TestSetAuthCallbackReconnectForceRefreshAndRetry`
- Marked step 36 as completed in `plans/26022603-convex-go-full-runtime-parity-plan/steps.json`.
- Prepended step learning in `plans/26022603-convex-go-full-runtime-parity-plan/learnings.md`.

## Validations run
- `go test ./...`

## Outcomes
- PASS: `go test ./...`
