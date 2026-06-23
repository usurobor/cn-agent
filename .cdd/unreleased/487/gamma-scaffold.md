---
cycle: 487
parent_issue: cnos#487
master_tracker: cnos#467 (Sub 5C of wake-orchestration wave — wave-goal-achievement cell)
cycle_branch: cycle/487
base_main_sha: ff47c3600e9d2b804a8f2d76eeff4544e6a13c4d
predecessor_subs:
  - cnos#485 (Sub 5A — renderer dispatch-shape extension; merged PR #488)
  - cnos#486 (Sub 5B — δ wake-invoked mode amendment; merged PR #489)
authored_by: γ@cdd.cnos (bootstrap-δ via δ-interface session)
date: 2026-06-23 (UTC)
shape: TWO-STAGE — Stage 1 (pre-merge: render + activate); Stage 2 (post-merge: smoke a real protocol:cds cell)
---

# γ-scaffold — cnos#487 (wake-cutover-B + dispatch smoke: activate cnos-cds-dispatch.yml + smoke a real protocol:cds cell)

> **Wave-goal-achievement cell.** This cycle is the load-bearing cell for cnos#467's 10-point completion definition. Two distinct wake classes (admin + dispatch) become observable in production substrate; a real `protocol:cds` cell runs end-to-end through the live dispatch wake; cnos#467 closes when the smoke verifies.

## §1. Goal

Ship the second wake class real in production. Render + activate the `cnos.cds`-owned dispatch wake (`cnos-cds-dispatch.yml`) by flipping cnos#483's `activation_state: declaration-only → live` and committing the rendered substrate artifact; then smoke a real `protocol:cds` cell end-to-end through the live dispatch wake → δ wake-invoked mode → γ/α/β → `.cdd/unreleased/{smoke-N}/` artifact set → PR → `status:review`. The two-wake architecture cnos#467 establishes becomes materially observable in production substrate for the first time.

## §2. Branch name

`cycle/487` (from `main@ff47c3600e9d2b804a8f2d76eeff4544e6a13c4d`).

## §3. Two-stage cycle shape (critical context)

cycle/487 differs from cycle/485 + cycle/486 because the **smoke is post-merge**:

### §3.1 — Stage 1 (pre-merge; the cycle/487 PR's work)

α flips cnos#483's manifest → renders `.github/workflows/cnos-cds-dispatch.yml` → verifies render against golden → potentially extends `install-wake-golden.yml` → opens PR → β reviews (per-AC + 4 OGs + per-CI-step bash-e audit + variable consistency + doctrine + cross-skill consistency + CI evidence + non-goal) → β converges → PR merges → workflow becomes LIVE in production.

### §3.2 — Stage 2 (post-merge; after the cycle/487 PR merges)

