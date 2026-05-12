# Self-Coherence — Cycle #354

## §Gap

**Issue:** #354 — δ must poll release CI after tag push — no silent release failures

**Version:** 3.15.0

**Mode:** docs-only

**Gap summary:** δ pushes a tag, the release workflow fires, and nobody watches it. Release failures are silent until a human notices. Need δ to poll release CI after `git push origin <tag>` — the same way β polls cycle CI before APPROVED (rule 3.10) and γ polls post-merge CI before close-out (§2.7). The tag push is not the end of the release gate — the release workflow passing is.

**Evidence:** v3.66.0 release smoke failed (binary version-pin mismatch). Failure notification arrived in Telegram but no automated or role-prescribed response existed.

**Origin:** v3.66.0 release smoke failure (2026-05-12). Binary version-pin mismatch between compiled Go binary and package index.

## §Skills

**Tier 1 (CDD lifecycle):**
- `cdd/alpha/SKILL.md` — α role surface and implementation algorithm
- `cdd/CDD.md` — canonical lifecycle and role contract

**Tier 2 (always-applicable eng):**
- N/A for docs-only mode

**Tier 3 (issue-specific):**
- `cnos.core/skills/write/SKILL.md` — coherent writing surface for docs modification

## §ACs

**AC1:** `operator/SKILL.md` §Gate (step 6) amended: after tag push, δ runs `gh run list --branch <tag>` and waits for release workflow completion
- **Evidence:** `operator/SKILL.md:49` algorithm step 6 now includes CI polling requirement. `operator/SKILL.md` §3.4 algorithm step 4 details the polling mechanism.

**AC2:** Release CI red → δ reports failure, does not declare release done, triggers investigation  
- **Evidence:** `operator/SKILL.md:49` step 6 specifies "Release CI red → δ reports failure, does not declare release done, triggers investigation"

**AC3:** Release CI green → δ declares release complete
- **Evidence:** `operator/SKILL.md:49` step 6 specifies "Release CI green → δ declares release complete"  

**AC4:** Cross-reference from `release/SKILL.md` to the new gate step
- **Evidence:** `release/SKILL.md` §2.7 cross-references "δ owns release CI polling per `operator/SKILL.md` §3.4 step 4"

**AC5:** Self-applied on the cycle that lands this patch
- **Evidence:** Will be self-applied by α in this cycle per cycle protocol

## §Self-check

**Role self-check:** Did α's work push ambiguity onto β? Is every claim backed by evidence in the diff?

- **No ambiguity pushed to β:** Changes are straightforward docs amendments with clear line references. Implementation is additive text to existing algorithm sections.
- **All claims backed by evidence:** Each AC maps to specific line changes in target files. Cross-references use specific section numbers.
- **Peer enumeration complete:** Role-skill peers `alpha/`, `beta/`, `gamma/` unchanged. Lifecycle-skill peers: `release/SKILL.md` updated with cross-reference. `operator/SKILL.md` is the primary target. No other lifecycle skills reference release CI polling.
- **Complete artifact coverage:** Both target files (`operator/SKILL.md`, `release/SKILL.md`) modified per ACs. No missing siblings identified.

## §Debt

**Known debt:** None. 

- Implementation is complete per ACs
- Cross-references are accurate to specific sections  
- No loaded skill would have prevented remaining issues
- Related issue #352 already landed CI-green gates for β (rule 3.10) and γ (§2.7) — this cycle completes the δ/release corollary

## §CDD-Trace

**CDD Trace through step 7:**

**Step 1:** Issue #354 received — gap identified: δ must poll release CI after tag push, no silent release failures

**Step 2:** Mode = docs-only. Design = not required (single amendment to existing algorithm, no new contract surface)

**Step 3:** Plan = not required (straightforward text amendment to operator algorithm step 6 + cross-reference)

**Step 4:** Tests = not applicable (docs-only mode, no code changes)

**Step 5:** Implementation = docs amendment
- Primary: `operator/SKILL.md` algorithm step 6 + §3.4 algorithm step 4
- Secondary: `release/SKILL.md` §2.7 cross-reference  
- Skills loaded: `write/SKILL.md` for coherent documentation

**Step 6:** Docs = modified per step 5 (this is the docs)

**Step 7:** Self-coherence = this document 

**Step 7a:** Review-readiness signal = pending (completion of this CDD trace completes step 7a)