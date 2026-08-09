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

3. **Composition is one reentrant mechanism, not a new construct — but α
   *proposes*, the runtime *executes*.** A parent decomposes by having its α
   emit *child contracts* (data); the **runtime/Drive** runs each child episode;
   the accepted child receipts become the parent's matter; parent-β reviews the
   composition against the parent contract. This is the settled industry answer
   — Railway-Oriented Programming (F# `Result` composition), the GoF Composite +
   Decorator patterns (uniform interface so a child stands in for the parent's
   input), the interpreter-over-data / Free-monad pattern (recursion lives in
   the *data + a recursive walker*, not a grammar), and Unix pipes (compose only
   through the typed payload — for us, the receipt).

   **Correction from Pi β (`msg-cn-pi-cnos-cell-runner-cases-review-31`, C1):**
   an earlier framing made parent-α *call the kernel on the child* directly
   (`α = extract ∘ Run ∘ resolve`, the `.NET Rx` `SelectMany`). That fuses
   *decomposition judgment* with *runtime execution* inside a cognition seat. A
   composite-α adapter may be a later convenience, but it must not be the
   normative boundary: **the runtime executes; α only proposes.** The surface
   `let!` (below) is sugar that *lowers to* "α emits child contract → runtime
   runs it → receipt → matter" — not to α invoking `Run` itself. This keeps
   child execution mechanical and auditable (route receipts, custody) rather
   than buried in α.

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
  let! impl = cds (issue, repo, $language)   # α PROPOSES this child contract;
                                             # runtime runs it; receipt → matter
  review with dispatch-review                 # β reviews the composition/choice
}
```

`let!` is surface sugar: it lowers to "parent-α emits the `cds` child contract →
the runtime executes that child episode → the accepted child receipt becomes
the parent's matter." It does **not** compile to α calling the kernel itself
(Pi β C1). Composition and single-episode work therefore share one runtime, not
two cell algorithms.

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
| `cellkernel` five-seat order | src/go `internal/cellkernel` | ⚠ **Case-0 sketch only** — needs Pi D1–D4 (see below) |
| Provider seam for rented cognition | src/go `internal/dispatch.Backend` | ✅ exists |
| **`#CellSpec` CUE schema (input contract)** | cnos.cdd | ❌ **new** |
| **cds params-domain overlay (`#CDSCellSpec`)** | cnos.cds | ❌ new |
| **JSON `spec.json` → `cellkernel.Spec` loader** | src/go `internal/cellspec` | ❌ new |
| **`cn cell run` (fill holes, resolve skills, Run)** | src/go `internal/cli` | ❌ new |
| **skill-path resolver (value → skill ref)** | cnos.cdd | ❌ new |
| **`main.cell` compiler (surface → spec.json)** | cnos.cdd | ❌ new (Go `participle`, or fork TSC `cm_surface.ml` shape) |
| **`cnos.cds/main.cell`** | cnos.cds | ❌ new |

## Kernel corrections required first (Pi β #31, D1–D4)

The current `internal/cellkernel` proves only Case-0 smoke; it is not yet a
reference. Four bounded corrections gate everything below:

- **D1 — honest closure.** Split `RunEpisode → EpisodeResult` (`terminal
  {accepted|degraded|rejected}` | `needs_repair` {typed request, parent stays
  open} | `malfunction` {error}) from a later `Drive` (bounded attempt loop).
  `invalid` (PASS+override, FAIL+accept) is a typed kernel error, never a
  returned closed cell.
- **D2 — no self-certification.** `Spec` carries **Contract + α + β only**; the
  kernel *owns* mechanical γ/V/δ (no injectable seat interfaces). V verifies
  *bindings* — contract/matter/review identity+hash, required evidence refs,
  runtime-produced route evidence is bound (not γ-authored), verdict/decision
  schema-valid — it does not merely mirror β.
- **D3 — evidence seam.** `AlphaResult{Matter; EvidenceRefs}`,
  `BetaResult{Review; EvidenceRefs}` (or a runtime recorder) so γ binds real
  refs. v0 refs are `{id, kind, ref, sha256, producer_execution_id}` (Pi Q4).
- **D4 — fail closed.** A nil/missing α or β returns a wrapped error before any
  seat runs; no panics.

## Phases (Pi's KISS ladder, `msg-…-review-31` D5; smallest-first)

**Phase K — correct the kernel.** Apply D1–D4 above; keep Case-0 smoke green and
add the negative tests Pi specified (self-certification blocked, nil-seat error,
invalid pair → error). *Prerequisite for a truthful CLI.*

**Phase 0 — the input contract.** Write `schemas/cdd/spec.cue` (`#CellSpec`:
`{version, contract (producer-attributed required_evidence), protocol_id,
profile, params, alpha:{skills}, beta:{skills}}`) and the cds overlay
(`#CDSCellSpec` pinning `protocol_id` + param domains). Hand-write a
`cds.spec.json` for a real issue; `cue vet` it. *Proves the data shape before
any Go or compiler.* (`budget` was removed as decorative per Pi #32 D5; it
returns only when it is actually enforced.)

