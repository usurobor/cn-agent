## Post-Release Assessment — 3.57.0

### 1. Coherence Measurement

- **Baseline:** 3.56.2 — α A, β A, γ A
- **This release:** 3.57.0 — α A-, β A-, γ A-
- **Level:** **L6** (cycle cap per CDD §9.1 — review rounds > 2 trigger fired). Diff-level is L6–L7 at the mechanism boundary: filesystem-as-authority for skill existence is now enforced in Go, not merely declared in `DESIGN-CONSTRAINTS.md §1`. `visibility: internal`, declared on 9 CDD sub-skills in 3.55.x, becomes effective at the runtime activation boundary for the first time this release. The cycle took 3 rounds and therefore does not earn the L7 cycle cap, even though the shipped mechanism eliminates a friction class (authority drift between manifest declaration and filesystem truth).
- **Delta:**
  - **α** regressed A → A-. Round-1 push dropped 81% of production trigger lines (inline YAML flow sequences) and shipped a fatal-on-conflict doctor check that the real corpus fails unanimously. A trivial local `cn deps restore && cn doctor` (or `cn kata run --class runtime`) at authoring time would have surfaced both. Rounds 2 and 3 were sharp — the separation of "internal skills shouldn't participate in conflict detection" (Bug A, correctness) from "should trigger overlaps be fatal at all" (Bug B, design) was the right decomposition, and both landed coherently.
  - **β** regressed A → A-. Reviews caught F1 / F2 / F3 / F4 with precise evidence and correctly named F2 and F4 as design-scope deferrals. The gap from A to A- is that β's round-1 §2.2.8 authority-surface audit cited OCaml as the reference authority but did not audit the parser against the actual production SKILL.md corpus — the consumer-side authority that was materially different from the OCaml reference because OCaml's parser was also blind to the corpus (manifest_skill_ids was empty post-ef53b939). The OCaml→Go severity mapping divergence (`cn_doctor.ml:46-49` vs the new Go mapping) also surfaced only after α landed round 3; it deserved an earlier β callout.
  - **γ** regressed A → A-. Ran the cycle cleanly — dispatch prompts, unblocking via the direct "downgrade to warning" direction on round 3, issue #264 filing. The A→A- drop is for the 3-round cycle length (over target 2) and for collapsing the F4 design decision into the PR rather than spawning a separate cycle. The collapse is legitimate per operator-override authority (CDD §1.4) but it does tax the cycle's economics rather than extending them cleanly.
- **Coherence contract closed?** Yes. All six ACs met:
  - **AC1** — `Discover` walks `<hub>/.cn/vendor/packages/<pkg>/skills/`; `TestDiscover_IgnoresManifestSkillsField` proves the manifest `sources.skills` field is never consulted. Filesystem presence is the only authority.
  - **AC2** — `ParseFrontmatter` handles `name`, `description`, `triggers`, `visibility`. Block lists (`triggers:\n  - a`) and inline flow sequences (`triggers: [a, b, c]`) both supported after the round-2 fix. CRLF, BOM, malformed lines all tested.
  - **AC3** — `BuildIndex` excludes `visibility: internal` and defaults absent visibility to public. `TestBuildIndex_ExcludesInternalVisibility` + `TestBuildIndex_AbsentVisibilityDefaultsPublic` prove both paths.
  - **AC4** — `Validate` / `ValidateSkills` return three issue kinds (missing/empty/conflict) and match OCaml's `Cn_activation.issue_kind` taxonomy. The round-3 refinement adds a visibility filter to conflict detection (mirrors OCaml-era `manifest_skill_ids` semantics) and a severity split in the doctor wiring (missing = fatal; empty + conflict = info, per the operator-ratified many-to-many activation model).
  - **AC5** — `PACKAGE-AUTHORING.md` §8 rewritten, "if it's not in the manifest the runtime doesn't know about it" line removed; `PACKAGE-ARTIFACTS.md` `sources.skills` example replaced with filesystem-walk + `visibility: internal` narrative; authoring checklist tightened to reference on-disk skill directories.
  - **AC6** — 33 activation tests (incl. 2 real-corpus integration tests that walk every `src/packages/*/skills/**/SKILL.md` into a temp hub) + 10 doctor tests.

### 2. Encoding Lag

