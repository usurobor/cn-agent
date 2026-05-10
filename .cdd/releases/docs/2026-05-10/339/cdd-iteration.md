---
cycle: 339
issue: "#339"
date: "2026-05-10"
---

# CDD Iteration — Cycle #339

This file records the `cdd-*-gap` findings that cycle #339 remediates. Both findings originated in cycle #335 (`beta-review.md` Round 2 and `cdd-iteration.md`). This cycle is the declared patch cycle for both.

Structure per `post-release/SKILL.md` Step 5.6b.

---

### F1: mechanical pre-merge closure-gate enforcement

- **Source:** `.cdd/releases/docs/2026-05-09/335/cdd-iteration.md` F2 ("cdd-tooling-gap: mechanical closure-gate reinforcement") — `next-MCA` disposition in cycle #335 pointing to this cycle as the patch
- **Class:** `cdd-tooling-gap`
- **Trigger:** Recurring failure mode — protocol-skip across three consecutive cycles (#331, #333, #334) each merging without close-out artifacts; soft gate at `gamma/SKILL.md` §2.10 demonstrably insufficient
- **Description:** `scripts/validate-release-gate.sh` was invoked only from `scripts/release.sh` (pre-tag, line 55) and only in `release` mode. No mechanical gate validated closure artifacts at merge time. δ could merge any cycle branch without `alpha-closeout.md`, `beta-closeout.md`, or `gamma-closeout.md` present. Three consecutive cycles demonstrated the soft constraint failed.
- **Root cause:** CDD's closure discipline relied entirely on role discipline (`gamma/SKILL.md` §2.10 table checked by γ in prose). No CI gate or pre-merge hook enforced the constraint. When cycles were urgent or treated as "docs-only," the soft gate was skipped.
- **Disposition:** `patch-landed`
- **Patch:** `ed982f6b` — extends `scripts/validate-release-gate.sh` with `--mode pre-merge`; validates `alpha-closeout.md`, `beta-closeout.md`, `gamma-closeout.md` for each triadic cycle dir (identified by `beta-review.md` presence per CDD.md §1.2); diagnostic names cycle number + filename; exits non-zero when any is absent.
- **Affects:** `scripts/validate-release-gate.sh` (extended); `src/packages/cnos.cdd/skills/cdd/release/SKILL.md` §3.8 (cross-references gate as enforcement surface); `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` §2.10 (soft gate now backed by mechanical check)

---

### F2: rubric closure-gate-failure handling and letter normalization

- **Source:** `.cdd/releases/docs/2026-05-09/335/beta-review.md` Round 2 §F7 ("honest-claim 3.13b — source-of-truth alignment across artifacts") — observational finding, operator disposition "deserves a small follow-on cnos issue on rubric closure-gate-failure handling"
- **Class:** `cdd-protocol-gap`
- **Trigger:** `release/SKILL.md` §3.8 rubric had no clause specifying that closure-gate failure forces `<C` regardless of per-axis math; `<C` (rubric label) was never reconciled with `C−` (used in CHANGELOG/PRA/alpha-closeout for cycles #331/#333/#334), creating silent dual-label drift
- **Description:** The §3.8 grading rubric mapped grades A → `<C` with per-axis numeric values and a geometric-mean formula. A cycle with no close-out artifacts could technically produce a C+ geometric mean if per-axis grades averaged out. The rubric provided no force-`<C` path for closure-gate failure. Additionally, two labels described the same disposition: the rubric used `<C` while operational artifacts used `C−`.
- **Root cause:** §3.8 was introduced by cycle #331 patch 5 (commit `b27fc15`) to make grading reproducible. The clause making closure-gate failure a hard floor was not included in that patch. The `<C` / `C−` split arose because prose artifacts drifted from the rubric label while the rubric label was not updated to acknowledge the equivalence.
- **Disposition:** `patch-landed`
- **Patch:** `0b56ff86` — amends `release/SKILL.md` §3.8 with (1) closure-gate override clause placed before geometric-mean instruction, declaring `C_Σ` forced to `<C` when any required close-out is absent at merge time; (2) letter normalization paragraph declaring `C−` as the operator-visible projection of `<C`, making both forms explicitly equivalent.
- **Affects:** `src/packages/cnos.cdd/skills/cdd/release/SKILL.md` §3.8 (amended)
