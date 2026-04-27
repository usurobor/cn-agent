# Writer Package

**Issue:** #278
**Version:** unset (not release-scoped; design-lock cycle)
**Mode:** MCA
**Active Skills:** `cnos.cdd/skills/cdd/design`, `cnos.core/skills/skill`, `cnos.core/skills/write`

---

## Problem

The four-essay doctrine sequence (CFA, EFA, JFA, IFA) ran four full triadic cycles using one stable agent pattern: roles α/β/γ, the `cnos.core/skills/write` skill governing the prose, cycle artifacts (cycle log, critiques, optional external observations) co-located with the essay. The pattern is repeatable. It exists nowhere except as evidence inside `docs/alpha/doctrine/`.

A future writer-shaped cycle (any cycle whose primary artifact is a written document with structural-coherence claims rather than running code) cannot inherit the pattern as constraint. It must reconstruct it from logs. Reconstruction risks the IFA-named failure mode `cold-author evidence refusal` (assuming inherited artifacts are unavailable or unneeded) and the IFA-named failure mode `soft inheritance` (importing a prior commitment in vocabulary while losing it in operative effect).

The named incoherence: the four-cycle archive is the only existing memory surface for the writer pattern, and a memory surface that requires reconstruction at each use is not yet inheritance — it is imitation candidate material.

Evidence:
- Four cycles ran (`docs/alpha/doctrine/coherence-for-agents/`, `ethics-for-agents/`, `judgment-for-agents/`, `inheritance-for-agents/`) — pattern is empirically stable across four uses.
- IFA itself argues that doctrine grows by composing cycles, and the composition is inspectable only when prior named failures are carried forward as constraints (`INHERITANCE-FOR-AGENTS.md` §"Doctrine grows by composing cycles").
- `docs/alpha/doctrine/README.md` §"Skill" notes the write skill is the one common surface across all four essays — but the surrounding triadic structure has no package home.

## Constraints

The package must respect contracts already in force:

- **`cnos.core/skills/write/SKILL.md` is the prose-discipline authority.** It governs all four cycles and every future writer cycle. Its frontmatter and rules apply to docs, issues, specs, READMEs, essays, reviews, and status updates — not just essays. The writer package must depend on it without owning it.
- **`cnos.cdd` is the canonical method package for code-shaped substantial cycles.** `CDD.md` §1.4 names α/β/γ roles for substantial development. The writer package and `cnos.cdd` are peer tools (issue Non-Goal #4); the writer package may not redefine, restate, or shadow CDD's selection rules, lifecycle, artifact contract, or dispatch format.
- **Doctrine essays at `docs/alpha/doctrine/` are frozen artifacts.** Each essay travels with its cycle log and critiques in the essay's own folder (doctrine README §"Cycle artifacts"). The writer package cannot move them (issue Non-Goal #2) and cannot copy them without violating `write/SKILL.md` rule 3.3 (one home per stable fact).
- **Issue Non-Goal #1 forbids writing `cnos.writer/SKILL.md` before design lock.** This design closes the design phase only; package authoring is the next cycle.
- **Issue Non-Goal #3 forbids generalizing the writer to non-cnos repos.** The package is cnos-scoped.
- **`cnos.core/skills/skill/SKILL.md` rule 3.10 requires source/packaged copy sync.** If any writer-related skill is added to `src/`, the package copy must stay in sync.

## Parallel Dependency: CTB Language Spec Reconciliation

Issue #278 §"Parallel dependency" requires that any conflict between this design and the parallel `cnos` agent language reference surface in the design doc, not silently. The reference is `docs/alpha/ctb/LANGUAGE-SPEC.md` (CTB Language Spec v0.1, draft-normative, dated 2026-04-26).

**LANGUAGE-SPEC's bearing on this design.** §0 names the signature surface defined in §2 — `name`, `description`, `artifact_class`, `kata_surface`, `governing_question`, `scope`, `visibility`, `triggers`, `inputs`, `outputs`, `requires`, `calls`, `calls_dynamic`, `runs_after`, `runs_before`, `excludes` — as fully realized today in `cnos.cdd` and **prescriptive** for cnos.core skills (which carry the older minimal frontmatter). §10 lists the well-formedness constraints a conformant skill must satisfy: every skill must declare the five required-fields per §2.1, must declare `scope` per §3.2, must declare `inputs` and `outputs` per §2.3, and public skills must declare `visibility: public` per §4.1.

**Reconciliation outcome: alignment with LANGUAGE-SPEC §2 prescriptive surface, no divergence.** Because this design claims peer status with `cnos.cdd` (which already realizes the §2 signature in full per LANGUAGE-SPEC §12.1–§12.2), the writer package's SKILL.md frontmatter targets the same fully-realized surface — not the older minimal shape. The §"Proposal" surface below names this explicitly. Specifically, every `SKILL.md` the authoring cycle creates declares all five §2.1 required fields, plus `scope`, plus `visibility`, plus structured `inputs`/`outputs`/`requires`, plus `calls` / `calls_dynamic` where dispatch occurs.

**Implication for the writer's relationship to `cnos.core/skills/write/SKILL.md`.** Per §5.2 of this design, `write/SKILL.md` stays in `cnos.core` and is referenced by the writer package, not owned by it. LANGUAGE-SPEC §12.3 names `write/SKILL.md` as an example where the spec is prescriptive, not yet realized — the skill's current frontmatter declares only `name`, `description`, `triggers`. The writer package must not author its dependency-by-reference in a way that requires `write/SKILL.md` to migrate first; the writer package depends on `write` at its current frontmatter shape today, and tracks `write`'s eventual migration as a separate matter (LANGUAGE-SPEC §0 says "Migration is not a precondition for adopting the spec elsewhere"). This is consistent with both documents.

**No divergence from LANGUAGE-SPEC is asserted by this design.** Future writer-package authoring cycles must verify each `SKILL.md` against LANGUAGE-SPEC §10 well-formedness and against §2 vocabulary. If a writer-specific need arises that the language spec does not provide for (e.g. a writer-only field like `governing_question` is declared but a writer-shape-specific field would help), the divergence must surface in the package's own normative spec at that point — not silently. This design names no such divergence today.

## Challenged Assumption

This design challenges the assumption that the triadic agent pattern is owned by `cnos.cdd` and that any other triadic work pattern must either run inside CDD or be ad-hoc. The doctrine cycles are evidence that a triadic pattern can be shaped by its material (prose with structural-coherence claims) and can produce its own lifecycle artifacts (cycle log, critiques, external observations) that differ from CDD's lifecycle artifacts (design, plan, tests, code, docs, gate, release). The challenged assumption is not that CDD is wrong; it is that triadic structure is a CDD-exclusive shape rather than a pattern of which CDD is one specialization and the writer is another.

The package therefore claims peer status with `cnos.cdd`, not derivation from it.

## Impact Graph

**Downstream consumers (who reads / loads / depends on the writer package):**
- Future writer-shaped cycles in cnos: any cycle whose primary artifact is a doctrine essay, a structural specification, a long-form review, or a written-form coherence claim. None exist yet outside the four doctrine cycles.
- The α / β / γ role sessions of those future cycles: each role loads the writer package's role skill at dispatch time, mirroring how CDD roles load `alpha/SKILL.md` etc.
- `docs/alpha/doctrine/README.md` §"Skill" — currently points to `write/SKILL.md`; if the writer package adds load-order or role context, this README is the natural place to link to the package without restating it.

**Upstream producers (what the writer package depends on):**
- `cnos.core/skills/write/SKILL.md` — the prose-discipline authority. Stays in `cnos.core` (see §5.2 below).
- The four doctrine essays and their cycle artifacts — referenced as primary evidence, not copied.
- `cnos.cdd` — peer reference for triadic-shape language already established (roles, dispatch, close-out, role separation). The writer package adopts the role names α/β/γ from `CDD.md` §1.4 to avoid coining a parallel vocabulary.

**Copies and embeddings (where the same fact must not duplicate):**
- The `write` skill rules must not be restated inside the writer package. Any rule the package needs is cited from `write/SKILL.md` by section number.
- The four-essay inherited failure modes (`docs/alpha/doctrine/README.md` §"Inherited failure modes" lines 38–46) must not be restated in the writer package's own SKILL.md or normative spec. They are cited by source path and section.
- Role definitions (α/β/γ) must not be redefined; the writer package adopts them from `CDD.md` §1.4 with a one-line note that the writer's β does not own merge/release (see §6 Divergences).

**Authority relationships:**
- On prose discipline: `write/SKILL.md` governs.
- On triadic role definitions: `cnos.cdd/skills/cdd/CDD.md` §1.4 governs.
- On the writer's lifecycle (cycle log, critiques, external observations, lock): the writer package's own normative spec governs (introduced by this design as `cnos.writer/skills/writer/WRITER.md`, analog of `CDD.md`).
- On individual essays' contents: the essay file at `docs/alpha/doctrine/{essay}/` governs; cycle logs are evidence about the essay, not authority over it.

## Scope Decisions

This section addresses the four scope questions named in issue #278 AC1. Each decision has three components: chosen path, alternative considered, structural reasoning.

### 5.1 Package timing — now vs. after one non-doctrine cycle

**Chosen path:** Package now. Mark the package as "scaffolded from the four doctrine cycles; contestable at first non-doctrine use." The package is created as a memory surface; the first non-doctrine writer cycle is the inheritance test, and the test cannot run if the surface does not exist when that cycle starts.

**Alternative considered:** Wait for one non-doctrine writer cycle to confirm the pattern generalizes, then package.

**Structural reasoning:** IFA's contestability commitment (`INHERITANCE-FOR-AGENTS.md` §"Inheritance prevents dogma by preserving contestability") establishes that inheritance carries forward named failure modes as contestable constraints, not as fixed authority. Withholding the package until a future cycle "verifies" the pattern would be reverse soft-inheritance: pretending to honor IFA's partial-inheritance discipline by treating the four cycles as too special to reuse, when IFA itself argues that they should be reused as constraints subject to test at the point of use. The package is contestable; the wait-and-see alternative substitutes a different incoherence (institutional knowledge stays in cycle logs requiring reconstruction at each use, which IFA names as cold-author evidence refusal risk per `IFA-cycle-log-gamma.md` §"IFA-Specific Risks Evaluated" → "Cold-Author Drift: Caught and repaired" bullet) for the one this issue closes.

### 5.2 Write skill location — stay in `cnos.core` vs. move to `cnos.writer`

**Chosen path:** `cnos.core/skills/write/SKILL.md` stays where it is. The writer package depends on it by reference; it does not own it.

**Alternative considered:** Move `write` to `cnos.writer/skills/writer/write/` so the writer package self-contains its prose discipline.

**Structural reasoning:** `write/SKILL.md` declares its scope as "docs, issues, specs, READMEs, essays, reviews, and status updates" (frontmatter `description`). It applies to any cnos prose, not only writer cycles. Moving it would couple a general skill to a specific agent pattern, restricting reuse and creating a duplication-or-divergence risk for non-writer prose (CDD design docs, package READMEs, issue text). The structural parallel is `cnos.cdd`: `design/SKILL.md` is a CDD-specific lifecycle phase (located inside `cnos.cdd/skills/cdd/design/`), but `write/SKILL.md` is general-purpose and lives in `cnos.core`. Writer cycles use `write` the same way CDD cycles use it — by reference, from `cnos.core`.

### 5.3 Inheritance discipline scope — writer-specific vs. general `cnos.core` candidate

**Chosen path:** The IFA-named inheritance failure modes (`soft inheritance`, `cold-author evidence refusal`) are flagged as a `cnos.core` candidate skill in this design but are not authored as a separate skill in this cycle. The writer package references them by citation to `INHERITANCE-FOR-AGENTS.md` and `IFA-cycle-log-gamma.md` §"Inherited Failure Modes" rather than restating them locally. Promotion of the inheritance discipline into a standalone `cnos.core/skills/inherit` (or equivalent) is deferred until at least one non-writer cycle (CDD or otherwise) demonstrates the same failure modes apply outside the writer's material.

**Alternative considered:** Encode `soft inheritance` and `cold-author evidence refusal` inside `cnos.writer` as writer-specific failure modes.

**Structural reasoning:** The IFA cycle separates inherited cross-cycle failure modes (binding on future cycles) from material-specific cycle lessons (recorded in the archive, not binding). `IFA-cycle-log-gamma.md` lists "Soft Inheritance" and "Cold-Author Evidence Refusal" under `## Inherited Failure Modes (Added for Future Cycles)` and lists "Single-Step Composition Claim" under `## IFA-Specific Cycle Lessons (Not Inherited)`. The writer-specific alternative would re-create the very mistake IFA names: importing a load-bearing, cross-cycle constraint as if it were material-specific, which IFA classifies as soft inheritance ("pretending to inherit while changing what the inheritance requires"). The chosen path keeps the constraint in its widest scope (citation accessible to any cnos cycle) while honoring the partial-inheritance rule (no premature generalization through standalone-skill creation when only writer evidence exists).

### 5.4 Package content — doctrine references vs. copies

**Chosen path:** Reference. The four doctrine essays and their cycle artifacts stay at `docs/alpha/doctrine/`. The writer package links to them and cites specific sections; it does not embed copies.

**Alternative considered:** Copy the four essays (or their inherited-failure-mode summaries) into `cnos.writer/skills/writer/` so the package is self-contained.

**Structural reasoning:** Three constraints converge. (1) `cnos.core/skills/write/SKILL.md` rule 3.3 ("Say a fact once, then point to it") forbids restating stable facts across files. (2) `cnos.cdd/skills/cdd/design/SKILL.md` §4.3 ("Enumerate copies and embeddings") names duplication as an impact-graph hazard. (3) Issue #278 Non-Goal #2 explicitly forbids migrating essays. Copies would risk drift (the package's copy of "soft inheritance" diverges from IFA's wording at the point of pressure), which is itself an instance of `soft inheritance`. References preserve one home per stable fact and let future cycles re-read the source artifact when the inherited constraint is contested.

## Divergences from cnos.cdd structure

Reference structure (`src/packages/cnos.cdd/`):
```
cn.package.json
commands/cdd-verify/
skills/cdd/
  SKILL.md          (loader entrypoint)
  CDD.md            (normative spec)
  alpha/ beta/ gamma/ operator/   (role skills)
  issue/ design/ plan/ review/ release/ post-release/   (lifecycle sub-skills)
```

Proposed writer structure (`src/packages/cnos.writer/`):
```
cn.package.json
skills/writer/
  SKILL.md          (loader entrypoint)
  WRITER.md         (normative spec)
  alpha/ beta/ gamma/   (role skills)
  cycle-log/ critiques/ external-observations/   (lifecycle sub-skills)
```

Each divergence is named with reasoning. Where no divergence applies, that is asserted explicitly.

**D1. No `commands/` directory.** `cnos.cdd` ships `cdd-verify` for cycle-artifact completeness checks. The writer package ships no commands in this cycle. Reasoning: writer-cycle correctness is judgment-bound (does the prose carry its argument? does β verify cross-cycle citations against source? does the closing line earn its body?). `CDD.md` §6.2 already classifies design coherence as judgment-bearing. A future `writer-verify` could mechanize artifact-presence checks (cycle-log present, critiques present, essay locked) but is not part of this design.

**D2. No `operator/` role skill.** `cnos.cdd` defines δ (operator) as a fourth boundary participant for platform gates (merge, tag, deploy, branch deletion). Writer cycles produce no merges, tags, or deploys (see D5). The operator role for writer cycles collapses to "host that routes between α/β/γ sessions." Reasoning: `CDD.md` §1.4 explicitly notes that δ "is not a fourth triad role — δ is the boundary between the triad-as-whole and the platform." The writer's platform surface is the filesystem and the doctrine folder; no platform action requires operator-as-role separation in the writer's lifecycle.

**D3. No `release/`, `gate/`, or `post-release/` lifecycle sub-skill.** `cnos.cdd` requires β to merge, tag, deploy, and triage post-release (steps 8–10 + γ steps 11–13 of `CDD.md` §4.1). Writer cycles end at "essay locked" (CFA-cycle-log §"Convergence" → "the cycle's most disciplined act was stopping"). β approves; the artifact is frozen in its folder; no further action is required. Reasoning: prose has no deployable surface. The writer's analog of CDD steps 8–10 collapses to one act — β's approve verdict — recorded in the critiques file. Adding release/gate/post-release stubs would import CDD-shaped overhead with no work for it to do.

**D4. New lifecycle sub-skills `cycle-log/`, `critiques/`, `external-observations/`.** These are first-class artifacts in writer cycles (every doctrine cycle has a cycle log; every doctrine cycle has critiques; EFA and the other cycles each had to decide whether to produce external-observations). `cnos.cdd` does not have direct analogs: `review/` is closer to `critiques/` but CDD reviews live as PR comments while writer critiques live as a versioned `.md` file inside the essay's folder; CDD has no canonical analog of `external-observations/` at all. Reasoning: the artifact shapes diverge because the material diverges. Code review iterates against tests and CI; prose review iterates against the governing question and the write skill, with rounds preserved in the cycle log so the inheritance is inspectable later (per IFA's contestability rule).

