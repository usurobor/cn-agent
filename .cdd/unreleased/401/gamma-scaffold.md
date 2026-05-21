# γ scaffold — cycle/401 (Phase 6 of cnos#366)

**Issue:** [cnos#401](https://github.com/usurobor/cnos/issues/401) — ε upscope: generic ε doctrine above CDD; cadence rule resolved.

**Parent:** cnos#366 (Phase 6; gates Phase 7 CDD.md rewrite). May run parallel with Phase 5.

**Mode:** design-and-build; γ+α+β-collapsed-on-δ. β-α-collapse acknowledged per Rule 7 implementation-contract conformance.

## Implementation contract (pinned by δ — issue body verbatim)

| Axis | Pinned value |
|---|---|
| Language | Markdown + CUE (if schema-side changes needed) |
| CLI integration target | N/A |
| Package scoping | Primary: chosen home for generic ε doctrine (ROLES.md / cnos.core/doctrine / cnos.protocol-iteration — α design decision). Secondary: `cnos.cdd/skills/cdd/epsilon/SKILL.md` (shrink), `schemas/cdd/receipt.cue` (verify watched fields), `cnos.cdd/skills/cdd/post-release/SKILL.md` (cadence rule). |
| Existing-binary disposition | `cdd/epsilon/SKILL.md` shrinks (CDD-specific instantiation only); `protocol_gap_count` + `protocol_gap_refs` already present in receipt.cue (cycle #388 Phase 2.5) — verify constraint |
| Runtime dependencies | None |
| JSON/wire contract | If `protocol_gap_count` semantics tighten, schemas/cdd/receipt.cue may need a constraint update; otherwise no contract change |
| Backward compat | Existing cdd-iteration.md files continue to validate; cadence rule tightens but existing artifacts comply (count ≥ 0 for all) |

## Surface

Files touched (anticipated):

1. **EDIT** `ROLES.md` — add new §4b "Generic ε (protocol-iteration role)" between §4a (five-layer chain) and §5 (cdd reference instantiation). Hosts the generic ε doctrine.
2. **EDIT** `src/packages/cnos.cdd/skills/cdd/epsilon/SKILL.md` — shrink to CDD-specific instantiation: gap class names (`cdd-skill-gap`, `cdd-protocol-gap`, `cdd-tooling-gap`, `cdd-metric-gap`), cdd-iteration.md artifact, cross-reference to generic ε at ROLES.md §4b.
3. **EDIT** `src/packages/cnos.cdr/skills/cdr/epsilon/SKILL.md` — change "kernel grammar inherited from cnos.cdd/skills/cdd/epsilon/SKILL.md" pointer to point at the new generic surface (ROLES.md §4b). Eliminates duplication of generic doctrine in cdr instantiation.
4. **EDIT** `src/packages/cnos.cdd/skills/cdd/post-release/SKILL.md` §5.6b — cadence rule resolved: "cdd-iteration.md required only when `protocol_gap_count > 0`". Backward-compat: existing empty-findings files stay valid.
5. **VERIFY** `schemas/cdd/receipt.cue` — confirm `protocol_gap_count` (int ≥ 0) and `protocol_gap_refs` ([...#ProtocolGapRef]) are required with consistency constraint `protocol_gap_count == len(protocol_gap_refs)`. Already present per cycle #388; no schema change anticipated unless a tightening surfaces.

Files NOT touched (out of scope per issue Non-goals):

- `cdd-iteration.md` per-finding shape (stays per `post-release/SKILL.md §5.6b`)
- Gap class names for CDD (`cdd-skill-gap` etc. stay)
- `gamma/SKILL.md` (Phase 5 territory)
- `CDD.md` rewrite (Phase 7 territory; this cycle produces the citable surface Phase 7 will reference)
- New `cnos.protocol-iteration` package (rejected per RECEIPT-VALIDATION.md §Q3 rationale; α design-notes records the rejection)

## AC oracle approach (issue body verbatim)

| AC | Oracle | Surface |
|----|--------|---------|
| AC1 | Generic ε surface exists at chosen home; declares ε role independent of CDD; names watched fields (`protocol_gap_count`, `protocol_gap_refs`) + generic gap class taxonomy. | ROLES.md §4b |
| AC2 | `rg "cdd-skill-gap|cdd-protocol-gap|cdd-tooling-gap|cdd-metric-gap" src/packages/cnos.cdd/skills/cdd/epsilon/SKILL.md` returns hits; file is non-empty (`wc -l ≥ 1`); cross-references generic ε. | cdd/epsilon/SKILL.md |
| AC3 | `post-release/SKILL.md §5.6b` text reflects "required only when `protocol_gap_count > 0`"; cross-references epsilon and ROLES.md §4b. | post-release/SKILL.md |
| AC4 | `cue vet` on generic fixture (`valid-generic-receipt.yaml`) passes; `protocol_gap_count == len(protocol_gap_refs)` constraint already enforced (line 110 of receipt.cue). Mismatched-count fixture FAILs vet. | schemas/cdd/receipt.cue |
| AC5 | `cnos.cdr/skills/cdr/epsilon/SKILL.md` cross-references the new generic ε surface (ROLES.md §4b) rather than `cnos.cdd/skills/cdd/epsilon/SKILL.md`; no generic-doctrine duplication. | cdr/epsilon/SKILL.md |
| AC6 | Sample existing `.cdd/unreleased/{N}/cdd-iteration.md` files validate under new schema (no schema change, so still valid). Cadence rule is on the **artifact requirement**, not on the artifact's shape; existing files do not need re-shaping. | existing cdd-iteration files |
| AC7 | Generic ε surface is at stable citable location (ROLES.md §4b) Phase 7 can reference for the compact ε doctrine in CDD.md rewrite. | path stability |

## Branch/identity

- Branch: `cycle/401` from `origin/main` (created)
- γ identity: `gamma@cdd.cnos`
- α identity: `alpha@cdd.cnos`
- β identity: `beta@cdd.cnos`
- δ identity: `operator@cdd.cnos` (post-merge to main)

## Dispatch shape

This is `cdd/operator/SKILL.md §5.2` (single-session δ-as-γ with γ+α+β collapsed on δ-as-agent). 7 ACs; mode is design-and-build with one new section in ROLES.md + four targeted edits to existing skill files + schema verification (no anticipated change). Per `release/SKILL.md §3.8` configuration-floor clause, γ-axis is capped at A- (γ/δ separation absent).

## Risks and forecasts

- **R1: Concurrent edits on `cdd/epsilon/SKILL.md` from cycle/400 (Phase 5).** δ contract notes possible conflict with cycle/400. At scaffold time, no `cycle/400` branch exists on `origin/main` (last main commit is `0d79c971` Merge cycle/397). If a cycle/400 surface lands before this merge, three-way merge or rebase resolves; the scopes are non-overlapping (Phase 5 is gamma/SKILL.md territory per issue Non-goals).
- **R2: ROLES.md numbering conflict.** Current §5 = "cdd as reference instantiation". δ contract Option A says "ROLES.md §5"; literal §5 is taken. Resolution: insert generic ε as **§4b** (between §4a five-layer-chain and §5 cdd-reference). Section number choice does not affect the AC1 oracle ("file/section exists; declares ε generically"), but does affect the citation in Phase 7. Recorded in design-notes.
- **R3: cdr/epsilon (post-#395) extension content might not fit generic surface.** cdr/epsilon contains CDR-specific trigger classes (`cdr-data-gate-gap`, `cdr-overclaim-gap`, etc.) which are correctly CDR-specific. The generic surface declares the **pattern** (each instantiation defines its own gap class names per its loss function); cdr's specific names are not generic doctrine. Cross-protocol verification confirms the pattern works.
- **R4: Empty-findings cycles under new cadence.** Many existing cdd-iteration.md files (e.g. #347, #393, #394, #395 has count > 0, #396) have empty findings. Under new rule, those would not be required. Backward-compat statement: rule is "required when > 0", not "forbidden when = 0". Existing files stay valid; future cycles may skip the file when count = 0.

## Plan order

1. ✅ Read all source files (gamma scaffold authored).
2. Draft `design-notes.md` (α design phase): home choice rationale, generic ε content structure, cross-protocol verification, cadence rule, schema verification plan.
3. Implement file changes (α build phase): ROLES.md §4b, cdd/epsilon/SKILL.md shrink, cdr/epsilon/SKILL.md cross-ref update, post-release §5.6b, schema verification.
4. Self-coherence sweep (`self-coherence.md`).
5. β-collapsed review (`beta-review.md`) — AC1–AC7 mechanical + β-rigor on generic surface declaration and cdr post-#395 citation.
6. Close-outs (α, β, γ).
7. cdd-iteration.md per closure-gate.
8. Merge to main.
9. Comment on cnos#366; update `.cdd/iterations/INDEX.md`.
10. Push & attempt branch delete.

## CDD Trace

| Step | Artifact | Skills loaded | Decision |
|------|----------|---------------|----------|
| 0 Observe | issue body, source skills (epsilon, post-release, cdr/epsilon, ROLES.md, RECEIPT-VALIDATION.md §Q3) | — | Inputs read; home choice pre-resolved by §Q3 to ROLES.md |
| 1 Select | cnos#401 | — | Phase 6 selected per δ dispatch |
| 2 Branch | `cycle/401` | cdd | Branched from `origin/main` per CDD §4.2 |
| 3 Bootstrap | `.cdd/unreleased/401/` | cdd | Cycle dir created |
| 4 Gap | this file | — | Named: ε doctrine is CDD-specific; multi-protocol future (cdr, cds, cdw) requires generic surface above CDD; cdd-iteration cadence has four-way drift across activation/epsilon/gamma/post-release skills |
