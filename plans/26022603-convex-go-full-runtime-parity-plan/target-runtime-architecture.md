# Target Runtime Architecture

This document defines the end-state layering for full runtime parity.

## Layering

1. **Protocol model + codec (`internal/protocol`)**
   - Owns typed protocol IDs and message variants.
   - Owns strict JSON wire encode/decode with deterministic malformed-input failures.
   - No runtime network concerns.

2. **Transport manager (`internal/sync`)**
   - Owns websocket lifecycle: dial/handshake, write queue, read loop, heartbeat, inactivity watchdog, reconnect backoff.
   - Emits typed protocol responses and connection-state events.
   - No query-state business logic.

3. **Base sync state (`internal/baseclient`)**
   - Owns query registry, subscriber mapping, auth identity version, observed timestamp, request completion semantics.
   - Applies transitions and tracks replay queue rebuild inputs.
   - No direct websocket I/O.

4. **Worker runtime (`internal/runtime` planned)**
   - Single owner of runtime transitions and ordering guarantees.
   - Multiplexes API commands, transport responses, and timers.
   - Enforces flush-before-select and communicate-during-send behavior.

5. **Public API facade (`convex`)**
   - Thin API methods (`Subscribe`, `Query`, `Mutation`, `Action`, auth/lifecycle) that route through worker commands.
   - Maintains user-facing channel/results ergonomics.
   - Avoids local bypass logic for sync semantics.

## Concurrency ownership

- Transport internals own websocket conn object and frame loops.
- Worker owns ordered state transitions and replay emission.
- Base state is mutated only on worker thread/goroutine.
- Public API methods enqueue commands; they do not mutate transport/base state directly.

## Control flow

### Outbound

API call -> Worker command -> Base state update (if needed) -> Protocol message encode -> Transport write queue -> Websocket.

### Inbound

Websocket frame -> Transport decode/classification -> Worker event -> Base state apply -> Request completion / subscription fan-out.

### Failure/reconnect

Transport protocol failure -> Worker reconnect command with reason + max observed timestamp -> Backoff reconnect -> Connect -> replay auth/query set/in-flight requests.

## Invariants

- Protocol encode/decode behavior is deterministic and strict.
- Worker is the only place that decides replay order.
- Mutation completion waits for visibility timestamp unless operation errors.
- Transition version mismatches trigger reconnect.
- Transition chunks must assemble before decode/apply.
