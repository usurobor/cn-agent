---
schema: cn.operator-review.v1
issue: 662
pr: 667
verdict: iterate
reviewer: human-operator
captured_by: kappa
captured_at: 2026-07-16 (UTC)
worker_pr_head_at_review: 4c3fc7f95eda139d7c07f6080fd83a41faa2ede8
findings_count: 6
---

# Operator-review — PR #667 (cnos#662 cycle/662 R2)

This artifact is **κ's operator-final-read verdict on PR #667**, materialized as typed durable
input for the cycle's R2 re-entry. κ recorded the verdict on cnos#662 (comment
2026-07-16T11:23:42Z); α/β/γ act on it. κ did not author the spec (κ≠α); κ's α owns the repair.

The body below transcribes κ's issue comment **verbatim**.

---

## Operator-final-read verdict on PR #667 — **ITERATE NARROWLY** (recorded by κ)

The architecture converges, but the normative spec has **six load-bearing blockers** in its most important promised output — the typed contracts future FSM/runner work depends on. β either missed or dismissed these as cosmetic; the human gate caught them. **PR #667 stays draft. Do NOT dispatch the CC ratification cell.** Return #662 through the shipped review-return path, repair, run a **fresh independent β R2**, then reconsider CC ratification.

### Blockers (α R2 repair contract)
1. **Canonicalize `cn.cell.contract.v1`.** §2's envelope and the worked instance disagree: `cell.id` vs `cell.issue`; `requested_output` **object** vs **list**; `constraints.non_goals` vs `cell.non_goals`; top-level `gates` vs `cell.gates`. **One exact shape must govern; the worked example must validate against it verbatim.**
2. **Reconcile intent.** §2 treats the issue *as* intent (`intent.source: issue`); §13 defines first-class `cn.intent.v1` produced before any cell exists. Both cannot be canonical. The contract must **reference or embed `cn.intent.v1`** (`intent_ref: {schema, id, carrier:{kind: github_issue, ref}}`); a GitHub issue is a **carrier/projection**, not intent identity.
3. **Split PC result by mode (tagged union).** §3.2 gives PC-D0 and PC-Wave the *same* result shape (`wave_ref` + `graph`), quietly turning every Planning Cell into a wave-producer. PC-D0 → `{mode:d0, artifact_ref, readiness: ready_for_coherence_review, requires_operator_gate}`; PC-Wave → `{mode:wave, wave_ref, graph, readiness: ready_for_wave_review, requires_operator_gate}`.
4. **Correct actor-collapse terminology.** Separate κ/α/β **activations** under one Sigma/GitHub identity are **hosting-identity collapse** (tracked by #664), **not** actor collapse (one activation performing two roles *inside one cell boundary*). Preserve that **κ did not author the spec** — protocol-level κ≠α held. Update §8 and the γ bootstrap declaration; stop calling every reuse of Sigma "actor collapse."
5. **Repair stale `§17` / `§17 Q6` references.** The document ends at **§16 Open questions** with five questions. Those refs carry unresolved policy/sequencing — fix them to the real targets.
6. **Fresh independent β R2 after the repairs.** δ-attestation was defensible for the two tiny R1 edits; it is **not** sufficient after substantive contract/schema corrections. Spawn a fresh independent β over the **entire revised contract surface**. The R0 lesson is "independent review matters," not "δ may substitute for β on the next substantive revision."

### Process
1. κ records this verdict *(done — this comment)*.
2. `cn cell return --verdict iterate` (review-return: `review → changes`).
3. α R2 owns the spec repair (blockers 1–5).
4. **Fresh independent β R2** reviews the entire revised contract surface (blocker 6).
5. γ updates the planning receipt.
6. Return #662 to `status:review` for a second operator-final-read.
7. **Do not dispatch CC ratification yet.**

### What is already good (do not lose)
WC/PC/CC as concrete telos classes of one CCNF kernel · CC↔ε reconciled · κ outside the cell · State-A / specified / illustrative-future honestly separated · shipped FSM + request-marker mechanics grounded · CC owns wave judgment, FSM owns wave transition · wave-level authorization (operator ≠ child scheduler) · mechanically-oriented class-specific V · bootstrap limitation disclosed. The spec is close; the defects are precisely in the typed contracts.

### Separate follow-up
The concurrent β firings exposed a real substrate **race** (two dispatch activations on one claimed cell). It helped here by accident; it can just as easily produce conflicting α matter or corrupted reconciliation. Filed as a **P1, distinct from #664/#665** — see below. *Do not rely on accidental duplicate β as a quality mechanism; a deliberate review quorum is the right model if multiple reviews are wanted.*

---

## R2 routing note (κ, control-plane)

- **verdict = iterate (narrowly).** Blockers 1–5 are α-repairable typed-contract-surface fixes consistent with the pinned D1–D10; blocker 6 is a process requirement (fresh independent β). None requires an unpinned architecture decision, so no operator-gate hold is opened at return time. If α discovers that any blocker cannot be repaired without an unpinned decision, α STOPs and returns a typed operator-gate hold to κ instead of guessing.
- **κ≠α preserved.** κ authored this operator-review; α owns and rewrites `CELL-RUNTIME-CLASSES.md`.
- **Substrate transparency.** The shipped `cn cell return`/`cn cell resume` label mutation depends on an authenticated `gh` CLI, which is not present in this interactive `cn-sigma` substrate. The `review → changes → in-progress` transition is effected via the equivalent control-plane label mutation and disclosed in the R2 claim comment on #662 (mirroring the prior cycle's interactive-substrate honesty). The `cn` binary built from `src/go/cmd/cn` this cycle; only its `gh`-dependent label step is substituted.
