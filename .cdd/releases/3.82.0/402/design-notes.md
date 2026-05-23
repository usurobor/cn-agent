<!-- sections: [intake, inventory-§§3-10, classification, operative-target-structure, cross-reference-map, hard-rule-verification, line-budget, risks] -->
<!-- completed: [intake, inventory-§§3-10, classification, operative-target-structure, cross-reference-map, hard-rule-verification, line-budget, risks] -->

# α design notes — cycle/402

## Intake

Rewrite `src/packages/cnos.cdd/skills/cdd/CDD.md` (1344 lines) into a compact CCNF-spine document (target ≤ 300 lines). The CCNF kernel from `COHERENCE-CELL-NORMAL-FORM.md` is the structural spine; software-specific realization is named and pointed at cds (which does not yet exist — see tracker [cnos#403](https://github.com/usurobor/cnos/issues/403)).

The essay `docs/gamma/essays/CCNF-AND-TYPED-TRUST.md` §"Wave 7 — final CDD.md rewrite" pins the design target: *"Rewrite `CDD.md` around CCNF once V is executable and domain evidence has homes in CDS/CDR."* Both preconditions hold (V works at `cn cdd verify --receipt <path>`; `schemas/cdd/` + `schemas/cds/` + `schemas/cdr/` all exist on `origin/main`; `cnos.cdr` v0.1 shipped per #376).

## Inventory of pre-cycle CDD.md §§3–10

Reading the pre-cycle file (1344 lines) end-to-end, the section structure is:

| § | Title | Lines (approx) | Substance |
|---|---|---|---|
| 0 | Purpose | 18–38 | Why CDD; the "model–reality gap" framing; the activation pointer |
| (Invocation) | Invocation model | 40–50 | cdd at global / role-local / task-local scope |
| 1 | Scope | 52–627 | §1.1 substantial change; §1.2 small change; §1.3 first principle; §1.4 large-file authoring rule; §1.5 roles + triadic rule + δ + dispatch/polling tables + cross-repo proposal lifecycle + STATUS state machine + γ/α/β/δ algorithm bodies + dispatch prompt format + named decision points; §1.6 coordination model; §1.6a/b/c re-dispatch prompt formats + timeout budget + commit checkpoints |
| 2 | Inputs | 629–668 | TSC table, encoding lag, doctor/status, last PRA |
| 3 | Selection Function | 671–745 | P0 override, op-infra override, assessment commitment, stale backlog re-eval, MCI freeze, weakest axis, max leverage, dependency order, effort tie-break, no-gap case |
| 4 | Development Lifecycle | 749–905 | §4.1 lifecycle 0–13 step table; §4.1a state machine S0–S12; §4.2 branch rule (`cycle/{N}`); §4.3 branch pre-flight; §4.4 skill loading tiers |
| 5 | Artifact Contract | 907–1097 | §5.0 terminology; §5.1 bootstrap; §5.2 ordered artifact flow; §5.3 artifact manifest (per-step format spec); §5.3a Artifact Location Matrix; §5.3b role/artifact ownership matrix; §5.4 CDD Trace; §5.5 supporting rules; §5.6 frozen snapshot rule |
| 6 | Mechanical vs Judgment Boundary | 1100–1132 | §6.1 mechanical; §6.2 judgment |
| 7 | Review | 1135–1151 | CLP terms/pointer/exit; reviewer ask list |
| 8 | Gate | 1154–1195 | Gate preconditions; §8.1 closure verification checklist F1–F10 |
| 9 | Assessment | 1198–1268 | PRA contents; §9.1 cycle iteration triggers + format + L5/L6/L7 cycle-level framework |
| 10 | Closure | 1271–1308 | §10.1 immediate outputs; §10.2 deferred outputs; §10.3 closure rule |
| 11 | Related documents | 1312–1321 | SKILL.md loader pointer; RATIONALE.md companion |
| 12 | Retro-packaging rule | 1324–1333 | Direct-to-main retro-snapshot rule |
| 13 | Non-goals | 1336–1344 | What CDD does not do |

## Classification — what moves, what stays, what dissolves into the kernel

Three dispositions per the issue body and the essay's Wave 5 / Wave 7 framing:

**A — Dissolves into the CCNF spine (no longer named as a CDD.md section).** The kernel of `COHERENCE-CELL-NORMAL-FORM.md` already names the recursion algorithm, the four outcomes, the two recursion modes, and the three scope-lift projections. CDD.md after the rewrite cites the kernel rather than re-deriving it. The pre-cycle's substrate-independent prose (the "model-reality gap" framing in §0, the triadic rule of §1.5, the dyad-plus-coordinator framing) becomes terse pointers to CCNF + ROLES.md §4 + the role SKILL.md files.

**B — Named as "pending cds extraction" with the tracker reference (still in CDD.md but explicitly named as software-specific realization).** The software-lifecycle substance — selection function (§3), lifecycle (§4), artifact contract (§5), mechanical/judgment boundary (§6), review CLP (§7), gate (§8), assessment (§9), closure (§10), retro-packaging (§12), non-goals (§13) — is software-realization detail. cnos.cds does not exist at cycle time (verified: `ls src/packages/` does not contain `cnos.cds`). The pinned contract says: keep this content in CDD.md with "pending cds extraction (tracked at [cnos#403](https://github.com/usurobor/cnos/issues/403))" notes; do NOT silently drop.

**C — Pointed-at-skill (the substance lives in a role/lifecycle skill already, CDD.md merely points).** The detailed γ/α/β/δ algorithms of §1.4–§1.6 already live at canonical depth in `gamma/SKILL.md`, `alpha/SKILL.md`, `beta/SKILL.md`, `operator/SKILL.md`. The CCNF-spine CDD.md names the roles and points at the SKILL files; the operational expansions live in their canonical homes.

**Per-section classification:**

| Source section | Class | Disposition in new CDD.md |
|---|---|---|
| §0 Purpose | A | First-paragraph pointer to CCNF; intro language compresses to "model-reality gap is closed by recursive closed-cell coherence" |
| Invocation model | C | One-line note pointing at `cdd/SKILL.md` for invocation surface |
| §1.1 substantial change / §1.2 small change | B | Lives in pending-cds-extraction §"Software-lifecycle realization" pointer; the substantial/small distinction is software-cycle realization detail |
| §1.3 first principle | A | Dissolves into the CCNF "model-reality gap" intro framing |
| §1.4 large-file authoring rule | B | Pending-cds-extraction (this is a software-authoring discipline rule) |
| §1.5 Roles + triadic rule + δ + dispatch/polling | C | One-line pointer at the role SKILL files; ROLES.md §1, §3, §4, §4a, §4b owns the generic doctrine |
| §1.5 Tracking / cross-repo / STATUS / γ-α-β algorithms / dispatch prompt format / named decision points | C | Pointed at the role SKILL files and `cross-repo/SKILL.md` for STATUS / `harness/SKILL.md` §5.4 for polling |
| §1.6 coordination model / §1.6a-c re-dispatch + budget + commits | C | Pointed at `operator/SKILL.md` + `harness/SKILL.md` |
| §2 Inputs | B | Pending-cds-extraction (TSC/lag/doctor are software-cycle inputs) |
| §3 Selection function | B | Pending-cds-extraction |
| §4 Lifecycle + branch rule + skill loading | B | Pending-cds-extraction; CDD.md notes "software-cycle lifecycle, branch mechanics, skill loading" as software-specific |
| §5 Artifact contract / location matrix / ownership matrix / CDD Trace / supporting rules / frozen snapshot | B | Pending-cds-extraction; the kernel knows about "receipt" only — `self-coherence.md`, `beta-review.md`, `RELEASE.md`, version directories are software-realization artifacts |
| §6 Mechanical vs judgment | B | Pending-cds-extraction (this is a software-realization framing) |
| §7 Review CLP | B | Pending-cds-extraction; the kernel names β.review as a typed step; CLP is the software-realization protocol |
| §8 Gate + §8.1 closure verification checklist | B | Pending-cds-extraction |
| §9 Assessment + §9.1 cycle iteration + L5/L6/L7 | B | Pending-cds-extraction (PRA is a software-cycle assessment surface; engineering-levels is a software-realization framework) |
| §10 Closure | B | Pending-cds-extraction |
| §11 Related documents | A | Subsumed into the Pointers section of the new CDD.md |
| §12 Retro-packaging | B | Pending-cds-extraction |
| §13 Non-goals | A | Dissolves into a brief "Non-goals" pointer + the kernel's substrate-independence rule |

## Operative target structure (the new CDD.md)

Per the issue body's pinned target plus the line-budget reality, the new CDD.md has the following section order. Targeted line ranges in parentheses (sum ≤ 300, including frontmatter, blank lines, and the verbatim kernel block):

1. **Frontmatter + title + one-paragraph statement of what CDD owns** (~25 lines)
2. **Kernel (CCNF)** — the five-step recursion equation block, verbatim from `COHERENCE-CELL-NORMAL-FORM.md` §Kernel lines 58–64, plus a brief naming of `closed_cellₙ` and the substrate-independence rule (~35 lines)
3. **Outcomes** — the four outcomes, verbatim from `COHERENCE-CELL-NORMAL-FORM.md` §Cell Outcomes lines 152–158, plus a brief note that `invalid` is non-terminal (~20 lines)
4. **Recursion modes** — the two modes (within-scope repair-dispatch; cross-scope projection), plus the override and reject sub-cases (~25 lines)
5. **Scope-lift** — the three projections, verbatim from `COHERENCE-CELL-NORMAL-FORM.md` §Scope-Lift lines 274–277 (~20 lines)
6. **Domain packages** — CDS (software, pending bootstrap per #403); CDR (research, v0.1 per #376); future c-d-X for other domains. Names the kernel-as-spine binding pattern. (~25 lines)
7. **Pointers** — canonical surfaces grouped (doctrine; schemas; roles; runtime substrate; realization peers). (~55 lines)
8. **Software-specific realization — pending cds extraction** — for each pre-cycle §§3–10 surface, name what lives there now and the migration target (cnos.cds via #403). (~70 lines)
9. **Hard rule** — restate the essay's hard rule, name the two preconditions, cite the satisfaction (V at #392; schemas at #388; cdr v0.1 at #376; cds extraction tracked at #403). (~15 lines)
10. **Non-goals** — what the kernel doctrine does not name (~10 lines)

Total budget: ~300 lines. The §"Software-specific realization" section is the budget pressure point; it carries one row per former §§3–10 surface plus the migration pointer. The kernel sections (1–5) are bound to verbatim CCNF source and cannot grow.

## Cross-reference map

Inventory from `rg "CDD.md §" src/packages/cnos.cdd/skills/cdd/` and the broader `rg "CDD\.md" src/packages/cnos.cdd/skills/cdd/`:

| Old citation | Citing skill file | New disposition |
|---|---|---|
| `CDD.md §1.6c(a)` (× 2) | `post-release/SKILL.md` | Re-cited as "the software-cycle lifecycle realization (pending cds extraction, [cnos#403](https://github.com/usurobor/cnos/issues/403))" — the timeout-budget heuristic lives in cds when extracted. Until then, cite continues to resolve as a software-realization pointer in CDD.md §"Software-specific realization — pending cds extraction" |
| `CDD.md §9.1` (× 3) | `post-release/SKILL.md`, plus inline references in body text | Same disposition — cycle-iteration triggers + L5/L6/L7 are software-realization (pending cds extraction); citation is updated to the new §"Software-specific realization" pointer |
| `CDD.md §5.3b` (× 1) | `release-effector/SKILL.md` | Same disposition (artifact ownership matrix is software realization) |
| `CDD.md §Tracking + §5.3a` (× 1) | `release/SKILL.md` | Same disposition (software-cycle artifact location) |
| `CDD.md §5.3a` (× 1) | `release/SKILL.md` | Same disposition |
| `CDD.md §Tracking` (× 5) | `gamma/SKILL.md`, `harness/SKILL.md`, `alpha/SKILL.md`, plus 2 CI workflow templates | Same disposition (polling protocol is software-realization) |
| `CDD.md §1.4 γ algorithm Phase 1 step 3a` (× 2) | `harness/SKILL.md`, `activation/SKILL.md` | Same disposition (γ algorithm steps are software-realization, in `gamma/SKILL.md` canonically — the citation upgrades to `gamma/SKILL.md` + the CDD.md pending-extraction pointer) |
| `CDD.md §1.6a` (× 4) | `operator/SKILL.md`, `alpha/SKILL.md` | Same disposition (re-dispatch prompts are software-realization) |
| `CDD.md §1.6c` (× 1) | `operator/SKILL.md` | Same disposition |
| `CDD.md §1.4` (× 4) | `operator/SKILL.md`, `review/SKILL.md`, `beta/SKILL.md`, `activation/SKILL.md` | Same disposition |
| `CDD.md §1.4 β algorithm step 8` (× 1) | `operator/SKILL.md` | Same disposition |
| `CDD.md §5.2` (× 1) | `alpha/SKILL.md` | Same disposition (canonical artifact order is software-realization) |
| `CDD.md §5.3a` (× 1) | `alpha/SKILL.md` | Same disposition |
| `CDD.md Phase 6 step 17` (× 1) | `release/SKILL.md` | Same disposition |
| `CDD.md §1.4 α step 10` (× 1) | `alpha/SKILL.md` | Same disposition |

**Resolution strategy:** The new CDD.md §"Software-specific realization — pending cds extraction" carries a sub-heading per former section family (selection / lifecycle / artifact-contract / review / gate / assessment / closure / etc.) plus an internal anchor. Citations from other skill files do NOT change *form* — they still read "CDD.md §X.Y" — they remain semantically valid because the new CDD.md retains the named software-realization surface under "pending cds extraction." This is the conservative resolution: cross-references continue to resolve; the next cycle (cds bootstrap, #403) will re-cite to the cds package.

**AC6 oracle check (planned):** After the rewrite, run `rg "CDD.md §" src/packages/cnos.cdd/skills/cdd/` and for each hit verify the cited section exists in the new CDD.md (either at its new §"Software-specific realization" subheading, or as one of the new §1–§5 kernel sections, or as one of the §6–§10 pointer/domain sections). β-collapsed runs this and records the verification.

## Hard rule verification

The essay's hard rule (`docs/gamma/essays/CCNF-AND-TYPED-TRUST.md` line 453): *"Do not finalize CDD.md until V works and domain evidence has somewhere else to live."*

**Precondition 1 — V is executable.** Verified at cycle intake by building `src/go/cmd/cn/main.go` and running `cn cdd verify --help`. Output includes `--receipt <path>` (the V validator dispatch with verdict emission; the canonical contract `V : Contract × Receipt → ValidationVerdict`). The Go implementation lives at `src/packages/cnos.cdd/commands/cdd-verify/run.go` and was shipped under cnos#392 (Phase 3).

**Precondition 2 — Domain evidence has homes.** Verified at cycle intake by `ls schemas/`:
- `schemas/cdd/` — generic schemas: `receipt.cue`, `contract.cue`, `boundary_decision.cue`, `validation_verdict.schema.json`, fixtures (Phase 2 + 2.5; #369 + #388)
- `schemas/cds/` — software domain: `receipt.cue` + fixtures (Phase 2.5; #388; **package itself pending bootstrap** per #403)
- `schemas/cdr/` — research domain: `receipt.cue` + fixtures (Phase 2.5; #388; full package shipped at #376)

**No domain-specific evidence requirements may appear in the new CDD.md.** β-collapsed verifies this at AC7 by `rg`-ing the new CDD.md for software-specific evidence vocabulary (test, code, package, ABI, CI, branch, deploy, release-effector) and confirming any hits are inside the §"Software-specific realization — pending cds extraction" section, NOT in the kernel / outcomes / recursion-modes / scope-lift / pointers sections.

## Line budget — refusal-condition watch

The pinned AC1 target is ≤ 300 lines. The structural budget above sums to ~300. The risk: the §"Software-specific realization" section needs to name every pre-cycle §§3–10 surface so that nothing is silently dropped. If a richer pointer per surface costs more than 70 lines total, two relaxations are available:

1. **Tighten the pointer.** Each former subsection gets one line (subsection name + "pending cds extraction at #403"). This reduces the section to ~30 lines but loses the breadcrumb each citing skill file uses. Acceptable if cross-reference resolution still works (i.e., the section retains the per-family anchors).
2. **Compress the kernel sections.** The CCNF sections must include the verbatim recursion equation but can compress the surrounding prose. Acceptable down to a floor of ~80 lines combined for §§Kernel + Outcomes + Recursion modes + Scope-lift.

Refusal condition (per dispatch): if AC6 (cross-references resolve) requires CDD.md > 300 lines, surface to operator. **First-pass plan:** target the structure above; tighten the §"Software-specific realization" section if the line budget runs over; only surface if both relaxations leave AC6 unsatisfied.

## Risks

1. **Cross-reference brittleness.** Citing skill files reference very specific subsection paths (e.g., `§1.4 γ algorithm Phase 1 step 3a`). The new CDD.md's §"Software-specific realization" section carries family-level anchors only (e.g., `§Lifecycle (pending cds extraction)`). A reader following an old citation lands on the family anchor, not on the named step. β-collapsed reviews whether this loss of granularity is acceptable; the alternative is keeping every §1.4 sub-anchor in the new CDD.md, which would blow the line budget.

   **Mitigation:** Each pending-cds-extraction subsection names what its content covers (e.g., "Selection function — P0 override, operational-infrastructure override, assessment commitment, weakest axis, max leverage, dependency order, no-gap case"). A reader following an old citation still finds the right concept; only the named anchor path is collapsed.

2. **CCNF verbatim drift.** The recursion equation block must match `COHERENCE-CELL-NORMAL-FORM.md` §Kernel lines 58–64 character-for-character. β-collapsed verifies this by `diff`-style comparison.

3. **AC1 line count.** If the line budget exceeds 300 after the first draft, I tighten the pointer section first, the kernel surrounding prose second, then surface to operator. The terminal-phase pressure is real; I do not silently exceed the budget.

4. **Domain leakage into the kernel.** The kernel sections must not contain "test", "code", "branch", "deploy", "CI", or other software-specific evidence vocabulary. The §"Software-specific realization" section is the only allowed home for that vocabulary.
