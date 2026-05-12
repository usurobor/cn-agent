# Self-Coherence Report — Issue #328

## Gap

**Problem statement:** CDD artifacts cannot be missing from active or released cycle directories without CI reporting the gap.

**Current incoherence:** The existing `cdd-verify` command validates individual release versions but not the repository-wide CDD artifact ledger. CI runs build/test checks but lacks systematic validation of CDD lifecycle artifacts across both `.cdd/unreleased/{N}/` and `.cdd/releases/{X.Y.Z}/{N}/` directories.

**Specific gaps:**
1. **Stale validator contract** - Current `cdd-verify` checks legacy paths like `.cdd/releases/{version}/alpha/CLOSE-OUT.md` instead of current issue-scoped paths `.cdd/releases/{version}/{issue}/alpha-closeout.md`
2. **Missing CI integration** - No continuous validation of CDD artifact completeness across the repository
3. **Incomplete coverage** - Pre-tag release gate only checks `.cdd/unreleased/` but not historical `.cdd/releases/` ledger
4. **No small-change detection** - Cannot distinguish between substantial triadic cycles (requiring full artifact set) and small-change cycles (with explicit collapse rules)
5. **Silent historical gaps** - Missing historical artifacts go undetected without systematic checking

**Impact:** CDD can produce "witness theater" where artifacts appear to exist but validation is incomplete, leading to:
- Incomplete cycle directories surviving after release
- Missing close-outs not detected until manual inspection
- RELEASE.md present but potentially incorrect for the version
- CDD lifecycle rules remaining prose without mechanical enforcement

**Expected behavior:** CI fails on unexcepted missing required artifacts, with clear diagnostics identifying the missing path and expected contract.

## Skills

**Mode:** MCA (Minimal Coherent Action)

**Tier 1 skills (always loaded):**
- `CDD.md` — canonical CDD lifecycle and artifact contract
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — α role contract and execution detail
- `src/packages/cnos.cdd/skills/cdd/SKILL.md` — CDD package loader

**Tier 2 skills (general engineering):**  
- `eng/tool` — shell command behavior, exit codes, diagnostics, local invocation
- `eng/test` — positive/negative fixtures and CI proof  
- `eng/ux-cli` — failure messages should tell contributors what artifact is missing and where

**Tier 3 skills (issue-specific):**
- `cnos.core/skills/design` — separate release-time gate from repo-wide artifact ledger validation
- `cdd/review` — verify lifecycle/checker coherence and avoid false closure

**Generation constraints applied:**
- Tool skill: Commands must provide clear exit codes, helpful diagnostics, and support local invocation
- Test skill: Must include both positive fixtures (complete cycles pass) and negative fixtures (incomplete cycles fail with clear messages)  
- UX-CLI skill: Error messages must identify the specific missing artifact path and expected contract
- Design skill: Must distinguish between release-time validation (scripts/validate-release-gate.sh) and repository-wide CI validation (this implementation)
- Review skill: Checker must validate artifact presence and basic structure, not semantic content truth

## ACs

**AC1: CDD package checker validates current artifact layout**
Evidence: Updated `src/packages/cnos.cdd/commands/cdd-verify/cn-cdd-verify` with `--all` and `--unreleased` modes that scan `.cdd/unreleased/{N}/` and `.cdd/releases/{X.Y.Z}/{N}/` directories (commits 4326d259, dabf6c05).

**AC2: Stale verifier contract is replaced or updated**  
Evidence: Checker now validates `.cdd/releases/{version}/{issue}/alpha-closeout.md` paths instead of legacy `.cdd/releases/{version}/alpha/CLOSE-OUT.md`. New `check_triadic_artifacts()` and `check_unreleased_cycle()` functions implement current contract (commit 4326d259).

**AC3: Triadic/substantial cycle requirements are enforced**
Evidence: `classify_cycle_type()` function detects triadic cycles (presence of `beta-review.md`) and enforces all 5 required artifacts. Test fixture `valid-triadic` passes with complete artifact set; `incomplete-triadic` fails with missing close-outs (commits 4326d259, f98da74e).

**AC4: Small-change collapse is explicit**
Evidence: Cycles without `beta-review.md` are classified as small-change. `check_small_change_artifacts()` applies CDD.md §1.2 collapse rules - `self-coherence.md` and `alpha-closeout.md` are optional. Test fixture `valid-small-change` demonstrates correct classification (commits 4326d259, f98da74e).

