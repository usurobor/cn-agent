# β Close-out — Cycle #395 (Sub 3 of #376)

**Cycle:** #395 — Sub 3 (#376): CDR role overlays
**Branch:** `cycle/395`
**Merge:** to be recorded post-merge by γ in `gamma-closeout.md`
**Date:** 2026-05-21
**β session:** β-collapsed-on-α (engineering-class collapse acknowledged in `beta-review.md §"β-α-collapse acknowledgement"`)
**Rounds:** R1 APPROVED — no fix-round

## §Review Summary

Reviewed `cycle/395` against AC1–AC6 by mechanical oracle (`rg`, `wc`, frontmatter parse) plus per-AC hit-by-hit classification for the forbidden-token oracles (AC3, AC4, AC6). All ACs PASS. No D or C findings raised.

## §Implementation Assessment

The 5 role overlay files (alpha/, beta/, gamma/, operator/, epsilon/) and the loader stub (cdr/SKILL.md) implement the issue's scope cleanly:
- Each overlay's preamble blockquote names "this is a CDR-specific extension of the generic cnos.cdd <role> doctrine" with the kernel-grammar-inherited-by-reference framing — making the AC2 "extends not replaces" property structurally visible in every file.
- Each overlay's "Persona / protocol / project boundary" section declares the layer-3 (protocol) position and refuses layer-1 (persona) and layer-4 (project) content — making AC4 structurally enforced not just textually absent.
- The operator/SKILL.md's "What research δ does NOT do" subsection mirrors `cnos.cdd/skills/cdd/operator/SKILL.md §6` structure and is the canonical home for the AC3 disavowal text — making AC3's intent ("no release/deploy/tag as δ verbs") readable in the role boundary itself, not just inferable from absence.
- The gamma/SKILL.md's "Selection inputs (research-specific)" and "Wave-issue quality gate (research-specific)" subsections replace the engineering γ's CHANGELOG-TSC-table / encoding-lag inputs with research-class inputs (open-claim ledger, construct-stability surface, measurement-disagreement surface, citation-debt surface), preserving the engineering γ's structural shape while diverging at the discipline profile.
- The epsilon/SKILL.md mirrors cnos.cdd ε's short-form structure (74 lines → 70 lines) and uses the same wording almost verbatim for §2 (the ε=δ collapse rule); only §1 (trigger classes) and the artifact-naming change, which is exactly the right surface for the divergence (per ROLES.md §4: "Independence, not headcount").

The loader stub (`cdr/SKILL.md`) is appropriately minimal: it satisfies AC5 (15 hits across the 5 role names) without overstepping into Sub 2's package-metadata scope. The preamble blockquote names the integration plan with Sub 2 explicitly.

## §Technical Review

**Surfaces audited:**
- 5 × role overlay SKILL.md (alpha, beta, gamma, operator, epsilon)
- 1 × loader SKILL.md (stub)
- 4 × cycle artifacts (γ-scaffold, design-notes, self-coherence, β-review)

**Discipline profile vs. matter type discrimination:** each overlay correctly identifies what changes (matter type, oracle, close-out artifact, cadence trigger) from the engineering exemplar vs. what is inherited unchanged (the algorithm structure, role-cell shape, independence rules, resumption protocol). The α=β prohibition is correctly carried forward to research with the engineering-class collapse explicitly noted as not-transferring.

**Cross-reference integrity:** all paths in cross-references resolve. `../CDR.md`, `../../../../cnos.cdd/skills/cdd/<role>/SKILL.md`, `../../../../../../ROLES.md`, `../../../../../../schemas/cdr/receipt.cue` — each path navigates from `src/packages/cnos.cdr/skills/cdr/<role>/SKILL.md` to the correct file.

**Frontmatter discipline:** each role file has the required keys; epsilon's frontmatter mirrors cnos.cdd ε's compact form (no `inputs`, `outputs`, `requires`, `calls` — those are loader concerns, not role-overlay concerns). The loader stub does carry `inputs`, `outputs`, `requires`, `calls` per the cnos.cdd loader exemplar.

## §Process Observations

**Single-session γ+α+β-collapsed-on-δ for c-d-X bootstrap continues to validate.** Sub 3 ships cleanly under the same dispatch pattern as cycles #390 (Sub 1), #388, #392, etc. The pattern is mature; the mechanical AC oracles + classification tables substitute reliably for separate-session β when the matter is docs-only protocol-document authoring.

**AC oracle disavowal-context pattern is now a 2-cycle empirical norm in the cnos.cdr family.** #390 surfaced it for `release|deploy|tag`; #395 surfaces it for the same plus `I am Rho|my voice` and `dispatch|polling`. ε iteration may want to file a structural carve-out for the AC template; current cycle records it as same-class-as #390 known pattern.

**No release-artifact authoring.** This Sub does not touch RELEASE.md, version files, or schemas. γ's PRA + release prep mechanics (per `cnos.cdd/skills/cdd/gamma/SKILL.md §2.6`) are not in scope for an in-flight cycle that has not yet hit release boundary.

## §Release Notes

Not applicable to in-cycle close-out. The cycle dir will move to `.cdd/releases/{X.Y.Z}/395/` at next release per `release/SKILL.md §2.5a`.

---

Filed by β on 2026-05-21.
