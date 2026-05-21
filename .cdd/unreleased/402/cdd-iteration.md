<!-- sections: [intake, finding-1, finding-2, summary] -->
<!-- completed: [intake, finding-1, finding-2, summary] -->

# cdd-iteration — cycle/402

## Intake

Per the closure gate, cdd-iteration.md is required when the cycle's close-out triage produces ≥1 finding with a `cdd-*-gap` class. This cycle's triage (gamma-closeout.md §"Close-out triage table") names two findings (α F1 and β F2) classified as `cdd-protocol-gap` (the same gap, surfaced independently by α and β). One iteration record per finding-class follows.

## Finding 1 — anchor-granularity collapse for legacy cross-references

**Source:** α F1 (alpha-closeout.md) + β F2 (beta-closeout.md) — same gap surfaced independently by α and β.

**Class:** `cdd-protocol-gap`.

**Trigger:** γ process-gap check (no §9.1 mechanical trigger fired; this gap surfaces structurally from the terminal-phase compression contract).

**Description:** The pre-cycle `CDD.md` carried very specific cross-reference anchors that other cdd/ skill files cited by exact section/step number (`§1.4 γ algorithm Phase 1 step 3a`, `§1.6a`, `§1.6c(a)`, `§5.2`, `§5.3a`, `§5.3b`, `§9.1`, `§Tracking`, `Phase 6 step 17`). The Phase 7 compression to a CCNF-spine doctrine surface cannot preserve that anchor granularity within the 300-line budget — the granular substrate-specific anchors live at canonical depth in the substrate-specific role / runtime-substrate SKILL.md files, not in the kernel doctrine. The compression therefore collapses all legacy anchors to family-level pointers inside the new CDD.md §"Software-specific realization — pending cds extraction" section.

**Affected surfaces:** `post-release/SKILL.md`, `release/SKILL.md`, `release-effector/SKILL.md`, `alpha/SKILL.md`, `beta/SKILL.md`, `gamma/SKILL.md`, `operator/SKILL.md`, `harness/SKILL.md`, `activation/SKILL.md`, `review/SKILL.md`, plus the 2 CI workflow templates under `activation/templates/`. 29 total cross-reference hits across these surfaces.

**Disposition:** patch-pending-followon. The durable fix is the cds extraction (tracked at [cnos#403](https://github.com/usurobor/cnos/issues/403)). When the cnos.cds package bootstraps, each citing skill file re-points its `CDD.md §X` anchor at the cds package directly (`cds/SKILL.md §X` or similar), and the legacy anchor forms are retired. Until #403 lands, the family-level resolution in CDD.md §"Software-specific realization" is the contract — every legacy anchor resolves to the family that owns its content.

**Patch-or-issue link:** [cnos#403](https://github.com/usurobor/cnos/issues/403) (filed as part of this cycle).

**Affected cdd skill file/section:** `CDD.md §"Software-specific realization — pending cds extraction"` (the named-quarantine section); operational expansions live in `gamma/SKILL.md`, `alpha/SKILL.md`, `beta/SKILL.md`, `operator/SKILL.md`, `harness/SKILL.md`, `release/SKILL.md`, `release-effector/SKILL.md`, `post-release/SKILL.md`, `activation/SKILL.md`, `review/SKILL.md`.

## Finding 2 — protocol gap: doctrine-compression cycles ship ahead of corresponding extraction cycles

**Source:** α F2 (alpha-closeout.md) — observation about the pattern of forward-referencing package bootstraps that lag the compression.

**Class:** `cdd-protocol-gap` (latent — surfaced here as a pattern observation; not blocking).

**Trigger:** γ process-gap check.

**Description:** Phase 7 of cnos#366 shipped the CDD.md compression without the corresponding cnos.cds package bootstrap (Wave 3 of the essay). The compression was shippable because the schemas already exist (`schemas/cds/` + fixtures per #388) and V can validate against them — the package bootstrap (skills, role overlays, CDS.md doctrine) is the missing piece. The compression cycle therefore filed a tracker issue (cnos#403) and quarantined the software-realization content in CDD.md under "pending cds extraction" naming.

This is a *protocol-gap* observation, not a *protocol-failure*. The pattern works: the compression cycle ships honestly, the tracker carries the named obligation forward, and the next extraction cycle closes the loop. But the pattern is worth naming because it surfaces a class of cross-cycle dependency: cycle A compresses doctrine that depends on cycle B (extraction) which lags A. Future doctrine-compression work may want a more formal "compression-ahead-of-extraction" contract — e.g., a known-debt schema in self-coherence.md that names "this content stays in surface X until extraction-cycle #Y lands."

**Affected surfaces:** the cdd protocol surface itself (the cycle-shape contract). The shape currently used (tracker issue + named-quarantine section + close-out finding) works; codifying it into the cdd protocol is a follow-on consideration.

**Disposition:** no-patch (pattern observation; not a protocol failure). The pattern is documented in this cycle's artifacts. If the cds extraction (#403) surfaces additional structural friction with the compression-ahead-of-extraction shape, a future cdd-iteration cycle may codify it. Until then, this is a one-observation pattern, not a recurring friction class.

**Patch-or-issue link:** none. The observation lives in this cdd-iteration.md row; if it recurs, a future cycle files a structural-iteration issue.

**Affected cdd skill file/section:** none directly (no skill patch). The observation belongs to the cycle-pattern surface, not to a specific skill file.

## Summary

Two findings in the `cdd-protocol-gap` class, both surfaced by the Phase 7 terminal compression:

1. **Anchor-granularity collapse** — patch-pending-followon at cnos#403; durable fix is the cds extraction.
2. **Compression-ahead-of-extraction pattern** — no-patch pattern observation; codification deferred unless the pattern recurs.

No `cdd-skill-gap`, `cdd-tooling-gap`, or `cdd-metric-gap` findings in this cycle. The cdd/design skill governed the rewrite cleanly; the cdd/issue/proof oracles caught the AC3 token-count subtlety on the first oracle pass.

The cycle's `cdd-protocol-gap` findings are pre-known structural consequences of the terminal-phase compression; neither indicates an unanticipated skill or tooling failure. The compression contract was pinned tightly enough that these gaps were anticipated (the issue body's "active design constraints" name the cds-extraction-as-follow-on shape; the dispatch named the tracker-issue-as-prerequisite shape).
