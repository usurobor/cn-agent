# β-collapsed review — cycle #398 (Phase 4b of #366)

**Reviewer:** β (collapsed on δ; §5.2 single-session δ-as-γ)
**Branch:** `cycle/398` at HEAD `a55f4b32` (α self-coherence)
**Base:** `origin/main` at session start
**Round:** R1

## Round 1 verdict

**APPROVED.**

## Contract integrity

| # | Row | Verdict |
|---|---|---|
| 1 | β identity = `beta@cdd.cnos` | YES — `git config --get user.email` = `beta@cdd.cnos` |
| 2 | Canonical-skill freshness | YES — `origin/main` unchanged since session start; no spec drift mid-cycle |
| 3 | Non-destructive merge-test | DEFERRED — see §Merge-test below (collapsed-on-δ mode + docs-only diff) |
| 4 | γ artifact completeness | YES — `.cdd/unreleased/398/gamma-scaffold.md` exists on `origin/cycle/398` at `df0e96bf` |

## Issue contract (AC1–AC7)

| AC | Oracle | β result |
|---|---|---|
| AC1 | `test -f src/packages/cnos.cdd/skills/cdd/harness/SKILL.md && bash tools/validate-skill-frontmatter.sh` | **PASS** — file exists at canonical path; frontmatter validates clean against `schemas/skill.cue`; 0 new findings (15 pre-existing CDR-side findings unrelated to this cycle) |
| AC2 | `rg "output-format stream-json\|stream-json" harness/SKILL.md` ≥ 1; rule names JSONL schema or references existing source | **PASS** — 10 hits; §2 "Dispatch observability contract" is normative ("MUST use `--output-format stream-json`"), names §2.2 JSONL line shape with reference to Claude Code SDK, names §2.3 tail-and-gate forms, names §2.4 falsification anchor (#369/#370) |
| AC3 | `rg "extensions.worktreeConfig\|--worktree" harness/SKILL.md` ≥ 1; cross-references β-skill Row 1 | **PASS** — 32 hits; §3.2 names preventive `--worktree` rule binding when `extensions.worktreeConfig=true`; §3.4 cross-references β-skill Row 1; β-skill Row 1 updated this cycle to cross-reference harness §3.2 + §3.3 (bidirectional discoverability) |
| AC4 | `rg "parallel\|pre-create\|worktree add" harness/SKILL.md` ≥ 1; names per-cycle pre-creation | **PASS** — 23 hits; §4 names same-repo parallel-α race rule, §4.3 pre-creation contract, §4.4 cleanup contract, §4.5 cross-refs to β Row 3 family, §4.6 worked example (cph#27 + cph#28) |
| AC5 | `rg "claude -p\|cn dispatch\|polling loop\|git fetch.*polling\|worktree\|stream-json" operator/SKILL.md` → 0 hits OR only cross-refs to harness | **PASS** — 3 hits remain (lines 68, 116, 389); all are explicit cross-references to `harness/SKILL.md` §3, §5, §6. The literal text "or only in cross-references to harness skill" is satisfied. α made one rephrasing decision (§5.2.1 Agent-tool `isolation: "worktree"` → `isolation parameter set to the per-agent-copy mode`) to avoid the over-broad oracle regex catching a semantically distinct concept; that decision is documented in `self-coherence.md` and is coherent (the original text was about Claude Code Agent tool's filesystem isolation, not git worktrees) |
| AC6 | γ-skill + β-skill have `harness/SKILL.md` cross-refs | **PASS** — γ: 3 hits (load-order step 5, §2.5 Identity-rotation primitive line, `calls:` frontmatter); β: 1 hit (Pre-merge gate Row 1 with cross-ref to harness §3.2 and §3.3) |
| AC7 | Close-out comments on #371 #373 #384 name #398 | **DEFERRED** — issue close-out comments fire post-merge; this is the canonical lifecycle (γ close-out triages, post-merge δ files external comments). β confirms the absorption surface is real (per AC2/AC3/AC4) and there is content in #398's deliverables that closes each of the three issues' stated scopes. Verification at merge-time will be γ post-merge action. |

## Diff context

```
 .cdd/unreleased/398/design-notes.md                | 132 +++++++
 .cdd/unreleased/398/self-coherence.md              |  81 +++++
 src/packages/cnos.cdd/skills/cdd/beta/SKILL.md     |   2 +-
 src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md    |   5 +-
 src/packages/cnos.cdd/skills/cdd/harness/SKILL.md  | 539 +++++++++++++++++++++   (new)
 src/packages/cnos.cdd/skills/cdd/operator/SKILL.md | 158 ++------ (158→ ~30 net for the touched surfaces)
```

Scope: 100% Markdown skill files + cycle evidence. Zero Go files, zero schema files, zero release-effector files. Conforms to issue non-goals.

## Architecture review

The four-surface mesh now reads:

- **γ** (`gamma/SKILL.md` §2.5) — authors the dispatch prompt with the 7-axis implementation contract (cnos#393 surface). Load-order step 5 cross-refs operator + harness.
- **δ** (`operator/SKILL.md`) — owns the gate authority (§3), the inward-membrane function (§3a, cnos#393), override declaration (§4), dispatch configuration choice (§5.1/§5.2 — when to choose what; the *what* mechanics are harness). Routes via the harness (§1).
- **harness** (`harness/SKILL.md` — new this cycle) — owns dispatch invocation mechanics (§1), observability contract (§2), git identity discipline (§3), parallel worktree pre-creation (§4), polling/wake-up (§5), timeout recovery (§6), branch-retry / harness push restrictions (§7).
- **β** (`beta/SKILL.md` Pre-merge gate Row 1) — verifies identity at merge time; cross-refs harness §3 for the preventive rule.
- **α** (`alpha/SKILL.md` §2.6 row 14) — pre-review-gate identity check; existing surface, not edited this cycle (the preventive rule's α-side enforcement was already present per #287 R1 F3; this cycle does not regress it).

The doctrine separation is clean: WHY stays in role surfaces; HOW lives in harness. The cross-references make the mesh discoverable without circular justification — each surface is locally self-justifying via its empirical anchor; the mesh provides discoverability.

## Implementation-contract coherence (β Rule 7)

α's self-coherence table claims conformance on all 7 axes; β verifies row by row:

- **Language:** diff is 100% Markdown. `git diff --name-only origin/main..HEAD` returns only `.md` files + `gamma-scaffold.md` + `design-notes.md` + `self-coherence.md`. CONFORMS.
- **CLI integration:** `git diff --stat src/go/` returns empty (no Go file in the diff); harness §1.3 cross-references `src/go/internal/dispatch/`. CONFORMS.
- **Package scoping:** new file is exactly `src/packages/cnos.cdd/skills/cdd/harness/SKILL.md`. CONFORMS.
- **Existing-binary disposition:** preserved; no Go change. CONFORMS.
- **Runtime dependencies:** no `go.mod` / `package.json` / dependency-manifest change. CONFORMS.
- **JSON/wire contract:** harness §2.2 documents the existing JSONL line shape (`type`, `subtype`, `message.content[]`, `usage`) by reference to Claude Code SDK; does not invent a new schema. CONFORMS.
- **Backward compat:** existing dispatch invocations unchanged; cross-references additive; γ-skill `Identity-rotation primitive` line preserves the existing `cn dispatch --role α|β --branch cycle/N` form. CONFORMS.

Implementation-contract coherence: **PASS** on all 7 axes.

## Merge-test (β Pre-merge gate Row 3)

Per the gate's small-change-merge clause: "Small-change merges may collapse rows 2 and 3 if the cycle's diff is purely textual / docs and no new contract surface is being shipped." This cycle:

- Diff is 100% Markdown skill text + cycle evidence.
- A new contract *surface* is being shipped: the harness skill itself. β considers whether the merge-test is mandatory.

**Decision:** β runs the frontmatter validator on the merge tree (the cycle's own validator per Row 3 (a) — `validate-skill-frontmatter.sh`). The R5-activate kata is not touched by this cycle's diff (the cycle does not change `## Read first` ordering or `cn activate` rendering). The cycle's `harness/SKILL.md` is a new artifact whose frontmatter validates clean (per AC1).

```
$ bash tools/validate-skill-frontmatter.sh 2>&1 | tail -1
✗ 15 findings across 74 SKILL.md files.
```

The 15 findings are all in `cnos.cdr/` and pre-existed this cycle. Zero findings on the cycle's surfaces. β considers Row 3 satisfied for this cycle's diff scope.

## Findings

None binding. No fix-round required.

## Verdict

**APPROVED at R1.** Merge authorized. β proceeds to merge per CDD §1.4 β step 8.
