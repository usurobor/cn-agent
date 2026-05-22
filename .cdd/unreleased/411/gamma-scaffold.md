# γ scaffold | cnos#411 — Sub 6 of #403

**Cycle:** cnos#411
**Mode:** design-and-build (β-α-collapse-on-δ; commits prefixed `α-411:` / `β-411:` / `γ-411:`)
**Branch:** `cycle/411` from `origin/main` @ `71b25672` (Merge cycle/410)
**Wave:** Sub 6 of #403 (parallel with Sub 7; #403 closes when both land)

## Coherence gap

Sub 5 of #403 (cnos#410) migrated the last canonical software-cycle
content into `cnos.cds/skills/cds/CDS.md`. As of `origin/main` HEAD, all
content groups named in CDD.md §"Software-specific realization — pending
cds extraction" (Inputs, Selection, Lifecycle, Roles-and-dispatch,
Coordination, Artifact contract, Mechanical, Review, Gate, Assessment,
Closure, Retro-packaging, Non-goals, Large-file) now live in CDS.md. The
"pending cds extraction" section in `CDD.md` is therefore a stale marker
— it still names the surfaces as "currently resident in this file" and
still says citing skill files "resolve their `CDD.md §X` citations against
this section family until #403 lands and re-points them at the cds
package." Two surfaces of incoherence:

1. **CDD.md surface.** The §"Software-specific realization — pending cds
   extraction" section (lines 122–141) is now redundant — the content
   it names is already in CDS.md. The section must be replaced by a
   pointer paragraph + concise pointer list to CDS sections.
2. **Citing skill files.** Eleven cdd skill files (and two template
   YAMLs) carry `CDD.md §X` citations of the form `§1.4`, `§1.6a`,
   `§1.6c`, `§5.2`, `§5.3a`, `§5.3b`, `§9.1`, `§Tracking`, `Phase 6 step
   17`. These anchor forms refer to pre-#402 CDD.md structure (now
   deleted). Per the closing paragraph at CDD.md line 141, those
   citations currently resolve via the "anchor convention" — the reader
   maps `§5.3a` → "Artifact contract" family bullet. Once Sub 6 deletes
   that closing paragraph, the citations become broken anchors with no
   resolution path. They must be re-pointed at CDS canonical homes.

## Scope (pinned by δ)

Three move-classes per #411:

### (i) CDD.md §"Software-specific realization — pending cds extraction" replacement

- Rename heading to `## Software-specific realization → cnos.cds` (no
  longer "pending").
- Body: one-paragraph pointer + concise list of CDS canonical sections.
- Replace the closing "Anchor convention for cross-references"
  paragraph with citations to CDS anchors.

### (ii) Cross-reference re-pointing

For each of the 11 cdd skill files (+ 2 template YAMLs) carrying
`CDD.md §X` or `CDD.md Phase N` citations:

- **Leave unchanged** if the citation refers to kernel content (§Kernel,
  §Outcomes, §Recursion modes, §Scope-lift, §Domain packages, §Pointers,
  §Hard rule, §Non-goals kernel-doctrine variant).
- **Re-point** if the citation refers to migrated software-cycle content
  (§1.4, §1.6, §1.6a, §1.6c, §5.2, §5.3a, §5.3b, §9.1, §Tracking, Phase
  6 step 17). Re-point at the canonical CDS section with the most
  specific sub-anchor available. Preserve granularity — do not collapse
  fine-grained §5.3a/§5.3b citations into a generic CDS pointer.

### (iii) Extraction-map status note (optional)

Add a top-of-file note in `cnos.cds/docs/extraction-map.md` recording
"Sub 6 (cnos#411) complete — CDD.md pending-cds markers replaced with
CDS pointers; cross-references re-pointed."

## Citation re-pointing table

(authoritative pin for α; covers every distinct anchor form in the
existing cdd skill files)

