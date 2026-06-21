# β review — cycle/476

## R1 (2026-06-21 — base origin/main: fcc5cdb9a533ad86e67524bcf05a33d2b4592e8a — head origin/cycle/476: 854102ceccd68c00abd3b2c606425a17cbadb840)

### Verdict: REQUEST CHANGES

One D-severity finding (F1 below): the CI step α added to discharge AC6's "CI mechanism in place" + AC2's negative-case proof is **broken** by a shell pipe-status bug. The renderer behavior is correct; the guard that proves the behavior is correct fails on every push, producing a permanent false-positive ::error and red CI on the very workflow α introduced. AC6 binds the **CI mechanism**, not just the local oracle; a broken-on-arrival CI mechanism does not satisfy AC6.

The renderer code itself is clean. ACs 1, 2 (behavior), 3, 4, 5, 7, 8 all pass on re-grep of the artifacts. AC6 partially passes (golden + idempotence verified; CI mechanism broken). Implementation contract conforms on every axis. Authority split is preserved. The renderer is deterministic and idempotent (proven by hash-equal in worktree merge-test).

The fix is one line in `.github/workflows/install-wake-golden.yml` (`set -o pipefail`). No other findings.

### Contract Integrity

| Check | Result | Notes |
|---|---|---|
| Status truth preserved | yes | self-coherence's §Gap correctly classifies mode = design-and-build, accepts γ pins, declares "no version bump" per γ scaffold |
| Canonical sources/paths verified | yes | wake-provider/SKILL.md @origin/main exists and is byte-identical on cycle/476 |
| Scope/non-goals consistent | yes | claude-wake.yml, wake-provider.json, prompt.md all byte-identical on cycle/476 |
| Constraint strata consistent | yes | renderer source carries clear authority-split comments matching contract |
| Exceptions field-specific/reasoned | yes | F1-F6 friction notes in α §Debt are scoped, routed correctly to γ closeout |
| Path resolution base explicit | yes | renderer derives CN_PACKAGE_ROOT from $0 fallback; manifest-override resolves prompt path relative to manifest dir per wake-provider/SKILL.md §2.1 |
| Proof shape adequate | partial | local oracles pass; CI proof is broken (F1) |
| Cross-surface projections updated | yes | self-coherence §CDD Trace lists every changed file, every commit, every section commit SHA |
| No witness theater / false closure | NO | α's claim "CI mechanism in place" (AC6 row) is contradicted by the actual CI run failing on the AC2 negative-case step that α authored. This is the exact "rule 3.13c wiring claim" surface — α documents a CI proof that doesn't fire correctly. Captured as F1 below. |
| PR body matches branch files | yes | PR #477 carries the same 6 files as `git diff --name-status origin/main..origin/cycle/476` |
| γ artifacts present (gamma-scaffold.md) | yes | `.cdd/unreleased/476/gamma-scaffold.md` exists on origin/cycle/476 at blob `a7bc2637...` |

### Pre-merge gate

