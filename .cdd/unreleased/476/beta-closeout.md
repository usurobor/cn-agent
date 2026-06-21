# β closeout — cycle/476

## Context

cnos#476 (cn-wake-install renderer, v0) merged to main as PR #477 on 2026-06-21. Merge SHA: `35380b3d0d37765be9803de7175553d844622ec0` (no-ff merge of `cycle/476` head `a3e20e3c` into `fcc5cdb9` origin/main). Three β review rounds: R1 RC (F1), R2 RC (F2), R3 APPROVE.

Files added by the cycle (on main as of merge):
- `src/packages/cnos.core/commands/install-wake/cn-install-wake` (renderer; 495 lines)
- `src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml` (golden fixture; 183 lines; sha256 `a912dd97d520b8660e4a53c2f306dad7e657b0c216bbc94f31e2958fdbd01c47`)
- `.github/workflows/install-wake-golden.yml` (CI workflow; 145 lines, 9 `run:` blocks; AC6 mechanism)
- `src/packages/cnos.core/cn.package.json` (4-line addition registering the new command)
- `.cdd/unreleased/476/self-coherence.md`, `.cdd/unreleased/476/beta-review.md`, `.cdd/unreleased/476/gamma-scaffold.md` (process artifacts)

## Review summary across R1 + R2 + R3

| Round | Verdict | Finding | Class | Resolution |
|---|---|---|---|---|
| R1 | REQUEST CHANGES | F1: AC2 negative-case CI step (line ~95 at R1 head) used `… 2>&1 \| tee /tmp/neg.log` without `set -o pipefail`, so the renderer's exit code was masked by `tee`'s exit 0; the `if` guard then took the success branch and printed an `::error::` for "accepted malformed manifest" on every push — constant-failure CI on the very workflow α introduced for AC6 | CI-correctness; wiring-claim defect (claimed proof did not fire) | α R2: added `set -o pipefail` to the step (`12f13045` implementation, `4913228f` readiness signal) |
| R2 | REQUEST CHANGES | F2: AC8 renderer-side authority audit step (lines ~115-125 at R2 head) ran `n=$(grep -ciE '…' cn-install-wake)` under `bash -e`; `grep -c` returns POSIX exit 1 when the count is zero (and zero is precisely the intended-success path here), so the command substitution killed the step before the `if [ "$n" != "0" ]` guard fired. Hidden at R1 by F1's earlier-step crash; unmasked once R2's pipefail fix let execution reach line 115 | CI-correctness; sibling of F1 in the same workflow (defect family: `bash -e` + step-internal exit-on-substitution-failure) | α R3: appended `\|\| true` to the grep substitution (`1224b5328ce13996d0a1ac499c6418e06d404733` implementation, `c55f706171e657e2d15e06f10c4c8b15fb6fda38` readiness signal); also performed a comprehensive bash-e audit of all 9 `run:` blocks in the workflow, documented in `self-coherence.md` §R3 fix |
| R3 | **APPROVE** | none | — | merged at `35380b3d0d37765be9803de7175553d844622ec0` |

## Implementation assessment

The renderer itself (`cn-install-wake`) was clean from R1. Across all three rounds, β found **zero** defects in the renderer logic, the manifest schema enforcement, the prompt-template inlining, the substrate (`anthropics/claude-code-action@v1`) shape, the authority split (manifest declares; renderer doesn't decide), the determinism (idempotence sha256-stable at `a912dd97…` across both R1 worktree merge-test and all three rounds of CI), or the golden fixture byte-equality. AC1, AC3, AC4, AC5, AC7, AC8 (substantive) all passed at R1; AC2 (behavior) passed at R1; AC6's golden + idempotence subclause passed at R1. α executed the design-and-build cleanly on the substantive surface.

The two findings (F1, F2) were both **CI-mechanism defects** in the workflow α authored to discharge AC6's "CI mechanism in place" subclause. Neither was a renderer defect. Both were shell-script CI-step correctness bugs of the same defect family: `bash -e` + step-internal exit-code propagation through pipelines (F1) or command substitutions (F2). Both required a one-line YAML fix.

The cycle's diff at merge is consistent with a clean R1 substantive implementation + two rounds of CI-mechanism iteration. The non-cycle-introduced ACs (1, 3, 4, 5, 7, 8 substantive) carried zero churn across rounds; only AC2 (CI proof on F1 step) and AC6 (CI proof on AC8 audit step, F2) churned, and only in the workflow file.

## Technical review

### What worked

- **Renderer authorship discipline**: α's R1 implementation was substantively correct on every renderer AC. The authority split (renderer doesn't encode role-decisions; AC8 invariant) was preserved both in source and by the renderer-side grep audit (zero leak strings in the renderer). Determinism was verifiable locally via two consecutive renders producing byte-identical output.
- **Golden-fixture mechanism**: AC6's "golden file + CI re-render + byte-diff" pattern is the correct invariant shape for this surface — it catches renderer drift, manifest drift, prompt-template drift, and non-determinism with a single CI assertion. The pattern is reusable for the next wake-provider rendered (cycle that adds Sub 4).
- **R2 unmasking**: F1's fix (R2) directly enabled F2's discovery (R3). The two defects are a sibling pair in the same defect family; R1 surfaced F1 by triggering the earlier step; R2's fix let execution reach F2's step. Without F1, F2 would have hidden until the next wake-provider was added.

