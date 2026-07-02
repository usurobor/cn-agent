---
artifact: self-coherence
cycle: 556
issue: https://github.com/usurobor/cnos/issues/556
branch: cycle/556
round: R0
author: alpha
completed:
  - Gap
---

# α self-coherence — cnos#556

## §Gap

**Issue:** cnos#556 — "issues: move Issue Pivot into `cnos.issues` package boundary."

**Mode:** `design-and-build` (per the issue's own first line). This is an
architecture-boundary refactor (CLI-dispatch vs package-owned-domain-logic),
not a pure mechanical move and not a pure design memo.

**Gap:** issue-board domain logic (fetch, taxonomy parse, effort
computation, HTML+JSON render) lived in generic kernel-internal Go
(`src/go/internal/issuesmap/`) with no package-level skill/doctrine surface.
The command boundary was already thin (`cmd_issues_map.go` was already
argument-parsing + one-line delegation), but the domain implementation had
no package home, and there was no `cnos.issues` package stating what owns
issue taxonomy, board mapping, and Issue Pivot generation.

**Operator clarifying comment (2026-07-02)** narrows the issue's original,
more permissive "Proposed package shape" language: this cycle creates the
`cnos.issues` package home, keeps `cn issues map` as a temporary built-in
shim (Option B, not Option A), does not implement generic package-command
discovery, and does not solve #216.

**Base:** cycle branched from `origin/main` at `4fe8e4333b36372f595201841fb76cc0c31acff4`
(per γ's scaffold §0 base-SHA note — the wake-invoked-δ pinned SHA
`2ae24f27...` was 2 bot-only commits stale at scaffold time; γ correctly
branched from live `main` instead).

**Precedent mirrored:** cnos#392 (`cdd-verify` co-located at
`src/packages/cnos.cdd/commands/cdd-verify/` as its own Go module, wired via
`require`+`replace` in `src/go/go.mod` and `go.work`'s `use` list) — this
cycle applies the exact same shape to `issues-map`.

**Implementation contract (pinned by δ via γ's scaffold §8, not
renegotiated):** Go (unchanged), `cn` subcommand / compiled-in kernel
dispatch (unchanged), package scoping = Go-source co-location under
`src/packages/cnos.issues/commands/issues-map/` (not declared as a package
command), single `cn` binary preserved, no new runtime dependencies,
JSON/wire contract preserved as-is, backward-compat invariants preserved.
All seven axes were populated in the scaffold; none required escalation.
