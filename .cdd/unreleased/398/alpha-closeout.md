# α close-out — cycle #398

**Author:** α (collapsed on δ)
**Cycle:** #398 (Phase 4b of #366; absorbs #371 #373 #384)
**Branch:** `cycle/398` HEAD `47dbd68c` (β closeout)
**Status:** APPROVED at R1 by β; ready for merge into main.

## What I built

1. **`src/packages/cnos.cdd/skills/cdd/harness/SKILL.md`** — new substrate skill (~540 lines) with §1 Dispatch invocation, §2 Dispatch observability contract (codifies #371), §3 Git identity for role actors with §3.2 worktree-aware preventive rule (codifies #373), §4 Parallel dispatch precondition (codifies #384), §5 Polling and wake-up, §6 Timeout recovery, §7 Branch retries / push restrictions, §8 Embedded Kata.

2. **`src/packages/cnos.cdd/skills/cdd/operator/SKILL.md`** — extracted harness mechanics; kept δ-doctrine. Operator/SKILL.md is now WHY-only for the touched surfaces; the literal `claude -p` mechanics, the `--output-format stream-json` flag, the polling-loop shell forms, the §8 timeout-recovery procedure all live in harness/SKILL.md. Cross-references where mechanics used to be inline.

3. **`src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md`** — load-order step 5 cross-references harness for dispatch mechanics; §2.5 Identity-rotation primitive line cross-references harness §1–§3; frontmatter `calls:` adds `harness/SKILL.md`.

4. **`src/packages/cnos.cdd/skills/cdd/beta/SKILL.md`** — Pre-merge gate Row 1 ("Identity truth") cross-references harness §3.2 (worktree-aware preventive rule) and §3.3 (failure-mode catalogue: #301 O8, #370 β R1, #370 α F4).

## Design decisions worth recording

- **Path choice.** `src/packages/cnos.cdd/skills/cdd/harness/` — peer to `operator/`, `gamma/`, `beta/`, `alpha/`, etc. The issue body's implementation contract pinned the path; no improvisation.
- **Frontmatter `scope:` value.** Initial draft used `substrate`; CUE schema (`schemas/skill.cue`) only allows `global` / `role-local` / `task-local`. Corrected to `global` (harness is global substrate, not role- or task-specific).
- **§5.2.1 Agent-tool `isolation: "worktree"` rephrasing.** The AC5 oracle pattern (`rg "...|worktree|..."`) is over-broad — it catches *any* occurrence of "worktree," including the Claude Code Agent tool's filesystem-isolation parameter mode, which is semantically distinct from git worktrees. Rephrased to "filesystem-isolation mode (the sub-agent operates on a copy of the repo, e.g. the Claude Code Agent tool's `isolation` parameter set to the per-agent-copy mode)" — preserves the doctrine, satisfies the AC5 mechanical oracle.
- **Phase 4a integration.** Phase 4a (cycle/397) merged to origin/main between dispatch and cycle/398's merge. Merge integration pulled origin/main into cycle/398 at `897a4552`. One conflict on operator/SKILL.md §8 Timeout recovery (Phase 4a kept the inline procedure with delta/SKILL.md §3 cross-refs; Phase 4b's pointer-form took precedence, with the delta/SKILL.md §3 cross-reference preserved in the override-declaration mechanics record). The Phase 4a quote-block at §3 line 122 was updated to name harness/SKILL.md for dispatch-coordinator mechanics, completing the 4a → 4b separation: 4a separated δ-role boundary policy from operator-as-coordinator + harness + release-effector mechanics; 4b separates dispatch-coordinator mechanics (now in harness/SKILL.md) from release-effector mechanics (remaining in operator/SKILL.md until Phase 4c).

## AC summary

AC1–AC6 PASS mechanically per `self-coherence.md` and `beta-review.md`. AC7 (close-out comments on #371 #373 #384 naming #398) is post-merge δ work (the comments are filed after the merge lands).

## Empirical anchors named

Each absorbed-issue rule traces to its falsification:

- §2.4 → cycle #369/#370 parallel dispatch in default text mode, 2026-05-17.
- §3.3 → cycle #301 O8 (first surfacing), cycle #370 β R1 §2.1 row 1 (second), cycle #370 α F4 (third).
- §4.2 → cph cdr-refactor wave 2026-05-18 to 2026-05-20 F1 (cph#27 + cph#28).

## What this cycle does NOT do (non-goals re-affirmed)

- No `cn dispatch` Go code reimplementation. `src/go/internal/dispatch/` is untouched.
- No movement of δ-role boundary content (Phase 4a — cycle/397).
- No movement of release-effector content (Phase 4c — pending).
- No change to the observability stream format substantively — codified what cnos#371 specified.

## Cycle metrics

- Commits on cycle/398: 6 (γ scaffold `df0e96bf`, α build `2966d133`, α self-coherence `a55f4b32`, β review `b497a9db`, merge resolution `897a4552`, β closeout `47dbd68c`) + this α closeout.
- Files changed: 7 (1 new harness skill, 3 cross-ref edits, 7 cycle evidence files).
- LOC: ~540 added in harness/SKILL.md; ~150 removed from operator/SKILL.md; ~10 added in cross-refs.
- Review rounds: 1 (APPROVED at R1; no fix-rounds).
