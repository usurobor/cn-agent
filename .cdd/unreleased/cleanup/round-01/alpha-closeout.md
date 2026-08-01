# α Close-out — Docs/Structure Remediation (Round 01)

**Role:** α (implementer)
**Scope:** DOCS/STRUCTURE only. Source code and package `SKILL.md`/`*.cue` behavior deferred to a later code pass.
**Input:** `beta-review.md` (F1–F7).
**Files changed:** 18 (all under root `*.md` + `docs/`). No commit — κ commits.

---

## Finding → action

| # | Sev | Status | What was done |
|---|-----|--------|---------------|
| F1 | HIGH | **FIXED** | Made README / CLI.md / OPERATOR.md tell one story. README is the shipped source of truth (unchanged). `CLI.md`: replaced stale `Status: Current / Date: 2026-03-06` header with a status that names the Go binary + README "What ships today" as the authoritative shipped set; tagged every non-shipped group **(planned)** — GTD, Agent Output, Agent Runtime, Orchestration, Thread Creation, Observability wholesale; per-line `(planned)` on non-shipped `deps`/`build`/hub commands with a "Shipped:" line each; banner marking the three OCaml implementation sections as the archived reference (active runtime is Go). `OPERATOR.md`: tagged the runtime-dependent surfaces `(planned)` (`cn logs` observing, Other-log-locations, Peers, `cn deps update`, troubleshooting log/daemon rows, Quick-Reference `cn logs`/`cn sync`); left shipped ops (`cn status`, `cn doctor`, `cn update`, `cn deps list/restore/doctor`, releasing) unmarked. |
| F2 | HIGH | **FIXED (docs surface) / DEFERRED (package files)** | Swept `src/agent/` → `src/packages/…` with correct package prefix across every active doc: `rules/INVARIANTS.md`, `cdd/README.md`, `cdd/RATIONALE.md`, `ENGINEERING-LEVELS.md`, `architecture/ARCHITECTURE.md`, `architecture/INVARIANTS.md`, `protocol/cn/THREAD-EVENT-MODEL.md`, `runtime/AGENT-RUNTIME.md`, `runtime/extensions/RUNTIME-EXTENSIONS.md`, `runtime/PROVIDER-CONTRACT-v1.md`, `runtime/ORCHESTRATORS.md`, `runtime/GO-KERNEL-COMMANDS.md`, `cli/SETUP-INSTALLER.md`, `cli/CLI.md`. Migration/design records (`BUILD-AND-DIST.md`, `CORE-REFACTOR.md`) got a completion banner + "Prior state" relabel instead of a blind swap, so their historical `src/agent/` mentions read correctly as pre-refactor. Occurrences inside `src/packages/**/SKILL.md` and `src/packages/**/README.md` (8 hits) DEFERRED → code pass. Genuine historical records (CHANGELOG, cycle log, frozen `plans/*`, `papers/MODULAR-*`, `PLAN-174`) left intact — see Residual. |
| F3 | MED | **FIXED** | `THESIS.md §17`: corrected `docs/alpha/doctrine/FOUNDATIONS.md` → `docs/papers/FOUNDATIONS.md`; the four bare doctrine names → `src/packages/cnos.core/doctrine/{CAP,COHERENCE,CA-CONDUCT,CBP,AGENT-OPS}.md` (all verified present). `§7.5`: removed the nonexistent `cnos.pm`; replaced the ad-hoc list with two illustrative packages + a pointer to the README package table (one home for the stable fact). |
| F4 | MED | **FIXED** | `THESIS.md` Abstract: added one line before the acronym cascade — "Each capitalized term below is defined at first use … and collected in the [Glossary]" — the minimum-text glossary pointer. Did not reorder the newcomer path or expand every acronym inline (avoids bloat; glossary already defines the cluster). |
| F5 | MED | **DEFERRED → code pass** | `ROLES.md` is cited by 157 root-relative `ROLES.md §N` references — 129 in docs/root, **28 inside `src/packages/**`** (out of scope). Moving the file would break the package-skill citations, which this pass may not edit. No docs-side edit made; keeping the file at root avoids introducing 129 broken doc refs mid-pass. The move + ref-rewrite must land as one coordinated code-pass change. |
| F6 | LOW | **NO CHANGE (rationale recorded)** | `docs/development/board/{index.html,board-data.json}` are generated (confirmed: `board/README.md` names the generator `cn issues map` and a `board-map` GitHub Action, marked do-not-edit). They are **intentionally committed** — served as a live preview via `raw.githack.com` and browsable in-tree. β's suggested gitignore+remove would break the documented live-preview links, so it was not applied. The existing do-not-edit README already answers "why is build output here." |
| F7 | LOW | **FIXED** | `OPERATOR.md` L5: `[README.md quickstart](README.md)` → `[README "Try it"](README.md#try-it)` — label now names the real section (`## Try it`) with a matching anchor. |

---

## Verification

### F2 — no stale `src/agent/` in the active docs surface

