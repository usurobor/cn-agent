# self-coherence — cycle/569

manifest:
  planned_sections: [Gap, Skills, ACs, Self-check, Debt, CDD Trace, Review-readiness]
  completed: [Gap, Skills, ACs, Self-check, Debt]

## §Gap

- **Issue:** [#569](https://github.com/usurobor/cnos/issues/569) — "cds/issues: FSM Phase 2 — authority flip (FSM applies labels; workers request transitions)"
- **Parent/wave:** #567 (master). Depends on #568 (Phase 1 read-only reconciler, closed/merged) and #570 (cell-kind taxonomy, closed/merged).
- **Mode:** design-and-build. γ's scaffold (`.cdd/unreleased/569/gamma-scaffold.md`) names one real open design decision left to α: the exact guard combination for the new `in-progress → review` proposal rule, and the 404-tolerance choice for the label-remove call. Both are resolved and documented below (§ACs, §Self-check).
- **Version/mode:** unreleased cycle, no release version cut yet; PR not opened by α per dispatch instructions (δ opens it after β converges).
- **Base:** cycle branch `cycle/569`, created by γ from `origin/main@0520235e1285c078eb3bc9d7eeba191b0413c53b`.
- **Binding scope constraint (operator comment, 2026-07-03T21:34:35Z on #569):** Phase 2 enforces FSM-controlled transitions **only** for the implementation-cell lifecycle already in `transitions.json` (`ready/todo/in-progress/review/changes`). No `cell_kind`-based enforcement branching; `FactSnapshot.CellKind` stays observed-only. Honored — see §ACs / scope-guardrail evidence below.

## Gap statement

Phase 1 (#568) gave the FSM the ability to *observe and reason about* issue state, read-only. It had no authority: workers (concretely, the δ-driven cds-dispatch wake) still wrote `status:*` labels directly, which is the root cause of the empty-review / stranded-in-progress / label-drift failure class the issue names (cnos#524 W4). This cycle moves that authority: the FSM applies status labels; workers produce matter (PR, commits, `REVIEW-REQUEST.yml`) and *request* transitions via a new guard-gated `--apply` flag on the existing `cn issues fsm evaluate` verb.

## §Skills

**Tier 1:** `CDD.md` (canonical lifecycle) + `alpha/SKILL.md` (this role surface).

**Tier 2 (`eng/*`, always-applicable):** `cnos.eng/skills/eng/SKILL.md` (coding bundle); `cnos.eng/skills/eng/go/SKILL.md` (Go conventions — package co-location per cnos#568/#556 precedent, dispatch-boundary rule INVARIANTS.md T-002 / `eng/go` §2.18, respected by keeping `cmd_issues_fsm.go` a thin doc-comment-only edit); `cnos.core/skills/write/SKILL.md`.

**Tier 3 (issue-specific):** `cnos.cdd/skills/cdd/issue/SKILL.md` (issue-pack contract); `cnos.cdd/skills/cdd/issue/proof/SKILL.md` (proof-plan discipline — AC oracle list was γ-authored, α proves each with fixture/CLI evidence below); `cnos.cdd/skills/cdd/issue/constraints/SKILL.md` (constraints discipline — the operator's cell-kind-deferral comment and the "no broad label redesign" non-goal are both hard constraints honored, not suggestions).

**Consumed but not loaded as role skills (read for context per §2.1.5 of `alpha/SKILL.md`, artifact enumeration):** `cnos.cds/skills/cds/fsm/{table,transitions}.json`-adjacent Go source (`table.go`, `snapshot.go`) to understand the guard-registry / rule-matching engine before extending its declarative data; `cnos.cds/orchestrators/cds-dispatch/SKILL.md`, `cnos.cdd/skills/cdd/delta/SKILL.md` §9.5–§9.6, `cnos.core/skills/agent/dispatch-protocol/SKILL.md` §2.4/§2.9 — the AC2 prose-edit targets named in the γ scaffold; `cnos.core/commands/install-wake/cn-install-wake` (the golden/live-workflow renderer, since editing cds-dispatch's SKILL.md body requires re-rendering `cnos-cds-dispatch.golden.yml` + `.github/workflows/cnos-cds-dispatch.yml` to stay in sync per `install-wake-golden` CI).

## §ACs

Implementation SHA (last implementation commit before this readiness signal): `f344c5c` (prose/golden). Go+JSON implementation landed at `901a5cb`.

### AC1 — `--apply` is guard-gated

**Invariant:** `--apply` mutates labels only when the evaluated transition's guards pass; idempotent on re-run.

**What was built:**
- `src/packages/cnos.issues/commands/issues-fsm/issuesfsm.go` — a new `--apply` bool flag on `evaluate` (not a new sub-verb; the flag's help text and package doc comment consume the "Label-write authority is Phase 2" forward reference). `runEvaluate` gates the mutation strictly on `dec.Outcome == "proposed" && dec.TargetState != ""` (i.e. only after `Evaluate` already matched a rule whose guards passed — table.go's `ruleMatches`). A `"blocked"` outcome sets a non-nil `applyErr` (nonzero exit) and never calls the write path; a `"valid"` outcome (or a `"proposed"` outcome with no direct status target, e.g. `propose_delta_recovery`) is a silent no-op.
- `src/packages/cnos.issues/commands/issues-fsm/fetch.go` — the label-write primitives: `ghRequest` (shared non-GET request builder, same auth-header idiom as the pre-existing `ghGetJSON`), `ghAddLabel` (POST `.../labels`), `ghRemoveLabel` (DELETE `.../labels/{name}`), and `applyStatusLabel` (the single mutation entry point: removes the old `status:{fromState}` label if present and different from the target, then adds `status:{toState}`). `githubAPIBase` is a package var (default `https://api.github.com`) so tests can point it at an `httptest.Server`.
- `src/packages/cnos.issues/commands/issues-fsm/decision.go` — `Decision` gains `ApplyAttempted bool` / `Applied bool`, additive-only in `Render` (see AC1's backward-compat sub-claim under AC5).

**Evidence (fixture + CLI, per the oracle):**
- Positive + idempotence, new review-request rule: `TestApply_ReviewTransitionAppliesOnGuardPassAndIsIdempotent` (`issuesfsm_test.go`) — first `Run()` call against `testdata/in-progress-review-request-with-matter.json` with `--apply --repo acme/widgets` against an `httptest.Server`: asserts `stdout` contains `applied: true` AND asserts the exact two requests received (`DELETE .../labels/status%3Ain-progress` then `POST .../labels`). Second `Run()` call, same issue, against `testdata/review-with-pr.json` (the post-mutation state: `status:review` + a PR, which the *pre-existing, unchanged* review-state rule evaluates as `outcome: valid`): asserts `applied: false` and **zero** requests reached the fake server — a real second call, not an assertion about intent.
- Positive + idempotence, pre-existing requeue rule: `TestApply_RequeueTransitionAppliesOnGuardPassAndIsIdempotent` — same two-call shape against `testdata/in-progress-dead-no-matter.json` → `testdata/todo.json` (the `todo` state's rule is unconditionally `valid/none`), proving the guard-gating generalizes beyond the new rule.
- Negative (blocked, zero writes): `TestApply_BlockedReviewRequestRefusesAndMutatesNothing` (new rule, `testdata/in-progress-review-request-no-matter.json`) and `TestApply_EmptyReviewStateBlocked` (the pre-existing, unchanged `review`-state empty-review rule, reused fixture `testdata/review-empty.json`) — both assert a non-nil `Run()` error (nonzero exit) and zero calls to the fake label-write server (the second test fails the fake server's handler itself via `t.Fatalf` if any request arrives, not just a counter check).
- Missing-`--repo`: `TestApply_MissingRepoErrors` — `--apply` with `$GITHUB_REPOSITORY` unset and no `--repo` errors even though the outcome is `proposed` (cannot silently skip the write).
- `applyStatusLabel` unit-level: `TestApplyStatusLabel_ToleratesNotFoundOnRemove`, `TestApplyStatusLabel_NoRemovalWhenFromStateEmpty` (see the 404-tolerance design decision below).
- Manual CLI smoke (not just unit tests) run against the built `cn` binary: `cn issues fsm evaluate --issue 601 --fixture testdata/review-empty.json` (no `--apply`) → `outcome: blocked`, exit 0, tail line `(read-only: no label was written; this is a proposal, not an action)` — byte-identical framing to Phase 1. `cn issues fsm evaluate --issue 701 --apply --fixture testdata/in-progress-review-request-no-matter.json --repo acme/widgets --token fake` → `outcome: blocked`, `applied: false`, **process exit code 1** (verified via `echo $?`).

**Design decision — guard combination for the new `in-progress → review` rule (γ left this open; α's call, documented per the scaffold's Friction notes):**

γ's non-binding suggestion was `any_true: [review_request_present, pr_has_commits]` with `all_false: [run_active]`. α used a different combination for two reasons:

1. **Not gated on `all_false: [run_active]`.** The primary caller of this rule (AC2's mechanism) is δ's *own* wake-invoked cycle requesting the transition synchronously via cds-dispatch step 7, *while its own workflow run is still executing* — `run_active` would read `true` (or the most-recent-run lookup would even see the currently-running run) at exactly the moment the request is made. Gating on `all_false: [run_active]` (mirroring the pre-existing *dead*-in-progress reconciliation rules, which are about a **different** use case — an external reconciler noticing a run has already died) would make AC2's entire "worker requests its own transition" mechanism inert in its main use case. This is the one real design tension the scaffold flagged; α resolved it by placing the new rule *before* the existing `all_true: [run_active] → valid/none` rule in the rules list, so a review request with matter fires regardless of run state.
2. **`all_true: [review_request_present]` (hard gate) instead of `any_true: [review_request_present, pr_has_commits]`.** AC2's own wording is "ensure REVIEW-REQUEST.yml **and** matter exist, then request" — a conjunction, not a disjunction. Requiring `REVIEW-REQUEST.yml` unconditionally (`all_true`), combined with `any_true: [pr_exists, pr_has_commits, branch_has_commits]` for the matter half, encodes exactly that conjunction and mirrors the evidence bar the *existing, unchanged* `review`-state validity rule already uses (`any_true: [pr_exists, branch_has_commits, review_request_present]`) — so once the transition lands, the destination state's own rule agrees the evidence is sufficient.
3. **Added an explicit second (blocked) rule** for `review_request_present == true` with no matter (`all_false: [pr_exists, pr_has_commits, branch_has_commits]`), rather than letting that case fall through to the pre-existing fallback rules (which would have silently proposed `propose_status_todo` or `propose_delta_recovery` instead — the *wrong* signal for a worker that explicitly asked for review). This is what makes AC3 mechanically true for the new rule specifically (see AC3 below), not just an artifact of the old `review`-state rule.

Neither new rule changes any Phase-1 fixture's outcome — both require `review_request_present == true`, which no Phase-1-era fixture sets while `in-progress` (verified: `TestAC4_DeadInProgressNoMatterProposesRequeue`, `TestAC5_DeadInProgressWithCommitsProposesDeltaRecovery`, `TestAC5_HealthyActiveInProgressIsValid` all still pass unmodified against the extended table).

**Design decision — 404-tolerance on label removal (γ flagged as an implementation-contract-adjacent call for α):** `applyStatusLabel`'s `ghRemoveLabel` call tolerates HTTP 404 (label already absent) as success, not an error. Rationale: idempotency across two live `--apply` calls mostly falls out of re-fetching live facts each call (the post-mutation state no longer matches the `proposed` rule) — but a manually-recovered state (operator already removed the old label by hand) or a partially-applied prior `--apply` (added the new label, crashed before removing the old one) must not turn a legitimate retry into a hard failure. `TestApplyStatusLabel_ToleratesNotFoundOnRemove` locks this. The add-label call is *not* given equivalent 404-tolerance because there is no analogous benign-404 case there — GitHub's add-labels endpoint is already idempotent for "label already present" (200/201, not 404), so `ghAddLabel` only accepts 200/201.

### AC2 — workers no longer own status labels

**Invariant:** workers produce matter and request transitions; the FSM applies status labels.

**What was built (prose-only, commit `f344c5c`):**
- `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` step 7: now reads "...requests the `status:in-progress → status:review` transition — δ ensures `REVIEW-REQUEST.yml` and the closeout matter exist, then runs `cn issues fsm evaluate --issue {N} --apply`... rather than writing the label directly." The "Closeout integrity preflight" intro + "No-deliverable rule" paragraphs, and the "Lifecycle transitions" table's β-converge row + its lead-in sentence, are updated the same way; the other three named transitions (claim, hard-block, release-back-to-queue) are explicitly called out as unchanged direct label writes — AC2's scope is the review request only.
- `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` §9.6's `status:review` return-token row: "Wake-observable mechanism" column now reads "a *request* for the `status:in-progress → status:review` transition... (`cn issues fsm evaluate --issue {N} --apply`... — δ does NOT write the label directly)".
- `src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md` §2.9: one added paragraph ("Mechanism (cnos#569 Phase 2)") stating the wake requests rather than writes the transition, and that the FSM guard is a second independent layer over the same deliverable-evidence invariant. §2.4 (the abstract lifecycle-transition diagram) was read and left unchanged — it names states and actors ("wake: claim operation", "cell completes") but never asserts a label-write mechanism, so it does not contradict AC2.
- `cnos-cds-dispatch.golden.yml` and `.github/workflows/cnos-cds-dispatch.yml` re-rendered via `cn-install-wake cds-dispatch` (twice: once to the golden default path, once with `--out .github/workflows/cnos-cds-dispatch.yml`) so both stay byte-identical to a fresh render of the edited SKILL.md — verified with `diff` (see AC5 evidence).

**Evidence:** `git diff origin/main...HEAD -- src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md src/packages/cnos.cdd/skills/cdd/delta/SKILL.md src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md` shows exactly the hunks above. Grepped all three files post-edit for `writes the label transition`, `then writes the label`, `sets \`status:review\`` (as a normative instruction, not a symptom description) — zero surviving hits asserting the wake/δ writes `status:review` as its own act outside the FSM-apply path. (`dispatch-protocol/SKILL.md`'s D12 troubleshooting entry still narrates the *historical failure symptom* "cell sets `status:review`" — left as-is since it is describing the pre-#569 failure mode being guarded against, not prescribing current behavior; it is followed immediately by "_Fix:_ §2.9" which now carries the corrected mechanism.)

### AC3 — empty-review is structurally impossible

**Invariant:** `status:review` cannot be reached without deliverable proof (PR, commits-beyond-base, or `REVIEW-REQUEST.yml`).

**What was built:** the new blocked rule under AC1 (`review_request_present == true`, no PR/commits/branch-commits → `outcome: blocked`), plus the *pre-existing, unmodified* `review`-state rule (empty `status:review` itself → `outcome: blocked`) — both are reachable and enforced under `--apply`.

**Evidence:**
- New in-progress rule: `TestAC569_InProgressReviewRequestNoMatterBlocked` (direct `Evaluate` call) + `TestApply_BlockedReviewRequestRefusesAndMutatesNothing` (through the CLI, nonzero exit, zero label-write calls to the fake server).
- Pre-existing review-state rule under the new `--apply` path: `TestApply_EmptyReviewStateBlocked` (reuses `testdata/review-empty.json`, the Phase-1 AC3 fixture, per the scaffold's explicit reuse instruction) — the fake server's handler itself calls `t.Fatalf` if any request arrives, the strongest form of "zero writes" assertion in this suite.
- `check-dispatch-closeout-integrity.sh --self-test` (the independent, Phase-1-era shell-level empty-review detector — unmodified by this cycle) still passes: `cnos#524 closeout-integrity self-test: empty-review detector correct`. AC3's own text explicitly allows satisfying it via "(a) the new `--apply` guard alone... (c) both"; α chose (a) — the FSM guard is the primary structural enforcement — and left the shell detector as an independent, unweakened second layer (see γ scaffold's "do not weaken or duplicate... without reason"; α found no reason to touch it, since it already worked and this cycle does not change its detection surface).

### AC4 — no broad label redesign

**Invariant:** existing labels/taxonomy remain stable.

**Evidence:** `git diff origin/main...HEAD | grep -E "^\+" | grep -iE "gh label create|label-doctrine|status:[a-z_-]+\"|protocol:[a-z_-]+\""` — every hit is either (a) this very AC4 oracle text being quoted verbatim from the γ scaffold (added to this repo's own scaffold file, not a taxonomy edit), (b) `url.PathEscape("status:in-progress")` test assertions (Go string literals asserting the *existing* label name is escaped correctly in a URL, not defining a new label), or (c) `"status:in-progress"` / `"status:todo"` inside the three new `testdata/*.json` fixtures (existing label names, test data, not label definitions). Zero `gh label create` additions; zero edits to any labels-definition surface.

### AC5 — current gates remain green

**Invariant:** the authority flip does not regress repo health.

**Evidence (all run locally on branch HEAD, commit `f344c5c` implementation content + this artifact's commits on top):**
- `cd src/go && go build ./... && go vet ./... && go test ./...` — all 14 `src/go` internal packages pass (`activate`, `activation`, `binupdate`, `cell`, `cli`, `discover`, `dispatch`, `doctor`, `hubinit`, `hubsetup`, `hubstatus`, `pkg`, `pkgbuild`, `restore`), `cmd/cn` has no test files (unchanged from baseline).
- `cd src/packages/cnos.issues/commands/issues-fsm && go build ./... && go vet ./... && go test ./... -race -v` — 27 top-level/subtest `--- PASS` lines, 0 failures (full list includes all pre-existing `TestAC1`–`TestAC8`-lineage tests unmodified in assertion shape, plus the new AC1/AC3/569-specific tests above).
- `cd src/packages/cnos.issues/commands/issues-map && go build ./... && go test ./...` and `cd src/packages/cnos.cdd/commands/cdd-verify && go build ./... && go test ./...` — both green (peer Go modules in the workspace, unaffected but re-verified per the harness-audit discipline).
- `./cn build --check` (I1, package/source drift) — `✓ All packages valid.`
- `diff docs/reference/schemas/protocol-contract.json tests/fixtures/protocol-contract.json` (I2) — no diff.
- I5 (skill-frontmatter, requires `cue` CLI not installed in this sandbox) not run directly; verified instead that every edited SKILL.md's diff hunks are well past the frontmatter block (`git diff` hunk headers `@@ -219,...`, `@@ -470,...`, `@@ -444,...` — all >200 lines into files whose frontmatter ends by line ~15), so frontmatter is untouched and I5 has no surface to regress.
- `cn cdd verify --unreleased` (I6) — exits 1 on **both** `origin/main` and this branch, with **byte-identical** warning/error output (`diff` of the two runs' `❌`/`⚠️` lines is empty) — the single failure (`issue #512` missing `self-coherence.md`) is pre-existing repo debt unrelated to this cycle, not a regression (see §Debt).
- `./scripts/ci/check-dispatch-closeout-integrity.sh` and `./scripts/ci/check-dispatch-repair-preflight.sh` — both green (see AC2/AC3 evidence above).
- `cn-install-wake cds-dispatch` (golden re-render) and `cn-install-wake cds-dispatch --out .github/workflows/cnos-cds-dispatch.yml` (live re-render) — `diff` between the two outputs is empty (byte-identical), satisfying `install-wake-golden`'s re-render-and-diff check. `cn-install-wake agent-admin` re-run for completeness (untouched wake) — reported "(unchanged)".
- Dispatch-boundary check (INVARIANTS.md T-002) — manually re-ran the CI grep against `internal/cli/cmd_*.go`: `cmd_issues_fsm.go`'s edit is a doc-comment-only change, no new imports; grep reports clean.
- `scripts/kata/run-all.sh` (binary-verify / package-verify's Tier 1 kata suite) — `KATA SUITE: all passed` (Kata 01 Deps/Help through Kata 06 Install, 31 total pass assertions across the 6 kata).

## §Self-check

**Did α's work push ambiguity onto β?**

- The one open design decision the scaffold flagged (guard combination for the new `in-progress → review` rule, including the deliberate deviation from γ's non-binding `all_false: [run_active]` suggestion) is resolved and reasoned through in §ACs/AC1 above, not left for β to adjudicate. β's job is to verify the reasoning holds, not to make the call itself.
- The 404-tolerance decision for label removal is likewise resolved and locked by a dedicated unit test (`TestApplyStatusLabel_ToleratesNotFoundOnRemove`), not left implicit.
- One residual gap: I5 (skill-frontmatter-check) could not be run locally because the `cue` CLI is not installed in this sandbox. Mitigated by checking that every edited SKILL.md's diff hunks land well past the frontmatter block (verified via `git diff` hunk headers), so there is no plausible frontmatter regression — but β (or CI) should still confirm I5 green on the actual PR, since this is the one AC5 gate α could not execute directly rather than reason about structurally.
- AC4's evidence is a grep-based enumeration, not an exhaustive semantic audit of "what counts as a taxonomy change." α re-read every grep hit individually (not just the hit count) and classified each; β should re-run the same grep independently per its own oracle rather than trust α's classification alone (per `alpha/SKILL.md` §"Do not claim structural closure without exhaustive enumeration").

**Peer enumeration (per `alpha/SKILL.md` §2.3):**

- **Sibling commands at the same layer:** `cn issues map` (`issues-map` package) and `cn cdd verify` (`cdd-verify` package) — both are Go modules in the same workspace, co-located under the same package-command pattern as `issues-fsm`. Neither consumes `transitions.json` or the FSM's `Decision`/`FactSnapshot` shapes, so neither needed a code change; both were re-built and re-tested green as part of AC5's harness audit (§ACs/AC5 evidence) rather than assumed unaffected.
- **Multiple writers/readers of the extended schema:** `transitions.json` (CDS-owned declarative data) has exactly one consumer, `table.go`'s `Evaluate`/`ruleMatches`/`guardFuncs` — unmodified this cycle (extension is data-only, per the scaffold's explicit "do not re-derive Phase-1 transition logic" guardrail). `FactSnapshot` (the JSON wire shape) has one producer path duplicated twice (`assembleLive` for live mode, `LoadFixture` for fixtures) and one consumer (`Evaluate`); neither producer changed shape this cycle — `review_request_present` already existed in Phase 1's `FactSnapshot`/fixture schema, the new rules only read it.
- **Prose ↔ rendered-substrate peers:** `cds-dispatch/SKILL.md` (source) has two rendered peers that must byte-match a fresh render — `cnos-cds-dispatch.golden.yml` and `.github/workflows/cnos-cds-dispatch.yml`. Both re-rendered and diffed identical (§ACs/AC2, AC5 evidence). `agent-admin`'s golden was re-run for completeness (not a peer of this change, since its SKILL.md was untouched) and confirmed "(unchanged)".
- **Intra-doc repetition (per §2.3's numeric/named-value drift rule):** the phrase "writes the label transition" / "sets `status:review`" as a *normative* (not historical-symptom) claim was grepped across all three touched prose files after editing, not just at the one site named in the γ scaffold — zero surviving hits (§ACs/AC2 evidence).

**Harness audit for schema-bearing changes (per `alpha/SKILL.md` §2.4):** `transitions.json` is schema-bearing CDS data; its one Go consumer (`table.go`) and its one non-Go "harness" (none exists — no shell/CI script parses `transitions.json` directly; `check-dispatch-closeout-integrity.sh` implements an *independent* empty-review predicate in shell, not a `transitions.json` reader, so it is a peer detector, not a harness of this schema) were both identified and left correctly alone. `FactSnapshot`'s JSON shape is unchanged (no new fields), so no fixture-writer or CI-emitted-example audit was needed beyond the three new `testdata/*.json` fixtures α authored to the existing schema.

**Every claim backed by evidence in the diff?** Yes — §ACs cites concrete test names, file paths, and command output for every AC; no claim above rests on "should work" without a corresponding test or manual run captured in this cycle.

## §Debt

- **I5 (skill-frontmatter-check) not executed locally.** The `cue` CLI is unavailable in this sandbox. Structural argument for why it should still pass (frontmatter blocks untouched — see §Self-check) is documented, but this is reasoning, not a run. β / CI should confirm.
- **Pre-existing I6 failure, unrelated to this cycle.** `cn cdd verify --unreleased` reports one failure (`.cdd/unreleased/512/self-coherence.md` missing) on **both** `origin/main` and `cycle/569` HEAD, byte-identical output before/after this cycle's changes (verified via diff of the two runs' warning/error lines). This is stale repo debt from an unrelated older cycle (#512), not introduced or touched by cnos#569. Not fixed in this cycle (out of scope — #512 is not named anywhere in this issue or scaffold, and "fixing" it would mean guessing at #512's actual missing artifact content, which α has no authority or context to author).
- **`assembleLive`'s live-mode path for `--apply` is not covered by an automated test that exercises the full live-fetch + live-apply round trip against a real GitHub API shape** (only the fixture-facts + fake-server-apply combination is tested, per the scaffold's own oracle: "or `--fixture` for hermetic tests"). This mirrors Phase 1's own test-coverage shape (`assembleLive`'s live GET path was likewise never exercised against a real network in Phase 1's test suite) — not a regression, but also not newly closed by this cycle. A future cycle wiring an end-to-end live smoke (e.g. against a disposable test repo) would strengthen this, but is out of this issue's scope (issue's own proof plan names fixture-based positive/negative cases, not a live-network integration test).
- **No PR opened.** Per dispatch instructions, α does not open the PR this cycle — δ does, after β converges.
- **γ's non-binding guard-combination suggestion was not adopted verbatim.** This is not "debt" in the technical-shortfall sense, but flagging explicitly per the process discipline: α deviated from γ's suggested `any_true`/`all_false: [run_active]` shape and documented why in §ACs/AC1. If β disagrees with the reasoning, this is the one design surface most likely to generate a fix-round finding.
