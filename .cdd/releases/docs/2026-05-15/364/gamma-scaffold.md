---
cycle: 364
issue: "https://github.com/usurobor/cnos/issues/364"
date: "2026-05-15"
dispatch_configuration: "§5.2 single-session δ-as-γ with Agent-based α/β"
base_sha: "d412a1e9"
scope: "doctrine articulation — add COHERENCE-CELL.md draft refactor doctrine + README pointer"
---

# γ Scaffold — #364

## Gap

CDD's coherence-cell refactor model (contract → α/β/γ cell → receipt → validator → δ boundary → accepted receipt as next-scope matter) is currently implicit in design discussion and scattered across role, operator, verifier, and foundational-coherence surfaces. Without a repo-local doctrine document, future refactors of `CDD.md`, `operator/SKILL.md`, `gamma/SKILL.md`, `cn-cdd-verify`, and the schema/receipt surfaces risk continuing the same surface smearing the model is meant to eliminate.

The δ/operator surface today fuses boundary policy with dispatch mechanics, polling, git identity setup, release execution, CI recovery, branch cleanup, and timeout recovery — membrane/substrate fusion. The γ surface today fuses cell-closure coordination with runtime supervision idioms (branch preflights, wake mechanics, polling, dispatch detail, monitoring formats).

This issue creates the doctrine surface needed to name and falsify those fusions before implementation refactors begin. It does not perform the splits.

**Empirical anchor (peer enumeration at scaffold time per `gamma/SKILL.md` §2.2a):**

```bash
# Verify COHERENCE-CELL.md does not already exist
test ! -f src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL.md  # passes (file absent)

# Verify no existing reference to coherence-cell doctrine in skills/
rg -l "COHERENCE-CELL|coherence cell" src/packages/cnos.cdd/skills/ 2>/dev/null  # no matches
```

Both confirm: the gap is real and additive, not consolidation.

## Mode

docs-only (substantial CDD cycle, single new artifact + README pointer + cycle evidence)

## ACs reference

All 17 ACs are defined in the issue body at https://github.com/usurobor/cnos/issues/364 — they are the binding contract for α/β/γ. α maps each AC to evidence in `self-coherence.md` §ACs.

## Acceptance posture summary (γ pre-flagged to α)

- **AC1**: explicit Status header naming "Draft refactor doctrine" + Does-Not-Own clause referencing `CDD.md`.
- **AC2**: structural four-way separation (roles / runtime substrate / validation / release/boundary effection). Biological framing allowed but not load-bearing — the document must make the structural prediction that fusing these into one skill is a design smell.
- **AC3**: full recursion equation + `V : Contract × Receipt × EvidenceGraph → ValidationVerdict`.
- **AC4**: `V` as predicate/capability invoked by δ, not ζ-as-agent. Trust = `contract + evidence + valid receipt`, not issued by γ, δ, or ε.
- **AC5**: overrides are degraded boundary actions, must be receipted; not equivalent to validity.
- **AC6**: `cn-cdd-verify` framed as current artifact-presence checker and future contract/receipt validator. Existing per-role closeouts may remain evidence inputs.
- **AC7**: illustrative receipt frontmatter with **unpinned** schema name (`cnos.cdd.receipt.example` or `cnos.cdd.receipt.v0-draft`). MUST NOT pin `v1`. Explicit "illustrative / not normative" label.
- **AC8**: operative β immune/firebreak doctrine. Required phrases: "α≠β is not bureaucracy", "contagion firebreak", "β is the cell's immune discrimination", "without independent β review", "degraded matter".
- **AC9**: review-independence evidence list (actor separation, author separation, β verdict/finding disposition, reviewed artifact refs, verdict-precedes-merge). Validator cannot prove semantic independence completely but can reject obvious counterfeit receipts.
- **AC10**: δ as boundary complex; runtime mechanics (claude -p, cn dispatch, git config/fetch, branch retries, timeout recovery, CI polling, shell mechanics) named **AS EXCLUDED SUBSTRATE** with explicit polarity. Required phrases include "must not contain runtime substrate" or "does not contain runtime substrate" or "outside δ role doctrine" or "outside delta role doctrine", plus "belong below δ" or "belong below delta" or "belong harness".
- **AC11**: γ/δ managerial-residue sweep rule (monitor/supervise/oversee/manage → name produced artifact/receipt/boundary decision; otherwise rename to observe/discriminate/route/validate/close/transport/release/repair-dispatch).
- **AC12**: ε as protocol evolution over accepted receipt streams; not ordinary in-cell metabolism.
- **AC13**: `protocol_gap_count` / `protocol_gap_refs` in every receipt; separate `cdd-iteration.md` required only when `protocol_gap_count > 0`.
- **AC14**: deferred landing order (1. add COHERENCE-CELL.md; 2. add `receipt.cue`/`contract.cue`; 3. refactor `cn-cdd-verify`; 4. split δ; 5. shrink γ; 6. move ε up).
- **AC15**: `src/packages/cnos.cdd/README.md` gains short pointer to the doctrine doc.
- **AC16**: branch diff contains COHERENCE-CELL.md + README pointer + `.cdd/unreleased/364/*` only. Forbidden: `ROLES.md`, `CDD.md`, `operator/SKILL.md`, `gamma/SKILL.md`, `epsilon/SKILL.md`, `activation/SKILL.md`, `cdd-verify/`, `.github/`, `*.cue`.
- **AC17**: `## Open Questions` section with all five questions, **left open** (not resolved): (1) when does V fire? (2) is V capability or command? (3) where does ε relocate? (4) how is an override receipted? (5) do per-role closeouts become evidence-graph inputs to V?

## Skills to load

Tier 1:
- `src/packages/cnos.cdd/skills/cdd/CDD.md`
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md`

Tier 2:
- `src/packages/cnos.eng/skills/eng/*` (always-applicable; markdown authoring)
- `src/packages/cnos.core/skills/write/SKILL.md` (every α output is a written artifact)

Tier 3:
- `src/packages/cnos.core/skills/design/SKILL.md` (boundary/refactor design — keep policy above detail, avoid surface smearing, distinguish role doctrine from runtime substrate, avoid premature canonicalization)

## Dispatch configuration

§5.2 single-session δ-as-γ via Agent tool: γ work executed directly by this session (acceptable per `operator/SKILL.md §5.2`). α and β dispatched via separate Agent tool invocations (subagent_type: claude) to preserve α≠β contagion firebreak. Three execution contexts maintained: this parent (γ), an α sub-agent, a β sub-agent.

Per `operator/SKILL.md §5.3` escalation criteria check: cycle has 17 ACs (≥7 trigger). However the work is docs-only with a single primary artifact and no cross-repo deliverables, expected fix rounds low, expected mid-cycle γ judgment calls low. δ has explicitly authorized §5.2 execution for this cycle. Proceeding with §5.2.

## Failure modes to guard against (γ-side)

1. **Surface smearing.** α must not touch forbidden surfaces. AC16 catches at oracle time.
2. **Doctrine pinning.** α must not write `cnos.cdd.receipt.v1`; use illustrative schema name. AC7 negative oracle catches.
3. **Resolving open questions.** α must seed five questions as inheritance. AC17 catches.
4. **α/β collapse.** γ dispatches α and β as **separate** Agent invocations.
5. **Premature canonicalization.** COHERENCE-CELL.md is draft doctrine; CDD.md remains authoritative.
6. **Metaphor-as-load-bearing.** AC2 checks structural separation, not literal biology wording.
