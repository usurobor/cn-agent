---
cycle: 343
issue: "cdd: Canonical git identity convention for cdd role actors ({role}@{project}.cdd.cnos)"
branch: cycle/343
mode: design-and-build
status: review-ready
---

# Self-Coherence — Cycle #343

## §Gap

**Issue:** #343 — The cdd identity convention `{role}@cdd.{project}` inverts the DNS namespace hierarchy. `cdd` is a protocol hosted by cnos, not a subdomain of the tenant project. The correct three-level form is `{role}@{project}.cdd.cnos`.

**Mode:** `design-and-build` — design is fixed in the issue body; α implements the patch in this cycle.

## §Skills

**Tier 1:**
- `alpha/SKILL.md` (this role surface)
- `CDD.md` (canonical lifecycle and role contract)

**Tier 2:** None — docs-only change, no code engineering skills required.

**Tier 3:**
- `src/packages/cnos.core/skills/write/SKILL.md` — every α output is a written artifact; applied as authoring constraint throughout.

## §ACs

### AC1 — Three-level identity convention named

**Claim:** `operator/SKILL.md` §Git identity for role actors prescribes `{role}@{project}.cdd.cnos` with a 2–4 sentence rationale citing DNS hierarchy and cnos-as-protocol-origin.

**Evidence:**
- `operator/SKILL.md` lines 53–75: new `## Git identity for role actors` section added.
- Rationale text: "DNS domains read broad-to-narrow right-to-left: `cnos` is the origin repository where the cdd protocol is defined and versioned, `cdd` is the protocol namespace inside cnos, and `{project}` is the tenant project running the protocol. The role name is the local part. This form makes the protocol's origin repo visible in every commit trailer and leaves namespace room for sibling protocols (`cnav`, `cnobs`) under the same cnos root." — 2 sentences covering DNS hierarchy + cnos-as-origin.
- `review/SKILL.md` §Review identity: replaced cycle #287 doctrine with cross-reference to `operator/SKILL.md` §Git identity for role actors. Old literal `beta@cdd.{project}` removed.
- `CDD.md` γ-algorithm step 1, Parameters block, α-algorithm step 2, β-algorithm step 2: all four identity setup lines updated to `{role}@{project}.cdd.cnos` with cnos-elision note and cross-reference.
- `alpha/SKILL.md` pre-review gate row 14: updated `alpha@cdd.{project}` → `alpha@{project}.cdd.cnos`.

**Oracle check:**
- `rg 'beta@cdd\.{project}' cdd/review/SKILL.md` → 0 hits ✓
- `rg '\.cdd\.cnos' cdd/operator/SKILL.md` → prescription site found ✓

### AC2 — Special case for cnos resolved

**Claim:** The elision form `{role}@cdd.cnos` is chosen as canonical for cnos-side actors, with rationale documented.

**Evidence:**
- `operator/SKILL.md` §Git identity: "When the project running the cycle is the cnos repo, the literal form would be `{role}@cnos.cdd.cnos` (redundant `cnos`). The canonical elision is `{role}@cdd.cnos`, which reads as 'the cdd protocol at cnos.' Existing cnos commit trailers already use this form; the redundancy adds no information."
- Worked example table row 2: `beta | cnos | beta@cdd.cnos | cnos actor — elision form`.
- One form chosen (not both); no fallback to "either is fine."

**Oracle check:**
- Worked example in prescription site shows cnos-side identity explicitly ✓
- No ambiguity between the two candidate forms ✓

### AC3 — Worked example table present

**Claim:** A 4-row table in `operator/SKILL.md` §Git identity covers required pairings including deprecated form.

**Evidence — table rows:**
| Row | Content | Satisfies |
|-----|---------|-----------|
| 1 | `alpha@tsc.cdd.cnos` — tsc project actor | AC3 row 1 |
| 2 | `beta@cdd.cnos` — cnos actor (elision) | AC3 row 2 (AC2 choice) |
| 3 | `gamma@acme.cdd.cnos` — hypothetical third project | AC3 row 3 |
| 4 | `beta@cdd.{project}` — **(deprecated)** cycle #287 form | AC3 deprecated row |

