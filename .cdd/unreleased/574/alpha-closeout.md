# α close-out — cycle/574

## Summary

cnos#574 ("fix(cds/fsm): harden review guards + remote-branch observation; correct wave-closure language, post-#567 remediation") shipped seven ACs against the #567 wave's FSM guard set:

- **AC1** — re-verified #567/#568/#569/#570 as CLOSED/COMPLETED with clean labels; public-render cross-check via `data-testid="header-state"` found no discrepancy.
- **AC2** — `review`-state validity tightened from `any_true:[pr_exists, branch_has_commits, review_request_present]` to `all_true:[review_request_present, pr_exists, pr_has_commits]`, with an unconditional second rule (no guard clause) replacing a naive `all_false` complement so every partial-evidence combination lands on a diagnostic `blocked` outcome with a populated `evidence_guards` list, instead of falling through to the engine's non-diagnostic generic fallback.
- **AC3** — `in-progress → review` rule 1 tightened to the same 3-guard `all_true`, dropping `branch_has_commits` as a lone qualifier; rule 2 simplified to `all_true:[review_request_present]` alone (unconditional given that gate); `run_active` non-gating on rules 1/2 preserved and directly tested.
- **AC4** — `observeRemoteBranch` added to `fetch.go` (GitHub-API branch/commit observation, option (b) over a workflow-level `git fetch` alternative), wired into `assembleLive`'s fallback path for remote-only branches; four new hermetic tests (fake-HTTP, no live network) including the cnos#368 regression pair (`outcome=="proposed"`, `action=="propose_delta_recovery"`, `TargetState != "todo"`).
- **AC5** — #567's wave-closure overclaim ("FSM owns status labels") corrected via an appended, dated comment (not an in-place edit) on the original closure comment; Phase 3 (claim/hard-block/release-back-to-queue routing) filed and deferred as #575, labeled `priority/deferred`.
- **AC6** — bidi/hidden-Unicode sweep of the full wave-touched file set (union of #571/#572/#573's changed files, reused from gamma-scaffold.md): zero matches.
- **AC7** — full local build/test/vet/race green across all four `go.work` modules; `TestSeam_CellKindNotEnforced` confirmed byte-identical to `origin/main` via brace-balanced function-body extraction.

Diff surface: `transitions.json` (2 states, rules 1–2 only), `fetch.go` (+1 function, +1 call site), `issuesfsm_test.go` (new tests + one pre-existing assertion update), 4 testdata fixtures (3 new, 1 edited). No workflow YAML touched. No `table.go`/evaluator-core change.

## What β confirmed independently

β's R0 review (`beta-review.md`) re-derived every AC from the diff itself rather than accepting α's prose:

- Re-ran the authenticated `gh issue view` calls for #567–#570 and the public-render selector check for AC1 — same result.
- Read `transitions.json`'s diff directly and mechanically traced the "unconditional second rule" design through `table.go`'s `ruleMatches`/`Evaluate` to confirm it actually produces the claimed diagnostic-blocked behavior for AC2/AC3, rather than accepting the claim as asserted.
- Read all four AC4 tests in full and independently confirmed each tests what it claims; independently agreed with the option-(b)-over-(a) design reasoning.
- Re-ran the AC5 GitHub-side checks (`gh issue view 567 --json comments`, `gh issue view 575`) and confirmed the correction comment, linking comment, and #575's labels/body live.
- Re-ran the exact AC6 `rg` sweep independently over the same file-set union — zero matches, matching α's result.
- Re-built/re-tested all four `go.work` modules independently with `-race`, and independently re-ran the `TestSeam_CellKindNotEnforced` byte-diff via a second, separately-executed extraction.
- Performed a non-destructive merge-test against `origin/main@80be1280` in a throwaway worktree: clean merge, zero conflicts.

Verdict: **APPROVED**, round 1 (R0 — first pass; the cell stranded before ever reaching β per δ's `REPAIR-PLAN.md`, so this was not a re-review of a rejected diff).

## Non-blocking findings (F1, F2)

**F1 — test-name imprecision.** `TestAssembleLive_RemoteOnlyBranchResolvesToDeltaRecovery`'s name implies an `assembleLive`-level integration test, but the test body calls `observeRemoteBranch` directly and hand-constructs a `FactSnapshot`, composing the result the way `assembleLive`'s fallback block would rather than invoking `assembleLive` itself. The function's own doc comment discloses this ("composes...as assembleLive's fallback block would"), so the drift is name-only — no behavior claim is hidden. β independently read the actual 5-line call site in `fetch.go` and confirmed it correct by direct inspection, since no test exercises it end-to-end. Filed as a follow-up polish note, not a required fix.

**F2 — `cn-cdd-verify` classifier gap, pre-existing and unrelated to this cycle's diff.** `classifyCycleType` (`ledger.go:425`) classifies a cycle as "small-change" whenever `beta-review.md` doesn't yet exist, without checking for `gamma-scaffold.md`'s presence. The small-change path applies a stricter section-presence check (`sectionPresent`'s literal `"## Gap"` line-prefix match) than the triadic path, which does not recognize this repo's own `§`-prefixed section-header convention (`## §Gap`, `## §Skills`, etc.) in `self-coherence.md`. This produced a spurious CI-red on I6 (`CDD artifact ledger validation`) for every commit between α's `self-coherence.md` push and β's first `beta-review.md` push. β reproduced this live (CI run 28693885842 and the 7 preceding it), confirmed the mechanism by adding a stub `beta-review.md` in a merge-tree copy (reclassifying the cycle to "triadic" turned the failure into a warning), and confirmed the condition self-resolved once β's own required `beta-review.md` and `beta-closeout.md` landed on the branch at commit `c6054952`. The classifier's own source is untouched by cycle/574's diff — this is a gap in `cn-cdd-verify` itself, observable only because this cycle's timeline happened to have a multi-commit window between α's and β's artifacts landing.

Neither finding is D-level; both are recorded in `beta-review.md` §Findings with disposition reasoning.

## Residual debt (from self-coherence.md §Debt, status at close-out)

1. **CI on the pushed head commit, not independently observed within α's own session — now resolved.** α's own session only ran local build/test/vet/race; the actual GitHub Actions run status on the pushed head commit was deferred to β per the pre-review gate's explicit allowance. β independently polled `gh run list --branch cycle/574` and confirmed all named jobs green (with the transient I6 exception explained by F2 above, itself confirmed green once β's own artifacts landed). Both β and δ (via the repair-preflight/claim-check work in this dispatch) have independently confirmed CI green at commit `c6054952`. Closed.
2. **AC6's file-set reuse from gamma-scaffold.md rather than independently re-derived from the three PRs (#571/#572/#573).** α reused γ's pre-gathered union verbatim, per the scaffold's own instruction, rather than re-running `gh pr view {571,572,573} --json files` from scratch. Spot-checked against a fresh `gh pr view 572 --json files` call and matched. Low-risk inherited assumption; not independently re-derived from all three PRs. Open, low-risk, no action taken.
3. **Phase 3's tracking issue (#575) is intentionally thin.** AC5 required only filing Phase 3, not fully speccing it; #575 names the three transitions in scope (claim, hard-block, release-back-to-queue) and a sketch of what a future design cell would need, deferring the fact-model design work to that future cell. Intentional, by scope guardrail 1 (Phase 3 not implemented this cycle).

## Pattern notes

- Same class as the AC2 friction-notes concern (any_true→all_true negation trap): a literal complement rewrite of a loosened guard set produces a non-diagnostic fallback for partial-evidence states unless the replacement rule is made unconditional. Two occurrences this cycle (AC2, AC3), same resolution both times.
- Cross-role artifact-timing gap: `cn-cdd-verify`'s cycle classifier reads artifact *presence* (`beta-review.md`) as a proxy for cycle *kind* (small-change vs. triadic), which is only accurate after all three roles have written to the branch — the classifier has no signal for "triadic cycle, β artifact simply not written yet." Surfaces affected: `src/packages/cnos.cdd/commands/cdd-verify` (`ledger.go:425` `classifyCycleType`, `sectionPresent`'s literal-header matching).
