# REPAIR-PLAN — cnos#628 R1 (repair_pass re-entry)

Written before any other artifact in this round is modified, per
`cds-dispatch/SKILL.md` §"Repair re-entry preflight" Step C.

## Rejected finding → planned repair

| Rejected finding | Source | Planned repair | Evidence target |
|---|---|---|---|
| PR #629's branch was not actually rebased onto current `main`: `git merge-base --is-ancestor main origin/sigma/cell-runtime-arch-note` failed, and the branch silently deleted already-landed doctrine — `delta/SKILL.md` §9.13 (cnos#639), `.cdd/unreleased/639/*` (8 files, ~1300 lines), and the four dispatch-substrate files (`cnos-cds-dispatch.yml`, `cds-dispatch/SKILL.md`, `cds-dispatch.golden.yml`, `dispatch-protocol/SKILL.md`) — none of which the PR's own "Changes" summary disclosed. | `.cdd/unreleased/628/beta-review.md` §R0 headline finding (non-converge, `verdict: iterate`); `status:blocked` bounce comment, issue #628, 2026-07-09T08:16:10Z | κ rebased the single reconcile commit onto current `main` (new head `68797cf9`, parent `44aa9f84`) outside this wake's own α role (the α-matter is externally authored per `gamma-scaffold.md` "Shape of this cell" — there is no in-band α to re-dispatch). This round's δ independently re-verifies the rebase rather than trusting κ's comment: re-run the ancestor check, re-diff `main...PR-branch` (three-dot) and confirm it shows only the 5 originally-intended doctrine files with no unrelated deletions, and confirm `delta/SKILL.md` §9.13 and all 8 `.cdd/unreleased/639/*` files and the 4 dispatch-substrate files are present and unchanged on the rebased branch. | `git merge-base --is-ancestor`; `git diff --stat origin/main...origin/sigma/cell-runtime-arch-note` (three-dot, merge-base diff — the two-dot diff is not the right comparison since `main` advanced by one unrelated board-map bot commit after the rebase); `git show <branch>:<path> \| grep` presence checks for the previously-deleted content. |
| (carried forward, not itself rejected) Per-AC1–AC8 substance — R0 found all satisfied on the merits, independent of the headline finding | `beta-review.md` §R0 per-AC table | Re-confirm the 5-file diff scope is unchanged in content from what R0 reviewed (same file set, same net `+149/-20`) so R0's per-AC findings carry over without re-deriving the doctrine from scratch (the issue body's own instruction: "It reviews that artifact; it does not re-derive the doctrine"). Spot-check the specific citations R0's per-AC table relies on (four-orthogonal-axes table, `--class` selector, `CELL-KINDS.md` re-heading banner + #614/#640 preservation, CM↔V enum-unchanged language, `CELL-RUNTIME.md` cross-links, escalation-predicate language) are still present verbatim. | `git diff --stat` file list + insertion/deletion counts match κ's rebase comment; grep hits for each per-AC citation on the rebased branch tip. |

## Scope discipline

This repair round does not re-derive the doctrine and does not modify PR #629's branch. Per the
issue body's Path-A framing, this cell's own matter remains the receipt set under
`.cdd/unreleased/628/` on `cycle/628`. If the independent re-verification above surfaces any new
deletion or content drift beyond what's mapped here, this round STOPS and returns to `iterate`
rather than converging — Step D below ("repair, do not re-certify") applies.
