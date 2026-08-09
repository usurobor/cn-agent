# Cell Runner — the case ladder

**Status:** Current reference for the single-episode cell runner. Companion to
the kernel doctrine `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md`
(CCNF) and the migration plan `CDS-CELL-MIGRATION.md`. Grounded in the code:
`src/go/internal/cellkernel`, `src/go/internal/cellspec`, and `cn cell run`.
This document reflects the kernel hardened through Pi β #31–#33; earlier
wording lives in Git history only.

## One kernel, every case

`RunEpisode(ctx, spec, …opts)` runs the five-step closure once, at one scope:

```
matter   := α.produce(contract)          α: open seat (may rent cognition)
review   := β.review(betaInput)          β: open seat (may rent cognition)
receipt  := γ.close(record)              γ: kernel-owned, mechanical
verdict  := V(record, receipt)           V: kernel-owned, mechanical
decision := δ.decide(receipt, verdict)   δ: kernel-owned, mechanical
```

Only **α and β** are open (customizable / cognition-rentable). **γ, V, δ are
kernel-owned functions** — not injectable, so a caller cannot substitute a γ
that certifies its own receipt. Every case on the ladder is the *same*
`RunEpisode`; only the α/β *fills* change. A repair driver and a composition
orchestrator invoke this kernel repeatedly — they are not extra seats.

## Outcomes

`RunEpisode` returns a `Closure` whose `Status` is terminal
(`accepted` | `degraded` | `rejected` | `simulated`) or non-terminal
(`needs_repair`, the parent stays open — the closure then carries a
`RepairRequest` derived from the verdict). A seat that returns an **error**
is a malfunction (the episode does not close). A review with `Pass=false` is
**contract-unmet** (closes `needs_repair`). Integrity failures fail closed
to `rejected`.

## Invariants the kernel enforces (FIDO/functional doctrine)

Governing rule (`msg-cn-pi-cnos-cell-runner-fido-functional-44`,
operator-ratified): **no mutable shared episode state.** Each station is a
pure-shaped function invoked with exactly the immutable data it needs,
returning one typed value. Structural isolation of the untrusted seats is the
primary safety mechanism — not the trusted runtime proving its own internal
steps to itself.

- **I1 — immutable seat scopes.** `α: AlphaInput → Result<AlphaOutput>`,
  `β: BetaInput → Result<BetaOutput>`. Each seat sees an isolated frozen
  contract copy; β additionally receives a fresh **projection** of sealed α
  output (copies) — never α's live scope, session, or a shared episode object.
- **I2 — sealed results.** The runtime seals each return (`sealAlpha`,
  `sealBeta`) before it can cross scope. Sealed values carry unexported state:
  no seat or external caller can construct or mutate one; a β that mutates its
  projection cannot reach the sealed original.
- **I3 — positional ownership.** The runtime knows a value came from α because
  it invoked α and received the return. Seats return candidates `{id, kind,
  text}` only — there are no producer roles, execution ids, hashes, verdicts,
  or status fields for a seat to forge. A required α artifact is satisfied only
  by an artifact sitting under `record.alpha`.
- **I4 — one record, one digest.** `compose` builds ONE immutable
  `EpisodeRecord` (identity, mode, resolved spec, contract, both stations,
  matter, review, policy); γ serializes it with ONE scope-lift digest.
  `VerifyClosure` is the single verification boundary: the digest recomputes,
  and `verdict←V(receipt)`, `decision←δ(verdict)`, `status←lift(...)` re-derive
  purely. No overlapping proof objects.
- **I5 — identity is per-invocation and fail-closed.** The identity tuple is
  minted through one error-returning op and must be non-empty and pairwise
  distinct before α runs.
- **I6 — typed failure routing.** Only `contract_unmet` may become
  `needs_repair`; integrity failures (`invalid_record`/`invalid_identity`) fail
  closed to `rejected` — never the α repair path. A stub run is
  non-authoritative `simulated` (exit 3), never accepted authority.
- **I7 — gated in CI.** `cell-schema.yml` vets the CUE schemas + actual
  `cn cell run` output (accepted / needs-repair / simulated) against a shared
  positive/negative corpus.
- **I8 — the kernel guards its own boundary.** Spec validation (non-empty id,
  unique/valid required refs, bounded cardinality); matter/review/artifact
  output size-bounded (per-item + aggregate); artifact text is explicit UTF-8
  (base64 is a future extension); cancellation checked between stations.

**Composition (future cases) is functional:** `map`/`traverse` a child cell
over immutable inputs, `zip` independent results, `bind` the next child
contract from a sealed previous result, `fold` a parent projection from an
append-only sequence of results. Children never write upward or sideways.

## The ladder (implementation order; Pi #32 D5)

- **Case 0 — empty.** `NoopAlpha` + `AcceptBeta`; no required evidence.
  Terminates `accepted`. The runner smoke reference (`cell-0`).
- **Case 1 — one-shot mechanical (bool).** α produces a bool; β **independently
  verifies** it from its bundle (non-tautological). `value=true → accepted`,
  `value=false → needs_repair`. No repair loop.
- **CLI 0 — local runner.** `cn cell run --contract <path|-> [--param k=v]`
  fills parameter holes, runs one episode, emits a generic
  `cnos.cellkernel.episode-closure.v0`, exit `0/1/2/3` (3 = `simulated`).
  Zero GitHub/network.
- **Case 2 — rented α, mechanical β.** First cognition behind α (a provider
  seam); β still mechanical. (Phase 3 / #717-F; held until CI + Pi converge.)
- **Case 3 — rented α and β.** Full single-episode CDS. V validates
  evidence/bindings; it never re-judges β's prose.
- **Case 4 — bounded repair driver.** A `Drive` loop invokes the same episode
  kernel under an attempt budget on `needs_repair`. Not a new algorithm.
- **Case 5 — composition.** Parent-α **proposes** child contracts; the
  **runtime executes** child episodes; accepted child receipts become parent
  matter; parent-β reviews the composition. α never calls the kernel itself.

Degraded/override, reject, budget exhaustion, and malfunction are **cross-cutting
outcome tests**, not new rungs.

## Parameters → skills (Unix typed holes)

A parameter is a typed hole resolved like `$PATH`: `skill`-kind values resolve to
a skill and splice into a seat via `$name`; `value`-kind values are scalars
passed to a builtin profile. Required vs optional-with-default = positional vs
flag; a closed `domain` makes a typo fail resolution. The invoker (CLI now, a
parent cell later) fills the hole; the cell body never changes.

## What the runner does not own (custody boundary)

The kernel owns no GitHub, ref, PR, branch, cursor, writer-locality, or custody
policy. `cn cell run` reads a local path or stdin and writes the closure JSON
to stdout. CLI and GitHub are invokers/projections over this same engine.
