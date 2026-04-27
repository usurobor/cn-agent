## Post-Release Assessment — 3.60.0

> Cycle **#235** (B2: `cn build --check` validates command entrypoints + skill SKILL.md presence; single-issue bundle). PR #276 squash-merged at `d814e16`; release commit `cce905e` by β; α close-out at `.cdd/releases/3.60.0/alpha/CLOSE-OUT.md`; β close-out at `.cdd/releases/3.60.0/beta/CLOSE-OUT.md` (incl. tag-push deferral to δ). Tag pushed by δ; `release.yml` completed green at 2026-04-26T23:11:33Z; GitHub release `3.60.0` published with 4-platform binaries + 5 tarballs + index.json + checksums.txt.

### 1. Coherence Measurement

- **Baseline:** 3.59.0 — α A-, β A, γ A-, **level L6** (diff L7; cycle cap L6)
- **This release:** 3.60.0 — α A, β A, γ A-, **level L6** (cycle cap L5)
- **Delta:**
  - **α** improved to A. Zero-finding cycle: 1 review round, 0 findings on the diff. Canonical parser (`pkg.ParseFullManifestData`) reused — no parallel parser introduced, directly applying the §2.17 diff-local generalization that 3.59.0's F1 surfaced. Contract-implementation confinement clean: empty-entrypoint and path-traversal rejection added beyond strict AC wording. Surface-enumeration test planning explicit: 9 tests cover pass + fail + exempt for both ACs. AC4 deliberate-break dry-run gave end-to-end CI-gate verification before push. α close-out is thorough: names polling-discipline failure honestly (webhook priority inversion, surface-selection error on `get_status` vs `get_check_runs` vs `get_reviews`), records four reusable patterns (pre-implementation peer enumeration, AC refinement via live-package check, copy+inject+run dry-run, and the §2.17 generalization applied successfully). The A reflects a clean diff cycle with honest process self-assessment.
  - **β** held A. Round-1 review walked the full `cdd/review/SKILL.md` discipline (§2.0 through §2.3); four notes (N1–N4) explicitly scoped as γ-deferred observations, not findings. β positive checks are precise: AC1 contract-implementation confinement verified with three mental edge cases on path-traversal; AC2 authority-surface alignment confirmed against `internal/activation/index.go`; AC3 negative-space coverage verified per `eng/test`; AC4 evidence depth accepted with structural reasoning. β's own polling miss (PR #276 existed but `list_pull_requests head=owner:branch` returned empty due to slash in branch name) was self-named as a CDD §Tracking gap with a concrete patch recommendation. β-step-8 tag-push deferral handled correctly per protocol. The A reflects sharp review, honest process self-assessment, and clean protocol execution.
  - **γ** holds A-. Not the cycle's γ (this is a δ-authored PRA). The A- captures: issue #235 was high-quality at dispatch (all γ §2.4 fields present, Tier 3 skills named, design constraints linked, non-goals stated); dispatch was right-sized (small substantial MCA). The minus reflects: issue body cited stale `docs/alpha/INVARIANTS.md` path (β's N4); harness branch naming (`claude/alpha-tier-3-skills-M8Vce`) continues the non-canonical CDD §4.2 pattern from 3.59.0 — same observation, still unresolved.
  - **Level — L6 (cycle cap L5).** The diff is L6: cross-surface coherence verified (build-time validator integrates with runtime activation/discover using the same authority surfaces; `pkg.ContentClasses`/`pkg.ClassSkills`/`pkg.ParseFullManifestData` shared across all three; no parallel authority introduced). The cycle cap is L5 because the §9.1 "avoidable tooling/environmental failure" trigger fired: both α and β independently missed each other's PR activity due to polling-discipline failures (α's webhook priority inversion, β's slash-in-branch filter failure). Operator intervention was needed in both directions. The diff itself was clean (zero findings), but the cycle's coordination required external correction. CHANGELOG provisional scores (α A, β A, γ A-, L6 cycle cap L5) **confirmed without revision**.
- **Coherence contract closed?** Yes. Tag `3.60.0` on origin at `cce905e`. `release.yml` completed green (4-platform matrix). GitHub release published with all expected assets. `cn build --check` validates entrypoints and skill paths on every I1 CI run. Pre-existing valid packages pass; injected breaks fail with named diagnostics and exit 1.

### 2. Encoding Lag

