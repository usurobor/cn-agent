# α Close-Out — cnos#514 docs Pass 4D — **SUPERSEDED**

**This close-out is superseded.** It was written on re-dispatch to finalize the rejected R0 implementation and does **not** reflect a compliant pass. The authoritative record is [`delta-repair.md`](delta-repair.md).

## Honest status of the α implementation

The R0 implementation moved the 10 bundles correctly **but**:
- abandoned the stub model — 41 active documents were `git mv`'d with **no redirect stub** at the old path (**AC3 FAIL**);
- edited the ring-fenced golden `docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md` (**AC11 FAIL**);
- left moved-doc version-snapshot links pointing into the empty new bundle dir (I4 risk);
- could not run `dune build && dune runtest` (**AC9 BLOCKED**, recorded by α/β as acceptable — incorrect).

Sound and retained: the move targets, the non-destructive `reference/schemas/` merge, `build.yml:223` path-only repoint, and the comment-only Go/OCaml citation repairs.

## Disposition

R0 was **rejected by δ**; re-dispatch **did not repair**; **manual δ repair** was applied directly to the branch (stubs at all 41 old paths, golden reverted, snapshot links repointed, receipt corrected). See [`delta-repair.md`](delta-repair.md) for the proofs and the residual AC9/CI gates. Held for operator review.
