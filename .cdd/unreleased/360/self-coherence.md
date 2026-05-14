---
cycle: 360
role: alpha
issue: "https://github.com/usurobor/cnos/issues/360"
date: "2026-05-14"
base_sha: "c77f34a4"
sections:
  planned: [Gap, Skills, ACs, Self-check, Debt, CDD-Trace, Review-readiness]
  completed: [Gap]
---

# α Self-coherence — #360

## §Gap

**Issue**: #360 — `cdd: review rule 3.11b — clarify exemption scope (sub-issue body required)`.

**Mode**: skill-patch on `review/SKILL.md` rule 3.11b. Source: tsc v0.10.0 wave `cdd-iteration` F2 (`.cdd/releases/0.10.0/49/cdd-iteration.md`).

**Incoherence**: rule 3.11b §Scope says exemptions are "documented in the issue" without naming *which* issue body counts. In tsc #49 a cycle had its exemption claim land as a comment on the parent/master issue; four independently-dispatched β subagents diverged — β@S1 read 3.11b literally and emitted D-severity RC, β@S2/S3/S4 accepted the master-issue comment as exemption and treated the protocol miss as a B non-blocker. Divergence was structural (the rule under-specified exemption discoverability), not careless. Until 3.11b names *which* issue body grants exemption, every cycle with a parent-comment exemption claim is a coin-flip on β verdict.

**Closure shape**: the rule must (a) name the sub-issue body (or any issue body γ links from the dispatch prompt) as the only authoritative exemption surface, (b) explicitly exclude parent-comment claims, and (c) name two recovery paths so a cycle that fires 3.11b RC has a concrete way forward without re-debating the rule.
