# Learnings

## 2026-02-26 Step 47 expand-live-integration-suite
- Context: Live coverage existed but bundled multiple flows into one linear test, with limited control for reconnect/auth-refresh validation.
- Learning: Scenario-based live subtests with explicit env-gated probes make reconnect/auth-refresh checks configurable without weakening default CI behavior.
- Impact on next steps: Final scaffold-zero and parity gates can reference a structured live suite with optional strict reconnect/auth expectations.

## 2026-02-26 Step 46 auth-api-through-worker
- Context: Auth APIs still mixed direct state mutation with out-of-band authenticate sends, which could diverge from worker ordering guarantees.
- Learning: Worker-dispatched auth update commands (including callback registration) keep token state, identity versioning, and authenticate replay on one ordered path while preserving pre-connect fallback semantics.
- Impact on next steps: Live integration expansion can now validate auth-refresh and reconnect behavior against the same worker runtime used by all public APIs.

## 2026-02-26 Step 45 watchall-clone-lifecycle
- Context: `WatchAll` lifecycle was still partly direct-path and clone/close semantics needed explicit regression coverage under worker ownership.
- Learning: Worker-routed watch registration/unregistration plus coherent multi-subscriber snapshot assertions catches lifecycle drift while preserving shared-runtime clone close idempotence.
- Impact on next steps: Auth API routing can reuse the same worker command lifecycle pattern with confidence around close behavior.

## 2026-02-26 Step 44 mutation-action-through-worker
- Context: Mutation/action request lifecycle was still initiated directly from API methods, outside worker command ownership.
- Learning: Creating pending requests via worker commands and returning typed request handles keeps enqueue, cancellation, and response completion on one deterministic control path.
- Impact on next steps: Watch/auth lifecycle wiring can now share the same command dispatch pattern used by query and request APIs.

## 2026-02-26 Step 43 subscribe-query-through-worker
- Context: Subscribe/query state mutations were still performed directly from API call paths, bypassing worker ownership.
- Learning: Routing subscribe/unsubscribe through typed worker commands keeps query-set state and outbound modify-query messages on the worker control path while preserving subscribe-first-value query semantics.
- Impact on next steps: Mutation/action/auth APIs can now be moved to the same worker command channel with shared ordering guarantees.

## 2026-02-26 Step 42 protocol-failure-reconnect-propagation
- Context: Failure-triggered reconnects existed, but payload completeness across failure classes and max-observed derivation were not explicitly enforced.
- Learning: Centralizing failure payload construction (`reason` + max of observed transition TS and mutation visibility TS) plus class-level propagation tests hardens reconnect correctness under protocol/transport errors.
- Impact on next steps: Worker-routed public APIs can rely on reconnect intent being complete and test-verified regardless of failure origin.

## 2026-02-26 Step 41 communicate-during-send
- Context: Outbound flush sends could stall worker progress if a send path blocked, risking starvation of inbound protocol processing.
- Learning: Wrapping sends in `sendWhileCommunicating` and selecting over transport/command channels while awaiting send completion removes starvation without violating flush ordering.
- Impact on next steps: Protocol failure propagation can now rely on inbound error visibility even when outbound traffic is under pressure.

## 2026-02-26 Step 40 flush-before-select
- Context: Reconnect replay messages and command/event handling shared the worker loop but had no explicit pre-select flush invariant.
- Learning: A dedicated outbound queue + wake signal with `flushOutboundBeforeSelect` guarantees replay/auth-critical payloads are emitted before normal event/command selection.
- Impact on next steps: Communicate-during-send can now be layered onto a single outbound send path without losing replay priority semantics.

## 2026-02-26 Step 39 worker-main-loop
- Context: Transport reads were handled by a dedicated listener goroutine, while worker command routing was not yet connected.
- Learning: Replacing listener reads with a worker event loop (command + transport select) establishes a single control point for runtime transitions and clean shutdown.
- Impact on next steps: Outbound flush-before-select and communicate-during-send behavior can now be implemented in one loop without cross-goroutine ordering drift.

## 2026-02-26 Step 38 worker-command-event-model
- Context: Runtime operations were coordinated by ad-hoc goroutines and direct method calls, which made ownership boundaries implicit.
- Learning: Explicit typed worker commands/events (API command, transport message/error/done) make cancellation and routing semantics unambiguous before loop refactors.
- Impact on next steps: Worker main loop implementation can consume one normalized event shape instead of branching on mixed channel payload conventions.

