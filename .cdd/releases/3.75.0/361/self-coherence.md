---
cycle: 361
issue: "https://github.com/usurobor/cnos/issues/361"
role: alpha
sections:
  planned: [Gap, Skills, ACs, Self-check, Debt, CDD-Trace, Review-readiness]
  completed: [Gap, Skills, ACs, Self-check, Debt, CDD-Trace]
---

# α Self-coherence — #361

## §Gap

**Issue**: cdd: review rules 3.3/3.4 — verdict-shape lint (no conditional APPROVEDs)

β@S4 in tsc #53 issued "APPROVED with 3 unresolved C findings + conditional language." Rules 3.3 (APPROVED = zero unresolved findings) and 3.4 (verdict before details) both forbid this shape — but the ban is implicit. No explicit lint enumerates invalid verdict shapes, names the conditional qualifiers that smuggle in unresolved findings, or rules out split verdicts.

**Mode**: skill-patch, single-file diff to `src/packages/cnos.cdd/skills/cdd/review/SKILL.md`.

## §Skills

- **Tier 1** — `CDD.md`, `cdd/alpha/SKILL.md` (this role surface). No lifecycle sub-skill needed (design/plan marked not-required: single-file skill patch, no impact graph beyond the file itself).
- **Tier 2** — none. The diff is normative-prose-only inside one role surface; no Tier 2 engineering bundle applies.
- **Tier 3** — `cnos.core/skills/write/SKILL.md` (per dispatch). Governs prose: front-load the point (3.4), one fact one home (3.3), specifics over abstractions (3.8), state what it is (3.1).

## §ACs

| # | AC | Status | Evidence |
|---|----|--------|----------|
| AC1 | review/SKILL.md §3.3/3.4 contains explicit verdict-shape lint rules: `APPROVED` + unresolved findings = invalid; `APPROVED` + "conditional/pending/modulo" qualifier = invalid | met | review/SKILL.md §3.4a bullet 1 (`APPROVED` + unresolved findings at any severity, cross-refs 3.3); §3.4a bullet 2 enumerates the six qualifier tokens `conditional`, `pending`, `modulo`, `subject to`, `assuming`, `provisional on`, `with follow-up` |
| AC2 | Recovery path documented: conditional becomes RC with conditions as required-fix findings | met | review/SKILL.md §3.4a Recovery paragraph: "any conditional or split shape becomes `REQUEST CHANGES`; each named condition is reformatted as a required-fix finding at its severity. The reviewer re-emits a single terminal verdict in the next round once findings clear." |
| AC3 | Reviewers must select one terminal verdict per round — no split verdicts | met | review/SKILL.md §3.4a bullet 3 (Split verdicts): "two terminal verdicts in one round … One round, one decision." Plus checklist line: "Verdict-shape lint passed (no `APPROVED` + unresolved findings; no conditional qualifier; no split verdict) per rule 3.4a." |

## §Self-check

**Did α push ambiguity onto β?** No. 3.4a names the three invalid shapes by exact form (`APPROVED` + unresolved findings; `APPROVED` + qualifier from a closed list of seven tokens; two terminal verdicts in one round) and names one recovery path (RC + conditions reformatted as findings + one terminal verdict next round). The checklist line gives β a single mechanical check.

**Is every claim backed by evidence in the diff?**

- "APPROVED + unresolved findings is invalid" — backed by 3.3 (`zero findings at any severity remain unresolved`) which 3.4a cross-references. No new claim, just made explicit.
- "Conditional qualifiers smuggle in unresolved findings" — backed by 3.3's prohibition: a `pending X` clause names X as unresolved, which 3.3 already forbids.
- "Split verdicts are invalid" — new constraint, derived from rule 3.4 ("Lead with APPROVED / REQUEST CHANGES" — singular, one decision).
- Derivation line cites tsc #53 β@S4 as the surfacing incident.

**Peer enumeration (per α §2.3).** The diff modifies one lifecycle skill (`review/`). Peer surfaces searched:

