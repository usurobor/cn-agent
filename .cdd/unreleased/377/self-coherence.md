# Self-Coherence — cycle #377

**Issue:** #377 — design(cdd): codify cross-repo coordination protocol (intake / outbound / bilateral)
**Branch:** `cycle/377`
**Mode:** design-and-build
**Active Skills:** `cdd/skill`, `cdd/design`, `cdd/gamma` (read-only), `cdd/post-release` (read-only), `cdd/issue` (read-only)
**Role collapse:** γ=δ=α=β collapsed on δ for a single design-and-build cycle under §5.2 wave-mode with γ=δ collapse permitted. α≠β within a session is structurally compromised but accepted for the design phase converging on the build per wave manifest precedent. Acknowledged explicitly here and in receipt artifacts.

## Gap

Cross-repo CDD coordination has run successfully eight times (counted across the empirical anchors below) but the protocol is partial, scattered across three skill files, and contains directional collisions that γ re-discovers each cycle.

Status before this cycle:
- protocol fragments live in `cdd/gamma/SKILL.md §2.1` + `§2.7`, `cdd/post-release/SKILL.md §5.6b`, `cdd/issue/SKILL.md`, and `.cdd/iterations/cross-repo/README.md`.
- the fragments contradict each other on path conventions (intake says `.cdd/iterations/proposals/*/` or `.cdd/proposals/{target}/*/`; actual practice ships at `.cdd/iterations/cross-repo/{target}/{slug}/`).
- five concrete gaps per the issue: path collision, two parallel state machines, ambiguous `landed`, no master/sub mapping, undocumented bundle file set / LINEAGE schema / feedback-patch format / archival rule.
- new observations from cnos session 2026-05-19: `drafted` event used by cn-sigma (not in 5-event vocabulary); bundle file set varies across anchors; `## Source Proposal` block convention with placeholders is cleaner than fresh-block insertion; hat-collapse attribution has no protocol form; inbound mirror at `cnos:.cdd/iterations/cross-repo/{source}/{slug}/` is dual-purpose; operator-pending bundles for non-existent repos; proposals shaped as issue comments.

This cycle codifies one canonical surface at `src/packages/cnos.cdd/skills/cdd/cross-repo/SKILL.md` that names the directional cases, canonical path, STATUS state machine, bundle file set, LINEAGE schema, feedback-patch format, archival rule, hat-collapse attribution, and known edge cases.

## Constraints

