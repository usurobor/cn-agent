# β Review — Newcomer / Structure Pass (Round 01)

**Reviewer:** β (independent adversarial newcomer review)
**Scope:** root `*.md`, `docs/` tree, top-level structure, navigation, links. Source code and package `SKILL.md` behavior are out of scope (deferred to a code pass).
**Lens:** `write` skill (signal/noise, one governing question per file, state what it *is*, say a fact once), `document` skill (every claim traced to source), L7 doc bar, CCNF cell doctrine.

---

## Verdict: **REVISE**

The newcomer surface is *above average* — the root `README.md` passes the headline test (a first-time reader knows what cnos is, who it's for, and can `curl | sh` within the first screen), the `docs/README.md` intent map is clean, and the I4/lychee markdown-link gate is holding (0 broken `[](…)` links across root + `docs/`). But the doc set **contradicts itself on what actually ships**, and a large stale-path cluster (`src/agent/` → `src/packages/`) sits *underneath* the link gate because those references are bare code spans, not markdown links. Both mislead a sharp newcomer on a critical path. Neither is polish; both must land before this surface is pristine.

## Newcomer impression (one paragraph)

I understood cnos from README's first screen and I trusted it. Then I opened the canonical `CLI.md` (marked **Status: Current**) and `OPERATOR.md` and found a full daemon/logs/peer/sync CLI presented as usable — which the README had just told me is *not shipped yet*. That is the exact README-vs-reference contradiction that erodes a newcomer's trust in the whole doc set. Following skill-file paths in the CDD and rules docs (`src/agent/…`) led nowhere, because the tree is now `src/packages/`. And the document the docs index tells me to read *first* (`THESIS.md`) opens with eleven undefined acronyms before defining any. The bones are good; the surface has drift the link gate can't see.

---

## Findings (ranked)

### F1 — HIGH — Shipped-vs-planned contradiction across README / CLI.md / OPERATOR.md

- **Where:** `README.md` L79–109 (declares `cn agent`, daemon/scheduler, inbox/outbox sync, peer sync, CN Shell as **planned, not shipped**) vs `docs/reference/cli/CLI.md` (header **Status: Current**, lists `cn agent`, `cn agent --daemon/--stdio`, `cn sync`, `cn logs …`, `cn peer …` as current commands) vs `OPERATOR.md` §1 Observing (`cn logs`), §2 (`cn peer`), §5 Quick Reference (`cn sync`, L230) — presented as operational, with only §"Start/Stop" and §"Configuration" carrying the "Target runtime (planned)" caveat.
- **Problem:** A newcomer who reads README learns the runtime is not shipped, then reads the *canonical CLI reference* and the operator manual and sees the same runtime as live, runnable commands. `cn logs`, `cn peer`, `cn sync` are in README's *planned* list but in CLI.md's *Current* list.
- **Why it violates the bar:** `document` skill — "every claim verified against source of truth"; the reference doc and README describe two different products. This is the tsc README-vs-STATUS failure the brief names explicitly.
- **Fix (α):** Pick the shipped set as the single source of truth (README's shipped table + `cn --help`). In `CLI.md` and `OPERATOR.md`, mark every not-yet-shipped command with a consistent **planned** tag (a column or inline marker), matching README. Correct `CLI.md`'s header (`Status: Current`, `Date: 2026-03-06` — 5 months stale vs `VERSION` 3.82.0) to reflect that it documents target + shipped, and say which is which.

### F2 — HIGH — 61 stale `src/agent/` path references (source moved to `src/packages/`)

- **Where:** `src/agent/` does not exist; the tree is `src/{go,packages}`. 61 references remain across `docs/`. Live, newcomer/governance-facing offenders:
  - `docs/development/rules/INVARIANTS.md` L31, L57, L65 — calls `src/agent/` "the **only** source of truth for doctrine, mindsets, skills" (a RULES doc asserting a dead path as canonical).
  - `docs/development/cdd/README.md` L50–56 — points readers to `src/agent/skills/cdd/…/SKILL.md`; the real files are at `src/packages/cnos.cdd/skills/cdd/…` (verified present).
  - `docs/development/ENGINEERING-LEVELS.md` L12, L201 — the L7 rubric doc itself.
  - `docs/development/cdd/RATIONALE.md` L198; `docs/reference/runtime/{RUNTIME-EXTENSIONS,ORCHESTRATORS,AGENT-RUNTIME,CORE-REFACTOR}.md`.
  - Remainder in `docs/development/plans/*` (historical plan records).
- **Problem:** Every path resolves to nothing. A newcomer chasing the CDD skill entrypoints or the "source of truth" rule lands on a 404.
- **Why it violates the bar:** `document` §3.1 "trace every claim to source"; §"no stale references to removed features." The I4/lychee gate does **not** catch these — they are bare backtick code spans, not `[](…)` links — so "links resolve" is only true for markdown-syntax links, not for the paths docs actually cite.
- **Fix (α):** Sweep live docs `src/agent/` → `src/packages/…` with the correct package prefix (`cnos.core`, `cnos.cdd`, `cnos.cds`, `cnos.eng`, per what the file references). Prioritize `rules/INVARIANTS.md`, `cdd/README.md`, `ENGINEERING-LEVELS.md`, `cdd/RATIONALE.md`, and the four `reference/runtime/` docs. For frozen `plans/*`, either leave with a single "paths historical; tree is now `src/packages/`" note or fix in-place. Consider adding a lychee/CI check that also validates inline code-span paths, not just markdown links (closes the gate gap this finding exposes).

### F3 — MED — `THESIS.md` §17 & §7.5 cite nonexistent paths and a nonexistent package, contradicting `docs/README.md`

- **Where:** `docs/THESIS.md` L703–709 §17 "Relationship to Other Documents" lists `docs/alpha/doctrine/FOUNDATIONS.md`, `CAP.md`, `COHERENCE.md`, `CA-CONDUCT.md`, `CBP.md`, `AGENT-OPS.md`. `docs/alpha/` does **not** exist; `FOUNDATIONS.md` is at `docs/papers/FOUNDATIONS.md`; the other four live at `src/packages/cnos.core/doctrine/`. §7.5 (L336–339) lists package `cnos.pm` — not in the tree (`src/packages/` has no `cnos.pm`; README's package table omits it too).
- **Problem:** The `docs/alpha/` prefix directly contradicts `docs/README.md` L6–12, which states the α/β/γ triad is "no longer a filing taxonomy … never as folders." The flagship thesis points at the exact folder scheme the docs system says was removed.
- **Why it violates the bar:** `document` §"drift"; `write` §3.13 "keep authority explicit" — a newcomer cannot tell which package list (THESIS §7.5 vs README) is authoritative.
- **Fix (α):** Correct §17 to real locations (drop `docs/alpha/`; point to `docs/papers/FOUNDATIONS.md` and `src/packages/cnos.core/doctrine/*`). In §7.5, replace the ad-hoc package list with the real set or a pointer to README's package table (one home for the stable fact).

### F4 — MED — `THESIS.md` is the docs index's #1 "Start here" but opens with ~11 undefined acronyms

- **Where:** `docs/README.md` L16 and `docs/concepts/README.md` L7 both send the first-time reader to `THESIS.md` first. `THESIS.md` Abstract (L28–38) fires CAP, MCA, MCI, MCP, CMP, CLP, CAA, CAR, AGENT-RUNTIME, TRACEABILITY, CDD before any is defined.
- **Problem:** The designated first read is a jargon wall. The reader must hold 11 undefined labels for pages before payoff.
- **Why it violates the bar:** `write` §2.1 "start with the point," §2.5/§3.8 "concrete nouns over vague labels." (README itself does the plain-language job well — the pattern to reuse exists in-repo.)
- **Fix (α):** Either (a) reorder the newcomer path so README → `concepts/` → THESIS, demoting THESIS from literal #1 for a cold reader; or (b) prepend a one-screen plain-language orientation to THESIS (borrow README's "cnos gives AI agents a Git-native home" framing) that lands the point before the acronym cascade. Do not expand every acronym inline — front-load the plain claim, then let the acronyms unfold.

### F5 — MED — Root-directory doctrine sprawl: `ROLES.md` (510 lines / 38 KB) misplaced at repo root

- **Where:** repo root holds `ROLES.md` (a 510-line role-ladder doctrine treatise) alongside standard meta files. Root convention is meta (`README`, `LICENSE`, `CONTRIBUTING`, `SECURITY`, `CODE_OF_CONDUCT`, `CHANGELOG`). `OPERATOR.md`, `RELEASE.md`, `SUSTAINABILITY.md` are also non-standard at root but defensible (operator entry, release notes consumed by the release workflow, funding). `ROLES.md` is pure deep doctrine.
- **Problem:** A newcomer scanning root for orientation meets a 38 KB doctrine document with no meta-file reason to be there.
- **Why it violates the bar:** `write` §1.2 "one file answers one question and points elsewhere"; structural signal/noise — root should orient, not host doctrine.
- **Fix (α):** Move `ROLES.md` under `docs/` (candidate: `docs/reference/governance/` or `docs/concepts/`) or into `src/packages/cnos.cdd/`, and update the root-relative `ROLES.md §N` citations (e.g. `CELL-KINDS.md` cites `ROLES.md §4b`). **Coordinate with the code pass** — several package `SKILL.md` files cite `ROLES.md` by root-relative path, so the ref updates cross the doc/code boundary. If it must stay at root, README should name why in one line.

### F6 — LOW — Large generated artifact committed inside the docs tree

- **Where:** `docs/development/board/index.html` (382 KB) + `board-data.json` (62 KB), both recently regenerated (Jul 31).
- **Problem:** A ~444 KB generated dashboard lives in a human-docs tree; it reads as build output, not documentation.
- **Fix (α):** Confirm whether it is generated; if so, gitignore it and generate on demand, or move it out of `docs/` into a `dist/`-style location. Low priority — it is gated behind the `development/` intent bucket and has a README.

### F7 — LOW — `OPERATOR.md` label drift

- **Where:** `OPERATOR.md` L6 — "For install: [README.md quickstart](README.md)."
- **Problem:** README has no "quickstart" section; the install section is titled "## Try it." Link resolves (no anchor) but the label names a section that isn't there.
- **Fix (α):** Rename the label to "Try it" or add a matching anchor. Trivial.

---

## Convergence assessment

- **Counts:** 2 HIGH, 3 MED, 2 LOW.
- **Distance from pristine:** Moderate. The architecture of the newcomer surface is sound (README headline, docs intent map, working markdown-link gate). The defects are *accuracy drift* — the doc set describes a partly-different product (F1) and cites a partly-different tree (F2) — plus one contradiction inside the flagship thesis (F3). None require restructuring the doc set; they require reconciling claims to source.
- **Another round needed?** **Yes.** F1 (shipped/planned reconciliation) and F2 (61-reference `src/agent/` sweep) are broad, cross-file edits where regression is easy; round 02 should re-verify that (a) no `src/agent/` path survives in live docs, (b) README / CLI.md / OPERATOR.md now agree on the shipped set, and (c) the F2 fix did not introduce new broken paths. F5 crosses into code-pass territory (package-skill `ROLES.md` refs) and needs cross-cell coordination, not a same-pass fix. Expect the surface to reach PRISTINE by round 02–03 if α executes F1–F4 cleanly.
