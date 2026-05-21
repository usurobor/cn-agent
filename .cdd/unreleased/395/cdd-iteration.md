# cdd-iteration — Cycle #395 (Sub 3 of #376)

**Cycle:** #395 — Sub 3 (#376): CDR role overlays (α/β/γ/δ/ε SKILL.md files for cnos.cdr)
**Merge:** `f891f08b`
**Closed:** 2026-05-21 (auto-closed via `Closes #395` in merge commit)
**Mode:** design-and-build (γ+α+β-collapsed-on-δ per breadth-2026-05-12 wave manifest precedent + cycles 375/377/378/388/390 empirical validation; engineering-class β-α-collapse acknowledged in `beta-review.md`)
**Rounds:** R1 APPROVED (no fix-round)
**ACs:** 6/6 PASS (with classification tables per oracle-vs-intent on AC3, AC4, AC6)

## §1 Findings dispositioned

### F1: AC oracle wording vs. intent — `rg`-only oracles overcount disavowal context

- **Source:** cycle 395 α-closeout F1 + β-review Obs-1 + design-notes §"Design source citations" (the disavowal idiom inherited from cnos.cdd).
- **Class:** `cdd-protocol-gap` (oracle-template friction; second-occurrence in the c-d-X-family bootstrap pattern).
- **Trigger:** ε process-gap check; surfaced by AC3 (6 release|deploy|tag hits, all disavowal/cross-reference) + AC4 (2 persona-pattern hits in CDR.md disavowal context) + AC6 ("dispatch" used as abstract role-function concept across role overlays).
- **Description:** The cnos#395 issue body's AC3, AC4, AC6 oracles use bare `rg <pattern>` count checks. Protocol-overlay skill files structurally inherit the parent-doctrine's "what the role does NOT do" disavowal idiom (cnos.cdd `operator/SKILL.md §6` is the canonical pattern). This means the literal `rg` count is non-zero, but every hit lands in disavowal context. The issue-prose carve-outs (AC3 explicit: "matches only in cross-references or 'what CDR is **not**' disclaimers — not in normative field declarations") acknowledge the gap; the oracle itself does not encode it.

  Cycle #390 (Sub 1) surfaced the same pattern for `release|deploy|tag` in CDR.md (6 hits, accepted via the cycle #390 cdd-iteration §2 classification: "all 6 hits classified as disavowing-context or label-tag (non-release-tag) usage"). Cycle #395 surfaces it for the same plus `I am Rho|my voice` (AC4) and `dispatch|polling` (AC6).
- **Root cause:** Protocol-overlay authoring inherently requires citing the parent doctrine's surfaces while declaring which ones do not transfer to the new discipline. The disavowal idiom is structural; a forbidden-token oracle that does not carve out disavowal context will overcount.
- **Disposition:** `next-MCA`. Patch shape: the cnos.cdd `issue/proof/SKILL.md` (AC oracle authoring skill) gains a "disavowal-context carve-out" template — when an AC oracle uses `rg <pattern>` to assert absence, the oracle includes a structural exception for hits in disavowal subsections (canonically titled "What the role does NOT do" or "X analogues that do not transfer"). The patch lands when a third occurrence triggers it; cnos.cdw or cnos.cda's first bootstrap is the likely third surface.
- **Issue filed:** none yet; suggest filing alongside cnos.cdw or cnos.cda bootstrap.
- **First AC for the eventual MCA:** `cnos.cdd/skills/cdd/issue/proof/SKILL.md` documents the disavowal-context exception as a structural component of forbidden-token AC oracles, with cycle #390 + #395 + (third occurrence) as empirical anchor.

### F2: Sub 2 loader-file integration mechanical conflict (resolved-in-merge)

