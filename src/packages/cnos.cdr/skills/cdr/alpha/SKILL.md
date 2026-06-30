---
name: alpha
description: α role in CDR. Produces research matter — claims, hypotheses, methods, datasets, analyses, reports — under the truth-preserving claim transmission discipline.
artifact_class: skill
governing_question: How does α turn a research gap into review-ready research matter without producing a claim stronger than the evidence supports?
parent: cdr
triggers:
  - alpha
scope: role-local
kata_surface: embedded
inputs:
  - research gap
  - wave and dispatch context
  - evidence and sources in scope
outputs:
  - review-ready research matter (claims, hypotheses, methods, datasets, analyses, reports)
  - truth-preserving claim transmission discipline applied
---

# Alpha (research α)

> **This is a CDR-specific extension of the generic cnos.cdd α doctrine.** The
> kernel grammar (role-cell shape, α-as-producer position in the scope
> ladder, the algorithm structure) is inherited by reference from
> [`cnos.cdd/skills/cdd/alpha/SKILL.md`](../../../../cnos.cdd/skills/cdd/alpha/SKILL.md).
> Only the **discipline profile** and the **matter type** diverge: this
> overlay specifies what α does under the research loss function
> (truth-preserving claim transmission under uncertainty per
> [`ROLES.md §4a.2`](../../../../../../ROLES.md#4a2-loss-function-distinction)),
> not the engineering loss function (artifact improvement under repairable
> feedback). The binding doctrinal contract is [`CDR.md`](../CDR.md).

## Core Principle

**Coherent research α work produces calibrated claims with evidence-aligned strength, names every method/data/result reference required by the typed receipt, and declares limitations before β asks.**

The failure mode is **overclaim**: a claim stated stronger than its evidence supports; an `observed`-status claim with under-specified `data_refs`; a `computed`-status claim with a missing reproduction record; an `inferred`-status claim presented as if it were `observed`. Per `CDR.md §0 Purpose` + Field 1 + Field 2, overclaim is the structural failure CDR must resist; research α is the first line of defence.

Research α does not merely write up findings. Research α owns the matter set up to β's review:
- the calibrated claim set (each carrying a `claim_status`)
- the method refs (script paths + commit SHA)
- the data refs (mount points + manifests + checksums + data-use compliance records)
- the result refs (output file paths)
- the limitations and the reproduction record
- the candidate `#CDRReceipt` typed surface (per [`schemas/cdr/receipt.cue`](../../../../../../schemas/cdr/receipt.cue))

## Load Order

When acting as research α:

1. Load [`CDR.md`](../CDR.md) as the canonical doctrinal contract (Field 1 names what α produces; Field 2 names what β audits; Field 6 names the actor-collapse constraint).
2. Load this file as the research-α role surface.
3. Load the generic [`cnos.cdd/skills/cdd/alpha/SKILL.md`](../../../../cnos.cdd/skills/cdd/alpha/SKILL.md) as the kernel-grammar reference — for the role-cell shape, the artifact-order discipline, the resumption protocol, and the "do not outsource authoring work to β" rule. The engineering-specific subsections (pre-merge gate row 3; harness audit for schema-bearing changes; SHA convention for readiness signal) do not apply mechanically; their research-discipline analogues are named in §"Algorithm" and §"Pre-receipt-emission gate" below.
4. Load the typed receipt schema [`schemas/cdr/receipt.cue`](../../../../../../schemas/cdr/receipt.cue) — `#CDRReceipt` is the typed γ close-out surface α's matter must populate.
5. Load Tier 2 and Tier 3 skills as required by the research gap (statistical method skills, dataset-handling skills, citation-management skills as the project binding declares them).

The detailed step sequence belongs in `CDR.md` (Sub 1) and per-project research-method skills (project-binding scope). This file owns research-α's execution detail at the protocol-overlay layer: what each step *means* under the research loss function, what evidence each step requires, and what gates the matter must pass before β audits.

## Matter type (per `CDR.md` Field 1)

Research α produces the matter classes named in `CDR.md` Field 1:

- **Claim** — a stated proposition about the world or about a system under study. Every claim α produces carries a `claim_status` calibration: `observed | computed | inferred | hypothesized | indeterminate` (per `schemas/cdr/receipt.cue` `#ClaimStatus`). A claim without a calibration is incomplete matter; α does not hand off uncalibrated claims to β.
- **Hypothesis** — a proposition not yet calibrated. α-matter when the hypothesis itself is the wave's deliverable (a "what should we measure next" wave). Hypotheses ship with explicit `claim_status: hypothesized`.
- **Method** — a procedure that produces evidence. α records each method as a `method_refs` entry: script path + commit SHA per `ROLES.md §4a.3`. Methods without commit-SHA pinning are incomplete.
- **Dataset** — a body of observations with manifest and provenance. α records each dataset as a `data_refs` entry: mount point + checksum + manifest path + source attribution + data-use policy compliance record. Datasets without manifest or checksum are incomplete.
- **Analysis** — the application of method to dataset, producing result-file artifacts. α records each analysis as a `result_refs` entry (output file paths).
- **Report** — the synthesis surface (field notes, wave reports, project briefs). Reports cite the typed `#CDRReceipt` that backs them; reports are α-matter only when the cycle's deliverable is the synthesis itself.

**Research α does not produce software artifacts as primary matter.** Code produced in service of research is software-*evidence* (a method ref with commit SHA), not α-authorship. A wave that ships code as its primary deliverable is structurally a CDS cycle (engineering protocol), not a CDR wave; per `CDR.md §"Field 1"` + `ROLES.md §4a.2`, the dispatch boundary routes such work to CDS.

## Algorithm

The kernel algorithm is inherited from the generic α doctrine (Receive → Produce → Prove → Gate → Review loop → Close-out). Below is the research-discipline overlay; each step replaces the engineering-discipline content while preserving the structural position.

1. **Receive** — take the wave dispatch; identify the research gap (an unanswered question, an unverified claim, a needed dataset, a recurring measurement disagreement per `CDR.md §"Field 4"`); load declared methods, datasets, citation set.
2. **Produce in matter order** — author in the canonical research order:
   1. method (the procedure that will generate the evidence) — recorded as `method_refs` with commit-pinned script path
   2. data acquisition / mount / verification — recorded as `data_refs` with manifest + checksum + data-use compliance
   3. analysis run — recorded as `result_refs` (output file paths)
   4. claim drafting — each claim labelled with a `claim_status` proportional to the evidence type
   5. limitations enumeration — explicit caveats in `limitations`
   6. reproduction draft — if `claim_status ∈ {observed, computed}`, candidate canonical command for β's reproduction-from-clean
   7. report or wave-receipt narrative — synthesis surface citing the candidate `#CDRReceipt`
3. **Prove (claim/evidence alignment self-check)** — for each claim, verify the stated strength is no stronger than the evidence supports: `observed` claims have a `data_refs` entry and (if appropriate) a reproduction record; `computed` claims have `method_refs` + `result_refs` and a reproduction record; `inferred` claims name the inferential step explicitly; `hypothesized` claims are labelled as such; `indeterminate` claims are labelled and the wave records what would change the determination.
4. **Gate (pre-receipt-emission)** — pass the gate below before emitting the candidate `#CDRReceipt` to β.
5. **Review loop** — if β audits and returns REVISE (per `CDR.md §"Field 3"` verdict vocabulary), fix the matter (typically: weaken a claim's status, add a limitation, add a missing data ref, re-run the reproduction), re-audit affected claims, re-emit.
6. **Close-out** — when β approves (GO or BOUNDED-GO), the receipt is the matter; α's close-out narrative is research-side observations (which claims survived, which were re-calibrated, what limitations emerged), not a re-statement of the receipt.

## Claim-status calibration discipline

Per `CDR.md §"Field 2"` and `schemas/cdr/receipt.cue` `#ClaimStatus`, α calibrates each claim before review. The discipline:

- **`observed`** — α directly measured the claimed value. Requires `data_refs` with mount + manifest + checksum; requires a reproduction record (β re-runs the producing command in a clean environment per `CDR.md §"Field 2"`).
- **`computed`** — α applied a method to data and obtained the claimed value. Requires `method_refs` (script + commit SHA), `data_refs` (input), `result_refs` (output); requires a reproduction record.
- **`inferred`** — α drew the claim from other claims by an inferential step. The inferential step is named explicitly (e.g. "by linear interpolation between observations X and Y"); the supporting claims are cited in `claim_refs`.
- **`hypothesized`** — α states the claim as a candidate awaiting measurement. No `data_refs` required; the wave records what would change the calibration.
- **`indeterminate`** — α cannot determine which status applies (evidence is insufficient). The wave records what would change the determination.

A claim with calibration weaker than its evidence supports is acceptable (epistemic conservatism). A claim with calibration *stronger* than its evidence is **overclaim** — the failure mode this skill exists to prevent.

## Evidence-ref completeness rule

For every claim α emits, the typed receipt's evidence refs must satisfy:

| claim_status | data_refs | method_refs | result_refs | reproduction |
|---|---|---|---|---|
| `observed` | required (≥1, with manifest + checksum) | optional (if a measurement procedure was used) | optional | required |
| `computed` | required (≥1, the input data) | required (≥1, script + SHA) | required (≥1, output files) | required |
| `inferred` | optional (the supporting observed/computed claims may carry their own data_refs) | optional | optional | not required for the inference itself |
| `hypothesized` | not required | not required | not required | not required |
| `indeterminate` | optional (whatever evidence was gathered before the determination failed) | optional | optional | not required |

α verifies this table for every claim before emitting the receipt. β re-verifies as part of the review oracle (Field 2: claim/evidence alignment).

## Pre-receipt-emission gate

Before signalling review-readiness (β audits the candidate receipt), α verifies:

1. Every claim carries a `claim_status` calibration.
2. Evidence refs per the table above are populated for every claim.
3. Reproduction record present when required (`claim_status ∈ {observed, computed}`); the canonical command is recorded; the output match is verified locally (β will re-run from clean).
4. Citations are resolvable and support the claims they back (no claim citing external work has an unresolvable reference or an over-reach interpretation).
5. Data-policy compliance recorded for every `data_refs` entry (consent, anonymisation, retention, redistribution rights per the project binding's policy file).
6. Limitations enumerated; the wave's known epistemic caveats are stated explicitly.
7. Method refs include commit SHA (no "current HEAD" or "branch tip" — SHA-pinned).
8. The candidate `#CDRReceipt` parses against `schemas/cdr/receipt.cue` (typed validation passes; `protocol_id: cnos.cdd.cdr.receipt.v1`).

This is the research analogue of the engineering α's pre-review gate. Engineering's "CI green on the head commit" row is replaced by row 8 (typed-receipt validation); the "polyglot re-audit" row is replaced by the evidence-ref completeness table.

## Role-discipline rules (research-specific)

### 1. Treat the claim_status enum as a generation constraint

Claim status is a discipline anchor authored *at claim drafting time*, not a label applied after review. Drafting a sentence first and then choosing a status post-hoc is structurally identical to engineering's "code first, then write the spec" — both invert the dependency that prevents overclaim.

### 2. Do not outsource calibration to β

A claim with under-specified evidence is α's failure, not β's. β audits whether the calibration matches the evidence; β does not calibrate.

### 3. Limitations declared, not omitted

If a claim has known epistemic caveats (sampling bias, measurement-floor effects, generalisation limits), α declares them in `limitations`. Asking β to find the limitation α already knew about is a form of overclaim via omission.

### 4. Reproduction record is not optional for observed/computed

`CDR.md §"Field 2"` and the typed schema both require a reproduction record for `claim_status ∈ {observed, computed}`. α drafts the canonical command and verifies the output match locally; β re-runs from clean.

### 5. No software-protocol verbs

α does not "ship" or "release" research matter. α emits a typed `#CDRReceipt` that β audits and γ closes; the gate verdict determines transmissibility per `CDR.md §"Field 3"`. The cadence is gate-transition-shaped (Field 4), not release-shaped.

## Persona / protocol / project boundary

This overlay declares **what research α does at the protocol layer**. It does not declare:

- **Who is doing the work** (persona — layer 1; lives in `<persona-hub>/spec/PERSONA.md`). The discipline profile of the persona enacting α — primary virtue, primary error, default tempo, refusal conditions — is the persona hub's concern, not this overlay's.
- **What concrete data α touches** (project — layer 4; lives in `<project>/.cdr/`). Mount points, dataset manifests, project-specific gate thresholds, citation databases, report templates — all project-binding.

If a persona's discipline profile and this overlay disagree, the discipline profile rules: the persona hub authors *who* α is; this overlay authors *what* α-as-CDR-role does. If a project's binding rules add stricter constraints (e.g. "for externally-published claims, β must be a distinct human reviewer"), the project binding rules: this overlay declares the protocol floor; projects may impose stricter floors via `<project>/.cdr/POLICY.md`.

## Resumption

The generic α doctrine's resumption protocol applies. Research α's resumption cases:

- **`.cdd/unreleased/{wave}/self-coherence.md`** (or the research-wave analogue) — α resumes mid-self-coherence; read existing sections; continue from next uncompleted section.
- **Fix-round appendices** — when β returns REVISE and α resumes to fix calibration / re-run reproduction / add data refs, append fix-round sections without rewriting completed sections.
- **`.cdd/unreleased/{wave}/alpha-closeout.md` analogue** — if α is re-dispatched for close-out and finds a partial close-out, read completed sections and continue.

**Never restart completed sections.** Committed sections represent settled research-α matter; resumption preserves that matter and continues forward.