## 2026-02-26 Step 37 replay-queue-rebuild
- Context: Reconnect replay used map iteration, so in-flight request replay order was non-deterministic and harder to validate.
- Learning: Rebuilding replay payloads from canonical sorted IDs (auth -> sorted query set -> sorted pending requests) gives deterministic reconnect behavior and blocks stale ordering regressions.
- Impact on next steps: Worker command/event modeling can assume stable replay ordering invariants during reconnect recovery.

## 2026-02-26 Step 36 auth-callback-refresh-semantics
- Context: Reconnect auth refresh calls existed, but callback failure handling and refresh-side auth sends could produce ambiguous behavior.
- Learning: Treating force-refresh callback failures as reconnect retry conditions and applying refreshed tokens without pre-reconnect sends makes auth replay deterministic.
- Impact on next steps: Replay-queue rebuild can now rely on stable, single-path authenticate replay after reconnect.

## 2026-02-26 Step 35 transition-chunk-assembly
- Context: `TransitionChunk` frames were treated as unsupported, which forced reconnect instead of applying valid chunked transitions.
- Learning: Enforcing strict per-transition chunk ordering (`partNumber` sequence), deterministic assembly, and decode-as-Transition checks gives safe chunk support without silent state corruption.
- Impact on next steps: Auth refresh/replay rebuild can assume both direct transitions and chunked transitions converge through the same validated apply path.

## 2026-02-26 Step 34 strict-transition-apply
- Context: Transition application accepted any start version, risking silent state divergence.
- Learning: Enforcing start/end version requirements and mismatch-triggered reconnect is essential to keep client/server state synchronized.
- Impact on next steps: Transition-chunk assembly can now feed into a strict version-checked apply pipeline.

## 2026-02-26 Step 33 localsyncstate-semantics
- Context: Local sync semantics needed stronger guarantees around dedupe/versioning/subscriber lifecycle and observed timestamp monotonicity.
- Learning: Exhaustive state-semantic tests are the fastest way to lock correctness before deeper transition/worker refactors.
- Impact on next steps: Transition apply/chunk assembly can now build on verified local state invariants.

## 2026-02-26 Step 32 requestmanager-semantics
- Context: Mutation replay/completion semantics were under-specified and allowed unresolved mutations to complete prematurely.
- Learning: Explicit `WaitingOnTS` tracking is required so mutation completion only occurs after response visibility timestamps are established.
- Impact on next steps: Local sync state and transition apply logic can now rely on deterministic request completion invariants.

## 2026-02-26 Step 31 transport-integration-tests
- Context: Individual transport unit tests existed, but lifecycle behavior (open/send/read/reconnect) needed integrated validation under race detector.
- Learning: A controllable websocket integration test plus race-gate execution exposed and fixed a real channel-close/send race in transport shutdown.
- Impact on next steps: Runtime orchestration work can build on a race-clean transport baseline.

## 2026-02-26 Step 30 websocket-state-callback-ordering
- Context: State callback emission was ad-hoc and could duplicate/reorder transitions across connect/reconnect cycles.
- Learning: Centralizing state emission with deduping (`emitState`) and callback-order tests locks in deterministic lifecycle ordering.
- Impact on next steps: Transport integration tests can assert full lifecycle sequences reliably.

## 2026-02-26 Step 29 reconnect-backoff-jitter-parity
- Context: Backoff jitter behavior diverged from Rust intent and used deterministic pseudo-jitter by default.
- Learning: Rust-style full-jitter exponential backoff with cap + deterministic RNG tests gives predictable parity and safer reconnect spread.
- Impact on next steps: Transport state callbacks can now rely on stable backoff semantics during reconnect cycles.

## 2026-02-26 Step 28 inactivity-watchdog
- Context: Heartbeats alone did not enforce reconnect when server-side silence persisted.
- Learning: Tracking last server response timestamp and failing with a deterministic `InactiveServer` signal provides a clean reconnect trigger.
- Impact on next steps: Reconnect/backoff ordering can now treat inactivity as a first-class transport failure reason.