- **Source:** cycle 395 α-closeout F2 + β-review Obs-2.
- **Class:** `cdr-protocol-integration` (cross-Sub merge surface; resolved at merge time).
- **Trigger:** ε process-gap check; surfaced by Sub 2 (#394) shipping its loader at `cdr/SKILL.md` with `delta/SKILL.md` as the δ-role path while δ pinned Sub 3's δ overlay at `operator/SKILL.md`.
- **Description:** Sub 2 shipped before Sub 3. Sub 2's loader listed five role-overlay sub-skills with the δ overlay at `delta/SKILL.md` (using role-letter naming). Sub 3's δ-pinned implementation contract specified the δ overlay at `operator/SKILL.md` (mirroring `cnos.cdd/skills/cdd/operator/SKILL.md`'s exemplar layout). The rebase-time loader integration kept Sub 2's richer loader content (`triggers`, `inputs`, `outputs`, "Cross-protocol relationship" section, "Conflict rule" section) and corrected the δ-overlay path to `operator/SKILL.md`. The corrected path is now consistent with: (i) the cnos.cdd exemplar, (ii) `CDR.md`'s cross-references, (iii) Sub 3's actual role-overlay file location.
- **Root cause:** The cnos#376 master and the per-Sub issue bodies did not explicitly pin a single δ-overlay directory name across Subs. Sub 2 chose role-letter (`delta/`); Sub 3 chose engineering-exemplar (`operator/`). With Sub 2 + Sub 3 dispatched in parallel, the naming convention surfaced as an integration cost paid at the cycle/395 → main merge.
- **Disposition:** `landed-in-merge`. Merge commit `f891f08b` ships the corrected loader path. The pattern is recognized for future c-d-X bootstrap waves: when sub-issues touch a shared file (here `cdr/SKILL.md`), the master should pin the file's surface conventions (directory naming for role overlays, in this case) before Sub dispatch. Recorded as cross-Sub coordination pattern.
- **Issue filed:** none — landed.
- **First AC for the prevention pattern:** not strictly an MCA; the prevention is a wave-coordination practice (the master issue for a multi-Sub bootstrap pins shared-surface conventions).

## §2 No-findings observations (informational)

- **AC3 release/deploy/tag classification revalidates #390's pattern.** 6 hits in operator/SKILL.md, all in disavowal context or cross-reference to cnos.cdd surfaces declared not-to-transfer. Same hit shape as #390 CDR.md's 6 hits. The cycle #390 ε iteration §2 classification ("all 6 hits classified as disavowing-context or label-tag (non-release-tag) usage") extends to cycle #395 unchanged.
- **AC4 persona-pattern hits in CDR.md (not in Sub 3's surface) are pre-existing disavowal.** The 2 hits the broad `rg src/packages/cnos.cdr/skills/cdr/` returns are in CDR.md (Sub 1's deliverable) where CDR.md *literally states* it does not author those patterns. Sub 3's authored surface (5 role overlays + loader) has 0 hits. The oracle-wording-vs-intent gap is the same as F1; the AC is satisfied on Sub 3's authored surface.
- **AC6 "dispatch" abstract-role-function usage is mandated by CDR.md.** CDR.md line 160 explicitly states: "operational mechanics (dispatch, polling, repo wiring) belong in role-overlay skills (Sub 3) and in persona/operator contracts." Sub 3's role overlays discuss the abstract concept of "dispatch" (γ→α/β handoff) as a doctrinal necessity. The AC6 oracle as literally written would flag any role file mentioning the abstract concept; the AC6 intent ("no runtime mechanics authored") is satisfied by 0 `claude -p`, 0 `gh issue`, 0 `cn dispatch`, 0 polling-loop authoring in normative position. Same class as F1; recorded; not a finding.
- **Identity rotation across rebased commits.** The rebase reset commit-author email to whichever `git config` was active at rebase time, scrambling the α/β/γ identity rotation that this cycle had attempted. The role-identity-is-git-observable property is satisfied in the artifact contents (each file is correctly attributed in its commit message); the git-log identity column is partially scrambled. Engineering-class single-session collapse cycle; no research claim transmitted; α §2.6 row 14 path (b) accepted with known-debt declaration. Recorded; not raised to a finding because the matter shipped is informational documentation, not research-claim transmission.
- **Sub 4 (#396) shipped before this Sub.** Sub 4's empirical-anchor doc at `src/packages/cnos.cdr/docs/empirical-anchor-cph.md` is now on origin/main; Sub 3's role-overlay files coexist cleanly. cnos#376 is now structurally complete (all four Subs shipped).

## §3 Trigger assessment (per `gamma/SKILL.md §2.8` table)

| Trigger | Fire condition | Fired? | ε note |
|---|---|---|---|
| Review churn | review rounds > 2 | **No** | R1 APPROVED on first pass. |
| Mechanical overload | mechanical ratio > 20% AND findings ≥ 10 | **No** | 4 α findings + 4 β observations + 1 γ-side observation = 9 total; below the ≥10 threshold; mechanical ratio is high (oracle-template friction is mechanical) but absolute count does not fire. |
| Avoidable tooling / environment failure | environment blocked the cycle | **Partial — non-protocol** | Signing-server rejected merges from `/tmp` paths and from fresh clones; resolved by merging from a named local branch within the cycle/395 worktree. Recorded; pattern recognized; no protocol patch needed. |
| Loaded-skill miss | a loaded skill should have prevented a finding | **No** | F1's oracle-template friction is a second-occurrence forward-looking next-MCA, not a retroactive skill miss; the existing skills do not yet carve out disavowal context, and that's the patch shape, not a miss. F2's cross-Sub integration is a wave-coordination cost recognized at merge time. |

No trigger fires that requires a Cycle Iteration entry to close.

## §4 INDEX update

Add to `.cdd/iterations/INDEX.md`:

```
| 395 | #395 | 2026-05-21 | 2 | 0 | 1 | 1 | .cdd/unreleased/395/cdd-iteration.md |
```

Findings: 2 (F1, F2). Patches: 0 immediate (F2 landed in the merge commit but is not a skill/spec patch). MCAs: 1 (F1 next-MCA on third occurrence). No-patch: 1 (F2 — landed-in-merge counts as resolved, not as protocol-patch-required). Path: `.cdd/unreleased/395/cdd-iteration.md` (will move to `.cdd/releases/<version>/395/cdd-iteration.md` at next release per `release/SKILL.md §2.5a`).

## §5 Skill-gap candidate disposition

- F1 (`cdd-protocol-gap`): next-MCA on third occurrence; patches `cnos.cdd/skills/cdd/issue/proof/SKILL.md` with the disavowal-context carve-out.
- F2 (`cdr-protocol-integration`): landed-in-merge; recorded as wave-coordination pattern (master issues for multi-Sub bootstraps pin shared-surface conventions before Sub dispatch).

## §6 Deferred outputs

- **Cycle-dir move.** `.cdd/unreleased/395/` → `.cdd/releases/<version>/395/` at next release. Standard release-time mechanic.
- **F1 next-MCA on third occurrence** — defer to cnos.cdw / cnos.cda bootstrap.
- **Lifecycle sub-skills not authored** (per issue Non-goals) — deferred to cds emergence.
- **Wave-coordination primitive for research γ** — may be addressable via Sub 4's cph-mapping content; re-evaluate when Sub 4 is read.

## §7 Next-MCA commitment

cnos#376 master has now shipped all four Subs (Sub 1 #390, Sub 2 #394, Sub 3 #395, Sub 4 #396); cnos#376 master can close. γ post-merge close-out comment on cnos#376 confirms Sub 3 shipped.

The next-MCA path is now the **persona hub** (`cn-rho`) bootstrap: cnos.cdr role overlays declare the protocol-layer role functions, but layer-1 (persona) content for the research persona is required for an actor to enact CDR work on a real research project. `cn-rho` bootstrap drafts are staged at `.cdd/iterations/cross-repo/cn-rho/bootstrap-2026-05-19/` per cnos#395 Non-goals.

Filed by ε on 2026-05-21.
