# α close-out — cycle/392

**Actor:** δ-as-agent (γ+α+β-collapsed-on-δ mode for cycle/392)
**Issue:** cnos#392 (Phase 3 remediation v2; supersedes #391)
**Branch:** `cycle/392`
**Mode:** design-and-build, γ+α+β-collapsed-on-δ

## Summary

V (Contract × Receipt → ValidationVerdict) ported from Python to Go and exposed as the `cn cdd-verify` kernel-tier subcommand. The Go source lives at the cnos.cdd package's command space (`src/packages/cnos.cdd/commands/cdd-verify/`), brought into the `cn` build via a repo-root `go.work` workspace. The Python `cn-cdd-validate-receipt` and bash `cn-cdd-verify` predecessors are removed; the cnos.cdd package manifest no longer declares `cdd-verify` (now a kernel command).

## Cycle execution quality

- **Implementation contract:** all 7 axes honored (language Go ✓; CLI integration target ✓; package scoping ✓; existing-binary disposition ✓; runtime deps ✓; JSON schema unchanged ✓; backward compat ✓).
- **eng/go SKILL:** every section §1.0 – §4.9 walked; no non-conformance.
- **ACs:** 8/8 oracles pass on first attempt (no R2 needed).
- **Tests:** 37/37 + 5/5 + 39/39 + 27 Go unit tests + go vet + go test -race — all green.
- **Diff scope:** matches gamma-scaffold expectations; no creep.

## Findings produced

α emitted 4 findings into cdd-iteration.md (`F1`–`F4`), all operator-specified at dispatch. ε will fold them into a coordinated patch covering α/β/γ skill updates and δ-Phase-4 scope refinement.

## Decisions log

| ID | Decision | Driver |
|---|---|---|
| D1 | Kernel-tier command (not package-vendored ExecCommand) | CLI integration target pin: "NOT a separate binary" — package-vendored entrypoint exec's a binary; kernel-tier compiles INTO cn. |
| D2 | Go workspace (`go.work` + replace directive) | Required to honor package-scoping pin (V source under `src/packages/cnos.cdd/commands/cdd-verify/`) while compiling INTO the cn binary (CLI pin). |
| D3 | YAML parsing via `cue export` subprocess | Required by runtime-deps pin (no new external Go deps). Reuses the Python predecessor's fallback path. |
| D4 | Remove `cdd-verify` entry from cnos.cdd `cn.package.json` | Required: kernel-tier owns the command name; the package's entrypoint script no longer exists. |
| D5 | Bash wrapper fully removed (not thin shim) | Operator preference per dispatch ("prefer full removal"). |
| D6 | Legacy ledger logic ported into Go (not deferred) | Required by backward-compat pin: `--unreleased`/`--all`/`--version`/`--pr` are existing operator surfaces. |

## Self-rating

- **Pattern (α):** A. Implementation contract honored axis-by-axis; eng/go SKILL walked section-by-section.
- **Relation (β):** A. β review APPROVED R1 with no findings. Backward compat preserved.
- **Process (γ):** A. δ-pinned implementation contract gave α an unambiguous frame; all decisions traceable to pins.

C_Σ ≈ A.

## Cycle artifacts (.cdd/unreleased/392/)

- gamma-scaffold.md
- design-notes.md
- self-coherence.md
- beta-review.md
- alpha-closeout.md (this file)
- beta-closeout.md
- gamma-closeout.md
- cdd-iteration.md
