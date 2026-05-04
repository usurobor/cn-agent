# Follow-up: Recast Activation Essays from Doctrine Lineage

Status: deferred follow-up to the activation-essays branch
([claude/activation-agency-essay-fWbHg](https://github.com/usurobor/cnos/tree/claude/activation-agency-essay-fWbHg))

## Context

The activation essays were drafted from product/architecture intuition and only later cited the agent doctrine ([CFA, EFA, JFA, IFA](../../docs/alpha/doctrine/README.md)) as ground. The current branch ships with a `Doctrine lineage` section in each essay that names the connection but does not derive the architecture from it.

A future cycle should rewrite the essays so the architecture is *derived* from doctrine rather than asserted and then cited.

## Acceptance criteria

For [`docs/alpha/essays/AGENT-CONTINUITY.md`](../../docs/alpha/essays/AGENT-CONTINUITY.md) and [`docs/alpha/essays/ACTIVATION-NOT-DEPLOYMENT.md`](../../docs/alpha/essays/ACTIVATION-NOT-DEPLOYMENT.md):

- AGENT-CONTINUITY starts from CFA and IFA, not from product intuition.
- Receipts are derived from CFA's claim/evidence/verdict attachment rule, not argued from system convenience.
- "False growth" is derived from CFA's comfort failure plus IFA's soft-inheritance failure mode.
- Kernel/Persona/Operator is explicitly inherited from `KERNEL.md`, not introduced as new taxonomy.
- ACTIVATION-NOT-DEPLOYMENT remains the comparison essay (MCP, OpenAI memory, LangGraph) but argues the activation/deployment split as application of CFA's structural commitments.
- No new doctrine is claimed; the essays are explicitly architectural application.

## Why deferred

The current essays are good enough as architectural orientation for new contributors. The recast would be a meaningful revision project — closer to a write-skill cycle than a patch — and should be undertaken when there is room to do it carefully. Shipping the lineage-cited version first gives readers the right pointers without blocking on the recast.

## Out of scope

- The recast is application work, not a new doctrine cycle. It does not need α/β/γ cycle artifacts.
- It does not change the placement: essays stay in `docs/alpha/essays/`, doctrine stays in `docs/alpha/doctrine/`.
- It does not require reopening IFA's closure stance.
