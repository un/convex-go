# Learnings

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
