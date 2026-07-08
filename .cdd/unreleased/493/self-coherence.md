§R0

## Gap

**Issue:** cnos#493 — "label-doctor: install + verify cnos.core canonical labels in repos" (live-discovery follow-up to cnos#491).

**Mode:** `design-and-build` (per issue header and γ scaffold).

**Branch:** `cycle/493`, based on `main@31d7ddfa53efb67f2d996b085f4d92986b585ef9` at γ scaffold time; rebased onto `origin/main@43b81295ff8172caceaf9c0795d03b5261a606ff` before this readiness signal (see §Review-readiness for the rebase record — `origin/main` had advanced via three automated `board-map` regeneration commits since scaffold time; the rebase was clean, no conflicts, diff unaffected).

**Incoherence being closed:** `src/packages/cnos.core/labels.json` (schema `cn.labels.v1`, 8 entries: 7 `status:*` lifecycle labels + `dispatch:cell`) declares the canonical cnos.core label set. The cnos repo SHOULD carry this full set after install, but no mechanism materializes or repairs it — labels get auto-created by GitHub with default colors/empty descriptions the first time an agent applies one (as happened with `status:review` during cnos#491), and nothing detects or corrects the drift. Two pieces of already-shipped runtime code assume this mechanism exists and name cnos#493 explicitly when it doesn't:

1. `src/go/internal/repoinstall/repoinstall.go`'s `ensureCanonicalDispatchLabels()` (the `cn repo install --dispatch cds` path) was a stub that unconditionally returned an error naming cnos#493.
2. `src/go/internal/cell/cell.go`'s `Returner.preflightTargetLabel()` (the `status:review → status:changes` review-return transition) has a runtime error whose fix-it text literally says "Run label-doctor before retrying" — a command that did not exist anywhere in the repo.

**What this cycle builds** (per the γ scaffold's `## α prompt` and 7-axis Implementation contract): a new Go module `src/packages/cnos.core/commands/label-doctor/` implementing a canonical-label audit+repair engine against `labels.json`, using dependency-free GitHub REST primitives (no `gh` CLI shellout, mirroring `cnos.issues/commands/issues-fsm/fetch.go`'s style); a thin `cn label-doctor` kernel subcommand wired per the dispatch-boundary convention; replacement of `ensureCanonicalDispatchLabels()`'s stub body with a real in-process call into the new package; `go.work` + CI wiring (module test step, CLI ergonomics smoke-test entry, a CI guard that actually exercises both "fails on drift" and "passes on clean" directions); and a committed live audit + repair artifact.

Full detail, the pinned 7-axis implementation contract, the per-AC oracle, and the source-of-truth table are in `.cdd/unreleased/493/gamma-scaffold.md` (γ's dispatch scaffold, read in full before any code was written).
