# REPAIR-PLAN — cycle/574

## Classification

`run_class: repair_pass` per the dispatch wake's Step A criteria (`cycle/574` exists with commits
beyond its merge-base with `origin/main`, and `.cdd/unreleased/574/` already carries
`gamma-scaffold.md` + `self-coherence.md`) — both conditions are mechanically true, so Steps B–E
of the repair re-entry preflight are mandatory before any further file write.

**This is NOT a rejection-driven repair.** Verified via `gh api repos/usurobor/cnos/issues/574/timeline`:
the issue's full label history is `todo → in-progress (01:52:25Z) → todo (03:38:42Z, release-back-to-queue)
→ in-progress (03:41:33Z, this claim)`. There is no `status:review` or `status:changes` event anywhere
in the timeline — the cell never reached β review, so nothing was ever reviewed and rejected. The
`todo` interval at 03:38:42Z is the dispatch wake's own **release-back-to-queue** event (per the
Lifecycle-transitions table) for a run that completed the substrate workflow but stranded the cell
before β — not an operator/reviewer rejection. The prior wake's continuation comment
(https://github.com/usurobor/cnos/issues/574#issuecomment-4880515100, 03:38:40Z) itself frames this
as "κ recovery of a stranded, incomplete cell," matching the timeline.

**Correction to the prior continuation comment's "remaining work" list.** That comment states AC5,
AC6, and AC1's recording remain outstanding. Direct inspection of `cycle/574`'s actual HEAD
(`897faad1`, pushed 02:16:24Z — i.e. *before* the 03:38:40Z comment was posted) shows this is stale:
`self-coherence.md` on the branch already carries full `MET` writeups for **AC1–AC7**, and the two
external artifacts AC5 claims to have produced are independently confirmed live:

- correction comment on #567: https://github.com/usurobor/cnos/issues/567#issuecomment-4880282434 (posted 02:11:06Z)
- Phase 3 tracking issue: https://github.com/usurobor/cnos/issues/575 (open, `priority/deferred`, linked from #567 via https://github.com/usurobor/cnos/issues/567#issuecomment-4880283288)

Both predate the continuation comment's claim that they remain outstanding. The continuation
comment's characterization is superseded by this observation; the row below is a completion-status
table (not a rejected-finding-to-repair map), since there is no rejection to repair against.

## Completion-status map (in place of a rejected-finding table — no rejection occurred)

| Item | Source | Status observed on `cycle/574`@`897faad1` | Evidence target |
|---|---|---|---|
| AC1 (authenticated wave-issue-state + discrepancy record) | self-coherence.md §ACs AC1 | MET | four `gh issue view` snapshots quoted inline; public-render check via `header-state` selector |
| AC2 (`review` guard tightened) | self-coherence.md §ACs AC2 | MET | TDD red→green narrated; `transitions.json` diff; test names given |
| AC3 (`in-progress→review` requires PR) | self-coherence.md §ACs AC3 | MET | TDD red→green narrated; `run_active` non-gating test added |
| AC4 (remote-branch observation) | self-coherence.md §ACs AC4 | MET | `observeRemoteBranch` in `fetch.go` + 4 hermetic tests |
| AC5 (wave-closure language + Phase 3 filed) | self-coherence.md §ACs AC5 | MET — externally reverified this session | #567 correction comment + #575 both confirmed live via `gh` (see above) |
| AC6 (bidi/hidden-Unicode sweep) | self-coherence.md §ACs AC6 | MET | `rg` sweep command + file set + zero-match result recorded |
| AC7 (gates green) | self-coherence.md §ACs AC7 | MET locally; **CI on pushed head commit not yet independently observed** (named as §Debt item 1) | local `go build`/`go test -race` output; `TestSeam_CellKindNotEnforced` byte-identity check |

## What actually remains (the real repair/continuation scope for this claim)

1. **β review** of the full α diff (AC2–AC7), independently re-derived — not yet run at all (no
   `beta-review.md` exists on the branch).
2. **Confirm CI is green on the pushed head commit** (`gh pr checks` / `gh run list --branch cycle/574`)
   — self-coherence.md §Debt item 1 names this as not independently observed within α's own session.
3. **Closeouts** (`alpha-closeout.md`, `beta-closeout.md`, `gamma-closeout.md`) — none exist yet.
4. **Open the PR** referencing `Closes #574` (not a bare `Refs`) — no PR exists yet for this branch.
5. **Request the `in-progress → review` transition** via `cn issues fsm evaluate --issue 574 --apply`
   once β converges and the closeout-integrity preflight's deliverable evidence is satisfied.

## Repair evidence (for the eventual closeout, per the wake prompt's mandatory block)

```yaml
repair_evidence:
  prior_rejection: "N/A — stranded run, not rejected. See timeline evidence above: no status:review
    or status:changes event ever occurred on this issue before this claim."
  repairs_required:
    - stranded-cell-recovery: "resume from cycle/574 HEAD (897faad1) without re-scaffolding or
      re-implementing AC1-AC7, which were already completed and are independently verifiable"
  repairs_completed:
    - stranded-cell-recovery: "REPAIR-PLAN.md (this file) verifies AC1-AC7 completion against live
      GitHub state (issue #575, #567 comments) before any further write; no re-implementation performed"
  repairs_not_completed: []
  delta_overrides: []
  new_state_differs_from_rejected: "N/A — no prior rejected state; this is the same branch continued
    from its stranding point (897faad1), not a rewritten state after rejection"
```
