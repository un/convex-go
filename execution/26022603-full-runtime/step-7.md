# Step 7 - Implement state version and timestamp encoding

## What changed
- Updated `internal/protocol/types.go` so `StateVersion.TS` is typed as `protocol.Timestamp` and serialized/deserialized with custom JSON logic.
- Added strict timestamp decode validation in `StateVersion.UnmarshalJSON` with deterministic error wrapping.
- Hardened timestamp decode error messages in `internal/protocol/codec.go`.
- Updated runtime and tests (`convex/client.go`, `convex/client_test.go`) to use typed state version timestamps.
- Expanded codec tests in `internal/protocol/codec_test.go` for state-version wire shape and invalid timestamp rejection.
- Marked step 7 as completed in `plans/26022603-convex-go-full-runtime-parity-plan/steps.json`.
- Prepended step learning in `plans/26022603-convex-go-full-runtime-parity-plan/learnings.md`.

## Validations run
- `go test ./...`

## Outcomes
- PASS: `go test ./...`
