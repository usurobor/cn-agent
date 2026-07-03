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
  - ACs
  - Self-check
  - Debt
  - CDD Trace
  - Review-readiness
---

# α self-coherence — cnos#556

## Gap

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

## Skills

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

## ACs

Implementation SHA (last implementation commit before this artifact's own
commits): `88156120` (directory move + wiring). Package-files commit:
`b39b44de`. Doc-update commit: `13198122`.

### AC1 — `cnos.issues` package exists

`src/packages/cnos.issues/cn.package.json` committed at `b39b44de`:

```json
{
  "schema": "cn.package.v1",
  "name": "cnos.issues",
  "version": "0.1.0",
  "kind": "package",
  "engines": {
    "cnos": ">=3.82.0"
  }
}
```

`cn build --check` (binary built from `src/go` at cycle-branch HEAD):

```
✓ cnos.cdd: valid
✓ cnos.cdd.kata: valid
✓ cnos.cdr: valid
✓ cnos.cds: valid
✓ cnos.core: valid
✓ cnos.eng: valid
✓ cnos.handoff: valid
✓ cnos.issues: valid
✓ cnos.kata: valid
✓ All packages valid.
```

`cnos.issues` is recognized and valid. **AC1: met.**

### AC2 — Package owns issue-board cognition

`src/packages/cnos.issues/SKILL.md` §"Core Principle" states explicitly:
"`cnos.issues` is the package home for issue-board cognition: it owns issue
taxonomy, board mapping, and Issue Pivot generation as a domain, distinct
from the CLI dispatch that invokes it," then names the three owned surfaces
(taxonomy, triage, board mapping) each pointing at its own sub-skill. This is
not an empty wrapper — it also carries the full command-dispatch-disposition
statement (§AC10 below) and an explicit non-goals list. **AC2: met.**

### AC3 — Public command remains `cn issues map`

Ran against a `cn` binary built from the cycle-branch HEAD, using the
relocated module's own fixture (network-independent, deterministic):

```
$ ./cn issues map --fixture src/packages/cnos.issues/commands/issues-map/testdata/issues.json --out /tmp/board-out
cn issues map: taxonomy gaps (advisory):
  missing kind/*: #547
  missing area/*: #547
cn issues map: wrote 5 open issues to /tmp/board-out (index.html, board-data.json, README.md)
```

Exit code 0; `index.html`, `board-data.json`, `README.md` all written. Also
confirmed with `--repo usurobor/cnos --out <tmpdir>` (live-API path) during
development — same success shape (issue count varies with live state, which
is expected). **AC3: met.**

### AC4 — CLI dispatch remains thin

`git diff` on `src/go/internal/cli/cmd_issues_map.go` (commit `88156120`):
only the import path changed (old `internal/issuesmap` → new
`packages/cnos.issues/commands/issues-map`) plus an explanatory doc-comment
expansion. No taxonomy/rendering/label logic added — `Run()` is still a
single-line delegation:

```go
func (c *IssuesMapCmd) Run(ctx context.Context, inv Invocation) error {
	return issuesmap.Run(ctx, inv.Args, inv.Stdin, inv.Stdout, inv.Stderr)
}
```

CI's own "Dispatch boundary check (INVARIANTS.md T-002)" step (scans every
`internal/cli/cmd_*.go` for banned domain-logic imports) reproduced locally:
`✓ Dispatch boundary: all cmd_*.go files are thin wrappers`. **AC4: met.**

### AC5 — GitHub Action calls the public command

`git diff --stat` on `.github/workflows/board-map.yml` across this entire
cycle: empty (file untouched). Read in full: it still only runs
`./cn issues map --repo "$GITHUB_REPOSITORY" --out docs/development/board`
after `go build -o "$GITHUB_WORKSPACE/cn" ./cmd/cn` from `working-directory:
src/go` — no package-internal path, no `scripts/board`, no Node. The
`go build` step from `src/go` picks up the new module transparently via
`go.work` (Go workspace mode walks up from cwd to find `go.work` at the repo
root) — confirmed locally by building exactly that way (see AC9). No reason
to edit the workflow was discovered; per γ's scaffold §6 guardrail #4, no
edit made, and none logged to `gamma-clarification.md` since none was
needed. **AC5: met (no change required, confirmed correct as-is).**

### AC6 — Package installation path is explicit

Because `cn issues map` stays a compiled-in kernel dispatch (not a declared
package command — see AC10), there is nothing to install: no `cn deps
restore` step exists in `board-map.yml` today, and none was added. The
workflow does not touch `.cn/vendor/` at all. `cnos.issues/cn.package.json`
has no `commands` object at all (confirmed: `grep -c "issues-map"
cn.package.json` → 0). **AC6: met** — stated explicitly here rather than
left implicit, per γ's oracle.

### AC7 — Current board output behavior is preserved

Before/after comparison, fixture-driven (deterministic, not live-API,
so the diff isolates code-motion effects from live-issue-state drift):

1. Built `cn` from `origin/main` (`4fe8e433`, pre-move) in a separate
   worktree; ran `cn issues map --fixture .../testdata/issues.json --out
   /tmp/board-out-pre`.
2. Built `cn` from `cycle/556` HEAD (post-move); ran the same fixture to
   `/tmp/board-out`.
3. `diff /tmp/board-out-pre/board-data.json /tmp/board-out/board-data.json`
   → **empty (byte-identical)**.
4. `diff /tmp/board-out-pre/index.html /tmp/board-out/index.html` →
   **empty (byte-identical)**.