α files a SEPARATE smoke issue (per §11 below) → live `cnos-cds-dispatch.yml` fires on label-event or scheduled sweep → wake claims smoke issue (`status:todo → status:in-progress` + claim comment) → δ wake-invoked mode (Sub 5B amendment, merged at PR #489) routes γ/α/β → `.cdd/unreleased/{smoke-N}/` accumulates canonical artifact set → cycle PR ships smoke cell's work → β converges → cell moves to `status:review` → α/β/γ closeouts on cnos#487 capture both Stage 1 + Stage 2 evidence → cnos#487 closes when smoke is end-to-end-verified.

### §3.3 — Implication for closeouts

The cycle/487 PR (Stage 1) can converge + merge BEFORE the smoke runs. Closeouts for cnos#487 specifically MUST capture both Stage 1 implementation evidence (PR merged; workflow live) AND Stage 2 smoke evidence (Actions run URL; claim comment record; artifact set; status transition). α/β/γ closeouts are authored post-Stage-2 (when the smoke completes) so they can capture both stages.

## §4. Touched files

### §4.1 — Stage 1 (cycle PR; in scope for α at R0)

| File | Action | Notes |
|---|---|---|
| `src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.json` | EDIT | Flip `activation_state: "declaration-only" → "live"`; update `activation_state_notes` to record preconditions HAVE landed; cite the merged Sub 5A PR #488 (renderer) + Sub 5B PR #489 (δ wake-invoked mode). |
| `.github/workflows/cnos-cds-dispatch.yml` | NEW | Rendered output from `cn install-wake cds-dispatch --out .github/workflows/cnos-cds-dispatch.yml` (manifest now `live` → no override flag needed; renderer should not refuse). |
| `src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml` | POTENTIAL EDIT | Verify post-flip render matches the synthetic-live golden Sub 5A pinned (current sha256 at base: `75e0406662ea12b3caa5ad01d8dbe6bdfb620d0c047fd910917a8e930855605e`). Sub 5A rendered with `--activation-state-override live` against a `declaration-only` manifest; the post-flip render is on a `live` manifest with no override. The renderer's `activation_state_effective` value should be identical (`live` either way), so the output SHOULD be byte-identical. If it differs, the divergence MUST be understood + the golden updated to match the production render. The cycle/485 closeouts (T-485-3) and §13 friction notes below cover this risk. |
| `.github/workflows/install-wake-golden.yml` | POTENTIAL EDIT | Sub 5A extended CI to render cds-dispatch with `--activation-state-override live`. Post-flip the manifest is `live` — the override flag is no longer required, but its continued presence does no harm (per the renderer's flag semantics: `--activation-state-override live` is the no-op-equivalent override). However: the AC5 negative-case (declaration-only refusal) step in install-wake-golden assumes `cn install-wake cds-dispatch` (no override) returns exit 3 — this CHANGES after the flip (will now succeed). The AC5 step MUST be either: (a) removed if no second `declaration-only` manifest exists post-cycle, OR (b) re-pointed at a synthetic fixture (`--manifest <test-fixture-path>`) that retains the declaration-only assertion. α picks per §6 design fork. |
| `.cdd/unreleased/487/gamma-scaffold.md` | THIS FILE | Authored by γ at R0. |
| `.cdd/unreleased/487/self-coherence.md` | NEW (α writes) | α's per-AC verification + design log + per-CI-step bash-e audit + variable consistency table + R[N] iteration sections. |
| `.cdd/unreleased/487/beta-review.md` | NEW (β writes) | β's per-AC verification + 4 OGs + per-CI-step audit + variable consistency + doctrine + cross-skill consistency + CI evidence + non-goal + verdict R[N]. |

### §4.2 — Stage 2 (post-merge; observational/evidence-gathering)

| Surface | Action | Notes |
|---|---|---|
| (the smoke issue) | NEW (α files post-merge) | Separate `protocol:cds + dispatch:cell + status:todo` issue per §11 design. NOT committed to the repo; lives as a GitHub issue. |
| `.cdd/unreleased/{smoke-N}/*` | NEW (smoke cell's δ writes via wake) | Canonical 6-7 artifact set (`gamma-scaffold.md`, `self-coherence.md` with §R[N], `beta-review.md` with §R[N], `alpha-closeout.md`, `beta-closeout.md`, `gamma-closeout.md`, optionally PRA). Lands on the smoke cell's own cycle branch. |
| `docs/gamma/smoke/cds-dispatch-smoke-YYYYMMDD.md` | NEW (if smoke task creates a doc) | Per §11 design — the smoke task's substantive deliverable IS a metadata doc recording the firing evidence. Lands on the smoke cell's cycle branch. |
| `.cdd/unreleased/487/alpha-closeout.md` | NEW (α writes post-Stage-2) | Captures both Stage 1 + Stage 2 evidence. |
| `.cdd/unreleased/487/beta-closeout.md` | NEW (β writes post-Stage-2) | Same. |
| `.cdd/unreleased/487/gamma-closeout.md` | NEW (γ writes post-Stage-2) | Same; triage of any friction from running the full path end-to-end. |

### §4.3 — α MUST NOT touch (operator-named guardrails; β verifies)

| Forbidden path | Reason |
|---|---|
| `.github/workflows/cnos-agent-admin.yml` | OG-2 — admin wake remains admin-only; Stage 1 MUST NOT regress the live admin wake. |
| `src/packages/cnos.core/orchestrators/agent-admin/*` | Same reason — admin wake manifest + prompt are not in scope. |
| `src/packages/cnos.cdd/skills/cdd/*/SKILL.md` | Sub 5B landed; not this cycle. |
| `src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md` | PR #466 / cnos#454 landed; not this cycle. |
| `src/packages/cnos.cds/orchestrators/cds-dispatch/prompt.md` | Sub 4 / cnos#483 landed; cycle SHOULD NOT need to touch the prompt body. |
| `src/packages/cnos.core/commands/install-wake/cn-install-wake` | Sub 5A landed; renderer is frozen for this cycle. |
| cnos#487 itself (as a smoke target) | OG-1 — bootstrap-δ-claimed (`status:in-progress`); smoke is a SEPARATE issue. |

## §5. AC-by-AC oracle (9 ACs from cnos#487; Stage 1 vs Stage 2 distinguished)

The 9 ACs are quoted abridged from cnos#487's body for traceability; the full body is the canonical source.

| AC | Stage | Oracle (runnable) |
|---|---|---|
| **AC1** — cnos#483's `activation_state` flipped to `"live"` + notes cite Sub 5A + 5B PRs | Stage 1 | `jq -r '.activation_state' src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.json` returns `live`; `jq -r '.activation_state_notes' ...` mentions `488` AND `489` (or `cnos#485` AND `cnos#486`). |
| **AC2** — `.github/workflows/cnos-cds-dispatch.yml` exists | Stage 1 | `test -e .github/workflows/cnos-cds-dispatch.yml && echo OK`. |
| **AC3** — Rendered = golden byte-identical | Stage 1 | `[ "$(sha256sum .github/workflows/cnos-cds-dispatch.yml \| cut -d' ' -f1)" = "$(sha256sum src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml \| cut -d' ' -f1)" ] && echo OK`. install-wake-golden CI extended to iterate over both goldens; both pass. |
| **AC4** — First claim fires green | Stage 2 | Actions run URL captured in the cycle/487 PR body's post-merge evidence section (or in α-closeout if added post-merge); job status = `success`. Verification: `gh run view <run-id> --json conclusion --jq .conclusion` returns `success`. |
| **AC5** — Claim sequence verified | Stage 2 | `gh issue view <smoke-N> --json labels --jq '.labels[].name'` shows `status:in-progress` (no `status:todo`) and preserves `protocol:cds` + `dispatch:cell`; `gh issue view <smoke-N> --json comments --jq '.comments[] \| select(.body \| contains("cds-dispatch"))'` shows a claim comment naming wake-name/run-id + `protocol: cds` + `head: <sha>`. |
| **AC6** — δ ran through γ/α/β to converge | Stage 2 | `ls .cdd/unreleased/<smoke-N>/` (on the smoke cell's cycle branch, OR on main post-merge of the smoke cell's PR) shows 6 artifacts: `gamma-scaffold.md`, `self-coherence.md`, `beta-review.md`, `alpha-closeout.md`, `beta-closeout.md`, `gamma-closeout.md` (PRA optional). Each is non-trivial (`wc -l` > 30 for scaffold; `grep -q '§R0' self-coherence.md` and `grep -q '§R0' beta-review.md`; `grep -qiE 'verdict.*(converge\|iterate)' beta-review.md`). |
| **AC7** — Cycle PR shipped | Stage 2 | PR URL referencing the smoke issue captured; `gh issue view <smoke-N> --json labels --jq '.labels[].name'` shows `status:review` (no `status:in-progress`). |
| **AC8** — No admin-wake regression | Stage 2 | Post-merge smoke: open issue titled `claude-wake AC3 re-verify <date>`; observe `cnos-agent-admin.yml` workflow run fires; `gh run list --workflow cnos-agent-admin.yml --json conclusion --jq '.[0].conclusion'` returns `success`. The admin wake does NOT claim any `dispatch:cell + protocol:cds` issue (admin's allowed-surfaces audit per `cnos-agent-admin.yml` prompt §"Allowed surfaces"; channel log shows admin's deferral to dispatch wake for cell-shaped directives). Close the smoke issue after evidence captured. |
| **AC9** — No cross-wake interference | Stage 2 | Smoke firing's job log (`gh run view <run-id> --log`) shows no concurrency-wait state attributable to `agent-admin-sigma`. Concurrency groups (`agent-admin-{agent}` vs `cds-dispatch-{agent}`) are per-protocol-scoped; both wakes can fire in the same hour without blocking. Empirical: both wakes' runs in the same hour observably succeed independently. |

**AC reconciliation note:** the operator's 5C directive (in the dispatch prompt) and cnos#487's issue body's 9 ACs may not align 1:1 (the issue body was drafted before the bootstrap-δ + two-stage shape was finalized). This scaffold reconciles them: AC1-AC3 + AC8/AC9 align with Stage 1 + post-merge admin regression check; AC4-AC7 are Stage 2's full path; AC8 also captures Stage 2's post-merge admin-wake regression check. β verifies each AC at its proper stage; AC4-AC7 are observed at Stage 2 (β re-engages post-merge to verify the smoke). See §13 friction note FN-6 on the two-stage β responsibility split.

## §6. Operator-named guardrails (OG-1 through OG-5)

Each OG carries an empirically β-checkable criterion. β verdicts must cite each OG's pass/fail explicitly.

### OG-1 — cnos#487 is not the smoke target

**Statement.** The cycle's substantive end-to-end smoke uses a SEPARATE issue from cnos#487. cnos#487 is already bootstrap-δ-claimed (`status:in-progress` + claim comment from `2026-06-23T05:22:05Z` per the GitHub API) so the live wake's verify gate REJECTS it (per `cnos.core/skills/agent/dispatch-protocol/SKILL.md` §2.2 step 3 verify gate: `status:todo` MUST be present; `status:in-progress` MUST NOT be present).

**β-checkable criterion.** α's α-prompt + the smoke issue design (§11) name the smoke as a SEPARATE issue. α-closeout (post-Stage-2) names the smoke issue number; it is NOT 487. cnos#487's lifecycle label at Stage 2 smoke firing time is `status:in-progress` (verifiable via `gh issue view 487 --json labels --jq '.labels[].name'`); the live cds-dispatch wake's firings observably do NOT claim cnos#487 (verifiable via the absence of a wake-fired claim comment on cnos#487 after Stage 2 firings; the only claim comment on cnos#487 is the bootstrap-δ one dated 2026-06-23).

### OG-2 — admin wake remains admin-only

**Statement.** Stage 1's workflow changes MUST NOT touch `.github/workflows/cnos-agent-admin.yml`, `src/packages/cnos.core/orchestrators/agent-admin/*`, or any other path that could cause an admin-wake regression. Stage 2's post-merge smoke procedure includes an admin-wake regression check (AC8).

**β-checkable criterion.** `git diff origin/main...HEAD --name-only` on the cycle/487 PR shows ZERO of: `.github/workflows/cnos-agent-admin.yml`, `src/packages/cnos.core/orchestrators/agent-admin/`. Verifiable mechanically:

```sh
git diff origin/main...HEAD --name-only | grep -E '(cnos-agent-admin|orchestrators/agent-admin)' && echo FAIL || echo OK
```

Should print `OK`. Stage 2: AC8 oracle exercised + admin wake observably continues firing post-merge.

### OG-3 — smoke exercises the real path, not bootstrap

**Statement.** The smoke MUST run through: label trigger (or scheduled sweep) → `cnos-cds-dispatch.yml` workflow → wake claims the smoke issue → δ wake-invoked mode → γ/α/β → `.cdd/unreleased/{smoke-N}/` artifacts → PR → `status:review`. The smoke MUST NOT be satisfied by running the bootstrap-δ path on the smoke issue. The whole point of 5C is proving the LIVE wake works.

**β-checkable criterion.** Stage 2 evidence captured in α-closeout: (a) Actions run URL is for `cnos-cds-dispatch.yml` (NOT for a manual bootstrap session); (b) claim comment author is the sigma bot (not the operator's account); (c) `.cdd/unreleased/{smoke-N}/` accumulation timestamps match the wake's firing window; (d) the smoke cell's cycle PR's commit-author attribution shows the wake's bot identity (`sigma@cnos.cn-sigma.cnos`); (e) no bootstrap-δ session ran the smoke cell's γ/α/β routing.

### OG-4 — failure handling explicit

**Statement.** If live dispatch misbehaves (double-claim, loop, regression of admin wake, claim of cnos#487 despite the bootstrap-δ lock), the cycle STOPS:
- Smoke issue moves to `status:blocked` with reason comment.
- If the workflow can double-claim or loop, disable/revert it (delete `.github/workflows/cnos-cds-dispatch.yml` from main via a follow-up PR; flip cnos#483's manifest back to `declaration-only`).
- Report from `.cdd` artifacts.
- Keep cnos#487 OPEN until repaired.
- Do NOT pretend the wave is complete; honest failure surfacing.

**β-checkable criterion.** α's α-prompt (§7 below) names each of the 4 failure symptoms (double-claim, loop, admin regression, cnos#487-claimed-despite-lock) with the specific action α takes for each. α-closeout (Stage 2) records which failure paths (if any) actually fired and the response. β-review post-Stage-2 verifies that the α response matched the §12 Rollback/disable plan.

### OG-5 — distinguish cap from failure

**Statement.** If smoke fails because of inherited I4/I5/I6-style CI debt → record as inherited cap (consistent with cnos#478's affordance). If failure is in claim, δ routing, artifact production, status transition, or any other surface this cycle introduces → it's a 5C failure (NOT a cap).

**β-checkable criterion.** α-closeout (Stage 2) classifies any failure observed: (a) failures attributable to claude-code-action limitations / inherited CI red status / substrate quirks unrelated to the wake-orchestration architecture → cap (operator-discretion to ship); (b) failures in the claim sequence, δ routing, artifact emission, or lifecycle transitions → 5C failure (block; OG-4 invoked). β-review verifies the classification: each failure has an explicit category + cited evidence (CI log line range, comment URL, or artifact-tree state).

## §7. α implementation prompt

**You are α (alpha) for cycle/487 — cds-dispatch render + activate + post-merge smoke (Sub 5C of cnos#467 wave-orchestration). Wave-goal-achievement cell.** This is a fresh Agent session; you have no prior context with this conversation. Read everything below before acting.

### α §1 — Your role

α implements. γ scaffolds (this document). β reviews. δ routes. Per operator directive: γ does not spawn α/β; δ dispatches every role. Your job is to: (Stage 1) flip the manifest + render + verify + open PR + signal review-ready; (Stage 2, post-merge) file the smoke issue + observe + capture evidence + classify outcome + close cnos#487.

### α §2 — Working directory + branch

`/home/user/cnos` — already cloned. Branch: `cycle/487` (created from `main@ff47c3600e9d2b804a8f2d76eeff4544e6a13c4d`). Commit + push to `cycle/487`. Open the PR against `main`.

### α §3 — Source-of-truth files to read BEFORE acting

In order:

1. **cnos#487 issue body** — via `mcp__github__issue_read` (method=get, owner=usurobor, repo=cnos, issue_number=487). The 9 ACs are canonical; this scaffold reconciles them with the two-stage shape (§5).
2. **This scaffold** — `.cdd/unreleased/487/gamma-scaffold.md` — your operational contract.
3. **cnos#483 cds-dispatch manifest** — `src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.json` — the file you'll flip.
4. **cnos#483 cds-dispatch prompt** — `src/packages/cnos.cds/orchestrators/cds-dispatch/prompt.md` — read-only context; you do NOT modify it.
5. **Sub 5A renderer** — `src/packages/cnos.core/commands/install-wake/cn-install-wake` — read-only context; you invoke it but do not modify it. Pay attention to the activation_state refusal logic (lines ~360-384) so you understand what changes when the manifest flips to `live`.
6. **Sub 5A's existing golden + CI guard** — `src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml` (current sha256 at base: `75e0406662ea12b3caa5ad01d8dbe6bdfb620d0c047fd910917a8e930855605e`) + `.github/workflows/install-wake-golden.yml`.
7. **Sub 5B δ wake-invoked mode** — `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` §9 (lines ~400-515) — read in full; the live wake will invoke this contract; you don't change it but its shape determines what evidence Stage 2 captures.
8. **Predecessor closeouts** — `.cdd/unreleased/485/{alpha,beta,gamma}-closeout.md` + `.cdd/unreleased/486/{alpha,beta,gamma}-closeout.md` — recommendations carry forward (per cycle/486 T-486-15 "γ reads predecessor closeouts as γ-template default"; α also benefits from them).
9. **Admin wake reference** — `.github/workflows/cnos-agent-admin.yml` — read-only; understand the shape so you can verify your Stage 1 diff does NOT touch it (OG-2).
10. **Dispatch protocol skill** — `src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md` — the contract the live wake operationalizes; you don't modify it; you reference it in your understanding of OG-1's bootstrap-δ lock.

### α §4 — Stage 1 implementation (pre-merge; the cycle/487 PR's work)

#### α §4.1 — Flip the manifest (AC1)

Edit `src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.json`:

- Change `"activation_state": "declaration-only"` → `"activation_state": "live"`.
- Update `"activation_state_notes"` to record that preconditions HAVE landed. Suggested updated text (you may refine):

> This provider is LIVE in production substrate as of cnos#487 (Sub 5C of cnos#467 wake-orchestration wave). Preconditions satisfied: (a) cnos#454 dispatch-protocol skill landed on main (PR #466); (b) cnos#467 Sub 5A landed the renderer extension consuming role:dispatch + protocol + selector + dispatch-shape output_contract fields + the issues_labeled_selector_match trigger (cnos#485 / PR #488); (c) cnos#467 Sub 5B landed the δ wake-invoked mode amendment in cnos.cdd's delta SKILL (cnos#486 / PR #489). The corresponding rendered substrate artifact lives at `.github/workflows/cnos-cds-dispatch.yml` per `cn install-wake cds-dispatch --out .github/workflows/cnos-cds-dispatch.yml`; the renderer no longer refuses (per wake-provider/SKILL.md §3.10 — `activation_state == "live"` is the production state).

Verify: `jq -r '.activation_state' src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.json` returns `live`; `jq -r '.activation_state_notes' src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.json | grep -E '488|489|cnos#485|cnos#486'` returns at least one match.

#### α §4.2 — Render the substrate artifact (AC2 + AC3)

Run: `./src/packages/cnos.core/commands/install-wake/cn-install-wake cds-dispatch --out .github/workflows/cnos-cds-dispatch.yml`

Expected exit code: 0. The renderer no longer refuses because `activation_state == "live"`. The output file should be byte-identical to the synthetic-live golden (`src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml`, sha256 `75e0406662ea12b3caa5ad01d8dbe6bdfb620d0c047fd910917a8e930855605e` at branch base).

Verify byte-identical:

```sh
sha256sum .github/workflows/cnos-cds-dispatch.yml src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml
```

Both should print the same hash. If they DON'T match, see §6 Design fork D1 below.

#### α §4.3 — Reconcile golden if render differs (Design fork D1)

**Design fork D1 — render-vs-golden divergence.** Sub 5A rendered the golden with `--activation-state-override live` against a `declaration-only` manifest. The post-flip render is on a `live` manifest with no override. The renderer's `activation_state_effective` value is `live` either way; the rendered output SHOULD be byte-identical. If it ISN'T:

- **Option A — Re-render the golden + commit:** the production render IS the new truth; the golden was synthetic. Run `./src/packages/cnos.core/commands/install-wake/cn-install-wake cds-dispatch` (default `--out` is the golden path) to regenerate the golden; commit both files; document the divergence in self-coherence §Design.
- **Option B — Investigate divergence + decide:** if the divergence is non-trivial (e.g. renderer behavior somehow depends on the `activation_state` source field, not just the effective value), STOP, surface to δ, do not commit a divergent state.

γ recommends Option A IF the divergence is purely cosmetic (whitespace, comment-string presence/absence) AND the renderer code path is the same. γ recommends Option B IF the divergence touches the workflow's substantive structure (different `if:` gate; different trigger; different concurrency group). α documents which option in self-coherence §Design.

#### α §4.4 — Update install-wake-golden CI (Design fork D2)

The Sub 5A CI guard (`.github/workflows/install-wake-golden.yml`) contains an AC5 step: "declaration-only refusal (cds-dispatch without override)". Post-flip, that step's assertion (`cn install-wake cds-dispatch` exits 3) FAILS because the manifest is now `live` (renderer exits 0).

**Design fork D2 — AC5 step disposition.**

- **Option A — Remove the AC5 step:** no second `declaration-only` manifest exists post-cycle (cds-dispatch was the only one). Removing the step is the minimal change. Risk: future declaration-only manifests (cdr-dispatch, observer roles per cnos#470) would need to re-introduce the assertion.
- **Option B — Re-point AC5 at a synthetic fixture:** create a test-fixture manifest under `src/packages/cnos.core/commands/install-wake/test-fixtures/` (or wherever convenient) with `activation_state: "declaration-only"`, point the AC5 step at it via `--manifest <fixture-path>`. Preserves the assertion. Risk: adds a new file; slight maintenance overhead.
- **Option C — Re-point AC5 at the agent-admin manifest with `--activation-state-override declaration-only`:** uses an existing manifest + the override flag in the OTHER direction (forcing a live manifest to be treated as declaration-only at render time). The renderer's `--activation-state-override` flag accepts any value; the refusal check is `!= "live"`. Risk: relies on the override flag's semantics — slightly brittle.

γ recommends **Option B** (synthetic fixture) as the cleanest — preserves the assertion's purpose (the renderer refuses declaration-only manifests), avoids surface contamination with the live manifest, and gives a clean fixture for any future declaration-only test. α may choose differently and document in self-coherence §Design with rationale.

Also update the install-wake-golden CI's `paths:` filters if your fork-D2 choice introduces a new fixture path that should trigger the workflow.

#### α §4.5 — Verify byte-identical render-vs-golden (CI guard)

After the fork-D1 + fork-D2 resolutions, run install-wake-golden CI locally (or via push + observe). The diff step should pass; the idempotence steps should pass; the substrate-shape steps should pass; the AC7/AC8 audit steps should pass; the AC5 step should pass per your fork-D2 choice.

#### α §4.6 — Document Stage 2 procedure in PR body

The cycle/487 PR body MUST contain a "Post-merge smoke procedure" section that names — for the operator — exactly what α will do at Stage 2 (or, if Stage 2 happens before α can return, what the operator should observe). This is operator-readable; it gives the operator a checklist.

Suggested PR body structure:

```
## Summary
Flips cnos#483's cds-dispatch manifest from declaration-only to live; renders + commits .github/workflows/cnos-cds-dispatch.yml; extends install-wake-golden CI for post-flip coverage. After merge, the live cds-dispatch wake claims a separately-filed smoke issue and runs it end-to-end through δ wake-invoked mode → γ/α/β → status:review.

Closes cnos#487 only after Stage 2 smoke completes (post-merge).

## Stage 1 (this PR — pre-merge)
- AC1: activation_state flipped to live (jq -r '.activation_state' ...)
- AC2: .github/workflows/cnos-cds-dispatch.yml exists
- AC3: sha256(rendered) == sha256(golden) — verified locally + install-wake-golden CI
- OG-2: git diff --name-only does NOT include cnos-agent-admin paths

## Stage 2 (post-merge)
1. File smoke issue with title "smoke: cds-dispatch end-to-end proof (cnos#487 Sub 5C)" + labels [dispatch:cell, protocol:cds, status:todo]
2. Wait for cnos-cds-dispatch.yml to fire (label-event responsive or next 15-min scheduled sweep)
3. Verify claim comment + label transition (AC4 + AC5)
4. Wait for δ → γ/α/β to converge; verify .cdd/unreleased/{smoke-N}/ accumulates (AC6)
5. Verify smoke cell PR opened + status:review reached (AC7)
6. File admin-wake regression check (AC8): open issue titled "claude-wake AC3 re-verify <date>"; observe cnos-agent-admin.yml fires green + does NOT claim any protocol:cds issue
7. Verify no cross-wake concurrency interference (AC9)
8. Close cnos#487 after all 9 ACs verified
```

#### α §4.7 — Open the PR + signal review-ready

After commits, push branch + open PR. Add a self-coherence file `.cdd/unreleased/487/self-coherence.md` with:

- Frontmatter: cycle, branch, base sha, R0 review-ready sha (commit hash where you signal ready).
- **§ACs section:** per-AC verification table (9 rows). For Stage 1 ACs (AC1-AC3 + AC8/AC9 Stage 1 portions), populate with actual verification output. For Stage 2 ACs (AC4-AC7, AC8/AC9 Stage 2 portions), mark "deferred to Stage 2 — post-merge evidence will be appended" + name the planned procedure.
- **§Design section:** record your fork-D1 + fork-D2 decisions with rationale.
- **§Per-CI-step bash-e audit:** populate the table per cnos#478 mechanical-injection format (see §9 below).
- **§Variable consistency table:** populate per cycle/486 T-486-1 (see §10 below).
- **§Self-check section:** confirm OG-1 through OG-5 by your own audit.
- **§Friction notes:** any FN-α-N items for triage in α-closeout.
- **§R0 review-ready signal:** explicit "α R0 review-ready at commit <sha>; β may engage" line.

Append the review-readiness signal as a discrete commit (per cycle/485 + cycle/486 pattern). δ reads the branch state to dispatch β.

### α §5 — Stage 2 implementation (post-merge; observation/evidence-gathering)

Stage 2 runs AFTER the cycle/487 PR merges. α re-engages (a new α session, or this α session if still active) to:

#### α §5.1 — File the smoke issue (concrete; per §11 below)

Use `mcp__github__issue_write` (or `gh issue create`) with:

- **Title:** `smoke: cds-dispatch end-to-end proof (cnos#487 Sub 5C)`
- **Labels:** `dispatch:cell, protocol:cds, status:todo`
- **Body:** per §11 below (verbatim).

Capture the issue number — call it `{smoke-N}`. Record it in α-closeout.

#### α §5.2 — Observe the wake firing

The live `cnos-cds-dispatch.yml` should fire on EITHER:

- The `issues: labeled` event (responsive path; fires within seconds when the `status:todo` label is applied — assuming label-application is one of the events triggering the wake).
- The next scheduled sweep (backstop; cron runs at `:08`, `:23`, `:38`, `:53` per the rendered cron slots — max 15-minute wait).

**Polling procedure.** After filing the smoke issue, poll for the workflow firing:

```sh
# Watch for new cnos-cds-dispatch.yml runs (poll every ~60s):
gh run list --workflow cnos-cds-dispatch.yml --limit 5 --json databaseId,status,conclusion,createdAt,event,headBranch
```

Wait UP TO 30 minutes (covers two scheduled sweeps). If no firing within 30 minutes, escalate per OG-4 (the wake may not be triggering; investigate the workflow's `on:` block and `if:` gate; if structural, follow §12 Rollback plan).

#### α §5.3 — Capture evidence at each phase

For AC4 (firing green):

```sh
gh run view <run-id> --json conclusion,databaseId,htmlUrl --jq '{conclusion, url: .htmlUrl}'
```

For AC5 (claim sequence):

```sh
gh issue view <smoke-N> --json labels --jq '.labels[].name'  # should include status:in-progress (NOT status:todo)
gh issue view <smoke-N> --json comments --jq '.comments[] | select(.body | contains("cds-dispatch") or contains("claimed_by"))'
```

For AC6 (δ ran through γ/α/β):

```sh
git fetch origin
# Find the smoke cell's cycle branch (the wake creates cycle/{smoke-N} or similar):
git branch -a | grep "cycle/${SMOKE_N}"
# Or list .cdd/unreleased/{smoke-N}/ on the branch:
git ls-tree -r origin/cycle/${SMOKE_N} -- ".cdd/unreleased/${SMOKE_N}/"
```

For AC7 (PR shipped + status:review):

```sh
gh pr list --search "cnos#${SMOKE_N}" --json url,state,headRefName --jq '.[]'
gh issue view <smoke-N> --json labels --jq '.labels[].name'  # should include status:review
```

For AC8 (admin-wake regression check):

```sh
# Open the AC3 re-verify issue:
gh issue create --title "claude-wake AC3 re-verify $(date -u +%Y-%m-%d)" --body "Post-merge regression check for cnos#487 Sub 5C. Admin wake should fire on this title; should NOT claim any protocol:cds cell. Auto-close after observation."
# Wait for cnos-agent-admin.yml firing:
gh run list --workflow cnos-agent-admin.yml --limit 3
# Verify it fired green:
gh run view <admin-run-id> --json conclusion --jq .conclusion
# Verify admin wake did NOT claim any protocol:cds cell:
gh issue list --label "protocol:cds,status:in-progress" --json number  # check none were transitioned by the admin wake in this window
# Close the AC3 re-verify issue:
gh issue close <ac3-reverify-N>
```

For AC9 (cross-wake non-interference):

```sh
# Inspect the smoke firing's job log for concurrency waits:
gh run view <smoke-run-id> --log | grep -iE "concurrency|wait|queued" || echo "no concurrency-wait state"
```

#### α §5.4 — Outcome classification

Per OG-5, classify the smoke outcome:

- **Success path:** all 9 ACs pass. α writes α-closeout (`.cdd/unreleased/487/alpha-closeout.md`) capturing Stage 1 + Stage 2 evidence + URLs + SHAs. δ dispatches β for β-closeout + γ for γ-closeout. δ closes cnos#487.
- **Cap path:** smoke fails due to inherited CI debt / claude-code-action limitations / substrate quirks unrelated to the wake-orchestration architecture. α-closeout names the cap class explicitly + cites the inherited-cap inventory item (I4/I5/I6 or equivalent). δ + operator decide whether to ship despite the cap (consistent with cnos#478 affordance).
- **5C failure path:** smoke fails in the claim sequence, δ routing, artifact emission, or lifecycle transitions. α invokes OG-4: smoke issue → `status:blocked` with reason; cycle/487 PR stays merged but cnos#487 stays OPEN; α-closeout names the failure + recommends rollback steps per §12.

#### α §5.5 — Closeouts (Stage 2-aware)

α writes `.cdd/unreleased/487/alpha-closeout.md` capturing both stages. Then δ dispatches β + γ for their closeouts. All three closeouts MUST be authored AFTER Stage 2 completes (or AFTER OG-4 invokes if a failure path triggers) so they can capture the full picture. δ then closes cnos#487 with a comment naming the closing PR + the smoke cell's cycle PR + the artifact set paths.

### α §6 — Iteration discipline

- Stage 1 R[N] iteration is standard: α implements at R0; if β returns `iterate`, δ re-dispatches α with β's findings; α implements R[N+1] against the same branch with §R[N] section in self-coherence + the relevant fix.
- Stage 2 is observation/evidence-gathering (no R[N] iteration unless smoke fails). If smoke fails (OG-4), α captures evidence + invokes rollback procedure; this is NOT an "α R[N+1]" event but a defer-to-δ-and-operator event.

### α §7 — What α MUST NOT do

(Recap from §4.3 + OGs.)

- Touch `.github/workflows/cnos-agent-admin.yml` or any agent-admin orchestrator file.
- Touch `src/packages/cnos.cdd/skills/cdd/*/SKILL.md` (Sub 5B landed).
- Touch `src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md` (cnos#454 landed).
- Touch `src/packages/cnos.cds/orchestrators/cds-dispatch/prompt.md` (cnos#483 landed; cycle SHOULD NOT need to touch the prompt).
- Touch `src/packages/cnos.core/commands/install-wake/cn-install-wake` (Sub 5A landed; renderer frozen).
- File the smoke issue with title containing "claude-wake" (that would trigger admin wake, not dispatch wake).
- Treat cnos#487 as the smoke target (OG-1; bootstrap-δ-claimed).
- Bypass the live wake by running bootstrap-δ on the smoke issue (OG-3; defeats 5C's purpose).
- Continue past OG-4 failure symptoms (stop + invoke rollback + keep cnos#487 OPEN).
- Manually re-do the smoke if it fails — the smoke is the LIVE wake's test; if it fails, the failure is real, not a retry opportunity.

### α §8 — Termination

α's Stage 1 work terminates when: cycle/487 branch pushed; cycle/487 PR opened against main; self-coherence.md committed with R0 review-ready signal; β can engage. α's Stage 2 work terminates when: AC4-AC7 + AC8 (post-merge admin regression) + AC9 are all verified (or OG-4 fires); α-closeout written; δ can dispatch β-closeout + γ-closeout.

## §8. β review prompt

**You are β (beta) for cycle/487 — cds-dispatch render + activate + post-merge smoke (Sub 5C of cnos#467 wave-orchestration). Wave-goal-achievement cell.** This is a fresh Agent session; you have no prior context with this conversation. Read everything below before acting.

### β §1 — Your role

β reviews. γ scaffolds (this document). α implements (§7 above). δ routes. Per operator directive: γ does not spawn α/β; δ dispatches every role. Your job is to: independently re-verify each AC, audit each OG, run the per-CI-step bash-e audit, populate the variable consistency table, check doctrine + cross-skill consistency + CI evidence + non-goal-touched-set, write `.cdd/unreleased/487/beta-review.md` with a verdict (converge or iterate) per round.

### β §2 — Working directory + branch

`/home/user/cnos` — already cloned. Branch: `cycle/487` (α implements here). PR: opens against `main`. Read α's branch HEAD; do NOT push commits to the branch (β writes only `.cdd/unreleased/487/beta-review.md`).

### β §3 — Source-of-truth files (same as α; read in order)

1. cnos#487 issue body.
2. This scaffold (`.cdd/unreleased/487/gamma-scaffold.md`).
3. α's `.cdd/unreleased/487/self-coherence.md` at the R[N] review-ready commit.
4. cnos#483 manifest (post-flip) + prompt (unchanged).
5. Sub 5A renderer (`cn-install-wake`; unchanged).
6. The rendered workflow (`.github/workflows/cnos-cds-dispatch.yml`).
7. The golden (`src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml`).
8. The CI guard (`.github/workflows/install-wake-golden.yml`; post α's edits).
9. Sub 5B δ wake-invoked mode skill (`src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` §9).
10. Predecessor closeouts (`.cdd/unreleased/{485,486}/{alpha,beta,gamma}-closeout.md`).
11. Admin wake reference (`.github/workflows/cnos-agent-admin.yml`).
12. Dispatch protocol skill (`src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md`).

### β §4 — Per-AC verification table (9 ACs; populate per round)

For each of the 9 ACs in §5 above, β independently runs the oracle (do NOT trust α's self-coherence verification; re-run from scratch). Record:

| AC | Stage | Oracle command | Result | Pass/fail | Notes |
|---|---|---|---|---|---|
| AC1 | 1 | `jq -r '.activation_state' src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.json` | (paste output) | ✅/❌ | (notes) |
| AC2 | 1 | `test -e .github/workflows/cnos-cds-dispatch.yml && echo OK` | (paste) | | |
| AC3 | 1 | `sha256sum` rendered vs golden | (paste both hashes) | | |
| AC4 | 2 | (deferred at R0; populated post-merge) | | | |
| AC5 | 2 | (deferred at R0; populated post-merge) | | | |
| AC6 | 2 | (deferred at R0; populated post-merge) | | | |
| AC7 | 2 | (deferred at R0; populated post-merge) | | | |
| AC8 | 2 | (deferred at R0; populated post-merge) | | | |
| AC9 | 2 | (deferred at R0; populated post-merge) | | | |

At R0 (Stage 1 review), only AC1-AC3 (and parts of OG-2/OG-5 that are Stage 1) can be verified. β converges Stage 1 ON ACs 1-3 + OGs 1-5 Stage 1 portions; the Stage 2 ACs are deferred to a post-merge β re-engagement (per §13 friction note FN-6).

### β §5 — Per-OG verification

For each OG (OG-1 through OG-5):

| OG | Statement | β-checkable criterion | Result | Pass/fail | Notes |
|---|---|---|---|---|---|
| OG-1 | cnos#487 is not the smoke target | α-prompt + smoke design name a SEPARATE issue; cnos#487 stays `status:in-progress` (bootstrap-δ lock) | (β reads self-coherence + smoke design + the issue API) | ✅/❌ | |
| OG-2 | admin wake remains admin-only | `git diff origin/main...HEAD --name-only \| grep -E '(cnos-agent-admin\|orchestrators/agent-admin)'` returns empty | (β runs the diff) | | |
| OG-3 | smoke exercises the real path, not bootstrap | Stage 2 evidence in α-closeout names Actions run URL + bot identity + wake-author signatures | (deferred to Stage 2) | | |
| OG-4 | failure handling explicit | α-prompt names each of 4 symptoms + actions; α-closeout records what fired | (β verifies α-prompt has the 4-symptom enumeration at R0; deferred Stage 2 verification) | | |
| OG-5 | distinguish cap from failure | α-closeout classifies outcomes with cited evidence | (deferred to Stage 2) | | |

### β §6 — Per-CI-step bash-e audit table (MANDATORY — cycle IS CI-touching)

Cycle/487 touches CI: it adds `.github/workflows/cnos-cds-dispatch.yml` (a CI workflow with `run:` blocks executed inside the claude-code-action carrier — see note below) and may modify `.github/workflows/install-wake-golden.yml` (definitely has `run:` blocks). β populates the audit table per cnos#478 mechanical-injection format.

**Note on claude-code-action `run:` exposure.** The rendered `cnos-cds-dispatch.yml` wraps its work in `anthropics/claude-code-action@v1`; the action's internal shell invocations are not directly `run:` blocks in the workflow YAML. β audits the entry-point steps (`actions/checkout@v4` `with:` parameters; the `claude-code-action@v1` `with:` parameters; any direct `run:` blocks if α adds them) and notes the indirection in the table.

Columns (per cnos#478 / cycle/485 + cycle/486):

| # | Step name | Line range | Command substitutions or pipelines | Guarded? | bash-e exit on intended-success input | Notes |
|---|---|---|---|---|---|---|
| 1 | (e.g. `actions/checkout@v4` in cnos-cds-dispatch.yml) | (line range in the rendered yml) | (none — `uses:` step) | (n/a) | (n/a) | Verify token binding matches expected pattern; verify no shell exposure. |
| 2 | (e.g. `claude-code-action@v1`) | | | | | The action's internal shell is opaque to this audit; β notes the indirection. |
| 3 | (install-wake-golden.yml steps modified by α — populate per α's fork-D2 choice) | | | | | |
| ... | | | | | | |

β re-runs the AC7/AC8 audit grep checks already in `install-wake-golden.yml` (the renderer-side authority audits; lines ~225-241 of the CI guard at base) — those should continue to pass since `cn-install-wake` is unchanged.

### β §7 — Variable consistency table (cycle/486 T-486-1)

Cycle/487 introduces multiple references to the same variables/fields across the manifest, the rendered YAML, and the CI guard. β populates the table per cycle/486 §2.1's "variable consistency table" pattern. Suggested anchor variables:

- `activation_state` value — appears in: (a) manifest field; (b) `activation_state_notes` referent; (c) renderer's `activation_state_effective` logic; (d) CI guard's AC5 step (post α's fork-D2 choice); (e) PR body's AC1 oracle.
- Cron slot values (`8 23 38 53`) — appear in: (a) renderer's `substrate_cron_slots`; (b) rendered admin wake YAML; (c) rendered cds-dispatch YAML; (d) potential timing-related notes in α-closeout (Stage 2 polling window).
- Concurrency group name (`cds-dispatch-{agent}`) — appears in: (a) manifest's `concurrency_intent.group`; (b) rendered YAML's `concurrency.group`; (c) install-wake-golden's substrate-shape check; (d) AC9's interference oracle.
- Selector labels (`dispatch:cell`, `protocol:cds`, `status:todo`) — appear in: (a) manifest's `selector.include`; (b) rendered YAML's `if:` gate; (c) smoke issue's labels (Stage 2); (d) AC5's claim sequence oracle; (e) OG-1's bootstrap-δ-lock reasoning.

For each variable: β walks each occurrence and confirms the value matches the expected meaning in its local context. Catches off-by-one / wrong-context / shape-drift errors that grep oracles miss (per cycle/486's §9.5 `§R[N-1]` off-by-one empirical case).

### β §8 — Doctrine consistency

Check for stale doctrine (per cnos#467's authoritative doctrine-correction at top):

```sh
grep -rE 'protocol:cdd|cdd-dispatch|cnos\.cdd as concrete protocol' src/packages/cnos.cds/orchestrators/cds-dispatch/ .github/workflows/cnos-cds-dispatch.yml .cdd/unreleased/487/
```

Should return zero matches (or matches only in cited legacy quotes, with clear context). If matches: 5C failure (β flags as iterate).

Verify:
- `cnos.cds` = concrete software-development protocol package
- `protocol:cds` = label qualifier
- `cds-dispatch` = wake name
- `cnos.cdd` = generic cell-runtime framework only

### β §9 — Cross-skill consistency

The rendered `cnos-cds-dispatch.yml`'s embedded prompt references:
- Sub 5B δ wake-invoked mode (`src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` §9) — verify the prompt's references match the LANDED §9 amendment shape (input contract, routing sequence, R[N] discipline, artifact contract, return tokens). Specifically check the prompt's "Invoke δ in wake-invoked mode" section (rendered lines ~115-128 of the golden) names contracts that exist in the merged §9.
- cnos#483 manifest — the rendered YAML's job-level `if:` gate contains exactly the labels in the manifest's `selector.include` (`dispatch:cell`, `protocol:cds`, `status:todo`). β re-reads `jq -r '.selector.include[]' src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.json` and confirms each label appears in the `if:` gate.
- cnos#454 dispatch-protocol — the prompt's claim sequence (rendered lines ~75-110 of the golden) matches dispatch-protocol/SKILL.md §2.2 step list (1-6).

### β §10 — CI evidence

Read the install-wake-golden CI status on the cycle/487 PR's HEAD:

```sh
gh pr checks <pr-number>
```

Verify install-wake-golden runs green (or document any inherited-cap failures that are NOT introduced by this cycle — cite the cap inventory). Any new CI red attributable to cycle/487's changes is a 5C failure (β flags iterate).

### β §11 — Non-goal verification

Cycle/487 has explicit non-goals from §13 + cnos#487 issue body + the dispatch prompt. β verifies α did NOT scope-creep into:

- CDR / CDW dispatch wakes (other protocols)
- NIM / OpenAI / alternative substrate carriers
- New label doctrine
- Broad runtime cleanup unrelated to cds-dispatch
- Unrelated CI repair (the inherited I4/I5/I6 caps stay caps)
- Private registration changes (PR #465 merged independently)
- Bundling Sub 6 (cycle-complete admin wake reading) into this cycle

Check via `git diff origin/main...HEAD --name-only` and verify no files outside §4.1's in-scope list appear.

### β §12 — Friction notes + verdict

Per round R[N], β appends a `§R[N]` section to `beta-review.md` with:

- Per-AC table (§4 above)
- Per-OG table (§5)
- Per-CI-step audit (§6)
- Variable consistency table (§7)
- Doctrine consistency result (§8)
- Cross-skill consistency result (§9)
- CI evidence (§10)
- Non-goal verification (§11)
- Friction notes (FN-β-N items)
- **Verdict: converge | iterate.** If iterate, list specific findings with line numbers + suggested resolutions.

β does NOT manufacture findings to look thorough (per cycle/485 closeout T11 — "do NOT manufacture findings" discipline). Honest converge is honest converge. β verdicts are NOT gated on inherited-cap CI failures (per cycle/485 T12 / β-closeout §6 second bullet — pre-existing reds traceable to a named inventory are converge-compatible).

### β §13 — Stage 2 re-engagement

After the cycle/487 PR merges + Stage 2 runs (smoke fires + δ routes + artifacts accumulate + smoke PR ships + status:review reached), δ MAY dispatch β again to verify AC4-AC7 + AC8/AC9 Stage 2 portions. β appends a `§Stage 2` section to `beta-review.md` with:

- AC4-AC9 verification (with run URLs, comment links, artifact tree state)
- OG-3 + OG-5 verification (smoke exercises real path; outcome classified)
- Verdict: converge | failure-cap | failure-5C

If Stage 2 re-engagement happens, this is a different β-review section from R[N] iterations on the Stage 1 PR. δ documents the routing in the closeouts.

### β §14 — Termination

β R[N] terminates when `beta-review.md §R[N]` is committed + pushed to cycle/487. β's Stage 2 re-engagement terminates when `beta-review.md §Stage 2` is committed + pushed to cycle/487 (or to the smoke cell's branch, per δ's routing decision).

## §9. Per-CI-step bash-e audit (γ pre-scan; α populates fully)

cycle/487 IS CI-touching. γ pre-scans the surfaces:

- The new `.github/workflows/cnos-cds-dispatch.yml` (rendered output; structure mirrors `cnos-agent-admin.yml` per Sub 5A's renderer); the active surface is the embedded `prompt:` body (which is the wake's prompt, audited at the prompt level, not the bash-e level) + the `with:` parameters on `actions/checkout@v4` + `anthropics/claude-code-action@v1`. No direct `run:` blocks in the rendered YAML.
- `.github/workflows/install-wake-golden.yml` (α modifies per fork-D2). The existing 12 steps (see Sub 5A's install-wake-golden) MUST be re-audited if α adds/modifies any.

**Pre-scan note.** The cycle/485 golden render command (`cn install-wake cds-dispatch --activation-state-override live`) is removed post-flip (renderer no longer needs the override since manifest is `live`). α's fork-D2 choice determines exact post-flip structure; α populates the full 3-column table in self-coherence §Per-CI-step.

## §10. Variable consistency table (γ pre-scan; α populates fully; β re-walks independently)

Pre-scan of variables that appear in multiple places:

| Variable | Occurrences | Expected consistency |
|---|---|---|
| `activation_state` field value | (a) manifest at `src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.json`; (b) `activation_state_notes` referent; (c) renderer's `activation_state_effective` runtime; (d) CI guard's AC5 step (fork-D2 dependent); (e) PR body's AC1 oracle | All occurrences see `live` post-flip; the renderer's refusal logic is bypassed; the CI guard's AC5 step is repointed per fork-D2. |
| Cron slots `8 23 38 53` | (a) renderer's `substrate_cron_slots` (line 423 of `cn-install-wake`); (b) `cnos-agent-admin.yml` `schedule:` block; (c) `cnos-cds-dispatch.yml` `schedule:` block; (d) Stage 2 polling window math (max 15-min wait) | Identical 4 cron values in all three YAML schedule blocks; α-closeout's Stage 2 procedure cites the 15-min sweep interval correctly. |
| Concurrency group `cds-dispatch-{agent}` → `cds-dispatch-sigma` | (a) manifest `concurrency_intent.group` (line 97 of wake-provider.json); (b) rendered YAML `concurrency.group` (line 30 of golden); (c) install-wake-golden substrate-shape check (line 139 of CI guard); (d) AC9 oracle | All occurrences resolve to `cds-dispatch-sigma` (with `{agent}` → `sigma`). AC9 verifies no cross-wake interference with `agent-admin-sigma` (admin's group). |
| Selector labels `dispatch:cell`, `protocol:cds`, `status:todo` | (a) manifest `selector.include` (lines 12-14 of wake-provider.json); (b) rendered YAML `if:` gate (line 35 of golden); (c) smoke issue labels (Stage 2); (d) AC5 oracle; (e) OG-1 lock reasoning | All 3 labels appear in `selector.include`; rendered `if:` gate contains `contains(github.event.label.name, '<label>')` for each; smoke issue applies all 3 at filing; AC5 oracle verifies `status:in-progress` after claim (NOT `status:todo`). |
| `bot_name`/`bot_id` (`sigma@cnos.cn-sigma.cnos` / `41898282`) | (a) `agent_bot_name`/`agent_bot_id` helpers in renderer (lines 127-138 of `cn-install-wake`); (b) rendered cnos-agent-admin.yml lines 46-47; (c) rendered cnos-cds-dispatch.yml lines 46-47; (d) Stage 2 OG-3 verification (claim comment author = bot, not operator) | Identical bot identity across both rendered wakes; Stage 2 evidence confirms claim comment is by `sigma`, not the operator. |

β re-walks this table at review time; updates with α's fork-D2 outcome for `activation_state`. Per cycle/486 T-486-1, the table catches semantic drift the grep oracles miss.

## §11. Smoke issue design (γ specifies; α files post-merge)

The smoke issue is intentionally TINY and LOW-RISK. The substantive work is the workflow firing + δ routing + artifact production, not the issue's task content.

### §11.1 — Concrete issue specification

**Title (verbatim — α uses this exact title):**
```
smoke: cds-dispatch end-to-end proof (cnos#487 Sub 5C)
```

**Labels (apply ALL three at filing):**
```
dispatch:cell
protocol:cds
status:todo
```

**Body (verbatim — α uses this body):**
```markdown
## Mode

`smoke-cell` — minimal-scope smoke cell proving the live `cnos-cds-dispatch.yml` workflow claims an issue end-to-end, δ routes γ/α/β, `.cdd/unreleased/{N}/` accumulates the canonical artifact set, and the cycle PR ships the cell's work to `status:review`. Filed per cnos#487 Sub 5C scaffold §11.

## Gap

cnos#487 (Sub 5C) flipped cds-dispatch's `activation_state` to `live` and committed the rendered substrate artifact at `.github/workflows/cnos-cds-dispatch.yml`. The renderer + manifest + workflow YAML are in production. What remains is verifying that the LIVE wake actually executes a real cell end-to-end. This issue is the smoke target — a separate, tiny, low-risk cell that the live wake claims, δ routes, γ/α/β execute, and the substrate-side firing produces the expected artifact set and lifecycle transitions.

## Impact

After this cell lands:
- A real protocol:cds cell has been claimed by `cnos-cds-dispatch.yml`'s first firing (post-cnos#487-merge), proving the wake's claim sequence per cnos#454 dispatch-protocol works in production substrate.
- δ wake-invoked mode (Sub 5B / cnos#486) has been exercised by a real wake firing (not bootstrap-δ), proving the §9 contract works end-to-end.
- `.cdd/unreleased/{this-issue-number}/` carries the canonical artifact set per cnos#483's `output_contract.artifact_class_taxonomy` (6 artifacts: γ-scaffold, self-coherence, β-review, α/β/γ closeouts).
- cnos#467's wave-goal-achievement is observable in production substrate for the first time.
- cnos#487 closes once the smoke verifies (the cnos#487-level closeouts capture this cell's evidence).

## Source of truth

| Surface | Canonical source | Notes |
|---|---|---|
| Live cds-dispatch wake | `.github/workflows/cnos-cds-dispatch.yml` (post cnos#487 merge) | The firing target |
| cds-dispatch manifest | `src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.json` | `activation_state: live` post cnos#487 |
| δ wake-invoked mode | `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` §9 | Routing contract |
| Dispatch protocol | `src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md` | Claim sequence |
| Issue-shape skill | `src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` | Cell contract this body conforms to |

## Constraints / non-goals

**In scope (the substantive task — intentionally tiny):**
- Create a single file `docs/gamma/smoke/cds-dispatch-smoke-YYYYMMDD.md` (replace YYYYMMDD with the cycle date) recording the workflow firing evidence:
  - Workflow run URL for the firing that claimed this issue
  - Claim comment quoted verbatim (wake-name/run-id + protocol + claimed_at + head SHA)
  - Per-R[N] artifact paths under `.cdd/unreleased/{this-issue-number}/`
  - Final β converge verdict
  - Status transition timestamps (`status:todo → status:in-progress` at claim; `status:in-progress → status:review` at converge)
- The file content IS the smoke task; the substantive value is the workflow + δ + γ/α/β producing it through the live path.

**Out of scope:**
- Any production code change
- Any other doc/test/skill change
- Multi-protocol exercises (this cell exercises `protocol:cds` only)
- Stress-testing the wake (one firing, one cell — that's the smoke)

## Acceptance criteria

1. **Live `cnos-cds-dispatch.yml` claims this issue.** Workflow run conclusion = `success`; claim comment present per the dispatch-protocol §2.2 step 4 format.
2. **δ routes γ/α/β successfully.** `.cdd/unreleased/{this-issue-number}/` contains: `gamma-scaffold.md`, `self-coherence.md` (with §R0 section minimum), `beta-review.md` (with §R0 verdict — converge or iterate), `alpha-closeout.md`, `beta-closeout.md`, `gamma-closeout.md`. PRA optional.
3. **Substantive smoke task complete.** `docs/gamma/smoke/cds-dispatch-smoke-YYYYMMDD.md` exists and contains all 5 evidence elements (run URL; claim comment quote; per-R[N] paths; converge verdict; status transition timestamps).
4. **Cycle PR opened + status:review reached.** A PR scoped to this issue references the cell; the issue's label set transitioned to `status:review` after β converge.
5. **No regressions.** The pre-existing `cnos-agent-admin.yml` continues to fire on its `claude-wake`-titled trigger (verifiable via the post-merge admin-wake regression check filed per cnos#487 AC8).

## Proof / rejection mechanism

**Proof of pass:** all 5 ACs above hold; the smoke is end-to-end-observable in production substrate.

**Proof of fail (any of):**
- The wake doesn't fire within 30 minutes of label application (label-event trigger broken AND scheduled sweep missed) — see cnos#487 OG-4 failure handling.
- The wake fires but doesn't claim (verify gate rejects this well-formed issue) — wake-prompt or claim-sequence regression.
- The wake claims but δ doesn't route (no γ-scaffold appears within the firing's horizon) — Sub 5B contract regression.
- δ routes but a role gets stuck (no convergence; no β verdict) — operational issue; investigate; possibly iterate Sub 5B.
- δ converges but the wake doesn't transition to status:review — wake-prompt regression.

## Cross-references

- **Activated by:** cnos#487 (Sub 5C; the cycle that rendered + activated the live wake)
- **Routing contract:** cnos#486 (Sub 5B; δ wake-invoked mode amendment)
- **Renderer:** cnos#485 (Sub 5A; dispatch-shape renderer extension)
- **Manifest:** cnos#483 (Sub 4; cds-dispatch declaration)
- **Claim protocol:** cnos#454 (dispatch-protocol skill)
- **Wave goal:** cnos#467 (master tracker)

Filed by α@cdd.cnos per cnos#487 Sub 5C scaffold §11 design.
```

### §11.2 — Why this design is the lowest-risk option that still exercises the full path

- **Substantive task is metadata.** The α-implementer's "work" is writing a single doc file recording observations. No production code changes. No test changes. No skill changes. Worst-case, the doc has a typo — recoverable via a follow-up trivial edit.
- **The substantive value is the path being exercised.** Workflow firing → claim sequence → δ routing → γ/α/β execution → artifact emission → PR → status:review. Each of these is non-trivial and is what 5C must prove.
- **The smoke doc is itself empirical evidence.** A reader looking at `docs/gamma/smoke/cds-dispatch-smoke-YYYYMMDD.md` post-merge has the workflow run URL + claim comment + artifact paths + transition timestamps — the same evidence cnos#487's closeouts will cite.
- **Alternative: smoke task could be a docstring change to a non-critical file.** γ considered this; rejected as marginally riskier (touching a doc file outside the smoke directory has higher chance of merge conflict with concurrent work). γ's call: the metadata-doc shape is the lowest-risk option that still exercises the full path.

## §12. Rollback/disable plan

Specific symptoms → specific actions. α's α-prompt names each; this section is the canonical procedure.

### §12.1 — Symptom: double-claim

**Description:** Two firings of `cnos-cds-dispatch.yml` both transition the smoke issue's labels `status:todo → status:in-progress`. Observable via two claim comments from different run IDs on the smoke issue, both within seconds.

**Cause classes:** (a) substrate concurrency group not serializing as expected (per-protocol concurrency layer 1 broken); (b) per-firing wake claimed without re-verify (post-claim re-read layer broken).

**Action (in order):**
1. Comment on the smoke issue naming the double-claim event.
2. Transition smoke issue to `status:blocked` with reason comment.
3. Delete `.github/workflows/cnos-cds-dispatch.yml` from main via an emergency follow-up PR (one-file delete; do NOT touch any other file).
4. In the same emergency PR or a follow-up, flip `cnos.cds/orchestrators/cds-dispatch/wake-provider.json`'s `activation_state` back to `declaration-only` and restore the original `activation_state_notes`.
5. Open a follow-up issue against cnos#454 dispatch-protocol skill OR cnos#485 Sub 5A renderer (whichever surface owns the concurrency broken layer) naming the bug.
6. Keep cnos#487 OPEN with a comment naming the rollback.
7. Do NOT pretend the wave is complete.

### §12.2 — Symptom: loop

**Description:** The wake fires repeatedly on the same issue without progress (multiple claim comments; cell never advances out of `status:in-progress`; no artifacts accumulate).

**Cause classes:** δ wake-invoked mode bug; substrate timeout horizon mis-estimated; claude-code-action restart loop.

**Action:** Same as double-claim (§12.1) — disable the workflow, flip manifest back, file follow-up issue against the suspected surface, keep cnos#487 OPEN.

### §12.3 — Symptom: admin-wake regression

**Description:** `cnos-agent-admin.yml` stops firing OR starts claiming cells (admin/dispatch boundary violation). Observable via: missed schedule firings of admin wake; admin-wake firing transitioning a `protocol:cds` cell's labels.

**Cause classes:** cycle/487's Stage 1 changes inadvertently broke admin wake (OG-2 violation); admin wake's prompt was modified by side-effect.

**Action:**
1. Verify OG-2 audit: `git diff origin/main...HEAD --name-only | grep -E '(cnos-agent-admin|orchestrators/agent-admin)'` — should be empty. If not empty, REVERT the cycle/487 PR (since the violation is in Stage 1, the merged PR is the cause).
2. If OG-2 audit is clean but admin wake still misbehaves, the cause is elsewhere (concurrency interference; environment change). Investigate; do NOT auto-revert.
3. Keep cnos#487 OPEN with the investigation findings.

### §12.4 — Symptom: cnos#487 claimed despite bootstrap-δ lock

**Description:** The live cds-dispatch wake claims cnos#487 itself (transitions `status:in-progress` away; posts a wake-fired claim comment). cnos#487 currently has `status:in-progress` from the bootstrap-δ claim — the dispatch-protocol verify gate (§2.2 step 3) should reject it because the verify requires `status:todo` and `status:in-progress` MUST NOT be present.

**Cause classes:** fundamental claim-protocol bug — the verify gate is broken; OR the dispatch-protocol skill's §2.2 step 3 is not being followed by the wake's prompt; OR labels drifted between bootstrap-δ claim and the firing.

**Action:**
1. STOP immediately.
2. This is a fundamental claim-protocol bug. Do NOT continue smoke testing.
3. Revert cycle/487 PR.
4. File a follow-up issue against cnos#454 (dispatch-protocol skill) AND cnos#483 (cds-dispatch prompt) naming the bug, with the wake firing's run URL + claim comment URL as evidence.
5. Keep cnos#487 OPEN; do NOT pretend the wave is complete.

### §12.5 — Symptom: smoke cell stalls

**Description:** Wake claimed the smoke issue (status:in-progress + claim comment present), but the smoke cell never converges. No β verdict appears within N minutes (estimate: 90 minutes covers a normal cycle/485/486-sized cell; 180 minutes covers worst-case).

**Cause classes:** δ wake-invoked mode misbehavior (Sub 5B contract gap); claude-code-action timeout; runner failure; γ/α/β role-skill bug; cell scope mis-estimated.

**Action:**
1. Investigate: read the wake's run log (`gh run view <run-id> --log`); read the smoke issue's comment trail; read the smoke cell's cycle branch if it exists.
2. If the bug is in δ wake-invoked mode, file a follow-up against cnos#486 (Sub 5B); meanwhile manually move the smoke issue to `status:blocked` with reason.
3. If the bug is in a role-skill, file a follow-up against the relevant role skill (γ/α/β) and surface to operator.
4. Keep cnos#487 OPEN until the smoke completes (with retries) or the failure is classified per OG-5.
5. Do NOT bootstrap-δ-rescue the smoke (that would defeat OG-3).

## §13. Admin-wake regression check procedure (post-merge; for AC8)

After cycle/487 merges + before declaring smoke complete (or as part of Stage 2's evidence-gathering):

1. **File the regression check issue.** Title: `claude-wake AC3 re-verify <date-utc>`. Body: brief; one paragraph naming the purpose. Labels: none required (the admin wake fires on title-match, not labels).
2. **Wait for cnos-agent-admin.yml firing.** Either: the `issues_opened_title_match` trigger fires immediately on the new issue's `opened` event; OR the next scheduled sweep claims it. Max wait: 15 minutes (next scheduled cron slot).
3. **Verify the admin wake fires green.** `gh run view <admin-run-id> --json conclusion --jq .conclusion` returns `success`.
4. **Verify the admin wake does NOT claim any protocol:cds cell.** Inspect the admin wake's firing log for any `gh issue edit --add-label "status:in-progress"` calls on any `protocol:cds` issue. The admin wake's prompt explicitly forbids cell execution (per `cnos-agent-admin.yml` prompt §"Admin-only boundary: MUST NOT execute cells"); the firing log should show admin work only (channel sync, status reporting). If admin wake claims a cell, that's an admin-wake regression — invoke §12.3.
5. **Document evidence in α-closeout (Stage 2).** Cite the admin firing's run URL + the AC3 re-verify issue's number + the result of step 4.
6. **Close the AC3 re-verify issue.** It's served its purpose; close with comment "AC8 verified; admin wake unmodified by cnos#487."

## §14. Friction notes (γ-side; for triage in γ-closeout post-Stage-2)

### FN-1 — Two-stage shape is novel for the wave

cycle/487 is the first cycle in the cnos#467 wave where the smoke is post-merge (Stage 2 happens AFTER the cycle PR ships). Predecessor cycles (470, 476, 485, 486) all had their substantive verification at R0/R[N] pre-merge. The two-stage shape introduces:
- α's α-prompt covers BOTH stages (heavy α-prompt: Stage 1 implementation + Stage 2 observation/evidence).
- β potentially re-engages post-merge to verify Stage 2 ACs (or β verifies only Stage 1 and α/operator handles Stage 2 verification — operator's call).
- Closeouts are authored AFTER Stage 2 (to capture both stages) rather than at converge.

**Triage suggestion:** if the wave produces more "smoke is post-merge" cycles (which seems possible for cdr-dispatch / cdw-dispatch / any future wake-cutover work), γ-template should formalize a "two-stage cycle" mode with sections for Stage 1 ACs, Stage 2 ACs, post-merge re-engagement procedure, and bilateral closeout discipline.

### FN-2 — Bootstrap-δ claim on cnos#487 is a self-application paradox

cnos#487 itself is `status:in-progress` (bootstrap-δ-claimed). The live wake activated by cnos#487's merge should claim issues matching `dispatch:cell + protocol:cds + status:todo`. cnos#487 carries `dispatch:cell + protocol:cds` but `status:in-progress` (not `status:todo`) — so the verify gate (per dispatch-protocol §2.2 step 3) REJECTS it. The structural lock is `status:in-progress`.

**Triage suggestion:** this is fine in v0 (the dispatch-protocol skill has this exact case in its `dispatch_label_drift` semantics — multi-status or wrong-status drift is rejected). For future wake-activation cycles (where the activating cell IS in the queue the activated wake would scan), the bootstrap-δ pattern + the structural lock pattern are the canonical workaround. Worth documenting in the γ-template OR operator-skill as a recognized "activator-cell self-application" anti-pattern.

### FN-3 — Smoke issue is filed by α (not γ); smoke happens post-merge

Per the operator directive "γ does not spawn α/β" + the two-stage shape: γ designs the smoke issue (§11 of this scaffold), α files it post-merge. This is correct per the operator-named role boundaries, but it means γ's "the smoke issue is concretely specified" is a paper artifact until α actually files it.

**Triage suggestion:** if Stage 2 fails because α mis-files the smoke issue (wrong title, wrong labels, wrong body), δ should be able to detect this from the wake's claim attempts (which would either silently skip or post a drift comment). γ's α-prompt §5.1 names the title + labels + body verbatim to minimize this risk. γ recommends bookmarking this scaffold's §11 as the canonical reference.

### FN-4 — Stage 2 observation is operationally complex

Stage 2 requires: filing an issue, polling for workflow firing (potentially 30+ minutes), inspecting workflow logs, inspecting issue comments, inspecting label states, inspecting branch state, inspecting PR state. The α-prompt §5 spells out the procedure; even so, the operational surface is substantial.

**Triage suggestion:** consider whether a future cycle should land a `cn dispatch-smoke <wake-name>` command that automates the observation procedure. Out of scope for cycle/487. Worth filing as a follow-up after Stage 2 completes successfully (sample size 1 — if the manual procedure works once, the automation may be deferrable; if it fails or is unreliable, the automation becomes higher priority).

### FN-5 — δ-side responsibility γ may miss

The two-stage shape splits δ's work too: pre-merge δ routes γ→α→β for Stage 1 (standard bootstrap-δ); post-merge δ would route α for Stage 2 (file smoke issue + observe + capture evidence) + potentially β for Stage 2 re-verification + γ for Stage 2-aware closeouts. The bootstrap-δ session that authored this scaffold IS the δ that will handle Stage 2 routing (or hands off to an operator session). γ flags this for explicit δ-side awareness.

**Triage suggestion:** δ's wake-invoked mode skill (Sub 5B / cnos#486) does NOT cover this case (it covers wake-fired δ for a single cell, not bootstrap-δ for a two-stage cell). γ recommends a follow-up amendment to cnos#486 OR operator-skill formalizing the "two-stage bootstrap-δ" mode after cycle/487 ships (sample size 1 → cycle/487 itself is the empirical case).

### FN-6 — β re-engagement at Stage 2 is structurally ambiguous

β converges at Stage 1 R0 (or R[N]) based on AC1-AC3 + OGs 1, 2, 4, 5 (the Stage-1-verifiable subset). Stage 2 ACs (AC4-AC9) are verified by α at Stage 2 + observed evidence in α-closeout. Does β re-engage post-merge to re-verify AC4-AC9? Two framings:

- **(i) β converges at Stage 1; Stage 2 AC4-AC9 are α-closeout deliverables verified by δ + operator.** Lighter β workload; less duplication; risk: AC4-AC9 verification quality depends on α-closeout discipline.
- **(ii) β re-engages at Stage 2 to re-verify AC4-AC9 independently.** Heavier β workload; structural separation preserved; risk: extra cycle latency.

γ recommends **(i)** — α's Stage 2 evidence-gathering + α-closeout cited evidence + δ's operator-facing summary is sufficient for the smoke's "operationally-confirmed" status. β's Stage 1 converge is the formal verdict for the cycle/487 PR; Stage 2 ACs are evidence-driven, not review-driven. α's α-prompt §5 spells out the evidence capture explicitly. If a future cycle wants stricter Stage 2 review, framing (ii) is the upgrade path.

### FN-7 — Renderer behavior change at activation_state flip is a load-bearing assumption

cycle/487 assumes: post-flip render output is byte-identical to the pre-flip synthetic-live render. The assumption rests on the renderer's `activation_state_effective` value being the only thing that differs between the two render paths, AND the `activation_state_effective` value being `"live"` in both cases. The cycle/485 implementation set `activation_state_effective` from EITHER the override (when supplied) OR the manifest's `activation_state` field (when not). Post-flip: no override, manifest `live` → `activation_state_effective = "live"`. Pre-flip (cycle/485 golden render): `--activation-state-override live`, manifest `declaration-only` → `activation_state_effective = "live"`. Same effective value → same output.

**Verification path:** α's Stage 1 §4.2 invokes the post-flip render + sha256 check. If sha256 differs, fork-D1 fires.

**Triage suggestion:** if fork-D1 fires (sha256 differs), it indicates the renderer has a code path that depends on the manifest's raw `activation_state` field independently of `activation_state_effective`. Worth investigating + potentially patching the renderer to use `activation_state_effective` consistently. Out of scope for cycle/487 (would be a Sub 5A regression cycle); record as a friction note for cycle/487's gamma-closeout if it actually fires.

## §15. Source-of-truth list (citations with file paths + commit SHAs)

All paths absolute on `cycle/487` at branch base `ff47c3600e9d2b804a8f2d76eeff4544e6a13c4d`. Commit SHAs are the merged main HEAD when this scaffold was authored.

| Citation | Path | Notes |
|---|---|---|
| cnos#487 issue body | (GitHub API) | 9 ACs canonical; this scaffold reconciles with two-stage shape |
| cnos#487 bootstrap-δ claim comment | (GitHub API; issue 487 comment 4775875858) | `status:in-progress` lock established 2026-06-23T05:22:05Z |
| cnos#467 master tracker | (GitHub API issue 467) | Wave doctrine; doctrine-correction at top |
| cnos#483 cds-dispatch manifest | `/home/user/cnos/src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.json` | At base; `activation_state: "declaration-only"` |
| cnos#483 cds-dispatch prompt | `/home/user/cnos/src/packages/cnos.cds/orchestrators/cds-dispatch/prompt.md` | Read-only; do NOT modify |
| Sub 5A renderer | `/home/user/cnos/src/packages/cnos.core/commands/install-wake/cn-install-wake` | Merged at PR #488; activation_state refusal at lines ~360-384 |
| Sub 5A golden | `/home/user/cnos/src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml` | sha256 at base: `75e0406662ea12b3caa5ad01d8dbe6bdfb620d0c047fd910917a8e930855605e` |
| Sub 5A admin wake (reference) | `/home/user/cnos/.github/workflows/cnos-agent-admin.yml` | sha256 at base: `fa6b8c0cd64fb626a5e1e991128cbb27fb883b6d1594914543032a2b0d2d3e72` |
| Sub 5A admin wake manifest (reference) | `/home/user/cnos/src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json` | Reference for live manifest shape |
| Sub 5A install-wake-golden CI | `/home/user/cnos/.github/workflows/install-wake-golden.yml` | At base; α modifies per fork-D2 |
| Sub 5B δ wake-invoked mode | `/home/user/cnos/src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` | §9 lines ~395-515; merged at PR #489 |
| cnos#454 dispatch-protocol skill | `/home/user/cnos/src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md` | Claim sequence §2.2; lifecycle §2.4 |
| cnos#479 / PR #481 cutover-A precedent | `/home/user/cnos/.github/workflows/cnos-agent-admin.yml` | Pattern this cycle mirrors |
| cnos#478 mechanical-injection | `.github/workflows/install-wake-golden.yml` (the audit format empirically) | Per-CI-step audit table format mandate |
| cycle/485 closeouts (predecessor) | `/home/user/cnos/.cdd/unreleased/485/{alpha,beta,gamma}-closeout.md` | T-485-3 (cross-package fallback); T-485-T11/T12 (β-skill amendments) |
| cycle/486 closeouts (predecessor) | `/home/user/cnos/.cdd/unreleased/486/{alpha,beta,gamma}-closeout.md` | T-486-1 (variable consistency table); T-486-7 (δ-direct R1 mode); T-486-12 (operator-final-read defense-in-depth) |
| Issue-shape skill | `/home/user/cnos/src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` | Cell-shaped issue contract (smoke issue conforms) |
| γ skill | `/home/user/cnos/src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` | γ role contract (scaffold + closeout) |
| α skill | `/home/user/cnos/src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` | α role contract |
| β skill | `/home/user/cnos/src/packages/cnos.cdd/skills/cdd/beta/SKILL.md` | β role contract |

---

— γ@cdd.cnos (bootstrap-δ via current session), 2026-06-23 (UTC)
