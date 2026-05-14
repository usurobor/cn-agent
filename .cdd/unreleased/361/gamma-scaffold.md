---
cycle: 361
issue: "https://github.com/usurobor/cnos/issues/361"
date: "2026-05-14"
dispatch_configuration: "§5.2 single-session δ-as-γ"
base_sha: "56202534"
scope: "skill-patch: review rules 3.3/3.4 verdict-shape lint"
---

# γ Scaffold — #361

## Gap

β@S4 in tsc #53 issued "APPROVED with 3 unresolved C findings + conditional language." Rules 3.3 and 3.4 both ban this shape, but the ban is implicit. No explicit verdict-shape lint exists.

Confirmed via `rg "3\.3|3\.4" src/packages/cnos.cdd/skills/cdd/review/SKILL.md` — rules exist but no explicit invalid-verdict-shape enumeration.

## α prompt

```
You are α. Project: cnos.
Load src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md and follow its load order.
Issue: gh issue view 361 --json title,body,state,comments
Branch: cycle/361
Tier 3 skills: src/packages/cnos.core/skills/write/SKILL.md

Git identity:
  git config user.name "alpha"
  git config user.email "alpha@cdd.cnos"

Scope: add verdict-shape lint rules to review/SKILL.md near 3.3/3.4. 3 ACs — small diff.
```

## β prompt

```
You are β. Project: cnos.
Load src/packages/cnos.cdd/skills/cdd/beta/SKILL.md and follow its load order.
Issue: gh issue view 361 --json title,body,state,comments
Branch: cycle/361

Git identity:
  git config user.name "beta"
  git config user.email "beta@cdd.cnos"
```
