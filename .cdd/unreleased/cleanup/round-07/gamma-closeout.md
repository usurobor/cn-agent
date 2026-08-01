# γ Closeout — Round 07 (verify α + independent convergence judgment)

**Role:** γ (verifier / closeout). **Branch:** `claude/repo-cleanup-newcomer`. **HEAD:** `b50e19bd`.
**Method:** Re-derived every claim from `git show b50e19bd` and the live tree — did not trust α's receipt. Edited only this file; no commit; no PR.
**Mandate:** de-historicize the **non-dot** reader surface to current-state only; git history is the archive. Dotdirs (`.cdd/`, `.cn-sigma/`, `.github/`) EXEMPT. Deferred to code pass: `src/**`, `schemas/**` (top-level CUE), `scripts/**`, `tests/**`, root `RELEASE.md`/`CHANGELOG.md`, root `ROLES/OPERATOR/SUSTAINABILITY` relocations, `docs/development/board/`, and 6 code-coupled docs. Standard: cnos write skill.

---

## 1. Per-check verification verdicts (all re-derived from the tree)

### 7 deletions — **VERIFIED**
`git show --name-status b50e19bd` lists exactly **7 `D`**, **2 `M`**, **1 `A`**. Each of the 7 confirmed absent from the working tree AND untracked in `git ls-files`:
- `docs/reference/runtime/CORE-REFACTOR.md`
- `docs/reference/packages/DESIGN-227-distribution-pipeline.md`
- `docs/reference/packages/SELF-COHERENCE-227.md`
- `docs/reference/runtime/PLAN-174-orchestrator-runtime.md`
- `docs/reference/schemas/DESIGN-LLM-SCHEMA.md`
- `docs/reference/schemas/DESIGN-LLM-SCHEMA-README.md`
- `docs/papers/RELEASE-LEVEL-CLASSIFICATION.md`

### Historical-vs-current spot-checks — **VERIFIED (no misclassification)**
Re-derived each deleted file's class from `git show b50e19bd` body + confirmed the stated current-state home survives via `git ls-files`:
- **CORE-REFACTOR.md** — self-labels "Design record (#182), partly landed"; `## Proposal`/`### Proposed Package Roles`/`### Target Structure`/`## Impact Graph` + unchecked AC-checklist for a package-runtime refactor over the dead pre-Go OCaml tree. Genuine historical design record. Stated homes present: `docs/reference/packages/PACKAGE-SYSTEM.md`, `docs/reference/runtime/RUNTIME-CONTRACT-v2.md`, `docs/reference/runtime/ORCHESTRATORS.md` — all EXIST. No live spec orphaned.
- **DESIGN-227-distribution-pipeline.md** — "Design: Prove the Distribution Pipeline End-to-End", Issue #227, all 6 ACs checked `[x]` (subject shipped). Home `docs/reference/packages/BUILD-AND-DIST.md` EXISTS. Historical.
- **SELF-COHERENCE-227.md** — "Self-Coherence Report — #227", `Branch: claude/alpha-cdd-cycle-227-…`, AC-verification/invariants tables. Pure per-cycle production record — the `reference/` analogue of the deleted doctrine-cluster / `.cdd/` review records. Nothing to migrate. Historical.
- **PLAN-174 / DESIGN-LLM-SCHEMA ×2 / RELEASE-LEVEL-CLASSIFICATION** — spot-confirmed from commit body: a #174 orchestrator-runtime plan citing ORCHESTRATORS.md as its design home; a shipped-v3.2.0 LLM-schema design record + its bundle-index README (home = `schemas/` + Go runtime); and a dated `**Date**: 2026-03-28` audit classifying "all 60 releases v0.1.0–v3.24.0" (repo now far past that). All historical; none the sole current-state home for a live spec. **α's classification holds for all 7 — zero current-state specs deleted.**

### KEEPs — the 6 code-coupled do-not-delete docs — **VERIFIED**
All present via `git ls-files --error-unmatch`: `docs/papers/ENGINEERING-LEVEL-ASSESSMENT.md`, `docs/reference/ctb/CTB-v4.0.0-VISION.md`, `docs/reference/ctb/LANGUAGE-SPEC-v0.2-draft.md`, `docs/development/cdd/DISPATCH-FAILURE-EVIDENCE.md`, `docs/development/design/cn-repo-install-MOCKS.md`, and **`docs/reference/packages/DESIGN-266-dist-out-of-git.md`** (esp.). None deleted.

