# γ scaffold — cycle #398 (Phase 4b of #366)

**Issue:** cnos#398 — harness substrate (observability + worktree + identity); absorbs cnos#371 #373 #384.
**Branch:** `cycle/398` from `origin/main`.
**Mode:** design-and-build; γ+α+β-collapsed-on-δ (single-session via Agent tool, per `operator/SKILL.md` §5.2).
**Wave:** none (single issue under #366 phase chain).

## Selection

CDD §3 decisive clause: parent #366 declares Phase 4b ready (4a/4c are siblings); this is the next-MCA committed by #393 close-out. Three named inputs (#371 #373 #384) are explicit Phase 4 inputs in #366 body and gate on Phase 4b shipping the harness substrate.

## Mode

design-and-build. The design half: choose harness skill path, partition operator/SKILL.md sections, define cross-references. The build half: author `harness/SKILL.md`, edit `operator/SKILL.md`, update γ + β cross-refs.

## Surfaces α expects to touch

- **NEW:** `src/packages/cnos.cdd/skills/cdd/harness/SKILL.md` (canonical path per §366 Phase 4 hint and existing sibling skill structure under `cdd/`)
- `src/packages/cnos.cdd/skills/cdd/operator/SKILL.md` (extract harness mechanics; keep operator-as-coordinator + dispatch authority WHY)
- `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` (cross-ref updates: γ dispatch mechanics → `harness/SKILL.md`)
- `src/packages/cnos.cdd/skills/cdd/beta/SKILL.md` Pre-merge gate Row 1 (cross-ref to harness identity discipline)
- `.cdd/unreleased/398/*` (γ scaffold, α design notes, α self-coherence, β review, α/β/γ close-outs, cdd-iteration)
- `.cdd/iterations/INDEX.md` (one row for #398)

## Implementation contract (pinned by δ; α MUST NOT improvise)

Per the dispatch prompt for cycle #398 (mirrors the issue body):

| Axis | Pinned value |
|---|---|
| Language | Markdown (skill files). No Go touches — `cn dispatch` is codified, not reimplemented. |
| CLI integration target | `cn dispatch` (existing Go subcommand under `src/go/internal/dispatch/`); contract codified, code untouched. |
| Package scoping | `src/packages/cnos.cdd/skills/cdd/harness/SKILL.md` (new); edits to existing `operator/SKILL.md`, `gamma/SKILL.md` cross-refs, `beta/SKILL.md` Row 1 cross-ref. |
| Existing-binary disposition | `cn dispatch` Go code preserved as-is. |
| Runtime dependencies | None new. |
| JSON/wire contract | `--output-format stream-json` JSONL contract per cnos#371 — codified as a δ-membrane primitive (the wire format itself is documented, not changed). |
| Backward compat | `cn dispatch` invocations continue to work; γ/β/δ skills' load orders unchanged structurally — only their cross-refs shift to point at `harness/SKILL.md` for HOW (mechanics) while continuing to point at `operator/SKILL.md` for the δ role's WHY (gate authority, override, route discipline). |

## AC oracle approach

Each AC is mechanical or ripgrep-checkable:

- **AC1:** `test -f src/packages/cnos.cdd/skills/cdd/harness/SKILL.md && head -1` shows valid frontmatter.
- **AC2:** `rg "output-format stream-json|stream-json" src/packages/cnos.cdd/skills/cdd/harness/SKILL.md` ≥ 1 hit; rule names JSONL line schema or references an existing schema source.
- **AC3:** `rg "extensions.worktreeConfig|--worktree" src/packages/cnos.cdd/skills/cdd/harness/SKILL.md` ≥ 1 hit; cross-references β-skill Row 1 (cycle #301 O8 anchor).
- **AC4:** `rg "parallel|per-cycle worktree|pre-created" src/packages/cnos.cdd/skills/cdd/harness/SKILL.md` ≥ 1 hit naming the per-cycle pre-creation requirement.
- **AC5:** `rg "claude -p|cn dispatch|polling loop|git fetch.*polling|worktree|stream-json" src/packages/cnos.cdd/skills/cdd/operator/SKILL.md` returns 0 hits OR only cross-refs to harness/SKILL.md.
- **AC6:** γ-skill + β-skill cross-refs to `harness/SKILL.md` exist (≥ 1 hit each via ripgrep).
- **AC7:** Close-out comments on #371 #373 #384 name #398 as the absorbing cycle.

## Empirical anchors

- cnos#371: parallel #369/#370 dispatch in default text mode (2026-05-17) — the falsification anchor that the observability rule already cites in `operator/SKILL.md` §1.2 and `gamma/SKILL.md` §2.5.
- cnos#373: cycle #370 β R1 §2.1 row 1 + α F4 + recurring #301 O8 — identity-leak through shared `.git/config` when `extensions.worktreeConfig=true`.
- cnos#384: cph cdr-refactor wave 2026-05-18–05-20 F1 — α-28 worktree-add race under shared `/root/cph` working tree.

## Diff scope expectation

- 1 new file: `src/packages/cnos.cdd/skills/cdd/harness/SKILL.md` (~400 lines).
- 1 edit: `src/packages/cnos.cdd/skills/cdd/operator/SKILL.md` (extract harness mechanics; keep §3 Gate authority WHY, §3a inward-membrane, §4 Override, §6 What δ does NOT do). Net reduction expected.
- 2 small cross-ref edits: `gamma/SKILL.md` (§1.2 dispatch flow → harness/SKILL.md), `beta/SKILL.md` (Pre-merge gate Row 1 → cross-ref harness/SKILL.md identity discipline).
- 7 cycle evidence files under `.cdd/unreleased/398/`.
- 1 row added to `.cdd/iterations/INDEX.md`.

## Non-goals (per issue body)

- Do NOT reimplement `cn dispatch` Go code.
- Do NOT move release effector content (Phase 4c — cycle/399).
- Do NOT move δ-role boundary content (Phase 4a — cycle/397).
- Do NOT change observability format substantively — codify what exists / what #371 specified.

## Coordination with siblings

Phase 4a (cycle/397) and Phase 4c (cycle/399) may touch the same `operator/SKILL.md`. Sections touched here:

- §1 (Route mechanics — `claude -p` invocation patterns)
- §1.2 (sequential α/β dispatch shell)
- §2.2 (issue polling loop)
- §5.1 (canonical multi-session dispatch shell)
- §5.2.1 (parent-session quiescence) — keep, this is δ/γ-as-coordinator discipline, not mechanics
- §8 (timeout recovery) — moves (this is harness session-handling)
- §Git identity for role actors — moves (this is harness identity setup)

Sections kept in `operator/SKILL.md`:

- §Core Principle (δ as boundary; WHY)
- §Algorithm step list (WHY only; HOW points at harness)
- §3 Gate (external actions δ holds — WHY of gate authority)
- §3a δ as inward membrane (cnos#393)
- §4 Override
- §5 Dispatch configurations (the §5.2 δ=γ collapse is δ doctrine, not harness mechanics)
- §6 What the operator does NOT do
- §7 Cycle lifecycle from operator's view
- §10 Wave Coordination (δ-as-wave-planner; not harness)

If cycle/397 (Phase 4a) intends to move §Algorithm into `delta/SKILL.md`, this cycle's edits to `operator/SKILL.md` integrate cleanly because Phase 4a moves the file wholesale (this cycle leaves what stays in operator a coherent residual). Merge conflict is resolvable by integration: each cycle touches different surfaces.

## Schedule

1. α design notes (`design-notes.md`) — partition decisions
2. α build: `harness/SKILL.md` + `operator/SKILL.md` edits + `gamma/SKILL.md` + `beta/SKILL.md` cross-refs
3. α self-coherence (with §Implementation contract conformance per β Rule 7)
4. β-collapsed review (AC1–AC7)
5. β-collapsed merge + β close-out
6. α close-out, γ close-out, cdd-iteration
7. Merge cycle/398 → main; comment + close #371 #373 #384
