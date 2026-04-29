---
name: missing-required-no-exception
description: Spec-required-but-exception-backed field missing without an exception.
governing_question: How do we prove the script rejects a missing required-or-excepted field?
triggers:
  - fixture
scope: task-local
---

# Missing required, not excepted

Negative fixture: `artifact_class`, `kata_surface`, `inputs`, `outputs` are
all missing AND there is no entry in `schemas/skill-exceptions.json`
covering this fixture's path. The script (not CUE) must fail this.
