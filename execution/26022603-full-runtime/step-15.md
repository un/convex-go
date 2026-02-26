# Step 15 - Implement strict client codec encode path

## What changed
- Added explicit client encode validation gate in `internal/protocol/codec.go` (`validateClientMessageForEncode`).
- `EncodeClientMessage` now fails fast on invalid variant shapes before JSON marshaling.
- Added protocol test coverage in `internal/protocol/codec_test.go` for encode wire-key expectations on authenticate payloads.
- Marked step 15 as completed in `plans/26022603-convex-go-full-runtime-parity-plan/steps.json`.
- Prepended step learning in `plans/26022603-convex-go-full-runtime-parity-plan/learnings.md`.

## Validations run
- `go test ./...`

## Outcomes
- PASS: `go test ./...`
