# D1 Lifecycle Audit — status truth for rounds 8–9 trimmed specs

Review-fix for CDD defect D1 (docs-cleanup PR #679). Rounds 8–9 stripped CDD
design-cycle apparatus (`Status: Draft` + unchecked ACs among it) from ~15 canonical
specs **by document shape**, without classifying each doc's lifecycle first. For
unimplemented proposals, the stripped `Status:` was the only signal that the doc
described a *desired* system. This audit restores an honest lifecycle status header +
tense for each doc, by evidence. Process apparatus (Issue/Mode/CDD Trace/File
Changes/AC checklists) was **not** re-added — that removal was correct.

Falsification gate: for every "not implemented" verdict, ran the grep that would
disprove it against `src/go/**` (and `src/packages/**`). OCaml is fully archived —
`git ls-files '*.ml'` returns **0 files** — so any doc citing an OCaml module as the
"current" implementation is citing a removed runtime.

## Per-doc table

| Doc | Lifecycle class | Evidence (falsification grep → result / src path) | Status header restored | Tense / contradiction fixes |
|---|---|---|---|---|
| `docs/reference/runtime/MEMORY.md` | Draft (proposed) | `grep memory_episodic\|memory_reflective\|memory_working src/go src/packages` → NONE. Surfaces exist (`threads/adhoc`,`threads/reflections`,`state/conversation.json`) via `hubinit.go`/`activate.go`; the *model* + zones do not. | Added `Status: Draft (proposed — not implemented)` + Implementation note (surfaces exist, model does not) | "defines the memory model…fixes the rule" → "proposes…would fix" |
| `docs/reference/runtime/HYBRID-LLM-ROUTING.md` | Draft (proposed) | `grep routing_policy\|cnos.llm\|dependency_depth\|remote_min_tokens src/go` → NONE | Upgraded `Status: Draft` → `Draft (proposed — not implemented)` + Implementation: not implemented | "defines the hybrid LLM routing layer" → "proposes…" |
| `docs/reference/runtime/POLYGLOT-PACKAGES-AND-PROVIDERS.md` | **Partially implemented** | Commands ship: `pkg.go` content classes incl `commands`, `CommandSource`/`CommandTier`/`CommandSpec` + precedence in `internal/cli/`. Providers: `grep capability-provider\|provider registry src/go` → not hosted. | Added `Status: Partially implemented` + Implementation split (commands yes / provider contract §5–6 no) | Added clause: command surface implemented, provider surface a normative target |
| `docs/reference/packages/PACKAGE-ARTIFACTS.md` | **Implemented** | `restore.go` (lockfile→index→HTTP/local fetch→SHA-256→tar extract), `cmd_deps.go`, `pkgbuild/build.go`; `commands` class in `pkg.go`; activation table in `activation/index.go` | Added `Status: Implemented` + Implementation paths | Removed stale OCaml "current" refs: `cn_runtime_contract.ml`, "27 references in cn_deps.ml" |
| `docs/reference/protocol/cn/MESSAGE-PACKET-TRANSPORT.md` | Draft (proposed) — **confirmed D1 example** | `grep cn.packet\|payload_sha256\|inbound-index\|refs/cn/msg\|materialize_branch src/go src/packages` → NONE | Added `Status: Draft (proposed — not implemented)` + Implementation: not implemented in Go | "It replaces branch-diff message discovery" → "It is intended to replace"; `cn_io.ml materialize_branch is the current…entry point` → "was the entry point in the archived OCaml runtime; the Go runtime has no message-packet transport yet"; opening "defines" → "specifies a proposed" |
| `docs/architecture/HUB-PLACEMENT-MODELS.md` | Draft (proposed) — **confirmed D1 example** | `grep placement.json\|hub_root\|workspace_root\|cn.hub_placement src/go` → NONE | Added `Status: Draft (proposed — not implemented)` + Implementation: not implemented in Go | "This change replaces it with" → "This design would replace it with"; "cnos has two explicit roots" → "cnos gains two explicit roots"; "defines" → "specifies a proposed" |
| `docs/reference/packages/BUILD-AND-DIST.md` | **Implemented** (was mislabeled Draft) | `cmd_build.go` + `pkgbuild/build.go` (dist tarballs + `index.json` + `checksums.txt`); `cmd_deps.go` + `restore.go` (restore) | Corrected `Status: Draft` → `Status: Implemented` + Implementation paths | Noted `cnos.transport.git`/Rust `git-cn` provider in the layout example is aspirational (does not exist) |
| `docs/reference/protocol/cn/THREAD-EVENT-MODEL.md` | Draft (proposed) | `grep thread_event\|thread-events\|root_event_id\|refs/cn/feed\|refs/cn/inbox src/go` → NONE. Depends on unimplemented Message Packet Transport. | Upgraded `Status: Draft` → `Draft (proposed — not implemented)` + Implementation: not implemented | Opening "cnos has validated packet transport" (false) → "cnos assumes validated packet transport (proposed…not yet implemented)"; "Define the semantic model" → "Propose…" |
| `docs/reference/runtime/GIT-CN-PACKAGE.md` | Draft (proposed) | `ls src/packages \| grep transport` → none; `grep git-cn\|cnos.transport.git src/go` → none | Upgraded `Status: Draft` → `Draft (proposed — not implemented)` + Implementation: not implemented | Header note: no package/command/provider exists |
| `docs/reference/runtime/ORCHESTRATORS.md` | Draft (proposed) | `grep cn.orchestrator\|effect.plan\|EffectPlan\|orchestrator.v1 src/go` → NONE | Upgraded `Status: Draft` → `Draft (proposed — not implemented)` + Implementation: not implemented (IR/CTB/activation-index); commands exist separately | — |
| `docs/reference/protocol/cn/PROTOCOL.md` | **Normative target** (was **falsely** `Implemented`) | Claimed "Implemented (v2 — `cn_protocol.ml`)" but 0 `.ml` files in tree; `grep sender_transition\|receiver_transition\|materialize_branch src/go` → none (only `cmd_issues_fsm.go`, an unrelated issue-state FSM) | Replaced false `Status: Implemented` with `Status: Normative target — not implemented in the Go runtime` + Implementation note (previously implemented in archived OCaml) | Added banner: all "Maps to current code" / "Done" OCaml references are the archived runtime, not the current Go implementation |
| `docs/reference/packages/EXTENSION-REGISTRY.md` | Draft (proposed) | Registry/trust/channels/bundles — none in `src/go`; depends on unimplemented Runtime Extensions | Upgraded `Status: Draft` → `Draft (proposed — not implemented)` + Implementation: not implemented | Opening "Runtime Extensions makes capability growth coherent" → "(a proposed design, not yet implemented) would make…"; noted dependency is also unimplemented |
| `docs/reference/runtime/extensions/RUNTIME-EXTENSIONS.md` | Draft (proposed) | `grep cn.extension\|extension host\|open op registry src/go` → no extension discovery/host/dispatch | Upgraded `Status: Draft — converged` → `Draft (proposed — not implemented)` + Implementation: not implemented | Header note: "native OCaml plugin loading" (§8.2) refers to the archived OCaml runtime |
| `docs/architecture/cognitive-substrate/COGNITIVE-SUBSTRATE.md` | **Normative — enforced** | Taxonomy enforced: `pkg.go` content classes (doctrine/mindsets/skills/extensions/commands/orchestrators/katas); `activation/index.go` skill loading; `doctor.go` validates 4-layer contract (identity/cognition/body/medium) | Changed `Status: Draft` → `Status: Normative — enforced` + Implementation paths | Noted §6.1 full Runtime-Contract emission / RUNTIME-CONTRACT-v2 remains a normative target beyond current emission |

## Summary counts (14 docs)

- **Implemented:** 2 — PACKAGE-ARTIFACTS, BUILD-AND-DIST
- **Partially implemented:** 1 — POLYGLOT-PACKAGES-AND-PROVIDERS (commands ship; provider contract not)
- **Normative target / enforced:** 2 — PROTOCOL (not implemented in Go; corrected false "Implemented"), COGNITIVE-SUBSTRATE (normative, enforced by package system + doctor)
- **Draft (proposed — not implemented):** 9 — MEMORY, HYBRID-LLM-ROUTING, MESSAGE-PACKET-TRANSPORT, HUB-PLACEMENT-MODELS, THREAD-EVENT-MODEL, GIT-CN-PACKAGE, ORCHESTRATORS, EXTENSION-REGISTRY, RUNTIME-EXTENSIONS
- **Historical:** 0

All 14 received a status correction. Five had **no** status header at all before this fix
(MEMORY, POLYGLOT, PACKAGE-ARTIFACTS, MESSAGE-PACKET-TRANSPORT, HUB-PLACEMENT-MODELS);
two had a **wrong** status (PROTOCOL falsely "Implemented"; BUILD-AND-DIST mislabeled
"Draft" when it ships); the rest had a bare "Draft" upgraded to explicit
"(proposed — not implemented)" with an implementation line.

## Docs left unchanged (checked, already honest)

- Spot-check of other `docs/reference/**` and `docs/architecture/**` specs: all carry a
  Status header already (ARCHITECTURE v2.0.0, CELL-RUNTIME "Proposed…not ratified",
  INVARIANTS Active, CAR Implemented, CAA, CONFIGURE-AGENT Draft, PROVIDER-CONTRACT-v1
  Draft, RUNTIME-CONTRACT-v2 Design, PACKAGE-AUTHORING Draft, PACKAGE-SYSTEM Draft,
  GIT-AS-THE-LOWEST-DURABLE-SUBSTRATE `status: RELEASE`, THREAD-API Implementing).
- `docs/architecture/DESIGN-CONSTRAINTS.md` — no Status header but it is a governing
  constraints reference (not a feature spec claiming implementation); honest as-is.
- `docs/reference/packages/DESIGN-266-dist-out-of-git.md` — no lifecycle header but
  retains its full CDD apparatus (`Issue:/Mode:/## Problem`), so its proposal lifecycle
  is already visible; not a present-tense-current overreach. Left as found (also outside
  the rounds 8–9 trimmed set).

## Notes on conservatism

- POLYGLOT and PROTOCOL were the hard calls. POLYGLOT is genuinely mixed (commands ship,
  providers don't) → labeled "Partially implemented" with the split spelled out rather
  than forcing one class. PROTOCOL documents a contract implemented only in the now-removed
  OCaml runtime → "Normative target, not implemented in Go" is the honest reading; the
  prior "Implemented" was the sharpest instance of the D1 defect.
- Where the choice was Draft vs Normative-target, chose the more conservative **Draft**
  (MESSAGE-PACKET-TRANSPORT, HUB-PLACEMENT, THREAD-EVENT, etc. all carry `## Proposal`/
  proposal framing and adoption is not evidenced).
