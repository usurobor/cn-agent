# 3.64.0 — CDD merge/release authority separation, README v4, #307 kata move

## Outcome

Coherence delta: C_Σ B+ (`α B+`, `β A-`, `γ B+`) · **Level:** `L5 (cycle cap: L5)`

Three changes ship: a spec-level CDD authority separation (β merges, δ releases), a full README rewrite for v4 convergence, and a small-change kata relocation (#307). The release also records the first live identity-rotation dispatch test.

## Why it matters

CDD's merge/release authority was ambiguous — β owned "merge, tag, deploy" as a bundle, and δ's role blurred between authorization and execution. This release separates the two: β is "Reviewer + Integrator" (merge into main requires no δ), δ owns the release boundary (preflight, tag, deploy, disconnect). A new Phase 5a (δ release-boundary preflight) gates the release without gating the merge. The README rewrite makes cnos legible to outsiders for the first time — thesis-first, layers diagram, honest shipped/not-shipped split.

#307 is notable not for the kata move itself (P3 mechanical) but as the first live test of the identity-rotation dispatch model: γ (Sigma/OpenClaw) dispatched α (Opus via `claude -p`) and β (Sonnet via `claude -p`) on a 2GB VPS. Ten implementation findings are recorded in the hub.

## Changed

- **CDD merge/release authority separation**: β is now "Reviewer + Integrator" — owns review + merge into main (no δ authorization required). δ owns the release boundary: preflight, tag, deploy, disconnect. New Phase 5a added to γ algorithm. Step 9 renamed from "Gate + merge" to "Merge" with "(β authority — no δ required)". Role key, default flow, and coordination loop updated. `operator/SKILL.md` gate table updated.
- **README rewrite**: thesis-first structure ("recurrent coherence system with Git as its lowest durable substrate"), layers diagram, coherent agents section with MCP/CAP/CLP/CDD operations table, network model, "Why Git?" section, honest shipped vs not-shipped split, v4 roadmap phase table, hub shape diagram, organized further-reading by concern.

## Added

- **#307 kata move**: Three embedded katas moved from `issue/SKILL.md` §5 to `src/packages/cnos.cdd.kata/katas/M5-issue-authoring/`. `issue/SKILL.md` frontmatter updated to `kata_surface: external` with `kata_ref`. `docs/gamma/cdd/KATAS.md` updated with M5 row.
- **δ release-boundary preflight** (Phase 5a in γ algorithm): After β merges and close-outs exist, γ requests δ preflight before closing the cycle. δ verifies merge commit, release artifacts, tag/deploy preconditions. Three outcomes: proceed, request changes, override.

## Validation

- Build CI: pass (main push)
- Release smoke: pending (this release)
- #307 review: 2 rounds (R1 RC → R2 Approved), clean merge