**AC5: Released and unreleased dirs are both checked**
Evidence: `--all` mode scans both `.cdd/unreleased/` and `.cdd/releases/` directories. `--unreleased` mode checks only active cycles. Repository-wide validation implemented in commits 4326d259, dabf6c05.

**AC6: Legacy exceptions are explicit**
Evidence: Added `--exceptions` flag with YAML exception file support. Missing artifacts that are exception-backed produce warnings instead of failures. `exceptions.example.yml` shows required fields: path, missing_artifacts, reason, repair_possible, follow_up (commit 4db68fc4).

**AC7: Release-note freshness is either checked or explicitly deferred**  
Evidence: **Option B chosen** - checker states in help text and README.md that "Release-note freshness (unless implemented separately)" is not validated. Focus is on artifact presence, not semantic freshness validation (commit dabf6c05).

**AC8: CI job exists**
Evidence: Added `cdd-artifact-check` job (I6) to `.github/workflows/build.yml` that runs `cn-cdd-verify --all` on push and pull request (commit 58bb3f15).

**AC9: Notify aggregation includes the new CI job**
Evidence: Updated notify job `needs:` list to include `cdd-artifact-check` and result aggregation loop to include `"I6:${{ needs.cdd-artifact-check.result }}"` (commit 58bb3f15).

**AC10: Local run documentation exists**
Evidence: Updated help text with repository-wide mode descriptions and examples. Added `README.md` with quick start guide, examples, and cycle classification explanation (commit dabf6c05).

**AC11: Tests cover positive and negative fixtures**
Evidence: Created test fixtures and `test-fixtures.sh` script covering:
- Complete triadic cycle (`valid-triadic`) → passes
- Missing artifacts (`incomplete-triadic`) → fails  
- Small-change cycle (`valid-small-change`) → correct classification
- Exception handling with `exceptions.example.yml` → warnings instead of failures
All tests pass per commit f98da74e.

**AC12: release.sh and CDD checker boundaries are clear**
Evidence: Help text and README.md document: "`scripts/validate-release-gate.sh` = pre-tag release gate; `cdd-verify` = repository ledger validation." Different purposes: release gate checks current release, CDD checker validates entire repository history (commit dabf6c05).

## Self-check

Did α's work push ambiguity onto β? **No** - all claims are backed by concrete evidence:
- Every AC maps to specific commits and file paths
- Test fixtures demonstrate actual checker behavior, not just prose claims
- CI integration is verifiable (job I6 exists in build.yml)
- Documentation includes concrete examples users can run

Is every claim backed by evidence in the diff? **Yes** - implementation is incremental and verifiable:
- CDD checker updates: commit 4326d259 + 4db68fc4 (core functionality)
- Test fixtures: commit f98da74e (positive/negative validation)  
- CI integration: commit 58bb3f15 (job + aggregation)
- Documentation: commit dabf6c05 (usage examples + README)

## Debt

**Known debt:**
- Exception file parsing uses simple regex instead of full YAML parser - works for documented format but may break on complex YAML edge cases
- Section validation in artifacts is basic pattern matching, not semantic validation of content quality
- No validation of cross-artifact consistency (e.g., issue ACs match self-coherence ACs)

**No additional debt** - all AC requirements implemented per specification.

## CDD Trace

| Step | Artifact | Skills loaded | Decision |
|------|----------|---------------|----------|
| 0 Observe | — | — | Observation inputs read from issue #328; CDD artifact validation gap selected |
| 1 Select | — | — | Selected gap: CI lacks systematic CDD artifact ledger validation |
| 2 Branch | cycle/328 | cdd | Branch existed; checked out from dispatch |
| 3 Bootstrap | .cdd/unreleased/328/ | cdd | Cycle directory created for self-coherence artifact |
| 4 Gap | self-coherence.md §Gap | — | Named incoherence: missing CI validation of CDD artifacts |
| 5 Mode | self-coherence.md §Skills | cdd + eng/tool + eng/test + eng/ux-cli + cnos.core/skills/design + cdd/review | MCA mode with T2/T3 generation constraints applied |
| 6 Artifacts | Updated cdd-verify command, test fixtures, CI job, README | eng/tool + eng/test + eng/ux-cli | Implementation complete: 4 commits covering functionality, tests, CI, docs |
| 7 Self-coherence | self-coherence.md complete | cdd | AC-by-AC self-check completed with evidence mapping |