### Link integrity (incl. frontmatter / nav) — **VERIFIED**
Independently ran `git grep -nI <basename>` over `docs/**` + `README.md` for all 7 deleted basenames:
- `CORE-REFACTOR`, `PLAN-174-orchestrator-runtime`, `DESIGN-LLM-SCHEMA`, `DESIGN-LLM-SCHEMA-README`, `RELEASE-LEVEL-CLASSIFICATION` → **ZERO hits.**
- `DESIGN-227-distribution-pipeline` / `SELF-COHERENCE-227` → **2 hits each, both inside `docs/reference/packages/DESIGN-266-dist-out-of-git.md` (lines 76, 237).** These are prose backtick-mentions in DESIGN-266's "impact inventory" bullets (`…describe cn build producing dist/. Unchanged…`), **not** `[..](..)` links or `related:` cites — they do not 404. DESIGN-266 is a HARD-DEFER code-coupled path AND itself queued for code-pass deletion (§5 item 7); editing it to drop the mentions would touch a defer path, which the mandate forbids. **Not a defect — correctly left; clears itself when the code pass deletes DESIGN-266.** No non-defer dangling remains.
- **AGENT-NETWORK.md repair landed:** `docs/concepts/AGENT-NETWORK.md` `## Related` block (lines 155–160) no longer carries the `CORE-REFACTOR.md` bullet; the sibling `- #182 — Core refactor (package-driven runtime)` bullet preserves the reference by issue key. Block is coherent, no dangling bullet.
- **papers/README.md repair landed:** the `## Engineering and release` section (line 48) retains `- [ENGINEERING-LEVEL-ASSESSMENT]`; the `RELEASE-LEVEL-CLASSIFICATION` nav row is gone. Heading + list valid.
- **Frontmatter `related:`:** no surviving `docs/papers/*.md` `related:` block cites `RELEASE-LEVEL-CLASSIFICATION` or any other deleted basename (confirmed by the zero-hit grep above, which covers `related:` lines).

### Defer-path integrity — **VERIFIED**
- `git show --name-only b50e19bd` minus (`docs/**`, `round-07/alpha-closeout.md`) → **empty.**
- Top-level `schemas/` **untouched** — no `^schemas/` path in the commit (only `docs/reference/schemas/` doc files deleted).
- No path under `src/`, `scripts/`, `tests/`, `.github/`, no other `.cdd/`/`.cn-sigma/`, no `RELEASE.md`/`CHANGELOG.md`/`ROLES.md`/`OPERATOR.md`/`SUSTAINABILITY.md`, no `docs/development/board/`. None of the 6 code-coupled do-not-delete docs touched.
- Working tree clean apart from this untracked `round-07/` receipt.

**All α-scoped checks: VERIFIED. No α↔tree discrepancy.** Receipt totals (7 D / 2 M / 1 A) reconcile exactly with `git show`.

---

## 2. Independent convergence hunt (exhaustive sweep of the surviving non-dot reader surface)

Swept the **entire** live non-dot reader surface — not just dirs prior rounds hit. Greps run (with dotdir/`src`/`schemas`/`scripts`/`tests`/board/deferred-root exclusions):
`Status:`(Draft|Proposed|Superseded|Complete|Implemented); `retired|deprecated|superseded|no longer|used to|formerly|was removed|has landed|partly landed|Design record|Self-Coherence Report|prove…end-to-end`; dated headers (`**Date**|**Period**|^Date:|20260\d{3}`); issue-keyed filenames (`DESIGN-\d|PLAN-\d|SELF-COHERENCE|-vN|AUDIT|SURVEY|RCA|-REFACTOR`); and the CDD design-cycle apparatus (`Mode: MCA`, `Engineering Level`, `## Acceptance Criteria`, `## Impact Graph`, `## Known Debt`, `## CDD Trace`). Coverage: `docs/{quickstart,concepts,guides,reference,development,architecture,papers,evidence}` + root `*.md`.