**Phase 1 — CLI 0: one local episode, honest receipt (no cognition).**
`internal/cellspec` loader (strict parse → kernel `Spec`), and `cn cell run
--contract <path|-> [--param k=v]` that fills holes and calls `RunEpisode`.
Emits a **generic** `cnos.cellkernel.episode-envelope.v0` with `execution_mode`
and `protocol_validated=false` — the declared `protocol_id` is provenance, never
a validated CDS claim. A `profile` (explicit; no default) selects the builtin
seat pair: `stub` (a non-authoritative **`simulated`** smoke run, exit 3) or
`bool` (a real mechanical episode where β **independently verifies** α's
artifact). **Zero GitHub/network.** Contract is frozen + hashed; evidence is
runtime-authenticated + producer-attributed + UTF-8 + size-bounded; β receives a
runtime-owned `BetaInput` review surface. The terminal **envelope re-verifies
whole** — `cellkernel.VerifyEnvelope` re-derives every field
(`verdict←V(receipt)`, `decision←δ(verdict)`, `status←(decision, mode)`), the
`resolved_spec` makes `resolved_spec_hash` recomputable, identity is
per-invocation and fail-closed, and integrity failures route to `rejected` not
`needs_repair`. A CI job (`.github/workflows/cell-schema.yml`) vets the CUE
schemas + actual `cn cell run` output against a shared positive/negative corpus
(Pi #31–#33 + PR-#718 β).

**Phase 2 — skill resolution.** The `$PATH`-like resolver (`value → skill ref`)
+ required/optional/default hole logic; `cn cell run` errors on unfilled
required holes with a Unix-style usage line. Seats load the resolved skills
(cognition still stubbed). *Proves parameters map to real skills.*

**Phase 3 — rented α + CDS profile.** First `internal/dispatch.Backend` adapter
behind the GitHub-free provider port; trivial escalation predicate ("compiled
implementation absent → rent α", Pi Q2). Expose `cn cds run --issue N --contract
<path|->` (`--issue N` is identity/output metadata only — no hidden `gh`). One
real bounded CDS episode closes locally. *This is #717/F.*

**Phase 4 — the surface + compiler (sugar, later).** `cn cell compile main.cell
→ spec.json` (Go `participle`, ~150 lines; CUE stays the independent oracle) +
author `cnos.cds/main.cell`. Pure sugar over Phase-1 JSON — **optional until it
earns its keep**, and explicitly *not* on the frozen S4/S6 critical path.

**Horizon (out of this pass):** `Drive` bounded-repair loop (Case 4); the
`cds-dispatch` parent (`let!` = α-proposes-child / runtime-executes, Case 5);
rewiring the GitHub wake to a thin adapter that fills holes + lands artifacts.

## Boundary the kernel never owns (Pi β #31, C3)

The episode kernel owns no GitHub, ref, PR, branch, cursor, writer-locality, or
custody policy. CLI and GitHub are **invokers/projections**. Persistence
defaults to content-like typed-receipt custody; CHAIN (`.cdd` publication)
remains an opt-in mechanism a CLI adapter may drive — never the kernel.

## Alignment with the frozen shipping plan

This migration instantiates, it does not reopen, the frozen plan
(`msg-cn-pi-cnos-shipping-plan-lock-27`): Phase K + Phase 1 = S4 core + S6 CLI
scaffolding; Phase 3 = the #717/F local CDS field proof (S7). The `.cell`
surface (Phase 4) is a CNOS addition layered *after* the mechanical CLI, not a
plan change. Repair, composition recursion, and the GitHub thin-adapter (S8)
follow only when their preceding cases have executable evidence.

## Ownership split

- **cnos.cdd** owns the *substrate*: the `#CellSpec` schema, the compiler, the
  loader, the runner (`cn cell compile` / `cn cell run`), and the skill-path
  resolver. It is the language + kernel.
- **cnos.cds** owns a *program written in it*: `main.cell`, the params-domain
  overlay, and its α/β skills. cds is one profile; cdw/cdr are siblings of the
  same shape.

### The abstract cell (cdd) vs. the concrete cells (cds/cdr/cdw)

The "abstract cell" is a **type cdd owns**, not a base program to inherit. It is
`schemas/cdd/spec.cue #CellSpec` (+ the kernel + the empty cell). Specializations
are unifications, exactly mirroring the receipt side:

```
#CDSCellSpec: cdd.#CellSpec & { protocol_id: "…cds…", params: {language: …} }
   ⟺  #CDSReceipt: cdd.#Receipt & { protocol_id, evidence_refs }
```

There is deliberately **no language-level inheritance** (no `cdd/main.cell` base
that cds `extends`): a cell has no shared *body* to inherit — its skills and
protocol are precisely what a specialization supplies, and the shared machinery
is the kernel. The abstraction lives in the type + the kernel.

cdd's canonical instance of the abstract cell is the **empty / identity cell**
(`schemas/cdd/fixtures/empty-cell-spec.json`) — the smallest well-formed
`#CellSpec`, the data analogue of the kernel's `EmptySpec`/`cell-0`, and the
runner's reference. It runs today (`cn cell run --contract …/empty-cell-spec.json`
→ accepted); at Phase 4 it becomes `cdd/main.cell`. The generic cdd cell and a
concrete cds cell close through the **same kernel**; only the overlay differs.
