# Documentation System

This document explains how the cnos docs tree is organized and where a given document lives.

**Version:** 3.82.0
**Date:** 2026-07-31

---

## 1. Layout

The docs tree is organized by **reader intent**: each top-level directory answers one kind of question a reader arrives with. A document lives in the directory whose question it answers.

| Directory | The reader asks | Holds |
|-----------|-----------------|-------|
| (root) | What is cnos? How do I read these docs? | `THESIS.md`, `README.md` |
| `quickstart/` | How do I run it now? | Runnable first experiences |
| `concepts/` | What is the mental model? | Doctrine, the coherence model |
| `guides/` | How do I do a task? | Task-oriented how-tos |
| `reference/` | What is the exact contract? | Canonical specs, APIs, CLI, schemas, protocol |
| `architecture/` | How do the parts fit? | System overview, invariants, constraints |
| `development/` | How is the work done? | CDD method, rules, checklists |
| `papers/` | Why is it built this way? | Essays, whitepapers, position papers |
| `evidence/` | Did it hold up? | Audits, RCAs, measurements, demo receipts |

`docs/README.md` is the navigation index over these directories. This document is the rule for what belongs in each.

The α/β/γ triad is **not** a filing taxonomy. It is a coherence measurement (role grammar and TSC analysis), recorded from a report, never written into folder names. There is no `docs/alpha/`, `docs/beta/`, or `docs/gamma/`; see [`GLOSSARY.md`](GLOSSARY.md).

---

## 2. Feature bundles

A **feature bundle** is a subdirectory that groups every document about one feature or subsystem, inside the intent directory that fits the feature. `reference/runtime/` and `reference/governance/` are bundles.

**The bundle contract:** every bundle has a `README.md` as its entrypoint, and that README names exactly one canonical spec as the normative source of truth.

The bundle's identity is its feature scope, so the directory uses the **kebab-case feature name** (`reference/runtime/extensions/`), not a version or issue number. Provenance — originating version, issue — goes in the bundle's `README.md`, not the directory name.

Create a bundle when a feature has a canonical spec plus at least one supporting document (a design narrative, a related plan, or a snapshot). A single-document feature (e.g. `architecture/DESIGN-CONSTRAINTS.md`) does not need one.

### Bundle README requirements

- Feature name and one-sentence purpose
- Which document is the canonical spec, by name and path
- A document map: every file in the bundle, one line each
- Reading order for a new reader

---

## 3. Document classes

Every document is one of these. The class decides which intent directory it lives in.

| Class | What it is | Lives in |
|-------|-----------|----------|
| Thesis | The whole, above the triad. One file. | `THESIS.md` (root) |
| Canonical spec | The single source of truth for its scope. Evolves in place. | `reference/`, `architecture/` |
| Reference | Stable lookup material (glossary, naming, conventions). | `reference/` |
| Guide | Task-oriented procedure connecting operator to system. | `guides/` |
| Concept / doctrine | The mental model and the doctrine essays. | `concepts/` |
| Paper | Essay, whitepaper, or position paper arguing a design. | `papers/` |
| Evidence | Audit, RCA, or model↔reality assessment. | `evidence/` |

A **canonical spec** evolves in place — it is never forked into parallel copies at the same level. It carries a version and date header, and absorbs design narratives after they ship.

---

## 4. Versioning

### Single version lineage

Every document uses **cnos release versions**. There is no independent per-document version lineage.

The `Version:` header records the cnos release in which the document was **last substantively changed**. `Version: 3.82.0` means "last updated as part of cnos v3.82.0." One question — "what does cnos v3.82.0 look like?" — then answers specs, governance, runtime, and docs with one number.

### When a version advances

| Change | Advances the version? |
|--------|-----------------------|
| Wording, examples, typos | Yes — next patch release |
| New section, additive capability | Yes — next minor release |
| Scope change, structural rewrite | Yes — next minor or major release |
| No change in a release | No — stays at the last release that touched it |

### Legacy version numbers

Some documents still carry pre-alignment numbers (e.g. THESIS 1.0.0). These re-version to a cnos release number the next time they are substantively updated. No bulk rename is required.

---

## 5. How documents evolve

### Supersession

When a canonical document is fully replaced, the replacement's version is the cnos release in which it first appears; it does not restart at 1.0.0. The superseded document records `Supersedes:` / `Superseded by:` in its header so a reader can trace the lineage. Do not keep parallel live versions of the same document.

### Frozen history

Released snapshots are **not** kept as folders on `main` — git history is the archive. To read a document as it stood at a past release, check out that release tag. A document on `main` is always the current one; its `Version:` header says which release it belongs to.

A completed, dated, or superseded record — a design decision, a cycle log, a finished plan — is **not kept on the reader surface**. Git history is its archive: to read it, check out the release or commit under which it was written. The working tree on `main` reflects only current state. Records that must persist as intentional history live in the dotdirs (`.cdd/`, `.cn-sigma/`, `.github/`), which are exempt from the current-state rule. A document that remains on the reader surface is, by that fact, current — its paths and contents are kept correct, not frozen.

---

## 6. Placement

When adding a document, ask **what question does the reader bring to it**, and place it in that intent directory (§1). Within the directory, add it to a feature bundle if one fits (§2); otherwise it stands alone.

If a document fits no intent directory, the taxonomy may need to evolve — update this document before creating new structure.

---

## 7. CI validation

These invariants should be enforced by CI:

- Every feature bundle directory has a `README.md` that names exactly one canonical spec.
- No document declares an α/β/γ folder path (`docs/alpha/`, `docs/beta/`, `docs/gamma/`).
- Internal links and backtick paths resolve to a file on disk.
