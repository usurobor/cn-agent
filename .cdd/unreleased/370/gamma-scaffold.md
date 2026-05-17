---
cycle: 370
issue: "https://github.com/usurobor/cnos/issues/370"
date: "2026-05-17"
dispatch_configuration: "§5.1 canonical multi-session (γ writes prompts; γ-acting-as-δ dispatches α and β sequentially via `claude -p`)"
base_sha: "704365d23378fcbfcf1e33679025809af6b81100"
parent_issue: 366
predecessor_cycle: 367
peer_cycle: 369
scope: "docs-only — author COHERENCE-CELL-NORMAL-FORM.md as the substrate-independent kernel companion to COHERENCE-CELL.md doctrine and RECEIPT-VALIDATION.md design surface"
---

# γ Scaffold — #370

## Gap

`COHERENCE-CELL.md` (#364, merge `32b126e4`) names the recursive coherence-cell *organism*: receipt rule, four-way structural separation, role-as-cell-function. `RECEIPT-VALIDATION.md` (#367, merge `37ac1c75`) names the parent-facing *receptor*: validation interface, verdict/decision distinction, δ-authoritative boundary. #369 (Phase 2, in flight, parallel) types the receptor as CUE schemas. #366's roadmap §"Recursion the implementation must instantiate" carries the algorithm in informal mathematical notation but is a roadmap, not a doctrine surface.

The *algorithm* the doctrine implies — the recursion that takes a closed cell at scope `n` and projects it as α-matter, β-discrimination, and γ-coordination at scope `n+1` — is implicit across these four surfaces but stated nowhere. Phase 3 (validator) would have to derive V's contract from a recursion no doctrine names. Phase 4 (δ split) would intuit δ's signature per cycle. Phase 6 (ε relocation) would move ε without a kernel statement of its signature. Phase 7 (`CDD.md` rewrite) would have to derive the kernel mid-rewrite — exactly the failure mode two-layer separation exists to prevent.

This cycle names the algorithm. It does not rewrite `CDD.md`; it does not shrink `gamma/SKILL.md`; it does not split δ; it does not move ε. It produces the kernel-layer doc that Phases 3, 4, 6, 7 will each cite as their input contract.

**Empirical anchor (peer enumeration at scaffold time per `gamma/SKILL.md` §2.2a):**

```bash
# Target file does not yet exist
test -f src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md
# → ABSENT

# No file in cnos.cdd/skills/cdd/ mentions "normal form" or "CCNF"
rg -i 'normal form|coherence-cell normal|CCNF' src/packages/cnos.cdd/skills/cdd/
# → no matches

# Repo-wide grep for "normal form" or the informal recursion claim
rg -i 'normal form|recursion the implementation must instantiate' src/packages/cnos.cdd/
# → no matches in package skills

# Predecessor + peer surfaces present and unedited at base SHA
ls -la src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL.md \
       src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md \
       src/packages/cnos.cdd/skills/cdd/CDD.md
# → all three present
```

All four checks confirm: the gap is real and additive. The new file lives next to `COHERENCE-CELL.md` and `RECEIPT-VALIDATION.md`; both predecessor surfaces remain unedited. Authoring framing is **additive** (kernel companion), not consolidation.

## Mode

docs-only (single new doctrine file under `src/packages/cnos.cdd/skills/cdd/` + cycle evidence under `.cdd/unreleased/370/`).

The cycle is doctrine-articulation within the docs-only modifier — α produces a kernel-layer doc that names a recursion algorithm. No `.cue`, no code, no role-skill edits, no `CDD.md` modification, no `COHERENCE-CELL.md` modification, no `RECEIPT-VALIDATION.md` modification.

## ACs reference

All 9 ACs are defined in the issue body at https://github.com/usurobor/cnos/issues/370 — they are the binding contract for α/β/γ. α maps each AC to evidence in `.cdd/unreleased/370/self-coherence.md` §ACs.

## Acceptance posture summary (γ pre-flagged to α)

The issue body is authoritative. The notes below are γ's compressed reading of each AC for handoff; they do not override the issue.

- **AC1 — file exists at canonical path.** `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md`. Sibling to `COHERENCE-CELL.md` and `RECEIPT-VALIDATION.md`. No alternative path is acceptable.
- **AC2 — draft-doctrine status + predecessor + peer citations.** Doc opens with explicit clause naming itself as **draft doctrine** companion at the kernel layer (not replacement of `COHERENCE-CELL.md`, not operational doctrine). Required tokens: literal `draft doctrine`, literal `COHERENCE-CELL.md`, literal `RECEIPT-VALIDATION.md`. Doc explicitly states `COHERENCE-CELL.md` remains predecessor doctrine and this doc names the algorithm the doctrine implies.
- **AC3 — kernel as five-step recursion with evidence-binding rule.** A `## Kernel` section states the five-step closure at scope `n`:
  1. `matterₙ := αₙ.produce(contractₙ)`
  2. `reviewₙ := βₙ.review(contractₙ, matterₙ)` — β consumes matter only
  3. `receiptₙ := γₙ.close(contractₙ, matterₙ, reviewₙ)` — γ binds evidence into the receipt at close-out
  4. `verdictₙ := V(contractₙ, receiptₙ, evidenceₙ)` — V dereferences evidence refs from the receipt
  5. `decisionₙ := δₙ.decide(receiptₙ, verdictₙ)` — δ decides on receipt + verdict only; does not re-read evidence
  Plus evidence-binding rule stated explicitly. β's signature must explicitly exclude evidence; δ's signature must explicitly exclude evidence.
- **AC4 — four cell outcomes pinned with verdict × decision preconditions.** A `## Cell Outcomes` section names exactly four: `accepted`, `degraded`, `blocked`, `invalid`. Each with verdict × decision precondition matching #369 AC4. `invalid` is **non-terminal** — δ must re-decide before the cycle terminates. Alignment with #369 AC4 explicitly noted.
- **AC5 — two recursion modes pinned.** A `## Recursion Modes` section names both: **within-scope** (repair-dispatch, same scope index, parent γₙ re-emits after child accept) and **cross-scope** (accept/release with verdict=PASS, or override with verdict≠PASS, scope index advances; closed cell projects as α-matter at `n+1`). Same recursion operator, different scope behavior.
- **AC6 — three scope-lift projections named with projection-not-renaming framing.** A `## Scope-Lift` section names all three: closed (αₙ, βₙ, γₙ) cell → αₙ₊₁ matter; δₙ boundary decision → βₙ₊₁-like discrimination; εₙ receipt-stream observation → γₙ₊₁-like coordination/evolution. Plus explicit framing: "projection under scope-lift, not flat role-renaming inside a single cell; βₙ and γₙ have no upward projection." Alignment with #369 AC3 explicitly noted.
- **AC7 — substrate-independence enforced in kernel sections (section-bounded oracle).** Kernel sections `## Kernel`, `## Cell Outcomes`, `## Recursion Modes`, `## Scope-Lift` use **exactly those headers**. AC7's awk-extract isolates them; word-bounded `rg -i '\b(github|cue|cn-cdd-verify|cn dispatch|claude|gh)\b'` against the extracted slice must return no matches. Substrate references live only in §Realization Layer or §Two-Layer Separation.
- **AC8 — two-layer separation declared with realization peers cited.** A `## Two-Layer Separation` (or equivalent) section names kernel layer (this doc) and realization layer with four explicit citations: `RECEIPT-VALIDATION.md`, `schemas/cdd/` (#369), `cn-cdd-verify` (Phase 3, deferred), `CDD.md` (Phase 7 rewrite, deferred). Framing: "kernel is the *what*; realization is the *how-on-this-substrate*."
- **AC9 — surface containment proved.** `git diff origin/main..HEAD --stat` shows exactly:
  - `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md` (new file, A)
  - `.cdd/unreleased/370/*.md` (cycle evidence — required)
  - **No other files.** Forbidden: any edit to `COHERENCE-CELL.md`, `RECEIPT-VALIDATION.md`, `CDD.md`, `schemas/cdd/`, `cn-cdd-verify`, `operator/SKILL.md`, `gamma/SKILL.md`, `epsilon/SKILL.md`, `ROLES.md`, CI workflows. Doc body restates non-goals.

## Skills to load

**Tier 1 (always-on for α/β CDD):**
- `src/packages/cnos.cdd/skills/cdd/CDD.md` — canonical algorithm
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` (α) / `beta/SKILL.md` (β)

**Tier 2 (always-applicable):**
- `src/packages/cnos.eng/skills/eng/*` — markdown authoring
- `src/packages/cnos.core/skills/write/SKILL.md` — every α output is a written artifact

**Tier 3 (issue-specific):**
- `src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` — form authority for any issue-pack reconciliation during cycle
- `src/packages/cnos.core/skills/design/SKILL.md` — kernel-articulation discipline: substrate-independence, two-layer separation, no-future-as-present, single source of truth, defer realization choices
- `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL.md` — predecessor doctrine; α reads as source-of-truth for receipt rule, four-way separation, role-as-cell-function — the organism this kernel describes algorithmically
- `src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md` — typed-interface design peer; α reads as source-of-truth for V's interface, verdict/decision distinction, δ-authoritative validation — the receptor this kernel positions at step 4–5

`design/SKILL.md` is the load-bearing Tier 3 skill: the work is **kernel articulation** — substrate-independent recursion algorithm that does not collapse into operational specification. Without it α may drift into operational details (citing `claude -p`, `cn-cdd-verify`, `.cue` syntax inside kernel sections) and fail AC7.

## Dispatch configuration

**§5.1 canonical multi-session.** γ produces α and β prompts; γ-acting-as-δ (Phase 4 of #366 has not landed yet) routes each role to a fresh `claude -p` session sequentially. No Agent-tool sub-agent dispatch from a parent session. Each role loads its skills fresh from disk on dispatch.

§5.3 escalation check: AC count = 9 (≥7, multi-session signal). Cycle is docs-only, ≤1 review round target. Proceeding with §5.1.

Dispatch order:
1. γ-acting-as-δ dispatches α with `.cdd/unreleased/370/alpha-codex-prompt.md`
2. α writes `COHERENCE-CELL-NORMAL-FORM.md`, commits to `cycle/370`, writes `self-coherence.md` with review-readiness signal, exits
3. γ-acting-as-δ dispatches β with `.cdd/unreleased/370/beta-codex-prompt.md`
4. β reviews on `cycle/370`; on APPROVE, merges into `main` per `beta/SKILL.md` §1, writes `beta-review.md` + `beta-closeout.md`, exits
5. γ-acting-as-δ runs release-boundary actions (none required for docs-only disconnect — cycle dir move + γ close-out per §3.4)

## Failure modes to guard against (γ-side)

1. **Substrate leakage in kernel sections.** AC7's awk-extract is mechanically strict. If α writes "γ binds evidence into the receipt at close-out via `cn-cdd-verify`-readable typed references" inside `## Kernel`, the rg check finds `cn-cdd-verify` and fails the AC. Substrate names belong only in §Realization Layer or §Two-Layer Separation. Watch for: `github`, `cue`, `cn-cdd-verify`, `cn dispatch`, `claude`, `gh`. The same word inside parentheses or backticks still matches.

2. **Section headers diverging from AC7 awk pattern.** AC7's awk script extracts blocks matching `^## (Kernel|Cell Outcomes|Recursion Modes|Scope-Lift)`. If α writes `## Kernel Algorithm`, `## The Kernel`, or `## Scope-lift` (lowercase l), the awk match silently fails and α may believe AC7 passes. Use **exactly** those four headers, capitalized, hyphenated.

3. **Recursion-mode collapse.** AC5 requires both modes named. The common failure is treating repair-dispatch as scope-advancing (it is not — same scope index `n`, parent γₙ re-emits a fresh `receiptₙ`) or treating accept/release as same-scope (it is not — scope index advances; closed cell becomes α-matter at `n+1`). Watch for: "recursion = scope advance" without distinguishing the within-scope sub-case.

4. **Three projections incomplete or framed as flat renaming.** AC6 requires three projections plus explicit projection-not-renaming framing. The most common drift is naming `δₙ → βₙ₊₁` without the "-like discrimination" qualifier and the "projection under scope-lift" disclaimer. A reader could parse a bare "δ becomes β at the next scope" as flat renaming and reach the wrong conclusion about what survives across scope-lift. The β/γ-have-no-upward-projection clause must be present and unambiguous.

5. **Four-outcome inconsistency with #369 AC4.** AC4 pins exactly four outcomes with preconditions that must match #369's verdict × action × transmissibility table. Particularly: `degraded` requires `verdict ≠ PASS ∧ decision = override` — not "any override". `invalid` is **non-terminal** — δ must re-decide. If α writes `invalid` as a terminal state ("cycle ends invalid"), the kernel diverges from #369 and Phase 2's CUE schemas cannot type-check against this doctrine.

6. **Evidence-binding rule omitted.** AC3 requires the rule stated explicitly. The failure mode is treating γ.close as "γ signs the matter" and V as "V validates the matter directly." This collapses receipt-as-typed-handoff. The rule: γ binds evidence into the receipt at close-out as typed references; V dereferences from the receipt; β never consumes evidence; δ never re-reads evidence. Each signature constraint must be visible in the prose.

7. **Two-layer separation absent or kernel/realization mingled.** AC8 requires the §Two-Layer Separation section to cite four realization peers. The failure mode is to treat the kernel as the operational spec — e.g., letting `cn-cdd-verify` invocation syntax appear in §Kernel "for concreteness." The kernel is the *what*; realization is the *how-on-this-substrate*. Realization examples may appear, but only inside non-kernel sections, and only as labeled examples that do not load-bear the algorithm.

8. **Angle-bracket placeholders rendering hazard.** GitHub's markdown renderer strips `<placeholder>` literals. The issue's active design constraints forbid them. Use `{placeholder}` or literal forms (`matterₙ`, `receiptₙ`).

9. **Operational length creep.** Target length 200–400 lines. Kernel docs are tighter than design docs. If α drifts past 400 lines, it is probably re-deriving `CDD.md` mid-cycle (Phase 7's failure mode this cycle exists to prevent).

## Branch and sequence

- Branch: `cycle/370` (created from `origin/main` at `704365d2`, pushed to origin)
- α target: write `COHERENCE-CELL-NORMAL-FORM.md` + `self-coherence.md` on `cycle/370`; signal review-readiness
- β target: review on `cycle/370`; on APPROVE, merge into `main` per `beta/SKILL.md` §1; write `beta-review.md` + `beta-closeout.md`
- γ close-out target: `gamma-closeout.md` after both α and β close-outs land

## Parallel cycle note

#369 (Phase 2 — schemas) runs in parallel; the two cycles do not block each other. #369 types the receptor's structure (verdict × action × transmissibility table, action enum, override shape, ε signal); #370 names the kernel algorithm. AC4 and AC6 of #370 reference #369's AC4 and AC3 respectively as alignment points — they are doctrine-level references, not blocking dependencies. If #369 evolves its verdict × action table mid-cycle, #370's kernel statement remains the load-bearing claim; #369 must align to the kernel, not the reverse.
