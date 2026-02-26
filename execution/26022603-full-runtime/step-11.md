# Step 11 - Implement ClientMessage Connect variant

## What changed
- Added strict `ClientMessage` custom JSON marshal/unmarshal logic for `Connect` in `internal/protocol/types.go`.
- Enforced required connect fields (`sessionId`, `connectionCount`) and session ID validation.
- Added connect defaulting behavior for missing `lastCloseReason` -> `unknown`.
- Added validation for `maxObservedTimestamp` wire encoding when present.
- Expanded protocol tests (`internal/protocol/codec_test.go`) for:
  - connect roundtrip and defaults
  - required-field rejection
  - invalid timestamp rejection
- Marked step 11 as completed in `plans/26022603-convex-go-full-runtime-parity-plan/steps.json`.
- Prepended step learning in `plans/26022603-convex-go-full-runtime-parity-plan/learnings.md`.

## Validations run
- `go test ./...`

## Outcomes
- PASS: `go test ./...`
