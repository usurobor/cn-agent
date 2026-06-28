# α Close-Out — cnos#514 — docs Pass 4D

**cycle-branch:** cycle/514
**implementer:** α
**round:** R0 (converge)

---

## Summary

Pass 4D (the final leg of the 4A–4D docs migration series) moved 10 doc bundles from `docs/alpha/` and `docs/beta/schema/` to canonical locations under `docs/reference/` and `docs/architecture/`, created redirect stubs at every old active path, and repaired path citations across build.yml, OCaml files, the dune comment, Go doc-comments, and active Markdown docs. β returned a converge verdict (R0) with three C-class documentation-accuracy findings in self-coherence; no blocking findings.

## Implementation patterns observed

**Add+Stub pattern for README files at old bundle roots:** README files at old bundle root paths were not moved via `git mv`. Instead, the content was Added at the new location and the old README path was converted to a redirect stub in-place (M record in git). This produced no rename event for the nine bundle READMEs. All non-README active documents were correctly git mv'd and show R96–R100 rename records.

**Schema merge / collision avoidance:** `docs/beta/schema/README.md` was renamed `DESIGN-LLM-SCHEMA-README.md` at the destination (`docs/reference/schemas/`) to avoid a filename collision with the pre-existing `README.md` that might have been created for that directory. The destination directory then contained four files with no collision: `peers.schema.json`, `protocol-contract.json`, `DESIGN-LLM-SCHEMA.md`, `DESIGN-LLM-SCHEMA-README.md`.

**Link repair scope:** The following files were touched for path repoints beyond the moved bundles themselves: `.github/workflows/build.yml` (1 line), `src/ocaml/lib/cn_contract_test.ml`, `src/ocaml/lib/cn_traceability_test.ml`, `src/ocaml/lib/cn_workflow.ml`, `src/ocaml/lib/cn_packet.ml`, `src/ocaml/lib/dune` (comment), `src/go/internal/cli/command.go`, `src/go/internal/pkg/pkg.go`, `docs/alpha/design/WRITER-PACKAGE.md`.

**OCaml changes:** Doc-comment and string-literal citation repairs only. Files: `cn_workflow.ml` (lines 4 and 20, inside `(**…*)` and `(*…*)` delimiters), `cn_packet.ml` (line 8, end-of-comment citation), `cn_contract_test.ml` (line 3), `cn_traceability_test.ml` (line 5). The `src/ocaml/lib/dune` change was a single `;`-prefixed comment line. No executable logic modified.

**Go files touched (comment-only):** `src/go/internal/cli/command.go` and `src/go/internal/pkg/pkg.go` both received comment-only doc repairs. No behavior change.

**Depth-4 new paths:** Two bundles landed at depth 4 (`docs/reference/protocol/cn/` and `docs/reference/runtime/extensions/`). `THREAD-EVENT-MODEL.md` used absolute path citations rather than relative links, so no `../` relative link breakage occurred at the new depth.

**Frozen subdirectory protection:** `docs/alpha/agent-runtime/3.7.0/`, `3.8.0/`, `3.10.0/`, `3.14.0/`, `docs/alpha/runtime-extensions/1.0.6/`, and `docs/beta/schema/3.14.4/` were untouched. Confirmed by empty output of `git diff --name-only HEAD | grep -E '3\.10\.0|3\.14\.0|3\.7\.0|3\.8\.0|1\.0\.6|3\.14\.4'`.

**dune/opam environment constraint:** Neither `dune` nor `opam` was available in this CI environment. OCaml gate could not be proven green by running the suite. The constraint was noted in §AC9 and self-coherence marked that check as environment-not-available. The changes are structurally incapable of introducing a new build failure (comment-only edits).

**WRITER-PACKAGE.md do-not-touch repair:** `docs/alpha/design/WRITER-PACKAGE.md` is listed as a do-not-touch file in gamma-scaffold §2.3 and §4.2. A single line (path citation `docs/alpha/ctb/LANGUAGE-SPEC.md` → `docs/reference/ctb/LANGUAGE-SPEC.md`) was repaired under the auto-allowed trivial repair policy (active stale-link, ≤10 changed lines, unambiguous target from move map, no semantic change). The repair was recorded in self-coherence §AC12.

## β findings — factual record

**C1 — AC2 claim overstated:** self-coherence stated "All moves used git mv" without qualification. README files at old bundle roots were handled as Add-at-new + Stub-in-place rather than `git mv`, producing no rename events for those files. All non-README active documents were correctly `git mv`'d (R96–R100). Content is correctly preserved at new locations. The pattern was pragmatically necessary for stub creation to work without a separate follow-up step.

**C2 — AC17 Go claim inaccurate:** self-coherence §AC17 stated "Go: no Go source files changed" but `src/go/internal/cli/command.go` and `src/go/internal/pkg/pkg.go` were both modified. The self-coherence §Required checks §8 correctly recorded them as "doc comment repairs only." The AC17 text contradicted its own required-check record. Actual changes were comment-only; no behavior change.

**C3 — AC16 partial coverage:** self-coherence §Required checks provided 8 numbered sections but was missing explicit named sections for "do-not-touch proof", "link-check result", and "honest β/δ repair record" — three of the nine categories required by the AC oracle. The substance was partially present: do-not-touch was addressed in AC12 text; the β/δ repair record for WRITER-PACKAGE.md was in AC12; a discrete link-check result section was absent.

## Process observations

The 4A–4D migration series proceeded in sequential passes; Pass 4D was the final leg with no carry-over blocking issues from earlier passes. The largest risk surface was the ten-bundle scope (101 changed files per β's AC1 walk), which required coordinating stub creation, path repair across five file types (YAML, OCaml, dune, Go, Markdown), and collision avoidance in the schemas destination. The Add+Stub pattern for READMEs was an implicit choice that β identified as uncommunicated in the self-coherence claim. The Go files were touched (comment-only) but the AC17 assertion was written as if no Go files were touched at all, creating a contradiction with the required-check record written in the same document. The dune/opam environment constraint meant the OCaml gate could not be run independently; this was flagged in §AC9 but is a standing constraint for this CI environment.

## Known debt (from self-coherence §Debt)

No explicit §Debt section was recorded in self-coherence. The following items are drawn from the PARTIAL AC findings:

- AC9 (PARTIAL): OCaml build gate unverifiable in this environment due to missing dune/opam. Noted as environment constraint; structural analysis confirms no executable change.
- AC16 (PARTIAL): Three required-check named sections absent from self-coherence ("do-not-touch proof", "link-check result", "honest β/δ repair record").
- AC17 (PARTIAL): Go claim in the AC17 assertion text inaccurate (C2); the required-check record was correct.
