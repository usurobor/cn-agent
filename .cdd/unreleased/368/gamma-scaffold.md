# γ scaffold — cycle/368

**Issue:** [cnos#368](https://github.com/usurobor/cnos/issues/368) — bug(cdd): cycle close-out can leave issue OPEN — Closes-keyword inconsistently applied, γ fallback skipped on docs-only disconnect

**Mode:** design-and-build (per issue body §Cycle scope sizing — design ~80% converged, exact insertion points confirmed below by γ peer-enumeration; not MCA because CDD.md's own citations have drifted and needed re-pointing before α can execute a known sequence).

**Branch:** `cycle/368`, created from `origin/main` at `d5bb2c20958998e236ab7c0d0a154ddc9ee319f2`.

**Priority:** label says `P1`; issue body says `Priority: P2`. Per the dispatch comment's own triage nit, this reconciles during the cycle — see Friction Notes item 1. Not a scaffolding blocker.

---

## 0. SHA-pin verification (γ pre-flight duty)

- Given main SHA (wake input #3): `e11cca164c3bf04f70ed7a913c1a9d2be1ea433a`
- Working-tree `git rev-parse HEAD` at scaffold start: `e11cca164c3bf04f70ed7a913c1a9d2be1ea433a` — **match confirmed.**
- `git fetch origin main && git rev-parse origin/main` at scaffold time: `d5bb2c20958998e236ab7c0d0a154ddc9ee319f2` — **origin/main has advanced 3 commits** beyond the pinned SHA (`303fdf5`, `a7c068a`, `d5bb2c2` — heartbeat log + two board-map regenerations). `git diff --stat e11cca1..d5bb2c20` touches only `.cn-sigma/logs/20260703.md` and `docs/development/board/{board-data.json,index.html}` — no overlap with this cycle's scope surfaces (`beta/SKILL.md`, `gamma/SKILL.md`, `CDS.md`, `operator/SKILL.md`, `review/SKILL.md`, `CDD.md`). Per gamma/SKILL.md's canonical-skill staleness check, a re-load is required only when the drift touches a role/canonical skill file loaded for this phase — it does not here. **Disposition: benign drift, non-blocking.** Branch created from current `origin/main` (`d5bb2c2`), not the stale pinned SHA, per the branch-rule doctrine (`cnos.cds/skills/cds/CDS.md` §"Development lifecycle" → §"Branch rule").

## 1. Branch pre-flight (γ-owned, `CDS.md` §"Branch pre-flight")

- `origin/cycle/368` does not yet exist — confirmed (`git ls-remote --exit-code --heads origin cycle/368` → exit 2).
- No stalled `.cdd/unreleased/368/` on `origin/main` — confirmed (`git ls-tree -r --name-only origin/main -- .cdd/unreleased/368/` → empty).
- Issue scope is declared in the body (Scope / In-scope / Out-of-scope / Deferred sections present).
- Issue state: `OPEN` (`gh issue view 368 --json state` → `OPEN`).
- Base SHA known: `d5bb2c20958998e236ab7c0d0a154ddc9ee319f2`.

## 2. Peer enumeration at scaffold time (binding per gamma/SKILL.md §2.2a)

Before trusting the issue's own Source-of-truth table, γ grepped every surface the issue cites and the surfaces γ expects α to touch. Two of the issue's own citations do not resolve as stated — both are named here explicitly rather than left for α to discover mid-cycle:

1. **CDD.md has been rewritten since the issue was filed.** `src/packages/cnos.cdd/skills/cdd/CDD.md` is now v4.0.0 (160 lines, kernel-only doctrine per cnos#366 Phase 7 terminal). The issue's cited lines (254, 261, 348 — the `Closes #N` protocol, the "If the convention is missed, γ closes via `gh issue close`" fallback, and the γ algorithm step) no longer exist in this file at all. `rg "If the convention is missed|gh issue close" src/packages/cnos.cdd/skills/cdd/CDD.md` → 0 matches (this makes the issue's own AC4 oracle line — `rg "If the convention is missed" CDD.md` returns 0 — **vacuously true today**, not because the doctrine was hardened, but because the whole instantiation-level chapter migrated out of CDD.md under the CDS split). `rg -r "gh issue close" src/packages/cnos.cdd/skills/cdd/ src/packages/cnos.cds/skills/cds/CDS.md` returns **zero literal instruction hits** anywhere in current canonical skill files — the only hit repo-wide is a forward citation in `harness/SKILL.md` line ~185 ("Close-token gap (#368) early-detection — γ's missing `gh issue close` surfaces in the stream, not only after process exit."), which is a citation of this very issue, not an implementation. **Conclusion: the close-fallback doctrine was dropped in the CDD→CDS migration, not merely under-enforced.** AC4's real target is `cnos.cds/skills/cds/CDS.md`, not `CDD.md`.
2. **AC1's stated surface does not exist.** The issue names `src/packages/cnos.cdd/skills/cdd/review/SKILL.md` pre-merge gate as AC1's surface. `review/SKILL.md` has no "Pre-merge gate" section (confirmed: its §3.11/§3.11b cover merge-instruction wording and the γ-artifact-completeness gate, not a structural table). The actual **Pre-merge gate** table (rows 1–4: identity truth / canonical-skill freshness / non-destructive merge-test / γ-artifact completeness) lives at `src/packages/cnos.cdd/skills/cdd/beta/SKILL.md` §"Pre-merge gate" (lines 76–87). Row 4 of that table is the closest existing precedent for exactly this kind of structural, fails-the-merge gate. AC1 should add a **row 5** there, not edit `review/SKILL.md`.
3. **AC2's surface is correctly cited.** `gamma/SKILL.md` §2.10 "Steps 13–15 — Close only after the closure gate passes" (lines 316–339) lists 14 closure-gate rows; none currently assert `gh issue view {N} --json state`. This is the right insertion point, confirmed by direct read.
4. **AC3's surface has a concrete existing home the issue didn't name.** The issue says "Operator SKILL.md or new δ-recovery section." `operator/SKILL.md` §8 "Timeout recovery" (lines 287–296) already codifies the δ-recovery doctrinal frame (citing `harness/SKILL.md` §6 for recovery mechanics). α should extend §8 (and/or `harness/SKILL.md` §6 mechanics) with the close-state assertion rather than inventing a new section.
5. No other in-flight or landed work overlaps this gap — `rg "issue.*state.*CLOSED|--json state" src/packages/cnos.cdd/ src/packages/cnos.cds/` outside the harness forward-citation returns nothing; the gap is real and unclaimed.

These five findings are additive framing per §2.2a — they retarget three of the issue's own AC surfaces to where the doctrine actually lives today, without changing what the ACs require.

## 3. Corrected source-of-truth table

| Claim / surface | Canonical source (issue's citation) | Corrected canonical source (γ peer-enumeration) | Status |
|---|---|---|---|
| Cycle merge + close protocol | `CDD.md` §S7, lines 254/261/348 | `cnos.cds/skills/cds/CDS.md` §"Development lifecycle" → §"State machine" row **S7** (~line 1026) + §"Gate" → §"Release-readiness preconditions" (~line 2295) | Migrated; issue's line numbers stale |
| β merge instruction | `review/SKILL.md` lines 132, 202 | `review/SKILL.md` §3.11 (merge-instruction wording) **for the keyword-wording rule**; `beta/SKILL.md` §"Pre-merge gate" rows 1–4 (~lines 76–87) **for the structural fails-the-merge gate AC1 actually wants** | Both exist; AC1 targets the latter |
| γ close-out responsibility | `gamma/SKILL.md` | `gamma/SKILL.md` §2.10 "Steps 13–15" (~lines 316–339) | Confirmed correct |
| δ-recovery path | (not cited) | `operator/SKILL.md` §8 "Timeout recovery" (~lines 287–296); mechanics in `harness/SKILL.md` §6 | Concrete home found |
| Falsification case (#367) | merge `37ac1c75`, γ close-out `704365d2` | unchanged | Shipped |
| Counter-evidence (working path) | #359/#360/#361/#364/#365 | unchanged | Shipped |
| Root-cause analysis | `cn-sigma:threads/adhoc/20260516-cdd-close-token-gap.md` | unchanged (external hub path, not verified by γ — outside repo) | Cited, not re-verified |

## 4. Per-AC oracle list (mechanical)

| AC | Invariant | Mechanical oracle | Corrected surface |
|---|---|---|---|
| AC1 | Merge subject must carry a GitHub close-keyword for the cycle's issue | New row 5 present in the Pre-merge gate table; row text requires verifying the planned merge subject matches `(Closes|Fixes|Resolves) #{N}` (case-insensitive) before push; `rg "Closes \|Fixes \|Resolves #" src/packages/cnos.cdd/skills/cdd/beta/SKILL.md` (within the Pre-merge gate section) → ≥1 match | `beta/SKILL.md` §"Pre-merge gate" |
| AC2 | γ close-out does not declare closure unless `gh issue view {N} --json state --jq .state` returns `CLOSED` | New row/bullet in `gamma/SKILL.md` §2.10 naming the assertion as hard-gate #15 (or folded into existing numbered list); `gamma-closeout.md` template gains a field recording the asserted state; `rg "gh issue view.*state\|issue view --json state" src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` → ≥1 match | `gamma/SKILL.md` §2.10 |
| AC3 | δ-recovery fires the same assertion before declaring a recovered cycle closed | `operator/SKILL.md` §8 (and/or `harness/SKILL.md` §6) names the assertion explicitly as part of the recovery procedure; `rg "gh issue view.*state" src/packages/cnos.cdd/skills/cdd/operator/SKILL.md src/packages/cnos.cdd/skills/cdd/harness/SKILL.md` → ≥1 match | `operator/SKILL.md` §8 / `harness/SKILL.md` §6 |
| AC4 | Doctrine states the assertion as a mandatory verification, not a conditional fallback | `rg "If the convention is missed" src/packages/cnos.cds/skills/cds/CDS.md` → 0 (no soft-fallback phrasing survives); CDS.md's S7 row and/or Gate section name the close-state assertion as mandatory language; `CDD.md` is not re-edited (kernel-only per its own non-goals) | `cnos.cds/skills/cds/CDS.md` (not `CDD.md`) |
| AC5 | A test/simulated docs-only-disconnect cycle following the patched skills cannot reproduce the #367 OPEN-state condition | Throwaway/simulated cycle recorded in `self-coherence.md` with before/after `gh issue view --json state` evidence; positive case closes, negative case (unpatched behavior) is shown to have left it open per the #367 precedent | Test evidence in `self-coherence.md` |

## 5. α prompt

```text
You are α. Project: cnos.
Load src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md and follow its load order.
Issue: gh issue view 368 --json title,body,state,comments --repo usurobor/cnos
Branch: cycle/368
Tier 3 skills: src/packages/cnos.core/skills/write/SKILL.md, src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md, src/packages/cnos.cdd/skills/cdd/beta/SKILL.md (Pre-merge gate section), src/packages/cnos.cds/skills/cds/CDS.md (Development lifecycle + Gate sections), src/packages/cnos.cdd/skills/cdd/operator/SKILL.md (§8 Timeout recovery)
Scaffold: .cdd/unreleased/368/gamma-scaffold.md on this branch — read §2 (peer enumeration) and §3 (corrected source-of-truth table) before editing; the issue body's own AC1/AC4 surface citations are stale, corrected surfaces are named there.
```

### Implementation contract (γ-pinned; unpinned rows named for δ)

| Axis | Pinned value | Pinned by |
|---|---|---|
| Language | Markdown (skill/doctrine prose only — no source code) | γ, per issue Cycle-scope-sizing row (a): "0 new modules; skill + doctrine edits" |
| CLI integration target | N/A — no new `cn` subcommand. Gate/assertion rows are agent-executed shell one-liners documented in skill prose, mirroring the existing pattern in `beta/SKILL.md` Pre-merge gate rows 1–4 and `gamma/SKILL.md` §2.10's 14 closure-gate rows (all prose-described, agent-run checks, not CLI additions). | γ, by repo-convention analogy — **flagged for δ confirmation**: if δ instead wants this codified as an executable check (e.g. a `cn cdd verify` extension or a shell script under `scripts/`), that is a real design fork the issue's AC1/AC2 oracles leave open ("checklist row or executable check") and should be re-pinned before α starts, not left to α's judgment |
| Package scoping | `src/packages/cnos.cdd/skills/cdd/beta/SKILL.md` (Pre-merge gate, new row 5); `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` (§2.10, new closure-gate row); `src/packages/cnos.cds/skills/cds/CDS.md` (§"Development lifecycle" → §"State machine" S7 row + §"Gate" section); `src/packages/cnos.cdd/skills/cdd/operator/SKILL.md` (§8 Timeout recovery) and/or `src/packages/cnos.cdd/skills/cdd/harness/SKILL.md` (§6 recovery mechanics) | γ, per corrected source-of-truth table §3 above |
| Existing-binary disposition | N/A — no binaries touched | γ |
| Runtime dependencies | None beyond the already-standard `gh` CLI (already a dependency of every β/γ role per existing skill prose) | γ |
| JSON/wire contract preservation | N/A — no wire-format change. Note (not a contract change): the `gamma-closeout.md` artifact template gains one additive field recording the asserted `gh issue view --json state` result; this is an artifact-shape addition, additive only, not a breaking wire change. | γ |
| Backward-compat invariant | All existing Pre-merge gate rows (1–4 in `beta/SKILL.md`) and all existing closure-gate rows (1–14 in `gamma/SKILL.md` §2.10) are preserved unchanged; new rows are additive only. `CDD.md` is not edited (its own non-goals forbid re-deriving instantiation-level doctrine it already delegated to CDS). | γ, per issue's Non-goals ("do not refactor docs-only disconnect path beyond close-token surface") |

**Unpinned axis flagged for δ enrichment before α dispatch:** *CLI integration target* — γ's default pin (prose-only, no new script/subcommand) is a judgment call by analogy to existing gate rows, not a hard repo convention citation. δ should confirm this reading or re-pin explicitly before routing the α prompt; if re-pinned, log in `.cdd/unreleased/368/gamma-clarification.md` per the mid-flight wire-format before α starts.

## 6. β prompt

```text
You are β. Project: cnos.
Load src/packages/cnos.cdd/skills/cdd/beta/SKILL.md and follow its load order.
Issue: gh issue view 368 --json title,body,state,comments --repo usurobor/cnos
Branch: cycle/368
```

β verifies (in addition to the standard Pre-merge gate and Rule 7 implementation-contract coherence check):
- Each of AC1–AC5's mechanical oracles (§4 above) actually resolves against the diff, not against α's prose description of the diff (β Role Rule 6 — anchor oracle evidence on code, not doc).
- AC4's oracle is checked against `cnos.cds/skills/cds/CDS.md`, not `CDD.md` (the issue's own oracle line cites `CDD.md`, which — per §2 peer enumeration above — no longer carries this doctrine at all; a literal re-run of the issue's stated oracle against `CDD.md` passes vacuously and proves nothing. β must not accept the vacuous pass as satisfying AC4 — β checks the corrected CDS.md surface).
- The new beta/SKILL.md Pre-merge gate row 5 does not regress rows 1–4 (backward-compat invariant above) — diff `beta/SKILL.md` and `gamma/SKILL.md` for byte-level preservation of existing rows outside the new insertions.
- AC5's falsification evidence in `self-coherence.md` is a genuine test/simulation, not an assertion — β Role Rule 6 applies here directly (the #367 precedent already showed doctrine-says-X / practice-does-Y drift once; a second doc-only assertion of "this fixes it" without executed evidence repeats the same class of failure this issue exists to close).

## 7. Scope guardrails

**In scope (per issue body):**
1. β-side structural lint: merge subject must carry a close-keyword for the cycle's parent issue — added to `beta/SKILL.md` §"Pre-merge gate" (corrected surface; issue said `review/SKILL.md`).
2. γ-side close-out self-check asserting `gh issue view {N} --json state` == `CLOSED` before the close-out commit, with an automatic `gh issue close {N}` + discrepancy record if not.
3. δ-recovery parity — the same assertion fires when δ stitches a SIGTERM-terminated cycle.
4. Doctrine update making the assertion a hard gate, not a conditional fallback — targets `CDS.md` (corrected; issue said `CDD.md` §S7 + γ algorithm, which no longer exist there).

**Out of scope (per issue body — do not let α drift into these):**
- Backfilling close state on already-merged-but-OPEN issues (separate cleanup scan).
- Restructuring the docs-only disconnect path beyond the close-token surface (belongs to cnos#366 Phase 4/5).
- Changing the merge mechanism itself (`git merge` stays; no `gh pr merge` adoption).
- Any edit to `CDD.md` kernel sections (§Kernel–§Scope-lift) — CDD.md's own non-goals forbid re-deriving instantiation doctrine there; if AC4 needs a CDD.md touch at all, it is at most a pointer-list correction, not new doctrine.

**Deferred (queue as separate follow-up, not this cycle):**
- One-time scan for orphan-OPEN issues with closed cycle dirs.

## 8. Friction notes

1. **P1/P2 label-vs-body priority mismatch.** Issue label is `P1`; issue body states `Priority: P2` with rationale ("issue-state drift; workaround is trivial; not release-blocking"). The dispatch comment already names this as a triage nit to reconcile during the cycle. γ does not resolve it here — noting it per the wake-invocation instruction; does not block scaffolding or dispatch. Whoever closes the cycle (γ close-out) should pick one and correct the GitHub label to match, with a one-line rationale.
2. **Issue's own doctrine citations are stale.** All three of the issue's CDD.md line citations (254/261/348) predate the cnos#366 Phase 7 rewrite of `CDD.md` to v4.0.0 (kernel-only, 160 lines). The instantiation-level close protocol migrated to `cnos.cds/skills/cds/CDS.md`. This is not a fabrication risk on γ's part — it is a direct `git`/`rg` finding, recorded in §2 and §3 above, and it changes AC1's and AC4's real target files. α should not "fix" `CDD.md` to restore the old fallback language; that would violate CDD.md's own non-goals.
3. **AC1's stated surface does not exist as named.** `review/SKILL.md` has no "Pre-merge gate" section; the actual table is in `beta/SKILL.md`. Corrected in §3/§4/§5 above.
4. **AC3's surface was under-specified in the issue** ("Operator SKILL.md or new δ-recovery section") but a concrete existing section (`operator/SKILL.md` §8 "Timeout recovery") already exists and is the natural extension point — named explicitly to save α a discovery round.
5. **Unpinned implementation-contract axis:** *CLI integration target* — γ's default pin is "no new CLI/script, prose-only gate rows," by analogy to the existing rows 1–4 / 1–14 patterns. This is a judgment call, not a hard convention citation, and is flagged in §5 for δ to confirm or override before α dispatch.
6. **origin/main drift during scaffolding** (§0 above) — benign; branch created from current `origin/main` (`d5bb2c2`), not the stale pinned SHA from the wake input. No re-load of canonical skills was required because the drift did not touch any skill file this cycle depends on.
7. **AC5 falsification risk.** The issue's own proof plan accepts "test run; evidence in self-coherence.md" for AC5 — this is inherently softer than the other ACs' `rg`-mechanical oracles. β should hold α to an actual executed simulation (§6 above), not a narrative claim, given that narrative-only claims are exactly the failure class (#367) this issue exists to close.

## 9. Cross-references

- #366 — parent roadmap (Phase 7 terminal CDD.md rewrite is the direct cause of friction notes 2–4 above; Phases 4/5 will revisit these surfaces further — out of scope here).
- #367 — falsification case.
- #362 — counter-example (correct `Closes #362` usage).
- `cn-sigma:threads/adhoc/20260516-cdd-close-token-gap.md` — root-cause analysis (external hub path, cited by issue, not independently re-verified by γ).

Filed by γ (wake-invoked dispatch, wake run `28663208946`) on 2026-07-03.