| Issue | Title | Type | Design | Impl | Lag |
|-------|-------|------|--------|------|-----|
| ~~#235~~ | `cn build --check` validates entrypoints + skill paths | feature (P1) | converged | **shipped this release** | — |
| #277 | Rewrite SOUL.md to skill form with UIE-V agent loop | docs/process (P1) | converged | not started | **new** |
| #275 | CTB Language Spec v0.2 | design | converged | not started | growing |
| #273 | Rebase-collision integrity guard | process (P1) | converged | not started | growing |
| #262 | `cn pack` packlist derivation (AC1/AC5 remaining) | feature | converged | partial (3.58.0) | low |
| #256 | CI as triadic surface | feature | converged | not started | growing |
| #255 | CDD 1.0 master | tracking | converged | in progress | low |
| #245 | cdd-kata 1.0 | feature | converged | stubs only | growing |
| #244 | kata 1.0 master | tracking | converged | in progress | low |
| #242 | .cdd/ Phase 1+2 layout | design | converged | not started | low |
| #241 | DID: triadic identity | design | converged | not started | growing |
| #240 | CDD triadic protocol | design | converged | partial | low |
| #218 | cnos.transport.git | design | converged | not started | growing |
| #216 | Migrate non-bootstrap kernel commands to packages | feature | converged | not started | growing |
| #199 | Stacked v3.39.0 + v3.40.0 PRAs | docs | converged | not started | growing |
| #193 | Orchestrator runtime | feature | converged | not started | growing |
| #192 | Runtime kernel rewrite | feature | converged | Phase 4 complete | low |

