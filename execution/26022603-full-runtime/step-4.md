# Step 4 - Add no-scaffold CI guard

## What changed
- Added `scripts/no_scaffold_guard.sh` to fail on forbidden runtime scaffold markers in core runtime paths.
- Added `scripts/no_scaffold_guard.allowlist` for explicit temporary exceptions.
- Wired guard into CI in `.github/workflows/ci.yml`.
- Marked step 4 as completed in `plans/26022603-convex-go-full-runtime-parity-plan/steps.json`.
- Prepended step learning in `plans/26022603-convex-go-full-runtime-parity-plan/learnings.md`.

## Validations run
- `./scripts/no_scaffold_guard.sh`
- `go test ./...`

## Outcomes
- PASS: `./scripts/no_scaffold_guard.sh`
- PASS: `go test ./...`
