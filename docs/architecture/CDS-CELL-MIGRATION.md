# CDS Cell Migration — from prose wake to a runnable data cell

Status: plan · Scope: `cnos.cdd` (substrate) + `cnos.cds` (program) · Focus:
**run a CDS cell from the CLI.** GitHub-issue mechanics (selector, claim, FSM
label transitions) are explicitly out of scope for this pass — they stay in the
existing wake and are wired back later.

## Goal

Today `cnos.cds` is expressed as a prose **wake** (`orchestrators/cds-dispatch/
SKILL.md`): a selector claims an issue and hands it to δ, which routes γ/α/β by
reading skill markdown. The cell only exists as text.

Migrate CDS so the cell exists as **data the kernel runs**:

```
main.cell  --compile-->  spec.json  --cue vet #CellSpec-->  cellkernel.Run
 (surface)   (cnos.cdd)    (IR)          (oracle)              (Go kernel)
```

This is the exact mirror of the pipeline we already ship on the *output* side
(γ emits a receipt → `cue vet #CDSReceipt` → V). We are adding the *input* side.

## Why — the reasoning behind every choice

This plan is the settled end of a long design dialogue. The reasoning matters
as much as the steps, because each decision closed a specific fork.

1. **The kernel is pure mechanics; it cannot hold a resolver.** `cellkernel.Run`
   executes the five-step closure at one scope and knows *nothing* about the
   work being done (which language, which skills, which repo). Therefore the
   cell spec must be **constructed before/around the run** — never resolved
   inside the kernel. This is the load-bearing separation the whole design
   protects.

2. **The spec is static data "all the way down" — but only one level at a
   time.** For a mono-lingual repo where the language falls out mechanically,
   the whole tree can be compiled up front (Case A). When a decision needs
   *judgment* (a multilingual repo where the issue demands Rust vs Go), the
   child spec does not exist until the parent's α runs — the tree unfolds one
   static level at a time (Case B). Each cell is static *when it runs*; deeper
   specs are produced at runtime by α cognition emitting data.