### Cleared as NOT residual (KEEP — current state)
- **Whole-doc historical-record class is now EXHAUSTED.** Zero surviving issue-keyed DESIGN/PLAN/SELF-COHERENCE records, dated audits, per-cycle logs, or proposed-unbuilt-future docs remain as standalone files on the reader surface. Every issue-keyed/`-vN` filename that survives is either a current-state spec (`AGENT-ACTIVATION-LOG-v0`, `PROVIDER-CONTRACT-v1` [`Doc-Class: canonical-spec`], `RUNTIME-CONTRACT-v2`), a template (`SELF-COHERENCE-TEMPLATE.md`), a kata (`docs/development/kata/B1-runtime-contract-v2-parity-review.md` — a training exercise, not a dated record), or a deferred code-coupled doc (`CTB-v4.0.0-VISION`, `LANGUAGE-SPEC-v0.2-draft`, `DESIGN-266`).
- `Status: Draft` stamps on ~25 living specs — governed maturity-marker convention (`DOCUMENTATION-SYSTEM.md`), not per-doc history. KEEP.
- `no longer / deprecated / superseded / removed` hits — overwhelmingly current-state *fact* ("daemon no longer requires TELEGRAM_TOKEN", "`templates` was removed in #321", MIGRATION/GLOSSARY/NAMING/TAXONOMY/TRIAGE deprecation *guidance*, doctrine-essay prose, `DOCUMENTATION-SYSTEM.md §Supersedes` lineage-stamp convention). Current-state. KEEP.
- Dated `20260…` tokens — all illustrative thread/trigger/boot IDs inside runtime/protocol examples, or the empirical proof-points in `AGENT-ACTIVATION-LOG-v0.md` (current-state justification for a convention). KEEP.
- `docs/development/kata/**` "you are reviewing a branch that claims…" framing — hypothetical kata task statements, current-state teaching artifacts. KEEP.

### NEW in-scope residual found — one class: **CDD design-cycle apparatus embedded in canonical reference/architecture specs**
Prior rounds cleared *whole-doc* historical records. This round's sweep surfaces a lighter, within-doc residual the loop pre-classified but never cleared: a set of docs that **ARE the canonical current-state spec for their subject** (so they cannot be deleted) but still carry the full CDD design-cycle apparatus — header provenance stamps (`**Issue:** #… / TBD`, `**Mode:** MCA`, `**Active Skills:** cdd/design…`, `**Engineering Level:** L7`), a `## Problem` framing, `## Impact Graph`, an **unchecked** `## Acceptance Criteria` checklist, `## Known Debt` (future-work), and a `## CDD Trace` process table ("Step 0 Observe / 1 Select / … / 6 Artifacts — *Design drafted; no implementation plan included yet*"). This is point-in-time design-process history and proposed-future sitting on the primary reader surface (`docs/reference/`, the canonical current-spec dir) — structurally identical to the 7 files α just deleted, differing only in that each is the sole home for its subject, so the correct fix is a **trim, not a deletion**.

**FILE — in-scope residual, docs-only, trim (no deletion; canonical specs stay):**
1. `docs/reference/runtime/MEMORY.md` — trailing `## Acceptance Criteria` (7 unchecked `[ ]`), `## Known Debt`, `## CDD Trace` ("no implementation plan included yet").
2. `docs/reference/runtime/HYBRID-LLM-ROUTING.md` — `**Issue:** TBD`/`**Mode:** MCA`/`**Engineering Level:** L7` header; `## Impact Graph`; unchecked `## Acceptance Criteria`; `## Known Debt`; `## CDD Trace` ("implementation planning deferred").
3. `docs/reference/runtime/POLYGLOT-PACKAGES-AND-PROVIDERS.md` — `**Issue:** #170`/`**Mode:** MCA` header; unchecked `## Acceptance Criteria`; `## Known Debt`; `## CDD Trace` ("implementation sequencing remains future work").
4. `docs/reference/packages/PACKAGE-ARTIFACTS.md` — titled `# Design: …`; `**Issue:** #167`/`**Mode:** MCA` header; `## Impact Graph`; unchecked `## Acceptance Criteria`; `## Known Debt`; `## CDD Trace`.
5. `docs/reference/protocol/cn/MESSAGE-PACKET-TRANSPORT.md` — `**Mode:** MCA`/`**Engineering Level:**` header; `## Impact Graph` (partial apparatus; no CDD Trace tail).
6. `docs/architecture/HUB-PLACEMENT-MODELS.md` — `**Issue:** #156`/`**Mode:** MCA` header; `## Impact Graph`; `## Acceptance Criteria`; `## Known Debt` (no CDD Trace/EngLevel).

