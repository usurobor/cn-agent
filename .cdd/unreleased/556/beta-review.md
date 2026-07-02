---
artifact: beta-review
cycle: 556
issue: https://github.com/usurobor/cnos/issues/556
branch: cycle/556
round: R0
author: beta
---

# β review — cnos#556

## §R0

**Base for this review:** `git fetch origin main` → `origin/main` = `4fe8e4333b36372f595201841fb76cc0c31acff4`
(unchanged since γ's scaffold-time snapshot; γ's §0 base-SHA note re-verified: two
board-map bot commits ahead of the wake-invoked-δ pinned SHA `2ae24f27...`, both
touching only `docs/development/board/{board-data.json,index.html}` — confirmed
benign by re-running `git diff --stat 2ae24f2..origin/main` myself).

**Cycle-branch head reviewed:** `origin/cycle/556` = `ff9635c307ac33141870e40cf9fb9858a5a96326`
(`git rev-parse origin/cycle/556`, matches local `cycle/556` after
`git fetch origin && git reset --hard origin/cycle/556`).

All findings below are from my own commands, not pasted from α's artifact. Where a
command matched α's claimed output, I note "confirms α" — I did not skip re-running
anything on the independent-verification list.

### AC1 — `cnos.issues` package exists — **CONFIRMED**

- `src/packages/cnos.issues/cn.package.json` exists, parses as valid JSON
  (`python3 -c "import json; json.load(open(...))"`), `schema: "cn.package.v1"`,
  `"name": "cnos.issues"`, `"version": "0.1.0"`, `"kind": "package"`,
  `"engines": {"cnos": ">=3.82.0"}`.
- Built `cn` myself from `src/go` at cycle-branch HEAD (`go build -o /tmp/cn_556
  ./cmd/cn`); ran `/tmp/cn_556 build --check` myself:
  ```
  ✓ cnos.issues: valid
  ✓ All packages valid.
  ```
- Cross-checked the version/engines convention against `cnos.cds/cn.package.json`
  and `cnos.handoff/cn.package.json` (both `0.1.0`/independent-versioned packages) —
  `cnos.issues` follows the same convention exactly, and `>=3.82.0` matches
  `git tag --sort=-v:refname | grep -v '^v' | head -1` = `3.82.0`.

### AC2 — Package owns issue-board cognition — **CONFIRMED**

Read `src/packages/cnos.issues/SKILL.md` in full myself (not just grepped). It
states, under "Core Principle": `cnos.issues` is the package home for issue-board
cognition: it owns issue taxonomy, board mapping, and Issue Pivot generation as a
domain, distinct from the CLI dispatch that invokes it" — then names three owned
sub-domains each with its own skill file, a full "Command-dispatch disposition"
section (see AC10), and a "Non-goals" list. Not an empty wrapper.

### AC3 — Public command remains `cn issues map` — **CONFIRMED**

Built `cn` myself from the cycle branch (same binary as AC1) and ran, myself:
```
$ /tmp/cn_556 issues map --fixture src/packages/cnos.issues/commands/issues-map/testdata/issues.json --out /tmp/board-out-beta
cn issues map: wrote 5 open issues to /tmp/board-out-beta (index.html, board-data.json, README.md)
exit=0
```
`index.html`, `board-data.json`, `README.md` all present and non-empty in the output
directory. Works from a package-relative fixture with no package-internal script or
path needed — same public command shape as before the move.

### AC4 — CLI dispatch remains thin — **CONFIRMED**

`git diff main...cycle/556 -- src/go/internal/cli/cmd_issues_map.go` (ran myself):
only the import path (`internal/issuesmap` → `packages/cnos.issues/commands/issues-map`)
and an expanded doc-comment changed. `Run()` is still a one-line delegation:
```go
func (c *IssuesMapCmd) Run(ctx context.Context, inv Invocation) error {
	return issuesmap.Run(ctx, inv.Args, inv.Stdin, inv.Stdout, inv.Stderr)
}
```
No taxonomy/rendering/label logic entered this file.

### AC5 — GitHub Action calls the public command — **CONFIRMED**

