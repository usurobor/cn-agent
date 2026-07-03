# self-coherence — cycle/368

## Gap

**Issue:** [cnos#368](https://github.com/usurobor/cnos/issues/368) — bug(cdd): cycle close-out can leave issue OPEN — Closes-keyword inconsistently applied, γ fallback skipped on docs-only disconnect.

**Mode:** design-and-build (per γ scaffold — design ~80% converged; exact insertion points confirmed by γ peer-enumeration before α execution).

**Priority:** reconciled to `P2` this cycle (label was `P1`; issue body stated `Priority: P2` with rationale "issue-state drift; workaround is trivial; not release-blocking"). GitHub label corrected via `gh issue edit 368 --repo usurobor/cnos --remove-label P1 --add-label P2`. Rationale: the issue body's own stated rationale is sound — a docs-only close-token gap with a one-line manual-close workaround is not release-blocking, even though it recurs.

**Version:** doctrine/skill-prose cycle — no code, no version bump. Implementation-contract axes pinned by γ (scaffold §5) and δ (`gamma-clarification.md`): Markdown/prose-only, no new CLI subcommand or script, no new binaries, no wire-format change (additive artifact field only).

**Corrected surfaces (per γ scaffold §2/§3, binding over the issue's own citations):**
- AC1 target: `beta/SKILL.md` §"Pre-merge gate" (not `review/SKILL.md` — that section doesn't exist there).
- AC2 target: `gamma/SKILL.md` §2.10 (confirmed correct by γ).
- AC3 target: `operator/SKILL.md` §8 + `harness/SKILL.md` §6 (concrete home found by γ; issue left this under-specified).
- AC4 target: `cnos.cds/skills/cds/CDS.md` (not `CDD.md` — the instantiation-level close protocol migrated out of `CDD.md` in the cnos#366 Phase 7 rewrite; `CDD.md` is kernel-only and its own non-goals forbid re-deriving this doctrine there).

Branch `cycle/368` created by γ from `origin/main` at `d5bb2c20958998e236ab7c0d0a154ddc9ee319f2`. α picked up at HEAD `6ccc2f7` (δ's implementation-contract pin) and implemented on top.

## Skills

**Tier 1:** `CDD.md` (canonical kernel, loaded and confirmed kernel-only / no touch needed); `alpha/SKILL.md` (this role surface, load order followed: dispatch intake §2.1, produce-in-order §2.2, peer enumeration §2.3, self-coherence §2.5, pre-review gate §2.6, request review §2.7).

**Tier 3 (per dispatch prompt / γ scaffold §5):**
- `src/packages/cnos.core/skills/write/SKILL.md` — skill-prose discipline; kept additions terse, matching the existing row/bullet style rather than inventing new structure.
- `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` — primary AC2 edit target; also read in full for §2.8–§2.11 context around the closure gate.
- `src/packages/cnos.cdd/skills/cdd/beta/SKILL.md` (Pre-merge gate section) — primary AC1 edit target.
- `src/packages/cnos.cds/skills/cds/CDS.md` (Development lifecycle + Gate sections) — primary AC4 edit target; read state-machine table (§"State machine") and Gate → Release-readiness preconditions + Closure verification checklist (F1–F10) in full to place the addition correctly and avoid colliding with the existing F-anchor numbering.
- `src/packages/cnos.cdd/skills/cdd/operator/SKILL.md` (§8 Timeout recovery) — AC3 edit target (doctrinal pointer).
- `src/packages/cnos.cdd/skills/cdd/harness/SKILL.md` (§6 Timeout recovery mechanics) — AC3 edit target (mechanics; γ scaffold flagged this as the concrete home for the actual recovery steps).

**Not loaded (per α load-order discipline):** β/γ role skills were read as edit targets (their content, not their role-frame) — α did not adopt β's or γ's judgment frame, only edited the doctrine they carry, per Rule 3.5.

## ACs

### AC1 — β pre-merge gate fails if merge subject lacks a close-keyword

**Surface:** `src/packages/cnos.cdd/skills/cdd/beta/SKILL.md` §"Pre-merge gate" (corrected surface — the issue named `review/SKILL.md`, which has no such section; γ scaffold §2 item 2 confirmed the real table lives in `beta/SKILL.md`).

**Change:** new row 5 ("Close-keyword presence") added after existing row 4, plus one prose sentence giving the row-5 oracle command, plus one clause appended to the "mandatory for substantial cycles" sentence marking row 5 as non-collapsible. Rows 1–4 and their text are byte-identical to `HEAD~5` (verified below).

**Evidence:**
```
$ rg -n "Closes|Fixes|Resolves #" src/packages/cnos.cdd/skills/cdd/beta/SKILL.md
86:| 5 | **Close-keyword presence** | ... | ...
88:**Row 5 oracle command:** `git log -1 --format=%s | grep -iE "(Closes|Fixes|Resolves) #{N}"` ...
```
≥1 match — oracle satisfied per γ scaffold §4 AC1 row.

**Backward-compat check:**
```
$ git diff origin/main..HEAD -- src/packages/cnos.cdd/skills/cdd/beta/SKILL.md | grep '^-' | grep -v '^---'
-The pre-merge gate is mandatory for substantial cycles. Small-change merges may collapse rows 2 and 3 if the cycle's diff is purely textual / docs and no new contract surface is being shipped, but row 1 (identity truth) is mandatory for every β-authored commit, full stop.
```
The only removed line is the closing sentence, which is re-added verbatim plus one appended clause (see diff hunk in the commit `c9e3b41`); rows 1–4 have zero removed lines — confirmed no row 1–4 text was altered.

### AC2 — γ close-out asserts `gh issue view` state == CLOSED before declaring closure

**Surface:** `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` §2.10 (confirmed correct by γ scaffold §2 item 3 — no retargeting needed).

**Change:** new closure-gate item 15 (hard-gate assertion + auto-close-and-record fallback), plus an additive field appended to the `gamma-closeout.md` "Contains:" list (the asserted issue-close state + discrepancy note when applicable). Items 1–14 are untouched (verified below).

**Evidence:**
```
$ rg -n "gh issue view.*state|issue view --json state" src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md
334:15. **γ MUST assert the parent issue's close state before declaring closure...** Run `gh issue view {N} --json state --jq .state`. ...
337:- write `.cdd/unreleased/{N}/gamma-closeout.md`. Contains: ... **and the asserted issue-close state from row 15 above** ...
```
≥1 match — oracle satisfied per γ scaffold §4 AC2 row.

**Backward-compat check:**
```
$ git diff origin/main..HEAD -- src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md | grep '^-' | grep -v '^---'
-- write `.cdd/unreleased/{N}/gamma-closeout.md`. Contains: cycle summary, close-out triage table, §"Cycle iteration triggers" assessment, cycle iteration, skill gap candidate dispositions, deferred outputs, hub memory evidence, and next MCA. ...
```
Only the `gamma-closeout.md` "Contains:" sentence changed (additive field appended, "and next MCA" → "next MCA, and the asserted issue-close state..."); the 14 numbered closure-gate items are untouched, new item 15 appended after item 14 without renumbering.

### AC3 — δ-recovery path includes the same close-state assertion

**Surface:** `src/packages/cnos.cdd/skills/cdd/operator/SKILL.md` §8 "Timeout recovery" (doctrinal pointer) + `src/packages/cnos.cdd/skills/cdd/harness/SKILL.md` §6 (mechanics) — both named explicitly by γ scaffold §2 item 4 as the concrete existing home, correcting the issue's own under-specified "Operator SKILL.md or new δ-recovery section."

**Change:** `harness/SKILL.md` gains new §6.6 "Close-state assertion parity on recovery closure," spelling out the same `gh issue view {N} --json state --jq .state` == CLOSED hard gate and auto-close-and-record fallback, scoped to the recovery path (δ runs it in γ's stead). `operator/SKILL.md` §8 gains one new bullet pointing at `gamma/SKILL.md §2.10` row 15 and `harness/SKILL.md §6.6`. No existing §8 bullets or §6 subsections (6.1–6.5) were altered.

**Evidence:**
```
$ rg -n "gh issue view.*state" src/packages/cnos.cdd/skills/cdd/operator/SKILL.md src/packages/cnos.cdd/skills/cdd/harness/SKILL.md
operator/SKILL.md:294:- **Close-state assertion parity.** ... `gh issue view {N} --json state --jq .state` MUST return `CLOSED` ...
harness/SKILL.md:491:... `gh issue view {N} --json state --jq .state` MUST return `CLOSED` before the closure declaration is written. ...
```
≥1 match in both files — oracle satisfied per γ scaffold §4 AC3 row.

### AC4 — CDS.md states the assertion as a mandatory hard gate, not a conditional fallback (CDD.md untouched)

**Surface:** `src/packages/cnos.cds/skills/cds/CDS.md` §"Development lifecycle" → §"State machine" row S7 (~line 1026) and §"Gate" → §"Release-readiness preconditions" (~line 2295) — corrected surface per γ scaffold §2 item 1: `CDD.md`'s cited lines 254/261/348 no longer exist (the whole instantiation-level close protocol migrated to CDS.md under the cnos#366 Phase 7 split); editing `CDD.md` to restore that language would violate its own non-goals.

**Change:**
- S7 row: "Required inputs" column now cites row 5's close-keyword check; "Failure / retry" column gains a clause naming the mandatory (not conditional) issue-state assertion and pointing at `gamma/SKILL.md §2.10` row 15.
- Release-readiness preconditions: new bullet "Issue close-state asserted CLOSED" stating the assertion is mandatory, that a `Closes #N` merge subject is necessary but not sufficient, and that the precondition does not hold on an unasserted or still-OPEN state.

**CDD.md disposition — confirmed no edit needed, not a skipped pointer fix:**
```
$ rg -n "If the convention is missed|gh issue close" src/packages/cnos.cdd/skills/cdd/CDD.md
(no output, exit 1)
```
`CDD.md` carries zero soft-fallback phrasing to harden and zero literal close-protocol instruction to begin with (confirmed independently by α, matching γ scaffold §2 item 1's finding). `CDD.md` §"Software-specific realization → cnos.cds" (read in full, lines ~124–151) already carries the correct pointer to `cnos.cds/skills/cds/CDS.md` as the software-cycle realization home — no stale pointer exists to correct. Per the constraint in the α prompt ("if AC4 needs any CDD.md touch at all, it is at most a pointer-list correction... and only if you find CDD.md currently points somewhere stale"), the check-first condition resolved to "already correct" — **zero lines of CDD.md were changed this cycle.**

**Evidence (mandatory-language oracle):**
```
$ rg -n "If the convention is missed" src/packages/cnos.cds/skills/cds/CDS.md
(no output, exit 1)
$ rg -n "mandatory verification|MUST have run|hard gate" src/packages/cnos.cds/skills/cds/CDS.md | grep -i "issue\|close"
1026:| **S7: β merged** | ... issue state MUST still be asserted `CLOSED` ... this is a mandatory verification, not a conditional fallback ...
2299:- **Issue close-state asserted CLOSED.** This is a mandatory verification, ...
2306:  `CLOSED`, γ MUST have run `gh issue close {N}` and recorded the ...
```
Oracle satisfied per γ scaffold §4 AC4 row: 0 soft-fallback hits, mandatory language present at both named surfaces.

**Backward-compat check:** `git diff origin/main..HEAD -- src/packages/cnos.cds/skills/cds/CDS.md` shows one modified line (S7 row — extended in place, not removed-and-rewritten in a way that drops existing content: "Merge conflict → β resolves in throwaway worktree" is preserved verbatim as the first clause of the "Failure / retry" cell) and one new bullet inserted after the existing "Merge commit present" bullet (which is untouched). The F1–F10 closure verification checklist enumeration and anchors were **not** touched or renumbered — this cycle adds language to the preconditions prose, not a new F-row, per the scope guardrail against introducing new mechanized checks.

### AC5 — Falsification against the #367 incident

**Invariant under test:** the patched doctrine (AC1 β gate + AC2 γ hard gate) would have caught, at one layer or the other, the exact condition that left #367 OPEN for ~24h.

**Method (explicit, per scaffold friction note 7 and β prompt guidance):** this is an **evidence-based simulation, not a literal replay** — #367's actual merge (`37ac1c75`) and γ close-out (`704365d2`) are historical events that already happened under the *old* (unpatched) doctrine and cannot be re-executed. What follows instead is a mechanical application of the new gate's exact oracle commands against the real, historically-documented input (the #367 merge subject string, quoted verbatim in the issue body's §Problem table) and against a real counter-example (#362's correctly-formed subject), run live in this session.

**Step 1 — simulate the new AC1 β pre-merge gate row 5 against #367's actual merge subject:**
```
$ SUBJECT="Merge cycle/367 — Design CDD contract/receipt validation surface (#367)"
$ echo "$SUBJECT" | grep -iE "(Closes|Fixes|Resolves) #367"
$ echo "exit: $?"
exit: 1
```
No match — grep exits 1. Under the new row 5, β would read this exit code, recognize the close-keyword is absent, and **refuse to merge** until the subject is corrected (row 5 text: "rewrite the planned subject to include the close-keyword before merging"). This is the first-layer catch that did not exist at #367's actual merge time.

**Step 1 counter-check — same simulation against #362's real (correct) subject, to show the gate discriminates rather than failing everything:**
```
$ SUBJECT2="Merge cycle/362: UIE communication gate (Closes #362)"
$ echo "$SUBJECT2" | grep -iE "(Closes|Fixes|Resolves) #362"
Merge cycle/362: UIE communication gate (Closes #362)
$ echo "exit: $?"
exit: 0
```
Match found — grep exits 0, row 5 passes, β proceeds. This is the falsifiability requirement: the same mechanical check that fails #367's subject passes #362's subject, so the gate is discriminating on the real invariant (presence of a close-keyword), not vacuously failing or passing every input.

**Step 2 — simulate the new AC2 γ hard-gate assertion against #367's actual post-merge state:**

Historical record (issue body §Problem table, γ scaffold §3 "Falsification case," unchanged/confirmed): #367 merged 2026-05-15 13:48 UTC with the subject above (no close-keyword); γ close-out landed at `704365d2` 13:54 UTC; the issue remained `OPEN` for ~24h until a manual `gh issue close` on 2026-05-16. Under the *old* doctrine, γ's close-out at `704365d2` had no step requiring `gh issue view --json state` to be checked — closure was declared without the assertion, which is exactly how an OPEN issue survived past a "closed" cycle.

Applying the new gamma/SKILL.md §2.10 item 15 to that same historical moment (13:54 UTC, immediately after `704365d2`): γ would run `gh issue view 367 --json state --jq .state`. Given the documented fact that the issue was OPEN at that moment (it stayed open until the next day), the assertion would observe `OPEN`, not `CLOSED`. Per item 15's mandatory text, γ **MUST** then run `gh issue close 367` immediately (at 13:54 UTC, not 24h later) and record the discrepancy in `gamma-closeout.md`. The patched doctrine collapses the ~24h drift window to effectively zero — closure and the state-correction happen in the same close-out commit.

**Negative control (unpatched behavior, for contrast):** re-reading `gamma/SKILL.md` at `origin/main` HEAD (pre-this-cycle) confirms item 15 did not exist — `git show origin/main:src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md | sed -n '316,340p'` ends at item 14, with no state-assertion step before the "Then: write gamma-closeout.md" instruction. This is the exact absence that let #367 close-out proceed without the check.

**Result:** both new gate layers, mechanically applied to the real #367 inputs, fire as designed — AC1's row 5 refuses the bare-reference subject at merge time; AC2's item 15 (had it existed at close-out time) would have caught and corrected the resulting OPEN state within the same close-out session rather than 24h later. This is falsifiable evidence (a command run against real historical strings, with exit codes shown, including a discriminating counter-example), not a narrative assertion that "this fixes it."

**Known limitation, disclosed:** this is a simulation against historical strings and documented facts, not a re-execution of the actual #367 GitHub merge/close API calls (those cannot be replayed without mutating real repo state). The Debt section below records this explicitly.

## CDD Trace

- **γ scaffold** (`gamma-scaffold.md`): γ branched `cycle/368` from `origin/main` at `d5bb2c20`, corrected the four AC surfaces off the issue's own citations (AC1 `beta/SKILL.md`, AC3 `operator`+`harness`, AC4 `CDS.md` not `CDD.md`), peer-enumerated the insertion points, and pinned the design-and-build mode (design ~80% converged).
- **δ clarification** (`gamma-clarification.md`): implementation-contract axes pinned — Markdown/prose-only, no new CLI/script/binary, additive artifact field only; α picked up at HEAD `6ccc2f7`.
- **α implement** (this artifact + the five skill/doctrine edits): AC1 `beta/SKILL.md` row 5; AC2 `gamma/SKILL.md` item 15; AC3 `harness/SKILL.md` §6.6 + `operator/SKILL.md` §8 bullet; AC4 `CDS.md` S7 + release-readiness precondition; AC5 falsification simulation. Every edit is additive (rows/items/bullets appended, no existing row/item text altered — backward-compat checks per AC above).
- **β review / merge**: not reached in the dispatched run — the substrate run failed at the close-out write-fence (`dispatch_activation_log_write_violation` on `.cn-sigma/logs/`); recovery carried to merge via the δ-recovery path (`operator/SKILL.md` §8 + `harness/SKILL.md` §6), which this cycle's own AC3 hardens. The close-state hard gate this cycle ships (AC2/AC3) applies to this cycle's own close-out.

## Debt

- **AC5 is a simulation, not a live replay.** The #367 falsification runs the new gate's oracle commands against #367's real historical merge-subject string and #362's counter-example (exit codes shown), not against a re-executed GitHub merge/close — those historical API events cannot be replayed without mutating real repo state. The evidence is a discriminating mechanical check (fails #367's subject, passes #362's), which is the strongest falsification available without time travel.
- **Dispatch-substrate defect (not a #368-content defect).** The dispatched run's agent work completed, but the run failed at the `.cn-{agent}/logs/` write-fence; the pushed `cycle/368` branch is clean of any `.cn-sigma/logs/` change (`git log origin/main..cycle/368 -- .cn-sigma/logs/` is empty). Worth a separate follow-up on the substrate fence + a stale-`status:in-progress` reaper (the #504-class failure mode that stranded this very cycle).
- **δ-recovery override.** This cycle reached main via δ-recovery rather than a normal β merge; per `operator/SKILL.md` §8 the recovery-effector merge is an implicit override and the γ-axis grade reflects it.
