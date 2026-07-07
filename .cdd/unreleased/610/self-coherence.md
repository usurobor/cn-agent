# α self-coherence — cnos#610

## §Gap

**Issue:** [cnos#610](https://github.com/usurobor/cnos/issues/610) — "cds-install Sub 3:
`cn repo install --dispatch cds` (dispatch layer)". Mode (issue header): `design-and-build`.
CDD cycle-sizing mode (per γ-scaffold): **substantial**.

**Gap (from the issue's Problem/Expected/Divergence):** no way to render a repo's dispatch
workflow through the installer; `--dispatch cds` unconditionally refused (Mock B4 guard,
pending #609). #609 (renderer generalization) merged via PR #619, so the guard is now
stale. Separately, the `cds-dispatch` wake's rendered prompt body still carries two
tenant-visible hardcoded-sigma leaks, confusing for a real third-party render (operator
#610 directive, folded into AC5 and the issue's Scope).

**What this cycle closes:** wires `cn repo install --dispatch cds` to the (now-merged)
#609 renderer with explicit caller identity, detects (but does not implement) the
cnos#493 canonical-label-install gap and reports it actionably, states the
`workflow`-scope PAT + never-pushes-main facts, and de-sigmas the two rendered-prompt-body
leaks while preserving `--agent sigma` operational compatibility.

**Depends on:** #608 (base installer, merged), #609 (renderer generalization, merged via
PR #619), cnos#493 (canonical label install, **still OPEN** — confirmed via
`gh issue view 493 --repo usurobor/cnos --json state` at authoring time; this cell does
**not** implement it, per Non-goals).

**Branch:** `cycle/610`. HEAD at review-readiness time: `f334807` (implementation SHA;
rebased onto `origin/main` at `f7e9aaa` — see §Review-readiness for the rebase record).

## §Skills

- **Tier 1:** `CDD.md` (canonical lifecycle/role contract) + `cnos.cdd/skills/cdd/alpha/SKILL.md`
  (this file's own role surface, followed for load order, peer enumeration, self-coherence
  discipline, and the pre-review gate).
- **Tier 2:** none newly loaded beyond the always-applicable `eng/*` the α load order names
  generically (Go is the only language touched besides one markdown prose edit; no new
  Tier 2 bundle was required beyond what the existing `repoinstall`/`cli` packages already
  assume).
- **Tier 3 (issue-specific, named by γ's scaffold):**
  - `cnos.core/skills/write/SKILL.md` — applied to the two prose edits in
    `cds-dispatch/SKILL.md` (front-load the point, one fix per leak, no filler) and to the
    help-text rewrite in `cmd_repo_install.go` (state the workflow-scope PAT fact once,
    point to it rather than repeating).
  - `docs/development/design/cn-repo-install-MOCKS.md` (Mock C + Mock E) — read in full;
    Mock C's transcript/invariants (C1–C6) and Mock E's tenant-portability invariants
    (E1–E4) are the design source of truth this cell's oracle table maps to.
  - `.cdd/unreleased/610/gamma-scaffold.md` — read in full before coding; its peer
    enumeration (confirming the #609 renderer's tenant-portable finalizer/scanner paths
    already exist; confirming cnos#493's mechanism is genuinely absent via `find`/`rg`) was
    reused rather than re-derived, per the scaffold's own instruction.

## §ACs

Implementation commits: `0b717f4` (Go wiring), `f8fd1c4` (tests), `f334807` (prose +
golden re-render). All command output below was captured against the pre-rebase tree
(content identical post-rebase — rebase touched no line this cycle wrote; confirmed by
re-running the full suite + a fresh manual e2e pass after rebase, see §Review-readiness).

### AC1 — dispatch render (Mock C1/C3)

**Oracle (issue):** `--dispatch cds …` writes exactly `.github/workflows/cnos-cds-dispatch.yml`
after a successful base install; base artifacts present.

**Automated evidence:** `go test ./src/go/internal/repoinstall/... -run TestRun_DispatchCds_RendersWorkflow_ThenSurfacesLabelGap -v`
→ PASS. Asserts `.cn/deps.json` + `.cn/deps.lock.json` exist, then
`.github/workflows/cnos-cds-dispatch.yml` exists at exactly that path, using fixture
tarballs that vendor the REAL `cn-install-wake` script + REAL `cds-dispatch/SKILL.md`
(not synthetic stand-ins).

**Manual e2e evidence (real `cn` binary, real `cn build`-produced package index, scratch
git repo, no unit-test mocking):**

```
$ /tmp/cn-alpha repo install --index .../dist/packages/index.json --dispatch cds
...
✓ wrote .cn/deps.json
✓ wrote .cn/deps.lock.json (3 package(s))
✓ restored cnos.cdd@3.82.0 / cnos.cds@0.1.0 / cnos.core@3.82.0
cn-install-wake: cds-dispatch → /tmp/cn-alpha-scratch/.github/workflows/cnos-cds-dispatch.yml (rendered)
✓ rendered .github/workflows/cnos-cds-dispatch.yml
...
✗ canonical dispatch labels not ensured: cnos#493 ...
exit=1
```
`find .cn/deps.json .cn/deps.lock.json .github` → all three present.

**AC1/AC3 relationship (read this before AC3 below):** AC1's oracle text says "after a
successful run"; AC3's oracle says the dispatch path "returns a named, non-silent error."
Both are true of the SAME observed behavior: base install + render complete successfully
(AC1's claim), and only after that does `Run` propagate the still-open cnos#493 gap as its
returned error (AC3's claim) — not a partial/rolled-back render. I verified both AC1 and
AC3 with the one test above rather than treating them as requiring separate scenarios;
see the code comment on `TestRun_DispatchCds_RendersWorkflow_ThenSurfacesLabelGap` for the
same reasoning inline. This is a judgment call the issue text does not fully disambiguate
— flagged explicitly, not papered over (see §Self-check).

**Verdict: MATCH.**

### AC2 — explicit identity (Mock C2)

**Oracle:** missing `--workflow-pat-secret`/agent/bot → fail early, nonzero exit, no
partial workflow file. Negative: a partial render fails.

**Automated evidence:** `TestRun_DispatchCds_MissingIdentity_FailsEarlyNoPartialWrite`
(package-level) + `TestRepoInstall_DispatchCds_MissingIdentity_CliWiring` (CLI-level,
through the noun-verb dispatcher) — both PASS. Assert: nonzero exit, error names
`--workflow-pat-secret`, `os.IsNotExist` on `.github/` (not just the `.yml` file — the
whole directory), and base-install artifacts (`.cn/deps.json`) still exist (C1 layering:
base install always precedes the dispatch-specific gate).

**Manual e2e evidence:**
```
$ /tmp/cn-alpha repo install --index .../index.json --dispatch cds --agent acme
...
✗ --workflow-pat-secret is required for --agent "acme" (no default substrate PAT-secret
  binding for non-sigma agents); pass --workflow-pat-secret <NAME> naming the GitHub
  Actions secret that holds this agent's workflow-scoped PAT
exit=1
$ ls .github
ls: cannot access '.github': No such file or directory
```

**Verdict: MATCH.** Sigma default (no `--agent`) does NOT trip this gate —
`TestRun_DispatchCds_SigmaDefault_NoIdentityFlagsRequired` proves the identity check is
scoped to non-sigma agents only, matching the renderer's own existing default-PAT-binding
convention for `sigma`.

### AC3 — labels ensured

**Oracle:** after dispatch install, the canonical dispatch labels exist via cnos#493; if
that mechanism is unavailable, an actionable error is returned (not a silent skip).

**Evidence that the mechanism is genuinely absent (not assumed):**
```
$ rg -n "cds-dispatch|install-wake|cnos-cds-dispatch" src/go   # (γ scaffold's own check, re-verified)
$ find src/go/internal src/packages/cnos.core/commands -iname '*label*'
$ gh issue view 493 --repo usurobor/cnos --json state
{"state":"OPEN"}
```
No `cn install cnos.core` / label-doctor command exists anywhere in `src/go` or
`src/packages/cnos.core/commands/`. This cell does **not** implement cnos#493 (Non-goals);
`ensureCanonicalDispatchLabels()` in `repoinstall.go` is a small, explicitly-documented
function that always returns the named error today — a compile-time-honest stand-in for
"the mechanism does not exist yet," not a runtime probe of something that could ever
succeed until #493 ships.

**Automated + manual evidence:** the same test/e2e run cited under AC1 — `err.Error()`
contains `"cnos#493"`, stderr carries the same diagnostic, and the render (from AC1) is
NOT rolled back or left partial.

**Verdict: MATCH** (not a silent skip; named, actionable, cites the exact dependency).

### AC4 — no substrate-binding leak + PR-only (Mock C4/C6)

**Oracle:** a non-sigma render contains no hardcoded `SIGMA_WORKFLOW_PAT`, bot id
(`41898282`), bot name (`sigma@…`), or concurrency group (`cds-dispatch-sigma`); the
command never pushes `main` and states the `workflow`-scope PAT requirement.

**Automated evidence:** `TestRun_DispatchCds_RendersWorkflow_ThenSurfacesLabelGap` +
`TestRepoInstall_DispatchCds_IdentityFlagsWireThrough` (CLI-level, real `cn` binary
subprocess) both grep the produced file for all four leak tokens (0 hits) and assert
stdout contains `"needs \`workflow\` scope"` and `"never pushes to main"`.

**Manual e2e evidence (real binary, real render, agent=acme):**
```
$ grep -inE 'sigma|SIGMA_WORKFLOW_PAT|41898282' .github/workflows/cnos-cds-dispatch.yml
NO LEAKS FOUND
$ grep -rn "git push|exec.Command(\"git\"" src/go/internal/repoinstall/*.go
NONE (repoinstall never shells to git)
```
"Never pushes to main" holds **by construction**, not by a runtime guard: `repoinstall.go`
never imports `os/exec` for `git`, never calls any git subcommand, and the renderer script
it invokes (`cn-install-wake`) does not push either (confirmed by reading the script:
its only side effect is writing the `--out` file). Stdout states:
`"  This changes .github/workflows/ — the installing token needs \`workflow\` scope."` and
`"  Dispatch never pushes to main (PR-only)."`.

**Verdict: MATCH.**

### AC5 — tenant prose-clean render (operator #610 directive)

**Oracle:** for a non-sigma agent, the rendered workflow contains no `"today: sigma"`,
`agent-admin-sigma`, or `cn-sigma:` leaks; `--agent sigma` still renders correctly
(compatibility preserved); a tenant render carrying confusing sigma prose fails the grep
(non-vacuous oracle).

**Positive (non-sigma) — automated + manual:** both cited tests above grep the acme
render for `"today: \`sigma\`"`, `"agent-admin-sigma"`, `"cn-sigma:"` → 0 hits. Manual e2e:
```
$ grep -n "today: \`sigma\`\|agent-admin-sigma\|cn-sigma:" .github/workflows/cnos-cds-dispatch.yml
NO PROSE LEAKS FOUND
```

**Positive (sigma) — sha256 evidence:**
```
$ sha256sum .github/workflows/cnos-cds-dispatch.yml   # from a fresh scratch repo, --agent sigma (default)
822bb9ec7119d3fc027c5a4b521a5046663df15278ef05603bd51686fd8c8b91
$ sha256sum src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml
822bb9ec7119d3fc027c5a4b521a5046663df15278ef05603bd51686fd8c8b91
```
Matches — **this cycle's own golden**, not the pre-cycle golden (see §Self-check for the
explicit backward-compat deviation this required, and why).

**Negative (non-vacuous oracle) — automated:**
`TestDispatchRenderer_ProseLeakGrep_CatchesPreFixSigmaPhrasing` constructs a synthetic
dispatch-shaped `SKILL.md` fixture carrying the EXACT pre-cnos#610 hardcoded-sigma
phrasing, renders it for a non-sigma agent through the real renderer, and asserts all
three leak strings ARE found (proving the grep is a real detector). Manual reproduction:
```
$ grep -nE 'today: `sigma`|agent-admin-sigma|cn-sigma:' /tmp/old-leak-render.yml
87: ... today: `sigma` ...
89: ... agent-admin-sigma ... cn-sigma:.cn-sigma/logs/20260624.md ...
```

**Verdict: MATCH**, with an explicit, documented deviation from the pinned byte-identical
floor for the line-296 fix only (see §Self-check).
