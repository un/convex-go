# Learnings

## 2026-02-26 Step 5 document-target-runtime-architecture
- Context: Parity implementation touches transport, state, and API simultaneously, so ownership boundaries must be explicit before refactors.
- Learning: Defining single-writer ownership (worker for state transitions, transport for websocket internals) prevents race-prone split responsibility.
- Impact on next steps: Protocol and transport refactors should enforce this boundary instead of adding more logic to `convex/client.go`.

## 2026-02-26 Step 4 add-no-scaffold-ci-guard
- Context: Runtime parity work needs an automated guard to prevent regressions back to scaffold markers.
- Learning: A marker-based guard with an explicit allowlist keeps CI strict while still allowing staged removal of known legacy paths.
- Impact on next steps: As each scaffold inventory row closes, remove its allowlist entry so step 48 can run with zero exceptions.

## 2026-02-26 Step 3 freeze-scaffold-inventory
- Context: Parity work spans many files; without a locked inventory, scaffold removal can be claimed without closure evidence.
- Learning: Assigning stable inventory IDs and mandatory closure metadata makes scaffold-zero audit mechanically checkable.
- Impact on next steps: Each runtime implementation step should close specific inventory rows with evidence references.

## 2026-02-26 Step 2 generate-gap-matrix
- Context: Existing Go runtime behavior is spread across protocol, transport, base state, and API layers.
- Learning: A function-level matrix prevents "hidden" runtime shortcuts by forcing each parity claim to cite a concrete Go location and Rust counterpart.
- Impact on next steps: Scaffold inventory and implementation commits should close matrix rows explicitly in step evidence.

## 2026-02-26 Step 1 pin-rust-runtime-baseline
- Context: Parity work needs one immutable Rust source-of-truth commit for protocol and transport behavior.
- Learning: Capturing both commit hash and concrete module list up front prevents accidental cross-version mixing during later codec and runtime changes.
- Impact on next steps: Gap analysis and scaffold inventory must explicitly trace each item back to `rust-baseline.md` references.

This file starts empty by design.

## Agent workflow rule
- Before any step starts, read this file.
- After each step, prepend new learnings at the top.
- Keep entries short, concrete, and action-oriented.

## Entry template
```
## YYYY-MM-DD Step <serial> <short-name>
- Context:
- Learning:
- Impact on next steps:
```

No learnings recorded yet.
