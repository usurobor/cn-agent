# β Close-Out — cnos#514 — docs Pass 4D

**cycle-branch:** cycle/514
**reviewer:** β
**round:** R0 (converge)

---

## Review summary

β reviewed α's implementation of cnos#514 (docs Pass 4D: move reference/runtime/schema bundles to docs/reference/ and docs/architecture/) and issued a converge verdict at R0 with three C-class findings and no blocking findings. All 17 ACs were evaluated; 14 passed outright, 3 were partial (C-class only).

## Implementation assessment

Scope accuracy is high: all 101 changed files are within scope, covering the 10 target bundles, redirect stubs at old paths, active link repairs across the repo, build.yml, OCaml/test files, and Go comment files. AC coverage is complete — no AC was missed or skipped, though three had self-coherence documentation inaccuracies rather than implementation defects. The path-repair work is substantially complete: all non-README active documents used git mv (R96–R100 renames confirmed); README files followed the necessary stub-in-place pattern. STOP-class gate compliance is clean: build.yml shows exactly one line changed (AC8), and golden/snapshot/fixture files are untouched (AC11).

## Technical review notes

- **AC8 STOP-class gate clean:** build.yml diff contains exactly one changed line (line 223: old schema path → new schema path). No other job, matrix, dependency, or cache change.
- **Schema merge approach:** The alpha/schemas bundle and beta/schema bundle were merged into docs/reference/schemas/. beta/schema/README.md was renamed to DESIGN-LLM-SCHEMA-README.md at the destination, correctly avoiding a filename collision with DESIGN-LLM-SCHEMA.md. The resulting four files in docs/reference/schemas/ are clean with no collision.
- **OCaml doc-comment changes:** All five OCaml/dune changes (cn_workflow.ml, cn_packet.ml, cn_contract_test.ml, cn_traceability_test.ml, dune) are comment-only — doc-comment preambles and citation strings. No executable logic was modified. dune was not available in this environment; the constraint is symmetric with α's environment and the changes are structurally incapable of introducing a new build failure.
- **README stub pattern (C1):** Nine bundle root README files were handled as Add-at-new-location + Stub-in-place rather than git mv. This is pragmatically correct (a bare git mv would require a separate stub creation step anyway), but self-coherence's unqualified "all moves used git mv" claim overstates the case.
- **Go comment repairs (C2):** src/go/internal/cli/command.go and src/go/internal/pkg/pkg.go were both modified (comment-only path citation repairs). Self-coherence §AC17 stated "Go: no Go source files changed," which contradicts the §Required checks §8 record that correctly notes them as "doc comment repairs only." No functional concern.
- **WRITER-PACKAGE.md do-not-touch auto-allowed repair (AC12):** docs/alpha/design/WRITER-PACKAGE.md is in the do-not-touch zone per gamma-scaffold §2.3 and §4.2, but the single-line change (stale path citation repair) satisfies every condition of the auto-allowed trivial repair policy. No δ override was required.

## Process observations

Review was completed in one round with no environment blockers beyond the expected dune/opam unavailability. All STOP-class gates were independently verifiable via git diff and grep. The AC5/AC13 ref sweep required layered analysis (redirect stubs, frozen-snapshot contexts, historical plan docs) but yielded a clean classification with no unaddressed stale refs in active documents. Self-coherence accuracy issues (C1, C2, C3) were documentation-layer concerns only; the underlying implementation was sound.

## C-class findings summary

**C1 — AC2 README git mv claim overstated:** self-coherence states "all moves used git mv" without qualification. Nine bundle root README files were handled via the Add-at-new + Stub-in-place pattern (no rename events). Content is correctly preserved at new locations; the claim is a documentation-accuracy issue, not an implementation defect.

**C2 — AC17 Go claim inaccurate:** self-coherence §AC17 states "Go: no Go source files changed," but two Go files (src/go/internal/cli/command.go, src/go/internal/pkg/pkg.go) were modified. The §Required checks §8 entry correctly records them as "doc comment repairs only." The discrepancy is within the self-assessment text; the actual changes are comment-only with no functional concern.

**C3 — AC16 missing explicit sections:** self-coherence §Required checks provides 8 numbered sections but is missing explicitly named sections for "do-not-touch proof," "link-check result," and "honest β/δ repair record" — three of the nine categories required by the AC oracle. Substance is partially covered inline (do-not-touch addressed in AC12 text; β/δ repair record present in AC12; link-check absent as a separate check), but the checklist structure does not fully match AC16 requirements.

## Release notes (pre-merge)

cnos#514 ships docs Pass 4D: ten reference/runtime/schema documentation bundles relocated from docs/alpha/ and docs/beta/ to docs/reference/ and docs/architecture/. Redirect stubs are in place at all old bundle root paths. All active cross-repo links (README.md, OPERATOR.md, docs/, src/, schemas/, build.yml) are repaired to new paths. No behavior changes: build.yml has one path update, OCaml and Go changes are comment-only, golden/snapshot/fixture files are untouched.
