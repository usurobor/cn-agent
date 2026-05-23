# α design notes — cycle #398 harness substrate

**Issue:** cnos#398 (Phase 4b of #366; absorbs #371 #373 #384)
**Author:** α (collapsed on δ)
**Date:** 2026-05-21

## 1. Harness skill canonical path

**Decision:** `src/packages/cnos.cdd/skills/cdd/harness/SKILL.md`.

**Rationale:**
- Issue body's pinned implementation contract names this path explicitly ("or equivalent location per cdd skill conventions"); the literal name is the canonical choice.
- Sibling skill structure under `src/packages/cnos.cdd/skills/cdd/` already hosts every role and lifecycle surface (`alpha/`, `beta/`, `gamma/`, `operator/`, `release/`, `review/`, `post-release/`, `issue/`, `cross-repo/`, `activation/`, `design/`, `plan/`, `epsilon/`). A `harness/` peer is the lowest-friction placement.
- #366 Phase 4 names the substrate "Harness / platform-driver substrate"; the skill name `harness` matches the doctrine surface exactly.
- Phase 4c will add `release/` effector content (overlap with the existing release/SKILL.md is intentional — `release/SKILL.md` covers β-side release artifacts; Phase 4c's δ-side release driver is a different surface, possibly `release-driver/` or extension to existing `release/`). This cycle does not touch that decision.

## 2. Source-section partition

Maps from `operator/SKILL.md` → destination:

| operator/SKILL.md surface | Destination | Rationale |
|---|---|---|
| §Core Principle (δ-as-boundary) | **stays in operator** | δ-role WHY; not harness mechanics. |
| §Git identity for role actors (table + bash snippets) | **moves to harness** | Mechanics of identity write — harness contract per #373. |
| §Algorithm step list | **stays in operator** (rewritten as WHY, points at harness for HOW) | Step *intent* is operator coordination; step *execution* (claude -p shell) is harness. |
| §1 Route, §1.1, §1.2 (claude -p invocation shells) | **moves to harness** | Pure dispatch mechanics — the existing observability rule (#371 already-landed prose) lives here. |
| §2 Wait, §2.1 (do not poll internal work) | **stays in operator** | δ-discipline (when to be silent); not harness mechanics. |
| §2.2 (issue + branch polling loop, git fetch reliability) | **moves to harness** | Polling/wake-up mechanics — harness substrate. |
| §3 Gate (external actions, table) | **stays in operator** | δ gate authority WHY; the *list of actions δ holds* is doctrine, not mechanics. |
| §3.4 Cut the release — `scripts/release.sh` mechanics | **stays in operator** (Phase 4c will move) | Out of scope for 4b; Phase 4c owns release-effector relocation. Keep here untouched. |
| §3.5 The tag is the signal | **stays in operator** (Phase 4c) | Same as §3.4. |
| §3a δ as inward membrane | **stays in operator** (Phase 4a will move) | Phase 4a owns δ-role boundary; do not touch here. |
| §4 Override | **stays in operator** | δ doctrine. |
| §5 Dispatch configurations (§5.1 canonical multi-session, §5.2 single-session δ=γ, §5.2.1 quiescence, §5.3 escalation) | **stays in operator** | Dispatch *configuration* is δ-doctrine (when to choose what); the *shell commands* inside §5.1 are mechanics — extract those, leave the doctrine. |
| §6 What the operator does NOT do | **stays in operator** | δ boundary. |
| §7 Cycle lifecycle from operator's view | **stays in operator** | δ-view table. |
| §8 Timeout recovery (worktree inspection + decision tree + override declaration + prevention) | **moves to harness** | Pure mechanics of recovering an exited `claude -p` session. δ doctrine (override declaration) cross-refs back to operator §4. |
| §10 Wave Coordination | **stays in operator** | Multi-cycle δ coordination is δ-doctrine, not harness. Wave manifest mechanics could later move, but issue body's "extract harness" list does not name waves; leave to Phase 4a/separate cycle. |

## 3. Absorptions

### #371 (observability)

Codify `--output-format stream-json --verbose` JSONL contract as a δ-membrane primitive in the harness skill. The contract is already present in `operator/SKILL.md` §1.2 prose and `gamma/SKILL.md` §2.5 (as `Identity-rotation primitive` line). Harness skill becomes the canonical home.

Harness section: **§Dispatch observability contract.** Names:
- The flag combination (`--output-format stream-json --verbose --permission-mode acceptEdits`).
- Why `--verbose` is required (Claude Code refuses `-p` + `stream-json` without it).
- Per-cycle JSONL log path convention (e.g. `/tmp/cycle-{N}.jsonl` or `.cdd/unreleased/{N}/dispatch.jsonl`).
- JSONL line shape (or reference to `claude` SDK docs).
- The rule that black-box dispatch violates δ's gate authority.
- Falsification anchor: #369/#370 parallel dispatch incident, 2026-05-17.

### #373 (worktree identity write)

Codify the `extensions.worktreeConfig=true` check + `--worktree` discipline.

Harness section: **§Git identity for role actors** (relocated from operator) with explicit worktree subsection. Names:
- The session-start check: `git config --get extensions.worktreeConfig`.
- When `true`: every identity write MUST use `--worktree`.
- The shared-config leak failure mode (cycle #301 O8 / #370 β R1 + α F4).
- Cross-reference: β-skill Row 1 (Pre-merge gate "Identity truth").

β-skill Row 1: gets a cross-reference line to `harness/SKILL.md §Git identity` so the doctrine is discoverable from both sides.

### #384 (parallel α worktree pre-creation)

Codify pre-created per-cycle worktrees for parallel α dispatch.

Harness section: **§Parallel dispatch precondition** with:
- The same-repo race condition (cph#27/#28 incident).
- The `git worktree add /path/to/{project}-cycle-{N}` requirement, one per parallel α.
- Per-worktree identity write (`git config --worktree user.email alpha@{project}.cdd.cnos`).
- Cleanup contract: `git worktree remove` after β merges.
- Cross-reference: β-skill Row 1 (test-merge worktree discipline — same family).
- Empirical anchor: cph cdr-refactor wave 2026-05-18–05-20 F1.

## 4. Cross-reference plan

### β-skill Row 1 (Pre-merge gate "Identity truth")

Current: cites cycle #301 O8 and prescribes re-assertion. After this cycle: add a line *"Identity-write doctrine: see `harness/SKILL.md` §Git identity — when `extensions.worktreeConfig=true`, all identity writes MUST use `--worktree`."*

### γ-skill §2.5 dispatch reference

Current: references `operator/SKILL.md` and the `cn dispatch` line. After: add reference to `harness/SKILL.md` for dispatch invocation mechanics + observability contract. Keep the `## Implementation contract` section unchanged (Phase 4a doctrine).

Current load-order item 5 says "load operator/SKILL.md — δ owns dispatch execution …". Add a cross-ref to harness for the *execution mechanics*.

### operator/SKILL.md → harness

Where mechanics moved out, leave a one-line pointer ("See `harness/SKILL.md` §X for the executable form") so the doctrine surface stays discoverable for someone reading operator first.

## 5. AC mapping plan

| AC | Surface | Oracle |
|---|---|---|
| AC1 | new file with valid frontmatter | `test -f` + `head -1` (`---`) |
| AC2 | harness §Dispatch observability contract | `rg "output-format stream-json"` ≥ 1 |
| AC3 | harness §Git identity — worktree subsection | `rg "extensions.worktreeConfig.*true\|--worktree"` ≥ 1 + cross-ref to β Row 1 |
| AC4 | harness §Parallel dispatch precondition | `rg "pre-create\|worktree add.*cycle"` ≥ 1 |
| AC5 | operator/SKILL.md cleaned | `rg "claude -p\|cn dispatch\|polling loop\|stream-json\|worktree"` → 0 non-cross-ref hits |
| AC6 | γ + β cross-refs | `rg "harness/SKILL.md" gamma/SKILL.md beta/SKILL.md` ≥ 1 each |
| AC7 | issue close-outs | post-merge `gh issue close` on #371 #373 #384 with absorption comment |

## 6. Frontmatter for harness skill

```yaml
---
name: harness
description: Harness / platform-driver substrate for CDD. Codifies dispatch mechanics, observability contract, worktree management, and identity discipline. Invoked by δ; referenced by γ and β.
artifact_class: skill
kata_surface: embedded
governing_question: How does the harness substrate execute dispatch, surface observability, isolate worktrees, and enforce identity — coherently — across role sessions?
visibility: internal
parent: cdd
triggers:
  - harness
  - dispatch-mechanics
  - worktree
  - observability
  - identity
scope: substrate
inputs:
  - δ's dispatch decisions
  - γ's α/β prompts
  - active cycle branch state
outputs:
  - dispatched role sessions (one per `claude -p` invocation)
  - per-cycle JSONL observability streams
  - per-cycle worktree directories (when parallel α applies)
  - per-actor git identity (worktree-local when `extensions.worktreeConfig=true`)
requires:
  - active CDD cycle
  - cycle branch exists on origin
calls:
  - operator/SKILL.md (δ doctrine — what gate the harness serves)
---
```

## 7. AC5 self-check — what `rg` patterns survive in operator after extraction

After extraction, operator may still contain:
- `cn dispatch` in §3a (Phase 4a-owned section that doesn't move here) — acceptable: this is an inward-membrane reference, not dispatch mechanics.
- `claude -p` in §Algorithm step list — must rewrite to remove the literal flag form, leaving only "δ dispatches γ via the harness substrate; see `harness/SKILL.md` for the invocation form."
- `worktree` in §3a — does not appear (§3a is about implementation contract, not worktrees).

Plan: after extraction, do a clean grep pass. If §3a's `cn dispatch` reference is the only `cn dispatch` remaining, that's acceptable (it's a doctrine reference to the *capability*, not a mechanics restatement). AC5's oracle should distinguish "mechanics hit" vs "cross-ref / doctrine pointer hit" — the issue's literal oracle is "0 hits OR only in cross-references to harness skill." So cross-references like "See `harness/SKILL.md` §X" do NOT count as failure.

To make AC5 trivially mechanical, I'll keep the operator §1/§1.2/§2.2 etc. clean (no `claude -p` literal, no shell snippets) and let §3a stand (its `cn dispatch` reference is a capability mention, not the dispatch shell). The AC5 oracle pattern `rg "claude -p|cn dispatch|polling loop|git fetch.*polling|worktree|stream-json"` will hit on §3a's `cn dispatch` mention — but the issue body explicitly says "or only in cross-references to harness skill." A capability reference in §3a is a cross-reference at heart; I'll surface the §3a `cn dispatch` retention as a *deliberate* cross-reference and note it in self-coherence.

Actually, looking more carefully: §3a's only `cn dispatch` mention is in the prose describing the inward-membrane function — it cites `gamma/SKILL.md` §2.5 Step 3b, not `cn dispatch` directly. Let me re-check.
