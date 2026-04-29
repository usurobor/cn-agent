<!-- sections: [context, friction-observed, proposed-stack, sequencing, disposition] -->
<!-- completed: [context, friction-observed, proposed-stack, sequencing, disposition] -->

# γ iteration proposal — after 3.61.0

**Author:** γ (`gamma@cdd.cnos`)
**Date:** 2026-04-29
**Anchor release:** 3.61.0 (cycle #283 closed at `fbaaa1a` — *"Cycle #283 closed. Next: #287."*)
**Status:** proposal — operator-pending; spawns issues #287 (filed), #288–#292 (pending operator confirmation)

## Context

Cycle #283 shipped a system-shaping MCA at L7 diff: triadic CDD coordination shifts from PR-mediated to artifact-driven. The cycle execution itself capped at L5 (two §9.1 trigger fires; `docs/gamma/cdd/3.61.0/POST-RELEASE-ASSESSMENT.md` §4b). The operator surfaced the cycle's friction explicitly: *"This was a messy cycle with agents not being able to communicate and confused about what to watch and what to do and what to check etc."* This document catalogs that friction and proposes the iteration stack to clean it up.

Some of the friction is already addressed by issues already filed (#286, #287). Some is not. This proposal names the gap and recommends six follow-up cycles to close it, with sequencing.

## Friction observed during #283

| # | Mess | Where surfaced | Already addressed? |
|---|---|---|---|
| M1 | β's `Monitor` `git fetch --quiet` silent transition drops (3 α-SHA transitions missed) | β close-out O3 | ✅ #287 AC8 |
| M2 | β's harness-pre-provisioned per-role branch with α-style instructions | β close-out O6 | ✅ #287 AC11 |
| M3 | γ glob-matched `'origin/claude/*'` against harness-encoded branch names | γ self-obs | ✅ #287 (γ creates `cycle/{N}` directly) |
| M4 | γ's Monitor timed out 5+ times despite `persistent: true` (harness ignored persistent) | γ in-cycle | ❌ harness gap — *not a CDD spec issue* |
| M5 | MCP issue-body cache staleness (γ read `.cn/` on first fetch, body had `.cdd/`) | γ in-cycle | ❌ no cache-bust rule yet |
| M6 | γ edited issue body twice mid-cycle, invalidated α's first commit | PRA γ-O1 | ⚠ partial via #286 AC6 named decision-points |
| M7 | γ committed scope expansion `eb48e17` *after* β R2 approval, then rolled back | PRA γ-O2 | ⚠ partial via #286 AC6 |
| M8 | β stream-idle timeout mid-`beta-closeout.md` write | operator surfaced | ❌ no resumption spec |
| M9 | "β tag at step 8 vs δ tag at step 17" — ambiguous tagging model | γ → β handoff | ❌ CDD spec ambiguity |
| M10 | β wrote provisional `L6 cycle cap` in CHANGELOG, γ had to revise to `L5` | this cycle's PRA | ❌ provisional/final scoring sequence not explicit |
| M11 | γ had to assemble TLDR manually (issue + branch tip + cycle-dir + open findings + gate state) | "TLDR current state" | ❌ no `cn cdd status` command |
| M12 | Operator was the relay channel for α↔β↔γ throughout the cycle | structural | ⚠ #286 (encapsulation, blocked on `cn dispatch` CLI) |
| M13 | Three filename conventions in one cycle (`{role}/{N}.md` → `{N}/{role}.md` → `{N}/`-descriptive) | β-O1 first-cycle pattern | ✅ explicit drop with reason (irreducible) |
| M14 | γ session branch (`claude/gamma-skill-issue-283-nuiUZ`) underspecified | γ in-cycle | ❌ spec doesn't formalize γ session branch |

## Proposed iteration stack

### #287 — γ creates the cycle branch ✅ filed

Direct response to #283 R1 F1's branch-discovery friction. 12 ACs across CDD.md (§1.4 algorithms + §Tracking + §4.2/§4.3 Branch rule), gamma/alpha/beta/operator SKILL.md, dispatch-prompt format. Self-application: implementing cycle uses `cycle/287`. Closes M1 (via AC8 `git fetch` reliability), M2 (β refusal expansion), M3 (γ creates named branch).

### #288 — Resumption protocol for partial CDD artifacts (proposed; closes M8)

§1.4 Large-file authoring rule expanded with a "Resumption" sub-clause + canonical section-manifest format at the top of any artifact ≥ 50 lines:

```markdown
<!-- sections: [a, b, c, d, e] -->
<!-- completed: [a, b] -->
```

Plus role-skill resumption sub-sections (alpha/beta/gamma) pointing at §1.4. Plus a future-harness `cn resume` directive that detects the manifest and re-spawns the role.

**Why a separate issue:** the failure is genus-level (any role, any artifact, any session timeout). #287's branch-rule changes don't touch it. Markdown-only, ~5 file edits.

### #289 — Disambiguate cycle-tag vs disconnect-tag, lock provisional/final scoring sequence (proposed; closes M9, M10)

CDD spec ambiguity exposed during this cycle's β→γ handoff: §1.4 Phase 2 step 8 says "β merges + tags + deploys"; Phase 6 step 17 says "δ cuts the disconnect release: bump version, tag, push." Two reads possible:

- **(A) β does not tag.** β commits release artifacts (CHANGELOG, RELEASE.md, version bump) but never executes `git tag`. δ creates the single tag per cycle as part of the disconnect release.
- **(B) Two-tag model is canonical.** β tags `{X.Y.Z}` at merge as cycle release; δ tags `{X.Y.Z+1}` (or `{X.Y.Z}.1`) as disconnect release.

Issue locks one model and updates §1.4 + `release/SKILL.md` + `operator/SKILL.md`. Same issue locks provisional-vs-final scoring sequence: β writes provisional CHANGELOG row with "*provisional, pending γ PRA*" in the level cell; γ updates to final score in the same commit as the PRA.

**Why a separate issue:** surgical 5-line CDD.md edit + 2-line `release/SKILL.md` + 2-line `post-release/SKILL.md` edit. Markdown-only.

### #290 — `cn cdd status N` CLI command (proposed; closes M11)

A single `cn` invocation that γ (or operator) runs to get a structured TLDR:

```
$ cn cdd status 283
issue:        #283 open    "Replace GitHub PR workflow with..."
branch:       cycle/283  →  fc50265 (α R2, +0/-0 since R1)
artifacts:    .cdd/unreleased/283/
              ├─ self-coherence.md      (α, R2, ready)
              ├─ beta-review.md         (β, R1+R2, APPROVED)
              ├─ gamma-clarification.md (γ, F1 resolved)
              └─ alpha-closeout.md      ❌ missing
gate state:   8/10 (α-closeout missing; merged-branch cleanup pending δ)
next:         γ writes closure declaration when α-closeout lands
```

Reads only; commits nothing. β-axis infrastructure cycle, separate from CDD spec changes. Naturally pairs with `cn dispatch` (#291 below) — both could be batched into a single "cn cdd commands" cycle if scope allows.

### #291 — `cn dispatch` CLI (proposed; unblocks #286, closes M12 partially)

The harness precondition for #286 (encapsulation) named in #286 AC4. γ invokes it to spawn α/β as sibling top-level sessions sharing only the cycle branch + artifact channel. Without this, γ can draft prompts but cannot autonomously dispatch — operator still pastes.

**Why a separate issue:** β-axis infrastructure work, identified in #286 Dependencies but not yet a filed issue. Hard precondition for #286 to land meaningfully.

### #292 — Issue-edit cache-bust convention + γ session branch formalization (proposed; closes M5, M14, partial M6/M7)

Two small spec items, naturally bundled:

1. **Issue-edit cache-bust:** when γ edits an issue body during an active cycle, γ writes a `gamma-clarification.md` entry on the cycle branch describing the edit (date + summary + which ACs / constraints / artifacts changed). Roles polling the cycle branch see the transition; their next intake re-fetches the issue from the live source rather than trusting cached state.

2. **γ session branch formalization:** spec names the γ session branch (e.g. `gamma/session-{N}` or just allows the harness-given `claude/...` shape) as γ's drafting surface for work that cannot land directly on the cycle branch. Default: γ commits directly to `cycle/{N}` for all cycle artifacts. The session branch is for pre-publication drafts only and must merge or be discarded by closure.

**Why a separate issue:** small surgical spec edits; CDD §Tracking gets the cache-bust rule, §1.4 / §4.2 get the session-branch clarification.

### #286 — Encapsulate α and β behind γ ✅ already filed

Already filed during this cycle. Hard precondition: #291 (`cn dispatch` CLI). Forward-looking; AC4 names the gap so the spec lands cleanly even if the CLI is built later.

### Out of CDD scope: harness gap (M4)

The harness silently ignores `persistent: true` and caps Monitors at 30min. γ re-armed 5+ times this cycle. **This is a Claude Code harness bug, not a CDD spec issue.** Surface to the harness team via whatever channel reaches them.

## Recommended sequencing

```
#287  γ creates the cycle branch  (next MCA — direct response to #283 R1 F1)
 ↓
#288  Resumption protocol         ┐
#289  Tag/scoring disambiguation  │ small-change cycles, parallel-dispatchable
#292  Cache-bust + session branch │ if operator dispatches in batch
                                  ┘
 ↓
#290  cn cdd status               ┐ infra cycles, β-axis, can pair
#291  cn dispatch                 ┘ if scope allows
 ↓
#286  Encapsulation (depends on #291)
```

#287 first: it's #283's direct child. #288 / #289 / #292 are small, parallel-safe, and unblock cleaner cycles immediately. #290 / #291 are the infra unlock for true encapsulation. #286 lands last when its precondition (#291) exists.

## Disposition

This document is γ's strategic record of what came out of #283 beyond the cycle's own scope. It is not itself a CDD cycle artifact. It does not commit any of the issues — operator confirms each before γ files.

Each filed issue (when filed) gets the standard shape of #283/#287 (Problem / Impact / ACs / Non-goals / Tier 3 / Constraints / Related artifacts / Dependencies / Priority / Work-shape). γ drafts each on operator confirmation; operator dispatches.

Updates to this document: when each proposed issue is filed, mark "✅ filed (#NNN)" inline. When all six are filed, this proposal moves to "executed" and stays as the historical record of the iteration sweep that followed 3.61.0.
