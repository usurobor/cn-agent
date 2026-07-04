# γ scaffold — cycle/574

## Issue reference / mode

- Issue: [#574](https://github.com/usurobor/cnos/issues/574) — "fix(cds/fsm): harden review guards + remote-branch observation; correct wave-closure language (post-#567 remediation)"
- Parent/wave: #567 (master, CLOSED/COMPLETED) — this cell is an operator-review remediation of the #567 wave (#568 Phase 1, #570 cell-kind doctrine, #569 Phase 2), all three merged and closed, but found **not fully clean**.
- Mode: **design-and-build** (per issue header). This is a bug-fix hardening cell, not greenfield design: AC2/AC3's exact target guard strings are already fully specified by the issue's own "Status truth" table, and AC1/AC5/AC6 are verification/correction tasks. The one place α does real design work is AC4 (fetch-vs-API choice for remote-branch observation) — see §Friction notes.
- Base SHA discrepancy (recorded per γ pre-flight, mirroring cycle/569's precedent): δ's wake-invoked input named current main SHA as **`b93ab0d75b68e269e749bdd26c6181abc958971e`**. At scaffold time, `git fetch origin main && git rev-parse origin/main` returned **`452191fe28bae8f7fbad53fe2010d4c122645342`** — two commits ahead. `git log b93ab0d..452191f` shows both intervening commits are `board-map: regenerate docs/development/board from live open issues`, touching only `docs/development/board/{board-data.json,index.html}` (doc-data auto-regeneration, not a CDD/CDS/handoff/dispatch-protocol/role-skill file). Per `gamma/SKILL.md` §2.5 branch pre-flight ("base SHA known") and the "Spec-staleness propagation" → "When not to propagate" clause, γ branched from the **actual current `origin/main` HEAD** (`452191fe`), and no staleness-propagation action is required.

## Cycle branch

`cycle/574`, created from `origin/main@452191fe28bae8f7fbad53fe2010d4c122645342`, pushed to `origin/cycle/574`.

## Surfaces α is expected to touch

CDS-owned declarative data (the table's designed extension point — no Go source change needed to tighten a guard):

- `src/packages/cnos.cds/skills/cds/fsm/transitions.json` — tighten the `review`-state rule set (AC2) and the `in-progress` state's `propose_status_review`/`block` rule pair (AC3). Exact current strings are quoted in full in §Source-of-truth and §Per-AC oracle list below — do not re-derive them, edit them in place.

Go (evaluator + CLI + observation), only if AC4 is resolved via fetch.go (option (b)):

- `src/packages/cnos.issues/commands/issues-fsm/fetch.go` — `assembleLive`'s branch-existence/commit-count block (lines ~62-73) currently shells out to **local** `git rev-parse --verify --quiet refs/heads/cycle/{N}` and `git rev-list --count main..cycle/{N}` only; there is no GitHub-API fallback/primary for branch/commit observation anywhere in this file. If α chooses option (b), this is the block to extend (reuse the existing `ghGetJSON` dependency-free pattern already used for PR/run observation lower in the same function — do not add a third-party GitHub client).
- `src/packages/cnos.issues/commands/issues-fsm/issuesfsm_test.go` — extend for AC2/AC3 (write the failing fixture FIRST per the issue's explicit TDD instruction), and for AC4 if a new fixture/scenario is added (a remote-only `cycle/{N}` case).
- `src/packages/cnos.issues/commands/issues-fsm/testdata/*.json` — existing fixtures relevant here: `review-empty.json` (already the empty-review negative case — verified content: `status:review`, `branch_exists:true, commits_beyond_base:0, pr_exists:false, review_request_present:false` — this fixture already fails AC2's *new* stricter rule the same way, since `pr_exists`/`pr_has_commits`/`review_request_present` are all false; the issue's AC2 fixture (`review_request_present:true, pr_exists:false`) is new — no existing fixture has `review_request_present:true` combined with no matter, so α must author it), `in-progress-review-request-with-matter.json` and `in-progress-review-request-no-matter.json` (already shipped by #573 — verify AC3's target rule doesn't flip the "no-matter" fixture's expected outcome away from blocked, and confirm the "with-matter" fixture still requires `pr_exists` specifically once the guard tightens from `any_true` to `all_true` — the with-matter fixture must be checked for whether it sets `pr_exists`/`pr_has_commits` or only `branch_has_commits`, since AC3 requires PR presence, not just a branch).

CDS/dispatch prose (AC5's named surfaces — read current state before editing, do not assume the overclaim survives everywhere):

- `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` and `src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md` §2.9 — **already read; already correctly scoped.** Both files already say δ/the wake "**requests**" the `in-progress → review` transition via `cn issues fsm evaluate --issue {N} --apply` and that "the FSM applies the label only when its guards pass" (cds-dispatch/SKILL.md lines 207, 217, 234, 243, 248; dispatch-protocol/SKILL.md line 447). Neither file contains the literal overclaim string "FSM owns status labels". α should re-grep before editing (`rg -i "FSM owns" src/`) rather than assume a fix is needed here — see §Friction notes on where the actual overclaim lives.
- The actual overclaim is a **GitHub comment**, not a repo file: #567's own wave-closure comment (2026-xx, posted by the prior γ/κ) is titled `## ✅ FSM wave complete — \`Workers produce matter. The FSM owns status labels.\`` and says in its body "**What changed for real:** the live cds-dispatch wake no longer self-writes status:review... Claim/block/requeue keep their pre-existing serialized-claim writes" — which is actually accurate prose in the body, but the bolded headline overclaims. α must correct or annotate this comment (`gh issue comment 567 --body ...` appending a correction, since GitHub comments by other authors typically cannot be edited by α's identity — verify edit permission first; if edit is not possible, post a follow-up correction comment referencing the original) so no shipped surface (issue thread included) still asserts the untrue absolute.
- File a Phase 3 tracking issue (deferred/held label) for routing claim/hard-block/release-back-to-queue through the FSM, linked from #567.

CI (AC6):

- Byte-level bidi/hidden-Unicode sweep over the wave-touched file set from #571/#572/#573 (union, deduplicated): `.github/workflows/build.yml`, `.github/workflows/cnos-cds-dispatch.yml`, `.github/workflows/install-wake-golden.yml` (verify if #571 touched this — not listed in the three PRs' file lists, include if the AC6 sweep is meant to cover CI workflows broadly), `go.work`, `src/go/cmd/cn/main.go`, `src/go/go.mod`, `src/go/internal/cli/cmd_issues_fsm.go`, `src/packages/cnos.cds/skills/cds/fsm/transitions.json`, `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md`, `src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml`, `src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md`, `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md`, `src/packages/cnos.cdd/skills/cdd/CDD.md`, `src/packages/cnos.cdd/skills/cdd/CELL-KINDS.md`, `src/packages/cnos.cdd/skills/cdd/issue/SKILL.md`, `src/packages/cnos.issues/SKILL.md`, `src/packages/cnos.issues/commands/issues-fsm/*.go`, `src/packages/cnos.issues/commands/issues-fsm/testdata/*.json`, `docs/reference/governance/GLOSSARY.md`, plus the six `.cdd/unreleased/{568,569,570}/*.md` closeout artifacts. #572 is the PR GitHub flagged with a hidden/bidi warning — sweep it first/specifically, then the rest of the union.

## Per-AC oracle list

### AC1 — authenticated wave-issue-state verification + discrepancy record

- **Invariant:** #567/#568/#569/#570 are CLOSED with no dispatch/active status labels; any public-render discrepancy is recorded, not shrugged.
- **γ pre-verified (2026-07-04, authenticated `gh issue view --json state,stateReason,labels`):**
  - #567: `state: CLOSED, stateReason: COMPLETED`, labels `[P1, area/cdd, area/cds, area/wake, kind/tracking, area/issues]` — no `status:*`/`dispatch:cell`/`protocol:*` label present.
  - #568: `state: CLOSED, stateReason: COMPLETED`, labels `[P1, kind/tooling, area/cds, area/issues]`.
  - #569: `state: CLOSED, stateReason: COMPLETED`, labels `[P1, area/cdd, kind/process, area/cds, area/wake]`.
  - #570: `state: CLOSED, stateReason: COMPLETED`, labels `[P1, area/cdd, kind/process, area/cds, area/coherence, area/wake, area/issues]`.
  - All four already clean — this AC's remaining work for α is the **public-render discrepancy check** (compare the unauthenticated/public HTML render of each issue against the authenticated API state above) and recording the result (clean or discrepant) in `self-coherence.md`, not re-verifying closure itself.
- **α proves via:** record the four authenticated states above (already gathered, reusable verbatim) in `self-coherence.md`, then check public render (`curl`/`WebFetch` the public issue page or note environment cannot render unauthenticated HTML and say so explicitly) and record discrepancy or absence of one.
- **β verifies via:** re-run `gh issue view {N} --json state,stateReason,labels` for all four and diff against `self-coherence.md`'s recorded values.

### AC2 — `review`-state validity requires deliverable matter, not review-request alone

- **Invariant:** `status:review` is valid only when `review_request_present AND pr_exists AND pr_has_commits`.
- **Current guard (verbatim, `transitions.json` `review` state, rule 1):**
  ```json
  {
    "any_true": ["pr_exists", "branch_has_commits", "review_request_present"],
    "outcome": "valid",
    "action": "none",
    "reason": "review evidence present (PR, commits beyond base, or REVIEW-REQUEST.yml); status:review is valid."
  }
  ```
  Rule 2 (blocked case) currently uses the symmetric `all_false` over the same three-guard set. **This is the guard AC2 requires α to tighten**: replace `any_true: [pr_exists, branch_has_commits, review_request_present]` with `all_true: [review_request_present, pr_exists, pr_has_commits]` on the valid rule, and correspondingly widen the blocked rule's condition to fire whenever the `all_true` set is not fully satisfied (not just when all three are false — a partial-evidence state, e.g. `pr_exists:true, pr_has_commits:false, review_request_present:true`, must also block, since the tightened invariant requires all three, not any).
- **α proves via:** write a fixture FIRST at `status:review` with `review_request_present:true, pr_exists:false, pr_has_commits:false` (the exact issue-specified case) and confirm it currently evaluates `valid` against unmodified `transitions.json` (FAIL), then tighten the guard and confirm it re-evaluates `blocked` (PASS). Also add/verify a partial-evidence fixture (`pr_exists:true, pr_has_commits:false, review_request_present:true`) blocks under the tightened rule, since a naive `all_false`→`negation-of-all_true` rewrite must be checked, not assumed correct.
- **β verifies via:** run the review-request-only fixture and the partial-evidence fixture through `evaluate` against the tightened table; confirm both `blocked`; confirm the existing `review-with-pr.json` fixture (PR + commits) still evaluates `valid` (no regression).

### AC3 — `in-progress → review` requires a PR (aligned with the prompt preflight)

- **Invariant:** `in-progress → review` is proposed only when `review_request_present AND pr_exists AND pr_has_commits`; branch-commits-only (no PR) does NOT qualify.
- **Current guard (verbatim, `transitions.json` `in-progress` state, rule 1 — added by #569 Phase 2):**
  ```json
  {
    "all_true": ["review_request_present"],
    "any_true": ["pr_exists", "pr_has_commits", "branch_has_commits"],
    "outcome": "proposed",
    "action": "propose_status_review",
    "target_state": "review",
    "repair_pass": false,
    "reason": "cnos#569 Phase 2: review requested (REVIEW-REQUEST.yml present) with deliverable matter (a PR and/or commits beyond base); in-progress -> review requested by the worker."
  }
  ```
  Rule 2 (blocked, immediately following) currently reads `all_true: [review_request_present]` + `all_false: [pr_exists, pr_has_commits, branch_has_commits]`. **This is the guard AC3 requires α to tighten**: change rule 1's `any_true: [pr_exists, pr_has_commits, branch_has_commits]` to `all_true: [review_request_present, pr_exists, pr_has_commits]` (drop `branch_has_commits` from qualifying alone), and widen rule 2's blocked condition symmetrically (branch-only, e.g. `branch_has_commits:true, pr_exists:false`, must now fall into blocked, not stay unmatched/fall through to the dead-run rules below).
- **Preserve exactly:** rules 3-5 further down the same `in-progress` state (`run_active` non-gating → `valid`/`none`; dead-run-with-matter → `propose_delta_recovery`; dead-run-no-matter → `propose_status_todo`) are untouched — the `_doc_phase2` block's stated rationale for NOT gating on `all_false: [run_active]` in rule 1 applies unchanged to the tightened version; do not add a `run_active` gate to rule 1 or 2.
- **α proves via:** write a fixture FIRST — `status:in-progress`, `review_request_present:true, branch_has_commits:true, pr_exists:false, pr_has_commits:false` — confirm it currently evaluates `proposed`/`propose_status_review` (FAIL against the new invariant), tighten, confirm it now evaluates `blocked` (PASS) or at minimum does not propose `review`. Existing fixtures `in-progress-review-request-with-matter.json`/`in-progress-review-request-no-matter.json` (from #573) must be re-checked against their actual field values (read them before assuming — the "with-matter" fixture may set only `branch_has_commits`, not `pr_exists`, in which case tightening this rule changes its expected outcome and the fixture's expectation must be updated, not left stale).
- **β verifies via:** run the branch-only fixture and confirm no `review` proposal; run a PR+commits fixture and confirm `review` still proposed (unchanged); confirm `run_active`'s non-gating behavior (rule 3) is unaffected — a fixture with `run_active:true` and `review_request_present:true` should still land on whichever rule matched before (verify rule order: rules 1/2 fire before rule 3 checks `run_active`, so a request+matter fixture proposes review regardless of `run_active`, exactly as `_doc_phase2` intends — confirm this ordering is preserved).

### AC4 — matter/branch observation is robust to remote-only branches

- **Invariant:** a dead run whose `cycle/{N}` exists on the remote (but not locally) is observed as having matter — never misclassified "no matter".
- **Current state (verbatim, `fetch.go` `assembleLive`, lines ~62-73):**
  ```go
  branch := "cycle/" + strconv.Itoa(issue)
  // --- local git: branch existence + commits beyond base ---
  if err := exec.CommandContext(ctx, "git", "rev-parse", "--verify", "--quiet", "refs/heads/"+branch).Run(); err == nil {
      snap.BranchExists = true
      out, err := exec.CommandContext(ctx, "git", "rev-list", "--count", "main.."+branch).Output()
      ...
  }
  ```
  This is **local-git-only** — no GitHub API call backs up or replaces this observation. Confirmed: `.github/workflows/cnos-cds-dispatch.yml`'s `actions/checkout@v4` step (line 37-39) takes no `ref:`/`fetch-depth:` override, i.e. it checks out only the default branch (main) at default depth — `refs/heads/cycle/*` is never fetched by this workflow. A scheduled re-poll of this workflow (the `cron:` triggers at lines 15-18) that lands on a dead run with a remote-only `cycle/{N}` (pushed by a *different* workflow run / sandbox that has since torn down) would see `BranchExists:false` and misclassify via the last `in-progress` rule (`propose_status_todo` — blind requeue, the #368 failure).
- **α proves via — the one real design decision (see §Friction notes):** choose (a) add a `git fetch origin refs/heads/cycle/*:refs/remotes/origin/cycle/*` (or similarly scoped) step to the workflow before the evaluator runs, and adjust `assembleLive`'s `rev-parse`/`rev-list` refs to check `refs/remotes/origin/cycle/{N}` (or check both local and remote refs), OR (b) add a GitHub-API-based branch/comparison call in `fetch.go` (e.g. `GET /repos/{repo}/branches/{branch}` + `GET /repos/{repo}/compare/main...{branch}` for commit count) alongside/replacing the local git calls. Either way, write a test/scenario proving a remote-only `cycle/{N}` (simulate: a branch that exists on origin but has no local ref) resolves to `propose_delta_recovery`, not `propose_status_todo`.
- **β verifies via:** if (a): confirm the workflow's fetch step is present and scoped correctly, and that `fetch.go`'s ref lookup matches what the fetch step actually populates (do not fetch `refs/heads/*` into local refs then check `refs/remotes/*`, or vice versa — ref-namespace mismatch would silently no-op). If (b): confirm the new API call path is exercised by a test using a mock/fake HTTP server (not just live GitHub, per AC1's `TestApply_*` pattern using `githubAPIBase` override) and that it correctly reports `BranchExists`/`CommitsBeyondBase` for a branch with no local ref.

### AC5 — wave-closure language corrected + Phase 3 filed (deferred)

- **Invariant:** no shipped surface claims "FSM owns status labels"; the true boundary is stated; Phase 3 is tracked, not silently implied done.
- **γ pre-checked:** `rg -i "FSM owns" src/ docs/ .cdd/ .github/` returns **no hits** — the doctrine files (`cds-dispatch/SKILL.md`, `dispatch-protocol/SKILL.md`, `delta/SKILL.md`) already carry the corrected "requests via FSM apply" / "FSM applies the label only when its guards pass" language from #569's own landing (confirmed: cds-dispatch/SKILL.md lines 207/217/234/243/248; dispatch-protocol/SKILL.md line 447). **The overclaim lives only in the #567 issue's own closing comment** (`gh issue view 567 --json comments`), whose bolded headline reads `## ✅ FSM wave complete — \`Workers produce matter. The FSM owns status labels.\`` — this is the surface α must correct. See §Friction notes for the edit-permission caveat.
- **α proves via:** correct or annotate the #567 comment (edit if author identity permits; otherwise post a dated follow-up comment stating the true boundary: "FSM gates the in-progress→review transition; claim/hard-block/release-back-to-queue remain direct wake writes"), and file the Phase 3 tracking issue (label `deferred`/`held` or repo equivalent — check `gh label list` for the actual deferred-label name used elsewhere in this repo before inventing one), linked from #567 (`gh issue comment 567` or a `Related to #567` reference in the new issue body).
- **β verifies via:** `gh issue view 567 --json comments` shows the correction; `gh issue list --search "Phase 3"` (or similar) confirms the new tracking issue exists, is linked, and carries a deferred/held label; re-run the `rg -i "FSM owns"` sweep to confirm still zero hits in shipped files.

### AC6 — hidden/bidi Unicode sweep of wave-touched files

- **Invariant:** no unexpected bidirectional/hidden Unicode control chars in the files the wave (#571/#572/#573) touched.
- **File set (γ-gathered via `gh pr view {571,572,573} --json files`):** the union is listed in full in §Surfaces above — reuse that list verbatim rather than re-deriving it from the PRs.
- **α proves via:** a byte-level sweep, e.g. `rg -nP '[\x{200B}-\x{200F}\x{2028}-\x{202E}\x{2060}-\x{2069}\x{FEFF}]' <file-list>` (or `LC_ALL=C grep -P` byte-range equivalent) over the file set above, specifically including `#572`'s files first (the PR GitHub flagged with a hidden/bidi warning) — record the result (clean, or matches found + explained) in `self-coherence.md`.
- **β verifies via:** independently re-run the same sweep command over the same file list and confirm the recorded result matches.

### AC7 — current gates remain green

- **Invariant:** the hardening does not regress repo health.
- **Named CI jobs (γ-confirmed present in `.github/workflows/build.yml` + `.github/workflows/install-wake-golden.yml`):** `package-source-drift` (I1), `protocol-contract-check` (I2), `link-check` (I4), `skill-frontmatter-check` (I5), `cdd-artifact-check` (I6), `dispatch-repair-preflight`, `dispatch-closeout-integrity`, `go` (Go build & test), `package-verify`, `binary-verify`, plus the `install-wake-golden` workflow's job(s). `TestSeam_CellKindNotEnforced` (confirmed present at `issuesfsm_test.go:678`) must remain unmodified and passing — α must NOT touch this test.
- **α proves via:** `go build ./... && go test ./... -race` green locally before requesting review; all named CI jobs green on the cycle branch's latest commit.
- **β verifies via:** `gh pr checks <PR>` (or `gh run list --branch cycle/574`) confirms all named jobs green; diff `issuesfsm_test.go` to confirm `TestSeam_CellKindNotEnforced`'s body is byte-identical to `origin/main`.

## Source-of-truth table

| Claim / surface | Canonical source | Status (γ-verified) |
|---|---|---|
| CDS transition table (guards to fix) | `src/packages/cnos.cds/skills/cds/fsm/transitions.json` | Current as of `origin/main@452191f`; exact `review`/`in-progress` rule text quoted in full above |
| Evaluator engine + guard registry | `src/packages/cnos.issues/commands/issues-fsm/table.go` | Current; `guardFuncs` registry already has every guard AC2/AC3 need (`review_request_present`, `pr_exists`, `pr_has_commits`, `branch_has_commits`) — no new guard function required |
| Live branch/PR/run observation | `src/packages/cnos.issues/commands/issues-fsm/fetch.go` (`assembleLive`) | Current; confirmed local-git-only for branch existence/commit count (lines ~62-73); PR/run observation already uses GitHub REST (`ghGetJSON`) |
| FSM workflow (branch fetch surface) | `.github/workflows/cnos-cds-dispatch.yml` | Current; confirmed `actions/checkout@v4` (lines 37-39) has no `ref:`/`fetch-depth:` override — `refs/heads/cycle/*` is not fetched |
| Closeout-integrity preflight (PR-required rule to align to) | `cnos.cds/orchestrators/cds-dispatch/SKILL.md` §"Closeout integrity preflight"; `cnos.core/skills/agent/dispatch-protocol/SKILL.md` §2.9 | Current; both already describe the FSM-request mechanism accurately (post-#569) |
| Wave-closure language to correct | #567 wave-completion comment | **Overclaim confirmed present** in the comment's bolded headline; doctrine files themselves are already clean (`rg -i "FSM owns"` → 0 hits in `src/`, `docs/`, `.cdd/`, `.github/`) |
| Authenticated wave issue state | #567/#568/#569/#570 | γ-verified CLOSED/COMPLETED, clean labels, 2026-07-04 (see AC1 for exact labels) |
| Cell-kind doctrine (must stay observed-only) | `cnos.cdd/skills/cdd/CELL-KINDS.md`; `TestSeam_CellKindNotEnforced` (`issuesfsm_test.go:678`) | Current; test confirmed present, must remain unmodified |
| Wave-touched file set (AC6 sweep target) | `gh pr view {571,572,573} --json files` | γ-gathered, listed in full in §Surfaces |
| Dispatch wire-format / 7-axis contract | `src/packages/cnos.handoff/skills/handoff/dispatch/SKILL.md` §2.3 | Current; canonical home for the implementation-contract block below |

## Scope guardrails (hard constraints)

1. **Do not implement Phase 3** (claim/hard-block/release-back-to-queue FSM ownership) — file the tracking issue only (AC5). Do not add FSM-routed logic for those three events in this cycle.
2. **`cell_kind` stays observed-only** — `TestSeam_CellKindNotEnforced` must remain byte-identical and passing; do not add any rule in `transitions.json` (or Go logic) that branches on `FactSnapshot.CellKind`.
3. **No new status labels / taxonomy change** — AC2/AC3/AC4 tighten existing guards and observation; they do not introduce new `status:*`/`protocol:*` labels.
4. **No external controller service; no wake-source-model change; no #216 package-command-discovery change; no Demo 0** (issue non-goals) — this stays a repo-local reconciler + declarative table + (at most) one workflow fetch step.
5. **Preserve the deliberate `run_active` non-gating** on the worker's own synchronous `in-progress → review` request (rules 1/2 of the `in-progress` state fire before rule 3's `run_active` check, regardless of tightening) — do not regress #569's main use case (a worker requesting its own transition while its run is still `in_progress`).
6. **Do not re-derive Phase-1/Phase-2 evaluator *logic*** (`Evaluate`/`ruleMatches`/`guardFuncs` in `table.go`) — AC2/AC3 are data-only edits to `transitions.json`; no new guard function is needed (verified: `pr_exists`, `pr_has_commits`, `review_request_present`, `branch_has_commits` all already exist in `guardFuncs`).
7. **STOP-and-ask trigger** (mirrors #569's operator note, still binding): if AC4's remote-branch observation requires broader GitHub API/token authority than already used elsewhere in this repo's wakes (`contents.write`/`issues.write`/existing read scopes), STOP and surface to operator before broadening scope.

## Friction notes

- **The one real design decision (issue's own cycle-scope-sizing factor (d)): AC4's fetch-vs-API choice.** The issue explicitly defers this to α: "either (a) the FSM workflow fetches `refs/heads/cycle/*` before invoking the evaluator, or (b) `fetch.go` observes remote branch/PR via GitHub API." Trade-offs α should weigh: (a) is a smaller diff (one workflow step + a ref-path tweak in two `exec.Command` calls) but only fixes the CI-workflow invocation path — any *other* caller of `assembleLive` (e.g. a local operator run, or a future non-workflow invocation) would still be local-git-only and could still misclassify a remote-only branch; (b) is a larger diff (new API calls in `fetch.go`, mirroring the existing `ghGetJSON`/PR-observation pattern already in the same function) but fixes the observation primitive itself regardless of caller, and is more consistent with how PR/run observation already works in this file (already API-based). γ's recommendation (non-binding): prefer (b) for correctness-regardless-of-caller, since `fetch.go` already has the exact idiom to extend (a `GET /repos/{repo}/branches/{branch}` call — 200 vs 404 gives `BranchExists`; a `GET /repos/{repo}/compare/main...{branch}` call's `.ahead_by` gives `CommitsBeyondBase`) and the workflow-fetch approach (a) leaves a latent gap for any future non-CI caller. α owns the final choice; if α instead picks (a), or a hybrid, document the reasoning in `self-coherence.md`. Per the issue's own "Cycle scope sizing" §Decision: if either option balloons the diff, split AC4 into a sub-issue and land AC2/AC3/AC5/AC6/AC7 first — this is the named split trigger, not a suggestion to avoid AC4 entirely.
- **AC5's overclaim location is a GitHub comment, not a repo file** — `rg`-based oracles (`rg -i "FSM owns (all )?status labels"`) will find nothing in the repo itself since the doctrine files are already clean; α must not conclude AC5 is already satisfied from a repo-only grep. The actual target is #567's closing comment. α's own GitHub identity may not have edit rights on a comment authored by a different identity (`usurobor`/prior κ session) — verify via `gh api` whether edit is possible; if not, the remediation is an appended correction comment, not an edit, and that distinction should be recorded rather than silently downgraded.
- **AC2/AC3 rule-tightening is not a pure `any_true`→`all_true` mechanical swap for the blocked/negative rule.** The current blocked rules use `all_false` over the same guard set as the valid/proposed rule's `any_true` — these are exact logical complements today. Once the positive rule becomes `all_true: [g1, g2, g3]`, its complement is "at least one of g1/g2/g3 is false", which is `NOT all_true`, not `all_false` (all three false is only one of the several ways the tightened rule can fail). If α leaves the blocked rule as `all_false`, a partial-evidence state (e.g. `review_request_present:true, pr_exists:true, pr_has_commits:false`) would match neither rule and fall through to "no rule matched" (`Evaluate`'s fallback: `outcome:"blocked", reason: "no rule in the transition table matched state..."`) — functionally still blocked, but with a generic/wrong `reason` string and no `MissingEvidence` list (since `EvidenceGuards` is only read for a rule that explicitly matched with `Outcome=="blocked"`). α should decide explicitly whether to rely on the engine's structural fallback (works, but produces a less informative message) or write an explicit widened blocked rule (better diagnostics) — document the choice in `self-coherence.md` rather than leave it as an accidental side effect.
- **Existing #569 in-progress fixtures may need their expected outcomes updated, not just re-verified.** `in-progress-review-request-with-matter.json` (shipped by #573) is presumed to satisfy the *old* `any_true` guard; α must read its actual field values before assuming it still resolves to `proposed`/`review` under the *tightened* `all_true: [review_request_present, pr_exists, pr_has_commits]` rule — if that fixture sets only `branch_has_commits:true` (not `pr_exists`), tightening the guard changes its correct expected outcome from `proposed` to `blocked`, and the existing test asserting `proposed` for it would need updating (this is an intentional behavior change per AC3, not a regression to avoid).

## Implementation contract (pinned by δ; α MUST NOT improvise)

| Axis | Pinned value |
|---|---|
| Language | Go (repo-native; `transitions.json` edits are JSON data, not Go, per the table's designed extension point) |
| CLI integration target | `cn issues fsm evaluate` / `cn issues fsm evaluate --apply` — no new subcommand, no flag rename; behavior changes come entirely from data (`transitions.json`) and (if AC4 picks option (b)) `fetch.go`'s observation internals, not the CLI surface |
| Package scoping | `src/packages/cnos.cds/skills/cds/fsm/transitions.json` (CDS-owned declarative data, AC2/AC3); `src/packages/cnos.issues/commands/issues-fsm/` for Go source (`fetch.go`, `issuesfsm_test.go`, `testdata/*.json` — AC4/test coverage); `.github/workflows/cnos-cds-dispatch.yml` only if AC4 picks option (a) or a hybrid; doctrine prose (`cds-dispatch/SKILL.md`, `dispatch-protocol/SKILL.md`, `delta/SKILL.md`) only if α's re-grep finds a surviving overclaim (γ's own re-grep found none — see §Source-of-truth) |
| Existing-binary disposition | Preserve — `cn issues fsm evaluate [--apply]` remains the single compiled-in kernel command; no new subcommand, no forked binary |
| Runtime dependencies | None new — if AC4 picks option (b), reuse `net/http` (stdlib) + the existing `ghGetJSON`/`ghRequest` dependency-free pattern already in `fetch.go`; no third-party GitHub client |
| JSON/wire contract preservation | Preserve `FactSnapshot`/fixture JSON shape as-is (existing fixtures must remain loadable unchanged in field names/types); `Decision`'s printed/rendered output field set is unchanged by this cycle (AC2/AC3 change which rule matches, not the `Decision` struct shape) |
| Backward-compat invariant | Every existing fixture whose current expected outcome is untouched by the invariant tightening (i.e. any fixture that doesn't specifically probe the review-request-alone / branch-only-no-PR cases) must continue to evaluate identically; `TestSeam_CellKindNotEnforced` must remain byte-identical; `cn issues fsm evaluate` (no `--apply`) remains exactly as read-only/side-effect-free as it is today — this cycle changes *which* rule matches, never the read-only/mutate-only-on-`--apply`-and-passing-guards invariant itself |

## α prompt

```text
You are α. Project: cnos (usurobor/cnos).
Load src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md and follow its load order.
Issue: gh issue view 574 --repo usurobor/cnos --json title,body,state,comments
Branch: cycle/574

Read .cdd/unreleased/574/gamma-scaffold.md on this branch FIRST — it names the
exact current (loosened) guard strings in transitions.json for AC2 and AC3,
the exact fetch.go block AC4 targets, the full wave-touched file list for
AC6's bidi sweep, the four authenticated wave-issue states for AC1, and the
one real open design decision (AC4: fetch-vs-API — γ's non-binding
recommendation is option (b), GitHub-API observation in fetch.go, but you own
the final call; document your choice and reasoning in self-coherence.md).

Work AC2 and AC3 in TDD order per the issue's explicit instruction: write the
failing fixture FIRST, confirm it fails against the current (unmodified)
transitions.json, THEN tighten the guard until it passes. Do not tighten the
guard before the fixture exists and is observed failing.

Read the Friction notes section carefully before touching transitions.json:
the blocked-rule for AC2/AC3 is not a mechanical any_true-to-all_true /
all_false swap — a naive rewrite can let a partial-evidence state fall through
to the engine's generic "no rule matched" fallback instead of an explicit,
diagnostic blocked rule. Decide explicitly and document the choice.

Re-grep before editing any doctrine prose for AC5 (rg -i "FSM owns" src/
docs/ .cdd/ .github/) — γ's own re-grep found zero hits in repo files; the
actual overclaim lives in GitHub issue #567's closing comment, not in a
tracked file. Check whether your identity can edit that comment before
assuming you can; if not, post an appended correction comment instead, and
file the Phase 3 tracking issue linked from #567.

Do NOT: implement Phase 3 itself; add any cell_kind-based branching to
transitions.json or Go logic (TestSeam_CellKindNotEnforced must stay
byte-identical); modify existing fixtures' expected outcomes without
explicit justification tied to the invariant change (see Friction notes on
in-progress-review-request-with-matter.json); introduce new status labels.

## Implementation contract (pinned by δ; α MUST NOT improvise)

| Axis | Pinned value |
|---|---|
| Language | Go + JSON data (transitions.json) |
| CLI integration target | `cn issues fsm evaluate [--apply]` — no new subcommand, no flag rename |
| Package scoping | `src/packages/cnos.cds/skills/cds/fsm/transitions.json` (data); `src/packages/cnos.issues/commands/issues-fsm/` (Go: fetch.go, issuesfsm_test.go, testdata/*.json); `.github/workflows/cnos-cds-dispatch.yml` only if AC4 picks option (a)/hybrid |
| Existing-binary disposition | Preserve — single compiled-in `cn issues fsm evaluate [--apply]` |
| Runtime dependencies | None new — stdlib net/http only, mirroring fetch.go's existing ghGetJSON/ghRequest pattern (if AC4 picks option (b)) |
| JSON/wire contract preservation | Preserve FactSnapshot/fixture shape and Decision struct fields as-is |
| Backward-compat invariant | Untouched fixtures evaluate identically; TestSeam_CellKindNotEnforced byte-identical; evaluate (no --apply) stays read-only |
```

## β prompt

```text
You are β. Project: cnos (usurobor/cnos).
Load src/packages/cnos.cdd/skills/cdd/beta/SKILL.md and follow its load order.
Issue: gh issue view 574 --repo usurobor/cnos --json title,body,state,comments
Branch: cycle/574

Read .cdd/unreleased/574/gamma-scaffold.md on this branch FIRST — walk the
per-AC oracle list (AC1-AC7) independently against α's diff and
self-coherence.md. In particular verify:
  - AC1: the four authenticated states recorded in self-coherence.md match a
    fresh `gh issue view {567,568,569,570} --json state,stateReason,labels`,
    and any public-render discrepancy claim (or explicit absence) is present;
  - AC2: a fixture with review_request_present:true and no PR/commits now
    evaluates blocked (re-run it against the tightened transitions.json); a
    partial-evidence fixture (pr_exists:true, pr_has_commits:false,
    review_request_present:true) also blocks — check whether it matched an
    explicit blocked rule (informative reason + MissingEvidence) or fell
    through to the engine's generic "no rule matched" fallback, and confirm
    that was a documented, deliberate choice in self-coherence.md, not an
    accident;
  - AC3: a branch-only (no PR) fixture no longer proposes review; a
    PR+commits fixture still proposes review; the run_active non-gating on
    rules 1/2 is unchanged (a request+matter fixture proposes review
    regardless of run_active); check whether any #573-era fixture's expected
    outcome changed and whether that change is justified against the new
    invariant rather than silently accepted;
  - AC4: the fetch-vs-API design decision is explicit and reasoned in
    self-coherence.md; a remote-only-branch scenario/test demonstrates
    propose_delta_recovery, not propose_status_todo; if option (a) was
    chosen, confirm the workflow's fetch step and fetch.go's ref lookup use
    consistent ref namespaces (no local/remote ref mismatch);
  - AC5: `rg -i "FSM owns" src/ docs/ .cdd/ .github/` still returns zero
    hits; the #567 comment overclaim is corrected or has an appended
    correction; a Phase 3 tracking issue exists, is labeled deferred/held,
    and is linked from #567;
  - AC6: the bidi/hidden-Unicode sweep covers the full wave-touched file set
    named in gamma-scaffold.md §Surfaces (union of #571/#572/#573 files),
    with #572's files checked specifically, and the result is recorded;
  - AC7: all named CI jobs (I1/I2/I4/I5/I6, install-wake-golden,
    dispatch-repair-preflight, dispatch-closeout-integrity, go,
    package-verify, binary-verify) are green on the cycle branch's latest
    commit, and TestSeam_CellKindNotEnforced is byte-identical to
    origin/main (diff it explicitly, don't just check it passes);
  - the scope guardrails are honored: no cell_kind-based enforcement
    branching added; no new status labels; Phase 3 itself not implemented;
  - the 7-axis implementation contract in gamma-scaffold.md is satisfied by
    the diff (Rule 7).
```
