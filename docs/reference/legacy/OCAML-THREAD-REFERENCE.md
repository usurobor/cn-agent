# OCaml thread reference (archived)

**Status:** archived reference implementation. Not an active build or runtime target on `main`.

The OCaml codebase was the original implementation of cnos: the CN thread model, the
typed protocol FSMs, the CN Shell / op executor, the N-pass orchestrator, and the first
`cn` CLI. The active runtime and CLI are now the **Go** implementation (`src/go/`, the
`cn` binary); the OCaml tree has been removed from `main` and frozen.

## Where it lives

- **branch:** `legacy/ocaml-thread-reference`
- **tag:** `ocaml-thread-reference-2026-06-29` (immutable archival baseline, commit `2d457976`)

Check out the branch or tag to read or build the OCaml implementation (`dune build`,
`dune runtest`). Nothing on `main` builds it.

## What the OCaml code is authoritative for

As a **reference for thread semantics** — the questions "how was the CN thread model
originally specified in code?" and "what did the typed FSMs enforce?":

- **Thread / Actor / Sender / Receiver FSMs** — `src/ocaml/protocol/cn_protocol.ml` (pure, compile-time-checked state machines).
- **Transport / packet model** — `src/ocaml/transport/`.
- **CN Shell, op executor, N-pass orchestration** — `src/ocaml/cmd/cn_shell.ml`, `cn_executor.ml`, `cn_orchestrator.ml`.
- **Pure-vs-IO split** — `src/ocaml/lib/` (pure parsers) vs `src/ocaml/cmd/` (IO). The Go port mirrors this split; several Go files cite their OCaml origin in provenance comments.

## What it is NOT

- **Not the active runtime.** The shipped runtime and CLI are Go (`cn`).
- **Not a build target on `main`.** `dune` / `opam` are not required for any mainline change; CI does not build OCaml.
- **Not editable as part of ordinary work.** Mainline changes must not modify OCaml implementation semantics. If a future change genuinely needs to alter the OCaml reference semantics, it belongs in the `legacy/ocaml-thread-reference` branch or a dedicated archival-reference update — never as an incidental docs/runtime/migration change.

## Relationship to the contract fixture

`tests/fixtures/protocol-contract.json` (an OCaml-era test fixture) is **retained on `main`**
because the I2 `protocol-contract-check` gate diffs it against the canonical
`docs/reference/schemas/protocol-contract.json`. It is data, not OCaml code. The OCaml
test tree is fully removed; this fixture was relocated from the old `test/cmd/` path to
`tests/fixtures/` (cnos H1) so no `test/` directory lingers at the root.

---

*OCaml remains part of cnos history and thread semantics, not part of the active mainline runtime.*
