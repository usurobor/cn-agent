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
`issuesmap_test.go`, `templates/`, `testdata/` — lives under this package at
[`commands/issues-map/`](commands/issues-map/), as its own Go module
(`module github.com/usurobor/cnos/packages/cnos.issues/commands/issues-map`),
mirroring the [cnos#392](https://github.com/usurobor/cnos/issues/392)
`cdd-verify` precedent exactly: wired into `src/go/go.mod` via a
`require`+`replace` pair, added to `go.work`'s `use (...)` list, and
dispatched by the thin `cmd_issues_map.go` shim (import + one-line
delegation only — no domain logic in the CLI file).

`cn issues map` is a **temporary built-in shim until
[#216](https://github.com/usurobor/cnos/issues/216)**: dispatch stays through
the compiled-in kernel registration, not the package-command exec-dispatch
mechanism ([`PACKAGE-SYSTEM.md` §7](../../../docs/reference/packages/PACKAGE-SYSTEM.md)).
`commands/issues-map/` under this package is **Go-source co-location**, not
a declared package command: `cn.package.json`'s `commands` object does not
list `issues-map` (matching `cnos.cdd/cn.package.json`, which likewise has
no `cdd-verify` entry despite `cdd-verify`'s Go source living under
`cnos.cdd/commands/cdd-verify/`). Declaring one would create a dead entry —
the kernel built-in always shadows a same-named package command silently
per the command-discovery precedence rule — and would misrepresent the
command as package-provided exec-dispatch when it is not.

A future #216 cycle may move `issues-map` onto real package-command
exec-dispatch: removing the kernel registration in `main.go`, adding a
`commands.issues-map` entry to this package's `cn.package.json` with a
`cn-issues-map` entrypoint script, and adding package-install plumbing to
the board-map GitHub Action. None of that is done by this package as it
stands — physical co-location of the domain source is not the same as
declared package-command dispatch.

## Non-goals

- Does not implement generic package-command discovery (out of scope; see #216).
- Does not fork the issue taxonomy or triage doctrine — it cites the
  canonical docs.
- Does not change board visual semantics, effort/priority/taxonomy
  semantics, or dispatch/wake/CDD/label behavior.
