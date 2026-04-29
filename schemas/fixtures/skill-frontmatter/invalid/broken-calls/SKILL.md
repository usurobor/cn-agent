---
name: broken-calls
description: Static calls target points at a missing SKILL.md.
governing_question: How do we prove the script rejects a calls target that does not exist?
triggers:
  - fixture
scope: task-local
calls:
  - does-not-exist/SKILL.md
---

# Broken calls

Negative fixture: `calls` references a path that does not exist under
the package skill root. The script (not CUE) must fail this.
