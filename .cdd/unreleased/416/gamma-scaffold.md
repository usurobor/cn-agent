# γ-scaffold — Cycle 416

**Issue:** [cnos#416](https://github.com/usurobor/cnos/issues/416) — Sub 2 of [cnos#404](https://github.com/usurobor/cnos/issues/404): Move `cross-repo/SKILL.md` (644 lines) into `cnos.handoff` — verbatim move + compatibility stub + cross-ref repair.
**Branch:** `cycle/416` from `origin/main @ 92a7442d`.
**Mode:** Collapsed γ+α+β on δ (per dispatch metadata).
**Class:** docs-only / package migration — **no behavioral redesign**.

## Pinned implementation contract (δ; see issue #416)

| Axis | Pinned value |
|---|---|
| Language | Markdown |
| CLI integration target | None new |
| Package scoping | New files at `cnos.handoff/skills/handoff/{HANDOFF.md, cross-repo/SKILL.md}`; replace `cnos.cdd/skills/cdd/cross-repo/SKILL.md` with compatibility stub; edit citations in `cnos.cdd/`, `cnos.cdr/`, `cnos.cds/`; update `cnos.handoff/docs/extraction-map.md` Sub 2 rows |
| Existing-binary disposition | N/A |
| Runtime dependencies | None |
| JSON/wire contract | N/A |
| Backward compat | Compatibility stub at old path preserves all unswept citations |

## Deliverables (D1–D5; full detail in issue body)

- **D1**: `cnos.cdd/skills/cdd/cross-repo/SKILL.md` → `cnos.handoff/skills/handoff/cross-repo/SKILL.md` (verbatim except: frontmatter `parent: cdd` → `parent: handoff`; `requires:` reworded to be package-agnostic; `calls:` paths get `cnos.cdd/skills/cdd/` prefix; STATUS-canonical-home declaration flipped from "canonical in `CDD.md`" → "canonical in this skill (§2.3)"; internal cross-references preserved as consumer references).
- **D2**: Replace old path with ≤ 50-line compatibility stub (verbatim content from issue body).
- **D3**: Author minimal `cnos.handoff/skills/handoff/HANDOFF.md` (50–150 lines; **NOT** a synthesis doc; verbatim shape from issue body). Replace `.gitkeep`.
- **D4**: Cross-reference repair in `cnos.cdd/skills/cdd/{gamma,post-release,epsilon,issue,operator,activation}/SKILL.md`, `CDD.md`; `cnos.cds/skills/cds/CDS.md`; `cnos.cdr/skills/cdr/`. Re-point any citation framing the old path as canonical authority → `cnos.handoff/skills/handoff/cross-repo/SKILL.md`. Leave historical / "current home" / migration-prose citations as-is (the compatibility stub catches them).
- **D5**: Update `cnos.handoff/docs/extraction-map.md` §1 rows to `Status: v0.1 migrated; canonical at cnos.handoff/skills/handoff/cross-repo/SKILL.md`.

## AC matrix (cnos#416 AC1–AC10)

| AC | Check | Plan |
|---|---|---|
| AC1 | 3 files exist (HANDOFF.md; new cross-repo/SKILL.md; stub at old path) | `test -f` each after α |
| AC2 | line count ≥ 600 + 6 keyword hits on `cross-repo state machine\|LINEAGE\|STATUS\|feedback patch\|bundle archival\|hat-collapse` | `wc -l` + `rg` |
| AC3 | stub ≤ 50 lines + ≥ 3 hits on `cnos.handoff\|handoff/cross-repo\|moved\|canonical` | `wc -l` + `rg` |
| AC4 | STATUS-canonical-home flipped: ≥ 1 hit on `STATUS vocabulary is canonical\|canonical home for the STATUS\|STATUS vocabulary lives here\|owns the STATUS vocabulary` in the moved file | `rg` |
| AC5 | HANDOFF.md is 50–150 lines | `wc -l` |
| AC6 | negative grep: no `cnos\.cdd/skills/cdd/cross-repo/SKILL\.md.*canonical\|cdd/cross-repo.*canonical` in cdd/cds/cdr skills | `rg -c` returns 0 |
| AC7 | ≥ 3 hits on `cnos.handoff/skills/handoff/cross-repo` across cdd/cds/cdr | `rg -c` |
| AC8 | ≥ 5 structural elements preserved verbatim (8 STATUS events; 4 directional cases; LINEAGE schema per case; feedback-patch header form; hat-collapse rules; known protocol edge cases) | manual spot-check + grep |
| AC9 | no schemas/handoff, no schemas/ccnf-o, no cdd-verify changes, no src/go changes, no scripts/release.sh changes | `git diff` |
| AC10 | extraction-map.md §1 rows updated with migrated Status | `rg` |

## Cross-reference repair targets (D4)

From grep against `origin/main @ 92a7442d`:

**`cnos.cdd/skills/cdd/`:**
- `gamma/SKILL.md` line 109 — cites `cdd/cross-repo/SKILL.md §2.3.3` for `drafted` direct-acceptance rule. **Canonical authority cite → re-point.**
- `gamma/SKILL.md` line 109 (second hit) — cites `cdd/cross-repo/SKILL.md` for protocol details (8-event vocab, transition graph, etc.). **Canonical → re-point.**
- `gamma/SKILL.md` line 318 — cites `cdd/cross-repo/SKILL.md §"STATUS state machine"` for master/sub rule. **Canonical → re-point.**
- `post-release/SKILL.md` line 338 — cites `cdd/cross-repo/SKILL.md` for cross-repo bundle creation; cites `cdd/cross-repo/SKILL.md §"Bundle archival rule"`. **Canonical → re-point.**
- `epsilon/SKILL.md` line 99 — cites `cross-repo/SKILL.md` (relative `../cross-repo/SKILL.md`) for cross-repo trace. **Canonical → re-point.**
- `epsilon/SKILL.md` line 137 — cites `cross-repo/SKILL.md` (relative) for iteration bundle structure. **Canonical → re-point.**
- `CDD.md` line 129 — refers in prose to "cross-repo proposals" as a coordination surface within CDS.md citation context; no `cdd/cross-repo` path cited as canonical. **Leave as-is** (the prose is descriptive; CDS.md is the cited home, not cdd/cross-repo).
- `issue/SKILL.md` line 259 — descriptive prose about cross-repo source proposals; no path citation. **Leave as-is.**
- `operator/SKILL.md` line 250 — descriptive prose; no path citation. **Leave as-is.**
- `activation/SKILL.md` lines 86–676 — descriptive prose about `.cdd/iterations/cross-repo/` directory + README; no path citation to `cdd/cross-repo/SKILL.md`. **Leave as-is.**

**`cnos.cds/skills/cds/CDS.md`:**
- Line 814 — cites `cnos.cdd/skills/cdd/cross-repo/SKILL.md §2.3.3` for `drafted` direct-acceptance rule. **Canonical → re-point.**
- Line 954 — link to `cnos.cdd/skills/cdd/cross-repo/SKILL.md` "the cross-repo proposal STATUS state machine that feeds the cross-repo-intake input surface". **Canonical → re-point.**
- Line 1506–1516 — prose paragraph titled "Operational realization location is open"; declares current operational home as `cnos.cdd/skills/cdd/cross-repo/SKILL.md`. **Historical / open-question framing post-#404-Sub-2 → update the prose to name the new canonical home (Sub 2 of #404 has now resolved the open question).** This is the substantive cite repair, not a leave-as-is.
- Line 1540 — link to `cnos.cdd/skills/cdd/cross-repo/SKILL.md §2.3` for STATUS state machine operational realization. **Canonical → re-point.**
- Line 1545 — prose "skill's current home is `cnos.cdd/skills/cdd/cross-repo/`; long-term home is open per #404". **Historical framing → update to reflect Sub 2 resolution.**
- Line 1547 — link to `cnos.cdd/skills/cdd/cross-repo/SKILL.md §2.1` for directional cases. **Canonical → re-point.**

**`cnos.cdr/`:** `rg` shows no hits on `cross-repo/SKILL\.md` in cnos.cdr. (The dispatch context mentions role overlays may cite cross-repo; verified absent.) **No edits required in cnos.cdr.**

**`cnos.cds/docs/extraction-map.md`:** lines 109, 111, 118, 283 — open-question framing about cross-repo long-term home. **This is a historical / "pre-Sub-2" doc; out of scope for D4 cite repair** (D4 targets cdd/cds/cdr **skill files**, not other packages' extraction maps). **Leave as-is.**

## Cycle plan (single-actor collapse)

Per the collapsed-roles pattern, all three role commits land on `cycle/416` from the same actor:

1. **γ commit** (this): scaffold filed; AC matrix declared; cross-reference targets enumerated.
2. **α commit**: D1–D5 land + `self-coherence.md` verifies ACs 1–10 mechanically.
3. **β commit**: `beta-review.md` runs the AC suite as adversary, attests PASS; closeouts (α/β/γ) + `cdd-iteration.md` courtesy stub + INDEX.md row land in the same commit.

Each commit message uses the role prefix per CDD §5.5b (`γ-416:`, `α-416:`, `β-416:`).

## Hard rules (acknowledgment)

- ✅ No behavioral redesign — STATUS vocab, directional cases, LINEAGE schemas, feedback-patch format, archival rule, hat-collapse all transport unchanged.
- ✅ No `schemas/handoff/` or `schemas/ccnf-o/` directories created.
- ✅ No `cn cdd verify` changes.
- ✅ No `src/go/` or `scripts/release.sh` changes.
- ✅ Old `cnos.cdd/skills/cdd/cross-repo/SKILL.md` is **not deleted** — becomes the compatibility stub.

## Closure condition

Cycle is closeable when ACs 1–10 PASS, branch `cycle/416` is pushed, and the operator-facing report is filed naming the merge instruction (`Closes #416`).
