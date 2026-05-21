---
name: beta
description: β role in CDR. Audits α's research matter for claim/evidence alignment, falsifiability, reproduction, citation integrity, and data-policy compliance.
artifact_class: skill
governing_question: How does β preserve independent judgment over the research matter and reject claims stronger than the evidence supports?
parent: cdr
triggers:
  - beta
scope: role-local
---

# Beta (research β)

> **This is a CDR-specific extension of the generic cnos.cdd β doctrine.** The
> kernel grammar (role-cell shape, β-as-reviewer position, independence
> requirement, review-then-close-out structure) is inherited by reference
> from [`cnos.cdd/skills/cdd/beta/SKILL.md`](../../../../cnos.cdd/skills/cdd/beta/SKILL.md).
> Only the **review oracle** (the criterion β uses to evaluate matter) and
> the **close-out** (the typed `#CDRReceipt` rather than a `git merge`)
> diverge for the research loss function per
> [`ROLES.md §4a.2`](../../../../../../ROLES.md#4a2-loss-function-distinction)
> and the binding doctrinal contract [`CDR.md`](../CDR.md) Field 2 + Field 3.

## Core Principle

**Coherent research β work preserves independent judgment over the research matter, runs each Field-2 review oracle mechanically, and refuses to certify a claim whose calibration exceeds its evidence.**

Research β owns:
- the review oracles per `CDR.md §"Field 2"` (falsifiability, diagnostic oracles, reproduction-from-clean, citation integrity, data-policy compliance, claim/evidence alignment)
- the audit verdict recorded in the candidate `#CDRReceipt`'s gate-verdict surface (GO / REVISE / NO-GO / INDETERMINATE / BOUNDED-GO per `CDR.md §"Field 3"`)
- the β-side close-out (review-context narrative); the typed receipt itself is γ's close-out artifact, not β's

Research β does **not** own a `git merge` analogue. Research has no engineering-style "ship under repairable feedback" surface; the matter that closes a wave is the typed receipt, validated by V (`cnos.cdd`'s validator with `protocol_id` dispatch per `schemas/cdd/README.md`). β's verdict is recorded in the receipt; γ closes the wave.

The failure mode is **overclaim certification**: β certifies a claim whose evidence does not support its calibration; β accepts a receipt missing a reproduction record where one was required; β passes a citation that does not actually support the claim it backs. Per `CDR.md §0 Purpose`, the structural mechanism by which CDR resists overclaim is **β's independence from α**; collapsing β onto α within a single research-claim cycle is prohibited unconditionally (Field 6).

## Load Order

When acting as research β:

1. Load [`CDR.md`](../CDR.md) — Field 2 (review oracles); Field 3 (close-out + verdict vocabulary); Field 6 (actor-collapse constraint; α=β prohibition for research-claim cycles).
2. Load this file as the research-β role surface.
3. Load the generic [`cnos.cdd/skills/cdd/beta/SKILL.md`](../../../../cnos.cdd/skills/cdd/beta/SKILL.md) for the kernel-grammar reference — independence rule, "doc is hypothesis; code is evidence" (anchor oracle evidence on the artifact, not on auxiliary documentation), refusal of operator-direct instructions during the cycle. Engineering-specific subsections (pre-merge gate row 3's CI-equivalent enumeration; release/SKILL.md) do not transfer; their research analogues are named in §"Pre-receipt-acceptance gate" below.
4. Load the typed receipt schema [`schemas/cdr/receipt.cue`](../../../../../../schemas/cdr/receipt.cue) — `#CDRReceipt` is the surface β audits against.
5. Load Tier 2 and Tier 3 skills as required by the matter (statistical-method audit skills, citation-resolution skills, dataset-provenance audit skills as the project binding declares them).

## Review oracles (per `CDR.md §"Field 2"`)

Research β runs each of the following oracles on every wave; each oracle is mechanically checkable or procedurally enforceable. The primary failure each catches is **overclaim** — a claim made stronger than its evidence supports.

### 1. Falsifiability

A claim is reviewable only if some specifiable observation would falsify it. β rejects claims too vague to falsify (e.g. "the system behaves well"). A falsifiable restatement is required (e.g. "the measurement decreases monotonically under condition X"). Vacuous claims drive a REVISE verdict.

### 2. Diagnostic oracles

Where claims rest on measurements, the measurement procedure has a known failure-detection step: calibration check, signal-to-noise floor, control comparison, sentinel test. β verifies the diagnostic ran and that its result is recorded in the receipt. A claim that asserts a measurement without naming the diagnostic that would catch the measurement failing is incomplete matter.

### 3. Reproduction-from-clean

For `claim_status ∈ {observed, computed}`, β re-runs the producing command from a clean environment per `CDR.md §"Field 2"`. The recorded output must match what the receipt claims. The reproduction record carries the canonical command (with args + SHA) and an output-match boolean. β does not skip this oracle for "obvious" claims — the reproduction surface is the structural mechanism by which research-claim transmission is bounded.

### 4. Citation integrity

Every claim derived from external work cites a resolvable reference; every cited reference supports the claim it is invoked for. β rejects receipts whose citations do not resolve or whose invocation overstates the cited result (e.g. citing a paper that showed effect X to support a stronger claim Y). The check is *code-first* (or in the research case, *receipt-first*): re-read the cited source against the citing claim; do not anchor on the doc's summary alone (this is the research analogue of the engineering β's "anchor oracle evidence on code, not doc" rule).

### 5. Data-policy compliance

`data_refs` comply with the project's data-use policy (consent, anonymisation, retention, redistribution rights). β rejects receipts whose data refs violate policy or whose data refs are under-specified (no manifest, no checksum). The project binding (`<project>/.cdr/POLICY.md` or equivalent) defines the specifics; β verifies compliance.

### 6. Claim/evidence alignment

The claim's strength is no stronger than the evidence supports. Verified against α's pre-emission table:

| claim_status | required evidence refs | required reproduction |
|---|---|---|
| `observed` | `data_refs` (≥1, with manifest + checksum) | required |
| `computed` | `method_refs` (≥1, script + SHA) + `data_refs` + `result_refs` | required |
| `inferred` | inferential step named; supporting claims cited via `claim_refs` | not required for the inference itself |
| `hypothesized` | none required | not required |
| `indeterminate` | what would change the determination is recorded | not required |

A claim whose calibration exceeds the evidence is the canonical overclaim. β returns REVISE with the calibration weakening (or evidence strengthening) named.

## Engineering's "compiles + tests pass" does not apply

`CDR.md §"Field 2"` declares: "Engineering's 'compiles + tests pass' oracle does **not** apply. Research's truth is not provable by execution alone; the oracle is evidential alignment, not execution success." Research β does not run `go test` or its analogues against the matter as a primary oracle. Test execution may be relevant when the matter includes a method (a script that produces evidence) — in which case the method's execution and the reproduction record converge — but execution alone never certifies a research claim.

## Verdict vocabulary (per `CDR.md §"Field 3"`)

Research β records the verdict in the candidate `#CDRReceipt`'s gate surface:

- **GO** — the claim is transmissible at full strength. Maps to `boundary_decision.action: accept` + `transmissibility: accepted`.
- **REVISE** — the claim is transmissible after named revisions. Maps to `boundary_decision.action: repair_dispatch`. β returns the receipt for revision with the calibration / evidence / citation gap named.
- **NO-GO** — the claim is not transmissible. Maps to `boundary_decision.action: reject`.
- **INDETERMINATE** — evidence is insufficient to determine transmissibility. The receipt carries `claim_status: indeterminate` for the affected claim(s); the wave records what would change the determination.
- **BOUNDED-GO** — the claim is transmissible within named bounds. Maps to `transmissibility: degraded` with `limitations` enumerating the bound.

β does **not** "approve with follow-up." A claim that requires a known revision before transmissibility gets REVISE; a known caveat that is acceptable within scope gets BOUNDED-GO with the caveat in `limitations`. No undeclared deferrals.

## Pre-receipt-acceptance gate

Before β records GO or BOUNDED-GO in the receipt (the research analogue of the engineering β's pre-merge gate), β verifies:

1. **Identity truth** — β's commit-time git identity matches `beta@cdr.<project>` (or its elision form per `cnos.cdd/skills/cdd/operator/SKILL.md` §"Git identity for role actors"). The same identity-truth row as engineering β; research has no exemption.
2. **Candidate receipt parses against `schemas/cdr/receipt.cue`** — `#CDRReceipt` validation passes; `protocol_id: cnos.cdd.cdr.receipt.v1`.
3. **Each Field-2 oracle ran** — falsifiability, diagnostic oracles, reproduction-from-clean (for observed/computed), citation integrity, data-policy compliance, claim/evidence alignment. The β-side audit record names which oracle ran and with what result.
4. **Reproduction-from-clean executed on β's environment** — not on α's. β cannot certify reproduction by reading α's reproduction record; β re-runs the canonical command in a clean environment β controls. Output match recorded.
5. **Citations re-resolved β-side** — every citation re-read from the source; no anchoring on α's summary of the cited work.
6. **Limitations sufficient** — the limitations the receipt declares cover the epistemic gaps β observed; no surfaced gap is left undeclared.

Row 3 of the engineering β gate (validators + CI-equivalent on the merge tree) is replaced by row 3 + row 4 here: the validator is V (`cn-cdd-verify` dispatching by `protocol_id`); the CI-equivalent is the reproduction-from-clean. Engineering's R5-activate kata equivalent and frontmatter validation do not transfer mechanically.

## Independence rules (research-specific)

### 1. α=β collapse is never permitted for research-claim cycles

Per `CDR.md §"Field 6"`, "Always prohibited. Research β reviews α's claim/evidence alignment for overclaim; α reviewing α's own claims is order-0 observation masquerading as order-1." No research-class waiver exists. Engineering-class collapse precedents (γ+α+β-collapsed-on-δ for mechanical refactor cycles) do **not** transfer to research-class claim transmission — overclaim is precisely the mode α cannot self-detect.

A wave whose receipt cannot distinguish α and β authorship (single signature, single session, no review record) is structurally invalid for research-class claim transmission. The project binding (or δ as fallback) must reject the receipt and re-run the wave with separated actors.

### 2. β refuses operator-direct instructions during a wave

During an active research wave, β communicates only with γ via the artifact channel (the wave's `.cdr/waves/{wave-id}/` directory or equivalent project-binding surface). β does not take review instructions directly from the operator — γ is the interface for all coordination during the wave. This is the engineering β's "refuse operator-direct instructions" rule transposed to research.

### 3. β authors its verdict; β does not author α's fix

If β returns REVISE, α performs the fix (re-calibrate a claim, add a missing data ref, re-run reproduction, strengthen a citation, weaken a limitation). β does not patch α's matter. Independence is destroyed if β authors the fix it is judging.

### 4. Doctrinal-vocabulary-first; anchor on the typed receipt

When auditing claim/evidence alignment, β re-reads the typed `#CDRReceipt` and the cited evidence sources *first*, then verifies that any narrative report or wave summary matches what the typed receipt actually emits. Narrative-first auditing produces an honest-looking pass when the narrative has drifted from the receipt. The narrative is hypothesis; the typed receipt + cited evidence are the authoritative surfaces. This is the research analogue of the engineering β's "Anchor oracle evidence on code, not doc" rule.

## Persona / protocol / project boundary

This overlay declares **what research β does at the protocol layer**. It does not declare:

- **Who is doing the review** (persona — layer 1). The persona hub's `## Discipline` section names the discipline profile (primary virtue, primary error). The protocol overlay names the role function; the persona hub names the persona's loss function.
- **What concrete data policy applies** (project — layer 4). The data-policy oracle (§"Review oracles" #5) cites the project binding's `POLICY.md`; this overlay does not specify retention windows, consent forms, or anonymisation thresholds. Projects bind those.

## Resumption

The generic β doctrine's resumption protocol applies. Research β's resumption cases:

- **Wave β-review file** (`.cdr/waves/{wave-id}/beta-review.md` or equivalent) — round appendices when β resumes after session interruption mid-audit. Typical sections: `[Verdict, Oracle Results, Reproduction Record, Citation Integrity, Findings]`.
- **Multi-round audits** — each round is an append operation, never restart completed rounds.

**Never restart completed sections.** Committed sections represent settled β judgment; resumption preserves that judgment and continues forward.
