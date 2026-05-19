<!--
section-manifest:
  planned: [review-summary, implementation-assessment, technical-review, process-observations, release-notes-for-gamma]
  completed: [review-summary, implementation-assessment, technical-review, process-observations, release-notes-for-gamma]
-->

# β close-out — #379 agent/activate skill

## Review summary

**Verdict:** APPROVED in R1. Zero findings at any severity. Single review round; no fix-round triggered.

**Cycle target round count (β/SKILL.md):** ≤2 rounds. Achieved in 1. No `§9.1 review-churn` trigger.

**Pre-merge gate (β/SKILL.md):**
- Row 1 (identity truth): worktree-local config inherited `alpha@cdd.cnos` at session start; re-asserted `beta@cdd.cnos` at worktree scope before any β-authored commit. Audit trail clean — `git log --format='%ae' origin/main..HEAD` on the cycle branch shows `gamma@cdd.cnos` (scaffold), `alpha@cdd.cnos` (eleven α commits), `beta@cdd.cnos` (R1 review verdict and the merge commit).
- Row 2 (canonical-skill freshness): `origin/main` SHA at session start `7a7f7152af649ea93cebe5f909fbf1397d809547` was unchanged at review-time. No skill re-load triggered.
- Row 3 (non-destructive merge-test): clean automatic merge in throwaway worktree at `/tmp/cnos-merge379/wt`; zero unmerged paths; `go test ./internal/activate/...` PASS on the merge tree. Worktree torn down; parent-repo identity preserved at `beta@cdd.cnos`.
- Row 4 (γ artifact completeness): `.cdd/unreleased/379/gamma-scaffold.md` present on cycle branch. Rule 3.11b satisfied; no exemption needed.
- Row N (test execution on cycle HEAD): `go test -count=1 ./internal/activate/...` PASS (29 tests, 0.042s); `go vet` clean.

**Merge:** `a3bf7892` on `main`, merge commit body carries `Closes #379` and a one-line summary of what shipped. No-ff merge so the cycle branch topology is preserved for γ's PRA.

## Implementation assessment

α's implementation maps issue-body ACs to concrete artifacts with no ambiguity left for β:

- **AC1–AC6 (skill artifact).** SKILL.md is well-formed, conforms to `cnos.core/skills/skill/SKILL.md` literally (frontmatter + Define/Unfold/Rules/Verify + embedded kata), and demonstrates its own formula (the skill is itself a coherent activation procedure that an AI body could follow). The §2.1 (human) and §4.1 (machine-readable) load orders are peers and match; the §2.2 capability matrix names three tiers with explicit load paths and marks tier (a) preferred with reasons; the §2.3 router template is paste-testable with a single `<HUB-URL>` placeholder; the §2.4 disambiguation cites both `cdd/activation` and this skill with exact paths.

- **AC7 (renderer evolution).** This is the load-bearing AC for the cycle's "everything is a skill" claim, and α landed the strong form. `writePrompt` no longer carries the section ordering as in-Go literals; it iterates a `readFirst []readFirstItem` slice produced by `loadActivateSkillOrdering`, which reads the skill's §4.1 marker-bounded block via `parseReadFirstOrderBlock`. The test `TestSkillIsSourceOfTruthForReadFirstOrder` is two-phase (vendor a fixture skill ordering kernel-before-persona, assert; then swap, assert the swap propagates) and includes a coherence assertion (`out1 != out2`). The fallback path is preserved and tested (`TestSkillFallback_NotVendored`), and the parser is unit-tested both happy-path and no-markers. The interpolation surface is explicit in §4.2 of the skill (token → hub-state-field table) and one-to-one with `renderReadFirstLine`'s switch cases.

- **Observable-output delta.** The `## Read first` ordering changed from persona/operator/kernel to kernel/CA-skills/persona/operator/hub-state/identity. The change is the intended consequence of AC3 and is documented in α's self-coherence §Debt 2. Existing tests (`TestReadFirstSection_OrderedSigma`, `TestReadFirstSection_InitOnlyOrdered`) were updated to the new canonical ordering — appropriate, not silent.

- **Surface containment.** Exactly five files in the diff (gamma-scaffold, self-coherence, activate.go, activate_test.go, SKILL.md), all in the allowed set. No edits to `cdd/activation/SKILL.md`, no hub README patches, no `cnos.xyz`, no `cn-doctor`, no `cmd_activate.go`. γ's failure mode 10 held.

## Technical review

**Strong test design.** The AC7 source-of-truth test does what the issue body literally asks (edit-to-output coherence) rather than the weaker file-exists oracle γ pre-flagged as failure mode 6. The two-phase pattern (write fixture, assert; rewrite fixture, assert + compare outputs) is also the right shape for catching future regressions where someone accidentally re-introduces a hardcoded ordering — the test would fail not because of the new code but because the swap wouldn't propagate.

