---
name: issues/map
description: Issue-board (Issue Pivot) generation domain. Use when reasoning about `cn issues map`'s fetch, taxonomy-parse, effort-computation, and render pipeline, or about the board's JSON/HTML output contract.
artifact_class: skill
kata_surface: none
governing_question: How does an open-issue set become a taxonomy-aware, effort-weighted board record set and a self-contained interactive board output?
visibility: public
triggers:
  - issues map
  - issue board
  - issue pivot
  - board-data.json
scope: task-local
inputs:
  - open GitHub issues (live via API, or fixture/stdin for offline use)
outputs:
  - board records (taxonomy, effort, effort_weight, unestimated)
  - index.html (Voronoi + Pivot views)
  - board-data.json
  - README.md (generated, do-not-hand-edit)
---

# Issue map

## Core Principle

The board is a **projection of the issue taxonomy**, not an independent data
model: every field a board record carries is derived from taxonomy labels
(`skills/taxonomy/SKILL.md`), never invented locally.

## Implementation location

The domain implementation lives at
[`src/go/internal/issuesmap/`](../../../../go/internal/issuesmap/) — **not**
under this package's directory. Per the operator's explicit "Go
implementation rule" on [cnos#556](https://github.com/usurobor/cnos/issues/556)
("Do not force Go implementation code into `src/packages/`. The active Go
implementation may remain under `src/go/internal/issuesmap/` during the
shim phase, but it must be treated as the implementation of the
`cnos.issues` domain."), this location is deliberate and pinned, not a
migration-in-progress. It is dispatched by the thin
`src/go/internal/cli/cmd_issues_map.go` (argument parsing + one-line
delegation only — no domain logic).

See `src/packages/cnos.issues/SKILL.md` §"Command-dispatch disposition" for
the full rationale, including the R1 reversion of an R0 attempt that
physically relocated this source under
`src/packages/cnos.issues/commands/issues-map/` (mirroring the
[cnos#392](https://github.com/usurobor/cnos/issues/392) `cdd-verify`
precedent) — that relocation contradicted the operator's instruction above
and was reverted.

## Command contract

```text
cn issues map --repo <owner/repo> --out docs/development/board
```

- fetches open issues (GitHub API, or `--fixture`/`--stdin` for offline use);
- parses `kind/*`, `area/*`, `P0`–`P3`, `status:*`, `dispatch:cell`,
  `protocol:*`, `effort/*` per the taxonomy;
- computes `effort`, `effort_weight` (S=1, M=2, L=4, XL=8), and
  `unestimated` (missing effort ≠ small — nominal weight 1, visually flagged);
- excludes `kind/tracking` from effort rollups by default;
- emits `index.html`, `board-data.json`, `README.md` into `--out`.

## Rule

- The CLI dispatcher never contains taxonomy, effort, or rendering logic —
  that is a dispatch-boundary violation (`docs/architecture/DESIGN-CONSTRAINTS.md` §3.2).
- Board output behavior (fields, weights, exclusions) is unchanged by moving
  the implementation's package home; a package-boundary refactor is not a
  license to change board semantics.