5. `diff /tmp/board-out-pre/README.md /tmp/board-out/README.md` → differs
   only in the intentional path-reference line:
   ```
   42c42,43
   < The generator is the Go command `cn issues map` (`src/go/internal/issuesmap/`).
   ---
   > The generator is the Go command `cn issues map`
   > (`src/packages/cnos.issues/commands/issues-map/`).
   ```
   This is the expected, non-semantic (prose path pointer only) change —
   not a board-visualization or data-model change. No `docs/board/`
   (pre-#545 path) resurrection. **AC7: met.**

### AC8 — No Node production generator

`find . -iname "*.mjs" -path "*board*"` → no hits (checked repo-wide).
`git diff --stat origin/main..HEAD` introduces no `package.json` /
`node_modules` on the production path. `.github/workflows/board-map.yml`
contains no `node`/`npm`/`npx` invocation (unchanged, confirmed by full
read). **AC8: met.**

### AC9 — CI remains green

Local reproduction of every job in `.github/workflows/build.yml` and
`.github/workflows/board-map.yml`'s relevant steps, run against
cycle-branch HEAD (`4fda6fc9` at time of this section):

- **`go` (Go build & test):** `go build ./...`, `go build -o cn ./cmd/cn`,
  `go test ./... -v`, `go vet ./...` — all pass, 0 failures (all 14
  `src/go` packages `ok`). Dispatch-boundary check: pass. Smoke test
  (`./cn help`) lists `issues-map` correctly.
- **New module** (`src/packages/cnos.issues/commands/issues-map`):
  `go build ./...`, `go vet ./...`, `go test ./... -v` — 4 test functions, 8
  pass lines (`TestToRecord_LabelParsingAndEffort` + 4 subtests,
  `TestEffortWeights`, `TestRun_Fixture`, `TestRun_Stdin`), 0 failures.
- **`binary-verify` / `package-verify`:** not run end-to-end locally (no
  Tier-1/Tier-2 kata harness invoked here), but the underlying `cn build`
  step they depend on was run directly (AC1) and passed.
- **`package-source-drift` (I1):** `cn build --check` → all packages valid
  (AC1 excerpt above).
- **`protocol-contract-check` (I2):** `diff docs/reference/schemas/protocol-contract.json
  tests/fixtures/protocol-contract.json` → identical (unaffected by this
  cycle's diff).
- **`link-check` (I4):** installed `lychee` v0.24.2 locally (`cargo
  install`), ran `lychee --offline --no-progress --hidden --config
  lychee.toml '**/*.md'` → initially found 3 broken relative links in the
  new `skills/taxonomy/SKILL.md` / `skills/triage/SKILL.md` (off-by-one
  `../` depth); fixed; re-ran → `🚫 0 Errors` across 1420 total / 437 unique
  links.
- **`skill-frontmatter-check` (I5):** installed `cue` v0.17.0 locally,
  ran `./scripts/ci/validate-skill-frontmatter.sh` (full repo) →
  `✓ 99 SKILL.md validated; no findings.` (includes the 4 new
  `cnos.issues` SKILL.md files).
- **`cdd-artifact-check` (I6):** built `cn`, ran `./cn cdd verify
  --unreleased --exceptions .cdd/exceptions.yml` from repo root 5
  consecutive times → `PASSED with warnings` every time (0 failed, ~79
  warnings, consistent with pre-existing baseline noise — e.g. cycle #512's
  known, exception-backed missing `self-coherence.md` per
  `.cdd/exceptions.yml`'s `cdd-unreleased-512-dispatch-repair-loop` entry).
  Also ran `src/packages/cnos.cdd/commands/cdd-verify/test-fixtures.sh`
  (unaffected by this cycle — no cdd-verify source touched).
- **`dispatch-repair-preflight`:** `./scripts/ci/check-dispatch-repair-preflight.sh`
  → pass (unaffected by this cycle).
- **`dispatch-closeout-integrity`:** `./scripts/ci/check-dispatch-closeout-integrity.sh`
  → pass, including its self-test (unaffected by this cycle).
- **`install-wake-golden.yml`:** not exercised locally (no install/activation
  surface touched by this cycle; the wake/golden contract is unrelated to a
  Go-domain relocation + package-metadata addition).

**Remote-CI update (post local-reproduction round):** pushing this
self-coherence artifact's own commits to `origin/cycle/556` surfaced a
genuine remote-only failure the local reproduction above had not caught —
`cdd-artifact-check (I6)` hard-failed on commits `16a0fcc6`, `4fda6fc9`,
`dc1e392c`, `e889e5ba` (root cause + fix documented in §Debt item 4: this
file's own "§"-prefixed section headers didn't match `ledger.go`'s literal
`"## Gap"`/`"## Skills"`/`"## ACs"`/`"## CDD Trace"` match, compounded by
this cycle being classified `"small-change"` — no `beta-review.md` yet —
which turns missing-section detection into a hard fail instead of a warn).
Fixed by renaming the headers; the fix commit `f3fd05bc` reached
**remote GitHub Actions run
[28627809025](https://github.com/usurobor/cnos/actions/runs/28627809025),
conclusion `success`, all 10 jobs green** (`go`, `binary-verify`,
`package-verify`, `package-source-drift` (I1), `protocol-contract-check`
(I2), `link-check` (I4), `skill-frontmatter-check` (I5),
`cdd-artifact-check` (I6), `dispatch-repair-preflight`,
`dispatch-closeout-integrity`), confirmed at 23:24:08 UTC via `gh run view
28627809025 --json jobs`. `install-wake-golden.yml` did not run for this
branch/HEAD (its trigger scope does not include `cycle/*` pushes — expected,
per this cycle's non-goals).
**AC9: met — remote CI green on `f3fd05bc`, independently re-checkable by
β via `gh run list --branch cycle/556`.**

### AC10 — Receipt records the #216 relationship

Stated verbatim (matching γ's scaffold §4 AC10 oracle and the operator's
Option-B default) in `src/packages/cnos.issues/SKILL.md` §"Command-dispatch
disposition":

> `cn issues map` remains a temporary built-in shim until #216 lands.
> Domain logic is isolated under this package's boundary
> (`src/packages/cnos.issues/commands/issues-map/`, its own Go module,
> mirroring the cnos#392 `cdd-verify` precedent), but dispatch stays through
> the compiled-in kernel registration, not through the package-command
> exec-dispatch mechanism (`PACKAGE-SYSTEM.md` §7).

This is **Option B**, not Option A: no claim of "package-provided through
documented command-discovery" is made anywhere in the diff or this
artifact. The kernel registration `reg.Register(&cli.IssuesMapCmd{})` at
`src/go/cmd/cn/main.go:47` is unchanged, and `cnos.issues/cn.package.json`
declares no `commands` object at all — both facts are consistent with the
Option-B statement and neither is contradicted anywhere else in the diff.
**AC10: met.**

## Self-check

**Peer enumeration.** The family of peers for a CLI-dispatch-boundary
change is "every `cmd_*.go` file in `src/go/internal/cli/`." CI's own
dispatch-boundary check (INVARIANTS.md T-002) scans that whole directory,
not just the changed file, and I reproduced it locally against the full
directory — it passed for all files, not only `cmd_issues_map.go`. The
`cn.package.json` `commands`-object peer question (does any other package
wrongly declare a co-located-but-undispatched directory as a command?) was
checked by reading `cnos.cdd/cn.package.json` (the precedent) directly:
it declares only `cdd-status`, not `cdd-verify`, confirming the same
distinction this cycle's `cnos.issues/cn.package.json` also holds (no
`issues-map` key). I did not find or need a third co-located-command
package to complete this peer set — `cdd-verify`/`cdd-status` and
`issues-map` are the only two instances of this pattern in the repo today
(confirmed: only `cnos.cdd/commands/` and `cnos.issues/commands/` exist
under `src/packages/*/commands/`).

**Harness audit.** The schema-bearing surfaces here are `board-data.json`'s
shape and the embedded HTML template's data-splice contract — both
unchanged (AC7's byte-identical diff is the proof). Non-Go harnesses
touching this surface: the GitHub Action (`board-map.yml`, read in full,
unchanged, confirmed still calls only the public command); the embedded
`//go:embed templates/board.template.html` directive (verified it still
resolves post-move — the fixture-driven `cn issues map` run in AC3/AC7
would fail immediately if it didn't, since `render()` calls
`templates.ReadFile` against the embedded FS); the test fixture
`testdata/issues.json` (moved with the code, `go test` from the new module
root passes, confirming the relative path still resolves).

**Did α push ambiguity onto β?** No known case identified. Every AC has a
locally-reproduced command/diff/log as evidence (§ACs above), not a bare
claim. AC9's CI-green claim is now backed by a concrete remote run URL and
conclusion (`28627809025`, `success`), not just local reproduction — but β
is still expected to independently re-check via `gh run list --branch
cycle/556` per pre-review-gate row 10's transient-row discipline, and this
artifact says so plainly rather than treating the one observed green run as
permanently authoritative.

**Would a loaded skill have prevented remaining debt?** Partially: the
`alpha/SKILL.md` §2.5 "§Section" prose convention, applied literally to
actual Markdown headers, is what produced the header-naming bug documented
in §Debt item 3. No single loaded skill states "when a skill instructs you
to name a self-coherence section '§X' in prose, the literal Markdown header
must still be plain '## X'" — this is the kind of doc/tool-convention gap
this artifact surfaces as a factual pattern (§Debt item 3's "Pattern worth
flagging"), not a recommendation for α to act on unilaterally.

## Debt

1. **Resolved during this round.** Remote GitHub Actions run
   [28627809025](https://github.com/usurobor/cnos/actions/runs/28627809025)
   on fix commit `f3fd05bc` came back `success` on all 10 `Build`-workflow
   jobs — including `binary-verify` and `package-verify` (their Tier-1/
   Tier-2 kata harnesses ran on the actual runner, resolving what would
   otherwise have been a "not exercised locally" debt item). β should still
   independently re-confirm via `gh run list --branch cycle/556 --json
   status,conclusion,name` per pre-review-gate row 10's transient-row
   discipline (this row was true at observation time — 23:24:08 UTC on
   `f3fd05bc` — and must be re-checked if the branch advances further).
2. **`install-wake-golden.yml` did not run for this branch.** Confirmed via
   `gh run list --branch cycle/556` (only the `Build` workflow appears) —
   its trigger scope does not include `cycle/*` pushes. Unrelated surface
   regardless (no install/activation change in this diff); named here for
   completeness per γ's AC9 oracle list rather than silently omitted.
3. **cdd-artifact-check (I6) genuinely failed on the remote GitHub Actions
   run for commits `16a0fcc6`, `4fda6fc9`, `dc1e392c`, `e889e5ba`** — this
   was not a fluke; it reproduced on every one of those pushes. Root cause,
   found by reading `gh run view --log-failed` and then
   `src/packages/cnos.cdd/commands/cdd-verify/ledger.go` directly:
   - This self-coherence.md's `## §Gap` / `## §Skills` / `## §ACs` /
     `## §CDD Trace` headers used a "§" section-mark prefix (matching this
     artifact's own prose convention, and γ's scaffold's convention for
     citing sections). But `ledger.go`'s `sectionPresent()` does an exact
     line match on the literal string `"## Gap"` (etc.) — `"## §Gap"` is
     not equal to, nor prefixed by, `"## Gap "`, so the check saw all four
     as missing.
   - Separately, `classifyCycleType()` classifies a cycle as
     `"small-change"` (not `"triadic"`) whenever `self-coherence.md` (or
     `alpha-closeout.md`) exists but `beta-review.md` does not — which is
     indistinguishable, by file presence alone, from a full triadic cycle
     that simply hasn't reached β yet (this cycle's actual state at R0).
     `checkSmallChangeArtifacts()` then calls `validateSections(...,
     forUnreleased=false)`, which turns any missing section into a hard
     `checkFail` rather than the `checkWarn` a genuinely in-progress
     triadic cycle gets via `checkUnreleasedTriadicArtifacts`'s
     `forUnreleased=true` path. So a mid-flight α round of a real triadic
     cycle is, today, one missing/misnamed section away from a hard CI
     failure that a same-shaped small-change cycle would only warn on.
   - **Fix applied (this session, working-tree edit before the section
     commits' content was otherwise finalized):** renamed all six headers
     to the literal form the parser checks (`## Gap`, `## Skills`, `## ACs`,
     `## Self-check`, `## Debt`, `## CDD Trace`) via `sed`. Re-ran `cn cdd
     verify --unreleased --exceptions .cdd/exceptions.yml` from repo root 4
     times post-fix → consistent `108 passed, 0 failed, 79 warnings,
     PASSED with warnings`. This is a genuine, reproducible fix, not a
     re-roll of a flake.
   - **Pattern worth flagging (not a recommendation — factual observation
     per `alpha/SKILL.md` §2.8 voice rule):** the `alpha/SKILL.md` §2.5
     instructions and this repo's own `gamma-scaffold.md` both write
     section names with a "§" prefix in prose (e.g. "§Gap", "§ACs"), but
     `ledger.go`'s literal-header-match convention has no "§". A future
     α following the doc-prose convention literally when naming actual
     Markdown headers reproduces this exact failure. Distinct from D1-class
     findings noted in prior cycles; first occurrence observed in this
     session.
4. **The earlier local `cn cdd verify` run that showed "1 failed" before
   this fix was diagnosed** was, separately, also affected by an
   incorrect working directory relative to `--exceptions .cdd/exceptions.yml`
   in one early invocation (resolved by `cd`-ing to repo root before
   invoking). That was a local-invocation artifact, distinct from the
   header-naming bug in item 4 above, and is fully resolved by always
   invoking from repo root (as all evidence in §ACs AC9 and this section
   now does).

No known debt beyond the above. In particular: no scope creep occurred (no
Node generator, no #216 solving, no board-visualization rewrite, no
taxonomy semantic change, no `.github/workflows/board-map.yml` edit) — see
§ACs AC5/AC6/AC8 for the concrete non-events.

## CDD Trace

1. **Design artifact:** not authored separately — γ's `gamma-scaffold.md`
   (committed by γ, present on branch) served the design-artifact role for
   this cycle: it names the concrete architectural shape (cnos#392
   precedent, package-scoping decision, scope guardrails) before any code
   was written. α did not need to author a competing design memo; the
   scaffold's §6/§8 constitute the "design decided" state per
   `alpha/SKILL.md` §2.2's "design artifact (when required) or explicit
   'not required'" — required and already supplied by γ.
2. **Coherence contract:** this file's §Gap (issue, mode, base SHA, pinned
   implementation contract).
3. **Plan:** not written as a separate artifact — the scaffold's §3
   (surfaces table) and §8 (α prompt, itemized 1–9) constitute a sufficient
   implementation sequence for a single-cell relocation; no additional
   sequencing artifact was needed (single atomic move + two additive
   commits, no multi-phase rollout).
4. **Tests:** `issuesmap_test.go` moved unmodified (100% file identity per
   `git diff --stat`'s rename detection); re-verified green from the new
   module root (`go test ./...` → 4 test funcs / 8 pass lines, 0 failures).
   No new test files added — behavior is preserved, not extended, so no new
   test surface was required (confirmed via the AC7 byte-identical output
   diff).
5. **Code:** commit `88156120` (directory move + wiring: new `go.mod`,
   `go.work` `use` entry, `src/go/go.mod` `require`+`replace`,
   `cmd_issues_map.go` import update, embedded README-string path update,
   package doc-comment update).
6. **Docs:** commit `b39b44de` (`cnos.issues/SKILL.md`,
   `cn.package.json`, `skills/map|taxonomy|triage/SKILL.md`) and commit
   `13198122` (`docs/development/board/README.md` path-reference update).

   **Full diff enumeration** (`git diff --stat origin/main..HEAD`, every
   entry accounted for above or in §ACs):
   - `.cdd/unreleased/556/gamma-scaffold.md` — γ's artifact (not α-authored;
     present on branch, cited throughout this file).
   - `.cdd/unreleased/556/self-coherence.md` — this artifact.
   - `docs/development/board/README.md` — commit `13198122` (AC-mapped
     above: source-of-truth-table row, §7 friction note).
   - `go.work` — commit `88156120` (wiring; AC1/AC9).
   - `src/go/go.mod` — commit `88156120` (wiring; AC1/AC9).
   - `src/go/internal/cli/cmd_issues_map.go` — commit `88156120` (AC4).
   - `src/packages/cnos.issues/SKILL.md` — commit `b39b44de` (AC2, AC10).
   - `src/packages/cnos.issues/cn.package.json` — commit `b39b44de` (AC1,
     AC6).
   - `src/packages/cnos.issues/commands/issues-map/fetch.go` (renamed,
     100% identity) — commit `88156120`.
   - `src/packages/cnos.issues/commands/issues-map/go.mod` (new) — commit
     `88156120` (AC1 wiring precedent).
   - `src/packages/cnos.issues/commands/issues-map/issuesmap.go` (renamed,
     94% identity — the embedded README string + doc comment updated) —
     commit `88156120` (AC3, AC7).
   - `src/packages/cnos.issues/commands/issues-map/issuesmap_test.go`
     (renamed, 100% identity) — commit `88156120` (AC3/AC9 test evidence).
   - `src/packages/cnos.issues/commands/issues-map/templates/board.template.html`
     (renamed, 100% identity) — commit `88156120` (`//go:embed` target,
     §Self-check harness audit).
   - `src/packages/cnos.issues/commands/issues-map/testdata/issues.json`
     (renamed, 100% identity) — commit `88156120` (used directly for the
     AC3/AC7 fixture-driven proof).
   - `src/packages/cnos.issues/skills/map/SKILL.md` — commit `b39b44de`
     (AC2).
   - `src/packages/cnos.issues/skills/taxonomy/SKILL.md` — commit
     `b39b44de` (AC2; link-depth fixed post-lychee in a follow-up amend
     within the same logical commit boundary — see §Debt item 4's sibling
     note: the taxonomy/triage link fix was folded into the same working
     tree before the package-files commit landed, so no separate "fix"
     commit exists for it).
   - `src/packages/cnos.issues/skills/triage/SKILL.md` — commit `b39b44de`
     (AC2).

7. **Self-coherence:** this file, written incrementally per §2.5 (one
   section per commit: `16a0fcc6` §Gap, `4fda6fc9` §Skills, `dc1e392c`
   §ACs, `e889e5ba` §Self-check, `e75f594e` §Debt, `a381a7e9` §CDD Trace),
   followed by two correction commits after remote CI surfaced the
   header-naming defect: `f3fd05bc` (section-header fix, §Debt item 3) and
   `987db05a` (remote-CI-green confirmation + debt-item resolution). Both
   correction commits' own diffs are self-documented in their commit
   messages and in §Debt above, per the "commit-message closure claims are
   a peer of the artifact they fix" rule (`alpha/SKILL.md` §2.3).

**Caller-path trace for new modules (pre-review-gate row 12).** The new
module `github.com/usurobor/cnos/packages/cnos.issues/commands/issues-map`
has exactly one non-test caller: `src/go/internal/cli/cmd_issues_map.go`'s
`Run()` method, which calls `issuesmap.Run(ctx, inv.Args, inv.Stdin,
inv.Stdout, inv.Stderr)` — the same call shape as before the move (the
import path is the only thing that changed). No new function or module was
added beyond the relocation itself (the package's exported surface —
`Run`, `readme`, `render`, `splice`, etc. — is unchanged from before the
move; `git diff` on `issuesmap.go` shows only the two path-reference string
edits and the header-comment expansion, not a new exported symbol).

**Test assertion count from runner output (pre-review-gate row 13).**
`go test ./... -v` from the new module root: 4 `func Test*` declarations,
8 `--- PASS` lines in the raw runner output (`TestToRecord_LabelParsingAndEffort`
+ its 4 named subtests, `TestEffortWeights`, `TestRun_Fixture`,
`TestRun_Stdin`) — counted directly from runner output, not hand-enumerated
from source.

**Commit-author identity (pre-review-gate row 14).** `git log -1
--format='%ae' HEAD` → `alpha@cdd.cnos` (the canonical `{role}@cdd.cnos`
elision form for the cnos project) — confirmed for every commit authored
in this round; no drift to correct.

**γ-artifact presence at the rule-3.11b surface (pre-review-gate row 15).**
`git cat-file -e origin/cycle/556:.cdd/unreleased/556/gamma-scaffold.md` →
present. **γ-artifact at canonical §5.1 path.**

## Review-readiness

**Round 1 | base SHA: `4fe8e4333b36372f595201841fb76cc0c31acff4` (origin/main
at cycle-branch time; re-confirmed not advanced as of this signal via `git
merge-base --is-ancestor origin/main HEAD`) | implementation SHA:
`987db05a6ced71006cfb4728951fecd6fa1a976c` | branch CI: green — remote run
[28627809025](https://github.com/usurobor/cnos/actions/runs/28627809025)
succeeded on fix commit `f3fd05bc`, and a follow-up run succeeded on
`987db05a` (confirmed via `gh run list --branch cycle/556`, both `status:
completed`, `conclusion: success`) at 23:27:27 UTC | ready for β.**

All ten pre-review-gate rows (`alpha/SKILL.md` §2.6) are satisfied:
rebase-currency (row 1, re-verified above), CDD Trace through step 7 (row
2), tests present (row 3), every AC has evidence (row 4), known debt
explicit (row 5, §Debt), schema/shape audit (row 6, §Self-check), peer
enumeration (row 7, §Self-check), harness audit (row 8, §Self-check), no
mid-cycle patch to re-audit yet (row 9 — this is R0, β has not yet
responded), branch CI green on head commit (row 10, above), artifact
enumeration matches diff (row 11, §CDD Trace step 6), caller-path trace
(row 12), test assertion count (row 13), commit-author identity (row 14),
γ-artifact presence (row 15).

β: please poll `.cdd/unreleased/556/beta-review.md` on `origin/cycle/556`
and re-verify AC1–AC10 independently per γ's β-prompt oracle list
(`gamma-scaffold.md` §9), rather than trusting this artifact's claims
alone — in particular AC1 (`cn build --check`), AC3 (build `cn` from this
branch yourself), AC9 (`gh run list --branch cycle/556`), and AC10 (read
`cnos.issues/SKILL.md` directly for the Option-B statement).

## §R1

**Round 1 | repair author: alpha | base for this round: `origin/cycle/556` at
`bdad537f5e5a719d6f76374da282eb6764fd014b` (β's R0 review commit, verdict
`converge`) | this round does not re-litigate β's R0 AC1–AC10 findings —
all ten remain independently verified and unchanged; this round repairs the
one thing β itself flagged (Finding 1) and δ subsequently overrode on.**

### What was overridden and why

δ's override comment on cnos#556 (posted after β's R0 `converge` verdict,
titled "δ override — R0 converge not accepted; R1 repair ordered") rejected
R0's `converge` verdict specifically on the axis β itself had flagged as
Finding 1 (a non-blocking "process/coherence" note in β's R0 review) but
did not block on. δ's own independent re-read of the operator's clarifying
comment on this issue (posted 2026-07-02T22:50:19Z, the same comment that
released this cell to `status:todo`) found it states, verbatim, under a
"## Go implementation rule" heading:

> Do not force Go implementation code into `src/packages/`. The active Go
> implementation may remain under `src/go/internal/issuesmap/` during the
> shim phase, but it must be treated as the implementation of the
> `cnos.issues` domain.

R0 (this file's `## Gap` through `## Review-readiness` sections above, all
authored in R0 and left unmodified by this append) implemented the
opposite: γ's R0 scaffold (`gamma-scaffold.md` §6 item 2) pinned physical
relocation of `src/go/internal/issuesmap/{issuesmap.go,fetch.go,
issuesmap_test.go,templates/,testdata/}` into
`src/packages/cnos.issues/commands/issues-map/` as a *binding* guardrail,
framed in the scaffold as "restated from the operator's clarifying
comment" — but on direct re-read (both by β in R0 Finding 1, and
independently by δ in the override comment), that framing was inaccurate:
the relocation was γ's own architectural judgment call (mirroring the
cnos#392 `cdd-verify` precedent), overriding the operator's specific,
on-topic, most-recent sentence rather than restating it. δ's override
holds that an explicit operator instruction on the exact axis in question
supersedes precedent-matching, and directs this R1 repair.

### Repair performed

1. **Reverted the physical code relocation.** `issuesmap.go`, `fetch.go`,
   `issuesmap_test.go`, `templates/board.template.html`, and
   `testdata/issues.json` are back at `src/go/internal/issuesmap/` (their
   pre-R0 location). Done via `git revert --no-edit` of the two R0 commits
   that performed and then referenced the relocation
   (`88156120` "issues-map: relocate domain implementation to cnos.issues
   package boundary" and `13198122` "docs: update board README path
   reference to cnos.issues package location"), reverted in reverse
   chronological order (`13198122` first, then `88156120`). Git recognized
   both file renames cleanly (`rename src/{packages/cnos.issues/commands/
   issues-map => go/internal/issuesmap}/*.go` etc.) — no manual conflict
   resolution was needed, and the revert is byte-for-byte equivalent to
   the pre-R0 file contents (verified via `git show 2ae24f27:...` diffed
   against the post-revert tree).
2. **Removed the now-unnecessary Go module wiring**, restored to its
   pre-R0 state by the same reverts: `src/packages/cnos.issues/commands/
   issues-map/go.mod` no longer exists (the `commands/` directory is gone
   entirely — confirmed via `find src/packages/cnos.issues -maxdepth 2
   -type d`, which now lists only `skills/{map,taxonomy,triage}`); the
   `require`+`replace` block for `issues-map` is gone from
   `src/go/go.mod`; `go.work`'s `use (...)` list no longer includes
   `./src/packages/cnos.issues/commands/issues-map` (only `./src/go` and
   `./src/packages/cnos.cdd/commands/cdd-verify` remain, exactly as
   pre-R0).
3. **Reverted `src/go/internal/cli/cmd_issues_map.go`'s import path** back
   to `github.com/usurobor/cnos/src/go/internal/issuesmap` (the same
   revert restored this). `Run()` remains the one-line delegation
   `return issuesmap.Run(ctx, inv.Args, inv.Stdin, inv.Stdout,
   inv.Stderr)` — still thin, no domain logic added or removed.
4. **Kept the `cnos.issues` doctrine/skill surfaces** —
   `cn.package.json`, `SKILL.md`, `skills/map/SKILL.md`,
   `skills/taxonomy/SKILL.md`, `skills/triage/SKILL.md` — all still
   present, per δ's explicit instruction that these are not part of the
   repair. **Updated `SKILL.md` and `skills/map/SKILL.md`** (commit
   `aa635c7a`) to state plainly that the Go implementation lives at
   `src/go/internal/issuesmap/`, not under `src/packages/cnos.issues/`,
   quoting the operator's "Go implementation rule" verbatim as the
   grounding, and naming the R0-relocation-then-R1-reversion history so a
   future reader does not mistake the current doctrine text for an
   unfinished migration. `skills/taxonomy/SKILL.md` and
   `skills/triage/SKILL.md` contained no path references to begin with
   (confirmed via grep) and needed no change.
5. **Reverted `docs/development/board/README.md:42`** back to
   `` The generator is the Go command `cn issues map`
   (`src/go/internal/issuesmap/`). `` — restored by the same revert of
   `13198122`.
6. **Re-grepped the whole repo** for `internal/issuesmap` and
   `cnos.issues/commands/issues-map` after all edits landed
   (`grep -rn "internal/issuesmap\|cnos.issues/commands/issues-map"
   src/ docs/ go.work`). Result: every live-code/doc reference now points
   at `src/go/internal/issuesmap/` (the import in `cmd_issues_map.go`,
   the embedded README-generator string inside `issuesmap.go` itself, the
   board `README.md` line, and the doctrine prose in the two updated
   `SKILL.md` files). The only remaining occurrences of
   `cnos.issues/commands/issues-map` are inside the two updated
   `SKILL.md` files' own prose, explicitly narrating the R0-relocation-
   then-R1-reversion history in the past tense — not live references to
   an existing path. No dangling reference to the removed module path
   exists anywhere outside `.cdd/unreleased/556/{gamma-scaffold,
   beta-review}.md`, which are R0-round historical artifacts left
   unmodified per this round's own instructions (append-only, do not
   rewrite prior rounds' artifacts).

### Build/test/CI re-verification

- `cd src/go && go build ./... && go vet ./... && go test ./...` — all
  clean; `go test ./...` reports `ok` for all 14 `src/go` packages
  including `github.com/usurobor/cnos/src/go/internal/issuesmap` (the
  restored package, tests unchanged and passing at their restored
  location).
- Built `cn` from `src/go` (`go build -o /tmp/cn_r1 ./cmd/cn`) and ran,
  myself: `/tmp/cn_r1 issues map --fixture internal/issuesmap/testdata/
  issues.json --out /tmp/board-out-r1` → `wrote 5 open issues to
  /tmp/board-out-r1 (index.html, board-data.json, README.md)`, exit 0.
  Also ran the exact form named in δ's repair contract against the live
  repo: `/tmp/cn_r1 issues map --repo usurobor/cnos --out
  /tmp/board-out-r1-live` → `wrote 93 open issues to /tmp/board-out-r1-live
  (index.html, board-data.json, README.md)`, exit 0.
- `/tmp/cn_r1 build --check` → `✓ cnos.issues: valid` alongside all other
  packages, `✓ All packages valid.` — `cnos.issues` is recognized as a
  doctrine-only package: `cn.package.json` has no `commands` object entry
  (confirmed: the file has no `commands` key at all, matching its pre-R0/
  R0 shape), and there is no `commands/` directory under
  `src/packages/cnos.issues/` at all after the revert.
- Remote CI on the R1 head commit (`aa635c7a80454e3c9dcb4793dda6d20ddd76e5dc`,
  pushed to `origin/cycle/556`): run
  [28628522058](https://github.com/usurobor/cnos/actions/runs/28628522058),
  `status: completed`, `conclusion: success`. All 10 required jobs green:
  `Go build & test`, `Binary verification`, `Package verification`,
  `Package/source drift (I1)`, `Protocol contract schema sync (I2)`,
  `Repo link validation (I4)` (validates the relative links added in the
  updated `SKILL.md` files resolve correctly — confirmed locally before
  push via direct path resolution, then confirmed green in CI),
  `SKILL.md frontmatter validation (I5)` (validates the two edited
  `SKILL.md` files' frontmatter is still well-formed — body-only edits,
  frontmatter untouched), `CDD artifact ledger validation (I6)`,
  `Dispatch repair-preflight guard (cnos#516)`, `Dispatch closeout-
  integrity guard (cnos#524)`.

### Commits this round

- `97a4b7d3` — `Revert "docs: update board README path reference to
  cnos.issues package location"`
- `f9707a9d` — `Revert "issues-map: relocate domain implementation to
  cnos.issues package boundary"`
- `aa635c7a` — `cnos.issues: update doctrine to reflect Go impl at
  src/go/internal/issuesmap/ (R1)`

### Review-readiness (R1)

**Round 1 | base SHA: `bdad537f5e5a719d6f76374da282eb6764fd014b` (β's R0
review commit) | implementation SHA: `aa635c7a80454e3c9dcb4793dda6d20ddd76e5dc`
| branch CI: green — remote run
[28628522058](https://github.com/usurobor/cnos/actions/runs/28628522058)
succeeded on the R1 head commit, all 10 required jobs `success` | ready for
β re-review.**

This repair does not reopen AC1–AC9 as previously verified by β in R0 (all
unaffected by this round — no board-output behavior, taxonomy, CLI dispatch
shape, or GitHub Action changed). It specifically repairs the disposition
underlying **AC10** (and β's R0 Finding 1): the receipt now honestly and
consistently states that the Go implementation lives at
`src/go/internal/issuesmap/`, not relocated into `src/packages/
cnos.issues/`, per the operator's explicit instruction, while
`cnos.issues` still exists as the domain's doctrine/skill home (AC1, AC2
unaffected). β should independently re-verify: (a) the file layout
(`find src/go/internal/issuesmap src/packages/cnos.issues -type f`), (b)
`go build/vet/test` from `src/go`, (c) `cn build --check` still recognizes
`cnos.issues` with no `commands` entry, (d) `cn issues map --repo
usurobor/cnos --out <tmpdir>` still works, (e) the repo-wide grep for
`internal/issuesmap` / `cnos.issues/commands/issues-map` shows no dangling
references, and (f) CI green on `aa635c7a` via `gh run list --branch
cycle/556`.

## §R2

**Repair round 2 | base: R1 rejected state `7cbd07b7` (branch tip at operator's `status:changes`) | repair SHA: `e2017b3f` | ready for β re-review.**

### What was rejected and why

The operator reopened #556 with `status:changes` (issue comment, 2026-07-03,
titled "Operator correction — #556 reopened"). The rejection: R1's repair
(triggered by a δ override) reverted R0's already-converged relocation of
the Go implementation into `src/packages/cnos.issues/commands/issues-map/`,
on the strength of an earlier operator comment ("the Go implementation may
remain under `src/go/internal/issuesmap/` during the shim phase"). The
operator's correction explicitly **withdraws** that earlier comment: *"That
was a misunderstanding of intent and is now withdrawn. The actual goal is a
clean package/kernel separation, including the Go code... Do not revert the
relocation this time — move the Go implementation into the package and keep
it there."*

### Repair taken (see `.cdd/unreleased/556/REPAIR-PLAN.md` for the full
finding→repair map)

1. `git revert --no-edit f9707a9d 97a4b7d3` (commits `df4bfd8b`, `4d9695f8`)
   — cleanly reinstates R0's relocation: `issuesmap.go`, `fetch.go`,
   `issuesmap_test.go`, `templates/`, `testdata/` move back to
   `src/packages/cnos.issues/commands/issues-map/` (git rename-detected, 0
   semantic diff on the moved files besides `issuesmap.go`'s embedded
   README string / doc-comment, matching R0's original diff exactly); the
   module's own `go.mod` is recreated; `go.work` regains the `use` entry;
   `src/go/go.mod` regains the `require`+`replace` pair; `cmd_issues_map.go`
   is back to importing the package path as a one-line thin shim;
   `docs/development/board/README.md`'s source-of-truth line points at the
   package path again.
2. Rewrote `src/packages/cnos.issues/SKILL.md` and
   `src/packages/cnos.issues/skills/map/SKILL.md` fresh (commit `e2017b3f`)
   rather than re-reverting R1's doctrine commit (`aa635c7a`, which was not
   itself part of the two reverted commits and so remained live, describing
   the now-false `src/go/internal/issuesmap/` location). New text states
   current truth only: implementation at
   `src/packages/cnos.issues/commands/issues-map/`, Go-source co-location
   mirroring the `cnos#392` `cdd-verify` precedent, kernel-dispatch thin
   shim (Option B, `#216` debt) — no R0/R1 round narration in the active
   skill prose, per the operator's "teach current truth" instruction (item
   5 of the repair contract).

### Build/test/CI re-verification

- `cd src/go && go build ./... && go vet ./... && go test ./...` — clean;
  13 packages `ok` (the `issuesmap` package is no longer part of this
  module — see next line).
- `cd src/packages/cnos.issues/commands/issues-map && go build ./... && go
  vet ./... && go test ./...` — clean; `ok` for the relocated module,
  tests unchanged and passing at the new location.
- Built `cn` from `src/go` (`go build -o /tmp/cn_r2 ./cmd/cn`) and ran,
  myself: `/tmp/cn_r2 build --check` → `✓ cnos.issues: valid` alongside
  all other packages, `✓ All packages valid.`
- `/tmp/cn_r2 issues map --fixture src/packages/cnos.issues/commands/
  issues-map/testdata/issues.json --out /tmp/board-out-r2` → `wrote 5 open
  issues to /tmp/board-out-r2 (index.html, board-data.json, README.md)`,
  exit 0.
- `/tmp/cn_r2 issues map --repo usurobor/cnos --out /tmp/board-out-r2-live`
  → `wrote 94 open issues to /tmp/board-out-r2-live (index.html,
  board-data.json, README.md)`, exit 0 (issue count differs from R1's 93
  only because new issues were filed on the live repo between rounds —
  expected, not a regression).
- Repo-wide grep `grep -rln "internal/issuesmap\|cnos.issues/commands/
  issues-map" . --exclude-dir=.git --exclude-dir=.cdd` → exactly 8 files,
  all pointing at the correct **new** `cnos.issues/commands/issues-map`
  path (`go.work`, `docs/development/board/README.md`, `src/go/go.mod`,
  `src/go/internal/cli/cmd_issues_map.go`, `src/packages/cnos.issues/
  SKILL.md`, `src/packages/cnos.issues/commands/issues-map/{issuesmap.go,
  go.mod}`, `src/packages/cnos.issues/skills/map/SKILL.md`). No dangling
  reference to the old `src/go/internal/issuesmap/` path remains anywhere
  live.
- `git diff --stat` from the cycle's true base (`4fe8e433`, the commit PR
  #557 was originally opened against) to this round's HEAD (`e2017b3f`)
  touches exactly the files the repair contract predicts: the package
  manifest/doctrine additions, the five relocated Go-domain files, `go.work`,
  `src/go/go.mod`, the recreated `commands/issues-map/go.mod`,
  `cmd_issues_map.go`, the board README path line, and this round's own
  `.cdd/unreleased/556/*.md` artifacts. `.github/workflows/board-map.yml`
  is untouched — confirmed via `git diff 4fe8e433...HEAD --
  .github/workflows/board-map.yml` (empty).

### Commits this round

- `35562a92` — `delta: REPAIR-PLAN for R2 repair round (cnos#556)`
- `df4bfd8b` — `Reapply "issues-map: relocate domain implementation to
  cnos.issues package boundary"` (revert of R1's `f9707a9d`)
- `4d9695f8` — `Reapply "docs: update board README path reference to
  cnos.issues package location"` (revert of R1's `97a4b7d3`)
- `e2017b3f` — `cnos.issues: fix doctrine to state current truth (R2)`

### repair_evidence

```yaml
repair_evidence:
  prior_rejection: "https://github.com/usurobor/cnos/issues/556 — operator status:changes comment, 2026-07-03, 'Operator correction — #556 reopened'"
  repairs_required:
    - finding-1: "Go implementation must physically live under src/packages/cnos.issues/commands/issues-map/, not src/go/internal/issuesmap/ — R1's revert was wrong relative to true operator intent"
    - finding-2: "cmd_issues_map.go must remain a thin shim over the package-owned implementation"
    - finding-3: "cnos.issues doctrine (SKILL.md, skills/map/SKILL.md) must state current truth (package-relocated), not the withdrawn src/go/internal/issuesmap/ narrative, and without R0/R1 narration in active prose"
    - finding-4: "board Action and board output behavior must remain unchanged"
    - finding-5: "cn build --check, go test ./..., and all required CI jobs must pass at the new HEAD"
    - finding-6: "receipt must explicitly state the package-command dispatch disposition (Option A vs B)"
  repairs_completed:
    - finding-1: "git revert of f9707a9d/97a4b7d3 (commits df4bfd8b, 4d9695f8) reinstates R0's relocation byte-for-byte; verified via file-layout find + go.work/go.mod diff"
    - finding-2: "cmd_issues_map.go re-verified as one-line import + delegation (see file contents in commit df4bfd8b)"
    - finding-3: "commit e2017b3f rewrites both SKILL.md files fresh, stating package-relocated truth only, no R0/R1 narration"
    - finding-4: "git diff 4fe8e433...HEAD -- .github/workflows/board-map.yml is empty; board-data.json/index.html generation logic unchanged (only the moved Go source, 0 semantic diff besides embedded string)"
    - finding-5: "local go build/vet/test green in both modules; cn build --check green; cn issues map fixture+live runs succeed; remote CI to be confirmed after push (see beta-review.md §R2)"
    - finding-6: "SKILL.md 'Command-dispatch disposition' section states Option B explicitly: package-owned implementation, kernel-dispatch thin shim, true package-command exec-dispatch remains #216 debt"
  repairs_not_completed: []
  delta_overrides: []
  new_state_differs_from_rejected: "git diff 7cbd07b7..e2017b3f --stat: 13 files changed, 114 insertions(+), 23 deletions(-) — the rejected R1 state had the Go implementation at src/go/internal/issuesmap/ with doctrine describing that location; the repaired state has it at src/packages/cnos.issues/commands/issues-map/ with doctrine describing that location, plus the go.work/go.mod wiring and REPAIR-PLAN.md"
```

### Review-readiness (R2)

This repair does not reopen AC1–AC9 as previously verified by β at R0
(board output, taxonomy, CLI dispatch shape, GitHub Action are unaffected —
this round is the physical-location repair only). It repairs AC10 and the
operator's explicit correction: the receipt now honestly states the Go
implementation lives at `src/packages/cnos.issues/commands/issues-map/`,
matching the operator's final, most-recent, on-topic instruction. β should
independently re-verify: (a) the file layout, (b) go build/vet/test in both
modules, (c) `cn build --check`, (d) `cn issues map` fixture + live, (e) the
repo-wide dangling-reference grep, (f) the diffstat footprint against the
true cycle base, and (g) CI green on the pushed head once available.
