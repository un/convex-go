# Step 16 - Implement strict client codec decode path

## What changed
- Added strict client decode envelope handling in `internal/protocol/codec.go`:
  - validates JSON shape before variant decode
  - requires message `type`
  - wraps variant decode errors with stable context
- Extended `ClientMessage` authenticate decode in `internal/protocol/types.go` with legacy compatibility fallback for old auth payload shape (`token`/`admin`/`actingAs`).
- Added protocol tests in `internal/protocol/codec_test.go` for:
  - strict envelope decode errors
  - legacy authenticate decode compatibility
- Marked step 16 as completed in `plans/26022603-convex-go-full-runtime-parity-plan/steps.json`.
- Prepended step learning in `plans/26022603-convex-go-full-runtime-parity-plan/learnings.md`.

## Validations run
- `go test ./...`

## Outcomes
- PASS: `go test ./...`