| Old anchor (CDD.md §X) | New canonical home (CDS §X) | Notes |
|---|---|---|
| `CDD.md §1.4` (role-identity-is-git-observable; α/β/γ algorithm references) | `cnos.cdd/skills/cdd/CDD.md §1.4` is NOT a current CDD.md anchor; the role-algorithm content was the pre-CCNF doctrine. The closest current home: kernel role names live in CCNF (`COHERENCE-CELL-NORMAL-FORM.md §Kernel`); software-cycle execution lives in CDS lifecycle sections. Re-point at the operational home: `cnos.cds/skills/cds/CDS.md §"Development lifecycle"` (general) or `§"Development lifecycle" → §"Step table"` for step-specific citations. For "Triadic rule" (γ scaffolds, α implements, β reviews-and-merges), point at `CDS.md §"Field 6: Actor collapse rule"` and `§"Development lifecycle"`. | Granularity preserved by selecting the most specific sub-anchor per citation. |
| `CDD.md §1.4 γ algorithm Phase 1 step 3a` (γ creates the cycle branch) | `cnos.cds/skills/cds/CDS.md §"Development lifecycle" → §"Branch rule"` (canonical branch-creation rule, names γ-owned creation from `origin/main`) | also covered by Step 2 in the §"Step table" |
| `CDD.md §1.4 β algorithm step 8` (β merge action) | `cnos.cds/skills/cds/CDS.md §"Development lifecycle" → §"Step table"` (Step 8: β review + merge) | β merge ownership is in §Step table row 8 + §Ownership matrix |
| `CDD.md §1.6` (sequential bounded dispatch model) | `cnos.cds/skills/cds/CDS.md §"Field 6: Actor collapse rule"` (names sequential bounded dispatch) and `CDS.md §"Development lifecycle" → §"State machine"` (the lifecycle the model produces); CDS.md itself still cites `CDD.md §1.6` for operational detail — that v0.1 overlay marker is left in place where it appears | CDS.md acknowledges this is a v0.1 overlay; that's not Sub 6's gap to close |
| `CDD.md §1.6a` (re-dispatch prompt formats) | `cnos.cds/skills/cds/CDS.md §"Coordination surfaces"` (the re-dispatch mechanism is part of γ-driven coordination); citing files may also cite the v0.1 overlay path `cnos.cdd/skills/cdd/operator/SKILL.md` for the actual prompt text | Per CDS.md F1 row at line 2333, the §1.6a anchor in cdd CDD.md is the v0.1 overlay — citing the v0.1 overlay is acceptable; what changes is the **path** to the v0.1 overlay (no longer `CDD.md §1.6a` since CDD.md no longer carries `§1.6a` after Sub 6) |
| `CDD.md §1.6c(a)` (timeout/dispatch budget heuristic) | `cnos.cds/skills/cds/CDS.md §"Field 6: Actor collapse rule"` (sequential bounded dispatch is named here); the timeout heuristic itself is in `operator/SKILL.md §5.2` / `cnos.cdd/skills/cdd/CDD.md §1.6c(a)` historically — but §1.6c is not a current anchor. Re-point at CDS §"Field 6" for the dispatch-model home and leave the heuristic-constants reference within the cdd skill file's own §5.2 | If post-release/SKILL.md needs the actual constants, the constants live in operator/SKILL.md §5.2 — point there |
| `CDD.md §5.2` (canonical artifact order: design → contract → ...) | `cnos.cds/skills/cds/CDS.md §"Artifact contract" → §"Ordered flow"` | Sub 4 migrated this content under "Ordered flow" |
| `CDD.md §5.3a` (Artifact Location Matrix; canonical paths; bare-tag rule) | `cnos.cds/skills/cds/CDS.md §"Artifact contract" → §"Location matrix"` | direct migration; CDS Location matrix names the bare-tag rule explicitly |
| `CDD.md §5.3b` (γ closure declaration / frozen-snapshot rule) | `cnos.cds/skills/cds/CDS.md §"Artifact contract" → §"Frozen snapshot rule"` | direct migration |
| `CDD.md §9.1` (cycle-iteration triggers; Level/Rounds columns) | `cnos.cds/skills/cds/CDS.md §"Assessment" → §"Cycle iteration triggers"` | direct migration |
| `CDD.md §Tracking` (cycle-state evidence; git fetch reliability rule; cycle branch as coordination surface) | `cnos.cds/skills/cds/CDS.md §"Coordination surfaces"` (top-level) or `§"Coordination surfaces" → §"Polling primitives"` for the git-fetch-reliability variant | preserve specificity where the sub-anchor applies |
| `CDD.md Phase 6 step 17` (δ tag creation; release-effector mechanics) | `cnos.cds/skills/cds/CDS.md §"Development lifecycle" → §"Step table"` (Step 10: δ release) | Step 10 in the canonical 0–13 ordering covers δ tag + release |

