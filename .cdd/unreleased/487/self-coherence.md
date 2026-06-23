---
cycle: 487
parent_issue: cnos#487
master_tracker: cnos#467 (Sub 5C — wave-goal-achievement cell)
cycle_branch: cycle/487
base_main_sha: ff47c3600e9d2b804a8f2d76eeff4544e6a13c4d
scaffold_sha: 30967958ee31796f1fa51522b0bba780436542d3
authored_by: α@cdd.cnos (dispatched by bootstrap-δ via δ-interface session)
date: 2026-06-23 (UTC)
shape: TWO-STAGE — Stage 1 (this document = R0); Stage 2 (post-merge; deferred ACs)
---

# α self-coherence — cnos#487 Sub 5C (cds-dispatch render + activate)

## §Design

### Design fork D1 — render-vs-golden divergence

**Resolution: byte-identical (no divergence; no golden update).**

Empirical verification at R0:
- Pre-flip golden sha256: `75e0406662ea12b3caa5ad01d8dbe6bdfb620d0c047fd910917a8e930855605e` (cycle base).
- Post-flip in-tree render (default `--out`, no override): identical sha256 — `cn-install-wake: cds-dispatch → ... (unchanged)`.
- Post-flip `--out .github/workflows/cnos-cds-dispatch.yml` render: identical sha256 `75e0406662ea12b3caa5ad01d8dbe6bdfb620d0c047fd910917a8e930855605e`.
- Three-way hash comparison passes.

**Rationale.** The renderer's only activation-state-dependent code path is the refusal gate (lines 370-384 of `cn-install-wake`): if `activation_state_effective != "live"`, exit 3; otherwise, render. The renderer does NOT emit the manifest's `activation_state` value into the YAML. The pre-cycle golden was rendered with `--activation-state-override live` against a `declaration-only` manifest, yielding `activation_state_effective = "live"`. The post-flip render uses no override against a `live` manifest, also yielding `activation_state_effective = "live"`. Same effective value → same render path → byte-identical output. This matches FN-7 of the γ scaffold ("renderer behavior change at activation_state flip is a load-bearing assumption") — empirically confirmed.

**Note on prompt body.** The prompt template `prompt.md` still carries its historical "Activation state: declaration-only" notice (the renderer inlines `prompt.md` verbatim). The prompt body is therefore included as-is in both the pre-flip and post-flip rendered YAML. The prompt is forbidden from modification this cycle (γ scaffold §4.3); the manifest's `activation_state_notes` carries an explicit note recording that the prompt's historical text will be scrubbed in a future cycle. The machine-readable manifest field is the authoritative source for the renderer's refusal gate, so this is a documentation-only artifact (no runtime impact).

### Design fork D2 — install-wake-golden AC5 step disposition

**Resolution: Option B (synthetic fixture).** Per γ recommendation.

