# γ close-out — cnos#628

**Issue:** #628 — S1 (doctrine): supersede cell-kind ontology with the CCNF-kernel + WC/PC/CC
deployment classes; settle CM↔V (sub of #627).
**cell_kind:** `doctrine`
**run_class:** `repair_pass` (R1). R0 was `first_pass`; R1 is a genuine `repair_pass` re-entry per
`cds-dispatch/SKILL.md` §"Repair re-entry preflight" Step A check 3 — R0's `beta-review.md`
carried a non-converge (`iterate`) verdict and the R0 `status:blocked` comment is a β bounce
naming the rejected finding.
**Shape:** this is **not** a standard γ-scaffolds → α-implements → β-reviews cell (see
`gamma-scaffold.md` "Shape of this cell"). The α-matter under review is externally authored (PR
#629, κ/Sigma-as-author) precisely because κ cannot self-review without violating the α≠β
firebreak. This cell's own matter is the independent β-review receipt set landed under
`.cdd/unreleased/628/` on `cycle/628`, distinct from PR #629's own branch.
**Rounds:** R0 (`verdict: iterate`) → R1 (`verdict: converge`). See `beta-review.md` and
`beta-closeout.md` for the full round history and repair-evidence block.
**Scope of this artifact:** PR-time close-out triage at the δ-contract convergence boundary
(`delta/SKILL.md §9.5`), per `gamma/SKILL.md §2.7`. Not the release-time closure gate
(`gamma/SKILL.md §2.10`) — no `RELEASE.md`, no move to `.cdd/releases/`, no tag, no PRA, no
issue-close assertion here.

## Cycle summary

R0 ran the independent β review this cell exists to produce, against PR #629's diff and AC1–AC8.
It found the substantive doctrine content sound but a dispositive blocker: the PR branch was not
actually rebased onto current `main` and silently deleted landed cnos#639 content undisclosed in
the PR's own summary. This cell moved to `status:blocked` (external dependency) since it has no
in-band α to re-dispatch and cannot repair PR #629 itself without collapsing the independence it
exists to preserve. κ rebased PR #629 externally (new head `68797cf9`) and re-applied
`status:todo`. R1 re-claimed under `run_class: repair_pass`, wrote `REPAIR-PLAN.md` before
touching any other artifact (per Step C), independently re-verified the rebase (not trusting
κ's comment — re-ran the ancestor check and the three-dot diff, confirmed all previously-deleted
content restored), spot-checked that AC1–AC8 substance was unchanged from R0's findings, and
converged.

## Close-out triage (CAP)

| Finding | Source | Type | Disposition | Artifact / commit |
|---|---|---|---|---|
| PR #629 stale base silently deleting cnos#639 content, undisclosed in the PR's own Changes summary | β, R0 (`beta-review.md §R0`) | correctness (dispositive) | **Resolved — repaired, independently re-verified, not overridden.** κ's rebase (68797cf9) restores all deleted content; R1 confirmed via `git merge-base --is-ancestor` + three-dot diff, not by trusting the PR body or κ's comment. See `beta-closeout.md`'s `repair_evidence` block. | `beta-review.md §R1`, `REPAIR-PLAN.md` |
| Same-`sigma`-account α/β firebreak gap: GitHub refused a formal "request changes" review on PR #629 since this wake and the PR author share one account | β, R0 | process / substrate | **One-off for this cell — flag for operator, do not fix here.** Non-goals explicitly exclude wakes/CDS-behavior changes from this cell's scope. Named as a candidate follow-up (separate bot identity for dispatch-wake execution) in both the R0 `status:blocked` comment and κ's rebase comment; not filed as a new issue by this cell — left for operator discretion. | R0 issue comment (2026-07-09T08:16:10Z); κ comment (2026-07-10T08:22:17Z) |
| `human_gate` / `waiting_human_gate` enum-naming overlap in `CELL-RUNTIME.md` | β, R0, carried to R1 | judgment / tidiness | **One-off — drop.** Non-blocking on AC6 at both rounds; a candidate follow-up tidy, not a gap in this cell's own scope. | `docs/architecture/CELL-RUNTIME.md` (on PR #629, not this cell's own diff) |
| AC7/AC8 verified only as far as diff-checkable / PR-body-text, not independently re-derivable beyond the PR's own claims | β, R0, carried to R1 | verification-depth limit | **One-off — drop.** Both non-blocking; the issue's own oracle discipline ("grep/citation checks against the landed doctrine") was applied everywhere it was mechanically applicable; AC7/AC8's remaining component (process/reconciliation narrative) is inherently a text-presence check, not a grep-verifiable fact. | `beta-review.md §R0` AC7/AC8 rows |

No finding in this cycle rose to an **immediate MCA**, a **project MCI**, or an **agent MCI** —
the one dispositive finding was fully repaired and independently re-verified within this cycle's
own repair round; the remaining three are one-off process/judgment/scope notes.

## Cycle-iteration trigger assessment (`CDD.md` step 10 / `gamma/SKILL.md §2.8`)

| Trigger | Fire condition | Observed this cycle | Fired? |
|---|---|---|---|
| Review churn | review rounds > 2 | 2 rounds (R0 iterate → R1 converge) | **No** — at the threshold, not over it |
| Mechanical overload | mechanical ratio > 20% **and** total findings ≥ 10 | 1 dispositive + 3 non-blocking findings total, well under the 10-finding floor | **No** |
| Avoidable tooling / environment failure | environment/tooling blocked the cycle in a way a guardrail could likely prevent | The R0→R1 gap was a genuine external-dependency block (stale PR branch), not a tooling/environment failure of this cell's own dispatch execution; the claiming wake's own local shallow-clone artifact (fixed via `git fetch --unshallow` before the FSM claim transition would resolve accurate facts) was a one-time local environment quirk of this interactive session standing in for a substrate firing, not a repo-level or substrate-level defect | **No** |
| Loaded-skill miss | a loaded skill should have prevented a finding but did not | No — the `cds-dispatch/SKILL.md` "Repair re-entry preflight" Step A classification worked exactly as designed: it correctly identified this R1 round as `repair_pass` (not `resumed_from_matter`, despite pre-existing branch/PR matter superficially resembling that shape) via the rejection-evidence discriminator (existing `beta-review.md §R0` non-converge verdict), which is precisely the ambiguity that section's ordered-checks discipline exists to resolve | **No** |

No trigger fired. Per `gamma/SKILL.md §2.9`, the four process-gap questions directly:

- **Did this cycle reveal a recurring friction?** One worth naming: the FSM's `run_active` guard
  (`transitions.json` `all_false: ["run_active"]`) checks the branch's *most recent* GitHub Actions
  workflow run generically (`GET .../actions/runs?branch=...&per_page=1`), not scoped to
  `cds-dispatch` runs specifically. A routine `Build` CI run triggered by this cycle's own claim
  push transiently blocked the `todo -> in-progress` transition until that CI run completed
  (~2 minutes). Not a bug — the guard is deliberately conservative — but worth naming as a known
  latency source for any dispatch-wake firing that pushes to `cycle/{N}` before requesting its own
  claim transition.
- **Was any gate too weak or too vague?** No. The Step A ordered-checks discriminator table
  correctly disambiguated `repair_pass` from `resumed_from_matter` here, which is exactly the case
  it was hardened to handle.
- **Did a role skill fail to prevent a predictable error?** No.
- **Did coordination burden show a better mechanical path?** No. The κ ↔ dispatch-wake handoff
  (block → external rebase → re-apply `status:todo` → re-claim) worked as the lifecycle-transition
  table intends: `status:blocked` is operator-authority to un-block, and κ used exactly that path.

**Disposition: no patch, no MCA.** The `run_active`/CI-latency observation above is named as a
process note, not urgent enough for a follow-up issue on its own — the guard did its job
correctly (prevented a claim transition racing an in-flight CI run on the same branch), it simply
cost a short wait.

## Next MCA

None committed from this cycle as a hard obligation. Two candidates named for operator discretion
(neither filed as a new issue by this cell, consistent with its own non-goals):

1. A separate bot identity for dispatch-wake execution (or an equivalent mechanism), so the α≠β
   firebreak this cell's subject matter concerns is structurally — not just procedurally —
   enforced when a dispatch wake reviews α-matter authored under the same account.
2. Scoping the FSM's `run_active` guard to the `cds-dispatch` workflow specifically (or otherwise
   distinguishing it from incidental CI runs like `Build`), to avoid a wait whenever a dispatch
   wake pushes claim-marker commits before requesting its own claim transition.