Read `.github/workflows/board-map.yml` in full myself. `git diff main...cycle/556
-- .github/workflows/board-map.yml` is empty (file untouched). The workflow builds
`cn` from `src/go` and runs only `./cn issues map --repo "$GITHUB_REPOSITORY" --out
docs/development/board` — no package-internal path, no `scripts/board`, no Node.

### AC6 — Package installation path is explicit — **CONFIRMED**

`grep -c "issues-map" src/packages/cnos.issues/cn.package.json` → 0 (no `commands`
key at all in the manifest, confirmed via direct JSON parse: `'commands' in d` →
`False`). No `cn deps restore` in `board-map.yml`, no `.cn/vendor/` mutation
anywhere in the diff. Receipt (`self-coherence.md` AC6, and `SKILL.md`'s
"Command-dispatch disposition" section) states the shim disposition honestly and
explicitly — matches AC10 below.

### AC7 — Current board output behavior is preserved — **CONFIRMED (independently reproduced, not just sanity-checked)**

I did not trust α's paste — I rebuilt both sides myself in a separate worktree:

1. `git worktree add /tmp/main-worktree-556 origin/main` (`4fe8e433`, the true
   pre-move base — **not** local `main`, which was stale at `2ae24f27`; I caught
   this myself when a first diff attempt showed spurious `board-data.json`/
   `index.html` deltas explained entirely by local `main` lagging `origin/main` by
   the same two bot commits γ's §0 note documents — re-ran against `origin/main`
   and the spurious deltas disappeared, leaving only the intended `README.md`
   line).
2. Built `/tmp/cn_premove` from that worktree's `src/go`.
3. Ran both `/tmp/cn_premove` (against the worktree's own pre-move
   `testdata/issues.json`) and `/tmp/cn_556` (against the relocated
   `testdata/issues.json`) with the same fixture, to separate output directories.
4. `diff board-data.json` → **empty (byte-identical)**.
5. `diff index.html` → **empty (byte-identical)**.
6. `diff README.md` → differs only in the expected path-reference line (old
   `src/go/internal/issuesmap/` vs new `src/packages/cnos.issues/commands/issues-map/`).

This exactly matches α's claimed result. No `docs/board/` (pre-#545 path)
resurrection anywhere in the diff.

### AC8 — No Node production generator — **CONFIRMED**

`find . -iname "*.mjs" -path "*board*"` → no hits. `grep -n "node\|npm\|npx"
.github/workflows/board-map.yml` → no hits. `git diff --stat main...cycle/556`
introduces no `package.json`/`node_modules`.

### AC9 — CI remains green at the FINAL head commit — **CONFIRMED**

`gh run list --branch cycle/556 --json status,conclusion,name,headSha` (ran
myself) shows the run at `headSha: ff9635c307ac33141870e40cf9fb9858a5a96326` —
which I confirmed via `git rev-parse origin/cycle/556` is the actual current
branch tip, one commit past the `987db05a` run α cited in self-coherence.md's
Review-readiness section. `gh run view 28628047010 --json jobs` on that exact
head SHA:

```
success  CDD artifact ledger validation (I6)
success  Package/source drift (I1)
success  SKILL.md frontmatter validation (I5)
success  Repo link validation (I4)
success  Dispatch repair-preflight guard (cnos#516)
success  Dispatch closeout-integrity guard (cnos#524)
success  Protocol contract schema sync (I2)
success  Go build & test
success  Binary verification
success  Package verification
```

All 10 required jobs green, all at the actual final HEAD — not a stale earlier
commit. (Minor observation, not a finding: α's self-coherence.md cited the green
run at `987db05a` and did not re-cite the subsequent `ff9635c3` "Review-readiness"
commit's own CI result — harmless since I independently verified the true final
head is also green, but worth noting for the closeout: the review-readiness commit
itself always lands one commit after the last-cited CI-green SHA in this protocol
shape, so β re-verifying at the literal branch tip, as required, is not optional
here.)

### AC10 — Receipt records the #216 relationship — **CONFIRMED, Option B, unambiguous**

