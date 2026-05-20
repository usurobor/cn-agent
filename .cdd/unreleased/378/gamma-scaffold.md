# γ-scaffold — cycle/378

**Issue:** [cnos#378](https://github.com/usurobor/cnos/issues/378) — "cdd: rule 3.11b discoverability under §5.2 wave-mode (wave-manifest as γ-artifact-of-record; cdd-protocol-gap)"

**Wave:** `.cdd/waves/2026-05-19-protocol-hygiene/manifest.md`

**Mode:** §5.2 wave-mode (δ-as-agent for the wave; γ=δ collapse permitted per wave manifest §5.2 + `operator/SKILL.md` §5.2). For this single skill-patch cycle, γ+α+β are further collapsed onto δ — `α ≠ β within a session is structurally compromised` per `operator/SKILL.md` §5.2 scope clause. This is acknowledged in each artifact (`self-coherence.md`, `beta-review.md`) and accepted only because (a) wave manifest precedent permits it for skill-patch class, (b) cycle scope is bounded (three additive skill-text edits, no logic change), and (c) the β-collapsed-on-δ self-review applies the AC1–AC4 oracle mechanically.

**Branch:** `cycle/378` (created from `origin/main` at SHA `c90171537ff5926afe33010d8f87abb49b1b975b`; pushed to `origin/cycle/378`; later rebased onto `origin/main` at `dd5a36d9` mid-cycle when `origin/main` advanced — see `self-coherence.md §Fix-round R1.1` for the rebase trace)

## Surfaces touched

Three SKILL files under `src/packages/cnos.cdd/skills/cdd/`:

1. **`review/SKILL.md`** — rule 3.11b "Exemption discoverability" clause: extend to recognize `.cdd/waves/{wave-id}/manifest.md` as a valid γ-artifact-of-record under §5.2 wave-mode, with the sub-issue → wave-manifest discoverability link (sub-issue body cites wave-id OR master tracking issue links to sub).
2. **`alpha/SKILL.md`** — §2.6 pre-review gate: add a row for γ-side artifact presence at the rule-3.11b literal path (canonical §5.1 `.cdd/unreleased/{N}/gamma-scaffold.md`) OR wave-manifest equivalent (§5.2 `.cdd/waves/{wave-id}/manifest.md`). α pre-empts the β rule-3.11b check.
3. **`operator/SKILL.md`** — §5.2 wave-mode (or §10 Wave Coordination): cross-reference the wave-manifest-as-γ-artifact convention so wave authors know what they're providing satisfies rule 3.11b discoverability.

## AC oracle approach (per AC)

- **AC1 (wave-manifest as γ-artifact-of-record under §5.2):** grep `review/SKILL.md` §3.11b text for `.cdd/waves/{wave-id}/manifest.md` AND `§5.2 wave-mode` (or equivalent phrasing). Discoverability requirement names sub-issue → wave-manifest link.
- **AC2 (α §2.6 row for γ-side artifact presence):** read `alpha/SKILL.md` §2.6; verify a row exists naming the γ-side artifact and rule 3.11b cross-reference; outcome ("present at canonical path" / "wave-manifest serves under §5.2" / "absent — rule 3.11b will fire RC") describable in `self-coherence.md §Review-readiness`.
- **AC3 (empirical anchor cited):** grep the patched skill text for `cph cdr-refactor` AND `2026-05-18` (or `cph#11`, `cph#12`–`#15`); verify the four-sub-uniform + three-distinct-β-substantive-read pattern is named. Per `gamma/SKILL.md` §2.2a precedent — additions to skill rules cite prior cycle evidence.
- **AC4 (no CI/runtime/release surface change):** `git diff origin/main..HEAD --stat` shows changes only under `src/packages/cnos.cdd/skills/cdd/review/SKILL.md`, `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md`, `src/packages/cnos.cdd/skills/cdd/operator/SKILL.md`, and `.cdd/unreleased/378/`. No CI workflow, validator, or CDD doctrine edit.

## Empirical anchor

**Cph cdr-refactor wave 2026-05-18** — master `usurobor/cph#11`; sub-cycles `cph#12`, `#13`, `#14`, `#15` all merged on `cph:main`.

- 4-of-4 sub-uniform §5.2 wave-mode configuration; zero per-sub `gamma-scaffold.md` files written; wave manifest carried γ-artifact-of-record duty.
- 3 distinct β substantive-reads of the same wave-manifest-as-γ-artifact configuration:
  - `cph#12`: `.cdd/unreleased/12/beta-review.md` §3.11b L133–158 — wave-manifest serves γ-artifact role; substantive read.
  - `cph#13`: `.cdd/unreleased/13/beta-review.md` §Contract Integrity row 11 — exemption text `α ≠ β as identities. γ = δ permitted.` in sub-issue body's last line.
  - `cph#14`: `.cdd/unreleased/14/beta-review.md` §2.5 L187–195 — cph-canonical-model reading; three-part justification.
  - `cph#15`: `.cdd/unreleased/15/beta-review.md` §2.0.0 row L29 — uniform §5.2 wave-exempt; β-closeout L55 names this as a wave-protocol patch candidate explicitly.
- α-side anchors: `cph#12/alpha-closeout.md` L38–40; `cph#14/alpha-closeout.md` L46–51.
- Wave-iteration: `usurobor/cph:.cdd/iterations/wave-2026-05-18.md` Finding F1.

## Expected diff scope

- 3 SKILL files modified (additive — no restructure):
  - `src/packages/cnos.cdd/skills/cdd/review/SKILL.md` — rule 3.11b "Exemption discoverability" extended with §5.2 wave-mode recognition; empirical anchor cited.
  - `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — §2.6 row added (γ-side artifact presence, rule 3.11b cross-reference).
  - `src/packages/cnos.cdd/skills/cdd/operator/SKILL.md` — §5.2 (and/or §10) cross-reference to wave-manifest-as-γ-artifact convention.
- 4 cycle artifacts under `.cdd/unreleased/378/`:
  - `gamma-scaffold.md` (this file)
  - `self-coherence.md` (α)
  - `beta-review.md` (β-collapsed-on-δ)
  - `alpha-closeout.md`, `beta-closeout.md`, `gamma-closeout.md`

## Cross-skill coherence as α discipline

The three edits must agree: §2.6 row, §3.11b clause, and §5.2 cross-reference all describe the same wave-mode discoverability path. Drift across files is the α failure mode this issue itself patches against. Internal-consistency check is a β-collapsed-on-δ binding gate.

## Standing wave permissions (per wave manifest)

- Push to cycle branch `cycle/378`: yes
- Push merge to `main` on β close-out: yes
- Auto-dispatch α fix-rounds on RC: yes (max 3)
- Tag/release: NO — wave is docs-only disconnect class; operator gates tag (manifest §Standing permissions)
- Branch delete after merge: yes

## Cross-cycle coordination notes

- Parallel cycles in this wave: #375 (γ-side pre-dispatch gate; touches `gamma/SKILL.md` or `CDD.md`) and #377 (cross-repo protocol; touches `gamma/`, `post-release/`, `cross-repo/`). Neither touches `review/SKILL.md`, `alpha/SKILL.md`, or `operator/SKILL.md` substantively — no expected file-level overlap with #378's surfaces.
- Last-merge-wins on any conflict; γ-as-δ resolves at merge time.
