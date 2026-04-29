---
name: bad-enum-scope
description: Scope is not in the allowed enum.
governing_question: How do we prove the schema rejects an invalid scope?
triggers:
  - fixture
scope: bogus
---

# Bad enum

Negative fixture: `scope: bogus` violates the scope enum.
