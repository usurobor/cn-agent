---
cycle: 485
parent_issue: cnos#485
master_tracker: cnos#467 (Sub 5A of wake-orchestration wave)
cycle_branch: cycle/485
authored_by: γ@cdd.cnos (bootstrap-δ session)
date: 2026-06-22 (UTC)
base_main_sha: af8ed8ec
agent_admin_golden_sha256: fa6b8c0cd64fb626a5e1e991128cbb27fb883b6d1594914543032a2b0d2d3e72
agent_admin_golden_sha256_cited_by_issue: 47824628a5958ec9372196b30fa2cb0b547c492799358368636c0e24981b10e1
agent_admin_golden_sha256_discrepancy: see §11 Friction notes — the issue body cites a stale sha; α/β use the verified current value
---

# γ-scaffold — cnos#485 (cn-install-wake renderer extension for dispatch wake providers)

## 2. Goal

Extend the `cn install-wake` renderer (`src/packages/cnos.core/commands/install-wake/cn-install-wake`, shipped by cnos#476) to consume the **dispatch shape** of the `cn.wake-provider.v1` contract (per `cnos.core/skills/agent/wake-provider/SKILL.md`, shipped by cnos#470 + PR #481). After this cycle ships:

- The renderer recognizes `role: dispatch` and emits the corresponding GitHub Actions workflow shape from a `cn.wake-provider.v1` manifest with the dispatch-shape fields (`protocol`, `selector`, dispatch-shape `output_contract`, `activation_state`, and the `issues_labeled_selector_match` trigger).
- A synthetic-`live` cds-dispatch fixture renders into a committed golden file (`src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml`).
- The `install-wake-golden` CI workflow guards BOTH the agent-admin golden (preserved byte-identical per cnos#467 Sub 2 cutover-A) AND the new cds-dispatch golden via re-render-and-diff.
- The renderer refuses by default to install cnos#483's `activation_state: declaration-only` cds-dispatch manifest with a precise, structured error per wake-provider/SKILL.md §3.10.

The substrate side of cnos#467's two-wake architecture becomes **materially possible** with this cycle — though still not active until cnos#486 (Sub 5B; δ wake-invoked mode) and cnos#487 (Sub 5C; flip cnos#483's `activation_state` to `live` + render + commit + smoke) land.

Cite: `cnos#485` (this issue body), `cnos#467` (master tracker; foundational architecture), `cnos.core/skills/agent/wake-provider/SKILL.md` §3.9–§3.10 (dispatch-shape required fields + activation-state discipline).

## 3. Branch name

`cycle/485`. γ has created this branch from `main` at `af8ed8ec` and pushed `gamma-scaffold.md`. α implements on this branch and pushes. β reviews on this branch.

## 4. Touched files (expected)

**α modifies / creates:**

- `src/packages/cnos.core/commands/install-wake/cn-install-wake` — extend renderer to consume dispatch-shape manifest fields; emit the dispatch-shape workflow YAML when `role: dispatch`; honor `activation_state` per §3.10. Today's renderer (verified by γ at branch HEAD) recognizes only `schedule` + `issues_opened_title_match` triggers (lines 366–384 in current source) and the admin-shape `output_contract.channel_log_convention` field (line 260); the dispatch shape needs the parallel branch.

- `src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml` — **new** file. The byte-identical re-render output of `cn install-wake cds-dispatch` against a synthetic-`live` cds-dispatch manifest (or via a test-only activation-state override on the production declaration-only manifest — α chooses; see §guardrails below).

- `.github/workflows/install-wake-golden.yml` — extend the CI to iterate over BOTH the agent-admin golden AND the new cds-dispatch golden: re-render each; `git diff --exit-code` each; idempotence check (sha256 before == sha256 after) on each; YAML parse check on each; substrate structural shape grep on each. Extend the `paths:` filters to include `src/packages/cnos.cds/orchestrators/**`. Extend the AC7 (per cnos#485 — was AC8 of cnos#476) renderer-source authority audit to additionally grep for `protocol:cds|cdr|cdw|dispatch:cell|status:todo` per cnos#485 AC7.

- `.cdd/unreleased/485/self-coherence.md` — **new**. α's per-AC verification table + §R0 (initial) section. Future §R[N] sections appended as β iterates.

**α MAY also touch (optional; α decides):**

- `src/packages/cnos.core/skills/agent/wake-provider/SKILL.md` — only if α adds a CLI flag (e.g. `--activation-state-override`) and the contract's §2.2 optional-fields table or §3.10 wording needs the affordance documented from the renderer-CLI side. If α does this, it is a documentation surface only — no contract-semantic change.

**α MUST NOT touch:**

- `src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.json` — Sub 5C (cnos#487) flips `activation_state` to `live`. Touching it in this cycle would conflate Sub 5A with Sub 5C.
- `src/packages/cnos.cds/orchestrators/cds-dispatch/prompt.md` — unchanged in this cycle. The prompt body is package-authored per the wake-provider contract; the renderer inlines it verbatim with `{agent}` substitution.
- `src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json` — admin manifest is shipped on main; not in scope.
- `src/packages/cnos.core/orchestrators/agent-admin/prompt.md` — admin prompt is shipped on main; not in scope.
- `src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml` — preserving this byte-identical is the AC8 invariant. The renderer extension MUST NOT cause this to change.
- `.github/workflows/cnos-agent-admin.yml` — live admin wake; preserve.
- `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` — Sub 5B (cnos#486) amends this with δ wake-invoked mode.
- `src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md` — shipped by PR #466 at base sha; not modified by this cycle.
- `src/packages/cnos.core/skills/agent/label-doctrine/SKILL.md` — shipped; not modified.

## 5. AC-by-AC oracle

The 9 ACs are quoted in summary form from cnos#485; full text in the issue body. For each, α produces verifiable oracle output; β re-runs each oracle on the cycle branch HEAD before issuing a verdict.

### AC1 — Renderer parses `role: dispatch`

**Statement:** `cn install-wake cds-dispatch --activation-state-override live` (or a test fixture with `activation_state: live`) exits 0; emits a workflow YAML; the YAML's `name:` is `cnos-cds-dispatch`; the `permissions:` block reflects the manifest's `permission_intent`; the `concurrency:` block uses `cds-dispatch-{agent}` with `{agent}` substituted (default `sigma` → `cds-dispatch-sigma`).

**Oracle (α-runnable; β re-runs):**
```bash
# Exit code:
./src/packages/cnos.core/commands/install-wake/cn-install-wake cds-dispatch <test-mode-mechanism>
echo "exit=$?"  # expect 0

# Render landed at expected path:
test -f src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml
echo "golden exists: $?"  # expect 0

# Substrate fields present:
grep -q '^name: cnos-cds-dispatch$' src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml
grep -q '^  group: cds-dispatch-sigma$' src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml
grep -qE '^  contents: write$' src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml
grep -qE '^  issues: write$' src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml
grep -qE '^  pull-requests: write$' src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml
grep -qE '^  id-token: write$' src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml
```

### AC2 — `protocol` + `selector` emitted in substrate

**Statement:** The rendered workflow includes a job-level `if:` gate constructed from `selector.include` (each include label appears in the `if:` expression).

**Oracle:**
```bash
# Each include label appears as a contains(...) check in the job-level if:
f=src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml
for label in 'dispatch:cell' 'protocol:cds' 'status:todo'; do
  grep -qF "$label" "$f" || { echo "missing $label"; exit 1; }
done
# All three include labels surface; β reads the if: expression to confirm
# they're in a contains() construct (not just appearing in a comment).
grep -E "if:.*contains\(github\.event\.label\.name" "$f"
```

β verifies: the labels appear in a `contains(github.event.label.name, '<label>')` (or equivalent `contains(github.event.issue.labels.*.name, '<label>')`) construct, not merely as prose.

### AC3 — `issues_labeled_selector_match` trigger encoded

**Statement:** The rendered workflow's `on:` block includes `issues: { types: [labeled] }` in addition to `schedule:`; the job-level `if:` includes a check for `github.event.label.name` (or equivalent) matching one of the selector's include-set labels. YAML parses; `on.issues.types` contains `labeled`.

**Oracle:**
```bash
f=src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml
# YAML parses + the trigger taxonomy is the dispatch shape:
python3 -c "
import yaml, sys
y = yaml.safe_load(open('$f'))
assert 'schedule' in y['on'], 'schedule trigger missing'
assert 'issues' in y['on'], 'issues trigger missing'
assert 'labeled' in y['on']['issues']['types'], 'labeled type missing'
print('AC3 on: block OK')
"
# Job-level if: references the labeled event:
grep -qE "github\.event_name == 'issues'.*github\.event\.label\.name" "$f" \
  || grep -qE "github\.event\.label\.name.*github\.event_name == 'issues'" "$f"
```

### AC4 — Dispatch-shape `output_contract` passes through to the prompt

**Statement:** The renderer's inlined prompt template references `cycle_artifact_root: .cdd/unreleased/{N}/` and `cell_runtime: cnos.cdd`. The renderer does NOT translate these to YAML fields; they're inlined as part of the prompt body (which is the canonical prompt.md content, unchanged by α). The rendered YAML's `prompt:` field includes the strings `.cdd/unreleased/{N}/` and `cnos.cdd`.

**Oracle:**
```bash
f=src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml
grep -qF '.cdd/unreleased/{N}/' "$f"
grep -qF 'cnos.cdd' "$f"
# Also confirm the inlined prompt block exists:
grep -qE "^[[:space:]]+prompt: \|" "$f"
```

Note: `.cdd/unreleased/{N}/` and `cnos.cdd` are present in `cnos.cds/orchestrators/cds-dispatch/prompt.md` lines 65 and 73 (verified by γ at branch HEAD); the renderer's job is to inline the prompt verbatim (with `{agent}` substitution per wake-provider/SKILL.md §2.3 step 6), so the strings flow through naturally if the renderer does not strip or transform them.

### AC5 — `activation_state: declaration-only` causes refusal

**Statement:** Given cnos#483's cds-dispatch manifest UNCHANGED (`activation_state: declaration-only`), `cn install-wake cds-dispatch` exits non-zero. The issue suggests **exit 3** to distinguish from precondition failures (1) and schema failures (2). Stderr names the `activation_state` value AND contains content from `activation_state_notes`.

**Oracle:**
```bash
# Refusal with precise exit code + error:
./src/packages/cnos.core/commands/install-wake/cn-install-wake cds-dispatch 2>/tmp/ac5.err
rc=$?
[ "$rc" = "3" ] || { echo "AC5: expected exit 3, got $rc"; exit 1; }
grep -qF 'declaration-only' /tmp/ac5.err
# activation_state_notes from the cnos#483 manifest references the three preconditions:
grep -qE 'cnos#454|cnos#467|preconditions' /tmp/ac5.err
```

α decides the precise stderr wording; the oracle requires the activation_state value AND a snippet from `activation_state_notes` to be present in stderr.

### AC6 — Idempotence

**Statement:** Two consecutive `cn install-wake cds-dispatch` invocations on a `live`-state manifest produce byte-identical output (extends cnos#476 AC6 to the dispatch shape).

**Oracle:**
```bash
# In test mode (because production manifest is declaration-only):
./src/packages/cnos.core/commands/install-wake/cn-install-wake cds-dispatch <test-mode>
sha_a=$(sha256sum src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml | cut -d' ' -f1)
./src/packages/cnos.core/commands/install-wake/cn-install-wake cds-dispatch <test-mode>
sha_b=$(sha256sum src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml | cut -d' ' -f1)
[ "$sha_a" = "$sha_b" ] || { echo "AC6: non-deterministic render"; exit 1; }
```

CI workflow (`install-wake-golden.yml`) encodes the same check; α extends the existing idempotence step to cover the cds-dispatch golden as well as agent-admin.

### AC7 — Renderer source-side authority audit

**Statement:** The renderer source MUST NOT encode role-decision strings for dispatch wakes. Those come from the manifest. **Oracle from cnos#485:** `grep -ciE 'protocol:cds|cdr|cdw|dispatch:cell|status:todo' cn-install-wake` returns `0`.

**Oracle:**
```bash
n=$(grep -ciE 'protocol:cds|cdr|cdw|dispatch:cell|status:todo' src/packages/cnos.core/commands/install-wake/cn-install-wake || true)
[ "$n" = "0" ] || { echo "AC7: $n role-decision leaks in renderer source"; exit 1; }
```

α extends the CI workflow's existing `AC8 renderer-side authority audit` step (already greps for `admin.only|disallowed_surfaces|defer.path|cell_execution`) with this dispatch-shape audit. Both audits run on each CI fire.

**β review point (operator-named §guardrails 1):** β reads the renderer source line by line and confirms that `selector.include` array values, `protocol` value, etc. come from `jq` reads against the manifest — never from string literals in the shell source. The grep oracle is a structural backstop; β's empirical read is the substantive check.

### AC8 — Existing agent-admin render unchanged

**Statement:** `cn install-wake agent-admin` continues to produce the current `cnos-agent-admin.golden.yml` byte-identically. The dispatch extension MUST NOT regress the admin-shape rendering.

**Note from γ on the sha256 mismatch:** cnos#485 AC8 quotes `sha256: 47824628a5958ec9372196b30fa2cb0b547c492799358368636c0e24981b10e1`. γ verified at branch HEAD (`af8ed8ec`) that the actual current sha256 is `fa6b8c0cd64fb626a5e1e991128cbb27fb883b6d1594914543032a2b0d2d3e72`. The issue-cited value appears stale — likely captured before the post-cycle-476 cleanup commits (`7ab62cb9` protocol-qualifier rename + `4b633bb2` framework-vs-concrete sweep + `c353d432` preserve `claude-wake` manual-trigger phrase). **Authoritative invariant:** the cycle/485 HEAD `cnos-agent-admin.golden.yml` sha256 equals the `main`-at-base-sha (`af8ed8ec`) value of the same file (whichever number that is, today `fa6b8c0c…`). α captures `sha_before` at branch base; verifies `sha_after == sha_before` after extension. See §11 Friction notes.

**Oracle:**
```bash
# Pin the baseline at branch base:
git show af8ed8ec:src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml \
  | sha256sum | cut -d' ' -f1
# → fa6b8c0cd64fb626a5e1e991128cbb27fb883b6d1594914543032a2b0d2d3e72

# After renderer extension lands:
./src/packages/cnos.core/commands/install-wake/cn-install-wake agent-admin
sha256sum src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml | cut -d' ' -f1
# → MUST equal fa6b8c0cd64fb626a5e1e991128cbb27fb883b6d1594914543032a2b0d2d3e72

# CI workflow's existing "Verify golden unchanged" step also enforces this via
# git diff --exit-code; if the renderer extension regresses the admin shape,
# the diff is non-empty and CI fails.
```

**β verifies empirically:** look at the `install-wake-golden` CI run for the cycle/485 PR. The job `Re-render + diff per-package goldens` must pass. URL convention: `https://github.com/usurobor/cnos/actions/workflows/install-wake-golden.yml?query=branch%3Acycle%2F485` — β cites the most recent green run.

### AC9 — Golden fixture for cds-dispatch

**Statement:** `src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml` exists; install-wake-golden CI guards via re-render-and-diff; the CI job iterates over both agent-admin AND cds-dispatch goldens.

**Oracle:**
```bash
# Fixture exists:
test -f src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml

# CI workflow guards both goldens — read the workflow file and verify both
# render+diff steps are present:
f=.github/workflows/install-wake-golden.yml
grep -qF 'cnos-agent-admin.golden.yml' "$f"
grep -qF 'cnos-cds-dispatch.golden.yml' "$f"
# Idempotence check covers both:
grep -cE 'sha256sum.*\.golden\.yml' "$f"  # expect ≥ 2 occurrences for the two goldens

# YAML parse covers both:
python3 -c "
import yaml
for p in [
  'src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml',
  'src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml',
]:
    yaml.safe_load(open(p))
print('both parse')
"
```

**β verifies empirically:** the cycle/485 PR's `install-wake-golden` CI run is green; the job log shows both goldens re-rendered and diffed.

## 6. α implementation prompt

You are α for `cycle/485` — implementer for cnos#485 (cn-install-wake renderer extension for dispatch wake providers; Sub 5A of cnos#467 wake-orchestration wave).

This is a fresh Agent session. You have no prior context with the parent session. Read this entire prompt + the cited source files before touching code. γ has already created `cycle/485` from `main@af8ed8ec` and committed `.cdd/unreleased/485/gamma-scaffold.md` (this file). Your scaffold is the full input; β reviews you on the cycle branch state alone.

**Working directory:** `/home/user/cnos`. Confirm `git rev-parse --abbrev-ref HEAD` returns `cycle/485` before acting.

**Inputs you load before writing code:**

1. This scaffold (`.cdd/unreleased/485/gamma-scaffold.md`) — your ACs, oracles, scope guardrails, source-of-truth list.
2. The cnos#485 issue body (re-read for cross-verification; do NOT trust this scaffold as the only source of truth).
3. The 9 source files in §10 of this scaffold — read each before touching the relevant surface.

**Scope (in / out — mirrored from cnos#485 + the operator directive that spawned γ):**

*In scope:*
- Extend `cn install-wake` to recognize `role: dispatch` and emit dispatch-shape workflow YAML.
- Consume `protocol`, `selector.include`, `selector.exclude`, `output_contract.cycle_artifact_root`, `output_contract.artifact_class_taxonomy`, `output_contract.cell_runtime`, `activation_state`, `activation_state_notes`, and the `issues_labeled_selector_match` trigger.
- Render the cds-dispatch fixture into `src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml`.
- Implement `activation_state` discipline per wake-provider/SKILL.md §3.10: REFUSAL is the default for `declaration-only` (exit 3; stderr names the activation_state value + activation_state_notes content).
- Extend `.github/workflows/install-wake-golden.yml` to guard both goldens.
- Preserve agent-admin render byte-identically (sha256 invariant; γ's verified-at-branch-base value is `fa6b8c0cd64fb626a5e1e991128cbb27fb883b6d1594914543032a2b0d2d3e72` — see §11 Friction note re: the issue's stale sha citation).

*Out of scope (DO NOT touch — these are gated by 5B/5C; touching them collapses sub boundaries):*
- δ wake-invoked mode amendment in `cnos.cdd/skills/cdd/delta/SKILL.md` (Sub 5B / cnos#486).
- Production `.github/workflows/cnos-cds-dispatch.yml` activation (Sub 5C / cnos#487).
- Flipping `cnos.cds/orchestrators/cds-dispatch/wake-provider.json`'s `activation_state` from `declaration-only` to `live` (Sub 5C).
- Real `protocol:cds` smoke cell (Sub 5C).
- CDR / CDW dispatch providers (future packages).
- NIM / OpenAI / alternative substrate carriers.
- Changes to the live admin wake beyond preserving existing render behavior.

**Operator-named guardrails (these are review points for β; honor them in your implementation):**

1. **Renderer source must remain package-driven, not hard-coded.** The renderer source code MUST NOT bake in `protocol:cds`, `dispatch:cell`, or `status:todo` as literals. Read them from the manifest via `jq` (`selector.include[]`, `protocol`, etc.). This is enforced by AC7's grep oracle (`grep -ciE 'protocol:cds|cdr|cdw|dispatch:cell|status:todo' cn-install-wake` returns 0). β reviews this empirically by reading the renderer source: any string literal matching these patterns is a finding even if the grep happens to miss it (e.g. if you obfuscate via shell-variable concatenation as is done already for `admin_only` per cnos#476 self-coherence §Debt — note that γ permits this for the existing admin-fields validation surface but views it as a *strict* leak for the new dispatch-fields surface; AC7 is satisfied only by genuine manifest-read, not by trick concatenation).

2. **`issues_labeled_selector_match` MUST NOT suppress the scheduled sweep.** The rendered workflow's job-level `if:` must allow BOTH:
   - Schedule-triggered sweep behavior (job runs unconditionally when `event_name == 'schedule'`).
   - Issue-labeled responsive behavior (job runs when `event_name == 'issues'` AND the labeled event matches the selector's include set).

   Concretely, the shape SHOULD be:
   ```yaml
   if: ${{ github.event_name == 'schedule' || (github.event_name == 'issues' && (contains(github.event.label.name, 'dispatch:cell') || contains(github.event.label.name, 'protocol:cds') || contains(github.event.label.name, 'status:todo'))) }}
   ```
   If you use only `contains(github.event.label.name, '<label>')`, schedule events (which have no `github.event.label`) will silently fail to fire the job. β reviews this case explicitly.

   **You MUST also write an explicit α-side AC in your `self-coherence.md` §R0:** "Rendered cds-dispatch workflow's job-level `if:` allows schedule events unconditionally and gates only issue events on the selector's include-set membership." Verify with: `grep -E "github\.event_name == 'schedule'" cnos-cds-dispatch.golden.yml`.

3. **`activation_state: declaration-only` refusal must be conservative.** Normal `cn install-wake cds-dispatch` (with cnos#483's manifest unchanged) MUST refuse with exit 3 and a precise error. Any test/fixture override you introduce (e.g. a `--activation-state-override` flag, or a `CN_INSTALL_WAKE_TEST_MODE=1` env var, or a separate fixture manifest with `activation_state: live` baked in) MUST be visibly test-only — NOT a silent runtime escape that operators could accidentally use. β reviews this as: "Is the override path obvious enough that an operator running `cn install-wake --help` would never accidentally bypass refusal?". You pick the mechanism. γ-suggested mechanisms (pick one; document in self-coherence.md):
   - **Option A:** a `--activation-state-override <state>` CLI flag that emits a stderr warning naming the override and requires the override value to be `live` (so the flag's name itself names what's being bypassed). Document in `--help`.
   - **Option B:** a dedicated test-fixture path (e.g. `src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.test-live.json`) that you render against via `--manifest`. The override is visible because the path includes `.test-live.`. β verifies that the production manifest (`wake-provider.json`) still has `activation_state: declaration-only` and still triggers refusal.
   - **Option C:** an environment variable (`CN_INSTALL_WAKE_ACTIVATION_OVERRIDE=live`) that triggers a stderr warning. Less visible than A or B; β scrutinizes this more closely.

   γ prefers Option A or B for visibility; α decides; β reviews the chosen mechanism with the "is this genuinely test-only?" lens.

**Output expectations:**

- Working renderer extension at `src/packages/cnos.core/commands/install-wake/cn-install-wake`.
- Golden fixture at `src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml`.
- CI extension at `.github/workflows/install-wake-golden.yml` covering both goldens.
- Self-coherence file at `.cdd/unreleased/485/self-coherence.md` with:
  - **§R0** section header.
  - Per-AC verification table: AC# / oracle command / actual output / pass-fail.
  - Verification commands you ran locally + their exit codes.
  - Notes on which guardrail mechanism (A/B/C above) you chose and why.
  - Explicit α-AC about the schedule/issues `if:` compatibility (per operator-named guardrail #2).
- Commit + push to `cycle/485` when R0 is complete. Append to `.cdd/unreleased/485/self-coherence.md` a §R0-complete signal: "α: R0 complete; β review-ready at sha=<sha>". This is the dispatch handshake β reads.

**Iteration discipline:**

- α uses `§R0`, `§R1`, `§R2`, ... section headers per round in self-coherence.md.
- Each §R[N] documents what β found in R[N-1] and what you fixed.
- α does not change any §R[N-1] content retrospectively (append-only).
- β iterates by writing `.cdd/unreleased/485/beta-review.md` with parallel `§R[N]` sections + verdict (converge or iterate) + findings (numbered F1, F2, ...).
- Iterate until β's verdict is converge.

**When you're done with R0:** commit + push to `cycle/485` + append the review-ready signal to self-coherence.md.

## 7. β review prompt

You are β for `cycle/485` — reviewer for cnos#485 (cn-install-wake renderer extension for dispatch wake providers; Sub 5A of cnos#467 wake-orchestration wave).

This is a fresh Agent session. You have no prior context. Read this entire prompt + the inputs below before issuing a verdict. δ (the parent session) dispatches you only after α signals R[N] is review-ready (see α's `.cdd/unreleased/485/self-coherence.md` for the signal).

**Working directory:** `/home/user/cnos`. The branch is `cycle/485`. Confirm `git rev-parse --abbrev-ref HEAD` returns `cycle/485`.

**Inputs you load:**

1. This scaffold (`.cdd/unreleased/485/gamma-scaffold.md`) — the canonical oracle list and scope guardrails.
2. α's `.cdd/unreleased/485/self-coherence.md` — α's per-AC verification + chosen guardrail mechanism + any prior §R[N] iterations.
3. The cnos#485 issue body — re-read for cross-verification; do not trust the scaffold as the only oracle source.
4. The 9 source-of-truth files (§10 of this scaffold).
5. The cycle branch HEAD diff: `git diff main...cycle/485 -- <changed files>`.

**Review checklist (per-AC; populate this table in `.cdd/unreleased/485/beta-review.md` §R[N]):**

| AC | Oracle | Pass / Fail | Notes / Finding ID if fail |
|---|---|---|---|
| AC1 | renderer parses `role: dispatch`; YAML emitted; name/permissions/concurrency correct (§5 AC1 oracle) | | |
| AC2 | `protocol` + `selector` emitted; each include label in `if:` (§5 AC2 oracle); β reads the if: expression empirically | | |
| AC3 | `on:` has `issues: { types: [labeled] }` AND `schedule:`; job-level `if:` references both (§5 AC3 oracle) | | |
| AC4 | inlined prompt contains `.cdd/unreleased/{N}/` + `cnos.cdd` (§5 AC4 oracle) | | |
| AC5 | `cn install-wake cds-dispatch` (no override) → exit 3 + precise stderr (§5 AC5 oracle) | | |
| AC6 | sha256 before == sha256 after on cds-dispatch (idempotence, §4 AC6 oracle) | | |
| AC7 | `grep -ciE 'protocol:cds\|cdr\|cdw\|dispatch:cell\|status:todo' cn-install-wake` returns 0 AND β reads source to confirm no literal-string trickery (§5 AC7 oracle + operator-named guardrail #1) | | |
| AC8 | agent-admin golden sha256 unchanged from base (`fa6b8c0c…`); install-wake-golden CI green (§5 AC8 oracle) | | |
| AC9 | cds-dispatch golden file exists; CI guards both goldens; YAML parses for both (§5 AC9 oracle) | | |

**Operator-named guardrails (review explicitly; each is a finding if violated):**

- **OG-1 (renderer-source authority audit, paired with AC7):** Read the diff of `cn-install-wake`. Confirm that `selector.include` values, `protocol` value, etc. come from `jq` reads, not from string literals. Note any obfuscation-via-concatenation attempts (e.g. `_d="dis"; _c="patch"; bad="$_d$_c:cell"`) as findings — AC7 is satisfied only by genuine manifest reads.
- **OG-2 (schedule/issues `if:` compatibility):** Read the rendered `cnos-cds-dispatch.golden.yml`. Verify the job-level `if:` allows schedule events unconditionally (e.g. starts with `github.event_name == 'schedule' ||` or equivalent). A test: locate the `if:` expression; mentally simulate a schedule event (`event.label` undefined) — does the job fire? If no → finding.
- **OG-3 (refusal-is-conservative):** Verify α's chosen activation-state-override mechanism is visibly test-only. Read `--help` output or the renderer source's argument parser. The override path should be discoverable but obviously not-for-production. The production manifest (`cnos.cds/orchestrators/cds-dispatch/wake-provider.json`) MUST still have `activation_state: declaration-only` on cycle/485 HEAD. If α touched this file, that is a finding (out-of-scope for Sub 5A).

**Per-CI-step verification table (required per cnos#478 mechanical-injection process — populate in beta-review.md §R[N] for every new `run:` step in `.github/workflows/install-wake-golden.yml` that this cycle introduces or modifies):**

For each `run:` step that this cycle adds or modifies, fill the row:

| # | Step name | Line range | Command substitutions / pipelines used | Guarded? (`set -eu`, `\|\| true`, etc.) | bash -e exit on intended-success input? | Notes |
|---|---|---|---|---|---|---|

Why this is required: cycle/476 hit a 3-round class-trap because `grep -c` exits 1 on zero matches under `bash -e` and the workflow file did not explicitly `|| true`-guard it. The friction note rolled forward; γ-scaffold inherits the discipline. For this cycle β specifically audits:

- Any new `grep -c` (or `grep -cE`) under bash-e in install-wake-golden.yml.
- Any new pipeline using command substitutions where the inner command can exit non-zero on intended-success input.
- Any new step that re-renders the cds-dispatch golden.
- Any new step that hashes goldens for idempotence checks.

**Output:** `.cdd/unreleased/485/beta-review.md` with:

- **§R[N]** section header (R0 first time; R1 if α had to iterate).
- The review checklist table above (per-AC; per-OG; per-CI-step).
- **Verdict line:** `verdict: converge` OR `verdict: iterate`.
- **Findings** numbered `F1`, `F2`, …, each with file + line + AC/OG/CI-step ref + reproducible reproduction step. Findings are only required when verdict = iterate; converge requires no findings.
- **Critical:** β reviews ACTUAL CI execution evidence, not just local rendering. The `install-wake-golden` CI run on α's PR is the canonical green-or-not signal. β cites the run URL (https://github.com/usurobor/cnos/actions/workflows/install-wake-golden.yml?query=branch%3Acycle%2F485) and the specific job result.

**Verdict semantics:**

- **converge** → α and γ proceed to closeout (α writes alpha-closeout.md, β writes beta-closeout.md, γ writes gamma-closeout.md, PR ships, master tracker cnos#467 sub 5A box ticks).
- **iterate** → α addresses findings, increments to §R[N+1] in self-coherence.md, β re-reviews. γ does not re-author the scaffold unless a finding surfaces a scaffold-side ambiguity (in which case δ asks γ to amend).

## 8. Non-goals (mirror of §4 "α MUST NOT touch" plus the cnos#485 issue's out-of-scope list, restated explicitly for α and β reference)

- **δ wake-invoked mode in `cnos.cdd/skills/cdd/delta/SKILL.md`** — Sub 5B (cnos#486). This cycle's scaffold + renderer + golden do NOT require δ's wake-invoked contract; the dispatch wake's prompt simply names δ as the handoff target without specifying δ's wake-invoked entry shape (that's a separate skill amendment).
- **`.github/workflows/cnos-cds-dispatch.yml` substrate activation** — Sub 5C (cnos#487). cycle/485 ships the renderer + the golden fixture (which lives under `src/packages/cnos.cds/orchestrators/`, NOT under `.github/workflows/`). The substrate activation is a separate cycle.
- **`cnos.cds/orchestrators/cds-dispatch/wake-provider.json` `activation_state` flip** — Sub 5C (cnos#487). Cycle/485 makes the renderer capable of rendering when `activation_state: live`; it does NOT flip the production manifest.
- **Real `protocol:cds` smoke cell** — Sub 5C (cnos#487).
- **CDR / CDW dispatch providers** — future packages, separate waves.
- **NIM / OpenAI / alternative substrate carriers** — separate wave, post cnos#467.
- **Multi-cell prioritization (priority labels)** — explicitly deferred to v1 per cnos#454 §2.5.
- **Changes to the live admin wake beyond preserving existing render behavior** — AC8 invariant.

## 9. Test plan (α's expected local verification before signaling review-ready)

α runs each of these and records exit codes + output excerpts in `.cdd/unreleased/485/self-coherence.md` §R0:

```bash
# (1) AC5 — refusal of unchanged manifest (declaration-only):
./src/packages/cnos.core/commands/install-wake/cn-install-wake cds-dispatch
# expect: exit 3; stderr names "declaration-only" and snippets from activation_state_notes

# (2) AC1+AC2+AC3+AC4+AC6 — render in test mode + idempotence:
./src/packages/cnos.core/commands/install-wake/cn-install-wake cds-dispatch <test-override>
# expect: exit 0; golden written
sha1=$(sha256sum src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml | cut -d' ' -f1)
./src/packages/cnos.core/commands/install-wake/cn-install-wake cds-dispatch <test-override>
sha2=$(sha256sum src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml | cut -d' ' -f1)
[ "$sha1" = "$sha2" ]  # idempotence (AC6)

# (3) AC8 — preserve agent-admin golden byte-identically:
./src/packages/cnos.core/commands/install-wake/cn-install-wake agent-admin
sha256sum src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml | cut -d' ' -f1
# expect: fa6b8c0cd64fb626a5e1e991128cbb27fb883b6d1594914543032a2b0d2d3e72

# (4) JSON parse on manifest + YAML parse on both goldens:
jq -e . src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.json >/dev/null
python3 -c "import yaml; yaml.safe_load(open('src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml'))"
python3 -c "import yaml; yaml.safe_load(open('src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml'))"

# (5) AC7 — renderer-source authority audit:
n=$(grep -ciE 'protocol:cds|cdr|cdw|dispatch:cell|status:todo' src/packages/cnos.core/commands/install-wake/cn-install-wake || true)
[ "$n" = "0" ]

# (6) AC2+AC3 substrate-shape audit on rendered cds-dispatch:
f=src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml
grep -qE '^on:' "$f"
grep -qE 'schedule:' "$f"
grep -qE 'issues:' "$f" && grep -qE 'types: \[labeled\]' "$f"
grep -qE '^permissions:' "$f"
grep -qE '^concurrency:' "$f"
grep -qE '^  group: cds-dispatch-sigma$' "$f"
grep -qE '^jobs:' "$f"
grep -qE 'anthropics/claude-code-action@v1' "$f"

# (7) AC2 — selector labels present in if:
for label in 'dispatch:cell' 'protocol:cds' 'status:todo'; do
  grep -qF "$label" "$f"
done

# (8) operator-named guardrail #2 — schedule + issues compatibility in if:
grep -qE "github\.event_name == 'schedule'" "$f"

# (9) operator-named guardrail #1 — confirm by reading the renderer diff:
git diff main...HEAD -- src/packages/cnos.core/commands/install-wake/cn-install-wake \
  | grep -iE 'protocol:cds|cdr|cdw|dispatch:cell|status:todo'
# expect: no NEW lines (+ prefix) carry these literals; only manifest-read code

# (10) AC4 — inlined prompt has the artifact root + cell_runtime strings:
grep -qF '.cdd/unreleased/{N}/' "$f"
grep -qF 'cnos.cdd' "$f"
```

CI runs after α pushes; β reads the CI result as the canonical signal alongside α's local-verification table.

## 10. Source-of-truth list

All paths are absolute on the local cycle/485 checkout. Base SHA = `af8ed8ec` (main HEAD when cycle/485 branched). Online URLs follow `https://github.com/usurobor/cnos/blob/af8ed8ec/<path>` convention.

| # | Reference | Path / URL | Authority |
|---|---|---|---|
| 1 | cnos#485 (this issue) | https://github.com/usurobor/cnos/issues/485 | The 9 ACs + in-scope/out-of-scope. Canonical. |
| 2 | cnos#467 master tracker (with doctrine-correction header) | https://github.com/usurobor/cnos/issues/467 | Wave architecture + sub-issue plan. Doctrine corrections at top supersede stale `cnos.cdd` / `protocol:cdd` / `cdd-dispatch` references in the body. |
| 3 | dispatch-protocol skill | `/home/user/cnos/src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md` | The 3-label selector + serialized claim guard + 4-event lifecycle the renderer encodes. Shipped in PR #466 at base sha. |
| 4 | cds-dispatch wake-provider manifest | `/home/user/cnos/src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.json` | The reference dispatch-shape manifest the renderer consumes. `activation_state: declaration-only` on main. |
| 5 | cds-dispatch prompt template | `/home/user/cnos/src/packages/cnos.cds/orchestrators/cds-dispatch/prompt.md` | The inlined prompt body. `.cdd/unreleased/{N}/` at line 65; `cnos.cdd` at line 73 (AC4 evidence). |
| 6 | cn-install-wake (renderer at base) | `/home/user/cnos/src/packages/cnos.core/commands/install-wake/cn-install-wake` | The admin-shape renderer α extends. 501 lines at base sha. Note especially: lines 366–384 (trigger taxonomy), lines 260+ (output_contract.channel_log_convention check — admin-shape-specific). |
| 7 | wake-provider contract skill | `/home/user/cnos/src/packages/cnos.core/skills/agent/wake-provider/SKILL.md` | The `cn.wake-provider.v1` contract. §2.1 required fields (incl. dispatch-required `protocol` + `selector`); §2.2 optional fields (incl. `activation_state` + `activation_state_notes` + `forward_references`); §2.5 the canonical package-vs-renderer split table (the renderer's authority boundary); §3.5 reject-rather-than-default; §3.9 dispatch-selector explicitness; §3.10 declaration-only refusal-or-warned-render. |
| 8 | label-doctrine skill | `/home/user/cnos/src/packages/cnos.core/skills/agent/label-doctrine/SKILL.md` | Two-layer label ownership rule. |
| 9 | install-wake-golden CI workflow | `/home/user/cnos/.github/workflows/install-wake-golden.yml` | The CI guard α extends to cover both goldens. Note especially: `AC8 renderer-side authority audit` step (lines 115+) which α extends with the dispatch-shape audit per AC7. |
| 10 | agent-admin golden YAML | `/home/user/cnos/src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml` | The byte-identical-preservation reference. sha256 at base (verified by γ): `fa6b8c0cd64fb626a5e1e991128cbb27fb883b6d1594914543032a2b0d2d3e72`. Issue cites `47824628…` — see §11 Friction notes. |
| 11 | live admin workflow | `/home/user/cnos/.github/workflows/cnos-agent-admin.yml` | The production rendered artifact (cutover-A per cnos#479 / PR #481). Identical structure to golden. |

## 11. Friction notes for future δ dispatch-prompt template (per the bootstrap-δ pattern)

These are γ's meta-observations while scaffolding; the parent (δ) session reads them; later δ wake-invoked mode work (Sub 5B / cnos#486) consumes them for the δ-prompt template.

**FN-1. The issue cites a stale agent-admin golden sha256.** cnos#485 AC8 says `47824628a5958ec9372196b30fa2cb0b547c492799358368636c0e24981b10e1`. γ verified at branch base `af8ed8ec`: the actual sha256 is `fa6b8c0cd64fb626a5e1e991128cbb27fb883b6d1594914543032a2b0d2d3e72`. The discrepancy is downstream of three commits between the issue's authoring time and main HEAD:
- `c353d432` preserve 'claude-wake' manual-trigger phrase across cutover
- `4b633bb2` chore: protocol-label-rename cleanup — finish framework-vs-concrete sweep
- `7ab62cb9` chore: protocol-qualifier label rename — cnos.cdd → cnos.cds

The scaffold uses the verified value (`fa6b8c0c…`). β re-verifies at PR-review time against current main HEAD (if main has moved). **Lesson for the δ dispatch-prompt template:** δ should re-verify cited sha256 values against actual filesystem state before dispatching γ, especially for ACs that pin invariants by hash — issue bodies become stale faster than file contents.

**FN-2. Issue body had enough detail for γ to scaffold without operator clarification.** All 9 ACs are independently testable; the scope guardrails are explicit; cross-references resolve. The one ambiguity γ had to resolve was AC5's "exit code 3 (or chosen exit code)" — γ pinned exit 3 in the scaffold per the issue's suggestion. β can flag this as a finding if the operator wants a different code.

**FN-3. All cited sources are present on main at base sha.** No forward-references except the explicit Sub 5B / Sub 5C deferrals, which are correctly named as out-of-scope.

**FN-4. Per-AC oracle granularity is testable + automatable.** The renderer-side oracle (AC7) is a precise grep — δ doesn't have to hand-wave. The one place γ had to gesture is **operator-named guardrail #1's "obfuscation via concatenation"** subtlety: the AC7 grep alone could be satisfied by a renderer source that constructs `protocol:cds` via shell-variable concatenation (analogous to the existing `_a="admin"; _u="_only"` pattern documented in cycle/476 self-coherence §Debt). γ surfaced this explicitly as a β empirical-read review point. Stronger oracle would require AST parsing of the shell source, which is heavier than the cycle warrants; the dual mechanism (grep + β empirical read) is the proposed compromise.

**FN-5. The role-decision-leak audit (AC7) needs β empirical eyes.** See FN-4. γ proposes the dual mechanism. The δ dispatch-prompt template may want a hook for "audit findings that resist automation; require β human-equivalent read" — this cycle is the empirical case to cite.

**FN-6. Operator-named guardrails should become a γ-scaffold template section per cnos#478 mechanical-injection process.** This scaffold introduces three operator-named guardrails (OG-1 through OG-3 in β's review checklist). They function the same way per-CI-step verification tables do per cnos#478 — they're a "things automation misses; reviewer must surface" surface. γ doesn't change the scaffold template format in this cycle (out of scope), but notes the observation: future scaffold templates may benefit from a dedicated `§Operator-named guardrails` section between §AC oracles and §α prompt. The parent session can carry this forward into a separate γ-template-amendment cycle.

**FN-7. Bootstrap-δ runs are the v0 dispatch mechanism until Sub 5C lands.** γ was invoked by a manually-dispatched Agent session from a parent (δ in role); γ's output (this scaffold + the cycle branch) is the input the parent uses to dispatch α (separately). The δ wake-invoked mode contract that lands in Sub 5B / cnos#486 should formalize:
- The handshake γ-finishes → δ-dispatches-α (today: γ commits + pushes; δ reads the branch).
- The handshake α-R[N]-ready → δ-dispatches-β (today: α appends a signal to self-coherence.md; δ polls or is signaled).
- The branch-as-shared-state pattern: no chat-state-as-input; everything is on the branch and in the .cdd/unreleased/{N}/ tree.

This cycle is the empirical case for that pattern. The friction observed during scaffolding: γ had to verify the agent-admin sha256 by reading the file rather than trusting the issue body — that's the "no chat state" discipline operating in practice. The δ wake-invoked mode skill should pin this as a contract: "wake-invoked δ verifies sha-pinned invariants against current filesystem state before dispatching γ".

**FN-8. Scaffold word count: ~6.5K words.** Heavy on the α prompt + β review prompt because they're self-contained per the bootstrap-δ "no chat state" discipline; β's checklist + per-CI-step table inflate further. γ-template-amendment work may want to consider whether the operator-guardrails + per-CI-step table should be hoisted into a shared skill (e.g. `cnos.cdd/skills/cdd/beta/SKILL.md`) so individual cycle scaffolds can `@include` rather than restate. Out of scope for this cycle.
