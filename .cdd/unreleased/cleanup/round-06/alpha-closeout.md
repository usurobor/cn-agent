# α Closeout — Round 06 (implementer: clear γ's Round-5 filed residuals)

**Role:** α (implementer). **Branch:** `claude/repo-cleanup-newcomer`. **HEAD at start:** `e3d33e4f` (clean tree, only the untracked Round-5 γ receipt present).
**Mandate:** de-historicize the **non-dot** reader surface to current-state only; git history is the archive. Dotdirs (`.cdd/`, `.cn-sigma/`, `.github/`) EXEMPT. Contract = γ Round-5 `gamma-closeout.md` §2 + §5 (the 6 clean-deletion groups, 16 files).
**Scope:** no commit, no PR. Edited only `docs/**`, `README.md`, and this receipt.

---

## 1. Deletions executed (16 files, all confirmed present via `git ls-files` before `git rm`)

**1 — unbuilt-future spec in the canonical current-spec dir (1 file):**
- `docs/reference/packages/PACKAGE-RESTRUCTURING.md` — `Status: Draft proposal — not the shipped structure`; proposed retiring cnos.pm + adding cnos.kernel/agent/hub, none landed.

**2 — doctrine production-record cluster (11 files), under `docs/concepts/doctrine/*/`:**
- `coherence-for-agents/CFA-cycle-log.md`, `CFA-critiques.md`
- `ethics-for-agents/EFA-cycle-log.md`, `EFA-critiques.md`, `EFA-external-observations.md`
- `inheritance-for-agents/IFA-cycle-log-dyad.md`, `IFA-cycle-log-gamma.md`, `IFA-critiques.md`
- `judgment-for-agents/JFA-cycle-log-dyad.md`, `JFA-cycle-log-gamma.md`, `JFA-critiques.md`
- The four `*-FOR-AGENTS.md` essays were KEPT (current-state doctrine).

**3–6 — dated point-in-time / plan-class docs under `docs/development/` (4 files):**
- `docs/development/cdd/ISSUE-CONSOLIDATION-ANALYSIS.md` — v3.58.0 backlog-triage snapshot (repo now 3.82.0).
- `docs/development/cdd/CDD-PACKAGE-AUDIT.md` — `**Date:** 2026-04-24` point-in-time "Converged Report" audit.
- `docs/development/design/ccnf-o-track-a1-survey.md` — `**Date:** 2026-05-23` in-flight track-dispatch survey (see §2).
- `docs/development/ENGINEERING-LANE-CLARITY.md` — `Status: Draft` design amendment whose targets already landed.

`git status` confirms exactly **16 D** entries. No code coupling found for any of the 16 (final grep over `src/**`/`schemas/**`/`scripts/**` returned zero; the only inbound hits anywhere were in exempt `.cdd/**` review records using dead `docs/gamma/`/`docs/alpha/` paths).

## 2. ccnf-o-track-a1-survey pinned-decision finding → **nothing live to migrate**

The survey pins five decisions (name = CCNF-O; 20-surface inventory; 6 higher-level forms classified universal; TSC-integration v0.1 parallelism; package = `cnos.ccnf-o`). All five govern an **unbuilt future package**. Confirmed absent from the live tree:
- `git ls-files 'src/packages/cnos.ccnf-o/**' 'schemas/ccnf-o/**'` → empty.
- No surviving doc references `cnos.ccnf-o` as current state (`git grep cnos.ccnf-o` over docs/src/schemas → 0 outside the deleted survey).

The survey states of itself: "It is not an implementation; it authors no schemas, no Go types, no validator changes." Its decisions pin shape for Tracks A2–A6/B1, which have not dispatched. There is therefore no canonical current-state doc that *would* own them — the thing they describe does not exist. Migrating them into `docs/reference/` or `docs/architecture/` would re-introduce unbuilt-future into the canonical dir, i.e. re-create the exact class being deleted (`PACKAGE-RESTRUCTURING.md`). Conclusion: **no live pinned decision to migrate**; git history archives the survey. If/when a track dispatches, its decision is re-derived into the actual package home.

## 3. Optional trims — both done (clean)

- **`docs/reference/conventions/AGENT-ACTIVATION-LOG-v0.md` §9** — reframed the `**Generalization:**` bullet, dropping the supersession/version-name narration ("the Sigma-specific `SIGMA-ACTIVATION-LOG-v0.md` name is superseded by…") to a current-state `**Scope:**` bullet ("covers any agent identity with a home hub and foreign activations; mechanics are the same for every adopter"). Convention itself unchanged. (The `cn-sigma:` field-writeup cites on lines 9/194/195 are dotdir refs, EXEMPT — left intact.)
- **Bare `**Date:**` frontmatter stamps** dropped from `docs/architecture/ARCHITECTURE.md`, `docs/architecture/cognitive-substrate/CAR.md`, `docs/architecture/security/SECURITY-MODEL.md`, `docs/concepts/AGENT-NETWORK.md`. Only the bare Date line removed in each; `Status`/`Supersedes`/`Author` kept. The OCaml-archival current-state facts in ARCHITECTURE/DESIGN-CONSTRAINTS were NOT touched. (AGENT-NETWORK: re-inserted the blank line before the `---` rule so `**Status:** Vision` does not render as a setext heading.)