| Issue | Title | Type | Design | Impl | Lag |
|-------|-------|------|--------|------|-----|
| #261 | Wire skill activation to filesystem discovery + frontmatter visibility | feature | converged | **shipped this release** | none |
| #264 | Deterministic package tarball builds | bug | converged | not started | **growing (new this cycle)** |
| #249 | .cdd/ Phase 1 audit trail | feature | converged | partially started | low |
| #256 | CI as triadic surface | feature | converged | not started | growing |
| #255 | CDD 1.0 master | tracking | converged | in progress | — |
| #244 | kata 1.0 master | tracking | converged | in progress | — |
| #245 | cdd-kata 1.0 (prove CDD adds value) | feature | converged | stubs only | growing |
| #243 | cnos.cdd.kata cleanup | feature | converged | partially done | low |
| #242 | .cdd/ directory layout | design | converged | not started | growing |
| #241 | DID: triadic identity | design | converged | not started | growing |
| #240 | CDD triadic protocol | design | converged | partial (Phase 1 in 3.55.0) | low |
| #238 | Smoke: release bootstrap | feature | converged | not started | growing |
| #235 | cn build --check validates entrypoints | feature (P1) | converged | not started | growing |
| #230 | cn deps restore version upgrade skip | bug (P1) | design exists | not started | growing |
| #218 | cnos.transport.git | design | converged | not started | growing |
| #216 | Migrate commands to packages | feature | converged | not started | growing |
| #193 | Orchestrator llm step execution | feature | converged | not started | growing |
| #192 | Runtime kernel rewrite (Go) | feature | converged | Phase 4 complete | low |

