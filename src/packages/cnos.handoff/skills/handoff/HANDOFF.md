# Coherence Handoff (HANDOFF)

cnos.handoff is the package that owns **wire formats and handoff doctrine** for work / state / messages moving between coherence cells, agents, roles, repos, and packages.

**Version:** 0.1 (Subs 2–5 of cnos#404)
**Status:** Draft — package contract + five substantive sub-surfaces (cross-repo doctrine; dispatch + implementation-contract; mid-flight rescue; artifact-channel rules; receipt-stream + INDEX.md aggregator). All five v0.1 sub-skills landed; Sub 6 (cross-reference cleanup + close tracker) remains.

---

## Package contract

cnos.handoff owns:
- Wire formats for handoff between activations of the same agent and between different agents
- Doctrine for inter-repo, inter-role, and inter-package work transport
- State machines for handoff lifecycles (proposal STATUS, mid-flight rescue, receipt-stream aggregation)
- Bundle, patch, and pointer formats for cross-side transport

cnos.handoff does NOT own:
- CCNF kernel (lives in cnos.cdd)
- Cycle lifecycle protocols (CDS / CDR own their respective per-domain lifecycles)
- Orchestration grammar (CCNF-O / cnos#405 Track A; composes cells into cycles/waves/roadmaps)
- Validators (cn cdd verify lives in cnos.cdd/commands)
- Runtime execution (cn runtime; Stage 4)

cnos.handoff is **consumed** by cnos.cdd / cnos.cdr / cnos.cds today; consumed by any future c-d-X protocol package; consumed by any future orchestration grammar.

---

## Sub-surfaces

### Landed (v0.1)

- **[skills/handoff/cross-repo/SKILL.md](cross-repo/SKILL.md)** — cross-repo state machine + bundle file sets + LINEAGE schemas + feedback-patch format + archival rule + hat-collapse attribution. Migrated from cnos.cdd/skills/cdd/cross-repo/SKILL.md in cnos#416 (Sub 2 of cnos#404).
- **[skills/handoff/dispatch/SKILL.md](dispatch/SKILL.md)** — dispatch-prompt template (γ / α / β) + 7-axis implementation-contract schema + δ-as-inward-membrane enrichment doctrine + four-surface mesh declaration. Synthesized from cnos.cdd/skills/cdd/gamma/SKILL.md §2.5 + cnos.cdd/skills/cdd/operator/SKILL.md §3a + cnos.cdd/skills/cdd/delta/SKILL.md §2 in cnos#417 (Sub 3 of cnos#404).
- **[skills/handoff/mid-flight/SKILL.md](mid-flight/SKILL.md)** — `gamma-clarification.md` rescue mechanism: file format + authoring role + reader role + resumption protocol + issue-edit cache-bust + spec-staleness propagation. Authored from cnos#391 empirical anchor in cnos#418 (Sub 4 of cnos#404).
- **[skills/handoff/artifact-channel/SKILL.md](artifact-channel/SKILL.md)** — `.cdd/unreleased/{N}/` α→β→γ sequential intra-cycle channel: directory pattern + per-role write ownership + sequential rule + frozen-snapshot rule + release-time move + cross-cycle aggregation pointer. Extracted from cnos.cds/skills/cds/CDS.md §"Artifact contract" wire-format invariants in cnos#418 (Sub 4 of cnos#404).
- **[skills/handoff/receipt-stream/SKILL.md](receipt-stream/SKILL.md)** — cross-cycle `cdd-iteration.md` receipt-stream + `.cdd/iterations/INDEX.md` aggregator: per-finding shape (six top-level fields + per-disposition sub-fields) + INDEX row format (eight pipe-separated columns) + cadence rule (required when `protocol_gap_count > 0`; courtesy empty-findings stub when 0) + finding-disposition vocabulary (`patch-landed` / `next-MCA` / `no-patch`) + cross-repo trace bundle invariant. Migrated from cnos.cdd/skills/cdd/post-release/SKILL.md §5.6b in cnos#419 (Sub 5 of cnos#404).

---

## Related documents

- `cnos.cdd/skills/cdd/CDD.md` (CCNF kernel; handoff consumes / is consumed by cdd)
- `cnos.cdr/skills/cdr/CDR.md` (research domain; consumes handoff)
- `cnos.cds/skills/cds/CDS.md` (software domain; consumes handoff)
- `cnos.cdd/skills/cdd/COHERENCE-CELL.md` and `COHERENCE-CELL-NORMAL-FORM.md` (kernel doctrine; handoff transports between cells)
- cnos#404 (parent tracker)
- cnos#405 (CCNF-O orchestration grammar — peer; composes cells over handoff's transport)
- `docs/gamma/essays/CCNF-AND-TYPED-TRUST.md`
- `docs/gamma/essays/DECREASING-INCOHERENCE.md`

---

## Non-goals

- Do NOT redesign STATUS / LINEAGE / case structure / feedback-patch format / dispatch-prompt template / 7-axis implementation-contract schema / gamma-clarification.md rescue mechanism / artifact-file names / channel directory pattern / per-finding shape / INDEX row format / cadence rule / finding-disposition vocabulary in v0.1 (Subs 2–5 / cnos#416 + cnos#417 + cnos#418 + cnos#419 are migration-only).
- Do NOT author CUE schemas under `schemas/handoff/` (deferred to a future cycle if needed; CCNF-O may type these from Track A).
- Do NOT compose cells into orchestration (CCNF-O / cnos#405 territory).
- Do NOT make handoff a runtime (Stage 4 territory).
