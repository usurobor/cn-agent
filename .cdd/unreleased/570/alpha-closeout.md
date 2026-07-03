# α close-out — cnos#570

**Issue:** #570 — cdd/cds: codify cell kinds and the matter each cell produces
**cell_kind:** `doctrine` (per `.cdd/unreleased/570/gamma-scaffold.md`)
**run_class:** `first_pass` — no prior rejection/repair; this cell has not been through `status:changes` before this cycle.
**Rounds:** R0 only. β returned `verdict: converge` on the first round; no repair round was required.

## Summary

Implemented the cell-kind taxonomy doctrine per AC1–AC9:

- New doctrine file `src/packages/cnos.cdd/skills/cdd/CELL-KINDS.md` (166 lines) naming all 10 cell kinds the issue requires (`issue_authoring`, `implementation`, `repair`, `recovery`, `wave`, `cleanup`, `audit`, `release`, `doctrine`, `experiment`), each with purpose / required input / matter produced / review surface / closeout-projection / FSM implication.
- A `cell kind -> valid matter -> allowed transition request` FSM-awareness table (AC7).
- A chosen `cell_kind` recording point (AC8): the `**cell_kind:** \`<kind>\`` line in `.cdd/unreleased/{N}/gamma-scaffold.md`, demonstrated as machine-readable (not just documented) by a new `parseCellKind` helper in `src/packages/cnos.issues/commands/issues-fsm/fetch.go`'s `assembleLive`, which populates `FactSnapshot.CellKind.Observed` / `.Source = "cdd_artifact"` as an **observation-only** fact — `table.go` (transition logic) was deliberately left untouched, and the pre-existing `TestSeam_CellKindNotEnforced` was re-verified byte-identical against `origin/main` to prove the wiring cannot change any transition decision.
- Cross-reference links from `CDD.md`, `cdd/issue/SKILL.md`, and `docs/reference/governance/GLOSSARY.md` into `CELL-KINDS.md` (AC9), including a new "an issue can itself be a cell's matter" paragraph in `issue/SKILL.md` (AC3).
- Two new Go tests (`TestParseCellKind`, `TestAssembleLive_ObservesCellKindFromGammaScaffold`) covering the new parse path.

Total diff: 8 files (1 new doctrine file, 3 doc files gaining a cross-reference, 1 Go file gaining an observation-only parse, 1 Go test file, plus `gamma-scaffold.md` and `self-coherence.md` under `.cdd/unreleased/570/`). No CUE schema, verifier enforcement, FSM-consumption, dispatch-label, or wake-behavior change was made — matches the issue's own Deferred/Non-goals lists.

## Friction log

- One judgment call was surfaced to β rather than resolved unilaterally: the `cellKindLinePattern` regex used to parse the `cell_kind:` line is intentionally permissive (no strict schema validation) rather than a hardened parse. This matches the issue body's own framing ("a regex or line-scan is fine," no mandatory verifier enforcement yet) and the scaffold's explicit deferral of CUE/verifier work, but was named explicitly in `self-coherence.md §Self-check` for β to confirm rather than assumed acceptable.
- The full named CI gate list (I1, I2, I4, I5, I6, install-wake-golden, dispatch-repair-preflight, dispatch-closeout-integrity, Go, Package, Binary) was not directly observable from the implementation sandbox. Only the local Go toolchain subset (`gofmt`, `go build`, `go vet`, `go test`) could be run and verified clean. This was disclosed as known debt in `self-coherence.md §Debt` rather than silently assumed green.
- No blocking friction otherwise. The γ scaffold's "Expected diff scope" and per-AC oracle table were specific enough that the diff landed with the predicted file set on the first pass — no scope renegotiation was needed mid-cycle.

## Observations

- The generic-kernel + typed-refinement precedent already established for CUE (`schemas/cdd/receipt.cue` → `schemas/cds/receipt.cue` via `#CDSReceipt: cdd.#Receipt & {...}`) transferred cleanly as the organizing shape for markdown doctrine authoring (one generic kernel, ten typed cell-kind refinements) — no new authoring pattern had to be invented for this cell.
- Structuring the "observation, not enforcement" guardrail as a machine-checkable byte-identical-test assertion (rather than a prose promise) made the AC10 no-behavior-change claim independently verifiable by β without re-deriving trust in α's diff-scope narrative.
- This cycle converged at R0 with no repair round — consistent with the `run_class: first_pass` classification recorded at dispatch time.

## Engineering-level reading

The diff footprint matched γ's scaffold-predicted "Expected diff scope" list exactly: one new doctrine file, three doc cross-reference additions, one observation-only Go change, and its companion tests. No AC required renegotiation of scope or an additional design artifact beyond the issue body and scaffold. The clean R0 convergence suggests the scaffold's per-AC oracle table (file existence / grep / byte-diff / named test) was specific enough to remove ambiguity about what "done" meant for a doctrine-first, no-behavior-change cell.