| # | Row | Result | Notes |
|---|---|---|---|
| 1 | Identity truth | pass | `git config --get user.email` = `beta@cdd.cnos`; merge-test worktree config set with `--worktree` flag |
| 2 | Canonical-skill freshness | pass | `git fetch --verbose origin main` → `origin/main` still at `fcc5cdb9a533ad86e67524bcf05a33d2b4592e8a` (matches session-start snapshot); no mid-cycle main drift |
| 3 | Non-destructive merge-test | **fail** | Merge tree builds cleanly (zero unmerged paths). Shell, JSON, YAML validators all pass. AC6 local idempotence + golden-fixture invariant: all three sha256 outputs (committed, render-1, render-2) equal `a912dd97d520b8660e4a53c2f306dad7e657b0c216bbc94f31e2958fdbd01c47`. Scope check (0 out-of-scope files), AC8 renderer-side grep (0), AC7 invariant (0 lines changed), Sub 2 declarations unchanged (0 lines changed) — all pass. **But: CI on the merge SHA (PR #477) carries the new `install-wake-golden / Re-render + diff per-package goldens` job at conclusion `failure` due to F1; this is the "CI mechanism in place" α claims for AC6.** |
| 4 | γ artifact completeness | pass | `gamma-scaffold.md` present on cycle/476 at canonical §5.1 path |

**Worktree note:** `/tmp/cnos-merge-476/wt` was created with `git config --worktree user.email beta@cdd.cnos` per β SKILL §"Pre-merge gate" row 1 discipline; the worktree was torn down after validation, no identity leak.

### CI status (rule 3.10)

| Workflow / Job | Conclusion on head 854102ce | Pre-existing on origin/main? | Verdict impact |
|---|---|---|---|
| `install-wake golden / Re-render + diff per-package goldens` | **failure** | NO — workflow added by this cycle | F1, D-severity, blocks merge |
| `Build / SKILL.md frontmatter validation (I5)` | failure | YES — already failing on origin/main@fcc5cdb9 (Build run 27888241115 job 82527033430) | pre-existing; tracked by cnos#473/#474/#475 per issue body cross-refs; not this cycle's concern |
| `Build / CDD artifact ledger validation (I6)` | failure | YES — same Build run on origin/main fails with `cn-cdd-verify: No such file or directory` (missing tool) | pre-existing; not this cycle's concern |
| `Build / Repo link validation (I4)` | failure | YES — same Build run on origin/main fails with stale link errors | pre-existing; not this cycle's concern |
| `Build / Go build & test` | success | — | green |
| `Build / Package verification` | success | — | green |
| `Build / Binary verification` | success | — | green |
| `Build / Package/source drift (I1)` | success | — | green |
| `Build / Protocol contract schema sync (I2)` | success | — | green |

The pre-existing I4/I5/I6 reds on origin/main are NOT in cycle/476's scope and do NOT block this cycle's merge per CDD doctrine (they pre-date the cycle, are tracked by named follow-up issues, and would block every cycle's merge regardless). What blocks merge is the *new* failure α introduced via the broken CI step.

### AC coverage

| AC | Status | Evidence (re-grepped from diff) | Notes |
|---|---|---|---|
| AC1 (form pin + rationale) | pass | γ scaffold §"Form choice" pins shell; α §Gap "Form-choice acknowledgment" accepts the pin with four structural reasons | no override |
| AC2 (consume `cn.wake-provider.v1`; reject malformed) | **pass behavior; fail CI proof** | Re-ran all 6 negative cases in merge-test worktree: every malformed input exits 2 with stderr naming the precise field (`required field "schema" missing`, `admin_only must be boolean (got string)`, `unsupported schema "cn.wake-provider.v999"`, `not valid JSON`, `required field "defer_path" missing`, `manifest not found: ...`); `--out` flag oracle passes (byte-identical with default-path golden). The CI guard that proves this behavior is broken (F1). |
| AC3 (renders to syntactically valid GHA workflow at γ-pinned path) | pass | YAML parses (`python3 -c "import yaml; yaml.safe_load(...)"` exit 0); all 6 structural greps return ≥ 1: `^on:`=1, `^permissions:`=1, `^concurrency:`=1, `^jobs:`=1, `anthropics/claude-code-action@v1`=1, `actions/checkout@v4`=1 |
| AC4 (admin-only boundary verbatim) | pass | Re-grepped golden: `MUST NOT execute cells under any circumstance`=1; `Cell execution`=1; `## Disallowed surfaces`=1; `## Wake termination`=1 (last heading of prompt.md → no silent truncation). α's widened oracle (extract `prompt: \|` block, strip 12-space prefix) gives char-count ratio 0.998 — well within 95-105%. The β-side widening is recorded per β SKILL Rule 6a (the issue's literal-grep AC oracle as written would falsely fail; α's widened oracle is the substantive measure). |
| AC5 (triggers + permissions match manifest) | pass | Per-trigger oracle: `schedule` → `on.schedule` with 4 cron entries `'8 * * * *'`, `'23 * * * *'`, `'38 * * * *'`, `'53 * * * *'`; `issues_opened_title_match` → `on.issues.types: [opened]` + `contains(github.event.issue.title, 'agent-admin')` in if-conditional. Per-permission oracle: all 4 manifest `permission_intent` entries → all 4 GHA `permissions:` keys (kebab-case correctly mapped: `pull_requests.write`→`pull-requests: write`; `id_token.write`→`id-token: write`). No extras; rendered `permissions:` block contains exactly 4 keys = manifest's 4-element array. |
| AC6 (golden + idempotence + CI) | **partial** | (a) Golden file exists at γ-pinned path. (b) Idempotence: rendered golden + render-1 + render-2 all `sha256sum`=`a912dd97d520b8660e4a53c2f306dad7e657b0c216bbc94f31e2958fdbd01c47`; `diff` exit 0; renderer reports `(unchanged)` on re-render so file-mtime is also preserved. (c) **CI mechanism: the workflow file exists but the `AC2 negative-case smoke` step fails on every push due to F1**, so the CI mechanism does not actually defend the AC. AC6 binds the CI mechanism, not just the local oracle — partial. |
| AC7 (`claude-wake.yml` byte-identical) | pass | `git diff origin/main HEAD -- .github/workflows/claude-wake.yml \| wc -l` = 0 |
| AC8 (package-vs-renderer authority split preserved) | pass | Renderer-side: `grep -ciE 'admin.only\|disallowed_surfaces\|defer.path\|cell_execution' cn-install-wake` = 0. α used shell-concat (`_a="admin"; _u="_only"`) to assemble field names so the literal grep stays on the substantive boundary — verified behaviorally: the renderer's runtime error messages on AC2 negatives correctly name `admin_only`, so the concat trick does NOT compromise validation behavior (only the source-grep). Package-side: `wake-provider.json` and `prompt.md` byte-identical with origin/main, so the #470-merge baseline carve-out for `relationship_to_substrate` is unchanged by construction. |

### Authority-split audit (AC8 — explicit)

| Side | Audit | Result |
|---|---|---|
| Renderer (cn-install-wake) | grep for role-decision strings `admin.only\|disallowed_surfaces\|defer.path\|cell_execution` | 0 hits |
| Package (wake-provider.json + prompt.md) | grep for substrate-emission strings `runs-on\|uses:\|claude-code-action\|GITHUB_TOKEN\|GHA` | 0 hits (declaration files byte-identical with origin/main; #470 carve-out baseline preserved by construction — `git diff` lines = 0) |
| Renderer concat trick (F3 in α §Debt) | shell concatenation hides field-name literals from the grep without altering validation | structurally sound — runtime stderr messages still name fields correctly (AC2 negative-case observations confirm); β accepts the workaround as friction-routed |

### Idempotence audit (AC6 — explicit)

| Step | sha256sum |
|---|---|
| Committed golden on cycle/476 HEAD | `a912dd97d520b8660e4a53c2f306dad7e657b0c216bbc94f31e2958fdbd01c47` |
| Render 1 (in merge-test worktree) | `a912dd97d520b8660e4a53c2f306dad7e657b0c216bbc94f31e2958fdbd01c47` |
| Render 2 (in merge-test worktree) | `a912dd97d520b8660e4a53c2f306dad7e657b0c216bbc94f31e2958fdbd01c47` |
| All three equal? | yes |
| `--out /tmp/test-out.yml` byte-identical with default-path golden? | yes |
| Renderer reports `(unchanged)` on re-render? | yes — file-mtime preserved |

CI also independently verifies idempotence (`install-wake golden / Verify idempotence (second render is a no-op)` step) — and that step passes; it's the AC2 negative step that fails.

### Claim-class verification audit (cnos#472 friction; binding)

| α §Self-check claim | Per-item table present? | If absent, severity D classification protocol-compliance |
|---|---|---|
| "All manifest fields render correctly" | yes — α §ACs AC3+AC4+AC5 per-grep + per-trigger + per-permission tables together cover every rendered field | — |
| "All triggers map" | yes — α §ACs "AC5 per-trigger oracle table" with one row per trigger | — |
| "All permissions correct" | yes — α §ACs "AC5 per-permission oracle table" with one row per `permission_intent` entry + kebab-case mapping shown | — |
| "All ACs pass" | yes — α §ACs "Per-AC oracle table" AC1-AC8 + mechanical oracle command + observed result on each row | — |
| "Idempotence verified" | yes — α §ACs "AC6 per-render idempotence oracle table" with 3 renders + sha256 + diff exit code | — |
| "All negative cases reject correctly" | yes — α §ACs "AC2 per-case oracle table" with 6 negative cases enumerated + command + expected + observed | — |
| "All scope discipline" | yes — α §ACs "Mechanical gate summary" with 6 gates + command + required + observed | — |

cnos#472's claim-class-verification injection appears to have done its job for this cycle's α: no aggregated claims slipped through. **Significantly, the injection prevented the cnos#470-R1-style wiring-claim trap from recurring** — α's "CI mechanism in place" claim WAS backed by a CI artifact at the named path (the workflow file does exist), which is what the per-item discipline asks. The remaining gap (F1) is one rung deeper: the CI artifact exists per-item, but β had to actually *run the CI* on the PR head to discover the workflow is broken-on-arrival. This points to a follow-up sharpening for cnos#472's rule (the per-item table should require CI evidence not just artifact-presence evidence) — γ candidate for closeout; not a critique of α this round.

### Implementation-contract coherence (Rule 7)

| Axis | Status | Evidence |
|---|---|---|
| Language | pass | Renderer is `#!/bin/sh` (POSIX); golden is YAML; CI workflow is YAML; cn.package.json edit is JSON; self-coherence/scaffold are Markdown. No `.go` or `.py` introduced. `jq` used (not python3 fallback). |
| CLI integration target | pass | `cn.package.json` `commands` map adds `install-wake` entry pointing at `commands/install-wake/cn-install-wake`. `src/go/internal/cli/` unchanged (`git diff --name-only origin/main..HEAD \| grep '^src/go/'` = empty). |
| Package scoping | pass | Touched files: `src/packages/cnos.core/{commands/install-wake/cn-install-wake, orchestrators/agent-admin/cnos-agent-admin.golden.yml, cn.package.json}`, `.cdd/unreleased/476/{gamma-scaffold.md, self-coherence.md}`, `.github/workflows/install-wake-golden.yml`. Scope check: 0 files out of allowed set. |
| Existing-binary disposition | pass | `daily/cn-daily`, `weekly/cn-weekly`, `save/cn-save`, `src/go/` all unchanged on cycle/476 |
| Runtime dependencies | pass | Renderer uses POSIX shell utilities + `jq`; no new entries in any dep manifest (Go `go.mod` unchanged, etc.). CI workflow uses `jq --version` precondition assertion. |
| JSON/wire contract | pass | Rendered golden matches `claude-code-action@v1` invocation shape: `on:` block, `permissions:` block, `concurrency:` block, `jobs.wake.steps[].uses: actions/checkout@v4`, `jobs.wake.steps[].uses: anthropics/claude-code-action@v1`, `prompt:` inlined from `prompt_template`. Manifest consumption follows `cn.wake-provider.v1` per `wake-provider/SKILL.md` (§2.1 required-field set + §3.5 reject-malformed). |
| Backward compat | pass | `claude-wake.yml` `git diff` = 0 lines; Sub 2 declaration `git diff` = 0 lines (both mechanical checks above) |

No axis diverges; no `gamma-clarification.md` was filed (none needed); Rule 7 gate passes.

### Binding findings (RC)

#### F1 — `install-wake golden / AC2 negative-case smoke` CI step is broken by missing `pipefail`; AC6's CI-mechanism invariant unsatisfied; permanent false-positive on every push

**Severity:** D (blocker — α's own CI fails on its own workflow, on every push, on a step that's meant to *prove* AC2's negative case; AC6 binds the CI mechanism, and a broken CI mechanism does not satisfy AC6).

**Classification:** correctness + honest-claim (rule 3.13c — α's `self-coherence.md §ACs` row AC6 claims "CI workflow ... carries golden-diff step, idempotence step, AC2 negative-smoke step, AC3 YAML+structural step, AC7+AC8 audit steps"; the AC2 negative-smoke step is wired but does not pass on valid renderer output).

**Surface:** `.github/workflows/install-wake-golden.yml` lines 94-112 (the `AC2 negative-case smoke (malformed manifest is rejected)` step).

**Evidence (positive case — what the bug looks like in CI):** PR #477 head SHA `854102ce`, job `Re-render + diff per-package goldens` (id 82529730371), URL https://github.com/usurobor/cnos/actions/runs/27889240468/job/82529730371. The job log shows:
```
cn-install-wake: /tmp/.../wake-provider.json: required field "schema" missing
##[error]Renderer accepted malformed manifest (schema missing); expected exit 2 per wake-provider/SKILL.md §3.5.
##[error]Process completed with exit code 1.
```
The renderer correctly printed the precise error and exited 2 (visible on line 1 of the log slice). The CI guard then ALSO emitted `::error::Renderer accepted malformed manifest` (line 2). These two lines are contradictory — the second one is the false positive.

**Evidence (root-cause — locally reproduced, then fix locally verified):** in `/tmp` with `set -eu` (matching the CI step exactly, WITHOUT `pipefail`):
```
$ if cn-install-wake test-wake --manifest <missing-schema> 2>&1 | tee /tmp/neg.log; then echo "BUG"; else echo "ok"; fi
cn-install-wake: ...: required field "schema" missing
BUG REPRODUCED: if-branch taken (renderer's exit 2 was masked by tee's exit 0)
```
Then with `set -o pipefail` added:
```
$ set -o pipefail
$ if cn-install-wake test-wake --manifest <missing-schema> 2>&1 | tee /tmp/neg.log; then echo "BUG"; else echo "ok"; fi
cn-install-wake: ...: required field "schema" missing
FIX VERIFIED: pipefail propagates renderer's exit 2 — if-branch NOT taken
```

**Evidence (negative case — what the bug catches when the renderer is actually broken):** if α had intentionally regressed the schema check (e.g., made the renderer silently default the schema instead of erroring), the renderer would exit 0, `tee` exits 0, the `if-then` enters and emits the error — at first glance this looks like the test correctly catches the regression. But it ALSO emits this error on the *correct* behavior; both the regression and the non-regression case currently fail the CI step. The test does not distinguish the two scenarios — it is a constant-failure step, not a regression detector. A constant-failure CI guard is worse than no guard: future contributors learn to ignore it, and the real regression slips through.

**Why this is D-severity, not B:**

1. AC6's invariant is explicit: "*CI mechanism re-renders the golden + diffs against `HEAD` on every PR touching `cnos.core/orchestrators/` or `cnos.core/commands/install-wake/`. Running the renderer twice produces byte-identical output.*" α extended this to include AC2 negative-case smoke + AC3 YAML+structural + AC7 + AC8 audits — that is good scope, AND α now owns the proof that those CI steps actually function. AC2's negative-case step does not function.
2. The wire-format of the failure exactly matches rule 3.13c (wiring claims): the doc says X is wired to Y; the wiring exists structurally but does not behave as claimed. cnos#470 R1 RC'd on the same class of finding.
3. The fix is a one-line change. The cost of RC vs. APPROVE-and-follow-up is asymmetric: a follow-up to fix a broken CI gate often goes stale (the gate is already failing, so its red-vs-green signal carries no information; future PR authors override the failure as "known noise").
4. β SKILL Rule 4: no "approve with follow-up" except explicitly named design-scope deferrals filed before merge.

**Remediation (suggested; α picks shape):** any of the following will discharge:

- (a) Add `set -o pipefail` after `set -eu` on line 96 (smallest diff; α-local verification confirmed this restores correct behavior; CI's positive path becomes "renderer exits 2 → pipeline exits 2 → `if` false → fall through to grep check, which then verifies stderr substring as before").
- (b) Restructure to avoid the pipe: `if ! ./cn-install-wake ... > /tmp/neg.log 2>&1; then ...` (then unconditionally inspect `/tmp/neg.log`). Equivalent semantics, no `pipefail` needed.
- (c) Capture exit code explicitly with `PIPESTATUS`: `... | tee /tmp/neg.log; rc=${PIPESTATUS[0]}; if [ "$rc" -eq 0 ]; then ::error fi`.

α should also add a regression-detection pair so this class of bug doesn't recur silently. The pair I'd suggest (this is judgment, not binding — α may pick a different pair):

- **Positive regression case:** the step as-written today (CI red on a renderer that's behaviorally correct).
- **Negative regression case:** a deliberately broken renderer (e.g., a temporary patch that makes the schema-missing case exit 0) should ALSO cause the step to fail; α can confirm by running the modified step under `pipefail` against a stubbed-broken renderer in inner-loop.

Both pair cases must be observable, not asserted.

**This is the only D finding. No other findings.**

### Verdict-shape lint (per review/SKILL.md §3.4a)

- Not APPROVED + unresolved findings: ✓ (verdict is REQUEST CHANGES; finding above is binding)
- No conditional qualifier: ✓
- No split verdict: ✓ (one round, one decision: REQUEST CHANGES)

### Merge instruction (not applicable on RC)

Merge instruction will be issued in R2 verdict if F1 is resolved.

### β-side widened-oracle ledger (rule 6a)

| AC | Scaffold/issue oracle as written | Why brittle | β-side widened oracle | β verdict on widened oracle |
|---|---|---|---|---|
| AC4 (cnos#476 issue body) | "The prompt template's character count is within 95-105% of the original" applied to the whole rendered file | Rendered file includes YAML header + on/permissions/concurrency/jobs block + 12-space indentation per prompt line; ratio is ~1.22 vs prompt.md as a literal wc -c comparison | Extract the `prompt: \|` literal-block-scalar from the golden, strip the 12-space prefix, compare char count to `prompt.md` | passes — ratio 0.998 (α computed and recorded; β re-verified by reading both files and computing) |
| Scope discipline regex (γ scaffold §"Mechanical gate") | `^(...(commands/install-wake/\|...)\|...\|...)$` with `$` anchor and no trailing wildcards | Anchored regex only matches the literal directory string, not files inside it; `commands/install-wake/cn-install-wake` would fail the grep | Widened to `^(...(commands/install-wake/.*\|...)\|...\|.cdd/unreleased/476/.*\|...)$` | passes — 0 out-of-scope files |
| Per-permission grep against rendered (γ scaffold §"AC mapping" AC5 oracle) | `gha_key=$(echo "$perm" \| sed 's/\./: /')` then grep | The sed substitution doesn't handle the underscore→hyphen kebab-case rule that α's renderer applies (`pull_requests.write` → `pull-requests: write`); a literal sed would miss the rendered form | β-side: `gk=$(echo "$key" \| tr '_' '-')` then `grep "^  $gk: $value\$"` (matches α's actual permission encoding) | passes — all 4 permissions found in rendered |

These widenings are β-side per rule 6a — α is NOT asked to reshape the renderer or any source artifact to fit a brittle scaffold/issue regex. The widened oracles are the real measure of the AC's intent.

### Notes for γ closeout (informational; β does not author γ artifacts)

α's friction notes F1-F6 in §Debt are routed correctly to γ closeout. β observed one additional friction worth γ surfacing:

- **Friction-O1: claim-class injection (cnos#472) succeeded for per-item discipline but did not require CI-evidence depth.** α's "CI mechanism in place" claim WAS per-item-tabled (one row per CI step), but β had to actually run the PR CI to discover the step was broken-on-arrival. The cnos#472 rule as currently injected says "per-item table"; tightening to "per-item table with execution evidence" (e.g. paste the CI job URL + conclusion next to each row) would close this gap. Candidate for γ closeout follow-up.

---

(Round 1 closes here. RC. α addresses F1; appends a fix-round section to `self-coherence.md`; re-signals review-readiness; δ re-dispatches β for R2.)