```
$ git grep -n 'src/agent/' -- docs/ '*.md' ':!.cdd/' \
    | grep -vE 'CHANGELOG.md|plans/|papers/MODULAR|PLAN-174|JFA-cycle-log|BUILD-AND-DIST|CORE-REFACTOR'
src/packages/cnos.core/skills/README.md:28 …
src/packages/cnos.core/skills/README.md:99 …
src/packages/cnos.core/skills/README.md:177 …
src/packages/cnos.core/skills/agent/agent-ops/SKILL.md:143 …
src/packages/cnos.core/skills/agent/agent-ops/SKILL.md:179 …
src/packages/cnos.core/skills/skill/SKILL.md:3 …
src/packages/cnos.core/skills/skill/SKILL.md:43 …
src/packages/cnos.core/skills/skill/SKILL.md:284 …
```

Every remaining hit is inside `src/packages/**` (package skill/README bodies) → **DEFERRED to code pass**. Zero stale `src/agent/` refs remain in root `*.md` or the `docs/` reader surface except the deliberately-labeled historical/migration records listed under Residual.

### F1 — three files agree on shipped vs planned

Shipped set (README "What ships today", authoritative): `help, init, setup, deps list/restore/doctor, status, doctor, build (+--check), update`.

- `CLI.md`: all commands outside that set now carry `(planned)`; header points to the README set as source of truth. ✓
- `OPERATOR.md`: `cn logs`, `cn peer`, `cn sync`, `cn agent*`, `cn deps update` all marked `(planned)`; shipped ops unmarked. ✓
- No command appears as "current/runnable" in CLI.md or OPERATOR.md while README lists it as planned. ✓

### F3/F4 — new link targets resolve

- `docs/papers/FOUNDATIONS.md` ✓, `src/packages/cnos.core/doctrine/{CAP,COHERENCE,CA-CONDUCT,CBP,AGENT-OPS}.md` ✓ (all present).
- `docs/reference/governance/GLOSSARY.md` ✓ (glossary pointer target).
- README anchors `#try-it`, `#what-ships-today`, `#how-the-project-is-organized` all present (`## Try it`, `## What ships today`, `## How the project is organized`). ✓
- `THESIS.md` no longer contains `docs/alpha/` or `cnos.pm`. ✓

### Changed-file count

`git status --porcelain` → **18 modified files**, all root `*.md` + `docs/`. README.md unchanged (already correct).

---

## Residual (for round 02 to re-check)

1. **Historical `src/agent/` refs left intact by design** (rewriting would falsify a dated record; the `write`/`document` skills treat changelog + logs as history): `CHANGELOG.md` (4), `docs/concepts/doctrine/judgment-for-agents/JFA-cycle-log-dyad.md` (a cycle log recording a 404), `docs/development/plans/{INVARIANT-HARDENING-v1,PLAN-package-system,PLAN-runtime-extensions,PLAN-v3.22.0-eng-lane-clarity,PR-docs-governance-v3.13.0}.md`, `docs/papers/MODULAR-ARCHITECTURE-REFACTOR.md` (dated `Status: Proposed` decision record), `docs/reference/runtime/PLAN-174-orchestrator-runtime.md`. If the cell wants zero stale paths even in frozen records, round 02 should add a one-line "paths historical" note per file — I judged that lower-signal than leaving the dated headers to speak for themselves.

2. **Labeled migration/design records** (`BUILD-AND-DIST.md`, `CORE-REFACTOR.md`) still contain `src/agent/` in their bodies, now framed by a completion banner. Confirm the framing reads correctly to a newcomer.

3. **F5 move** is unstarted — round 02 (or the code pass) must relocate `ROLES.md` and rewrite all 157 `ROLES.md §N` citations (129 docs/root + 28 `src/packages/**`) atomically.

4. **NEW drift, outside β's F1–F7 (flag for a future β round, not remediated here):** a broad `docs/alpha/` folder-path cluster survives across `reference/governance/DOCUMENTATION-SYSTEM.md`, `reference/ctb/*`, `evidence/design/WRITER-PACKAGE.md`, `papers/COHERENCE-SYSTEM.md`, `architecture/HUB-PLACEMENT-MODELS.md`, and others — the same "α/β/γ as folders" scheme `docs/README.md` says was retired. Likewise a broad `cnos.pm` cluster (`CAR.md`, `PACKAGE-RESTRUCTURING.md`, `ORCHESTRATORS.md`, `COHERENCE-SYSTEM.md`). F3 only scoped the THESIS occurrences; the wider sweep is a new finding.

5. **F1 code-adjacent drift not touched** (out of docs scope): `CLI.md`'s OCaml file tree and `ARCHITECTURE.md`/`THREAD-EVENT-MODEL.md` OCaml module names (`cn_*.ml`) describe the archived runtime; only the path/status framing was corrected. The Go runtime's real module surface is a code-pass concern. `SETUP-INSTALLER.md` `SOUL.md` template path could not be verified against any file under `src/packages/` — flagged inline, confirm against the Go `cn setup` source.
