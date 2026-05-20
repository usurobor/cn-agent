# Self-Coherence — cycle/378

## §Gap

**Issue:** [cnos#378](https://github.com/usurobor/cnos/issues/378) — "cdd: rule 3.11b discoverability under §5.2 wave-mode (wave-manifest as γ-artifact-of-record; cdd-protocol-gap)".

**Mode:** §5.2 wave-mode (wave `2026-05-19-protocol-hygiene`). γ=δ collapse permitted per wave manifest. **γ+α+β-collapsed-on-δ** for this single skill-patch cycle — `α ≠ β within a session is structurally compromised but acceptable for skill-patch class per wave manifest precedent` (acknowledged in dispatch). The β-collapsed-on-δ self-review at §beta-review.md applies AC1–AC4 mechanically; this is the strongest available β discipline under collapse.

**Gap (paraphrased):** `review/SKILL.md` rule 3.11b "Exemption discoverability" — silent on §5.2 wave-mode. Strict-literal reading would RC every sub of a wave-internal §5.2 configuration where the wave manifest at `.cdd/waves/{wave-id}/manifest.md` carries the γ-artifact-of-record duty in lieu of per-sub `.cdd/unreleased/{N}/gamma-scaffold.md`. The cph cdr-refactor wave 2026-05-18 (master `cph#11`; subs `cph#12, #13, #14, #15`) demonstrated this end-to-end: zero per-sub scaffolds; wave manifest carried γ duty; β produced three distinct substantive-read justifications across four subs — not a β skill issue, a discoverability gap.

Parallel α-side gap: `alpha/SKILL.md` §2.6 pre-review gate had no row for γ-side artifact presence at all, so α could not pre-empt the rule 3.11b check.

## §Skills

**Tier 1 (always-loaded for δ-as-agent under γ+α+β collapse):**
- `cdd/CDD.md` (canonical lifecycle and role contract; γ=δ collapse rule context)
- `cdd/alpha/SKILL.md` (α role surface — patch target)
- `cdd/beta/SKILL.md` (β role surface — applied for self-review under collapse)
- `cdd/gamma/SKILL.md` (γ role surface — §2.2a empirical anchor precedent)
- `cdd/operator/SKILL.md` (δ role surface — patch target; §5.2 wave-mode definition)
- `cdd/review/SKILL.md` (β review skill — rule 3.11b primary patch surface)

**Tier 2:** None — this is a skill-text patch, no engineering toolchain involved.

**Tier 3 (per dispatch):**
- `cdd/review/SKILL.md` rule 3.11b (primary patch surface)
- `cdd/alpha/SKILL.md` §2.6 (pre-review gate patch surface)
- `cdd/operator/SKILL.md` §5.2 + §10 (wave-manifest cross-reference patch surface)
- `cdd/CDD.md` §1.4 (γ=δ collapse context; read-only)

**Wave manifest:** `.cdd/waves/2026-05-19-protocol-hygiene/manifest.md` (loaded as γ-artifact-of-record per §5.2 — this is exactly the configuration the patch codifies).

## §ACs

### AC1: Wave-manifest recognized as γ-artifact-of-record under §5.2

**Invariant (issue body):** Rule 3.11b's exemption-discoverability clause names the wave-manifest path as a valid γ-artifact-of-record for any sub of that wave under §5.2 mode.

**Oracle:** `review/SKILL.md` rule 3.11b text mentions `.cdd/waves/{wave-id}/manifest.md` AND `§5.2 wave-mode` (or equivalent §5.2-mode phrasing). The discoverability requirement names what makes the link from sub-issue → wave-manifest auditable.

**Evidence:** `review/SKILL.md` L141 (new clause (ii) under §3.11b Exemption discoverability) — names `.cdd/waves/{wave-id}/manifest.md`, names `§5.2 wave-mode`, names both discoverability link paths: (a) sub-issue body cites wave-id, (b) master tracking issue links to sub. L142 adds recovery path (c) for §5.2 wave-mode. L143 extends the `beta-review.md §Artifact completeness` documentation requirement.

**Status:** PASS.

### AC2: α §2.6 row for γ-side artifact presence

**Invariant:** `alpha/SKILL.md` §2.6 pre-review gate enumerates the γ-side artifact whose presence rule 3.11b examines — canonical §5.1 path OR wave-manifest under §5.2.

**Oracle:** A row appears in α §2.6 naming the γ-side artifact and rule 3.11b cross-reference. Outcome describable in `self-coherence.md §Review-readiness`.

**Evidence:** `alpha/SKILL.md` row 15 (new) under §2.6 — enumerates §5.1 canonical / §5.2 wave-mode / absent configurations; cross-references `review/SKILL.md` §3.11b clause (ii); names the executable check (`git cat-file -e` / `git ls-tree -r`); names the recording format for §Review-readiness (`γ-artifact at canonical §5.1 path` / `wave-manifest serves under §5.2 (wave-id: {wave-id}; discoverability: ...)` / `absent — ...`).

**Status:** PASS.

### AC3: Empirical anchor cited

**Invariant:** The skill texts cite the cph cdr-refactor wave (2026-05-18) as the empirical anchor — four subs uniformly under §5.2, three distinct β substantive reads of the same configuration.

**Oracle:** Cycle / wave reference appears in the patched skill text (per `gamma/SKILL.md` §2.2a precedent).

**Evidence:** All three patched files cite the cph wave:
- `review/SKILL.md` L141 — full citation with per-sub β anchors (`cph#12 beta-review.md §3.11b L133–158`; `cph#13 §Contract Integrity row 11`; `cph#14 §2.5 L187–195`; `cph#15 §2.0.0 row L29`), names the three-distinct-β-read pattern, and cross-references `cph#15 β-closeout L55` (formal recommendation) and `usurobor/cph:.cdd/iterations/wave-2026-05-18.md` Finding F1.
- `alpha/SKILL.md` L259 — cites α-side anchors (`cph#12 alpha-closeout.md L38–40`, `cph#14 alpha-closeout.md L46–51`); names the α-side surface of the gap.
- `operator/SKILL.md` L319 — cites wave + subs + four-of-four pattern + three-distinct-β-read pattern + wave-iteration F1.

**Status:** PASS.

### AC4: No CI / runtime / release surface change

**Invariant:** Skill-patch cycle; no new CI workflow, no validator change, no CDD doctrine edit.

**Oracle:** `git diff origin/main..HEAD --stat` shows changes only under:
- `src/packages/cnos.cdd/skills/cdd/review/SKILL.md`
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md`
- `src/packages/cnos.cdd/skills/cdd/operator/SKILL.md`
- `.cdd/unreleased/378/*.md`

**Evidence:** `git diff --stat origin/main..HEAD` at HEAD `64c96317` (current implementation SHA; rebase chain: pre-R1.1 `561c19af` → post-R1.1 `f56afd2d` → post-R1.2 `64c96317` — see §Fix-round R1.1 and §Fix-round R1.2) shows only the three SKILL files. The γ-scaffold and self-coherence commits add only `.cdd/unreleased/378/*.md`. No CI workflows, no validators, no `cdd/CDD.md` edit.

**Status:** PASS.

## §Self-check

- **Did α push ambiguity onto β?** No. Each AC has direct line-level evidence; the three edits use identical phrasing for the wave-mode discoverability path (`.cdd/waves/{wave-id}/manifest.md`, §5.2 wave-mode, (a)/(b) discoverability paths), so cross-skill drift is preempted by construction.
- **Is every claim backed by evidence in the diff?** Yes — line numbers cited above resolve to the patched text.
- **Cross-skill coherence as α discipline.** The three edits describe the same wave-mode discoverability path: review §3.11b clause (ii) defines it as β's gate; alpha §2.6 row 15 defines it as α's pre-gate check; operator §5.2 defines it as the wave-author's responsibility before sub-dispatch. The chain — wave author → α pre-check → β gate — is internally consistent.
- **Non-goal preservation.** Issue body §Non-goals (1) `gamma-scaffold.md` canonical §5.1 path not relocated — confirmed; (2) rule 3.11b's binding β-side gate preserved for §5.1 — confirmed (clause (i) carries the canonical §5.1 surface); (3) no new validator — confirmed; (4) `## Protocol exemption` not required in every wave-sub body — confirmed (path (c) is canonical for §5.2, path (b) remains escape valve); (5) cnos#375 not subsumed — confirmed (γ-side pre-dispatch gate is the orthogonal axis; this cycle's surfaces are review/alpha/operator, not gamma).
- **Skill-class peer enumeration** per `alpha/SKILL.md` §2.3. Role-skill peers: `review`, `alpha`, `operator` — all three patched. Lifecycle-skill peers that might encode the rule 3.11b contract operationally: `gamma` (γ creates the γ-artifact), `beta` (β reads the gate, but the gate text is in `review/SKILL.md`), `release` / `post-release` (no direct dependency — rule 3.11b is a β-gate, not a release-gate). `gamma/SKILL.md` is patched in parallel cycle #375 for the γ-side pre-dispatch gate — the orthogonal axis. β re-audit confirms no peer drift in `gamma/SKILL.md` from this cycle's changes (the §5.2 wave-mode discoverability path lands at `operator/SKILL.md` §5.2, which γ-as-wave-planner reads at wave-author time per the new clause).
- **Honest-claim verification (rule 3.13).** (a) Reproducibility: the cph wave anchor is sourced from the issue body's `## Source` section, which itself cites `usurobor/cph:.cdd/waves/cdr-refactor-2026-05-18/{manifest.md,wave-closeout.md,receipt.md}` + `usurobor/cph:.cdd/iterations/wave-2026-05-18.md` + per-sub `cph:.cdd/unreleased/{12,13,14,15}/beta-review.md` and `alpha-closeout.md` files. These are cited authoritatively in the issue and reproduced verbatim in the patched skill text; no novel claims added. (b) Source-of-truth alignment: `γ-artifact-of-record`, `§5.2 wave-mode`, `§5.1 canonical dispatch`, `.cdd/waves/{wave-id}/manifest.md` all match operator/SKILL.md §5.2 and §10 canonical definitions. (c) Wiring claims: the cross-skill cross-references (review→alpha, alpha→review, operator→review, review→operator) resolve — each names the partner skill and the section in the partner skill where the contract lives.

## §Debt

None outright. Two acknowledged constraints (not debt):

1. **α=β=γ collapse within session.** The wave manifest permits γ=δ collapse for skill-patch class; this cycle further collapses γ+α+β-onto-δ. `α ≠ β within a session is structurally compromised` per `operator/SKILL.md` §5.2 scope clause. Accepted because (a) wave-manifest precedent, (b) cycle scope bounded (three additive skill-text edits, no logic change), (c) the β-collapsed-on-δ self-review (`beta-review.md`) applies the AC1–AC4 oracle mechanically and surfaces cross-skill coherence as the binding gate. The β-collapse is named in `beta-review.md`.

2. **§5.2 wave-mode itself is the cycle's mode.** This cycle runs under exactly the configuration it patches — the wave manifest at `.cdd/waves/2026-05-19-protocol-hygiene/manifest.md` carries the γ-artifact-of-record duty for cycle #378. The sub-issue body (cnos#378) does not currently cite the wave-id — the wave-manifest links to #378 via the master `## Issues` table (path (b) per the new rule 3.11b clause (ii)). The dispatch prompt itself names the wave-manifest path explicitly, which makes the discoverability link auditable at the dispatch surface. This cycle is therefore self-validating: the patch lands a rule that this cycle is itself the first instance of post-cph satisfying.

## §CDD Trace

| Step | Name | Phase | Producer | Evidence |
|---|---|---|---|---|
| 0 | Observe | γ | δ-as-γ | wave manifest read; issue #378 read via `mcp__github__issue_read` |
| 1 | Select | γ | δ-as-γ | issue #378 selected per wave manifest queue |
| 2 | Branch | γ | δ-as-γ | `cycle/378` created from `origin/main` SHA `c9017153`; pushed to `origin/cycle/378`; rebased onto `origin/main:dd5a36d9` (R1.1); rebased again onto `origin/main:8e118320` after parallel cycle #375 merged (R1.2) |
| 3 | Bootstrap | α | δ-as-α | git identity `alpha@cdd.cnos` set; `.cdd/unreleased/378/` created |
| 3a | γ-scaffold | γ | δ-as-γ | `.cdd/unreleased/378/gamma-scaffold.md` at `f5ab0e35` (current; chain `66d02c8a→26303a3d→f5ab0e35`) — γ-side artifact present per `gamma/SKILL.md` §2.5 step 3a + #375 preventive — even though wave-manifest also serves under §5.2, the per-cycle scaffold is authored as belt-and-suspenders per wave invariant #3 |
| 4 | Gap + ACs | α | δ-as-α | §Gap above; AC1–AC4 mapped above |
| 5 | Mode | α | δ-as-α | §5.2 wave-mode under γ+α+β-on-δ collapse (named §Gap, §Debt) |
| 6 | Artifacts | α | δ-as-α | three SKILL files patched at `64c96317` (current; chain `561c19af→f56afd2d→64c96317`): `review/SKILL.md` L141–143, `alpha/SKILL.md` L255–259 (row 15), `operator/SKILL.md` L319 |
| 7 | Self-coherence | α | δ-as-α | this file |

## §Review-readiness

**Round:** 1
**Implementation SHA:** `64c96317` (current; rebase chain `561c19af→f56afd2d→64c96317` — see §Fix-round R1.1 and §Fix-round R1.2)
**Base SHA at observation:** `8e118320` (origin/main current; advance chain `c9017153→dd5a36d9→8e118320`; cycle/378 rebased twice — R1.1 onto `dd5a36d9`, R1.2 onto `8e118320` after parallel cycle #375 merged)
**Branch CI:** N/A (skill-patch, docs-only disconnect class per wave manifest §wave-level invariants 1; no CI workflows triggered by these paths)

**§2.6 row 15 outcome (the row this cycle adds):**
- `γ-artifact at canonical §5.1 path` — `.cdd/unreleased/378/gamma-scaffold.md` present at `origin/cycle/378:f5ab0e35` (current; chain `66d02c8a→26303a3d→f5ab0e35`); verifiable via `git ls-tree -r origin/cycle/378 .cdd/unreleased/378/gamma-scaffold.md`
- AND `wave-manifest serves under §5.2 (wave-id: 2026-05-19-protocol-hygiene; discoverability: master-tracking-link)` — the wave manifest at `.cdd/waves/2026-05-19-protocol-hygiene/manifest.md` exists on `origin/main` and lists #378 in the `## Issues` table (path (b) per the new clause).
- Both halves hold; rule 3.11b satisfied via either configuration.

Ready for β (collapsed-on-δ).

## §Fix-round R1.1 — mid-cycle rebase (α-internal, pre-β-handoff)

Not a β-RC fix-round — α-internal pre-handoff state update per `alpha/SKILL.md` §2.6 row 1 (transient row: cycle branch rebased onto current `origin/main`). Surfaced by α's `git diff --stat origin/main..HEAD` check between implementation commit and β handoff.

**Trigger.** Between branch creation (`c9017153`, T0) and β handoff prep (T1), `origin/main` advanced by one commit: `dd5a36d9 chore: ignore .claude/{settings.local.json,worktrees/}`. Surfaced by `git diff origin/main..HEAD --stat` showing an unexpected `-4 lines` entry for `.gitignore` that no commit on `cycle/378` had authored.

**Action.** `git rebase origin/main` (clean, no conflicts); `git push --force-with-lease origin cycle/378` (lease-checked because cycle/378 is the only branch carrying the work).

**SHA rewrite map** (pre-rebase → post-rebase):
- γ-scaffold commit: `66d02c8a → 26303a3d`
- α implementation commit: `561c19af → f56afd2d`
- α self-coherence commit (this file's first revision): `453144aa → 7e64a673`

**Re-stamping (per `alpha/SKILL.md` §2.6 SHA-citations-across-rebase paragraph, path (ii) reactive resolution applying §2.3 intra-doc rule).** Grep target: `c9017153|66d02c8a|561c19af|453144aa` across `.cdd/unreleased/378/*.md`. 5 sites found and re-stamped:
- `gamma-scaffold.md:9` — branch creation SHA preserved (it is the historical fact of branch creation, not a current pointer); appended cross-reference to this §Fix-round R1.1 section so readers know about the rebase.
- `self-coherence.md:78` — AC4 evidence: `561c19af → f56afd2d` (post-rebase); pre-rebase value preserved as cross-reference.
- `self-coherence.md:105` — CDD Trace step 2: branch-creation SHA `c9017153` preserved as historical fact; rebase target `dd5a36d9` appended.
- `self-coherence.md:107` — CDD Trace step 3a: γ-scaffold SHA `66d02c8a → 26303a3d`; pre-rebase preserved.
- `self-coherence.md:110` — CDD Trace step 6: implementation SHA `561c19af → f56afd2d`; pre-rebase preserved.
- `self-coherence.md:116` — §Review-readiness Implementation SHA: `561c19af → f56afd2d`; pre-rebase preserved.
- `self-coherence.md:117` — §Review-readiness Base SHA: `c9017153 → dd5a36d9`.
- `self-coherence.md:121` — §Review-readiness row 15 outcome: γ-scaffold SHA `66d02c8a → 26303a3d`; pre-rebase preserved.

**Post-restamp grep verification** (run before β handoff): `grep -n 'c9017153\|66d02c8a\|561c19af\|453144aa' .cdd/unreleased/378/*.md` returns only sites where the pre-rebase SHA is preserved as historical cross-reference (each accompanied by "post-rebase" / "pre-rebase" annotation per §2.3 intra-doc-repetition discipline) — no bare claims to invalidated SHAs.

**Post-restamp implementation SHA** for §Review-readiness: this fix-round section is itself committed *after* re-stamping; the post-rebase implementation SHA `f56afd2d` (the SKILL-file edit commit) remains the stable "implementation SHA" per `alpha/SKILL.md` §2.6 SHA convention (last implementation commit before the readiness-signal commit). The self-coherence commit that lands this §Fix-round R1.1 section is the readiness-signal commit; β reads HEAD via polling.

**β handoff.** Re-signaled. Branch HEAD post-restamp: see the commit that lands this fix-round section.

## §Fix-round R1.2 — second mid-cycle rebase (α-internal, parallel cycle #375 merge)

Not a β-RC fix-round — α-internal pre-handoff state update per `alpha/SKILL.md` §2.6 row 1.

**Trigger.** Between R1.1 commit and β handoff prep, `origin/main` advanced again — parallel cycle #375 (the orthogonal axis of this wave per wave manifest) closed and merged: commits `7ef98eb4` (γ-375), `c4d29344` (α-375 patching `gamma/SKILL.md` §2.5 Step 3b), `feebd45c` (β-375 APPROVED + close-outs), `8e118320` (merge). Surfaced by `git diff --stat origin/main..HEAD` showing unrelated `.cdd/unreleased/375/` and `gamma/SKILL.md` entries.

**Action.** `git rebase origin/main` (clean; no file-level conflict — #375 touched `gamma/SKILL.md` §2.5 Step 3b only, which this cycle does not touch; this cycle's `alpha/SKILL.md` and `operator/SKILL.md` changes do not overlap with #375); `git push --force-with-lease origin cycle/378`.

**Cross-cycle coordination note.** Wave manifest §wave-level invariants 4 anticipates this case: "If two cycles touch the same file, β raises a cross-cycle finding at review and γ coordinates merge order." No file-level overlap occurred (cycle #375 patched `gamma/SKILL.md`; cycle #378 patched `review/SKILL.md`, `alpha/SKILL.md`, `operator/SKILL.md`), so this is the no-conflict path — α rebases and the merge order is mechanical (#375 first, #378 second).

**SHA rewrite map** (R1.1 SHAs → R1.2 SHAs):
- γ-scaffold: `26303a3d → f5ab0e35`
- α impl:    `f56afd2d → 64c96317`
- α self-coh:`7e64a673 → 23903ab4`
- α R1.1:    `80b51d59 → fe22da53`

**Re-stamping.** Live citation sites updated (sites that point at *current* state — implementation SHA, base SHA, γ-scaffold presence, AC4 evidence, CDD Trace rows). The §Fix-round R1.1 narrative retains its R1.1-era SHAs verbatim because that section describes the R1.1 event as it occurred (`66d02c8a→26303a3d`, etc. are accurate as descriptions of *what happened during R1.1*, even though those R1.1 post-rebase SHAs no longer exist on the branch after R1.2). This preserves the audit narrative without overloading R1.1's voice with R1.2's SHAs.

**Post-restamp grep verification** (run before β handoff): `grep -n 'c9017153\|66d02c8a\|561c19af\|453144aa\|dd5a36d9\|26303a3d\|f56afd2d\|7e64a673\|80b51d59' .cdd/unreleased/378/*.md` returns only sites accompanied by chain-cross-reference annotations (e.g. "chain `66d02c8a→26303a3d→f5ab0e35`") or sites inside §Fix-round R1.1's R1.1-event description — no bare claims to invalidated SHAs at the live-state surface.

**Cycle resilience observation.** Two mid-cycle rebases (R1.1 single-commit advance, R1.2 four-commit advance from parallel cycle) is within the wave's expected behavior — three parallel skill-patch cycles dispatched concurrently per the wave manifest means rebase churn is a structural feature of §5.2 wave-mode, not a friction signal. The §2.6 row 1 + §2.3 intra-doc + §2.6 SHA-citations-across-rebase chain handles this cleanly; no protocol gap surfaced.

**β handoff.** Re-signaled. Branch HEAD post-R1.2: see the commit that lands this fix-round section.

