# Issue Consolidation Analysis — 2026-04-25

> Author: α (`alpha@cdd.cnos`) under task `claude/consolidate-issue-bundles-mOkM7`.
> Inputs: open-issue snapshot (79 issues, GitHub MCP `list_issues` at 13:50:25+03:00); 3.57.0 + 3.58.0 PRAs (encoding lag tables, MCI freeze status); `.cdd/unreleased/{α,β,γ}/268.md` (just-merged cycle close-outs).
> Purpose: propose triad-executable bundles, surface stale/superseded issues, and rank execution order under the active MCI freeze.
> Method: each open issue is mapped to a subsystem; cross-issue dependency claims (`Depends on`, `Related`, `Parent`, `Subsumes`) are folded into a text dependency graph; bundles consolidate items that share files / contracts / cycle scope.

## State of play (1-paragraph summary)

The MCI freeze has continued for four consecutive releases (3.55.0 → 3.58.0). Latest PRA (3.58.0) shows **11 growing-lag items**, well above the 3-issue freeze threshold. #268 (CDD package convergence) just merged at `bfddcf22` with close-outs filed under `.cdd/unreleased/`; γ flagged **#272** (close-out filename convention) as the next MCA candidate. Of the 79 open issues, ~22 carry pre-Go-port (OCaml-era) framings or were superseded by shipped work and should be closed without execution. The remaining ~57 cluster into ~14 coherent bundles, only ~6 of which fit cleanly inside one CDD triad cycle as currently scoped.

---

## 1. Issue inventory

Subsystem keys: **CDD-method**, **package-runtime**, **kata**, **CI/proof**, **identity**, **memory/cog**, **threads/issues**, **comms (telegram)**, **diagnostics**, **hub-placement**, **vision**, **bug-OCaml**, **process-tiny**, **stale**.

Lag column reads the 3.58.0 PRA encoding-lag table where present; "—" = not in the active table; "shipped" = closed by the named cycle but the issue itself is still open.

