# RCA: cnos#524 W4 dispatch cell set status:review with no deliverable

**Date:** 2026-06-30
**Severity:** Medium
**run_class:** manual_delta_repair (remediation by manual δ, not a dispatched cell)
**Caused by:** the cnos#524 W4 empty-review incident (this RCA).

## Summary

During the final phase (W4) of the cnos#524 wake-as-skill migration, a `cds-dispatch`
cell claimed issue #524 (`status:todo → status:in-progress`), ran ~25 minutes to a
"success" conclusion, and set `status:review` — while producing **no pull request, no
commits on any branch, no closeout artifacts, and no STOP/BLOCKED comment**. The
`status:review` label was a false-complete state, mechanically indistinguishable from a
finished cell. No W4 work landed.

This document is both the RCA and the receipt for the remediation guard (RCA action A2).

## Timeline (UTC)

| # | Time | Event |
|---|------|-------|
| 1 | 17:54 | operator directive "Go W4"; #524 flipped `status:review → status:todo` |
| 2 | 17:54 | Actions dispatch run `28464981342` starts (checkout main `90b9e812`) |
| 3 | 17:59 | cell posts claim comment → `status:in-progress` |
| 4 | 18:19 | `claude-code-action` step ends "success"; post-work write-fence green |
| 5 |  --   | #524 ends at `status:review` — but no PR, no commits, no closeout, no STOP |
| 6 |  --   | κ status check finds the empty review; RCA opened (no remediation pre-sign-off) |

## Five Whys

1. Why is W4 not delivered? → The cell's run produced no PR and pushed no commits.
2. Why no PR/commits? → The cell's branch/push/PR step did not land (hypothesis: a stale
   pre-existing `cycle/524` branch — not deleted after the #529 merge — disrupted the
   create-branch→push→open-PR flow; W1–W3 always started from no branch).
3. Why was a stale `cycle/524` present? → The W3 PR (#529) merge did not delete its head
   branch; post-merge branch cleanup was non-deterministic across phases.
4. Why did the failure surface as a silent `status:review` instead of an error? → The
   dispatch protocol let a cell set `status:review` with no mechanically-verified
   deliverable. Same class as W3's false-CONVERGE-on-red — closeout was self-asserted,
   not gated.
5. **Root cause** → Two coupled process gaps: (R1) non-deterministic post-merge branch
   lifecycle leaves a stale `cycle/<N>` that collides with the next same-issue phase; and
   (R2) no closeout-integrity gate — `status:review` was trusted without proof of a PR +
   diff, so an empty/aborted run looked identical to a finished one.

## Actions

```yaml
actions:
  - action: "A1 — delete cycle/<N> after its PR merges; if deletion fails, the next phase STOPs before dispatch and reports the stale branch"
    owner: operator/kappa
    status: in-progress    # stale cycle/524 deleted 2026-06-30; auto-delete policy pending
  - action: "A2 — closeout-integrity guard: a dispatch cell may set status:review ONLY with deliverable proof (PR + commits beyond base + branch differs + closeout artifacts + receipt names evidence); else STOP/BLOCKED. dispatch-protocol §2.9/§4.9/D12 + cds-dispatch prompt 'Closeout integrity preflight' + scripts/ci/check-dispatch-closeout-integrity.sh + build.yml job."
    owner: kappa
    status: done          # this change
  - action: "A3 — kappa treats status:review with no PR/diff as an anomaly, not a completion"
    owner: kappa
    status: done
  - action: "A4 — W4 header rule: do not freeze the stale JSON/prompt header; W4 updates the generated header to SKILL.md attribution and accepts a header-only golden diff"
    owner: operator/W4-cell
    status: pending        # applied when W4 is re-dispatched
```

## Remediation receipt — closeout-integrity guard (RCA action A2)

Extends cnos#516 (repair integrity); this adds **deliverable integrity**.

| AC | Requirement | Evidence |
|---|---|---|
| AC1 | Dispatch protocol requires deliverable proof before `status:review` | `dispatch-protocol/SKILL.md` §2.9 + §4.9 + failure mode D12 |
| AC2 | CDS dispatch prompt requires PR+diff proof before closeout | `cds-dispatch/prompt.md` + `SKILL.md` body §"Closeout integrity preflight" (verbatim-identical for W2/W3 parity); re-rendered golden + live workflow carry it |
| AC3 | A guard catches the empty-review failure mode (`status:review` + no PR + no commits) | `scripts/ci/check-dispatch-closeout-integrity.sh` — `closeout_violation()` detector + presence guard across protocol/prompt/SKILL/golden/live |
| AC4 | Synthetic negative proof | `--self-test`: `(review,no,no)`→violation, `(review,yes,no)`→violation, `(review,yes,yes)`→ok, `(in-progress,no,no)`→ok |
| AC5 | install-wake-golden remains green | re-rendered cds-dispatch golden + live byte-identical (sha256); W3 parity `render(JSON+prompt) == render(SKILL.md)` holds |
| AC6 | dispatch-repair-preflight remains green | `check-dispatch-repair-preflight.sh` unchanged + green |
| AC7 | I1/I2/I4/I5/I6/Go/Package/Binary remain green | I1 `cn build --check` all valid; I5 95 SKILL.md no findings; rest on CI |
| AC8 | Receipt states this was caused by the W4 empty-review RCA | this document |

**Boundary:** the guard changes the dispatch wake's contract (intended), so the
`cds-dispatch` golden + live workflow change by the new prompt section. No other wake
(`agent-admin`) changed; `schemas/skill.cue` unchanged; no `wake-provider.json`/`prompt.md`
deletion (that is W4).

## Anti-pattern check (per cnos.eng/rca)

- Not "human error" → the system allowed `status:review` with no deliverable; A2 designs around it.
- Not "be more careful" → a mechanical guard + self-test, not exhortation.
- Action items present and owned.
