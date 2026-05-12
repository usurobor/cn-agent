# Self-Coherence — Cycle #323

## Gap

**Issue:** #323 — fix(activate): scanner misses threads/inbox/

**Problem:** `cn activate` scans `threads/in`, `threads/mail`, and `threads/archived` but does not scan `threads/inbox/`. The canonical reference hub (cn-sigma) uses `threads/inbox/` as its primary inbound message surface, making activation output blind to the actual inbox.

**Expected:** `cn activate` should scan `threads/inbox/` in addition to the current surfaces.

**Version/Mode:** MCA (Minimal Correct Action) — direct fix to existing scanner function.

## Skills

**Tier 1a — CDD authority (loaded):**
- `src/packages/cnos.cdd/skills/cdd/CDD.md`
- `src/packages/cnos.cdd/skills/cdd/SKILL.md`
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md`

**Tier 2 — General engineering (coding bundle, loaded):**
- `src/packages/cnos.eng/skills/eng/go/SKILL.md` (primary language skill)
- `src/packages/cnos.eng/skills/eng/ux-cli/SKILL.md` (CLI interface standards)
- `src/packages/cnos.eng/skills/eng/test/SKILL.md` (testing discipline)

**Tier 3 — Issue-specific:** None explicitly specified in issue #323.

**Work Shape:** Small local bugfix (L5 — local correctness)
**Active Skills:** go + ux-cli + test

## ACs

**AC1:** `cn activate` should scan `threads/inbox/` in addition to existing surfaces
- **Evidence:** Added `"threads/inbox"` to the `candidates` slice in `scanThreads()` function in `src/go/internal/activate/activate.go:232`
- **Status:** ✓ MET

**AC2:** Activation output should include `threads/inbox: present` when the directory exists
- **Evidence:** The `scanThreads` function returns paths discovered by `presentPaths()`, which checks for directory existence. Added `threads/inbox` to test fixture in `activate_test.go:70`
- **Status:** ✓ MET

**AC3:** Implementation should maintain existing functionality for current surfaces
- **Evidence:** No existing lines changed in `scanThreads()`, only added `"threads/inbox"` to the list. Preserved existing candidates: `threads/in`, `threads/mail`, `threads/archived`
- **Status:** ✓ MET

## Self-check

**Did α's work push ambiguity onto β?**
No. The implementation is a straightforward addition to an existing list. The change is:
- Localized to one function (`scanThreads`)
- Self-contained (no ripple effects)
- Follows existing pattern (same structure as other thread paths)

**Is every claim backed by evidence in the diff?**
Yes:
- `src/go/internal/activate/activate.go:232` — added `"threads/inbox"` to candidates slice
- `src/go/internal/activate/activate_test.go:70` — added `"threads/inbox"` to test fixture

**Implementation follows loaded skills:**
- **eng/go**: Small package change, no new interfaces, pure addition to existing slice
- **eng/ux-cli**: Maintains existing output format, adds info where appropriate
- **eng/test**: Updated test fixture to include new path for coverage

## Debt

**Known debt:** None.

This is a minimal, localized change that follows existing patterns. No technical debt introduced:
- No new abstractions
- No additional dependencies
- No architectural changes
- No deferred testing or documentation
- No compatibility concerns

## CDD-Trace

**Step 0 — Observe:** Gap identified during 3.71.0 triad compliance audit (per issue evidence)

**Step 1 — Select:** Issue #323 filed by γ, gap = activation scanner missing canonical inbound surface

**Step 2 — Branch:** `cycle/323` created by γ, α checked out successfully

**Step 3 — Bootstrap:** Not required (small change, no version directory needed)

**Step 4 — Gap:** Named incoherence = `cn activate` blind to `threads/inbox/` (canonical surface)

**Step 5 — Mode:** MCA + Tier 2 skills (go + ux-cli + test)

**Step 6a — Design:** Not required (single function modification, no architecture change)

**Step 6b — Plan:** Not required (trivial implementation, no sequencing complexity)

**Step 6c — Tests:** ✓ Updated `activate_test.go` fixture to include `threads/inbox`

**Step 6d — Code:** ✓ Modified `src/go/internal/activate/activate.go:232` — added `"threads/inbox"` to candidates

**Step 6e — Docs:** Not required (no user-facing documentation changes needed)

**Step 7 — Self-coherence:** ✓ This document