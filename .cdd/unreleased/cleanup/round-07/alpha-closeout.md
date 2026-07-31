# α Closeout — Round 07 (execute γ Round-6 convergence-hunt matter)

**Role:** α (implementer). **Branch:** `claude/repo-cleanup-newcomer`. **HEAD at start:** `e941ea7c` (clean tree).
**Mandate:** de-historicize the **non-dot** reader surface to current-state only; git history is the archive. Dotdirs EXEMPT.
**Contract:** γ Round-6 closeout §2/§5 — the newly-found residual class *issue-keyed DESIGN/PLAN/SELF-COHERENCE design-and-plan records living in `docs/reference/`* + one dated audit in `docs/papers/`. 7 docs-only clean deletions (the 8th, DESIGN-266, is code-coupled → code pass).
**Result:** 7 files deleted, 2 inbound-link repairs, 0 KEEP-as-misclassification, 0 defer paths touched. No commit, no PR.

---

## 1. Per-file judgment (opened and read each before deleting)

All 7 live in `docs/reference/` or `docs/papers/` — the canonical current-spec dirs — so each was verified to be a genuine *historical record* (issue-keyed design/plan record for a shipped subject, per-cycle record, or dated audit) with a separate current-state home that EXISTS. None was the sole current-state home for a live/unbuilt spec. No γ misclassification.

