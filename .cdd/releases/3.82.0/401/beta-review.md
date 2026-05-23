# β review — cycle/401 (Phase 6 of cnos#366)

**Reviewer:** β-collapsed-on-δ (per δ contract; β-α-collapse acknowledged).
**Round:** R1.
**Verdict:** APPROVED.
**β-rigor:** AC1 generic-surface declaration; AC5 cdr/epsilon (post-#395) actually cites the new generic surface; AC6 existing cdd-iteration.md files still validate.

## R1 review oracle pass

### AC1 — Generic ε surface

**β check:** the ROLES.md §4b text declares ε generically, without privileging cdd in normative position.

```
$ grep -ni "cdd\|cdr\|cds\|cdw\|cda" ROLES.md | grep -n "^[0-9]*:" | wc -l
```

Spot-read §4b.1 (lines ~270–278): names "α/β/γ/δ", "the protocol itself", "code or prose or analysis", "across packages". No CDD-specific normative declaration in §4b.1.

§4b.2 (lines ~282–289): names `protocol_gap_count`, `protocol_gap_refs`, `#ProtocolGapRef`. Cites `schemas/cdd/receipt.cue` as the generic kernel schema (factual — receipt.cue *is* the generic kernel per its own header per cycle #388). No CDD-normative declaration.

§4b.3 (lines ~293–303): pattern `{protocol}-{axis}-gap` is generic. cdd's four classes and cdr's six classes appear as labeled **examples** of the generic pattern. Both examples are present (not just cdd), and both cite their respective per-protocol epsilon SKILL.md files. The text "A protocol that declares no gap classes has no machinery for protocol-iteration" is generic and forward-applicable. β verdict: AC1 generic-surface declaration is correctly generic. PASS.

§4b.4 (iteration artifact + cadence rule): "For each instantiation, ε produces an iteration artifact when …". Generic. cdd and cdr both cited as labeled examples. The cadence rule is stated generically ("required only when `protocol_gap_count > 0`") and each instantiation is shown to bind a path. PASS.

§4b.5 (MCA discipline): generic three-branch shape (ship-now / next-MCA / no-patch). The "no-patch requires justification per CDD §13a for the cdd instantiation, and equivalent in other instantiations" clause correctly names CDD as one binding without making the rule CDD-specific. PASS.

§4b.6 (ε=δ collapse): the generic argument from §4 hats-vs-actors is referenced; not re-derived. No CDD-specific language. PASS.

§4b.7 (instantiation declaration shape): five-item declaration obligation, with cdd and cdr as labeled examples. PASS.

**AC1 verdict:** PASS — generic surface declares ε generically; cdd and cdr appear only as labeled examples of the generic pattern.

### AC2 — CDD instantiates ε

```
$ rg -c "cdd-skill-gap|cdd-protocol-gap|cdd-tooling-gap|cdd-metric-gap" \
    src/packages/cnos.cdd/skills/cdd/epsilon/SKILL.md
6
$ wc -l src/packages/cnos.cdd/skills/cdd/epsilon/SKILL.md
143
```

Six hits. Non-empty. Verified the file's substantive content is now CDD-specific instantiation:

- Header pointer block (lines 19–31): cites ROLES.md §4b as the kernel.
- §1 (lines 33–95): CDD-specific gap-class definitions, receipt-stream location, iteration-artifact path, cadence rule (citing ROLES.md §4b.4), aggregator, cross-repo trace.
- §2 (lines 97–119): CDD-specific operating point for ε=δ collapse.
- §3 (lines 121–135): cross-references.

No generic doctrine remains in the file (the previous §2 ε/δ doctrine was generic; it now lives in ROLES.md §4b.6). The CDD-specific operating-point discussion in the new §2 names the cnos operating point (one operator, handful of cycles per day, ε=δ collapse default), which is CDD-instantiation-specific and correctly remains here. PASS.

### AC3 — cdd-iteration.md cadence resolved

**Four-way drift check:**

```
$ rg "every cycle.*cdd-iteration|writes.*every cycle|written on every" \
    src/packages/cnos.cdd/skills/cdd/{epsilon,activation,gamma,post-release}/SKILL.md
(no matches)
```

No "every cycle" language remains in any of the four surfaces. All four now state "required when `protocol_gap_count > 0`" or equivalent. Drift closed.

```
$ rg "protocol_gap_count > 0|protocol_gap_count > 0" \
    src/packages/cnos.cdd/skills/cdd/{epsilon,activation,gamma,post-release}/SKILL.md | wc -l
```

Multiple hits across all four files. PASS.

### AC4 — Schema watched-fields verified

Per α self-coherence verification (reproduced):

```
$ cue vet -c schemas/cdd/*.cue schemas/cdd/fixtures/valid-generic-receipt.yaml -d '#Receipt'; echo $?
0
$ cue vet -c schemas/cdd/*.cue /tmp/mismatched-count.yaml -d '#Receipt'; echo $?
1
$ cue vet -c schemas/cdd/*.cue /tmp/mismatched-refs.yaml -d '#Receipt'; echo $?
1
```

Schema enforces the constraint. No change needed. Disposition: schema-unchanged.

β additional check: confirmed `schemas/cdd/receipt.cue` line 110 carries the consistency constraint `protocol_gap_count: len(protocol_gap_refs)`. The constraint is a CUE unification with the bare-field declaration on line 107 (`protocol_gap_count: int & >=0`). Unification of two declarations of the same field is the CUE idiom for combining a type-bound with a value-derivation; both hold simultaneously. When a fixture asserts a third value, all three unify, and inconsistency fails vet. PASS.

### AC5 — Cross-protocol reusability

**β-rigor check 1:** cdr/epsilon (post-#395) actually cites the new generic surface.

```
$ grep -c "ROLES.md §4b" src/packages/cnos.cdr/skills/cdr/epsilon/SKILL.md
2
$ grep -c "cnos.cdd/skills/cdd/epsilon" src/packages/cnos.cdr/skills/cdr/epsilon/SKILL.md
0
```

cdr/epsilon's opening pointer (header block, lines 16–28) targets `ROLES.md §4b`; line 34 (§1 first paragraph, on the cadence rule) cites `ROLES.md §4b.4`. The previous CDD-sibling pointer (`cnos.cdd/skills/cdd/epsilon/SKILL.md`) is removed. The CDR instantiation now inherits from the shared generic surface, not from the engineering sibling. PASS.

**β-rigor check 2:** no duplication of generic doctrine in cdr/epsilon.

Reading cdr/epsilon's substantive content:
- §1: six cdr-specific gap classes (correctly CDR-specific) + cadence rule sentence updated to align with generic rule (necessary internal consistency; not duplication).
- §2 ε/δ relationship: re-derives the small-protocol-regime argument for the research operating point. β observation: this *is* duplication of generic doctrine (ROLES.md §4b.6 makes the same argument in generic form). However, the research-specific framing ("research-claim transmission discipline," "overclaim class," "waves" instead of "cycles") gives the argument a research-discipline operating point that the generic doctrine deliberately does not commit to. β verdict: this is *re-framing*, not duplication; the research-specific instantiation of a generic argument is what an instantiation file is for. PASS with note.
- §3 persona/protocol/project boundary: CDR-specific.

**Refusal condition check:** "cdr/epsilon (post #395) has content that doesn't fit the generic surface → surface; either revise generic surface OR mark as cdr-specific extension." β verdict: every substantive piece of cdr/epsilon fits one of the three categories (a) instantiation specifics that the generic surface declares to live in the instantiation, (b) research-discipline re-framings of generic arguments (legitimate instantiation work), or (c) project/persona boundary (correctly CDR-specific). No refusal needed. PASS.

### AC6 — Existing cdd-iteration.md files validate

**β-rigor check:** sample of existing files passes under new schema/rule.

Sample read (just to confirm files are non-empty markdown with expected shape):

- `.cdd/unreleased/395/cdd-iteration.md` (83 lines, F1/F2 findings, dispositions, §2 no-findings observations, §3 trigger assessment) — valid markdown structure per post-release §5.6b per-finding shape.
- `.cdd/unreleased/396/cdd-iteration.md` (26 lines, "**None.**" findings, project-binding mapping cycle) — valid markdown structure; under new rule, would not have been required (empty findings); under prior rule, required. Backward-compat clause makes the existing file remain valid.
- `.cdd/releases/3.78.0/379/cdd-iteration.md` (63 lines) — valid.
- `.cdd/releases/docs/2026-05-17/369/cdd-iteration.md` (86 lines) — valid.

The schema does not validate `cdd-iteration.md` content (it validates `receipt.yaml`); the iteration file is referenced via `protocol_gap_refs[].ref`. The cadence-rule change does not alter what counts as a valid iteration file; it alters when one is *required*. Existing files remain valid by the same rule that made them valid before (markdown shape per post-release §5.6b per-finding template). PASS.

**Refusal condition check:** "existing cdd-iteration.md files fail to validate under new schema constraints → surface; the tightening is too aggressive; relax." No schema-side tightening occurred; no schema constraint was added in this cycle. The refusal does not fire. PASS.

### AC7 — Phase 7 prerequisite

`ROLES.md §4b` is a stable citable location. The §4b numbering will survive future additions to ROLES.md because §4b is positioned between two stable sections (§4a five-layer chain; §5 cdd reference) and adding sections between them is unlikely. Phase 7's CDD.md rewrite can cite `ROLES.md §4b` or `ROLES.md §4b.4` (for the cadence rule specifically) and the citation will remain valid. PASS.

## β observations (non-blocking, not findings)

**Obs-1:** §4b.3 enumerates cdd's four classes and cdr's six classes as examples. This is technically forward-binding for cdr (the cdr class names live in cdr/epsilon/SKILL.md §1, which is the per-protocol-overlay declaration). β verdict: forward-binding is appropriate because §4b is itself the generic surface that names what instantiations declare; citing the two known instantiations is an example-of-the-pattern, not normative duplication.

**Obs-2:** the post-release §5.6b pre-publish gate row (line 411) was updated to cite ROLES.md §4b.4 with the relative path `../../../../../../ROLES.md`. β verified the path: from `src/packages/cnos.cdd/skills/cdd/post-release/SKILL.md`, that resolves correctly to repo-root `ROLES.md`. ✓

**Obs-3:** the cdr/epsilon §1 cadence-rule sentence was updated beyond the strict δ contract scope ("cite the generic surface"). β verdict: this update is *necessary for internal consistency* — leaving the previous "written on every wave" sentence would have meant cdr/epsilon's header (claiming inheritance of the generic cadence rule) contradicts cdr/epsilon's §1 body. The scope expansion is one sentence, narrow, and internally-required. Acceptable.

**Obs-4:** activation/SKILL.md §22 was tightened; the severity scale (D/C/B/A + info) and auto-spawn MCA trigger (≥3 same-axis in 5 consecutive cycles) are preserved unchanged. β verified by reading the third and fourth paragraphs of §22 after the edit: both paragraphs survive intact. ✓

**Obs-5:** gamma/SKILL.md §2.10 row 14 was a single-line edit dropping the "Empty-findings cycles still write cdd-iteration.md" parenthetical. β verified: row 14 now reads consistently with the new rule and cites ROLES.md §4b.4 + epsilon/SKILL.md + activation/SKILL.md. The single-line scope does not overlap with Phase 5 (which is broader gamma/SKILL.md territory).

## Trigger assessment

| Trigger | Fire condition | Fired? | β note |
|---|---|---|---|
| Review churn | review rounds > 2 | **No** | R1 APPROVED. |
| Mechanical overload | mechanical ratio > 20% AND findings ≥ 10 | **No** | 0 binding findings; 5 advisory observations. |
| Avoidable tooling / environment failure | environment blocked the cycle | **No** | `cue` available; vet ran cleanly; no environmental friction. |
| CI red on merge commit | CI fails post-merge | **N/A** | merge not yet executed. |
| Loaded skill failed to prevent a finding | skill underspecified | **No** | the cadence rule resolution was the named work, not a skill failure. |

No §9.1 triggers fired.

## Verdict

**R1 APPROVED.** AC1–AC7 PASS mechanically. β-rigor checks (AC1 generic-surface declaration; AC5 cdr cross-reference verified; AC6 existing files validated) all pass. No binding findings. Five non-blocking observations recorded above. No fix-round needed.

The cycle is ready for close-out.
