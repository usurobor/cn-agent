---
name: minimal-fixture
description: Minimal positive fixture — hard-gate plus the four spec-required-but-exception-backed fields.
artifact_class: skill
kata_surface: embedded
governing_question: How do we prove the schema accepts a minimum-conformant skill?
triggers:
  - fixture
scope: task-local
inputs:
  - any input
outputs:
  - any output
---

# Minimal fixture

Used by `tools/validate-skill-frontmatter.sh --self-test` (#301 AC8) to
prove the schema accepts a skill with only the hard-gate fields.
