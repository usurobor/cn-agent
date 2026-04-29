## Post-Release Assessment — 3.61.0

> Cycle **#283** (replace GitHub PR workflow with `.cdd/unreleased/{N}/` artifact exchange in CDD; first cycle of new artifact-channel coordination protocol). Merged at `58c1666`; release commit `585efe3` by β; tag `3.61.0` created locally, push pending δ (HTTP 403 in β's env). α close-out at `.cdd/releases/3.61.0/283/alpha-closeout.md` is **pending α** at PRA write time. β close-out at `.cdd/releases/3.61.0/283/beta-closeout.md` (incl. tag-push deferral to δ + 4 §9.1 trigger candidates + 2 process-gap candidates surfaced as PRA inputs). γ writes this PRA ahead of the tag and ahead of α's close-out per the disconnect-release model in `CDD.md` §1.4 Phase 6 (β's 3.61.0 release ships the cycle's MCA; δ's later disconnect release crystallizes post-cycle patches + this PRA's emitted skill patches).

### 1. Coherence Measurement

- **Baseline:** 3.60.0 — α A, β A, γ A-, **level L6** (cycle cap L5; trigger: avoidable tooling failure on polling discipline)
- **This release:** 3.61.0 — α A-, β A, γ A-, **level L5** (diff L7 attempted; cycle cap L5 by §9.1 trigger fires below)
- **β's provisional CHANGELOG row** (release commit `585efe3`): `C_Σ A-, α A-, β A, γ A-, L6 cycle cap, L7 diff`. **γ revises cycle cap to L5** — see §4b Cycle Iteration. The diff itself is structurally L7 (eliminates a friction class for all future cycles), but the cycle's execution missed L6 cleanly because lifecycle-skill drift (F2/F3/F4) reached review. Per `CDD.md` §9.1 ("If cross-surface drift reached review, L6 was not met"), cycle level = lowest cleanly-cleared = L5.
- **Delta from 3.60.0:**
  - **α** held A- (provisional pending α close-out). The cycle's diff was substantial (~9 markdown files restructured around the new artifact-channel surface) and α applied 3.60.0's polling-discipline lessons partially: synchronous baseline + transition-only emission landed in `CDD.md` §Tracking and gamma/operator polling snippets. The minus reflects three things: (a) round-1 cross-surface drift (F2/F3/F4 in `release/`, `post-release/` — α's re-audit per `alpha/SKILL.md` §2.6 row 9 missed lifecycle-skill peer surfaces while catching role-skill peer surfaces); (b) α's first commit (`aac0607`) was structurally rewritten when γ simplified the issue body mid-cycle (issue-edit-mid-cycle pattern, β's observation #2); (c) self-application required a same-cycle correction (β's R1 verdict cherry-pick from β-side branch onto α's branch as `8d78514` to make this cycle exemplify its own one-branch rule). The cycle still cleared review in 2 rounds with 0 mechanical findings, so A- not B.
  - **β** held A. R1 review walked the full review skill discipline; 4 findings, all judgment, all in 9 markdown files; mechanical ratio 0%. β narrowed cleanly between rounds. β correctly refused harness implementation instructions at intake (`beta/SKILL.md` §1) and cherry-pick mechanism preserved role-separation audit trail. β surfaced its own polling-failure observation (Monitor missing α-branch transitions; suspected `git fetch --quiet` silent failure) honestly in close-out as a §9.1 input — a process-debt self-observation worth more than zero findings would be without it.
  - **γ** held A- (provisional). Issue #283 quality at dispatch was acceptable (12 numbered ACs, non-goals stated, design constraints linked, work-shape declared, dependencies stated) but had a `.cn/`-vs-`.cdd/` typo that γ caught only on the second read (initial fetch returned a stale-cached body with `.cn/unreleased//`; the actual issue text already had `.cdd/`). γ's mid-cycle issue rewrites (twice — `{N}/{role}.md` → `{N}/` + descriptive filenames) invalidated α's first commit and forced an α-side rework. γ resolved β R1 F1 cleanly via `gamma-clarification.md` on the cycle branch, exemplifying the new artifact-channel coordination. γ's *attempted* mid-cycle scope expansion (γ-creates-branch protocol, ACs 13–17 added to issue body) was correctly rolled back when the operator chose to preserve cycle hygiene over bundling — but the attempted expansion is itself a γ process miss (named decision-point: "scope expansion mid-cycle" should have paused for operator before edits, not after committing local work). The minus reflects: issue-edit-mid-cycle + scope-expansion-after-β-approval, both γ-side process drift caught and corrected, neither blocking close.
  - **Level — L5 (cycle cap L5).** The diff is L7 in intent (eliminates PR-coordination ceremony for all future cycles — a system-shaping MCA). But the cycle execution missed L6 because cross-surface drift reached review (F2/F3/F4 in lifecycle skills). Per §9.1, cycle level caps at L5. The L7 diff is real and ships; the cycle that produced it was L5-execution-quality.
- **Coherence contract closed?** Structurally yes. The cycle's diff implements the protocol the cycle operates under, and the cycle directory at `.cdd/releases/3.61.0/283/` is the integration test: `self-coherence.md` (α R1+R2) + `beta-review.md` (β R1+R2) + `gamma-clarification.md` (γ F1 resolution, two sections) + `beta-closeout.md` (β handoff) all live on one branch, arrived on `main` in one merge commit, with role separation preserved through git authorship + committer fields. Zero PR was created during the cycle; β's intake polled the cycle branch directly (after one operator-surfaced polling-failure exception). α's close-out is the one missing artifact at PRA write time (see §6 CDD Closeout for resolution path).

### 2. Encoding Lag

| Issue | Title | Type | Design | Impl | Lag |
|-------|-------|------|--------|------|-----|
| ~~#283~~ | Replace GH PR workflow with `.cdd/unreleased/{N}/` artifact exchange | feature/process (P1) | converged | **shipped this release** | — |
| #287 | γ creates the cycle branch — α and β only check out (forthcoming) | process (P1) | converged | not started | **new** |
| #286 | Encapsulate α and β behind γ — γ as autonomous in-cycle coordinator | process | converged | not started | **new** |
| #285 | CTB-ize cnos.core + cnos.cdd: agent-module model + composition operators | design | in progress | not started | new |
| #284 | SOUL as agent type: express SOUL as type constraint, not module | design | in progress | not started | new |
| #281 | Redo writer-package design after SOUL.md migrates to skill form | design | converged | not started | new |
| #280 | Decide cdd-shape triad-internal comms formally adopts git-only as default | design | in progress | not started | new |
| #277 | Rewrite SOUL.md to skill form with UIE-V agent loop | docs/process (P1) | converged | not started | growing |
| #275 | CTB Language Spec v0.2 | design | converged | not started | growing |
| #273 | Rebase-collision integrity guard | process (P1) | converged | not started | growing |
| #256 | B4: CI as triadic surface | feature | converged | not started | growing |
| #255 | CDD 1.0 master | tracking | converged | in progress | low |
| #245 | B5: cdd-kata 1.0 | feature | converged | stubs only | growing |
| #244 | kata 1.0 master | tracking | converged | in progress | low |
| #242 | B3: `.cdd/` Phase 1+2 layout | design | converged | not started | low |
| #241 | B6b: DID — triadic identity | design | converged | not started | growing |
| #240 | B6b: CDD triadic protocol | design | converged | partial | low |
| #218 | B8b: cnos.transport.git | design | converged | not started | growing |
| #216 | B7: package-provided commands Phase 4 | feature | converged | not started | growing |
| #193 | Orchestrator runtime | feature | converged | not started | growing |
| #192 | Runtime kernel rewrite | feature | converged | Phase 4 complete | low |

(Older items #190 / #189 / #186 / #181 / #175 / #170 / #168 / #156 / #154 / #153 / #149 / #101 / #100 / #87 remain open and stale; not promoted to the active lag table.)

**Growing-lag count:** 9 (3.60.0 PRA had 8; #283 closed but #284/#285/#286/#287/#280/#281 net add to the active set). Net direction: still flat-to-up. The freeze is holding the long-stale items steady but new design surfaces opened during 3.60.x continue to land in the table.

**MCI/MCA balance:** **Freeze continues** — 9 growing-lag items, well over the 3-issue threshold.

**Selection note for next cycle:** #287 (γ creates the cycle branch) is the natural next MCA — direct response to #283 β-R1 F1's branch-discovery friction (β's verdict landed on the harness-given branch instead of the cycle branch), the same five surfaces are touched, change is markdown-only, and γ's rolled-back mid-cycle expansion shows the work is well-understood. #286 (encapsulation) sits behind #287 + a separate `cn dispatch` CLI cycle (its hard precondition is harness-side work, see #286 Dependencies).

### 3. Process Learning

**What went wrong:**

- **First-cycle-of-new-protocol friction is structural** (β observation #1). The cycle implements the protocol it operates under. α's first signal (`1aaf9fb`) was made under a stale AC text shape (`{N}/{role}.md`) and reworked when γ surfaced the issue-body simplification mid-cycle. β's first verdict (`1ceb99c`) landed on the harness-given β-side branch (`claude/implement-beta-skill-loading-ZXWKe`) because the harness's branch instruction predated γ's F1 resolution naming the cycle branch as canonical. Both behaviors required a same-cycle correction (α's `1aaf9fb` rework, β's verdict cherry-pick at `8d78514`). Pattern: a new coordination spec cannot be tested by a cycle that pre-dates it; the *implementing* cycle is always one half-step behind its own output. This is irreducible, but the corrections themselves can be lighter-weight (see process-gap suggestions below).

- **Issue-edit-mid-cycle invalidates committed α work** (β observation #2). γ rewrote the issue body twice during α's first checkpoint window — once to switch the path layout from `{N}/{role}.md` to `{N}/` with role-distinguished filenames; once earlier to fix a `.cn/`-vs-`.cdd/` typo (γ initially read a stale-cached body). α stopped, re-read, and reworked in `1aaf9fb`. Documented in α's `self-coherence.md` Known debt #4 and γ's `gamma-clarification.md`. Compounded with β's Monitor polling failure (next bullet).

- **β's `Monitor` polling silently dropped α-branch transitions** (β observation #3). During the round-2 dispatch window, α's branch advanced through three SHAs without β's Monitor emitting `branch-update:` events. β's task IDs `b6vala0kx`, `b2m54i3kr`, `b3ak6xcyg`, `beu5utmvj` all timed out without preceding transition events. Suspected mechanism: `git fetch --quiet origin` returning silently with stale data (auth/network flake masked by `--quiet`), so the per-iteration refs comparison saw no transition. Operator surfaced the activity manually with "A did something — didn't you see?" The new-spec polling shape is structurally correct; the `git fetch` reliability layer underneath is an implicit dependency the spec does not name. γ's own Monitor across this cycle saw the same timeout pattern (re-armed 5+ times before any α-side transition surfaced). Same root pattern as 3.60.0's polling failures (different mechanism, same class).

- **Authority-surface conflicts cluster in lifecycle skills, not role skills** (β observation #4). F2 (`release/SKILL.md`), F3 (`post-release/SKILL.md`), F4 (`post-release/SKILL.md`) — three of four R1 findings landed in lifecycle skills (`release/`, `post-release/`). Zero R1 findings landed in α's, β's, or γ's role skills, despite the diff modifying all three. α's re-audit per `alpha/SKILL.md` §2.6 row 9 was structurally correct in naming the audit, but the audit's checklist did not distinguish role-skill peer surfaces (which α covered) from lifecycle-skill peer surfaces (which α missed). Pattern: lifecycle skills are downstream of role skills (Tier 1c pairing), and downstream surfaces drift more than upstream when an upstream contract changes.

- **γ-side process drift on issue editing and scope expansion** (γ self-observation). γ edited the issue body twice during the cycle (typo fix + path simplification — both legitimate "clarify ambiguity in scope" actions per `gamma/SKILL.md` §2.5) and *then* attempted a third, scope-expansion edit (γ-creates-branch, ACs 13–17) after β had already approved on the original 12 ACs. The third edit was rolled back when the operator chose cycle hygiene over bundling. First two edits were γ-correct but compounded with β's polling failure to invalidate α's first commit. Third edit was γ-wrong: "scope expansion mid-cycle" is a named decision-point per the encapsulation spec γ just filed (#286 §1.4 AC6), and γ should have paused for operator *before* editing local files. γ has the rule in the issue it just filed; γ does not yet have it in the executable spec. Recursion of the kind #283 itself produced.

**What went right:**

- **The cycle is its own integration test for the new protocol.** β observation #1 framed this; γ confirms structurally. `.cdd/releases/3.61.0/283/` carries `self-coherence.md` + `beta-review.md` + `gamma-clarification.md` + `beta-closeout.md` — every coordination artifact arrived on `main` through one merge commit, with role separation preserved by git authorship + committer fields. Zero PR was opened during the cycle. This is what the spec promised, run end-to-end on its own implementation.

- **Cherry-pick preserves role-separation audit trail** (β observation #5). β's R1 verdict `1ceb99c` (author `beta@cdd.cnos`) was cherry-picked onto α's branch as `8d78514` (committer `alpha@cdd.cnos`, author preserved). The audit trail correctly shows β authored the review and α landed it on the cycle branch. `review/SKILL.md` §7.1 names this mechanism. The cherry-pick mechanism extends the contract to mid-cycle migrations between branches, which is exactly what "first cycle of new protocol" required.

- **β's role-conflict refusal worked as designed** (β observation #6). β was dispatched onto a pre-set harness branch with α-style instructions ("DEVELOP / COMMIT / PUSH"). β refused at intake per `beta/SKILL.md` §1, reported the conflict to the operator as a status report, continued β intake. The harness's α-style instruction was not honored.

- **β narrowed cleanly between rounds.** R1: 4 findings (1 D structural + 3 C/B authority-surface). R2: 0 findings. Mechanical ratio 0% across both rounds. Review-rounds = 2 (at §9.1 threshold but not exceeding).

- **β surfaced its own polling failure honestly.** Rather than noting the Monitor miss as a private observation, β included it in close-out as a §9.1 input. γ does not have to dig for it.

- **γ-creates-branch scope expansion was correctly rolled back.** γ proposed bundling, operator declined, γ executed the discard cleanly (`git reset --hard origin/...`, issue body reverted to 12-AC form, follow-up issue #287 filed). Cycle hygiene preserved. γ's process miss (editing before pausing) was caught at the operator interface; the rollback itself was clean.
