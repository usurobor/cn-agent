# α close-out — cnos#614

**Issue:** #614 — cdd: amend CELL-KINDS (#570) — add `planning` cell kind + mandatory terminal learning / ε-observations section.
**cell_kind:** `doctrine` (per `.cdd/unreleased/614/gamma-scaffold.md`)
**run_class:** `first_pass` — see the explicit note below on why this is not `repair_pass` despite pre-existing branch/artifact matter.
**Rounds:** R0 only. β returned `verdict: converge` on the first round; no repair round was required.

## run_class disclosure (resumed-cell transparency)

`.cdd/unreleased/614/CLAIM-REQUEST.yml` and `gamma-scaffold.md` pre-existed this firing (written 2026-07-07T05:03 by a prior firing that then hit a write-fence false-positive — root cause: a concurrent `agent-admin` heartbeat commit misattributed to the worker, fixed in #625). Mechanically, this pattern-matches the cnos#516 repair-re-entry detection heuristic ("cycle/{N} branch already exists with commits beyond base" + "pre-existing `.cdd/unreleased/{N}/` artifacts"). It is classified `first_pass`, not `repair_pass`, because: (1) there is no prior `status:changes` history on #614 — the issue never reached β/δ content review before this firing; (2) there are no prior β findings, δ findings, or bounce comments to repair against — `gamma-scaffold.md` is unreviewed dispatch matter, not rejected work; (3) the operator's 2026-07-08 "Dispatch authorized" comment explicitly frames this as "a clean resume," not a repair, and directs continuing from the scaffold as valid prior matter. A `repair_evidence` block would have nothing genuine to point at (no rejected findings exist); asserting one would be a counterfeit closeout, not a transparent one. This paragraph is the disclosure the cnos#516 preflight exists to force.

## Summary

Implemented Amendments 1–3 per the issue body and γ scaffold's α prompt, resuming from the existing `gamma-scaffold.md`:

- **Amendment 1:** new `### 11. \`planning\`` cell kind in `CELL-KINDS.md` (Purpose / Required input / Matter produced / Lifecycle α-β-γ / Role-vs-model note / Relationship to neighbors / Worked example / FSM implication), mirroring the shape of the existing 10 entries. Heading bumped "10 cell kinds" → "11 cell kinds". New `planning` row added to the FSM-awareness table.
- **Amendment 2:** new `## Mandatory terminal learning section` doctrine block in `CELL-KINDS.md`, stating every terminal cell closeout must carry a `learning`/`epsilon_observations` section (field schema: `observations`, `process_deltas`, `reusable_patterns`, `followups`, `operator_burden`), applying to all cell kinds, framed as ε-captured-by-γ, with verifier enforcement explicitly deferred.
- **Amendment 3:** `gamma/SKILL.md` frontmatter description, Core Principle, and §2.7 `gamma-closeout.md` authoring bullet all updated to state explicitly that γ binds learning into the receipt at closeout.
- `CDD.md`'s `CELL-KINDS.md` pointer updated (adds `planning` to the enumeration; names the mandatory learning-section rule).
- `schemas/cdd/receipt.cue` gained a documentation-only comment cross-referencing the mandatory learning section (no field/type change — independently verified via `cue vet` against both the generic and CDS-unified schema/fixture pairs, both exit 0).

Pre-work: rebased `cycle/614` onto current `main` (base was 73 commits stale at time of claim; rebase was clean, zero conflicts) so the PR carries a minimal, current diff rather than a 73-commit-stale one.

Total diff (this cell's own contribution): 4 files (`CELL-KINDS.md`, `gamma/SKILL.md`, `CDD.md`, `schemas/cdd/receipt.cue`), all `.md`/`.cue`. Zero `.go`, `transitions.json`, CI workflow, or label files touched — matches the issue's AC5 and Out-of-scope/STOP conditions exactly.

## Friction log

- The prior firing's infra false-positive (write-fence misattribution) meant this firing had to first establish that the pre-existing branch/scaffold were valid matter to resume from, not stale/conflicting state — resolved by reading the operator's explicit "Dispatch authorized" comment before touching any file, per the resumed-cell disclosure above.
- No blocking friction on the amendment implementation itself. The γ scaffold's source-of-truth table and α prompt were specific enough (four named files, per-file instructions) that the diff landed on the first pass with no scope renegotiation.
- `cue` was not preinstalled in the sandbox; installed via `go install cuelang.org/go/cmd/cue@latest` to independently verify the `receipt.cue` comment-only claim rather than asserting it without running the tool.

## Observations

- The doctrine-first precedent #570 established (one generic kernel, N typed cell-kind refinements, each with the same six/eight-field shape) transferred cleanly to adding an 11th kind — no new authoring pattern needed.
- Placing the "Mandatory terminal learning section" as its own top-level doctrine block (rather than folding it into the `planning` entry) was the right structural call per the issue's own Amendment 2 framing ("A general rule across ALL cell kinds") — nesting it under `planning` would have implied narrower scope than intended.
- This cell is itself an early consumer of the doctrine it lands: `gamma-closeout.md` for this cycle carries a `learning` section as a dogfooding demonstration (not yet mechanically required, but consistent with the change's own spirit — see γ closeout).

## Engineering-level reading

The diff footprint matched the γ scaffold's predicted scope exactly: four doctrine/schema files, no code. The resumed-cell handling (rebase + explicit run_class disclosure) added process overhead beyond the amendment work itself but is the correct, transparent way to handle a cell whose branch/artifacts pre-date this firing — asserting a silent `first_pass` without disclosure, or a fabricated `repair_pass` with an empty `repair_evidence` block, would both have been worse than the explicit paragraph above.
