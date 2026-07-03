# γ scaffold — cnos#570

**Issue:** #570 — cdd/cds: codify cell kinds and the matter each cell produces
**Mode:** design-and-build (doctrine-first; deferred CUE/verifier/FSM-enforcement per issue's own Deferred list)
**Base SHA:** `3a511121c4cf50fa4024494ae00f7bca6fae47dd` (origin/main at scaffold time)
**Branch:** `cycle/570`
**cell_kind:** `doctrine` (this cell produces CDD doctrine + an observation-wiring change; no dispatch/renderer/FSM-enforcement behavior change)

## Surfaces α is expected to touch

1. **New doctrine file** — `src/packages/cnos.cdd/skills/cdd/CELL-KINDS.md`. Defines the 10 cell kinds named in the issue body (issue_authoring, implementation, repair, recovery, wave, cleanup, audit, release, doctrine, experiment), each with purpose / input / matter / review surface / closeout-projection / FSM implication, per AC1–AC7.
2. **Recording point (AC8)** — `cell_kind: <kind>` recorded as a line in `.cdd/unreleased/{N}/gamma-scaffold.md` (this file demonstrates the convention: see the `**cell_kind:**` line above). Doctrine names this as the first canonical recording point per the issue's own recommendation ("record it in the CDD receipt / self-coherence surface... first pass, no mandatory verifier enforcement").
3. **Observation wiring (per the operator's wave-note comment on #570, dated 2026-07-03T22:01:22Z)** — `src/packages/cnos.issues/commands/issues-fsm/fetch.go`'s `assembleLive`: parse a `**cell_kind:** \`<kind>\`` (or `cell_kind: <kind>`) line out of `.cdd/unreleased/{N}/gamma-scaffold.md` when present, and populate `snap.CellKind.Observed` + `snap.CellKind.Source = "cdd_artifact"`. This is **observation only** — `snap.normalizeCellKind()` must still run unchanged, and `table.go` (transition rules) MUST NOT be touched. `TestSeam_CellKindNotEnforced` (issuesfsm_test.go:417) locks this — it must continue to pass unmodified, proving the wiring cannot change any transition decision.
4. **Doctrine links (AC9)** — add a pointer to `CELL-KINDS.md` from:
   - `src/packages/cnos.cdd/skills/cdd/CDD.md`
   - `src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` (state that an issue can itself be a cell's matter — issue_authoring cell kind)
   - `docs/reference/governance/GLOSSARY.md`

## AC oracle approach

| AC | Oracle | How α proves it |
|---|---|---|
| AC1 | `CELL-KINDS.md` exists, defines ≥10 named kinds | file exists; grep each kind name |
| AC2 | each kind entry has purpose/input/matter/review-surface/closeout/FSM-implication | table or per-kind subsection with all fields present |
| AC3 | issue-authoring cell defined, issue is matter | explicit statement + α/β/γ lifecycle note |
| AC4 | implementation cell distinguished, maps to CDS dispatch | explicit statement + cross-ref to `cnos.cds` dispatch doctrine |
| AC5 | wave cell defined as parent cell | explicit "projection, not role-renaming" statement |
| AC6 | cleanup cell defined | explicit TSC/cleanup-issue/cleanup-PR matter statement |
| AC7 | FSM-awareness section: `cell kind -> valid matter -> allowed transition request` table | table present in `CELL-KINDS.md` |
| AC8 | recording point chosen | doctrine names `.cdd/unreleased/{N}/gamma-scaffold.md`'s `cell_kind:` line; fetch.go wiring demonstrates it is machine-readable |
| AC9 | links from CDD.md, issue/SKILL.md, GLOSSARY.md | grep each file for a `CELL-KINDS.md` reference |
| AC10 | no behavior change to dispatch/renderer/FSM-enforcement | diff review: only `fetch.go`'s observation path touched, `table.go` untouched, `TestSeam_CellKindNotEnforced` passes unmodified |
| AC11 | CI gates green | CI run on PR |

## Expected diff scope

- 1 new doctrine file (`CELL-KINDS.md`)
- 3 doc files gaining a cross-reference link (`CDD.md`, `cdd/issue/SKILL.md`, `GLOSSARY.md`)
- 1 Go file gaining an observation-only parse (`fetch.go`)
- Optionally: a small Go test added to `issuesfsm_test.go` asserting the new parse path populates `Observed`/`Source` without affecting `Evaluate()` output (companion to the existing `TestSeam_CellKindNotEnforced`)
- This scaffold + `self-coherence.md` + `beta-review.md` + closeouts under `.cdd/unreleased/570/`

## Empirical anchor

Direct precedent for the generic-kernel + typed-refinement shape: `schemas/cdd/receipt.cue` (generic `#Receipt`) refined per-domain via CUE unification in `schemas/cds/receipt.cue` (`#CDSReceipt: cdd.#Receipt & {...}`). The operator's design-note comment on #570 (2026-07-03T21:28:28Z) confirms this is the intended Phase-B substrate for a future typed `cell_kind` field — out of scope for this doctrine-first cell.

## Scope guardrails

- Do NOT implement `table.go` transition consumption of `cell_kind`.
- Do NOT touch dispatch labels, wake behavior, or renderer YAML.
- Do NOT require existing cells to backfill `cell_kind`.
- Do NOT start Demo 0, CUE schema, or verifier enforcement.

## Friction notes

None yet — first pass.
