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
