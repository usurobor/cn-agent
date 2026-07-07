# γ clarification (retroactive, R1) — `docs/guides/INSTALL-CDS.md` scoping ratification

**Cycle:** #610 — `cn repo install --dispatch cds` (dispatch layer)
**Trigger:** β round-0 review Finding 3 (C, contract) — Rule 7 (implementation-contract
coherence) drift: `docs/guides/INSTALL-CDS.md` was touched by α but is not named in the
pinned "Package scoping" row of `.cdd/unreleased/610/gamma-scaffold.md` §"Implementation
contract".
**Authored by:** α, in the R1 fix round, per β's explicit instruction ("author
`.cdd/unreleased/610/gamma-clarification.md` ... ratified retroactively in R1 per β's
finding"). This is a documentation-of-record artifact filling a procedural gap that β
found — it is not a new γ-side decision made mid-round; it is the missing ratification
of a scope addition α already made and disclosed.

## What was touched

`docs/guides/INSTALL-CDS.md` (commit `a9c194f`, "cnos#610: fix stale
docs/guides/INSTALL-CDS.md dispatch-cds prose").

## Why it was in-bounds

γ's pinned "Package scoping" row (`gamma-scaffold.md` §"Implementation contract") names:

- `src/go/internal/repoinstall/repoinstall.go` (+ `_test.go`)
- `src/go/internal/cli/cmd_repo_install.go` (+ `_test.go`)
- `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` (+ its golden re-render)

`docs/guides/INSTALL-CDS.md` is not in that list. α's peer enumeration (`alpha/SKILL.md`
§2.3 — "sibling commands / providers / ops at the same layer" and "multiple renderers or
projections of the same fact") found that this guide stated, verbatim, the same
pre-#609/#610 stale claim as `cmd_repo_install.go`'s help text: "`cn repo install
--dispatch cds` ... is not implemented yet ... fails explicitly with a clear error and
writes nothing." That claim was true before #609 (renderer generalization) and #610
(this cycle) landed, and is now false — a third sibling surface (beyond
`cmd_repo_install.go` help text and `cds-dispatch/SKILL.md` prose) making the identical
now-false claim. Per `alpha/SKILL.md` §2.3, once one peer of a "multiple renderers of the
same fact" family is identified as needing a fix, every peer must be enumerated and either
fixed or explicitly exempted — leaving this one stale would have been a rule-3.13
honest-claim violation in a shipped user-facing doc (an operator following this guide
literally would believe `--dispatch cds` still fails unconditionally, which is no longer
true after this cycle).

α self-disclosed this addition at implementation time in `self-coherence.md` §CDD Trace
("`docs/guides/INSTALL-CDS.md` | α | peer doc fix | §2.3 peer enumeration (not in the
δ-pinned package-scoping list...)") — the gap β's Finding 3 identified is purely
procedural: no `gamma-clarification.md` existed to formally extend the pinned row to
cover this file, even though the fix itself was disclosed and is judged correct/necessary
by β ("The fix itself is correct and needed... the gap is procedural").

## Disposition

**Ratified.** `docs/guides/INSTALL-CDS.md` is retroactively added to the "Package scoping"
implementation-contract row for cycle #610, alongside the three surfaces γ originally
named. No revert or content change to the doc fix is required — β independently verified
its content against live CLI output and found it accurate. This file exists so the pinned
contract row and the actual diff scope agree in the artifact record, closing Rule 7's
procedural requirement for any diff outside a pinned axis.

— α (R1 fix round, on behalf of the missing γ-side ratification; see β's `beta-review.md`
§R0 Finding 3 for the full evidentiary citation)
