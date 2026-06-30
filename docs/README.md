# cnos Documentation

cnos is a Git-native coherence system for accountable human + AI work.

> **Navigation note.**
> The documentation surface is organized by reader intent:
> `quickstart/`, `concepts/`, `guides/`, `reference/`, `development/`,
> `architecture/`, `papers/`, and `evidence/`.
>
> Frozen history is not kept on `main` — git history is the archive. The
> α/β/γ triad is no longer a filing taxonomy; it is kept only as a coherence
> measurement (role grammar and TSC analysis), never as folders.

## Start here

1. [THESIS.md](./THESIS.md) — what cnos is; the whole, above the triad.
2. [concepts/](./concepts/README.md) — the mental model: why coherence is the root concept.
3. [Dumb Models, Smart Cells](papers/DUMB-MODELS-SMART-CELLS.md) — why models are bounded executors, not authorities.
4. [reference/](./reference/README.md) — the CN protocol and why Git is the lowest durable substrate.

## I want to…

| Goal | Read |
|---|---|
| Understand cnos | [`THESIS.md`](./THESIS.md) → [`concepts/`](./concepts/README.md) → [`papers/`](./papers/README.md) |
| Try it / activate an agent | [`quickstart/`](./quickstart/README.md) |
| Build with cnos | [`reference/`](./reference/README.md) |
| Understand cells & receipts | [`development/`](./development/README.md) → [`papers/`](./papers/README.md) |
| Understand trust & coherence | [Dumb Models, Smart Cells](papers/DUMB-MODELS-SMART-CELLS.md) |
| Contribute | [`guides/`](./guides/README.md) → [`development/`](./development/README.md) |
| Find canonical specs | [`reference/`](./reference/README.md) |
| Find audits & evidence | [`evidence/`](./evidence/README.md) |

## Directory map (reader intent)

- [`quickstart/`](./quickstart/README.md) — runnable first experiences
- [`concepts/`](./concepts/README.md) — mental model and doctrine
- [`guides/`](./guides/README.md) — task-oriented how-tos
- [`reference/`](./reference/README.md) — canonical specs, APIs, CLI, schemas
- [`architecture/`](./architecture/README.md) — how the system fits together
- [`development/`](./development/README.md) — CDD method, rules, plans, checklists
- [`papers/`](./papers/README.md) — essays, whitepapers, position papers
- [`evidence/`](./evidence/README.md) — audits, RCAs, measurements, demo receipts

## Coherence metadata, not coherence folders

Documents may declare classification in frontmatter — `doc_type`, `status`,
`canonical`, `owner`, `supersedes`. These are *knowable* facts about a
document. TSC coherence axes (α/β/γ) are a **measurement**, recorded from a
TSC report, and must never be hand-written into frontmatter (a hand-authored
score is an invented score).

The filesystem is organized for readers. The triad is kept for measurement.

---

*The legacy α/β/γ reading model and the per-axis document tables now live in
each intent index above, pointing at the same files. See
[`reference/governance/DOCUMENTATION-SYSTEM.md`](reference/governance/DOCUMENTATION-SYSTEM.md)
for the documentation-system rules (versioned directories are frozen;
corrections ship as new versions or superseding notes).*