## Mandatory terminal learning section

```yaml
learning:
  observations:
    - "A repair_pass re-entry with an externally-authored α-matter (this cell's own shape) still fits the Steps B-E contract cleanly: REPAIR-PLAN.md's rejected-finding -> planned-repair row just names 'independently re-verify what an external actor (kappa) already fixed' as the planned repair, rather than 'route alpha to fix it' — the contract doesn't assume alpha does the repairing, only that delta verifies against a written plan before re-certifying."
    - "The FSM's run_active guard reads the branch's single most-recent workflow run via the GitHub REST API, not filtered by workflow name -- a routine Build CI run triggered by this cycle's own claim-marker push was enough to transiently block the todo -> in-progress transition until that unrelated run completed."
    - "A shallow git clone silently corrupts commits_beyond_base / claim_request_present observation (merge-base fails to find a common ancestor within the shallow window, and CDDArtifacts is read from the local working tree's checked-out branch, not the remote cycle/{N} branch content) -- git fetch --unshallow and checking out cycle/{N} locally before running cn issues fsm evaluate are both required for accurate facts in an interactive session standing in for a substrate firing."
  process_deltas:
    - "None required in this cell's own scope (doctrine-only, no schema/wake/CDS changes per its own non-goals). Both named follow-up candidates above are left for operator discretion, not executed here."
  reusable_patterns:
    - "For a review-of-externally-authored-matter cell shape, alpha-closeout.md can stand in as an explicit pointer document (asserting nothing about correctness) rather than a true implementer self-assessment, keeping the six-artifact canonical set complete without fabricating an alpha round that didn't happen."
    - "Independently re-deriving a repair's resolution (re-running the exact verification command that produced the original rejected finding, against the new state) rather than trusting the rejecting party's or the repairing party's own account of the fix is the direct, mechanical form of 'repair, do not re-certify' (cnos#516)."
  followups:
    - issue: "(unfiled candidate) separate bot identity or equivalent mechanism for dispatch-wake execution, so the alpha!=beta firebreak is structurally enforced when a dispatch wake reviews same-account alpha-matter"
    - issue: "(unfiled candidate) scope the FSM's run_active guard to cds-dispatch workflow runs specifically, or otherwise distinguish them from incidental CI (e.g. Build) runs on the same branch"
  operator_burden:
    - "One block -> external-rebase -> re-claim round-trip (2026-07-09 block, 2026-07-10 kappa rebase + re-apply status:todo). This is the lifecycle-transition table's intended operator-authority path for status:blocked, not excess burden -- the alternative (this wake rebasing PR #629 itself) would have collapsed the independence this cell exists to preserve."
```

