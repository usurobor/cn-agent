# γ Closeout — Round 08 (verify α trim + independent convergence judgment)

**Role:** γ (verifier / closeout). **Branch:** `claude/repo-cleanup-newcomer`. **HEAD:** `b32e2e72`.
**Method:** Re-derived every claim from `git show b32e2e72` / `git diff b32e2e72^ b32e2e72` and the live tree — did not trust α's receipt. Edited only this file; no commit; no PR.
**Mandate:** de-historicize the **non-dot** reader surface to current-state only; git history is the archive. Dotdirs EXEMPT. Deferred code pass out of scope. Standard: cnos write skill.

**Headline:** α's 11 edits are all VERIFIED-correct with **zero lost substance and zero α↔tree discrepancy** — every claim in the receipt reconciles with the tree. But the round does **NOT converge**: my independent sweep found a real, previously-uncleared residual **class of 4 canonical specs** carrying the same CDD design-cycle apparatus in a **numbered/embedded form** (`## 0. Coherence Contract`, `## N. Acceptance Criteria`, `## N. File Changes`, `## N. Known Debt`) that every prior round's unnumbered-heading greps — and α's Round-8 residual grep — silently skipped. One of them (`THREAD-EVENT-MODEL.md`) was in α's Round-8 target set but mis-scoped as a stamp-only *sibling*, so α (faithfully following the R7 contract) removed only its `Issue:` line and left the full apparatus. **Verdict: REVISE — Round 9 needed.**

---

## 1. Per-check verification verdicts (all re-derived from the tree)

### Defer-path integrity — **VERIFIED**
`git show --name-status b32e2e72` = **1 `A`** (`round-08/alpha-closeout.md`) + **11 `M`**, all 11 under `docs/`. `git show --name-only` minus (`docs/**`, receipt) → **empty**. **Zero deletions** (all 11 are edits — matches "TRIM round"). Diffstat: 122 ins / 571 del total; minus the 85-line receipt = **37 doc ins / 571 del**, exactly α's claim. No `.cdd/` (bar the receipt), `.cn-sigma/`, `src/`, `schemas/`, `scripts/`, `tests/`, `.github/`, no root `RELEASE/CHANGELOG/ROLES/OPERATOR/SUSTAINABILITY`, no `docs/development/board/`, none of the 6 code-coupled do-not-delete docs, not the CDD `design` template. Untouched.

