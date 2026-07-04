# α Close-out — cnos#584

## Summary

Clean R0 cycle. γ scaffolded, α implemented, β converged with verdict APPROVE and zero findings on the first round — no fix round, no re-dispatch loop. Diff: 4 files (`gamma-scaffold.md`, `self-coherence.md`, `CDD.md`, `beta/SKILL.md`), all `.md`, 430 insertions / 2 deletions. All 4 ACs independently re-verified by β against the diff (not against α's transcription); Rule 7 implementation-contract conformance held on all 7 pinned axes.

## What worked

- **Scaffold pre-specification narrowed α's design surface.** γ's scaffold pinned the canonical-home candidate set to exactly `{CDD.md, CELL-OF-CELLS.md}` and pre-specified the AC3 resolution path (prose-only qualification, no behavior change) via friction note 2. This left α with one genuine judgment call — which of the two candidate files, and how to word the correction — rather than an open design space, which is consistent with why design/plan were marked "not required" (§Debt #3) without β objecting to that call.
- **Oracle-first self-coherence let β's re-verification converge without friction.** Every AC in `self-coherence.md` cited a literal command + literal output rather than a narrative claim. β re-ran the same commands independently (§Architecture in `beta-review.md`) and got matching results on all four ACs. No AC required a second look.
- **Disclosed-not-fixed debt was accepted as disclosed-not-fixed.** Both out-of-scope findings named in α's §Debt were explicitly sanity-checked by β and confirmed genuinely out of scope rather than treated as gaps α should have closed — the disclosure norm held under independent review, not just under α's own self-assessment.

## Canonical-home decision — held up under independent review

α chose `src/packages/cnos.cdd/skills/cdd/CDD.md` over `docs/papers/CELL-OF-CELLS.md` for the new "Mechanism and cognition" section, on three grounds: trivial Tier-1 discoverability (loaded by all three roles already), precedent for domain-flavored content in `CDD.md`'s non-kernel sections, and single-file footprint matching the contract's "(1) exactly one canonical doctrine home."

β re-derived the discoverability argument independently rather than accepting it on α's say-so: it grepped `alpha/SKILL.md`, `gamma/SKILL.md`, and `beta/SKILL.md` load orders itself, confirmed all three name `CDD.md` first, and concluded "the 'already-loaded framework surface' requirement is satisfied by construction" — the same reasoning α gave, reached by an independent path. No disagreement with the placement choice was raised. The decision held up under adversarial re-derivation, not just unchallenged acceptance.

## Debt disclosed this cycle (not fixed, per scope)

Two items, both confirmed by β as genuinely out-of-scope-for-this-cycle rather than fixed or waved through silently:

1. **`CDD.md` line 104 Roles-pointer phrasing** (`β (reviews and merges; merge is β's authority)`) carries the same collapsed-authority pattern the AC3 correction fixed in `beta/SKILL.md` itself, but sits outside the new section and outside the pinned Package-scoping axis for this cycle's `CDD.md` edit. β confirmed touching it would itself have been a Rule-7 violation (unpinned scope expansion) in a cycle whose contract was doctrine-only and tightly scoped.
2. **No loaded skill states the literal `self-coherence.md` header-string contract** `cn cdd verify`'s ledger checker enforces (`## Gap`, not `## §Gap`) — discovered via a red CI run (I6) mid-cycle, fixed in-cycle (commit renaming all six section headers), confirmed green on branch head. β confirmed this as resolved-in-cycle, with the residual ask being a documentation improvement to `alpha/SKILL.md` §2.5 so a future α doesn't rediscover the same gap the same way.

Both items are named for γ triage; disposition is not α's call.

## Pattern note

Same class as `#283`'s skill-class-peer distinction (role-skill peers vs. lifecycle-skill peers) in spirit: this cycle's AC3 audit set (`gamma/SKILL.md`, `alpha/SKILL.md`, `beta/SKILL.md`, `dispatch-protocol/SKILL.md`, `cds-dispatch/SKILL.md`) mixed role skills and mechanism/orchestrator skills in one enumerated list; the one required correction landed in a role skill (`beta/SKILL.md`), and the audit's report format (per-file, edited-or-not-edited-with-reason) surfaced it without incident. No new occurrence of the #283 failure mode this cycle — noted because the audited-file set spans both classes.