**D5. No version directory bootstrap.** `cnos.cdd` §5.1 requires the first diff on a substantial branch to create `docs/{tier}/{bundle}/{X.Y.Z}/`. Writer cycles produce essays that lock at convergence and live forever in their named folder (e.g. `docs/alpha/doctrine/coherence-for-agents/`). They carry no `X.Y.Z` because they are not released as a versioned bundle. Reasoning: doctrine README §"Cycle artifacts" already encodes the convention "every cycle's artifacts live in the essay's folder." A version-directory layer would invent a release boundary the artifact does not have.

**D6. No CHANGELOG / RELEASE.md / tag.** Direct consequence of D3 + D5. Writer cycles do not enter the CHANGELOG ledger. Reasoning: same as D3 — no deployable surface, no release event.

**D7. No `pre-review` gate analog.** `cnos.cdd/alpha/SKILL.md` §2.6 requires a 10-row gate (rebased, CDD trace through step 7, schema audit, peer enumeration, harness audit, polyglot re-audit, CI green). Writer cycles have a different α-side gate, drawn from the cycle artifacts: source load verified (no cold-author evidence refusal), governing question explicit, draft does not claim unavailability without a verify step, prior-cycle inherited commitments cited rather than restated. Reasoning: the gate must match the failure modes the material actually has. CDD's gate matches code-shape failures (CI, schema drift, peer surfaces). The writer's gate must match prose-shape failures (cold-author refusal, soft inheritance, mimicked closure, comfort failure at the writing level).

