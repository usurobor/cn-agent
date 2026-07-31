# β Review — Newcomer / Structure Pass (Round 02)

**Reviewer:** β (independent adversarial newcomer review)
**Scope (unchanged):** root `*.md`, `docs/` tree, top-level structure, navigation, links. `src/**` (Go/OCaml), `scripts/`, and package `SKILL.md`/`*.cue` behavior are OUT OF SCOPE → code pass.
**Lens:** `write` skill (signal/noise, one governing question per file, say a fact once, state what it *is*), `document` skill (every claim traced to source), L7 doc bar.
**Method:** Re-derived γ's Round-02 target list from the live tree. Ran every `git grep`/`ls`/`find` myself; confirmed each stale target does/doesn't exist before filing.

---

## Verdict: **REVISE** (one round from PRISTINE on the docs surface)

Round 1 moved the needle. The two Round-1 HIGHs held: **no `src/agent/` regressed** into the live reader surface (every survivor is a dated CHANGELOG/plan/banner record), and README / CLI.md / OPERATOR.md still tell one shipped-vs-planned story. The remaining defects are the residual cluster γ predicted — a stale *folder taxonomy* and a stale *package name* cited across the reference surface — plus one authoritative governance doc that still teaches the abolished scheme as current. All are accuracy drift, not restructuring. None block by themselves; together they still mislead a sharp newcomer on the "how are the docs / packages organized?" path.

## Newcomer impression (has Round 1 moved the needle?)

Yes. The shipped/planned contradiction that broke my trust in Round 1 is gone — CLI.md and OPERATOR.md now agree with README, and the CDD/rules skill paths I chased last round resolve. But when I went one level deeper — "how is this docs tree organized, and what packages exist?" — I hit the same class of drift at the reference layer. `docs/README.md` told me plainly the α/β/γ folders were retired, then pointed me to `DOCUMENTATION-SYSTEM.md` as "the single source of truth for the documentation taxonomy," and that doc describes `docs/alpha/`, `docs/beta/`, `docs/gamma/` bundles as the live system. The flagship architecture doc `CAR.md` (**Status: Implemented**) lists `cnos.pm` as a shipped role pack; it isn't in `src/packages/`. These are the same trust-erosion pattern as Round-1 F1/F3, one layer down. The bones are now solid; the reference surface has residual drift the link gate can't see (bare code-span paths, not `[](…)` links).

---

## Findings (ranked, continuing F-numbering)

### F8 — HIGH — `DOCUMENTATION-SYSTEM.md` teaches the retired α/β/γ folder taxonomy as the live doc system, while cited as the canonical authority

- **Where:** `docs/reference/governance/DOCUMENTATION-SYSTEM.md` (**Version 3.13.0**, no "superseded" marker) — §1.1 "every document has a dominant ontological character (α, β, γ)"; the taxonomy table maps `alpha/` / `beta/` / `gamma/` directories (L22–24); §1.2 "a bundle lives as a subdirectory within its dominant axis (usually `alpha/`)"; the bundle scheme `docs/alpha/{scope}/` (L35, 47, 59–60, 89, 134, 171, 202, 248, 274) — 12 `docs/alpha/` occurrences total.
- **Contradicts, on the newcomer's own path:**
  - `docs/README.md` L10–12: "The α/β/γ triad is no longer a filing taxonomy … never as folders." Then L58–62 routes the reader **to this very doc** "for the documentation-system rules."
  - `docs/reference/governance/README.md` L9/L17/L25 calls it "**the single source of truth** for the documentation taxonomy," "Canonical spec," "**start here** for the full doc system."
  - `docs/reference/governance/GLOSSARY.md` L18: "There is no `docs/alpha/`, `docs/beta/`, or `docs/gamma/`; those directories were retired."
- **Why it violates the bar:** `document` "every claim verified against source of truth" + `write` §3.13 "keep authority explicit" — the designated authority for how the docs work describes a scheme the repo abolished. A newcomer who wants to understand or contribute docs is explicitly sent here and taught the wrong model. This is the F1/F3 README-vs-reference failure re-instanced at the governance layer, and it is the *root* of the whole `docs/alpha/` cluster (12 of the ~21 live hits live here).
- **Fix (α):** This is a **rewrite, not a path-swap** — do not sed `alpha/`→`papers/`. Rewrite §1 (Taxonomy) and §1.2 (bundles) to describe the actual intent-directory scheme (`quickstart/ concepts/ guides/ reference/ architecture/ development/ papers/ evidence/`) that `docs/README.md` and `GLOSSARY.md` already state, keeping the still-valid parts (feature-bundle grouping, version-freeze policy, superseding-notes rule) but re-homed onto real directories. Update the header `Version`/`Date` to reflect the rewrite. One governing question for the file: "how the current docs tree is organized and how documents evolve." When done, re-grep `docs/alpha/` here → 0.

