---
cycle: 339
issue: "#339"
role: alpha
provisional: true
---

# Alpha Close-out — Cycle #339 [provisional — pending β outcome]

This close-out is written at review-readiness time per `alpha/SKILL.md` §2.8 provisional close-out fallback. It will be updated or confirmed after β merge.

---

## Cycle Summary

Cycle #339 addresses two `cdd-*-gap` findings from cycle #335:

- **F2** (cdd-tooling-gap): mechanical pre-merge closure-gate enforcement — the `gamma/SKILL.md` §2.10 closure gate was prose-only; no script or CI check prevented merging a cycle branch without close-out artifacts.
- **F7** (cdd-protocol-gap): §3.8 rubric did not specify how a closure-gate-missing cycle should be graded; the rubric's `<C` label was also not reconciled with the `C−` label used in prose artifacts.

Both are closed in one docs-only cycle (§2.5b disconnect). Deliverables:

1. `scripts/validate-release-gate.sh` extended with `--mode pre-merge` (AC1) — commit `ed982f6b`.
2. `release/SKILL.md` §3.8 amended with closure-gate override + letter normalization (AC2) — commit `0b56ff86`.

Mode: `docs-only` (§2.5b); no version bump; no tag. Disconnect is the merge commit.

---

## Friction Log

- **eng/writing/SKILL.md not found.** The Tier 3 writing skill declared in the issue (`src/packages/cnos.eng/skills/eng/writing/SKILL.md`) did not exist at that path. The §3.8 prose amendment was authored against the existing §3.8 style. No content constraint was missed — the amendment is consistent with adjacent rubric paragraphs — but the load failure is disclosed. The issue's skills-to-load list cited a file that doesn't exist in the repo.

- **Script extension vs new script.** The issue offered two shapes: extend `validate-release-gate.sh` or add `validate-pre-merge-closure.sh`. α chose extension (option A / preferred shape) because the artifact-presence loop is already correct and parameterizing on mode avoids duplication. No friction in the decision; issue preference was clear.

- **Provisional close-out constraint.** Dispatch protocol requires alpha-closeout.md to be written before exit (§2.8 fallback). The standard post-merge re-dispatch path is not available in this session. The provisional marker is explicit.

---

## Observations

**Pattern: missing Tier 3 skill at dispatch-declared path.** The issue declared `cnos.eng/skills/eng/writing/SKILL.md` as a Tier 3 skill, but the file was not found at that path. This is the second time a Tier 3 skill has been declared by an issue but not found at dispatch time (first observed: [if applicable]). The issue's `## Skills to load` section should be validated against the repo before dispatch.

**Pattern: script extension as the correct shape for mechanical reinforcement.** The artifact-presence check in `validate-release-gate.sh` (lines 50–72) was exactly the right abstraction for the pre-merge gate. Extending via `--mode` rather than duplicating kept the diff minimal and made the behaviour immediately visible in the existing invocation documentation.

**Pattern: `<C` vs `C−` drift across artifact types.** The §3.8 rubric used `<C`; CHANGELOG/PRA/alpha-closeout used `C−`. Neither was wrong — they described the same disposition — but the gap was invisible in each individual artifact. The fix (declare equivalence explicitly) is minimal; the cost of not fixing it was continued ambiguity at each grading event.

**Pattern: recursive coherence as first-activation test.** The gate introduced by this cycle was immediately exercised against this cycle's own close-out triple (required by AC4). β running `scripts/validate-release-gate.sh --mode pre-merge` before merge approval is the canonical first-activation test for any gate. The dispatch prompt made this explicit; it was not inferred.

---

## Engineering-Level Reading

The core work is straightforward: a flag on a shell script and two paragraphs of prose. The structural value is not in the code complexity but in the enforcement gap it closes.

The pre-merge gate converts the closure-gate from a role-discipline constraint to a machine-verifiable constraint. `git merge` can still happen without running the gate (the gate is not wired into a CI workflow in this cycle), but a δ operator or β who runs the script gets an immediate, actionable diagnostic. The friction of running it is low; the benefit is not having to reconstruct close-outs retroactively.

The §3.8 amendment is similarly structural: it makes the math irrelevant when the gate fails. Without the override, a cycle missing all close-outs could technically score C+ (if per-axis grades averaged out). The override short-circuits that by making the classification precede the computation — a cleaner model for how grading should work when there's a hard constraint.

The letter normalization is editorial bookkeeping: two labels for one disposition is one label too many. `C−` (used operationally) maps to `<C` (rubric label). Declaring the equivalence prevents future annotators from treating the difference as meaningful.
