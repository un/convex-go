# Step 10 - Implement typed AuthenticationToken variants

## What changed
- Refactored `internal/protocol/types.go` `AuthenticationToken` into strict `Admin` / `User` / `None` variants.
- Added constructors and accessors:
  - `NewAdminAuthenticationToken`
  - `NewUserAuthenticationToken`
  - `NewNoAuthenticationToken`
- Implemented strict JSON marshal/unmarshal with malformed token rejection.
- Added backward-compatible decode support for legacy `impersonating` alias mapped to `actingAs`.
- Expanded protocol tests (`internal/protocol/codec_test.go`) for token roundtrips and compatibility decode behavior.
- Marked step 10 as completed in `plans/26022603-convex-go-full-runtime-parity-plan/steps.json`.
- Prepended step learning in `plans/26022603-convex-go-full-runtime-parity-plan/learnings.md`.

## Validations run
- `go test ./...`

## Outcomes
- PASS: `go test ./...`
