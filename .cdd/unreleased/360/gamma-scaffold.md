---
cycle: 360
issue: "https://github.com/usurobor/cnos/issues/360"
date: "2026-05-14"
dispatch_configuration: "§5.2 single-session δ-as-γ"
base_sha: "c77f34a4"
scope: "skill-patch: review rule 3.11b exemption scope clarification"
---

# γ Scaffold — #360

## Gap

`review/SKILL.md` rule 3.11b exemption "documented in the issue" is ambiguous — does master-issue comment suffice, or sub-issue body only? Four β subagents diverged on this in tsc #49.

Confirmed via `rg "3\.11b" src/packages/cnos.cdd/` — rule exists in review/SKILL.md, exemption scope not specified.

## α prompt

```
You are α. Project: cnos.
Load src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md and follow its load order.
Issue: gh issue view 360 --json title,body,state,comments
Branch: cycle/360
Tier 3 skills: src/packages/cnos.core/skills/write/SKILL.md

Git identity:
  git config user.name "alpha"
  git config user.email "alpha@cdd.cnos"

Scope: patch review/SKILL.md rule 3.11b to clarify exemption must be in sub-issue body, not parent comment. 3 ACs — small diff.
```

## β prompt

```
You are β. Project: cnos.
Load src/packages/cnos.cdd/skills/cdd/beta/SKILL.md and follow its load order.
Issue: gh issue view 360 --json title,body,state,comments
Branch: cycle/360

Git identity:
  git config user.name "beta"
  git config user.email "beta@cdd.cnos"
```
