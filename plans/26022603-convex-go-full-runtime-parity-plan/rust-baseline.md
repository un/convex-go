# Rust Runtime Baseline

- Repository: `/Users/omar/code/convex-rs`
- Baseline commit: `95ec2b33f568f99b053eaa101ddce6c94cb71f57`
- Captured on: `2026-02-26`

## Primary parity reference modules

- `sync_types/src/types/mod.rs`
  - canonical protocol type system (typed IDs, variants, state/request semantics)
- `sync_types/src/types/json.rs`
  - strict wire codec encode/decode behavior and compatibility defaults
- `sync_types/src/timestamp.rs`
  - timestamp model and base64 wire encoding semantics
- `src/sync/mod.rs`
  - transport protocol interfaces and reconnect request model
- `src/sync/web_socket_manager.rs`
  - websocket lifecycle (open/send/read/heartbeat/inactivity/reconnect/backoff)

## Baseline usage rules for this run

- Every parity decision in Go runtime changes must cite one of the files above in step evidence.
- If Rust baseline changes during this run, update this file and record the diff in the corresponding step evidence.
- Fixture import work (later steps) must track this exact commit hash so drift is attributable.
