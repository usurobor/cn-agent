---
cycle: 360
role: alpha
issue: "https://github.com/usurobor/cnos/issues/360"
date: "2026-05-14"
base_sha: "c77f34a4"
sections:
  planned: [Gap, Skills, ACs, Self-check, Debt, CDD-Trace, Review-readiness]
  completed: [Gap, Skills]
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
