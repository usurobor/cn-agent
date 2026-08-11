# cnos cell-runner — walking-skeleton spike (proof-of-life)

**Status: EXPLORATORY SPIKE, not canonical runtime.** This is a proof-of-life for the
three-cell agentic loop (CC → PC → WC). It is built *against* the R12 cell-runtime doctrine
design (PR #672) to learn from running, not to pre-empt the authored doctrine (#672 WCs) or the
implementation wave (#627 S2–S8). When the doctrine is authorized, this spike is reconciled to it
(or discarded). It invents no policy; where it must choose, it picks the minimal mechanical option
and says so here.

## What it proves

One turn of the loop, mechanically, git-native, no human in the machinery:

```
CC.measure(target)  → CM object + judgment{request_planning}   (incoherence detected)
  → PC.plan(judgment,CM) → wave{ one WC contract }             (relation graph)
    → WC.execute(contract) → α produce → β review → γ close → V validate → δ decide → receipt
      → CC.measure(target) → judgment{coherent}                (incoherence cleared)
  → loop terminates: incoherence 1 → 0, receipt chain closed.
```

The three cells are distinguished ONLY by output telos (the R12 doctrine's classification rule):
**WC → artifact, PC → relation_graph, CC → judgment.** They run the *same* kernel.

## Non-negotiables (operator directives)

- **Go for code, CUE for schemas. NO Python.** `go.work` module + `cue.mod`.
- Reuse cnos idioms: a **declarative FSM transition table + a package-agnostic Go evaluator**
  (mirroring `src/packages/cnos.issues/commands/issues-fsm/`), and receipts shaped after
  `schemas/cdd/receipt.cue`. Do not hardcode a state name in a switch/if-chain.
- Fail closed. `cue vet` every emitted object. `gofmt`/`go vet` clean.

## Architecture

### Kernel FSM (declarative)
`kernel/transitions.json` — the CCNF cell kernel as a transition table over states
`{init, produced, reviewed, closed, validated, decided, held}` with guards evaluated by a generic
evaluator (`kernel/engine.go`). Transitions:

| from | guard | to | step |
|---|---|---|---|
| init | contract_present | produced | α produce |
| produced | artifact_present | reviewed | β review |
| reviewed | review_approved | closed | γ close (compute receipt) |
| reviewed | review_rejected | held | δ hold (stop) |
| closed | receipt_bound | validated | V validate |
| validated | v_pass | decided | δ accept |
| validated | v_fail | held | δ hold (stop) |

The evaluator is generic: it reads the table + a `FactSnapshot` (booleans produced by the actors
and the workspace) and picks the first matching transition, exactly like the CDS issues-fsm engine.

### Actor seam (the "assisted" boundary)
```go
type Actor interface {
    // Alpha produces the artifact; Beta reviews; Gamma closes (computes receipt fields);
    // V validates mechanically; Delta decides; Measure is the CC's CM+judgment producer.
}
```
v0 ships **deterministic stub actors** (`actor/stub/`) so the loop runs end-to-end with zero
external calls — proving the *machinery*. Agent-backed actors (real α/β/CC judgment) implement the
same interface next (v1); nothing in the kernel/loop changes.

### Cells = kernel + telos
A "cell" is the kernel run under a telos:
- **WC**: `requested_output.kind = artifact`; V checks the artifact against acceptance predicates.
- **PC**: `requested_output.kind = relation_graph`; α emits a wave (one WC contract) from the CM.
- **CC**: `requested_output.kind = judgment`; the Measure actor produces a CM object, δ emits a
  judgment (`coherent | request_planning | request_working`).

### CM object (the shared cnos↔tsc seam)
`schema/cm.cue` — `#CM{ target, alpha_score, defects[], provenance, measured_at_ref }`. This is the
**runner/provider-produced measurement object** (NOT the CC's judgment; the CC *consumes* it). It is
the exact seam a tsc-provided TSC measurement plugs into (WC-2 of #672). v0's CM provider is a
mechanical checker over the toy target; a tsc/agent provider drops in behind the same schema.

### Receipts
`schema/receipt.cue` — a minimal runner receipt (id, cell, telos, transitions[], evidence_refs,
verdict, boundary_decision, produced_at_ref) shaped after `schemas/cdd/receipt.cue`. γ computes it
from the transition evidence; it is not free-text authored. Every cell writes one to
`<workspace>/receipts/<cell>.receipt.yaml`; the loop chains them (ε = the receipt stream).

## Toy task (real, mechanical, decreasing-incoherence)
`testdata/toy/` — a target markdown file missing a required `## Coherence` section.
- CC.measure: section absent → defect `missing-required-section` → judgment `request_planning`.
- PC.plan: emits a WC contract "add the `## Coherence` section to the target".
- WC.execute: α appends the section; β approves (section text matches the required shape);
  γ closes; V mechanically asserts the section is present and well-formed; δ accept.
- CC.re-measure: defect cleared → `coherent` → loop ends.

## v0 acceptance (the spike is "alive")
`go run ./cmd/cell-runner --workspace <tmp> --task testdata/toy` exits 0 and:
1. prints a loop transcript showing incoherence `1 → 0` across the CC→PC→WC→CC turn;
2. writes CUE-valid CM, wave, and per-cell receipts under `<workspace>/`;
3. `cue vet` of every emitted object against `schema/*.cue` passes;
4. a `--fail-injection` flag (e.g. β rejects, or V fails) drives the kernel to `held` with a stop
   receipt and non-zero exit — proving the FSM branches and fails closed;
5. `gofmt -l` empty, `go vet ./...` clean, `go test ./...` green.

## Explicitly deferred (NOT in v0)
Agent-backed actors; multi-WC waves / parallel branches; the wave FSM (WC-3b) beyond one node;
real CM scoring beyond defect-count; operator-authorization gate wiring; self-hosting on a real
cnos change; tsc adoption. Each is a named next increment, not a silent gap.

## Build notes

Points where the spec was underspecified and the build took the minimal mechanical option:

- **Guard/step timing in the kernel.** The transition table lists a `step` per edge. The
  implemented semantics: the guard of the edge *out of* a state is satisfied by the step run *into*
  it (the prior step), and taking an edge runs its `step`, which produces the fact the next edge's
  guard reads. `contract_present` is the α-input fact, set true before the run (the contract is
  present at handover). This makes the evaluator a pure first-match projection with no lookahead,
  matching the issues-fsm engine.
- **`step` vs guard registries.** Two fixed generic vocabularies, mirroring issues-fsm's single
  `guardFuncs` map: `guardFuncs` (fact predicates) in `kernel/table.go` and `stepFuncs` (actor
  dispatch) in `kernel/actor.go`. Neither switches on a state name; `Drive` dispatches steps through
  the registry. Adding a state/edge is a `transitions.json` edit.
- **Terminal states are declared** in the table (`"terminal": ["decided","held"]`) rather than
  inferred as "states with no outgoing edge", so `Drive`'s stop condition stays table-driven and a
  genuinely stuck non-terminal state is still an explicit fail-closed error.
- **"Measure" actor method.** DESIGN lists Measure alongside α/β/γ/V/δ. To keep the kernel uniform
  (all cells run the identical six steps), Measure is realized as the **CC actor's α step** (it
  produces the CM); the CC's **δ** emits the judgment. No extra interface method — the `kernel.Actor`
  seam is exactly α/β/γ/V/δ (with δ split into `DeltaAccept`/`DeltaHold` for the two boundary
  actions), so PC/WC don't carry a no-op Measure.
- **Emitted objects are JSON, not YAML.** DESIGN sketches `*.receipt.yaml`. To honor the stdlib-only
  constraint (no third-party YAML dependency) objects are written as `*.json` (JSON is a strict YAML
  subset and `cue vet` reads it directly). Files: `cm/<turn>.cm.json`, `wave/pc.wave.json`,
  `receipts/<turn>.receipt.json`, `judgment/<turn>.judgment.json`.
- **CC runs twice**, so receipts are disambiguated by turn (`cc-1`, `pc`, `wc`, `cc-2`) rather than
  the bare `<cell>` name DESIGN sketches — otherwise the re-measure receipt would clobber the first.
- **`alpha_score`** is scored mechanically as `1.0 - min(1, defect_count)` (so 0.0 with the one
  defect, 1.0 when clear); incoherence is the defect count. A richer provider refines the score
  behind the same `#CM` schema.
- **Fault injection targets the WC turn** (the natural α/β/V gate for an artifact). `beta-reject`
  drives `reviewed → held`; `v-fail` drives `validated → held`. Both write a `not_transmissible`
  stop receipt and exit non-zero.
- **Receipt schema is a minimal spike shape**, not `cnos.cdd.receipt.v1`: it reuses that schema's
  structural verdict × action → transmissibility derivation (the fail-closed gate) but types only the
  fields the skeleton needs (`id, cell, telos, transitions, evidence_refs, validation,
  boundary_decision, transmissibility, verdict, produced_at_ref`).
