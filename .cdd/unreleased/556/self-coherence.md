---
artifact: self-coherence
cycle: 556
issue: https://github.com/usurobor/cnos/issues/556
branch: cycle/556
round: R0
author: alpha
completed:
  - Gap
  - Skills
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

## §Skills

**Tier 1:** `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` (this file's
own load order, §2.1–§2.8, §3.6).

**Tier 2 (always-applicable `eng/*`):** `src/packages/cnos.core/skills/write/SKILL.md`
(prose discipline for the new SKILL.md doctrine files and this artifact).

**Tier 3 (issue-specific, per γ's scaffold §8 α-prompt):**

- `docs/architecture/DESIGN-CONSTRAINTS.md` §3 (§3.1 command naming,
  §3.2 CLI-dispatch/domain-package boundary — the rule this cycle exists to
  honor and keep honoring through the move).
- `docs/reference/packages/PACKAGE-SYSTEM.md` §1.1 (content classes —
  `commands` class shape) and §7 (command discovery precedence — built-in
  always shadows; the reason `commands/issues-map/` must NOT be declared in
  `cn.package.json`'s `commands` object).
- `docs/reference/runtime/GO-KERNEL-COMMANDS.md` (bootstrap-kernel target
  set; confirms `issues-map` is not in the target set today, i.e. the
  built-in-shim disposition is intentional, #216-shaped, not this cycle's
  job).
- `docs/development/issues/TAXONOMY.md` and `docs/development/issues/TRIAGE.md`
  (cited, not forked, by the new `skills/taxonomy/SKILL.md` and
  `skills/triage/SKILL.md`).

**Not loaded (correctly out of scope):** β/γ role skills (α does not load
peer-role skills per §2.1 load-order rule); no `eng/{language}`-specific
skill beyond stdlib-Go conventions already established by the moved code
(no new language, no new dependency).
