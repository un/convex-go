# Step 18 - Implement strict server codec decode path

## What changed
- Added strict server decode envelope handling in `internal/protocol/codec.go`:
  - requires top-level message `type`
  - wraps decode failures with variant-specific context
- Tightened server variant decode semantics in `internal/protocol/types.go`:
  - rejects `ActionResponse` payloads that include `ts`
  - validates `MutationResponse.ts` wire encoding when present
- Added protocol tests (`internal/protocol/codec_test.go`) for:
  - decode envelope errors
  - null optional-field handling for auth errors
  - malformed action-response timestamp rejection
- Marked step 18 as completed in `plans/26022603-convex-go-full-runtime-parity-plan/steps.json`.
- Prepended step learning in `plans/26022603-convex-go-full-runtime-parity-plan/learnings.md`.

## Validations run
- `go test ./...`

## Outcomes
- PASS: `go test ./...`
