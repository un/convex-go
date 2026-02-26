# Step 20 - Import Rust fixture vectors

## What changed
- Added deterministic fixture import script: `scripts/import_rust_protocol_fixtures.py`.
- Imported Rust-linked fixture vectors into `internal/protocol/testdata/rust_fixture_vectors.json`.
  - Includes Rust source commit hash and source file metadata.
  - Includes extracted auth compatibility vectors and timestamp vectors.
- Added fixture usage documentation: `internal/protocol/testdata/RUST_FIXTURES.md`.
- Marked step 20 as completed in `plans/26022603-convex-go-full-runtime-parity-plan/steps.json`.
- Prepended step learning in `plans/26022603-convex-go-full-runtime-parity-plan/learnings.md`.

## Validations run
- `go test ./...`

## Outcomes
- PASS: `go test ./...`