### What recurred

The same class of `bash -e`-semantics CI-correctness defect recurred between F1 (R1→R2) and F2 (R2→R3) **in the same workflow file**, despite β's R1 review notes naming "per-CI-step execution evidence" as the next-level sharpening for cnos#472. Both findings were preventable by a per-step bash-e simulation discipline. The discipline was named in prose at R1 (β's "Notes for γ closeout") and again at R2 (β's "Note on cnos#472 effectiveness" + α R2's "Friction note for cnos#472 sharpening"), but R2 still shipped F2 because the prose discipline was not mechanically enforced as a scaffold-section α had to populate.

α R3 produced the full 9-row bash-e audit table in `self-coherence.md` §R3 fix — the shape γ should lift verbatim into the cnos#472 scaffold template — but only after F2 was found. The discipline arrived after the defect.

### Process observations

The R1+R2+R3 trajectory is the **cycle's signature finding**. Each round, the next-level sharpening was named in prose by both β (review notes) and α (self-coherence friction notes). Each subsequent round still shipped an instance of the very class the prose had named. This is the empirical case for "prose discipline insufficient → mechanical scaffold injection required" — the discipline must live as a scaffold subsection γ must populate per cycle, not as prose guidance α must remember to apply.

## Final cnos#472 sharpening recommendation (for γ closeout / PRA)

Both β R3 and α R3 converge on the same recommendation: γ closeout for 476 should file a cnos#472 follow-up that does **both**:

1. **Amend the γ scaffold template** to inject (as a populated subsection, not prose guidance) a 3-column per-CI-step table for any cycle touching CI workflows:
   - (i) per-`run:`-block `bash -e` substitution-failure-mode audit (command substitutions / pipelines / commands that could exit non-zero, guard mechanism `\|\| true` / `set -o pipefail` / `if !`, empirically-observed `bash -e` exit on intended-success input)
   - (ii) per-step CI execution evidence (job URL + conclusion + which assertion the step proves) populated once the cycle's PR CI has run
   - (iii) per-step assertion-fires verification (the step actually exercises its assertion on at least one observed input)
   The table format is α's R3 audit table in `self-coherence.md` §R3 fix, byte-liftable.

2. **Amend the α SKILL** to require α populate the table before signaling review-readiness (R[N]) on any cycle touching CI workflows, with the `bash -e` simulations as authorial evidence (not just claim).

The friction is "prose discipline insufficient → mechanical template injection required." The empirical conclusion across cycle/476's three rounds is that the discipline must be enforced by the act of populating a scaffold section, not the act of remembering prose guidance.

## Cycle health summary

- **Substantive ACs (1, 3, 4, 5, 7, 8 substantive)**: clean from R1.
- **Behavioral ACs (2)**: clean from R1; CI proof restored R2.
- **CI-mechanism AC (6)**: clean substantively from R1; CI mechanism restored R2 (F1) then R3 (F2).
- **Authority split**: preserved across all rounds.
- **Idempotence**: sha256-stable at `a912dd97…` across R1 worktree merge-test, R2 CI, R3 CI.
- **Scope discipline**: R3 only touched the CI workflow + self-coherence; R2 only touched the CI workflow + self-coherence; no out-of-scope churn at any round.
- **Pre-existing CI red builds (I4, I5, I6)**: inherited from origin/main, untouched by this cycle's diff, do not block merge per β SKILL rule 3.10.

## Closure

Merged at `35380b3d0d37765be9803de7175553d844622ec0` on 2026-06-21. β exits.
