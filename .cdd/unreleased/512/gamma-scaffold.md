# gamma-scaffold — cnos#512 — docs Pass 4C

**cell:** cnos#512 — docs Pass 4C: move gamma/cdd, gamma/rules, gamma/plans, gamma/checklists → development/ homes (+ cdd-verify source path)
**protocol:** cds
**cycle-branch:** cycle/512
**base-sha:** 5c12ceba998b434427385913b3f1df24a83cf25e
**mode:** design-and-build (git mv + redirect stubs + link repoint + Go source path update; Go build/test is a hard gate)
**collapse-authority:** β-α-collapse-on-δ (4B precedent — mechanical AC oracle; docs + single Go path-default change)

---

## Source of truth

| Surface | Path | Status |
|---|---|---|
| Move map + sub-cell order | `.cdd/unreleased/508/pass4-move-map.md`, `pass4-subcell-order.md` | Shipped via #509 |
| Go source coupling | `src/packages/cnos.cdd/commands/cdd-verify/run.go:59,161` | Current: `"docs/gamma/cdd"` literal |
| Reader-intent target structure | `docs/README.md` portal + `docs/development/README.md` | Shipped (Pass 1) |
| Stub + accept-with-repair precedent | 4B (#511, commit `5c12ceba`) | Shipped |
| Golden-bound do-not-move | `pass4-move-map.md` §Golden-Bound Bundle Alert | `gamma/conventions/AGENT-ACTIVATION-LOG-v0.md` — NOT in 4C scope |

---

## Surfaces α will touch

### Bundle 1 — docs/gamma/cdd/ (active root-level files only)

Move each file below with `git mv`; leave a stub at the old path:

| Old path | New path |
|---|---|
| `docs/gamma/cdd/CDD-PACKAGE-AUDIT.md` | `docs/development/cdd/CDD-PACKAGE-AUDIT.md` |
| `docs/gamma/cdd/CDD.md` | `docs/development/cdd/CDD.md` |
| `docs/gamma/cdd/DISPATCH-FAILURE-EVIDENCE.md` | `docs/development/cdd/DISPATCH-FAILURE-EVIDENCE.md` |
| `docs/gamma/cdd/GATE-TEMPLATE.md` | `docs/development/cdd/GATE-TEMPLATE.md` |
| `docs/gamma/cdd/ISSUE-CONSOLIDATION-ANALYSIS.md` | `docs/development/cdd/ISSUE-CONSOLIDATION-ANALYSIS.md` |
| `docs/gamma/cdd/KATA-FRAMEWORK.md` | `docs/development/cdd/KATA-FRAMEWORK.md` |
| `docs/gamma/cdd/KATAS.md` | `docs/development/cdd/KATAS.md` |
| `docs/gamma/cdd/OVERVIEW.md` | `docs/development/cdd/OVERVIEW.md` |
| `docs/gamma/cdd/PLAN-TEMPLATE.md` | `docs/development/cdd/PLAN-TEMPLATE.md` |
| `docs/gamma/cdd/POST-RELEASE-EPOCH-v3.12.md` | `docs/development/cdd/POST-RELEASE-EPOCH-v3.12.md` |
| `docs/gamma/cdd/POST-RELEASE-EPOCH-v3.14.md` | `docs/development/cdd/POST-RELEASE-EPOCH-v3.14.md` |
| `docs/gamma/cdd/RATIONALE.md` | `docs/development/cdd/RATIONALE.md` |
| `docs/gamma/cdd/README.md` | `docs/development/cdd/README.md` |
| `docs/gamma/cdd/SELF-COHERENCE-TEMPLATE.md` | `docs/development/cdd/SELF-COHERENCE-TEMPLATE.md` |

**Frozen — STAY (do not move, do not rewrite):**
- All `docs/gamma/cdd/<X.Y.Z>/` version snapshot subdirs (~60, e.g. `3.13.0/`…`3.81.0/`)
- `docs/gamma/cdd/docs/` — historical POST-RELEASE-ASSESSMENT records (PRAs)

### Bundle 2 — docs/gamma/rules/ (active files only)

| Old path | New path |
|---|---|
| `docs/gamma/rules/INVARIANTS.md` | `docs/development/rules/INVARIANTS.md` |
| `docs/gamma/rules/README.md` | `docs/development/rules/README.md` |
| `docs/gamma/rules/RULES.md` | `docs/development/rules/RULES.md` |
| `docs/gamma/rules/CAR-implementation-plan.md` | `docs/development/rules/CAR-implementation-plan.md` |
| `docs/gamma/rules/INVARIANT-HARDENING-v1.md` | `docs/development/rules/INVARIANT-HARDENING-v1.md` |
| `docs/gamma/rules/PLAN-package-system.md` | `docs/development/rules/PLAN-package-system.md` |
| `docs/gamma/rules/PLAN-runtime-contract-v2.md` | `docs/development/rules/PLAN-runtime-contract-v2.md` |
| `docs/gamma/rules/PLAN-runtime-extensions.md` | `docs/development/rules/PLAN-runtime-extensions.md` |
| `docs/gamma/rules/documenting.md` | `docs/development/rules/documenting.md` |
| `docs/gamma/rules/engineering.md` | `docs/development/rules/engineering.md` |
| `docs/gamma/rules/functional.md` | `docs/development/rules/functional.md` |
| `docs/gamma/rules/ocaml.md` | `docs/development/rules/ocaml.md` |
| `docs/gamma/rules/testing.md` | `docs/development/rules/testing.md` |

**Frozen — STAY:** `docs/gamma/rules/3.14.5/`

### Bundle 3 — docs/gamma/plans/ (all active files; no versioned snapshots)

| Old path | New path |
|---|---|
| `docs/gamma/plans/CAR-implementation-plan.md` | `docs/development/plans/CAR-implementation-plan.md` |
| `docs/gamma/plans/INVARIANT-HARDENING-v1.md` | `docs/development/plans/INVARIANT-HARDENING-v1.md` |
| `docs/gamma/plans/PLAN-package-system.md` | `docs/development/plans/PLAN-package-system.md` |
| `docs/gamma/plans/PLAN-runtime-contract-v2.md` | `docs/development/plans/PLAN-runtime-contract-v2.md` |
| `docs/gamma/plans/PLAN-runtime-extensions.md` | `docs/development/plans/PLAN-runtime-extensions.md` |
| `docs/gamma/plans/PLAN-v3.10.0-runtime-contract.md` | `docs/development/plans/PLAN-v3.10.0-runtime-contract.md` |
| `docs/gamma/plans/PLAN-v3.13.0-docs-governance.md` | `docs/development/plans/PLAN-v3.13.0-docs-governance.md` |
| `docs/gamma/plans/PLAN-v3.22.0-eng-lane-clarity.md` | `docs/development/plans/PLAN-v3.22.0-eng-lane-clarity.md` |
| `docs/gamma/plans/PLAN-v3.6.0.md` | `docs/development/plans/PLAN-v3.6.0.md` |
| `docs/gamma/plans/PLAN-v3.7.0-scheduler.md` | `docs/development/plans/PLAN-v3.7.0-scheduler.md` |
| `docs/gamma/plans/PLAN-v3.8.0-n-pass-bind.md` | `docs/development/plans/PLAN-v3.8.0-n-pass-bind.md` |
| `docs/gamma/plans/PLAN-v3.8.0-syscall-surface.md` | `docs/development/plans/PLAN-v3.8.0-syscall-surface.md` |
| `docs/gamma/plans/PLAN.md` | `docs/development/plans/PLAN.md` |
| `docs/gamma/plans/PR-docs-governance-v3.13.0.md` | `docs/development/plans/PR-docs-governance-v3.13.0.md` |
| `docs/gamma/plans/TRACEABILITY-implementation-plan.md` | `docs/development/plans/TRACEABILITY-implementation-plan.md` |
| `docs/gamma/plans/issue-41-pass-b-wiring.md` | `docs/development/plans/issue-41-pass-b-wiring.md` |
| `docs/gamma/plans/documenting.md` | `docs/development/plans/documenting.md` |
| `docs/gamma/plans/engineering.md` | `docs/development/plans/engineering.md` |
| `docs/gamma/plans/functional.md` | `docs/development/plans/functional.md` |
| `docs/gamma/plans/ocaml.md` | `docs/development/plans/ocaml.md` |
| `docs/gamma/plans/testing.md` | `docs/development/plans/testing.md` |

### Bundle 4 — docs/gamma/checklists/ (all files)

| Old path | New path |
|---|---|
| `docs/gamma/checklists/documenting.md` | `docs/development/checklists/documenting.md` |
| `docs/gamma/checklists/engineering.md` | `docs/development/checklists/engineering.md` |
| `docs/gamma/checklists/functional.md` | `docs/development/checklists/functional.md` |
| `docs/gamma/checklists/ocaml.md` | `docs/development/checklists/ocaml.md` |
| `docs/gamma/checklists/testing.md` | `docs/development/checklists/testing.md` |

### Go source — run.go lines 59 and 161

File: `src/packages/cnos.cdd/commands/cdd-verify/run.go`

- Line 59: `a.Bundle = "docs/gamma/cdd"` → `a.Bundle = "docs/development/cdd"`
- Line 161: `--bundle <path>        bundle-relative PRA dir (default: docs/gamma/cdd)` → `(default: docs/development/cdd)`

These are pure literal path-default changes; no behavior change, no new imports, no logic change.

---

## Per-AC oracle

| AC | What β checks | Mechanical verification |
|---|---|---|
| AC1 (Scope) | Only 4C bundle files appear in diff; no 4D/4E/conventions files touched | `git diff --name-status main...HEAD \| grep -v '^[RM].*docs/\(gamma\|development\)/\(cdd\|rules\|plans\|checklists\)' \| grep -v '.cdd/unreleased/512' \| grep -v 'run\.go'` returns empty |
| AC2 (History) | Each moved file is a `git mv` rename, not a delete+add | `git diff --name-status main...HEAD \| grep '^R'` covers all 38 active files |
| AC3 (Stubs) | Every old active path has a moved-notice stub | `for f in docs/gamma/cdd/CDD.md docs/gamma/cdd/README.md ... ; do [ -f "$f" ] && grep -q "Moved" "$f"; done` (all 38 old paths have stub content) |
| AC4 (Frozen) | Snapshot dirs (`X.Y.Z/`), `gamma/cdd/docs/`, `gamma/rules/3.14.5/`, `.cdd/` evidence untouched | `git diff --name-status main...HEAD \| grep -E 'gamma/cdd/[0-9]+\.[0-9]+\.[0-9]+/\|gamma/cdd/docs/\|gamma/rules/3\.14\.5/'` returns empty |
| AC5 (Do-not-touch) | `gamma/conventions/`, golden files, `.github/workflows/`, 4D/4E bundle paths absent from diff | `git diff --name-only main...HEAD \| grep -E 'gamma/conventions\|alpha/protocol\|alpha/agent-runtime\|alpha/package-system\|alpha/cli\|alpha/ctb\|alpha/schemas\|beta/schema\|beta/architecture\|alpha/security\|alpha/cognitive-substrate\|\.github/workflows'` returns empty |
| AC6 (Source update) | run.go line 59 contains `"docs/development/cdd"` and line 161 contains `(default: docs/development/cdd)` | `grep -n 'docs/gamma/cdd' src/packages/cnos.cdd/commands/cdd-verify/run.go` returns empty |
| AC7 (Active links) | All active markdown links and src/path citations to 4C paths repointed to new homes | Active-surface stale-reference grep (see AC9 oracle) returns only stubs + frozen citations |
| AC8 (Relative links) | Relative links inside moved files recomputed for new depth (`docs/development/X/` vs `docs/gamma/X/`) | Inspect `git diff main...HEAD -- docs/development/` for link paths; no `../gamma/` targets remain in active files |
| AC9 (Stale-ref classification) | ALL remaining references to `docs/gamma/(cdd\|rules\|plans\|checklists)` and relative forms classified | `git grep -nE 'docs/gamma/(cdd\|rules\|plans\|checklists)\|\.\.\/gamma\/(cdd\|rules\|plans\|checklists)\|\.\.\/\.\.\/gamma\/(cdd\|rules\|plans\|checklists)' -- .` — every hit is labeled: redirect-stub / frozen-historical / intentionally-deferred / external-literal-history |
| AC10 (Goldens) | No `*.golden.yml` in diff | `git diff --name-only main...HEAD \| grep -i golden` returns empty |
| AC11 (Go build/test gate) | `go build ./...` and `go test ./...` pass with no new failures after run.go change | Run both; compare against known pre-existing failures (I4/I5/I6); any new red stops the cell |
| AC12 (Existing reds) | Any pre-existing failures are identical to known inherited reds | Check `go test ./...` failures against I4/I5/I6 baseline; new failures block |
| AC13 (No prose edits) | Moved files differ from originals only in relative link paths | `git diff main...HEAD -- docs/development/` shows only link-repair diff lines (no prose reflow, no frontmatter) |
| AC14 (Hidden/bidi) | No hidden/bidi/object-replacement chars in changed files incl. CDD receipt files | `LC_ALL=C git show --name-only main...HEAD \| xargs -I{} sh -c 'grep -Pn "[\x00-\x08\x0b\x0c\x0e-\x1f\x7f\xc2\x80-\xef\xbb\xbf]" {} && echo {}'` — no hits |
| AC15 (Receipt honesty) | Receipt records any β/δ repair; no false PASS claim | Read `.cdd/unreleased/512/beta-review.md` §R0; any accept-with-repair event is documented |

---

## Scope guardrails

**α MUST NOT touch:**
- `docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md` — golden-bound, not in 4C
- Any `*.golden.yml` file anywhere in the repo
- `.github/workflows/` directory (any file)
- Versioned snapshot dirs: `docs/gamma/cdd/<X.Y.Z>/` (~60 dirs), `docs/gamma/rules/3.14.5/`, `docs/gamma/cdd/docs/` PRAs
- `.cdd/` versioned snapshots, release evidence, historical audit records (except `.cdd/unreleased/512/` artifacts for this cell)
- Any 4D/4E bundle: `docs/alpha/protocol/`, `docs/alpha/agent-runtime/`, `docs/alpha/package-system/`, `docs/alpha/cli/`, `docs/alpha/ctb/`, `docs/alpha/schemas/`, `docs/alpha/runtime-extensions/`, `docs/beta/schema/`, `docs/beta/architecture/`, `docs/alpha/security/`, `docs/alpha/cognitive-substrate/`
- `docs/alpha/design/`, `docs/gamma/essays/`, `docs/gamma/smoke/`, `docs/gamma/kata/`, `docs/gamma/design/` — deferred, not in this cell
- Any behavior change to `run.go` beyond the literal path-default string replacements on lines 59 and 161
- No new frontmatter, no new labels, no prose reflow in moved docs
- No new imports or module changes in Go code

---

## Friction notes

### Relative-link depth changes

Moving from `docs/gamma/X/` (depth 3 from repo root; depth 2 from `docs/`) to `docs/development/X/` (same depth) means the relative link depth is IDENTICAL. Both old and new locations are at `docs/<triad>/<bundle>/file.md`. No mechanical depth shift.

However, any file that previously linked to `../gamma/rules/` from within `docs/gamma/cdd/` will now be at `docs/development/cdd/` and must link to `../rules/` (within development) rather than `../gamma/rules/` or `../../development/rules/`. α must audit all cross-bundle relative links within the 4C bundles themselves.

### Stale-reference grep (canonical form for AC9)

```bash
git grep -nE 'docs/gamma/(cdd|rules|plans|checklists)|\.\./gamma/(cdd|rules|plans|checklists)|\.\./\.\./gamma/(cdd|rules|plans|checklists)' -- .
```

The relative forms `../gamma/` and `../../gamma/` MUST be included — omitting them caused the 4B AC7 over-claim (per the issue body). Every hit must be classified:
- `redirect-stub` — the stub file itself (expected)
- `frozen-historical` — inside a frozen snapshot dir or `.cdd/` archive
- `intentionally-deferred` — path in a doc explicitly noting the move is pending
- `external-literal-history` — CHANGELOG or release notes citing the old path by name

No unclassified active-surface hit may remain.

### Go build/test gate (non-optional for design-and-build mode)

The change to `run.go` makes this a design-and-build cell. Go build/test is a hard gate — not a no-op as in 4B (which touched no Go). Run:

```bash
go build ./...
go test ./...
```

A new failure (any failure not in the I4/I5/I6 pre-existing baseline) stops the cell. Document the pre-existing baseline explicitly in the receipt.

### cdd-verify package narrower test

Also run the narrower package-level test if it exists:

```bash
go test ./src/packages/cnos.cdd/commands/cdd-verify/...
```

### docs/development/ already exists

`docs/development/README.md` is already present (Pass 1). α must not overwrite or remove it. The four new subdirs (`cdd/`, `rules/`, `plans/`, `checklists/`) are new under `docs/development/`.

### ~91 active-surface references (18 in src/)

The issue body notes ~91 active references across the 4 bundles. The 18 `src/` citations are the higher-risk surface (beyond run.go:59,161) — α must grep src/ for `docs/gamma/(cdd|rules|plans|checklists)` and repoint each active citation. Comments in Go/OCaml source files referencing old paths are active citations and must be updated.

---

## α prompt

```
You are α (implementer) executing cycle/512 of the CDD protocol.

Issue: cnos#512 — docs Pass 4C: move gamma/cdd, gamma/rules, gamma/plans, gamma/checklists → development/ homes (+ cdd-verify source path)
Mode: design-and-build
Cycle branch: cycle/512
Base SHA: 5c12ceba998b434427385913b3f1df24a83cf25e
Scaffold: .cdd/unreleased/512/gamma-scaffold.md

## 7-axis implementation contract

- Language: Go (run.go path-default update) + Markdown (document moves)
- CLI integration target: N/A — no new CLI surface; run.go change is a literal path-default string replacement only (lines 59 and 161)
- Package scoping: `src/packages/cnos.cdd/commands/cdd-verify/` for Go; `docs/development/{cdd,rules,plans,checklists}/` for docs
- Existing-binary disposition: Coexist — binary behavior unchanged; the `--bundle` default string value moves from `docs/gamma/cdd` to `docs/development/cdd`; all flag behavior and logic is preserved
- Runtime dependencies: Go toolchain only; no new runtime deps, no new imports
- JSON/wire contract: N/A — no wire contract involved in this change
- Backward-compat invariant: Stubs at every old active path (`docs/gamma/cdd/CDD.md`, etc.) redirect to new location. The run.go default is a path-move only — operators passing `--bundle docs/gamma/cdd` explicitly will still work (the stub exists); the default now points at the canonical new home

## What to do

### Step 1 — Create destination directories

```bash
mkdir -p docs/development/cdd docs/development/rules docs/development/plans docs/development/checklists
```

### Step 2 — git mv each active file

Move all 38 active files (see gamma-scaffold.md §Surfaces α will touch for full file-by-file table) with `git mv`.

Do NOT move:
- `docs/gamma/cdd/<X.Y.Z>/` snapshot dirs (~60)
- `docs/gamma/cdd/docs/` PRA records
- `docs/gamma/rules/3.14.5/`

### Step 3 — Write redirect stubs at old paths

For each moved file, write a stub at the old path:

```markdown
# Moved

This file has moved to [`docs/development/<bundle>/<filename>`](../../development/<bundle>/<filename>).
```

Adjust relative path depth as needed. Stage each stub with `git add`.

### Step 4 — Update run.go (lines 59 and 161)

File: `src/packages/cnos.cdd/commands/cdd-verify/run.go`
- Line 59: change `"docs/gamma/cdd"` to `"docs/development/cdd"`
- Line 161: change `(default: docs/gamma/cdd)` to `(default: docs/development/cdd)`

These are the ONLY changes to run.go. Do not change any logic, imports, or other content.

### Step 5 — Repoint active links and src/ citations

Grep for all active references to the four 4C bundle paths (incl. relative forms):

```bash
git grep -nE 'docs/gamma/(cdd|rules|plans|checklists)|\.\./gamma/(cdd|rules|plans|checklists)|\.\./\.\./gamma/(cdd|rules|plans|checklists)' -- .
```

For every hit that is NOT:
- a stub file at the old path
- inside a frozen snapshot dir or `.cdd/` archive
- in a `.cdd/unreleased/` historical record

Update the reference to point at the new `docs/development/` path. Include `src/` Go/OCaml comment citations.

### Step 6 — Recompute relative links inside moved files

Within the moved files themselves, any relative link that pointed to `../gamma/X/` or similar must be recomputed for the new location (`docs/development/<bundle>/` rather than `docs/gamma/<bundle>/`). Both are at the same structural depth; cross-bundle links within the gamma triad will change from `../rules/` (if already triad-relative) or from `../../gamma/rules/` (if repo-relative-ish) to `../rules/` (within development/).

### Step 7 — Run required checks

```bash
git status --short
git diff --name-status main...HEAD
git diff --check
# stale-reference grep (include relative forms — omitting them caused 4B AC7 over-claim):
git grep -nE 'docs/gamma/(cdd|rules|plans|checklists)|\.\./gamma/(cdd|rules|plans|checklists)|\.\./\.\./gamma/(cdd|rules|plans|checklists)' -- .
# classify every remaining hit per AC9
# Go build/test (hard gate — new failure stops the cell):
go build ./...
go test ./...
go test ./src/packages/cnos.cdd/commands/cdd-verify/...
# Golden protection — expect NO output:
git diff --name-only main...HEAD | grep -i golden
```

### Step 8 — Write receipt and set status:review

Write `.cdd/unreleased/512/alpha-closeout.md` with:
- name-status proof (paste output of `git diff --name-status main...HEAD`)
- active stale-reference grep (paste classified output)
- do-not-touch proof (golden no-change grep output)
- Go build/test result (pass or documented pre-existing failure baseline vs I4/I5/I6)
- honest record of any β/δ repair performed

Commit all changes with:
```
α R0 implementation: cycle/512 (cnos#512, docs Pass 4C)
```

Then push and set issue status to `status:review`.

## AC oracle reference

See gamma-scaffold.md §Per-AC oracle for the full 15-AC table.

## δ review policy (operator-set)

Accept-with-repair only if: link-repair only, ≤10 lines, target unambiguous from move map, no Go/source/behavior change beyond literal path move, no golden, no frozen record, no do-not-touch path, receipt records the miss.

Stop for operator if: Go build/test goes green→red, a golden needs re-render, `gamma/conventions` is implicated, a frozen record would need editing, the source-path update changes behavior beyond a literal path move, or stale references exceed trivial repair scope.
```

---

## β prompt

```
You are β (reviewer) executing review R0 for cycle/512 of the CDD protocol.

Issue: cnos#512 — docs Pass 4C
Mode: design-and-build
Cycle branch: cycle/512
Scaffold (AC oracle): .cdd/unreleased/512/gamma-scaffold.md §Per-AC oracle
α receipt: .cdd/unreleased/512/alpha-closeout.md

## Your task

Walk each of the 15 ACs from gamma-scaffold.md §Per-AC oracle independently. For each AC:
1. Run the mechanical check specified in the oracle column
2. Record: PASS / FAIL / REPAIR-NEEDED with evidence
3. If REPAIR-NEEDED and within accept-with-repair limits (below), apply the repair and record it

Write your verdict to `.cdd/unreleased/512/beta-review.md §R0` using this structure:

```markdown
# β review R0 — cnos#512 — docs Pass 4C

**cycle-branch:** cycle/512
**reviewer:** β
**round:** R0

## AC verdicts

| AC | Verdict | Evidence summary |
|---|---|---|
| AC1 | PASS/FAIL | ... |
| AC2 | PASS/FAIL | ... |
...
| AC15 | PASS/FAIL | ... |

## Required-checks output

### git diff --name-status main...HEAD
<paste output>

### stale-reference grep (incl. relative forms)
<paste classified output>

### golden protection
<paste: git diff --name-only main...HEAD | grep -i golden>

### go build/test
<paste output or "PASS (no new failures vs I4/I5/I6 baseline)">

## Repairs performed (if any)
<list or "none">

## Overall verdict
PASS / FAIL / ACCEPT-WITH-REPAIR

## Operator stop? (yes/no)
<yes if any stop condition triggered; no if within accept-with-repair>
```

## Required checks to run

```bash
git status --short
git diff --name-status main...HEAD
git diff --check
# stale-reference grep — include relative forms (omitting them caused 4B AC7 over-claim):
git grep -nE 'docs/gamma/(cdd|rules|plans|checklists)|\.\./gamma/(cdd|rules|plans|checklists)|\.\./\.\./gamma/(cdd|rules|plans|checklists)' -- .
# classify every remaining hit per AC9
# Go build/test (hard gate):
go build ./...
go test ./...
go test ./src/packages/cnos.cdd/commands/cdd-verify/...
# Golden protection — expect NO output:
git diff --name-only main...HEAD | grep -i golden
```

## Accept-with-repair limits (operator-set)

Accept-with-repair only if ALL of:
- Link-repair only (no Go/source/behavior change)
- ≤10 lines total across all repairs
- Target unambiguous from move map (gamma-scaffold.md §Surfaces α will touch)
- No golden, no frozen record, no do-not-touch path implicated
- Receipt records the miss

Stop for operator (do NOT self-repair) if:
- Go build/test goes green→red
- A golden needs re-render
- `gamma/conventions` is implicated
- A frozen record would need editing
- The source-path update changes behavior beyond a literal path move
- Stale references exceed trivial repair scope (>10 lines)
```
