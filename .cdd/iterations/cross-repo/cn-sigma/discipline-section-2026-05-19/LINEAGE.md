# Feedback patch — `cn-sigma:spec/PERSONA.md` discipline section

## Purpose

Add a `## Discipline` section to `cn-sigma/spec/PERSONA.md` per `ROLES.md §4a.1`. Sigma is the canonical engineering persona; this section encodes Sigma's loss function (action-biased, shippable artifact under repairable feedback) as a standing artifact.

## Source

- **Trigger:** `ROLES.md §4a.1` requires every persona hub to carry a discipline-profile section. Sigma's existing `spec/PERSONA.md` predates §4a.1 and is missing the section.
- **Reference doctrine:** `ROLES.md §4a` (five-layer enforcement chain), §4a.1 (discipline-profile fields), §4a.2 (engineering loss function)
- **Sibling reference:** `cnos:.cdd/iterations/cross-repo/cn-rho/bootstrap-2026-05-19/PERSONA.md` (the new Rho persona hub carries the same section shape for the research loss function)

## Bundle contents

- `PERSONA-discipline.patch` — unified diff appending the section to cn-sigma's PERSONA.md (operator may relocate the block if a more natural insertion point exists)
- `LINEAGE.md` — this file

## Application

cnos session cannot write to `usurobor/cn-sigma` (MCP scoped to `usurobor/cnos`). Operator applies:

```
cd usurobor/cn-sigma
git apply <path-to-PERSONA-discipline.patch>
git commit -m "persona: add discipline section per ROLES.md §4a.1"
git push
```

## Verification after application

- `rg "^## Discipline" spec/PERSONA.md` returns one match.
- The six fields (Primary virtue, Primary error, Default tempo, Claim/artifact boundary, Refusal conditions, Receipt requirements) all appear under the section.
- Cross-references to `ROLES.md §4a.1`, `ROLES.md §4a.2`, and `cn-sigma/spec/OPERATOR.md` are present.

## Disposition

Status: **drafted** (operator-pending — patch application is the gate).

Filed by γ on 2026-05-19.