Trim = strip the design-cycle apparatus (`Issue:`/`Mode:`/`Active Skills:`/`Engineering Level:` header stamps; `## Impact Graph`; unchecked `## Acceptance Criteria`; `## Known Debt` futures; `## CDD Trace`), reframe any `## Problem` opening into a current-state Purpose, and **keep** the spec body + governed markers (`Version:`, `Status:`, `Doc-Class:`, `Owns:`/`Does-Not-Own:`). No file deleted; no code coupling (none of the 6 is referenced from `src|schemas|scripts|tests` — grepped clean). γ Round-6 §55/§96 pre-classified the MEMORY.md tail as a "low-priority trim," but under the mandate standard ("current-state only") an unchecked AC-checklist + a CDD Trace process table + `Issue: TBD` are not current state — they are residual history on the live surface, and they recur as a **class of 6** (not a one-off), so this is a genuine blocking residual, not an optional nicety.

**Out-of-scope (noted, not this round's matter):** `docs/reference/packages/DESIGN-266-dist-out-of-git.md` carries the same full apparatus but is a HARD-DEFER code-coupled path already queued for code-pass deletion (§5). Handle there.

---

## 3. Convergence verdict

**REVISE — Round 8 needed.**

α's Round-7 matter is fully VERIFIED with zero regressions and zero non-defer dangling links: all 7 deletions confirmed, all 7 confirmed historical (no current-state spec lost), both link repairs (AGENT-NETWORK Related bullet, papers/README nav row) landed and coherent, the only surviving inbound hits are defer-path prose mentions inside DESIGN-266, and no defer path was touched. **The whole-doc historical-record class — the loop's core purpose — is now exhausted on the reader surface.**

But the non-dot reader surface does **not yet reflect only current state**: the exhaustive sweep found a real, previously-uncleared residual **class of 6 canonical reference/architecture specs** carrying the CDD design-cycle apparatus (`Issue:`/`Mode: MCA` provenance stamps, unchecked `## Acceptance Criteria`, `## Known Debt` futures, `## CDD Trace` process tables). That apparatus is point-in-time design-process history and proposed-future — the exact class the mandate forbids — living in `docs/reference/` and `docs/architecture/`, not in dotdirs and not in the deferred set. It cannot be deleted (each doc is the canonical home for its subject) but must be **trimmed** to current state. This is a bounded, terminal, docs-only matter: once Round 8 trims these 6 (leaving DESIGN-266 to the code pass), the reader surface converges — there is no other outstanding class.

---

## 4. Round-07 cell completion statement

The Round-07 coherence cell is **closed**:
- **α** produced the matter (`b50e19bd`: 7 deleted, 2 edited) and the receipt (`round-07/alpha-closeout.md`), executing the γ Round-6 §2/§5 contract exactly (the 7 docs-only design/plan/self-coherence records + dated audit).
- **γ** (this closeout) re-derived every claim from `git show b50e19bd` and the live tree: all α-scoped checks VERIFIED, no α↔tree discrepancy, receipt totals reconcile.

Matter landed on the branch; receipt verified; cell closed with an explicit REVISE hand-forward. Per COHERENCE-CELL doctrine, this closed cell is trusted because its claims validate against the tree, not because a role approved it — and it emits the Round-8 α-matter (the 6-doc apparatus trim) plus the carried-forward code-pass agenda below.

---

## 5. Round-8 α-matter + final code-pass handoff agenda

**Round 8 (docs-only, TRIM — no deletions; from my §2 hunt) — the terminal reader-surface matter:**
- Trim the CDD design-cycle apparatus from these 6 canonical specs, preserving each current-state spec body + governed `Version:`/`Status:`/`Doc-Class:`/`Owns:` markers:
  1. `docs/reference/runtime/MEMORY.md`
  2. `docs/reference/runtime/HYBRID-LLM-ROUTING.md`
  3. `docs/reference/runtime/POLYGLOT-PACKAGES-AND-PROVIDERS.md`
  4. `docs/reference/packages/PACKAGE-ARTIFACTS.md`
  5. `docs/reference/protocol/cn/MESSAGE-PACKET-TRANSPORT.md`
  6. `docs/architecture/HUB-PLACEMENT-MODELS.md`
- Apparatus to strip per doc: `**Issue:**`/`**Mode:** MCA`/`**Active Skills:**`/`**Engineering Level:**` header stamps; `## Impact Graph`; unchecked `## Acceptance Criteria`; `## Known Debt` (futures); `## CDD Trace`. Reframe `## Problem` openings to a current-state Purpose where present. Also sweep sibling `reference/` specs γ R6 §55 flagged as carrying lighter `Issue:`/`Mode:` stamps (`BUILD-AND-DIST`, `PROTOCOL`, `THREAD-EVENT-MODEL`, `GIT-CN-PACKAGE`, `ORCHESTRATORS`) for the same header-stamp trim. No code coupling (grepped clean over `src|schemas|scripts|tests`); no inbound-link repair needed (no deletions). Adjudicate per-doc as α did in R7.

**Deferred code pass (carried forward — complete to-do; unchanged unless noted):**
- **6 code-coupled do-not-delete docs + their defer-path refs:**
  - `docs/reference/packages/DESIGN-266-dist-out-of-git.md` — delete + repair coupling `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md`. Deleting it also clears the DESIGN-227 / SELF-COHERENCE-227 prose mentions at its own lines 76 & 237 (the only surviving non-defer-adjacent references to those two deleted files). No separate action.
  - `docs/papers/ENGINEERING-LEVEL-ASSESSMENT.md` — repair `cnos.eng/skills/eng/README.md:9,277`.
  - `docs/reference/ctb/CTB-v4.0.0-VISION.md` — repair `emoji-language/SKILL.md:70` + `ctb/README.md`, `LANGUAGE-SPEC.md`, `SEMANTICS-NOTES.md`, `ORCHESTRATORS.md`.
  - `docs/reference/ctb/LANGUAGE-SPEC-v0.2-draft.md` — repair `schemas/README.md:184` + `ctb/README.md`, `papers/ACTIVATION-NOT-DEPLOYMENT.md`.
  - `docs/development/cdd/DISPATCH-FAILURE-EVIDENCE.md` — repair `release-effector/SKILL.md:79`.
  - `docs/development/design/cn-repo-install-MOCKS.md` — repair `cn-install-wake:40,1070,1082`.
- **`docs/development/design/` directory lifecycle:** after the code pass deletes `cn-repo-install-MOCKS.md`, the dir holds only its README with an empty survey map — decide whether to remove dir + README.
- **`RELEASE.md` / `CHANGELOG.md` release-gate coupling** (incl. dead `docs/beta/evidence/rca/` cite at `CHANGELOG.md:1336`).
- **`docs/development/board/`** — kanban snapshot, still present/untouched.
- **Root relocations** `ROLES.md` / `OPERATOR.md` / `SUSTAINABILITY.md` → `docs/`.
- **src/doctrine narrative to reconcile (don't silently keep):**
  - `docs/reference/governance/DOCUMENTATION-SYSTEM.md §5 "Frozen history"` + §"Supersedes/Superseded by" — still encodes the "leave completed/`Superseded` docs in place with stale paths" doctrine that THIS mandate reverses. Doctrine conflict; the loop/doctrine owner must reconcile it or every round re-litigates the frozen-bucket ruling. (The `Supersedes:`/`Superseded by:` lineage-stamp convention on `PACKAGE-SYSTEM.md`, `CAR.md` etc. is legitimate under *current* doctrine → KEEP unless the doctrine changes.)
  - `src/packages/cnos.cdd/skills/cdd/design/SKILL.md:257` — stale illustrative `PLAN-package-system.md` example string (harmless, non-link; refresh when touching the file).
  - NEW: the CDD `design` template (`cdd/design/SKILL.md`, CDD.md §5.4) prescribes the `## CDD Trace` / Engineering-Level / AC sections for design artifacts — the source of the Round-8 apparatus. When the code pass touches the CDD skills, decide whether design artifacts promoted into `docs/reference/` should be stripped of this apparatus at promotion time (a process fix that prevents the residual from recurring).