Implemented via:
- New directory `src/packages/cnos.core/commands/install-wake/test-fixtures/declaration-only/` with a minimal `wake-provider.json` (`activation_state: "declaration-only"`, `activation_state_notes` mentioning `cnos#454` + `cnos#467` + `preconditions` so the AC5 stderr-content grep checks pass) + a trivial `prompt.md` (required by renderer's prompt-template existence check).
- `.github/workflows/install-wake-golden.yml` AC5 step re-pointed at the fixture via `--manifest src/packages/cnos.core/commands/install-wake/test-fixtures/declaration-only/wake-provider.json`. All three assertions preserved (exit 3, `declaration-only` literal in stderr, `cnos#454|cnos#467|preconditions` regex in stderr).

**Rationale (why B, not A or C):**
- **A (remove the step)** loses the renderer's refusal-gate smoke. Future declaration-only manifests (cdr-dispatch, observer roles per cnos#470) would need to re-introduce it; better to keep the smoke alive on a stable fixture now.
- **C (invert with `--activation-state-override declaration-only` on a live manifest)** relies on the override flag's brittle "any non-live value → refusal" semantics. The flag's documented purpose is to bypass refusal, not to induce it; using it inverted couples the test to renderer internals.
- **B (synthetic fixture)** preserves the assertion's purpose (renderer refuses declaration-only), avoids contamination of the production cds-dispatch surface, and gives a stable fixture for future declaration-only test reuse. The fixture is admin-shape (minimal required fields; simpler than dispatch shape).

### OG self-checks

- **OG-1 (cnos#487 not the smoke target):** Verified via API at R0 — cnos#487 carries `status:in-progress` (bootstrap-δ lock), so the live wake's verify gate (per `dispatch-protocol/SKILL.md` §2.2 step 3) would reject it. The smoke target is a SEPARATE issue α files post-merge (Stage 2 per γ §11). The cycle PR body documents this.
- **OG-2 (admin wake untouched):** Verified via `git diff --name-only origin/main...HEAD | grep -E '(cnos-agent-admin|orchestrators/agent-admin)'` returns empty (only `.cdd/unreleased/487/gamma-scaffold.md` and the in-scope files appear). The two re-render steps of `agent-admin` in install-wake-golden still produce its pinned sha256 `fa6b8c0cd64fb626a5e1e991128cbb27fb883b6d1594914543032a2b0d2d3e72`.
- **OG-3 (smoke exercises real path, not bootstrap):** Stage 2; deferred. α's Stage 2 procedure (γ §11 issue, observe `cnos-cds-dispatch.yml` firing, capture evidence) names the live path explicitly; will NOT bootstrap-δ-rescue the smoke.
- **OG-4 (failure handling explicit):** Stage 2; γ scaffold §12 names the 5 failure symptoms + actions. α's Stage 2 will engage rollback per §12 if any fires.
- **OG-5 (cap vs failure):** Stage 2; classification will be recorded in α-closeout post-Stage-2.

## §AC-by-AC verification

| AC | Stage | Oracle | Result | Pass/fail |
|---|---|---|---|---|
| **AC1** | 1 | `jq -r '.activation_state' ...` → `live`; `jq -r '.activation_state_notes' ... \| grep -oE '488\|489\|cnos#485\|cnos#486' \| sort -u` returns 4 matches | `live` + `488`, `489`, `cnos#485`, `cnos#486` all present | PASS |
| **AC2** | 1 | `test -e .github/workflows/cnos-cds-dispatch.yml` | exists | PASS |
| **AC3** | 1 | sha256 of rendered vs golden | both `75e0406662ea12b3caa5ad01d8dbe6bdfb620d0c047fd910917a8e930855605e` | PASS |
| **AC4** | 2 | `gh run view <run-id> --json conclusion --jq .conclusion` = `success` | deferred (post-merge) | DEFERRED |
| **AC5** | 2 | `gh issue view <smoke-N>` shows `status:in-progress`; claim comment naming wake-name/run-id/protocol/head | deferred | DEFERRED |
| **AC6** | 2 | `.cdd/unreleased/<smoke-N>/` contains 6 artifacts (γ-scaffold, self-coherence §R0, β-review §R0, α/β/γ closeouts) | deferred | DEFERRED |
| **AC7** | 2 | PR URL referencing smoke issue; smoke issue at `status:review` | deferred | DEFERRED |
| **AC8** | 2 | post-merge `claude-wake AC3 re-verify <date>` issue triggers `cnos-agent-admin.yml`; conclusion=success; admin does NOT claim any `protocol:cds` cell | deferred | DEFERRED |
| **AC9** | 2 | smoke firing's job log shows no concurrency-wait state attributable to `agent-admin-sigma` | deferred | DEFERRED |

## §Per-CI-step bash-e audit

`.github/workflows/install-wake-golden.yml` (modified by α at R0) — all `run:` steps run under `/usr/bin/bash -e {0}` per GHA semantics.

| # | Step name | Line range | Command substitutions / pipelines | Guarded? | bash-e exit on intended-success input | Notes |
|---|---|---|---|---|---|---|
| 1 | `Verify jq present` | 39-43 | `jq --version` (single command) | no pipelines | exit 0 | unchanged from base |
| 2 | `Re-render agent-admin wake (cnos.core)` | 45-47 | `cn-install-wake agent-admin` (single command) | no pipelines | exit 0 (unchanged sha) | unchanged from base |
| 3 | `Re-render cds-dispatch wake (cnos.cds)` | 49-57 | `cn-install-wake cds-dispatch` (single command; no override) | no pipelines | exit 0 (unchanged sha) | **MODIFIED** — dropped `--activation-state-override live` flag (manifest is live now) |
| 4 | `Verify goldens unchanged` | 59-71 | `git diff --exit-code -- a b` | inside `if !` | exit 0 (no drift) | unchanged from base |
| 5 | `Verify live cds-dispatch workflow matches golden (sha256)` | 73-92 | `live_sha=$(sha256sum ... \| cut -d' ' -f1)` × 2; `[ "$a" != "$b" ]` | `set -eu` + plain assignment | exit 0 (match) | **NEW** — substitutions cannot exit non-zero (sha256sum + cut on existing files always 0); guard is the `[ ]` test |
| 6 | `Verify idempotence — agent-admin` | 94-107 | `sha_before=$(sha256sum ... \| cut ...)` × 2; `cn-install-wake agent-admin` | (no `set -eu`; relies on bash-e default) | exit 0 (match) | unchanged from base |
| 7 | `Verify idempotence — cds-dispatch` | 109-123 | `sha_before=$(sha256sum ... \| cut ...)` × 2; `cn-install-wake cds-dispatch` | (no `set -eu`; relies on bash-e default) | exit 0 (match) | **MODIFIED** — dropped `--activation-state-override live` flag |
| 8 | `Verify YAML parses (both goldens)` | 125-137 | python3 heredoc; `yaml.safe_load` on each path | inside python | exit 0 (parses) | unchanged from base |
| 9 | `Verify substrate structural shape — agent-admin` | 139-151 | `for needle in ...; do grep -qE ... done` | `set -eu`; explicit `if ! grep` | exit 0 (all match) | unchanged from base |
| 10 | `Verify substrate structural shape — cds-dispatch` | 153-177 | `for needle in ...; do grep -qE ... done`; `grep -qE` schedule | `set -eu`; explicit `if ! grep` | exit 0 (all match) | unchanged from base |
| 11 | `AC5 — declaration-only refusal (synthetic fixture)` | 179-213 | `cn-install-wake test-declaration-only --manifest ...` with `\|\| rc=$?` capture; `grep -qF`; `grep -qE` | `set -u`; explicit `rc=$?` capture; `if [ ]` guards | exit 0 (rc=3, stderr matches both greps) | **MODIFIED** — re-pointed at `src/packages/cnos.core/commands/install-wake/test-fixtures/declaration-only/wake-provider.json`; assertions unchanged |
| 12 | `AC2 negative-case smoke (malformed manifest is rejected)` | 215-234 | `mktemp -d`; heredoc echos; `cn-install-wake ... \| tee /tmp/neg.log` in `if !`; `grep -q` | `set -eu` + `set -o pipefail`; `if !` on the pipeline | exit 0 (pipeline fails because renderer exit 2; pipefail makes `if !` true; grep finds "schema" message) | unchanged from base — verified with pipefail locally |
| 13 | `AC8 + AC7 renderer-side authority audit` | 236-274 | `n=$(grep -ciE ... \|\| true)`; `[ "$n" != "0" ]` | explicit `\|\| true` to defang `grep -c` exit 1 on zero matches | exit 0 (both counts = 0) | unchanged from base; pattern documented inline |

**Subordinate workflow file `.github/workflows/cnos-cds-dispatch.yml`** (NEW; rendered).
- Two `uses:` steps (`actions/checkout@v4` with `token`; `anthropics/claude-code-action@v1` with `with:` parameters including a literal-block `prompt:` body).
- NO direct `run:` blocks in the rendered YAML.
- The action's internal shell exposure is opaque to this audit (it lives inside `claude-code-action@v1`); audit covers the entry-point `with:` parameters only. `claude_code_oauth_token` + `github_token` are bound to repo secrets (`CLAUDE_CODE_OAUTH_TOKEN`, `SIGMA_WORKFLOW_PAT`) per the renderer's contract; same shape as the live `cnos-agent-admin.yml`.

## §Variable consistency table

| Variable | Occurrences | Expected consistency | Verified |
|---|---|---|---|
| `activation_state` value | (a) manifest at `wake-provider.json` line 7: `"live"`; (b) `activation_state_notes` (line 8): "is LIVE in production"; (c) renderer's `activation_state_effective` runtime (`cn-install-wake` line 367): copies `activation_state_declared` since no override; (d) CI guard's AC5 step (`install-wake-golden.yml` line 195): runs against the synthetic fixture instead, which has `"declaration-only"`; (e) PR body's AC1 oracle: `jq -r '.activation_state'` returns `live` | All occurrences agree: production manifest is `live`; synthetic fixture is `declaration-only`; renderer + CI dispatch correctly per fixture | ✅ |
| Cron slots `8 23 38 53` | (a) renderer's `substrate_cron_slots` (`cn-install-wake` line 423); (b) `cnos-agent-admin.yml` schedule block; (c) `cnos-cds-dispatch.yml` schedule block (lines 16-19); (d) Stage 2 polling window math: max 15-min wait between sweeps | All four slots identical in both rendered YAMLs; α's Stage 2 procedure cites the 15-min interval correctly | ✅ |
| Concurrency group `cds-dispatch-{agent}` → `cds-dispatch-sigma` | (a) manifest `concurrency_intent.group` (`wake-provider.json` line 97): `cds-dispatch-{agent}`; (b) rendered YAML `concurrency.group` (line 30 of golden + workflow): `cds-dispatch-sigma`; (c) install-wake-golden substrate-shape check (line 162): `'^  group: cds-dispatch-sigma$'`; (d) AC9 oracle: verifies no interference with `agent-admin-sigma` (admin's group) | All resolve to `cds-dispatch-sigma` (with `{agent}` → `sigma` via `agent_variable.default`); CI shape check pinned on the literal | ✅ |
| Selector labels `dispatch:cell`, `protocol:cds`, `status:todo` | (a) manifest `selector.include` (lines 12-14); (b) rendered YAML `if:` gate (line 35): three `contains(github.event.label.name, '<label>')` clauses joined by `\|\|`; (c) smoke issue labels (Stage 2; γ §11 names verbatim); (d) AC5 oracle: claim transitions `status:todo → status:in-progress`; (e) OG-1 lock reasoning: cnos#487 carries `status:in-progress` so verify gate rejects | All three labels in `selector.include`; all three in the `if:` OR-chain; smoke issue applies all three at filing (per γ §11) | ✅ |
| `bot_name`/`bot_id` (`sigma@cnos.cn-sigma.cnos` / `41898282`) | (a) `agent_bot_name`/`agent_bot_id` helpers (`cn-install-wake` lines 127-138); (b) rendered cnos-agent-admin.yml lines 46-47; (c) rendered cnos-cds-dispatch.yml lines 46-47; (d) Stage 2 OG-3 verification (claim comment author = bot) | Identical bot identity across both rendered wakes | ✅ |

## §Self-check (OG-1 through OG-5)

- ✅ OG-1 — cnos#487 carries `status:in-progress` (lock confirmed via API at R0); smoke target is SEPARATE issue α files post-merge per γ §11. Cycle PR body names the boundary explicitly.
- ✅ OG-2 — `git diff --name-only origin/main...HEAD` does NOT include `cnos-agent-admin.yml` or `orchestrators/agent-admin/`. agent-admin's golden sha256 unchanged (`fa6b8c0cd64fb626a5e1e991128cbb27fb883b6d1594914543032a2b0d2d3e72`).
- ⏳ OG-3 — Stage 2; α's procedure (γ §11) targets the live wake's firing, not bootstrap-δ.
- ⏳ OG-4 — Stage 2; γ §12 names 5 symptoms + actions; α will engage rollback if any fires.
- ⏳ OG-5 — Stage 2; α-closeout will classify any failure as cap (inherited I4/I5/I6 etc.) or 5C failure (cycle-introduced).

## §Friction notes (FN-α-N for triage)

### FN-α-1 — Prompt body retains "declaration-only" historical text

The rendered workflow at `.github/workflows/cnos-cds-dispatch.yml` carries a "⚠️ Activation state: declaration-only." block in the prompt body (line 61 of the file), inherited from `prompt.md`. The prompt is forbidden from modification this cycle (γ §4.3 + cnos#483/PR landed). The manifest's `activation_state_notes` records that the prompt's historical text will be scrubbed in a future cycle; the machine-readable manifest field is the authoritative source for the renderer's refusal gate, so runtime behavior is correct. But the prompt text the wake firing reads at execution time is inconsistent with its activation state.

**Triage suggestion (gamma-closeout):** file a follow-up sub for the cnos#467 wave to scrub the prompt's historical block (single-file edit to `src/packages/cnos.cds/orchestrators/cds-dispatch/prompt.md` + re-render of the golden + workflow). Low risk; out of scope for cnos#487.

### FN-α-2 — Synthetic fixture lives under cnos.core/commands/install-wake/

The fixture lives under `src/packages/cnos.core/commands/install-wake/test-fixtures/declaration-only/`. Justification: the fixture exercises the renderer's refusal gate (renderer is owned by cnos.core); putting test fixtures next to the command they exercise mirrors common Go-style colocated-tests conventions. Alternative would be a dedicated `tests/` subtree; rejected as out-of-pattern for this cycle.

**Triage suggestion (gamma-closeout):** the existing renderer command directory only contains `cn-install-wake`; cycle/487 adds `test-fixtures/` as a sibling. If future cycles add multiple fixture classes (per cycle/485 + cycle/486 patterns), a `tests/fixtures/` taxonomy may need formalization.

### FN-α-3 — D1 byte-identical hypothesis empirically confirmed at R0

γ scaffold FN-7 noted that the byte-identical render-vs-golden assumption was load-bearing. Empirical verification at R0 confirms the hypothesis: same `activation_state_effective` → same render path → identical bytes (sha256 `75e0406662ea12b3caa5ad01d8dbe6bdfb620d0c047fd910917a8e930855605e`). This is a positive friction note for the renderer's design: the gate-vs-output separation is clean.

**Triage suggestion (gamma-closeout):** consider promoting this invariant to renderer-doctrine ("activation_state affects gate, not output") so future contributors don't accidentally couple them.

### FN-α-4 — Two-stage cycle's PR-merge-then-smoke shape

Cycle/487 ships its PR (Stage 1) before the smoke runs (Stage 2). The PR's verification is purely on the rendered/CI-side; the production verification waits for merge + wake firing. β converges on Stage 1 only; α/operator handle Stage 2 evidence-gathering. Per γ FN-1 + FN-6, this is the first time the wave produces this shape.

**Triage suggestion (gamma-closeout):** if Sub 5C succeeds end-to-end, the γ-template should formalize a "two-stage cycle" mode for future wake-activation cycles. If it fails (OG-4 fires), the rollback procedure (γ §12) needs to be exercised + closeouts will document the rollback as the closing evidence rather than the smoke evidence.

### FN-α-5 — install-wake-golden `paths:` filter extension

I extended the workflow's `paths:` to include `.github/workflows/cnos-cds-dispatch.yml` so that a hand-edit of that file (which should not happen, but if it did) re-triggers the golden-diff check. The test-fixtures directory is already covered by `src/packages/cnos.core/commands/install-wake/**`.

## §R0 review-ready signal

**α R0 review-ready at commit `a24bedaa` (this self-coherence file shipped in that commit; the review-ready signal lands as a discrete follow-up commit per cycle/485 + cycle/486 pattern so δ can detect it cleanly).** β may engage. All Stage 1 ACs (AC1–AC3) verified locally; OG-2 audit clean; doctrine consistency clean; install-wake-golden CI steps re-verified locally (including pipefail behavior of AC2 negative step). Stage 2 ACs (AC4–AC9) deferred per scaffold §3.2.

α's Stage 2 procedure is described in the cycle PR body's "Stage 2 (post-merge)" section + γ scaffold §11. δ will re-dispatch α after operator merges this PR.
