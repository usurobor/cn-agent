<!-- sections: [intake, gap, mode, active-skills, impact, ac-evidence, cdd-trace, hard-rule-satisfied, cross-reference-map, debt, review-readiness] -->
<!-- completed: [intake, gap, mode, active-skills, impact, ac-evidence, cdd-trace, hard-rule-satisfied, cross-reference-map, debt, review-readiness] -->

# α self-coherence — cycle/402

## Intake

Issue: [cnos#402](https://github.com/usurobor/cnos/issues/402) — Phase 7 (#366): CDD.md rewrite — compress to CCNF spine.
Branch: `cycle/402` (from `origin/main @ 8c3d573b`).
Mode: design-and-build.

## Gap

`CDD.md` (1344 lines pre-cycle) hides software-specific lifecycle assumptions inside doctrine that the essay `CCNF-AND-TYPED-TRUST.md` shows must be substrate-independent. The kernel's recursion equation lives at `COHERENCE-CELL-NORMAL-FORM.md` but is not the spine of `CDD.md`; software-cycle realization (selection, lifecycle, artifact contract, gate, assessment, closure) is interleaved with what should be a generic protocol surface. The essay's hard rule (*"Do not finalize CDD.md until V works and domain evidence has somewhere else to live"*) is now satisfied (V at #392; schemas at #388; cdr v0.1 at #376), so the rewrite can ship.

## Mode

Design-and-build. The contract is pinned in the issue body's operative target; α executes per the pin, surfacing only refusal conditions if they appear.

## Active skills

Tier 1a: CDD.md, cdd/SKILL.md, alpha/SKILL.md (collapsed-α voice; γ+α+β-collapsed-on-δ per dispatch).
Tier 1b: cdd/design (terminal compact-doc shape), cdd/issue/proof (AC oracles).
Tier 2: cnos.eng/eng/document, cnos.core/skills/skill (for doctrine-doc authoring constraints).

## Impact graph

Files modified:
- `src/packages/cnos.cdd/skills/cdd/CDD.md` — 1344 lines → 159 lines (rewrite per CCNF spine)

Files not modified (despite cross-references resolving against the rewritten CDD.md):
- All other cdd/ skill files (alpha/, beta/, gamma/, delta/, epsilon/, harness/, release-effector/, operator/, post-release/, activation/, review/, release/, issue/, design/, plan/) — their `CDD.md §X` citations continue to resolve because the new CDD.md retains the named software-realization surfaces under §"Software-specific realization — pending cds extraction", with the legacy anchor forms enumerated in the closing paragraph.

Tracker issue filed as part of this cycle:
- [cnos#403](https://github.com/usurobor/cnos/issues/403) — cnos.cds bootstrap; extracts software-lifecycle realization from CDD.md when the cds package exists.

## AC-by-AC evidence

| AC | Statement | Evidence | Verdict |
|---|---|---|---|
| AC1 | `wc -l src/packages/cnos.cdd/skills/cdd/CDD.md` ≤ 300 | `wc -l` returns 159 (from 1344; ~11.8% of pre-cycle) | PASS |
| AC2 | CCNF kernel five-step equation in first ~50 lines | All 5 symbols (`αₙ.produce`, `βₙ.review`, `γₙ.close`, `V(contractₙ, receiptₙ)`, `δₙ.decide`) on CDD.md lines 18–22 | PASS |
| AC3 | CDS / CDR / c-d-X named | CDS: 3 lines × multiple occurrences; CDR: 2 lines; c-d-X: 2 lines; 17 total token occurrences | PASS |
| AC4 | Pointers section names canonical surfaces | 32 lines match the pointers oracle (COHERENCE-CELL.md, COHERENCE-CELL-NORMAL-FORM.md, RECEIPT-VALIDATION.md, ROLES.md, schemas/, role SKILL files, runtime substrate SKILL files, cnos.cdr, cnos.cds) — ≥10 required | PASS |
| AC5 | Every pre-cycle §§3–10 section accounted for | §"Software-specific realization — pending cds extraction" names: Inputs, Selection, Lifecycle, Roles and dispatch, Coordination surfaces (Tracking), Artifact contract, Mechanical, Review, Gate, Assessment, Closure, Retro-packaging, Non-goals, Large-file — all 14 family-level rows. Tracker issue #403 cited. | PASS |
| AC6 | All `CDD.md §` cross-references resolve | 29 hits across other cdd/ skill files; each anchor form (§1.4, §1.6a, §1.6c, §5.2, §5.3a, §5.3b, §9.1, §Tracking, Phase 6 step 17) named in §"Software-specific realization" closing paragraph; family-level resolution preserved | PASS |
| AC7 | Hard rule satisfied | §"Hard rule" section in new CDD.md cites V executable (#392) + domain evidence homes (#388, #376, #403). β-collapsed verifies no software-specific evidence vocabulary in the kernel sections. | PASS |

## CDD Trace

| Step | Artifact | Decision |
|---|---|---|
| 0 Observe | issue body, essay, CCNF source | Cycle #402 selected as terminal phase of #366 |
| 1 Select | `.cdd/unreleased/402/gamma-scaffold.md` | Gap named: CDD.md hides software assumptions; rewrite per CCNF spine |
| 2 Branch | `cycle/402` from `origin/main` @ 8c3d573b | Branch created via standard cycle/{N} convention |
| 3 Bootstrap | — | Not required (doctrine doc rewrite; no version-directory bundle) |
| 4 Gap | self-coherence.md §Gap | CDD.md not spined on CCNF; software realization interleaved |
| 5 Mode | self-coherence.md §Mode | design-and-build per pinned contract |
| 6 Artifacts | design-notes.md + CDD.md rewrite | Design notes precede build; CDD.md compressed 1344 → 159 lines |
| 7 Self-coherence | this file | AC1–AC7 mapped to evidence above |
| 7a Pre-review | manifest sections + diff against origin/main | β-collapsed runs AC sweep next |

## Hard rule satisfied

The essay's hard rule (`docs/gamma/essays/CCNF-AND-TYPED-TRUST.md` line 453): *"Do not finalize `CDD.md` until V works and domain evidence has somewhere else to live."* Both preconditions verified at this cycle:

**Precondition 1 — V is executable.** Built `src/go/cmd/cn/main.go` at cycle intake (cleanly, no errors); ran `./cn cdd verify --help`; output includes `--receipt <path>` (dispatch into V validator with verdict emission). The Go implementation lives at `src/packages/cnos.cdd/commands/cdd-verify/run.go`; shipped under [cnos#392](https://github.com/usurobor/cnos/issues/392) (Phase 3 of #366). The canonical contract `V : Contract × Receipt → ValidationVerdict` (per essay §"Wave 2") is reachable through the operator-facing wrapper.

**Precondition 2 — Domain evidence has homes.** Verified by `ls schemas/` at cycle intake:
- `schemas/cdd/` — generic kernel schemas (`contract.cue`, `receipt.cue`, `boundary_decision.cue`, `validation_verdict.schema.json`, fixtures); shipped under [cnos#369](https://github.com/usurobor/cnos/issues/369) + [cnos#388](https://github.com/usurobor/cnos/issues/388).
- `schemas/cds/` — software domain schemas (`receipt.cue` + fixtures); shipped under [cnos#388](https://github.com/usurobor/cnos/issues/388). The cds package itself is pending bootstrap, tracked at [cnos#403](https://github.com/usurobor/cnos/issues/403).
- `schemas/cdr/` — research domain schemas (`receipt.cue` + fixtures); shipped under [cnos#388](https://github.com/usurobor/cnos/issues/388). The cdr package itself is at v0.1, shipped under [cnos#376](https://github.com/usurobor/cnos/issues/376).

**No domain-specific evidence requirements appear in the kernel sections of the new CDD.md.** The kernel (lines 17–23), outcomes (42–46), recursion modes (54–63), scope-lift (66–72), domain packages (76–84), pointers (86–117), and hard rule (140–145) sections contain no software-specific evidence vocabulary (`test`, `code`, `branch`, `deploy`, `CI`, `release`). Software-realization vocabulary is quarantined inside §"Software-specific realization — pending cds extraction" (lines 119–137) and §"Non-goals" (lines 147–158, naming what the kernel does NOT do).

β-collapsed verifies this rg-style at AC7 in `beta-review.md`.

## Cross-reference map

The new CDD.md preserves family-level anchors for every legacy citation form used by other cdd/ skill files. Below: per legacy anchor → resolution in new CDD.md.

| Legacy anchor (citing skill) | Cited content | Resolution in new CDD.md |
|---|---|---|
| `CDD.md §1.4` (operator, review, beta, activation) | γ/α/β triadic-rule + algorithms | §"Software-specific realization" → "Roles and dispatch (§Roles)" family row |
| `CDD.md §1.4 γ algorithm Phase 1 step 3a` (harness, activation) | γ branch-creation step | Same family row; operational expansion at `gamma/SKILL.md` (per Pointers) |
| `CDD.md §1.4 β algorithm step 8` (operator) | β merge + push step | Same family row; operational expansion at `beta/SKILL.md` |
| `CDD.md §1.4 α step 10` (alpha) | α close-out re-dispatch | Same family row; operational expansion at `alpha/SKILL.md` |
| `CDD.md §1.6a` (operator, alpha) | re-dispatch prompt formats | Same family row (Roles and dispatch / sequential bounded dispatch) |
| `CDD.md §1.6c` / `§1.6c(a)` (operator, post-release) | initial-dispatch sizing, timeout budget | Same family row |
| `CDD.md §5.2` (alpha) | canonical artifact order | §"Software-specific realization" → "Artifact contract (§Artifacts)" family row |
| `CDD.md §5.3a` (alpha, release) | Artifact Location Matrix | Same family row (named explicitly in row 6: "Artifact Location Matrix") |
| `CDD.md §5.3b` (release-effector) | role/artifact ownership matrix | Same family row (named explicitly in row 6) |
| `CDD.md §9.1` (post-release × 3) | cycle iteration triggers | §"Software-specific realization" → "Assessment (§Assessment)" family row (names §9.1 cycle iteration explicitly) |
| `CDD.md §Tracking` (gamma, harness, alpha × 2, release, 2 CI templates) | polling protocol | §"Software-specific realization" → "Coordination surfaces (§Tracking)" family row |
| `CDD.md Phase 6 step 17` (release) | δ disconnect-release tag | §"Software-specific realization" → "Roles and dispatch (§Roles)" family row (names δ at the external boundary) |

The closing paragraph of §"Software-specific realization" enumerates the legacy anchor forms explicitly so a reader following any pre-cycle citation can trace which family hosts the content. Once [cnos#403](https://github.com/usurobor/cnos/issues/403) extracts the cds package, the citing skill files will re-point each anchor at the cds package directly; until then this family-level resolution is the contract.

**Note on anchor granularity (acknowledged debt).** Some legacy anchors are very specific (`§1.4 γ algorithm Phase 1 step 3a` resolves to a single executable step). The new CDD.md collapses these to family-level pointers (Roles and dispatch); the operational step lives in `gamma/SKILL.md`. A reader following the legacy anchor reaches the family that owns the content; the operational expansion is one click away in the role SKILL file. This is the conservative resolution chosen over breaking the line budget. The cds extraction (#403) is the durable fix.

## Debt

1. **Anchor-granularity collapse for legacy cross-references.** Documented above. The cds extraction tracker (#403) is the durable fix; this cycle ships the conservative family-level resolution.

2. **CCNF-X formalization remains a follow-on.** Per the issue body's active design constraints, this cycle does NOT formalize the CDD orchestration grammar (mode enum, sizing predicate, master+sub graph, dispatch-prompt schema, findings state machine). The compact CCNF-spine doc shipped here remains prose with the recursion equation as the kernel. CCNF-X is the next direction beyond this cycle.

3. **cnos.cds package bootstrap remains pending.** Tracker filed at [cnos#403](https://github.com/usurobor/cnos/issues/403). Until that cycle lands, software-lifecycle realization (selection / lifecycle / artifact-contract / gate / assessment / closure / retro-packaging) lives in CDD.md's §"Software-specific realization — pending cds extraction" section. No silent drops.

## Review-readiness

| Field | Value |
|---|---|
| Round | 1 |
| Base SHA | `8c3d573b` (origin/main at cycle start) |
| Head SHA | (set at the commit landing this file) |
| Branch CI state | docs-only changes; CI green expected on no-op runs |
| Status | ready for β-collapsed AC sweep |

β-collapsed runs AC1–AC7 against the new CDD.md with particular rigor on AC7 (hard-rule preconditions explicitly verified in `beta-review.md`).