3. **Composition is one reentrant mechanism, not a new construct.** A subcell
   becomes an α via an adapter `α = extract ∘ Run ∘ resolve` (the `.NET Rx`
   `SelectMany` / F# `let!`). The parent's α calls the kernel on the child; the
   kernel is reentrant and stays oblivious. This is the settled industry answer
   — Railway-Oriented Programming (F# `Result` composition), the GoF Composite +
   Decorator patterns (uniform interface so a subcell stands in for α), the
   interpreter-over-data / Free-monad pattern (recursion lives in the *data + a
   recursive walker*, not a grammar), and Unix pipes (compose only through the
   typed payload — for us, the receipt). Because composition is closed under the
   α interface, nothing in the kernel is ever re-implemented to recurse.

4. **The static-vs-dynamic "spec expression" dilemma was never a language
   problem.** We weighed compiling everything to a static resolved spec against
   a dynamic expression language carrying `resolve` at runtime. TSC's
   `cm-language` (the F#-shaped `coh` surface) settled it by example: it is a
   *compiler to a CUE-validated data IR* with **no runtime** — it names
   providers but executes nothing. The dynamism we feared (`α = resolve |>
   coding_cell`) is not a language operator; it is α cognition *emitting data*
   that the kernel then loops over. So: **spec is data, `resolve` is a library
   function, recursion is the kernel's loop.** No runtime expression language.

5. **We adopt `cm-language`'s *shape and toolchain architecture*, not its
   *semantics*.** cm is measure-only by construction (`forbid compile, admit,
   authorize, repair, self_authorize`) — the categorical opposite of a cell
   whose α *produces* and whose δ *repair-dispatches*. So we do not adopt the
   language. We adopt its form: the compact F#-computation-expression surface
   (`NAME (params) -> ReturnType { … }`, `let!`/`and!` binders) compiled to a
   CUE-validated data IR with **CUE as the independent oracle**. It is the most
   compact legible way we have seen to express a bounded, typed, composable
   workflow; the noise of raw YAML (`mode: rented` everywhere) is exactly what
   it removes.

6. **No OCaml.** cm chose OCaml because TSC *is* an OCaml toolchain. Our kernel
   is Go; introducing OCaml means a second toolchain and a subprocess/FFI seam
   between compiler and kernel — pure friction. A ~150-line Go `participle`
   grammar (or Starlark to bootstrap) targets `cellkernel.Spec` natively.

7. **Compile to JSON + `cue vet` — because it mirrors what we already run.** The
   receipt path is already `γ emits → cue vet #CDSReceipt → V`. Adding the
   input-side mirror (`compile → cue vet #CellSpec → Run`) means both ends share
   one compile-vet-execute shape. That symmetry is why the remaining work is
   three small artifacts, not a new architecture.

8. **Parameters are Unix typed holes — because that is the seam that lets us
   start simple and grow.** Filling `language` from the CLI now, from a parent
   cell (mechanically or cognitively) later, changes only *who fills the hole* —
   never the cell. Unix already solved value→implementation resolution (`$PATH`)
   and required/optional inputs (positional vs flag); we borrow it wholesale.

9. **The four protocols already differ in only two things** (verified against
   `schemas/cds/receipt.cue` and `schemas/cdr/receipt.cue`: both are
   `cdd.#Receipt & { protocol_id, evidence_refs }`). A protocol *is*
   `(protocol_id, evidence-key set, skills)`. So one generic cell + a thin
   overlay expresses cdd/cds/cdr/cdw — the surface differs only in the header
   protocol and the two skill lines. This is what makes "CDS is one cell the
   runner runs" concrete rather than aspirational.

## Two settled decisions

1. **Entrypoint file:** `main.cell` (per package). Mirrors `main.go`; the file
   holds cell definitions, the CLI runs one by name (`cn cell run cds`).
2. **Parameters are Unix-shaped typed holes.** See "Parameters → skills" below.

## The surface (F#-inspired, same shape as TSC `coh`/cm — our own vocabulary)

The CDS cell, concretely, grounded in `schemas/cds/receipt.cue`:

```
cell cds (issue, repo, language: skill, style?: skill = functional)
    -> cnos.cdd.cds.receipt.v1 {
  produce with eng, $language, $style     # α: implement
  review  with $language, cds-review       # β: review
}                                          # γ/V/δ: mechanical kernel defaults
```

- The header's `-> cnos.cdd.cds.receipt.v1` selects the receipt schema; V
  dispatches on it and cue-vets against `#CDSReceipt`.
- The five per-cycle artifacts (`self_coherence, alpha_closeout, beta_review,
  beta_closeout, gamma_closeout` + `diff`) are **not authored** — they are the
  protocol's evidence contract; **γ binds them mechanically**.
- cdr / cdw are the *same body shape*, differing only in header protocol and
  the two skill lines (cdr overlay exists; cdw schema not yet shipped).

A cell whose α is a subcell (the parent/dispatch case, later) is one `let!`:

```
cell cds-dispatch (issue, repo, language: skill) -> cnos.cdd.cds.receipt.v1 {
  let! impl = cds (issue, repo, $language)   # α = SelectMany over the child cell
  review with dispatch-review                 # β reviews the choice, not the code
}
```

## Parameters → skills (the resolution model)

A parameter is a **typed hole**, resolved exactly like a Unix command resolves
argv against `$PATH`:

| Concern | Unix analogue | Here |
|---|---|---|
| required vs optional | positional arg vs flag-with-default | `language: skill` (required) vs `style?: skill = functional` |
| value → implementation | `ls` → `/bin/ls` via `$PATH` | `"ocaml"` → the `ocaml` skill via the **skill path** |
| domain / typo check | — | CUE constrains `language: "go"|"ocaml"|"rust"`; bad value fails `vet` |
| splice into the seat | — | `$language` in `produce with` |
| pass down to a child | wrapper `exec cmd "$@"` | `let! impl = cds(issue, repo, $language)` |

Four layers, matching the rest of the architecture:

1. **Surface** declares holes (`language: skill`, `style?: skill = functional`)
   and splices them (`$language`).
2. **CUE `#CellSpec`** carries the holes and (via the cds overlay) constrains
   their domains — a missing *required* hole or an out-of-domain value fails
   vet. Mirror of how `#CDSReceipt` overlays `#Receipt` with `evidence_refs`.
3. **Runner** fills holes from the CLI (`--language go`), resolves each value to
   a concrete skill ref against the skill path, loads it into the seat.
4. **Kernel / seats** receive *loaded skills*. α/β never see the string
   `"ocaml"` — they get the resolved skill. (Seats consume skills, never
   resolve them — same top-down rule as everywhere.)

**Staging of who fills the hole (the cell never changes):**
- *Now (CLI bootstrap):* the runner fills every hole from CLI flags.
- *Later:* a parent cell fills `language` mechanically (`language_of(repo)`) or
  cognitively (`triage(issue)`). Only the *filler* changes; resolution
  (value→skill) and the `cds` cell body are untouched.

## Piece inventory

| Piece | Owner | Status |
|---|---|---|
| Receipt schema (output contract) | cnos.cdd + cds | ✅ shipped |
| `cellkernel.Run` (five-step closure) | src/go `internal/cellkernel` | ✅ built (empty cell green) |
| Provider seam for rented cognition | src/go `internal/dispatch.Backend` | ✅ exists |
| **`#CellSpec` CUE schema (input contract)** | cnos.cdd | ❌ **new** |
| **cds params-domain overlay (`#CDSCellSpec`)** | cnos.cds | ❌ new |
| **JSON `spec.json` → `cellkernel.Spec` loader** | src/go `internal/cellspec` | ❌ new |
| **`cn cell run` (fill holes, resolve skills, Run)** | src/go `internal/cli` | ❌ new |
| **skill-path resolver (value → skill ref)** | cnos.cdd | ❌ new |
| **`main.cell` compiler (surface → spec.json)** | cnos.cdd | ❌ new (Go `participle`, or fork TSC `cm_surface.ml` shape) |
| **`cnos.cds/main.cell`** | cnos.cds | ❌ new |

## Phases (smallest-first; each independently demonstrable)

**Phase 0 — the input contract.** Write `schemas/cdd/spec.cue` (`#CellSpec`:
`{contract, protocol_id, params, alpha:{skills}, beta:{skills}, budget?}`) and
the cds overlay (`#CDSCellSpec` pinning `protocol_id` + param domains).
Hand-write a `cds.spec.json` for a real issue; `cue vet` it. *Proves the data
shape before any Go or compiler.*

**Phase 1 — run it from the CLI (stub cognition).** `internal/cellspec` loader
(spec.json → `cellkernel.Spec`), and `cn cell run <spec.json> --language go`
that fills holes and calls `cellkernel.Run` with **stub α/β** (echo producer /
accept reviewer). Prints the CCNF trace + receipt. *Proves the execution path
end-to-end with no compiler and no cognition.*

**Phase 2 — skill resolution.** The `$PATH`-like resolver (`value → skill ref`)
+ the required/optional/default hole logic; `cn cell run` errors on unfilled
required holes with a Unix-style usage line. Seats load the resolved skills
(cognition still stubbed). *Proves parameters map to real skills.*

**Phase 3 — the surface + compiler.** `cn cell compile main.cell → spec.json`
(Go `participle` grammar, ~150 lines; CUE stays the independent oracle). Author
`cnos.cds/main.cell`. Now `cn cell run cds --language go` runs the surface
directly. *The compiler is pure sugar over Phase-1 JSON — optional until it
earns its keep.*

**Horizon (out of this pass):** real α/β via `internal/dispatch.Backend`
(rented cognition); the `cds-dispatch` parent cell (`let!` composition); wiring
the existing GitHub wake to fill holes + land artifacts instead of routing δ by
prose.

## Ownership split

- **cnos.cdd** owns the *substrate*: the `#CellSpec` schema, the compiler, the
  loader, the runner (`cn cell compile` / `cn cell run`), and the skill-path
  resolver. It is the language + kernel.
- **cnos.cds** owns a *program written in it*: `main.cell`, the params-domain
  overlay, and its α/β skills. cds is one profile; cdw/cdr are siblings of the
  same shape.
