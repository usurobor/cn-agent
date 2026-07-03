# self-coherence — cycle/368

## §R0

### §Gap

**Issue:** [cnos#368](https://github.com/usurobor/cnos/issues/368) — bug(cdd): cycle close-out can leave issue OPEN — Closes-keyword inconsistently applied, γ fallback skipped on docs-only disconnect.

**Mode:** design-and-build (per γ scaffold — design ~80% converged; exact insertion points confirmed by γ peer-enumeration before α execution).

**Priority:** reconciled to `P2` this cycle (label was `P1`; issue body stated `Priority: P2` with rationale "issue-state drift; workaround is trivial; not release-blocking"). GitHub label corrected via `gh issue edit 368 --repo usurobor/cnos --remove-label P1 --add-label P2`. Rationale: the issue body's own stated rationale is sound — a docs-only close-token gap with a one-line manual-close workaround is not release-blocking, even though it recurs.

**Version:** doctrine/skill-prose cycle — no code, no version bump. Implementation-contract axes pinned by γ (scaffold §5) and δ (`gamma-clarification.md`): Markdown/prose-only, no new CLI subcommand or script, no new binaries, no wire-format change (additive artifact field only).

**Corrected surfaces (per γ scaffold §2/§3, binding over the issue's own citations):**
- AC1 target: `beta/SKILL.md` §"Pre-merge gate" (not `review/SKILL.md` — that section doesn't exist there).
- AC2 target: `gamma/SKILL.md` §2.10 (confirmed correct by γ).
- AC3 target: `operator/SKILL.md` §8 + `harness/SKILL.md` §6 (concrete home found by γ; issue left this under-specified).
- AC4 target: `cnos.cds/skills/cds/CDS.md` (not `CDD.md` — the instantiation-level close protocol migrated out of `CDD.md` in the cnos#366 Phase 7 rewrite; `CDD.md` is kernel-only and its own non-goals forbid re-deriving this doctrine there).

Branch `cycle/368` created by γ from `origin/main` at `d5bb2c20958998e236ab7c0d0a154ddc9ee319f2`. α picked up at HEAD `6ccc2f7` (δ's implementation-contract pin) and implemented on top.