| #   | Title (short)                                                          | Labels         | Subsystem        | Lag      |
| --- | ---------------------------------------------------------------------- | -------------- | ---------------- | -------- |
| 4   | Daily threads should persist to hub (not workspace)                    | bug            | bug-OCaml        | —        |
| 5   | Require review commit before merging to main                           | enhancement    | process-tiny     | —        |
| 6   | Add `cn thread new` command                                            | enhancement    | threads/issues   | —        |
| 7   | Add processed thread archive                                           | enhancement    | threads/issues   | —        |
| 8   | Daily log rotation for `cn.log`                                        | enhancement    | bug-OCaml        | —        |
| 9   | Design thread-based task management                                    | enhancement    | threads/issues   | —        |
| 10  | Update ROADMAP.md to reflect current state                             | bug            | stale            | —        |
| 11  | Replace blanket exception catches (OCaml)                              | bug            | bug-OCaml        | —        |
| 12  | Replace mutable ref with fold in inbox.ml (OCaml)                      | bug            | bug-OCaml        | —        |
| 14  | Ensure only non-secret info gets committed                             | —              | process-tiny     | —        |
| 15  | Create forget skill                                                    | —              | memory/cog       | —        |
| 16  | `cn doctor`: warn if peers exist but no recent sync                    | —              | diagnostics      | —        |
| 17  | Installer: configure heartbeat for `cn in` wake                        | —              | bug-OCaml        | —        |
| 18  | Create troubleshooting skill                                           | —              | memory/cog       | —        |
| 20  | Eliminate silent failures in `cn agent` daemon                         | bug            | bug-OCaml        | —        |
| 30  | Inject timestamp into agent context (temporal awareness)               | enhancement    | diagnostics      | —        |
| 34  | Runtime coherence loop: CMP/MCP/CAP/CLP explicit in runtime            | enhancement    | memory/cog       | —        |
| 35  | Daemon mode periodic reflection cycles                                 | enhancement    | memory/cog       | —        |
| 39  | Pi→Sigma peer messaging clone-not-found                                | bug            | bug-OCaml        | —        |
| 42  | Telegram media parity with OpenClaw                                    | enhancement    | comms (telegram) | —        |
| 43  | Interrupt mechanism (stale-queue races user)                           | enhancement    | comms (telegram) | —        |
| 44  | Streaming output to Telegram                                           | —              | comms (telegram) | —        |
| 45  | Native cnos issue tracking (move off GH Issues)                        | —              | threads/issues   | —        |
| 52  | Cognitive Substrate: Role determines skill loadout                     | —              | memory/cog       | —        |
| 59  | `cn doctor` deep self-model consistency validation                     | —              | diagnostics      | —        |
| 68  | Agent self-diagnostics — full runtime introspection                    | —              | diagnostics      | —        |
| 69  | First-wake setup experience (Her-style)                                | —              | vision           | —        |
| 70  | Personality archetypes                                                 | —              | vision           | —        |
| 71  | Voice messages — TTS + Telegram sendVoice                              | —              | vision           | —        |
| 77  | P2P agent network with compute marketplace                             | enhancement    | vision           | —        |
| 79  | Projection surfaces (Custom GPT, X, web)                               | enhancement    | vision           | —        |
| 84  | CA mindset missing reflection requirements for continuity              | mindset        | memory/cog       | —        |
| 87  | Claude Code scheduled agent to run `cn` steps via cron                 | —              | hub-placement    | —        |
| 94  | `cn cdd`: mechanize CDD invariants without replacing judgment          | —              | CDD-method       | —        |
| 96  | Align docs taxonomy / bundle contracts / navigation                    | docs           | stale            | —        |
| 100 | Memory as first-class retention faculty                                | —              | memory/cog       | —        |
| 101 | Normalize the skill corpus (skill / runbook / reference)               | —              | memory/cog       | —        |
| 124 | Agent asks permission for in-scope work despite autonomy defaults      | bug,coherent…  | memory/cog       | —        |
| 132 | Rename skill categories `agent/`→`self/`, `ops/`→`hub/`                | —              | memory/cog       | —        |
| 135 | #74 Phase 2: per-pass logging in N-pass bind loop                      | —              | bug-OCaml        | —        |
| 138 | Remove cron installation management from `cn`                          | bug            | bug-OCaml        | —        |
| 141 | Peer sync: empty-content branches cause permanent limbo                | bug            | bug-OCaml        | —        |
| 142 | Peer sync FidoNet parity (receipts/relay/dead-letter)                  | —              | bug-OCaml        | —        |
| 148 | `check_binary_version_drift` parses wrong binary output                | bug            | bug-OCaml        | —        |
| 149 | UIE must include skill loading as part of Understand                   | enhancement    | memory/cog       | —        |
| 153 | Thread event model: canonical semantic layer                           | —              | threads/issues   | —        |
| 154 | Thread event model Phase 1                                             | —              | threads/issues   | —        |
| 156 | Attached hub placement: split hub_root and workspace_root              | —              | hub-placement    | —        |
| 162 | Modular CLI commands                                                   | —              | package-runtime  | —        |
| 168 | L8 candidate: method-shaping leverage — track evidence                 | —              | CDD-method       | —        |
| 170 | Orchestrators: mechanical workflow IR                                  | —              | package-runtime  | —        |
| 175 | #170.4: CTB → orchestrator IR compiler                                 | —              | package-runtime  | —        |
| 180 | Package system doc incoherence (β git-native vs α artifact-first)      | —              | stale            | —        |
| 181 | Package source: gh-pages index as consumer-facing endpoint             | —              | package-runtime  | —        |
| 182 | Core refactor: package-driven runtime substrate (OCaml)                | —              | stale            | —        |
| 186 | Package restructuring: lean cnos.core + domain-aligned packages        | —              | package-runtime  | —        |
| 189 | Native change proposals & review mechanics (PR-equivalents as threads) | —              | threads/issues   | —        |
| 190 | Agent web surface (browsable hub state)                                | —              | threads/issues   | —        |
| 192 | Runtime kernel rewrite OCaml → Go                                      | —              | package-runtime  | low      |
| 193 | Orchestrator runtime: `llm` step execution + bindings                  | —              | package-runtime  | growing  |
| 199 | Stacked v3.39.0 + v3.40.0 PRAs                                         | —              | stale            | growing  |
| 216 | Migrate non-bootstrap kernel commands to packages                      | —              | package-runtime  | growing  |
| 218 | `cnos.transport.git` provider (Rust)                                   | —              | package-runtime  | growing  |
| 230 | `cn deps restore` version-upgrade silent skip                          | P1             | package-runtime  | growing  |
| 235 | `cn build --check` validates entrypoints + skill paths                 | P1             | package-runtime  | growing  |
| 238 | Smoke: release bootstrap compatibility check                           | enhancement    | CI/proof         | growing  |
| 240 | CDD triadic protocol — git transport, `.cdd/` artifacts                | enhancement    | CDD-method       | low      |
| 241 | DID — Dissociated Identity for multi-personality agents                | enhancement    | identity         | growing  |
| 242 | `.cdd/` directory layout, lifecycle, retention                         | enhancement    | CDD-method       | low      |
| 243 | `cnos.cdd.kata` cleanup (M0/M4, remove runtime, honest stubs)          | enhancement    | kata             | low      |
| 244 | kata 1.0 master                                                        | enhancement    | kata             | low      |
| 245 | cdd-kata 1.0 (prove CDD adds value)                                    | enhancement    | kata             | growing  |
| 249 | `.cdd/` Phase 1 audit trail                                            | enhancement    | CDD-method       | low      |
| 255 | CDD 1.0 master                                                         | —              | CDD-method       | low      |
| 256 | CI as triadic surface (α produce / β verify / γ attest)                | —              | CI/proof         | growing  |
| 262 | `cn pack`: derive packlist from content-class model                    | enhancement    | package-runtime  | low      |
| 269 | Per-project notification channels for agent observability              | —              | comms (telegram) | —        |
| 271 | GitHub App for agent identity                                          | —              | identity         | —        |
| 272 | Close-out filename convention drift                                    | process,P2     | CDD-method       | new      |

**Counts:** CDD-method 7 · package-runtime 11 · kata 3 · CI/proof 2 · identity 2 · memory/cog 9 · threads/issues 8 · comms (telegram) 4 · diagnostics 4 · hub-placement 2 · vision 5 · bug-OCaml 11 · process-tiny 2 · stale 4. Sum = 79.

**Growing-lag (per 3.58.0 PRA):** 11 — #256 #245 #241 #238 #235 #230 #218 #216 #199 #193 + #268 (now shipped/closed-pending). MCI freeze remains in force.

---

## 2. Dependency graph (text)

Read `A → B` as "A depends on or is unblocked by B". Edges drawn from issue bodies' explicit `Depends on` / `Parent` / `Subsumes` / "blocked by" sections (not from speculative coupling).

**CDD-method chain (mostly converged in design, partially shipped):**

```
#168 (L8 evidence tracker, meta) ──────── independent
#272 (close-out filename) ── exposes ──> #268 (shipped) — needs operator decision
#249 (Phase 1 audit trail) ─ in use ──> shipped pattern (.cdd/releases/{v}/{role}/N.md)
   └─ depends on ──> #242 (.cdd/ layout)
#240 (triadic protocol) ─ depends on ──> #242
#242 ─ depends on ──> #241 (identity model — DID)
#256 (CI triadic) ─ consumes ──> #242 Phase 3 + #244 (kata 1.0)
#255 (CDD 1.0 master) ─ tracks ──> {#249 #242 #240 #245 #256 #241 #244}
#94 (mechanize CDD) ─ partial-ship via ──> cn-cdd-verify shipped in #268
```