**Growing-lag count:** 10 (up from 9 in 3.56.2 — #261 closed, #264 opened). Still well over the 3-issue freeze threshold.

**MCI/MCA balance:** **Freeze MCI** — 10 issues at "growing" lag, continuing the freeze from 3.56.0 / 3.56.1 / 3.56.2. No new design work until the MCA backlog drops below threshold.

**Rationale:** The freeze continues to produce the intended pressure. This cycle closed #261 (a real authority-drift gap) and introduced #264 (a CI-observability gap caught by the #261 review). Net growing-lag change: +1, but the introduced item is a concrete bug with a one-function fix and should close quickly as the next MCA. The pattern "pick the next-smallest growing-lag bug with a clear fix path, ship it, move on" continues to work — this is the fourth consecutive cycle applying it.

### 3. Process Learning

**What went wrong:**

- **3-round cycle, above target ≤2.** Round 1 pushed a parser that silently dropped 81% of the production trigger lines and a doctor check that would fail unanimously on the real corpus. Both were findable at authoring time with a trivial local flow: `cn build && cd <hub> && cn deps lock && cn deps restore && cn doctor`, or the full Tier-2 kata. α ran neither before the round-1 push.
- **Round 2 fix surfaced a new class of failure (F4).** Once the inline-list parser started reading production triggers, the real cnos.core corpus turned up 9 overlapping public-skill trigger keywords plus 12 internal-parent-orchestrator conflicts. β framed the resolution as design-scope-outside-#261 and expected a follow-up issue; α instead investigated and landed a γ-ratified inline resolution across two independent sub-bugs. This was the right outcome, but it inflated the cycle to a 3rd round.
- **OCaml→Go severity divergence surfaced late.** `cn_doctor.ml:46-49` treated `Trigger_conflict` as fatal and Missing/Empty as warn; the new Go mapping inverts that (Missing fatal; Empty + Conflict warn). β's round-3 approval noted this divergence but didn't call it out in rounds 1 or 2 — both α and β had read the OCaml reference in round 1 but neither grepped `cn_doctor.ml` for the severity mapping. Downstream tooling that greps for `✗ skill activation` on conflicts will see `○ skill activation N warning(s)` instead. The change is defensible and operator-ratified, but the gap between "parser parity with OCaml" (β's round-1 check) and "doctor-severity parity with OCaml" (not checked) is a real review-skill gap.
- **I3 red since 3.56.2, not fixed this cycle.** The `Package dist/source sync` check has been red for at least two consecutive releases on a pre-existing main defect (`createTarGz` captures live ModTime + gzip writer time). β correctly scoped it out as #264 rather than block this PR, but the CI gate has been unreliable long enough that every cycle now has to explain "I3 deferred per #264" in every review body. Shipping #264 as the next MCA closes this.

**What went right:**

- **α's round-3 investigation was sharp.** Separating Bug A (internal-skills-in-conflict, mechanism bug) from Bug B (fatal-on-overlap, design question) was exactly the right decomposition. The visibility filter in `ValidateSkills` mirrors OCaml-era `manifest_skill_ids` semantics precisely; the severity split is a γ-directed doctrine call. Both landed coherently with matching tests.
- **The real-corpus integration tests are the MCA.** `TestDiscover_RealCoreSkills_HaveTriggers` copies an unmodified production SKILL.md into a temp hub and asserts non-empty triggers; `TestValidate_RealCorpus_NoEmptyTriggers` walks every `src/packages/*/skills/**/SKILL.md` into a temp hub and asserts no empty-trigger issues surface. This converts the authoring-time integration gap from a discipline question into a mechanical `go test` failure. Future parser/validator changes that break any real SKILL.md now fail unit tests, not CI integration tests. This is an L6 test-shape MCA that shipped in the PR itself.
- **F2 scope discipline held.** The non-deterministic-tar defect in `pkgbuild/build.go:185` was not a #261 regression. β correctly named it deferred-by-design-scope; α filed #264 before merge with a concrete root cause and patch shape. No scope creep into this cycle.
- **§7.0 "no follow-up" rule held through the cycle.** Every finding (A/B/C/D) was resolved on-branch or explicitly deferred by design scope with an issue filed. No "approved with follow-up" anti-pattern.
- **Operator-override authority used cleanly.** γ directed the round-3 "downgrade to warning" path for F4, collapsing a doctrine-level decision into the PR rather than spawning a separate cycle. Per CDD §1.4 this is legitimate; α's round-3 commit message explicitly records "this is the architecturally-significant path the user chose", making the override auditable rather than implicit.

**Skill patches — disposition:**

**No skill patch shipped this cycle.** The failure modes were application gaps and review-skill gaps that do not warrant spec-level patches:

1. F1 authoring-time gap (α didn't run kata locally before round-1 push): application gap. `eng/go` already prefers CI-verification-before-review; a new rule saying "run kata against the installed corpus when the diff touches activation" is too narrow. The MCA is the integration-test pattern that shipped in the PR, not a skill rule.
2. β round-1 gap (§2.2.8 authority-surface audit cited OCaml reference but not production corpus): review-skill gap. The sharpening would be "when reviewing a parser, the real consumer corpus is an authority surface alongside the reference implementation". This is a legitimate candidate MCA — recording here for γ triage rather than shipping this session, because the candidate patch would benefit from reviewing two or three cycles' worth of parser-review patterns before codifying.

**Active skill re-evaluation:**

- `eng/go` §2.17 Parse/Read purity boundary: **worked as written.** `frontmatter.go` has zero `os`/`io` imports; `index.go` is the IO wrapper. β independently grepped.
- `eng/test` §3.10 regression-test-by-name: **worked as written.** F1–F5 / B1–B3 / V1–V3 test names encode ACs. Round-2 `TestParseFrontmatter_InlineTriggers{Supported,Whitespace,Malformed}` names encode the regression class.
- Review skill §2.0 Issue Contract gate: **worked as written.** β produced AC coverage + Named Doc Updates + CDD Artifact Contract tables before reading the diff on all three review rounds.
- Review skill §2.2.8 authority-surface conflict: **partial gap.** β audited (c) Go parser against (a) OCaml parser and concluded "matches OCaml v1 behavior" — accurate but incomplete. The missing audit is (c) against (b) the production SKILL.md corpus, which was materially different from the OCaml reference because OCaml's parser was also blind to the corpus post-ef53b939. Candidate MCA: extend §2.2.8 to name "consumer corpus" as an explicit authority surface alongside "canonical doc vs executable skill". Deferred to γ triage.
- Review skill §2.2.1a input-source enumeration for filters/validators: **worked as written.** β enumerated the frontmatter parser's input sources (inline SKILL.md in tests; real SKILL.md in production) and correctly identified the corpus gap once the F1 fix exposed it.
- `eng/design-principles` (architecture check): **worked as written.** The round-1 review filled the §2.2.14 table in full; "Interfaces remain truthful" correctly went red on F1, driving the RC verdict.

**CDD improvement disposition:** no patch landed this cycle. Justification: F1's failure mode was an application gap with adequate existing spec; β's review-skill gap is real but narrow enough to warrant observation-before-patching. The MCA that closes F1's recurrence class is the real-corpus integration test pattern — already shipped in the PR diff. This closes the §9.1 self-learning loop with a concrete mechanism rather than an explicit-why-not.

### 4. Review Quality

| Metric | Value | Target | Status |
|---|---|---|---|
| PRs this cycle | 1 (#263) | — | — |
| Avg review rounds | 3.0 | ≤2 code | **over target — §9.1 trigger fires** |
| Superseded PRs | 0 | 0 | met |
| Finding breakdown | 0 mechanical / 4 judgment / 4 total | — | — |
| Mechanical ratio | 0% (0/4) | <20% threshold | met (and below 10-total-finding threshold) |
| Action | §9.1 cycle iteration mandatory — see §4b | — | — |

**Round-by-round:**

- **Round 1 (7283741 + 969dcc1 → RC).** β posted REQUEST CHANGES with full §2.0 Issue Contract (6 ACs, 2 named doc updates, 4 CDD artifact rows) + full §2.2.14 Architecture Check + 3 findings:
  - F1 (D, judgment): `ParseFrontmatter` silently drops inline YAML trigger lists — 42 of 52 production SKILL.md files use inline form. Reproduced locally: `cn doctor` on a freshly restored hub reports 19 spurious empty-triggers issues; CI `kata-tier2` (job 72823211280) red.
  - F2 (C, judgment): `Package dist/source sync (I3)` failure is a pre-existing main defect (non-deterministic `createTarGz` + `gzip.NewWriter`). Reproduced on `origin/main`. Out of #261 scope — file separate issue.
  - F3 (B, judgment): `flushList` only materialises `triggers` buffer; rename + switch dispatch for future block-list fields.
- **Round 2 (9cfd3c9 → conditional APPROVE).** α fixed F1 + F3 + filed #264 for F2. β verified all three locally (`parseInlineList` supports inline + whitespace-trimmed + empty-list-nil + malformed-bracket paths; `flushList` dispatches via switch on `pendingListKey`; #264 body carries concrete root cause + proposed fix + byte-identity regression test in ACs). **But** the F1 fix surfaced F4: 9 real trigger-keyword conflicts across public cnos.core skills that now correctly reach the doctor check. β named F4 as design-scope deferred: resolution requires doctrine on skill/trigger precedence outside #261 scope. Conditional approval: merge requires F4 follow-up issue filed.
- **Round 3 (7e76798 → APPROVE).** α did not file the F4 issue; instead investigated further and landed a two-part γ-directed resolution: (A) exclude internal skills from conflict detection (mirrors OCaml-era `manifest_skill_ids` semantics; two new tests `TestValidate_InternalDoesNotConflictWithPublic` + `TestValidate_TwoInternalDoNotConflict`), (B) severity split in doctor wiring (`IssueMissingSkill` stays fatal; `IssueEmptyTriggers` + `IssueTriggerConflict` downgrade to `StatusInfo` — the many-to-many activation model). β verified locally: `cn doctor` on a hub with `cnos.core@3.57.0 + cnos.kata@0.2.0 + cnos.cdd.kata@0.3.0` exits rc=0 with `○ skill activation 9 warning(s)`; `cn kata run --class runtime` passes all 4 runtime katas incl. R3-doctor-broken pre-break baseline. β's round-3 comment surfaced the OCaml `cn_doctor.ml:46-49` severity divergence as a note, not a blocker. Merge: squash, commit `41485ad` on main.

**Finding-count note:** 4 findings total, all judgment, 0 mechanical. Below 10-finding threshold so no process issue filed for mechanical ratio. Review identity was shared-GitHub per review §7.1 — all three β reviews were posted as comments rather than native `--approve`/`--request-changes` because α and β share the `usurobor` account. The review bodies carry all the structure that native review state would have encoded.

### 4a. CDD Self-Coherence

- **CDD α (artifact integrity):** **3/4.** Round-3 PR had all required artifacts: CDD Trace in PR body, 33 activation + 10 doctor tests, no-ocaml-modified constraint honored, `visibility: internal` consumed at runtime. Round-1 artifact integrity was lower (parser didn't read production format) but that is measured on the final diff per the scoring rules. Drop from 4: round-1 self-coherence did not include a real-corpus integration test that would have caught F1.
- **CDD β (surface agreement):** **3/4.** PACKAGE-AUTHORING.md ↔ PACKAGE-ARTIFACTS.md ↔ new Go activation package ↔ CHANGELOG ↔ RELEASE.md ↔ this assessment all agree. Drop from 4: OCaml `cn_doctor.ml` severity mapping divergence surfaced only in round-3 approval, not in earlier rounds — a review-side surface-agreement gap.
- **CDD γ (cycle economics):** **3/4.** 3 review rounds (over target 2). F4 collapsed into PR rather than spawning a separate cycle (operator-override legitimate, but taxes cycle economics). Issue #264 filed cleanly. Tag-push HTTP 403 deferred to γ/operator per β step 8 — environmental, not a γ failure. Drop from 4: cycle length.
- **Weakest axis:** γ (cycle economics — review rounds).
- **Action:** no skill patch. The MCA that eliminates F1's recurrence class (real-corpus integration tests) shipped in the PR. A review-skill sharpening candidate exists (§2.2.8 consumer-corpus as authority surface) but is deferred to γ triage rather than patched this session.

### 4b. Cycle Iteration (mandatory — §9.1 trigger fired)

**Triggers fired:**

- [x] review rounds > 2 (actual: 3)
- [ ] mechanical ratio > 20% — no (0% over 4 total, and below 10-finding threshold)
- [ ] avoidable tooling/environmental failure — no (I3 red is pre-existing main defect deferred per #264; HTTP 403 tag push is documented sandbox constraint per β step 8)
- [x] loaded skill failed to prevent a finding — yes (F1: `cn kata run --class runtime` would have surfaced the inline-list gap at authoring time)

**Friction log:**

Round 1 pushed (a) a parser that silently dropped 81% of production trigger lines because it only implemented the OCaml v1 block-list grammar, and (b) a fatal-on-conflict doctor check that the real corpus fails unanimously. Neither was tested against the actual installed cnos.core SKILL.md corpus before push. α fixed both across rounds 2 and 3 without operator-visible scope slippage — but two pre-merge surprises cost the cycle its ≤2-round target.

**Root cause:**

**Authoring-time integration gap.** The parser was tested against synthetic frontmatter snippets; the doctor check was tested against a clean temp hub with synthetic skills. Neither was exercised against the full `src/packages/` corpus before push. This is an application gap, not a spec gap — α had access to both `cn kata run --class runtime` (which now gates in CI) and a trivial local `cn deps restore && cn doctor` loop.

**Skill impact:**

- `eng/go` — no patch: the skill already prefers CI-verification-before-review; a narrow rule for "activation-adjacent changes" would over-specify.
- `eng/test` — no patch: §3.10 regression-test-by-name worked for every test that was written; the gap is which tests α chose to write, not how those tests were structured.
- Review skill §2.2.8 — **candidate MCA deferred to γ triage.** Extending "authority surfaces" to explicitly name the production consumer corpus (alongside canonical doc / executable skill / runtime prompt / issue AC). Deferred because codifying this benefits from two-to-three parser-review cycles of observation before the rule shape is clear.

**MCA — landed in the PR itself:**

- `TestDiscover_RealCoreSkills_HaveTriggers` — copies an unmodified production SKILL.md from `src/packages/cnos.core/skills/agent/ca-conduct/SKILL.md` into a temp hub and asserts `Discover` returns non-empty triggers.
- `TestValidate_RealCorpus_NoEmptyTriggers` — walks every `src/packages/*/skills/**/SKILL.md` into a temp hub (52 files copied), runs `Validate`, asserts zero `IssueEmptyTriggers`.

Both tests live in `src/go/internal/activation/index_test.go` and run as part of `go test ./...`. Any future parser or validator change that breaks any real SKILL.md now fails unit tests at authoring time, not integration tests at CI time. This converts the authoring-time integration gap from a discipline question ("remember to run the kata") into a mechanical test failure.

"Won't repeat" without a mechanism is not an MCA; these two tests are the mechanism.

**Cycle level: L6** (cycle cap per §9.1 review-rounds trigger).

- **L5 local correctness:** met on the round-3 diff; **not met** at the cycle level because mechanical/design defects (inline-list gap, conflict-severity question) reached review rather than being caught at authoring time.
- **L6 system-safe:** met on the round-3 diff. Cross-surface coherence across docs (PACKAGE-AUTHORING, PACKAGE-ARTIFACTS), runtime (Go activation + doctor), tests (unit + real-corpus integration), and operator truth (`cn doctor` output) all align.
- **L7 system-shaping:** the mechanism eliminates the "authority drift between manifest declaration and filesystem truth" class, but the cycle itself took 3 rounds — cycle level caps at L6 per §9.1.

Justification: The shipped mechanism is strictly L6–L7-shaped (filesystem-as-authority enforced in code; many-to-many activation model made explicit in doctor wiring; real-corpus integration tests as auto-catching regression surface for future cycles). The cycle execution missed L5 cleanly (two authoring-time defects reached review) and therefore caps at L6.

### 5. Production Verification

**Scenario:** `visibility: internal` on a CDD sub-skill must exclude that skill from the public activation index at runtime, and must not cause `cn doctor` to flag spurious trigger conflicts with its public parent orchestrator.

**Before this release:** `visibility: internal` was declared on 9 CDD sub-skills (`cdd/design`, `cdd/review`, `cdd/release`, `cdd/post-release`, `cdd/alpha`, `cdd/beta`, `cdd/gamma`, `cdd/issue`, `cdd/plan`) since 3.55.x but had no runtime effect. OCaml's `Cn_activation.validate` iterated only manifest-declared skill IDs, and `sources.skills` had been removed from manifests in ef53b939 — so `manifest_skill_ids` returned `[]` for every installed package. The `visibility: internal` field was a docstring, not a contract.

**After this release:** The Go `BuildIndex` walks the filesystem, parses each SKILL.md's frontmatter, and excludes any skill with `visibility: internal` from the returned activation index. `ValidateSkills` also excludes internal skills from trigger-conflict claimant sets (an internal sub-skill sharing a trigger keyword with its public parent is not an ambiguity because the runtime cannot route to internal skills). `cn doctor` runs `activation.Validate` and reports unreadable SKILL.md files as `StatusFail` but empty triggers and public-public overlaps as `StatusInfo` warnings.

**How to verify:**

```bash
# 1. Build from the released source
git checkout 3.57.0
go -C src/go build -o /tmp/cn ./cmd/cn

# 2. Create a fresh hub and install the released packages
rm -rf /tmp/verify-hub && mkdir /tmp/verify-hub && cd /tmp/verify-hub
/tmp/cn init h >/dev/null && cd cn-h
/tmp/cn setup >/dev/null
cat > .cn/deps.json <<'JSON'
{
  "schema": "cn.deps.v1",
  "profile": "engineer",
  "packages": [
    {"name": "cnos.core", "version": "3.57.0"},
    {"name": "cnos.cdd",  "version": "3.57.0"},
    {"name": "cnos.kata", "version": "0.2.0"},
    {"name": "cnos.cdd.kata", "version": "0.3.0"}
  ]
}
JSON
/tmp/cn deps lock && /tmp/cn deps restore

# 3. Run doctor
/tmp/cn doctor

# Expected: rc=0, ✓ All checks passed
# Expected: ○ skill activation <N> warning(s) for public-public overlaps
# Expected: no [conflict] between cdd (public) and cdd/review (internal) or similar
```

**Result:** **pass, verified locally on 7e76798 before merge** (the pre-release branch SHA; the released tag `3.57.0` lives at commit `3357409` which adds only version-bump + CHANGELOG + RELEASE.md + dist rebuild on top of the merged PR content).

Observed output:
- `✓ skill activation` checks run as part of doctor.
- `○ skill activation 9 warning(s)` on a hub with cnos.core installed — the 9 expected public-public overlaps (`alignment`, `boundary`, `coherence`, `drift`, `onboard`, `reflect`, `self-check`, `sync`, `verify`). No spurious internal-vs-public conflicts.
- `cn kata run --class runtime` passes all 4 runtime katas (R1 command dispatch, R2 author→build→install round-trip, R3 doctor-broken detection, R4 cn-status).
- `TestBuildIndex_ExcludesInternalVisibility` (unit test) confirms the index-builder side.

The targeted coherence delta — `visibility: internal` enforced at the runtime activation boundary rather than being a docstring — is proven by both the unit-test path (`BuildIndex` excludes internal skills from the index) and the integration path (`TestValidate_RealCorpus_NoEmptyTriggers` walks the full `src/packages/` tree). **Pass.**

### 6. CDD Closeout

| Step | Artifact | Skills loaded | Decision |
|------|----------|--------------|----------|
| 11 Observe | local `cn doctor` + `cn kata run --class runtime` against restored hub on 7e76798 / 3357409 | post-release | runtime matches design — internal skills filtered from conflicts, overlaps warn-not-fatal, `cn doctor` exits rc=0 |
| 12 Assess | this file (`docs/gamma/cdd/3.57.0/POST-RELEASE-ASSESSMENT.md`) | post-release, writing | assessment completed |
| 13 Close | β close-out at `.cdd/releases/3.57.0/beta/261.md` + release commit `3357409` + merge commit `41485ad` | post-release, cdd | cycle closed; deferred outputs (#264 tarball determinism, tag push to γ, branch cleanup to γ) recorded concretely in §7 |

### 6a. Invariants Check

cnos maintains `docs/alpha/DESIGN-CONSTRAINTS.md` as the canonical invariants surface. Constraints touched this cycle:

| Constraint | Touched? | Status |
|---|---|---|
| §1 Source of truth (each fact in one place) | yes | **tightened** — filesystem (`<pkg>/skills/<id>/SKILL.md`) is now the single authority for skill existence; manifest `sources.skills` field is never consulted by Go runtime |
| §2.1 One package substrate (`cn.package.v1`) | no | preserved |
| §2.2 Source / artifact / installed clarity | yes | preserved — `.cn/vendor/packages/` is the only installed-state surface read; source is `src/packages/`; artifact is `dist/packages/` |
| §3.1 Git-style subcommands | no | preserved |
| §3.2 Dispatch boundary (pure in `internal/`, IO wrappers in `cli/`) | yes | preserved — `activation/frontmatter.go` is pure; `activation/index.go` is the IO wrapper; `cli/cmd_doctor.go` remains a thin dispatcher |
| §4.1 Surface separation (skills/commands/orchestrators/providers distinct) | yes | preserved — the new activation check runs only over `skills/`; commands and orchestrators are unaffected |
| §4.2 Registry normalization | no | preserved |
| §5.0 OCaml deprecated | yes | preserved — zero files in `src/ocaml/` or `test/cmd/*.ml` modified; the OCaml path is now correctly superseded rather than extended |
| §6.1 Reason to change | yes | preserved — `activation` package owns discovery + index + validate; `doctor` consumes via a thin wrapper; no convenience-bucket smear |
| §6.2 Policy above detail | yes | preserved — kernel decides public/internal via `BuildIndex` filter; package authors declare only `visibility:` |
| §6.3 Degraded-path visibility | yes | **preserved-tightened** — trigger overlaps now visible as `○` warnings with specific messages ("trigger X claimed by: a, b") rather than the earlier invisible (OCaml path blind) or misleading ("skill X has no triggers" when triggers actually existed) framings |

No constraint was revised without explicit naming. §1 and §6.3 are tightened; all others preserved.

### 7. Next Move

**Next MCA:** **#264** — Deterministic package tarball builds: strip file ModTime and gzip header time.

**Owner:** α (next dispatch)
**Branch:** pending
**First AC:** `cn build` twice on the same clean source tree produces byte-identical `dist/packages/*.tar.gz` under `git diff`.
**MCI frozen until shipped?** Yes — lag table shows 10 growing-lag items, freeze continues.
**Rationale:** The `Package dist/source sync (I3)` CI check has been red for at least two consecutive releases on a pre-existing main defect (`createTarGz` captures live ModTime + gzip writer time). Fixing it is (a) concrete — a one-function patch in `src/go/internal/pkgbuild/build.go` with the exact shape already drafted in #264's body, (b) CI-gating — I3 goes green → cross-cycle trust in `dist/packages/` reproducibility, (c) a precondition for any external mirror or audit trail that needs "this tarball was produced from that source". Closing it next closes the persistent CI-trust gap, and is adjacent to #261's work (both live in the Go kernel's build/restore chain).

The natural candidate after #264 remains **#230** (`cn deps restore` version-upgrade skip) — the `restoreOne` skip-if-installed shortcut means a lockfile bump from `cnos.core@1.0.0` → `cnos.core@2.0.0` is silently skipped because the directory already exists. That's the next link in the manifest → lockfile → installed authority chain once the dist chain is reproducible.

**Closure evidence (CDD §10):**

- **Immediate outputs executed:**
  - ✅ Release artifacts committed to main at `3357409` — VERSION=3.57.0, cn.json, 3 package manifests, CHANGELOG ledger row, RELEASE.md, dist/packages (5 rebuilt tarballs + checksums + index).
  - ✅ Main pushed to origin (`1d156ae..3357409`).
  - ❌ Tag push `git push origin 3.57.0` returned HTTP 403 (sandbox environment constraint; reproduced via 2 retries and independent `git ls-remote` verification showing the tag is not on the remote). **Deferred to γ/operator** per CDD §1.4 β step 8. Local tag `3.57.0` exists at `3357409`; γ must run `git push origin 3.57.0` from an environment that can reach the remote.
  - ❌ Branch cleanup `git push origin --delete claude/alpha-skill-go-implementation-LJVHx` also returned HTTP 403 (same sandbox constraint). **Deferred to γ/operator** per release skill §2.6a.
  - ✅ Issue #264 filed by α in round 2 with concrete root cause + proposed fix + byte-identity regression test in ACs — now the committed next MCA.
  - ✅ Issue #261 auto-closed on squash-merge (the commit body contained `Closes #261`).
  - ✅ Post-release assessment (this file) committed to `docs/gamma/cdd/3.57.0/POST-RELEASE-ASSESSMENT.md`.
  - ✅ β close-out at `.cdd/releases/3.57.0/beta/261.md` (next write).
- **Deferred outputs:**
  - **#264** — deterministic tarball build (α, next cycle). Explicit next MCA.
  - **Corpus trigger-conflict disambiguation** — whether to disambiguate the 9 overlapping public-skill triggers in cnos.core, introduce activation precedence, or formalize the many-to-many activation model as a dedicated doctrine doc. γ triage candidate, not β-scope. Not a blocker for any subsequent cycle — the current runtime semantic (overlaps are informational warnings) is coherent as shipped.
  - **Review skill §2.2.8 consumer-corpus authority-surface sharpening** — candidate MCA for a future review-skill patch, deferred to γ triage after observing two to three more parser-review cycles.
  - **OCaml `cn_doctor.ml` severity divergence note in a doctrine surface** — currently the many-to-many activation model lives only in `src/go/internal/doctor/doctor.go` line 107–109 comments. Capturing it in a doctrine doc (candidate: `docs/alpha/package-system/PACKAGE-ARTIFACTS.md` activation section, or a new `docs/alpha/ACTIVATION-MODEL.md`) is a future-cycle docs task.

**Immediate fixes (executed in this session):**

- Release artifacts committed + main pushed (this session).
- Post-release assessment committed (this file).
- β close-out committed (`.cdd/releases/3.57.0/beta/261.md`, next write).

No CDD skill patches shipped this cycle — rationale in §3 disposition and §4b skill impact.

### 8. Hub Memory

The cnos-kernel repository does not maintain a hub-memory surface distinct from `.cdd/releases/`. Prior releases treated the β close-out under `.cdd/releases/<version>/beta/<issue>.md` as the adhoc-thread analog (cross-session continuity for the Go-port / skill-activation / CDD-evolution threads), and the post-release assessment under `docs/gamma/cdd/<version>/POST-RELEASE-ASSESSMENT.md` as the daily-reflection analog.

- **Daily reflection analog:** this assessment, at `docs/gamma/cdd/3.57.0/POST-RELEASE-ASSESSMENT.md`. Committed in the same commit as §7 closure.
- **Adhoc thread update:** β close-out at `.cdd/releases/3.57.0/beta/261.md` (next write). This release advances two ongoing threads: (a) **Go kernel rewrite** (#192) — Phase 4 is effectively complete, and the activation path is now the second-largest OCaml→Go port after `cn deps restore` in 3.43.0; (b) **Authority clarity** — filesystem-as-authority for skill existence closes a gap that was open since ef53b939 removed `sources.skills` from manifests.

**No external hub-memory repository exists for this project**; noting this explicitly so future assessments don't silently skip a surface that is supposed to be there. If an adhoc-thread repo is introduced (#242 "cdd directory layout" or follow-up), this section's obligation extends to cover it.

---

Signed: β (`beta@cdd.cnos`) · 2026-04-24 · release commit `3357409` · merge commit `41485ad` · assessment `docs/gamma/cdd/3.57.0/POST-RELEASE-ASSESSMENT.md`
