---
artifact: gamma-scaffold
cycle: 556
issue: https://github.com/usurobor/cnos/issues/556
branch: cycle/556
base_sha: 4fe8e4333b36372f595201841fb76cc0c31acff4
round: R0
author: gamma
---

# γ R0 Scaffold — cnos#556

## 0. Base-SHA note (read first)

The wake-invoked-δ input contract (`cnos.cdd/skills/cdd/delta/SKILL.md` §9.2 input #3) pinned
`current main SHA = 2ae24f274eaa232097b9180899d4fe1eb9d83e2b` at claim time. At scaffold time
`origin/main` had advanced two commits to `4fe8e4333b36372f595201841fb76cc0c31acff4`. Diff
inspected (`git diff --stat 2ae24f2..origin/main`): both commits are `board-map: regenerate
docs/development/board from live open issues` (bot-authored, `board-map.yml`'s own scheduled/
event-triggered run), touching only `docs/development/board/board-data.json` and
`docs/development/board/index.html`. No canonical skill file, no `src/go/internal/issuesmap/`
file, no `src/packages/*/cn.package.json`, and no workflow file changed. This is benign,
expected drift (the exact generated-output surface this issue is about to relocate the *owner*
of, not the *content* of) — **not** a canonical-skill staleness event under
`gamma/SKILL.md`'s "Canonical-skill staleness check." `cycle/556` was branched from the live
`origin/main` HEAD (`4fe8e433...`) per `gamma/SKILL.md` §2.5's actual branch-pre-flight rule
("create `cycle/{N}` from `origin/main`"), not the stale pinned SHA. Both SHAs are recorded here
for audit; α should treat `4fe8e433...` as the cycle's base.

## 1. Issue

cnos#556 — "issues: move Issue Pivot into `cnos.issues` package boundary."
`gh issue view 556 --repo usurobor/cnos --json title,body,state,comments` is the source of
truth; this scaffold does not paraphrase it, only indexes it for α/β.

## 2. Mode

**design-and-build** (per the issue's own first line: `Mode: design-and-build`). This is an
architecture-boundary refactor (CLI-dispatch vs package-owned-domain-logic), not a pure
mechanical move and not a pure design memo — α both decides the concrete package shape (within
the guardrails below) and lands it as a working, CI-green diff.

## 3. Surfaces α is expected to touch

Confirmed by direct `Read`/`Grep`/`Bash` inspection of the current tree (not assumed):

| Surface | Current state (verified) | Expected α action |
|---|---|---|
| `src/go/internal/issuesmap/` (`issuesmap.go` 335 lines, `fetch.go` 68 lines, `issuesmap_test.go`, `templates/board.template.html`, `testdata/issues.json`) | Domain logic: fetch, taxonomy parse, effort computation, HTML+JSON render. `issuesmap.go` embeds `templates/board.template.html` via `//go:embed templates/board.template.html` (relative path — the whole subtree must move together, not be split). | **Relocate** the entire directory's contents (Go source + `templates/` + `testdata/`) to `src/packages/cnos.issues/commands/issues-map/`, as its own Go module (see §"Package scoping" in the Implementation Contract below and cnos#392 precedent in §7). |
| `src/go/internal/cli/cmd_issues_map.go` (31 lines) | Already thin: `Spec()` returns `CommandSpec{Name: "issues-map", Source: SourceKernel, Tier: TierKernel}`; `Run()` is a one-line delegation to `issuesmap.Run(...)`. Already compliant with `docs/architecture/DESIGN-CONSTRAINTS.md` §3.2 ("`cli/` owns dispatch. Domain packages own logic. No `cmd_*.go` file may contain domain logic beyond argument parsing and dispatch."). | **Preserve** the thin shape; only update the import path to the relocated module. Do not add domain logic here. This satisfies AC4 — it is already true, α's job is to keep it true through the move. |
| `src/go/cmd/cn/main.go:47` | `reg.Register(&cli.IssuesMapCmd{})` — the kernel command registration. | **No change.** `cn issues map` stays a compiled-in kernel command; this line is the evidence AC3 ("public command remains `cn issues map`") stays true. |
| `src/go/go.mod` | Single-module today for `internal/issuesmap`. cnos#392 precedent already lives here: a `require`/`replace` pair for `github.com/usurobor/cnos/packages/cnos.cdd/commands/cdd-verify => ../packages/cnos.cdd/commands/cdd-verify`. | **Add** an analogous `require`/`replace` pair for the new `cnos.issues` module. |
| `go.work` (repo root) | `use (./src/go, ./src/packages/cnos.cdd/commands/cdd-verify)`. | **Add** `./src/packages/cnos.issues/commands/issues-map` to the `use` list. |
| `src/packages/cnos.issues/` | **Does not exist yet** (confirmed: `find src/packages/cnos.issues -type f` returned nothing before this cycle). | **Create**: `cn.package.json`, `SKILL.md`, `skills/map/SKILL.md`, `skills/taxonomy/SKILL.md`, `skills/triage/SKILL.md`, `commands/issues-map/` (the relocated Go module — see guardrail on what this directory is and is not, §6). |
| `.github/workflows/board-map.yml` | Already calls the public command only: `./cn issues map --repo "$GITHUB_REPOSITORY" --out docs/development/board` (built via `go build -o "$GITHUB_WORKSPACE/cn" ./cmd/cn` in the same job — no `cn deps restore`, no package-install step). | **No change expected.** AC5 and AC6 are already satisfied by the *existing* workflow shape: it calls the public command and does not touch package internals, and it does not need install plumbing because the command stays a compiled-in kernel dispatch, not a package-command-content-class dispatch. α should verify this is still true after the move (the built `cn` binary just needs the new module to build cleanly) rather than editing the workflow. If α finds a reason the workflow *must* change, that is new information — escalate via `gamma-clarification.md` before editing it (§8.2 authoring-order discipline, `delta/SKILL.md` §8, is not triggered here since this is not a new remote-runner artifact, but the workflow is still a δ-class effect surface and unplanned edits should not be silent). |
| `docs/development/board/README.md:42` | States: `` The generator is the Go command `cn issues map` (`src/go/internal/issuesmap/`). `` | **Update** the path reference to the new location once the move lands. |
| `docs/architecture/DESIGN-CONSTRAINTS.md` §3 | §3.1 "Git-style subcommands" confirms `cn issues map` = noun `issues` + verb `map`, hyphenated registry key `issues-map` — already the shape in use. §3.2 is the dispatch-boundary rule this cycle exists to honor. | No change expected; cited as source of truth only. |

## 4. Per-AC oracle list

For each AC: what evidence α must produce, and what β must independently check.

### AC1 — `cnos.issues` package exists
- **α evidence:** `src/packages/cnos.issues/cn.package.json` committed, schema `cn.package.v1`, valid JSON, `"name": "cnos.issues"`. Run `cn build --check` (or the local equivalent — `./cn build --check` from a locally built binary) and paste the output (or the relevant excerpt) into `self-coherence.md`.
- **β check:** File exists, JSON parses, `cn build --check` (β re-runs it independently) passes and specifically recognizes `cnos.issues`.

### AC2 — Package owns issue-board cognition
- **α evidence:** `src/packages/cnos.issues/SKILL.md` (or equivalent top-level entrypoint) states explicitly that `cnos.issues` owns issue taxonomy, board mapping, and Issue Pivot generation — not just a name, an actual doctrine statement naming what it owns and where the implementation currently lives (built-in shim vs package-command).
- **β check:** Read the file; confirm it is not an empty wrapper (positive/negative pair from the issue: "it states what it owns" vs "empty wrapper with no domain contract").

### AC3 — Public command remains `cn issues map`
- **α evidence:** After the move, `cn issues map --repo usurobor/cnos --out <tmpdir>` still runs and regenerates board output (α runs this locally against a built `cn` binary and records the output/diff-vs-before in `self-coherence.md`).
- **β check:** β independently builds `cn` from the cycle branch and re-runs the same command; confirms it works without needing a package-internal path or script.

### AC4 — CLI dispatch remains thin
- **α evidence:** `git diff` on `src/go/internal/cli/cmd_issues_map.go` shows only an import-path change; no taxonomy/rendering/label logic added.
- **β check:** Read the diff; confirm no domain logic entered any `cmd_*.go` file. This is the diff-level oracle named in the issue (`Surface: Go diff`).

### AC5 — GitHub Action calls the public command
- **α evidence:** `.github/workflows/board-map.yml` diff is empty (or, if α determines a change is truly required, the diff still only invokes `cn issues map ...` and never `src/packages/cnos.issues/...` or a Node generator — α states which case applies and why in `self-coherence.md`).
- **β check:** Read the workflow; confirm it shells only into `./cn issues map ...`, no package-internal path, no `scripts/board`.

### AC6 — Package installation path is explicit
- **α evidence:** Because the command stays a compiled-in kernel dispatch (not package-command-content-class dispatch — see §6 guardrail), there is no package-install step to add; α states this explicitly in `self-coherence.md` rather than leaving it implicit, and confirms the workflow does not ad-hoc-copy package files or mutate `.cn/vendor/` manually (it doesn't touch `.cn/vendor/` at all today).
- **β check:** Confirm the workflow has no undeclared local-state assumptions and that the receipt (see AC10) names the shim disposition honestly.

### AC7 — Current board output behavior is preserved
- **α evidence:** Before/after comparison: α runs `cn issues map` from `main` (pre-move) and from `cycle/556` (post-move) against the same fixture/live data and diffs `board-data.json` — the diff should be empty or limited to expected non-semantic noise (timestamps, if any). α pastes the diff (or "empty diff") into `self-coherence.md`. No `docs/board/` (the old, pre-#545 path) resurrection.
- **β check:** β independently reproduces the before/after comparison.

### AC8 — No Node production generator
- **α evidence:** `find . -iname "*.mjs" -path "*board*"` (or equivalent) shows no new Node generator on the production path; `git diff --stat` shows no new `package.json`/`node_modules` production dependency.
- **β check:** β independently confirms no Node generator was introduced and the workflow does not invoke `node`/`npm`/`npx`.

### AC9 — CI remains green
- **α evidence:** All CI jobs pass on the cycle branch: `Go build & test` (`go`), `Binary verification` (`binary-verify`), `Package verification` (`package-verify`), `Package/source drift (I1)` (`package-source-drift`), `Protocol contract schema sync (I2)` (`protocol-contract-check`), `Repo link validation (I4)` (`link-check`), `SKILL.md frontmatter validation (I5)` (`skill-frontmatter-check`), `CDD artifact ledger validation (I6)` (`cdd-artifact-check`), `Dispatch repair-preflight guard` (`dispatch-repair-preflight`), `Dispatch closeout-integrity guard` (`dispatch-closeout-integrity`), plus `install-wake-golden.yml`. α names the exact green run URL(s) in `self-coherence.md`.
- **β check:** β re-verifies via `gh run list --branch cycle/556 --json status,conclusion,name` (or equivalent) that every job is green on the merge-candidate commit, independent of α's claim.

### AC10 — Receipt records the #216 relationship
- **α evidence:** `self-coherence.md` (and later `alpha-closeout.md`) states explicitly and unambiguously: **"`cn issues map` remains a temporary built-in shim until #216 lands. Domain logic is isolated under the `cnos.issues` package boundary (`src/packages/cnos.issues/commands/issues-map/`), but dispatch is still through the compiled-in kernel registration, not through the package-command-content-class exec-dispatch mechanism (`PACKAGE-SYSTEM.md` §7)."** This is Option B from the issue's AC10 oracle, and matches the operator's clarifying-comment default. α must not claim Option A (package-provided through documented command-discovery) — that is false given the architecture pinned in §6/§7 below.
- **β check:** β confirms the receipt states Option B in those or equivalent unambiguous words, and that the receipt is honest against the actual diff (no claim of "package-provided" while the kernel registration in `main.go` is unchanged).

## 5. Source-of-truth table

(Adapted from the issue's own "Source of truth" table; unchanged claims are carried forward,
with γ's verification notes appended.)

| Claim / surface | Canonical source | Status | γ verification note |
|---|---|---|---|
| Command boundary | `docs/architecture/DESIGN-CONSTRAINTS.md` §3.2 | Current | Read directly: "`cli/` owns dispatch. Domain packages own logic. No `cmd_*.go` file may contain domain logic beyond argument parsing and dispatch." `cmd_issues_map.go` already conforms. |
| Command naming | `docs/architecture/DESIGN-CONSTRAINTS.md` §3.1 | Current | Read directly: git-style `cn <noun> <verb>`, hyphenated registry key. `cn issues map` → `issues-map` already conforms. |
| Package content classes | `docs/reference/packages/PACKAGE-SYSTEM.md` §1.1 | Current | Read directly: 7 content classes incl. `commands` (directory-tree copy, discovered by directory presence). |
| Package command discovery | `docs/reference/packages/PACKAGE-SYSTEM.md` §7 | Current | Read directly: 3-tier precedence — built-in (authoritative, silently shadows) → repo-local (`.cn/commands/cn-<name>`) → vendored package commands (exec-dispatch via `cn.package.json`'s declared `commands` object). **`cn issues map` sits at tier 1 (built-in) today and this cycle keeps it there** — see §6 guardrail below. |
| Kernel-is-a-package direction | `docs/reference/runtime/GO-KERNEL-COMMANDS.md` §1.2–1.3 | Current | Read directly: the target bootstrap-kernel set is `help, init, setup, deps, build, doctor, status, update`. `issues-map` and `cdd-verify` are *not* in that target set — meaning the long-run direction is for both to eventually migrate, but that migration is #216-shaped, not this cycle. |
| Package source/install model | `docs/reference/packages/BUILD-AND-DIST.md` | Current | Not read in depth this scaffold; α should confirm `src/packages/` vs `dist/` vs `.cn/vendor/` framing matches PACKAGE-SYSTEM.md §2 before writing `cn.package.json`. |
| Lean kernel / package commands | #216 | Open | This cycle does not solve it; see §6. |
| Issue taxonomy | `docs/development/issues/TAXONOMY.md` | Current | Read directly: `kind/*` (exactly one), `area/*` (one or more), priority (`P0`–`P3`/`priority/deferred`), `status:*`, `protocol:*`, `dispatch:cell`. This is what `cnos.issues/skills/taxonomy/SKILL.md` should *cite*, not fork. |
| Board docs/output | `docs/development/board/` | Current | Generated output directory; unchanged by this cycle except the one path-reference update in `README.md`. |
| Current command (paths) | `src/go/internal/issuesmap/` + `src/go/internal/cli/cmd_issues_map.go` | Current, **this is the refactor boundary** | Read directly (full file contents): confirmed shapes are exactly as described in §3 above. |
| Current board Action | `.github/workflows/board-map.yml` | Current | Read directly (full file contents): already calls `cn issues map --repo "$GITHUB_REPOSITORY" --out docs/development/board`; builds `cn` from source in the same job (`go build -o "$GITHUB_WORKSPACE/cn" ./cmd/cn`); no package-install step exists today. |
| cnos#392 co-location precedent | `src/go/go.mod`, `go.work`, `src/packages/cnos.cdd/commands/cdd-verify/go.mod`, `src/go/internal/cli/cmd_cdd_verify.go` | Current, **the concrete architectural template for this cycle** | Read directly: `cdd-verify` is a `SourceKernel`/`TierKernel` built-in whose Go domain package lives at `src/packages/cnos.cdd/commands/cdd-verify/` as its own Go module (`module github.com/usurobor/cnos/packages/cnos.cdd/commands/cdd-verify`), wired into `src/go/go.mod` via `require` + `replace`, and into `go.work`'s `use` list. `src/packages/cnos.cdd/cn.package.json`'s `"commands"` object declares only `cdd-status` — **not** `cdd-verify** — confirming the directory is Go-source co-location, not a declared package-command-content-class entry. |
| `cn build --check` / `cn doctor` permissiveness re: undeclared `commands/<id>/` | `src/go/internal/pkgbuild/build.go` (`CheckOne`, `checkCommandEntrypoints`, `FindContentClasses`, `DerivePacklist`), `src/go/internal/discover/discover.go` (`ScanPackageCommands`), `src/go/internal/doctor/doctor.go` (`ValidateCommands`) | Current, verified by direct source read | Both tools are forward-check-only against the manifest's *declared* `commands` map; neither reverse-checks directory presence against the manifest. An undeclared `commands/<id>/` subdirectory (no `cn-<id>` entrypoint, no manifest key) passes cleanly — this is the exact live behavior keeping `cnos.cdd`'s `cdd-verify/` subdirectory CI-green today via the `package-source-drift` (I1) job. |

## 6. Scope guardrails (binding — restated from the operator's clarifying comment)

The operator's clarifying comment on the issue (2026-07-02) **narrows and, where it conflicts,
supersedes** the original issue body's more permissive "Proposed package shape" language. These
are hard constraints, not suggestions:

1. **Do not implement generic package-command discovery in this issue.** The 3-tier precedence
   mechanism (`PACKAGE-SYSTEM.md` §7) already exists and is out of scope to redesign. This cycle
   does not flip `cn issues map` onto that mechanism.
2. **`cn issues map` may remain a temporary built-in shim** — and per the grounding in §3/§5
   above, γ's recommendation (not a soft suggestion — pin this unless α's own discovery
   surfaces a concrete blocker, in which case escalate via `gamma-clarification.md` rather than
   improvising) is that it **does** remain a built-in shim this cycle, with its Go domain source
   relocated to co-locate with the new package per the cnos#392 precedent.
3. **Only add `commands/issues-map/` if it is actually executable.** It is — but *not* through
   the package-command exec-dispatch mechanism (`PACKAGE-SYSTEM.md` §7); it is executable
   because `src/go/internal/cli/cmd_issues_map.go` (compiled into the `cn` binary, registered in
   `src/go/cmd/cn/main.go:47`) imports it directly as a Go module and calls `issuesmap.Run(...)`
   on every invocation of the kernel-dispatched `cn issues map` command. **Do not** declare
   `issues-map` under `cnos.issues/cn.package.json`'s `"commands"` object — doing so would
   create a dead, unreachable declaration (the kernel built-in always shadows it silently per
   §7's precedence rule) and would itself be "a fake package command directory that is not
   used" from the declaration side, even though the underlying Go code is genuinely executed.
   The distinction to hold in mind throughout: **Go-source co-location under `commands/<id>/`
   ≠ a declared package command.** `cdd-verify` is the live, CI-green proof this distinction is
   real and already load-bearing in this repo.
4. **The Action must call the public `cn issues map` command only** — it already does
   (`.github/workflows/board-map.yml`); do not make it call package internals, and do not add
   `cn deps restore` / package-install plumbing (there is nothing to install — the command
   dispatches via the kernel binary, not via `.cn/vendor/packages/`).
5. **No Node production generator.** Not present today; do not introduce one.
6. **The receipt must honestly state the disposition** — Option B (temporary built-in shim,
   tied to #216) per AC10 above. Do not claim Option A.
7. **Do not solve #216.** No generic command-discovery redesign, no bootstrap-kernel-set
   migration for `issues-map` or any other command, no Node/package-install plumbing beyond what
   is already present.
8. **Non-goals carried forward verbatim from the issue:** no standalone external GitHub Action
   extraction, no GitHub Pages enablement, no board-visualization rewrite, no
   effort/priority/taxonomy semantic changes, no dispatch/wake/CDD/label behavior changes, no
   Demo 0.

## 7. Friction notes

Grounded in direct file reads this scaffold session, for α/β to watch for:

- **The pinned base SHA was stale by 2 bot-only commits at scaffold time** (see §0). Not a
  blocker; branched from live `origin/main`. If α's session also observes `origin/main` has
  moved further by the time α starts (e.g. another `board-map.yml` regeneration fired), that is
  expected background noise from the very system this issue is about to re-home — it does not
  require a rebase of `cycle/556` unless it touches files α is editing.
- **`cn.package.json` version convention split exists across the repo.** Some packages track the
  repo release version (`cnos.cdd`, `cnos.core`, `cnos.eng` — all `3.82.0` at scaffold time, the
  current release per `git tag --sort=-v:refname`); others version independently starting at
  `0.1.0`/`0.2.0`/`0.3.0` (`cnos.cds`, `cnos.cdr`, `cnos.handoff`, `cnos.kata`). `cnos.issues` is
  a brand-new package like `cnos.cds`/`cnos.cdr`/`cnos.handoff` — follow that convention:
  `"version": "0.1.0"`, `"engines": {"cnos": ">=3.82.0"}` (or the current tag at α's
  implementation time if it has advanced — check `git tag --sort=-v:refname | grep -v '^v' |
  head -1` fresh rather than trusting this note's number blindly).
- **The `//go:embed templates/board.template.html` directive is relative to the Go source
  file's own directory.** Move `templates/` and `testdata/` together with the `.go` files as one
  atomic directory move — do not split them across old/new locations, and do not `git mv` the
  `.go` files first and the subdirectories second in a way that leaves a moment where the embed
  target is missing (a single commit doing the whole directory move is simplest and matches how
  cnos#392 moved `cdd-verify`).
- **`issuesmap_test.go` references `testdata/issues.json` by relative path** (`--fixture
  testdata/issues.json` in at least two subtests) — this moves for free if the whole directory
  moves together; verify `go test ./...` still passes from the new module root after the move.
- **`docs/development/board/README.md:42` is the only prose reference outside the code itself**
  to the old `src/go/internal/issuesmap/` path (confirmed via repo-wide grep for
  `internal/issuesmap`) — update it, and do a final grep before closing to make sure no other
  reference was missed (a second reviewer's grep is cheap insurance; β should re-run the same
  grep independently rather than trusting α's claim).
- **`go.work` and `src/go/go.mod` both need edits, in the same commit as the directory move** —
  the cnos#392 precedent's own commit message notes "the workspace file (`../../go.work`) and
  this replace directive both ensure the local module is used; the replace makes `go mod tidy`
  work without network resolution." Follow the same shape exactly: `module
  github.com/usurobor/cnos/packages/cnos.issues/commands/issues-map` in the new `go.mod`; a
  `require`/`replace` pair in `src/go/go.mod` pointing `=> ../packages/cnos.issues/commands/
  issues-map`; the new path added to `go.work`'s `use (...)` list. Run `go mod tidy` (or
  equivalent workspace sync) in both modules and confirm `go build ./...` / `go vet ./...`
  succeed from `src/go/` before considering the move done — this is exactly what CI's `go` job
  (`Go build & test`) will check.
- **The `Framework command surface` step in `build.yml`'s `package-verify` job** exercises
  `kata-list`/`kata-run`/`kata-judge` — unrelated to this cycle's surface, but it is evidence CI
  does walk installed package commands at runtime for *some* packages; α does not need to touch
  it, just should not be surprised it exists when reading `build.yml` in full.
- **Do not conflate this cycle's "package scoping" decision with #216.** The `commands/
  issues-map/` directory this cycle creates is a Go-module co-location choice within the current
  built-in-dispatch architecture; it is not a step toward package-command exec-dispatch, and
  should not be described as such in any artifact. If a future #216 cycle later *does* flip
  `issues-map` onto real package-command exec-dispatch, it will need to (a) remove the kernel
  registration in `main.go`, (b) add a `commands.issues-map` entry to `cnos.issues/
  cn.package.json` with a `cn-issues-map` entrypoint script, and (c) add package-install
  plumbing (`cn deps restore` or equivalent) to `board-map.yml` — none of which this cycle does.

---

## 8. α prompt

```text
You are α. Project: cnos (github.com/usurobor/cnos).
Load src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md and follow its load order.
Issue: gh issue view 556 --repo usurobor/cnos --json title,body,state,comments
Branch: cycle/556
Scaffold: .cdd/unreleased/556/gamma-scaffold.md (read this in full before coding — it names the
  concrete files/dirs you touch, the per-AC oracle list, the scope guardrails, and the
  cnos#392 precedent this cycle mirrors)
Tier 3 skills: src/packages/cnos.core/skills/write/SKILL.md, docs/architecture/DESIGN-CONSTRAINTS.md (§3),
  docs/reference/packages/PACKAGE-SYSTEM.md (§1.1, §7), docs/reference/runtime/GO-KERNEL-COMMANDS.md,
  docs/development/issues/TAXONOMY.md, docs/development/issues/TRIAGE.md

## Scope guardrails (binding — see gamma-scaffold.md §6 for full grounding)

1. Do NOT implement generic package-command discovery. The 3-tier command-resolution mechanism
   (PACKAGE-SYSTEM.md §7: built-in > repo-local > vendored package commands) is out of scope to
   redesign.
2. `cn issues map` stays a compiled-in kernel command (SourceKernel/TierKernel), registered in
   src/go/cmd/cn/main.go:47. Do not remove or change this registration.
3. Relocate the Go domain source from src/go/internal/issuesmap/ to
   src/packages/cnos.issues/commands/issues-map/ as its own Go module, mirroring the cnos#392
   precedent at src/packages/cnos.cdd/commands/cdd-verify/ EXACTLY:
   - new go.mod: `module github.com/usurobor/cnos/packages/cnos.issues/commands/issues-map`
   - add to go.work's `use (...)` list
   - add require + replace pair in src/go/go.mod (mirror the cdd-verify require/replace block)
   - update src/go/internal/cli/cmd_issues_map.go's import path only — it MUST stay a thin
     dispatcher (argument parsing + one-line delegation). No domain logic in any cmd_*.go file.
   - move templates/ and testdata/ together with the .go files in the same commit (the
     //go:embed directive and the test's relative testdata path both depend on this)
4. Do NOT declare `issues-map` under cnos.issues/cn.package.json's "commands" object. That
   object is for real package-command-content-class exec-dispatch entries (see cdd-status for
   the live example of that shape). `commands/issues-map/` in this cycle is Go-source
   co-location only — it is executed via the kernel's compiled-in dispatch, not via package-
   command discovery. Declaring it would create a dead entry (silently shadowed by the kernel
   built-in per PACKAGE-SYSTEM.md §7's precedence rule) — exactly the "fake unused command
   directory" the operator's clarifying comment on the issue forbids.
5. Create: src/packages/cnos.issues/cn.package.json (schema cn.package.v1, name cnos.issues,
   version 0.1.0, engines >= the current repo tag — check `git tag --sort=-v:refname | grep -v
   '^v' | head -1` fresh), SKILL.md (states cnos.issues owns issue taxonomy, board mapping, and
   Issue Pivot generation, and states the command's current dispatch disposition honestly —
   built-in shim, not package-command-provided, tied to #216), skills/map/SKILL.md,
   skills/taxonomy/SKILL.md (cites docs/development/issues/TAXONOMY.md, does not fork it),
   skills/triage/SKILL.md (cites docs/development/issues/TRIAGE.md).
6. Do NOT edit .github/workflows/board-map.yml unless you discover a concrete reason it must
   change (it already calls only `cn issues map --repo "$GITHUB_REPOSITORY" --out
   docs/development/board` with no package-internal calls and no `cn deps restore`). If you find
   a reason, log it in gamma-clarification.md before editing rather than editing silently.
7. Update docs/development/board/README.md:42's path reference; re-grep the repo for
   `internal/issuesmap` before closing to confirm no other reference was missed.
8. No Node production generator. No standalone external GitHub Action extraction. No GitHub
   Pages enablement. No board-visualization rewrite. No effort/priority/taxonomy semantic
   changes. No dispatch/wake/CDD/label behavior changes. No Demo 0. Do not solve #216.
9. Walk all 10 ACs in the issue body; gamma-scaffold.md §4 names the exact oracle/evidence for
   each. AC10 in particular: your receipt (self-coherence.md, and later alpha-closeout.md) MUST
   state Option B verbatim-equivalent: "cn issues map remains a temporary built-in shim until
   #216 lands; domain logic is isolated under the cnos.issues package boundary but dispatch
   stays through the compiled-in kernel registration, not package-command exec-dispatch."

## Implementation contract (pinned by δ; α MUST NOT improvise)

| Axis | Pinned value |
|---|---|
| Language | Go (unchanged — the existing domain implementation is already Go; no new language introduced) |
| CLI integration target | `cn` subcommand (`cn issues map`), compiled-in kernel dispatch — unchanged. Registry key `issues-map` (hyphenated, per DESIGN-CONSTRAINTS.md §3.1); user-facing surface `cn issues map` (space-separated noun+verb) |
| Package scoping | Move `src/go/internal/issuesmap/{issuesmap.go,fetch.go,issuesmap_test.go,templates/,testdata/}` → `src/packages/cnos.issues/commands/issues-map/` as its own Go module (module path `github.com/usurobor/cnos/packages/cnos.issues/commands/issues-map`), wired via `require`+`replace` in `src/go/go.mod` and added to `go.work`'s `use` list — mirrors the cnos#392 `cdd-verify` precedent exactly. `src/go/internal/cli/cmd_issues_map.go` stays at its current path, thin, import updated. This directory is NOT declared in `cnos.issues/cn.package.json`'s `commands` object (Go-source co-location, not package-command-content-class dispatch) |
| Existing-binary disposition | Preserve — single `cn` binary, unchanged; `IssuesMapCmd` registration in `src/go/cmd/cn/main.go:47` unchanged; no new binary, no package-command-exec-dispatch shim script |
| Runtime dependencies | None (Go stdlib only — unchanged from today; no new external dependency introduced by the relocation) |
| JSON/wire contract preservation | Preserve as-is — `board-data.json` shape, the embedded HTML template's data-splice contract, and the `cn issues map --repo <owner/repo> --out <dir>` CLI flag shape are all unchanged. Before/after board output must diff empty (or non-semantic-noise-only) per AC7 |
| Backward-compat invariant | Existing rules preserved; new content additive. `cn issues map`'s flags, taxonomy parsing (`kind/*`, `area/*`, `P0`–`P3`, `status:*`, `dispatch:cell`, `protocol:*`, `effort/*`), effort/`unestimated` computation, `kind/tracking` effort-rollup exclusion, and output paths are all unchanged; only the Go source's on-disk location and the new `cnos.issues` package metadata/skill/doctrine surfaces are additive |
```

---

## 9. β prompt

```text
You are β. Project: cnos (github.com/usurobor/cnos).
Load src/packages/cnos.cdd/skills/cdd/beta/SKILL.md and follow its load order.
Issue: gh issue view 556 --repo usurobor/cnos --json title,body,state,comments
Branch: cycle/556
Scaffold: .cdd/unreleased/556/gamma-scaffold.md (read in full — §4 is your per-AC oracle
  checklist; §6 is the binding scope-guardrail list; §7 is the friction-notes list you should
  independently re-verify rather than trust)

## AC oracle list to walk independently (see gamma-scaffold.md §4 for full detail per AC)

AC1  — src/packages/cnos.issues/cn.package.json exists, valid, `cn build --check` recognizes it
        (re-run independently; do not trust α's pasted output alone).
AC2  — src/packages/cnos.issues/SKILL.md (or equivalent) states what the package owns; not an
        empty wrapper.
AC3  — `cn issues map --repo usurobor/cnos --out <tmpdir>` works from a binary YOU build on the
        cycle branch.
AC4  — git diff on any cmd_*.go file (esp. cmd_issues_map.go) contains ONLY argument-parsing/
        dispatch changes — no taxonomy/rendering/label logic.
AC5  — .github/workflows/board-map.yml calls only `cn issues map ...`; no package-internal path,
        no scripts/board, no standalone Node generator.
AC6  — no `cn deps restore` / package-install plumbing was added unless independently justified;
        no ad-hoc .cn/vendor/ mutation; receipt states the shim disposition honestly (see AC10).
AC7  — before/after board-data.json diff is empty or non-semantic-noise-only; no docs/board/
        (pre-#545 path) resurrection.
AC8  — no Node production generator anywhere on the production path (grep for .mjs under
        board-related paths; grep the workflow for node/npm/npx).
AC9  — every CI job green on the merge-candidate commit: go (Go build & test), binary-verify,
        package-verify, package-source-drift (I1), protocol-contract-check (I2), link-check (I4),
        skill-frontmatter-check (I5), cdd-artifact-check (I6), dispatch-repair-preflight,
        dispatch-closeout-integrity, install-wake-golden — verify via `gh run list --branch
        cycle/556 --json status,conclusion,name` yourself, do not trust α's claim alone.
AC10 — self-coherence.md (and alpha-closeout.md at convergence) states Option B unambiguously:
        temporary built-in shim tied to #216, domain logic isolated under cnos.issues, dispatch
        still via compiled-in kernel registration (not package-command exec-dispatch). Flag as a
        REQUEST CHANGES (severity D, classification implementation-contract per
        alpha/SKILL.md §3.6 / beta/SKILL.md Rule 7) if the receipt instead claims or implies
        Option A (package-provided through documented command-discovery) — that would be false
        given the pinned Package-scoping axis, and false-positive receipts are exactly the
        drift class cnos#389/#391 anchor.

## Additional verification specific to this cycle (beyond the standard AC walk)

- Confirm `src/packages/cnos.issues/cn.package.json`'s `"commands"` object does NOT contain an
  `issues-map` key. If it does, that is a scope-guardrail violation (gamma-scaffold.md §6 item
  3) — a declared-but-shadowed dead entry, or worse, evidence of an attempt to actually flip
  command discovery (out of scope per the operator's clarifying comment on the issue).
- Confirm `src/go/cmd/cn/main.go`'s `reg.Register(&cli.IssuesMapCmd{})` line is unchanged
  (Existing-binary-disposition axis: preserve).
- Confirm the relocated Go module's `go.mod`, `go.work`, and `src/go/go.mod`'s require/replace
  block mirror the cnos#392 `cdd-verify` shape (module path pattern, replace directive shape).
- Independently re-grep the repo for `internal/issuesmap` to confirm no stray reference was
  missed (gamma-scaffold.md §7 flags this as a friction point α should also check — verify both
  did it).
- Confirm the `//go:embed templates/board.template.html` directive still resolves (i.e.
  `templates/` moved with the Go source, not left behind or duplicated).
```
