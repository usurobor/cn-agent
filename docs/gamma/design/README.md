# Design — Decision Bundle

Survey + decision documents that pin shape for downstream implementation cycles. Each doc here is a **decision artifact** — not an essay (position paper), not an implementation. The class is "survey + decision" per `cnos.core/skills/design/SKILL.md §1.0`.

These docs sit alongside `docs/gamma/essays/` (position papers) and the role-skill / handoff / cdd / cds / cdr packages (implementation surfaces). A doc lands here when:

- The cycle is a survey / decision cycle (`Mode: explore` per `cnos.cdd/skills/cdd/issue/SKILL.md`).
- The output is a pinned-decisions matrix or design-touchpoint enumeration that downstream cycles dispatch against.
- The output is not yet a CUE schema, a Go type, a role-skill rule, or a kernel doctrine edit (those land in their respective package homes once the design pin holds).

---

## Document Map

| Document | Class | Description |
|----------|-------|-------------|
| [ccnf-o-track-a1-survey.md](./ccnf-o-track-a1-survey.md) | Survey + decision (v0.1) | Track A1 of [cnos#405](https://github.com/usurobor/cnos/issues/405): name-pick (CCNF-O), 20-surface inventory matrix, higher-level form classification, TSC integration v0.1 scope, package-location pin (`cnos.ccnf-o`), and refined dispatch-brief shapes for Tracks A2–A6 + B1. |

---

## Relationship to other doc directories

- **`docs/gamma/essays/`** — position papers. Read essays for the *why* and the long-form vision; read design docs for the *what and where* a downstream cycle must pin.
- **`docs/gamma/cdd/`** — CDD-class docs and historical post-release assessments (PRAs). Read for archival precedent; do not edit historical records.
- **`docs/alpha/essays/`** — foundational position papers (C≡ / TSC / CTB / cnos stack). Cited from γ-class designs.
- **`src/packages/*/skills/`** — executable doctrine (the *how*). Design docs here gate skill / schema authoring there.
