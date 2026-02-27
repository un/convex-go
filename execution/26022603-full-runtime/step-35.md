# Step 35 - Implement TransitionChunk assembly

## What changed
- Implemented `TransitionChunk` runtime handling in `convex/client.go`.
  - Added per-transition chunk buffers keyed by `transitionId`.
  - Enforced strict ordering (`partNumber` must match next expected part).
  - Enforced `totalParts` consistency across all chunks for a transition.
  - Assembled chunks only when complete, decoded assembled payload through `protocol.DecodeServerMessage`, and required decoded type `Transition`.
  - Routed assembled transition through existing strict `handleTransition` path.
  - Cleared chunk buffers on protocol failure to avoid stale partial state across reconnects.
- Added integration coverage in `convex/client_test.go`:
  - `TestTransitionChunkAssemblyAppliesTransition`
  - `TestTransitionChunkOutOfOrderTriggersReconnect`
- Marked step 35 as completed in `plans/26022603-convex-go-full-runtime-parity-plan/steps.json`.
- Prepended step learning in `plans/26022603-convex-go-full-runtime-parity-plan/learnings.md`.

## Validations run
- `go test ./...`

## Outcomes
- PASS: `go test ./...`
