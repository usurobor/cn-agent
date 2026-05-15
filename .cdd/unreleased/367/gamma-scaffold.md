---
cycle: 367
issue: "https://github.com/usurobor/cnos/issues/367"
date: "2026-05-15"
dispatch_configuration: "§5.1 canonical multi-session (γ writes prompts; δ dispatches α and β sequentially via claude -p)"
base_sha: "ffdd77acab2fcfa7670b4c2d77f1dc305fcff76b"
parent_issue: 366
predecessor_cycle: 364
scope: "design-only cycle — author RECEIPT-VALIDATION.md as the parent-facing validation receptor surface for COHERENCE-CELL.md doctrine"
---

# γ Scaffold — #367

## Gap

`COHERENCE-CELL.md` (landed at #364 merge `32b126e4`) names the recursive coherence-cell model and declares `V : Contract × Receipt × EvidenceGraph → ValidationVerdict` as the load-bearing predicate of the doctrine. The doctrine names the organism but does not specify the parent-facing receptor. Five Open Questions seeded at AC17 of #364 carry forward the design tension:

- Q1 — when does `V` fire?
- Q2 — capability or command?
- Q3 — where does ε relocate?
- Q4 — override receipt shape?
- Q5 — role closeouts as evidence-graph inputs or parent-facing artifacts?

Without resolving these five questions and freezing the validation interface (input shape, output shape, invocation contract), the remaining phases of #366 are speculative:

- Phase 2 (`.cue` schemas) — receipt fields designed without their consumer's interface
- Phase 3 (`cn-cdd-verify` refactor) — refactor target is an undefined predicate
- Phase 4 (δ split) — no receptor to thin the membrane around
- Phase 5 (γ shrink) — leans on Phase 4
- Phase 6 (ε relocation) — partially blocked because ε's protocol-iteration role intersects evidence-graph inputs
- Phase 7 (`CDD.md` rewrite) — remains aspirational without the receptor

The deep invariant from #366 (carried via #364):

> A closed cell is not trusted because a higher role approved it. A closed cell is trusted because its receipt validates against its contract.

This cycle names the receptor. It does not thin the membrane. The membrane (δ split etc.) is Phase 4 work.

**Empirical anchor (peer enumeration at scaffold time per `gamma/SKILL.md` §2.2a):**

```bash
# Target file does not yet exist (additive)
test -f src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md
# → absent

# No prior file under cnos.cdd/skills/cdd/ matches receipt|validation
ls src/packages/cnos.cdd/skills/cdd/ | grep -iE 'receipt|valid'
# → no matches

# Repo-wide grep for the target filename
rg -l 'RECEIPT-VALIDATION' src/packages/cnos.cdd/
# → no matches

# Validation-interface terms occur exclusively in predecessor doctrine
rg -l 'Validation Interface|ValidationVerdict|validation interface|validator interface' src/packages/cnos.cdd/
# → src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL.md
```

All four checks confirm: the gap is real and additive. `COHERENCE-CELL.md` is the predecessor doctrine surface and the only file naming `ValidationVerdict`. Authoring framing is **additive**, not consolidation. The new file lives next to the doctrine; the doctrine remains unedited.

## Mode

docs-only (single new design artifact under `src/packages/cnos.cdd/skills/cdd/` + cycle evidence under `.cdd/unreleased/367/`).

The cycle is design-only within the docs-only modifier — α produces a design surface, not implementation. No `.cue`, no code, no role-skill edits, no `CDD.md` modification.

## ACs reference

All 9 ACs are defined in the issue body at https://github.com/usurobor/cnos/issues/367 — they are the binding contract for α/β/γ. α maps each AC to evidence in `.cdd/unreleased/367/self-coherence.md` §ACs.

## Acceptance posture summary (γ pre-flagged to α)

The issue body is authoritative. The notes below are γ's compressed reading of each AC for handoff; they do not override the issue.

- **AC1 — file existence at canonical path.** `src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md`. Sibling to `COHERENCE-CELL.md` and `CDD.md`. No alternative path is acceptable.
- **AC2 — draft-design status + predecessor citation.** Doc opens with an explicit clause naming itself as **draft design** (not doctrine, not implementation spec), citing `COHERENCE-CELL.md` as predecessor, and stating the binding-under-Phase-3 condition. Required tokens: literal `draft design`, literal `COHERENCE-CELL.md`, plus an explicit phrase tying the doc's binding-status to Phase 3 implementation of `V`.
- **AC3 — Q1 (when does V fire?).** Single authoritative firing point: **δ-boundary validation** after γ emits the receipt and before parent acceptance/release. γ-preflight may be allowed but must be **explicitly non-authoritative**. Required: ordering rule that parent scope cannot treat the cell as one until δ-boundary invocation produces a verdict and δ records the boundary decision. Avoid "pre-merge OR post-merge" framings that leave authority diffuse.
- **AC4 — Q2 (capability or command?).** Exactly one chosen answer: predicate / capability in package manifest / separate command under `cnos.cdd/commands/` / registry provider. Minimal invocation sketch. Do not enumerate options without choosing.
- **AC5 — Q3 (ε relocation target).** Single named target: `ROLES.md` / `cnos.core/doctrine/` / new `cnos.protocol-iteration` / other-named. Doc explicitly states the **move is Phase 6 work** (this cycle names the target; it does not move ε).
- **AC6 — Q4 (override receipt shape).** Override is a **degraded boundary action** that **must be receipted** and **never substitutes for `V = PASS`** in downstream consumers. Choose exactly one shape: structured `override:` block / sibling artifact / degraded-state flag with metadata. Required fields: actor, justification, original validation verdict, failed predicates being overridden, degraded-state flag, downstream-consumer detection rule. The override never rewrites or masks the original `ValidationVerdict`.
- **AC7 — Q5 (role closeouts as evidence-graph inputs).** Each of the five files (`self-coherence.md`, `beta-review.md`, `alpha-closeout.md`, `beta-closeout.md`, `gamma-closeout.md`) gets a named role: evidence-graph input / receipt-derivation source / deprecated / unchanged. Receipt artifact identified as the parent-facing surface. Derivation rule: receipt is computed at γ close-out from the five files.
- **AC8 — Validation interface frozen.** A `## Validation Interface` section names input contract, output contract, invocation contract — all three at doctrine level. Mandatory verdict/decision distinction:
  - `ValidationVerdict`: emitted by `V`; describes whether the receipt satisfies the contract/evidence predicates.
  - `BoundaryDecision`: emitted by δ; records accept / reject / repair-dispatch / override / release decision.
  - Rule: `ValidationVerdict = PASS` enables normal δ acceptance; δ override is possible only as degraded boundary action and never rewrites the validation verdict.
  - Minimal illustrative example — prose plus example block; **NOT `.cue` syntax**; explicit "schema syntax is Phase 2" disclaimer.
- **AC9 — Surface containment.** `git diff origin/main..HEAD --stat` shows exactly:
  - `src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md` (new file, status `A`)
  - `.cdd/unreleased/367/*.md` (cycle evidence — required)
  - **No other files.** Forbidden: any `*.cue`; any edit under `commands/cdd-verify/`, `operator/SKILL.md`, `gamma/SKILL.md`, `epsilon/SKILL.md`, `CDD.md`, `ROLES.md`, `COHERENCE-CELL.md`. Doc body also restates non-goals.

## Skills to load

**Tier 1 (always-on for α/β CDD):**
- `src/packages/cnos.cdd/skills/cdd/CDD.md` — canonical algorithm
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` (α) / `beta/SKILL.md` (β)

**Tier 2 (always-applicable):**
- `src/packages/cnos.eng/skills/eng/*` — markdown authoring
- `src/packages/cnos.core/skills/write/SKILL.md` — every α output is a written artifact

**Tier 3 (issue-specific):**
- `src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` — form authority for any issue-pack reconciliation during cycle
- `src/packages/cnos.core/skills/design/SKILL.md` — interface-freeze discipline: single source of truth, no-future-as-present, explicit deferral of implementation choices, distinguish doctrine from spec
- `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL.md` — predecessor doctrine; α reads as source-of-truth for `V`, receipt rule, and the four-way separation that this design must extend coherently

`design/SKILL.md` is the load-bearing Tier 3 skill: the work is to **freeze an interface at doctrine level**, which is precisely the discipline `design` encodes (interface > implementation, single source of truth, defer implementation choices). Without it α may drift into schema authorship.

## Dispatch configuration

**§5.1 canonical multi-session.** γ produces α and β prompts; δ routes each role to a fresh `claude -p` session sequentially. No Agent-tool sub-agent dispatch from a parent session. Each role loads its skills fresh from disk on dispatch; spec-staleness propagation is automatic across role boundaries.

§5.3 escalation check: AC count = 9 (≥7, multi-session signal). Cycle has no cross-repo deliverables and no fix-rounds expected (docs-only convention is ≤1 review round). The ≥7 AC signal alone is sufficient to nominate §5.1. δ's dispatch protocol already maps to §5.1 ("δ will dispatch α and β sequentially"); proceeding without operator override.

Dispatch order (δ executes):
1. δ dispatches α with `.cdd/unreleased/367/alpha-codex-prompt.md`
2. α writes `RECEIPT-VALIDATION.md`, commits to `cycle/367`, writes `self-coherence.md` with review-readiness signal, exits
3. δ dispatches β with `.cdd/unreleased/367/beta-codex-prompt.md`
4. β reviews, writes `beta-review.md`, merges to main (or returns RC for fix-round), writes `beta-closeout.md`, exits
5. δ holds release-boundary gates per operator authorization scope

## Failure modes to guard against (γ-side)

1. **Resolving open questions without choosing.** α must not list options for any of Q1–Q5 without picking one. AC3–AC7 negative oracle is "Question described but unanswered" / "multiple answers without choosing" / "TBD". α should write each `## Q{n}` section with a single chosen answer at the head and rationale beneath.

2. **Schema-pinning the interface.** AC8 forbids `.cue` syntax in illustrative examples. The interface must be prose plus a minimal example block (yaml-ish, json-ish, or prose-tabular is fine). The disclaimer "schema syntax is Phase 2" must appear in §Validation Interface.

3. **Verdict/decision collapse.** A common failure mode is to fuse `ValidationVerdict` (V's output) with `BoundaryDecision` (δ's act). AC8 explicitly requires both names and the rule that an override never rewrites the verdict. If α treats override as "V emits OVERRIDE-PASS" the cycle has failed.

4. **Surface drift.** AC9 catches any file outside the cycle envelope at `git diff --stat`. Forbidden surfaces enumerated above. Particularly easy to drift: editing `COHERENCE-CELL.md` to "fix a small thing while we're here". Don't. The predecessor doctrine is frozen by this cycle's contract.

5. **γ-preflight as authority substitute.** Q1 is delicate: γ-preflight may exist, but the doctrine line is δ-boundary authority. If α writes Q1 such that γ-preflight is "sufficient when V = PASS", the cycle has weakened the receptor. The receptor lives at the δ boundary; γ may invoke `V` as informational, not as gate.

6. **ε relocation overreach.** Q3 names the target. Phase 6 ships the move. α must not produce a relocation patch, just name the target with rationale.

7. **Override pathway substituting for PASS.** Q4 is delicate: an override is a degraded action. The override receipt must be detectable by downstream consumers — i.e. parent scope must be able to tell from the receipt that the cell closed under override and not under `V = PASS`. AC6 requires the **downstream-consumer detection rule** explicitly.

8. **AC granularity vs document structure (#364 PRA observation, carried forward).** AC3–AC7 are five mechanically uniform "answer Q{n}" ACs. Resist writing them as a checklist with one-line answers; the design surface is the prose that justifies the choice. Embed required answer tokens in operative reasoning, not in checkboxes.

## Branch and sequence

- Branch: `cycle/367` (created from `origin/main` at `ffdd77ac`)
- α target: write `RECEIPT-VALIDATION.md` + `self-coherence.md` on `cycle/367`; signal review-readiness
- β target: review on `cycle/367`; on APPROVE, merge into `main` per `beta/SKILL.md` §1; write `beta-review.md` + `beta-closeout.md`
- γ close-out target: `gamma-closeout.md` after both α and β close-outs land
