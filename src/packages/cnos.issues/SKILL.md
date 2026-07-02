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

## Command-dispatch disposition (read this before assuming package-command discovery)

The user-facing command `cn issues map` is a **compiled-in kernel command**
(`Source: SourceKernel`, `Tier: TierKernel`), registered directly in
`src/go/cmd/cn/main.go` and dispatched by `src/go/internal/cli/cmd_issues_map.go`.

`cn issues map` remains a **temporary built-in shim until [#216](https://github.com/usurobor/cnos/issues/216) lands**.
Domain logic is isolated under this package's boundary
(`src/packages/cnos.issues/commands/issues-map/`, its own Go module,
mirroring the [cnos#392](https://github.com/usurobor/cnos/issues/392)
`cdd-verify` precedent), but dispatch stays through the compiled-in kernel
registration, not through the package-command exec-dispatch mechanism
([`PACKAGE-SYSTEM.md` §7](../../../docs/reference/packages/PACKAGE-SYSTEM.md)).

Concretely, `commands/issues-map/` under this package is **Go-source
co-location**, not a declared package command: `cn.package.json`'s
`commands` object does not list `issues-map`. Declaring it there would
create a dead entry — the kernel built-in always shadows a same-named
package command silently per the command-discovery precedence rule — and
would misrepresent the command as package-provided when it is not. This
package is honest about that distinction rather than papering over it.

A future [#216](https://github.com/usurobor/cnos/issues/216) cycle that
flips `issues-map` onto real package-command exec-dispatch would need to (a)
remove the kernel registration in `main.go`, (b) add a `commands.issues-map`
entry to this file's `cn.package.json` with a `cn-issues-map` entrypoint
script, and (c) add package-install plumbing to the board-map GitHub Action.
None of that is done by this package as it stands — this cycle
([cnos#556](https://github.com/usurobor/cnos/issues/556)) only relocates
the domain implementation to this package boundary.

## Non-goals

- Does not implement generic package-command discovery (out of scope; see #216).
- Does not fork the issue taxonomy or triage doctrine — it cites the
  canonical docs.
- Does not change board visual semantics, effort/priority/taxonomy
  semantics, or dispatch/wake/CDD/label behavior.
