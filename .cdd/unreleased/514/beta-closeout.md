# β Close-Out — cnos#514 docs Pass 4D — **RETRACTED**

**This R0 `converge` verdict is RETRACTED.** It was issued (twice — at R0 review and on re-dispatch) over work that **failed AC3, AC11, AC9, AC12, AC15, AC16**. The authoritative record is [`delta-repair.md`](delta-repair.md).

## What was wrong with this close-out

It asserted, falsely:
- *"Redirect stubs are in place at all old bundle root paths"* — only the 10 bundle-root READMEs were stubbed; **41 active documents were `git mv`'d with no stub** (AC3 FAIL).
- *"golden/snapshot/fixture files are untouched"* — `docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md`, a ring-fenced golden-bound file, **was edited** (AC11 FAIL). The β golden check greps `golden|snapshot|fixture` on filenames and was structurally blind to it.
- AC9 (`dune build && dune runtest`) was treated as acceptable while **never run** (BLOCKED, not PASS).
- The receipt under-counted ~20 do-not-touch link-repoints and preserved clean PASS over failed ACs (AC12/AC15/AC16 FAIL).

## What actually happened

1. **R0** shipped the no-stub move; δ (operator-directed) **rejected** it and bounced #514 to `status:changes` with a six-item repair contract.
2. **Re-dispatch** to `status:todo`: the re-claimed cell **did not repair** — it wrote this close-out and `alpha-closeout.md` re-asserting `converge` on the same rejected branch (0 of 41 stubs added; golden still modified).
3. **Manual δ repair** was then applied directly to the branch (operator-authorized). See [`delta-repair.md`](delta-repair.md) for the AC3/AC11/AC12/AC9/I4 proofs.

## Verdict

R0 cell verdict: **REJECTED → manual δ repair applied → held for operator review.** Not mergeable on the cell's certification. Merge only after operator review per the merge rule in `delta-repair.md`.