Read `src/packages/cnos.issues/SKILL.md` directly (not via α's paste). Its
"Command-dispatch disposition" section states, verbatim:

> `cn issues map` remains a **temporary built-in shim until #216 lands**. Domain
> logic is isolated under this package's boundary
> (`src/packages/cnos.issues/commands/issues-map/`, its own Go module, mirroring
> the cnos#392 `cdd-verify` precedent), but dispatch stays through the compiled-in
> kernel registration, not through the package-command exec-dispatch mechanism
> (`PACKAGE-SYSTEM.md` §7).

This is unambiguously Option B. No claim of Option A ("package-provided through
documented command-discovery") appears anywhere in `SKILL.md` or
`self-coherence.md`. Cross-checked against the facts that would falsify an
Option-A claim: `main.go:47`'s `reg.Register(&cli.IssuesMapCmd{})` is unchanged
(confirmed: `git diff main...cycle/556 -- src/go/cmd/cn/main.go` is empty), and
`cn.package.json` has no `commands` object at all. Both are consistent with
Option B and neither is contradicted. **Hard gate passes — no REQUEST CHANGES
required on this axis.**

## Additional independent verification (gamma-scaffold.md §9 "beyond the standard AC walk")

- **`cn.package.json` "commands" object does NOT contain `issues-map`** —
  confirmed (`'commands' in d` → `False`; the manifest has no `commands` key at
  all, which is even stricter than "key present but not listing issues-map").
- **`src/go/cmd/cn/main.go`'s `reg.Register(&cli.IssuesMapCmd{})` unchanged** —
  confirmed, `git diff main...cycle/556 -- src/go/cmd/cn/main.go` is empty.
- **Relocated module's `go.mod`/`go.work`/`src/go/go.mod` mirror the `cdd-verify`
  shape** — diffed both shapes directly, not skimmed:
  - `cdd-verify/go.mod`: `module github.com/usurobor/cnos/packages/cnos.cdd/commands/cdd-verify` /
    `go 1.24`. `issues-map/go.mod`: `module github.com/usurobor/cnos/packages/cnos.issues/commands/issues-map` /
    `go 1.24`. Identical shape, only the module path differs as expected.
  - `go.work`'s `use (...)` list: `./src/go`, `./src/packages/cnos.cdd/commands/cdd-verify`,
    `./src/packages/cnos.issues/commands/issues-map` — new entry added, existing
    entries untouched.
  - `src/go/go.mod`: the new `require`/`replace` pair for `issues-map` is
    structurally identical to the existing `cdd-verify` pair (same comment style,
    same `v0.0.0-00010101000000-000000000000` placeholder version, same
    `replace ... => ../packages/...` shape). Confirmed by direct `grep -A2 -B2`
    on both blocks side by side.
- **Repo-wide grep for `internal/issuesmap`** — ran myself:
  `grep -rn "internal/issuesmap" .` (excluding `.git/`) → only prose references
  inside `.cdd/unreleased/556/{gamma-scaffold,self-coherence}.md` (historical/
  narrative mentions of the old path, correctly describing what moved and from
  where) — **zero references in code, docs outside the CDD artifacts, or
  comments**. `docs/development/board/README.md:42` was correctly updated (verified
  by reading the file directly).
- **`//go:embed templates/board.template.html` resolves** — confirmed two ways:
  (a) `go build ./...` from the new module root succeeds (would fail immediately
  if the embed target were missing), and (b) directly listed
  `src/packages/cnos.issues/commands/issues-map/templates/board.template.html` on
  disk — present, moved with the `.go` files as a single directory, not split or
  duplicated. `testdata/issues.json` also confirmed co-located.
- **`cd src/go && go build ./... && go vet ./... && go test ./...`** — ran myself:
  `go build ./...` clean (no output = success); `go vet ./...` clean; `go test ./...`
  — all `src/go` packages `ok` (14 packages, matching α's count).
- **`go test ./...` from `src/packages/cnos.issues/commands/issues-map/`** — ran
  myself: `go build ./... && go vet ./... && go test ./... -v` — 4 test functions,
  8 pass lines (`TestToRecord_LabelParsingAndEffort` + 4 subtests,
  `TestEffortWeights`, `TestRun_Fixture`, `TestRun_Stdin`), 0 failures. Matches
  α's claimed count exactly.

## Findings

### Finding 1 — γ-scaffold-level tension with the operator's explicit "Go implementation rule" (process/coherence, non-blocking)

**Severity: C (process-coherence, informational — not an implementation-contract
violation by α, and not blocking this review).**

The operator's clarifying comment on the issue (read directly via `gh issue view
556 --json comments`, not via any paraphrase) contains an explicit "Go
implementation rule":

> Do not force Go implementation code into `src/packages/`. The active Go
> implementation may remain under `src/go/internal/issuesmap/` during the shim
> phase, but it must be treated as the implementation of the `cnos.issues` domain.

γ's scaffold §6 nonetheless pins, as a **binding** guardrail, relocating the Go
domain source into `src/packages/cnos.issues/commands/issues-map/` — the opposite
of what this operator sentence permits/prefers — justified by the cnos#392
`cdd-verify` precedent and the "fake unused command directory" concern. γ's
scaffold labels this section "restated from the operator's clarifying comment,"
but on direct reading of the actual comment, this specific axis is γ's own
architectural judgment call overriding one explicit operator sentence, not a
restatement of it. No `gamma-clarification.md` was filed on the branch flagging
this specific divergence for δ/operator visibility before dispatch.

**Why I am not blocking on this:** β's Rule 7 verification is against the pinned
Implementation Contract in γ's dispatch prompt (which α followed exactly and
which δ had the opportunity to enrich/block per its own §2 authority), not
against re-litigating γ's scaffold decisions against the raw issue text. The
operator's own "Revised success condition" list (in the same clarifying comment)
does not name Go-source location as a success criterion — it only requires the
receipt to honestly state shim-vs-package-provided disposition, which α's receipt
does (AC10, confirmed above). The relocation is also defensible on its own
technical merits (mirrors an already-shipped, CI-green precedent; keeps the
distinction between "Go-source co-location" and "declared package command"
concrete and auditable, which arguably serves the operator's underlying intent
better than a bare pinned instruction anticipated). I am naming this as a finding
per Rule 3 ("stale references and authority conflicts...reviewable incoherence")
rather than silently passing over it, but I am not treating it as grounds for
REQUEST CHANGES against α's diff, since α's diff conforms 100% to what was
actually pinned, and the receipt is honest about the choice made.

**What γ/δ should know for the PRA:** this is worth a one-line note in the cycle
assessment — a future cycle should either (a) have γ file an explicit
`gamma-clarification.md` when a scaffold's binding guardrail knowingly overrides
a specific operator sentence rather than merely narrowing permissive language, or
(b) have δ's dispatch-prompt review catch and flag the divergence at enrichment
time. Neither happened here; the outcome was fine, but the protocol path that
should have surfaced it did not fire.

No other findings. Every AC oracle, every scope guardrail, every friction note in
gamma-scaffold.md §7, and every "beyond the standard AC walk" item in §9 was
independently re-verified using my own commands and matched α's claims (or, in
AC7's case, initially diverged due to my own stale local `main` ref before I
caught and corrected the base, at which point it matched exactly).

## Verdict

verdict: converge

All 10 ACs independently confirmed. All scope guardrails held (no generic
package-command-discovery implementation, no workflow edit, no Node generator, no
declared `issues-map` command, no #216-solving, no board-visualization or
taxonomy/effort semantic change). CI green at the true final head commit
(`ff9635c3`). Implementation-contract axes (Rule 7) all conform to what γ pinned
and δ had the opportunity to enrich. `go build`/`go vet`/`go test` all pass in both
`src/go` and the new `src/packages/cnos.issues/commands/issues-map` module,
verified independently. Finding 1 is a process-coherence observation for γ's PRA,
not a merge blocker — α's diff and receipt are both faithful to the contract that
was actually pinned, and the receipt is honest.

This artifact records the R0 review verdict only; merge/close-out is a separate
dispatched step per this session's own task scope.
