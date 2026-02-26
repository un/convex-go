# Learnings

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
