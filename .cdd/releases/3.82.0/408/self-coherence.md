<!-- sections: [Gap, Skills, ACs, Self-check, Debt, CDD-Trace, Review-readiness] -->
<!-- completed: [Gap, Skills, ACs, Self-check, Debt, CDD-Trace, Review-readiness] -->

# Self-coherence — cycle/408

## §Gap

**Issue:** [cnos#408](https://github.com/usurobor/cnos/issues/408) — Sub 3 of #403: Migrate §Selection function + §Development lifecycle to CDS (B-lite thin extract).

**Mode:** MCA (substantial cycle; multi-section canonical-content migration; B-lite extract per the issue's pinned scope ruling).

**Dispatch shape:** γ+α+β collapsed on δ (per the breadth-2026-05-12 wave manifest precedent; CDS §"Six-field instantiation contract" Field 6 — actor-collapse rule for skill/docs-class cycles).

**Coherence gap closed by this cycle:**

Before this cycle, the canonical CDS rules for **selection** (how δ picks the next cycle's gap) and **lifecycle** (the 0–13 step structure, S0–S12 state machine, branch rule, pre-flight, tier loading) had **no canonical home** in the `cnos.cds` package:

- `CDS.md` (post-#407) declared the six-field instantiation contract but explicitly delegated Selection and Lifecycle detail downstream (Field 4's "Sub-3-vs-Field-4 line" and Field 6's "Sub-3-vs-Field-6 line" both name Sub 3 of #403 as the surface that ships the detail).
- `cnos.cdd/skills/cdd/CDD.md` lines 122–143 carried the surfaces as bullet-name markers in §"Software-specific realization — pending cds extraction" but no rule text — the actual canonical rule statements had been quarantined out at the #402 CCNF spine rewrite.
- The actual canonical text lived only in pre-#402 CDD.md (commit `8f06a606^`) §3 + §4, which was no longer the citable canonical home after the spine rewrite.
- Operational mechanics continued to live in `cnos.cdd/skills/cdd/{gamma,alpha,beta}/SKILL.md` but those files cited `CDD.md §3` and `CDD.md §1.4` anchors that no longer resolved to rule statements.

After this cycle, CDS.md is the canonical home for both surfaces; operational realization stays in the cdd role/runtime skills as the v0.1 overlay (per the B-lite scope ruling).

## §Skills

**Tier 1a (CDS / kernel authority):**
- `src/packages/cnos.cds/skills/cds/CDS.md` (the file under edit; read end-to-end at pre-edit state).
- `src/packages/cnos.cdd/skills/cdd/CDD.md` (read; not edited — hard rule AC6).
- `src/packages/cnos.cds/skills/cds/SKILL.md` (loader; read).

**Tier 1b (lifecycle phase + design):**
- `src/packages/cnos.cdd/skills/cdd/design/SKILL.md §3.2` "one source of truth" — the discipline this cycle is a canonical-content move under.
- `src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` — AC interpretation.
- `src/packages/cnos.cdd/skills/cdd/review/SKILL.md` — β self-review CLP form.

**Tier 1c (β closure bundle):**
- `src/packages/cnos.cdd/skills/cdd/release/SKILL.md §3.8` configuration-floor cap on γ-axis + β-axis at A- for collapsed-role cycles.

**Tier 2 (general engineering):**
- `src/packages/cnos.eng/skills/eng/writing` family — markdown authoring.
- `src/packages/cnos.core/skills/design` — architecture-level reasoning (section placement, B-lite scope adherence judgment).

**Tier 3 (issue-specific):**
- `cnos.cds/docs/extraction-map.md` rows 1 + 2 — destination commitments.
- `cnos.cdd/skills/cdd/gamma/SKILL.md`, `alpha/SKILL.md`, `beta/SKILL.md` — source for the operational mechanics being cited (read; not edited per hard rule AC6).
- `cnos.cdr/skills/cdr/CDR.md` — structural sibling (read for sub-section pattern emulation).
- pre-#402 `cnos.cdd/skills/cdd/CDD.md` (commit `8f06a606^`) §3 + §4 — source for the canonical rule text (read via `git show 8f06a606^:src/packages/cnos.cdd/skills/cdd/CDD.md`).

## §ACs

Per cnos#408 §"Acceptance criteria", verified at HEAD `fb401f99` (α-408 commit):

### AC1: §Selection function section exists in CDS.md

```text
$ grep -n "^## Selection function" src/packages/cnos.cds/skills/cds/CDS.md
778:## Selection function
```

PASS — exactly 1 match. Section is non-empty (spans lines 778–961 in the post-α state; ~184 lines).

### AC2: §Development lifecycle section exists in CDS.md

```text
$ grep -n "^## Development lifecycle" src/packages/cnos.cds/skills/cds/CDS.md
963:## Development lifecycle
```

PASS — exactly 1 match. Section is non-empty (spans lines 963–1219; ~257 lines).

### AC3: §Selection covers all 10 named rules

Each rule name appears as a `### ` sub-heading in §Selection:

| Rule | Sub-heading present | Line |
|---|---|---|
| P0 override | ✓ | line 833 |
| Operational-infrastructure override | ✓ | line 845 |
| Assessment-commitment default | ✓ | line 861 |
| Stale-backlog re-evaluation | ✓ | line 869 |
| MCI freeze check | ✓ | line 887 |
| Weakest-axis rule | ✓ | line 896 |
| Maximum-leverage rule | ✓ | line 910 |
| Dependency order | ✓ | line 917 |
| Effort-adjusted tie-break | ✓ | line 924 |
| No-gap case | ✓ | line 931 |

Mechanical check: `for r in <rules>; do grep -c "^### $r" CDS.md; done` returns `1` for each of the 10 rule names. PASS.

### AC4: §Lifecycle covers all 5 named components

```text
$ grep -n "^### Step table\|^### State machine\|^### Branch rule\|^### Branch pre-flight\|^### Skill loading tiers" src/packages/cnos.cds/skills/cds/CDS.md
977:### Step table
1007:### State machine
1040:### Branch rule
1079:### Branch pre-flight
1104:### Skill loading tiers
```

PASS — all 5 components present as sub-headings with rich content:
- Step table — 14-row table (Steps 0–13) with `# / Step / Owner / Purpose / Required output` columns.
- State machine — 13-row table (S0–S12) with `State / Owner / Required inputs / Required outputs / Next state / Failure-retry` columns.
- Branch rule — `cycle/{N}` canonical, γ-from-`origin/main` ownership, γ-session-branch rule, legacy-shape warn-only.
- Branch pre-flight — 5 verified checks γ runs before creating the branch.
- Skill loading tiers — 1a / 1b / 1c / 2 / 3 with explicit lists per tier.

### AC5: Operational-realization pointer present in each new section

```text
$ grep -n "^### Operational realization" src/packages/cnos.cds/skills/cds/CDS.md
933:### Operational realization
1191:### Operational realization
```

PASS — both new top-level sections end with an "### Operational realization" sub-section. §Selection's pointer cites `cnos.cdd/skills/cdd/gamma/SKILL.md §2.1, §2.2, §2.2a`, `alpha/SKILL.md §2.1`, and `cross-repo/SKILL.md`. §Lifecycle's pointer cites `gamma/SKILL.md`, `alpha/SKILL.md`, `beta/SKILL.md`, `delta/SKILL.md`, `harness/SKILL.md`, `release-effector/SKILL.md`, `operator/SKILL.md` — at least one cdd role-skill file per pointer.

### AC6: cnos.cdd untouched (hard rule)

```text
$ git diff origin/main..HEAD -- src/packages/cnos.cdd/
(empty)
$ git diff origin/main..HEAD -- src/packages/cnos.cdd/ | wc -l
0
```

PASS — zero lines.

### AC7: cnos.cdr untouched (hard rule)

```text
$ git diff origin/main..HEAD -- src/packages/cnos.cdr/
(empty)
$ git diff origin/main..HEAD -- src/packages/cnos.cdr/ | wc -l
0
```

PASS — zero lines.

### AC8: extraction-map.md Status column updated

Verified:

```text
$ grep -n "^\*\*Status:\*\*" src/packages/cnos.cds/docs/extraction-map.md
62:**Status:** **v0.1 migrated; canonical at [`CDS.md §"Selection function"`](../skills/cds/CDS.md)** ...
85:**Status:** **v0.1 migrated; canonical at [`CDS.md §"Development lifecycle"`](../skills/cds/CDS.md)** ...
```

Both rows 1 (§Selection function) and 2 (§Development lifecycle) carry an explicit `**Status:** **v0.1 migrated; canonical at ...**` line that names the CDS canonical path. Other rows (§Coordination surfaces / §Artifact contract / §Mechanical / §Review / §Gate / §Assessment / §Closure / §Retro-packaging / §Non-goals / §Large-file) untouched — grep returns only the 2 rows. PASS.

### AC9: No deep role rewrites

```text
$ ls src/packages/cnos.cds/skills/cds/
CDS.md
SKILL.md
lifecycle
selection
```

No new directory under `skills/cds/{alpha,beta,gamma,delta,epsilon,operator}/`. Only `selection/` and `lifecycle/` (per D3's permission) — each with one `SKILL.md` file of ≤40 lines:

```text
$ wc -l src/packages/cnos.cds/skills/cds/selection/SKILL.md src/packages/cnos.cds/skills/cds/lifecycle/SKILL.md
  39 src/packages/cnos.cds/skills/cds/selection/SKILL.md
  40 src/packages/cnos.cds/skills/cds/lifecycle/SKILL.md
```

PASS — both within the ≤40-line cap.

## §Self-check

Did α's work push ambiguity onto β? **No.** Each AC has a mechanical oracle (grep + line count + git-diff line count). The B-lite scope ruling is verified empirically: the canonical rule statements moved into CDS.md (10 rules, 14-row step table, 13-row state machine, branch rule, pre-flight, 5-tier loading); operational realization stays in the cdd role skills as the v0.1 overlay (cited from each new section's "Operational realization" sub-section). No CCNF kernel doctrine is restated — every kernel reference (CCNF, recursion modes, scope-lift, evidence-binding) is by citation.

Is every claim backed by evidence in the diff? **Yes.** Each AC's PASS row above carries a mechanical command + its output. The canonical-rule source pin (pre-#402 CDD.md commit `8f06a606^` §3 + §4) is named in the α commit message and in `gamma-scaffold.md §"Source mining map"`.

**Preserve granular anchors check.** The pre-#402 CDD.md used anchors like `CDD.md §3.1 P0 override`, `CDD.md §3.7 maximum leverage`, `CDD.md §4.1 lifecycle steps`, `CDD.md §4.1a state table`, `CDD.md §4.2 branch rule`, `CDD.md §4.3 branch pre-flight`, `CDD.md §4.4 skill loading`. The new CDS.md sub-headings preserve equivalent fine-grained anchors: `CDS.md §"Selection function" → "P0 override"`, `→ "Maximum-leverage rule"`, etc.; `CDS.md §"Development lifecycle" → "Step table"`, `→ "State machine"`, `→ "Branch rule"`, `→ "Branch pre-flight"`, `→ "Skill loading tiers"`. A future cycle citing the granular anchor lands at the canonical CDS home (per the "Preserve granular anchors" doctrine touchpoint in the dispatch contract).

**No CCNF kernel duplication check.** The new sections cite kernel doctrine by reference only:
- §Lifecycle cites `COHERENCE-CELL-NORMAL-FORM.md §Recursion Modes` for the within-scope-vs-cross-scope distinction.
- §Selection does not introduce new kernel content; it cites existing CDS Field 4 (δ cadence) and Field 5 (ε cadence) for the pairing.
- Both sections cite the existing cdd v0.1 overlay for operational mechanics; neither restates the cdd algorithms.

**Pointer-only check (negative — must NOT be pointer-only).** Each new section ships actual canonical content: 10 rule paragraphs in §Selection; 14-row step table + 13-row state table + branch rule + pre-flight + tier structure in §Lifecycle. A reader who reads only the CDS.md sections has the canonical statement of what each rule says and how the lifecycle composes — the "Operational realization" pointer is an *additional* surface, not the *primary* surface. This satisfies the B-lite "positive" criterion: a reader can cite a CDS path as the canonical home for the surface family.

**Full v1 rewrite check (negative — must NOT attempt v1 rewrite).** No new files under `skills/cds/{alpha,beta,gamma,delta,epsilon,operator}/`; the operational realization pointer routes to existing cdd role skills as the v0.1 overlay. The optional `selection/SKILL.md` + `lifecycle/SKILL.md` thin overlays are ≤40 lines each, contain no novel operational content, and explicitly delegate to the cdd overlay. This satisfies the B-lite "negative" criterion: no role-overlay deep rewrite is attempted.

## §Debt

**Cross-cut from §Roles family (acknowledged, not closed):** The pre-#402 CDD.md's §4 lifecycle content cross-referenced §1.5 Roles content (the triadic rule, dyad-plus-coordinator framing, dispatch model, dispatch-prompt formats) which extraction-map row 2 also lists. Per the B-lite scope and the issue's explicit non-goal "Do NOT author full role overlays under `skills/cds/{alpha,beta,gamma,delta,epsilon,operator}/`", the role-side mechanical detail (γ dispatch-prompt format with the 7-axis `## Implementation contract` table; the §1.6 sequential bounded dispatch contract; §1.6a/§1.6b re-dispatch prompt forms; §1.6c initial dispatch sizing) remains in the cnos.cdd v0.1 overlay and is cited by reference from the §Development lifecycle → "Operational realization" pointer. This is intentional, declared in extraction-map row 2's updated Status note ("**§Roles cross-cut content ... remains pending — ... deferred per the B-lite scope ruling**"), and is the **Sub-?-vs-§Roles** open question routed to a post-#403 cycle when v1 role overlays are authored.

**Filesystem `.cdd/` → `.cds/` rename (acknowledged, not closed):** The new CDS.md §Development lifecycle uses `.cdd/unreleased/{N}/` paths because the current project binding has not yet performed the filesystem rename. Per `docs/extraction-map.md §14` open questions, the rename is a separate post-#403 coordination problem; CDS.md §Lifecycle does not pre-rename paths it cannot observe.

**No `cdd-skill-gap` / `cdd-protocol-gap` / `cdd-tooling-gap` / `cdd-metric-gap` surfaced during α work.** The cycle proceeded under the pinned implementation contract; no skill loaded by α failed to prevent a finding (because no β finding has yet been authored at this point — β self-review follows).

## §CDD Trace

| Step | Status | Evidence |
|---|---|---|
| 0. Observe | ✓ | γ-scaffold.md §"Peer enumeration" verified §Selection + §Lifecycle absent at d9829412; extraction-map rows 1+2 read; CDR.md structural sibling read. |
| 1. Select | ✓ | Issue #408 selected by δ; decisive clause = cnos#403 wave's Sub-3 dispatch under B-lite scope ruling. |
| 2. Branch | ✓ | `cycle/408` created from `origin/main` at `d9829412`; pre-flight passed (no stalled `.cdd/unreleased/408/`; origin/cycle/408 did not exist; issue open with full ACs). |
| 3. Bootstrap | ✓ | `.cdd/unreleased/408/gamma-scaffold.md` authored and committed at `d5c93e1e` before α dispatch (β rule 3.11b pre-empted). |
| 4. Gap | ✓ | §Gap above names the incoherence: canonical Selection + Lifecycle rules had no canonical home in cnos.cds. |
| 5. Mode | ✓ | MCA; collapsed γ+α+β on δ for skill/docs-class cycle per CDS §Field 6. |
| 6. Artifacts | ✓ | CDS.md (2 new sections + manifest update + version-status update); extraction-map.md (Status column rows 1+2); selection/SKILL.md + lifecycle/SKILL.md thin overlays. α commit `fb401f99`. |
| 7. Self-coherence | ✓ | this file (§Gap, §Skills, §ACs, §Self-check, §Debt, §CDD Trace authored; §Review-readiness below). |

## §Review-readiness

| Round 1 | base SHA: `d9829412` | head SHA: `fb401f99` | branch CI: not required (markdown-only diff; CDS package CI not yet wired in for `cnos.cds` per #406 deferral) | ready for β |

Per `alpha/SKILL.md §2.6` pre-review gate, transient rows re-validated immediately before signal:

- Row 1 (branch rebased onto current `origin/main`): `git rev-parse origin/main` = `d9829412`; cycle/408 base = `d9829412`; rebase needed: no.
- Row 10 (branch CI green): CDS package does not yet have CI wired in per cnos#406; manual pre-review verification = grep + diff + wc oracles per AC1–AC9 above; all PASS. Pre-existing kernel-side CI (cn binary build, schema validation, kata pass) does not exercise cnos.cds markdown content; the diff does not touch any surface that CI checks.

β self-review follows directly (β-α collapsed on δ; same agent reads α's work under β identity in `beta-review.md`).
