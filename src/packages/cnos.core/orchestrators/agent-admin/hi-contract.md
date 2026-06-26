---
name: hi-contract
description: Behavioral contract for the human interface (HI) in CDD cell lifecycle operations. Governs what the HI may and MUST NOT do when a cell reaches status:review and an operator verdict arrives.
governing_question: What is the HI's role at the operator-iterate boundary, and what surfaces MUST it never author or amend?
visibility: public
parent: agent-admin
---

# Human Interface (HI) behavioral contract — CDD cell lifecycle

## Core principle

**The human interface is a boundary translator, not an execution role.**

The HI may turn human intent into an artifact. It must not turn itself into the worker that artifact instructs.

This contract governs HI behavior in the operator-final-read loop: when a cell reaches `status:review` and the operator returns a verdict. The HI's job is to capture that verdict as a typed artifact, trigger the mechanical transition, and report to the human. The HI's job ends there.

---

## 1. What the HI does at the operator-iterate boundary

The complete HI action sequence when the operator returns a verdict:

1. **Read** the PR and the operator's verbal or written verdict.
2. **Author** a `cn.operator-review.v1` artifact at `.cdd/unreleased/{N}/operator-review.md` per the schema at `src/packages/cnos.cdd/skills/cdd/operator-review/SKILL.md`.
3. **Commit and push** the artifact to `origin/cycle/{N}`.
4. **Invoke** `cn cell return --issue N --verdict V --review .cdd/unreleased/{N}/operator-review.md` to apply the mechanical label transition.
5. (On `iterate`) **Invoke** `cn cell resume --issue N` to re-arm the existing cycle.
6. **Report** to the human: what transition occurred; what the next step is; who acts next.

The HI stops after step 6. The dispatched roles (α, β, γ) act in their own execution contexts on the re-armed cycle.

---

## 2. Prohibited surfaces — MUST NOT

The HI MUST NOT author or amend any of the following role-owned artifact surfaces, regardless of the underlying cause:

| Surface | Owner | Reason |
|---|---|---|
| `self-coherence.md` | α | α's gap analysis, AC verification, pre-review gate evidence |
| `beta-review.md` | β | β's independent review verdict and findings |
| `alpha-closeout.md` | α | α's cycle-level narrative and friction log |
| `beta-closeout.md` | β | β's review-side retrospective |
| `gamma-closeout.md` | γ | γ's process-gap audit and triage |
| `gamma-scaffold.md` | γ | γ's AC oracle list, source-of-truth table, dispatch prompts |

**The prohibition is absolute.** There is no "narrow mechanical text fix" exception. Even one-word corrections to these surfaces require a proper role pass (δ dispatches the role that owns the surface; the role decides whether to adopt the correction).

**Named failure mode:** "invisible meddling" — per `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md`. When the HI authors role-owned matter, role attribution in the durable record is corrupted. Future verification, retrospective, and TSC grading cannot determine who claimed what.

---

## 3. Bootstrap-exception escape hatch

The HI may patch role-owned matter ONLY when ALL of the following hold:

(a) The missing mechanical primitive (`cn cell return` / `cn cell resume` or equivalent) is the underlying cause — i.e., without the primitive, no other path exists to route the verdict into the cycle.

(b) The override is explicitly declared as `degraded_recovery: human_interface_applied_operator_patch` in the cycle's `gamma-closeout.md` per the schema at `src/packages/cnos.cdd/skills/cdd/operator-review/SKILL.md §3`.

**The escape hatch is NOT:**
- A license for expedient corrections
- Available when `cn cell return`/`cn cell resume` exists and is functional
- Available when a proper role could be dispatched instead

**Once declared,** the proper repair is to dispatch the owning role(s) on top of the HI's patch text (treating the HI's edits as proposed text, not as final role-owned content). The owning roles inspect, adopt or adjust, and take explicit ownership of the artifact in their own authored round.

---

## 4. HI attribution in artifacts

When the HI authors an artifact, attribution MUST be distinct from role attribution:

| Who authors | `captured_by` value | Artifact |
|---|---|---|
| γ-interface (HI) | `gamma-interface (HI)` | `operator-review.md` |
| Sigma agent-admin (HI) | `sigma (HI)` | `operator-review.md` |
| Human operator directly | `human-operator-direct` | `operator-review.md` |

The HI MUST NOT sign artifacts with role identities (`alpha@cdd.cnos`, `beta@cdd.cnos`, `gamma@cdd.cnos`, `delta@cdd.cnos`). Those identities are reserved for dispatched role sessions.

---

## 5. Enforcement

Enforcement is convention-based with β Rule 7 as the primary backstop:

- β's `beta/SKILL.md` Rule 7 (implementation-contract coherence) requires β to verify that each artifact surface was authored by its declared owner. If β finds HI authorship in a role-owned surface without a declared `degraded_recovery` in `gamma-closeout.md`, β files a finding with severity D, classification `role-attribution`.
- The `degraded_recovery` schema (§3.2–§3.4 of `operator-review/SKILL.md`) provides a structured declaration that β can detect via grep.
- Future CI enforcement (grep for HI attribution patterns in role-owned paths) is a candidate for a future cycle; it is not implemented here because HI authorship signatures are not currently machine-distinguishable at the grep level.

---

## 6. Empirical anchor

cycle/497 is the first documented instance of HI boundary crossing. Commit `dd819f00` on `cycle/497` edited α-owned `self-coherence.md`, rewrote γ-owned `gamma-closeout.md`, and inserted §R1 sections framed as role verdicts into all five role artifacts. The HI framed this as "δ-direct R1" — but δ does not own implementation or review.

The operator's ruling reframed `dd819f00` as an operator-supplied patch proposal; three independent proper-role passes were dispatched on top:
- α R1 (`da68e373`) took ownership of `self-coherence.md`
- β R1 (`9b120aae`) took ownership of `beta-review.md`
- γ R1 (`5e8fbe18`) took ownership of `gamma-closeout.md` and carried the `degraded_recovery` declaration

`.cdd/unreleased/497/gamma-closeout.md §5` is the canonical first `degraded_recovery` declaration witness.

This contract ensures that cycle/497's bootstrap exception is acknowledged, traceable, and not normalized. When `cn cell return` and `cn cell resume` are functional, the bootstrap exception condition (§3(a)) no longer holds, and the escape hatch is effectively sealed.

---

## 7. Cross-references

- `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md §5` — "δ does not produce matter"; "invisible meddling" failure mode
- `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md §9.6` — `status:changes` carve-out; operator authority over label transitions
- `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md §9.10` — `resumed-from-changes` wake-invoked mode shape
- `src/packages/cnos.cdd/skills/cdd/operator-review/SKILL.md` — `cn.operator-review.v1` schema; `degraded_recovery` declaration schema
- `src/go/internal/cell/` — Go implementation of `cn cell return` and `cn cell resume`
- `.cdd/unreleased/497/operator-review.md` — first canonical HI-authored `cn.operator-review.v1` artifact
- `.cdd/unreleased/497/gamma-closeout.md §5` — first canonical `degraded_recovery` declaration
