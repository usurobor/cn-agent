---
cycle: 362
issue: "https://github.com/usurobor/cnos/issues/362"
date: "2026-05-14"
dispatch_configuration: "§5.2 single-session δ-as-γ"
base_sha: "cfe143f7"
scope: "skill-patch: add UIE communication gate to cap/SKILL.md"
---

# γ Scaffold — #362

## Gap

`cap/SKILL.md` §UIE has no communication gate distinguishing questions from instructions. Agent treats "what is X?" as "fix X" — understanding is invisible to operator, action fires before model verification.

Confirmed via `rg "question|answer|communication" src/packages/cnos.core/skills/agent/cap/SKILL.md` — no question-handling gate exists.

## α prompt

```
You are α. Project: cnos.
Load src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md and follow its load order.
Issue: gh issue view 362 --json title,body,state,comments
Branch: cycle/362
Tier 3 skills: src/packages/cnos.core/skills/write/SKILL.md

Git identity:
  git config user.name "alpha"
  git config user.email "alpha@cdd.cnos"

Scope: add UIE communication gate to cap/SKILL.md — question vs instruction distinction, answer-before-act rule, named failure mode. 3 ACs — small diff to one file.
```

## β prompt

```
You are β. Project: cnos.
Load src/packages/cnos.cdd/skills/cdd/beta/SKILL.md and follow its load order.
Issue: gh issue view 362 --json title,body,state,comments
Branch: cycle/362

Git identity:
  git config user.name "beta"
  git config user.email "beta@cdd.cnos"
```
