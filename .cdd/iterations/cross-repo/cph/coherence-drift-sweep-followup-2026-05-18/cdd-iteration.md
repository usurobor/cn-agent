# ε cdd-iteration — cph coherence-drift-sweep-followup wave (2026-05-18)

**Source cycle:** cph#21 (master) + cph#{22,23,24,25} (subs), wave artifacts at `cph:.cdd/waves/coherence-drift-sweep-followup-2026-05-18/`
**Source branch:** `cph:claude/review-repo-coherence-PNbjQ`
**ε actor:** cnos γ acting as cph ε per role-collapse rule (`epsilon/SKILL.md §2`); ε work authored from cnos session, ε output emitted as feedback artifact for cph γ to apply on receipt
**ε output date:** 2026-05-19

---

## Hat-collapse acknowledgment

This iteration is authored by an actor playing γ in cnos and ε in cph during the same session. Per `epsilon/SKILL.md §2` ("ε and δ are often the same actor") and `ROLES.md §4` (role collapse rules), the ε hat is permissible here because (i) the operator explicitly delegated ε; (ii) the cph wave-closeout and receipt explicitly defer master closure and protocol-gap dispositions to ε; (iii) the ε work is reflective (read close-outs, decide protocol patches) and does not collapse with α / β responsibilities for the in-flight cycle.

The substantive ε work below is independent of cnos γ work; the hat-collapse is recorded for honest attribution rather than to override role separation.

---

## §1 Findings dispositioned

### F1: AC oracle regex brittleness — anchor evidence on code, not doc

- **Source:** cph wave-closeout §"cross-sub debt item 3" + receipt §4 cdd-protocol-gap candidates
- **Class:** `cdd-skill-gap` + `cdd-protocol-gap`
- **Trigger:** review pattern across cycles — cph#16 wave (precursor) missed F7 (`quality_flag="short"`) by anchoring AC1 oracle on doc; followup R1 re-grepped code first and surfaced the miss + widened the regex β-side
- **Description:** Dispatcher-prompt AC oracle regexes are authored before α writes the code and cannot anticipate every literal shape α emits (single-line vs multi-line conditionals, escape quirks). When β anchors on doc text rather than code, doc-drift produces an honest-looking pass. The structural fix is β-side discipline: re-grep code first, widen the regex β-side when the prompt regex misses a real literal α emits, record both forms in the verdict.
- **Root cause:** No canonical β skill rule named "anchor evidence on code, not doc"; each wave's β manifest had to restate the discipline. The codification gap let the precursor wave's β anchor on doc text without challenge.
- **Disposition:** `patch-landed` — Rule 6 added to `cnos:src/packages/cnos.cdd/skills/cdd/beta/SKILL.md` ("Anchor oracle evidence on code, not doc") consolidating both this finding and F2 below. Patch is on cnos branch `claude/file-cnos-cdr-issue-fi9Ld` (commit forthcoming).
- **Patch:** `cnos:src/packages/cnos.cdd/skills/cdd/beta/SKILL.md` §"Role Rules" §6 (added) + §6a "Widen brittle oracle regexes β-side" (added)
- **Affects:** every cnos / cph / sibling-repo cycle that vendors cnos cdd; future cycles inherit the discipline without per-wave manifest restatement
- **Cross-repo:** cnos branch `claude/file-cnos-cdr-issue-fi9Ld` carries the patch; cph γ folds this F1 entry into cph's own `cdd-iteration.md` once accepted

### F2: Name-overpromise as a β-axis review class

- **Source:** cph wave-closeout §"cross-sub debt item 10" + receipt §4
- **Class:** `cdd-skill-gap` + `cdd-protocol-gap`
- **Trigger:** review pattern across cycles — F7 (`quality_flag="short"` promised but never emitted) and F10 (`extract_shape` named but only marking persistence) share a structural class: a name encoding a promise the code does not fulfill
- **Description:** Function names, variable names, and string literals can encode promises (richer behavior, specific outputs, distinguishing flags) that the implementation does not fulfill. The drift between interface-level promise and substance-level delivery is a distinct review class worth naming, with two fix paths: rename to match (cosmetic narrowing) or implement what the name promises (substantive widening). The choice is α's; β's job is to surface that the choice is being made.
- **Root cause:** Same as F1 — no canonical β skill rule naming the class; cph#21 wave-closeout proposed the rule but it had no home. β re-derived the class per-wave.
- **Disposition:** `patch-landed` — folded into Rule 6 above as §6b "Surface name-overpromise as a finding class" with empirical-anchor citation to F7 + F10. Same patch as F1.
- **Patch:** `cnos:src/packages/cnos.cdd/skills/cdd/beta/SKILL.md` §6b (added)
- **Affects:** same as F1
- **Cross-repo:** same as F1

### F3: CHANGELOG-scope policy (cph#21 AC6 / F12)

