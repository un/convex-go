# Step 32 - Finalize RequestManager semantics

## What changed
- Hardened request semantics in `internal/baseclient/request_manager.go`:
  - duplicate-add guard
  - mutation `WaitingOnTS` tracking to prevent premature completion
  - completion result reporting from `ApplyTransition`
  - helper accessors for testability
- Added `internal/baseclient/request_manager_test.go` covering:
  - mutation visibility-gated completion
  - action immediate completion
  - deterministic replay ordering behavior
- Marked step 32 as completed in `plans/26022603-convex-go-full-runtime-parity-plan/steps.json`.
- Prepended step learning in `plans/26022603-convex-go-full-runtime-parity-plan/learnings.md`.

## Validations run
- `go test ./...`

## Outcomes
- PASS: `go test ./...`
