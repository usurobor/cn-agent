<!-- sections: [issue, mode, surfaces, peer-enumeration, scope-boundary, ac-oracle, design-source-citations, diff-scope, dispatch-config, tier3-skills] -->
<!-- completed: [issue, mode, surfaces, peer-enumeration, scope-boundary, ac-oracle, design-source-citations, diff-scope, dispatch-config, tier3-skills] -->

# γ Scaffold — Cycle #390

**Date:** 2026-05-21
**Issue:** [#390](https://github.com/usurobor/cnos/issues/390) — Sub 1 (#376): CDR six-field instantiation contract + architecture-choice declaration
**Parent:** [#376](https://github.com/usurobor/cnos/issues/376) (CDR v0.1 master; Sub 1 owns AC6 + AC7)
**Branch:** `cycle/390`
**Base SHA:** `417b6227` (origin/main HEAD at session start; ε for #388)
**γ identity:** gamma / gamma@cdd.cnos
**Dispatch config:** γ+α+β-collapsed-on-δ (single Claude Code session; per breadth-2026-05-12 wave manifest precedent, validated empirically across cycles 375/377/378/388). β-α-collapse acknowledged: this is a docs-only-shaped design-and-build cycle whose AC oracles are mechanical (`rg`, `wc`, section-presence) against canonical sources.

---

## Issue

**Gap:** `cnos.cdr` does not yet exist as a protocol-skill package. `ROLES.md §3` mandates the six-field instantiation contract every c-d-X protocol must declare. cnos#388 (Phase 2.5) decided option (a) split (common kernel + per-protocol overlay) for schemas; the same architectural reading applies to skills (cnos#376 AC7), but no document records that inheritance for the CDR skill bundle. Without (i) the six-field contract declared in CDR-shape research-loss-function language and (ii) the architectural-choice section recording (a)-inheritance, downstream Sub 2 (package skeleton), Sub 3 (role overlays), and Sub 4 (empirical-anchor doc) cannot author against a stable shape.

**Goal:** Ship `src/packages/cnos.cdr/skills/cdr/CDR.md` declaring:
1. The six instantiation-contract fields (matter type, review oracle, γ close-out artifact, δ cadence, ε iteration cadence, actor collapse rule) — each in research-loss-function language per `ROLES.md §4a.2`.
2. An `## Architecture choice` section recording option (a) inheritance from cnos#388 + the (b)-rejected rationale (1)–(5) by reference.
3. A `## Persona, Protocol, Project` section naming the three doctrinal layers with their canonical homes.
4. An `## Empirical anchor` section citing `usurobor/cph` and declaring shape-compatibility (detailed mapping deferred to Sub 4).

**Priority:** P2 — gates Sub 2/3/4 of #376.

**Work-shape:** small-medium docs-only cycle. Single new file (`src/packages/cnos.cdr/skills/cdr/CDR.md`) plus cycle evidence under `.cdd/unreleased/390/`. The package directory exists; package metadata (`cn.package.json`, package-level `README.md`, `SKILL.md`) is **not** in scope (Sub 2).

**Mode:** design-and-build (γ+α+β collapsed). The design half consolidates §3 mapping to research-loss-function language and the (a)-inheritance rationale; the build half authors the CDR.md file.

---

## Surfaces γ expects α to touch

1. `src/packages/cnos.cdr/skills/cdr/CDR.md` — **new** (single deliverable; ~250-400 lines).
2. `.cdd/unreleased/390/gamma-scaffold.md` — this file.
3. `.cdd/unreleased/390/design-notes.md` — design-half record (six-field shape draft + architectural-choice rationale + persona/protocol/project boundary structure + empirical-anchor citation strategy).
4. `.cdd/unreleased/390/self-coherence.md` — α self-coherence per AC1–AC6 oracle mapping.
5. `.cdd/unreleased/390/beta-review.md` — β-collapsed review (mechanical AC re-check).
6. `.cdd/unreleased/390/{alpha,beta,gamma}-closeout.md` — close-out artifacts.
7. `.cdd/unreleased/390/cdd-iteration.md` — closure-gate (per `post-release/SKILL.md §5.6b`, empty findings still produces the file).
8. `.cdd/iterations/INDEX.md` — append cycle 390 row.

**Out-of-scope surfaces (explicitly):**
- `src/packages/cnos.cdr/cn.package.json` — Sub 2.
- `src/packages/cnos.cdr/README.md` — Sub 2.
- `src/packages/cnos.cdr/skills/cdr/SKILL.md` — Sub 2.
- `src/packages/cnos.cdr/skills/cdr/{alpha,beta,gamma,delta,epsilon}/SKILL.md` — Sub 3.
- A full mapping of `cph`'s `.cdr/` artifact set to CDR.md fields — Sub 4.
- Persona content for `cn-rho` (staged at `.cdd/iterations/cross-repo/cn-rho/`).

---

## Peer enumeration (§2.2a — grep evidence)

- "`ROLES.md §3` declares the six-field instantiation contract" — confirmed: `grep -n "Field 1 — Matter type\|Field 6 — Actor collapse rule" ROLES.md` returns lines 105 + 146 (six bolded Field N entries at lines 105, 114, 123, 131, 138, 146).
- "`ROLES.md §4a.2` distinguishes engineering (CDS) from research (CDR) loss functions" — confirmed: `grep -n "Loss-function distinction\|truth-preserving claim transmission" ROLES.md` returns lines 217 + 222.
- "`ROLES.md §4a.3` sketches the CDR receipt field set" — confirmed: `grep -n "claim_refs\|data_refs\|method_refs\|result_refs\|claim_status\|limitations\|reproduction" ROLES.md` returns matches at lines 244-252 + 261.
- "`schemas/cdr/receipt.cue` carries the CDR receipt skeleton" — confirmed: file exists (75 lines); declares `#CDRReceipt: cdd.#Receipt & {...}` with `protocol_id: "cnos.cdd.cdr.receipt.v1"`, four required evidence-ref lists (claim/data/method/result with min-length-one), `#ClaimStatus` enum, optional limitations + `#Reproduction`.
- "`schemas/cdd/README.md` records the (a) decision and cross-references cnos#376 AC7" — confirmed: `grep -n "Architectural choice\|cnos#376 AC7" schemas/cdd/README.md` returns lines 59 + 165.
- "`docs/gamma/essays/CCNF-AND-TYPED-TRUST.md` §'Separate persona, protocol, and project' exists" — confirmed: line 260.
- "essay §'Wave 4 — CDR bootstrap' exists" — confirmed: line 405 (`Create cnos.cdr as the research overlay, anchored in cph but not defined by cph. CDR owns research role overlays. Rho owns the research persona. cph owns the project binding.`).
- "`src/packages/cnos.cdr/` does not yet exist as a populated package" — confirmed: directory is empty (just-created `skills/cdr/` subdir; no `cn.package.json`, no `README.md`, no SKILL.md files).
- "`src/packages/cnos.cdd/skills/cdd/CDD.md` is the reference instantiation" — confirmed: file exists (1344 lines); structural exemplar for the persona/protocol/project + six-field shape.
- "cph empirical anchor exists at `usurobor/cph`" — referenced in #376 + #388 design surfaces; cph repo not present in this worktree, treated as design-source by citation (path-only, no prose embedding).

No claim of "X does not exist" asserted without ls/grep evidence.

---

## Scope boundary

**In scope (Sub 1):** authoring `CDR.md` with §Architecture choice, §Persona/Protocol/Project, §Six-field contract (six numbered subsections), §Empirical anchor. Creating the `src/packages/cnos.cdr/skills/cdr/` directory if not present. Cross-references to ROLES, schemas/cdr/, cnos#388, cnos#376, the essay.

**Out of scope:**
- Package skeleton (`cn.package.json`, package `README.md`, root `SKILL.md`) — Sub 2 of #376.
- Role-overlay skills (`alpha/SKILL.md` … `epsilon/SKILL.md` for CDR) — Sub 3 of #376.
- Empirical-anchor doc mapping cph's `.cdr/` artifact set — Sub 4 of #376.
- Persona drafting for `cn-rho` — separate hub work.
- Editing `ROLES.md` — CDR instantiates; ROLES governs.
- Editing `schemas/cdr/receipt.cue` — already shipped (#388); cite, do not modify.
- V implementation / `cn-cdd-verify` changes — Phase 3 of #366.

**Sister-cycle awareness:**
- cnos#376 — parent master. Sub 1 close-out comment confirms Sub 1 shipped, AC6 + AC7 met, Sub 2/3/4 unblocked.
- cnos#388 — architectural choice upstream; cite as decision-record source.
- cnos#366 — Phase 2.5 + Phase 3 context; not directly affected by this sub.

---

## AC oracle approach

α verifies each AC using the oracles embedded in the issue body. β re-runs the same mechanical checks against the branch HEAD.

**AC1 — Six fields declared per ROLES.md §3:**
- `rg "^### Field" src/packages/cnos.cdr/skills/cdr/CDR.md | wc -l` returns 6.
- `rg -c "TBD" src/packages/cnos.cdr/skills/cdr/CDR.md` returns 0.
- Each `### Field N` subsection uses research-loss-function language: claim/evidence/data/method/reproduction/falsifiability for AC2 (review oracle); claim_refs/data_refs/method_refs/result_refs/claim_status/limitations/reproduction for AC3 (γ close-out); gate-transition-shaped (not release/tag/deploy) for AC4 (δ cadence); receipt-stream review over protocol gaps for AC5 (ε cadence); α≠β always for AC6 (actor collapse).
- `rg "release|deploy|tag" src/packages/cnos.cdr/skills/cdr/CDR.md` returns matches only in cross-references or "what CDR is **not**" disclaimers — not in normative field declarations.

**AC2 — Architecture-choice section declared:**
- `grep -n "^## Architecture choice" src/packages/cnos.cdr/skills/cdr/CDR.md` returns ≥1 hit.
- Section names option (a) ("common constitution + per-protocol procedures") as chosen; names option (b) ("common protocol-agnostic skill + domain overlay" or equivalent) as rejected.
- Cross-references cnos#388 + `schemas/cdd/README.md §"Architectural choice"` for the upstream decision.
- Rationale cites: types-first, clarity at protocol boundary, mechanical generic-vs-domain separation, generalizes to future c-d-X, decision-once-applied-twice.

**AC3 — Persona/protocol/project boundary section declared:**
- `grep -n "^## Persona, Protocol, Project" src/packages/cnos.cdr/skills/cdr/CDR.md` returns ≥1 hit.
- All three layers named: persona (canonical home `cn-rho/spec/`), protocol overlay (canonical home `cnos.cdr/skills/cdr/`), project binding (canonical home `<project>/.cdr/`, exemplar cph).
- States: "Rho may play δ but is not CDR"; "CDR defines α/β/γ/δ/ε overlays independent of any persona hub"; "a project binds those roles via its own `.cdr/` directory."

**AC4 — Empirical anchor validation:**
- `grep -ni "cph\|usurobor/cph" src/packages/cnos.cdr/skills/cdr/CDR.md` returns ≥1 hit.
- Section declares the six-field shape compatible with cph's CDR practice (PROJECT.md, `.cdr/waves/`, field reports, dataset docs, gate verdicts including bounded GO / partial GO / INDETERMINATE).
- Defers detailed mapping to Sub 4 with explicit pointer.

**AC5 — No surface fusion:**
- `rg "polling|dispatch|claude -p|gh issue" src/packages/cnos.cdr/skills/cdr/CDR.md` returns 0 hits in normative sections (matches only in "out of scope" / "Sub 2/3 owns this" disclaimers, if any).
- Persona-identity prose absent: `rg -i "I am Rho|my voice|as Rho" src/packages/cnos.cdr/skills/cdr/CDR.md` returns 0 hits.
- Project-specific gates not authored (cph-local gate names should not appear in CDR.md normative text).

**AC6 — Sub 2/3/4 inheritance ready:**
- Sub 2 scope (package skeleton: `cn.package.json`, `README.md`, `SKILL.md`) can be authored citing CDR.md as the doctrinal anchor without re-deriving the six fields.
- Sub 3 scope (role overlays) can author each role's `SKILL.md` citing CDR.md's `### Field N` subsections for its review oracle, γ close-out, etc.
- Sub 4 scope (empirical-anchor doc) can map cph artifacts to the §Empirical anchor section's named field set.
- Read-through check: walk each Sub's expected scope mentally against CDR.md; no re-derivation surfaces.

---

## Design-source citations

The cycle's design half cites:

1. `ROLES.md §3` (lines 98-156) — six-field instantiation contract. Mandatory shape. CDR.md mirrors the six-field structure precisely.

2. `ROLES.md §4a` (lines 177-263) — five-layer enforcement chain. Persona/protocol/project boundary section in CDR.md instantiates layers 1–4 (CDR is layer 3; receipts are layer 5 and are typed by `schemas/cdr/`).

3. `ROLES.md §4a.2` (lines 217-227) — loss-function distinction. CDR optimises for truth-preserving claim transmission under uncertainty; primary failure mode is overclaim. Each field declaration in CDR.md uses this loss function as its discipline anchor.

4. `ROLES.md §4a.3` (lines 228-254) — receipt-enforces-discipline; the CDR receipt field set sketch (`claim_refs`, `data_refs`, `method_refs`, `result_refs`, `claim_status`, `limitations`, `reproduction`). CDR.md's Field 3 (γ close-out) cites `schemas/cdr/receipt.cue` as the typed surface; the field-set descriptor mirrors §4a.3.

5. `ROLES.md §7` (lines 309-335) — naming convention for new c-d-X protocols. `cdr` letter is reserved for research (investigation, synthesis, citation).

6. `docs/gamma/essays/CCNF-AND-TYPED-TRUST.md` §"Separate persona, protocol, and project" (lines 260-283) — the doctrinal source for the persona/protocol/project chain in CDR.md. Sigma/CDS vs Rho/CDR loss-function distinction inherited.

7. `docs/gamma/essays/CCNF-AND-TYPED-TRUST.md` §"Wave 4 — CDR bootstrap" (lines 405-407) — declares the wave-shape: "Create `cnos.cdr` as the research overlay, anchored in cph but not defined by cph. CDR owns research role overlays. Rho owns the research persona. cph owns the project binding." CDR.md realises this declaration.

8. `schemas/cdd/README.md §"Architectural choice"` (lines 59-167) — the decision-record from cnos#388 (option (a) for schemas; rationale 1-5). CDR.md's `## Architecture choice` cites this section as the upstream decision; rationale is inherited by reference, not re-derived.

9. `schemas/cdr/receipt.cue` (75 lines) — the typed CDR receipt skeleton (`#CDRReceipt` definition). CDR.md's Field 3 (γ close-out artifact) names this file as the typed surface.

10. `src/packages/cnos.cdd/skills/cdd/CDD.md` — the reference instantiation; structural exemplar. CDR.md borrows the heading shape (top preamble + numbered sections) without copying software-protocol language.

---

## Expected diff scope

| Surface | Expected delta |
|---|---|
| `src/packages/cnos.cdr/skills/cdr/CDR.md` | new (~280-380 lines): preamble + Architecture choice + Persona/Protocol/Project + Six-field contract (six `### Field N` subsections) + Empirical anchor + cross-references |
| `.cdd/unreleased/390/gamma-scaffold.md` | new (this file; ~240 lines) |
| `.cdd/unreleased/390/design-notes.md` | new (~250-350 lines): L7-style design record (architectural choice statement; six-field rough prose; persona/protocol/project boundary structure; empirical-anchor citation strategy) |
| `.cdd/unreleased/390/self-coherence.md` | new (~150-200 lines): gap, mode, AC1–AC6 mapping, CDD Trace, review-readiness |
| `.cdd/unreleased/390/beta-review.md` | new (~100-150 lines): β-α-collapse acknowledgment, AC1–AC6 mechanical re-check, APPROVE verdict |
| `.cdd/unreleased/390/alpha-closeout.md` | new (~80 lines): cycle summary, findings (factual only), debt |
| `.cdd/unreleased/390/beta-closeout.md` | new (~80 lines): merge evidence, review evidence, β-side findings |
| `.cdd/unreleased/390/gamma-closeout.md` | new (~120-150 lines): triage, post-merge cross-references, cnos#376 close-out comment plan |
| `.cdd/unreleased/390/cdd-iteration.md` | new (~30-60 lines): closure-gate per `post-release/SKILL.md §5.6b`; findings (likely empty or 1-2 minor) |
| `.cdd/iterations/INDEX.md` | +1 row for cycle 390 |

Total: ~1300-1700 line net change across ~10 files. Docs-only-shaped. No code, no schemas modified.

---

## Dispatch configuration

**γ+α+β collapsed on δ.** Single Claude Code session. The β-α-collapse is acknowledged: α ≠ β within a session is structurally compromised but acceptable for **docs-only-shaped contract-authoring class** because β oracles are mechanical:
- Section-presence (`grep -n "^## "` for required headings).
- Field count (`rg "^### Field" CDR.md | wc -l`).
- Forbidden-token absence (`rg "polling|dispatch|TBD|release|deploy"` returns ∅ in normative sections).
- Cross-reference resolution (cited paths exist on the branch).
- Read-through inheritance check (Sub 2/3/4 mental walk).

No subjective judgment β could provide beyond α's verification.

Pre-flight check result:
```
γ pre-dispatch check — 2026-05-21:
  origin/cycle/390 exists: YES (just created from origin/main HEAD 417b6227)
  .cdd/unreleased/390/gamma-scaffold.md exists locally: YES (this file)
  issue #390 is open: YES
  branch pre-flight: PASS (base SHA 417b6227 matches origin/main HEAD)
  peer enumeration: PASS — all referenced surfaces confirmed by grep/ls
  scope boundary: documented (Sub 1 only; Sub 2/3/4 deferred; Sub-Sub-non-goals enumerated)
  cross-repo intake: n/a (cph cited path-only; no prose embedding)
  issue quality gate: PASS (six ACs with oracles, architectural inheritance specified)
  dispatch config: γ+α+β-collapsed on δ; docs-only-shaped; mechanical β gates
  refusal conditions: none triggered at scaffold time
  timeout budget: n/a (synchronous session)
```

---

## Tier 3 skills

Named explicitly:
- `src/packages/cnos.cdd/skills/cdd/design/SKILL.md` — design-half discipline; architectural-choice authoring shape.
- `src/packages/cnos.cdd/skills/cdd/issue/contract/SKILL.md` — contract section authoring (the six-field shape is structurally a contract).
- `src/packages/cnos.cdd/skills/cdd/issue/proof/SKILL.md` — AC oracle authoring; β re-checks the issue ACs.
- `src/packages/cnos.cdd/skills/cdd/post-release/SKILL.md` §5.6b — `cdd-iteration.md` cadence; closure-gate.

Tier 1 (read-only doctrinal):
- `ROLES.md §3` + `§4a` + `§4a.2` + `§4a.3` + `§7` — instantiation grammar; binding.
- `docs/gamma/essays/CCNF-AND-TYPED-TRUST.md` §"Separate persona, protocol, and project" + §"Wave 4 — CDR bootstrap" — doctrinal source for the boundary section.
- `schemas/cdd/README.md §"Architectural choice"` — decision-record for the (a) inheritance.
- `schemas/cdr/receipt.cue` — typed γ close-out surface.
- `src/packages/cnos.cdd/skills/cdd/CDD.md` — structural exemplar (heading shape only; do not copy software-protocol verbs).
