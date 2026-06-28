# beta-review — cnos#514

## §R0

### Verdict

converge

### Findings

No blocking findings. Three minor (C) documentation-accuracy findings in self-coherence:

**C1 — AC2 claim overstated:** README files at old bundle roots were not moved via `git mv`. The pattern used was Add-at-new-location + Stub-in-place (M at old path, A at new path). All non-README active documents were correctly git mv'd (R96–R100). The content IS correctly preserved at new locations; only the rename-detection chain for README files is absent. This is pragmatically necessary for the stub pattern to work (a bare `git mv README.md` would need a follow-up stub creation anyway), but the self-coherence claims "All moves used git mv" without qualification. Nine bundles have README-at-old-path stubs and Added READMEs at new paths.

**C2 — AC17 Go claim inaccurate:** self-coherence §AC17 states "Go: no Go source files changed" but `src/go/internal/cli/command.go` and `src/go/internal/pkg/pkg.go` were both modified. The §Required checks §8 correctly records them as "doc comment repairs only" — the self-assessment text at AC17 contradicts its own required-check record. The actual changes are comment-only (no behavior change). No functional concern.

**C3 — AC16 partial coverage:** self-coherence §Required checks provides 8 numbered sections but is missing explicit named sections for: "do-not-touch proof", "link-check result", and "honest β/δ repair record" — three of the nine categories required by the AC oracle. The substance is partially covered: do-not-touch is addressed in the AC12 text; the β/δ repair record for WRITER-PACKAGE.md is in AC12; link-check is absent as a separate check. Coverage is sufficient to understand the implementation but does not fully match the AC16 checklist.

### AC walk

