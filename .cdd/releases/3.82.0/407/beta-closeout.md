# β close-out — cycle/407 (Sub 2 of cnos#403)

β-level findings from the cycle.

## Findings

**None binding.** β R1 APPROVED. All nine ACs PASS mechanically; AC4 no-duplication audit PASS post-categorization; doctrine-coherence-with-cdr-template structural diff PASS; name-overpromise audit on Sub 3/4/5 references PASS.

## Non-binding observations (β-rigor surface)

**β-Obs-1: AC4 mechanical sliding-window audit is a portable β-rigor check for contract-authoring cycles.** The 50-char sliding-window scan (case-insensitive, whitespace-collapsed) against the named kernel files produced a categorizable hit table that made the AC4 self-check mechanical rather than judgment-bound. β recommends Sub 3's `beta/SKILL.md` codify the AC4 sliding-window audit as a pre-review check for any cycle whose matter is doctrinal-contract authoring (any cycle whose deliverable is a `*.md` file under `skills/cds/`, `skills/cdr/`, or future `skills/cdw/` directories that cites kernel files).

Filing as **β-recommended candidate `cds-skill-gap`** for Sub 3's β-skill authoring: the AC4 audit procedure (50-char sliding-window scan with hit-categorization table) should be a documented β-rigor procedure. Not surfaced as binding gap (α caught all duplications before β-review; β verified the categorization independently); recorded for ε observation.

**β-Obs-2: The per-Field Sub-N-vs-Field-M scope-discipline lines in CDS.md are a structural improvement over CDR.md's per-Field shape.** CDR.md's six fields do not include explicit Sub-N-vs-Field-M boundary statements (CDR has only 4 sub-deliverables; the boundary is less load-bearing). CDS.md's six fields each include the line, making the contract-vs-operational boundary mechanical. β assessment: the lines are the per-cycle dispatch surface — Sub 3/4/5's γ reads CDS.md Field N and dispatches against the named taxonomy + the named "Sub N owns operational detail" line, without needing to re-derive which Sub owns what content. The lines are doctrine-coherence improvements over CDR's shape; not over-promising, but more rigorous scope-discipline.

Filing as **β-recommended candidate `cdr-protocol-gap`** for a future CDR.md edit: add per-Field Sub-N-vs-Field-M scope-discipline lines to CDR.md's six fields (analogous to CDS.md's). This is a CDR-side concern (not a CDS concern), so it does not surface as a CDS gap; recorded as a cross-protocol cross-reference observation for the future CDR Sub 5+ work (if any).

**β-Obs-3: Field 3's "relationship between the two surfaces" paragraph (narrative close-outs + typed `#CDSReceipt`) is the load-bearing innovation that distinguishes CDS's Field 3 from CDR's.** CDR's Field 3 declares only the typed receipt; CDS's Field 3 declares both surfaces because the engineering substrate is mid-transition (pre-Phase 3 of cnos#366, V is wired in but most CDS cycles still author narrative close-outs as the primary surface). β verified the composition statement: "the typed receipt's `evidence_refs` point at the close-out files; V dereferences the refs to validate the receipt against the cycle's contract." This is accurate per `schemas/cds/receipt.cue`'s declared `evidence_refs` structure (read at α-time; the schema's `evidence_refs.{self_coherence, beta_review, alpha_closeout, beta_closeout, gamma_closeout, evidence_root}` set matches CDS.md Field 3's narrative artifact set). Acceptable; not a finding.

**β-Obs-4: Field 5's `cds-*-gap` rename declaration is a doctrine commitment Sub 5 must honor when authoring `cds/epsilon/SKILL.md` and the `cds-iteration.md` per-finding shape.** β verified the rename is mechanical (per ROLES.md §4b.3 `{protocol}-{axis}-gap`); the four classes (skill/protocol/tooling/metric) are preserved across the rename; the transition is described as organic. Sub 5 will need to author the `cds-iteration.md` per-finding shape using the new names; the file-rename (`cdd-iteration.md` → `cds-iteration.md`) is a separate coordination item (currently noted in CDS.md Field 5 as Sub-5-territory; consistent with `docs/extraction-map.md §14` open coordination questions). Acceptable; downstream-Sub coordination is documented.

