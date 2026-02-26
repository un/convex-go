# Step 22 - Add codec fuzzing gate

## What changed
- Added fuzz targets in `internal/protocol/codec_fuzz_test.go`:
  - `FuzzDecodeClientMessage`
  - `FuzzDecodeServerMessage`
- Seeded fuzzing inputs from:
  - valid encoded protocol messages
  - `testdata/malformed_corpus.json`
  - `testdata/rust_fixture_vectors.json` (client decode vectors)
- Fuzz gates exercise decode -> optional re-encode paths and assert panic-free behavior.
- Marked step 22 as completed in `plans/26022603-convex-go-full-runtime-parity-plan/steps.json`.
- Prepended step learning in `plans/26022603-convex-go-full-runtime-parity-plan/learnings.md`.

## Validations run
- `go test ./...`
- `go test -race ./...`

## Outcomes
- PASS: `go test ./...`
- PASS: `go test -race ./...`
