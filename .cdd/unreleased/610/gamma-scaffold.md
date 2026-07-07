# γ scaffold — cnos#610

## Issue

[cnos#610](https://github.com/usurobor/cnos/issues/610) — "cds-install Sub 3:
`cn repo install --dispatch cds` (dispatch layer)". Protocol `cds`. Refs #607
(master tracker). Depends on #608 (base installer, merged) + #609 (renderer
generalization, merged via PR #619) + cnos#493 (canonical label install,
**still OPEN, P1** — this cell does not implement it; AC3 names the actionable
error path for its absence).

Issue already passed the γ issue-quality gate (§2.4 of
`cnos.cdd/skills/cdd/gamma/SKILL.md`) before this scaffold was authored:
Problem/Impact/Source-of-truth/Scope/AC1–AC5/Proof-plan/Parity-requirement/
Non-goals/Closure are all present; ACs are numbered and independently
testable; non-goals exist; Tier 3 design surface
(`docs/development/design/cn-repo-install-MOCKS.md`) is linked; priority is
implicit in the wave sequencing (Sub 3 of the `cn repo install` wave, #607).
No prompt-only constraints are hiding outside the issue — the two issue
comments (Mock E3 scope extension; AC5 tenant-prose-clean addition) are both
already folded into the issue body's Scope/AC5 sections quoted above, so α
does not need the comment thread to get the full contract — the issue body
alone is complete and dispatchable.

## Mode

**substantial** (CDD cycle-sizing sense, distinct from the issue's own
`design-and-build` MCA-mode header field). Rationale against the five-factor
heuristic (`cdd/issue/SKILL.md` §"Cycle scope sizing"):

- 5 ACs (AC1–AC5), at the low end of "typical" band but touching a genuinely
  new CLI-flag surface, a cross-package renderer-integration point, and a
  prose-cleanup pass — the AC-count-alone heuristic undercounts this cell's
  difficulty the same way cnos-tsc cycle 24 did (issue/SKILL.md §"Empirical
  anchor").
- (a) New code surface: no new module: **extends** `repoinstall.Args`/`Run`
  and `cmd_repo_install.go`'s flag surface. Not a split signal alone.
- (b) Cross-module breadth: touches `src/go/internal/repoinstall/`,
  `src/go/internal/cli/`, and `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md`
  (a different package tree entirely) — 3 module/package boundaries. This is
  the split-signal factor, but the three surfaces are **not independently
  shippable** (AC1/AC2/AC4 require the SKILL.md prose fix to pass AC5's
  no-sigma-leak oracle in the same render), so splitting would violate
  master+subs' "independent shippability" rule. Kept whole; the five-factor
  table's factor (e) reading is "no" — the CLI wiring and the prose cleanup
  are two dimensions of the *same* render, not separable deliverables.
- (c) Lifecycle span: single design→code→test phase; no separate design
  cycle needed (design surface already converged at Mock C/Mock E3 in
  `cn-repo-install-MOCKS.md`).
- (d) MCA-precondition stability: design (Mock C + Mock E3) is committed and
  stable at `docs/development/design/cn-repo-install-MOCKS.md`; the renderer
  it calls (#609/PR #619) is merged and stable. Issue mode is correctly
  `design-and-build` (Mock E3 is design-adjacent — deciding whether the
  finalizer-invokes-installed-`cn` behavior needs any *new* wiring here, since
  #609 already built the tenant-portable finalizer path into the renderer
  itself — see §Peer enumeration below).
- Decision: **keep whole**, per (e) non-independent-shippability.

## Peer enumeration (per gamma/SKILL.md §2.2a)

Before framing §Gap, checked for prior partial coverage:

- `find src/go/internal/repoinstall src/go/internal/cli -iname '*.go'` +
  `rg -n "Dispatch|dispatch" src/go/internal/repoinstall/*.go` (excluding
  `_test.go`): confirms `repoinstall.go` already has the `--dispatch none|cds`
  flag parsed and a `validateDispatch` stub that **unconditionally fails**
  `cds` with `"--dispatch cds requires generalized wake renderer support
  (#609)"` — this is exactly the Mock B4 guard the issue's Scope calls "base
  install internals (#608, merged)" and marks **out of scope** to re-litigate;
  α's job is to make `validateDispatch`/`Run` route to the (now-merged)
  renderer instead of unconditionally refusing.
- `rg -n "cds-dispatch|install-wake|cnos-cds-dispatch" src/go` found no Go
  code that already invokes the renderer — the wiring point named in AC1/AC2
  genuinely does not exist yet. Negation confirmed empirically (not just
  asserted).
- **Mock E3 partial-closure finding (additive framing):** `cn-install-wake`
  (the #609 renderer, `src/packages/cnos.core/commands/install-wake/cn-install-wake`)
  **already implements** the tenant-portable "Install cn" step + the
  installed-`cn`-invoking finalizer/scanner steps for `agent != "sigma"`
  (lines ~1035–1071 and ~1316–1346 of that script, citing "cnos#609, Mock
  E2/E4"). Mock E3's invariant ("finalizer/engine steps invoke the installed
  `cn`, runnable in a repo with no `src/go`") is **already satisfied by the
  renderer's existing output** for any non-sigma agent. This cell's Mock-E3
  obligation is therefore **not** "build the tenant-portable finalizer" (already
  built) — it is "prove it end-to-end through `cn repo install --dispatch
  cds`'s own call path" (the renderer has its own CI-level proof via
  `install-wake-golden.yml`'s E2/E4 steps; this cell's proof plan must show
  the *installer's* invocation of the renderer reaches the same tenant-portable
  branch, not re-derive the renderer logic). α should not re-implement
  anything the renderer already does; α's surface is the *plumbing* from
  `cn repo install --dispatch cds` into the existing renderer entrypoint.

## Surfaces α is expected to touch

1. **`src/go/internal/repoinstall/repoinstall.go`** (primary surface):
   - `Args`/`Options` gain new fields for the identity flags the renderer
     already accepts: `Agent`, `WorkflowPatSecret`, `BotName`, `BotID` (names
     chosen to mirror `cn-install-wake`'s own flag names 1:1 — do not invent
     new flag vocabulary).
   - `ParseArgs` parses `--agent`, `--workflow-pat-secret`, `--bot-name`,
     `--bot-id` (two-token `--flag value` form per existing convention in this
     file).
   - `validateDispatch` (or its call site in `Run`) stops unconditionally
     refusing `cds` — AC1/AC2 require it to attempt the render.
   - `Run` sequences: base install (`applyInstall`, unchanged) completes
     first (AC1/C1) → if `Dispatch == "cds"`, resolve identity (AC2: missing
     `--workflow-pat-secret` for a non-sigma `--agent` fails **before any
     write**, matching the renderer's own existing fail-early behavior at
     `cn-install-wake` lines ~695–707 — α should surface that failure early
     rather than shelling out and letting the renderer's own die() be the
     only gate, per Mock C2 "no partial render") → invoke the vendored
     renderer (`<RepoRoot>/.cn/vendor/packages/cnos.core/commands/install-wake/cn-install-wake`,
     which `applyInstall`'s `restore.Restore` call already placed on disk)
     via `os/exec`, targeting the `cds-dispatch` wake manifest that ships
     under the vendored `cnos.cds` package, with `--out
     .github/workflows/cnos-cds-dispatch.yml` and the resolved identity flags
     → ensure canonical dispatch labels (AC3: cnos#493's `cn install
     cnos.core` / label-doctor mechanism **does not exist on this branch's
     base** — confirmed OPEN via `gh issue view 493`; α's job is to detect
     its absence and return a named, actionable error, e.g. "canonical
     dispatch labels not ensured: cnos#493 label-install mechanism is not
     yet available; labels must be applied manually until it ships" — **not**
     to build the mechanism, and **not** to silently skip it).
   - State the `workflow`-scope PAT requirement + "never pushes to `main`"
     in the command's stdout/help text (AC4/C6 — these are asserted facts
     about the command's own behavior, not new mechanics: `cn repo install`
     already never touches git remotes anywhere in `applyInstall`, so this is
     a documentation/stdout-message obligation, not a new guard to build).

2. **`src/go/internal/cli/cmd_repo_install.go`**: extend `repoInstallHelp`
   with the new flags (`--agent`, `--workflow-pat-secret`, `--bot-name`,
   `--bot-id`) and correct the now-stale help text ("`cds` — NOT implemented
   until #609; fails explicitly" — #609 is merged, this line must change) and
   wire `Args` → `Options` pass-through for the four new fields.

3. **`src/go/internal/repoinstall/repoinstall_test.go`**: new tests for
   AC1–AC5 (see oracle table below); the existing dispatch-related tests
   (`TestRun_DispatchCds_FailsWithNoPartialWrite` at line ~567) assert the
   *old* unconditional-failure behavior and must be updated to the new
   contract (still covering the genuine failure path — missing identity,
   AC2 — but no longer asserting `cds` fails unconditionally).

4. **`src/go/internal/cli/cmd_repo_install_test.go`**:
   `TestRepoInstall_DispatchCds_CliWiring` (line 258) currently asserts
   `--dispatch cds` **always** fails and names `#609` in stderr — this
   assertion is now false (that guard is exactly what this cell removes) and
   must be rewritten to cover the CLI-level pass-through of the new flags
   instead.

5. **`src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md`** (AC5 —
   prose-clean render, the second dimension of this cell). Exactly two
   tenant-visible hardcoded-sigma spots exist in the **rendered prompt body**
   (i.e. before the `## Responsibilities (body reference)` heading at line
   354, which is the appendix boundary `cn-install-wake`'s `skill_body()`
   strips — confirmed by reading the extraction regex; everything after line
   354, e.g. the `## Agent variable` / `## Concurrency notes` sections, is
   NOT rendered into the prompt and needs no change):
   - **Line 101**: `"You substrate-execute as \`{agent}\` (today: \`sigma\`;
     future: per-package bot accounts per cnos#449 follow-up)."` — the
     `(today: sigma; ...)` parenthetical is the leak; the `{agent}` token
     itself is already correctly parameterized and must stay.
   - **Line 296** (inside `## Disallowed surfaces`): two leaks in one
     sentence — `"the admin wake's \`agent-admin-sigma\` concurrency group"`
     and the empirical-motivator citation `"the 2026-06-24 mixed log entries
     (\`cn-sigma:.cn-sigma/logs/20260624.md\`, ...)"`. Both name a concrete
     sigma-specific historical incident/identity that is confusing in a
     tenant's own rendered wake.
   - **Invariant α must preserve**: `--agent sigma` still renders
     byte-identical to the existing golden
     (`src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml`)
     — AC5's positive case, and also `install-wake-golden.yml`'s existing E4
     regression guard. Since the renderer's `sed "s/{agent}/${agent}/g"`
     substitution already makes `{agent}` render as `sigma` when
     `--agent sigma`, the cleanest fix converts the line-101 leak to use the
     `{agent}` token directly (dropping the `(today: sigma; ...)`
     parenthetical entirely, or rephrasing it agent-neutrally) — this
     preserves byte-identical sigma output automatically. The line-296 leaks
     are historical/incident-specific and cannot be `{agent}`-templated
     honestly; α's design choice (rephrase agent-neutrally in the rendered
     body vs. move the specific historical citation below the "body
     reference" appendix boundary) is α's to make, verified by AC5's
     negative-case oracle, not pre-decided here.
   - This is a `cnos.cds`-owned package file, not `src/go/` — the same cell,
     different package tree, per the "not independently shippable" call
     above.

## Per-AC oracle approach

| AC | Oracle (concrete verification) |
|---|---|
| AC1 | `go test ./src/go/internal/repoinstall/...` — new/updated test drives `Run` with `Dispatch: "cds"` + a fixture index carrying `cnos.core` (which vendors the renderer script) + `cnos.cds` (which vendors the `cds-dispatch` manifest), asserts `.github/workflows/cnos-cds-dispatch.yml` exists at exactly that path after a successful run, and that base-install artifacts (`.cn/deps.json`, `.cn/deps.lock.json`) are also present. |
| AC2 | Table-driven test: omit `--workflow-pat-secret` with a non-sigma `--agent` → nonzero exit, no `.github/workflows/` directory created at all (`os.IsNotExist` check, mirroring the existing `TestRepoInstall_DispatchCds_CliWiring` pattern at cmd_repo_install_test.go:272). Positive case: full identity flags supplied → success. |
| AC3 | Test asserts that with cnos#493's mechanism absent (true today — no `cn install cnos.core` / label-doctor command exists anywhere under `src/go` or `src/packages/cnos.core/commands/`, confirmed by `find`/`rg` above), the dispatch-install path returns a named, non-silent error surfacing the cnos#493 dependency — not a clean exit that silently skipped label creation. |
| AC4 | Reuse the leak-audit pattern already proven in `.github/workflows/install-wake-golden.yml`'s "AC2 (cnos#609, Mock C4)" steps (scoped `grep -inE 'sigma\|SIGMA_WORKFLOW_PAT\|41898282'` against a non-sigma render) — apply the same grep against the file `cn repo install --dispatch cds --agent acme ...` produces end-to-end (not just the renderer's own direct-invocation golden), plus a `git log`-equivalent / stdout assertion that no `git push` to `main` ever occurs and the PAT-scope statement appears in stdout. |
| AC5 | (a) Render for a non-sigma agent through the full `cn repo install --dispatch cds` path; grep the two named leak strings (`today: sigma` / `agent-admin-sigma` / literal `cn-sigma:`) and confirm absence in the *tenant* render. (b) Render with `--agent sigma`; sha256-compare against the existing committed golden (`cnos-cds-dispatch.golden.yml`) — must be byte-identical (positive case / backward-compat). (c) Negative case: a fixture SKILL.md carrying the old hardcoded-sigma prose (or the pre-fix source, in a before/after test) must fail the leak-grep, proving the oracle is real and not vacuously true. |

**Empirical anchor for the review-discipline this cell inherits:** the
existing `install-wake-golden.yml` CI job is the direct precedent for AC4's
oracle shape (scoped grep, not case-insensitive blanket grep, per its own
"AC2 (scoped grep)" step comment at line ~818 — a blanket `grep -i sigma`
would false-positive on the legitimate `agent-admin-sigma` prose in *other*
wakes' own descriptions of themselves). α should follow that same
scoped-token discipline for AC4; AC5's oracle is different in kind (it is
about *this* wake's own rendered prompt body, not cross-wake false
positives) and should name the exact two leak strings above, not a blanket
`sigma` grep, to avoid AC5 vacuously failing on the legitimate `{agent}` →
`sigma` substitution result elsewhere in a sigma-targeted render.

## Expected diff scope

Small-to-moderate: two Go files + two Go test files in
`src/go/internal/repoinstall/` + `src/go/internal/cli/`, plus one prose edit
(two locations) in `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md`.
No new package, no new binary, no schema change. Renderer golden
(`cnos-cds-dispatch.golden.yml`) will need a re-render if the SKILL.md prose
changes touch the rendered body — α must re-run
`./src/packages/cnos.core/commands/install-wake/cn-install-wake cds-dispatch`
and commit the regenerated golden alongside the SKILL.md edit, or
`install-wake-golden.yml`'s own CI will fail on drift.

## Dispatch shape

Standard γ→α→β sequential dispatch (`cnos.cds/skills/cds/CDS.md` §"Field 6"
— no actor collapse; this is substantive software work, not a
skill/docs-class cycle).

---

## α prompt

```text
You are α. Project: usurobor/cnos.
Load src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md and follow its load order.
Issue: gh issue view 610 --repo usurobor/cnos --json title,body,state,comments
Branch: cycle/610
Tier 3 skills: src/packages/cnos.core/skills/write/SKILL.md, docs/development/design/cn-repo-install-MOCKS.md (Mock C + Mock E), .cdd/unreleased/610/gamma-scaffold.md
```

## Implementation contract (pinned by δ; α MUST NOT improvise)

| Axis | Pinned value |
|---|---|
| Language | Go (matching the existing `repoinstall`/`cli` packages); the one non-Go touch is a prose edit to an existing `SKILL.md` (markdown), not new code in another language. |
| CLI integration target | `cn` subcommand — extend the existing `cn repo install` command (`cmd_repo_install.go` + `repoinstall` package). No new subcommand, no standalone binary. |
| Package scoping | `src/go/internal/repoinstall/repoinstall.go` (+ `repoinstall_test.go`) and `src/go/internal/cli/cmd_repo_install.go` (+ `cmd_repo_install_test.go`); prose surface at `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` (+ its golden `cnos-cds-dispatch.golden.yml` re-render). |
| Existing-binary disposition | Extend — `cn` remains one binary; `repo-install` remains one command. Do not fork a separate installer binary or a separate dispatch-render binary. |
| Runtime dependencies | None new. The renderer this cell wires into (`cn-install-wake`) already requires `jq` + `python3`+`pyyaml`; this cell adds no new runtime dependency — it invokes the existing vendored script via `os/exec`. |
| JSON/wire contract preservation | Preserve as-is. `.cn/deps.json` / `.cn/deps.lock.json` schema (`cn.deps.v1` / `cn.lock.v2`) is untouched; the existing `--dispatch none` output and base-install-only stdout/JSON shape must not change. `--dispatch cds`'s NEW output (the render success/failure messages) is additive, not a breaking change to any existing consumer. |
| Backward-compat invariant | `--agent sigma` render must remain byte-identical to the committed golden (AC5 positive case) — this is the hard backward-compat floor. Existing `--dispatch none` behavior, existing base-install diff shape (Mock A/B), and the existing `install-wake-golden.yml` CI job's assertions must all continue to pass unchanged. |

## β prompt

```text
You are β. Project: usurobor/cnos.
Load src/packages/cnos.cdd/skills/cdd/beta/SKILL.md and follow its load order.
Issue: gh issue view 610 --repo usurobor/cnos --json title,body,state,comments
Branch: cycle/610
```

β's independent AC-oracle walk (reference: the same design doc + issue the α
prompt cites — β reads the issue and `docs/development/design/cn-repo-install-MOCKS.md`
directly, not α's rationale):

1. **AC1/C1/C3** — run the new/updated `repoinstall` tests; independently
   confirm `.github/workflows/cnos-cds-dispatch.yml` is the *only* new file
   the dispatch path writes, and that base-install artifacts precede it.
2. **AC2/C2** — independently construct a missing-identity invocation and
   confirm fail-early with no partial `.github/workflows/` directory at all.
3. **AC3** — confirm cnos#493 is still open (`gh issue view 493 --repo
   usurobor/cnos --json state`) at review time; confirm the shipped error
   path names the dependency rather than silently no-op'ing.
4. **AC4/C4/C5/C6** — run the scoped-grep leak audit (mirroring
   `install-wake-golden.yml`'s own AC2/E2/E4 steps) against a fresh non-sigma
   render produced through `cn repo install --dispatch cds` itself (not just
   the renderer's direct-invocation golden); confirm no `git push`/`git
   remote` write path exists anywhere in the new code; confirm the
   `workflow`-scope PAT statement appears in stdout/help.
5. **AC5** — grep the rendered non-sigma prompt body for the two named leak
   strings (`today: sigma` parenthetical; `agent-admin-sigma` /
   `cn-sigma:...` citation); sha256-diff the `--agent sigma` render against
   the committed golden for byte-identical backward compat; confirm the
   golden file was actually re-committed if the SKILL.md prose changed (drift
   would otherwise be caught by `install-wake-golden.yml` CI, but β should
   not rely on CI alone to catch a `git diff --exit-code` failure that could
   also be checked directly in review).
6. **Rule 7 (implementation-contract coherence)** — confirm the diff matches
   every pinned axis above; particularly confirm no new package, no new
   binary, and no runtime-dependency addition crept in.
7. **Parity requirement** — confirm the closeout carries `mock_parity` rows
   for C1–C6 with `missed: 0`, plus the issue's own explicit AC5
   tenant-prose-clean row, per the issue's "Parity requirement" section.
