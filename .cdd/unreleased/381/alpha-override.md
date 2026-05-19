# α-direct override — cnos#381 (parser: read-first-order block)

**Date:** 2026-05-19
**Issue:** [cnos#381](https://github.com/usurobor/cnos/issues/381) — parser matches docstring marker mention before real block
**Cycle:** small-change override (no β review, no γ scaffold)
**Override authority:** operator (usurobor) — explicit direction "α-direct, fix it and ship it" during cnos 3.79.0 disconnect phase
**Branch:** `cycle/381` → merged to main → released as 3.80.0

## Override declaration (per operator/SKILL.md §4)

1. **What is overridden:** the full triadic dispatch (γ scaffold + α implement + β review + close-out) for cnos#381. Replaced with α-direct implementation + single commit on `cycle/381`.
2. **Why:** cnos#381 is a 1-AC parser bugfix with the recommended fix (option a — markers on their own line) named in the issue body. The change touches one function (~15 LOC delta) and adds two regression tests. The full cycle ceremony for a parser fix this small and this well-scoped is not load-bearing — β review would have nothing structural to verify beyond the test additions.
3. **New state:** α writes the parser change + two regression tests on `cycle/381`. Operator merges to main (or α merges per override scope). Release shipped as 3.80.0 via `scripts/release.sh`. Cycle artifact set collapses to this override.md only (validate-release-gate.sh small-change branch, no `beta-review.md`).

## Implementation

### Parser change

`src/go/internal/activate/activate.go::parseReadFirstOrderBlock` converted from substring-based (`strings.Index(content, marker)`) to line-based scan. The marker must appear as its own line content after `strings.TrimSpace` for the line to qualify as the begin or end marker. Effect: backtick-wrapped mentions of the marker in prose (e.g. the §4.1 docstring describing the format) no longer match ahead of the real block.

```go
// before (3.79.0)
bi := strings.Index(content, readFirstOrderBeginMarker)
// ... finds the docstring mention at line 299 instead of the real block at line 301

// after (3.80.0)
lines := strings.Split(content, "\n")
for i, raw := range lines {
    trimmed := strings.TrimSpace(raw)
    if begin < 0 && trimmed == readFirstOrderBeginMarker { begin = i }
    else if begin >= 0 && trimmed == readFirstOrderEndMarker { end = i; break }
}
```

The parsing of numbered list items inside `lines[begin+1:end]` is unchanged from 3.79.0.

### Regression tests

`src/go/internal/activate/activate_test.go`:

- `TestParseReadFirstOrderBlock_IgnoresDocstringMention` — fixture with a §4.1-style preamble (`Format: between \`<!-- read-first-order:begin -->\` and \`<!-- read-first-order:end -->\`, ...`) BEFORE the real block. Pre-fix: parser matched the docstring mention, parsed the prose between marks as the block, returned `ok=false`. Post-fix: parser scans line-by-line, only matches the real block, returns 3 items.
- `TestParseReadFirstOrderBlock_IndentedMarker` — fixture with leading whitespace on the marker lines. Confirms `TrimSpace` handles the leading-whitespace case, so markers indented inside list items or quotes still match.

Existing tests `TestParseReadFirstOrderBlock_HappyPath` + `TestParseReadFirstOrderBlock_NoMarkers` continue to pass unchanged.

## Verification

Local: `go test -count=1 ./internal/activate/` — full activate suite green (4 parser tests + 25 other tests).

Live: built `/tmp/cn-381` from `cycle/381` HEAD; ran against `/root/cn-sigma` (vendored cnos.core@3.79.0):
- 3.79.0 binary: stderr contains the false `→ activate skill not vendored ...` warning
- 3.80.0 binary: stderr does NOT contain the warning
- stdout from both: bytes-equal (fallback ordering already matched the skill on commit, so the rendered prompt was correct in both — the fix only removes the misleading warning)

Operator-visible projection: an operator running `cn activate HUB_DIR` against a hub with a current vendored cnos.core@3.80.0 or later no longer sees the contradictory warning. The skill-as-source-of-truth invariant is now restored in observable output.

## Known scope-collapse implications (per operator/SKILL.md §3.8 release/SKILL.md grading)

- No β review = no APPROVED verdict. The cycle's grade per `release/SKILL.md §3.8` configuration-floor clause caps at the override level. For this 1-AC mechanical patch with named recommended fix, that's acceptable.
- No γ-closeout = no PRA, no cdd-iteration.md. Coverage gap: the test for the docstring-mention case (which a γ-flagged failure mode would have anticipated) is added in this same commit, so the failure mode is closed at the test layer. Future cycles touching `parseReadFirstOrderBlock` should re-run `TestParseReadFirstOrderBlock_IgnoresDocstringMention` to prevent regression.

## Filed by

α-direct under operator override. cn-sigma δ session of 2026-05-19; closes cnos#381 by shipping the parser fix in 3.80.0.
