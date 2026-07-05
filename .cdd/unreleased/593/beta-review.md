# beta-review — cnos#593 — R0

## Verdict: **converge**

## Independent verification performed (not a re-read of α's claims)

- Re-ran `cd src/packages/cnos.issues/commands/issues-fsm && go test ./... -v` — all `TestScan_*` + the full pre-existing suite pass (no regression to the Phase 1–3 `evaluate` tests).
- Re-ran `cd src/go && go build ./... && go test ./... && go vet ./...` — clean.
- Independently re-ran the install-wake-golden.yml AC7/AC8 leak-audit greps against `cn-install-wake` — zero matches for both the dispatch-shape (`protocol:cds|cdr|cdw|dispatch:cell|status:todo`) and admin-shape (`admin_only|disallowed_surfaces|defer_path|cell_execution`) sets.
- Independently re-rendered both wakes (`cn-install-wake agent-admin`, `cn-install-wake cds-dispatch`) and confirmed: (a) `agent-admin`'s golden is byte-unchanged (this cycle doesn't touch it), (b) `cds-dispatch`'s golden and the live `.github/workflows/cnos-cds-dispatch.yml` are sha256-identical, (c) a second consecutive render of `cds-dispatch` reports "(unchanged)" — idempotence holds.
- `git diff --stat` confirms the change surface matches the scaffold's "Surfaces touched" table exactly: `scan.go`/`issuesfsm.go`/`scan_test.go`/two new fixtures (new logic + tests), `cmd_cell.go` (the hub-gate fix), the renderer + two regenerated artifacts. Zero diff to `transitions.json`, `table.go`, `decision.go`, `CellResumeCmd`, `CellReturnCmd` — confirmed via targeted `git diff` reads, not just a description.
- Read `scan.go` end-to-end line by line for AC8 (this is the one α's own self-coherence flagged as most important to verify skeptically, mirroring the #591 β prompt's instruction not to trust a clean grep alone): confirmed every mutation path funnels through either `applyStatusLabel` (imported from the same package, the identical function `evaluate --apply` calls) or the injected `RunFinalize`/`opts.RunFinalize` subprocess call (which itself has zero label-mutation code path — verified by reading `cell.Finalizer.Finalize` in `cell.go`, which α's cycle did not modify). No `gh issue edit --add-label/--remove-label` call exists anywhere in the new code.
- Ran `./cn cdd verify --unreleased --exceptions .cdd/exceptions.yml` independently — 0 failed, consistent with α's report.
- Sandboxed the hub-gate fix myself (fresh `mktemp -d` git repo, no `.cn/`, one committed CDD-artifact commit) and confirmed pre-fix behavior would have been "no hub found" (by inspecting `discoverHub()`'s `.cn`-directory-only check, unchanged by this cycle) and post-fix behavior reaches real matter-detection (fails only on the missing `origin` remote in the synthetic repo, which is the *correct* next failure for a repo with no real remote — not a symptom of the bug being un-fixed).

## AC-by-AC walk (independent, not trusting self-coherence.md as a substitute)

- **AC1**: confirmed — `cn issues fsm scan --help` exits 0 and documents the six-case behavior; `Run()`'s switch dispatches `"scan"` to `runScan`.
- **AC2**: confirmed — `TestScan_LiveRunActive_NoOp` asserts zero finalize/comment calls for `run_active: true`.
- **AC3**: confirmed — `TestScan_DeadNoMatter_RequeuesToTodo` asserts the label pair (`DELETE status:in-progress` + `POST status:todo`) and exactly one comment.
- **AC4**: confirmed — `TestScan_DeadWithMatterNoPR_FinalizesAndComments` asserts finalize is invoked and, on the "opened draft PR for" stdout marker, exactly one comment.
- **AC5**: confirmed — `TestScan_DeadWithPRNoReviewRequest_FinalizeNoOpNoComment` asserts finalize is *still* invoked (correct — it must keep confirming the checkpoint persists) but zero comments, using the "idempotent no-op" stdout marker to distinguish. This is the one case I scrutinized hardest, since it's the one most likely to accidentally spam comments on a multi-hour stranded cell across many sweep ticks — the design (parse the finalizer's own idempotent-vs-new wording) is sound and independently reproducible.
- **AC6**: confirmed — `TestScan_ReviewRequestWithMatter_AppliesReview` asserts the `review` target-state reconciliation and zero finalize calls (correct: the label transition alone is the reconciliation here, no matter to checkpoint beyond what's already there).
- **AC7**: confirmed — `TestScan_StaleReviewNoProof_BlockedNoComment` asserts zero mutation and zero comment for a blocked review.
- **AC8**: confirmed (see "Independent verification" above).
- **AC9**: confirmed — `TestScan_Idempotent` explicitly asserts the SECOND label-request count is zero and the comment count does not grow, while explicitly asserting finalize IS re-invoked (documenting that this is intentional, not an oversight).
- **AC10**: confirmed — five fixtures map 1:1 to the issue's own named strand classes (two new: `scan-died-after-pr-before-review-request.json`, `scan-died-after-review-request-before-review.json`; three reused from the pre-existing #568/#569/#574 fixture set, which is the right call — no need to duplicate fixtures that already exist and already carry their own historical context in their filenames).
- **AC11**: confirmed for every gate this sandbox can run (Go, I1, I6, install-wake-golden re-render+diff+idempotence+leak-audit, dispatch-repair-preflight, dispatch-closeout-integrity, I2, Tier-1 kata). I4/I5 tool-unavailability is accurately disclosed, not silently skipped, and the scope argument (no SKILL.md/frontmatter edits, no new markdown links beyond one already-`--offline`-exempt GitHub issue URL) holds up under my own read of the diff.

## Findings

**None severity-blocking.** Two minor observations, neither changing the verdict:

1. `scan.go`'s `RunScan` returns a non-nil aggregate error whenever ANY scanned issue reports a per-issue error (e.g., a transient GitHub API hiccup on one issue among ten). In the rendered workflow, this step has no `if: always()` and isn't the last step in the job — a single flaky issue would fail the whole "Mechanical recovery scanner" step for that firing. This is arguably *correct* (surfacing transient failures rather than swallowing them) and matches the finalizer step's own error-surfacing convention; noting it here as a design choice worth revisiting only if it causes operational noise in practice, not as a defect to fix in this cycle.
2. The renderer comment block for the new step is verbose relative to the finalizer's; this is consistent with — not worse than — the existing convention (the finalizer step's own comment block is comparably long), so not a nit worth requesting changes over.

## Design-decision spot-check

α's "no re-observation after finalize" choice (self-coherence.md, "Design choices") is correct: I re-derived the table logic and independently confirmed `propose_delta_recovery`'s guard (`all_false: [run_active], any_true: [branch_has_commits, pr_exists, pr_has_commits]`) is unaffected by `pr_exists` flipping from false to true after finalize runs — the rule still matches either way, and no other rule in the `in-progress` transition list can fire without `review_request_present` being true, which this scanner never sets. Confirmed by reading `transitions.json` directly, not by trusting the scaffold's restated claim.

## Verdict rationale

Every AC has a runnable, independently-re-run test. The one genuinely novel discovery this cycle made (the `cn cell finalize` hub-gate bug) is backed by a real, reproducible sandbox demonstration, not asserted. The blast radius of the bug fix is minimal and verified (only `CellFinalizeCmd` touched). All named CI gates this environment can run are green. **Converge.**
