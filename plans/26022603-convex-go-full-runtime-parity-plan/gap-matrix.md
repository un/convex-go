# Gap Matrix: Go vs Rust Runtime Baseline

Baseline: `95ec2b33f568f99b053eaa101ddce6c94cb71f57`

## Protocol model and codec

| Area | Go status | Rust reference | Gap | Target behavior |
| --- | --- | --- | --- | --- |
| IDs and versions | Primitive aliases in `internal/protocol/types.go` | `sync_types/src/types/mod.rs` | Weak typing and broad optional fields | Explicit typed wrappers and variant structs |
| Query set modifications | `QuerySetModification` has string `Type` | `QuerySetModification::{Add,Remove}` | Runtime-only validation, no compile-time variant safety | Strong Add/Remove variants with strict decode rules |
| State modifications | `StateModification` with optional fields | `StateModification::{QueryUpdated,QueryFailed,QueryRemoved}` | Missing strict per-variant required fields | Typed union with deterministic malformed rejection |
| Client message codec | `json.Marshal` passthrough | `sync_types/src/types/json.rs` | No compatibility defaults / missing-field checks | Strict encoder/decoder with stable errors |
| Server message codec | `json.Unmarshal` passthrough | `sync_types/src/types/json.rs` | Malformed unions accepted too loosely | Strict variant decoding and null-vs-missing semantics |

## Transport lifecycle

| Area | Go status | Rust reference | Gap | Target behavior |
| --- | --- | --- | --- | --- |
| Open + handshake | `openConn` direct dial/write | `WebSocketInternal::new` | Limited handshake error detail | Structured handshake failure propagation |
| Send path | Direct `WriteMessage` in `Send` | Worker request queue | No queued writer/ack semantics | Ordered write queue with backpressure contract |
| Read path | Basic read loop, no frame classes | `work` server frame switch | Close/decode classes not separated | Classified failures + protocol failure propagation |
| Heartbeat/inactivity | Missing ping ticker/watchdog | `HEARTBEAT_INTERVAL` / inactivity threshold | Reconnect not triggered on inactivity | Deterministic ticker/watchdog reconnect path |
| Backoff | `baseclient.Backoff` exists but not transport-owned | `convex_sync_types::backoff::Backoff` | Reconnect policy not transport-owned | Transport-integrated jittered backoff and reset |

## Base sync state and request semantics

| Area | Go status | Rust reference intent | Gap | Target behavior |
| --- | --- | --- | --- | --- |
| Request completion | Mixed logic in `convex/client.go` | Request manager semantics | Mutation visibility and action completion split across layers | Dedicated request manager with deterministic replay and completion |
| Transition apply | Inline handling in `handleTransition` | Versioned transition machine | No start/end version mismatch reconnect | Strict start/end checks + reconnect trigger |
| Chunk assembly | Explicitly unsupported (`transition chunk unsupported`) | `TransitionChunk` server variant | Missing production chunk path | Buffered assembly with ordering checks |
| Auth refresh | Callback used but state flow mixed | force-refresh on reconnect intent | Identity version/replay semantics not centralized | Auth refresh coordinated with replay rebuild |

## Runtime orchestration and API surface

| Area | Go status | Rust reference intent | Gap | Target behavior |
| --- | --- | --- | --- | --- |
| Worker architecture | Monolithic `Client` lock orchestration | Dedicated worker loop model | API calls and transport handling tightly coupled | Worker command/event loop owns ordering and state transitions |
| Flush-before-select | Not guaranteed | Rust `select_biased!` behavior | Replay or control messages can be delayed | Flush critical outbound queue before normal select |
| Communicate-during-send | Blocking send path can starve reads | Rust keeps reading while awaiting send completion | Potential starvation/deadlock windows | Continue processing inbound protocol while send in flight |
| Public APIs | Query/Mutation/Action wired directly to transport | Worker-managed runtime | Local shortcut behavior remains in API layer | API methods route through worker only |

## Validation and parity gates

| Area | Go status | Rust reference intent | Gap | Target behavior |
| --- | --- | --- | --- | --- |
| Malformed corpus | Limited decode coverage | strict decode test corpus | Missing broad malformed vectors | Add corpus + stable error assertions |
| Fixture conformance | No imported Rust vectors | shared fixture parity | Drift can go undetected | Fixture import + CI conformance gate |
| Fuzzing | No codec fuzz gate | robust decode coverage | Panic/regression risk | Fuzz targets seeded with corpus + fixtures |
| Live runtime parity | Integration tests are narrow | deployment-backed parity checks | reconnect/auth-refresh parity unproven | Expanded live suite + final parity gate |

## Execution notes

- The matrix is intentionally function-level actionable and maps directly to steps 3-48.
- Step evidence files must cite the row(s) they closed.
