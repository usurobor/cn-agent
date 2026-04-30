# β Close-out — #307

skill(cdd/issue): move issue katas to cnos.cdd.kata package

---

## Review context

**Dispatch mode:** identity-rotation (single agent; roles committed under separate git identities).

**R1 (β R1):** `b2cc9f9b` `beta@cdd.cnos` — REQUEST CHANGES; 1 finding (F1).
**α fix-round:** `e71831dc` `alpha@cdd.cnos` — KATAS.md M5 row + M0–M5 reference update.
**R2 (β R2):** `8e465c1f` `beta@cdd.cnos` — APPROVED; F1 resolved; no new findings.
**Merge commit:** `55cd3c20` on `main` — `Closes #307`.

**origin/main SHA at each review pass:**
- R1 base: `396d9982af2c2318400fb5c9eb400178d488f99b` (fetched synchronously before R1)
- R2 base: `396d9982af2c2318400fb5c9eb400178d488f99b` (fetched synchronously before R2; no new main commits between R1 and R2)

---

## Narrowing pattern

**R1 → R2:** One finding (F1), B-severity, judgment/contract type. The implementation was coherent in R1; the single gap was a stale secondary index (KATAS.md) with an unnamed debt omission. F1 was addressed in one commit; R2 was narrowly scoped to the fix. No D or C findings across either round.

---

## Merge evidence

- **Merge commit SHA:** `55cd3c20e56aa20f52311b593d0b58775f3bcfea`
- **Branch merged:** `cycle/307` (head `8e465c1f` at merge time)
- **Target:** `main` (was `396d9982` before merge)
- **Method:** `git merge --no-ff cycle/307 -m "Closes #307: ..."` — real merge commit, not fast-forward; issue auto-close via `Closes #307` in merge commit message.
- **Post-merge main state:** clean; all cycle branch commits present on main.

---

## β-side findings

**F1 — KATAS.md stale catalog (B, judgment/contract):** `docs/gamma/cdd/KATAS.md` Tier 3 table listed M0–M4 only; M5-issue-authoring bundle (created in this cycle) was not added. The `self-coherence.md` §Debt section stated "No known α debt new this cycle" without naming the gap — omission was silent rather than explicit. Pattern: when a new M-series bundle is created, the KATAS.md catalog requires a corresponding row addition; this was missed at α stage and caught at R1. The precedent cycle (#304) did not add a new M-series row because M2 pre-existed; the distinction (existing vs. new bundle) was not part of α's pre-review gate checklist at the time.

**Cycle structure observation:** The issue/SKILL.md §5 body was ~35 lines of embedded kata content. Post-move, issue/SKILL.md carries only a pointer and the kata content is discoverable from `cn kata list`. The separation is now consistent with review/SKILL.md (moved in #304). The pattern holds: eight remaining cdd lifecycle skills carry `kata_surface: embedded` or `none`; none of those were in scope for this cycle per the issue §Non-goals.

---

## CI state

Branch CI not configured on cycle branches (per self-coherence.md §Review-readiness; CI gates on main push). CI job I5 (`skill-frontmatter-check` via `tools/validate-skill-frontmatter.sh` + `cue vet schemas/skill.cue`) will run on the post-merge main push. The `kata_ref` key passes through the open-schema trailing `...` (same key shape as `review/SKILL.md`, which already ships with `kata_ref`).

---

## Handoff to γ

Merge committed. β close-out written. α close-out pending (α's responsibility per CDD.md §1.4 α step 10). γ writes the post-release assessment after both close-outs exist on main.
