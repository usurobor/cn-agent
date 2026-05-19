# 3.80.0

## Outcome

Coherence delta: C_Σ **override** (α-direct; no β review; per operator/SKILL.md §3.8 configuration-floor) · **Level:** `L4`

A false `→ activate skill not vendored ...; using built-in ordering` warning emitted by `cn activate` since 3.78.0 — even when the skill IS vendored — is closed. The parser now matches the read-first-order block markers structurally (line-level identity after `TrimSpace`) instead of by substring, so docstring mentions of the marker syntax in §4.1 prose no longer match ahead of the real block.

## Why it matters

cnos 3.79.0 left a small but operator-visible contradiction in observable output. cn-sigma δ ran `cn activate /root/cn-sigma` against a hub with `cnos.core@3.79.0` vendored — including the canonical `agent/activate/SKILL.md` at the path the renderer reads — and the renderer still emitted the warning that the skill was not vendored. Functionally harmless (the in-Go fallback ordering matched the skill on commit, so rendered output was correct), but the warning undermined the cycle/379 claim that the skill is the source of truth: an operator inspecting stderr would conclude the renderer fell back to in-Go literals despite the skill being present.

Root cause: `parseReadFirstOrderBlock` used `strings.Index(content, "<!-- read-first-order:begin -->")` to find the marker. The shipped SKILL.md describes the marker format in its §4.1 prose using literal backtick-wrapped markers (`Format: between \`<!-- read-first-order:begin -->\` and \`<!-- read-first-order:end -->\`, ...`). The parser matched the docstring mention at line 299 ahead of the real block at line 301, then tried to parse the prose between marks as a numbered list, returned `ok=false`, and the fallback fired. The existing parser tests used a clean fixture with no docstring mention, so the failure mode was invisible to the test suite.

This release fixes the parser and locks in the regression with two new tests.

## Added

- **#381** — `src/go/internal/activate/activate.go::parseReadFirstOrderBlock` converted from substring-based to line-based scan. Markers must appear on their own line after `TrimSpace` to qualify. Backtick-wrapped or otherwise embedded mentions in prose no longer match.
- **#381** — `src/go/internal/activate/activate_test.go` — `TestParseReadFirstOrderBlock_IgnoresDocstringMention` (fixture mirrors the shipped SKILL.md §4.1 shape: docstring preamble + real block + trailing-prose mention; pre-fix returns `ok=false`, post-fix returns the 3 items from the real block) and `TestParseReadFirstOrderBlock_IndentedMarker` (markers with leading whitespace match via `TrimSpace`).

## Changed

- The parser's `ok=true` outcome on the shipped SKILL.md. 3.78.0/3.79.0 returned `ok=false` and fell back to the in-Go `canonicalReadFirstOrdering`. 3.80.0 returns `ok=true` and uses the skill block as the source of truth. **Rendered prompt is identical** because the fallback already matched the skill — the change is in which code path produces the ordering, not what ordering is produced. Existing tests that inspect the rendered output continue to pass without modification.

## Unchanged

- Default `cn activate HUB_DIR` (no flag) behavior — stdout content identical to 3.79.0.
- `--claude` and `--codex` spawn paths — unchanged.
- `agent/activate/SKILL.md` — no edits. The skill prose still describes the marker format using backtick-wrapped marker mentions; that is now correctly tolerated by the parser.
- `cn activate`'s observable output for any hub without a vendored cnos.core@3.78.0+ — same fallback path, same rendered prompt, same warning.

## Cycle shape

α-direct override per `operator/SKILL.md §4`. Operator directed "α-direct, fix it and ship it" during the cnos 3.79.0 disconnect phase. Scope: 1 AC (recommended fix named in issue body), ~15 LOC parser change + 2 regression tests. No γ scaffold, no β review, no γ-closeout. The cycle's small-change branch in `validate-release-gate.sh` collapses artifact requirements; this release ships with `.cdd/unreleased/381/alpha-override.md` as the sole cycle artifact, naming the override declaration + the fix narrative.

Grade implication per `release/SKILL.md §3.8` configuration floor: cycle grade caps below A (no β verdict to anchor against). Acceptable for a P3 mechanical bugfix where the failure mode is functionally inert and the fix is unambiguous.

## Empirical anchor

Surfaced 2026-05-19 by cn-sigma δ while bumping cn-sigma's `.cn/deps.json` to cnos 3.79.0 and verifying `cn activate /root/cn-sigma` against the freshly-vendored skill. Both `--claude` and `--codex` demos succeeded; the parser warning was the only operator-visible noise on a hub that should have been fully on the canonical path. Logged in cn-sigma `threads/reflections/daily/20260519.md` §"cnos#381 — read-first-order parser bug discovered during deps verification".