**β-Obs-5: Field 6's β-α-collapse-on-δ permission is the first canonical doctrine declaration of this collapse pattern.** β verified: the empirical anchor (cycles #388 through #407) has been operating under this pattern via the breadth-2026-05-12 wave manifest precedent + dispatch-shape convention; no prior doctrine surface named the pattern with stated conditions. CDS.md Field 6 names it with three conditions (matter is structural-mirror, migration, or contract-authoring; no novel substantive code; mechanical correctness verifiable from the spec) and the configuration-floor consequence (γ-axis and β-axis capped at A- per release/SKILL.md §3.8). β assessment: the declaration is the doctrine-shape that the empirical practice has been operating without; making it canonical is doctrine-coherence improvement. The exclusion clause ("This collapse is not permitted for substantive software work") is sharp enough to be β-enforceable. Acceptable; β-confirmable.

**β-Obs-6: CDS.md Field 4's three-time-scale cadence (per-cycle, per-release, per-wave) is a structural CDS-vs-CDR divergence accurately framed.** β verified the three time-scales map to engineering practice: per-cycle is the unit of work (branch → commit → merge); per-release is the bundling unit (PRA + RELEASE.md + release tag); per-wave is the multi-cycle coordination unit (cnos#366 phases, cnos#376 cdr wave, cnos#403 cds wave). CDR.md Field 4 declares only per-wave because research has no release; the divergence is correctly attributed to the loss-function difference. The "Cadence triggers" sub-list (GO unblocks dispatchables; REVISE triggers in-flight repair; NO-GO triggers follow-up cycle; OVERRIDE triggers debt-cycle; release-readiness triggers per-release δ work) is exhaustive per the boundary-decision vocabulary CDS Field 3 declares. Acceptable; doctrine-coherence-with-#CDSReceipt confirmed.

**β-Obs-7: The §"Architecture choice" rationale (5) renaming from "Decision-once-applied-twice" (CDR) to "Decision-once-applied-thrice" (CDS) is a numerical update reflecting empirical reality.** β verified: cnos#388 (schemas) → cnos#376 (research skills) → cnos#403 (engineering skills) is indeed the third application of the option-(a) split decision. The renaming is accurate; not over-promising; reflects the post-cnos#376 state at the time CDS.md is authored. Acceptable.

## Skill-loading effectiveness

| Skill | Loaded for | Effective? |
|---|---|---|
| `cdd` | β-cycle convention | ✓ Provided the β review template (per-AC oracle check + β-rigor checks + observation/finding categorization + trigger assessment + verdict). |
| `cdd/design §3.2` | One-source-of-truth principle | ✓ The principle was the load-bearing check for AC4; β verified the categorization of the 11 remaining 50+char overlaps against the principle's "citation vs duplication" line. |
| `cdr-Sub-1 precedent` (read-only) | Structural template | ✓ The structural diff of CDS.md against CDR.md was the load-bearing β-rigor check for doctrine-coherence-with-cdr-template. |

No β-skill failure.

## Configuration-floor declaration

Per `release/SKILL.md §3.8`, β-axis is capped at A- (β-α collapse on δ acknowledged). β's R1 APPROVED verdict respects the cap: the β verdict is a *quality verdict*, not an *independent-actor verification*. For substantive software work, a non-collapsed β would be required; for skill/docs-class cycles per the breadth-2026-05-12 wave manifest precedent, the collapse is acceptable and the configuration-floor cap is the structural acknowledgment.

## Closure

β-side closed. R1 APPROVED. No binding findings; seven non-blocking observations recorded above. γ close-out proceeds.