**D8. β role does not own merge/release.** `cnos.cdd/CDD.md` §1.4 names β "Reviewer + Releaser." Writer β is "Judge" only. Reasoning: D3 + D5.

**D9. γ may be implicit in dyad cycles.** CFA-cycle-log §"What broke the oscillation" records that "the meta-observation arrived inside β's critique rather than from a separate observer, but it functioned as γ regardless." Some writer cycles run as α/β with γ-shape work absorbed into β's review surface; others (JFA, IFA) split γ out into a separate `cycle-log-gamma.md` file. `cnos.cdd` requires γ for substantial cycles. Reasoning: the writer's third surface can be a separate observer or an external observation; what is load-bearing is that some surface outside the producer/judge dyad name a pattern they cannot exit (per CFA). This divergence does not weaken the triadic structure; it acknowledges the dyad-without-meta failure mode CFA discovered and lets the third surface take whatever local form serves it.

**D10. Triad-internal communication is git-only — per-cycle subdir + per-role file.** All role-to-role communication inside the triad — α drafts, β verdicts, γ cycle observations, external-observations evaluation records, α dispositions — lives in committed files. The convention is "per-cycle subdir, per-role file." Two concrete forms apply:

- **CDD-shape cycles** (e.g. this design's authoring cycle, issue #278) use `cnos.cdd`'s existing convention `.cdd/unreleased/{role}/{N}.md` mid-cycle, promoted to `.cdd/releases/{X.Y.Z}/{role}/CLOSE-OUT.md` per CDD.md §5.3a Artifact Location Matrix at release time. Each per-issue file is a living record: it accumulates round-by-round review history during the cycle and becomes the role's close-out at the end. (Verified by inspection: `.cdd/unreleased/alpha/268.md`, `.cdd/unreleased/beta/268.md`, `.cdd/unreleased/gamma/268.md` are the in-flight pattern; `.cdd/releases/3.59.0/alpha/CLOSE-OUT.md` etc. are the post-release pattern. Historical practice through 3.58.0 used `.cdd/releases/{X.Y.Z}/{role}/{N}.md` filenames — e.g. `.cdd/releases/3.58.0/alpha/262.md` — but §5.3a standardized on `CLOSE-OUT.md` at 3.59.0 and the per-issue release-form is legacy.)
- **Writer-shape cycles** (future cycles using the `cnos.writer` package) use the doctrine-folder convention: `docs/{...}/{essay}/<essay>-cycle-log.md` (α-side), `docs/{...}/{essay}/<essay>-critiques.md` (β-side), optional `docs/{...}/{essay}/<essay>-cycle-log-gamma.md` (γ-side), optional `docs/{...}/{essay}/<essay>-external-observations.md` (third-surface). Each essay's folder IS its per-cycle subdir; the role-prefix is in the filename rather than a sub-directory. Established by the four doctrine cycles.

GitHub PR comments, issue comments, and review threads are operator-facing scaffolding, not triad-internal authority. `cnos.cdd` historically defaulted β review output to PR comments (`gh pr view --comments` polling per `CDD.md` §Tracking; `review/SKILL.md` §7.1). The writer diverges by elevating the committed file to canonical authority. Reasoning: (a) **inheritance** — LB9 already requires cycle artifacts to travel with the essay; PR comments do not travel, they live on a platform whose API can change, deprecate, or paywall. The doctrine cycles already practiced this discipline (cycle-log, critiques, external-observations are committed files; LB11 names the inheritance). The CDD package's `.cdd/unreleased/` directory practices the same discipline for CDD cycles. (b) **soft-inheritance prevention** — a verdict committed as a file at the point of pressure is harder to silently soften than a PR comment that can be edited or buried in scrollback. (c) **inspectability** — any future cycle reads the cycle's record directly from the repo with no platform dependency, satisfying LANGUAGE-SPEC §9.4 reproducibility for the cycle's communication record. (d) **effect-boundary cleanliness** — per LANGUAGE-SPEC §9.2, the language returns plans, the runtime executes them; PR comments cross the effect boundary onto a third-party platform, git commits stay on the data side. The writer's communication surface stays inside the language layer.

**No divergence — explicitly:**

- **Loader entrypoint pattern.** Writer adopts `skills/writer/SKILL.md` as loader and `WRITER.md` as normative spec, mirroring `cnos.cdd/skills/cdd/SKILL.md` + `CDD.md` exactly. Same conflict rule (`WRITER.md` governs on disagreement).
- **Role names α/β/γ.** Adopted from `CDD.md` §1.4 verbatim; no parallel vocabulary.
- **Source/packaged-copy sync rule.** `skill/SKILL.md` rule 3.10 applies unchanged.
- **Frontmatter discipline for skills.** `skill/SKILL.md` rule 3.1 applies unchanged.
- **Triadic role separation.** α produces, β judges, γ observes the cycle. The dyad cannot collapse into one agent for substantial cycles. Adopted from `CDD.md` §1.4 with the dyad-with-implicit-γ relaxation in D9.

## Load-Bearing Constraints from Cycle Artifacts

This section names the constraints that govern the writer package's structure, with citations to the cycle artifacts that establish each one. Per IFA, these are inherited as contestable constraints — the package surfaces them at the point of use, and a future writer cycle may reject or revise an inherited constraint by naming the structural reason and leaving the rejection inspectable.

**LB1. The dyad cannot reliably exit oscillation without a third surface.**
Citation: `docs/alpha/doctrine/coherence-for-agents/CFA-cycle-log.md` §"What broke the oscillation": "β's third critique named the pattern explicitly… The naming was the third surface. Once it existed, neither role could continue the oscillation without acknowledging it."
Constraint on the package: γ (or an external-observation surface) must be declared structural, not optional ceremony. The writer's triadic shape derives from this CFA discovery, not from imitation of CDD. `WRITER.md` must require that every writer cycle name where its third surface lives — separate γ session, external observations, or β explicitly absorbing the γ-shape work in its critique.

**LB2. Stopping is the most disciplined act; further iteration after structural convergence is comfort failure.**
Citation: `docs/alpha/doctrine/coherence-for-agents/CFA-cycle-log.md` §"Convergence": "The cycle's most disciplined act was stopping… continuing to choose between micro-variants past that point would have been the comfort failure in its purest form."
Constraint: the writer package's closure rule must be structural ("no remaining structural grievance"), not aesthetic ("the prose could be tighter") and not budget-bound ("two rounds maximum"). β's approve verdict in `critiques/` must be attached to a named convergence, not a tolerance threshold.

**LB3. External observations enter as findings to evaluate, not as verdicts to execute.**
Citation: `docs/alpha/doctrine/ethics-for-agents/EFA-cycle-log.md` §"Final round — two surgical reinsertions" + `INHERITANCE-FOR-AGENTS.md` §"Inheritance is not imitation": "The cycle did not imitate openness by obeying every reviewer. It inherited boundary discipline by evaluating each observation against the governing question and the write skill."
Constraint: when an `external-observations/` artifact exists in a writer cycle, the role (typically α with β verification) must record which observations were incorporated, which were rejected, and the structural reason in each case. Mass acceptance and mass rejection are both failure modes. `WRITER.md` must require this evaluation record.

**LB4. Inherited commitments hold without re-derivation.**
Citation: `docs/alpha/doctrine/judgment-for-agents/JFA-cycle-log-gamma.md` §"No findings hidden under approval" → "EFA's no-fixed-order commitment held. The essay denies fixed ordering at four sites and never reintroduces it through softer phrases" + `INHERITANCE-FOR-AGENTS.md` §"Doctrine grows by composing cycles".
Constraint: the writer package must reference inherited failure-mode lists by citation rather than restating them. If a future writer cycle's β finds the inherited commitment violated, β cites the source artifact and section, not the package's own restatement. This is the package's primary anti-drift mechanism: there is no second wording for the constraint to drift toward.

**LB5. Inheritance is partial — prepared attention, not prediction.**
Citation: `docs/alpha/doctrine/inheritance-for-agents/IFA-cycle-log-gamma.md` §"IFA-Specific Cycle Lessons (Not Inherited)": "This is a material-specific failure mode of essays making composition claims, and is recorded here for the archive rather than as a binding constraint on unrelated future cycles."
Constraint: `WRITER.md` must distinguish (a) cross-cycle inherited failure modes (binding on every writer cycle) from (b) material-specific cycle lessons (recorded in the archive but not binding on unrelated future cycles). The same partition that IFA's γ log uses for its own cycle must be the partition the package uses for the doctrine archive as a whole. Without this partition, soft inheritance becomes structurally unavoidable: a future cycle either imports too much (importing material-specific lessons as binding) or too little (treating cross-cycle constraints as material-specific).

**LB6. Soft inheritance is the central failure mode of multi-cycle work.**
Citation: `docs/alpha/doctrine/inheritance-for-agents/INHERITANCE-FOR-AGENTS.md` §"Inheritance can soften into drift" + `IFA-cycle-log-gamma.md` §"Inherited Failure Modes": "Pretending to inherit a prior commitment while softening its phrasing at the point of pressure (e.g., converting a structural constraint into a general preference)."
Constraint: the package itself is an act of inheritance from four cycles. Its structure must surface each inherited commitment at the point of pressure (β's review surface), not in ceremonial preamble. `WRITER.md` must require β to cite the relevant inherited failure mode whenever a finding rhymes with one of them; a finding that rhymes with `soft inheritance` but is recorded only as "phrasing tightened" is itself an instance of soft inheritance.

**LB7. Cold-author evidence refusal — assuming inherited artifacts are unavailable without verifying.**
Citation: `docs/alpha/doctrine/inheritance-for-agents/IFA-cycle-log-gamma.md` §"IFA-Specific Risks Evaluated" → "Cold-Author Drift: Caught and repaired" bullet: "α's initial failure to utilize available JFA artifacts was a strict failure of the contract. β's refusal to accept the artifact's absence forced the necessary derivation."
Constraint: the writer package's α role skill must require α to verify access to every named cycle artifact before treating it as unavailable. β's role skill must include "refuse to accept disclaimed-as-unavailable evidence without proof of verification attempt." This constraint operates at the dispatch boundary (α's intake step) and at the review boundary (β's verification step).

**LB8. β verifies cross-cycle claims against source artifacts before approval.**
Citation: `docs/alpha/doctrine/inheritance-for-agents/IFA-critiques.md` §"Verdict 02 — approve" lists six explicit verifications β performed: "JFA dyad round structure → matches JFA-critiques exactly; Round counts → match JFA-cycle-log-gamma verbatim in substance; No-fixed-order denial at multiple sites → matches JFA-cycle-log-gamma verification subsection; α-1's procedure-substitute phrase → appears verbatim in JFA-cycle-log-dyad; Case-capture and terminal-debt repairs → match JFA-critiques findings 2–4; Both composition steps → match CFA-cycle-log, EFA-cycle-log, EFA-external-observations, JFA-cycle-log-gamma."
Constraint: `WRITER.md` must include "verify cross-cycle claims against named source artifacts" as a β gate. β's approval is not valid without this verification record.

**LB9. Cycle artifacts travel with their essay; cycle artifacts are evidence about the essay, not authority over it.**
Citation: `docs/alpha/doctrine/README.md` §"Cycle artifacts": "Each essay travels with the record of the cycle that produced it. Every cycle's artifacts live in the essay's folder."
Constraint: the writer package's filesystem layout convention must keep `<essay>.md`, `<essay>-cycle-log.md`, `<essay>-critiques.md`, optional `<essay>-cycle-log-gamma.md`, optional `<essay>-external-observations.md` co-located. The package does not migrate the doctrine essays (Non-Goal #2); it formalizes the convention so future writer cycles use the same layout.

**LB10. Each essay has a single governing question.**
Citation: `cnos.core/skills/write/SKILL.md` §1.3 + every doctrine essay's "Governing question" frontmatter line (e.g. `INHERITANCE-FOR-AGENTS.md` line 9: "How does a doctrine grow across cycles without becoming dogma?").
Constraint: `WRITER.md` must require α's first artifact in any cycle to name the governing question explicitly and freeze it before drafting begins. Drift in the governing question mid-cycle is itself a finding β must hold for. Per CFA-cycle-log, the governing question is what the dyad's oscillation ultimately bends toward losing.

**LB11. Triad-internal communication in the doctrine cycles was practiced as git-committed files, not as platform comments.**
Citation: `docs/alpha/doctrine/README.md` §"Cycle artifacts": "Each essay travels with the record of the cycle that produced it. Every cycle's artifacts live in the essay's folder." Concrete artifacts: `CFA-cycle-log.md`, `CFA-critiques.md`, `EFA-cycle-log.md`, `EFA-critiques.md`, `EFA-external-observations.md`, `JFA-cycle-log-dyad.md`, `JFA-cycle-log-gamma.md`, `JFA-critiques.md`, `IFA-cycle-log-dyad.md`, `IFA-cycle-log-gamma.md`, `IFA-critiques.md`. Every β verdict, α disposition, γ observation, and external-observation evaluation in the four-cycle archive is a committed file, not a PR/issue comment.
Constraint: `WRITER.md` must encode this as the writer's communication surface (D10). The lifecycle sub-skills `cycle-log/`, `critiques/`, `external-observations/` (Proposal §"Proposal") are not optional artifacts on top of a PR-comment workflow — they are the canonical communication channel. PR/issue comments may exist as operator-visible scaffolding but carry no triad-internal authority. This constraint, combined with LB9 (cycle artifacts travel with the essay), is what makes inheritance inspectable across platform changes.

## Proposal

A new package `cnos.writer` that packages the triadic writer pattern as a peer of `cnos.cdd`.

**Communication surface (cross-cutting principle).** Triad-internal communication is git-only by design (D10, LB11). The cycle-log, critiques, optional gamma cycle log, and optional external-observations files are the canonical surfaces; every role-to-role exchange in a writer cycle commits there. GitHub PR/issue comments are not the triad's coherence record. `WRITER.md` must declare this as a package-level aspect (LANGUAGE-SPEC §7) with the lifecycle sub-skills as the per-role applications.

The package contains:

- **`cn.package.json`** — package manifest with `schema: cn.package.v1`, `name: cnos.writer`, peer engines pin to current cnos version. No commands declared in this cycle.
- **`skills/writer/SKILL.md`** — loader entrypoint. Frontmatter declares the full LANGUAGE-SPEC §2 signature surface: `name`, `description`, `artifact_class: skill`, `kata_surface: embedded`, `governing_question`, `triggers`, `scope: global`, `visibility: public`, `inputs`, `outputs`, `requires`, and a `calls` list (the role skills + lifecycle sub-skills below). This mirrors the shape `cnos.cdd/skills/cdd/SKILL.md` declares (LANGUAGE-SPEC §12.1). Body declares load order: load `WRITER.md`, then the active role skill, then `cnos.core/skills/write/SKILL.md`. Body declares the conflict rule: if `WRITER.md` and `SKILL.md` disagree, `WRITER.md` governs.
- **`skills/writer/WRITER.md`** — normative spec. Sections cover: triadic role definitions (adopted from `CDD.md` §1.4 by reference, with the writer-specific divergences D8 and D9 from this design), lifecycle (intake → draft → critique rounds → external observations evaluation if any → lock), cycle artifacts (cycle-log, critiques, optional gamma cycle log, optional external observations), inherited failure modes (cited from doctrine README §"Inherited failure modes" and IFA γ log §"Inherited Failure Modes (Added for Future Cycles)" — never restated), gates (α-side gates per LB7+LB10; β-side gates per LB1+LB2+LB6+LB8; γ-side or third-surface requirement per LB1).
- **`skills/writer/alpha/SKILL.md`** — α role skill. Mirrors the shape of `cnos.cdd/skills/cdd/alpha/SKILL.md` but adapted for prose: source-load verification (LB7), governing-question freeze (LB10), inherited-commitment citation discipline (LB4). Frontmatter declares the LANGUAGE-SPEC §2 signature surface with `scope: role-local`, `visibility: internal`, and a `calls` set listing the lifecycle sub-skills α reaches for. This mirrors LANGUAGE-SPEC §12.2 (the α role skill is the spec's first-class `calls_dynamic` example). Whether the writer's α uses `calls_dynamic` (analogous to issue.tier3_skills) is a question the authoring cycle decides; this design does not foreclose either choice.
- **`skills/writer/beta/SKILL.md`** — β role skill. Verification of cross-cycle citations against source artifacts (LB8), inherited-failure-mode citation in findings (LB6), structural-grievance-only stopping rule (LB2), no-merge / no-tag closure shape (D8). Optionally absorbs γ-shape work when no γ session is dispatched (D9).
- **`skills/writer/gamma/SKILL.md`** — γ role skill. Cycle-observer voice; produces `cycle-log-gamma.md`-shaped artifact when present. Notes that γ may be absent in dyad cycles per D9.
- **`skills/writer/cycle-log/SKILL.md`** — lifecycle sub-skill for producing cycle-log artifacts. Format: round-by-round failure modes named, what broke each oscillation, what carried forward.
- **`skills/writer/critiques/SKILL.md`** — lifecycle sub-skill for β verdicts. Format: per-verdict finding list, α repair disposition, β cross-cycle verification record (LB8).
- **`skills/writer/external-observations/SKILL.md`** — lifecycle sub-skill for cycles that receive outside review. Format: per-observation evaluation (incorporated / rejected / structural reason), per LB3.

`docs/alpha/doctrine/README.md` gains one link to `cnos.writer/skills/writer/SKILL.md` under §"Skill" with one sentence: "The triadic writer pattern is packaged at `cnos.writer`." No other doctrine surface changes.

The package authoring itself is **not** part of this design's deliverable. Issue #278 Non-Goal #1 forbids writing `cnos.writer/SKILL.md` before design lock. This design closes the design phase. A subsequent issue / cycle authors the package per this design.

## Leverage

- A future writer-shaped cycle starts from a named pattern, not from log reconstruction.
- The IFA-named failure modes (`soft inheritance`, `cold-author evidence refusal`) become package-cited constraints accessible at the role-skill level, not buried in cycle archives.
- The third-surface-required structural finding from CFA travels forward as a hard constraint on writer-cycle dispatch, preventing the dyad-without-meta failure mode from recurring.
- `cnos.cdd` retains its scope (code-shaped substantial cycles); the writer pattern stops accumulating in `cnos.cdd`'s lifecycle skills.

## Negative Leverage

- The package is scaffolding from one essay-shaped material. The first non-doctrine writer cycle may surface divergences this design did not anticipate. Per LB5, those divergences are evidence to inspect, not failures of the package; but they will require a subsequent revision cycle.
- Two peer triadic packages (`cnos.cdd` and `cnos.writer`) increase the surface where shared concepts (α/β/γ, cycle log, role separation) can drift out of sync. Mitigation: writer adopts CDD's role names and references rather than redefining them; the source-of-truth is `CDD.md` for the role contract, `WRITER.md` for writer-specific divergences.
- One additional package to keep src/packaged-copy synced (`skill/SKILL.md` rule 3.10).

## Alternatives Considered

| Option | Pros | Cons | Decision |
|---|---|---|---|
| Package now (chosen) | Memory surface exists when first non-doctrine cycle starts; honors IFA contestability | Scaffolding from one material may need revision after first non-writer use | Selected |
| Wait for one non-doctrine cycle | Generalization is verified before package shape locks | Inheritance via reconstruction risks cold-author evidence refusal at every reuse | Rejected per §5.1 |
| Move `write` to `cnos.writer` | Package is self-contained | Couples general prose discipline to a specific agent pattern | Rejected per §5.2 |
| Encode inheritance discipline inside `cnos.writer` | Single home for IFA failure modes | Imports cross-cycle constraint as material-specific — instance of soft inheritance | Rejected per §5.3 |
| Copy doctrine essays into the package | Self-contained inheritance evidence | Drift risk + violates `write` rule 3.3 + violates Non-Goal #2 | Rejected per §5.4 |
| Make writer a sub-skill of `cnos.cdd` | One triadic package | Misrepresents the writer pattern as a CDD specialization when D3+D5 show distinct shape | Rejected per Challenged Assumption |
| Author the package now (in this cycle) | Faster shipping | Violates Non-Goal #1; design must lock before authoring | Rejected per Non-Goal #1 |

## Process Cost / Automation Boundary

This design adds one package directory and the obligation to author it in a subsequent cycle. The authoring cycle's α produces the SKILL.md / WRITER.md / role skills / lifecycle sub-skills; β reviews against the constraints listed in LB1–LB10; γ verifies the package mirrors `cnos.cdd`'s loader-entrypoint discipline.

Automation boundary:
- Mechanical: package presence in `src/packages/cnos.writer/`, frontmatter validity per `skill/SKILL.md` rule 3.1, source/packaged-copy sync per rule 3.10. A future `writer-verify` command could mechanize cycle-artifact presence (cycle-log + critiques required, external-observations + gamma cycle log optional).
- Judgment-bearing: every LB1–LB10 constraint. Each requires a reviewer to recognize a failure mode and cite the source artifact. No tool can decide whether a phrase is soft inheritance.

## Non-goals

Adopted from issue #278:

- Writing `cnos.writer/SKILL.md` before design lock.
- Migrating the four essays from `docs/alpha/doctrine/`.
- Generalizing the writer to non-cnos repos.
- Replacing `cnos.cdd` (writer and CDD are peer tools).

This design adds:

- Authoring `cnos.core/skills/inherit/` (or any new core skill encoding inheritance discipline) before a non-writer cycle exhibits the same failure modes. This is deferred per §5.3.
- Writing a `writer-verify` command in this cycle. Deferred per §"Process Cost".
- Coining new role names. The writer adopts α/β/γ from `CDD.md` §1.4 verbatim.
- Re-deriving the inherited failure modes inside the writer package. Per LB4, they are cited from doctrine, not restated.

## File Changes

This design produces no source changes other than itself. The package authoring cycle's file changes (subsequent cycle, separate issue) will be:

- **Create:** `src/packages/cnos.writer/cn.package.json`
- **Create:** `src/packages/cnos.writer/skills/writer/SKILL.md`
- **Create:** `src/packages/cnos.writer/skills/writer/WRITER.md`
- **Create:** `src/packages/cnos.writer/skills/writer/alpha/SKILL.md`
- **Create:** `src/packages/cnos.writer/skills/writer/beta/SKILL.md`
- **Create:** `src/packages/cnos.writer/skills/writer/gamma/SKILL.md`
- **Create:** `src/packages/cnos.writer/skills/writer/cycle-log/SKILL.md`
- **Create:** `src/packages/cnos.writer/skills/writer/critiques/SKILL.md`
- **Create:** `src/packages/cnos.writer/skills/writer/external-observations/SKILL.md`
- **Edit:** `docs/alpha/doctrine/README.md` §"Skill" — add one sentence linking to `cnos.writer`.

This design itself creates one file:

- **Create (this cycle):** `docs/alpha/design/WRITER-PACKAGE.md`

## Acceptance Criteria

Per issue #278:

- [x] **AC1** — Design doc addresses four scope questions with three components each (chosen path / alternative considered / structural reasoning):
  - §5.1 Package timing (now vs. after one non-doctrine cycle) — three components present.
  - §5.2 Write skill location (cnos.core vs. cnos.writer) — three components present.
  - §5.3 Inheritance discipline scope (writer-specific vs. cnos.core candidate) — three components present.
  - §5.4 Package content (references vs. copies) — three components present.
- [x] **AC2** — Design doc includes explicit section "Divergences from cnos.cdd structure" (§6) listing each divergence (D1–D10) with reasoning, plus an explicit "No divergence" subsection. Verifiable by reading §6 against `src/packages/cnos.cdd/skills/cdd/` directory listing.
- [x] **AC3** — Design doc names load-bearing constraints from four cycle artifacts (LB1–LB11) with citations to the cycle logs. Citations reference: CFA-cycle-log (LB1, LB2), EFA-cycle-log (LB3), JFA-cycle-log-gamma (LB4), IFA + IFA-cycle-log-gamma + IFA-critiques (LB5, LB6, LB7, LB8), doctrine README (LB9, LB11), write/SKILL.md + essay frontmatter (LB10). Per IFA's argument, citations preserve contestability — a future writer cycle reads the source artifact, not the package's restatement.

## Known Debt

- `cnos.core/skills/inherit/` (or equivalent) remains a candidate, not authored. The first non-writer cycle that exhibits soft inheritance or cold-author evidence refusal becomes the trigger to author it.
- `writer-verify` command not authored in this cycle. The package's mechanical correctness checks are deferred.
- Dispatch protocol for writer cycles (analog of CDD.md's "γ dispatch prompt format") is not specified in this design. The authoring cycle's α defines it inside `WRITER.md`.
- The relationship between writer cycles and cnos issue numbers is not specified. The doctrine cycles ran without issue numbers; future writer cycles may adopt issue dispatch but the package does not require it. This is a deliberate looseness, contestable at first issue-driven writer cycle.

## CDD Trace

| Step | Artifact | Skills loaded | Decision |
|------|----------|---------------|----------|
| 0 Observe | — | — | Read issue #278; surveyed `docs/alpha/doctrine/` (four cycles complete); read `cnos.cdd` package structure; confirmed `cnos.core/skills/write/` and `cnos.core/skills/skill/` are current. Selected signal: writer pattern exists only in cycle artifacts; no package home. |
| 1 Select | — | — | Selected gap: package the triadic writer as `cnos.writer`. |
| 2 Branch | `claude/cnos-issue-278-Xh37v` | cdd | Branch designated by dispatch. (No remote in this environment; branch state held locally.) |
| 3 Bootstrap | — | cdd | Not required — design-only cycle, no version directory bootstrap. Per §6 D5, writer cycles do not carve `docs/{tier}/{bundle}/{X.Y.Z}/`. |
| 4 Gap | this artifact §"Problem" | — | Named incoherence: writer pattern exists only as cycle-archive evidence; future writer cycles must reconstruct the pattern, risking IFA-named cold-author evidence refusal and soft inheritance. |
| 5 Mode | this artifact (header) | `cdd/design`, `core/skills/skill`, `core/skills/write` | MCA — answer is in the system (four cycle artifacts establish the pattern); package is the lighter-weight closure that preserves contestability per IFA. |
| 6 Artifacts | `docs/alpha/design/WRITER-PACKAGE.md` | `cdd/design`, `core/skills/skill`, `core/skills/write` | Design artifact produced. Plan not required (single artifact, no implementation sequencing). Tests/code/docs not required (design-only cycle). |
| 7 Self-coherence | this artifact §"Acceptance Criteria" | cdd | All three ACs (AC1 four-question × three-component; AC2 divergence section + no-divergence assertion; AC3 cited constraints) addressed by inspection of the artifact's own sections §5, §6, §7. Known debt explicit in §"Known Debt". |
| 7a Pre-review | this artifact | cdd | Design-locked draft. Branch `claude/cnos-issue-278-Xh37v` rebased onto current `origin/main` (HEAD `42084b8` at observation time). α-1 HEAD `766680b` carried CI 7/7 success per β round-1 review. α-2 HEAD is the PR head SHA at the moment β round-2 review is requested; CI re-validated on α-2 HEAD before re-request (alpha skill §2.6 transient rows). PR #279 carries CDD trace through step 7a. |