- **Source:** cph wave-closeout §"F12" + cph receipt §3 escalation #2
- **Class:** `cdd-skill-gap` (cph-side `CDR.md` doctrine, not a cnos doctrine concern)
- **Trigger:** γ process-gap check — cph wave-closeout deferred to ε whether qualitative coherence waves merit retroactive CHANGELOG entries
- **Description:** The current `CDR.md` §"Changelog rule" reads "Each meaningful research wave gets a CHANGELOG.md entry with: empirical state (REVISE / GO / STOP, or no-change with reason) [...]". Two readings: (a) narrow — only empirical research waves warrant entries; coherence-drift sweeps describe code coherence, not empirical state, so no entry; (b) broad — both coherence-drift waves moved the repo's coherence surface and merit `0.1.1-cdr` + `0.1.2-cdr` entries.
- **Root cause:** Wording "meaningful research wave" is undefined; both readings are defensible without further constraint. The changelog's own header ("a coherence ledger, not a release-note style log") reads more narrowly than the rule text.
- **ε decision:** Adopt reading **(a)**. The body of the rule lists fields that only an empirical wave can fill ("empirical state (REVISE / GO / STOP)", "α / β / γ / C_Σ if measured"). Forcing those fields on a coherence wave produces a hollow entry ("no-change with reason" + "pending — coh unavailable" + zero coherence delta on empirical state). The narrow reading matches both the field list and the changelog's self-description as a coherence-of-empirical-state ledger. No retroactive entries for the two coherence-drift waves.
- **Disposition:** `patch-landed` (substance authored here, application pending cph γ) — patch text for `CDR.md` §"Changelog rule" emitted as `cdr-changelog-rule.patch` in this bundle for cph γ to apply
- **Patch:** `cph:CDR.md` §"Changelog rule" — replace "Each meaningful research wave gets a..." with "Each empirical research wave (one that produces a REVISE / GO / STOP transition or moves the empirical-state ledger) gets a..."; add explicit exclusion clause "Coherence-only waves (drift sweeps, doc reconciliation, protocol-internal refactors) do not produce CHANGELOG entries; their record lives in `.cdd/waves/` and `.cdd/iterations/`."
- **Affects:** cph CDR doctrine; future cph coherence-only waves
- **Cross-repo:** patch lives in cph; emitted as feedback patch in this bundle

### F4: Master closure recommendations (cph#21, cph#16)

- **Source:** cph wave-closeout §"Master Closure Recommendations" + receipt §3 escalation #1
- **Class:** none (operational recommendation, not a protocol gap)
- **ε decision:** Recommend cph γ close **both** cph#21 and cph#16.
  - cph#21: With F3 resolved (reading (a) adopted), AC6 is met by the doctrinal tightening rather than by retroactive CHANGELOG entries. ACs 1–8 are met. The 10 cross-sub debt items are either dispositioned in this iteration (F1, F2, F3), pre-existing deferrals (C_Σ baseline, `coh` PATH), or in-cycle housekeeping. Close cph#21 with a closure comment that links the wave-closeout, this `cdd-iteration.md`, and the four sub closeouts.
  - cph#16: The precursor wave deferred cph#16 closure to ε. cph#16's substantive findings carried into cph#21's followup wave, which dispositioned them (cleanly, 0 fix-rounds, all 4 subs APPROVE R1). cph#16's deferred ACs are now addressed via cph#21. Close cph#16 with a closure comment citing cph#21 as the resolution path.
- **Disposition:** `next-MCA` (cph γ effects close; this is recommendation, not effection)
- **Issue filed:** none (recommendation only; cph γ owns the close)
- **First AC for cph γ:** verify F1+F2 patch lands on cnos main (or carry as named debt with the cnos branch SHA), then close cph#21 with the four cited linkages, then close cph#16 with the cph#21 linkage

### F5: Operator-gated branch merge (claude/review-repo-coherence-PNbjQ → main)

- **Source:** cph wave-closeout §"Master Closure Recommendations" + receipt §3 escalation #1, paragraph 4
- **Class:** none (operator authority, not ε)
- **ε decision:** Out of ε's authority. Surface to operator. Recommendation: merge after F3 patch lands on cph (cdr-changelog-rule.patch applied) so the merged main reflects the tightened doctrine.
- **Disposition:** `no-patch` (operator gate; ε is not a release-driver)
- **Reason:** Per `ROLES.md §1` row 5 and `epsilon/SKILL.md §1`, ε iterates protocol; ε does not drive merges or tag releases. The operator (δ in cph) holds merge authority.

---

## §2 Trigger assessment (per gamma/SKILL.md §2.8 table)

