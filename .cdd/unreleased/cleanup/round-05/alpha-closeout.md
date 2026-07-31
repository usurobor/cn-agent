# α Closeout — Round 05 (docs-only legacy-narrative removal)

**Role:** α (implementer). **Branch:** `claude/repo-cleanup-newcomer`.
**Scope:** the docs-only, zero-code-coupling subset of β findings C2, C3, C6, C7, C8, S4.
**Result:** 46 files deleted, 18 files edited, all resulting broken links repaired, 3 items flagged for the code pass. No defer-list / code path touched.

---

## 1. Falsification gate (pre-flight)

Every deletion target was confirmed present with `git ls-files` before removal. Every "no inbound refs remain" claim below is backed by the `git grep` in §7. One discrepancy with the prompt's "ZERO code coupling / already dependency-verified" premise was found and acted on per the gate's "trust the live tree" rule — see §6 (three targets carry inbound references from defer-list paths and were **not** deleted).

---

## 2. C2 — `docs/development/plans/` (16 files) — DELETED, directory gone

Method: read each plan's `Status:` and its named target, then checked whether the target's canonical current-state doc exists in `docs/reference/*` or `docs/architecture/*`. The repo is at **VERSION 3.82.0**; every plan is version-keyed to v3.3–v3.22 (long shipped). The "Draft" stamps are stale metadata never flipped — corroborated by internal evidence (the syscall plan states "v3.6.0 … is shipped"; the docs-governance plan says "retroactive — plan created after implementation"; the package-system plan's header points to its retired→current reference doc). **No plan describes a not-yet-built target, so no live design was orphaned** — each design already has a canonical home. All 16 deleted; none relocated.

| Plan | Stamp | Current-state home (already exists) | Decision |
|---|---|---|---|
| `CAR-implementation-plan.md` | Implemented | `architecture/cognitive-substrate/CAR.md` | deleted |
| `PLAN.md` (CN Shell v3.3) | Complete | `reference/runtime/AGENT-RUNTIME.md` | deleted |
| `PLAN-v3.22.0-eng-lane-clarity.md` | Complete | `development/ENGINEERING-LEVELS.md` | deleted |
| `issue-41-pass-b-wiring.md` | Complete | `reference/runtime/AGENT-RUNTIME.md` (N-pass) | deleted |
| `TRACEABILITY-implementation-plan.md` | Implementing (1-9 done) | `architecture/security/TRACEABILITY.md` | deleted |
| `PLAN-package-system.md` | Draft (v3.18) | `reference/packages/PACKAGE-SYSTEM.md` | deleted |
| `PLAN-runtime-extensions.md` | Draft (#73) | `reference/runtime/extensions/RUNTIME-EXTENSIONS.md` | deleted |
| `PLAN-runtime-contract-v2.md` | (Issue #62) | `reference/runtime/RUNTIME-CONTRACT-v2.md` | deleted |
| `PLAN-v3.10.0-runtime-contract.md` | Draft | `reference/runtime/RUNTIME-CONTRACT-v2.md` | deleted |
| `PLAN-v3.6.0.md` (Output Plane Sep) | Draft | `reference/runtime/AGENT-RUNTIME.md` (shipped v3.6.0) | deleted |
| `PLAN-v3.7.0-scheduler.md` | Draft | `reference/runtime/AGENT-RUNTIME.md` (v3.7.0 changelog) | deleted |
| `PLAN-v3.8.0-n-pass-bind.md` | Draft | `reference/runtime/AGENT-RUNTIME.md` (v3.8.0) | deleted |
| `PLAN-v3.8.0-syscall-surface.md` | Draft | `reference/runtime/AGENT-RUNTIME.md` (v3.8.0) | deleted |
| `PLAN-v3.13.0-docs-governance.md` | (retroactive) | `reference/governance/DOCUMENTATION-SYSTEM.md` | deleted |
| `PR-docs-governance-v3.13.0.md` | (PR desc) | `reference/governance/DOCUMENTATION-SYSTEM.md` | deleted |
| `INVARIANT-HARDENING-v1.md` | Active | `architecture/INVARIANTS.md` / `development/rules/INVARIANTS.md` | deleted |

`docs/development/plans/` no longer exists.

---

## 3. C3 — historical assessments / audits / runs

**Deleted (12):**
- `docs/development/cdd/POST-RELEASE-EPOCH-v3.12.md`, `POST-RELEASE-EPOCH-v3.14.md` — retroactive epoch assessments.
- `docs/evidence/AUDIT.md` — dated (2026-02-11) point-in-time docs audit.
- `docs/evidence/rca/**` (8 dated incident write-ups).
- `docs/evidence/design/WRITER-PACKAGE.md` — self-labeled "Superseded … remains as cycle-input evidence, not binding design."
- `docs/evidence/smoke/cds-dispatch-smoke-20260623.md` — dated smoke receipt.
- `docs/development/kata/runs/**` (15 run/score/DELTA logs).

**`docs/papers/ENGINEERING-LEVEL-ASSESSMENT.md` — NOT deleted, FLAGGED for the code pass** (defer-path coupling; see §6).

**KEPT in `docs/evidence/` and why:** `docs/evidence/README.md` only. The directory stays because `src/packages/cnos.cdd/skills/cdd/CELL-KINDS.md` (defer) names `docs/evidence` as a live cell-closure target, and `docs/README.md` + `docs/quickstart/README.md` nav to it. Deleting the directory would break a defer-path reference I cannot repair. The README was rewritten to state what the surface **is** (a live closure surface) and to drop its now-dead "Available now: AUDIT / rca" index and its future-promise / "leave DRAFT" narration.

**KEPT in `docs/development/kata/` and why (templates + definitions, 9 files):** `TEMPLATE-KATA-PACKET.md`, `TEMPLATE-RUN-RECORD.md`, `TEMPLATE-SCORE-SHEET.md`, `KATA-EVALUATION.md`, and the five kata definitions `A1-open-op-registry-conflict.md`, `A2-extension-registry-engine-compatibility.md`, `B1-runtime-contract-v2-parity-review.md`, `C1-capability-growth-boundary.md`, `C2-browser-capability-and-ecosystem-boundary.md`. These are current-state doctrine (the reusable practice definitions and scoring rubric); only the dated `runs/` execution logs were point-in-time and removed. Confirmed the kept files reference no `runs/` path.

---

## 4. C6 — future / superseded docs

**Deleted (1):** `docs/papers/MODULAR-ARCHITECTURE-REFACTOR.md` (`Status: Proposed` refactor decision record).

**NOT deleted, FLAGGED for the code pass (defer-path coupling; see §6):**
- `docs/reference/ctb/CTB-v4.0.0-VISION.md`
- `docs/reference/ctb/LANGUAGE-SPEC-v0.2-draft.md`

---

## 5. S4 / C7 / C8 — narrative & directory edits

**S4 — DELETED:** `docs/reference/legacy/` (whole dir; `OCAML-THREAD-REFERENCE.md`). The root README Legacy paragraph is the one home for the fact.

**C7 — retirement narrative → positive current fact (edited):**
- `docs/reference/packages/BUILD-AND-DIST.md` — removed the entire "Migration from the prior layout" section (blockquote "This migration has landed. `src/agent/` is gone…" plus Prior state / Why this breaks / Target state / What disappears / Migration steps / Sequencing). The current authoring layout is already stated positively above under "Directory Structure → Source (authored)"; the doc now ends on "The package declares compatibility. The runtime enforces it."
- `docs/reference/governance/GLOSSARY.md` — dropped "those directories were retired in the docs cleanup and their contents now live under intent directories (e.g. …)"; kept the current rule "There is no `docs/alpha/`, `docs/beta/`, or `docs/gamma/`."
- `docs/reference/governance/DOCUMENTATION-SYSTEM.md` — CI-rule line: dropped "— those were retired," kept the rule "No document declares an α/β/γ folder path (…)." Also removed the three `development/plans/` references left dangling by C2 (the `development/` layout row's "plans", the bundle-README "Links to related plans in `development/plans/`" requirement, and the "Plan" document-class row whose home was the deleted dir).

**C8 — inert authoring-intent lines (edited):**
- `docs/papers/CCNF-AND-TYPED-TRUST.md` — removed the "### Wave 0 — Design essay" block ("Land this document under `docs/gamma/essays/`…"); Waves 1-7 (real system roadmap) kept.
- `docs/papers/DECREASING-INCOHERENCE.md` — removed the "### Wave 0 — Land this essay" block ("Add this file under `docs/gamma/essays/`…"); Waves 1-6 kept.

---

## 6. Flagged for the code pass (deviations from the prompt's deletion list)

The prompt framed this subset as "ZERO code coupling (already dependency-verified)." The live tree contradicts that for three C3/C6 targets: each has an inbound reference from a **defer-list** path (`src/**` or `schemas/**`) that I am forbidden to edit. Per the mandate rule ("if a doc-side edit would require changing a defer path to stay correct, STOP that edit and note it for the code pass") and the "all links must work" constraint, I did **not** delete them — deleting would leave a dangling reference I cannot repair. Keeping them also keeps all their many editable-doc inbound refs valid, so no repair was needed.

| Target (NOT deleted) | Inbound ref in defer path | Code-pass action |
|---|---|---|
| `docs/papers/ENGINEERING-LEVEL-ASSESSMENT.md` | `src/packages/cnos.eng/skills/eng/README.md:9,277` | delete doc + update skill README |
| `docs/reference/ctb/CTB-v4.0.0-VISION.md` | `src/packages/cnos.core/skills/agent/emoji-language/SKILL.md:70` (markdown link) | delete doc + repoint/remove skill link; also repair editable refs in `ctb/README.md`, `LANGUAGE-SPEC.md`, `SEMANTICS-NOTES.md`, `reference/runtime/ORCHESTRATORS.md` |
| `docs/reference/ctb/LANGUAGE-SPEC-v0.2-draft.md` | `schemas/README.md:184` | delete doc + update schemas README; also repair editable refs in `ctb/README.md`, `papers/ACTIVATION-NOT-DEPLOYMENT.md` |

**Two more code-pass notes (no live-doc breakage, defer paths untouched):**
- `src/packages/cnos.cdd/skills/cdd/design/SKILL.md:257` — an illustrative example string ("… Plan: `PLAN-package-system.md` …"), not a link. Now references a deleted basename; harmless but stale.
- `docs/reference/governance/DOCUMENTATION-SYSTEM.md §5 "Frozen history"` (esp. the clause "a completed plan … is left in place under its dated or `Superseded` header. Its stale internal paths are not corrected") is the exact "frozen = leave" doctrine this round's mandate reverses. It is doctrine, not a dangling link, and outside my enumerated C7 scope — flagged for reconciliation by the loop/doctrine owner rather than edited here.

---

## 7. Inbound-link repair list (ref → fix) and zero-dangling evidence

All repairs are in editable docs (root README + `docs/**`); no defer path was edited.

| Deleted target | Inbound ref repaired | Fix |
|---|---|---|
| `OCAML-THREAD-REFERENCE.md` | `README.md` Legacy ¶ | dropped dead link; kept branch/tag fact |
| " | `docs/architecture/ARCHITECTURE.md:13` | dropped link; kept branch/tag fact |
| " | `docs/architecture/DESIGN-CONSTRAINTS.md:47` | dropped "— see [..]" clause |
| " | `docs/reference/cli/CLI.md:175` | dropped link; kept branch pointer |
| `AUDIT.md` | `docs/architecture/ARCHITECTURE.md:272` | removed related-docs row |
| " | `docs/architecture/security/TRACEABILITY.md:743` | removed bullet |
| " | `docs/evidence/README.md` | rewrote README (dropped index) |
| `POST-RELEASE-EPOCH-v3.12/14` | `docs/development/cdd/README.md:43,44` | removed both rows |
| `MODULAR-ARCHITECTURE-REFACTOR.md` | `docs/papers/README.md:42` | removed bullet |
| " | `docs/reference/runtime/CORE-REFACTOR.md:529` | removed bullet |
| plans (C2) | `docs/development/README.md` | removed `plans/` bullet |
| " | `docs/reference/governance/DOCUMENTATION-SYSTEM.md:22,48,64` | removed layout "plans", bundle-README plans requirement, Plan doc-class row |
| " | `docs/reference/runtime/AGENT-RUNTIME.md:24,46,1537` | dropped 3 dead PLAN links, kept spec pointers |
| " | `docs/reference/runtime/README.md:25-34` | removed whole "Related Plans" section |
| " | `docs/development/rules/INVARIANTS.md:200` | removed dead `INVARIANT-HARDENING-v1.md` bullet |

**Final zero-dangling grep** (`git grep` over the tree, excluding `.cdd/`, `.cn-sigma/`, `CHANGELOG.md`, `RELEASE.md`, and code `rca`-skill paths):

```
$ git grep -nE "OCAML-THREAD-REFERENCE|reference/legacy|MODULAR-ARCHITECTURE-REFACTOR|
   POST-RELEASE-EPOCH|evidence/AUDIT|evidence/rca|WRITER-PACKAGE|evidence/design|
   evidence/smoke|cds-dispatch-smoke|kata/runs|development/plans|plans/PLAN|PLAN-v3.|
   PLAN-package-system|PLAN-runtime-extensions|PLAN-runtime-contract-v2|PR-docs-governance|
   CAR-implementation-plan|TRACEABILITY-implementation-plan|INVARIANT-HARDENING|
   issue-41-pass-b-wiring|A1-cold-001|A1-DELTA|B1-selected-001|C1-DELTA|C1-cold-001"
→ src/packages/cnos.cdd/skills/cdd/design/SKILL.md:257   (illustrative example, defer path — flagged, not a link)
```

The only surviving hit is the flagged defer-path example above. Zero dangling references in the live reader surface. Confirmed no new dead link introduced: no doc links into the removed BUILD-AND-DIST "Migration" section or the removed runtime "Related Plans" section; `docs/evidence/` still resolves (README kept). `CHANGELOG.md`/`RELEASE.md` mentions of deleted basenames are prose (not `[](…)` links) in defer files — left untouched by mandate.

---

## 8. Hard-defer confirmation

`git status --short` shows every changed path under `docs/**` or root `README.md`. **Untouched:** `.cdd/**` (except this receipt, the mandated write location), `.cn-sigma/**`, root `RELEASE.md` `CHANGELOG.md` `ROLES.md` `OPERATOR.md` `SUSTAINABILITY.md`, `docs/development/board/`, `src/**`, `schemas/**`, `scripts/**`, `.github/**`, `tests/**`, `install.sh`, `cn.json`.

**Totals:** 46 deleted (16 plans + 2 epoch + AUDIT + 8 rca + WRITER-PACKAGE + smoke + 15 kata runs + MODULAR-ARCHITECTURE-REFACTOR + OCAML-THREAD-REFERENCE), 18 edited, 3 deletions flagged for the code pass. Not committed; no PR opened.