**Kata chain:**

```
#244 (kata 1.0 master) ─ tracks ──> {#243 #245 #238 #235 #256 (#254 closed) (#253 closed)}
#243 (cnos.cdd.kata cleanup) ─ blocks ──> #245 (cdd-kata 1.0)
#245 ─ depends on ──> #244 framework
#238 (smoke) ─ independent test surface, complements #244/#245
#235 (cn build --check) ─ independent of katas, parallel quality gate
```

**Package-runtime chain:**

```
#192 (Go kernel rewrite) ─ Phase 4 complete; remaining ──> #216 + #193
#216 (migrate kernel commands to packages) ─ depends on ──> #192 Phase 4 + #235
#262 (cn pack packlist) ─ partially shipped ──> #266 (closed); deferred ACs ride forward
#230 (deps restore version skip) ─ independent one-function fix; complements #181
#181 (gh-pages package source) ─ depends on ──> #180 (stale, see closures)
#218 (cnos.transport.git Rust) ─ depends on ──> #192 + Phase 4 + provider contract
#193 (orchestrator llm step) ─ OCaml-era; depends on ──> #182 (stale) + Go orchestrator port
#170 (orchestrator IR master) ─ subsumes ──> #173 (closed) #174 (closed) #175
#175 (CTB → IR compiler) ─ depends on ──> #174 (closed) — re-eval after Go orchestrator
#162 (modular CLI commands) ─ effectively shipped via Go-port + #167 (closed); see closures
#186 (lean cnos.core + domain pkgs) ─ depends on ──> #182 (stale)
```

**Identity / multi-agent chain:**

```
#241 (DID) ─ first consumer ──> #240 (triadic protocol), #242 (cdd dir)
#271 (GitHub App) ─ near-term ──> #45 / #189 / #218 long-term replace it
#156 (attached hub) ─ Steps 1–3 landed in PR #163; Steps 4–10 remain
   └─ Step 5+ unblocks ──> #87 (Claude Code scheduled) and multi-agent ACs
```

**Threads / native issues / web:**

```
#45 (native issue tracking) ─ umbrella for ──> {#189 (PRs+reviews) #190 (web) #6 #7 #9}
#153 (thread event model master) ─ subsumes ──> #154 (Phase 1)
#189 ─ closes ──> review-identity gap (was #108) and converges with #271
#190 ─ presentation layer over ──> {#45 #189 #181}
#180 (pkg system doc incoherence) — claimed shipped (PR #204) per #192 body; see closures
```

**Diagnostics / cognition:**

```
#30 (temporal `now`) ─ prereq for ──> {#34 #35 #68}
#34 (runtime coherence loop) ─ uses ──> #30
#35 (daemon reflection cadence) ─ uses ──> #30 + #100
#100 (memory faculty) ─ adjacent to ──> #84 (mindset reflection req) + #15 (forget)
#59 (deep cn doctor) ─ extends ──> shipped doctor; complements #68 anomaly layer
#68 (self-diagnose) ─ depends on ──> #59 + runtime contract surfaces
#149 (UIE skill loading) ─ mindset patch, independent
#84 (CA mindset reflection) ─ mindset patch, independent
#101 (skill corpus normalize) + #132 (rename agent→self/ops→hub) ─ pair: rename then conform
#52 (Role as load-bearing) ─ independent classification
#124 (permission-asking) ─ model-level; outside cnos source
```

**Comms (telegram):**

```
#42 (media parity) ─ Phase 1 (vision/D-sev) ─ subsumes ──> {#44 #43 partial}
#71 (TTS voice) ─ depends on ──> #42 outbound media + network access
#269 (project notification channels) ─ independent observability surface
#39 (Pi→Sigma OCaml peer bug) ─ stale (Go runtime / #218 supersedes)
```

**Vision (post-substrate):**

```
#79 (projection surfaces) ─ enabled by ──> #45 + git-public substrate
#77 (P2P network) ─ enabled by ──> {#218 #100 #45 #189}
#69 (first-wake setup) ─ enabled by ──> #62 (closed) + #65 (closed)
#70 (personality archetypes) ─ enabled by ──> #69
#71 (voice) — see comms
```

**Hub-placement:**

```
#156 Steps 4–10 ─ unblocks ──> sandboxed-agent runtime (Claude Code on web etc.)
#87 ─ converges with ──> #156 attached-hub flow
```

**Bug-OCaml cluster (mostly stale — see §5 closures):**

```
#11 #12 #20 #135 #138 #141 #142 #148 #4 #17  ─ all OCaml-runtime fixes superseded by #192 Go port
#39 ─ same (peer transport will be #218)
#16 ─ doctor-warn (carries forward as Go-doctor enhancement; folds into #59)
#8 ─ log rotation (folds into Go runtime ops surface)
```

**Cross-cluster meta-edges:**

- #244 (kata) ⇄ #255 (CDD): both reference #245 and #256.
- #45 (native issues) ⇄ #271 (GH App): #271 is the near-term identity hack the long-term #45/#189 surface eliminates.
- #240 (triadic protocol) ⇄ #156 (attached hub) ⇄ #241 (DID): three angles on "one agent, multiple roles, multiple workspaces".
- #181 (gh-pages) ⇄ #218 (git transport): both rework "where does a package come from"; should be sequenced not parallelized.

---

## 3. Proposed bundles

