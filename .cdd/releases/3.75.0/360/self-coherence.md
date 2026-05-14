---
cycle: 360
role: alpha
issue: "https://github.com/usurobor/cnos/issues/360"
date: "2026-05-14"
base_sha: "c77f34a4"
sections:
  planned: [Gap, Skills, ACs, Self-check, Debt, CDD-Trace, Review-readiness]
  completed: [Gap, Skills, ACs, Self-check, Debt, CDD-Trace]
---

# α Self-coherence — #360

## §Gap

**Issue**: #360 — `cdd: review rule 3.11b — clarify exemption scope (sub-issue body required)`.

**Mode**: skill-patch on `review/SKILL.md` rule 3.11b. Source: tsc v0.10.0 wave `cdd-iteration` F2 (`.cdd/releases/0.10.0/49/cdd-iteration.md`).

**Incoherence**: rule 3.11b §Scope says exemptions are "documented in the issue" without naming *which* issue body counts. In tsc #49 a cycle had its exemption claim land as a comment on the parent/master issue; four independently-dispatched β subagents diverged — β@S1 read 3.11b literally and emitted D-severity RC, β@S2/S3/S4 accepted the master-issue comment as exemption and treated the protocol miss as a B non-blocker. Divergence was structural (the rule under-specified exemption discoverability), not careless. Until 3.11b names *which* issue body grants exemption, every cycle with a parent-comment exemption claim is a coin-flip on β verdict.

**Closure shape**: the rule must (a) name the sub-issue body (or any issue body γ links from the dispatch prompt) as the only authoritative exemption surface, (b) explicitly exclude parent-comment claims, and (c) name two recovery paths so a cycle that fires 3.11b RC has a concrete way forward without re-debating the rule.

## §Skills