## 4. Inbound-link repairs (ref → fix)

| Surviving file | Inbound ref to deleted target | Fix |
|---|---|---|
| `docs/papers/CELL-OF-CELLS.md:16` | `related:` YAML cite → `docs/development/design/ccnf-o-track-a1-survey.md` | **removed** that `related:` line (the known trap). YAML block re-verified valid. |
| `docs/development/design/README.md:17` | Document-Map row → `ccnf-o-track-a1-survey.md` (its only row) | replaced the now-empty survey+decision map with a current-state note: past surveys are archived in git history; a survey's pins land in the owning package/schema home when a downstream cycle dispatches. (`cn-repo-install-MOCKS.md` remains in the dir — different class, deferred to code pass — so the survey+decision map is legitimately empty.) |
| `docs/concepts/doctrine/inheritance-for-agents/INHERITANCE-FOR-AGENTS.md:15–17` | `Related:` links (GitHub URLs) → CFA/EFA/JFA cycle-logs + critiques + external-observations (the deleted cluster) | **removed** those three bullets; kept the three surviving essay-to-essay links (lines 12–14). Essay prose is self-contained (it narrates the cycle findings inline), so no argument breaks. |
| `docs/concepts/doctrine/README.md` §"Cycle artifacts" (lines 25–34) | per-folder prose + "read the cycle logs alongside the essays" describing the deleted cluster | rewrote the section to current state: the essays carry the doctrine; the named failure modes are inherited in the next section (`[Inherited failure modes](#inherited-failure-modes)` — anchor verified present); the per-cycle production records are archived in git history. No doctrine lost — the distilled failure modes were already enumerated inline in the "Inherited failure modes" section. |

**Essay back-link check (all four):** grepped `cycle-log|critiques|external-observations` in each `*-FOR-AGENTS.md`. Only `INHERITANCE-FOR-AGENTS.md` back-linked the deleted siblings (repaired above). `COHERENCE-`, `ETHICS-`, `JUDGMENT-FOR-AGENTS.md` had **no** back-links — untouched.

**Not a dangling ref (noted):** `ISSUE-CONSOLIDATION-ANALYSIS.md:452` carried a self-citation to its own old `docs/gamma/cdd/…` path — internal to the deleted file, disappears with it. `docs/development/cdd/README.md` does **not** list either deleted cdd doc in its Document Map — no dangle there.

## 5. Zero-dangling evidence (the falsifying grep)

```
git grep -nE "PACKAGE-RESTRUCTURING|CFA-cycle-log|CFA-critiques|EFA-cycle-log|EFA-critiques|\
EFA-external-observations|IFA-cycle-log|IFA-critiques|JFA-cycle-log|JFA-critiques|\
ISSUE-CONSOLIDATION-ANALYSIS|CDD-PACKAGE-AUDIT|ccnf-o-track-a1-survey|ENGINEERING-LANE-CLARITY" \
  -- 'docs/**' 'README.md'
→ (no output) ZERO dangling
```

Out-of-scope hits (noted, not chased, all EXEMPT): only `.cdd/**` review records (e.g. `.cdd/releases/3.61.0/268/*`, `.cdd/releases/3.82.0/421/*`, `.cdd/unreleased/508/pass4-classification.md`) referencing the deleted files via dead `docs/gamma/`/`docs/alpha/` paths. No hit in `src/**`, `schemas/**`, `scripts/**`, or `CHANGELOG.md` for any of the 16 — confirming all clean docs-only deletions.

## 6. Defer-path integrity

`git status` over `src/ schemas/ scripts/ tests/ .github/ .cn-sigma/ install.sh cn.json` (excluding this round's receipt dir) → **no defer path touched**. The five code-coupled/code-pass docs on the do-not-delete list (`ENGINEERING-LEVEL-ASSESSMENT.md`, `CTB-v4.0.0-VISION.md`, `LANGUAGE-SPEC-v0.2-draft.md`, `DISPATCH-FAILURE-EVIDENCE.md`, `cn-repo-install-MOCKS.md`) are all still present, untouched. No edit required changing a defer path to stay link-correct.

## 7. Flagged for the code pass

- **`docs/development/design/` directory lifecycle.** After this round it holds `README.md` + the deferred `cn-repo-install-MOCKS.md`. When the code pass deletes MOCKS (per γ §5), the dir will hold only its README with an empty survey+decision map — the code pass should decide whether to remove the directory (+ README) entirely at that point. Flagged, not acted on (MOCKS is a defer path).

## 8. Changeset summary

- **16 deletions** (D), **9 doc edits** (M): `ARCHITECTURE.md`, `CAR.md`, `SECURITY-MODEL.md`, `AGENT-NETWORK.md` (Date trims), `doctrine/README.md` + `INHERITANCE-FOR-AGENTS.md` (cluster link repairs), `design/README.md` + `CELL-OF-CELLS.md` (ccnf-survey link repairs), `AGENT-ACTIVATION-LOG-v0.md` (supersession trim), plus this receipt add.
- Four doctrine essays: three untouched, IFA edited only to drop dead back-links. All defer paths untouched. Zero dangling in the reader surface.
