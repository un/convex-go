# Step 47 - Expand live integration suite

## What changed
- Reworked `integration/integration_test.go` into scenario-based live subtests:
  - `query`
  - `subscribe`
  - `mutation` (env-gated)
  - `action` (env-gated)
  - `auth-refresh-callback` (env-gated with optional strict expectation)
  - `reconnect-probe` (env-gated with optional strict expectation)
- Added shared live config loader and client factory:
  - `loadLiveConfig`
  - `newLiveClient`
  - supports reconnect probe duration and strict expectation toggles via env vars.
- Marked step 47 as completed in `plans/26022603-convex-go-full-runtime-parity-plan/steps.json`.
- Prepended step learning in `plans/26022603-convex-go-full-runtime-parity-plan/learnings.md`.

## Validations run
- `go test ./...`
- `CONVEX_INTEGRATION=1 go test ./integration -count=1`

## Outcomes
- PASS: `go test ./...`
- PASS: `CONVEX_INTEGRATION=1 go test ./integration -count=1`
