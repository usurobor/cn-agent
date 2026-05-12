---
cycle: 352
issue: "#352"
type: self-coherence
date: 2026-05-12
---

# Self-Coherence — Cycle #352

## §Gap

**Issue:** #352 — cdd/review+gamma: Add CI-green gate — β refuses merge on red CI, γ PRA verifies post-merge

**Version/Mode:** docs-only, design-and-build

**Selected gap:** CDD protocol has no CI-green verification mechanism. β can APPROVE on red CI; γ can close-out without verifying post-merge CI status. This creates a latent failure mode where structurally correct workflows (readable by β) can still fail at runtime, and neither role catches it. Issue emerged from tsc #36 cycle which shipped a CI gate but had no mechanism for verifying its *own* CI ran green post-merge.

**Coherence target:** Add symmetric CI-green gates on both sides of merge:
1. β refuses APPROVED without CI green on review SHA  
2. γ PRA polls CI post-merge before close-out

**Mode justification:** docs-only appropriate because this is rule addition to existing skills, not new code surfaces. design-and-build because it requires coordination across 4 skill files with a cohesive verdict-rule + close-out-rule pair.