Each bundle is sized for one CDD triad cycle (one issue → one α/β pair → one release) unless explicitly marked **multi-cycle**. Bundles that exceed one cycle are split with a named split point.

Size legend: **immediate** = ship as a small-change-direct (no triad needed); **small** = single cycle, ≤ ~150 lines diff; **substantial** = single cycle, mostly docs/code under 600 lines; **multi-PR** = needs ≥2 cycles (split below).

### B1 — Distribution chain honesty *(growing-lag impl, MCI-shrinking)*

- **Coherence theme:** make the manifest → lockfile → installed substrate trust itself again. Today a lockfile bump is silently skipped (#230) and there is no check that the released binary can actually bootstrap a fresh hub (#238).
- **Constituents:** **#230** (`cn deps restore` version-skip P1, one-function fix) → **#238** (smoke release-bootstrap CI job).
- **Size:** **substantial** (single cycle). #230 is small (~30 lines + tests); #238 adds one CI workflow + one shell script. Coherent because both touch the same install authority chain (lockfile vs vendor vs released binary), and #238 will fail loudly if #230 ever regresses.
- **Depends on:** none active; both are inside the Go kernel surface that #192 already shipped.
- **Subsumes / closes by consolidation:** —.

### B2 — Manifest ↔ filesystem honesty *(growing-lag impl, MCI-shrinking)*

- **Theme:** every manifest claim becomes a verifiable filesystem fact at build time.
- **Constituents:** **#235** (`cn build --check` validates command entrypoints + skill paths against manifest, P1).
- **Size:** **small** (single cycle). The fix is one validator pass in `pkgbuild/build.go::CheckOne` + tests. Already named in 3.58.0 PRA as the natural pick after #230.
- **Depends on:** —. Independent of B1 but a natural pair (both close drift between declared and installed state).

### B3 — `.cdd/` Phase 1 + 2 formalization *(low-lag spec; folds CI prep)*

- **Theme:** lock down the `.cdd/releases/{X.Y.Z}/{role}/N.md` pattern that is already in use, so future cycles inherit it as a verifiable contract rather than as practice.
- **Constituents:** **#249** (Phase 1 audit trail) + **#242** (directory layout, lifecycle, retention) + **#272** (close-out filename convention — operator decision A/B/C, then matrix + verifier sync).
- **Size:** **substantial** (single cycle, but only if the operator has already picked the #272 option before α dispatch). The spec-text edits are small; the `cn-cdd-verify` fixture work is already in place from #268.
- **Depends on:** #272 needs the operator decision named in the issue body before α can land it.
- **Subsumes / closes by consolidation:** finishes the #249 + #242 partial implementations the 3.58.0 PRA records as "low lag, partially shipped."

### B4 — CI as triadic surface (Phase 3) *(growing-lag impl)*

- **Theme:** CI jobs are renamed to α-produce / β-verify / γ-attest; per-run proof bundles land under `.cdd/runs/`.
- **Constituents:** **#256** (CI triadic) → consumes **B3** outputs (the Phase 1/2 layout) and **#244** kata framework.
- **Size:** **multi-PR**. Split: **B4a** = workflow rename + α-produce job + γ-attest receipts (one cycle). **B4b** = per-run JSON evidence + README rollups + dist/.cdd/ split enforcement (one cycle).
- **Depends on:** **B3** (must land first so the artifact paths are pinned).

### B5 — Kata 1.0 closure *(growing-lag impl + design)*

- **Theme:** ship the executable proof that CDD adds value over ad-hoc development.
- **Constituents:** **#243** (cleanup: M0/M4, remove runtime, honest stubs) → **#245** (cdd-kata 1.0: M0–M4 prompts + rubrics + first comparison run) → **#244** (master tracking; closes when #243 + #245 land).
- **Size:** **multi-PR**. Split: **B5a** = #243 cleanup (small). **B5b** = #245 prompt + rubric authoring + judge wiring + first real run (substantial). #244 closes by tracking, not by separate work.
- **Depends on:** **B4** for CI integration of the kata verdict in `.cdd/runs/{run}/beta/kata-certification.json` (otherwise B5 ships stand-alone but without the CI hook).

### B6 — Triadic protocol (git transport) + DID identity *(post-freeze design)*

- **Theme:** make role handoffs run on artifact presence rather than operator paste; each role commits with a distinct identity.
- **Constituents:** **#240** (CDD triadic protocol) + **#241** (DID identity) + **#271** (GitHub App as the near-term platform shim until git-native identity exists).
- **Size:** **multi-PR**. Split: **B6a** = #271 GitHub App (operational unblock for native review events + tag push + branch delete; closes the 4-cycle HTTP 403 pattern). **B6b** = #240 + #241 design + first execution model running end-to-end on a real cycle (substantial design + impl).
- **Depends on:** **B3** (artifact location matrix), **B4** (proof-bundle surface), **B5** (kata as sanity check on a triadic run).
- **Closes by consolidation:** #271's purpose is subsumed by #240/#241 long-term; ship #271 as deliberate scaffolding with a documented sunset path.

### B7 — Lean kernel + package commands *(growing-lag impl)*

- **Theme:** finish the boundary the Go port (#192) drew between bootstrap kernel and package-provided commands.
- **Constituents:** **#216** (migrate `status`/`doctor`/`build`/`setup` to `cnos.core/commands/`) — implicitly depends on Phase 4 package-command discovery in Go (sub-issue of #192 not yet filed as its own issue).
- **Size:** **multi-PR**. Split: **B7a** = Phase 4 discovery in Go (file as new issue if not already tracked under #192). **B7b** = mechanical extraction of the four commands and CI verification that `cn status` works identically whether kernel-built or package-provided.
- **Depends on:** **B2** (`cn build --check` would catch missing entrypoints during the migration).

### B8 — Package source + transport *(growing-lag impl + design)*

- **Theme:** decouple the consumer endpoint from GitHub Releases; later, decouple from GitHub entirely.
- **Constituents:** **#181** (gh-pages package source — near-term consumer endpoint) + **#218** (`cnos.transport.git` Rust provider — long-term polyglot transport).
- **Size:** **multi-PR**. Split: **B8a** = #181 (substantial, single cycle, narrow CI workflow + restore-path change). **B8b** = #218 (multi-cycle on its own — first Rust package, first provider, polyglot surface).
- **Depends on:** **B1** (deterministic + version-honest restore must hold before adding a second source).
- **Closes by consolidation:** #180 (β doc claims git-native, α doc says artifact-first) is recorded by #192 as already shipped via PR #204; see §5 closures.

### B9 — Memory + reflection cadence + UIE *(post-freeze cognition)*

- **Theme:** make persistent retention + scheduled self-review part of the runtime, not of the human's prompting.
- **Constituents:** **#100** (memory faculty) + **#84** (CA mindset reflection requirement, mindset patch) + **#35** (daemon periodic reflection) + **#149** (UIE skill loading rule). Adjacent tiny: **#15** (forget skill).
- **Size:** **multi-PR**. Split: **B9a** = #84 + #149 mindset/skill text patches (immediate, two-line edits each, can ship same commit as B3). **B9b** = #100 memory faculty + skill (substantial, single cycle). **B9c** = #35 daemon scheduler + #15 forget skill (substantial).
- **Depends on:** B9c needs #30 (`now` injection — already-shipped via OCaml; verify Go port).

### B10 — Self-diagnostics + temporal + doctor depth *(post-freeze cognition)*

- **Theme:** the agent can answer "what is my current state?" with one structured op and one CLI command, both honest by construction.
- **Constituents:** **#68** (full self-diagnose op + `cn diagnose` CLI) + **#59** (doctor deep semantic checks) + **#30** (timestamp injection — verify Go-port equivalent) + **#16** (warn if peers + no recent sync) + **#34** (CMP/MCP/CAP/CLP runtime vocabulary).
- **Size:** **multi-PR**. Split: **B10a** = #59 + #16 (small, doctor enhancements that re-use the activation pipeline shipped in #261). **B10b** = #68 + #30 verification (substantial, new op + CLI surface). **B10c** = #34 trace-event vocabulary (substantial, touches runtime trace schema; can defer indefinitely).
- **Depends on:** **B7** (because doctor migrates to a package command), **B9** (memory faculty exposes some of the diagnostic surfaces).

### B11 — Telegram media + interrupt + streaming *(parity with OpenClaw)*

- **Theme:** close the inbound/outbound media gap so cnos can replace OpenClaw for first-class users.
- **Constituents:** **#42** (media parity, 6 phases internally) + **#44** (streaming output) + **#43** (interrupt) + **#71** (TTS voice).
- **Size:** **multi-PR**. Split: **B11a** = #42 Phase 1 (vision: inbound photo + caption + LLM mixed content). **B11b** = #42 Phase 2 + #43 (reply/forward context + interrupt). **B11c** = #42 Phase 3 + #44 + #71 (outbound media + streaming + TTS).
- **Depends on:** Go runtime parity with OCaml Telegram surface (verify before scoping).

### B12 — Native threads + native issues + web surface *(post-Go-port substrate)*

- **Theme:** replace GitHub Issues / PRs / Reviews / web with thread-based equivalents on the cnos substrate.
- **Constituents:** **#45** (native issues) + **#189** (native PRs/reviews) + **#190** (web surface) + **#153** (thread event model master) + **#154** (Phase 1) + **#6** (`cn thread new`) + **#7** (processed archive) + **#9** (thread-based task system) + **#269** (project notification channels).
- **Size:** **multi-PR (large)**. Split: **B12a** = #154 Phase 1 (event-store substrate); **B12b** = #45 + #6 + #7 + #9 (issue surface as threads); **B12c** = #189 (PR + review surface); **B12d** = #190 (web surface); **B12e** = #269 (notifications).
- **Depends on:** **B6** (DID + identity) for review identity, **B8a** (gh-pages) for #190 hosting story.

### B13 — Hub placement + multi-agent + scheduled cron *(operational)*

- **Theme:** finish the work started in PR #163 so attached-hub mode is fully usable for sandboxed agents and multi-agent setups.
- **Constituents:** **#156** Steps 4–10 (currently 7 of 10 steps remain) + **#87** (Claude Code scheduled cron, partially superseded by γ-on-the-web pattern).
- **Size:** **multi-PR**. Split: **B13a** = #156 Steps 4–6 (`cn attach`, executor placement-aware, runtime contract integration). **B13b** = #156 Steps 7–10 (sync, deps, doctor, init --attached). **B13c** = #87 rescoped against current scheduling primitives.
- **Depends on:** —. Independent of the freeze; impl-side, eligible whenever capacity exists.

### B14 — Skill corpus normalize + rename *(post-Go-port hygiene)*

- **Theme:** every artifact under `skills/` has an explicit class (skill / runbook / reference / deprecated) and conforms to the meta-skill where it claims skill status; category names match what's inside.
- **Constituents:** **#101** (normalize corpus) + **#132** (rename `agent/`→`self/`, `ops/`→`hub/`) + **#52** (Role as load-bearing classification) + **#18** (troubleshooting skill) + **#15** (forget skill, also in B9).
- **Size:** **multi-PR**. Split: **B14a** = #132 rename (small, mechanical, but cross-surface — needs `cn build --check` from B2 first to catch path drift). **B14b** = #101 classification pass + Phase 1 quick-win rewrites. **B14c** = #101 mid-range + boundary cases + #52 + #18.
- **Depends on:** **B2** (so the rename can be validated mechanically).

### B15 — Vision tier *(post-substrate, evidence-bound)*

- **Theme:** projection surfaces and P2P network — load-bearing only after substrate stabilizes.
- **Constituents:** **#79** (projection surfaces) + **#77** (P2P network) + **#69** (first-wake setup) + **#70** (personality archetypes) + **#168** (L8 evidence tracker, meta).
- **Size:** **deliberately deferred**. None should be picked while MCI freeze is in force.
- **Depends on:** **B6** + **B8** + **B12** + **B9**.

### B16 — Process-tiny patches *(immediate, ship as direct commits)*

- **Theme:** small spec/skill text edits that don't deserve a triad cycle each.
- **Constituents:** **#5** (require review commit before merging, skill text patch) + **#14** (gitignore + pre-commit hook for secrets, mechanical) + **#84** (CA mindset reflection — also in B9a) + **#149** (UIE skill loading — also in B9a).
- **Size:** **immediate** — bundle into a single direct-to-main "skill text sweep" commit, mirroring the 4-patch sweep in the 3.58.0 PRA.
- **Depends on:** —. Eligible to ship under MCI freeze because they're all MCI-class skill patches that derive from observed compliance gaps, not new design.

---

## 4. Execution order

**Ordering rule (per CDD §3 + 3.58.0 PRA freeze rationale):** while MCI freeze is active, prioritize bundles that **shrink the growing-lag count first**, then bundles that **unblock other bundles**, then bundles that **earn substrate leverage** (L7-shaped). Defer net-new design bundles until lag count drops below the 3-issue threshold.

| Order | Bundle                                                       | Type                | Why now                                                                                                                                                                                            | Lag delta when shipped                                                                  |
| ----- | ------------------------------------------------------------ | ------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------- |
| 1     | **B16** process-tiny patches                                 | immediate-direct    | All four patches are MCI-class; can ship in one commit alongside the next PRA without consuming a triad cycle. Mirrors the 3.58.0 sweep precedent.                                                | none (MCI, not lag-counted)                                                               |
| 2     | **B3** `.cdd/` Phase 1 + 2 + #272                            | substantial         | Highest-leverage MCI work: locks in the close-out filename + matrix + verifier loop the 3.58.0 PRA flagged. Must precede B4. **Operator decision on #272 needed before α dispatch.**              | -3 (#249 + #242 → "shipped"; #272 closes immediately)                                     |
| 3     | **B1** distribution chain honesty (#230 + #238)              | substantial         | Two growing-lag P1/feature items in one cycle, both inside the kernel surface #192 already shipped. Adjacent code paths.                                                                          | -2 (#230 + #238 → "shipped")                                                              |
| 4     | **B2** `cn build --check` (#235)                             | small               | One-function P1 fix that pays for itself by guarding B7 + B14 from regressions.                                                                                                                   | -1 (#235 → "shipped")                                                                     |
| 5     | **B4a** CI triadic — workflow rename + α-produce + γ-attest   | substantial         | Consumes B3 outputs; halves #256.                                                                                                                                                                  | partial #256                                                                              |
| 6     | **B5a** + **B4b** — kata cleanup (#243) + CI proof bundles    | substantial         | Pair: #243 is the cnos.cdd.kata cleanup blocker for #245; B4b lights up the kata-certification.json surface that #245 needs.                                                                      | -1 (#243 → "shipped"; #256 fully closes)                                                  |
| 7     | **B5b** cdd-kata 1.0 (#245 + #244 closes)                     | substantial         | M0–M4 prompt + rubric + first comparison run. Unblocks #244 master tracking close.                                                                                                                 | -2 (#245 → "shipped"; #244 closes)                                                        |
| 8     | **B6a** GitHub App (#271)                                     | substantial         | Operational unblocker: native review events + tag push + branch delete. Closes the 4-cycle HTTP 403 pattern β #267 close-out flagged. Not a growing-lag item but high friction return.            | none (operational)                                                                        |
| 9     | **B7a + B7b** lean kernel + package commands (#216)           | multi-PR (2 cycles) | Earns L7 boundary move (kernel bootstrap = 4 commands; everything else is package-provided).                                                                                                       | -1 (#216 → "shipped")                                                                     |
| 10    | **B8a** gh-pages package source (#181)                        | substantial         | Decouples consumer from Releases; depends on B1's deterministic build chain holding.                                                                                                                | partial growing list (post-#180 close, see §5)                                            |
| 11    | **B9a** mindset/skill text (B16 already shipped these)        | —                   | Already covered by B16; no separate bundle needed.                                                                                                                                                  | —                                                                                          |
| 12    | **B6b** triadic protocol + DID (#240 + #241)                  | substantial design  | After kata + audit trail + identity surface (B3/B5/B6a) are stable, design ships against real evidence. **Only after lag count drops below 3.**                                                   | -2 (#240 + #241 → "shipped")                                                              |
| 13    | **B10a** doctor depth (#59 + #16) + #34 verify                | small               | Cleanup; folds into B7-era doctor work.                                                                                                                                                              | none (low-lag items)                                                                      |
| 14    | **B9b** memory faculty (#100)                                  | substantial         | Substantial design but evidence-bound; defer until B10a + freeze relaxed.                                                                                                                            | none                                                                                      |
| 15    | **B8b** cnos.transport.git Rust (#218)                         | multi-PR (2–3 cyc)  | First polyglot package; first provider; depends on B7 + B8a + provider-contract surfaces.                                                                                                            | -1 (#218 → "shipped")                                                                     |
| 16    | **B7-adjacent** `cn pack` deferred ACs (#262 carry)            | small               | β #265 close-out's deferred AC1 + AC5 + AC3. Land alongside B14a rename when convenient.                                                                                                             | none                                                                                      |
| 17    | **B14a** category rename + **B14b** corpus normalize           | multi-PR            | After B2 + B7 land (so the rename is mechanically validated), do the L6 cross-surface rename, then the classification pass.                                                                          | partial #101                                                                              |
| 18    | **B12a** thread event Phase 1 (#154)                           | substantial         | Substrate prerequisite for B12 chain; can run after kernel work calms.                                                                                                                              | none                                                                                      |
| 19    | **B11a** Telegram media Phase 1 (#42 vision)                   | substantial         | High operational return; verify Go runtime parity with OCaml Telegram first.                                                                                                                          | none                                                                                      |
| 20    | **B13a** hub-placement #156 Steps 4–6                          | substantial         | Unblocks sandboxed-agent flows; eligible whenever capacity exists.                                                                                                                                   | none                                                                                      |
| 21    | **B10b–c**, **B9c**, **B11b–c**, **B12b–e**, **B13b–c**        | substantial→large   | Sequence by then-current evidence and freeze status; all post-substrate.                                                                                                                            | none                                                                                      |
| 22    | **B15** vision tier (#79 #77 #69 #70 #168)                     | deferred            | Do not pick during freeze. Eligible only after B6b + B12 ship.                                                                                                                                       | none                                                                                      |

**Rationale highlights:**

1. **B16 first, immediately, by direct commit.** The 3.58.0 PRA already established the pattern of bundling small skill-patch sweeps into a single non-cycle commit. #5 + #14 + #84 + #149 are all MCI text edits with clear derive-from-evidence paths.

2. **B3 next as a triad cycle**, because (a) it captures the operator decision #272 forces, (b) it converts two of the 11 growing-lag items (#249 + #242) into "shipped", (c) it pins the artifact surface every later bundle (B4, B5, B6, B12) depends on. Operator must pick option A/B/C in #272 before α dispatch — that decision sets the cycle's first AC.

3. **B1 + B2 next** because together they remove three growing-lag items (#230, #238, #235), all inside the Go kernel surface that's already stabilized, and both are "next-smallest growing-lag bug with a clear fix path" — the pattern the 3.58.0 PRA explicitly endorses.

4. **B6a (#271 GitHub App) earns its slot out of growing-lag order** because the harness-403 + shared-identity friction is now over the 4-cycle threshold β #267 close-out flagged. Cost: one operational cycle. Return: native review events + branch-cleanup unblock for every subsequent cycle.

5. **B6b (#240 + #241) waits for evidence.** The 3.58.0 PRA's MCI freeze rationale says "no new substantial design work until the MCA backlog drops below threshold." B6b is design-class; running it before B1 + B2 + B3 + B4 + B5 ship would re-trigger the freeze rationale.

6. **B7 (lean kernel + package commands) is the next L7 candidate** after the audit-trail + CI surface stabilizes. Pairs naturally with B8a (gh-pages) since both rework the "where commands and packages come from" boundary.

7. **B12 + B15 are deliberately last.** Both depend on substrate that is still shifting; landing them earlier would burn cycles re-doing them when the substrate moves.

**Freeze exit signal:** when growing-lag count drops below 3 (currently 11), the freeze lifts and B6b / B9b / B12 design bundles become eligible without operator override.

---

## 5. Recommended closures (close without work)

Each entry: issue, reason, evidence path. Closure is a triage-only action — γ can ship them in one batch comment under the same release as B16.

| #   | Reason                                       | Evidence                                                                                                                                                                                                                                                                                  |
| --- | -------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 180 | **Superseded — already shipped.**            | #192's "Prerequisites" section lists `[x] #180 doc coherence shipped (PR #204)`. The β doc was either retired or rewritten at PR #204 merge time. The issue was never closed by `Closes #180` in the squash body, so it stayed open as a paperwork artifact. **Close as superseded.**       |
| 182 | **Superseded by #192 (Go port complete).**   | #182 is the OCaml-era core-refactor master. #192 records "Phase 4 complete" in the 3.58.0 lag table, having implemented the substrate boundary #182 sketched. Remaining ACs ride into B7. **Close, redirect remaining ACs to #216 (already split off).**                                  |
| 199 | **Stale — 18 versions overdue.**             | Stacked v3.39.0 + v3.40.0 PRAs. We are now at 3.58.0 with PRAs landing per-release. Backfilling two-year-old assessments has no operator return; the encoding-lag table on every subsequent PRA already covers what those assessments would have said. **Close as stale.**                |
| 11  | **Obsolete — OCaml deprecated.**             | "Replace blanket `with _ -> None` catches" targets `cn.ml`, `inbox.ml`, `peer_sync.ml` — all under `src/ocaml/` which §5.0 of `DESIGN-CONSTRAINTS.md` declares deprecated and the 3.57.0 + 3.58.0 PRAs explicitly preserve as "correctly superseded rather than extended." **Close.**       |
| 12  | **Obsolete — OCaml deprecated.**             | Same pattern — `inbox.ml` mutable ref → fold. OCaml deprecated. **Close.**                                                                                                                                                                                                                |
| 20  | **Obsolete — OCaml runtime superseded.**     | Targets `cn_hub.ml::log_action`, `cn_runtime.ml::run_daemon`, `cn_config.ml`. Go runtime supersedes the entire daemon path. Any diagnostic gap that survived the port is captured by B10a (#59 + #16). **Close.**                                                                          |
| 39  | **Obsolete — OCaml peer-sync.**              | "Pi→Sigma `clone not found`, 17 undelivered messages" — OCaml `cn sync`. Replaced by Go transport substrate; long-term replacement is #218. **Close, fold remaining concern into #218.**                                                                                                  |
| 141 | **Obsolete — OCaml peer-sync.**              | `cn_mail.ml` empty-content branches limbo. OCaml. **Close, fold into #218.**                                                                                                                                                                                                              |
| 142 | **Obsolete — OCaml peer-sync.**              | FidoNet parity for OCaml peer-sync. Receipts/relay/dead-letter all become #218 design concerns once the Rust transport ships. **Close, fold into #218.**                                                                                                                                  |
| 148 | **Obsolete — OCaml binary.**                 | `check_binary_version_drift` parses output from `cn --version` in the OCaml daemon. Go runtime path is different. Verify Go-side behavior; if equivalent gap exists, file fresh under B10a. **Close.**                                                                                    |
| 138 | **Effectively done.**                        | AC1–AC4 of the four ACs are checked. The unchecked "CI green" was a v3.x-era box; CI has been green for many cycles since. **Close as done.**                                                                                                                                              |
| 135 | **Obsolete — OCaml runtime trace.**           | `cn_ulog.ml`, `cn_runtime.ml` per-pass logging. OCaml runtime. Re-spec under #34 (or B10c) if the Go trace surface needs it. **Close.**                                                                                                                                                   |
| 4   | **Likely-fixed (verify, then close).**         | "Daily threads should persist to hub" — likely resolved by hub/workspace split work in #156 Steps 1–3. γ should verify with one `git log -- threads/daily/` check; if persistence holds, close.                                                                                              |
| 8   | **Obsolete — OCaml log path.**                 | `/etc/logrotate.d/cn` for OCaml `cn.log`. Go runtime logging surface is different. Re-file under B10c if rotation is still needed.                                                                                                                                                          |
| 10  | **Stale — ROADMAP.md doesn't exist as named.** | `ROADMAP.md` referencing v1.6.0 / Protocol v1. The current encoding-lag tables in PRAs serve the roadmap function. **Close as stale.**                                                                                                                                                    |
| 16  | **Fold into B10a.**                            | "warn if peers but no recent sync" — restate against Go doctor + activation surface. Filing intent preserved in B10a; close-and-redirect.                                                                                                                                                  |
| 17  | **Stale — installer + heartbeat both moved.**  | `cn-cron`, OpenClaw heartbeat assumption. Cron management was removed by #138 itself; heartbeat surface has moved. **Close, re-file if a Go-era equivalent gap surfaces.**                                                                                                                |
| 162 | **Effectively shipped.**                       | "Modular CLI commands — add commands without modifying cnos core" — addressed by Phase 4 package-command discovery in Go (work tracked under #192 + #216). **Close, redirect to #216.**                                                                                                  |
| 170 | **Stale — re-evaluate after Go orchestrator.** | Master for OCaml-era orchestrator IR work. Sub-issues #173, #174 are closed; #175 + #193 remain stale-pending. The Go runtime does not yet have an orchestrator surface. **Close master as stale; #175 + #193 stay open as "OCaml-era, re-spec on Go orchestrator port."**               |
| 175 | **Stale-pending — keep open with note.**       | CTB → orchestrator IR compiler. Depends on Go orchestrator surface that doesn't exist yet. **Keep open; relabel as "stale-pending #192 successor" so it doesn't show in growing-lag.**                                                                                                    |
| 193 | **Stale-pending — keep open with note.**       | OCaml `llm` step execution. Same story — Go orchestrator port re-opens scope. **Keep open; relabel to remove from growing-lag once the OCaml-era pin is acknowledged.** Currently the 3.58.0 lag table inflates by counting it.                                                          |
| 87  | **Partially overtaken.**                       | "Claude Code scheduled cron" — γ-on-the-web pattern (per the 3.58.0 PRA's polling/Monitor work) covers most of the scheduled-agent intent. The remaining `/schedule` design surface is real but small. **Reduce scope to the unmet portion or close in favor of B13c restatement.**       |
| 168 | **Keep, but flag low.**                        | L8 evidence tracker is a meta-issue by design. The 3.58.0 PRA itself is L7-shaped without L8 framing; evidence is accumulating but not yet "5+ clear examples." **Keep open; do not bundle.**                                                                                              |

**Net closure:** 18 issues recommended for outright closure (`#180 #182 #199 #11 #12 #20 #39 #141 #142 #148 #138 #135 #4 #8 #10 #16 #17 #162`); 4 issues recommended for relabel-don't-close (`#170 #175 #193 #87`); 1 issue (`#168`) keeps current state with low-priority flag.

**After closures:** 79 → 61 open issues. Growing-lag table cleanup: removing #199 (closed) and relabeling #193 (out of lag) drops 11 → 9 even before any new B-cycle ships.

---

## 6. Summary view

| Metric                                         | Today                                       | After §5 closures + B16 ship      | After B3 + B1 + B2 ship              |
| ---------------------------------------------- | ------------------------------------------- | --------------------------------- | ------------------------------------- |
| Open issues                                    | 79                                          | 61                                | 55                                    |
| Growing-lag count                              | 11                                          | 9                                 | 4                                     |
| MCI freeze status                              | active (≥3)                                  | active                            | active (still ≥3, but trajectory clear) |
| Bundles eligible to dispatch                  | B16, B3, B1, B2, B6a, B7a (queue)            | same                              | + B4, B5a, B6a accelerates             |
| First L7 surface still pending                  | none clean — every L7 candidate carries dependency | B3 done unlocks B4 / B6     | B7a unlocks lean-kernel L7 move        |

---

Signed: α (`alpha@cdd.cnos`) · 2026-04-25 · branch `claude/consolidate-issue-bundles-mOkM7` · doc at `docs/gamma/cdd/ISSUE-CONSOLIDATION-ANALYSIS.md`.
