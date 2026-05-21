# α design-notes — cycle/401 (Phase 6 of cnos#366)

**Mode:** design-and-build; α design phase before build phase. β-α-collapse acknowledged.

## §1 Home choice decision

**Decision:** Generic ε doctrine lives at **`ROLES.md` §4b** (new section between §4a five-layer chain and §5 cdd reference instantiation).

### Decision rationale

The home choice is **pre-resolved** by `src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md §Q3` (cycle #367's converged design surface), which states verbatim: "ε relocates to **`ROLES.md`** — the generic role-scope ladder doctrine at the repo root. The CDD-specific instantiation (`src/packages/cnos.cdd/skills/cdd/epsilon/SKILL.md`) is rewritten in Phase 6 as a thin pointer to the generic doctrine in `ROLES.md`, retaining only the CDD-instantiation details ε needs (where receipt streams live, what protocol gaps look like in CDD, what `cdd-iteration.md` is). **This cycle names the target; Phase 6 ships the move.**"

That cycle's §Q3 considered all three δ-contract candidates and chose `ROLES.md` with three rationales (reproduced as the operative ground for this cycle's home choice):

| Candidate | Authority scope | Result |
|---|---|---|
| **`ROLES.md`** | Generic role-scope ladder doctrine | **CHOSEN** — ε is a role function of the ladder pattern; instantiations inherit it. |
| `cnos.core/doctrine/` | Core runtime substrate | Rejected — role function and runtime substrate are distinct surfaces per `COHERENCE-CELL.md §Structural Prediction`; fusing them is the surface-smearing the structural prediction guards against. `cnos.core` is the runtime-substrate home (kernel, dispatch, harness primitives); ε is a role. |
| New `cnos.protocol-iteration` package | Dedicated subsystem | Rejected — per `design/SKILL.md §3.10` (prefer package/install cohesion over topic labeling), a new package boundary is justified when there is a real install/use unit to separate. ε's matter is doctrine evolution — patches that land across `cnos.cdd`, `cnos.core`, `cnos.eng` depending on where the gap surfaced. No install/use unit "the protocol-iteration runtime" exists; inventing one creates a package whose only purpose is to host one role's doctrine — the false-package-boundary failure §3.10 names. |

Reaffirming the choice for cycle/401:

1. **ε is a role function, not a runtime substrate.** `ROLES.md` already owns the generic α/β/γ/δ doctrine; ε's generic doctrine joining it preserves a single home for the ladder pattern.
2. **`ROLES.md` already cites ε as a generic role.** §1 names it as the iteration role; §2 names it as order-4 observation; §3 Field 5 names "ε iteration cadence"; §4 ε=δ collapse discussion. The skeleton is already there — Phase 6 fills in the body.
3. **The δ contract option A literal text "ROLES.md §5" is unavailable** (§5 is "cdd as reference instantiation"). The natural insertion point is **§4b** between §4a (five-layer enforcement chain) and §5 (cdd reference). This puts the generic ε declaration alongside the other generic-grammar pieces (role ladder §1, scope escalation §2, instantiation contract §3, hats-vs-actors §4, five-layer chain §4a). δ's contract grants α design authority over home choice "with rationale"; section number is a placement detail consistent with the spirit of "ROLES.md §5" (operator named the file, not literally §5 vs §4b).
4. **No operator preference question to surface.** The refusal-condition "home choice has multiple defensible rationales and operator preference matters" does NOT fire because the design decision was pre-converged in #367 and the rationale is recorded in code (RECEIPT-VALIDATION.md §Q3) on main. Re-litigating it would discard work the operator already converged.

### Section placement detail

The new section is **§4b — Generic ε (protocol-iteration role)**, placed immediately after §4a.4. §5 ("cdd as the reference instantiation") shifts in interpretation but keeps its number — it is still the first instantiation example, now appearing after both the abstract grammar (§§1–4a) and the generic ε declaration (§4b).

This placement also makes downstream Phase 7 (CDD.md rewrite) citation stable: "see `ROLES.md §4b` for the generic ε doctrine" is a path that does not move when later sections are added.

