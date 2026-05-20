# γ scaffold — cycle #375

## Issue

usurobor/cnos#375 — "γ-side pre-dispatch gate for gamma-scaffold.md (rule 3.11b symmetry; cdd-protocol-gap)".

Class: skill-patch (docs-only disconnect). Wave member: 2026-05-19 protocol hygiene (see `.cdd/waves/2026-05-19-protocol-hygiene/manifest.md`).

## Mode

§5.2 wave-mode with γ=δ collapse permitted (per `operator/SKILL.md` §5.2 + `cdd/CDD.md` §1.4). This δ-as-agent invocation plays γ+α+β-collapsed-on-δ for the cycle. α≠β within a session is structurally compromised but acceptable for the skill-patch class per the wave manifest precedent.

## Surfaces touched

Canonical surface (one chosen, not both, per AC1 oracle):

- `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` §2.5 Step 3b "Subscribe and dispatch α and β" — pre-dispatch γ-scaffold existence gate inserted as the first sub-section before any prompt is produced.

Rationale for choice: §2.5 Step 3b is the canonical "dispatch α" surface; the gate must bind to *dispatch*, not to branch creation (Step 3a) or to general CDD doctrine (`CDD.md` §1.4 step 3 is too coarse — it doesn't mention `gamma-scaffold.md`). Putting the gate in §2.5 Step 3b co-locates it with the existing dispatch flow and the α/β prompt templates that follow.

## Empirical anchor

Cycle #369 β R1 D1 — rule 3.11b binding fire on missing `gamma-scaffold.md`. Recovery path (a) at `227d2373` (γ-369: scaffold — recovery path (a) for β R1 D1) → β R2 APPROVED at `4e179db6` → merged at `ff54f2a0` (γ-369: merge cycle/369 — Phase 2 schema-typing of receptor).

Verified locally:
- `git log --oneline ff54f2a0` → "γ-369: merge cycle/369 — Phase 2 schema-typing of receptor (CUE schemas + fixtures)"
- `git log --oneline 227d2373` → "γ-369: scaffold — recovery path (a) for β R1 D1 (rule 3.11b)"

Both SHAs present in repo history; empirical anchor satisfiable.

## AC oracle approach (paraphrased from issue body, AC1–AC4)

- **AC1 — Pre-dispatch check named in canonical surface as binding gate.** Oracle: `gamma/SKILL.md` §2.5 Step 3b contains a sub-section that names "`.cdd/unreleased/{N}/gamma-scaffold.md`" as a pre-dispatch existence check, with binding gate language ("γ MUST NOT proceed to Step 3b α dispatch until …") rather than advisory prose.
- **AC2 — Symmetry with rule 3.11b documented.** Oracle: the new sub-section explicitly references `review/SKILL.md` rule 3.11b and frames the γ-side check as its dual (γ pays cost once at scaffold-time; β pays cost once-per-RC-round when γ forgets).
- **AC3 — Empirical anchor cited.** Oracle: the new sub-section names cycle #369, the β R1 D1 finding, the recovery-path commit `227d2373`, β R2 APPROVED, and the merge `ff54f2a0` (verbatim SHAs).
- **AC4 — No CI / runtime / release surface change.** Oracle: `git diff origin/main..HEAD --stat` after the cycle shows changes only under `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` and `.cdd/unreleased/375/`. No edits to `CDD.md`, no edits to `review/SKILL.md`, no CI workflows, no validator code, no release scripts.

## Expected diff scope

- `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` — one new sub-section inserted into §2.5 Step 3b ("Pre-dispatch γ scaffold check"), section-level addition, ~15–25 lines, no surrounding restructure.
- `.cdd/unreleased/375/gamma-scaffold.md` — this file.
- `.cdd/unreleased/375/self-coherence.md` — α artifact (§Gap, §AC mapping, §CDD Trace, §Review-readiness signal).
- `.cdd/unreleased/375/beta-review.md` — β verdict (APPROVE or RC with rounds).
- `.cdd/unreleased/375/alpha-closeout.md`, `beta-closeout.md`, `gamma-closeout.md` — close-outs per `cdd/CDD.md` §1.4.

No file outside the above set is touched by this cycle.

## Non-goals (per issue body)

- Do NOT modify `review/SKILL.md` rule 3.11b (already exists, already binding on β-side).
- Do NOT relocate `gamma-scaffold.md` outside `.cdd/unreleased/{N}/`.
- Do NOT introduce a new validator surface; the check is a SKILL-level gate, not a `cn-cdd-verify` check.
- Do NOT require an entire new γ pre-flight phase; the check fits inside Step 3b.
- Do NOT edit `CDD.md` doctrine more broadly than naming the gate (and in fact this cycle does not edit `CDD.md` at all — `gamma/SKILL.md` is the canonical surface choice).

## Wave-level cross-cycle note

Per the wave manifest, cycle #377 may also touch `gamma/SKILL.md` (§2.1 + §2.7 for cross-repo coordination). This cycle's edit is scoped to §2.5 Step 3b. If #377 merges first, conflict resolution at merge time integrates both patches (different sections, no overlap expected).
