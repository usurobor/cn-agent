# β review capture (envelope) — #671 (external, revision-pinned)

> **Purpose.** Make the β role **observable in git**, and bind the β evidence
> by content. The β judgment for the #671 Planning Cell was performed by a
> **genuinely external reviewer of a different model lineage** ("Codex",
> posting on PR #672 as the repo owner account `usurobor`), content-bound to
> the exact reviewed SHA. κ (Sigma-at-repo) **captures** this external
> evidence into the ledger; κ does **not** author the β judgment.
>
> **This file is the κ-authored capture envelope (metadata only).** The
> **verbatim β payload is a separate file** — `beta-review-source-5076109763.md`
> — which is **byte-identical to the source comment** and carries **no κ
> metadata**. This separation is the R17 repair of external-β finding
> [REQUIRED] "The β artifact is not the full verbatim review γ claims to bind"
> (comment 5076629728).

## Bound β evidence identity (immutable)

| Field | Value |
|---|---|
| Verdict | **CONVERGE** |
| Reviewed matter | branch `wave/671-cell-runtime-doctrine` @ `614829a4682e148d98c70371e600ffdc3fa6386e` (R16) |
| Reviewer | external, different-lineage ("Codex"), posting as `usurobor` (OWNER) |
| Source comment | [PR #672 comment 5076109763](https://github.com/usurobor/cnos/pull/672#issuecomment-5076109763) |
| Created | `2026-07-25T01:38:50Z` |
| **Raw-body SHA-256 (authoritative)** | **`75cdb9b67619570c668574ca27efc4107709f8c5a3f52199484033695ccfe3ae`** |
| **Raw-body length** | **9,894 bytes** |
| Verbatim payload file | `.cdd/unreleased/671/beta-review-source-5076109763.md` |
| `origin/main` observed by reviewer | `41a86cef72437cf1d8a7800aaa96e5a01e305d78` |

## Payload fidelity (mechanically checkable)

The committed payload file `beta-review-source-5076109763.md` is a **byte-exact
reproduction of the source comment's raw body**:

```
$ sha256sum .cdd/unreleased/671/beta-review-source-5076109763.md
75cdb9b67619570c668574ca27efc4107709f8c5a3f52199484033695ccfe3ae  ...
$ wc -c   .cdd/unreleased/671/beta-review-source-5076109763.md
9894 ...
```

- **All nine source `##` sections are present, in order:** Verdict ·
  Reconstructed intent · What Claude produced · Findings · Wave graph
  assessment · Invariant audit · Mechanical-enforceability audit · Operator
  decisions · Recommended next action. (Plus the `# External β Review` title,
  the `**Review target:**` line, and the `### [OBSERVATION]` subsection.)
- **Line endings / trailing byte:** LF line endings, **no trailing newline**
  (this is the form that reproduces the authoritative raw-body hash exactly).
  A fresh reviewer re-fetching comment 5076109763's raw body and hashing it
  gets `75cdb9b6…` / 9,894 bytes, matching this file byte-for-byte.

## What this fixes vs the prior (voided) capture

The prior `beta-review.md` (commit `8ebe5e6b`) captured only the terminal
verdict + Findings + Recommended-next-action **excerpt**, mislabeled it a
"full verbatim body," and computed a hash of the κ **wrapper** (`840188d6…`)
rather than the source body — then the prior γ closeout claimed that wrapper
hash "Matches comment 5076109763," which it did not. R17 corrects all three:
(1) the full body is preserved byte-exact in a payload file; (2) the bound hash
is the **authoritative raw-body** `75cdb9b6…`; (3) κ metadata lives only here,
outside the payload. γ (re-dispatched under R17) binds the raw-body hash, not a
wrapper.

## Sole open item

The only open item in the β review is the **[OBSERVATION]** — a stale README
historical sentence (`README.md:408-411`) — for which the reviewer's own
required repair is **"None for this boundary. Do not mutate the converged R16
matter."** Carried as debt (`self-coherence.md` `## Known debt` item 3),
unrepaired to preserve the R16 freeze.