- **AC1 PASS** — All 101 changed files are within scope: 10-bundle doc paths, stubs at old paths, active link repairs across the repo (README.md, OPERATOR.md, docs/, src/, schemas/), build.yml, OCaml/test files, dune file, CDD skill files, and Go comment files. No out-of-scope content changes observed.
- **AC2 PARTIAL** — Non-README active documents: all confirmed as R96–R100 renames (git mv). README files at old bundle roots: handled as Add-at-new + Stub-in-place (no rename events). Content correctly preserved at new locations. See C1 finding above.
- **AC3 PASS** — Stubs confirmed at all 10 old bundle root paths: protocol, agent-runtime, runtime-extensions, package-system, cli, ctb, schemas (alpha), schema (beta), security, cognitive-substrate. All stubs cite new locations correctly. runtime-extensions and beta/schema stubs correctly note frozen snapshots remain.
- **AC4 PASS** — `git diff --name-only main...HEAD | grep -E '3\.10\.0|3\.14\.0|3\.7\.0|3\.8\.0|1\.0\.6|3\.14\.4'` returns empty. No `.cdd/` files outside `unreleased/514/` modified.
- **AC5/AC13 PASS** — Absolute sweep: remaining refs to old paths are all classified as redirect stubs, frozen-snapshot contexts (3.14.4/, code blocks in DOCUMENTATION-SYSTEM.md), or historical plan docs (PLAN-v3.13.0, CHANGELOG.md, docs/gamma/). Relative sweep (scoped to the 10 bundle paths): returns empty for active references. Broader relative sweep returns only frozen-snapshot contexts (docs/beta/architecture/3.14.4/) and historical plans pointing to frozen subdirs (3.8.0/, 3.10.0/).
- **AC6 PASS** — build.yml:223, cn_contract_test.ml:3, cn_traceability_test.ml:5, cn_workflow.ml:4+20, cn_packet.ml:8 all show new paths. Runtime fixture strings ("protocol-contract.json" without directory prefix) correctly left unchanged.
- **AC7 PASS** — Files at depth-4 new locations (docs/reference/protocol/cn/ and docs/reference/runtime/extensions/): THREAD-EVENT-MODEL.md uses no relative ../ links (confirmed no relative link output). extensions/README.md similarly clean.
- **AC8 PASS (STOP-class)** — `git diff main...HEAD -- .github/workflows/build.yml` shows exactly one line changed: line 223, `docs/alpha/schemas/protocol-contract.json` → `docs/reference/schemas/protocol-contract.json`. No other job, matrix, dependency, or cache change.
- **AC9 PARTIAL** — dune and opam not available in this environment (confirmed independently). OCaml changes verified as doc-comment/string-literal citation only: cn_workflow.ml (lines 4, 20 — `(** ... *)` block and prose comment), cn_packet.ml (line 8 — end-of-comment citation), cn_contract_test.ml:3, cn_traceability_test.ml:5 (doc-comment preambles). No executable logic modified. dune file change: single comment line in src/ocaml/lib/dune. Environment constraint is acceptable.
- **AC10 PASS** — `docs/reference/schemas/` contains exactly: `peers.schema.json`, `protocol-contract.json`, `DESIGN-LLM-SCHEMA.md`, `DESIGN-LLM-SCHEMA-README.md`. No filename collision. beta/schema/README.md correctly renamed to DESIGN-LLM-SCHEMA-README.md at destination.
- **AC11 PASS (STOP-class)** — `git diff --name-only main...HEAD | grep -Ei 'golden|snapshot|fixture'` returns empty. `test/cmd/protocol-contract.json` verified unchanged.
- **AC12 PASS** — `docs/alpha/design/WRITER-PACKAGE.md` is in the do-not-touch zone (§2.3, §4.2 of gamma-scaffold). The change (1 line: `docs/alpha/ctb/LANGUAGE-SPEC.md` → `docs/reference/ctb/LANGUAGE-SPEC.md`) meets all auto-allowed trivial repair criteria: active stale-link repair only; 2 changed lines (1 removed + 1 added); target unambiguous from move map (ctb → reference/ctb/); no workflow semantic change; no source behavior change; no golden; no frozen record; no schema collision. Self-coherence records the repair at AC12. No δ override needed under the auto-allowed policy.
- **AC13 PASS** — See AC5/AC13 above. All remaining stale refs classified.
- **AC14 PASS** — Sampled ORCHESTRATORS.md and PROTOCOL.md diffs: new file additions (git mv'd content), no prose edits visible. All modified-in-place docs show only path string replacements.
- **AC15 PASS** — No hidden/bidi/object-replacement characters detected; all changes are plain ASCII path string replacements.
- **AC16 PARTIAL** — Required checks §1–8 present. Missing as explicit named sections: "do-not-touch proof", "link-check result", "honest β/δ repair record". See C3 finding. All other check categories present with actual output.
- **AC17 PARTIAL** — No new reds in any gate. STOP-class gates (build.yml, golden protection) confirmed clean. OCaml changes comment-only. Go changes comment-only (see C2 finding on self-assessment inaccuracy). Schema content unchanged.

### AC9 note

The OCaml changes are confirmed doc-comment/citation only. cn_workflow.ml: changes are inside the `(** ... *)` doc-comment block (module header) and a prose comment line — both are `(* ... *)` or `(** ... *)` delimiters with no executable content. cn_packet.ml: single comment citation at line 8. test files: doc-comment preambles only. The dune file change is a single `;`-prefixed comment line. No executable logic was modified in any of the four OCaml files or the dune file. The dune-not-available environment constraint is acceptable: the gate cannot be proven green independently, but the constraint is symmetric with α's environment and the changes are structurally incapable of introducing a new build failure.

### AC12 evaluation

`docs/alpha/design/WRITER-PACKAGE.md` is explicitly listed as a do-not-touch file in gamma-scaffold §2.3 ("outside named 10 bundles; not in scope") and §4.2. However, the change (single line, path citation to `docs/alpha/ctb/LANGUAGE-SPEC.md` → `docs/reference/ctb/LANGUAGE-SPEC.md`) satisfies every condition of the auto-allowed trivial repair policy: active stale-link repair only; ≤10 lines changed (2 total); target unambiguous from the move map; no workflow semantic change; no source behavior change; no golden; no frozen record; no schema collision; receipt records the miss and repair (self-coherence AC12). The repair is auto-allowed. No δ override listing required.
