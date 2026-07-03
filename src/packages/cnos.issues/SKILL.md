---
name: issues
description: cnos.issues package entrypoint. Use when reasoning about issue taxonomy, the Issue Pivot / issue-board (`cn issues map`), or the current command-dispatch disposition of issue-board tooling.
artifact_class: skill
kata_surface: none
governing_question: What owns issue taxonomy, board mapping, and Issue Pivot generation, and how is that ownership realized through the current command-dispatch boundary?
visibility: public
triggers:
  - issues
  - issue board
  - issue pivot
  - board map
  - taxonomy
  - triage
scope: global
inputs:
  - open GitHub issues (labels, state, comments)
  - the issue taxonomy (kind/area/priority/effort/status/protocol/dispatch labels)
outputs:
  - issue-board records (taxonomy-aware, effort-weighted)
  - static Issue Pivot board output (index.html, board-data.json, README.md)
calls:
  - skills/map/SKILL.md
  - skills/taxonomy/SKILL.md
  - skills/triage/SKILL.md
---

# cnos.issues

## Core Principle

`cnos.issues` is the package home for issue-board cognition: it owns issue
taxonomy, board mapping, and Issue Pivot generation as a domain, distinct
from the CLI dispatch that invokes it.

`cnos.issues` owns:

- **issue taxonomy** — the label vocabulary an issue is judged against
  (`skills/taxonomy/SKILL.md`, citing [`docs/development/issues/TAXONOMY.md`](../../../docs/development/issues/TAXONOMY.md) as the canonical source, not forking it);
- **triage doctrine** — how a raw issue becomes coherently labeled and
  dispatchable (`skills/triage/SKILL.md`, citing [`docs/development/issues/TRIAGE.md`](../../../docs/development/issues/TRIAGE.md));
- **board mapping / Issue Pivot generation** — the domain logic that fetches
  open issues, derives taxonomy-aware records, computes effort, and renders
  the interactive board (`skills/map/SKILL.md`).

## Command-dispatch disposition

The user-facing command `cn issues map` is a **compiled-in kernel command**
(`SourceKernel` / `TierKernel`), registered in `src/go/cmd/cn/main.go` and
dispatched by the thin `src/go/internal/cli/cmd_issues_map.go` (argument
parsing + delegation only).

The Go domain implementation — `issuesmap.go`, `fetch.go`,
`issuesmap_test.go`, `templates/`, `testdata/` — lives at
[`src/go/internal/issuesmap/`](../../go/internal/issuesmap/), not under this
package. That path **is** the implementation of the `cnos.issues` domain —
own it conceptually as such, even though it is not physically nested here.

`cn issues map` is a **temporary built-in shim until
[#216](https://github.com/usurobor/cnos/issues/216)**: dispatch stays through
the compiled-in kernel registration, not the package-command exec-dispatch
mechanism ([`PACKAGE-SYSTEM.md` §7](../../../docs/reference/packages/PACKAGE-SYSTEM.md)).
Accordingly, `cn.package.json` declares **no** `commands` entry and there is
no `commands/issues-map/` directory — declaring one would misrepresent the
command's dispatch path and the code's location.

A future #216 cycle may move `issues-map` onto real package-command
exec-dispatch (placing the domain source under this package, adding a
`commands.issues-map` entry, and package-install plumbing to the board
Action).

## Non-goals

- Does not implement generic package-command discovery (out of scope; see #216).
- Does not fork the issue taxonomy or triage doctrine — it cites the
  canonical docs.
- Does not change board visual semantics, effort/priority/taxonomy
  semantics, or dispatch/wake/CDD/label behavior.
