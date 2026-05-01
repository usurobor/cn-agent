## Gap

**Issue:** #325 — cdd: systematic lifecycle audit and refactor for role, artifact, and release coherence
**Version/mode:** MCA — lifecycle coherence patch (spec/skill refactor; no runtime implementation)
**Priority:** P1

CDD defines the development lifecycle for substantial cnos work. Recent cycles exposed gaps where CDD itself was ambiguous or internally inconsistent:

- The coordination model (sequential bounded dispatch vs. persistent polling) is not declared, so each cycle must infer it
- α close-out timing is broken in bounded dispatch: α is told to write close-out "after β merge" but α has already exited
- No role/artifact ownership matrix exists: a mechanical lookup by artifact yields no single source for owner + verifier + preconditions + failure consequences
- γ closure gate does not make `gamma-closeout.md` itself a precondition for δ tag/release
- Small-change path does not state which close-outs and release gates collapse
- No failure-mode regression surface exists to catch the known patterns (stale unreleased dirs, missing close-outs, missing RELEASE.md, δ tag before γ closure)

The work is entirely spec/skill: edits to CDD.md and all role skills. No runtime implementation.

Design artifact: `.cdd/unreleased/325/alpha-design.md`

---

## Skills

**Tier 1a (CDD authority):**
- `CDD.md` — canonical lifecycle and role contract
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — α role surface

**Tier 1b (lifecycle sub-skills):**
- `src/packages/cnos.cdd/skills/cdd/design/` — not loaded separately; design principles applied via Tier 3 `cnos.core/skills/design`

**Tier 2 (general engineering):**
- Writing bundle applies: all outputs are written artifacts

**Tier 3 (issue-specific):**
- `src/packages/cnos.core/skills/write/SKILL.md` — all CDD edits are written artifacts; write skill governs prose quality, one-governing-question-per-doc discipline, stable-facts-once rule
- `src/packages/cnos.core/skills/design/SKILL.md` — lifecycle refactor is an authority-boundary and state-machine design problem; design skill governs: one reason to change per boundary, policy above detail, interfaces truthful, no premature canonicalization
- `src/packages/cnos.eng/skills/eng/test/SKILL.md` — regression fixtures and lifecycle checklists must prove invariants, not just examples; negative space is mandatory; test the failure-mode regression surface

---

## ACs

### AC1: CDD lifecycle model is explicit

**Status: MET**

Evidence: `CDD.md` §1.6 "Coordination model" added (commit `e91e72af`). States explicitly:
- Current model: sequential bounded role invocations via `claude -p`
- Each role exits when phase is complete
- In-session polling applies during each role's active session window (not background daemon)
- Positive test: reader can answer "α exits after signaling review-readiness; re-dispatch is used for fix rounds and close-out"
- Negative test: no longer possible to read CDD and conclude α stays alive indefinitely waiting for β approval

`CDD.md` §Tracking updated with: "Polling — in-session only. Polling applies during each role's active session. A role that has completed its phase and exited is not polling."

### AC2: Role/artifact ownership matrix exists

**Status: MET**

Evidence: `CDD.md` §5.3b "Role/artifact ownership matrix" added (commit `e91e72af`). Contains:
- All 8 required artifact rows (self-coherence.md, beta-review.md, alpha-closeout.md, beta-closeout.md, gamma-closeout.md, RELEASE.md, .cdd/releases/{X.Y.Z}/{N}/, POST-RELEASE-ASSESSMENT.md)
- Columns: Artifact | Owner | Written when | Verified by | Required before | Missing means
- gamma-closeout.md row explicitly states: "Missing means: δ must not tag; γ has not declared closure"
- Cross-referenced from lifecycle state table (§4.1a)

### AC3: α close-out timing is executable

**Status: MET**

Evidence:
- `CDD.md` §1.4 α algorithm step 10 rewritten (commit `e91e72af`): specifies re-dispatch mechanism and provisional fallback
- `CDD.md` §1.6a: re-dispatch prompt formats added (α close-out re-dispatch, α fix-round re-dispatch)
- `alpha/SKILL.md` §2.8 rewritten (commit `41fede0f`): re-dispatch path + provisional fallback both documented
- `gamma/SKILL.md` §2.7 updated (commit `40b71d67`): "γ must explicitly request this re-dispatch"
- `operator/SKILL.md` algorithm step 5 added (commit `d7324165`): "Re-dispatch α for close-out — mandatory when γ requests it"
- `CDD.md` §4.1a lifecycle state table (commit `e91e72af`): S8 state "α close-out" — owner: α (re-dispatched), failure: "Re-dispatch unavailable → provisional close-out at review-readiness (declare as debt)"

No cycle can reach γ closure with missing alpha-closeout.md: §5.3b row says "Missing means: γ closure gate blocks; γ requests δ to re-dispatch α"

### AC4: β close-out and merge boundary are executable

**Status: MET**

Evidence:
- `beta/SKILL.md` already correctly states β owns merge and β close-out; δ owns tag/release
- `release/SKILL.md` Core Principle corrected (commit `b9a107d6`): "β owns: review approval outcome, git merge into main, and β close-out. δ owns tag/release/deploy"
- `post-release/SKILL.md` §Who corrected (commit `b9a107d6`): "β owns git merge into main and the β close-out. δ owns tag/release/deploy"
- All three surfaces now agree: β merges, δ tags

### AC5: γ closure gate is complete

**Status: MET**

