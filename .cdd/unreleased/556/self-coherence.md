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
