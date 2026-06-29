# OCaml thread reference

The OCaml runtime was the original/reference implementation of the CN thread model.
It is **no longer the active runtime implementation** on `main`.
Active runtime and CLI work lives in Go (`src/go/`, the `cn` binary).

The archived OCaml reference is preserved at:

- branch: `legacy/ocaml-thread-reference`
- tag: `ocaml-thread-reference-2026-06-29`

Use it as a reference artifact for thread semantics, **not** as an active build target.

See [`docs/reference/legacy/OCAML-THREAD-REFERENCE.md`](../../docs/reference/legacy/OCAML-THREAD-REFERENCE.md)
for what the OCaml code was authoritative for and what it is not.
