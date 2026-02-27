# Step 42 - Implement protocol-failure reconnect propagation

## What changed
- Hardened reconnect payload propagation in `convex/client.go`:
  - `onProtocolFailure` now guards nil errors and always emits a concrete reconnect reason.
  - reconnect max-observed timestamp now uses `maxObservedTimestampLocked`, combining observed transition timestamp with pending mutation visibility timestamps.
  - reconnect path now supports injected reconnect function (`reconnectFn`) for deterministic propagation tests.
- Added failure propagation tests in `convex/client_test.go`:
  - `TestProtocolFailureReconnectPayloadIncludesReasonAndMaxObservedTimestamp`
  - `TestFailureClassesPropagateReconnectReason`
  - validates reason propagation for auth/fatal/unknown/transport protocol failures.
- Marked step 42 as completed in `plans/26022603-convex-go-full-runtime-parity-plan/steps.json`.
- Prepended step learning in `plans/26022603-convex-go-full-runtime-parity-plan/learnings.md`.

## Validations run
- `go test ./...`
- `go test -race ./...`

## Outcomes
- PASS: `go test ./...`
- PASS: `go test -race ./...`