- **One canonical surface per protocol concern** (`CDD.md` doctrine-vs-skill rule). Cross-repo coordination is one protocol concern; one skill owns it.
- **Empirical anchor before doctrine** (`cdd/issue/SKILL.md §"Empirical anchor"` precedent). The state machine and schemas must validate retroactively against the 8 existing precedents without contradiction.
- **No surface fusion** (`COHERENCE-CELL.md` AC2 invariant). The cross-repo skill names protocol shape (states, paths, schemas); it does not author runtime mechanics (how γ scans, polls, or writes — those stay in gamma's role surface).
- **Master/sub mapping must be explicit.** The protocol answers "what fires `landed` when a proposal accepted as a master fans out into N independently-landing subs?" — not leaves it to γ judgment.
- **Doctrine lives in cnos, not in bundles.** Any rule that applies to all cross-repo bundles lives in `cdd/cross-repo/SKILL.md`. Bundle READMEs describe only their bundle's specific context.
- **Wave manifest invariants apply.** Changes only under named directories. No CI / runtime / release surface change. Empirical anchors must be cited.
- **No expansion of vocabulary or phases.** Per issue non-goals: do not add a new STATUS event or bundle-state phase; reconcile within the existing five events + three phases.

## Challenged Assumption

The pre-cycle assumption was that cross-repo coordination doctrine could live as fragments — one piece in each role skill, with `.cdd/iterations/cross-repo/README.md` carrying bundle convention. This cycle challenges that assumption: cross-repo coordination is **one protocol concern** that demands **one canonical surface**, and the fragments served only as scaffolding pending the canonical skill.

The challenge is structural, not stylistic — fragment-based doctrine produces directional ambiguity (the empirically-observed failure mode), where γ invents bundle shape per cycle because no single surface answers "which case am I in?".

## Impact Graph

**Downstream consumers:**
- `cdd/gamma/SKILL.md §2.1` (intake) — was carrying inline directional doctrine; now references the new skill.
- `cdd/gamma/SKILL.md §2.7` (close-out) — was carrying inline `landed`-event ambiguity; now references the new skill including master/sub rule.
- `cdd/post-release/SKILL.md §5.6b` (cross-repo trace) — was carrying inline outbound-only doctrine; now references the new skill for both outbound and bilateral cases.
- `.cdd/iterations/cross-repo/README.md` — was carrying bundle convention with three phases unrelated to STATUS events; now aligns phases with the canonical state machine and cross-references the new skill.

**Upstream producers:**
- `cdd/skill/SKILL.md` (meta-skill) — provides frontmatter shape, kata embedding, parent declaration; the new skill conforms to its rules 3.1–3.10.
- `cdd/design/SKILL.md` — provides the design-and-build mode; this cycle's design half produced `design-notes.md` with retroactive validation against 8 anchors before the build half.
- 8 empirical anchors under `.cdd/iterations/cross-repo/` — retroactive validators; their existing shape constrains the protocol.

**Copies and embeddings:**
- The new skill is the canonical home for cross-repo doctrine. No other surface should carry the doctrine after this cycle. Verified by `rg` (see AC6 oracle).
- The `## Source Proposal` block in `cdd/issue/SKILL.md` is **unchanged** — that fragment is the target-side issue-body integration and belongs in `issue/`.

**Authority relationships:**
- Per `CDD.md` §1.4: role skills add execution detail without changing canonical step order. The new skill is a **protocol-doctrine sub-skill** under `parent: cdd`; it does not redefine roles or lifecycle, only the cross-repo coordination domain.
- If `cdd/cross-repo/SKILL.md` and `gamma/SKILL.md §2.1/§2.7` disagree after this cycle, `cdd/cross-repo/SKILL.md` governs (cross-referenced in the gamma sections).

## Proposal

A new skill `src/packages/cnos.cdd/skills/cdd/cross-repo/SKILL.md` codifies:

1. **Directional cases** — four cases (a inbound 1:1/master+sub, b outbound, c bilateral, d operator-pending with three sub-shapes).
2. **Canonical path** — `.cdd/iterations/cross-repo/{counterpart-repo}/{slug}/`, direction-agnostic.
3. **STATUS state machine** — 5 events (`submitted | accepted | modified | rejected | landed`), `drafted` reconciled as `submitted` synonym, transition graph with `accepted → modified` permitted (post-filing refinement), terminal `landed` and `rejected`, master/sub `landed` rule (per-sub + master-close).
4. **Bundle-state phases** — `open | converging | closed` mapped to STATUS events.
5. **Bundle file set** — per directional case; case (a) carries STATUS + LINEAGE + optional FEEDBACK.patch + optional ISSUE.md; case (b) carries LINEAGE; case (c) carries LINEAGE + cdd-iteration.md + N feedback patches; case (d) carries LINEAGE + substantive content.
6. **LINEAGE.md schema** — per case; named sections required per case in §2.6.
7. **Feedback-patch format** — unified diff with header (From / Date / Subject), prose context, apply command (`git apply`), then the diff. Canonical filename `FEEDBACK.patch` for STATUS-only patches; `<descriptor>.patch` for substantive patches.
8. **Bundle archival rule** — asymmetric: target-side mirrors preserved indefinitely (audit artifacts); source-side bundles archivable post-receipt.
9. **Hat-collapse attribution** — recorded in `LINEAGE.md §"<role> actor"` and (for case c) in `cdd-iteration.md §"Hat-collapse acknowledgment"`. Boundary-preserving — names actor-level fact without merging role-level identities or authorities.
10. **Known protocol edge cases** — `drafted` synonym, repository rename mid-coordination, post-filing refinement, asymmetric source posture, dual-purpose case-c bundle, proposal-as-issue-comment, target-repo-doesn't-exist.
11. **Cross-references** — to gamma, post-release, issue, epsilon skills.

Three existing skills are updated to reference the new skill instead of carrying inline doctrine:
- `cdd/gamma/SKILL.md §2.1` (intake) — reduced from 12-line directional doctrine to 6-line one-paragraph summary + cross-reference.
- `cdd/gamma/SKILL.md §2.7` (close-out) — reduced from inline `landed`-ambiguity prose to 3-line summary + cross-reference (now naming master/sub rule).
- `cdd/post-release/SKILL.md §5.6b` — reduced from 5-line outbound-only inline doctrine to 4-line summary + cross-reference covering outbound + bilateral.

`.cdd/iterations/cross-repo/README.md` is updated to:
- name the path as direction-agnostic ({counterpart-repo} instead of {target}).
- align bundle-state phases with the canonical state machine.
- cross-reference the new skill.
- name the four directional cases.
- list all existing bundles with their case classification.

The `## Source Proposal` block in `cdd/issue/SKILL.md` is **unchanged**.

## Leverage

- **γ re-discovery cost is eliminated.** A future γ acting on any cross-repo event finds one canonical surface naming the directional case, bundle shape, STATUS rules, archival rule. No more "judgment outside protocol coverage" per cycle.
- **Bundle shape stability across cycles.** New bundles match §2.5 file-set; new LINEAGE.md match §2.6 schema. Drift over time is bounded.
- **Tooling becomes feasible.** With the protocol surface stable, a future `cn cross-repo` CLI (status check, lineage validator, archival helper) can be authored against a fixed contract. This cycle does not ship the tool; it makes the tool authorable.
- **Master/sub patterns supported.** Future research-protocol or supercycle proposals that fan out into waves have a documented `landed` rule.
- **Cross-repo state cannot go stale.** With `landed` emitter and master/sub rule defined, source-side STATUS ledgers cannot sit at `accepted` indefinitely after target work lands.

## Negative Leverage

- **One more skill to maintain.** The new skill is +600 lines under `src/packages/cnos.cdd/skills/cdd/`. Future cycles must keep it in sync with cross-repo practice (low cost given the protocol's now-explicit shape).
- **Strict shape compliance now required.** Future bundles that diverge from §2.5/§2.6 are findings rather than judgment calls — small adjudication cost shifted from γ to β.
- **No multi-hop coverage.** The first time a future cycle hits multi-hop (A → B → C) or cross-org, this skill extends. Carried as named debt in §Scope out-of-scope; not a regression vs. the pre-cycle state which also did not cover multi-hop.

## Alternatives Considered

| Option | Pros | Cons | Decision |
|---|---|---|---|
| Keep doctrine in `gamma/SKILL.md` as inline fragments | Minimal change; no new file | Continues directional ambiguity; fragments contradict each other on path; doctrine archives with bundle README | Rejected — issue identifies the failure mode |
| Add `cross-repo` as a section inside `gamma/SKILL.md` | One canonical home; no new file | Bloats `gamma/`; mixes role doctrine with protocol doctrine; violates "one skill, one domain" (skill §2.1) | Rejected |
| New skill `cdd/cross-repo/SKILL.md` (this) | Domain-isolated; one canonical home; aligns with `cdd/skill/` rule 2.1 + `issue/` precedent | One more file to maintain | **Selected** |
| Author the CLI surface (`cn cross-repo`) in this cycle | Tooling + doctrine in one cycle | Doctrine not yet stable; would re-couple them; issue explicitly defers CLI | Rejected per issue non-goal |
| Expand STATUS vocabulary to 6 events (`drafted` as distinct) | Recognizes observed practice directly | Issue non-goal explicitly forbids new events | Rejected — reconcile within 5 events |

## Process Cost / Automation Boundary

- **New artifact:** `cdd/cross-repo/SKILL.md`. Owned by γ (the role that invokes it most often); maintained per `cdd/skill/SKILL.md` meta-skill rules.
- **Light update cost on existing skills:** `gamma/SKILL.md §2.1+§2.7`, `post-release/SKILL.md §5.6b`, `.cdd/iterations/cross-repo/README.md`. One-line summary + cross-reference each.
- **No automation in this cycle.** Future automation candidates (per leverage above): `cn cross-repo status` (scan all bundles, list phases), `cn cross-repo validate` (check LINEAGE schema per case, STATUS event sequence), `cn cross-repo archive` (apply §2.8 rule to source-side bundles). These are deferred until after one more cross-repo cycle validates protocol stability.
- **Automation boundary:** STATUS sequence-validation is mechanical (could be linted); case-identification is judgment (γ must read the bundle's directional shape and decide); archival is judgment (γ decides when the predicate is satisfied).

## Non-goals

- No multi-hop coordination (A → B → C). Bilateral only.
- No cross-org coordination. Same-org assumed.
- No CLI tooling (`cn cross-repo` deferred to a separate issue once the protocol is stable across ≥1 more cycle).
- No change to `cdd/issue/SKILL.md` `## Source Proposal` block — that fragment is clean and stays.
- No new STATUS event (the `drafted` reconciliation is a synonym, not a new event).
- No new bundle-state phase.
- No retroactive cleanup of existing bundles (in-flight coordination stays as-is; this cycle codifies the protocol so the next one runs cleaner).
- No deletion of historical bundles (e.g., `gait-support-paths/bootstrap-cdr/` stays as audit artifact).

## File Changes

**Create:**
- `src/packages/cnos.cdd/skills/cdd/cross-repo/SKILL.md` — new canonical skill (~620 lines).
- `.cdd/unreleased/377/gamma-scaffold.md` — γ scaffold (committed).
- `.cdd/unreleased/377/design-notes.md` — α design notes with retroactive validation against 8 anchors.
- `.cdd/unreleased/377/self-coherence.md` — this file.
- `.cdd/unreleased/377/beta-review.md` — β-collapsed-on-α self-review (pending).
- `.cdd/unreleased/377/alpha-closeout.md` — α close-out (pending).
- `.cdd/unreleased/377/beta-closeout.md` — β close-out (pending).
- `.cdd/unreleased/377/gamma-closeout.md` — γ close-out (pending).

**Edit:**
- `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` — §2.1 + §2.7 reduced to one-line summary + cross-reference.
- `src/packages/cnos.cdd/skills/cdd/post-release/SKILL.md` — §5.6b `**Cross-repo trace.**` paragraph reduced to summary + cross-reference.
- `.cdd/iterations/cross-repo/README.md` — direction-agnostic path; phase-to-STATUS mapping; new skill cross-reference; case classification per existing bundle.

**Unchanged (per non-goals):**
- `src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` — `## Source Proposal` block stays.

## Acceptance Criteria

- [x] **AC1.** `src/packages/cnos.cdd/skills/cdd/cross-repo/SKILL.md` exists with standard CDD frontmatter (`name`, `description`, `artifact_class`, `governing_question`, `visibility: internal`, `parent: cdd`, `triggers`, `scope`); ≥2 reference matches in `gamma/` + `post-release/`. Oracle: read frontmatter; `rg "cdd/cross-repo" src/packages/cnos.cdd/skills/cdd/{gamma,post-release}/SKILL.md`. Pass: file exists, frontmatter complete, gamma + post-release each cite the new skill.
- [x] **AC2.** STATUS state machine: 5-event vocabulary + `drafted` reconciliation + transition graph + emitters per event + master/sub rule for `landed`. Oracle: read §2.3 of new skill; verify table/graph + master/sub rule named. Pass: all components present in §2.3.1–§2.3.5.
- [x] **AC3.** Bundle file set + LINEAGE.md schema per directional case; both existing precedents validate retroactively without contradiction. Oracle: read §2.5 + §2.6 of new skill; design-notes.md §3.6 + §4.5 validation tables show zero contradictions across all 8 anchors (not just the two precedents named in the issue). Pass: §2.5/§2.6 carry per-case tables; design-notes retroactive validation shows zero contradictions.
- [x] **AC4.** Path-collision resolved: one canonical path per directional case (here: one canonical path direction-agnostic); intake-scan in `gamma/SKILL.md §2.1` updated to match. Oracle: `rg "\.cdd/iterations/proposals|\.cdd/proposals/\{target\}" src/packages/cnos.cdd/skills/cdd/` returns zero or deprecation-only. Pass: legacy paths removed from gamma; canonical path named in exactly one place (the new skill) and referenced from gamma.
- [x] **AC5.** Feedback-patch format + bundle archival rule codified in new skill; doctrine removed from bundle READMEs (none had it; `.cdd/iterations/cross-repo/README.md` carries bundle-shape + cross-reference, not protocol doctrine). Oracle: new skill §2.7 + §2.8 carry the doctrine; `rg "archive|delete.*source-side" .cdd/iterations/cross-repo/*/README.md` returns no doctrine (no bundle has a README; the top-level README does not carry archival doctrine). Pass: format + rule in skill, not in bundles.
- [x] **AC6.** Existing fragments in `gamma/SKILL.md §2.1 + §2.7` and `post-release/SKILL.md §5.6b` reference the new canonical surface; `issue/SKILL.md` `## Source Proposal` block stays as-is. Oracle: `rg "cdd/cross-repo" src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md src/packages/cnos.cdd/skills/cdd/post-release/SKILL.md` returns ≥3; `git diff origin/main -- src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` is empty. Pass: ≥3 cross-references; no diff on `issue/SKILL.md`.

All ACs are mechanically verifiable and map to specific files/checks per `cdd/design/SKILL.md` rule 6.2.

## Known Debt

- **Multi-hop and cross-org coordination** — protocol does not cover A → B → C or different GitHub orgs. When the next cycle hits multi-hop, this skill extends. Carried as named debt in §Scope out-of-scope.
- **CLI surface deferred** — `cn cross-repo status / validate / archive` candidates are deferred until protocol stability is validated across ≥1 more cross-repo cycle. Not filed as an issue this cycle; ε iteration may file as project MCI.
- **Bundle READMEs from pre-protocol cycles** — none observed to carry protocol doctrine (the issue's concern was anticipatory based on a pattern that did not materialize in the empirical anchors examined; the `gait-support-paths/bootstrap-cdr` README that the issue body cited as carrying doctrine was archived in a previous cleanup — the current state has only LINEAGE.md). No action required.
- **Per-cycle protocol-observation forwarding** — case (c) bundles may carry `## Protocol observations forwarded` sections (per §4.3); these need a cadence for codification into this skill. Suggested cadence: at each ε iteration, surface unfolded observations to the next γ for evaluation.

## CDD Trace

| Step | Artifact | Skills loaded | Decision |
|------|----------|--------------|----------|
| 0 Observe | wave manifest, issue #377, 8 empirical anchors | gamma (read-only) | wave-dispatched cycle; mode design-and-build |
| 1 Select | wave manifest | gamma (read-only) | decisive clause: wave manifest §"Issues" row #377 |
| 4 Gap | this file §Gap, design-notes.md | gamma (read-only), design | named: cross-repo protocol is partial, fragmentary, directionally ambiguous; failure mode "directional ambiguity" |
| 5 Mode | this file header | design | design-and-build; active skills: cdd/skill, cdd/design, cdd/gamma, cdd/post-release, cdd/issue |
| 6 Artifacts | design-notes.md (design), self-coherence.md (this file), cross-repo/SKILL.md (build), gamma+post-release+README updates (extraction) | skill (meta), design | design phase ran first with retroactive validation against 8 anchors; build phase wrote the skill and updated references |
| 7 Review-ready | this file §Acceptance Criteria + beta-review.md (pending) | β (collapsed-on-α self-review) | review-ready signal: this file complete; β-collapsed review runs next |

## Review-readiness signal

**Review-ready: yes.**

α phase outputs complete:
- design-notes.md (design half)
- cross-repo/SKILL.md (build half — canonical surface)
- gamma/SKILL.md §2.1 + §2.7 edits (extraction)
- post-release/SKILL.md §5.6b edit (extraction)
- .cdd/iterations/cross-repo/README.md update (alignment)
- self-coherence.md (this file)

All 8 empirical anchors validate retroactively against the new protocol with zero contradictions (per design-notes §3.6, §4.5).

All 6 ACs have been authored against and self-checked.

Next: β-collapsed-on-α self-review (beta-review.md) verifying AC1–AC6 with particular attention to AC3's retroactive validation discipline.

## Fix rounds

### R1 → RC (β R1, 2 binding findings)

β R1 issued REQUEST CHANGES with two binding findings:

- **B-1.** STATUS vocabulary in `cdd/cross-repo/SKILL.md` enumerated 5 events; `CDD.md` §"Cross-repo proposal lifecycle" lines 234–243 ships an 8-event vocabulary. Per `cdd/SKILL.md` §"Conflict rule", CDD.md governs on artifact contract. The 5-event count in the issue body was an under-read of CDD.md.
- **B-2.** Legacy paths `.cdd/iterations/proposals/{slug}/` and `.cdd/proposals/{target}/{slug}/` were named as first-class in `CDD.md` lines 215–225. The new skill ❌-marked them; same-cycle contradiction with the canonical source.

### R2 fix

α R2 applied the following corrections:

1. **`CDD.md` §"Cross-repo proposal lifecycle"** updated to name the canonical source-side path as `{source-repo}:.cdd/iterations/cross-repo/cnos/{slug}/` (matching the cnos-side mirror shape) and cross-reference `cdd/cross-repo/SKILL.md` as the canonical surface for the protocol mechanics (transitions, emitters, master/sub, archival). Legacy first-class layouts removed.

2. **`cdd/cross-repo/SKILL.md §2.3.1`** vocabulary expanded from 5 to 8 events: `drafted | submitted | accepted | modified | landed | rejected | withdrawn | revised/corrected`, verbatim with CDD.md meanings.

3. **`cdd/cross-repo/SKILL.md §2.3.2`** transition graph rebuilt to include `drafted`, `withdrawn`, audit events, and the `accepted → modified` post-filing-refinement transition observed in cn-sigma.

4. **`cdd/cross-repo/SKILL.md §2.3.3`** added: `drafted → accepted` direct-acceptance path (replaces the R1 `drafted`-as-`submitted`-synonym), permitted when source delegates filing-authority via LINEAGE record.

5. **`cdd/cross-repo/SKILL.md §2.3.4`** emitters table extended with `drafted` (source role), `withdrawn` (source γ), `revised/corrected` (event-originator).

6. **`cdd/cross-repo/SKILL.md §2.4`** bundle-state phase mapping updated: `open` includes both `drafted` and `submitted`; `closed` includes `withdrawn`. Audit events do not change phase.

7. **`cdd/cross-repo/SKILL.md §2.10`** known edge case "drafted as synonym" replaced with "drafted → accepted direct acceptance" (cn-sigma anchor).

8. **`cdd/cross-repo/SKILL.md §3.3`** rule expanded: source STATUS must not remain at `submitted` OR `drafted` post-decision.

9. **`cdd/cross-repo/SKILL.md`** kata reasoning step 2 + step 4 STATUS handling updated to reflect direct-acceptance.

10. **`.cdd/iterations/cross-repo/README.md`** phase-mapping table updated to include `drafted` in `open` and `withdrawn` in `closed`; audit-events note added.

11. **`.cdd/unreleased/377/design-notes.md`** §2.1, §2.2, §2.3, §2.4, §2.7 updated to reflect R2 vocabulary; §2.3a (R1 reconciliation discussion) removed since R2 transition graph permits `accepted → modified` natively.

The wave manifest invariant 1 path "src/packages/cnos.cdd/skills/cdd/" includes CDD.md, so the CDD.md update is within wave scope. Cross-cycle impact: none — #375 touches gamma §2.5, #378 touches review/alpha/operator skills; this cycle's CDD.md update touches only the cross-repo lifecycle section (lines 213–245).

### Re-running AC oracles after R2 fix

| AC | Result after R2 |
|---|---|
| AC1 | PASS (unchanged; frontmatter + 3 references) |
| AC2 | **PASS** — 8-event vocabulary matches CDD.md canonical source; transition graph + emitters + master/sub rule all present |
| AC3 | PASS — retroactive validation now zero contradictions including cn-sigma anchor |
| AC4 | **PASS** — CDD.md updated to canonical path; legacy paths removed (now appear only as ❌ deprecation examples in the new skill) |
| AC5 | PASS (unchanged) |
| AC6 | PASS — references now ≥3 (gamma:2 + post-release:1); `issue/SKILL.md` still unchanged; CDD.md cross-references the new skill |

α R2 review-ready. Pending β R2.
