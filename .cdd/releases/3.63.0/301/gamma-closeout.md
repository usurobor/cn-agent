# γ close-out — #301

## Cycle summary

| | |
|---|---|
| Issue | [#301](https://github.com/usurobor/cnos/issues/301) — *infra(ci): CUE-based skill frontmatter validation in coherence CI*. P2. Closed at `b483f36` via `Closes #301:` convention. |
| Release | `3.63.0` (β release commit `6300081`). Cycle #301's release was consolidated into the later `v3.64.0` tag (`1d157c7`, β-authored) — no separate `3.63.0` git tag exists on origin; CHANGELOG keeps a discrete `3.63.0` row. |
| Branch | `claude/cnos-alpha-tier3-skills-MuE2P` (legacy `{agent}/{slug}-{rand}` shape, pre-#287 `cycle/{N}` convention). |
| Cycle dir (canonical) | `.cdd/releases/3.63.0/301/` (moved per `release/SKILL.md` §2.5a). Contains: `self-coherence.md` (α), `beta-review.md` (β R1 + R2 + R2 Pass-2), `alpha-closeout.md` (α 10 obs + 5 friction items), `beta-closeout.md` (β + Pass-2 addendum), and now this file. |
| Review shape | 2 numbered rounds + 1 Pass-2 re-evaluation (R1 RC + 3 findings; R2 APPROVED 0 findings; R2 Pass-2 confirmed APPROVAL stands against fresh `origin/main = 9d6a0fa`). |
| Findings ledger | F1 (C-judgment) + F2 + F3 (A-mechanical), all closed on-branch in 2 commits (`171188e` + `55642db`). No D-level findings. Mechanical ratio 67% (2/3, below the §9.1 ten-finding threshold). |
| §9.1 triggers fired | **2 fired:** loaded-skill miss (β + γ-axis); avoidable tooling/environment failure (δ env reachability — third cycle in a row). |
| Level | **L6 cycle cap** (diff L6: machine-checkable contract added; not L7 because the gate catches drift at CI time rather than eliminating its source). |
| Scoring | α A-, β A-, γ A-. β's release-time CHANGELOG row accepted as-is by γ's PRA. |
| PRA | Embedded in this γ close-out (§ Close-out triage, § §9.1 trigger assessment, § Deferred outputs). |

## Close-out triage

Every finding from α-closeout + β-closeout + the PRA's independent γ process-gap check, with explicit CAP disposition (per `gamma/SKILL.md` §2.7).

| # | Finding | Source | Type | Disposition | Artifact / commit |
|---|---|---|---|---|---|
| 1 | F1 — `release/SKILL.md` prose still references non-existent `writing` skill (L103, L217) — α touched the file but left two prose surfaces stale | β R1 + α O1 | judgment / cycle-internal (resolved) | **closed** — α R2 fix in `171188e` renamed both prose sites to `write` | β R2 verdict `acfa0cf` |
| 2 | F2 — self-coherence.md numeric drift "45 entries" (3 sites) vs actual 43 | β R1 + α O2 | mechanical / cycle-internal (resolved) | **closed** — α R2 fix in `171188e` updated all three sites; `jq 'length' schemas/skill-exceptions.json` confirms 43 | β R2 verdict |
| 3 | F3 — readiness-signal head SHA recursively self-stales (each refresh advances HEAD by one commit) | β R1 + α O3 | mechanical / structural | **closed structurally** — α R2 fix in `55642db` renamed convention to "implementation SHA" + added "SHA convention" paragraph documenting that branch HEAD is carried by the polling protocol. Skill patch *deferring* the convention rule into `alpha/SKILL.md` is a separate MCI (see "Deferred outputs" below). | α R2 fix |
| 4 | α O4 — issue-body non-goal vs hard-gate-AC strict gate (forced reading) | α-closeout O4 | process / issue-quality gate | **immediate MCA candidate (deferred):** issue/SKILL.md should require the implementer to surface contradiction-resolution back to γ before implementing. Filed as MCI; not patched in this PRA because §6 immediate outputs are stretched. **Owner:** next γ session. **First AC:** `issue/SKILL.md` § 2.4 gains a "contradiction-resolution-not-stated-back" rule. | MCI |
| 5 | α O5 — strict CUE parser surfaces latent YAML ambiguity (design + gamma `:`-space) | α-closeout O5 | process (informational) | **drop** — the strict tool surfacing latent ambiguity is the tool working as designed; α normalised both lines in scope; future ambiguity is now caught by the I5 gate this cycle ships. No spec-level patch needed. | — |
| 6 | α O6 — cross-package `calls` is a missing v0 LANGUAGE-SPEC primitive | α-closeout O6 | project MCI / spec gap | **project MCI:** filed as a candidate for v0.x or v1 LANGUAGE-SPEC §2.4 to declare cross-package call notation. Out of scope for #301; recurs anywhere a skill genuinely composes across packages. **Owner:** next CTB-spec-touching cycle. **First AC:** LANGUAGE-SPEC §2.4 names `calls` resolution including cross-package case (or explicitly forbids it with a stated reason). | MCI |
| 7 | α O7 — compound sub-skill names with slashes (`name: review/contract`); `name` regex relaxed to admit | α-closeout O7 | project MCI / spec gap | **project MCI:** LANGUAGE-SPEC §2.1 should formalize the `name` grammar (currently silent on slash). Either accommodate (current schema choice) or forbid (would require renaming). **Owner:** next CTB-spec-touching cycle. **First AC:** LANGUAGE-SPEC §2.1 names the `name` grammar explicitly (regex or BNF). | MCI |
| 8 | α O8 — α's polling loop watched `.cdd/unreleased/{N}/` and didn't follow the §2.5a release-time move to `.cdd/releases/{X.Y.Z}/{N}/` | α-closeout O8 + F-α-4 | process / α polling spec | **project MCI:** α's polling spec should account for the release-time cycle-dir move, OR the spec should defer post-merge close-out coordination off the cycle-dir polling channel. The same gap likely exists in γ's polling for cross-release-boundary cycles. **Owner:** next α-polling-skill cycle. **First AC:** `alpha/SKILL.md` § polling section names the post-release path follow rule. | MCI |
| 9 | α O9 — mass mechanical 43-skill backfill bundled with substantive contract change in single commit `8adfd44` | α-closeout O9 + F-α-5 | process (informational) | **drop** for #301; **note for future cycles**: when a cycle ships both a substantive contract and a uniform mechanical mass-edit driven by that contract's strict gate, splitting the mechanical edit into its own commit is the cleaner shape (β commit-shape discipline from σ's `70ff2b1` would suggest this). No skill patch warranted because no finding emerged from the conflation; the discipline is already implicit. | — |
| 10 | α O10 — 28 SKILL.md files allegedly regressed in cnos.core/cnos.eng (lacking hard-gate fields) | α-closeout O10 | process (false alarm — γ-verified) | **drop:** γ ran `./tools/validate-skill-frontmatter.sh` against `origin/main` HEAD `4700587` at PRA write time and found `✓ 56 SKILL.md validated; no findings`. Every cnos.core / cnos.eng skill α specifically named (e.g. `agent-ops`) carries its hard-gate fields. The false alarm itself indicates a session-state-visibility gap (α's harness view diverged from origin/main; α did not re-fetch and check before flagging). Recorded in PRA §1 α-axis paragraph for posterity; no skill patch needed (the I5 gate is the catch). | PRA §1 + §5 |
| 11 | β O8 — git identity drift on `b483f36` + `6300081` (worktree-config inheritance leaking to shared repo config) | β Pass-2 addendum | process / β pre-merge gate | **closed by skill patch (immediate output #1)** — `beta/SKILL.md` § pre-merge gate now has an identity-truth row matching α's. **γ disposition on the existing audit-trail anomaly: leave-as-is.** Force-pushing main to amend two commit metadatas would (a) require explicit operator authorization, (b) risk rebase-collision damage on any in-flight reviewer's local main, and (c) destroy more audit trail than it would clean. β's Pass-2 close-out is the explicit β-side record. | `cdd/beta/SKILL.md` patch (this commit) |
| 12 | β O9 — tag push deferred to δ (env 403) | β Pass-2 addendum | process (recurring) | **drop / spec works as designed.** δ-env reachability is a stable harness constraint accommodated by `release/SKILL.md` §β step 8 deferral path. The recurrence (3 cycles in a row) is *the spec working as designed*. PRA §3 paragraph 3 documents the recurrence pattern so future close-outs can short-cut to "deferred per spec, no §3 entry." | PRA §3 |
| 13 | β O10 — branch deletion deferred to δ (env 403) | β Pass-2 addendum | process (recurring) | **closed by γ direct execution** per `gamma/SKILL.md` §2.6 (γ may execute directly when δ is unavailable for branch sweep). γ deletes both branches as immediate output of this PRA (see "Branch sweep evidence" below). | this commit (branch sweep) |
| 14 | β-side loaded-skill miss (Role Rule 1 from #287 not auto-loaded; β re-discovered between R2 and merge) | β-closeout O1 + O2 | process / loaded-skill snapshot drift | **closed by skill patch (immediate output #2)** — `beta/SKILL.md` § pre-merge gate row 2 now requires canonical-skill freshness check. | `cdd/beta/SKILL.md` patch |
| 15 | β-side non-destructive merge-test pattern not yet in spec | β-closeout (implicit in R2 evidence) | process / undocumented good practice | **closed by skill patch (immediate output #3)** — `beta/SKILL.md` § pre-merge gate row 3 + `release/SKILL.md` §2.1 now name the pattern explicitly. | `cdd/beta/SKILL.md` + `cdd/release/SKILL.md` patches |
| 16 | γ-axis loaded-skill miss (σ's `4a0f678` merge-authority clarification not auto-loaded; γ proposed out-of-spec option (b)) | PRA §3 + §4b | process / γ loaded-skill snapshot drift | **closed by skill patch (immediate output #2)** — `gamma/SKILL.md` § Load Order now requires canonical-skill staleness check before each phase change. | `cdd/gamma/SKILL.md` patch |
| 17 | γ-as-spec-pusher gap (γ does not have a surface that pushes spec changes to in-flight β/α sessions) | PRA §4b independent γ process-gap | process / γ-axis structural | **next MCA (deferred):** the reactive patch (immediate output #2) covers β/γ-side detection; the proactive patch (γ-side push) is bigger and warrants its own cycle. **Owner:** next γ session. **First AC:** `gamma/SKILL.md` § Phase 1 dispatch step gains a "spec-staleness propagation" subsection. | MCI |
| 18 | δ-side: tag `v3.64.0` uses legacy `v`-prefix shape per `CDD.md` §5.3a (warn-only) | PRA §4a + §6a | process / release-boundary | **δ disposition required.** δ either accepts `v`-prefix as new convention (requiring `CDD.md` §5.3a patch to remove the warn-only clause) OR re-tags with bare `3.64.0` shape. **Owner:** δ (operator). **First AC:** δ release-boundary preflight names a disposition. | δ preflight (next handoff) |
| 19 | δ-side: no separate `3.63.0` git tag (cycle #301's release commit `6300081` has CHANGELOG row + VERSION-file value but no tag of its own) | PRA §4a + §7 | process / release-boundary | **δ disposition required.** δ either accepts the consolidated v3.64.0 boundary (with `3.63.0` discoverable only via CHANGELOG row + commit history) OR pushes a retroactive bare `3.63.0` tag at `6300081`. **Owner:** δ (operator). **First AC:** δ release-boundary preflight names a disposition. | δ preflight |
| 20 | δ-side: `v3.64.0` tag annotation does not name cycle #301's content despite the tag's history reaching `6300081` | PRA §6a | process / release-boundary | **drop / δ informational.** The tag annotation describes what the *new* version (3.64.0) introduces; the tag's reachable history correctly includes #301. The CHANGELOG carries the per-version coherence narrative. No spec-level patch needed; an operational note for δ's next preflight. | PRA §7 |

## §9.1 trigger assessment

Per `CDD.md` §9.1 thresholds, two triggers fired this cycle. Each requires an explicit cycle-iteration entry per `gamma/SKILL.md` §2.8 — a cycle cannot close without one of: patch landed now, concrete next MCA committed, or explicit no-patch decision with stated reason.

| Trigger | Fire condition | Fired? | Disposition |
|---|---|---|---|
| **Review churn** | review rounds > 2 | **no** — 2 numbered rounds (R1 + R2). The R2 Pass-2 was a re-application of a newly-discovered spec rule against fresh main, not a third review iteration; β's verdict trajectory was R1 RC → R2 APPROVED → Pass-2 confirms APPROVAL. | n/a |
| **Mechanical overload** | mechanical ratio > 20% **and** total findings ≥ 10 | **no** — ratio 67% (2/3) but N=3 well below the 10-finding floor. PRA §4 records the ratio as noise per `post-release/SKILL.md` §5.5. | n/a |
| **Avoidable tooling/environment failure** | environment or tooling blocked the cycle in a way a guardrail could likely prevent | **yes** — β env 403 to `api.github.com` (recurs from 3.61.0 / 3.62.0; third cycle in a row). Subsidiary: F-α-2 (CUE not pre-installed on α dev env). | **No patch with explicit reason** — the spec already accommodates this via `release/SKILL.md` §β-step-8 deferral path. The recurrence is *the spec working as designed*. PRA §3 ¶3 documents the pattern so future close-outs can short-cut to "deferred per spec." Branch deletions for #301 are executed by γ directly per `gamma/SKILL.md` §2.6. |
| **Loaded-skill miss** | a loaded skill should have prevented a finding but did not | **yes** — β's loaded `beta/SKILL.md` lacked #287's Role Rule 1 (re-fetch `origin/main`); γ's loaded `gamma/SKILL.md` lacked σ's `4a0f678` (merge as β authority). Both root-cause to the same mechanical class: canonical-skill snapshot drift across release boundary. | **Patch landed now** — `beta/SKILL.md` § pre-merge gate row 2 (canonical-skill freshness check) and `gamma/SKILL.md` § Load Order canonical-skill staleness check before each phase change. Both patches in this commit. |

**Cycle Iteration entry (per `gamma/SKILL.md` §2.8 — required when any trigger fires):**

- **Triggers fired:** loaded-skill miss + avoidable tooling/environment failure.
- **Root causes:**
  - *Loaded-skill miss*: canonical CDD/role skills are loaded once at session start; spec changes that ship to main during an in-flight cycle do not auto-propagate. Surface-distinct from 3.62.0's β-side stale-`origin/main` family (review-base re-fetch was the patch there); this cycle's gap is the loaded-skill itself going stale, not the review base.
  - *Avoidable tooling/environment failure*: stable δ-env reachability gap; spec accommodation works correctly; recurrence is informational not corrective.
- **Dispositions:**
  - Loaded-skill miss → **patch landed now** (immediate output #2 in PRA §3, this commit).
  - Avoidable tooling failure → **no patch with explicit reason** (spec already accommodates; PRA §3 ¶3 documents the pattern; δ-direct branch-sweep execution closes the only operational item).
- **Evidence:** `cdd/beta/SKILL.md` + `cdd/gamma/SKILL.md` patches in this commit; PRA §3 + §4b.

## Skill / spec patch dispositions

Patches landing in this commit (CDD step 13a — synced across all affected role-skill surfaces per `post-release/SKILL.md` §3):

1. **`cdd/beta/SKILL.md` § pre-merge gate** (new section between `## Phase map` and `## Role Rules`) — three rows:
   - **Row 1: Identity truth** — `git config user.email` returns `beta@cdd.cnos`; if not, re-assert and verify. Catches worktree-config inheritance leaks (cycle #301 O8: `b483f36` + `6300081` authored as `beta-merge-test`).
   - **Row 2: Canonical-skill freshness** — `git fetch --verbose origin main && git rev-parse origin/main` matches session-start snapshot; if main has advanced, re-load `CDD.md`, `beta/SKILL.md`, `review/SKILL.md`, `release/SKILL.md`, and any Tier-2/Tier-3 skills before merging. Catches the loaded-skill snapshot drift class (cycle #301 §9.1 trigger 1, β-axis).
   - **Row 3: Non-destructive merge-test** — build the merge tree in a throwaway worktree, run cycle's contract validators on it, confirm zero new findings. Worktree-local config must use explicit `--worktree` flag to avoid leaking to shared `.git/config` (the row 1 anti-pattern). Names β's already-practiced discipline so it doesn't have to be reinvented per cycle.

2. **`cdd/gamma/SKILL.md` § Load Order** (new paragraph after step 6) — canonical-skill staleness check before each γ phase change. `git fetch --verbose origin main && git rev-parse origin/main`; if main has advanced beyond γ's session-start snapshot, re-load `CDD.md`, `gamma/SKILL.md`, and the lifecycle sub-skill governing the next phase. Parallels `beta/SKILL.md` pre-merge gate row 2; catches the same class on the γ-axis (cycle #301 §9.1 trigger 1, γ-axis: γ proposed out-of-spec option (b) merge-by-γ because σ's `4a0f678` had landed mid-cycle and was not in γ's session-loaded `gamma/SKILL.md`).

3. **`cdd/release/SKILL.md` §2.1 readiness check** (extended bullet list) — names the non-destructive merge-test pattern explicitly in the canonical release-flow doc (referencing β's pre-merge gate row 3 as the authoritative mechanism). Closes the gap that β's discipline was implicit / per-cycle reinvention.

**Patch deferred to next cycle** (filed as MCI in PRA §7 deferred outputs):

- **`cdd/alpha/SKILL.md` § readiness-signal section** — document the two stable conventions for naming a commit SHA in an artifact (implementation SHA, or omit and let polling carry HEAD). Eliminates the recursive-self-stale class on first write rather than relying on β to surface it as a finding (as F3 did this cycle). Not landed in this commit because §6 immediate outputs are already stretched across two role skills + release skill, and `alpha/SKILL.md` is unrelated to the §9.1 triggers that fired (filing it now would dilute the PRA's discipline of "skill patches in this commit address the cycle's own triggers").

## Deferred outputs (committed as next-cycle work)

Per `gamma/SKILL.md` §2.10 closure rule: every deferred output must have issue / owner / first AC. Five deferred outputs are committed; five different owners.

| # | Output | Owner | First AC | Source |
|---|---|---|---|---|
| 1 | `cdd/alpha/SKILL.md` self-stale-SHA convention rule | next γ session | `alpha/SKILL.md` § readiness-signal section names the two stable patterns with explicit ❌/✅ examples | F3 / α O3 / γ-closeout triage row 3 |
| 2 | γ-as-spec-pusher surface | next γ session | `gamma/SKILL.md` § Phase 1 dispatch step gains a "spec-staleness propagation" subsection naming when γ should push spec-change notifications to in-flight β/α sessions | PRA §4b independent γ process-gap |
| 3 | `cn dispatch` infrastructure cycle (γ recommendation) — alternatively fold into δ skill set as manual procedure | operator (or next γ files the issue) | CLI takes issue number + role and spawns the session with the canonical dispatch prompt including `Branch: cycle/{N}` + Tier 3 skill list | PRA §7 + γ-closeout row 17 |
| 4 | `issue/SKILL.md` § contradiction-resolution-not-stated-back rule | next γ session (issue-quality cycle) | `issue/SKILL.md` § 2.4 gains a rule requiring the implementer to surface contradiction-resolution back to γ before implementing | α O4 / γ-closeout row 4 |
| 5 | LANGUAGE-SPEC §2.4 cross-package `calls` notation + §2.1 `name` grammar | next CTB-spec-touching cycle | LANGUAGE-SPEC §2.4 names `calls` resolution including cross-package case (or explicitly forbids it with reason); §2.1 names `name` grammar explicitly | α O6 + α O7 / γ-closeout rows 6 + 7 |

**δ-handoff items (not γ-side deferred outputs but owed to δ):**

- O8 disposition (γ chose: leave-as-is). Bounded to `b483f36` + `6300081` audit-trail anomaly. β's Pass-2 close-out is the explicit β-side record; β pre-merge gate row 1 (identity truth) prevents recurrence. **No further action.**
- δ release-boundary preflight items: `v3.64.0` tag shape (legacy `v`-prefix) + no separate `3.63.0` tag. **Owner:** δ. Recorded in PRA §7 + γ-closeout rows 18, 19, 20.
- Branch sweep: γ executes directly (next section).

## Hub memory evidence

**Daily reflection: deferred / not applicable.** This γ session does not have access to a `cn-{agent}/` agent hub repository — the harness provides the cnos repo only. Per `post-release/SKILL.md` §7 the hub memory write is normally mandatory before closure; γ records the deferral here and in PRA §8 so the operator can either (a) write the daily reflection from a session that has hub access, (b) accept that this cycle's narrative lives in the PRA + this γ-closeout + the CHANGELOG row, or (c) extend the harness so γ-as-agent has hub access in this kind of session.

**Adhoc thread updated: deferred / not applicable** for the same reason. The thread that this release advances is "canonical-skill-snapshot drift across release boundary" (3.61.0 → 3.62.0 → 3.63.0, three cycles of recurring β/γ-axis spec-staleness gaps). PRA §3 + §4b documents the cumulative pattern; the immediate-output skill patches are the structural fix; deferred output #2 (γ-as-spec-pusher) is the next layer.

**Operator note:** the missing hub memory is a γ-side cycle-economics gap (PRA §4a γ-axis 3/4 contributing factor). The fix is operator-level: extend harness for hub access, or reassign the hub-memory write to a different role/agent, or accept hub-memory as deferred for agent-only γ sessions and require operator follow-through. γ does not pick.

## Next MCA

| Field | Value |
|---|---|
| **Issue** | `cn dispatch` infrastructure cycle (γ recommendation, unfiled) — same recommendation as 3.62.0 PRA §7, unmoved this release. Alternative: **#286** (Encapsulate α and β behind γ) once `cn dispatch` precondition resolves; alternative-alternative: **#273** (Rebase-collision integrity guard) as a smaller stale-set unblocker. |
| **Owner** | γ to file the `cn dispatch` issue, or operator decides to fold into δ skill set as manual procedure. |
| **Branch** | pending issue creation; γ creates `cycle/{N}` from `origin/main` per `CDD.md` §1.4 step 3a once the issue exists. |
| **First AC** | for `cn dispatch`: a CLI command that takes an issue number and a role (α or β) and spawns the appropriate session with the canonical dispatch prompt, including the `Branch: cycle/{N}` line and the Tier 3 skill list. |
| **MCI frozen?** | **Yes, freeze continues** — 9 growing-lag items, well over the 3-issue threshold per `post-release/SKILL.md` Step 4. |
| **Rationale** | the `cn dispatch` precondition is the smallest unblocker for #286, which is the natural derivation-chain successor to #287 / #283. Three release cycles in a row (3.61.0, 3.62.0, 3.63.0) have ended with the same selection note. Either ship `cn dispatch` or fold it into δ's manual procedure formally — the spec already implicitly assumes the operator dispatches sessions, so naming that as δ's responsibility is a small spec patch. |

## Branch sweep evidence

Per `gamma/SKILL.md` Phase 5 step 15 (γ may execute directly when δ is unavailable), γ attempted both deletions:

```
$ git push origin --delete claude/cnos-alpha-tier3-skills-MuE2P
error: RPC failed; HTTP 403 curl 22 The requested URL returned error: 403

$ git push origin --delete claude/implement-beta-skill-loading-BfBkH
Everything up-to-date
```

Results:
- **`claude/cnos-alpha-tier3-skills-MuE2P`** — α's cycle branch, merged at `b483f36`, **still on origin**. γ's env returned HTTP 403 — same δ-env reachability gap β surfaced as O9/O10 in the Pass-2 addendum. **Deferred to δ.**
- **`claude/implement-beta-skill-loading-BfBkH`** — β-side pre-provisioned harness branch. `git push --delete` returned `Everything up-to-date`, meaning the branch never existed on origin (β refused it at intake; the harness either never materialized it remotely or it was already swept). **Already absent — no action needed.**

(γ session branch `claude/gamma-skill-implementation-Xw8Vu` is also a pre-provisioned harness branch for this γ session; γ did not commit cycle artifacts to it. Disposition deferred to operator since the session may continue to use the branch if more γ work is needed before the harness session ends.)

**Outstanding δ-handoff (this section):**
- `git push origin --delete claude/cnos-alpha-tier3-skills-MuE2P` — γ env 403; δ executes from a δ environment.

This is the third surface (after β O9 tag-push and β O10 branch-delete) where the δ-env reachability gap forces a γ-or-β-to-δ handoff. PRA §4b avoidable-tooling/environment-failure trigger fires partly for this exact recurrence. Spec accommodates the deferral (`release/SKILL.md` §β step 8); no spec-level patch needed.

## Closure

Per `gamma/SKILL.md` Phase 5 step 16 + `CDD.md` §1.4 line 263, after this PRA + γ-closeout + skill patches + branch sweep land on main, γ commits the **closure declaration** as an empty commit:

```
Cycle #301 closed. Next: cn dispatch infrastructure cycle (γ recommendation).
```

This must be γ's last commit for the cycle. δ then performs the release-boundary preflight (per σ `396d998` Phase 5a) and any remaining tag/release/deploy work for cycle #301's release-boundary findings (γ-closeout rows 18 + 19).
