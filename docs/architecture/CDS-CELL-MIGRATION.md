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