## Explicitly out of scope for this artifact

Per this cell's own PR-time closeout scope (see header), the following release-time steps were
**not** performed and are deferred to the actual repo-release pass: `RELEASE.md` authoring; moving
`.cdd/unreleased/628/` → `.cdd/releases/{X.Y.Z}/628/`; tagging; post-release assessment (PRA);
hub-memory update; merged-branch cleanup; parent-issue close-state assertion. Issue #628 remains
`OPEN` as of this closeout, correctly so (this cell's own PR #646 has not merged; nor has the
reviewed PR #629).

---

## Deliverable evidence (dispatch-wake closeout-integrity preflight, cnos#524)

```yaml
deliverable_evidence:
  pr: "#646 (cycle/628 -> main)"
  head_sha: "43dc6f0b8bd4179061c993a29ccd9cb6341cf3da"
  base_sha: "a706aa861897932a3d0da91558af0c7d15596a1b"
  commits_beyond_base: 5
  closeout_artifacts: [gamma-scaffold.md, self-coherence.md, beta-review.md, alpha-closeout.md, beta-closeout.md, gamma-closeout.md]
```

`head_sha` above is the branch state immediately before this file's own commit (5 of 6 closeout
artifacts present + the R0/R1 review and repair-plan commits); this file's commit is the next one
on `cycle/628`, completing the six-artifact set before requesting `status:review`. All five items
of the §Closeout integrity preflight are satisfied: PR #646 exists and references `#628`
(`Refs #628`, opened by the mechanical cell finalizer per cnos#591); `cycle/628` HEAD differs from
base by 5 commits (6 once this commit lands); the branch exists and diverges from base; all six
required `.cdd/unreleased/628/` artifacts are present (this commit completes the set); this block
names the PR number and head/base SHAs as evidence.
