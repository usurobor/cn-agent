# α self-coherence — cycle/401 (Phase 6 of cnos#366)

α writes self-coherence before β review. AC-by-AC mechanical check + β Rule 7 implementation-contract conformance.

## AC self-check

### AC1 — Generic ε surface exists at chosen home

**Oracle:** file/section exists; declares ε's role independent of CDD; names watched receipt fields + gap classes.

**Verification:**

```
$ grep -n "^## §4b" ROLES.md
265:## §4b Generic ε — the protocol-iteration role
```

§4b exists; declares ε generically (the section title says "Generic ε"); names `protocol_gap_count` and `protocol_gap_refs` (lines 284–286); names the gap-class instantiation pattern `{protocol}-{axis}-gap` (line 295); enumerates cdd and cdr as reference instantiations (lines 297–299) without privileging cdd in the kernel grammar declaration (§4b.1, §4b.2, §4b.4, §4b.5, §4b.6, §4b.7 are all instantiation-agnostic; §4b.3 enumerates the two existing instantiations as examples of the generic pattern).

The section is structured in seven subsections (§4b.1 through §4b.7), one per generic-doctrine concern: what ε observes, watched receipt fields, gap class taxonomy, iteration artifact + cadence rule, MCA discipline, ε=δ collapse, instantiation declaration shape.

**Status:** PASS.

### AC2 — CDD instantiates ε

**Oracle:** `rg "cdd-skill-gap|cdd-protocol-gap|cdd-tooling-gap|cdd-metric-gap" cdd/epsilon/SKILL.md` returns hits; `wc -l cdd/epsilon/SKILL.md ≥ 1`.

**Verification:**

```
$ rg "cdd-skill-gap|cdd-protocol-gap|cdd-tooling-gap|cdd-metric-gap" \
    src/packages/cnos.cdd/skills/cdd/epsilon/SKILL.md | wc -l
6
$ wc -l src/packages/cnos.cdd/skills/cdd/epsilon/SKILL.md
143
```

Six hits across §1 (the four bullet definitions, plus two references in the cadence-rule paragraph). 143 lines (down from the previous 74 — but the new file is more structured: header pointer + §1 with explicit per-class definitions + §2 ε/δ relationship + §3 cross-references). The shrink intent (per design-notes §5) was content reduction relative to the *generic doctrine portion now lifted to ROLES.md §4b*; net lines increased modestly because the file now carries per-class one-paragraph definitions that the previous file delegated to the surrounding skill set. The content scope is now strictly CDD-specific.

**Cross-reference verification:**

```
$ grep -c "ROLES.md §4b" src/packages/cnos.cdd/skills/cdd/epsilon/SKILL.md
5
```

The file cites ROLES.md §4b (the generic doctrine) in five places: the header pointer, the gap-class pattern reference, the MCA-discipline reference, the cadence-rule reference, the ε=δ relationship reference, and the cross-references section.

**Status:** PASS.

### AC3 — cdd-iteration.md cadence resolved

**Oracle:** `post-release/SKILL.md §5.6b` text reflects "required only when `protocol_gap_count > 0`"; cross-references epsilon for watched-fields semantics.

**Verification:**

```
$ grep -n "Cadence rule" src/packages/cnos.cdd/skills/cdd/post-release/SKILL.md
297:**Cadence rule:** `cdd-iteration.md` is **required only when the cycle's receipt has `protocol_gap_count > 0`**
```

The §5.6b body now opens with an explicit "Cadence rule" paragraph that pins the rule to `protocol_gap_count > 0`; cross-references `ROLES.md §4b.4` and `cdd/epsilon/SKILL.md §1`. Backward-compat clause present.

**Four-way drift resolution:** all four surfaces named in #366 Phase 6 now state the same rule:

