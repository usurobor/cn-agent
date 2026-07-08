# cnos#493 — label-doctor live audit + repair record (AC1)

This is the AC1 "audit complete" artifact: a live-repo enumeration of all
8 `src/packages/cnos.core/labels.json` (schema `cn.labels.v1`) entries,
their state at cycle start, the repair actions taken, and their state
after repair. All commands below were run against the real
`usurobor/cnos` repo using the built `cn label-doctor` binary (via
`cn label doctor`, the noun-verb realization of the registered
`label-doctor` kernel command — see
`src/go/internal/cli/cmd_label_doctor.go`'s doc comment) plus one
targeted `gh api` correction (see "Residual gap" below).

## Before (matches γ scaffold's live audit table exactly)

`cn label doctor --dry-run` output, run before any repair:

```
label-doctor: auditing 8 canonical label(s) against usurobor/cnos
  drifted  status:backlog       color: ededed -> ededed; description: "" -> "Well-formed scope but not yet refined to dispatch readiness."
  drifted  status:ready         color: ededed -> c2e0c6; description: "" -> "Spec'd; ACs converged; awaiting operator authorization."
  drifted  status:todo          color: ededed -> 0e8a16; description: "" -> "Operator-authorized; dispatch wake may claim."
  drifted  status:in-progress   color: ededed -> fbca04; description: "" -> "Claimed by a dispatch wake; cycle running."
  drifted  status:review        color: ededed -> 5319e7; description: "" -> "Cell complete; awaiting external/operator review of the receipt/result."
  drifted  status:changes       color: ededed -> d93f0b; description: "" -> "External review requested changes; issue must be revised before re-dispatch."
  missing  status:blocked       color=b60205 (will be created)
  drifted  dispatch:cell        color: ededed -> 1d76db; description: "" -> "This issue is an executable implementation cell. Eligible to be claimed by a dispatch wake when paired with status:todo and a matching protocol:{id}."
(dry-run: no mutations performed)
✗ label-doctor: drift detected (8 finding(s) not canonical) — rerun without --dry-run to repair
```

7 of 8 canonical labels drifted (color mismatch, all with an empty live
description vs. canonical), 1 (`status:blocked`) missing outright —
exactly reproducing `.cdd/unreleased/493/gamma-scaffold.md`'s "Live
audit" table (§"Live audit (real findings, not the issue's narrower
claim)").

## Repair run (`cn label doctor`, no `--dry-run`)

The first repair attempt (before the two bug fixes described below)
applied `status:backlog` through `status:blocked` (7 labels, in manifest
order) and then hard-failed on `dispatch:cell`'s PATCH with:

```
✗ label-doctor: repair "dispatch:cell": github api update label "dispatch:cell": HTTP 422:
  {"message":"Validation Failed","errors":[{"resource":"Label","code":"custom","field":"description",
  "message":"description is too long (maximum is 100 characters)"}],...}
```

This surfaced two real defects in the shipped tool, both fixed on this
branch (see commit `baab9367` "fix two bugs found live against
usurobor/cnos"):

1. **Doctor's apply loop aborted entirely on the first per-label
   failure.** `dispatch:cell` happens to sort last in `labels.json`, so
   this run got lucky — but the bug meant any earlier-sorted label
   failing would have silently blocked every later one too. Fixed to
   attempt every finding independently and join per-label errors into
   one returned error (`Result.Applied` / `Result.Failed` now both
   populated). Regression test:
   `TestDoctor_Apply_ContinuesPastPerLabelFailure`.
2. **`ghCreateLabel` tolerated ANY 422 as "already exists."** GitHub
   also returns 422 for unrelated validation failures (this exact
   over-length-description case, had the label been missing rather
   than drifted) — silently swallowing a genuine failure as a false
   success. Fixed to inspect the 422 body's `errors[].code` and only
   tolerate `"already_exists"`. Regression test:
   `TestGhCreateLabel_422WithoutAlreadyExistsCode_IsHardError`.

## After

`cn label doctor --dry-run` output, run after the repair run plus one
targeted `gh api` color-only correction for `dispatch:cell` (see
"Residual gap" below):

```
label-doctor: auditing 8 canonical label(s) against usurobor/cnos
  match    status:backlog       color=ededed
  match    status:ready         color=c2e0c6
  match    status:todo          color=0e8a16
  match    status:in-progress   color=fbca04
  match    status:review        color=5319e7
  match    status:changes       color=d93f0b
  match    status:blocked       color=b60205
  drifted  dispatch:cell        color: 1d76db -> 1d76db; description: "" -> "This issue is an executable implementation cell. Eligible to be claimed by a dispatch wake when paired with status:todo and a matching protocol:{id}."
(dry-run: no mutations performed)
✗ label-doctor: drift detected (1 finding(s) not canonical) — rerun without --dry-run to repair
```

**7 of 8 canonical labels are now fully canonical** (color AND
description byte-equal to `labels.json`, case-insensitive on hex) —
`status:backlog`, `status:ready`, `status:todo`, `status:in-progress`,
`status:review` (AC3's originally-cited finding — color confirmed
`5319e7`), `status:changes`, `status:blocked` (created).

## Residual gap: `dispatch:cell`'s description (real, live-verified, out of scope to silently fix)

`labels.json`'s `dispatch:cell` entry has a 149-byte description:

> "This issue is an executable implementation cell. Eligible to be
> claimed by a dispatch wake when paired with status:todo and a
> matching protocol:{id}."

GitHub's label API caps descriptions at **100 characters**
(`docs.github.com/rest/issues/labels#update-a-label`); any PATCH or
POST carrying this value 422s with `"description is too long (maximum
is 100 characters)"`. This is a genuine defect in the canonical
manifest's own data — not a bug in `label-doctor` — discovered live,
not assumed.

`label-doctor`'s color+description PATCH is a single atomic call
(GitHub validates the whole payload before applying anything), so a
PATCH carrying the correct color alongside the too-long description
fails the color update too. To get the live repo as close to canonical
as achievable without either (a) editing `labels.json`'s content
(out of scope — the pinned JSON/wire-contract axis says the tool reads
the manifest as-is, no schema/data change in scope for this issue) or
(b) adding retry/partial-PATCH complexity to the shipped tool for what
is currently a single-label edge case, this was corrected with one
targeted, manual, color-only PATCH:

```
gh api -X PATCH repos/usurobor/cnos/labels/dispatch:cell -f color=1d76db
```

`dispatch:cell`'s color is now `1d76db` (canonical); its description
remains empty live and cannot be set to the canonical value until
`labels.json`'s `dispatch:cell.description` is shortened to ≤100
characters. This is named as known debt in `self-coherence.md` §Debt —
a small, well-scoped follow-up against `labels.json` itself (owned by
cnos.core's label doctrine, not this issue) would close it completely.

## Update (resumption pass, 2026-07-08): residual gap closed

The "residual gap" above is now resolved and this section records how,
rather than silently rewriting the account above (per `alpha/SKILL.md`
§2.3's commit-message/intra-doc-repetition rule — the fact appeared at
three sites: this file's "Before"/"After" transcripts and "Residual gap"
prose, plus `self-coherence.md`'s AC2 evidence; all three are addressed
here and in a corresponding `self-coherence.md` correction note rather
than edited in place).

Commit `0135b30b` ("fix over-100 dispatch:cell label description
(GitHub API limit)") shortened `labels.json`'s `dispatch:cell.description`
from the 149-byte string quoted above to a 92-byte string ("Executable
cell; claimable by a dispatch wake with status:todo and a matching
protocol:{id}."), and updated `label-doctrine/SKILL.md`'s doc table to
match. That commit landed on the branch but its downstream effect — the
*live* `usurobor/cnos` label was never re-repaired against the new,
now-installable description — was not yet re-verified against the real
repo.

Re-verified live during this resumption pass:

```
$ cn label doctor --repo usurobor/cnos --dry-run
...
  drifted  dispatch:cell        color: 1d76db -> 1d76db; description: "" -> "Executable cell; claimable by a dispatch wake with status:todo and a matching protocol:{id}."
✗ label-doctor: drift detected (1 finding(s) not canonical) — rerun without --dry-run to repair

$ cn label doctor --repo usurobor/cnos
... (exit 0, no error)

$ gh label list --repo usurobor/cnos --json name,color,description | jq -c '.[] | select(.name=="dispatch:cell")'
{"color":"1d76db","description":"Executable cell; claimable by a dispatch wake with status:todo and a matching protocol:{id}.","name":"dispatch:cell"}

$ cn label doctor --repo usurobor/cnos --dry-run
... (all 8 labels report "match", exit 0)
```

**All 8/8 canonical labels are now fully canonical live** (color AND
description byte-equal to `labels.json`), including `dispatch:cell`.
AC2 is fully MET, not partially MET — see the corresponding correction
note in `self-coherence.md` §ACs. The "known debt" this section
originally pointed at (`labels.json`'s `dispatch:cell` description
being un-settable via the API) is resolved, not merely worked around;
see `self-coherence.md` §Debt for the residual debt that legitimately
remains (the manual `gh api` color-only PATCH performed earlier in this
same cycle, before commit `0135b30b` landed, was itself superseded by
the full `cn label doctor` repair above — no separate cleanup needed,
the tool's own repair overwrote it).
