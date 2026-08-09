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

`RunEpisode` returns an `EpisodeResult` whose `Status` is terminal
(`accepted` | `degraded` | `rejected`) or non-terminal (`needs_repair`, the
parent stays open). A seat that returns an **error** is a malfunction (the
episode does not close). A review with `Pass=false` is **contract-unmet**
(closes `needs_repair`). An inconsistent (verdict, decision) pair is a typed
`ErrInvalidClosure` — never a returned closed cell.

## Invariants the kernel enforces

- **I1 — trust by receipt, not role.** γ/V/δ are mechanical; cognition is
  rentable only at α/β. The trust surface is deterministic.
- **I2 — the contract is frozen.** At episode start the runtime deep-copies +
  hashes the contract; each seat gets an isolated copy; V/γ bind the frozen
  snapshot. A seat cannot mutate the terms it is judged against.
- **I3 — evidence is runtime-authenticated.** Seats return candidate
  `{id, kind, bytes}` only. The runtime stamps producer role + execution id +
  content digest, creates the `sha256:` ref, and inlines the bytes. α cannot
  mint β's evidence; the canonical `beta_review` is minted by the runtime from
  the actual review.
- **I4 — γ is self-proving.** The runtime builds an authoritative
  `EpisodeRecord`; γ binds it; V recomputes every hash from the receipt's own
  content *and* compares to the record. A γ bug that rewrote any field fails V.
- **I5 — identity is per-invocation.** Each run gets a distinct episode id and
  α/β execution ids (injectable only for deterministic tests). A separate
  `resolved_spec_hash` binds the exact executable spec, so runs differing only
  in resolved input differ in the receipt.
- **I6 — the receipt re-verifies out of process.** `VerifyReceipt` re-checks a
  serialized receipt alone (hashes, content-addressing, uniqueness, required
  evidence + producer authority) — the check a parent runs on what it received.
- **I7 — the kernel guards its own boundary.** A direct `Spec` is validated
  (non-empty id, unique/valid required refs, bounded cardinality); output is
  size-bounded; cancellation is checked between α, β, and closure.

## The ladder (implementation order; Pi #32 D5)

- **Case 0 — empty.** `NoopAlpha` + `AcceptBeta`; no required evidence.
  Terminates `accepted`. The runner smoke reference (`cell-0`).
- **Case 1 — one-shot mechanical (bool).** α produces a bool; β **independently
  verifies** it from its bundle (non-tautological). `value=true → accepted`,
  `value=false → needs_repair`. No repair loop.
- **CLI 0 — local runner.** `cn cell run --contract <path|-> [--param k=v]`
  fills parameter holes, runs one episode, emits a generic
  `cnos.cellkernel.episode-receipt.v0`, exit `0/1/2`. Zero GitHub/network.
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
policy. `cn cell run` reads a local path or stdin and writes a receipt to
stdout. CLI and GitHub are invokers/projections over this same engine.