## 2026-02-26 Step 27 heartbeat-ping-ticker
- Context: Transport had no periodic liveness signal, leaving long-lived idle connections unmanaged.
- Learning: A dedicated heartbeat loop with websocket ping control frames and deterministic failure surfacing closes that liveness gap cleanly.
- Impact on next steps: Inactivity watchdog can now use heartbeat cadence as the timing backbone for reconnect triggers.

## 2026-02-26 Step 26 websocket-read-loop-classification
- Context: Reader loop surfaced raw websocket/decode errors without consistent failure classes.
- Learning: Explicit classification for close frames, read failures, frame-type mismatches, and protocol decode failures improves reconnect reason quality and testability.
- Impact on next steps: Heartbeat/watchdog and reconnect logic can branch on clearer transport failure signals.

## 2026-02-26 Step 25 websocket-write-queue-loop
- Context: Direct websocket writes coupled API send latency to socket I/O and made ordering guarantees implicit.
- Learning: A dedicated write queue + writer loop preserves send order and defines clear backpressure behavior (`queue full` error) without blocking on network writes.
- Impact on next steps: Reader/heartbeat/watchdog logic can now run independently from outbound write pressure.

## 2026-02-26 Step 24 websocket-connect-message-path
- Context: Connect payload correctness across open/reconnect cycles was implicit and unverified.
- Learning: Transport-level tests that capture first-frame connect messages lock in session metadata, reason propagation, and observed-timestamp encoding semantics.
- Impact on next steps: Writer/read-loop work can assume connect handshake metadata is stable per connection attempt.

## 2026-02-26 Step 23 websocket-dial-handshake
- Context: Transport dial failures lacked actionable handshake context and weakly surfaced context cancellation.
- Learning: Capturing HTTP handshake status/body and explicit cancellation wrapping makes connection setup failures diagnosable and deterministic.
- Impact on next steps: Connect/send/read loops can now rely on richer failure classification for reconnect logic.

## 2026-02-26 Step 22 codec-fuzzing-gate
- Context: Strict unit/corpus checks reduce regressions, but fuzzing is needed to guard against panic-class edge cases.
- Learning: Seeding fuzz targets from malformed corpus + Rust fixtures gives high-signal coverage without sacrificing determinism in normal test runs.
- Impact on next steps: Transport/runtime refactors should keep protocol entry points panic-free under arbitrary inputs.

## 2026-02-26 Step 21 fixture-conformance-gate
- Context: Imported fixtures are useful only if they actively gate protocol behavior in normal CI runs.
- Learning: Embedding Rust fixture conformance checks in `go test` keeps drift detection continuous and avoids a separate, easy-to-skip gate.
- Impact on next steps: Fuzzing can now seed from the same fixtures to broaden protocol robustness coverage.

## 2026-02-26 Step 20 import-rust-fixture-vectors
- Context: Protocol conformance needed baseline-linked fixtures instead of hand-authored local-only examples.
- Learning: A deterministic import script tied to the Rust source file and commit hash makes fixture updates auditable and repeatable.
- Impact on next steps: Conformance gate can now fail on drift between Go codec behavior and imported Rust vectors.

## 2026-02-26 Step 19 malformed-protocol-corpus-tests
- Context: Strict codecs need regression-resistant malformed-input coverage beyond ad-hoc unit tests.
- Learning: A shared malformed corpus file enables stable, scalable decode-failure assertions for both client and server protocol paths.
- Impact on next steps: Fixture import and fuzzing can reuse this corpus as deterministic seeds.

## 2026-02-26 Step 18 strict-server-codec-decode
- Context: Server decode lacked a stable envelope and did not validate some variant-specific optional-field edge cases.
- Learning: Adding a typed envelope and variant-aware null/field checks yields deterministic malformed-union handling without breaking optional-null compatibility.
- Impact on next steps: Malformed protocol corpus and fuzz gates can now assert stable decode failure classes.

## 2026-02-26 Step 17 strict-server-codec-encode
- Context: Server encode path still relied on struct marshaling behavior without a dedicated codec-level validation gate.
- Learning: An explicit server encode validator aligns encode-time failures with decode-time strictness and catches invalid response unions before wire emission.
- Impact on next steps: Add a symmetric strict decode envelope for server messages to stabilize malformed error reporting.