- Table renders in valid markdown ✓
- Table is reachable from CDD.md: CDD.md §Parameters cross-references `operator/SKILL.md §Git identity for role actors` ✓
- No contradictory examples elsewhere in `cdd/` (all remaining `@cdd.` hits are either canonical elision form `@cdd.cnos` or the deprecated form inside migration/history blocks) ✓

### AC4 — Migration paragraph in post-release

**Claim:** `cdd/post-release/SKILL.md` gains an "Identity migration" subsection.

**Evidence:**
- `post-release/SKILL.md` §Identity migration added (lines before Kata section).
- Content: names cycle #343 as cutover, date 2026-05-11, history immutable, forward-only switch, transition-window tolerance.
- Word count: 77 words ≤ 80 ✓

**Oracle check:**
- `post-release/SKILL.md` contains "## Identity migration" subsection ✓
- No expectation of rewriting historical commits ✓

### AC5 — Patch-landing cycle uses new form

**Claim:** All α commits on cycle/343 use `alpha@cdd.cnos` (cnos elision form per AC2).

**Evidence:**
- `git config user.email` → `alpha@cdd.cnos` (set at session start)
- `git log --format='%ae' cycle/343` → γ commits: `gamma@cdd.cnos`; α implementation commit: `alpha@cdd.cnos`
- The implementation commit `0757c9be` was authored as `alpha@cdd.cnos` ✓

**Oracle:** `git log --format='%ae' cycle/343` shows `alpha@cdd.cnos` and `gamma@cdd.cnos` — all cnos elision form ✓

## §Self-check

- Did α push ambiguity onto β? No. Every AC is mapped to concrete evidence in the diff. The peer set was fully enumerated (6 files with `@cdd.` references; all accounted for). The one ambiguity — whether elision form `@cdd.cnos` triggers the AC1 negative oracle — was resolved by confirming the oracle targets the deprecated `beta@cdd.{project}` form specifically, not the new canonical elision.
- Is every claim backed by evidence? Yes — each AC section names the specific file and line range.
- Did α identify all peer surfaces? Yes. Enumerated: `review/SKILL.md`, `alpha/SKILL.md`, `beta/SKILL.md`, `operator/SKILL.md`, `CDD.md`, `issue/contract/SKILL.md`. Updated: first 5. Exempted: `issue/contract/SKILL.md` (uses `alpha@cdd.cnos` which is already the canonical cnos elision form — no change needed).

## §Debt

None known. `beta/SKILL.md` pre-merge gate row 1 retains `beta@cdd.cnos` (already the canonical cnos elision form); no update needed.

## §CDD-Trace

| Step | Action | Skill | Decision |
|------|--------|-------|----------|
| 1 | Cycle branch created | gamma/SKILL.md §2.5a | cycle/343 from origin/main |
| 2 | Scaffold created | gamma/SKILL.md §2.5a | .cdd/unreleased/343/ |
| 3 | α prompt written | gamma/SKILL.md §2.5b | dispatched to δ |
| 4 | β prompt written | gamma/SKILL.md §2.5b | dispatched to δ |
| 5 | α loaded skills | alpha/SKILL.md §2.1 | Tier 1: CDD.md + alpha/SKILL.md; Tier 3: write/SKILL.md |
| 6 | α implemented AC1–AC5 | alpha/SKILL.md §2.2 | Updated: operator/SKILL.md (new §Git identity section), review/SKILL.md (§Review identity cross-ref), alpha/SKILL.md (pre-review gate row 14), CDD.md (4 identity lines), post-release/SKILL.md (§Identity migration). Callers: CDD.md §Parameters → operator/SKILL.md (cross-ref added) |
| 7 | α wrote self-coherence | alpha/SKILL.md §2.5 | AC-by-AC mapping with evidence; peer enumeration complete; no debt |

## Review-readiness signal

status: review-ready

α R1 — implementation SHA: `0757c9be` — branch CI: docs-only, no executable tests. All 5 ACs have evidence. Peer set fully enumerated. No known debt.
