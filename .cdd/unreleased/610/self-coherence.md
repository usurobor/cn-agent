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

## §Self-check

**Did α push ambiguity onto β?** Two genuine judgment calls were required where the issue
text and γ's scaffold under-determine the answer; both are written up in
`.cdd/unreleased/610/alpha-clarification-needed.md` with the exact reasoning, the
assumption made, and the narrow fix if the assumption is wrong — rather than silently
picking one reading and hoping β doesn't notice:

1. Whether AC3's "actionable error" means `Run()` returns non-nil (my reading) or a
   printed warning with a 0 exit. I implemented the stronger (nonzero-exit) reading, since
   it is the more defensible parse of "an actionable error is returned" and matches the
   Proof plan's parallel framing of "missing labels" alongside "missing identity" as
   negative cases.
2. Whether AC5's prose-clean requirement can override the pinned byte-identical backward-
   compat floor for the one leak (line 296's historical citation) that γ's own scaffold
   says cannot be templated honestly. I chose prose-clean, moving the specific citation to
   the non-rendered appendix; `--agent sigma` is byte-identical to *this cycle's own*
   golden (committed), not to the pre-cycle golden.

**Is every claim backed by evidence in the diff?** Yes — every AC row above cites either
an automated test name (all runnable via `go test ./src/go/internal/repoinstall/...
./src/go/internal/cli/...`) or a manual e2e transcript captured against the real `cn`
binary + a `cn build`-produced real package index + a scratch git repo (not against unit
fixtures alone). The AC5 negative case specifically constructs a before/after fixture
(the synthetic pre-fix-phrasing `SKILL.md`) rather than asserting "the grep would have
caught it" without running it — per the α role's explicit rule against exactly that
overclaim.

**Known-debt-relevant self-check:** `ensureCanonicalDispatchLabels()` is intentionally a
static, always-erroring function (not a runtime probe) because there is genuinely nothing
to probe yet — cnos#493 does not exist as any command, flag, or file in this repo. This is
disclosed as debt (§Debt) rather than dressed up as a "check."

## §Debt

1. **cnos#493 is unimplemented (by design/Non-goals).** `--dispatch cds` cannot fully
   succeed (exit 0) today — it always ends by naming the cnos#493 gap. This is the
   issue's own explicit contract (AC3), not an oversight, but it means the Mock C
   transcript's fully-green console experience is not yet achievable end-to-end. Tracked
   upstream at cnos#493 (open, P1); this cell does not own closing it.
2. **Two interpretive judgment calls** (AC1/AC3 relationship; AC5-vs-byte-identical
   precedence) are documented in `alpha-clarification-needed.md` rather than silently
   resolved. Both are reversible with a narrow, named code change if β/δ/operator disagree
   with the reading I chose.
3. **`--dry-run --dispatch cds` is new behavior, not covered by any AC.** I extended
   `printPlan` to state the pending render (previously it silently said "Dispatch: none"
   even when `--dispatch cds` was requested, which pre-dated this cycle as a latent
   dry-run inaccuracy masked by the old unconditional-failure guard). No AC names dry-run
   + cds explicitly; I verified it manually (see the dry-run e2e transcript in the
   `alpha-clarification-needed.md`-adjacent manual runs) but there is no dedicated
   automated test for this specific combination. Low risk: the change is additive and the
   `--dispatch none` dry-run branch is untouched (verified byte-identical via
   `TestRun_DryRun_WritesNothing`, unchanged).
4. **The renderer's identity/label vocabulary (agent/PAT-secret/bot-name/bot-id) is not
   validated beyond presence.** E.g. a malformed `--bot-id` (non-numeric) is passed through
   to the renderer as-is; the renderer itself does not validate its shape either. Out of
   scope per the pinned "reuse the #609 renderer" contract axis (no new validation logic
   invented on top of what the renderer already enforces).
5. **No test asserts the exact stdout identity/PAT-secret display lines
   (`identity: <agent>` / `pat secret: <name>`) beyond substring checks for the
   PAT-scope/never-pushes-main facts.** These display lines are informational, not
   contract-bearing (no AC names their exact wording), so I did not pin them with a golden
   string test; a future prose tweak to these two lines would not need a test update.

## §CDD Trace

**Step 6 — artifact enumeration (matches `git diff --stat origin/main..HEAD`):**

| File | Authored by | Class | Maps to |
|---|---|---|---|
| `.cdd/unreleased/610/CLAIM-REQUEST.yml` | γ (pre-existing on branch, not α) | dispatch marker | cycle bookkeeping, not an AC |
| `.cdd/unreleased/610/gamma-scaffold.md` | γ (pre-existing on branch, not α) | scaffold | α's work order (read, not authored) |
| `.cdd/unreleased/610/alpha-clarification-needed.md` | α | clarification note | §Self-check items 1–2 |
| `.cdd/unreleased/610/self-coherence.md` | α | self-coherence | this file |
| `.github/workflows/cnos-cds-dispatch.yml` | α | rendered substrate (live) | AC5 (re-rendered after the SKILL.md prose fix; kept in sync with the golden per `install-wake-golden.yml`'s sha256 check) |
| `src/go/internal/cli/cmd_repo_install.go` | α | code (CLI flag wiring + help text) | AC1/AC2/AC4 (flag pass-through); stale-help-text fix named in the dispatch prompt |
| `src/go/internal/cli/cmd_repo_install_test.go` | α | tests | AC1–AC5 (CLI-level: `TestRepoInstall_DispatchCds_RendererNotVendored_CliWiring`, `TestRepoInstall_DispatchCds_MissingIdentity_CliWiring`, `TestRepoInstall_DispatchCds_IdentityFlagsWireThrough`) |
| `src/go/internal/repoinstall/repoinstall.go` | α | code (domain logic) | AC1/AC2/AC3/AC4 (`runDispatchCds`, `ensureCanonicalDispatchLabels`, `resolveDispatchAgent`, `dispatchWorkflowPath`, `validateDispatch` rewrite, `printPlan` dispatch-awareness) |
| `src/go/internal/repoinstall/repoinstall_test.go` | α | tests | AC1–AC5 (package-level: `TestRun_DispatchCds_RendersWorkflow_ThenSurfacesLabelGap`, `TestRun_DispatchCds_SigmaDefault_NoIdentityFlagsRequired`, `TestRun_DispatchCds_MissingIdentity_FailsEarlyNoPartialWrite`, `TestRun_DispatchCds_RendererNotVendored_FailsWithNoPartialWrite`, `TestDispatchRenderer_ProseLeakGrep_CatchesPreFixSigmaPhrasing`, rewritten `TestValidateDispatch`) |
| `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` | α | prose | AC5 (the two named leaks) |
| `src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml` | α | golden re-render | AC5 (kept in sync with the SKILL.md prose edit) |
| `docs/guides/INSTALL-CDS.md` | α | peer doc fix | §2.3 peer enumeration (not in the δ-pinned package-scoping list; a third sibling surface — beyond `cmd_repo_install.go`'s help text — stating the same stale "`--dispatch cds` fails until #609" fact; grepped for the phrase repo-wide and fixed the one hit outside `.cdd/unreleased/608/` and `.cdd/unreleased/610/`, which are historical/self-authored records, not live docs) |

Every file in the diff is accounted for above; none is unmentioned. **Implementation-
contract note:** `docs/guides/INSTALL-CDS.md` is not named in δ's pinned "Package scoping"
row. I judged this an in-bounds peer-enumeration fix (alpha/SKILL.md §2.3 — a sibling doc
surface describing the same mechanism, not an architectural/package-scoping change to the
code) rather than a contract violation; flagging it explicitly per rule 3.6 rather than
silently expanding scope.

**Caller-path trace for new functions (repoinstall.go):**

- `runDispatchCds` — called from `Run` (repoinstall.go), itself called from
  `RepoInstallCmd.Run` (cmd_repo_install.go), the real `cn repo install` entrypoint. Not a
  dead function; exercised by every dispatch test above plus the manual e2e runs.
- `ensureCanonicalDispatchLabels` — called from `runDispatchCds`.
- `resolveDispatchAgent` — called from both `Run` (for the post-render "Dispatch: cds
  (agent: …)" stdout line — currently unreached on the error path, see §Debt item 1 — and
  from `runDispatchCds` (the actual identity-gate + renderer-arg resolution)).
- `dispatchWorkflowPath` — called from `runDispatchCds` to build the fixed `--out` path.

**Steps 0–5 (design/plan/tests/code/docs):** design = γ's scaffold (already-converged Mock
C/E design surface; no separate α design artifact required — γ's scaffold explicitly
states the design surface is stable). Plan = the γ scaffold's "Surfaces α is expected to
touch" section, followed directly. Tests/code/docs authored in that order per commit
history (`0b717f4` code, `f8fd1c4` tests were actually committed test-after-code — noted
as a sequencing deviation from the canonical tests-before-code order; both were written
together in practice and verified passing before either commit landed, so no code shipped
untested at any point on the branch).

## §Review-readiness

| Pre-review gate row (alpha/SKILL.md §2.6) | Status |
|---|---|
| 1. Cycle branch rebased onto current `origin/main` | **Done, re-verified at signal time.** Base at dispatch: `48d561e` (γ's HEAD). `origin/main` advanced to `f7e9aaa` (one unrelated `board-map` regen commit, no overlap with any file this cell touches — confirmed via `git diff --stat 922cc5c..origin/main` against the touched paths, 0 hits) while α worked; rebased via `git rebase origin/main` (clean, no conflicts) and force-pushed with `--force-with-lease`. Re-checked immediately before this section: `git merge-base --is-ancestor origin/main HEAD` → true, at `origin/main = f7e9aaad34793dbea80c603e315e0ecc0760fdfa`, observed 2026-07-07 ~02:30 UTC. |
| 2. CDD Trace through step 7 | Done (§CDD Trace above). |
| 3. Tests present | Done — 26 tests in `internal/repoinstall`, 39 in `internal/cli`, all passing (§ACs + runner output below). |
| 4. Every AC has evidence | Done (§ACs: AC1–AC5, each with automated test name + manual e2e transcript). |
| 5. Known debt explicit | Done (§Debt, 5 items). |
| 6. Schema/shape audit | `.cn/deps.json` (`cn.deps.v1`) / `.cn/deps.lock.json` (`cn.lock.v2`) schemas untouched — no field added, no writer changed; `pkg.Manifest`/`pkg.Lockfile` types untouched in this diff (confirmed: `git diff origin/main..HEAD -- src/go/internal/pkg/` is empty). |
| 7. Peer enumeration | Done — role-skill/lifecycle-skill split not applicable (no role skills touched); CI-workflow-file-comment class checked (`install-wake-golden.yml` grepped for the stale "#609"/"fails explicitly" phrasing — 0 hits, that file was never wrong); doc-surface peer enumeration found and fixed `docs/guides/INSTALL-CDS.md` (a third sibling surface beyond the help text, see §CDD Trace note). |
| 8. Harness audit | The one non-Go harness this cell's contract touches is the renderer script `cn-install-wake` (shell) — NOT modified in this diff (confirmed: `git diff origin/main..HEAD -- src/packages/cnos.core/commands/install-wake/` is empty); α only invokes it via `os/exec`, per the pinned "reuse the renderer" contract axis. The golden + live-workflow YAML outputs it produces are the only generated artifacts, both re-rendered and verified in sync (§ACs AC5). |
| 9. Polyglot re-audit | Languages in this diff: Go (`go build`/`go vet`/`go test` all clean, see below) + Markdown (`SKILL.md`, `INSTALL-CDS.md` — both read back for structure/cross-reference correctness) + YAML (golden + live workflow — both parsed with `python3 -c "import yaml; yaml.safe_load(...)"`, both pass; both re-verified in sync via sha256). No shell was touched, so no shell-specific audit applies. |
| 10. Branch CI green | **No CI runner available in this environment** (no GitHub Actions access from this sandbox). Stated explicitly, per row 10's fallback: β must wait for `install-wake-golden.yml` + the Go test workflow to go green on this branch's HEAD before merge. All checks that workflow runs were reproduced manually above (§ACs; the "Re-render both" / structural-shape / AC7-AC8-leak-audit / AC1-AC2/E2-E4 steps all pass locally against this HEAD). |
| 11. Artifact enumeration matches diff | Done (§CDD Trace step 6 table; 12 files, all named). |
| 12. Caller-path trace for new modules | Done (§CDD Trace: `runDispatchCds`, `ensureCanonicalDispatchLabels`, `resolveDispatchAgent`, `dispatchWorkflowPath`, each with a named non-test caller). |
| 13. Test assertion count from runner output | `go test ./internal/repoinstall/... ./internal/cli/... -v` → 26 + 39 = **65 test functions, 65 PASS, 0 FAIL** (full `--- PASS:` list captured; not manually enumerated). |
| 14. α git identity | `git log -1 --format='%ae' HEAD` → `alpha@cdd.cnos` (canonical elision form). Configured at session start before any commit; no retroactive rebase-for-identity was needed. |
| 15. γ-artifact-of-record presence | **§5.1 canonical dispatch** — `.cdd/unreleased/610/gamma-scaffold.md` present at the literal path on `origin/cycle/610` (confirmed: it was read in full before any code was written, and is unmodified by α). |

**Final verification commands (paste-able):**
```
cd src/go && go build ./... && go vet ./... && go test ./internal/repoinstall/... ./internal/cli/...
```
All green at HEAD `a5ed0f2` (implementation SHA — this section's own commit will advance
HEAD past it, per the SHA-citation convention in alpha/SKILL.md §2.6).

**Ready for β.**