Evidence:
- `gamma/SKILL.md` §2.10 has 12-row closure gate (pre-existing)
- Row 1: alpha-closeout.md; Row 2: beta-closeout.md; Row 11: RELEASE.md; Row 12: cycle dir move
- `gamma/SKILL.md` §2.10 "Then:" updated (commit `40b71d67`): gamma-closeout.md explicitly labeled as "closure declaration artifact. δ must not tag/release until gamma-closeout.md exists on main."
- `CDD.md` §5.3b ownership matrix states gamma-closeout.md "Required before: δ tag/release; Missing means: δ must not tag; γ has not declared closure"

### AC6: δ release-boundary preflight is placed correctly

**Status: MET (pre-existing + reinforced)**

Evidence:
- `CDD.md` §1.4 γ algorithm Phase 5a: preflight after β merge + close-outs, before γ closure declaration
- `CDD.md` §4.1a state table (commit `e91e72af`): S10 (δ preflight) requires S7 (β merged) and S9 (γ triaging) complete; precedes S11 (γ closing) and S12 (δ disconnect)
- `operator/SKILL.md` lifecycle table updated (commit `d7324165`): "δ preflight" row placed between release prep and closure

### AC7: `.cdd/unreleased` movement is owned and timed

**Status: MET (pre-existing + reinforced)**

Evidence:
- `gamma/SKILL.md` §2.6: γ owns the move before requesting δ tag (pre-existing)
- `gamma/SKILL.md` §2.10 row 12: closure gate blocks until moved (pre-existing)
- `CDD.md` §8.1 F4 checklist row (commit `e91e72af`): "Stale .cdd/unreleased/{N}/ after release" — explicit negative test
- `CDD.md` §5.3b ownership matrix: `.cdd/releases/{X.Y.Z}/{N}/` row — "Required before: δ tag/release; Missing means: Stale unreleased dir; γ closure gate blocks until moved"

### AC8: RELEASE.md ownership and timing are explicit

**Status: MET (pre-existing + reinforced)**

Evidence:
- `gamma/SKILL.md` §2.6 already: "Write RELEASE.md — per release/SKILL.md §2.5" (pre-existing)
- `release/SKILL.md` §3.7: "RELEASE.md must exist before tag" (pre-existing)
- `CDD.md` §5.3b: RELEASE.md "Owner: γ; Required before: δ tag/release; Missing means: δ must not tag"
- `CDD.md` §8.1 F5 checklist row: "Missing RELEASE.md before tag" — explicit negative test
- `CDD.md` §1.2 small-change table: "RELEASE.md: Required before tag/release — no exception"

### AC9: Polling-era text is reconciled

**Status: MET**

Evidence:
- `CDD.md` §1.6 declares current model as sequential bounded dispatch (commit `e91e72af`)
- `CDD.md` §1.6 explicitly states: "Not current: persistent polling role sessions where α, β, γ remain alive indefinitely between phases"
- `CDD.md` §Tracking updated: "Polling — in-session only. Polling applies during each role's active session." (commit `e91e72af`)
- All polling mechanics in §Tracking remain valid as in-session polling descriptions
- `operator/SKILL.md` §1 describes sequential dispatch with re-dispatch for each additional phase
- No conflict remains between polling language and sequential dispatch language

### AC10: Role skill peer audit is complete

**Status: MET**

Peer set — role-skill peers:
- `alpha/SKILL.md`: Updated §2.8 close-out timing ✅
- `beta/SKILL.md`: Confirmed agreement — β merge authority, no tag ownership ✅ (no change needed)
- `gamma/SKILL.md`: Updated §2.7 (α close-out re-dispatch) and §2.10 (gamma-closeout.md as closure declaration) ✅
- `operator/SKILL.md`: Updated Algorithm (steps 4-6) and lifecycle table ✅

Peer set — lifecycle-skill peers:
- `release/SKILL.md`: Fixed Core Principle (β/δ authority split) ✅
- `post-release/SKILL.md`: Fixed §Who (β/δ authority split) ✅
- `review/SKILL.md`: Confirmed no changes needed (verdict/merge section already correct) ✅

All peers either updated or explicitly confirmed no change needed.

### AC11: CDD has a failure-mode regression surface

**Status: MET**

Evidence: `CDD.md` §8.1 "Closure verification checklist" (commit `e91e72af`) — 10 rows covering:
- F1: Missing alpha-closeout.md (positive + negative test defined)
- F2: Missing beta-closeout.md
- F3: Missing gamma-closeout.md before δ tag
- F4: Stale .cdd/unreleased/{N}/ after release
- F5: Missing RELEASE.md before tag
- F6: δ tag before γ closure declaration
- F7: δ tag before δ preflight
- F8: α close-out relies on already-exited α
- F9: Polling-required role exits early
- F10: PRA missing after release

Positive test: "A substantial cycle that reaches S12 (δ disconnect) with all 10 rows above passing has closed correctly."
Negative test: "Any row that fails when γ reaches state S11 (γ closing) must block closure."

This is a reviewable checklist embedded in CDD — inspectable and runnable per-cycle without separate tooling.

### AC12: Small-change path remains coherent

**Status: MET**

Evidence: `CDD.md` §1.2 small-change section expanded (commit `e91e72af`) with explicit artifact collapse table. For each artifact:
- self-coherence.md: Optional
- alpha-closeout.md: Optional (one-liner if no findings)
- beta-closeout.md: Not applicable
- gamma-closeout.md: Not applicable unless γ coordinated
- beta-review.md: Not applicable
- RELEASE.md: **Required — no exception**
- .cdd/releases/{X.Y.Z}/{N}/ move: **Required**
- PRA: **Required after every release**
- Version directory / bootstrap stubs: Not required

The table explicitly states: "A small-change cycle may not use the small-change path to avoid RELEASE.md, cycle-directory movement, or the post-release assessment."
