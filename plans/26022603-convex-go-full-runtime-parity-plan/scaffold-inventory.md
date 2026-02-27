# Runtime Scaffold Inventory

Status key: `open` | `closed`

| Inventory ID | File | Function/Area | Scaffold symptom | Target replacement | Status | Closed in step | Evidence |
| --- | --- | --- | --- | --- | --- | --- | --- |
| SI-001 | `internal/protocol/types.go` | `ClientMessage` / `ServerMessage` | Single broad structs with optional fields | Strict variant unions and typed payload structs | closed | 18 | `execution/26022603-full-runtime/step-18.md` |
| SI-002 | `internal/protocol/types.go` | `QuerySetModification` | String-tag variant with optional fields | Explicit Add/Remove variants | closed | 8 | `execution/26022603-full-runtime/step-8.md` |
| SI-003 | `internal/protocol/types.go` | `StateModification` | String-tag state updates with weak validation | Explicit QueryUpdated/Failed/Removed variants | closed | 9 | `execution/26022603-full-runtime/step-9.md` |
| SI-004 | `internal/protocol/codec.go` | `EncodeClientMessage`/`DecodeClientMessage` | Direct JSON marshal/unmarshal passthrough | Strict encode/decode with deterministic errors | closed | 16 | `execution/26022603-full-runtime/step-16.md` |
| SI-005 | `internal/protocol/codec.go` | `EncodeServerMessage`/`DecodeServerMessage` | Direct JSON marshal/unmarshal passthrough | Strict encode/decode with null-vs-missing rules | closed | 18 | `execution/26022603-full-runtime/step-18.md` |
| SI-006 | `internal/sync/websocket_manager.go` | open/send/reconnect loops | No write queue and limited reconnect semantics | Worker-owned transport lifecycle with ordered queue | closed | 31 | `execution/26022603-full-runtime/step-31.md` |
| SI-007 | `internal/sync/websocket_manager.go` | `readLoop` | Limited close/decode classification | Protocol failure classes and reconnect propagation | closed | 42 | `execution/26022603-full-runtime/step-42.md` |
| SI-008 | `internal/sync/websocket_manager.go` | heartbeat/inactivity | No ping ticker or inactivity watchdog | Deterministic heartbeat and inactivity reconnect | closed | 28 | `execution/26022603-full-runtime/step-28.md` |
| SI-009 | `convex/client.go` | `handleServerMessage` | Transition chunks rejected unconditionally | Transition chunk assembly and decode pipeline | closed | 35 | `execution/26022603-full-runtime/step-35.md` |
| SI-010 | `convex/client.go` | orchestration | Monolithic lock-based runtime orchestration | Dedicated worker command/event loop | closed | 41 | `execution/26022603-full-runtime/step-41.md` |
| SI-011 | `convex/client.go` | `Subscribe`/`Query` | Local orchestration path bypassing worker model | Worker-managed subscribe/query semantics | closed | 43 | `execution/26022603-full-runtime/step-43.md` |
| SI-012 | `convex/client.go` | `Mutation`/`Action` | Inline pending request handling | Request manager + worker completion semantics | closed | 44 | `execution/26022603-full-runtime/step-44.md` |
| SI-013 | `convex/client.go` | auth APIs | Auth refresh and replay logic not centralized | Worker-routed auth lifecycle and replay rebuild | closed | 46 | `execution/26022603-full-runtime/step-46.md` |
| SI-014 | `integration/integration_test.go` | live gates | Narrow integration coverage | Deployment-backed reconnect/auth-refresh parity suite | closed | 47 | `execution/26022603-full-runtime/step-47.md` |

## Closure policy

- A row can be marked `closed` only after implementation is merged and validated.
- Every closure must include the step number and evidence file path.
- Step 48 must leave no `open` rows.
