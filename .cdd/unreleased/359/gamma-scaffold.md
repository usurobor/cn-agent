---
cycle: 359
issue: "https://github.com/usurobor/cnos/issues/359"
date: "2026-05-14"
dispatch_configuration: "§5.2 single-session δ-as-γ"
base_sha: "23e28e45"
scope: "skill-patch: clarify §5.2 δ↔γ collapse scope"
---

# γ Scaffold — #359

## Gap

`operator/SKILL.md §5.2` says "single-session δ-as-γ via Agent tool" without specifying which role-pair collapses. Observed misread in tsc #49 wave-1: all roles fused into single subagent, violating §1.4 Triadic rule.

Confirmed via `rg "δ-as-γ\|delta-as-gamma\|δ=γ" src/packages/cnos.cdd/` — §5.2 exists in `operator/SKILL.md` but the collapse scope is ambiguous.

## α prompt

```
You are α. Project: cnos.
Load src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md and follow its load order.
Issue: gh issue view 359 --json title,body,state,comments
Branch: cycle/359
Tier 3 skills: src/packages/cnos.core/skills/write/SKILL.md

Git identity:
  git config user.name "alpha"
  git config user.email "alpha@cdd.cnos"

Scope: patch operator/SKILL.md §5.2 to explicitly state δ↔γ only collapse, γ↔α↔β stays separate per §1.4. 3 ACs — small diff.
```

## β prompt

```
You are β. Project: cnos.
Load src/packages/cnos.cdd/skills/cdd/beta/SKILL.md and follow its load order.
Issue: gh issue view 359 --json title,body,state,comments
Branch: cycle/359

Git identity:
  git config user.name "beta"
  git config user.email "beta@cdd.cnos"
```
