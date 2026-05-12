---
cycle: 344
role: alpha
status: in-progress
---

# Self-Coherence — Cycle #344

## §Gap

**Issue:** #344 — `cdd: New skill cdd/activation/SKILL.md — bootstrap cdd in an existing repo (CI, notifications, secrets, identity)`

**Version / mode:** docs-only (design-and-build). Design converged in issue body. No code changes; all deliverables are Markdown skill files and cnos marker files.

**Branch:** `cycle/344` (created by γ at dispatch; α checks out, never creates).

**Gap being closed:** cdd ships skills for all lifecycle phases but has no canonical "how to turn cdd on in this repo" skill. New-tenant onboarding is implicit — every tenant re-derives the bootstrap sequence (`.cdd/` scaffold, version pin, labels, identity, CI, notifications, secrets). This issue closes that gap with a single authoritative `cdd/activation/SKILL.md` covering §1–§24.

**Cycle A scope (this cycle):** prose-only — skill authoring + cross-references + cnos self-activation marker files. Reference notifier impl (Cycle B) and tsc adoption (Cycle C) are out of scope.

**ACs in scope:** A.AC1 through A.AC6 (6 total). 37 open questions; all decided by α.

---

## §Skills

**Tier 1:**

- `src/packages/cnos.cdd/skills/cdd/CDD.md` — canonical lifecycle and role contract
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — α role surface (load order, algorithm, pre-review gate)

**Tier 2:**

- `src/packages/cnos.core/skills/write/SKILL.md` — writing standards (declared in dispatch Tier 3 bundles)
- `src/packages/cnos.core/skills/skill/SKILL.md` — skill authoring standards (declared in dispatch Tier 3 bundles)

**Tier 3 (issue-specific):**

- `src/packages/cnos.cdd/skills/cdd/operator/SKILL.md` — §5 dispatch configurations (§8 reference), §"Git identity for role actors" (§7 reference); cycle #342 + #343 context
- `src/packages/cnos.cdd/skills/cdd/release/SKILL.md` — §3.8 grading floor (referenced by dispatch declaration §8)
- `src/packages/cnos.cdd/skills/cdd/review/SKILL.md` — §3.13 honest-claim categories (§14), severity scale D/C/B/A (§22)
- `src/packages/cnos.cdd/skills/cdd/post-release/SKILL.md` — cdd-iteration flow; cross-reference target
- `src/packages/cnos.cdd/skills/cdd/epsilon/SKILL.md` — ε role; cdd-iteration.md work product (§22 reference); cycle #345 context
- `ROLES.md` (repo root) — §1 row 5: ε iterates the δ-discipline (§22 reference)

**Loaded artifacts (cross-cycle context):**

- cycle #339: `scripts/validate-release-gate.sh` — canonical CI artifact check example (§9 reference); confirmed to exist at `scripts/validate-release-gate.sh`
- cycle #342: `operator/SKILL.md §5` — dispatch configurations §5.1 / §5.2
- cycle #343: `operator/SKILL.md §"Git identity for role actors"` — `{role}@{project}.cdd.cnos` canonical form
- cycle #345: `epsilon/SKILL.md` — ε role definition; `ROLES.md §1` row 5

**Design artifact:** Not required — design converged in issue body (24 sections specified, 37 OQs with recommendations).

**Plan:** Not required — implementation sequencing is single-file plus marker files; issue body specifies exact section structure.
