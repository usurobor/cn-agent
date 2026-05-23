# v3.82.0 — CCNF package-architecture baseline

## Outcome

Coherence delta: the cnos protocol crosses a stable architecture boundary. CDD is the compact CCNF kernel; CDS, CDR, and cnos.handoff are v0.1-complete peer packages around it; V (`cn cdd verify`) is executable; domain evidence has homes. Per-package READMEs now reflect v0.1-complete reality; the canonical CCNF spine (`CDD.md` at 160 lines) is unchanged in this release.

This is a **release-hygiene cycle**, not a feature cycle. The bump is documentation + version-string; the architecture boundary was achieved by the four prior waves (#402 / #403 / #376 / #404) and the surrounding essays + sigma activation; v3.82.0 marks the boundary as the **package-architecture baseline** under which all future protocol evolution must justify itself.

## Why it matters

The pre-v3.82.0 state had four entangled risks structurally present:

1. **Stale READMEs misrepresenting v0.1 status.** `cnos.cds/README.md` claimed "v0.1 skeleton" + "Pending Sub 2" + "Forthcoming surfaces" while CDS.md was a 3,588-line canonical doctrine surface with all seven subs of #403 closed. `cnos.cdr/README.md` claimed "Role overlays (Sub 3) and empirical-anchor doc (Sub 4) are in flight" while both had shipped. Readers landing on these packages got materially wrong status signals — the failure mode is **package-discovery drift**: a contributor evaluating whether CDS/CDR is ready to depend on would conclude "no, it's a skeleton" when the reality is "yes, v0.1 is complete and load-bearing."

2. **Protocol-expansion momentum without a release boundary.** The #402 / #403 / #376 / #404 waves landed in rapid succession; the #405 (CCNF-O + TSC) roadmap was already filed; Track A1 (#421) had just closed. Without a release boundary the natural next move would be to dispatch Track A2 + Track B1 — but the operator's 2026-05-22 directive ruled that **the current architecture is stable enough to deserve a tagged baseline, and field evidence is more valuable than further theoretical refinement**. v3.82.0 is the structural pause point.

3. **CCNF kernel + per-protocol realization + handoff wire-format separation lacked a citable version.** The kernel/realization/wire-format three-way split is what allows CDS, CDR, and future c-d-X protocols to evolve independently. Naming v3.82.0 as the "CCNF package-architecture baseline" is the first version that pins this split as the canonical structure rather than as an in-progress reorganisation.

4. **No declared stop condition.** Without one, the protocol-evolution thread would keep pulling. The Stop condition at the bottom of this release declares the pause and names what comes next (field application + memory-return testing, not more theory).

## Includes

This release packages the following landed work into the v3.82.0 baseline:

### Protocol packages (v0.1 complete)

- **CDD — compact CCNF kernel ([cnos#402](https://github.com/usurobor/cnos/issues/402); closed).** `CDD.md` is the 160-line CCNF spine. Pre-#402 software-specific content was extracted to CDS; pre-#404 handoff-wire-format content was extracted to cnos.handoff. The kernel is now the generic recursive coherence-cell algorithm and nothing else.
- **CDS — Coherence-Driven Software v0.1 ([cnos#403](https://github.com/usurobor/cnos/issues/403); closed via [cnos#411](https://github.com/usurobor/cnos/issues/411)).** Seven-sub wave: #406 (package skeleton + extraction map), #407 (CDS.md = 3,588-line six-field instantiation contract), #408 (§Selection + §Lifecycle migration + selection/ thin overlay), #409 (§Coordination + §Artifact migration), #410 (§Mechanical + §Review + §Gate + §Closure + §Assessment + §Retro + §Non-goals + §Large-file migration + lifecycle/ thin overlay), #411 (CDD.md "pending cds extraction" marker sweep), #412 (empirical-anchor-cdd.md). Per-role overlays (`skills/cds/{alpha,beta,gamma,delta,epsilon}/`) deferred to v1.
- **CDR — Coherence-Driven Research v0.1 ([cnos#376](https://github.com/usurobor/cnos/issues/376); closed).** Four-sub wave: #390 (CDR.md = 616-line six-field instantiation contract + architecture-choice declaration), #394 (package skeleton + loader), #395 (five role overlays at `skills/cdr/{alpha,beta,gamma,operator,epsilon}/`), #396 (empirical-anchor-cph.md mapping `usurobor/cph` practice).
- **cnos.handoff — Inter-agent / inter-activation / inter-repo wire-format doctrine v0.1 ([cnos#404](https://github.com/usurobor/cnos/issues/404); closed via [cnos#420](https://github.com/usurobor/cnos/issues/420)).** Six-sub wave: #415 (package skeleton + extraction map), #416 (cross-repo doctrine — 644-line wholesale move), #417 (dispatch + 7-axis implementation-contract + δ-inward-membrane), #418 (mid-flight rescue + artifact-channel rules), #419 (cdd-iteration receipt-stream + INDEX.md aggregator), #420 (cross-reference sweep + v0.1-complete tightening).

### Essays + steering

- **`docs/gamma/essays/CCNF-AND-TYPED-TRUST.md`** — typed-trust precursor essay; positions CCNF as the substrate trust grammar.
- **`docs/gamma/essays/DECREASING-INCOHERENCE.md` ([cnos#414](https://github.com/usurobor/cnos/issues/414); closed).** Steering essay; per-shipment artifact contract for measuring coherence delta direction.

### Cross-repo activation

- **Sigma activation bundle staged + operator-applied ([cnos#413](https://github.com/usurobor/cnos/issues/413); closed).** Bundle at `cnos:.cdd/iterations/cross-repo/cn-sigma/sigma-activation-2026-05-22/`; operator-applied to the `cn-sigma` main branch.

### Roadmap (filed, not executed)

- **CCNF-O / TSC steering roadmap ([cnos#405](https://github.com/usurobor/cnos/issues/405); open).** Track A1 (CCNF-O survey + name-pick + sub-issue queue) shipped at [cnos#421](https://github.com/usurobor/cnos/issues/421); the survey doc at `docs/gamma/design/ccnf-o-track-a1-survey.md` (382 lines) carries five pinned decisions (name = CCNF-O; 20-surface inventory matrix; 6 higher-level forms classified; Track B1 dispatches in parallel with Track A2; package = `cnos.ccnf-o`). **Tracks A2–A6 + Track B1–B6 are explicitly NOT executed in this release.** They are dispatchable post-tag but deliberately deferred — see Stop condition.

## Does NOT include

The following surfaces were considered and explicitly excluded from v3.82.0's scope:

- **TSC report attachment.** No `TSCReport` schema; no attachment surface; no CHANGELOG-vs-receipt cross-check. Deferred to Track B1 of #405 (dispatchable post-tag).
- **`IssueProposal.v1` schema.** No `schemas/ccnf-o/issue-proposal/`; no proposal-pipeline CUE; no proposal-to-issue compiler. Deferred to Track A3 + B2 of #405.
- **`RiskPolicy.v1` schema.** No `schemas/ccnf-o/risk-policy/`; no risk-axis declaration; no policy-gate enforcement. Deferred to Track B5 of #405.
- **CCNF-O schemas.** No `schemas/ccnf-o/` directory of any kind. The 20-surface inventory matrix in `docs/gamma/design/ccnf-o-track-a1-survey.md` lists what Tracks A2–A6 will type, but none are typed in v3.82.0.
- **Field-trial results.** No empirical CDS/CDR field-cycle results. The empirical-anchor docs (`cnos.cds/docs/empirical-anchor-cdd.md`, `cnos.cdr/docs/empirical-anchor-cph.md`) map *current* practice onto v0.1 structure; they do not report *future* field-trial outcomes.
- **Wave-executor / new runtime / autonomy work.** No L5/L6 autonomous-loop runtime; no wave-executor scheduler; no `cn cdd verify` extensions beyond the current implementation; no harness changes; no release-effector changes.

## Added

- **CDD — compact CCNF kernel** ([#402](https://github.com/usurobor/cnos/issues/402)): `CDD.md` reduced to the 160-line CCNF spine; software-specific content extracted to CDS; handoff-wire-format content extracted to cnos.handoff.
- **CDS v0.1 — Coherence-Driven Software** ([#403](https://github.com/usurobor/cnos/issues/403); subs #406–#412): canonical doctrine surface + extraction map + operational sub-area overlays + empirical-anchor doc.
- **CDR v0.1 — Coherence-Driven Research** ([#376](https://github.com/usurobor/cnos/issues/376); subs #390/#394/#395/#396): canonical doctrine surface + loader + five role overlays + empirical-anchor doc.
- **cnos.handoff v0.1 — wire-format doctrine** ([#404](https://github.com/usurobor/cnos/issues/404); subs #415–#420): loader + HANDOFF.md + five per-surface sub-skills + extraction map.
- **`docs/gamma/essays/DECREASING-INCOHERENCE.md`** ([#414](https://github.com/usurobor/cnos/issues/414)): steering essay; per-shipment artifact contract.
- **`docs/gamma/essays/CCNF-AND-TYPED-TRUST.md`**: typed-trust precursor essay.
- **`docs/gamma/design/`** new directory with README ([#421](https://github.com/usurobor/cnos/issues/421)): home for survey + decision-class documents.
- **`docs/gamma/design/ccnf-o-track-a1-survey.md`** ([#421](https://github.com/usurobor/cnos/issues/421)): 382-line CCNF-O Track A1 survey; five pinned decisions; six refined sub-issue paragraphs; ten deferred questions.
- **Sigma activation bundle** ([#413](https://github.com/usurobor/cnos/issues/413)): staged at `cnos:.cdd/iterations/cross-repo/cn-sigma/sigma-activation-2026-05-22/`; operator-applied to cn-sigma main.

## Changed

- **`src/packages/cnos.cds/README.md`** ([#422](https://github.com/usurobor/cnos/issues/422); this release): replaced "v0.1 skeleton" + "Pending Sub 2" + "Forthcoming surfaces" + "in flight" framing with v0.1-complete status; landed sub-issue references (#406–#412); wave-shape narrative.
- **`src/packages/cnos.cdr/README.md`** ([#422](https://github.com/usurobor/cnos/issues/422); this release): replaced "v0.1 skeleton" + "Role overlays (Sub 3) and empirical-anchor doc (Sub 4) are in flight" framing with v0.1-complete status; landed sub-issue references (#390/#394/#395/#396); wave-shape narrative.
- **`VERSION`** ([#422](https://github.com/usurobor/cnos/issues/422); this release): bumped `3.81.0` → `3.82.0`.

## Removed

None. v3.82.0 removes no surfaces, no skills, no schemas. The release is additive (the four package families and essays already landed across #402–#421) plus documentation-coherence (this cycle's README rewrites) plus the version-bump boundary.

## Validation

This is a docs + version-bump release. The validation surfaces are:

- **Pre-tag (β / γ side, this cycle):** 11 acceptance criteria mechanical pass-set per [cnos#422](https://github.com/usurobor/cnos/issues/422); see `.cdd/releases/3.82.0/422/self-coherence.md` post-merge for the AC1–AC11 pass-set; CCNF kernel byte-identical to origin/main (`git diff origin/main -- src/packages/cnos.cdd/skills/cdd/{CDD.md,COHERENCE-CELL.md,COHERENCE-CELL-NORMAL-FORM.md}` → 0 lines); no new schemas, packages, skill content, or runtime/harness changes.
- **Post-tag (δ side, operator):** `scripts/release.sh` produces the bare `3.82.0` tag (no `v` prefix per `cnos.cdd/skills/cdd/release/SKILL.md §2.6` "tag naming convention"); release CI builds the binaries; `cn --version` reports `3.82.0`; `scripts/check-version-consistency.sh` passes (VERSION = cn.json = package manifests).
- **Post-deploy:** the release is a no-op in functional surface — no new ops, no new skills, no new runtime entry points; what changed is documentation coherence. The validation is therefore *documentation-coherence validation*: a reader landing on `cnos.cds/README.md` or `cnos.cdr/README.md` gets accurate v0.1-complete status; the `Status` section names the closed wave date (2026-05-22); the `Package Structure` section lists landed surfaces with sub-issue references.

What was proven coherent by the validation: every v0.1-complete claim in cds/cdr READMEs is grounded in a closed sub-issue and a landed file. No anticipatory "Forthcoming" / "Pending" / "in flight" language remains as a current-state claim.

## Known Issues

- **`src/packages/cnos.cds/skills/cds/SKILL.md` carries stale "Pending Subs 3–5" / "v0.1 caveat" wording** in its frontmatter description, "Load order" §"v0.1 status", "Role overlays" section, and "v0.1 caveat" section. The wording anticipates per-role overlays (`alpha/beta/gamma/delta/epsilon`) that did not ship in v0.1 (the wave landed operational sub-area overlays — `lifecycle/`, `selection/` — instead, per the extraction map's Sub 3/4/5 split). The SKILL.md is left unchanged in this release because [cnos#422](https://github.com/usurobor/cnos/issues/422) AC10 prohibits protocol-skill content changes beyond README/VERSION/RELEASE.md. This is a **post-v0.1 follow-up** for a separate cycle (likely the cycle that lifts CDS to v0.2 with actual per-role overlays, or a small cleanup cycle that updates the SKILL.md framing to match what shipped). It is **not blocking** for the v3.82.0 baseline because the canonical doctrine in `CDS.md` is authoritative and the loader's `Conflict rule` declares `CDS.md` governs in any conflict.
- **CCNF-O / TSC steering roadmap ([cnos#405](https://github.com/usurobor/cnos/issues/405)) remains open.** Tracks A2–A6 + B1–B6 are filed but not executed. This is deliberate; see Stop condition.

## Stop condition

**After the v3.82.0 tag is pushed, pause protocol evolution.** The current package-architecture (CDD kernel + CDS/CDR domain realizations + cnos.handoff wire-format) is the baseline against which all future protocol changes must justify themselves. The next phase is:

1. **Field application of CDS** — actual software cycles (not protocol cycles) run under the v0.1 CDS doctrine; the per-role overlays Subs 3–5 anticipated are likely to be authored only when field evidence shows what they need to enforce.
2. **Field application of CDR** — actual research cycles run under the v0.1 CDR doctrine; the empirical-anchor doc's cph mapping gets tested against real wave evidence.
3. **Memory-return testing of cnos.handoff** — the cross-repo + dispatch + receipt-stream wire formats get exercised across multiple agents and activations; the failure modes that surface inform what handoff v0.2 needs.

**Do NOT dispatch [cnos#405](https://github.com/usurobor/cnos/issues/405) Tracks A2 / B1 or any other #405 sub** until the field-evidence pause has produced concrete signal that a typed orchestration grammar (CCNF-O) or a TSC report attachment surface is the most-leverage next move. The Track A1 survey at `docs/gamma/design/ccnf-o-track-a1-survey.md` is *available* for the post-pause cycle to dispatch against; it is *not* a dispatch instruction in this release.

The CDD's §Selection function will, post-tag, naturally prefer field-trial cycles over further protocol-extension cycles because field-trial evidence is the highest-leverage gap for a freshly-baselined protocol. The Stop condition is the structural expression of that selection-preference at the release boundary.

---

**Release cycle:** [cnos#422](https://github.com/usurobor/cnos/issues/422) — Release-hygiene v3.82.0 (CCNF package-architecture baseline).
**Tag (operator-side, post-merge):** `3.82.0` (bare; no `v` prefix).
**Release-effector:** `scripts/release.sh` per `cnos.cdd/skills/cdd/release/SKILL.md §2.6` + `release-effector/SKILL.md`.