### Substance preservation — the 6 primary specs — **VERIFIED (no lost substance)**
Read `git diff b32e2e72^ b32e2e72 -- <file>` for each; every removed block is genuine process apparatus, and each authority/mechanism fact that lived *only* in apparatus was confirmed restated in a kept body section on the live tree:
- **MEMORY.md** — removed `Issue/Mode/Active Skills/Engineering Level` stamps, `## Problem` gap-narration (→ 1-para `## Purpose`), `## File Changes`, `## Acceptance Criteria` (7 `[ ]`), `## Known Debt` (3), `## CDD Trace`. Body intact: `## Current Evidence` (L73), `### Core memory semantics` (L110) — so the Problem's "core vs optional memory" facts survive. AC9 pointer repaired (`#156 … (agent memory stays in hub, tagged by workspace)`); superseded-narration dropped from `## Related`. No spec substance lost.
- **HYBRID-LLM-ROUTING.md** — removed stamps, `## Impact Graph`, `## File Changes`, `## Acceptance Criteria` (AC1–10), `## Known Debt` (4), `## CDD Trace`. The one Impact-Graph authority nuance ("commands/orchestrators may declare routing hints, but the runtime makes the final decision") is **confirmed folded into Proposal §1** in the diff. `## Tuning Strategy` kept. No loss.
- **POLYGLOT-PACKAGES-AND-PROVIDERS.md** — removed stamps, `## Problem`(→`## Purpose`, keeping Go/Rust/packages breakdown + two-shapes-avoided), `## File Changes`, `## Acceptance Criteria` (8 bullets), `## Known Debt` (4), `## CDD Trace`, trailing "the next useful artifact is…" chatter. `## Decision` + Proposal §1–9 + Alternatives kept. No loss.
- **PACKAGE-ARTIFACTS.md** — de-historicized title (`# Design: …`→`# …`), removed stamps, `## Problem` incident forensics (#155/`cn_deps.ml` 27-refs/#162 → 2-bullet `## Purpose`), `## Impact Graph`, `## File Changes`, `## Acceptance Criteria` (AC1–10), `## CDD Trace`. Index authority restated at L267 ("The index is the package-resolution authority"); precedence authority at L337. Challenged Assumption + Proposal + Migration kept. No loss.
- **MESSAGE-PACKET-TRANSPORT.md** — removed stamps, `## Problem` (Pi/Sigma incident + `get_branch_files` narration → `## Purpose` + `### Failure modes it defends against`, 8 modes kept verbatim), `## Impact Graph`, `## 15. Acceptance Criteria`. Envelope authority restated at L224 ("`packet/envelope.json` is the **sole** transport-authoritative manifest") + frontmatter-ignored rule L258/260; mechanism restated in kept §13 "Why This Is Bulletproof". No loss.
- **HUB-PLACEMENT-MODELS.md** — removed stamps, `## Problem` (root-ownership lists + "why raw submodule isn't enough" → `## Purpose` keeping 3 problem classes as "problems the model prevents"), `## Impact Graph`, `## File Changes`, `## Acceptance Criteria` (10 `[ ]`). Root-ownership lists confirmed restated in kept `## 4. Root Semantics` (L146: 4.1 hub_root, 4.2 workspace_root, run-against lists L180/L194). No loss.

### Limitations handling — **VERIFIED**
- PACKAGE-ARTIFACTS: `## Known Debt`→`## Limitations`, **4 genuine current limits kept** verbatim (flat index / no third-party hosting / SHA-256-only signing / exact-match versioning).
- MESSAGE-PACKET: `## 16. Known Debt`→`## 15. Limitations`, **3 kept** (dedup GC undesigned / attachments unsupported / `cn.packet.v2` undesigned); dropped only the 1 pure migration-compat item (owned by kept §14 Migration Plan).
- HUB-PLACEMENT: `## Known Debt`→`## Limitations`, **2 kept** (multi-workspace attachments unsupported / CI backend quirks unhandled); dropped only the migration-sequencing item.
- MEMORY / HYBRID / POLYGLOT: Known Debt removed entirely — re-read each; all dropped items are pure future-work (restore-map "may become useful later", recall/search "worth packaging later", zone-name "tightening pass"; token-estimation/depth/confidence tuning; provider-protocol-spec/index-format-"later"/first-migration-"future work"/surfacing-"follow-up"). No current limitation dropped.

### Coherence — **VERIFIED**
- MESSAGE-PACKET numbered headings now run **12 → 13 → 14 → 15 (Limitations) → 16 (Summary)**, contiguous; `git grep '§1[4-7]|section 1[4-7]'` → **0 dangling refs** to the renumbered/removed sections.
- MEMORY `## Related` block coherent, AC9 pointer repaired, no dangling bullet.
- Residual-apparatus grep (`^**Issue|Mode|Active Skills|Engineering Level|Author|Contributors|Date:`, `^## Acceptance Criteria|CDD Trace|Impact Graph|Known Debt|File Changes`) across **all 11 edited docs → clean**.
- No fragment/anchor link in the repo targets any removed heading (α's grep re-confirmed for the removed-heading set).

### Sibling sweep — **PARTIALLY VERIFIED (one mis-scope surfaced, see §2)**
Stamp-only trims confirmed correct and content-preserving:
- **BUILD-AND-DIST.md** — removed `**Issue:** #219, #186`; kept Version/Status. ✔
- **GIT-CN-PACKAGE.md** — removed `**Issue:** #218`; kept Version/Status/Parent. ✔
- **ORCHESTRATORS.md** — removed `**Issue:** #170`; kept Version/Status/Doc-Class/Canonical-Path/Owns. Carries `## 0. Purpose` + a lean `## 1. Problem` (current-state motivation: "cnos currently has X but lacks Y") and NO AC/Impact/Debt/Trace stamps — a defensible KEEP, though the `## 1. Problem` is mildly redundant with §0 Purpose. ✔ (noted nit, non-blocking)
- **PROTOCOL.md** — removed `Author/Contributors/Date`; **kept `**Status:** Implemented (v2 — cn_protocol.ml matches this spec)`** as the intended code-correspondence marker. ✔ (title retains `# Design:` prefix — left because sibling scope is stamp-only, not retitle; minor, non-blocking.)
- **THREAD-EVENT-MODEL.md** — removed `**Issue:** #153`; kept Version/Status/Purpose/Related. **✗ SCOPE DEFECT:** this doc is NOT a light-stamp sibling — it carries the **full CDD apparatus** (`## 0. Coherence Contract` [Gap/Named incoherence], `## 20. File Changes` [create/edit plan incl. `src/…` paths], `## 21. Acceptance Criteria` [11 unchecked `[ ]`], `## 22. Known Debt` [6 future items]). α followed the R7 stamp-only sibling contract exactly, so this is a **contract mis-scope**, not an α execution error — but the net effect is that a full-apparatus doc remains under-trimmed on the reader surface. See §2.

**α-scoped conclusion:** everything α actually did is VERIFIED — no lost substance, receipt totals reconcile exactly, no defer path touched. The only gap is **coverage/scope**, addressed in §2.

---

## 2. Independent convergence hunt — NEW in-scope residual found (a class of 4)

Swept the entire live non-dot reader surface with the mandated apparatus/history signatures **plus the numbered and embedded variants prior rounds never grepped for** (`^## [0-9]+\. (Acceptance Criteria|File Changes|Known Debt|Impact Graph|CDD Trace)`, `## 0. Coherence Contract`, `### α / β / γ target`, `**Mode:** MC*`, `Smallest coherent intervention`, `Failure if skipped`).

**Why prior rounds (and α's R8 residual grep) missed it:** every earlier sweep matched *unnumbered* headings (`^## Acceptance Criteria`) and *header* stamps (`**Mode:** MCA`). This residual wears the apparatus as **numbered body sections** (`## 21. Acceptance Criteria`) and as a **`## 0. Coherence Contract`** opening block whose fields are `### Gap` / `### Mode` / `### α / β / γ target` / `Smallest coherent intervention` / (`Expected effect` / `Failure if skipped`). Same semantic class as the `## Problem` + `**Mode:** MCA` + `## CDD Trace` α just trimmed — different surface syntax, so the greps sailed past it.

### FILE — in-scope residual, docs-only, TRIM (canonical specs stay; no deletion):

1. **`docs/reference/protocol/cn/THREAD-EVENT-MODEL.md`** — `## 0. Coherence Contract` (Gap/Named incoherence); `## 20. File Changes` (create/edit plan, `src/transport/cn_thread_event.ml` etc.); `## 21. Acceptance Criteria` (11 unchecked `[ ]`); `## 22. Known Debt` (6). *Touched in R8 but only the `Issue:` stamp removed — the mis-scoped sibling.*
2. **`docs/reference/packages/EXTENSION-REGISTRY.md`** — `## 0. Coherence Contract` (Gap / `### Mode` MCA / `### α / β / γ target` / Smallest coherent intervention); `## 20. Acceptance Criteria` ("This is done when: 1…6"). *Never in any round's target set.*
3. **`docs/reference/runtime/extensions/RUNTIME-EXTENSIONS.md`** — `## 0. Coherence Contract` (Gap [cites #67] / `### Mode` MCA / `### α / β / γ target` / Smallest coherent intervention); `## 15. Acceptance Criteria` ("This is done when: 1…7"). *Never fully trimmed.*
4. **`docs/architecture/cognitive-substrate/COGNITIVE-SUBSTRATE.md`** — trailing `## Coherence Contract for This Document` (`**Gap:**` / `**Mode:** MCI` / `**Scope:**` / `**Expected effect:**` α-β-γ / `**Failure if skipped:**`). *Never in any round's target set.*

All 4 are **canonical current-state specs** (governed `Version:`/`Status:`) — the sole home for their subject, so the fix is a **trim, not a deletion**, applying α's exact R8 recipe: strip `## 0. Coherence Contract`/trailing Coherence Contract (reframe to a current-state `## Purpose`), `## N. File Changes`, `## N. Acceptance Criteria`, `## N. CDD Trace`; convert `## N. Known Debt`→`## Limitations` per-item (keep genuine current limits, drop future-work); **keep** the spec body + design-rationale sections (Leverage/Negative Leverage/Alternatives/Non-goals/Process Cost/Migration) + governed markers. No code coupling (none referenced from `src|schemas|scripts|tests` for a link that a heading-trim would break); no deletions, so no inbound-link repair. Renumber remaining sections contiguously and re-grep `§N` cross-refs (as α did for MESSAGE-PACKET).

The scoping grep is **bounded to exactly these 4** — no numbered `## N. Impact Graph`/`## N. CDD Trace` exists elsewhere in-scope, and no fifth doc carries the Coherence-Contract/α-β-γ/Mode-MC signatures.

### Cleared as NOT residual (KEEP — current state)
- `SELF-COHERENCE-TEMPLATE.md` `## Acceptance Criteria Check` / `## Known Debt` — a **template** (placeholder scaffolding); templates are current-state tooling and the CDD design template is explicitly exempt. KEEP.
- `## 0. Coherence Contract` as a *governed convention*? **Checked — it is not.** No doctrine/template prescribes it; it appears on exactly these 4 docs as leftover per-artifact design framing, not a required footer. So trimming it does not violate a convention. (Contrast the `Supersedes:`/`Superseded by:` lineage stamps, which ARE a governed convention → KEEP.)
- `PACKAGE-SYSTEM.md:9` "(retired — redirect stub)" — current-state fact about a redirect stub, KEEP.
- "Triadic Self-Coherence" / TSC references (LINEAGE, ORIGIN, SEMANTICS-NOTES, GIT-AS-LOWEST-DURABLE-SUBSTRATE, papers) — the named methodology/system, current-state doctrine, KEEP.
- ORCHESTRATORS `## 1. Problem`, PROTOCOL `# Design:` title prefix — noted nits (§1), non-blocking.

**Out-of-scope (noted):** DESIGN-266 and the 6 code-coupled docs still carry apparatus but are HARD-DEFER to the code pass. Handle there.

---

## 3. Convergence verdict

**REVISE — Round 9 needed.**

α's Round-8 matter is fully VERIFIED with zero regressions and zero lost substance: all 11 edits confirmed apparatus-only, all authority/mechanism facts restated in kept body sections, limitations correctly preserved-and-renamed, numbering coherent, no defer path touched, receipt totals reconcile exactly. **The whole-doc historical-record class (R1–7) and the header-stamp CDD apparatus on the 6 primaries (R8) are exhausted.**

But the non-dot reader surface does **not yet reflect only current state**: a bounded, terminal residual **class of 4 canonical specs** still carries the CDD design-cycle apparatus in numbered/embedded form (`## 0. Coherence Contract` with `### Mode`/`### α / β / γ target`, numbered `## N. Acceptance Criteria`/`File Changes`/`Known Debt`). It is the exact class the mandate forbids, living in `docs/reference/` and `docs/architecture/` — not dotdirs, not the deferred set. `THREAD-EVENT-MODEL.md` proves it is real and reachable: it was in the R8 target set but the R7 contract mis-scoped it as a light-stamp sibling, so it was under-trimmed. This is docs-only, needs no deletions, and once Round 9 trims these 4 (same recipe), the reader surface converges — the scoping grep confirms no fifth doc and no other apparatus class remains.

---

## 4. Round-08 cell completion statement

The Round-08 coherence cell is **closed**:
- **α** produced the matter (`b32e2e72`: 11 edited, 37 doc-ins / 571 del) and the receipt (`round-08/alpha-closeout.md`), executing the γ Round-07 §5 contract (6-primary apparatus trim + 5-sibling stamp sweep) **exactly and correctly** — zero lost substance.
- **γ** (this closeout) re-derived every claim from `git show b32e2e72` and the live tree: all α-scoped checks VERIFIED, no α↔tree discrepancy, receipt totals reconcile. Independent sweep surfaced a residual class of 4 the R7 sweep had missed → explicit REVISE hand-forward.

Matter landed; receipt verified; cell closed. Per COHERENCE-CELL doctrine this closed cell is trusted because its claims validate against the tree — and it emits the Round-9 α-matter (the 4-doc apparatus trim) plus the carried-forward code-pass agenda below.

## Overall docs-pass status (rounds 1–8)
**NOT yet converged — one bounded docs-only round remains (R9).** Rounds 1–7 exhausted the whole-doc historical-record/design-record/dated-audit class; Round 8 cleared the header-stamp CDD apparatus on the 6 primary specs + 5 siblings. The single outstanding docs-only class is the **numbered/embedded CDD apparatus in 4 canonical specs** (§2). After R9 trims those 4, the non-dot reader surface is current-state-only modulo the code-pass-deferred set.

---

## 5. Round-9 α-matter + FINAL consolidated code-pass handoff agenda

**Round 9 (docs-only, TRIM — no deletions) — the terminal reader-surface matter:**
Trim the CDD design-cycle apparatus (numbered/embedded form) from these 4 canonical specs, preserving each spec body + rationale sections + governed `Version:`/`Status:`/`Doc-Class:`/`Owns:` markers; apply α's exact R8 recipe (reframe Coherence Contract→Purpose, drop File Changes/Acceptance Criteria/CDD Trace, Known Debt→Limitations per-item, renumber contiguously, re-grep `§N` cross-refs):
1. `docs/reference/protocol/cn/THREAD-EVENT-MODEL.md` — `## 0. Coherence Contract`; `## 20. File Changes`; `## 21. Acceptance Criteria`; `## 22. Known Debt` (Known Debt has real current gaps — reflections-outside-event-model, subscription UX deferred, feed pagination, chain/nostr locators, retention/GC undesigned — judge per-item: current-limitation → Limitations, pure-future/migration → drop).
2. `docs/reference/packages/EXTENSION-REGISTRY.md` — `## 0. Coherence Contract`; `## 20. Acceptance Criteria`. (Keep §18 Alternatives Considered, §19 Migration Strategy, §21 Summary.)
3. `docs/reference/runtime/extensions/RUNTIME-EXTENSIONS.md` — `## 0. Coherence Contract`; `## 15. Acceptance Criteria`. (Keep §13 Alternatives, §14 Migration Strategy, §16 Summary, §17 Marketplace-Readiness.)
4. `docs/architecture/cognitive-substrate/COGNITIVE-SUBSTRATE.md` — trailing `## Coherence Contract for This Document` (Gap/Mode: MCI/Scope/Expected effect/Failure if skipped). Keep §12 Compatibility and Migration (current-state normativity statement).
Optional low-risk nits while in-doc (adjudicate, don't force): reframe ORCHESTRATORS `## 1. Problem` (redundant with §0 Purpose); drop `# Design:` prefix on PROTOCOL title. No code coupling; no inbound-link repair.

**Deferred code pass (carried forward — complete, actionable to-do list; unchanged unless noted):**

*Code-coupled do-not-delete docs — delete each + repair the listed coupling refs:*
- [ ] `docs/reference/packages/DESIGN-266-dist-out-of-git.md` — delete; repair coupling `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md`. Deleting it also clears the only surviving DESIGN-227 / SELF-COHERENCE-227 prose mentions (its own lines 76 & 237) — no separate action.
- [ ] `docs/papers/ENGINEERING-LEVEL-ASSESSMENT.md` — repair `cnos.eng/skills/eng/README.md:9,277`.
- [ ] `docs/reference/ctb/CTB-v4.0.0-VISION.md` — repair `emoji-language/SKILL.md:70` + `ctb/README.md`, `LANGUAGE-SPEC.md`, `SEMANTICS-NOTES.md`, `ORCHESTRATORS.md`.
- [ ] `docs/reference/ctb/LANGUAGE-SPEC-v0.2-draft.md` — repair `schemas/README.md:184` + `ctb/README.md`, `papers/ACTIVATION-NOT-DEPLOYMENT.md`.
- [ ] `docs/development/cdd/DISPATCH-FAILURE-EVIDENCE.md` — repair `release-effector/SKILL.md:79`.
- [ ] `docs/development/design/cn-repo-install-MOCKS.md` — repair `cn-install-wake:40,1070,1082`.

*Directory / root lifecycle:*
- [ ] `docs/development/design/` — after `cn-repo-install-MOCKS.md` is deleted, the dir holds only its README + empty survey map; decide remove dir + README.
- [ ] `RELEASE.md` / `CHANGELOG.md` release-gate coupling (incl. dead `docs/beta/evidence/rca/` cite at `CHANGELOG.md:1336`) — de-historicize under the release-automation contract, not by hand.
- [ ] `docs/development/board/` — kanban snapshot, still present/untouched; regenerate-from-live or drop per board doctrine.
- [ ] Root relocations `ROLES.md` / `OPERATOR.md` / `SUSTAINABILITY.md` → `docs/`.

*src/doctrine narrative to reconcile (don't silently keep):*
- [ ] `docs/reference/governance/DOCUMENTATION-SYSTEM.md §5 "Frozen history"` + §"Supersedes/Superseded by" — still encodes the "leave completed/`Superseded` docs in place with stale paths" doctrine that THIS mandate reverses. **Doctrine conflict** — the loop/doctrine owner must reconcile it or every round re-litigates the frozen-bucket ruling. (The `Supersedes:`/`Superseded by:` lineage-stamp *convention* on `PACKAGE-SYSTEM.md`/`CAR.md` etc. is legitimate under *current* doctrine → KEEP unless doctrine changes.)
- [ ] **CDD `design`/`alpha` skill templates** (`src/packages/cnos.cdd/skills/cdd/design/SKILL.md`, CDD.md §5.4) — prescribe `## CDD Trace` / Engineering-Level / Acceptance Criteria / **Coherence Contract** (`Gap`/`Mode`/`α-β-γ target`) for design artifacts: **this is the upstream source of the R8 AND R9 apparatus.** When the code pass touches the CDD skills, decide whether artifacts promoted into `docs/reference/` should be stripped of this apparatus at promotion time — a process fix that prevents the residual recurring. (`cdd/design/SKILL.md:257` also has a stale illustrative `PLAN-package-system.md` example string — harmless, refresh when touched.)