---

### F9 — MED — `docs/alpha/` stale folder-path sweep across the live reference surface (batch for α)

`docs/alpha/` does not exist (`ls` → No such file or directory). Outside F8's `DOCUMENTATION-SYSTEM.md`, the live reader surface still cites it in **7 files / 9 hits**. Every target below was verified present on disk. These are bare code-span paths — invisible to the lychee link gate. Group and batch:

| File:line | Stale cite | Correct target (verified exists) |
|---|---|---|
| `docs/architecture/INVARIANTS.md:6` | `**Canonical-Path:** docs/alpha/architecture/INVARIANTS.md` | `docs/architecture/INVARIANTS.md` (self-referential dead path — highest signal in this batch) |
| `docs/architecture/HUB-PLACEMENT-MODELS.md:468` | `docs/alpha/HUB-PLACEMENT-MODELS.md` | `docs/architecture/HUB-PLACEMENT-MODELS.md` (self-reference) |
| `docs/papers/COHERENCE-SYSTEM.md:601` | `docs/alpha/doctrine/FOUNDATIONS.md` | `docs/papers/FOUNDATIONS.md` |
| `docs/reference/ctb/CTB-v4.0.0-VISION.md:352` | `docs/alpha/doctrine/SKILL-ARCHITECTURE.md` | `docs/papers/SKILL-ARCHITECTURE.md` |
| `docs/reference/ctb/CTB-v4.0.0-VISION.md:377` | `docs/alpha/doctrine/COHERENCE-FOR-AGENTS.md` | `docs/concepts/doctrine/coherence-for-agents/COHERENCE-FOR-AGENTS.md` |
| `docs/reference/ctb/LANGUAGE-SPEC.md:305,511` | `docs/alpha/doctrine/SKILL-ARCHITECTURE.md` | `docs/papers/SKILL-ARCHITECTURE.md` |
| `docs/reference/ctb/SEMANTICS-NOTES.md:50` | `docs/alpha/doctrine/SKILL-ARCHITECTURE.md` | `docs/papers/SKILL-ARCHITECTURE.md` |
| `docs/development/design/README.md:25` | `docs/alpha/essays/` | `docs/papers/` (foundational position papers now live there) |

- **Why it violates the bar:** `document` §"no stale references to removed features"; a newcomer chasing the CTB architectural-argument pointers or the INVARIANTS canonical-path lands on a dead folder. Same class as Round-1 F2/F3, at the reference layer.
- **Fix (α):** Apply the table above verbatim — no guessing required. The two self-referential canonical-path headers (INVARIANTS, HUB-PLACEMENT) are the sharpest (a doc naming its own dead location); fix those first.
- **Judgment calls left for α (do NOT blind-swap):**
  - `docs/papers/EXECUTABLE-SKILLS.md:398` — a bibliography citation `[1] … docs/alpha/protocol/WHITEPAPER.md`. No `WHITEPAPER.md` exists anywhere in the repo. This is a *published paper's reference list*, not a live nav link → treat as frozen (leave), or resolve only if α can locate the renamed target. Flagged, not filed for swap.

---

### F10 — MED — `cnos.pm` nonexistent-package sweep across the live reference surface (batch for α)

`src/packages/cnos.pm` does not exist. The real set is `cnos.core, cnos.cdd, cnos.cdd.kata, cnos.cdr, cnos.cds, cnos.eng, cnos.handoff, cnos.issues, cnos.kata` (README package table, L144–151; `ls src/packages/`). `cnos.pm` was a v3.9-era PM package (CHANGELOG L900/907) that no longer ships. Live reader-surface hits:

| File:line | Context | Fix direction |
|---|---|---|
| `docs/architecture/cognitive-substrate/CAR.md:41,105` | **Status: Implemented** doc lists `cnos.pm` as a role pack + shows a `cnos.pm/` tree with `pm/follow-up`, `pm/issue`, `pm/ship` skills | Replace the illustrative role-pack with a package that ships (e.g. `cnos.eng` or `cnos.cds`), or point to README's package table for the live set. An "Implemented" doc must not show a phantom package. |
| `docs/papers/COHERENCE-SYSTEM.md:266` | §5.5 "Package articulations" lists `cnos.core / cnos.eng / cnos.pm / future org packs` as fact | Drop `cnos.pm` or replace with a real role pack; the "future org packs" line already covers the open-ended case. |
| `docs/reference/runtime/ORCHESTRATORS.md:183,191,493` | JSON examples with `"package": "cnos.pm"` | Repoint the example package to one that ships (`cnos.cds`/`cnos.eng`). Runtime-reference examples that name a nonexistent package read as untested. |

