# γ scaffold — cycle/407 (Sub 2 of cnos#403)

**Issue:** [cnos#407](https://github.com/usurobor/cnos/issues/407) — Sub 2 of [cnos#403](https://github.com/usurobor/cnos/issues/403): Author cnos.cds/CDS.md — six-field software instantiation contract.

**Parent:** [cnos#403](https://github.com/usurobor/cnos/issues/403) (cnos.cds bootstrap tracker — extract-by-reference v0.1 wave; this is Sub 2 of 7; gates Subs 3–5).

**Mode:** design-and-build; γ+α+β collapsed on δ (per the breadth-2026-05-12 wave manifest precedent for skill/docs-class cycles; mirrors cycle/406 Sub 1 dispatch shape).

## Implementation contract (pinned by δ — issue body verbatim)

| Axis | Pinned value |
|---|---|
| Language | Markdown |
| CLI integration target | None new |
| Package scoping | `src/packages/cnos.cds/skills/cds/CDS.md` (canonical home) + minor edits to `src/packages/cnos.cds/skills/cds/SKILL.md` to remove the "forthcoming" note |
| Existing-binary disposition | N/A |
| Runtime dependencies | None |
| JSON/wire contract | N/A (markdown only) |
| Backward compat | `cnos.cdd` is **NOT modified**. "Pending cds extraction" markers in CDD.md stay until Sub 6. `cnos.cdr` is **NOT modified**. |

## Surface

Files created or edited (D1 + D2 per #407):

1. **NEW** `src/packages/cnos.cds/skills/cds/CDS.md` — six-field software instantiation contract; mirrors `cnos.cdr/skills/cdr/CDR.md` section structure; ~500–700 lines target.
2. **EDIT** `src/packages/cnos.cds/skills/cds/SKILL.md` — remove "forthcoming" / "Sub 2 lands it" note about CDS.md; add CDS.md as Step 1 in Load order; update `## Rule` section to name CDS.md as normative source for the six fields + boundary + anchor; v0.1 caveat section either removed or repurposed.
3. **DELETE** `src/packages/cnos.cds/skills/cds/.gitkeep` (optional per AC9 — CDS.md now occupies the directory; placeholder no longer load-bearing).

Cycle artifacts authored under `.cdd/unreleased/407/`:

4. `gamma-scaffold.md` (this file).
5. `self-coherence.md` — α self-coherence pass; AC-by-AC mechanical check.
6. `beta-review.md` — β-collapsed-on-δ review; AC verification + AC4 no-duplication audit + doctrine-coherence-with-cdr-template check.
7. `alpha-closeout.md` — α-level findings (likely empty; this is contract authoring).
8. `beta-closeout.md` — β-level findings.
9. `gamma-closeout.md` — γ-level closure summary; finding dispositions.
10. `cdd-iteration.md` — courtesy empty-findings stub per cycle/401 rule if `protocol_gap_count == 0`.
11. `.cdd/iterations/INDEX.md` — row appended for cycle 407 (per the courtesy convention).

Files NOT touched (per #407 Non-goals + design discipline):

- `src/packages/cnos.cdd/**` — explicitly out of scope; AC7 hard rule (mechanical `git diff` empty).
- `src/packages/cnos.cdr/**` — explicitly out of scope; AC8 hard rule (mechanical `git diff` empty).
- `schemas/cds/**` — already exists per #388; CDS.md cites it once. No edits.
- `src/packages/cnos.cds/skills/cds/{alpha,beta,gamma,delta,epsilon}/SKILL.md` — Subs 3–5 territory; AC9 verifies no new files under those paths.
- `src/packages/cnos.cds/docs/extraction-map.md` — Sub 1 deliverable; CDS.md cites it but does not re-author it.

## Design judgment calls (contract-shape)

The build half (markdown authoring) is mechanical. The design half — what the six fields actually say for software-class matter — requires judgment calls. Each judgment call below is recorded so β review can audit it.

### Judgment 1 — Field 1 (Matter type): the scope of "software artifact"

CDR.md Field 1 enumerates research-class matter (claims, hypotheses, methods, datasets, analyses, reports). CDS.md Field 1 must enumerate software-class matter without leaking into Subs 3–5's lifecycle territory.

**Decision:** declare matter type as: source code, tests, documentation (skills / doctrine / READMEs), CI definitions, schemas (CUE / JSON-Schema), runtime contracts (manifest files, CLI command shapes), design notes / RFCs. Each item is named as a *type of artifact* (what β reviews against, what γ closes out about) — not as *a lifecycle step that produces it* (which is Sub 3 territory) and not as *evidence schemas* (which are `schemas/cds/` territory).

The Sub-3-vs-Field-1 line: Field 1 says "α produces code, tests, and docs"; Sub 3 will say "the 0–13 step lifecycle that produces code, tests, and docs". Field 1 owns the artifact taxonomy; Sub 3 owns the per-step procedure.

### Judgment 2 — Field 2 (Review oracle): what counts as the oracle

CDR.md Field 2 names six research oracles (falsifiability, diagnostic, reproduction-from-clean, citation integrity, data-policy compliance, claim/evidence alignment). CDS.md Field 2 must name the software-engineering equivalent — the failures β catches structurally, not "what β does step-by-step".

**Decision:** declare six software oracles:

1. **Compilation / type-check / build** — the build succeeds, types resolve, the artifact links.
2. **Tests pass** — the cycle's tests pass (cycle-touched suites at minimum; full suite where mechanical).
3. **AC verification** — each AC has a mechanical or read-check oracle; β verifies each oracle ran and recorded its result.
4. **No regressions** — pre-existing tests pass; no measured-coherence degradation on adjacent surfaces.
5. **Implementation-contract coherence (β Rule 7 per cnos#393)** — α's matter satisfies the pinned axes (language, scoping, runtime deps, JSON/wire, backward compat).
6. **Evidence-binding (per CCNF kernel)** — the receipt carries typed refs to the cycle's artifacts; γ binds the evidence at close-out per the kernel's evidence-binding rule.

The Sub-5-vs-Field-2 line: Field 2 names the oracles (what β catches); Sub 5 will name the CLP form, the reviewer-ask-list shape, and the per-finding disposition workflow. Field 2 owns the oracle taxonomy; Sub 5 owns the review-procedure operational detail.

### Judgment 3 — Field 3 (γ close-out artifact): the per-cycle artifact set vs the typed receipt

CDR.md Field 3 names the typed `#CDRReceipt` schema as the load-bearing γ artifact, with narrative wave reports as separate human-facing surfaces. CDS.md Field 3 must name the per-cycle close-out artifact set (`self-coherence.md`, `beta-review.md`, `alpha-closeout.md`, `beta-closeout.md`, `gamma-closeout.md`, `cdd-iteration.md`) AND the typed `#CDSReceipt` schema as the typed surface.

**Decision:** declare both surfaces. The narrative close-out files are the cycle-local close-out artifact set (per CDD §5.5b — the convention that pre-existed CCNF); the typed `#CDSReceipt` (per `schemas/cds/receipt.cue`) is the parent-facing typed surface that ties the close-out files to V (the validator). The two surfaces compose: the typed receipt carries `evidence_refs` pointing at the close-out files; V dereferences the refs to validate the receipt.

The cycle/401 rule (`cdd-iteration.md` required only when `protocol_gap_count > 0`; courtesy stub when count is 0) is the CDS-specific cadence rule for the iteration artifact, inherited from the cdd/post-release/SKILL.md authoring lineage.

The Sub-4-vs-Field-3 line: Field 3 names the artifact set + the typed receipt; Sub 4 will detail the artifact contract (the Artifact Location Matrix, the role/artifact ownership matrix, the CDS Trace format, the frozen-snapshot rule, the bootstrap procedure). Field 3 owns the artifact taxonomy; Sub 4 owns the per-artifact operational detail.

### Judgment 4 — Field 4 (δ cadence): per-cycle vs release-shaped

CDR.md Field 4 declares the cadence as gate-transition-shaped (not release-shaped); CDR explicitly rejects the engineering "cut a release" framing. CDS.md Field 4 must declare the cadence as release-shaped (because CDS *does* release; CDR's rejection was specifically because research has no release).

**Decision:** declare the cadence as per-cycle, with the δ boundary decision gating merge (per the CCNF kernel's step 5). The cycle is the unit; merge is the boundary effection; release is the cross-cycle bundling (a release contains one or more merged cycles, packaged for downstream consumers). The inward-membrane discipline (per cnos#393 δ-as-architect: implementation-contract enrichment at dispatch) runs at the open of each cycle; the outward-membrane discipline runs at merge.

Cadence triggers:
- **Per-cycle:** each cycle opens when δ dispatches a gap and closes when δ records the boundary decision on the receipt.
- **Per-release:** when the unreleased bundle reaches a coherent shipping point (PRA-readiness, RELEASE.md author-ability), δ tags a release; the release-effector mechanics handle the tag/build/deploy.
- **Per-merge:** the merge into main is the boundary effection per cycle; cycles that fail the boundary do not merge.

The Sub-3-vs-Field-4 line: Field 4 names the cadence (per-cycle + release-shaped); Sub 3 will detail the lifecycle state machine S0–S12, the branch rule, the skill-loading tiers. Field 4 owns the cadence taxonomy; Sub 3 owns the per-step lifecycle operational detail.

### Judgment 5 — Field 5 (ε iteration cadence): the engineering gap class taxonomy

CDR.md Field 5 declares six research-class gap classes (missing data gates, overclaiming, unreproducible numbers, weak citation discipline, recurring oracle ambiguity, construct drift). CDS.md Field 5 must declare the engineering-class gap classes per ROLES.md §4b.3 reference instantiation.

**Decision:** declare the four cdd-class gap names as the CDS gap classes (since these are the gap classes the empirical anchor — cnos cycles #364 through #406 — has been using):

1. **`cds-skill-gap`** — procedural skill underspecified or wrong (a CDS skill file did not give the operator enough to avoid the finding).
2. **`cds-protocol-gap`** — CDS doctrine itself drifted (the contract — CDS.md — needs an edit).
3. **`cds-tooling-gap`** — tooling absent, wrong, or unavailable (a mechanical check the cycle needed did not exist or did not work).
4. **`cds-metric-gap`** — measurement missing or wrong (a coherence-delta measurement the cycle relied on did not produce the right signal).

Naming-class change note: the current empirical-anchor cycles use `cdd-*-gap` because they predate the cds extraction. CDS.md declares the new `cds-*-gap` names as the canonical CDS taxonomy; the rename is mechanical (it does not change which findings the classes catch). The transition happens organically: future CDS cycles use `cds-*-gap`; pre-#403 cycles' `cdd-*-gap` markers stay as historical record.

The cadence rule (`cds-iteration.md` required when `protocol_gap_count > 0`; courtesy empty stub when 0) is inherited from the cycle/401 cdd-cadence rule; the rename `cdd-iteration.md` → `cds-iteration.md` is mechanical and is the Sub-5-territory operational detail.

The Sub-5-vs-Field-5 line: Field 5 names the gap classes + the cadence rule; Sub 5 will detail the per-finding shape, the iteration aggregator path (`.cdd/iterations/INDEX.md` → `.cds/iterations/INDEX.md` post-rename), the MCA discipline. Field 5 owns the class taxonomy; Sub 5 owns the per-class operational detail.

### Judgment 6 — Field 6 (Actor collapse rule): the engineering-class collapse table

CDR.md Field 6 declares α=β never permitted (because research's overclaim failure mode is precisely what α cannot self-detect) and names two safe collapses (γ=δ for small projects; ε=δ until receipt-stream volume warrants split). CDS.md Field 6 must declare the engineering-class collapse table.

**Decision:** declare four collapse rules:

1. **α=β prohibited for substantive software work.** Same structural reasoning as CDR: order-1 review of α's order-0 production requires β's frame to be outside α's; same actor cannot hold both within a single cycle.
2. **β-α-collapse-on-δ permitted for skill/docs-class cycles** (per breadth-2026-05-12 wave manifest precedent). When the matter is structural-mirror or migration (no novel substantive code; mechanical correctness verifiable from the spec), δ may dispatch γ+α+β collapsed on a single agent. The release/SKILL.md §3.8 configuration-floor clause then caps β-axis and γ-axis at A-.
3. **γ=δ permitted for small-project regimes** (same as CDR). The collapse is safe when γ's coordination authority does not compromise β's independence — γ sees both α and β but does not author either's matter.
4. **ε=δ permitted until receipt-stream volume warrants split** (same as CDR generic doctrine per ROLES.md §4b.6). The collapse is one of the *safe* collapses per ROLES.md §4 hats-vs-actors.

Project-specific stricter floors permitted (same pattern as CDR Field 6): a downstream engineering project may impose stricter rules in its `<project>/.cds/POLICY.md` (or `.cdd/POLICY.md` until the rename lands).

The Sub-3-vs-Field-6 line: Field 6 names the collapse table; Sub 3 will detail the dispatch model (§1.6 sequential bounded; §1.6a/§1.6b/§1.6c re-dispatch forms). Field 6 owns the actor-shape taxonomy; Sub 3 owns the per-shape operational detail.

### Judgment 7 — Architecture choice: rephrase from CDS side, do not verbatim-copy CDR

Per #407 active design constraint, the §"Architecture choice" section may reuse the (a)-vs-(b) framing from CDR but framed from CDS's side. Do not verbatim-copy CDR.md prose.

**Decision:** structure the section as:
- **The decision:** option (a) common constitution + per-protocol procedures; show the same ASCII layout CDR uses but with `cnos.cds` as the per-protocol example (paralleling cnos.cdr in the same role).
- **Option (b) rejected:** state the five reasons in CDS-specific language (e.g. "Engineering α and research α share α's verb 'produces' but diverge sharply on what counts as production discipline" — rephrased as "Engineering α and research α share α's verb 'produces' but the failure modes diverge sharply: engineering's stalled-loop failure is α-internal; research's overclaim failure is α-blind").
- **Design source:** same source citations as CDR (CCNF essay; ROLES.md §4a.2; schemas/cdd/README.md §"Architectural choice"; cnos#388 / #376 / #403 / #406).
- **Cross-reference:** cite CDR.md §"Architecture choice" as the parallel record.

The verbatim-copy prohibition is not just stylistic — it is the per-protocol authorship discipline. CDR and CDS each record the inheritance in their own voice; future cdw/cda will do the same.

### Judgment 8 — Empirical anchor: cnos as the anchor

CDR.md Empirical anchor cites `usurobor/cph` (an external research repo) as the CDR anchor and defers the surface-by-surface mapping to Sub 4. CDS.md Empirical anchor cites `usurobor/cnos` itself (this repo) as the CDS anchor — every cnos cycle since #364 has been a CDS cycle in practice (artifact-improvement-under-repairable-feedback discipline applied to engineering work).

**Decision:** cite cnos cycles #364 (COHERENCE-CELL.md doctrine landing) through #406 (cnos.cds skeleton + extraction map) as the empirical anchor. Name representative milestones: #366 (CCNF roadmap), #376 (cdr wave), #393 (δ-as-architect / Rule 7), #402 (CCNF spine rewrite), #406 (cds skeleton). Name `docs/empirical-anchor-cdd.md` as the Sub 7 deliverable that will perform the surface-by-surface mapping (mirroring cnos.cdr/docs/empirical-anchor-cph.md per cycle/406 README §"Forthcoming surfaces").

The naming `empirical-anchor-cdd.md` (not `empirical-anchor-cnos.md`) is intentional: the artifact maps the *current `.cdd/` cycle artifacts* (which are the historical anchor) to the *CDS six-field surfaces* (which are the target). Sub 7 may rename to `empirical-anchor-cnos.md` post-`.cdd/`-to-`.cds/`-filesystem-migration; the v0.1 naming follows the cycle/406 README precedent.

## AC oracle approach (issue body verbatim)

| AC | Oracle | Surface |
|----|--------|---------|
| AC1 | `grep "^## " skills/cds/CDS.md` returns the seven required top-level headings (Purpose, Architecture choice, Persona Protocol Project, Six-field instantiation contract, Empirical anchor, Related documents, Non-goals). | `grep` + file existence |
| AC2 | `grep -c "^### Field [1-6]:" skills/cds/CDS.md` returns 6. | `grep -c` |
| AC3 | §0 Purpose names both "artifact-improvement-under-repairable-feedback" (or equivalent ROLES.md §4a.2 phrasing) and "truth-preservation-under-uncertainty" (or equivalent). | `grep` + read-check |
| AC4 | No prose paragraph in CDS.md duplicates ≥ 50 contiguous characters from CDD.md / COHERENCE-CELL.md / COHERENCE-CELL-NORMAL-FORM.md / RECEIPT-VALIDATION.md / ROLES.md. Architecture-choice section references the schema README rather than restating the five-reason rationale. | Manual review (β audit) |
| AC5 | `skills/cds/SKILL.md` Step 1 of "Load order" loads CDS.md. The "forthcoming" / "Sub 2 lands it" note about CDS.md is removed. `## Rule` section names CDS.md as the normative source. | read-check |
| AC6 | `## Empirical anchor` cites cnos cycles #364–#406 (or representative subset). Names `docs/empirical-anchor-cdd.md` as Sub 7 deliverable. | read-check |
| AC7 | `git diff origin/main..HEAD -- src/packages/cnos.cdd/` returns empty. | `git diff` |
| AC8 | `git diff origin/main..HEAD -- src/packages/cnos.cdr/` returns empty. | `git diff` |
| AC9 | No files under `src/packages/cnos.cds/skills/cds/{alpha,beta,gamma,delta,epsilon,operator}/`. The `.gitkeep` placeholder may stay or be removed. | `ls` + `find` |

## Branch/identity

- Branch: `cycle/407` from `origin/main` (HEAD == `987acd04` Merge cycle/406; created).
- Worktree: `worktree-agent-ad82368de52f4ae9d` (harness-created).
- γ identity: `gamma@cdd.cnos`
- α identity: `alpha@cdd.cnos`
- β identity: `beta@cdd.cnos`
- δ identity: `operator@cdd.cnos` (post-merge to main; merge by operator, not this agent per dispatch).

(Note: identity strings retain `cdd.cnos` suffix per the current convention; the `cds.cnos` rename is post-#403-wave territory and not in this sub's scope.)

## Dispatch shape

This is `cdd/operator/SKILL.md §5.2` (γ+α+β-collapsed-on-δ-as-agent). Mode is design-and-build (contract-shape judgment = design; markdown authoring = build); 9 ACs; ε is in-cycle per the cycle/401 cadence rule. Per `release/SKILL.md §3.8` configuration-floor clause, γ-axis and β-axis are capped at A- (γ/δ separation absent; β-α collapse acknowledged per Rule 7).

## Risks and forecasts

- **R1: AC4 no-duplication audit is the load-bearing rigor check.** The temptation when authoring a contract document that mirrors CDR.md's structure is to also mirror CDR.md's *prose*. Each section of CDS.md must be authored in CDS's voice, citing the kernel rather than restating it. β audit is required.
- **R2: name-overpromise in Sub 3/4/5 references.** CDS.md will reference Sub 3/4/5 destinations frequently (per Judgments 1–6). The references must be framed as "Sub N owns operational detail; CDS.md owns the contract" — not "Sub N will say X" (which would over-commit Sub N's authoring).
- **R3: scope creep into lifecycle/artifacts/review/gate detail.** The temptation to "while I'm here, just sketch the artifact contract" is the failure mode per #407 active design constraint. The Sub-N-vs-Field-M lines in Judgments 1–6 are the discipline; β verifies no operational detail leaked.
- **R4: gap-class rename (cdd-*-gap → cds-*-gap) coordination.** Field 5's gap-class rename is a substantive doctrine decision. β audits: is this a CDS.md-Field-5-territory decision or a Sub-5-territory decision? α's reading: the gap-class *names* are Field-5-territory (they are the taxonomy); the per-class operational *workflow* is Sub-5-territory. The naming pattern follows ROLES.md §4b.3 `{protocol}-{axis}-gap` convention, so `cds-*-gap` is mechanically derivable.
- **R5: receipt-shape commitment.** Field 3 commits CDS to the `#CDSReceipt` schema at `schemas/cds/receipt.cue`. The schema already exists per #388; CDS.md cites it. The Field 3 description must not re-state the schema's required-field set (those live in the schema file); it names the load-bearing artifact and points at the schema for the typed shape.

## Plan order

1. ✅ Read all source files (#407 body, #403 body, cnos.cdr/CDR.md template, cnos.cds/SKILL.md current state, extraction-map.md, ROLES.md §3 + §4a + §4a.2 + §4b, schemas/cdd/README.md §Architectural choice, CDD.md, CCNF/COHERENCE-CELL/RECEIPT-VALIDATION preambles for citation-only).
2. ✅ Branch `cycle/407` from `origin/main` HEAD (`987acd04`).
3. ✅ Author gamma-scaffold.md (this file).
4. Author D1 (`skills/cds/CDS.md`) — section by section per the large-file authoring rule (manifest header + completed-marker resumption protocol).
5. Author D2 (SKILL.md edits) — remove forthcoming note; add CDS.md as Step 1; update Rule section.
6. Optional: remove `.gitkeep` (per AC9; either is acceptable).
7. Mechanical AC check (`grep` count for headings; `grep` count for fields; phrase presence; `git diff` empty for cdd/cdr; `find` for role overlays).
8. α commit (role tag `α-407`).
9. Self-coherence sweep (`self-coherence.md`).
10. β-collapsed review (`beta-review.md`) — AC4 no-duplication audit is the rigor focus.
11. Close-outs (α, β, γ).
12. cds-iteration courtesy stub (or `cdd-iteration.md` — naming-conservation: the iteration *file* keeps `cdd-iteration.md` for this cycle since Sub 5 has not landed the rename to `cds-iteration.md`).
13. INDEX.md row (per cycle/401 courtesy convention).
14. β+γ+cdd-iteration commit (role tag `β-407` blended with `γ-407`).
15. γ-INDEX commit (role tag `γ-407`).
16. Push `cycle/407` to origin.
17. Report back to operator with branch name, commits, AC summary, and merge instruction. **Do NOT merge to main.**

## CDD Trace

| Step | Artifact | Skills loaded | Decision |
|------|----------|---------------|----------|
| 0 Observe | #407 body, #403 body, CDR.md, ROLES.md, schemas/cdd/README.md, CDD.md, extraction-map.md, cycle/406 close-outs | cdd, design, issue/contract, issue/proof | Inputs read; implementation contract pinned in #407 body |
| 1 Select | cnos#407 | — | Sub 2 of #403; first cycle of the post-Sub-1 wave; gates Subs 3–5 |
| 2 Branch | `cycle/407` | cdd | Branched from `origin/main` (HEAD `987acd04`) per CDD §4.2 / #407 dispatch |
| 3 Bootstrap | `.cdd/unreleased/407/` | cdd | Cycle dir created |
| 4 Gap | this file | — | Named: CDS.md does not exist; Subs 3–5 cannot dispatch against destinations without it; SKILL.md's "forthcoming" note describes a v0.1 state that #407 closes |
