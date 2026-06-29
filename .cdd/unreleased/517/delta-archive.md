# δ archival receipt — cnos#517: archive OCaml runtime as legacy CN thread reference

**Author:** δ (Sigma), operator-authorized 2026-06-29. **Manual δ** (not dispatched; cnos#516 open).
**Branch:** cycle/517. **Base:** main @ `2d457976`.

> **This is an archival move, NOT a semantic rewrite of the CN thread model (AC10).** No `.ml`/`.mli`
> logic was edited. The OCaml tree was relocated to a frozen reference (branch + tag) and removed from
> `main`; the only content edits are a stub, a reference doc, and three framing updates (README,
> CONTRIBUTING, DESIGN-CONSTRAINTS). The thread semantics the OCaml code expresses are unchanged and
> remain readable on the legacy branch/tag.

## AC verification

| AC | Status | Evidence |
|---|---|---|
| AC1 — legacy branch preserves OCaml tree | **PASS** | `legacy/ocaml-thread-reference` pushed @ `2d457976` (full tree, pre-removal). `git ls-remote` confirms. |
| AC2 — immutable tag at archival baseline | **PARTIAL — operator action needed** | Tag `ocaml-thread-reference-2026-06-29` **could not be pushed**: this environment's git proxy rejects all `refs/tags/*` pushes (`remote end hung up`), and the MCP toolset has no tag/ref/release-creation call. The branch preserves the bytes; the immutable tag must be created by the operator (one-liner below). |
| AC3 — main no longer treats OCaml as active runtime | **PASS** | `src/ocaml/**`, `dune-project`, `cn_agent.opam`, all `test/**/*.ml{,i}` + `dune`, `test/cram/**` removed (119 files). No `.ml/.mli/dune/opam` remain on `main`. |
| AC4 — `src/ocaml/` only a README stub | **PASS** | `src/ocaml/README.md` is the sole remaining file, pointing to branch/tag + the reference doc. |
| AC5 — legacy reference doc | **PASS** | `docs/reference/legacy/OCAML-THREAD-REFERENCE.md` states what the OCaml code is authoritative for (thread FSMs, transport, CN Shell, pure/IO split) and what it is not (not active runtime, not a build target, not editable in ordinary work). |
| AC6 — README no longer implies OCaml active | **PASS** | README src-tree line updated: "archived OCaml thread reference (stub → legacy/ocaml-thread-reference)". Go = active CLI/runtime. |
| AC7 — gates no longer require dune for mainline | **PASS** | `CONTRIBUTING.md` test step `dune runtest` → `go test ./...`. No CI job builds OCaml (verified: go/binary-verify/package-verify use the Go `cn` binary; I1/I2/I4/I5/I6 don't touch dune). The 4D-era "OCaml gate (AC9)" is retired for mainline work. |
| AC8 — no frozen records rewritten | **PASS** | No `.cdd/` evidence, release records, or versioned snapshots touched (this receipt aside). |
| AC9 — active links resolve | **PASS (pending CI confirm)** | No markdown links pointed at `src/ocaml/` or `test/*.ml` (only backtick code-spans), so removal creates no broken links. New doc links (stub→reference, DESIGN-CONSTRAINTS→reference) verified to resolve. I4 expected to carry only the inherited baseline. |
| AC10 — receipt records archival-not-rewrite | **PASS** | This receipt. |

## The one active coupling — handled

I2 (`protocol-contract-check`, `build.yml:223`) diffs `test/cmd/protocol-contract.json`, a JSON fixture
inside the OCaml test tree. It was **intentionally retained on `main`** (data, not OCaml code), so I2
stays green with **zero workflow change**. Recorded in the reference doc.

## Deferred to 4E (deliberate, not a miss)

`docs/beta/architecture/ARCHITECTURE.md` and `docs/architecture/security/SECURITY-MODEL.md` describe the
system via `cn_*.ml` backtick code-spans (the thread-semantics reference). These are not links and do not
break. Rewriting those narratives with the settled "active = Go, reference = OCaml" truth is **4E's job** —
which is exactly why this archival runs first.

## Operator action required (AC2)

Create the immutable tag from an environment that permits tag pushes (or via the GitHub UI / a release):

```
git tag ocaml-thread-reference-2026-06-29 2d457976dfb204ae487a1a2ad2f0cf99dd62d1de
git push origin ocaml-thread-reference-2026-06-29
```

(`2d457976` is the archival baseline = current `main` head at filing, with the full OCaml tree intact via
`legacy/ocaml-thread-reference`.)