| Surface | Class | Verdict-shape statement? | Action |
|---------|-------|--------------------------|--------|
| `review/SKILL.md` | lifecycle (the target) | yes — 3.3, 3.4, 3.4a (new) | updated |
| `beta/SKILL.md` | role | references "review verdict (RC or A)" structurally; delegates shape to `review/SKILL.md` | no update — authority is at `review/` |
| `CDD.md` §S6 | canonical | names verdict at the lifecycle level ("RC / A"); does not constrain shape | no update — shape lives in `review/SKILL.md` |
| `release/`, `post-release/`, `alpha/`, `activation/` | mixed | use `provisional`/`pending` in *release-scoring* and *close-out fallback* contexts, not verdict-shape | no update — different referent (grep confirmed: `release/SKILL.md:100` = CHANGELOG TSC row, `alpha/SKILL.md:296` = closeout fallback) |

Single-home rule satisfied: verdict-shape lint is stated once in `review/SKILL.md` §3.4a; siblings continue to reference review/SKILL.md as authority.

**Intra-doc repetition check (α §2.3).** Within `review/SKILL.md`, verdict-shape facts appear in 3.3 (APPROVED+unresolved invalid), 3.4 (verdict-before-details), 3.4a (lint enumeration), and the checklist. 3.4a explicitly cross-references 3.3 ("contradicts 3.3"); the checklist line cites "rule 3.4a." No drift between sites.

**Write skill compliance.** §3.4a leads with the point ("Three invalid verdict shapes — each auto-RC"), uses specifics (six named qualifier tokens, not "such as"), gives one Recovery paragraph (one fact, one home), cites the deriving incident.

## §Debt

None. The diff is a single-file normative-prose patch with no new module, no schema change, no harness, no test surface, no language tooling. AC1, AC2, AC3 are all met with text-only evidence in the same file. The provisional alpha-closeout fallback (α §2.8) is not invoked: this cycle uses the standard re-dispatch path.

Pre-#283 form `.cdd/releases/{X.Y.Z}/alpha/CLOSE-OUT.md` is warn-only per α §2.8 and not exercised here.

## §CDD-Trace

| Step | Artifact | This cycle |
|------|----------|------------|
| 1 | design | not required — single-file normative-prose patch, no impact graph beyond the file |
| 2 | coherence contract | §Gap above (issue #361; tsc #53 β@S4 surfaced "APPROVED with 3 unresolved C findings + conditional language") |
| 3 | plan | not required — one rule block + one checklist line; no sequencing |
| 4 | tests | not required — normative prose. β consuming 3.4a + the checklist line during verdict emission is the enforcement path |
| 5 | code | not required — no executable code in the diff |
| 6 | docs | `src/packages/cnos.cdd/skills/cdd/review/SKILL.md` — added rule 3.4a (Verdict-shape lint, 11 lines) and checklist entry (1 line); +13 lines. `.cdd/unreleased/361/gamma-scaffold.md` — γ scaffold authored before α dispatch. `.cdd/unreleased/361/self-coherence.md` — this file |
| 7 | self-coherence | this file |

**Artifact enumeration (α §2.6 row 11).** `git diff --stat origin/main..HEAD` lists three files; all three are named at step 6 above.

**Caller-path trace (α §2.6 row 12).** N/A — no new module or function. The rule is consumed by β reading `review/SKILL.md` during verdict emission; the checklist entry added in the same diff is the call site.

**Test assertion count (α §2.6 row 13).** N/A — no test surface.

## Review-readiness

**Ready for β.** Implementation complete at `34aa2b90`. All 3 ACs addressed. Self-coherence sections complete. Pre-review gate rows pass:

- Base SHA `56202534` == `origin/main` at dispatch ✅
- All commits on `cycle/361`, authored as `alpha@cdd.cnos` ✅
- `.cdd/unreleased/361/self-coherence.md` present with all sections ✅
- Implementation diff: new §3.4a verdict-shape lint in `review/SKILL.md` ✅

Note: review-readiness written by δ (operator override) after α session timeout. α completed all substantive work; only this signal section was missing.