1. **`docs/reference/runtime/CORE-REFACTOR.md`** — *History.* Self-labels line 11: "**Design record (#182), partly landed.**" Problem→Proposal→`### Proposed Package Roles`/`### Target Structure`/`## Impact Graph`/AC-checklist (all `[ ]`) for a package-driven-runtime refactor. Describes the pre-Go OCaml tree (`src/agent/`, `src/cmd/`, `.ml` modules) that the #192 Go kernel rewrite has since replaced — doubly historical (proposal + dead code layout). Current-state homes exist: `PACKAGE-SYSTEM.md` (content classes/roles), `RUNTIME-CONTRACT-v2.md` (activation_index/commands/orchestrators surfaces), `ORCHESTRATORS.md` (orchestrator surface). No live design orphaned. **Deleted.**
2. **`docs/reference/packages/DESIGN-227-distribution-pipeline.md`** — *History.* "Design: Prove the Distribution Pipeline End-to-End", Issue #227, Problem/Proposal/select-table, **all 6 ACs checked `[x]`** (subject shipped). Canonical current-state home = `BUILD-AND-DIST.md` (the doc cites it as the spec it proves against; that spec survives in `docs/reference/packages/`). Current behavior self-evident in `src/go/internal/{pkg,restore,pkgbuild}/`. **Deleted.**
3. **`docs/reference/packages/SELF-COHERENCE-227.md`** — *History.* "Self-Coherence Report — #227", `Branch: claude/alpha-cdd-cycle-227-vCB7e`, AC-verification table (all met), invariants-check table. A per-cycle production/self-coherence record — the `docs/reference/` analogue of the deleted doctrine cluster / `.cdd/` review records. No canonical-spec content; nothing to migrate. **Deleted.**
4. **`docs/reference/runtime/PLAN-174-orchestrator-runtime.md`** — *History.* "Plan: #174 — Orchestrator IR Runtime", Gap/Staging (Stage A/B/C, AC mapping). Cites `ORCHESTRATORS.md §7–8` as its *design reference* — the spec home exists and survives. References the pre-Go OCaml modules (`cn_orchestrator.ml`, `cn_build.ml`). Subject shipped v3.36.0 (per its own #174 note in CORE-REFACTOR). **Deleted.**
5. **`docs/reference/schemas/DESIGN-LLM-SCHEMA.md`** — *History.* "Design: Structured LLM Request Schema", Problem/Decision/`## Module changes` design record. Subject shipped v3.2.0 ("Structured LLM schema, multi-turn messages"). References pre-Go OCaml (`cn_llm.ml`, `cn_context.ml`). Subject home = `schemas/` (the JSON schema dir) + current Go runtime. **Deleted.**
6. **`docs/reference/schemas/DESIGN-LLM-SCHEMA-README.md`** — *History.* Bundle-index ("Schema — Feature Bundle") whose entire Document Map / Canonical-Spec pointer is the single design doc in #5 above; carries a "Version History" row (v3.14.4 bundle migration #89). It is the index *for* the deleted design record and nothing else — no independent current-state content. Deleting #5 without it would leave a README pointing at nothing. **Deleted** (pair).
7. **`docs/papers/RELEASE-LEVEL-CLASSIFICATION.md`** — *History.* Dated point-in-time audit: `**Date**: 2026-03-28`, `**Scope**: All 60 releases … v0.1.0–v3.24.0` (repo is far past that now), per-version evidence tables. Same class as the previously-deleted `evidence/AUDIT.md` / `CDD-PACKAGE-AUDIT.md`. Its rubric sources (`ENGINEERING-LEVELS.md`, `CDD.md §9.1`, `ENGINEERING-LEVEL-ASSESSMENT.md`) survive. **Deleted.**

**No KEEP / no misclassification:** every target is a shipped-issue design record, a per-cycle record, or a dated audit — not a current-state canonical spec. γ's classification holds for all 7.

---

## 2. Inbound-link repair list (surviving reader surface)

Method: `git grep -n <basename>` over `docs/** README.md` before and after deletion. Nav/index READMEs checked directly: `docs/reference/runtime/README.md` Document Map does **not** list CORE-REFACTOR or PLAN-174; `docs/reference/packages/README.md` does not exist; `docs/reference/schemas/` has no README. So the only inbound links were:

| Deleted file | Inbound ref | Fix |
|---|---|---|
| CORE-REFACTOR.md | `docs/concepts/AGENT-NETWORK.md:161` — `## Related` bullet `- CORE-REFACTOR.md — Architecture that enables this vision` | **Removed the dangling bullet.** The sibling bullet `- #182 — Core refactor (package-driven runtime)` already carries the same reference by issue key, so the Related list stays coherent (no bullet left pointing at nothing). |
| RELEASE-LEVEL-CLASSIFICATION.md | `docs/papers/README.md:49` — nav row `- [RELEASE-LEVEL-CLASSIFICATION](./RELEASE-LEVEL-CLASSIFICATION.md)` under "Engineering and release" | **Removed the nav row.** Section retains `- [ENGINEERING-LEVEL-ASSESSMENT]`; heading and list stay valid. |
| DESIGN-227 / SELF-COHERENCE-227 / PLAN-174 / DESIGN-LLM-SCHEMA / DESIGN-LLM-SCHEMA-README | none on the non-defer reader surface | no repair needed |

**Papers `related:` frontmatter check:** `RELEASE-LEVEL-CLASSIFICATION.md` was NOT cited from any `docs/papers/*.md` `related:` block — the only inbound was the README nav row above. (It cited others *outward*, but nothing cited it inward except the README.)

---

## 3. Falsification gate — zero-dangling evidence

`git grep -n <basename> docs README.md` after deletion:

- `CORE-REFACTOR` → **ZERO**
- `PLAN-174-orchestrator-runtime` → **ZERO**
- `DESIGN-LLM-SCHEMA` → **ZERO** (covers both the doc and its README)
- `DESIGN-LLM-SCHEMA-README` → **ZERO**
- `RELEASE-LEVEL-CLASSIFICATION` → **ZERO**
- `DESIGN-227-distribution-pipeline` → 2 hits, **both inside `docs/reference/packages/DESIGN-266-dist-out-of-git.md`** (lines 76, 237)
- `SELF-COHERENCE-227` → 2 hits, **both inside `docs/reference/packages/DESIGN-266-dist-out-of-git.md`** (lines 76, 237)

The only surviving reader-surface hits are inside **DESIGN-266**, which is:
- a **HARD-DEFER do-not-touch** path in this round's mandate (code-coupled ← `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md`), and
- itself **queued for code-pass deletion** (γ §5 item 7), and
- the hits are **prose backtick-mentions** in "impact inventory" bullets ("…describe `cn build` producing dist/. **Unchanged**…"), **not** navigational `[..](..)` links or `related:` cites — they do not 404 and do not render as broken links.

→ Editing DESIGN-266 to drop these mentions would require touching a defer path to keep my deletions cite-clean, which the mandate forbids ("If an edit would need a defer path to stay link-correct, STOP and flag for the code pass"). **Flagged below.** When the code pass deletes DESIGN-266 alongside its SKILL.md coupling, these two mentions vanish with it. No non-defer dangling remains.

Out-of-scope hits (noted, not defects): `.cdd/**` (EXEMPT dotdir — review/audit records), `CHANGELOG.md` (deferred, code-pass/release-gate owned). No `src/**`/`schemas/**`/`scripts/**`/`tests/**` hits for any deleted basename.

---

## 4. Defer-path integrity

`git status --short` full change set:

```
 M docs/concepts/AGENT-NETWORK.md
 M docs/papers/README.md
D  docs/papers/RELEASE-LEVEL-CLASSIFICATION.md
D  docs/reference/packages/DESIGN-227-distribution-pipeline.md
D  docs/reference/packages/SELF-COHERENCE-227.md
D  docs/reference/runtime/CORE-REFACTOR.md
D  docs/reference/runtime/PLAN-174-orchestrator-runtime.md
D  docs/reference/schemas/DESIGN-LLM-SCHEMA-README.md
D  docs/reference/schemas/DESIGN-LLM-SCHEMA.md
```

7 deletions + 2 doc edits, **all under `docs/`**. No `src/**`, `schemas/**`, `scripts/**`, `tests/**`, `.github/**`, no `.cn-sigma/**`, no other `.cdd/**` (only this receipt), no `RELEASE.md`/`CHANGELOG.md`/`ROLES.md`/`OPERATOR.md`/`SUSTAINABILITY.md`, no `docs/development/board/`. None of the 6 code-coupled do-not-delete docs touched. **Defer set intact.**

---

## 5. Flags for the code pass

- **DESIGN-266 (already on the code-pass list, γ §5 item 7):** delete `docs/reference/packages/DESIGN-266-dist-out-of-git.md` + repair its coupling `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md`. When deleted, it also clears its own lines 76 & 237 prose mentions of the now-deleted `DESIGN-227-distribution-pipeline.md` and `SELF-COHERENCE-227.md` (this round left them because DESIGN-266 is a hard-defer path). No separate action needed beyond deleting DESIGN-266.
- No new code coupling discovered for the 7 files deleted this round (`git grep` over `src|schemas|scripts|tests` was clean for each basename).

---

## 6. Round-07 completion statement

The γ Round-6 convergence-hunt matter (7 docs-only records) is **executed in full**: all 7 confirmed historical (no misclassification, no KEEP), deleted; both surviving-surface inbound links repaired (AGENT-NETWORK Related bullet, papers/README nav row); zero non-defer dangling proven by grep; defer set untouched. The one remaining residual of this class on the reader surface is the code-coupled DESIGN-266, correctly left for the code pass. Matter ready for γ Round-7 verification.
