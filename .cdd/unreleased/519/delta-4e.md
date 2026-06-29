# δ receipt — cnos#519: docs Pass 4E — architecture docs → docs/architecture/

**Author:** δ (Sigma), operator-authorized 2026-06-29. **Manual δ** (not dispatched; cnos#516 open).
**Branch:** cycle/519. **Base:** main @ `64945dc3`.

> Migration + narrow truth-alignment pass. The architecture **theory** (relations, module
> layering, FSMs) is unchanged; the only content edit is a status banner reframing OCaml as the
> archived reference and Go `cn` as active.

## AC verification

| AC | Status | Evidence |
|---|---|---|
| AC1 Scope | PASS | Only `ARCHITECTURE.md` + `PACKAGE-SYSTEM.md` moved `beta/architecture/` → `architecture/`. |
| AC2 History | PASS (stub-model) | `git mv` used; per the kept-stub-model policy the old paths then receive stubs, so git records A(new)+M(old-stub) rather than R — consistent with 4B/4C. |
| AC3 Stubs | PASS | `# Moved` stubs at both old active paths → `../../architecture/<file>`. |
| AC4 Frozen | PASS | `docs/beta/architecture/3.14.4/` untouched (not in changed set). |
| AC5 Reader index | PASS | `docs/architecture/README.md`: ARCHITECTURE link → local; PACKAGE-SYSTEM added. |
| AC6 Old bundle index | PASS | `docs/beta/architecture/README.md` rewritten as a redirect/index → `../../architecture/`, with the frozen `3.14.4/` note preserved. |
| AC7 Active links | PASS | 5 inbound links repointed (see below). Post-sweep: 0 active links to old paths remain. |
| AC8 Relative links | PASS | Moved files: 7 (ARCHITECTURE) + 2 (PACKAGE-SYSTEM) `../../` → `../`; local link-check = 0 broken. |
| AC9 Architecture truth | PASS | Banner at top of moved `ARCHITECTURE.md`: active = Go `cn`; `cn_*.ml` = archived CN thread reference (links `docs/reference/legacy/OCAML-THREAD-REFERENCE.md`, branch+tag); `dune`/OCaml not a mainline gate. Body theory left intact (no broad rewrite). |
| AC10 No code/workflow | PASS | `git diff --name-only` shows no `src/`, `test/`, `.github/`, package, or binary changes. |
| AC11 Golden protection | PASS | No goldens; no `docs/gamma/conventions/` touched. |
| AC12 I4 | PASS (pending CI confirm) | Local broken-link check across moved + repointed files = 0; expect inherited baseline only. |
| AC13 No prose drift | PASS | Only the AC9 banner added; no other prose/semantic edits in the moved docs. |
| AC14 Hidden/bidi | PASS | Sweep over all changed files (incl. this receipt) clean. |
| AC15 Receipt honesty | PASS | This receipt; do-not-touch repoints listed; no clean PASS over a repaired AC. |

## AC7 inbound link repoints

In-scope (architecture home):
- `docs/architecture/README.md` — ARCHITECTURE link → local `ARCHITECTURE.md`; added `PACKAGE-SYSTEM.md`.
- `docs/architecture/security/SECURITY-MODEL.md:158` — `../../beta/architecture/ARCHITECTURE.md` → `../ARCHITECTURE.md`.

**Do-not-touch files repointed under δ override (pure active-link repoints only, recorded per AC15):**
- `docs/evidence/AUDIT.md:18` — link text + target `../beta/architecture/ARCHITECTURE.md` → `../architecture/ARCHITECTURE.md`.
- `docs/reference/cli/CLI.md:244` — `../../beta/architecture/ARCHITECTURE.md` → `../../architecture/ARCHITECTURE.md`.
- `docs/reference/runtime/AGENT-RUNTIME.md:2366` — `../../beta/architecture/ARCHITECTURE.md` → `../../architecture/ARCHITECTURE.md`.

Exempt (classified, not edited):
- `docs/beta/governance/DOCUMENTATION-SYSTEM.md:297` — a historical migration-tracking table cell (`beta/ARCHITECTURE.md | beta/architecture/ARCHITECTURE.md | Migrated (#89)`), plain-text path citations, not a markdown link. Left as historical record.

## δ finding — PACKAGE-SYSTEM.md is a retired redirect, not an architecture paper

`docs/beta/architecture/PACKAGE-SYSTEM.md` is a **retirement notice** (since #180) already pointing to
`docs/reference/packages/PACKAGE-SYSTEM.md` as the authority — not the "architecture-level package-system
paper" the spec described. Moved per AC1 (it is an active file at an old path being emptied); the moved
file at `docs/architecture/PACKAGE-SYSTEM.md` remains a retirement breadcrumb (its two links depth-fixed).
No merge/collision with `docs/reference/packages/PACKAGE-SYSTEM.md` (distinct path; the retired file already
defers to it). **Flagged for operator**: if you'd rather not feature a retired redirect in the reader-intent
architecture home, deleting it (leaving only the old-path stub → reference/packages) is a trivial follow-up.

## Frozen / not-moved

`docs/beta/architecture/3.14.4/` stays frozen. No release records, snapshots, goldens, or code touched.

## Residual

I4/I2/Go/Package/Binary via PR CI — expect I4 inherited baseline only; no new code so Go/I1/I2/Package/Binary
unaffected. Merge under operator review.