| Trigger | Fire condition | Fired? | ε note |
|---|---|---|---|
| Review churn | review rounds > 2 | **No** | All 4 subs APPROVE R1 (0 fix-rounds). Wave was clean. |
| Mechanical overload | mechanical ratio > 20% AND findings ≥ 10 | **Borderline** | 10 cross-sub debt items, but the wave's findings are policy/doctrine-shaped, not mechanical. The 20% mechanical threshold does not fire on inspection; the count threshold is met by aggregating across subs but each sub's findings were modest. No mechanization patch landed; the closer-fitting trigger is the next one. |
| Avoidable tooling / environment failure | environment blocked the cycle | **Yes** | `coh` not on PATH blocked first numeric C_Σ baseline (carried from cph#11, cph#16, this wave). Disposition: `cdd-tooling-gap`; out of repo scope; harness-side. ε notes the recurring deferral as a candidate for an MCI in cph's `INDEX.md`. |
| Loaded-skill miss | a loaded skill should have prevented a finding | **Yes** | The cph wave β did *not* have a loaded `cdd/beta/SKILL.md` rule for code-first oracle anchoring — that gap is exactly what F1+F2 address. The miss did not produce a fix-round (β re-derived the discipline manually), but codification removes the per-wave manifest restatement cost. Patch landed (F1, F2). |

---

## §3 No-patch decisions

None this cycle. All findings produced either a `patch-landed`, `next-MCA`, or operator-gate disposition. Silence is not triage — no findings dropped.

---

## §4 INDEX update

`cph:.cdd/iterations/INDEX.md` should add one row for this cycle:

```markdown
| Cycle | Issue | Date | Findings | Patches | MCAs | No-patch | Path |
|-------|-------|------|----------|---------|------|----------|------|
| coherence-drift-sweep-followup-2026-05-18 | cph#21 | 2026-05-19 | 5 | 3 | 1 | 0 | .cdd/waves/coherence-drift-sweep-followup-2026-05-18/cdd-iteration.md |
```

Patches counted: F1+F2 land as one cnos β-skill patch (counted as 1), F3 lands as one cph CDR patch (counted as 1, application pending cph γ) — totals 3 patch artifacts across findings, 2 destination repos.

Cross-repo trace: this iteration produces a cnos β-skill patch from a cph cycle finding. Per `cdd/post-release/SKILL.md` §5.6b, the cph-side bundle at `cph:.cdd/iterations/cross-repo/cnos/cph-21-followup-beta-anchor/` should track outbound; the cnos-side mirror lives at `cnos:.cdd/iterations/cross-repo/cph/coherence-drift-sweep-followup-2026-05-18/` (this bundle).

---

## §5 Skill-gap candidate dispositions

The two `cdd-protocol-gap` candidates surfaced by the cph wave are dispositioned as F1 + F2 above (consolidated into a single Rule 6 patch in `cdd/beta/SKILL.md`). No additional skill-gap candidates surfaced from the close-out review.

---

## §6 Deferred outputs

None deferred from ε. F4 (master closes) is the only `next-MCA`, and cph γ is the named actor with a concrete first AC. F5 (branch merge) is operator-gated, named explicitly.

---

## §7 Cross-repo protocol observations forwarded to cnos#377

This cph→cnos iteration surfaced two more cross-repo protocol gaps for cnos#377 to fold in:

1. **ε output direction for cph→cnos.** ε's `cdd-iteration.md` for a cph cycle whose patches land in cnos has no canonical home — cph's own `.cdd/waves/<wave>/` is the natural place for the wave's iteration record, but the cnos-side patch trace lives in cnos's own iteration channels. The bundle at `cnos:.cdd/iterations/cross-repo/cph/<slug>/` carries both (the iteration content for cph γ to apply, and the patch trace for cnos's record), but the protocol does not codify this dual-purpose bundle shape.

2. **Hat collapse documentation requirement.** An actor playing γ in repo A and ε for repo B in the same session has no protocol-defined attribution form. The hat-collapse acknowledgment block at the top of this file is an ad hoc convention. cnos#377 should consider naming this as a required artifact for cross-repo ε work.

---

## §8 Next-MCA commitment

cph γ:
1. Apply `cdr-changelog-rule.patch` (in this bundle) to `cph:CDR.md` §"Changelog rule"
2. Verify F1+F2 cnos β-skill patch lands on cnos main (track cnos branch `claude/file-cnos-cdr-issue-fi9Ld` → main merge)
3. Close cph#21 with closure comment citing wave-closeout, this cdd-iteration.md, sub closeouts
4. Close cph#16 with closure comment citing cph#21 as the resolution path
5. Update `cph:.cdd/iterations/INDEX.md` with the row in §4 above
6. Optionally write the cph-side cross-repo outbound trace at `cph:.cdd/iterations/cross-repo/cnos/cph-21-followup-beta-anchor/LINEAGE.md` per §5.6b
7. Surface F5 (merge claude/review-repo-coherence-PNbjQ → main) to operator after F3 lands

ε's role for this wave concludes after applying F1+F2 cnos-side and emitting this bundle. cph γ holds the remaining workitems.
