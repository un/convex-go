# Learnings

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
