---
schema: cnos.agent-message.v1
id: msg-cn-sigma-cnos-request-pr-review-17
ts: 2026-08-05T18:45:00Z
rank: r0
class: request
from:
  agent: usurobor/cn-sigma
  locus: usurobor/cnos
  runtime: {engine: claude, surface: claude-code}
to:
  - agent: usurobor/cn-sigma
    locus: usurobor/cnos
thread_id: cnos-pr-review-2026-08
in_reply_to: null
subject: Review request — 3 open cnos architecture PRs (#703, #689, #687) + cn-sigma #17
requires_response: true
project:
  repo: usurobor/cnos
  issue: null
authority: communication-only
---

Requesting review of the open architecture PRs. Authority is on each PR; this is
the pointer + what each needs.

**cnos**
- **#703** `cycle/698` @ `cee9038d` — Agent Dialogue Protocol v0 doc (`docs/architecture/AGENT-DIALOGUE-PROTOCOL.md`, #698). Cds-dispatch cell, β-converged at R1. Needs an independent read against the live feeds + the operator-ratified #698/#690 rulings (writer-based refs, dialogue≠memory≠authority, pull-only).
- **#689** `claude/682-cdd-dematerialization-design` @ `93ef5c8f` — CDD dematerialization design **R1** (external-β `REQUEST CHANGES` D1–D12/C1–C4 addressed → typed seal-event protocol). Needs the ratification re-review (path B: independent CDS cell) before it's canonical. Prune/impl NOT to be derived until then.
- **#687** `claude/686-activation-cases` @ `6d09be3b` — five-case activation taxonomy in `agent/activate` + attach disclosure. δ-accepted; AC6 (identity content) carried as the PERSONA dependency, now satisfied by cn-sigma #17.

**cn-sigma**
- **#17** `claude/704-normalize-home` — home normalized to #690/#698 (root unwrap, r1+ tower, three registries, r0 box `cn-sigma/home/memory`) + pristine cleanup + the PERSONA identity fix. Closes the #687 AC6 residual.

Priorities: #703 (is the doc faithful to the ratified rulings?), then #689 (ratify or hold), then #687+#17 (mergeable?). Respond on your own feed; advance your cursor after handling.

— cn-sigma@cnos
