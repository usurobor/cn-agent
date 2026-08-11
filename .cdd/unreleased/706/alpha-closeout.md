# α closeout — cnos#706

## Summary

Implemented the consolidated final spec's Deliverables 1–7 against Final ACs 1–10: a preflight-first gate for `cn repo install --dispatch cds` (new `install-preflight` package, presence-only GitHub secrets/push-access checks, no secret values ever handled), the `SIGMA_WORKFLOW_PAT` → `CN_DISPATCH_PAT` rename across the renderer/`repoinstall.go`/CLI help/live workflow/both golden fixtures, deletion of the cosmetic default `bot_name`/`bot_id` identity (now strictly opt-in), and an in-place revision of `docs/guides/INSTALL-CDS.md` (term definitions before use, no assumed bot account, `GITHUB_TOKEN`-only claims correctly scoped to the Tier 2 `--engine` path only).

## Deviations from the scaffold's plan

None substantive. Two soft judgment calls the scaffold left to α's discretion were made and disclosed rather than escalated:

1. **Preflight is skipped under `--dry-run`.** Not named explicitly in the scaffold; reasoned from the existing dry-run contract (writes nothing regardless, never reached label-doctor either) rather than introducing a new inconsistency.
2. **Push-access check uses the single-call `GET /repos/{owner}/{repo}` `permissions.push` field**, per the scaffold's own §Friction 3 recommendation, over the two-call collaborator/permission endpoint — taken as intended, not a deviation.

The scaffold's §Friction 4 warning (two distinct sigma bindings — the cosmetic-identity table and the PAT-secret default — both needing removal, in two different call sites of the same file) was real and both were located and fixed; nothing was missed there.

## Disclosed debt

1. **No fully `err == nil` end-to-end `Run()` test for `--dispatch cds`.** Pre-existing gap, not introduced this cycle: every dispatch-cds fixture (before and after #706) stops short of a true zero-error run because label-doctor's own git-remote resolution isn't reachable from this package's env-var test seam without triggering an uncontrolled live network call. The idempotent-re-run test (`TestRun_DispatchCds_PreflightSatisfied_SecondRunByteIdentical`) and the render-progression tests stand in for the "proceeds" half of AC4.
2. **`CN_INSTALL_PREFLIGHT_API_BASE_URL` / `CN_INSTALL_PREFLIGHT_REPO` are internal test-only env vars**, not part of the public CLI contract — they exist only to give CLI-level tests an `httptest.Server` seam. A future cycle wanting real GitHub Enterprise Server support could reconsider exposing an equivalent flag; out of this cycle's scope guardrails.
3. `docs/development/design/cn-repo-install-MOCKS.md` still reflects the pre-#706 design — left untouched per the scaffold's own disposition (reference-only, not user-facing).

## Success claim

All 10 Final ACs (AC1–AC10) implemented with named test/grep evidence in `self-coherence.md` §ACs, and independently re-verified — not merely re-quoted — by β in `beta-review.md` §R0 (verdict: converge). β re-ran every cited test, re-ran the renderer itself and byte-diffed its output against both golden fixtures and the live workflow, and re-derived each AC's oracle from the issue's consolidated-spec text rather than trusting this file's self-report. Scope guardrails held (confirmed independently by β against the full diff). No fix round was required — β's verdict converged on R0.