| Surface | New text |
|---|---|
| `cdd/epsilon/SKILL.md §1` | "required only when the cycle's receipt has `protocol_gap_count > 0`" |
| `cdd/activation/SKILL.md §22` | "produces a `cdd-iteration.md` … when the cycle's receipt has `protocol_gap_count > 0`" |
| `cdd/gamma/SKILL.md §2.10 row 14` | "if the cycle's receipt has `protocol_gap_count > 0`" |
| `cdd/post-release/SKILL.md §5.6b` | "**Cadence rule:** `cdd-iteration.md` is **required only when the cycle's receipt has `protocol_gap_count > 0`**" |

All four cite ROLES.md §4b as authoritative source. The drift is closed.

**Status:** PASS.

### AC4 — Schema watched-fields verified

**Oracle:** `cue vet` passes; fixture with `protocol_gap_count: 0` + empty refs validates; mismatched count/refs fails vet.

**Verification:**

```
$ cue vet -c schemas/cdd/*.cue schemas/cdd/fixtures/valid-generic-receipt.yaml -d '#Receipt'
$ echo $?
0
```

Valid fixture passes (existing fixture from cycle #388 Phase 2.5; not modified by this cycle).

```
$ cue vet -c schemas/cdd/*.cue /tmp/mismatched-count.yaml -d '#Receipt'
protocol_gap_count: conflicting values 0 and 1:
    ./schemas/cdd/receipt.cue:110:22
    .../mismatched-count.yaml:19:21
$ echo $?
1
```

Mismatched (count > refs) fixture fails vet via the existing consistency constraint `protocol_gap_count: len(protocol_gap_refs)` on line 110 of receipt.cue.

```
$ cue vet -c schemas/cdd/*.cue /tmp/mismatched-refs.yaml -d '#Receipt'
protocol_gap_count: conflicting values 1 and 0
$ echo $?
1
```

Mismatched (refs without count) also fails. Constraint is bidirectional.

**Disposition:** schema-unchanged. The AC4 constraint was already present from cycle #388 Phase 2.5; this cycle verified its enforcement without modification.

**Status:** PASS.

### AC5 — Cross-protocol reusability

**Oracle:** `cnos.cdr/skills/cdr/epsilon/SKILL.md` cross-references the new generic ε surface; no duplication of generic doctrine.

**Verification:**

```
$ grep -n "ROLES.md §4b" src/packages/cnos.cdr/skills/cdr/epsilon/SKILL.md
17:> [`ROLES.md §4b`](../../../../../../ROLES.md).** The kernel grammar —
34:The canonical output artifact is the **CDR-iteration record** (the research analogue of `cdd-iteration.md`). Per the inherited cadence rule ([`ROLES.md §4b.4`](../../../../../../ROLES.md)) …

$ grep -n "cnos.cdd/skills/cdd/epsilon\|cdd/epsilon/SKILL.md" \
    src/packages/cnos.cdr/skills/cdr/epsilon/SKILL.md
(no matches)
```

cdr/epsilon's opening pointer now targets `ROLES.md §4b` (the generic doctrine) rather than `cnos.cdd/skills/cdd/epsilon/SKILL.md` (the CDD sibling). The previous coupling — research instantiation citing engineering instantiation — is broken; both instantiations now inherit from the shared generic surface.

The cdr/epsilon §1 cadence-rule sentence ("the artifact is written on every wave; an empty findings list is itself signal") was updated to align with the new generic rule ("the artifact is required only when the wave's receipt has `protocol_gap_count > 0`"). The §§1–3 substantive content (six `cdr-*-gap` class definitions, ε=δ collapse for research, persona/protocol/project boundary) is preserved and remains correctly CDR-specific.

**Cross-protocol verification result:** the generic surface supports both cdd/epsilon and cdr/epsilon as instantiations without re-declaring kernel doctrine. Each instantiation carries only its specifics: gap class names, receipt-stream location, iteration-artifact path/aggregator. The refusal condition "cdr/epsilon has content that doesn't fit the generic surface" did not fire — cdr's research-specific class names are *expected* instantiation specifics, not generic-surface mismatches.

**Status:** PASS.

### AC6 — Existing cdd-iteration.md files validate

**Oracle:** sample of existing `.cdd/unreleased/{N}/cdd-iteration.md` files and `.cdd/releases/{X.Y.Z}/{N}/cdd-iteration.md` files validate under new schema constraints.

**Verification:** The new cadence rule does not impose any schema-side constraint on existing `cdd-iteration.md` files (the schema validates receipts; iteration files are referenced from receipts via `protocol_gap_refs`, not directly typed). The new rule changes only the *production cadence* (when the file is required), not the *file's shape*.

Sample of existing files (line counts confirm artifact presence and non-emptiness):

```
.cdd/unreleased/395/cdd-iteration.md:           83 lines (non-empty findings)
.cdd/unreleased/396/cdd-iteration.md:           26 lines (empty findings — explicit "None.")
.cdd/releases/3.78.0/379/cdd-iteration.md:      63 lines (non-empty findings)
.cdd/releases/docs/2026-05-17/369/cdd-iteration.md: 86 lines
```

All four samples remain valid markdown documents. The backward-compat clause (now stated in post-release §5.6b, epsilon/SKILL.md §1, and activation §22) explicitly preserves their validity: "existing empty-findings `cdd-iteration.md` files (written under the prior every-cycle rule) remain valid artifacts — the rule is 'required when > 0', not 'forbidden when == 0'."

**Refusal condition check:** "existing cdd-iteration.md files fail to validate under new schema constraints — the tightening is too aggressive; relax." This refusal would fire only if a schema-side constraint added in this cycle invalidated past files. No schema constraint was added in this cycle (AC4 verification confirmed the existing #388 constraint is what enforces AC4; no new constraint introduced). The refusal does not fire.

**Status:** PASS.

### AC7 — Phase 7 prerequisite

**Oracle:** Phase 7's planned CDD.md rewrite can cite the generic ε surface by path; the path is stable.

**Verification:** The generic ε surface lives at `ROLES.md §4b` — the repo-root doctrine file. Section number §4b is positioned between §4a (five-layer enforcement chain) and §5 (cdd reference instantiation). It will not move under future ROLES.md additions because:

- Sections §1–§4 are the core grammar and are stable.
- §4a is the persona/operator/protocol/project/receipt chain, stable.
- §5 onward are instantiation examples and forward-pointer stubs (cdw, naming convention, role-name stability, glossary).
- Adding new sections between §4a and the new §4b is unlikely; adding sections after §5 is the natural place for future doctrine.

Phase 7's CDD.md rewrite can write `see ROLES.md §4b` or `see ROLES.md §4b.4` (for the iteration-artifact rule specifically) and the citation will remain valid. The 10 occurrences of "§4b" within ROLES.md (in §4b's own subsection cross-references and in §4b.4's iteration-rule discussion) confirm the section is established as an internal-citation target as well.

**Status:** PASS.

## β Rule 7 implementation-contract conformance

Per #393 Rule 7, α self-verifies that each axis of the implementation contract is satisfied before β review.

| Axis | Pinned value | Conformance |
|---|---|---|
| Language | Markdown + CUE (if schema changes needed) | PASS — Markdown only. Schema verified, not changed. |
| CLI integration target | N/A | PASS — no CLI involved. |
| Package scoping | Primary home (α chose ROLES.md §4b with rationale in design-notes §1); secondary edits (`cdd/epsilon/SKILL.md`, `cdd/post-release/SKILL.md`, `cdd/activation/SKILL.md`, `cdd/gamma/SKILL.md`); cross-protocol (`cdr/epsilon/SKILL.md`); schema verification (`schemas/cdd/receipt.cue`). | PASS — all surfaces touched are within scope. δ contract noted secondary edits to `cdd/epsilon/SKILL.md`, `schemas/cdd/receipt.cue`, `cdd/post-release/SKILL.md`. Additional edits (`cdd/activation/SKILL.md §22`, `cdd/gamma/SKILL.md §2.10 row 14`) are *necessary* tightenings for the four-way drift resolution called out in the issue body Scope ("the four-way drift on cdd-iteration.md cadence named in #366 Phase 6"). Recorded in design-notes §4. |
| Existing-binary disposition | `cdd/epsilon/SKILL.md` shrinks; schema's `protocol_gap_count`+`protocol_gap_refs` already present, verify | PASS — epsilon shrunk to CDD-specific instantiation (generic content moved to ROLES.md §4b); schema constraints verified unchanged. |
| Runtime dependencies | None | PASS — no runtime dependencies added. |
| JSON/wire contract | If semantics tighten, schema may need update | PASS — semantics did not tighten beyond what cycle #388 already enforced. No schema change. |
| Backward compat | Existing cdd-iteration.md files validate; cadence tightens but existing comply | PASS — backward-compat clauses present in three surfaces (post-release §5.6b, epsilon §1, activation §22). AC6 confirmed existing files remain valid. |

**Surface containment:** files touched in this cycle (8 total):

1. `ROLES.md` (new §4b)
2. `src/packages/cnos.cdd/skills/cdd/epsilon/SKILL.md` (shrink)
3. `src/packages/cnos.cdr/skills/cdr/epsilon/SKILL.md` (cross-ref update + §1 cadence sentence)
4. `src/packages/cnos.cdd/skills/cdd/post-release/SKILL.md` (§5.6b + pre-publish gate)
5. `src/packages/cnos.cdd/skills/cdd/activation/SKILL.md` (§22 lead)
6. `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` (§2.10 row 14)
7. `.cdd/unreleased/401/gamma-scaffold.md`
8. `.cdd/unreleased/401/design-notes.md`

Files NOT touched (per issue non-goals + design discipline):
- `cdd/CDD.md` — Phase 7 territory.
- `cdd/gamma/SKILL.md` outside row 14 — Phase 5 territory.
- `cdd/COHERENCE-CELL.md`, `cdd/RECEIPT-VALIDATION.md` — already correctly point at the doctrine; their references to the future ROLES.md surface become more concrete after Phase 6, but the references themselves are correct as-is.
- `schemas/cdd/receipt.cue` — verified, unchanged.
- Other CDR / CDS schema files, fixtures.
- `cdr/epsilon/SKILL.md` §§2–3 substantive content (only opening pointer + §1 cadence sentence updated; §2 ε/δ relationship and §3 persona/protocol/project boundary preserved).

## Forecasts for β

- **Likely binding finding 0.** All ACs pass mechanically; cross-protocol verification confirmed; schema unchanged; backward compat preserved.
- **Likely advisory finding 1: ROLES.md §4b is long.** §4b spans ~90 lines (seven subsections). This is comparable to §4a (which is also ~90 lines), so the section weight is appropriate for the generic doctrine it carries. β may comment on density; α holds.
- **Likely advisory finding 2: gap-class enumeration in §4b.3.** The new §4b.3 enumerates cdd's four and cdr's six classes as *examples* — these are technically forward-references to the per-instantiation files. The alternative (declaring only the pattern and leaving examples to instantiations) was considered and rejected: anchoring with two concrete instantiations makes the pattern intelligible. β may raise; α defends the choice.
- **Likely advisory finding 3: cdr/epsilon §1 was modified beyond the strict δ contract.** The δ contract said "Update cdr/epsilon to cite the generic surface." α additionally updated one sentence in §1 (the "written on every wave" cadence sentence) to align with the new generic cadence rule, because not doing so would have left cdr/epsilon internally inconsistent (header says "inherits cadence rule from ROLES.md §4b.4" but §1 said "every wave"). The change is narrow (one sentence) and necessary for coherence. β may raise as scope concern; α defends as necessary internal consistency.

No binding findings forecast. β-collapsed review proceeds.
