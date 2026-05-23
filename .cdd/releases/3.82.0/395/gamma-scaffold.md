<!-- sections: [issue, mode, surfaces, peer-enumeration, scope-boundary, ac-oracle, design-source-citations, diff-scope, dispatch-config, tier3-skills] -->
<!-- completed: [issue, mode, surfaces, peer-enumeration, scope-boundary, ac-oracle, design-source-citations, diff-scope, dispatch-config, tier3-skills] -->

# γ Scaffold — Cycle #395

**Date:** 2026-05-21
**Issue:** [#395](https://github.com/usurobor/cnos/issues/395) — Sub 3 (#376): CDR role overlays (α/β/γ/δ/ε SKILL.md files for cnos.cdr)
**Parent:** [#376](https://github.com/usurobor/cnos/issues/376) (CDR v0.1 master; Sub 3 owns the role-overlay surface)
**Branch:** `cycle/395`
**Base SHA:** `e531dba0` (origin/main HEAD at session start; ε for #392)
**γ identity:** gamma / gamma@cdd.cnos
**Dispatch config:** γ+α+β-collapsed-on-δ (single Claude Code session; per breadth-2026-05-12 wave manifest precedent, validated empirically across cycles 375/377/378/388/390/392). β-α-collapse acknowledged: this is a docs-only design-and-build cycle whose AC oracles are mechanical (`rg`, `wc`, section-presence, cross-reference resolution) against canonical sources.

---

## Issue

**Gap:** Sub 1 (#390) shipped `src/packages/cnos.cdr/skills/cdr/CDR.md` declaring the six-field instantiation contract for the CDR protocol. The contract is doctrinal (what CDR is); per-role procedural skill (how the α/β/γ/δ/ε actors operate against the contract) is not yet authored. Without operationally-usable per-role skill files, downstream dispatchers (Rho via cn-rho hub; cph or other research projects binding CDR) have no overlay to load when they instantiate a role for CDR work — they would have to read CDR.md and infer the role-level procedural shape from the doctrinal contract alone. This is the surface the engineering protocol (CDD) addresses via `cnos.cdd/skills/cdd/{alpha,beta,gamma,operator,epsilon}/SKILL.md`; the CDR protocol needs the symmetric surface.

**Goal:** Author the 5 role-overlay SKILL.md files at:
- `src/packages/cnos.cdr/skills/cdr/alpha/SKILL.md` — research α: produce claims, methods, datasets, analyses, reports
- `src/packages/cnos.cdr/skills/cdr/beta/SKILL.md` — research β: audit claim/evidence/data-provenance + falsifiability + reproduction
- `src/packages/cnos.cdr/skills/cdr/gamma/SKILL.md` — research γ: coordinate the research wave, close into the typed research receipt
- `src/packages/cnos.cdr/skills/cdr/operator/SKILL.md` — research δ: enforce the data-mounted gate; record gate verdicts (GO / REVISE / NO-GO / INDETERMINATE / BOUNDED-GO); open next wave on gate-transition
- `src/packages/cnos.cdr/skills/cdr/epsilon/SKILL.md` — research ε: observe research-protocol gaps (overclaiming, unreproducible numbers, weak citation, construct drift) and patch CDR

Each role file extends the corresponding `cnos.cdd/skills/cdd/<role>/SKILL.md` generic doctrine; the kernel grammar is inherited by reference, the discipline profile is research-loss-function-specific (per `ROLES.md §4a.2`).

Additionally, stage an update to the loader at `src/packages/cnos.cdr/skills/cdr/SKILL.md` (Sub 2's deliverable, not yet shipped) so that the 5 role overlays are referenced as sub-skills. If Sub 2 (#394) lands before Sub 3, the loader update is a follow-on edit; if Sub 3 lands first, the loader file is authored here as a minimal stub that Sub 2 will extend. The decision is recorded in §dispatch-config below.

**Priority:** P2 — gates research-class CDD-style dispatch (any actor playing a CDR role needs the per-role skill to load).

**Work-shape:** small-medium docs-only cycle. Five new role SKILL.md files + one loader SKILL.md (stub) + cycle evidence under `.cdd/unreleased/395/`. No code, no schemas modified.

**Mode:** design-and-build (γ+α+β-collapsed on δ). The design half maps each cnos.cdd generic role surface to the research-loss-function extension (what counts as α matter; what β audits; what γ closes into; what δ gates on; what ε reads). The build half authors the five files plus the loader edit.

---

## Surfaces γ expects α to touch

1. `src/packages/cnos.cdr/skills/cdr/alpha/SKILL.md` — **new** (~120-180 lines).
2. `src/packages/cnos.cdr/skills/cdr/beta/SKILL.md` — **new** (~120-180 lines).
3. `src/packages/cnos.cdr/skills/cdr/gamma/SKILL.md` — **new** (~120-180 lines).
4. `src/packages/cnos.cdr/skills/cdr/operator/SKILL.md` — **new** (~120-180 lines).
5. `src/packages/cnos.cdr/skills/cdr/epsilon/SKILL.md` — **new** (~80-120 lines; cnos.cdd's epsilon exemplar is short by design).
6. `src/packages/cnos.cdr/skills/cdr/SKILL.md` — **new** (stub loader, ~40-60 lines; Sub 2 will extend with package metadata). AC5 requires the loader to reference the 5 role overlays.
7. `.cdd/unreleased/395/gamma-scaffold.md` — this file.
8. `.cdd/unreleased/395/design-notes.md` — design-half record (per-role research-function-vs-software-function differential; persona/protocol/project boundary; cross-reference plan).
9. `.cdd/unreleased/395/self-coherence.md` — α self-coherence per AC1–AC6 mapping.
10. `.cdd/unreleased/395/beta-review.md` — β-collapsed review (mechanical AC re-check).
11. `.cdd/unreleased/395/{alpha,beta,gamma}-closeout.md` — close-out artifacts.
12. `.cdd/unreleased/395/cdd-iteration.md` — closure-gate per `post-release/SKILL.md §5.6b`.
13. `.cdd/iterations/INDEX.md` — append cycle 395 row.

**Out-of-scope surfaces (explicitly):**
- Lifecycle sub-skills (`issue/`, `design/`, `plan/`, `review/`, `release/`, `post-release/` for CDR) — deferred to cds emergence per issue body.
- Persona drafting for `cn-rho` — separate hub work.
- Project-specific gates in cph `.cdr/` — Sub 4 + cph repo.
- Empirical-anchor mapping (Sub 4 of #376).
- Modifying `CDR.md` (Sub 1's deliverable; cite, do not modify).
- Modifying `ROLES.md` (binding; cite, do not modify).
- Modifying `cnos.cdd/skills/cdd/<role>/SKILL.md` (exemplars; cite as parent doctrine, do not extend).
- Package metadata (`cn.package.json`, `README.md`) — Sub 2.

---

## Peer enumeration (§2.2a — grep evidence)

- "`src/packages/cnos.cdr/skills/cdr/CDR.md` exists" — confirmed: `ls src/packages/cnos.cdr/skills/cdr/CDR.md` returns the file (31,739 bytes). Sub 1 (#390) shipped it.
- "`src/packages/cnos.cdr/skills/cdr/{alpha,beta,gamma,operator,epsilon}/` do not yet exist" — confirmed: `ls src/packages/cnos.cdr/skills/cdr/` returns only `CDR.md`. All five role subdirs are absent.
- "`src/packages/cnos.cdr/skills/cdr/SKILL.md` loader does not yet exist" — confirmed: same `ls` returns only `CDR.md`. Sub 2 (#394) has not shipped.
- "`src/packages/cnos.cdr/cn.package.json` does not exist" — confirmed: `ls src/packages/cnos.cdr/` returns only `skills/`. Sub 2 has not shipped.
- "`src/packages/cnos.cdd/skills/cdd/{alpha,beta,gamma,operator,epsilon}/SKILL.md` exist as exemplars" — confirmed: all five files present and read in full. α (396 lines), β (196 lines), γ (705 lines), operator/δ (678 lines), epsilon (74 lines).
- "`ROLES.md §4a.2` declares the loss-function distinction" — confirmed: `grep -n "Loss-function distinction\|truth-preserving claim transmission" ROLES.md` returns 217 + 222.
- "`docs/gamma/essays/CCNF-AND-TYPED-TRUST.md §"Separate persona, protocol, and project"` exists" — confirmed: line 260 (per #390 scaffold).
- "Sub 2 (#394) is open and unshipped" — confirmed via `mcp__github__issue_read`: state `open`; no commits referencing #394 in `git log --oneline origin/main` since #392 merge.
- "INDEX.md has 18 rows including cycle 392 (the most recent ε-tracked cycle)" — confirmed: file ends at the 392 row.

No claim of "X exists" or "X does not exist" asserted without ls/grep evidence.

---

## Scope boundary

**In scope (Sub 3):** authoring the 5 role-overlay SKILL.md files plus a minimal loader stub at `cdr/SKILL.md` that references the 5 overlays. Cycle evidence + ε iteration artifact. INDEX row.

**Out of scope:**
- Lifecycle sub-skill files (issue/, design/, plan/, review/, release/, post-release/) for CDR — deferred to cds emergence per issue Non-goals.
- Persona content for `cn-rho` (CDR is layer 3; persona is layer 1; they are doctrinally separate per CDR.md §"Persona, Protocol, Project").
- Project-specific content for cph (layer 4; lives in `<project>/.cdr/`).
- Empirical-anchor mapping (Sub 4).
- Modifying CDR.md or ROLES.md.
- Modifying cnos.cdd role skills (they are cited as parent doctrine; never edited from cnos.cdr).
- Authoring `cn.package.json` or `README.md` (Sub 2).
- Authoring V validator dispatch for CDR (already shipped via `protocol_id` dispatch in `schemas/cdd/README.md`; cited, not modified).
- Authoring runtime mechanics (dispatch, polling, git/CI plumbing) — per cnos#376 AC4 / cnos#395 AC6.

**Sister-cycle awareness:**
- cnos#376 — parent master. Sub 3 close-out comment confirms role overlays shipped.
- cnos#390 — Sub 1; CDR.md contract. Cite, do not modify.
- cnos#394 — Sub 2; package skeleton. Possible conflict point on `skills/cdr/SKILL.md` loader (both Sub 2 and Sub 3 touch the same file). Strategy: this Sub authors a minimal loader stub naming the 5 overlays; if Sub 2 lands first, Sub 3 extends Sub 2's loader by appending the role-overlay references. If Sub 3 lands first, Sub 2's later cycle will extend the stub with package-level concerns (frontmatter polish, cross-reference to CDR.md, package metadata pointers). The cycle/395 → main merge integrates whatever Sub 2 has shipped at merge time.
- cnos#388 — architectural inheritance; cite only.

---

## AC oracle approach

α verifies each AC using the oracles embedded in the issue body. β-collapsed re-runs the same mechanical checks against branch HEAD.

**AC1 — Five role SKILL.md files exist, non-empty, frontmatter valid:**
- `for f in src/packages/cnos.cdr/skills/cdr/{alpha,beta,gamma,operator,epsilon}/SKILL.md; do test -s "$f" || echo "fail:$f"; done` returns no output.
- Each frontmatter parses (YAML, contains `name`, `description`, `artifact_class`, `governing_question`, `parent: cdr`, `triggers`, `scope`).
- File-size sanity: each > 2KB (substantive overlay, not stub).

**AC2 — Each role overlay extends the cnos.cdd generic doctrine:**
- `rg "cnos.cdd/skills/cdd/(alpha|beta|gamma|operator|epsilon)/SKILL.md" src/packages/cnos.cdr/skills/cdr/{alpha,beta,gamma,operator,epsilon}/SKILL.md` returns ≥5 hits (one per file, each citing its corresponding generic doctrine).
- Each file's body explicitly declares "this is a CDR-specific extension of the generic cnos.cdd <role> doctrine; the kernel grammar (the algorithm structure, the role-cell shape) is inherited by reference; only the discipline profile and the matter type diverge for the research loss function."

**AC3 — Research-loss-function language throughout:**
- `rg "release|deploy|tag" src/packages/cnos.cdr/skills/cdr/operator/SKILL.md` returns 0 hits in normative sections (matches only in cross-references to cnos.cdd or in disavowal context like "δ does NOT 'cut a release'").
- Each role file's core principle / algorithm uses research vocabulary: claims, evidence, data, methods, reproduction, falsifiability, citation, claim_status, transmissibility, gate verdicts (GO/REVISE/NO-GO/INDETERMINATE/BOUNDED-GO).
- α produces research matter (claims/methods/data/analyses/reports); β audits oracles (claim/evidence/falsifiability/reproduction/citation); γ closes into `#CDRReceipt`; δ gates on data-mountedness; ε reads receipt-stream for overclaim/unreproducibility/citation-drift/construct-drift.

**AC4 — Persona/protocol/project boundary:**
- `rg "I am Rho|my voice|my temperament|/opt/" src/packages/cnos.cdr/skills/cdr/` returns 0 hits.
- No persona-identity prose; no project-specific gate names (cph-local paths, cph-specific datasets); no project-specific gate names embedded.
- Each role file explicitly declares (or relies on CDR.md's declaration via cross-reference) that this is a protocol-layer overlay, not a persona or project artifact.

**AC5 — SKILL.md loader updated to reference role overlays:**
- `rg "alpha/SKILL.md|beta/SKILL.md|gamma/SKILL.md|operator/SKILL.md|epsilon/SKILL.md" src/packages/cnos.cdr/skills/cdr/SKILL.md` returns ≥5 hits.
- The loader stub (if Sub 2 has not shipped) names the 5 role files explicitly; if Sub 2 has shipped, this Sub appends the role-skills list to Sub 2's loader without breaking Sub 2's package-level cross-references.

**AC6 — No surface fusion:**
- `rg "polling|dispatch|claude -p|gh issue|cn dispatch" src/packages/cnos.cdr/skills/cdr/*/SKILL.md` returns 0 hits in normative sections. Cross-references to cnos.cdd (which does discuss dispatch) are permitted as path strings only.
- No role overlay file authors runtime mechanics (subprocess invocation, branch creation, CI integration, scripts/release.sh equivalent).
- No release-driver effection (no tag commands, no merge instructions, no branch-cleanup procedures).

---

## Design-source citations

The cycle's design half cites:

1. **`src/packages/cnos.cdr/skills/cdr/CDR.md` (Sub 1, #390)** — binding contract. Each role overlay maps directly to one or more CDR.md §Field N subsections. α overlays to Field 1 (matter type); β to Field 2 (review oracle); γ to Field 3 (γ close-out artifact); δ to Field 4 (δ cadence); ε to Field 5 (ε iteration cadence). Field 6 (actor collapse) applies across all five.

2. **`ROLES.md §4a.2`** — loss-function distinction. Each role overlay declares its function in research-loss-function terms (truth-preserving claim transmission under uncertainty; primary failure mode is overclaim). Not engineering language (no "ship under repairable feedback").

3. **`ROLES.md §4a`** — five-layer enforcement chain. Each role overlay names its layer position (layer 3 — protocol overlay) and refuses to author layer 1 (persona) or layer 4 (project) content.

4. **`ROLES.md §3`** — six-field instantiation contract. CDR.md instantiates §3 for CDR; role overlays operate against that instantiation.

5. **`docs/gamma/essays/CCNF-AND-TYPED-TRUST.md §"Separate persona, protocol, and project"`** — doctrinal source for the three-layer separation.

6. **`src/packages/cnos.cdd/skills/cdd/{alpha,beta,gamma,operator,epsilon}/SKILL.md`** — structural exemplars. Each CDR role overlay borrows the heading shape (Core Principle / Load Order / Algorithm / Rules / Embedded Kata or its analogue) without copying the engineering-loss-function content. The frontmatter shape mirrors cnos.cdd's; the body diverges per discipline.

7. **`schemas/cdr/receipt.cue`** — typed γ close-out surface. The γ overlay cites `#CDRReceipt` and its required fields (claim_refs, data_refs, method_refs, result_refs, claim_status, reproduction).

---

## Expected diff scope

| Surface | Expected delta |
|---|---|
| `src/packages/cnos.cdr/skills/cdr/alpha/SKILL.md` | new (~150 lines) |
| `src/packages/cnos.cdr/skills/cdr/beta/SKILL.md` | new (~150 lines) |
| `src/packages/cnos.cdr/skills/cdr/gamma/SKILL.md` | new (~150 lines) |
| `src/packages/cnos.cdr/skills/cdr/operator/SKILL.md` | new (~150 lines) |
| `src/packages/cnos.cdr/skills/cdr/epsilon/SKILL.md` | new (~100 lines) |
| `src/packages/cnos.cdr/skills/cdr/SKILL.md` | new (~50 lines; stub loader referencing the 5 overlays) |
| `.cdd/unreleased/395/gamma-scaffold.md` | new (this file; ~250 lines) |
| `.cdd/unreleased/395/design-notes.md` | new (~250-350 lines) |
| `.cdd/unreleased/395/self-coherence.md` | new (~150-200 lines) |
| `.cdd/unreleased/395/beta-review.md` | new (~100-150 lines) |
| `.cdd/unreleased/395/alpha-closeout.md` | new (~80 lines) |
| `.cdd/unreleased/395/beta-closeout.md` | new (~80 lines) |
| `.cdd/unreleased/395/gamma-closeout.md` | new (~120-150 lines) |
| `.cdd/unreleased/395/cdd-iteration.md` | new (~50-80 lines) |
| `.cdd/iterations/INDEX.md` | +1 row for cycle 395 |

Total: ~1900-2400 line net change across ~15 files. Docs-only-shaped. No code, no schemas modified.

---

## Dispatch configuration

**γ+α+β collapsed on δ.** Single Claude Code session. The β-α-collapse is acknowledged: α ≠ β within a session is structurally compromised but acceptable for **docs-only-shaped role-overlay-authoring class** because β oracles are mechanical:
- Frontmatter parse (5 files; YAML keys present).
- Cross-reference resolution (`rg "cnos.cdd/skills/cdd/<role>/SKILL.md"` ≥5).
- Forbidden-token absence (`rg "release|deploy|tag"` 0 in operator/SKILL.md normative; `rg "I am Rho|/opt/"` 0; `rg "polling|dispatch|claude -p|gh issue|cn dispatch"` 0).
- Loader cross-reference (`rg "alpha/SKILL.md|..."` ≥5 in cdr/SKILL.md).
- Read-through Sub-2 / Sub-4 inheritance check (do the overlays leave room for Sub 2's package metadata + Sub 4's empirical-anchor mapping without rewriting?).

No subjective β judgment that α cannot self-verify against mechanical oracles. This is the same class as cycle #390 (Sub 1 contract authoring) and validated by precedent.

**Sub 2 loader-file collision strategy:** the loader file `src/packages/cnos.cdr/skills/cdr/SKILL.md` does not yet exist on `origin/main`. This Sub authors a minimal stub naming the 5 overlays. If Sub 2 (#394) ships before this Sub merges, the merge integrates Sub 2's loader (assumed to be richer in package-level concerns) by appending the 5 role-overlay references. If Sub 3 ships first, Sub 2's later cycle will extend the stub. Tracked as known integration point in §scope-boundary above.

Pre-flight check result:
```
γ pre-dispatch check — 2026-05-21:
  origin/cycle/395 exists: YES (just created from origin/main HEAD e531dba0)
  .cdd/unreleased/395/gamma-scaffold.md exists locally: YES (this file)
  issue #395 is open: YES
  branch pre-flight: PASS (base SHA e531dba0 matches origin/main HEAD)
  peer enumeration: PASS — all referenced surfaces confirmed by grep/ls
  scope boundary: documented (Sub 3 only; Sub 2 loader collision strategy named)
  cross-repo intake: n/a (cph cited path-only; no prose embedding)
  issue quality gate: PASS (six ACs with mechanical oracles)
  dispatch config: γ+α+β-collapsed on δ; docs-only-shaped; mechanical β gates
  refusal conditions: CDR.md exists and is readable; no refusal triggered
  timeout budget: n/a (synchronous session)
```

---

## Tier 3 skills

Named explicitly:
- `src/packages/cnos.core/skills/skill/SKILL.md` — skill-authoring meta-skill (frontmatter discipline, governing-question pattern).
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — exemplar α role overlay (engineering); extend, do not copy.
- `src/packages/cnos.cdd/skills/cdd/beta/SKILL.md` — exemplar β role overlay (engineering); extend, do not copy.
- `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` — exemplar γ role overlay (engineering); extend, do not copy.
- `src/packages/cnos.cdd/skills/cdd/operator/SKILL.md` — exemplar δ role overlay (engineering); extend, do not copy.
- `src/packages/cnos.cdd/skills/cdd/epsilon/SKILL.md` — exemplar ε role overlay (engineering); extend, do not copy.
- `src/packages/cnos.cdd/skills/cdd/design/SKILL.md` — design-half discipline.
- `src/packages/cnos.cdd/skills/cdd/issue/proof/SKILL.md` — AC oracle authoring.

Tier 1 (read-only doctrinal):
- `src/packages/cnos.cdr/skills/cdr/CDR.md` (Sub 1 deliverable) — binding contract.
- `ROLES.md §4a.2` + `§4a` + `§3` — loss-function + chain + contract.
- `docs/gamma/essays/CCNF-AND-TYPED-TRUST.md §"Separate persona, protocol, and project"` — boundary source.
- `schemas/cdr/receipt.cue` — typed γ close-out surface.