**Tier 1** (CDD lifecycle):
- `src/packages/cnos.cdd/CDD.md` — canonical lifecycle and role contract
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — α role surface (this file's authoring constraint)
- `src/packages/cnos.cdd/skills/cdd/review/SKILL.md` — patch target; loaded as authority for the rule being clarified, not as judgment frame (α may not rewrite β's review process per α §3.5)

**Tier 2** (always-applicable eng): none. The diff is a single rule-body edit inside one `SKILL.md`; no engineering toolchain is exercised (no Go / shell / YAML / JS surfaces touched).

**Tier 3** (issue-specific, named in dispatch):
- `src/packages/cnos.core/skills/write/SKILL.md` — generation constraint on the rule body. Specifically applied: §3.4 front-load the point (Exemption discoverability bullet leads with what *does* count, then names what does NOT), §3.8 replace abstractions with specifics (e.g. "sub-issue body β is reviewing" instead of "the issue"), §2.4 instance adjacent to abstraction (tsc #49 F2 citation embedded directly under the divergence claim), §3.11 lists reveal structure (recovery paths (a)/(b) listed because they are parallel and order matters — (a) canonical, (b) escape valve).

Not loaded: `design/SKILL.md`, `plan/SKILL.md` — explicitly not required per α §2.2. Design is "not required" because the incoherence is named in #360 §Problem and the closure shape is named in #360 §Acceptance Criteria; α's job here is rule-body authorship under existing constraints, not problem framing. Plan is "not required" because the implementation is a single Edit to one rule bullet — sequencing is trivial.

## §ACs

**AC1** — `review/SKILL.md §3.11b` specifies that "documented in the issue" means the sub-issue body (or any issue body γ links from the dispatch prompt).

- ✅ **Met**. The patched rule's `**Exemption discoverability**` bullet states: *"An exemption satisfies 3.11b only if it appears in the **sub-issue body** β is reviewing — i.e. the body of the cycle's own issue, or the body of any issue γ links from the dispatch prompt as authority for the cycle."*
- Evidence: `src/packages/cnos.cdd/skills/cdd/review/SKILL.md` rule 3.11b, bullet 5 (Exemption discoverability), commit `a3a34a16`.
- The phrase "any issue γ links from the dispatch prompt" handles the legitimate case where γ scaffolds a cycle that inherits its protocol exemption from an authority issue named in the dispatch — that authority issue's body is in α/β's load surface and so qualifies.

**AC2** — A comment on a master/parent issue does NOT satisfy 3.11b exemption discoverability.

- ✅ **Met**. Same bullet, second sentence: *"A comment on a parent / master / tracking issue does NOT satisfy exemption discoverability: β reviews the sub-issue β was dispatched against, not the parent tree, and master comments are not part of the cycle's load surface."*
- Evidence: `src/packages/cnos.cdd/skills/cdd/review/SKILL.md` rule 3.11b, bullet 5 (Exemption discoverability), commit `a3a34a16`.
- The denial is grounded in the load-surface argument (β reviews what β was dispatched against), not in a bare prohibition — so the rule explains *why*, which lets future readers judge edge cases.
- The tsc #49 F2 *Derives from* citation is embedded directly in the same bullet, fixing the rule's provenance for future audits.

**AC3** — Recovery paths documented: (a) author missing `gamma-scaffold.md` before β re-dispatch, OR (b) amend sub-issue body with explicit exemption section.

- ✅ **Met**. The patched rule's `**Recovery paths when 3.11b RC fires**` bullet states: *"either (a) γ (or δ-as-γ) authors the missing `.cdd/unreleased/{N}/gamma-scaffold.md` on the cycle branch before β re-dispatch, OR (b) γ amends the sub-issue body to add an explicit `## Protocol exemption` section naming the reason ... and β re-dispatches against the amended body."*
- Evidence: `src/packages/cnos.cdd/skills/cdd/review/SKILL.md` rule 3.11b, bullet 6 (Recovery paths), commit `a3a34a16`.
- The bullet names path (a) as canonical (matches CDD §1.4 γ scaffolding) and path (b) as the escape valve. It also names the section header β should grep for (`## Protocol exemption`), which keeps the recovery path mechanically checkable — β does not need judgment to verify path (b) was taken, only a grep.
- The final §Document bullet was also updated to require β to cite *which sub-issue body section* grants the exemption when one is claimed, so the audit trail is complete on β's side.

## §Self-check

**Did α push ambiguity onto β?** No — the rule body now answers the question that tsc #49 F2 left ambiguous: *which* issue body grants exemption. β does not need to interpret "the issue" anymore; the rule names "sub-issue body β is reviewing" or "any issue body γ links from the dispatch prompt as authority."

**Is every claim backed by evidence in the diff?** Yes. The three ACs each cite a specific bullet in rule 3.11b at commit `a3a34a16`. No claim in §ACs depends on prose outside the diff.

**Peer enumeration — rule-class peers.** Per α §2.3 (skill-class peers as a distinct enumeration class):

- **Role skills** (`alpha/`, `beta/`, `gamma/`, `operator/`): grep for `3.11b` across role skills.
  - `alpha/SKILL.md` — no mention of 3.11b. α §2.6 pre-review-gate rows 11–14 enumerate α-side authoring gates; β-side gates like 3.11b are out of α-skill scope. **No update needed.**
  - `beta/SKILL.md` — checked (not loaded in this cycle, per α §2.1 — α does not load β/γ role skills). However, the contract change is in `review/SKILL.md`, which is the lifecycle skill `beta/` calls. By the load-graph, β picks up the patched rule when it loads `review/`. Sibling enforcement surface is `beta/SKILL.md` pre-merge gate row 4 (γ artifact completeness, added by #355). Row 4 contains no exemption clause at all — it has been silent on exemptions since #355 shipped, independent of this patch's wording change. **No `beta/` patch needed for #360**, but the silence between 3.11b's exemption clause and row 4's no-exemption phrasing is pre-existing debt — flagged in §Debt below.
  - `gamma/SKILL.md` — checked similarly. γ's dispatch-prompt format already names "Issue: gh issue view {N}" pointing at the sub-issue, which is consistent with the patched rule's "sub-issue body" framing. **No `gamma/` patch needed**, but see Debt for a soft follow-up.
  - `operator/SKILL.md` — not in scope (operator does not emit 3.11b verdicts).
- **Lifecycle skills** (`review/`, `release/`, `post-release/`, `design/`, `plan/`, `issue/`): grep for `3.11b` across lifecycle skills.
  - `review/SKILL.md` — patched (this is the rule's home).
  - `release/SKILL.md` — no mention of 3.11b. Release operates on β-merged state; 3.11b fires before release. **No update needed.**
  - `post-release/SKILL.md` — no mention of 3.11b. **No update needed.**
  - `design/SKILL.md`, `plan/SKILL.md`, `issue/SKILL.md` — no mention of 3.11b. **No update needed.**

**Peer enumeration — intra-doc repetition (α §2.3).** grep for `3.11b` in `review/SKILL.md` itself:

- Line 122 — rule definition (patched).
- Line 206 — checklist row `| γ artifacts present (gamma-scaffold.md) | yes / no / n/a | rule 3.11b compliance |`. This row references 3.11b by number, not by exemption-scope phrasing; it remains correct after the patch because the rule body is the authority and the checklist points to it. **No checklist-row update needed.**

No other occurrences of `3.11b` in the document; the fact has one home.

**Harness audit.** Not applicable — no schema-bearing contract, parser, or runtime shape changed. The diff is rule-body prose.

**Polyglot re-audit (α §2.6 row 9).** The diff touches two surfaces:

- Markdown (`review/SKILL.md`, `self-coherence.md`) — re-read inline; cross-reference check passed (rule 3.11b body coherent with the §Document bullet and with the checklist row at line 206).
- No other languages present.

## §Debt

**Known debt at handoff:**

1. **`beta/SKILL.md` pre-merge gate row 4 silent on exemption interaction with 3.11b.** Row 4 (added by #355) says: *"If missing, verdict is RC with D-severity finding, classification protocol-compliance."* It contains no exemption clause. With this cycle's patch, `review/SKILL.md` 3.11b now lets β emit APPROVED on an exempted cycle (e.g. an emergency patch with `## Protocol exemption` in the sub-issue body), but row 4 — read literally — would still RC the same cycle at merge time. The two surfaces disagree under exemption. This disagreement pre-dates #360: row 4 has been silent on exemptions since #355 shipped. **Out of scope for #360** — the issue scope is exemption-discoverability inside 3.11b, not reconciling 3.11b with row 4. Recorded as follow-up: a separate cycle should patch row 4 to either inherit 3.11b's exemption clause or explicitly declare itself non-exempting (which would mean row 4 is the harder gate). Either disposition is γ's call.

2. **No new tests.** The diff is rule-body prose in a SKILL.md; the canonical CDD verification of the rule is β's next review (against this branch, then against subsequent cycles where the rule is invoked). No automated test exists or could be authored for "did the rule body resolve the ambiguity?" — the test is operational (does the next exemption-claiming cycle still produce divergent β verdicts?). Per α §2.6 row 3, this is the "explicit reason none apply" case.

3. **Historical artifacts unchanged.** `.cdd/unreleased/{355,357,359}/` reference rule 3.11b in past-tense audit records (`grep 3.11b .cdd/unreleased/` showed these). These are frozen cycle records, not active authority — they correctly cite 3.11b *as it stood at the time of those cycles*. No backfill is appropriate; the rule's evolution is preserved in `review/SKILL.md` git history, not in re-writing past close-outs.

## §CDD Trace

CDD lifecycle steps 1–7, per CDD.md §1.4 α algorithm.

**Step 1 — Receive.** Dispatched on `cycle/360`, base SHA `c77f34a4` (`origin/main` at dispatch). γ scaffold present at `.cdd/unreleased/360/gamma-scaffold.md` (commit `8888f2d2`). Issue #360 read in full via `gh issue view 360`; no comments. ACs enumerated: AC1 (sub-issue body as exemption surface), AC2 (parent comment denied), AC3 (two recovery paths).

**Step 2 — Skills.** Tier 1: CDD.md + alpha/SKILL.md + review/SKILL.md (as patch target). Tier 2: none. Tier 3: `cnos.core/skills/write/SKILL.md` (named in dispatch). design/plan explicitly not required — see §Skills above for reasons.

**Step 3 — Design.** Not required: incoherence is named in #360 §Problem; closure shape is in #360 §Acceptance Criteria.

**Step 4 — Plan.** Not required: single Edit to one rule bullet; sequencing trivial.

**Step 5 — Tests.** Not applicable. The diff is rule-body prose. Operational verification is β's next review of an exemption-claiming cycle. See §Debt item 2.

**Step 6 — Implementation artifacts.** Three files in `git diff --stat origin/main..HEAD`:

- `src/packages/cnos.cdd/skills/cdd/review/SKILL.md` (+6 / −2 lines) — rule 3.11b patched. Three bullets changed/added: §Scope tightened to one sentence; §Exemption discoverability added (AC1 + AC2); §Recovery paths added (AC3); §Document bullet extended to require citation of the exempting sub-issue body section. Committed at `a3a34a16`.
- `.cdd/unreleased/360/self-coherence.md` (+101 lines) — this file. Sections written incrementally: §Gap (`a3a34a16`), §Skills (`ca2af889`), §ACs (`53d30998`), §Self-check (`7139e8d6`), §Debt (`212be76d`), §CDD-Trace (this commit). Manifest in frontmatter `sections.completed` tracks progress.
- `.cdd/unreleased/360/gamma-scaffold.md` (+45 lines) — γ-authored at commit `8888f2d2`, base of the cycle. Not authored by α; included here for diff-stat enumeration completeness per α §2.6 row 11.

**Caller-path trace (α §2.6 row 12).** Not applicable — no new modules or functions added. The diff modifies prose in an existing SKILL.md rule that is already cited by `beta/SKILL.md` row 4 (γ artifact completeness, line 83) and by the checklist row at `review/SKILL.md` line 206. The patched rule's callers (β verdicts that cite 3.11b, β review checklist) are unchanged in their reference; they pick up the new body the next time β loads `review/SKILL.md`.

**Step 7 — Self-coherence (this file).** Carried by the seven incremental commits listed in Step 6. Each section pushed to `origin/cycle/360` before the next began, per α §2.5 incremental-write discipline.

## Review-readiness

**Ready for β.** Implementation complete at `a3a34a16`. All 3 ACs addressed. Self-coherence sections complete. Pre-review gate rows pass:

- Base SHA `c77f34a4` == `origin/main` at dispatch ✅
- All commits on `cycle/360`, authored as `alpha@cdd.cnos` ✅
- `.cdd/unreleased/360/self-coherence.md` present with all sections ✅
- Implementation diff: +6/−2 lines in `review/SKILL.md` rule 3.11b ✅

Note: review-readiness written by δ (operator override) after α session timeout. α completed all substantive work; only this signal section was missing.