## Files touched

- `src/packages/cnos.cdd/skills/cdd/CDD.md` — replace lines 122–141 §section + closing paragraph
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — 4 citations
- `src/packages/cnos.cdd/skills/cdd/beta/SKILL.md` — 1 citation
- `src/packages/cnos.cdd/skills/cdd/operator/SKILL.md` — 6 citations
- `src/packages/cnos.cdd/skills/cdd/release/SKILL.md` — 3 citations
- `src/packages/cnos.cdd/skills/cdd/release-effector/SKILL.md` — 1 citation
- `src/packages/cnos.cdd/skills/cdd/review/SKILL.md` — 1 citation
- `src/packages/cnos.cdd/skills/cdd/harness/SKILL.md` — 2 citations
- `src/packages/cnos.cdd/skills/cdd/activation/SKILL.md` — 3 citations
- `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` — 1 citation
- `src/packages/cnos.cdd/skills/cdd/post-release/SKILL.md` — 4 citations
- `src/packages/cnos.cdd/skills/cdd/activation/templates/github-actions/cdd-cycle-on-merge.yml` — 1 citation (in YAML comment)
- `src/packages/cnos.cdd/skills/cdd/activation/templates/telegram-notifier/cdd-notify.yml` — 1 citation (in YAML comment)
- `src/packages/cnos.cds/docs/extraction-map.md` — optional Status note (top of file)

## Acceptance criteria (AC1–AC9 per #411)

| AC | Mechanical check | Notes |
|----|------------------|-------|
| AC1 | `grep -c "pending cds extraction" src/packages/cnos.cdd/skills/cdd/CDD.md` = 0 | section renamed |
| AC2 | CDD.md has one-paragraph pointer + per-section list to CDS.md | manual inspect |
| AC3 | CDD.md kernel sections substantively unchanged | manual inspect; diff shows only line 7, the §"Software-specific realization" section, and possibly the §Domain packages CDS bullet + §Hard rule final paragraph |
| AC4 | `grep -rn "CDD\.md §" src/packages/cnos.cdd/skills/cdd/ \| grep -v "CCNF\|Kernel\|Outcomes\|Recursion\|Scope-lift\|Domain packages\|Pointers\|Hard rule\|Non-goals$"` returns 0 lines (or only reviewed false-positives) | mechanical |
| AC5 | `grep -rn "cnos.cds/skills/cds/CDS\.md" src/packages/cnos.cdd/skills/cdd/` returns ≥ 5 hits | mechanical |
| AC6 | `git diff origin/main..HEAD -- src/packages/cnos.cdr/` returns 0 lines | mechanical |
| AC7 | `git diff origin/main..HEAD -- src/packages/cnos.cds/skills/cds/CDS.md` returns 0 lines; only `cnos.cds/docs/extraction-map.md` top note allowed | mechanical |
| AC8 | Spot-check ≥ 5 re-pointed citations: target anchor exists at cited CDS path | mechanical grep |
| AC9 | CDD.md still parses; markdown structure preserved | manual inspect |

## Skills loaded

- Tier 1a: `cdd/CDD.md` (kernel doctrine)
- Tier 1b: `cdd/design/SKILL.md` (boundary discipline)
- Tier 3: this issue's citation re-pointing table above

## Branching and commits

- α-411 commits: heading + body replacement of §"Software-specific realization" section in CDD.md; cross-reference updates in 11 skill files + 2 YAML templates; optional Status note in extraction-map.md
- β-411 commits: review verdict (CLP form) at `.cdd/unreleased/411/beta-review.md`
- γ-411 commits: gamma-closeout + INDEX.md + courtesy cdd-iteration stub

Author identity per role:
- α: `alpha@cnos.local` (worktree-local for collapsed mode)
- β: same agent (collapsed); rounds recorded in `beta-review.md`
- γ: same agent (collapsed); closeout in `gamma-closeout.md`

Wave will close #403 once Sub 6 + Sub 7 land on main.
