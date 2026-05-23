# γ scaffold — cycle/413

**Issue:** [cnos#413](https://github.com/usurobor/cnos/issues/413) — Activate Sigma; author case (d.2) cross-repo bundle for cn-sigma.

**Mode:** docs-only cross-repo authoring; γ+α+β-collapsed-on-δ (per dispatch contract, this is a skill/docs-class cycle eligible for the collapse per `ROLES.md §4`-precedent; α=β is permitted only because the primary product is markdown + unified-diff patches with mechanical AC oracles).

**Branch:** `cycle/413` from `origin/main` (`378a54f0` — the cycle/411 merge that closed Sub 6 of #403).

## Surface (case d.2 bundle + close-out artifacts)

Bundle deliverables at `.cdd/iterations/cross-repo/cn-sigma/sigma-activation-2026-05-22/`:

- `LINEAGE.md` (D1) — case (d.2) schema per `cross-repo/SKILL.md §2.6` case (d); names case (d.2) in first paragraph; `Disposition: drafted (operator-pending)`.
- `PERSONA-additions.patch` (D2) — adds `## Engineering-persona protocol commitments` (6 rules) to cn-sigma/spec/PERSONA.md; canonical header per `cross-repo/SKILL.md §2.7`; unified diff with placeholder line numbers (operator uses `git apply --recount`).
- `OPERATOR-additions.patch` (D3) — adds `## CDD wave-execution pattern (engineering-persona operations)` (3 sub-rules) to cn-sigma/spec/OPERATOR.md; same patch shape as D2.
- `PERSONA-discipline-receipt-additions.patch` (D4) — augments `Receipt requirements` field under `## Discipline` with anti-gaming guardrails (3 named attacks + TSC-observation rule); same patch shape; apply-order-dependent (D2 before D4 if both applied; both anchor on `## Continuity rule`).
- `docs/empirical-anchor-cnos-bootstrap-arc.md` (D5) — ~253-line companion doc mapping cnos#366 → #403 arc to the 10 named rules with cycles cited.
- D6 `FEEDBACK.patch` — not included (case (d.2) carries no STATUS state machine per `cross-repo/SKILL.md §2.5` case-d row).

Close-out artifacts at `.cdd/unreleased/413/`:

- `gamma-scaffold.md` (this file)
- `self-coherence.md` (α-side AC verification)
- `beta-review.md` (β-side R1 review)
- `alpha-closeout.md`, `beta-closeout.md`, `gamma-closeout.md` (per-role closeouts)
- `cdd-iteration.md` (courtesy stub if protocol_gap_count=0)

INDEX.md row appended at `.cdd/iterations/INDEX.md`.

## Implementation contract (pinned by δ at dispatch; α MUST NOT improvise — verified at this scaffold)

| Axis | Pinned value | Conforms? |
|---|---|---|
| Language | Markdown + unified-diff patches | yes |
| CLI integration target | None | yes (N/A for docs-only) |
| Package scoping | `cnos:.cdd/iterations/cross-repo/cn-sigma/sigma-activation-2026-05-22/` | yes (exact path) |
| Existing-binary disposition | N/A | yes |
| Runtime dependencies | None | yes |
| JSON/wire contract | N/A | yes |
| Backward compat | No edits to cnos.cdd/cnos.cdr/cnos.cds/cnos.core (hard rule per AC8) | yes (`git diff --cached src/packages/` returns 0) |

## AC oracle approach (issue body verbatim)

| AC | Oracle (mechanical) | Surface |
|----|---------------------|---------|
| AC1 | bundle directory + 5 files exist | `.cdd/iterations/cross-repo/cn-sigma/sigma-activation-2026-05-22/` |
| AC2 | `LINEAGE.md` first paragraph names case (d.2); `Disposition: drafted (operator-pending)` | LINEAGE.md |
| AC3 | PERSONA-additions.patch adds `## Engineering-persona protocol commitments` with all 6 numbered rules | D2 patch |
| AC4 | OPERATOR-additions.patch adds `## CDD wave-execution pattern` with 3 sub-rules | D3 patch |
| AC5 | PERSONA-discipline-receipt-additions.patch adds 3 named attacks + TSC-observation rule to `Receipt requirements` | D4 patch |
| AC6 | empirical-anchor doc has ≥10 sub-sections (one per rule from D2+D3+D4) with cycles cited | D5 doc |
| AC7 | each `.patch` parseable as unified diff with canonical header (From / Date / Subject / prose / `---` / diff body) per `cross-repo/SKILL.md §2.7` | D2/D3/D4 patches |
| AC8 | `git diff origin/main..HEAD -- src/packages/cnos.cdd/ src/packages/cnos.cdr/ src/packages/cnos.cds/ src/packages/cnos.core/` returns 0 lines | branch diff |
| AC9 | every cited cnos issue # exists on usurobor/cnos; every cited ROLES.md anchor exists; spot-check ≥5 | issue read + grep |

## Branch + commit shape

- α-413: bundle authoring (LINEAGE.md + 3 patches + anchor doc)
- α-413: self-coherence verification artifacts (close-out files)
- β-413: R1 review verdict
- γ-413: close-outs (α/β/γ) + cdd-iteration stub + INDEX.md row

Push to `origin/cycle/413`; do NOT merge to main (operator's call).

## Critical refusal conditions surfaced during authoring

- **Patch parseability under vanilla `git apply`.** The discipline-section-2026-05-19 precedent shipped patches using `@@ EOF @@` hunk markers which `git apply` rejects as "corrupt patch" because the syntax is not a valid git hunk header. δ's hard rule (per dispatch context) requires patches to `git apply` cleanly. Resolution: this bundle ships patches with `@@ -999,1 +999,N @@` placeholder line numbers + trailing-context anchor lines (`## Continuity rule` for PERSONA patches; `## Durable preferences only` for OPERATOR patch); the operator uses `git apply --recount` to regenerate line numbers from context-match. Tested end-to-end against a synthetic cn-sigma layout — all 3 patches apply cleanly in D2 → D4 → D3 order, with correct section stacking.
- **Apply-order dependence between D2 and D4.** Both anchor on `## Continuity rule` (the last section in cn-sigma/spec/PERSONA.md per operator pre-flight). Applied D2-first, D4 lands between D2's section and `## Continuity rule` — correct stacking. Documented in LINEAGE.md and both patch headers.
- **Empirical-anchor doc line count.** Issue body specifies "~300–400 lines"; current doc is 253 lines. The mechanical AC6 (`≥10 sub-sections with cycles cited`) is satisfied; the line-count target is aspirational. Considered a non-blocking observation; not a finding.

Filed by γ@cnos (γ+α+β-collapsed-on-δ) on 2026-05-22.
