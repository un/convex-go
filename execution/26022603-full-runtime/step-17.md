# Step 17 - Implement strict server codec encode path

## What changed
- Added explicit server encode validation gate in `internal/protocol/codec.go` (`validateServerMessageForEncode`).
- `EncodeServerMessage` now rejects malformed transitions, transition chunks, and invalid response unions before serialization.
- Added validation for mutation-response timestamp wire encoding.
- Expanded protocol tests in `internal/protocol/codec_test.go` to assert encode-time failure behavior.
- Marked step 17 as completed in `plans/26022603-convex-go-full-runtime-parity-plan/steps.json`.
- Prepended step learning in `plans/26022603-convex-go-full-runtime-parity-plan/learnings.md`.

## Validations run
- `go test ./...`

## Outcomes
- PASS: `go test ./...`