(Older items #190 / #189 / #186 / #181 / #175 / #170 / #168 / #156 / #154 / #153 / #149 / #101 / #100 / #87 remain open and stale; not promoted to the active lag table.)

**Growing-lag count:** 8 (same as 3.59.0 PRA — #235 closed, #277 opened). Net direction: flat. The MCA freeze continues to hold the backlog steady but is not shrinking it.

**MCI/MCA balance:** **Freeze continues** — 8 growing-lag items, well over the 3-issue threshold.

### 3. Process Learning

**What went wrong:**

- **Polling discipline gap — same shape, fifth consecutive cycle.** Both α and β independently missed each other's activity. α treated webhook as primary and `get_reviews` was never re-polled after PR-open. β's `list_pull_requests head=owner:branch` filter silently failed because the branch name contained a slash. Operator intervention required in both directions. CDD §Tracking names the rule ("all roles must track via periodic polling, not event subscriptions"); neither role followed it. The recommended §Tracking patch from β's close-out (broad open-PR list + client-side scan on any new-branch transition event) is the right MCA.

- **Branch-naming non-canonicality (recurring from 3.59.0).** `claude/alpha-tier-3-skills-M8Vce` — CDD §4.2 canonical format is `{agent}/{issue}-{scope}`. Issue number absent from branch name. Same harness pattern. Both β's and γ's Monitor globs missed the branch because they filtered on issue number. Root cause is upstream (harness ↔ CDD §4.2 contract); the per-cycle workaround is the broad-filter + client-side scan recommended in the §Tracking patch.

- **Issue body stale doc-path.** #235 body cited `docs/alpha/INVARIANTS.md`; canonical is `docs/alpha/architecture/INVARIANTS.md`. Did not affect implementation. γ-side issue-quality observation.

**What went right:**

- **Zero-finding diff.** α applied 3.59.0's lessons directly: canonical parser reused (§2.17 generalization), pre-implementation peer enumeration across three consuming surfaces, AC4 deliberate-break dry-run, surface-enumeration test planning. This is the inheritance pattern working — a named failure mode from one cycle prevented the same class of failure in the next.

- **β caught its own process failure and proposed the MCA.** Rather than noting the filter bug and moving on, β's close-out contains a specific §Tracking patch recommendation with exact query forms. The fix addresses both the slash-in-branch-name issue and the invalid `closes:`/`refs:` search qualifiers.

- **β-step-8 tag-push deferral clean again.** Same pattern as 3.59.0. Release artifacts on main, tag deferred to δ, δ pushed, release pipeline fired. Protocol works.

- **α's four reusable patterns.** Pre-implementation peer enumeration, AC refinement via live-package check, copy+inject+run dry-run, and diff-local §2.17 generalization — all applied successfully this cycle. The close-out records them explicitly for future inheritance.

### 4. CDD Self-Coherence

**Per-finding triage:**

β's four notes (N1–N4) deferred to γ:

| # | Surface | Disposition |
|---|---------|-------------|
| N1 | `pkgbuild.PackageManifest` pre-existing parallel parser | Known debt. Eliminating it ripples through `DiscoveredPackage.Manifest`'s static type. Not this cycle's scope. Remains as documented tech debt for a future `pkgbuild` cleanup cycle. |
| N2 | `BuildOne` does not gate on `CheckOne` | Pre-existing. CI runs `--check` separately so the gate exists in CI; locally opt-in. Architecture decision — `build` and `check` are separate concerns. No action needed. |
| N3 | `os.Stat` follows symlinks | Pre-existing. Symlink entrypoint to a regular file outside pkgDir would pass. Low risk in current usage (packages authored in-repo, not installed from untrusted sources). Can be addressed when/if package sources become untrusted. No action now. |
| N4 | Issue body stale doc-path | γ-side quality gap. Immediate fix not available (issue is closed). Note for γ to check doc-path freshness at dispatch time in future cycles. |

### 4b. Cycle Iteration

- **Triggered by:** avoidable tooling/environmental failure (polling discipline gap); loaded skill failed to prevent a finding (CDD §Tracking was loaded but did not prevent the polling miss)
- **Root cause:** CDD §Tracking names polling as primary over webhooks, but the rule is stated for session-start baseline only. The per-event handling rule (broad sync on transition events) is not stated. Both α and β trusted narrow filters that silently returned empty.
- **Disposition:** patch §Tracking now — extend the synchronous-baseline rule from session-start to per-event handling; add the broad-filter + client-side scan form β's close-out recommends; note that `head=owner:branch` is unreliable for branches with slashes and that `closes:`/`refs:`/`fixes:` are not valid GitHub search qualifiers.
- **Evidence:** β close-out §Cycle findings recommended patch; α close-out §Cycle Iteration friction log items 1–3 + root cause; same pattern across 3.59.0 and 3.60.0 cycles (N=5+ independent occurrences across all three roles).

### 5. Release Verification

- Tag: `3.60.0` on origin at `cce905e`
- `release.yml`: completed green (4-platform matrix) at 2026-04-26T23:11:33Z
- GitHub release: published with `cn-{linux-x64,linux-arm64,macos-x64,macos-arm64}`, `checksums.txt`, 5 tarballs, `index.json`
- `release-smoke.yml`: last run 2026-04-26T12:10:42Z (pre-tag; this may be from 3.59.3). **Note:** no `release-smoke.yml` run appears for the 3.60.0 release event. This should be verified — either the smoke fired and completed before the run-list query, or it did not fire on this release's `published` event.
- CI on main (post-release commits): green as of 2026-04-27T09:17:38Z

### 6. §Tracking Patch (Immediate MCA)

Per β's close-out recommendation and the cycle-iteration disposition above, the following patch is due:

**In CDD.md §Tracking, extend the baseline rule:**

> On any new-branch transition event (not just session-start), do a synchronous broad open-PR list (`gh pr list --state open` or `mcp__github__list_pull_requests state=open` with no `head`/`search` filter) and scan client-side for `(head.ref == new-branch) OR (body matches (?i)\b(closes|fixes|resolves|refs)\s*#N\b)`. The MCP `head=owner:branch` filter is unreliable for in-repo branches whose name contains a slash. The `gh pr list --search` qualifier set does not include `closes:`/`refs:`/`fixes:`; those query strings always return empty.

This extends the existing baseline rule from session-start to per-event handling. Same structural shape, broader application scope.

### 7. Next Move

**MCI freeze continues** (8 growing-lag items). Selection per §3.3 (PRA commitment → highest-leverage stale MCA):

The natural next pick from the lag table is **#262 AC1/AC5** (`cn pack` rename + `--list`/`--dry-run` — finishes the 3.58.0 partial). It is the smallest remaining MCA that closes an open partial.

Alternatively, **#273** (rebase-collision integrity guard, P1) is a process-safety MCA that prevents silent content loss — a different class of risk than #262's feature completion.

**#277** (SOUL.md rewrite) is P1 docs/process, not a code MCA — it can run in parallel with a code cycle if operator capacity permits.

The §Tracking patch from §6 should be landed as immediate output before the next cycle dispatches, so the polling discipline fix is in place before α and β start their next coordination.

**Stale-backlog re-evaluation** (per 3.59.0 PRA §7, CDD §3.4): still overdue. The 14+ items below the active lag table have carried across 3+ cycles. Should be slotted when no §3.3-default MCA forces a different next move.
