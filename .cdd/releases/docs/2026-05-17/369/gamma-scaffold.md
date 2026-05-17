---
cycle: 369
issue: "https://github.com/usurobor/cnos/issues/369"
date: "2026-05-17"
dispatch_configuration: "claude -p subprocess (sequential single-session); γ acting as δ for boundary actions (Phase 4 δ split not yet landed per #366); permissions bypassed for this run"
base_sha: "704365d23378fcbfcf1e33679025809af6b81100"
parent_issue: 366
predecessor_cycle: 367
scope: "Schema-typing of the frozen receptor interface from #367: three CUE schemas (contract / receipt / boundary_decision) + README + fixture corpus + cycle evidence under .cdd/unreleased/369/"
late_authored: true
late_authored_note: "Scaffold authored on cycle/369 after β R1 D1 (review/SKILL.md rule 3.11b). Recovery path (a). All §Gap claims here re-derive from issue #369 body and peer-enumeration commands; nothing in this file overrides α's review-ready artifacts at 6835197d."
---

# γ Scaffold — #369

## Gap

`schemas/skill.cue` + `schemas/README.md` + `schemas/fixtures/skill-frontmatter/` established the schema-layer pattern at the repo-root `schemas/` surface (CUE owns shape; README owns prose; fixtures live under `schemas/{subsystem}/fixtures/`). `src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md` (landed at #367 merge `37ac1c75`) froze the parent-facing validation interface at doctrine level: named `ValidationVerdict` vs `BoundaryDecision`, named `Override` as a degraded-boundary-action receipt, named role-closeouts as receipt-derivation sources, and explicitly deferred schema syntax to Phase 2.

The receptor is **named** but **not yet typed**. Without typed schemas:

- Phase 3 (`cn-cdd-verify` refactor) has no declarative input/output contract — the validator would re-derive shape mid-cycle, exactly the failure mode #366 exists to prevent.
- Phase 4 (δ split) has no typed `BoundaryDecision` to organize around.
- Phase 6 (ε relocation) has no typed `protocol_gap_count` / `protocol_gap_refs` signal to read across receipt streams; relocation would force a `v2` bump if ε's signal is not pinned at `receipt.v1` now.
- Phase 7 (`CDD.md` rewrite) has no typed surface to reference.

The deep invariant from #366 (carried via #364):

> A closed cell is not trusted because a higher role approved it. A closed cell is trusted because its receipt validates against its contract.

Phase 1 named the receptor. **Phase 2 types it.**

**Empirical anchor (peer enumeration at scaffold time per `gamma/SKILL.md` §2.2a):**

```bash
# Target directory does not yet exist (additive)
test -d schemas/cdd/
# → absent at scaffold time

# No prior file under schemas/ matches the three schema names
ls schemas/ | grep -iE 'contract|receipt|boundary'
# → no matches

# Grep schemas/ for the target schema identifiers
rg -l 'cnos\.cdd\.(contract|receipt|boundary_decision)\.v1' schemas/
# → no matches

# Schema-bearing predecessor pattern present (additive precedent)
test -f schemas/skill.cue && test -d schemas/fixtures/skill-frontmatter/
# → present
```

All four checks confirm: the gap is real and **additive**, not consolidation. The two grep hits in `schemas/skill.cue:14` ("`review/contract` and `review/issue-contract`") and `schemas/README.md` (generic word "contract" in unrelated prose) are not the proposed surface — they reference the existing `review/` skill, not the CDD receptor types.

## Mode

**docs + schema, no code, no validator behavior change.**

Three new CUE schemas + one new README + one fixture corpus under `schemas/cdd/`, plus required cycle evidence under `.cdd/unreleased/369/`. No edits under `src/packages/cnos.cdd/commands/cdd-verify/`, no edits to `operator/SKILL.md`, `gamma/SKILL.md`, `epsilon/SKILL.md`, `ROLES.md`, `CDD.md`, `COHERENCE-CELL.md`, or `RECEIPT-VALIDATION.md` (sole permitted touch on the last: an optional trailing-pointer line, declined by α).

The cycle is single-phase: schema-typing. Disconnect class is **docs-only** (no version tag, no VERSION bump, no RELEASE.md). The merge commit is the disconnect signal, mirroring #367 precedent (commit `37ac1c75`).

## ACs reference

All 10 ACs are defined in the issue body at https://github.com/usurobor/cnos/issues/369 — they are the binding contract for α/β/γ. α maps each AC to evidence in `.cdd/unreleased/369/self-coherence.md §ACs`.

## Acceptance posture summary (γ's compressed reading)

The issue body is authoritative. Notes below are γ's compressed reading for handoff; they do not override the issue.

- **AC1 — three CUE schemas exist + compile.** Files at `schemas/cdd/{contract,receipt,boundary_decision}.cue`. `cue vet -c=false` succeeds (schema-only compile; `-c=false` allows incomplete schema values). No alternative paths.
- **AC2 — schema IDs pinned at v1.** Literal tokens `cnos.cdd.contract.v1`, `cnos.cdd.receipt.v1`, `cnos.cdd.boundary_decision.v1` appear in each schema's header. `v1` is the schema-artifact version; **not** a CDD-protocol version bump.
- **AC3 — recursive scope-lift invariant typed in README + schemas.** `## Scope-Lift Invariant` section in README names all three projections explicitly as projection-under-scope-lift (closed cell at n → α-matter at n+1; δₙ boundary decision → βₙ₊₁-discrimination; εₙ receipt-stream → γₙ₊₁-coordination). Structural typing: `receipt.cue` references `#ValidationVerdict`, `#BoundaryDecision`, `#Transmissibility`; `#Transmissibility` is a CUE-typed function of verdict × action, **not** declarative metadata.
- **AC4 — `#Transmissibility` and `#BoundaryDecision` enforced structurally.** Action enum pinned to exactly five values: `accept | release | reject | repair_dispatch | override` (no `none`, `missing`, `preflight`). `#Override` required iff `action == override` (both directions enforced). Verdict × action × transmissibility mapping table is enforced by CUE, not prose. The two invalid rows (PASS+override, FAIL+accept) must `_|_` at vet time.
- **AC5 — `#Receipt` required-field rule.** `boundary_decision: #BoundaryDecision` required (non-null); `protocol_gap_count: int & >=0` required; `protocol_gap_refs: [...#ProtocolGapRef]` required (`#ProtocolGapRef` is structured: `id`, `source` enum, `ref` — not opaque string); consistency constraint `protocol_gap_count == len(protocol_gap_refs)` structurally enforced; seven evidence refs required (`evidence_root_ref`, `self_coherence_ref`, `beta_review_ref`, `alpha_closeout_ref`, `beta_closeout_ref`, `gamma_closeout_ref`, `diff_ref`), optional `ci_refs: [...string]`.
- **AC6 — one valid receipt fixture vets.** `cue vet -c -d '#Receipt' …  schemas/cdd/fixtures/valid-receipt.yaml` exits 0. Fixture demonstrates PASS / accept / accepted with all required fields populated; passes positively, not by absent-field omission.
- **AC7 — three doctrine-load-bearing invalid fixtures.** Each fails for its named load-bearing reason:
  - `invalid-override-masks-verdict.yaml` — exercises #367 AC6 (override polarity); FAIL+override+`transmissibility: accepted` must fail per AC4 mapping.
  - `invalid-fail-no-boundary-decision.yaml` — exercises #367 AC3 (δ-authoritative-ordering); receipt with FAIL verdict and no `boundary_decision` must fail per AC5 required-field rule.
  - `invalid-gamma-preflight-authoritative.yaml` — exercises #367 AC3 (γ-preflight-non-authoritative); receipt asserting γ-preflight as authoritative and omitting `boundary_decision` must fail per AC5.
- **AC8 — `cn-cdd-verify` behavior unchanged.** `git diff origin/main..HEAD -- src/packages/cnos.cdd/commands/cdd-verify/` empty. Zero tolerance.
- **AC9 — surface containment.** Diff stat limited to `schemas/cdd/` (new directory) + `.cdd/unreleased/369/*.md` (cycle evidence). Forbidden: any new package, any CI workflow, any edit to the listed role skills / doctrine / CDD.md / verifier. The optional trailing-pointer exception on `RECEIPT-VALIDATION.md` is exactly one line and must not re-litigate doctrine.
- **AC10 — `cue vet` invocation pattern documented.** Both invocations (schema-only compile + fixture validation) appear in `schemas/cdd/README.md` in literal multi-line form. Phase 3 inherits these unchanged.

## Skills to load

**Tier 1 (always-on for α/β CDD):**
- `src/packages/cnos.cdd/skills/cdd/CDD.md` — canonical algorithm
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` (α) / `beta/SKILL.md` (β)

**Tier 2 (always-applicable):**
- `src/packages/cnos.eng/skills/eng/*` — markdown authoring
- `src/packages/cnos.core/skills/write/SKILL.md` — every α output is a written artifact

**Tier 3 (issue-specific):**
- `src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` — form authority for any issue-pack reconciliation during cycle
- `src/packages/cnos.core/skills/design/SKILL.md` — single-source-of-truth, declarative-not-implementation, defer implementation decisions
- `src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md` — semantic owner (read-only this cycle; the schemas reference it via CUE header comment + README §Files)
- `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL.md` — doctrine (read-only this cycle; informs §Scope-Lift Invariant)
- `schemas/skill.cue` + `schemas/README.md` — schema-layer pattern precedent

No new schema skill is introduced; `schemas/skill.cue` is the precedent and form authority for CUE-shape patterns.

## Dispatch configuration

**`claude -p` subprocess (sequential single-session).** γ produces α and β prompts; γ routes each role to a fresh `claude -p` invocation sequentially. No nested Agent-tool dispatch. Each role loads its skills fresh from disk on dispatch; spec-staleness propagation is automatic across role boundaries.

`cn dispatch` is not yet implemented (see `gamma/SKILL.md` §2.5); raw `claude -p --dangerously-skip-permissions --allowedTools "Read,Write,Edit,Bash,Glob,Grep"` is used in its place for this cycle.

Dispatch order:
1. γ creates `cycle/369` from `origin/main` at `704365d2`, pushes to origin.
2. γ dispatches α with `/tmp/alpha-369-prompt.txt` (minimal: role + issue + branch + Tier 3 skills + cue installation note).
3. α writes schemas/cdd/* + `self-coherence.md` with review-readiness, exits.
4. γ dispatches β with `/tmp/beta-369-prompt.txt`.
5. β reviews, writes `beta-review.md`, exits. (β does **not** merge in this cycle — γ-acting-as-δ holds the merge boundary because Phase 4 δ split has not landed.)
6. γ merges on β APPROVE; γ writes `gamma-closeout.md`; γ moves `.cdd/unreleased/369/` → `.cdd/releases/docs/2026-05-17/369/`; γ closes #369.

Cycle runs in parallel with `cycle/370` — independent, no coordination.

## Failure modes to guard against (γ-side)

1. **Declarative-not-structural transmissibility.** AC4 is the cycle's load-bearing risk: if `transmissibility` becomes a free string the author sets, the verdict × action mapping is no longer enforced and Phase 3's validator inherits an ambiguous surface. The mapping must be CUE-computed from `validation.verdict` × `boundary_decision.action`, with `_|_` literals on the two invalid combinations.

2. **`#Override` polarity collapse.** `#Override` must be required **iff** `action == "override"`. If only one direction is enforced (e.g. "required when override" but "permitted when non-override"), authors can attach override blocks to accepting boundary decisions and the doctrine line "override is a degraded boundary action that never substitutes for `V = PASS`" leaks.

3. **`boundary_decision` made optional.** AC5 requires it as a non-null field. If left optional, the "no δ decision recorded" and "γ-preflight-only" doctrine-load-bearing invalid cases become structurally valid receipts and the δ-authoritative-ordering freeze from #367 AC3 leaks.

4. **ε signal deferred to v2.** Phase 6 (ε relocation) reads `protocol_gap_count` / `protocol_gap_refs` across receipt streams. If these are not pinned in `receipt.v1` now, ε's move requires a `v2` bump — exactly the failure mode the issue's §Impact enumerates. `#ProtocolGapRef` must be a structured object (`id`, `source` enum, `ref`), not an opaque string.

5. **Evidence refs untyped.** Phase 3's validator needs the input shape declared. The seven required role-artifact refs (`evidence_root_ref`, `self_coherence_ref`, `beta_review_ref`, `alpha_closeout_ref`, `beta_closeout_ref`, `gamma_closeout_ref`, `diff_ref`) must be typed `string` (non-empty) and required at schema level. `ci_refs` is optional.

6. **Recursion as flat role-renaming.** README §Scope-Lift Invariant must frame the three projections as **projection-under-scope-lift**, not as renaming roles inside a single cell. The risk wording: "a `#Receipt` at scope n is also a `#BetaReview` at n+1" — this is *not* what the doctrine says. The correct framing: an entire **closed cell** at scope n appears as α-matter to the parent at scope n+1; the parent runs its own δ/ε surfaces over that.

7. **Surface drift.** `git diff --stat origin/main..HEAD` must show only paths under `schemas/cdd/` and `.cdd/unreleased/369/`. Particularly easy to drift: editing `COHERENCE-CELL.md` "while we're here," or pre-emptively touching `cn-cdd-verify/` to "wire up" the new types. Neither is in this cycle's scope.

8. **`cue vet` invocation drift.** Both invocations (schema-only compile via `-c=false`, fixture validation via `-c -d '#Receipt'` with all three schema files listed) must be documented in literal multi-line form in `schemas/cdd/README.md`. If only the schema-only form ships, or if the fixture form is rendered as a broken single-line string, Phase 3 inherits ambiguity at the worst possible time.

## Branch and sequence

- **Branch:** `cycle/369` (created from `origin/main` at `704365d2`; pushed at scaffold time)
- **α target:** write `schemas/cdd/{contract,receipt,boundary_decision}.cue`, `schemas/cdd/README.md`, `schemas/cdd/fixtures/{valid-receipt,invalid-override-masks-verdict,invalid-fail-no-boundary-decision,invalid-gamma-preflight-authoritative}.yaml`, and `.cdd/unreleased/369/self-coherence.md` on `cycle/369`; signal review-readiness in `self-coherence.md §Review-readiness`.
- **β target:** review on `cycle/369`; write `.cdd/unreleased/369/beta-review.md` with verdict (APPROVE / RC / REJECT). β does **not** merge this cycle (γ holds the merge boundary because Phase 4 δ split has not landed).
- **γ merge target:** on β APPROVE, γ (acting as δ) performs a standard merge commit into `main` (no force push, no fast-forward-only requirement); writes `gamma-closeout.md`; moves `.cdd/unreleased/369/` → `.cdd/releases/docs/2026-05-17/369/`; closes #369 via `gh issue close 369`.

## Process note (late-authored scaffold)

This scaffold was authored **after** β R1 returned RC with binding D1 (`review/SKILL.md` rule 3.11b: missing `gamma-scaffold.md` on the cycle branch). γ should have authored this file **before** dispatching α — that is the canonical order per the §Skills + §Failure modes pattern observed across cycles 357–367. Path (a) recovery is in flight: this file lands now; β R2 re-verifies the single artifact-presence row and re-checks CI on the new HEAD. The cycle is materially review-ready at α's commit `6835197d` per β R1's contract-integrity table (10/11 rows yes; the lone "no" is this artifact). No re-implementation by α is required.

This is a γ-axis process gap and a candidate for PRA cycle-iteration triage at close-out time. The trigger: rule 3.11b discoverability — γ's CDD lifecycle does not block dispatch on the gamma-scaffold.md existence check. Candidate next-MCA: add a pre-dispatch row in `gamma/SKILL.md §2.5` Step 3a/3b or in CDD.md step 3 that mechanically blocks Step 3b on the artifact presence. Disposition to be decided in γ close-out.