- **Special case (do NOT swap):** `docs/reference/packages/PACKAGE-RESTRUCTURING.md` (5 hits) is a **Status: Draft** design doc (#186) that legitimately *proposes retiring* `cnos.pm`. Its `cnos.pm` mentions are intentional. BUT its "Target Structure" table proposes `cnos.kernel`, `cnos.agent`, `cnos.hub` — none of which shipped — so the whole doc describes an unrealized plan while sitting under `reference/packages/` (which reads as authoritative). Recommend α add a one-line status banner ("Draft proposal, not the shipped structure; see README package table for the live set") rather than editing the body. Flagged as a fresh newcomer trap, filed as sub-item not a swap.
- **Why it violates the bar:** `document` §"drift"; `write` §3.13 — a newcomer cannot tell whether `cnos.pm` is a real package. Same class as Round-1 F3's `cnos.pm` line, now swept across the surface.

---

### F11 — LOW / POLICY — ratify "leave frozen historical records unannotated" (γ's open policy question)

γ (and α residual #1) asked whether frozen records (CHANGELOG, cycle logs, `plans/*`, dated design/paper records) should each get a one-line "paths historical" note. **Recommendation: LEAVE them; add no per-file note.** Rationale (write skill §3.14 "brevity is earned," §3.3 "say a fact once"): annotating dozens of dated records adds a restated caveat to every one for zero newcomer benefit — a `Status: Proposed`/`Superseded`/dated header already signals "this is history." The lower-noise coherent rule:

> A doc is "frozen" (leave it) iff it carries a dated/superseded/status header **and** is not linked from the live nav as current authority. Only the **live reader surface** gets corrected.

This rule is what separates F8 from the leave-list: `DOCUMENTATION-SYSTEM.md` fails the second clause (it *is* cited as current authority), so it is live and must be fixed — it is not "frozen." Records that pass both clauses and stay untouched: `CHANGELOG.md`, `docs/concepts/doctrine/**/JFA-cycle-log-dyad.md`, `docs/development/plans/*`, `docs/papers/MODULAR-ARCHITECTURE-REFACTOR.md` (Status: Proposed), `docs/reference/runtime/PLAN-174-orchestrator-runtime.md`, `docs/reference/packages/DESIGN-266-dist-out-of-git.md` (design record citing a line number in a frozen INVARIANTS), `docs/evidence/design/WRITER-PACKAGE.md` (header: "Superseded by post-SOUL.md redo … remains as cycle-input evidence" — this is why its 25 `docs/alpha/` hits and its `cnos.writer` references are left alone). The two Round-1 banner docs (`BUILD-AND-DIST.md`, `CORE-REFACTOR.md`) read correctly under their completion banners — confirmed.

---

## Regression re-verify (Round-1 F1/F2)

- **`src/agent/` in live docs:** `git grep -n 'src/agent/' -- docs/ '*.md' ':!.cdd/' | grep -v 'src/packages/'` → every hit is CHANGELOG / `plans/*` / `papers/MODULAR` / `PLAN-174` / `JFA-cycle-log` / `BUILD-AND-DIST` + `CORE-REFACTOR` banner bodies. **Zero new stale `src/agent/` on the live surface. NO REGRESSION.**
- **Shipped-vs-planned:** README / CLI.md / OPERATOR.md unchanged since Round 1; still agree. **NO REGRESSION.**
- **No new broken paths** introduced by the Round-1 α sweep (spot-checked the F2/F3 targets; all resolve).

**Round-1 regressed? NO.**

---

## Converge-vs-defer split (per the brief)

- **Converges the surface this round (docs-only, α can execute):** F8 (rewrite DOCUMENTATION-SYSTEM.md), F9 (8-row path table), F10 (3 files + 1 banner), F11 (policy ratification — a decision, no edit).
- **Defer → code pass (out of docs scope):** F5 `ROLES.md` move + 121 `src/packages/**` ref-rewrites (γ carried); the 8 `src/agent/` hits inside `src/packages/**/{README,SKILL}.md`; `cnos.issues` absent from README's package table (a table-vs-tree check that touches package identity — verify in code pass); `CLI.md`/`ARCHITECTURE.md` OCaml module-name surface (Round-1 α residual #5).

---

## Convergence assessment

- **Counts:** 1 HIGH (F8), 2 MED (F9, F10), 1 LOW/policy (F11).
- **Distance from pristine:** **One round.** Only F8 requires prose judgment (a bounded rewrite of one file); F9/F10 are enumerated, target-verified swaps a competent α lands in a single pass. With F8+F9+F10 done and F11 ratified, the docs reader surface reaches **PRISTINE**. The only cross-boundary item (F5) is legitimately the code pass's, not this surface's — so the docs surface can converge without waiting on it.
- **Another β round?** One verification round (γ or β) to confirm F8's rewrite is coherent and the two sweeps hit zero, then close. Trajectory matches γ's Round-01 forecast: pristine by round 02–03.