**Parser hygiene.** `parseReadFirstOrderBlock` is line-oriented, accepts both em-dash (U+2014) and ASCII `--` separators, requires a digit prefix to skip prose lines, and returns `(nil, false)` on missing markers — letting the caller pick a fallback rather than coupling parsing failure to renderer panic. Token vocabulary is reserved-forward-compatible: unknown tokens emit `(unknown token: X)` and continue, which is the right call for a renderer that must stay live so the rest of the prompt reaches the body.

**Fallback duplication is honest debt.** `canonicalReadFirstOrdering` mirrors §4.1 as in-Go constants. α flagged this in self-coherence §Debt 3 and named the durable fix (`//go:embed` the skill into the binary in a future cycle). Today's compromise — skill canonical on drift, fallback only on missing vendor, fallback path tested — is acceptable for this cycle's scope.

**Identity-truth observation reused.** The worktree-local config leak β had to recover from at session start is the same root failure mode β/SKILL.md row 1 documents from cycle #301; α encountered an identical leak (worktree carried γ's identity from the scaffold session) and used the α/SKILL.md §2.6 row 14 path (a) rebase-with-reset-author. Both recoveries succeeded; the audit trail is clean.

## Process observations (for γ PRA)

1. **Round count = 1; design-and-build cycle landed in a single review pass.** This is the lower end of the β/SKILL.md target band (≤2 nominal, 2 acceptable, 3 triggers §9.1). γ's scaffold (peer enumeration + ten pre-flagged failure modes + AC posture) is load-bearing for the single-round outcome — α had no ambiguity to resolve mid-cycle, and the failure-mode catalogue gave β a concrete checklist to run rather than open-ended skepticism. Worth surfacing as a positive process observation: γ scaffolding paid off here in measurable round-count.

2. **AC7 was the highest-leverage AC and α treated it that way.** The cycle's "everything is a skill" claim hinges on the source-of-truth test being strong. α wrote it strong (two-phase edit, coherence assertion) on the first pass — no β nudge required. This is the kind of γ-flagged failure mode (file-read-oracle masquerading as source-of-truth) that bites cycles where the test only proves the file is readable; α dodged it by writing the test the issue body literally asked for, not a weaker version that would also pass.

3. **Observable-output delta surfaced honestly.** α named the `## Read first` ordering delta in self-coherence §Debt 2, listed the consumers affected (string-matchers within the section) vs unaffected (section-header parsers), and updated the existing tests to the new canonical order. This is the right shape for an intended-but-observable change — surface it as debt to track, not as a hidden side-effect.

4. **Cross-repo proposal faithfulness.** The issue body's §Source Proposal cites `usurobor/cn-sigma:.cdd/iterations/cross-repo/cnos/agent-activate-skill/` bundle commit `1a4e25f75` as carrying the AC4 capability-matrix refinement. β did not fetch the source bundle (per brief: issue body is authoritative). α's SKILL.md §2.2 matches the issue body's AC4 specification — three tiers with explicit load paths and tier (a) preferred — so the citation is faithful at the level of what landed.

5. **No `gamma-clarification.md` mid-cycle.** Issue body was binding throughout; no cache-bust required. Indicator of well-scoped initial issue.

## Release notes for γ / δ

**Cycle complete from β's side.** Merge commit `a3bf7892` on `main`. Cycle branch `cycle/379` not deleted (δ's call at the release boundary). `.cdd/unreleased/379/` not moved (δ's call at the release boundary per `release/SKILL.md` §2.5a). Version not bumped, no tag created, CHANGELOG not updated for release (all δ).

**γ inputs for the PRA:**
- All seven ACs PASS; zero findings of any severity; single review round.
- α surfaced five debt items in self-coherence §Debt; β agreed with α's framing on all five and added no debt items of its own.
- Process observations 1–5 above name what γ's scaffolding contributed to the smooth outcome.
- Cross-repo source-proposal disposition `accepted` per issue body §Source Proposal; γ may want to confirm the cn-sigma side acknowledges acceptance and update STATUS on the source bundle.

**δ inputs for the release boundary:**
- Observable-output delta in `cn activate`'s `## Read first` section (intentional per AC3; documented in α's self-coherence §Debt 2). Consumers parsing by section header are unaffected; string-matchers within the section will see different content. Worth a release-note mention.
- Renderer fallback `canonicalReadFirstOrdering` duplicates the skill's §4.1 ordering as in-Go constants. Mitigations in place; future-cycle direction is `//go:embed` per α's §Debt 3.
- Hub README router template ships in the skill (§2.3) but no hub has adopted it yet; per-hub README patches are deferred to per-hub cycles per the issue body §Scope.

β's last action is this close-out file + the merge commit + the cycle-branch push. γ takes the PRA; δ takes the release boundary.
