# α close-out — cycle/406 (Sub 1 of cnos#403)

**Issue:** [cnos#406](https://github.com/usurobor/cnos/issues/406) — Bootstrap cnos.cds package skeleton + extraction map.
**Mode:** design-and-build; γ+α+β collapsed on δ. β-α collapse acknowledged.
**Rounds:** R1 APPROVED (no fix-rounds).
**ACs:** 7/7 PASS.

## Surface delivered

Five files created in the new `src/packages/cnos.cds/` package; six cycle artifacts authored under `.cdd/unreleased/406/`; one INDEX.md row to append.

### Package skeleton (D1)

1. **`src/packages/cnos.cds/cn.package.json`** — manifest verbatim from `cnos.cdr/cn.package.json` with `name: cnos.cds`. Schema `cn.package.v1`; version `0.1.0`; kind `package`; engines `cnos >= 3.81.0`.
2. **`src/packages/cnos.cds/README.md`** — 64 lines; mirrors `cnos.cdr/README.md` section shape (`# CDS — Coherence-Driven Software`, "What CDS Does", "Package Structure", "Forthcoming surfaces", "Quick Start", "Status", "License"). Declares CDS as software-development realization of CCNF; cites kernel docs (CDD.md, COHERENCE-CELL.md, COHERENCE-CELL-NORMAL-FORM.md) as pointer-only; cites `cnos.cdr` as sibling research realization; cites #388 architecture inheritance; cites #403 as tracker.
3. **`src/packages/cnos.cds/skills/cds/SKILL.md`** — 121 lines; mirrors `cnos.cdr/skills/cdr/SKILL.md` loader pattern. Frontmatter has the six required fields per AC1 (name, description, artifact_class, governing_question, triggers, calls). `calls:` lists `CDS.md` + 5 role overlays as advisory targets per the cdr-Sub-1 precedent. Sections: Load order, Rule, Role overlays, Cross-protocol relationship, Conflict rule (one rule beyond cdr's set: kernel/realization-conflict rule), v0.1 caveat (new section documenting the Sub-1-without-Sub-2 status and pointing readers at `docs/extraction-map.md`).
4. **`src/packages/cnos.cds/skills/cds/.gitkeep`** — placeholder until Sub 2 lands CDS.md.
5. **`src/packages/cnos.cds/docs/extraction-map.md`** — 275 lines; 89 |-prefixed table rows across 12 surface-group tables + §13 coverage-verification table + §14 open-questions. The load-bearing artifact of Sub 1; Subs 3–5 dispatch against it.

### Cycle artifacts

6. `.cdd/unreleased/406/gamma-scaffold.md` — γ scaffold.
7. `.cdd/unreleased/406/self-coherence.md` — α self-coherence pass.
8. `.cdd/unreleased/406/beta-review.md` — β-collapsed review (APPROVED R1).
9. `.cdd/unreleased/406/alpha-closeout.md` — this file.
10. `.cdd/unreleased/406/beta-closeout.md` — β-level findings (none).
11. `.cdd/unreleased/406/gamma-closeout.md` — γ-level closure summary.
12. `.cdd/unreleased/406/cdd-iteration.md` — courtesy empty-findings stub per cycle/401 rule (`protocol_gap_count: 0`).
13. `.cdd/iterations/INDEX.md` — row 406 appended per the cycle/401 courtesy convention.

## Design decisions recorded

**D1: Manifest is verbatim from cdr (name swap only).** #406 pins the manifest shape verbatim. α executed the contract literally. No deviation considered or made.

**D2: README section shape parallels cdr/README.md.** The section headings match cdr's (What X Does / Package Structure / Forthcoming surfaces / Quick Start / Status / License); content is CDS-specific (software-development realization, engineering loss function, sibling-to-cdr framing). The mirror discipline is the #406 implementation-contract directive "mirror cnos.cdr v0.1 shape exactly".

**D3: SKILL.md adds two doctrine-coherence improvements beyond cdr's shape.**

1. New `## v0.1 caveat` section documenting that CDS.md is forthcoming (Sub 2) and pointing readers at `docs/extraction-map.md` for interim navigation. This is necessary because cds Sub 1 (this cycle) explicitly defers CDS.md to Sub 2 per #406 D1; cdr Sub 1 did not need this because cdr/CDR.md was shipped in the same Sub 1.
2. New Conflict rule: "If CDS.md and cnos.cdd (CCNF kernel) disagree on the kernel grammar, cnos.cdd governs." This makes the kernel/realization-conflict rule explicit (cdr's set leaves it implicit). Consistent with the kernel-realization hierarchy declared in the Cross-protocol relationship section.

Both additions are documented in self-coherence.md §AC4 and beta-review.md §Obs-1/Obs-2.

**D4: Role-overlay directory naming uses `delta/` directly rather than cdr's `operator/`.** cdr's legacy `operator/` naming is per the cdr SKILL.md note that "the directory is named `operator/` to mirror the engineering doctrine's `cnos.cdd/skills/cdd/operator/SKILL.md` exemplar; the role itself is δ per CDR.md Field 4". cds ships the canonical role-letter naming (`delta/`) directly. This is consistent with the broader cycle-394+ trend toward role-letter naming and avoids inheriting cdr's legacy workaround.

**D5: Extraction-map adds §13 coverage-verification table and §14 open-questions.** §13 makes Sub-6's marker-sweep mechanical (read CDD.md marker → §13 lookup → destination table → verify migrated → remove marker). §14 records 5 migration-coordination questions that Subs 3–5 will face (`.cdd/`→`.cds/` rename; cross-repo/harness/operator SKILL location; release-effector location) but do not invalidate the per-row destinations. Both additions exceed the strict AC5 minimum; both are Sub-3-through-Sub-6-enabling.

**D6: 4 of the 14 CDD.md markers are folded into related tables (Inputs into §1, Roles into §2, Large-file into §12). §13 makes the folding explicit.** Coordination surfaces is its own §3 (not folded). The folding decision was editorial: Inputs is "what selection reads", which is selection-adjacent; Roles is "who runs the lifecycle", which is lifecycle-adjacent; Large-file is a small operational discipline rule that fits inside the §Mechanical / §Artifacts family. β verified the folding decisions are defensible.

## AC verification summary

| AC | Oracle | Result | Note |
|----|--------|--------|------|
| AC1 | files exist + README ≥ 50 + SKILL.md valid frontmatter | PASS | 4 files; 64 README lines; 6 frontmatter fields present |
| AC2 | `cn build --check` reports `cnos.cds: valid`; no `src/go/` diff | PASS | Discovery works unchanged |
| AC3 | README declares software-realization of CCNF; cites kernel; cites CDR sibling; cites #388 inheritance | PASS | All 4 sub-clauses verified |
| AC4 | SKILL.md mirrors cdr structural shape; `calls:` advisory targets per cdr precedent | PASS | All 6 cdr sections present; 2 additions documented as doctrine-coherence improvements |
| AC5 | extraction-map row count ≥ 12; destinations under cnos.cds/; subs name 3/4/5 | PASS | 89 rows; all destinations resolve; all sub assignments correct |
| AC6 | `git diff` on cdd empty; CDS.md absent; no role overlays | PASS | Hard rule satisfied; 0-byte cdd diff |
| AC7 | cdr cross-references to CDS now resolve to existing peer | PASS | Informational; reciprocating "X is not Y-with-different-words" loop closed |

## Cross-protocol verification result

**cdr cross-references work cleanly.** cdr/SKILL.md's "CDR is not CDS-with-different-words" statement now resolves to an existing peer package (cnos.cds) that itself reciprocates the distinction in its own Cross-protocol relationship section. The cross-reference loop is closed. cdr is not modified in this cycle (AC6 hard rule); the AC7 informational check verifies the citation is no longer a forward-reference to a non-existent package.

## Schema-change disposition

**Unchanged.** `schemas/cds/` already exists per #388 Phase 2.5. CDS.md (Sub 2) will cite `schemas/cds/receipt.cue` as the typed γ close-out surface; no schema-side migration is required in Subs 3–5; no schema change in this Sub 1.

## Backward compatibility verified

cnos.cdd is entirely untouched (AC6). All 14 "pending cds extraction" markers in CDD.md remain in place; Sub 6 will sweep them once Subs 3–5 land. Existing cross-references from `gamma/SKILL.md`, `alpha/SKILL.md`, `release/SKILL.md`, `post-release/SKILL.md`, `harness/SKILL.md`, `operator/SKILL.md`, `activation/SKILL.md` to `CDD.md §X` anchors continue to resolve to the named family in CDD.md §"Software-specific realization" — no breakage. Sub 6 re-points those citations at the new CDS surfaces once Subs 3–5 migrate the content.

## Commit shape (anticipated)

This cycle ships (anticipated) commits:

1. `α-406: bootstrap cnos.cds package skeleton + extraction map (Sub 1 of #403)` (already landed; 5 files in src/packages/cnos.cds/)
2. `β+γ+cdd-iteration-406: AC1–AC7 PASS; close-outs; courtesy empty-findings cdd-iteration` (next; 6 cycle-artifact files in .cdd/unreleased/406/)
3. `γ-406: INDEX.md row for cycle/406 (0 findings; courtesy)` (next; 1 row appended to .cdd/iterations/INDEX.md)

Identity rotation: α work as `alpha@cdd.cnos`; β review + γ scaffold/close-outs as `beta@cdd.cnos` and `gamma@cdd.cnos` respectively; merge by `operator@cdd.cnos` (post-push, operator authority).

## Round-1 work summary

- 1 round, R1 APPROVED.
- 0 binding findings.
- 6 non-blocking β observations (Obs-1 through Obs-6 in beta-review.md).
- Mechanical ratio: undefined (no β findings; ratio computed only when findings ≥ 10).
- Sub 2 unblocked (CDS.md authoring can dispatch against the extraction map's destination commitments).
- Subs 3–5 unblocked (per-surface tables with sub assignments name what each sub owns).
- Sub 6 enabled (extraction-map §13 makes marker-sweep mechanical).
- Sub 7 unblocked (cnos.cds package exists; the empirical-anchor doc has a placement home at `docs/empirical-anchor-cdd.md` named in README "Forthcoming surfaces").

## Receipt obligation

Per CDD §5.5b, the cycle's parent-facing artifact is the merged branch + closed issue. This cycle's α-level receipt:

- Branch: `cycle/406` (pushed to origin; merge pending operator action).
- Issue: cnos#406 (closes when cycle/406 merges to main with `Closes #406` keyword in merge commit).
- ε artifact: `cdd-iteration.md` courtesy empty-findings stub (`protocol_gap_count: 0`).
- INDEX.md row: appended per the cycle/401 courtesy convention.