## 2026-02-26 Step 16 strict-client-codec-decode
- Context: Decode errors were not consistently classified, and legacy authenticate payloads needed compatibility handling.
- Learning: Wrapping decode with an envelope check (`type` first) gives stable malformed-input errors, while targeted legacy fallback preserves backward compatibility.
- Impact on next steps: Server encode/decode paths should adopt the same stable error-envelope approach.

## 2026-02-26 Step 15 strict-client-codec-encode
- Context: Client encode still depended mostly on implicit struct marshaling behavior.
- Learning: Adding an explicit encode validator in the codec path gives stable early failures and keeps wire-shape guarantees centralized.
- Impact on next steps: Apply the same explicit gate pattern to decode, including compatibility and malformed-input classification.

## 2026-02-26 Step 14 servermessage-remaining-variants
- Context: Response/error server variants still allowed malformed unions and unknown message types to slip through decode.
- Learning: Strict per-variant decode checks (including success/error union requirements) make protocol failures deterministic and actionable.
- Impact on next steps: Codec-level strictness now supports implementing explicit encode/decode gates and malformed corpus testing.

## 2026-02-26 Step 13 servermessage-transition-and-chunk-variants
- Context: Transition and chunk payloads were decoded permissively, which could hide malformed stream state.
- Learning: Strict required-field validation on transition/chunk variants catches framing errors early and gives deterministic failure paths.
- Impact on next steps: Apply strict variant validation to the remaining server message unions (responses/errors/ping).

## 2026-02-26 Step 12 clientmessage-remaining-variants
- Context: Non-connect client messages were still loosely marshaled, allowing invalid shapes through encode/decode paths.
- Learning: Full variant-specific validation in `ClientMessage` custom JSON logic catches malformed requests at the protocol boundary.
- Impact on next steps: Server-message variants should be tightened with the same explicit required-field checks.

## 2026-02-26 Step 11 clientmessage-connect-variant
- Context: Connect payload decoding accepted missing required fields and inconsistent defaults.
- Learning: Variant-specific custom JSON logic allows strict required-field checks while still applying protocol-compatible defaults (`lastCloseReason = unknown`).
- Impact on next steps: Remaining client message variants can be made strict by extending the same custom decode path.

## 2026-02-26 Step 10 typed-authentication-token-variants
- Context: Auth token payloads were loose field bags, making Admin/User/None handling ambiguous and brittle.
- Learning: A strict token variant model with compatibility aliases (`impersonating` -> `actingAs`) preserves backward compatibility without sacrificing validation.
- Impact on next steps: Client and server message unions can now embed token variants directly instead of flattening token fields.

## 2026-02-26 Step 9 typed-state-modifications
- Context: Transition apply logic depended on weakly typed state-modification payloads, which accepted malformed field combinations.
- Learning: Explicit QueryUpdated/QueryFailed/QueryRemoved variants make error handling and transition apply paths simpler and safer.
- Impact on next steps: Remaining protocol unions (auth/client/server messages) should follow the same strict variant strategy.

## 2026-02-26 Step 8 typed-query-and-queryset-modifications
- Context: Query-set updates used a single permissive struct, which allowed malformed Add/Remove payload combinations.
- Learning: Modeling Add/Remove as explicit variants with custom JSON (de)serialization gives strict required-field enforcement while preserving wire compatibility.
- Impact on next steps: Apply the same variant-modeling pattern to state modifications and auth/message unions.

## 2026-02-26 Step 7 state-version-and-timestamp-encoding
- Context: `StateVersion.ts` was treated as an unvalidated string, so malformed wire timestamps could leak into runtime state.
- Learning: Moving timestamp validation into `StateVersion.UnmarshalJSON` centralizes correctness and removes per-call decode branching.
- Impact on next steps: Protocol decode paths can trust typed `StateVersion` objects and focus on variant semantics.

## 2026-02-26 Step 6 strict-protocol-identifiers
- Context: Runtime code used ad-hoc numeric casts (`uint64`<->`uint32`) at many protocol boundaries.
- Learning: Centralized conversion helpers with overflow checks remove silent truncation risk and make identifier misuse testable.
- Impact on next steps: All upcoming codec/model changes should rely on protocol conversion helpers instead of direct casts.

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
