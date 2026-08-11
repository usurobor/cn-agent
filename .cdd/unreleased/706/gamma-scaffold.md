# γ scaffold — cnos#706

## Cell header

| Field | Value |
|---|---|
| Issue | [usurobor/cnos#706](https://github.com/usurobor/cnos/issues/706) |
| Protocol | `cds` |
| Base SHA (`main`) | `7f249ddbb50f230d5d41287b6554ab17b5a1d1d5` |
| Wake run id | `cds-dispatch-manual-706` (manual dispatch-wake invocation for this session; not a real GH Actions run id) |
| Branch | `cycle/706` |
| `run_class` | `first_pass` — no prior `cycle/706` branch, no prior PR, no `.cdd/unreleased/706/` artifacts beyond the claim marker (`CLAIM-REQUEST.yml`) existed at claim time; no rejection evidence (no `status:changes` history, no `operator-review.md`, no bounce comments); verified by δ before dispatch (issue comment 6: "no prior branch/PR/artifacts, no rejection evidence"). Confirmed independently by γ: `git ls-remote --heads origin cycle/706` returned nothing before this branch was created. |
| Mode | `design-and-build` (see §Mode rationale) |

### Mode rationale

The consolidated final spec (issue comment "CONSOLIDATED FINAL SPEC … build from this") is a fully converged design — every open design question (secret name, bot identity, GITHUB_TOKEN-sufficiency) was investigated and resolved across the refinement thread. But per `cdd/issue/SKILL.md` §"MCA preconditions", MCA requires the design to live at a **stable docs path** (`docs/{tier}/{bundle}/{X.Y.Z}/DESIGN.md`), not inside an issue-comment thread. No such path exists for #706. This does not block the cycle — the design is converged, just not filed at an MCA-qualifying path — so mode is `design-and-build`, not MCA. α should treat the consolidated-spec comment as authoritative design input without needing to re-derive it.

### Cycle scope sizing (five-factor check, `issue/SKILL.md` §"Five-factor split-decision heuristic")

10 ACs sits in the "at-edge" (8–10) band, which requires this check + written justification if kept whole.

| Factor | Reading | Splitting signal? |
|---|---|---|
| (a) New code surface | ~1 new package (secrets/push-access preflight, mirroring `label-doctor`'s shape) + edits to existing renderer/CLI | No (< 2 new modules) |
| (b) Cross-module breadth | `repoinstall.go`, `cn-install-wake` (bash renderer), the live rendered workflow, 2 golden fixtures, `docs/guides/INSTALL-CDS.md`, 2 test files | Yes (≥3 modules touched) |
| (c) Lifecycle span | Design already converged (issue thread); this cycle is build + doc update, not multi-phase design→code→docs serialization | No |
| (d) MCA-precondition stability | Design stable (consolidated, operator-ruled); no further design churn expected | No |
| (e) Independent shippability | The rename+bot-deletion (AC6/AC9) and the preflight-ordering (AC1–4) and the docs (AC5/AC7/AC8/AC10) are each independently shippable | Yes |

Two factors fire ("yes" on ≥2 is a strong split indicator"). **Decision: keep whole**, not split into master+subs. Justification: (1) the operator already dispatched this as one consolidated cell after the full refinement thread (issue comment: "Dispatching now" against the 10-AC consolidated list) and δ has already claimed it as a single `first_pass` cell — splitting now would require an issue restructure mid-claim, disproportionate to the gain; (2) all 10 ACs serve one governing gap (a first-time user cannot follow the install) and are tightly coupled around the same rendered artifact (`.github/workflows/cnos-cds-dispatch.yml`) and the same command (`cn repo install --dispatch`) — splitting would fragment one coherent user journey across multiple PRs touching the same file, increasing merge-conflict risk rather than reducing it. γ records this transparently rather than silently overriding δ's already-executed claim.

---

## Governing gap (restated from consolidated spec)

The CDS install was built by/for someone who already had the whole setup (a bot, tokens, mental model). It surfaces operator-only gates *last* instead of *first*, invents a non-existent "bot" concept, and names secrets opaquely. Fix: ask the operator for what only they can provide, explain exactly how to get it, before doing anything — and stop inventing a bot.

**Authoritative source:** the issue's LAST comment, "CONSOLIDATED FINAL SPEC (for dispatch, 2026-08-05) — build from this" — explicitly: *"Where it differs from the issue body, this wins."* Build from its Deliverables 1–7 and Final ACs 1–10. Do NOT build from the stale original issue body (superseded design: it never mentions deleting the cosmetic bot or renaming to `CN_DISPATCH_PAT`).

---

## Source-of-truth table

| Claim / surface | Canonical source | Status |
|---|---|---|
| `cn repo install --dispatch` domain logic (`Run`, `runDispatchCds`, `ensureCanonicalDispatchLabels`, `applyInstall`) | `src/go/internal/repoinstall/repoinstall.go` | Shipped (cnos#608/#610); this cycle amends |
| CLI wrapper / flag parsing / help text | `src/go/internal/cli/cmd_repo_install.go` (`repoInstallHelp` const, `RepoInstallCmd.Run`) | Shipped; this cycle amends `repoInstallHelp` for AC5 |
| Existing domain tests (pattern precedent — extend, don't discard) | `src/go/internal/repoinstall/repoinstall_test.go` | Shipped; contains a test whose assertion AC2 will invert (see §Friction 1) |
| Existing CLI tests | `src/go/internal/cli/cmd_repo_install_test.go` | Shipped |
| Dispatch-workflow renderer (bash) | `src/packages/cnos.core/commands/install-wake/cn-install-wake` — `agent_bot_name()`/`agent_bot_id()` (lines 194–206), sigma-default PAT-secret binding (lines 741–756), `github_token`/`bot_name`/`bot_id` YAML emission (lines 1274–1276) | Shipped (cnos#609/#613); this cycle amends |
| Live rendered workflow (migration target, AC6/AC9) | `.github/workflows/cnos-cds-dispatch.yml` — currently: `token`/`GH_TOKEN`/`github_token` all bound to `secrets.SIGMA_WORKFLOW_PAT` (lines 40, 65, 88, 424); `bot_name: "sigma@cnos.cn-sigma.cnos"` / `bot_id: "41898282"` (lines 89–90) | Shipped; must be regenerated via the renderer after the renderer changes, not hand-edited |
| Golden fixtures (CI-diffed, regenerate via renderer) | `src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml`, `src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml` — both reference `SIGMA_WORKFLOW_PAT` today | Shipped; regenerate, don't hand-edit |
| Golden-diff CI gate (re-renders + byte-diffs the above) | `.github/workflows/install-wake-golden.yml` | Shipped; will fail if goldens are hand-edited out of sync with the renderer |
| GitHub REST primitive precedent (dependency-free `net/http`, name-only wire types, no `gh` CLI shellout, no third-party client) | `src/packages/cnos.core/commands/label-doctor/github.go` (`ghRequest`, `ghListLabels` — the closest existing "list via GET, presence check" shape) | Shipped (cnos#493); model, don't import directly (separate module) |
| `owner/repo` resolution from git remote (precedent for resolving the *installing* repo's target, distinct from `repoinstall.Options.Repo` which names the cnos *release source*) | `src/packages/cnos.core/commands/label-doctor/resolve.go` (`resolveRepoFromGitRemote`) | Shipped; model |
| In-process cross-module package call precedent (exactly the shape a new preflight package should follow) | `src/go/internal/repoinstall/repoinstall.go`'s `ensureCanonicalDispatchLabels()` → `labeldoctor.Doctor(ctx, labeldoctor.Options{...})`, plus `go.work`'s `use (...)` block registering `label-doctor` as a linked module | Shipped; model |
| Package module registration (a new preflight package needs an entry here) | `go.work` (repo root) | Shipped; amend if a new module is added |
| Existing quickstart doc (revise in place — it already exists and is comprehensive but reflects the pre-#706 design) | `docs/guides/INSTALL-CDS.md` | Shipped; AC5/AC7/AC8/AC9/AC10 require substantial revision, not a fresh doc |
| Design-mocks doc (naming convention: "Mock A/B/C/…") | `docs/development/design/cn-repo-install-MOCKS.md` | Shipped; reference only, not required reading for implementation |

---

## Per-AC oracle list (Final ACs 1–10, consolidated spec)

**AC1 — preflight runs before any label/render/commit.**
Oracle: read `repoinstall.Run()` in `src/go/internal/repoinstall/repoinstall.go`. Today (lines ~316–334) `applyInstall(...)` (writes `.cn/deps.json`, restores packages) runs unconditionally *before* the `if opts.Dispatch == "cds"` branch that calls `runDispatchCds` (labels/render). After the fix, for `opts.Dispatch == "cds"`, the new preflight check must be the first thing that happens — strictly before `applyInstall`. Mechanical test oracle: a new test (e.g. `TestRun_DispatchCds_PreflightRunsBeforeBaseInstall`) that stubs the secrets/permission API to report "missing," calls `Run` with `Dispatch: "cds"`, and asserts `.cn/deps.json` does **not** exist under `opts.RepoRoot` afterward (proving base install never ran) and that no label-doctor / render call was attempted.

**AC2 — missing prerequisites → non-zero exit + actionable message; no partial deploy artifacts.**
Oracle: same fixture as AC1. Assert `Run` returns a non-nil `error`; assert the error/stderr text names exactly which secret(s)/permission are missing; assert neither `.cn/` nor `.github/workflows/cnos-cds-dispatch.yml` exist under `opts.RepoRoot` afterward. This directly supersedes the current assumption in `TestRun_DispatchCds_RendererNotVendored_FailsWithNoPartialWrite` (`repoinstall_test.go:579`), which today asserts `.cn/deps.json` **does** still exist after a dispatch-path failure — see §Friction 1.

**AC3 — presence verified without the CLI ever receiving secret values.**
Oracle: `GET /repos/{owner}/{repo}/actions/secrets` is GitHub's documented presence-only endpoint — it returns `name`/`created_at`/`updated_at`, never a value. Grep the new preflight source for any field/variable that could hold a secret value (`grep -rn "value\|Value" <new-package>/*.go` should show no secret-value handling — only `name` fields, mirroring `ghLabel`'s name-only wire struct in `label-doctor/github.go`). Test oracle: an httptest fixture returning secret objects `{name, created_at, updated_at}` (no `value` key at all, matching GitHub's real response shape) decodes and drives the preflight correctly — proving the code path never expects or reads a value field.

**AC4 — prerequisites present → install proceeds; re-run resumes cleanly.**
Oracle: happy-path test with an httptest fixture reporting both `CLAUDE_CODE_OAUTH_TOKEN` and `CN_DISPATCH_PAT` present and push-access true; assert `Run()` completes with no error (mirrors the existing happy-path precedent `TestRun_DispatchCds_SigmaDefault_NoIdentityFlagsRequired`, `repoinstall_test.go:1048`). Idempotency oracle: call `Run()` twice against the same `opts.RepoRoot` with prerequisites present both times; assert the second run's file set is byte-identical to the first (mirrors the existing idempotency precedent `TestRun_Idempotent_ByteIdenticalArtifacts`, `repoinstall_test.go:412`).

**AC5 — one-page quickstart exists, linked from the install command's help/output.**
Oracle: `docs/guides/INSTALL-CDS.md` exists (confirmed) and, after revision, contains no stale references (see AC8/AC9 oracles below for the specific greps). Link-wiring oracle: `grep -n "INSTALL-CDS" src/go/internal/cli/cmd_repo_install.go` — today this returns **zero matches** (confirmed by reading `repoInstallHelp`); α must add a reference (path or full GitHub URL) to `repoInstallHelp` and/or the preflight failure message. Test oracle: extend `cmd_repo_install_test.go` to assert `--help` output contains the quickstart reference string.

**AC6 — workflow-PAT secret named `CN_DISPATCH_PAT` in the rendered template.**
Oracle: after the fix, `grep -rn "SIGMA_WORKFLOW_PAT" src/packages/cnos.core/commands/install-wake/cn-install-wake src/go/internal/repoinstall/repoinstall.go src/go/internal/cli/cmd_repo_install.go` returns zero matches; `grep -n "CN_DISPATCH_PAT" src/packages/cnos.core/commands/install-wake/cn-install-wake` shows the new default binding replacing the sigma-only default at line 754. Per the migration note (issue comment, "CONSOLIDATED FINAL SPEC" → "Migration note"), the live rendered artifact `.github/workflows/cnos-cds-dispatch.yml` (today referencing `SIGMA_WORKFLOW_PAT` at lines 40/65/88/424) and both golden fixtures (`cnos-cds-dispatch.golden.yml`, `cnos-agent-admin.golden.yml`) must be **regenerated via the renderer**, not hand-edited — the `install-wake-golden.yml` CI job re-renders and byte-diffs them.

**AC7 — preflight message includes, per prerequisite: what/why/exact acquisition steps/Settings link.**
Oracle: unit test asserting the preflight failure text for a missing `CLAUDE_CODE_OAUTH_TOKEN` contains `"claude setup-token"` and a Settings-path reference; for a missing `CN_DISPATCH_PAT` contains `"fine-grained"`, the four scopes (`Contents`, `Issues`, `Pull requests`, `Workflows`), and `"Settings"` → `"Developer settings"` → `"Personal access tokens"` → `"Fine-grained tokens"`. Grounded verbatim in the operator's final wording (issue comment "Decisions (operator, 2026-08-05)"): *"a fine-grained Personal Access Token on your own GitHub account, scoped to this one repo with Contents + Issues + Pull requests + Workflows = write... Create it at Settings → Developer settings → Personal access tokens → Fine-grained tokens."*

**AC8 — quickstart defines every term before use; states plainly there is no bot account today; no step assumes a pre-existing account.**
Oracle: `grep -in "bot" docs/guides/INSTALL-CDS.md` — every "bot" mention must follow (in doc reading order) an explicit definition. Negative oracle: `grep -n "sigma@cnos.cn-sigma.cnos\|41898282" docs/guides/INSTALL-CDS.md` must return zero matches after the fix (today it references both as if a real bindable identity — Tier 3 runbook section, "for the default sigma agent the renderer expects SIGMA_WORKFLOW_PAT"). The doc's Tier 3 runbook (currently instructing the operator to supply `--bot-name`/`--bot-id`) needs rewriting per the operator's final ruling (comment "Decisions … delete the bot").

**AC9 — rendered workflow contains no cosmetic `bot_name`/`bot_id` by default; commits authored by the token's account.**
Oracle: after regenerating, `grep -n "bot_name\|bot_id" .github/workflows/cnos-cds-dispatch.yml` returns zero matches (today: lines 89–90). This requires removing (or making non-default/opt-in) the `agent_bot_name()`/`agent_bot_id()` lookup + emission in `cn-install-wake` (lines 194–206, 1275–1276) — the consolidated spec's Deliverable 3 reads as outright deletion for the default path ("the renderer/template injects **no** `bot_name`/`bot_id` by default … never a first-run requirement"), not merely leaving the sigma table present-but-unused. Test oracle: extend the existing prose-leak-grep precedent `TestDispatchRenderer_ProseLeakGrep_CatchesPreFixSigmaPhrasing` (`repoinstall_test.go:1128`) — or add a sibling test — asserting a fresh render never emits `bot_name:`/`bot_id:` keys absent an explicit `--bot-name`/`--bot-id` flag.

**AC10 — bot-less path is the documented default; dedicated bot is the future upgrade (cnos#449/#702); any `GITHUB_TOKEN`-only minimal claim is verified before being documented.**
Oracle: the operator's own investigation (issue comment "Decisions … The token is a PAT, not GITHUB_TOKEN — verified rationale") already discharges the "verify before documenting" requirement — it concludes a PAT **is** required for the agent/dispatch tier (install-time workflow-write + GitHub's loop-prevention blocking `GITHUB_TOKEN`-triggered downstream workflows). α's job is to carry that already-verified conclusion into the doc, not re-derive it, and must not assert an unverified "`GITHUB_TOKEN` suffices for dispatch" claim anywhere. Oracle: `grep -n "GITHUB_TOKEN" docs/guides/INSTALL-CDS.md` — every `GITHUB_TOKEN`-only claim in the revised doc must stay scoped to the Tier 2 `--engine` path (which already correctly documents `GITHUB_TOKEN`-only operation, per the existing "Tier 2 runbook" section) and must never claim it suffices for Tier 3 (agent/dispatch). The "Tenant secrets, by tier" table's Tier 3 row + the Tier 3 runbook must present the two-own-account-secrets flow (`CLAUDE_CODE_OAUTH_TOKEN` + `CN_DISPATCH_PAT`) as the default, with dedicated-bot explicitly called out as future/deferred (cnos#449 / #702), not a required step.

---

## α prompt

**Branch:** `cycle/706`

Implement the consolidated final spec's Deliverables 1–7 against the 10 ACs above. Work directly in `src/go/internal/repoinstall/`, `src/go/internal/cli/cmd_repo_install.go`, `src/packages/cnos.core/commands/install-wake/cn-install-wake`, the rendered/golden YAML artifacts, and `docs/guides/INSTALL-CDS.md`.

1. **Preflight-first (AC1/AC2/AC4).** Add a preflight check that runs, for `--dispatch cds` only, before `applyInstall` and before `runDispatchCds`'s render/label steps. On missing prerequisites: return a named, actionable, non-zero error; write **no** files (no `.cn/`, no `.github/workflows/…`). On present prerequisites: proceed exactly as today. Make re-running after prerequisites are satisfied resume cleanly (no special "resume" state needed — the existing idempotent-install property already gives you this, as long as the preflight is a pure read-check with no side effects on failure).

2. **Presence-only checks, never secret values (AC3/AC7).** Two checks:
   - **Secrets exist by name**: `GET /repos/{owner}/{repo}/actions/secrets`, check `CLAUDE_CODE_OAUTH_TOKEN` and `CN_DISPATCH_PAT` (see #4 below for the rename) are present by name. This endpoint never returns values — decode only `name`/`created_at`/`updated_at`.
   - **Push access**: verify the token being used to install has push access to the target repo. `label-doctor/github.go`'s `ghRequest` + `resolveRepoFromGitRemote` (`label-doctor/resolve.go`) are the closest precedent for the dependency-free `net/http` idiom and target-repo resolution — mirror that shape rather than introducing a `gh` CLI shellout or a third-party GitHub client. See §Friction 3 for a concrete recommendation on which endpoint to use for the push-access check.
   - Recommended package shape: a new sibling Go module under `src/packages/cnos.core/commands/` (own `go.mod`, added to `go.work`'s `use (...)` block), called in-process from `repoinstall.go` exactly as `ensureCanonicalDispatchLabels()` calls `labeldoctor.Doctor(...)` today. This is a strong precedent, not a rigid mandate — if a same-module `internal/repoinstall` file turns out simpler given how tightly this couples to `Run()`, that is a reasonable α call; just don't invent a third pattern.

3. **Preflight message content (AC7).** For each missing prerequisite, state what it is, why it's needed, and the exact acquisition steps, using the operator's own final wording (quoted verbatim in the AC7 oracle above) — don't paraphrase away the specifics (scope names, the exact Settings path, `claude setup-token`).

4. **Rename the secret everywhere (AC6).** `SIGMA_WORKFLOW_PAT` → `CN_DISPATCH_PAT`, across: `cn-install-wake`'s default binding, `repoinstall.go`'s display value, `cmd_repo_install.go`'s help text, the live rendered `.github/workflows/cnos-cds-dispatch.yml`, and both golden fixtures. Regenerate the rendered/golden YAML via the renderer (`cn install-wake cds-dispatch` / `cn install-wake agent-admin`) — do not hand-edit generated YAML; the `install-wake-golden.yml` CI job will byte-diff it.

5. **Delete the cosmetic bot identity by default (AC9).** Remove (or make strictly opt-in via explicit `--bot-name`/`--bot-id`, never defaulted) the `agent_bot_name()`/`agent_bot_id()` sigma-only lookup table and its YAML emission in `cn-install-wake`. A fresh render must carry no `bot_name`/`bot_id` keys. Regenerate the rendered/golden artifacts.

6. **Write the one-page quickstart (AC5/AC8/AC10).** `docs/guides/INSTALL-CDS.md` already exists and covers most of the ground — **revise it in place**, don't create a second doc. Update: the Tier 3 runbook (drop the "create a bot" framing entirely per AC8/AC9), the secret name everywhere, the "Tenant secrets, by tier" table, and add explicit term definitions (PAT, repo secret, default branch, "bot" — stating plainly there is no bot account to create today) before first use of each term.

7. **Wire the quickstart into the install command's help/output (AC5).** Add a reference (path or GitHub URL) to `docs/guides/INSTALL-CDS.md` in `repoInstallHelp` (`cmd_repo_install.go`) and/or the preflight failure message.

The 7-axis implementation-contract discipline (`cnos.handoff/skills/handoff/dispatch/SKILL.md`) is not required for this cycle — no axis is undecidable in a way that blocks starting. One soft design choice is flagged in §Friction 3 (which GitHub endpoint to use for the push-access check) — resolve it yourself and record the choice + rationale in `self-coherence.md`; it is not a hard blocker requiring escalation. If you hit a genuinely undecidable axis while implementing, STOP and escalate to γ/operator rather than improvising past it — but none is expected here.

Write `.cdd/unreleased/706/self-coherence.md` per `cdd/alpha/SKILL.md` — canonical section headers (`## Gap`, `## Skills`, `## ACs`, `## CDD Trace`, `## Self-check`, `## Debt`), not decorated `## §X` forms (see `gamma/SKILL.md` §2.5's binding reminder on this — a prior cycle hard-FAILed `cn cdd verify` over exactly this).

---

## β prompt

**Branch:** `cycle/706`

Review against the 10 ACs above independently — walk each AC's oracle yourself; don't rely on α's self-report.

1. **AC1/AC2 — preflight-before-mutation.** Run (or trace through) the missing-prerequisites path yourself. Confirm zero filesystem writes occur under the install target when prerequisites are absent — not just "no `.github/workflows/`" (the pre-existing bar) but also **no `.cn/deps.json`** (the new, stricter bar this cycle sets). If α left the old `TestRun_DispatchCds_RendererNotVendored_FailsWithNoPartialWrite` assertion (`.cn/deps.json` still exists) unchanged, that is a correctness bug against AC2, not acceptable legacy.
2. **AC3 — no secret values ever handled.** Grep the new/changed preflight code yourself for any value-shaped field, log line, or error message that could carry a secret value. Confirm the wire type decoding the secrets-list API response has no `value`/`Value` field at all (mirror `ghLabel`'s name-only shape).
3. **AC6/AC9 — migration completeness.** Confirm `SIGMA_WORKFLOW_PAT` and `bot_name`/`bot_id` are gone from: `cn-install-wake`, `repoinstall.go`, `cmd_repo_install.go`, the live `.github/workflows/cnos-cds-dispatch.yml`, and **both** golden fixtures (`cnos-cds-dispatch.golden.yml`, `cnos-agent-admin.golden.yml`) — not just the live workflow. Verify the golden fixtures were regenerated via the renderer (re-run it yourself and diff), not hand-edited to merely look right.
4. **AC5/AC7/AC8/AC10 — docs.** Confirm `docs/guides/INSTALL-CDS.md` was revised in place (not duplicated), every "bot" reference is preceded by a definition, no hardcoded `sigma@cnos.cn-sigma.cnos`/`41898282` remains, the preflight-message wording in code matches (or is a faithful paraphrase of) the doc, and no `GITHUB_TOKEN`-suffices claim leaks outside the Tier 2 (`--engine`) scope.
5. **AC4 — happy path + idempotency.** Confirm a fresh test exercises the present-prerequisites path end-to-end, and that a second `Run()` call with prerequisites still present produces no further diff.
6. **Tests exist and pass.** Run the full `src/go` test suite (and the new preflight package's own tests, if it's a separate module) — don't just check that new tests were *added*; run them.
7. **Scope discipline.** Confirm nothing in §Scope guardrails below leaked into the diff (no secret-value handling, no runtime-permissions-model change, no dedicated-bot creation flow, no `--engine` behavior change beyond incidental cleanup).

---

## Scope guardrails (explicitly OUT of scope this cycle)

- **No secret value collection, storage, or logging** — presence-by-name only, ever.
- **No change to the workflow's runtime security model** — the default-branch merge gate and the existing secrets gates stay exactly as they are; dispatch remains PR-only.
- **No one-command full automation of the three operator gates** (secrets, push access, merge-to-default-branch) — they stay operator-only by design. The fix is *ordering* (ask first) and *explanation* (say how), not elimination.
- **No dedicated-bot-account creation flow.** Bot-less is the default. A real dedicated bot account is future/opt-in work tracked at cnos#449 / #702 — not built, scaffolded, or stubbed in this cycle.
- **No `--engine` (PAT-free mechanical FSM tier) behavior change** beyond whatever incidental secret-name/bot-cosmetic cleanup naturally falls out of the rename+delete. The engine tier already needs neither secret; don't add new engine-tier scope.
- **No new GitHub App / installer surface.** `docs/development/design/cnos-installer-github-app.md` is explicitly future work (referenced in `INSTALL-CDS.md` as "coming") — not this cycle.
- **No operator-facing re-verification against `usurobor/cmp`.** The consolidated spec's "Test target" section names a post-merge operator action, not a cycle deliverable.

---

## Friction notes

1. **Existing test assertion will flip.** `TestRun_DispatchCds_RendererNotVendored_FailsWithNoPartialWrite` (`repoinstall_test.go:579`, in the current codebase before this cycle's changes) explicitly asserts `.cn/deps.json` **still exists** after a dispatch-path failure ("Base install artifacts, by contrast, are unaffected"). AC2's stricter "no partial deploy artifacts" bar means this assumption is now wrong for the *missing-prerequisites* failure mode specifically (it may still hold for *other* failure modes further down the pipeline, like the renderer-not-vendored case that test covers, since preflight would have already passed by then) — α needs to read this test carefully and adjust deliberately, not leave it silently contradicting AC2, and not break it by accident either. Both failure modes (missing prerequisites vs. renderer-not-vendored) are real and distinct; make sure the new preflight-failure test is additive, and that the renderer-not-vendored test's own assertion still correctly reflects that preflight already passed by the time that specific failure fires.

2. **No existing precedent for GitHub secrets-presence or push-access checks.** Confirmed via `grep -rn "actions/secrets\|collaborators.*permission"` across `src/go` and `src/packages` — zero hits. This is genuinely new surface. `label-doctor/github.go`'s `ghRequest`/`ghListLabels` (dependency-free `net/http`, no `gh` CLI, no third-party client) is the strongest existing shape to mirror; flagged so α doesn't reach for a different HTTP client library or shell out to `gh`.

3. **Push-access check: recommend the simpler endpoint.** The issue's Deliverable 2 says "verify push access via the permission API," and the original issue body named the collaborator/permission endpoint (`GET /repos/{owner}/{repo}/collaborators/{username}/permission`) — but that requires first resolving *which* username the PAT belongs to (an extra `GET /user` call). Since AC9 deletes the separate-bot-identity concept, the check is really just "does the token used for install have push access to this repo," which `GET /repos/{owner}/{repo}` already answers directly via its `permissions.push` field for the authenticated caller, in one call. Recommending the single-call form; not escalating as a blocking undecidable axis — α should pick one, document the choice + one-line rationale in `self-coherence.md`, and move on.

4. **Two distinct "sigma" bindings, both need to go, don't fix one and miss the other.** `cn-install-wake` has (a) the sigma-only `agent_bot_name()`/`agent_bot_id()` cosmetic-identity table (AC9's target) and (b) a *separate* sigma-only default for the PAT secret name itself (`workflow_pat_secret="SIGMA_WORKFLOW_PAT"` when `--agent sigma`, AC6's target). These are two different call sites in the same file (~lines 194–206 vs. ~741–756) that could each be fixed independently while missing the other — both must land.

5. **Mode is `design-and-build`, not MCA** — no `DESIGN.md`/`PLAN.md` exists at a stable docs path for #706; the converged design lives in the issue-comment thread instead. Not a blocker, just the correct mode label per `issue/SKILL.md`'s MCA preconditions (design must be at a stable path, not in-issue).

6. **A pre-existing `CLAIM-REQUEST.yml` was already staged in the working tree** at `.cdd/unreleased/706/CLAIM-REQUEST.yml` when γ began this session (git status showed it as `A` — staged, uncommitted, on `main`, before `cycle/706` was cut). γ carried it into `cycle/706` via `git switch -c` but, per the explicit dispatch instruction to commit "only" this scaffold file, left it uncommitted and did not add it to the scaffold commit. It documents the FSM claim transition that issue comment 6 confirms already happened server-side (`status:todo → status:in-progress`). δ/α should confirm whether it needs a separate commit onto `cycle/706` (it is not overwritten by this scaffold's commit either way).
