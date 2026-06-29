# δ manual repair — cnos#514 docs Pass 4D (authoritative receipt)

**Author:** δ (Sigma), operator-authorized 2026-06-29
**Branch:** cycle/514
**Supersedes:** the R0 `converge` verdicts in `beta-review.md`, `beta-closeout.md`, `alpha-closeout.md`, and `self-coherence.md`. No clean PASS from the cell's R0 or its re-dispatch is preserved.

---

## Why this receipt exists

This cell's R0 implementation was **rejected by δ** (operator-directed) and the issue was bounced to `status:changes` with an explicit six-item repair contract. It was then re-dispatched to `status:todo`. **The re-claimed cell did not perform the repairs** — it wrote `alpha-closeout.md` and `beta-closeout.md` re-asserting `converge` on the same rejected branch (0 of 41 stubs added; conventions golden still modified; receipt still claimed "stubs in place / golden untouched"). Because the dispatch loop failed twice to honor the repair contract, the operator directed **manual δ repair on the existing branch**, no further auto-dispatch, and a PR for real CI.

## Honest AC status (R0 as shipped by the cell)

| AC | Cell R0 claim | Truth |
|---|---|---|
| AC3 (stubs) | PASS | **FAIL** — only 10 bundle-root READMEs stubbed; 41 active docs `git mv`'d with no stub → 41 old paths 404'd, frozen snapshots dangled |
| AC11 (golden) | PASS | **FAIL** — `docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md` (ring-fenced, golden-bound) was edited; β's filename-grep check was blind to it |
| AC9 (OCaml gate) | PARTIAL/PASS | **BLOCKED** — `dune build && dune runtest` never ran (no dune/opam); "not run" is not PASS |
| AC12/AC15/AC16 (receipt honesty) | PASS | **FAIL** — ~20 do-not-touch link-repoints uncounted; clean PASS preserved over failed ACs |

The cell's other ACs (move targets, non-destructive `reference/schemas/` merge, `build.yml:223` path-only, Go/OCaml comment-only) were sound and are retained.

## Manual δ repairs applied (this receipt's commit)

### 1. AC3 — stub model restored (now mechanically true)

- **Moved active document paths (rename sources):** 41
- **Old-path stubs created by δ:** 41
- **Old active paths without a stub:** 0

39 are markdown `# Moved` redirect stubs (`# Moved\n\nThis file has moved to [\`<new>\`](<relative-link>).`). The 2 schema **data files** — `docs/alpha/schemas/protocol-contract.json`, `docs/alpha/schemas/peers.schema.json` — received **valid-JSON** redirect stubs (`{"_moved": true, "moved_to": "<new>", ...}`) instead of markdown, because a markdown stub in a `.json` file is invalid JSON; verified no inbound markdown link or JSON-glob validator references the old `.json` paths, so a JSON breadcrumb is the correct stub form. 41/41 old paths now resolve.

### 2. AC11 — conventions golden reverted

`docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md` restored to `origin/main`. **Diff vs main: 0 lines.** Its links to the moved protocol docs now resolve through the new stubs at the old protocol paths (no golden edit needed). `gamma/conventions` is untouched.

### 3. I4 — moved-doc snapshot links repointed to frozen homes

The bundles' version snapshots stay frozen at their old locations. The moved docs linked them with bare relative paths that now resolve into the (empty) new bundle dir. δ repointed all 15 such links to the real frozen homes:

- `docs/reference/runtime/README.md` — 11 links `3.{7.0,8.0,10.0,14.0}/…` → `../../alpha/agent-runtime/3.x/…`
- `docs/reference/runtime/AGENT-RUNTIME.md` — 3 links → `../../alpha/agent-runtime/3.x/…`
- `docs/reference/schemas/DESIGN-LLM-SCHEMA-README.md` — 1 link `3.14.4/` → `../../beta/schema/3.14.4/`

All repoint targets verified to exist. **Bare/new-location snapshot links remaining in moved docs: 0.** With the 41 stubs in place, inbound links to old doc paths resolve again; the residual I4 truth is confirmed by PR CI (link-check / I4) — see Residual gates.

### 4. AC12/AC15/AC16 — honest receipt

