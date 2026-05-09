---
cycle: 335
issue: "#335"
branch: "cycle/335-cdd-retro-closeout"
date: "2026-05-09"
---

# CDD Iteration — Cycle #335

This file lists the `cdd-*-gap` findings that cycle #335 addresses. This cycle is itself a `cdd-protocol-gap` triage event — its entire scope is remediation of two consecutive protocol skips.

---

### F1: protocol-skip pattern — cycles closing without close-out artifacts

- **Source:** γ-triage / issue #335 — post-merge audit of `cnos:main`
- **Class:** `cdd-protocol-gap`
- **Trigger:** §9.1 trigger: "loaded skill failed to prevent a finding" — the closure gate introduced by cycle #331 was not enforced for cycle #331 itself, nor for cycle #333 which followed immediately
- **Description:** Cycles #331 and #333 each produced zero close-out artifacts despite: (a) the γ/SKILL.md §2.10 closure gate requiring them, (b) cycle #331 having just introduced the `cdd-iteration.md` requirement (F6 in cycle #331's own cdd-iteration.md). The pattern repeated across two consecutive cycles, confirming a structural enforcement gap: the closure gate exists in the spec but no mechanical check prevents a merge without it.
- **Root cause:** The closure gate (`gamma/SKILL.md` §2.10) is γ-enforced but not operator-enforced. The δ role (operator/merge authority) can merge without verifying that `gamma-closeout.md` exists. The spec requires δ to check this (`gamma-closeout.md` is the closure declaration artifact per §5.3b), but no mechanical pre-merge CI gate enforces it. Two successive cycles demonstrate that the soft constraint is not sufficient.
- **Disposition:** `patch-landed`
- **Patch:** All 18 retroactive artifacts in `.cdd/releases/docs/2026-05-09/{331,333,335}/`, plus INDEX.md, LINEAGE.md, PRA, and CHANGELOG rows.
- **Affects:** `cdd/gamma/SKILL.md` §2.10 (closure gate enforcement)

### F2: closure gate needs mechanical reinforcement

- **Source:** F1 root cause analysis
- **Class:** `cdd-tooling-gap`
- **Trigger:** Recurring failure mode — same protocol-skip across two consecutive cycles
- **Description:** The closure gate at `gamma/SKILL.md` §2.10 is a soft constraint (γ checks a table of rows). No CI gate or pre-merge hook verifies that `.cdd/unreleased/{N}/gamma-closeout.md` exists before a cycle branch merges. δ merged both #332 and #333 without this check.
- **Root cause:** CDD's closure discipline relies entirely on role discipline, not mechanical enforcement. When cycles are urgent or treated as "just docs," the soft gate is skipped.
- **Disposition:** `next-MCA` — requires a cnos issue to add a mechanical pre-merge check (CI or operator-side script) that verifies closure-gate artifacts exist on cycle branches before merge.
- **Affects:** `cdd/gamma/SKILL.md` §2.10, `cdd/operator/SKILL.md` (δ pre-merge checklist)
