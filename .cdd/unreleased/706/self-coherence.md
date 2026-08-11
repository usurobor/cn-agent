## Gap

Issue: [usurobor/cnos#706](https://github.com/usurobor/cnos/issues/706) — "CDS install preflight-first: ask for operator prerequisites before doing anything, stop inventing a bot."

Version/mode: `design-and-build` (per γ scaffold — the design is converged in the issue's "CONSOLIDATED FINAL SPEC" comment, but not filed at a stable `docs/{tier}/{bundle}/{X.Y.Z}/DESIGN.md` path, so MCA preconditions are not met; this is not a blocker, just the correct mode label).

Governing gap (restated from the consolidated spec, authoritative over the original issue body per its own "where it differs, this wins" clause): `cn repo install --dispatch cds` did agent-doable work first (render/labels/commit) and surfaced operator-only gates (secrets, push access, merge-to-default-branch) last, invented a non-existent "bot" concept, and named the workflow-PAT secret opaquely (`SIGMA_WORKFLOW_PAT`, hardcoding the agent name). Fix: ask the operator for what only they can provide, explain exactly how to get it, before doing anything — and stop inventing a bot.

This cycle (R0) implements the consolidated spec's Deliverables 1–7 against Final ACs 1–10, working from `cycle/706` (cut from `main` HEAD `7f249ddbb50f230d5d41287b6554ab17b5a1d1d5`) on top of γ's scaffold (`.cdd/unreleased/706/gamma-scaffold.md`, commit `3bf1b2d`).

## Skills

Active skills this cycle (Tier 1/2/3, per `cnos.cdd/skills/cdd/alpha/SKILL.md`):

- `cdd/alpha/SKILL.md` §2.1 (dispatch intake — branch checkout, reading the scaffold as the authoritative contract), §2.5 (this file's incremental, canonical-header-form write discipline), §2.6 (pre-review gate, applied before signaling review-readiness).
- No `cnos.handoff/skills/handoff/dispatch/SKILL.md` 7-axis contract discipline — γ's scaffold explicitly waives it for this cycle ("no axis is undecidable in a way that blocks starting").
- `eng/go` conventions (implicit, not a loaded skill file this session, but followed): cli/-boundary compliance (`cli/cmd_repo_install.go` stays a thin wrapper; all domain logic in `internal/repoinstall`), dependency-free `net/http` GitHub REST idiom (mirrored from `label-doctor/github.go`, not imported — separate go.work module).

## ACs

Final ACs 1–10 per the consolidated spec (`.cdd/unreleased/706/gamma-scaffold.md` §"Per-AC oracle list"). Evidence below cites tests by name (all runnable via `cd src/go && go test ./internal/repoinstall/... ./internal/cli/...` and `go test ./src/packages/cnos.core/commands/install-preflight/...` from repo root) and greps I ran directly against branch HEAD.

**AC1 — preflight runs before any label/render/commit.**
PASS. `Run()` (`src/go/internal/repoinstall/repoinstall.go`) calls `runPreflight(ctx, opts)` immediately after the "✓ Git repository root" stdout line, strictly BEFORE `resolveIndex` (no network call for the release index), `applyInstall` (no `.cn/` write), and `runDispatchCds` (no render, no label-doctor call) — gated on `opts.Dispatch == "cds" && !opts.DryRun`. Evidence: `TestRun_DispatchCds_PreflightRunsBeforeAnythingElse_NoPartialArtifacts` (`repoinstall_test.go`) passes a deliberately unreadable `IndexPath` and asserts the returned error is the preflight error, not an index-read error naming that path — proving preflight fires before `resolveIndex` ever touches it.

**AC2 — missing prerequisites → non-zero exit + actionable message; no partial deploy artifacts.**
PASS. Same test as AC1 additionally asserts: `err != nil`; the error names `CLAUDE_CODE_OAUTH_TOKEN`, `CN_DISPATCH_PAT`, and `push access`; `.cn/` does not exist under `opts.RepoRoot`; `.github` does not exist; stdout carries no `"rendered"` / `"label-doctor"` line. `TestRun_DispatchCds_RendererNotVendored_FailsWithNoPartialWrite`'s pre-existing assertion (`.cn/deps.json` DOES exist) is deliberately left unchanged, with a new comment explaining why it is still correct: that test's fixture now satisfies preflight first (`setPreflightSatisfiedEnv`), so its failure point is genuinely downstream of preflight (the renderer-not-vendored check), where base install has already legitimately run — a distinct, later failure mode from AC2's missing-prerequisites mode (γ's Friction note 1).

**AC3 — presence verified without the CLI ever receiving secret values.**
PASS. `installpreflight.ghSecret` (`src/packages/cnos.core/commands/install-preflight/github.go`) has exactly 3 fields (`Name`, `CreatedAt`, `UpdatedAt`) — no `Value`/`value` field anywhere in the package. `TestGhSecret_HasNoValueField` asserts this structurally via reflection (fails if a future edit adds a value-shaped field). `TestGhListSecrets_PresenceOnlyResponse_DecodesCorrectly` decodes GitHub's exact real response shape (`{name, created_at, updated_at}`, no `value` key) and confirms presence-detection still works. `grep -rn "value" src/packages/cnos.core/commands/install-preflight/*.go` (excluding `_test.go`) returns zero secret-value-shaped hits — confirmed by running the grep directly.

**AC4 — prerequisites present → install proceeds; re-run resumes cleanly.**
PASS. `TestRun_DispatchCds_SigmaDefault_NoIdentityFlagsRequired` and `TestRun_DispatchCds_RendersWorkflow_ThenSurfacesLabelGap` (both with `setPreflightSatisfiedEnv`) prove preflight passing lets execution proceed into the render. `TestRun_DispatchCds_PreflightSatisfied_SecondRunByteIdentical` proves the re-run-resumes-cleanly half directly: calling `Run` twice against the same `repoRoot` with prerequisites satisfied produces byte-identical `.cn/deps.json` and rendered-workflow output both times. A fully `err == nil` end-to-end `Run()` call was not achieved (see §Debt #1 — label-doctor's own hardcoded API base is not reachable from this package's test seam); the idempotency + progression evidence above stands in for it.

**AC5 — one-page quickstart exists, linked from the install command's help/output.**
PASS. `docs/guides/INSTALL-CDS.md` revised in place (not recreated — `git log --follow` shows continuous history). `cmd_repo_install.go`'s `repoInstallHelp` const now contains a `SEE ALSO: docs/guides/INSTALL-CDS.md` line plus two inline mentions (`grep -n "INSTALL-CDS" src/go/internal/cli/cmd_repo_install.go` → 3 matches, was 0 before this cycle). `TestRepoInstall_HelpFlag` extended to assert `--help` output contains `"INSTALL-CDS.md"`. `formatPreflightFailure` (the preflight failure message itself) also references it.

**AC6 — workflow-PAT secret named `CN_DISPATCH_PAT` in the rendered template.**
PASS. `grep -rn "SIGMA_WORKFLOW_PAT" src/packages/cnos.core/commands/install-wake/cn-install-wake src/go/internal/repoinstall/repoinstall.go src/go/internal/cli/cmd_repo_install.go` → zero matches (confirmed directly; a follow-up commit `c272e53` scrubbed two rename-history comments that still carried the literal old string). `grep -n "CN_DISPATCH_PAT" src/packages/cnos.core/commands/install-wake/cn-install-wake` shows the new default binding. Live `.github/workflows/cnos-cds-dispatch.yml` and both golden fixtures (`cnos-cds-dispatch.golden.yml`, `cnos-agent-admin.golden.yml`) regenerated via the renderer; re-running the renderer a second time reports "unchanged" for all three (idempotent), and the live workflow's sha256 matches the golden's sha256 (`eb8c6294...` — verified by direct `sha256sum` comparison).

**AC7 — preflight message includes, per prerequisite: what/why/exact acquisition steps/Settings link.**
PASS. `preflightPrerequisiteDoc` (`repoinstall.go`) carries the operator's verbatim final wording for both secrets. `TestFormatPreflightFailure_ContainsOperatorWording` asserts the rendered message (exercised through the real `Run()` call, not `formatPreflightFailure` in isolation) contains `"claude setup-token"`, `"fine-grained"`, `"Contents"`, `"Issues"`, `"Pull requests"`, `"Workflows"`, `"Settings"`, `"Developer settings"`, `"Personal access tokens"`, and `"Fine-grained tokens"` — matching the operator's exact final wording quoted in the issue's "Decisions (operator, 2026-08-05)" comment.

**AC8 — quickstart defines every term before use; states plainly there is no bot account today; no step assumes a pre-existing account.**
PASS. `docs/guides/INSTALL-CDS.md`'s `## Terms` section (defining PAT, repo secret, default branch, "bot") now sits immediately after the audience line — before the "two layers" paragraph that first uses "PAT"/"default branch", and before every "bot" mention in the doc. Verified directly: `grep -in "bot" docs/guides/INSTALL-CDS.md` → the first hit is inside the Terms section itself (line 29 of the final doc); no earlier occurrence exists (I confirmed this with an `awk` line-number comparison against the definition's line number). `grep -in "sigma@cnos.cn-sigma.cnos\|41898282\|SIGMA_WORKFLOW_PAT" docs/guides/INSTALL-CDS.md` → zero matches. The Tier 3 runbook no longer instructs "create a bot" anywhere.

**AC9 — rendered workflow contains no cosmetic `bot_name`/`bot_id` by default; commits authored by the token's account.**
PASS. `grep -n "bot_name\|bot_id" .github/workflows/cnos-cds-dispatch.yml` → zero matches (previously lines 89–90). `cn-install-wake`'s `agent_bot_name()`/`agent_bot_id()` lookup table deleted outright — the render step now conditionally emits `bot_name:`/`bot_id:` only when the caller passes `--bot-name`/`--bot-id` explicitly (`if [ -n "$bot_name" ]; then ... fi`). `TestRun_DispatchCds_SigmaDefault_NoIdentityFlagsRequired` (extended) renders a fresh sigma-default workflow and asserts neither key, nor the old cosmetic-identity strings, appear. `TestRun_DispatchCds_BotFlags_StillOptIn` confirms `--bot-name`/`--bot-id` still land in the render when explicitly passed (opt-in preserved, not accidentally deleted).

**AC10 — bot-less path is the documented default; dedicated bot is the future upgrade (cnos#449/#702); any `GITHUB_TOKEN`-only minimal claim is verified before being documented.**
PASS. `docs/guides/INSTALL-CDS.md`'s "Tenant secrets, by tier" table's Tier 3 row + the Tier 3 runbook present the two-own-account-secrets flow as the default, with the dedicated bot explicitly named as future/deferred (cnos#449 / cnos#702), not a required step. `grep -n "GITHUB_TOKEN" docs/guides/INSTALL-CDS.md` → every `GITHUB_TOKEN`-only claim stays scoped to the Tier 2 (`--engine`) section (confirmed by reading each hit in context); the doc carries forward the operator's own already-verified PAT-vs-`GITHUB_TOKEN` rationale from the issue thread (workflow-write at install time; loop-prevention blocking `GITHUB_TOKEN`-triggered downstream workflows at runtime) rather than re-deriving or asserting a new unverified claim.

## Self-check

Did α's work push ambiguity onto β? Two soft judgment calls are named explicitly rather than left implicit:

1. **Preflight is skipped entirely under `--dry-run`.** `Run()` gates the new preflight call on `opts.Dispatch == "cds" && !opts.DryRun`. Rationale: `--dry-run` writes nothing regardless (so AC2's "no partial artifacts" bar is trivially satisfied) and the existing dry-run path never reached label-doctor either — extending the exemption to preflight keeps dry-run's existing "no live credential required" property intact rather than introducing a new inconsistency. No test exercises `--dry-run --dispatch cds` (none existed pre-cycle either); this is a scoped, disclosed choice, not a gap I'm asking β to discover independently.
2. **Push-access check uses the single-call `GET /repos/{owner}/{repo}` `permissions.push` field**, per γ's Friction note 3 recommendation, not the two-call collaborator/permission endpoint. Rationale (recorded here per the scaffold's explicit instruction): since AC9 deletes the separate bot-identity concept, "push access" reduces to "does the installing token have push access to this repo" — which `permissions.push` answers directly in one round trip, with no need to first resolve which account the token belongs to via a separate `GET /user` call.

Every AC claim above cites a specific test name or a grep command I ran directly against branch HEAD (not paraphrased or assumed); §CDD Trace step 6 below enumerates every file in the diff, matching `git diff --stat origin/main..HEAD` exactly (pre-review gate row 11).

## Debt

1. **No fully-successful (`err == nil`) end-to-end `Run()` test for `--dispatch cds`.** Every dispatch-cds test in `repoinstall_test.go`, before AND after this cycle, ends at label-doctor's own "could not resolve target repo" error (repoRoot has no git "origin" remote in these fixtures) — this is pre-existing, unrelated to #706. My new preflight tests reuse this same pattern (env-var-based repo/API-base overrides that bypass git entirely for MY check, deliberately NOT git-initializing repoRoot, so label-doctor's OWN downstream resolution is left exactly as before — see `repoinstall_test.go`'s `setPreflightSatisfiedEnv` doc comment for the full rationale, including why git-initializing repoRoot to satisfy MY preflight would have been actively harmful: it would ALSO make label-doctor's resolution succeed and attempt a real, uncontrolled network call). A true "prerequisites present → install proceeds with ZERO error" test would additionally require overriding label-doctor's own hardcoded `githubAPIBase`, a different package's white-box-only test seam not reachable from this package. `TestRun_DispatchCds_PreflightSatisfied_SecondRunByteIdentical` evidences the idempotent-re-run half of AC4 without needing that.
2. **`CN_INSTALL_PREFLIGHT_API_BASE_URL` / `CN_INSTALL_PREFLIGHT_REPO` are internal test-only env vars, not part of the public CLI contract.** They exist solely so CLI-level tests (which have no Options-level seam into `RepoInstallCmd.Run`) can point preflight at an `httptest.Server`. A future cycle wanting real GitHub Enterprise Server support could reconsider exposing an equivalent as an actual `--flag`, but that's out of this cycle's scope guardrails.
3. **`docs/development/design/cn-repo-install-MOCKS.md`** still references the pre-#706 design (per γ's scaffold, marked "reference only, not required reading for implementation") — left untouched, consistent with the scaffold's own disposition of that file. It is a historical design doc, not user-facing.

## CDD Trace

Step 7 of `cdd/alpha/SKILL.md`'s CDD Trace (through step 7 — the pre-review gate and review-readiness signal are appended as separate commits below, per §2.5/§2.6/§2.7).

1. Issue: cnos#706, consolidated final spec (issue comment "CONSOLIDATED FINAL SPEC (for dispatch, 2026-08-05)").
2. Mode: `design-and-build`, per γ scaffold.
3. Contract: γ's `.cdd/unreleased/706/gamma-scaffold.md` (this cycle's authoritative implementation contract — per-AC oracles, source-of-truth table, α prompt; committed by γ at `3bf1b2d`, not authored by α — listed here for completeness against `git diff --stat origin/main..HEAD`).
4. Base SHA: `main` HEAD `7f249ddbb50f230d5d41287b6554ab17b5a1d1d5` (cycle branch cut point, per γ scaffold header; confirmed still current — `git merge-base --is-ancestor origin/main HEAD` succeeds, no rebase needed as of this writing).
5. Design: converged in-issue-thread (no stable `DESIGN.md` path exists; MCA preconditions not met, not a blocker).
6. Diff (every file in `git diff --stat origin/main..HEAD`, mapped to the AC(s) it evidences):
   - `.cdd/unreleased/706/CLAIM-REQUEST.yml` — carried forward from γ's session (Friction note 6), no AC.
   - `.cdd/unreleased/706/gamma-scaffold.md` — γ's own artifact (committed before α's dispatch), no AC.
   - `.cdd/unreleased/706/self-coherence.md` — this file.
   - `go.work` — registers the new `install-preflight` module.
   - `src/packages/cnos.core/commands/install-preflight/{go.mod,preflight.go,github.go,resolve.go,preflight_test.go,github_test.go,resolve_test.go}` — new package (AC1/AC2/AC3/AC4).
   - `src/go/internal/repoinstall/repoinstall.go` — `runPreflight`, `requiredSecretNames`, `preflightPrerequisiteDoc`, `formatPreflightFailure`, the `Run()` call-site insertion, and the `SIGMA_WORKFLOW_PAT` → `CN_DISPATCH_PAT` display-value rename (AC1/AC2/AC3/AC4/AC6/AC7).
   - `src/go/internal/repoinstall/repoinstall_test.go` — preflight test fixtures (`setPreflightSatisfiedEnv`/`setPreflightMissingEnv`) applied to every pre-existing dispatch-cds test, plus new tests: `TestRun_DispatchCds_PreflightRunsBeforeAnythingElse_NoPartialArtifacts` (AC1/AC2), `TestFormatPreflightFailure_ContainsOperatorWording` (AC7), `TestRun_DispatchCds_BotFlags_StillOptIn` + extended `TestRun_DispatchCds_SigmaDefault_NoIdentityFlagsRequired` (AC9), `TestRun_DispatchCds_PreflightSatisfied_SecondRunByteIdentical` (AC4).
   - `src/go/internal/cli/cmd_repo_install.go` — help text: preflight gate description, `CN_DISPATCH_PAT` rename, `docs/guides/INSTALL-CDS.md` links (AC5/AC6).
   - `src/go/internal/cli/cmd_repo_install_test.go` — `setPreflightSatisfiedEnv` fixture applied to all three CLI-level dispatch-cds tests; `TestRepoInstall_HelpFlag` extended (AC5).
   - `src/packages/cnos.core/commands/install-wake/cn-install-wake` — secret rename, deletion of `agent_bot_name()`/`agent_bot_id()`, conditional `bot_name:`/`bot_id:` emission (AC6/AC9).
   - `.github/workflows/cnos-cds-dispatch.yml`, `src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml`, `src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml` — regenerated via the renderer, confirmed idempotent + sha256-matched (AC6/AC9).
   - `docs/guides/INSTALL-CDS.md` — revised in place (AC5/AC7/AC8/AC10).
7. Caller-path trace for the new module: `installpreflight.Check` is called from `src/go/internal/repoinstall/repoinstall.go`'s `runPreflight` (non-test caller), which is called from `Run()`'s `if opts.Dispatch == "cds" && !opts.DryRun` block — a real, reachable production call path, not a dead module.

## Review-readiness

Pre-review gate (§2.6), re-validated immediately before this signal:

1. **Cycle branch rebased onto main.** `git fetch origin main && git merge-base --is-ancestor origin/main HEAD` succeeds — `origin/main` is still at `7f249ddbb50f230d5d41287b6554ab17b5a1d1d5` (unchanged since the branch was cut), no rebase needed.
2. **CDD Trace through step 7** — present above.
3. **Tests present.** 38 test functions in `internal/repoinstall` (0 failing), 53 in `internal/cli` (0 failing), 15 in `packages/cnos.core/commands/install-preflight` (0 failing) — counts taken directly from `go test -v -count=1` runner output (`grep -c '^--- PASS'`), not manually enumerated. Full `cd src/go && go test ./...` is green across all 15 packages.
4. **Every AC has evidence** — §ACs above, each citing a named test or a grep I ran directly.
5. **Known debt is explicit** — §Debt above (3 items).
6. **Schema/shape audit** — `ghSecret`'s wire shape was verified against GitHub's real documented response shape (no `value` field), not assumed.
7. **Peer enumeration** — the "two sigma bindings" (bot-identity table + PAT-secret default) were both located and both fixed; `SIGMA_WORKFLOW_PAT` and `bot_name`/`bot_id` greps were run across every file the scaffold named (renderer, repoinstall.go, cmd_repo_install.go, live workflow, both golden fixtures) — not just the first hit.
8. **Harness audit** — not applicable (no CI-emitted example/fixture derived from this schema exists outside the golden-fixture mechanism already exercised).
9. **Post-patch re-audit** — not applicable (no fix-round yet; this is the initial R0 submission).
10. **Branch CI is green on the head commit.** Verified directly via `gh api repos/usurobor/cnos/commits/<head-sha>/check-runs` at the moment of this signal: all 11 check runs (`Go build & test`, `CDD artifact ledger validation (I6)`, `Package verification`, `Binary verification`, `Repo link validation (I4)`, `Workflow + design-template parse guard`, `Protocol contract schema sync (I2)`, `SKILL.md frontmatter validation (I5)`, `Dispatch closeout-integrity guard`, `Dispatch repair-preflight guard`, `Package/source drift (I1)`) report `completed`/`success` at commit `79ed10b028f17ea196a6a6c3956d1fd582767732` as of 2026-08-06T01:31:30Z. (Earlier commits in this cycle's history show CI red — `cn cdd verify`'s `## ACs`/`## CDD Trace`/`## Self-check`/`## Debt` section check failed on every commit before this file carried those sections, per §2.5's incremental-write discipline; this is expected and resolved by this final self-coherence commit.)
11. **Artifact enumeration matches diff** — §CDD Trace step 6 above enumerates every file in `git diff --stat origin/main..HEAD` (verified by direct comparison, 20/20 files).
12. **Caller-path trace for new modules** — §CDD Trace step 7 above.
13. **Test assertion count from runner output** — row 3 above, pasted from real `go test -v` output, not estimated.
14. **α's commit author email** — `git log -1 --format='%ae' HEAD` → `alpha@cdd.cnos`, matching the canonical role pattern; verified for every α commit on this branch (`git log --format='%h %ae' origin/main..HEAD`), all consistent.
15. **γ-side artifact presence** — `git cat-file -e origin/cycle/706:.cdd/unreleased/706/gamma-scaffold.md` succeeds: γ-artifact at canonical §5.1 path.

**Post-signal CI confirmation.** Per the SHA convention this file uses ("implementation SHA" — the last implementation commit before the readiness-signal commit itself, `79ed10b028f17ea196a6a6c3956d1fd582767732`), the readiness-signal commit that follows (adding this §Review-readiness section) necessarily advances HEAD past that SHA. Re-checked directly: the signal commit itself is ALSO green — `gh api repos/usurobor/cnos/commits/3e06c20.../check-runs` reports all 11 check runs `success` (checked 2026-08-06T01:33:54Z). Every commit on this branch as of writing is CI-green.

## Review-readiness | round 1 | base SHA: 7f249ddbb50f230d5d41287b6554ab17b5a1d1d5 | implementation SHA: 79ed10b028f17ea196a6a6c3956d1fd582767732 | branch CI: green at 2026-08-06T01:31:30Z | ready for β

R0 is ready for β review. All 10 Final ACs are implemented with named evidence in §ACs above; the pre-review gate (§2.6) passed on every row as of this signal; branch CI is green on the head commit at the moment of signaling (not a stale earlier check). Known debt is disclosed in §Debt (none of it blocks any AC). β should independently walk each AC's oracle per the γ scaffold's β prompt rather than relying on this self-report, per that prompt's own instruction.