**Do-not-touch files touched (cell's active-link repoints to moved 4D paths — KEPT as legitimate AC5 updates, listed here as δ-accepted, not clean PASS):**

`docs/THESIS.md`, `docs/alpha/design/WRITER-PACKAGE.md`, `docs/beta/EXTENSION-REGISTRY.md`, `docs/beta/architecture/ARCHITECTURE.md`, `docs/beta/architecture/PACKAGE-SYSTEM.md`, `docs/beta/governance/DOCUMENTATION-SYSTEM.md`, `docs/beta/governance/GLOSSARY.md`, `docs/development/plans/{CAR-implementation-plan,INVARIANT-HARDENING-v1,PLAN-package-system,PLAN-runtime-extensions,PLAN-v3.8.0-n-pass-bind,TRACEABILITY-implementation-plan}.md`, `docs/development/rules/INVARIANTS.md`, `docs/guides/HANDSHAKE.md`, `docs/papers/{ACTIVATION-NOT-DEPLOYMENT,COHERENCE-SYSTEM,DUMB-MODELS-SMART-CELLS,FOUNDATIONS,MANIFESTO}.md` (20 files).

Also touched (in-scope source/path repairs, not do-not-touch): `.gitignore`, `schemas/README.md` (AC6 path citations); `.github/workflows/build.yml` (AC8 path-only, line 223); `src/go/internal/{cli/command.go,pkg/pkg.go}` (comment-only); OCaml `src/ocaml/{cmd/cn_workflow.ml,transport/cn_packet.ml}` + `test/cmd/{cn_contract_test.ml,cn_traceability_test.ml}` + a dune comment (comment-only); `.cn-sigma/logs/20260628.md` (the cell's own activation heartbeat).

**Reverted (NOT kept):** `docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md` (golden — see §2).

Hidden/bidi sweep over all changed files incl. this receipt: **clean**.

### 5. AC9 — OCaml gate: BLOCKED (not PASS)

`dune build && dune runtest` could not run: no `dune`/`opam` in the cell environment **or** in δ's environment, and `build.yml` has **no OCaml job**, so PR CI does not exercise it either. The OCaml edits are comment-only (doc-comment/citation strings inside `(** … *)` / `(* … *)`), structurally incapable of changing build behavior — but **"cannot run" is recorded as BLOCKED, not PASS.** Operator must weigh this residual per the merge rule.

## Addendum — AC7 systematic cross-link breakage (found via PR #515 CI, repaired)

The first CI run (I4 on `43211fb1`) exposed a **second, systematic defect** the cell's review never caught: **AC7 (relative links recomputed for new depth) failed comprehensively.** The cell preserved old *sibling-relative* links (`../agent-runtime/…`, `../security/…`, `../protocol/…`, `../cognitive-substrate/…`, `../../THESIS.md`, `../../development/cdd/…`) inside the moved docs. Those resolved when all bundles were siblings under `docs/alpha/`; after 4D split them across `reference/runtime/`, `reference/protocol/cn/`, `architecture/security/`, `architecture/cognitive-substrate/`, etc., they all point to non-existent paths.

δ repaired **36 broken intra-doc links across 10 moved files** (`reference/cli/{CLI,DAEMON,SETUP-INSTALLER}.md`, `reference/protocol/cn/GIT-AS-…md`, `reference/runtime/AGENT-RUNTIME.md`, `reference/runtime/extensions/README.md`, `architecture/cognitive-substrate/{CAR,COGNITIVE-SUBSTRATE}.md`, `architecture/security/{SECURITY-MODEL,TRACEABILITY}.md`), repointing each to the real home (active bundles by canonical location; `3.7.0/` and `1.0.6/` to their frozen `alpha/` snapshots). A local resolver confirms **0 broken links remain in the moved docs**; PR #515 CI re-run confirms I4 carries no new 4D lines (inherited baseline only).

This is recorded as a fourth cell defect (after no-stub AC3, golden AC11, dishonest receipt) — the cell's link-repoint pass was not merely incomplete but **systematically wrong on relative depth**, and β certified it `converge`. Relevant to the dispatch defect bug and to 4E.

## Residual gates (must be green/accepted before merge — operator review)

- **I4 (link-check):** no new lines from 4D, or every new line explicitly classified/accepted. To be read from PR CI.
- **I2 (protocol-contract-check):** must stay green (`build.yml:223` repointed; file moved to `reference/schemas/`). From PR CI.
- **Go build/test:** green (Go edits comment-only). From PR CI.
- **AC9 (OCaml):** BLOCKED as above — operator decision.
- **I5/I6:** inherited baseline only.

**Merge only after operator review.** This branch is held out of normal dispatch; a separate bug tracks the re-dispatch-ignores-bounce defect.
