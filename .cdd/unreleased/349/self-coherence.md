# Self-Coherence — Cycle #349

## §Gap

**Issue:** #349 — Activation §11a: operator access flow for notification channels

**Version/Mode:** docs-only, design-and-build

**Gap closed:** cnos #344's activation skill §10 prescribes transport-agnostic notification interface and §11 prescribes secrets (bot token, chat ID), but stops at the channel. The operator-access flow — how a new operator joins the channel the bot posts to — was unspecified. This cycle adds §11a to cdd/activation/SKILL.md with invite-link convention, channel scope guidance, operator removal hygiene, and sample message shapes that notification adapters implement.

**Incoherence resolved:** Without operator-access prescription, tenants ship bot-side wiring successfully but experience onboarding friction for every new operator: manual invite-link generation, ad-hoc DM distribution, inconsistent link storage, and incomplete operator-removal (registry updated but channel access remains). §11a provides the canonical path that eliminates per-tenant improvisation.

## §Skills

**Tier 1 (CDD core):**
- src/packages/cnos.cdd/skills/cdd/CDD.md — canonical lifecycle
- src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md — α role surface and algorithm
- issue/SKILL.md — loaded implicitly for AC interpretation 

**Tier 2 (always-applicable eng):**
- eng/* bundle per src/packages/cnos.eng/skills/eng/README.md

**Tier 3 (issue-specific):**
- src/packages/cnos.core/skills/write/SKILL.md — loaded per dispatch prompt, applied as generation constraint for docs-only prose authoring

**Active constraints:** All sections in §11a authored under write/SKILL.md constraints (front-load the point, one paragraph one move, cut throat-clearing, brevity earned not chopped). Sample message shapes use transport-agnostic variable syntax per AC3 requirement.