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

### §Skills

**Tier 1:** `CDD.md` (canonical kernel, loaded and confirmed kernel-only / no touch needed); `alpha/SKILL.md` (this role surface, load order followed: dispatch intake §2.1, produce-in-order §2.2, peer enumeration §2.3, self-coherence §2.5, pre-review gate §2.6, request review §2.7).

**Tier 3 (per dispatch prompt / γ scaffold §5):**
- `src/packages/cnos.core/skills/write/SKILL.md` — skill-prose discipline; kept additions terse, matching the existing row/bullet style rather than inventing new structure.
- `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` — primary AC2 edit target; also read in full for §2.8–§2.11 context around the closure gate.
- `src/packages/cnos.cdd/skills/cdd/beta/SKILL.md` (Pre-merge gate section) — primary AC1 edit target.
- `src/packages/cnos.cds/skills/cds/CDS.md` (Development lifecycle + Gate sections) — primary AC4 edit target; read state-machine table (§"State machine") and Gate → Release-readiness preconditions + Closure verification checklist (F1–F10) in full to place the addition correctly and avoid colliding with the existing F-anchor numbering.
- `src/packages/cnos.cdd/skills/cdd/operator/SKILL.md` (§8 Timeout recovery) — AC3 edit target (doctrinal pointer).
- `src/packages/cnos.cdd/skills/cdd/harness/SKILL.md` (§6 Timeout recovery mechanics) — AC3 edit target (mechanics; γ scaffold flagged this as the concrete home for the actual recovery steps).

**Not loaded (per α load-order discipline):** β/γ role skills were read as edit targets (their content, not their role-frame) — α did not adopt β's or γ's judgment frame, only edited the doctrine they carry, per Rule 3.5.
