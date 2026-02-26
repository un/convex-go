# Rust Fixture Vectors

This directory contains protocol conformance fixture vectors imported from the Rust baseline.

## Source baseline

- Repo: `/Users/omar/code/convex-rs`
- Source file: `sync_types/src/types/json.rs`
- Commit hash is recorded in `rust_fixture_vectors.json`.

## Regeneration

Run:

```bash
./scripts/import_rust_protocol_fixtures.py
```

The script re-extracts vectors from the Rust baseline and overwrites `rust_fixture_vectors.json` deterministically.
