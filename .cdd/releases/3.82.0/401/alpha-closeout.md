# α close-out — cycle/401 (Phase 6 of cnos#366)

**Issue:** [cnos#401](https://github.com/usurobor/cnos/issues/401) — ε upscope; cadence rule resolved.
**Mode:** design-and-build; γ+α+β-collapsed-on-δ. β-α-collapse acknowledged.
**Rounds:** R1 APPROVED (no fix-rounds).
**ACs:** 7/7 PASS.

## Surface delivered

Six files edited; two cycle artifacts created; one schema verified-unchanged:

1. **`ROLES.md`** — NEW §4b "Generic ε — the protocol-iteration role" (seven subsections: ε's observation domain, watched receipt fields, gap class taxonomy, iteration artifact + cadence rule, MCA discipline, ε=δ collapse, instantiation declaration shape).
2. **`src/packages/cnos.cdd/skills/cdd/epsilon/SKILL.md`** — rewritten as the CDD-specific instantiation: header pointer to ROLES.md §4b; §1 enumerates the four `cdd-*-gap` classes + receipt stream + iteration artifact path + cadence rule; §2 CDD-specific ε=δ operating point; §3 cross-references.
3. **`src/packages/cnos.cdr/skills/cdr/epsilon/SKILL.md`** — header pointer retargeted from `cnos.cdd/skills/cdd/epsilon/SKILL.md` to `ROLES.md §4b`; §1 cadence-rule sentence aligned with new generic rule. §§2–3 substantive content preserved.
4. **`src/packages/cnos.cdd/skills/cdd/post-release/SKILL.md`** — §5.6b body and pre-publish gate row updated with the new cadence rule "required only when `protocol_gap_count > 0`"; cites ROLES.md §4b.4.
5. **`src/packages/cnos.cdd/skills/cdd/activation/SKILL.md`** — §22 lead paragraph updated to the conditional rule; severity scale + auto-spawn MCA trigger preserved.
6. **`src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md`** — §2.10 row 14 updated; the empty-cycles-still-write parenthetical dropped.
7. `.cdd/unreleased/401/gamma-scaffold.md`, `.cdd/unreleased/401/design-notes.md`, `.cdd/unreleased/401/self-coherence.md`, `.cdd/unreleased/401/beta-review.md` — cycle artifacts.

Schema: `schemas/cdd/receipt.cue` verified to enforce `protocol_gap_count: int & >=0` + `protocol_gap_refs: [...#ProtocolGapRef]` + consistency `protocol_gap_count == len(protocol_gap_refs)`. No change needed (constraints from cycle #388 Phase 2.5 already satisfy AC4).

## Design decisions recorded

**Home choice for generic ε: `ROLES.md §4b`** (new section between §4a.4 and §5).

Rationale (full version in `design-notes.md §1`):

1. Pre-resolved by `src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md §Q3` (cycle #367's converged design surface) — that cycle considered all three δ-contract candidates and chose `ROLES.md` with explicit rationale.
2. ε is a role function, not a runtime substrate; `ROLES.md` is the doctrine home for the role-scope ladder.
3. A new `cnos.protocol-iteration` package adds boundary without adding substitutability — rejected per `design/SKILL.md §3.10`.
4. `cnos.core/doctrine/` is the runtime-substrate home (kernel, dispatch, harness primitives); putting ε under `cnos.core/doctrine/` would smear role doctrine onto the runtime layer.
5. δ contract's "ROLES.md §5" literal text is unavailable (§5 is "cdd as reference instantiation"); §4b is the natural placement (after §4a five-layer chain, before §5 cdd reference).

No operator preference question was surfaced — the design decision was pre-converged in #367 and recorded on main; the refusal condition "home choice has multiple defensible rationales and operator preference matters" does not fire.

## Cross-protocol verification result

**cdr/epsilon (post-#395) cross-references work cleanly.** The generic surface (ROLES.md §4b) supports both `cdd/epsilon/SKILL.md` and `cnos.cdr/skills/cdr/epsilon/SKILL.md` as instantiations without re-declaration. Each instantiation carries only:

- its gap class names (cdd: 4 classes; cdr: 6 classes),
- its receipt-stream location (cdd: `.cdd/releases/`; cdr: project-binding-dependent),
- its iteration-artifact path (cdd: `.cdd/...`; cdr: project-binding),
- its aggregator (cdd: `.cdd/iterations/INDEX.md`; cdr: project-binding equivalent).

No content from either instantiation has been orphaned. cdr/epsilon's CDR-specific six gap classes (cdr-data-gate-gap, cdr-overclaim-gap, cdr-reproduction-gap, cdr-citation-gap, cdr-oracle-ambiguity-gap, cdr-construct-drift-gap) and the research-discipline operating point (ε=δ collapse argument re-framed for research) remain in cdr/epsilon as instantiation-specifics.

## Schema-change disposition

**Unchanged.** Cycle #388 Phase 2.5 added the `protocol_gap_count` + `protocol_gap_refs` fields with the consistency constraint `protocol_gap_count == len(protocol_gap_refs)`. AC4 oracle verified: valid fixture passes vet; mismatched-count fixture fails vet; mismatched-refs fixture fails vet. The constraint already satisfies AC4 without modification.

The δ contract's "if `protocol_gap_count` semantics tighten, schemas/cdd/receipt.cue may need a constraint update; otherwise no contract change" — α verifies: no tightening required, no change.

## Backward compatibility verified

Existing `cdd-iteration.md` files (e.g. `.cdd/unreleased/396/cdd-iteration.md` with empty findings; `.cdd/unreleased/395/cdd-iteration.md` with two findings; `.cdd/releases/3.78.0/379/cdd-iteration.md` with three findings) all remain valid artifacts under the new rule. The new rule is "required when `protocol_gap_count > 0`", not "forbidden when `protocol_gap_count == 0`". The backward-compat clause is stated explicitly in three surfaces (post-release §5.6b, epsilon §1, activation §22).

The refusal condition "existing cdd-iteration.md files fail to validate under new schema constraints" did not fire because no schema constraint was added in this cycle.

## Commit shape

This cycle ships (anticipated) commits:

1. γ scaffold + design-notes (γ+α blended in collapsed mode)
2. ROLES.md §4b — generic ε doctrine
3. cdd/epsilon shrink to instantiation
4. cdr/epsilon cross-ref retarget
5. post-release/SKILL.md §5.6b cadence rule
6. activation + gamma cadence alignments (single squash with post-release)
7. self-coherence + β-review
8. close-outs + cdd-iteration

Identity rotation: γ scaffold + close-outs as `gamma@cdd.cnos`; α work as `alpha@cdd.cnos`; β review as `beta@cdd.cnos`; merge by `operator@cdd.cnos`.

## Round-1 work summary

- 1 round, R1 APPROVED.
- 0 binding findings.
- 5 non-blocking β observations (recorded in `beta-review.md` Obs-1 through Obs-5).
- Mechanical ratio: undefined (no β findings; ratio computed only when findings ≥10).
- Phase 7 unblocked (citable surface ROLES.md §4b stable).
- Phase 5 (gamma/SKILL.md) not blocked by this cycle's gamma row-14 single-line edit; scopes do not overlap.
