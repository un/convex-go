# Step 19 - Add malformed protocol corpus tests

## What changed
- Added malformed protocol corpus file: `internal/protocol/testdata/malformed_corpus.json`.
- Added corpus-driven test in `internal/protocol/codec_test.go` (`TestMalformedProtocolCorpus`) covering:
  - missing required fields
  - invalid unions
  - unknown variants
  - client and server decode paths
- Added error substring assertions for failure classification stability.
- Marked step 19 as completed in `plans/26022603-convex-go-full-runtime-parity-plan/steps.json`.
- Prepended step learning in `plans/26022603-convex-go-full-runtime-parity-plan/learnings.md`.

## Validations run
- `go test ./...`

## Outcomes
- PASS: `go test ./...`
