---
name: full-fixture
description: Full positive fixture — every reserved frontmatter field used at least once.
artifact_class: skill
kata_surface: embedded
governing_question: How do we prove the schema accepts every reserved field together?
visibility: public
triggers:
  - fixture
  - full
scope: global
inputs:
  - input one
  - input two
outputs:
  - output one
requires:
  - prereq one
calls:
  - sub/SKILL.md
calls_dynamic:
  - source: issue.tier3_skills
    constraint: "skills under eng/"
runs_after:
  - other-skill
runs_before:
  - later-skill
excludes:
  - mutually-exclusive-skill
parent: full-fixture
---

# Full fixture

Used by `tools/validate-skill-frontmatter.sh --self-test` (#301 AC8) to
prove the schema accepts every reserved field, an unknown package-local
key (`parent:`), and a `calls` target whose existence is checked relative
to the package skill root.
