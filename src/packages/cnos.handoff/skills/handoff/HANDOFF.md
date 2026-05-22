# Coherence Handoff (HANDOFF)

cnos.handoff is the package that owns **wire formats and handoff doctrine** for work / state / messages moving between coherence cells, agents, roles, repos, and packages.

**Version:** 0.1 (Sub 2 of cnos#404)
**Status:** Draft — package contract + first substantive sub-surface (cross-repo doctrine). Forthcoming sub-surfaces (dispatch, mid-flight, artifact-channel, receipt-stream) land in Subs 3–5.

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

### Forthcoming (Subs 3–5 of cnos#404)

- `skills/handoff/dispatch/SKILL.md` — dispatch-prompt template + implementation-contract schema (per cnos#393). Migration Sub 3.
- `skills/handoff/mid-flight/SKILL.md` — gamma-clarification.md mid-flight rescue mechanism (cnos#391 anchor). Migration Sub 4.
- `skills/handoff/artifact-channel/SKILL.md` — `.cdd/unreleased/{N}/` α→β→γ intra-cycle handoff rules. Migration Sub 4.
- `skills/handoff/receipt-stream/SKILL.md` — `cdd-iteration.md` + INDEX.md ε feed. Migration Sub 5.

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

- Do NOT redesign STATUS / LINEAGE / case structure / feedback-patch format in v0.1 (cnos#416 is migration-only).
- Do NOT author CUE schemas under `schemas/handoff/` (deferred to a future cycle if needed; CCNF-O may type these from Track A).
- Do NOT compose cells into orchestration (CCNF-O / cnos#405 territory).
- Do NOT make handoff a runtime (Stage 4 territory).
