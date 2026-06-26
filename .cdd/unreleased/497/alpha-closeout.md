# alpha-closeout — cnos#497

cycle: 497
role: alpha
base_sha: c82750d24381b878c30cf80f09b0ccf4e50494e5

## Artifacts produced

- `.cdd/unreleased/497/self-coherence.md` — decision artifact (Model B verdict + Q1–Q7 answers)

## Diff summary

Only `.cdd/unreleased/497/` touched. No code. No renames. No migrations.

## Retrospective

The decision cell is a clean mode: analysis-to-text, no ambiguous implementation surface. The five structural arguments for Model B held up under β review with no findings. The ledger-vs-namespace distinction is the load-bearing insight; the remaining arguments reinforce it from orthogonal angles (write-locality from cnos#496, wave-master O(1), mechanical-orchestration simplicity).

## Debt

None. The decision is complete. 497B and 497C do not file.

---

## R1 note — α's proper role-pass on operator-supplied patch proposal

cycle: 497
role: alpha
round: R1
input: `.cdd/unreleased/497/operator-review.md` (HI's typed translation; 6 findings O1–O6)
patch_proposal_reviewed: `dd819f00` (HI-authored; reframed by operator as patch proposal pending α's adoption)

### What the operator finding actually was

The human-operator final-read on PR #499 returned `iterate (narrowly)`. The verdict left Model B's substance intact and named six findings that β R0's mechanically-scoped AC oracle did not surface. Read against `operator-review.md`, the findings are *not* "narrow mechanical text fixes" (the HI's framing in `dd819f00`); they are six distinct semantic corrections at different doctrinal levels:

- O1: asymptotic-complexity correctness (`O(1)` was the wrong bound for the receipt walk)
- O2: canonical typed-identity anchoring (`protocol_id` belongs in the typed receipt, not in operational corroborating surfaces)
- O3: citation accuracy against the canonical filesystem (no "wherever it lives" in a decision artifact)
- O4: CDD doctrine — cell closure ≠ boundary acceptance (γ matter)
- O5: CDS actor-collapse rule completeness — collapse declared, not just acknowledged (γ matter)
- O6: γ doctrine — every finding gets explicit triage (γ matter)

The first three live in α's `self-coherence.md`; the last three live in γ's `gamma-closeout.md`. α R1 owns the first three; γ R1 owns the last three.

### How adopting the operator findings improves the decision

The materially strongest correction is O2: anchoring canonical protocol identity to the typed receipt's `protocol_id` field (CDS pin: `cnos.cdd.cds.receipt.v1`). The R0 framing put protocol identity on operational surfaces (labels, scaffold metadata, PR body); the operator finding moves it onto the typed-runtime field the Go verifier already dispatches on. The decision artifact now states a clean rule the long-arc trajectory needs: **path = ledger owner; `protocol_id` = concrete protocol.** That rule is the typed-runtime-compatible framing; runtime does not infer protocol from prose or directory names.

O1 strengthens the structural argument's accuracy (the Model B advantage is constant *root count*, not constant *work*; that distinction matters when the argument is invoked against future protocol additions). O3 removes uncertainty from canonical citations (a decision about canonical ownership cannot itself carry "wherever it lives"). O4–O6 cleanly sharpen γ's closeout's doctrinal posture (γ R1 owns).

None of the corrections affect the verdict. The five structural arguments survive; Q1–Q7 still receive substantive answers; "497B and 497C never file" stands.

### α's adopt/adjust/reject decisions

After inspecting `dd819f00`'s textual edits to `self-coherence.md` §R0 against the operator's expected_change in `operator-review.md`:

- **O1: adopt with completeness extension** (Argument 4 + Q6 patches accepted as α's text; α additionally fixes the residue in `CDD Trace` Action-taken bullet — HI missed line 116's "wave-master O(1) resolution" which is now "wave-master single-root / constant-root-count discovery").
- **O2: adopt verbatim** (Verdict + Argument 1 + Q1 patches accepted as α's text).
- **O3: adopt verbatim** (cross-references patches accepted; both target paths verified to exist on the current tree).
- **O4, O5, O6: not α's matter** — surface in γ-owned `gamma-closeout.md`; α R1 carries no §R0 edits for these (the patches in `dd819f00` for those files are γ's to adopt/adjust/reject in γ R1).

The decision artifact at `.cdd/unreleased/497/self-coherence.md` after this R1 pass is α's work: §R0 text reflects α's adopted corrections; §R1 is α's own analysis (replacing the HI-authored §R1 in `dd819f00`).

### What the boundary-violation episode teaches about role discipline

The HI absorbed the operator findings inline as commit `dd819f00`, framing them as "narrow mechanical text fixes" and editing role-owned matter (α's `self-coherence.md` §R0 text + a new §R1; β's `beta-review.md` §R1; γ's `gamma-closeout.md` rewrite + α's R1 closeout note + β's R1 closeout note). The operator's ruling reframed `dd819f00` as an operator-supplied patch proposal and dispatched proper role passes.

The discipline lesson for α: even when the per-finding text edit looks small, the aggregate of six findings spanning three doctrinal levels (architectural anchoring; doctrinal application; per-role-doctrine completeness) is a role-pass-shaped workload, not a mechanical absorption. "Looks like a one-liner" is not a role-bypass authorization. The proper repair sequence is exactly what is now in flight: HI files `operator-review.md` (typed durable input); α/β/γ run independent R1 role passes; the HI's `dd819f00` text is reviewed as a patch proposal, adopted/adjusted/rejected by the role owning each surface, and explicitly taken-ownership-of in the role's R1 record.

The framing pressure that drove the HI to absorb inline was the absence of a mechanical review-return primitive (no `cn cell return` / `cn cell resume`). Without it, the live wake exited at `status:review` with no way to wake back into the cell for narrow iterations; the HI's only available actions were either inline-absorb or escalate-to-fresh-cell. Both options are wrong for "operator-final-read narrow iterate"; the correct option needs runtime support that does not exist yet.

### Carry-forward

**cnos#500** records the missing primitive: `cdd/review-return: mechanically route operator iterate verdicts back into an existing cell`. The doctrinal output of this episode is that primitive's spec; until it lands, the proper-role-pass-on-top sequence used here (HI files operator-review.md; roles R1) is the manual stand-in. α's carry-forward observation: the role-pass-on-top sequence works but is heavy; once the primitive lands, narrow operator iterates should route mechanically into the live (or re-spawned) cell rather than through HI's hands.

— α@cdd.cnos, cycle/497 R1, 2026-06-26 (UTC)