## §2 Generic ε content structure

The §4b body declares:

1. **Role definition.** ε observes the receipt streams emitted by closed coherence cells across many cycles; ε's matter is the protocol itself. Inside the cell, α/β/γ/δ do ordinary metabolism; ε observes from outside, across cells, and patches the protocol when receipt-stream patterns reveal structural gaps. (Inherits framing from `COHERENCE-CELL.md §ε as Protocol Evolution`.)
2. **Watched receipt fields.** Every receipt (per `schemas/cdd/receipt.cue#Receipt`) carries:
   - `protocol_gap_count: int & >=0` — the count of protocol gaps observed in this cycle.
   - `protocol_gap_refs: [...#ProtocolGapRef]` — typed structured references to those gaps.
   - Consistency constraint: `protocol_gap_count == len(protocol_gap_refs)`.
3. **Gap class taxonomy (generic pattern).** Each instantiation declares its own gap classes named to its loss function. The **pattern** is `{protocol}-{axis}-gap` where axis names the protocol-internal surface the gap touches. Examples:
   - **cdd (engineering, repairable-feedback regime):** `cdd-skill-gap` (skill underspecified), `cdd-protocol-gap` (CDD doctrine itself drifted), `cdd-tooling-gap` (tooling absent or wrong), `cdd-metric-gap` (measurement missing or wrong).
   - **cdr (research, truth-preserving regime):** `cdr-data-gate-gap`, `cdr-overclaim-gap`, `cdr-reproduction-gap`, `cdr-citation-gap`, `cdr-oracle-ambiguity-gap`, `cdr-construct-drift-gap` (per cnos#395).
   - cdw, cda, future: each instantiation declares its own class names against its discipline profile (per `ROLES.md §4a.2`).

   The generic doctrine does **not** mandate four classes (cdd has four; cdr has six). The doctrine mandates: classes are named in the instantiation; classes are typed against the instantiation's loss function; ε reads `protocol_gap_refs` typed against the instantiation's class set.
4. **Relationship to δ (the parsimony argument).** ε and δ may be the same actor in small-protocol regimes. The structural separation is not a headcount requirement; it is observational: δ operates the cycle sequence; ε observes whether the protocol itself is learning. Same actor can wear both hats when protocol-iteration volume does not justify dedicated specialization. Separation becomes warranted when iteration volume crowds out operator work or when the longitudinal view requires sustained attention a single operator cannot provide. (Inherits framing from current `cdd/epsilon/SKILL.md §2` and current `cdr/epsilon/SKILL.md §2`.)
5. **Iteration artifact (cadence rule).** Each instantiation produces an iteration artifact when its receipt's `protocol_gap_count > 0`. The receipt is the always-present record of *whether* an iteration artifact is required; the iteration artifact (e.g. `cdd-iteration.md`, `cdr-iteration.md`) is the conditional record of *what was found*. (Inherits framing from `COHERENCE-CELL.md §ε Artifact Rule` and `RECEIPT-VALIDATION.md §Q3.3`.)
6. **Instantiation pattern (declaration shape).** A protocol instantiating the ladder declares (per §3 Field 5 + this section):
   - The instantiation's gap class names (e.g. `cdd-skill-gap`, …).
   - The location of receipt streams ε reads (e.g. `.cdd/releases/{X.Y.Z}/{N}/`).
   - The iteration artifact name (e.g. `cdd-iteration.md`).
   - The iteration aggregator location (e.g. `.cdd/iterations/INDEX.md`).
   - The MCA discipline branches (ship-now / next-MCA / no-patch with reason).

## §3 Cross-protocol verification

The generic surface must support both `cdd/epsilon/SKILL.md` and `cdr/epsilon/SKILL.md` (post-#395) as instantiations without re-declaration. Verification by cross-reading:

**cdd/epsilon:** declares four `cdd-*-gap` class names; reads `.cdd/releases/{X.Y.Z}/{N}/` receipt streams; writes `cdd-iteration.md`; aggregates at `.cdd/iterations/INDEX.md`; MCA discipline = ship-now / next-MCA / drop. Each of these declarations fits the generic instantiation pattern (§2 item 6). ✓

**cdr/epsilon (post-#395):** declares six `cdr-*-gap` class names; reads project-binding-dependent wave-surface receipt streams; writes a CDR-iteration record (project-binding-dependent path); aggregates per project binding; MCA discipline = ship-now / next-MCA / drop. Each declaration fits the generic pattern. ✓

**The ONE detail that does not transfer literally:** cdr/epsilon §3 (current) explicitly states that the storage path is project-bound (the protocol-overlay layer declares the artifact's existence and shape, not its file path) while cdd/epsilon writes to a CDD-fixed path (`.cdd/...`). This is **not** a defect in the generic surface — the generic surface should not mandate a fixed file path. The cdd-fixed path is a CDD-specific instantiation choice (cdd is a single-project bootstrap binding); cdr's project-binding-dependent path is a CDR-specific instantiation choice (cdr is project-overlay over many possible repos). Both are valid instantiations of the generic pattern. The generic surface declares: "the iteration artifact has a stable per-protocol name and a per-instantiation path." ✓

**No surface revision needed.** The generic ε surface as drafted supports both instantiations without overreach. cdr/epsilon's "extension" content (the six gap class names + the project-binding storage discussion) is correctly CDR-specific and stays in `cdr/epsilon/SKILL.md`. The generic surface is the kernel; instantiations carry their specifics.

**Refusal condition not fired:** "cdr/epsilon (post #395) has content that doesn't fit the generic surface" does not apply. cdr/epsilon has content the generic surface declares to live in the instantiation, by design.

## §4 Cadence rule resolution

**Current state (drift across four files):**

| Surface | Current cadence rule |
|---|---|
| `cdd/epsilon/SKILL.md §1` | "writes `cdd-iteration.md` on every cycle. An empty findings list is itself signal" |
| `cdd/activation/SKILL.md §22` | "Every cycle produces a `cdd-iteration.md`" |
| `cdd/gamma/SKILL.md §2.10 row 14` | "if the close-out triage produced ≥1 finding tagged `cdd-*-gap`, `cdd-iteration.md` exists" + parenthetical: "Empty-findings cycles … still write `cdd-iteration.md` per `epsilon/SKILL.md §1` and `activation/SKILL.md §22`, but the INDEX row is optional for cycles with no findings" |
| `cdd/post-release/SKILL.md §5.6b` | "When the close-out triage in `gamma-closeout.md` produces ≥1 finding tagged `cdd-*-gap`, γ authors `.cdd/unreleased/{N}/cdd-iteration.md`" + `activation/SKILL.md §22` reference for "every cycle produces" |

**Drift summary:** epsilon and activation say every-cycle; post-release says conditional; gamma row 14 says conditional with parenthetical empty-cycles-also-write. RECEIPT-VALIDATION.md §Q3.3 (#367 design) states: "*cycle artifacts no longer require an ε-produced `cdd-iteration.md` unconditionally; the receipt's `protocol_gap_count == 0` is the no-gap signal, and `cdd-iteration.md` is required only when `protocol_gap_count > 0`*".

**Resolution (per δ contract):** Tighten to a single rule: **`cdd-iteration.md` required only when `protocol_gap_count > 0`**.

**Files to update for the resolution:**

1. `cdd/post-release/SKILL.md §5.6b` — δ contract names this surface. Update the rule text and pre-publish gate.
2. `cdd/epsilon/SKILL.md §1` — currently says "writes `cdd-iteration.md` on every cycle"; conflicts with new rule. Update during epsilon shrink (§5 below) to match: "writes `cdd-iteration.md` when `protocol_gap_count > 0`".
3. `cdd/activation/SKILL.md §22` — currently says "Every cycle produces a `cdd-iteration.md`". This conflicts. Update §22 to match (within this cycle's scope, because the four-way drift surface includes activation by name).
4. `cdd/gamma/SKILL.md §2.10 row 14` — parenthetical "Empty-findings cycles still write `cdd-iteration.md`" conflicts with new rule. Update to drop the parenthetical (Phase 5 is gamma territory, but this one-line consistency tightening prevents resurrecting the drift; if Phase 5 has not yet landed this file is otherwise untouched, and the edit is a single-line change to row 14's parenthetical, which does not overlap with Phase 5's gamma scope).

**Backward compatibility:** The new rule is "required when > 0", not "forbidden when = 0". Existing `cdd-iteration.md` files with empty findings (e.g. #347, #393, #394, #396) remain valid artifacts. They are no longer *required* under the new rule, but they are not invalid. Future cycles may skip the file when count = 0. AC6 oracle passes by construction.

**Refusal condition not fired:** "post-release §5.6b cadence rule conflict: existing cycles … wouldn't have been required to write." Verified by reading the new rule wording: it is conditional on the *presence* of findings, not *absence*. Existing files are not invalidated.

**Note on activation §22 severity scale + auto-spawn MCA trigger.** The `activation/SKILL.md §22` text additionally specifies a severity scale (D/C/B/A + `info`) and an auto-spawn MCA trigger (≥3 same-axis findings in 5 consecutive cycles). These features attach to findings, not to file-existence. They are preserved unchanged when §22's "every cycle produces" lead clause is tightened to "produced when `protocol_gap_count > 0`". Cycles with zero findings produce no file; cycles with findings continue to apply severity and the sliding-window trigger.

## §5 cdd/epsilon shrink — content plan

The current `cdd/epsilon/SKILL.md` has two sections (§1 cdd-side scope, §2 ε's relationship to δ). The shrink keeps only the CDD-specific instantiation. New shape:

- **Frontmatter:** unchanged (artifact_class, parent, scope, governing_question, triggers).
- **Header pointer:** "This is a CDD-specific instantiation of the generic ε doctrine declared in [`ROLES.md §4b`](../../../../../../ROLES.md). The kernel grammar (ε-as-protocol-reviewer, watched receipt fields, MCA discipline, ε=δ collapse rule) is inherited by reference. Only CDD-specific gap classes, the iteration artifact path, and the receipt-stream location diverge." (parallels cdr/epsilon's header pointer pattern.)
- **§1 ε's CDD-side scope.** CDD-specific gap class names: `cdd-skill-gap`, `cdd-protocol-gap`, `cdd-tooling-gap`, `cdd-metric-gap`. Iteration artifact path: `.cdd/unreleased/{N}/cdd-iteration.md` (during cycle) → `.cdd/releases/{X.Y.Z}/{N}/cdd-iteration.md` (after release). Aggregator: `.cdd/iterations/INDEX.md`. Cadence rule cross-reference: "see `cdd/post-release/SKILL.md §5.6b`" — required only when `protocol_gap_count > 0`. Updates the lead from "every cycle" to the new conditional rule.
- **§2 ε's relationship to δ.** Shortened — the generic argument lives in ROLES.md §4b; this section carries only the CDD-specific operating point (small-protocol regime, ε=δ collapse normal in cnos's current scale, role becomes warranted-to-separate when …).

Net effect: 50% reduction in line count, with the surface-area-loss covered by the new ROLES.md §4b section.

## §6 cdr/epsilon cross-reference update

The current `cdr/epsilon/SKILL.md` opens with: "This is a CDR-specific extension of the generic cnos.cdd ε doctrine. The kernel grammar … is inherited by reference from [`cnos.cdd/skills/cdd/epsilon/SKILL.md`](…). Only the **trigger classes** … and the **output-artifact location** … diverge."

**Problem:** cdr inherits from a sibling CDD skill, not from generic doctrine. After Phase 6's upscope, generic doctrine lives at `ROLES.md §4b`. cdr's pointer should target the generic surface, not the CDD instantiation.

**Update:** Change the opening paragraph to:

> This is a CDR-specific instantiation of the generic ε doctrine declared in [`ROLES.md §4b`](../../../../../../ROLES.md). The kernel grammar (ε-as-protocol-reviewer in the scope ladder, watched receipt fields `protocol_gap_count` + `protocol_gap_refs`, the MCA discipline, the ε=δ collapse rule for small-protocol regimes) is inherited by reference. Only the **trigger classes** (research-failure classes per [`CDR.md`](../CDR.md) Field 5, not engineering-failure classes) and the **output-artifact location** (the CDR-iteration artifact under the project binding's wave surface) diverge for the research loss function per [`ROLES.md §4a.2`](../../../../../../ROLES.md#4a2-loss-function-distinction).

§§1–3 of cdr/epsilon are otherwise unchanged: §1 enumerates the six CDR-specific gap classes (correctly instantiation-specific); §2 discusses ε=δ collapse for research (mirrors the generic argument, project-specifically); §3 declares the persona/protocol/project boundary (CDR-specific).

**No content lost.** The "kernel grammar" content cited above no longer needs to be re-declared in cdr/epsilon — it lives in ROLES.md §4b — but cdr/epsilon never re-declared it; it only cited it. The cross-reference target moves; the substantive content does not.

## §7 Schema verification plan

`schemas/cdd/receipt.cue` already has (per cycle #388 Phase 2.5):

```cue
protocol_gap_count: int & >=0
protocol_gap_refs: [...#ProtocolGapRef]
protocol_gap_count: len(protocol_gap_refs)  // consistency constraint
```

**AC4 oracle:** `cue vet` passes on the schema (already passes per `schemas/cdd/fixtures/valid-generic-receipt.yaml`); a fixture with mismatched count/refs fails vet.

**Test plan:**

1. Run `cue vet schemas/cdd/*.cue schemas/cdd/fixtures/valid-generic-receipt.yaml -d '#Receipt'` — expect exit 0.
2. Author a temporary in-memory fixture asserting `protocol_gap_count: 1` + `protocol_gap_refs: []` — expect exit non-zero. (The consistency unification `protocol_gap_count: len(protocol_gap_refs)` makes `1` unify with `0`, which fails.)
3. Author an in-memory fixture with one gap ref + `protocol_gap_count: 0` — expect exit non-zero.

If all three pass per expectation, the schema already enforces the AC4 oracle and no schema change is needed. Disposition: **schema-unchanged**; AC4 satisfied by existing constraints from cycle #388.

**Refusal condition not fired:** "existing cdd-iteration.md files fail to validate under new schema constraints" requires a schema change in this cycle that imposes constraints incompatible with existing artifacts. No schema change is planned in this cycle; existing artifacts continue to validate.

## §8 Phase 7 prerequisite

AC7: the generic ε surface is at a citable location Phase 7 (CDD.md rewrite) can reference. After Phase 6 lands, Phase 7 can write in the rewritten CDD.md: "see `ROLES.md §4b` for the generic ε doctrine" — a stable path that does not change as ROLES.md grows new sections after §4b (§5+ shift in number context, not in citation target).

Section anchor: `#4b-generic-ε-protocol-iteration-role` (the GFM-generated anchor for the header text). The header text is canonical; downstream citations should use the section number, not the anchor, to remain stable under future header re-wording.

## §9 Build-phase summary

Six edits in order:

1. `ROLES.md` — insert new §4b "Generic ε (protocol-iteration role)" between §4a.4 and §5.
2. `src/packages/cnos.cdd/skills/cdd/epsilon/SKILL.md` — shrink per §5 above; update cadence-rule lead.
3. `src/packages/cnos.cdr/skills/cdr/epsilon/SKILL.md` — opening cross-reference points at `ROLES.md §4b`.
4. `src/packages/cnos.cdd/skills/cdd/post-release/SKILL.md` §5.6b — cadence rule resolved.
5. `src/packages/cnos.cdd/skills/cdd/activation/SKILL.md` §22 — cadence rule tightened to match (single-paragraph edit on lead clause; severity scale + sliding-window unchanged).
6. `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` §2.10 row 14 — parenthetical updated (single-line edit) to drop empty-cycles-still-write language.

Then schema-verification run (no edit, just `cue vet`); β review; close-outs; cdd-iteration.md.
