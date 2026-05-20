# β review — cycle #377

**Issue:** #377 — design(cdd): codify cross-repo coordination protocol
**Branch:** `cycle/377`
**Reviewer:** β-collapsed-on-α self-review (acknowledged; α≠β within session is structurally compromised but accepted under §5.2 wave-mode role-collapse precedent — manifest invariant + role-collapse note in self-coherence.md).

## Round 1

### Mechanical AC oracle results

| AC | Oracle | Result |
|---|---|---|
| AC1 | frontmatter (`name`, `description`, `artifact_class`, `governing_question`, `visibility: internal`, `parent: cdd`, `triggers`, `scope`); ≥2 references in `gamma/` + `post-release/` | **PASS** — frontmatter complete; 3 references (2 in gamma + 1 in post-release) |
| AC2 | STATUS state machine: vocabulary + transition graph + emitters + master/sub `landed` rule | **RC** — vocabulary contradicts CDD.md (see B-1) |
| AC3 | Bundle file set + LINEAGE schema per case; both precedents validate retroactively without contradiction | **PASS** — design-notes §3.6, §4.5 show zero contradictions across all 8 anchors |
| AC4 | One canonical path per directional case; intake-scan in `gamma/SKILL.md §2.1` updated to match | **RC** — CDD.md still names the legacy paths as first-class (see B-2) |
| AC5 | Feedback-patch format + bundle archival rule codified; doctrine removed from bundle READMEs | **PASS** — format + rule in skill §2.7 + §2.8; bundle READMEs do not exist in any current bundle (no doctrine inversion) |
| AC6 | Existing fragments cite the new canonical surface; `issue/SKILL.md` Source Proposal block stays | **PASS** — 3 references; `issue/SKILL.md` unchanged (verified by `git diff origin/main`) |

**Verdict R1: REQUEST CHANGES** — 2 binding findings.

### Binding findings

#### B-1: STATUS vocabulary contradicts `CDD.md` §"Cross-repo proposal lifecycle"

**Severity:** binding (AC2 fails).
**Class:** wiring / honest-claim hybrid — the new skill claims a 5-event vocabulary but CDD.md (the canonical authority) ships an 8-event vocabulary.

**Evidence:**

`src/packages/cnos.cdd/skills/cdd/CDD.md` lines 234–243 names the cross-repo STATUS vocabulary:

```
- `drafted` - source has written the proposal but has not requested target action.
- `submitted` - source requests target intake. This is the only event required for target intake.
- `accepted` - target will act substantially as proposed and has a target reference.
- `modified` - target accepts the gap but changes scope, split, wording, implementation, proof, or patch application materially.
- `landed` - target work merged or otherwise became target truth.
- `rejected` - target declines the proposal as target work.
- `withdrawn` - source retracts the request.
- `revised` / `corrected` - optional audit events for post-submission revisions or corrections.
```

That is **8 events** (drafted, submitted, accepted, modified, landed, rejected, withdrawn, revised/corrected), not the 5 the new skill codifies.

Per `cdd/SKILL.md` §"Conflict rule": "If a role skill and `CDD.md` disagree on ordered steps, selection rules, or artifact contract, `CDD.md` governs." Vocabulary is artifact contract. CDD.md governs.

The new skill's `drafted`-as-synonym reconciliation is **factually wrong against the canonical source**: `drafted` is a distinct event in CDD.md meaning "source has written the proposal but has not requested target action" — semantically different from `submitted` ("source requests target intake"). Treating them as synonyms collapses two states the canonical source separates.

The issue body's count of "five events" was an under-read of CDD.md by the issue author. The new skill cannot ratify the under-read; it must align with the canonical source.

**Required α R2 fix:**

1. Update `cdd/cross-repo/SKILL.md §2.3.1` to enumerate all 8 events with the meanings from CDD.md (verbatim where possible).
2. Remove §2.3.2 `drafted` reconciliation as synonym; replace with the CDD.md semantics (`drafted` is a distinct pre-`submitted` event).
3. Update §2.3.3 transition graph to include `drafted → submitted` (or `drafted → withdrawn`) and `* → withdrawn`, plus `* → revised|corrected` as optional audit events.
4. Update §2.3.4 emitters table for the new events.
5. Re-verify the 8 anchors against the expanded vocabulary. Specifically: cn-sigma's `drafted sigma cn-tau@8f153c15e` is a **true `drafted` event**, not a `submitted` synonym; the cn-sigma γ filing without an intermediate `submitted` event is permitted because `drafted → accepted` is legal per CDD.md (`submitted` is only required for target intake, but cn-sigma also delegated authority to cnos γ to act on the `drafted` content). Document the transition.

#### B-2: Legacy paths still named as first-class in CDD.md; new skill claims AC4 resolution

**Severity:** binding (AC4 fails).
**Class:** wiring — the new skill says "canonical path is `.cdd/iterations/cross-repo/{counterpart-repo}/{slug}/`" but `CDD.md` §"Cross-repo proposal lifecycle" lines 215–225 still names `.cdd/iterations/proposals/{slug}/` and `.cdd/proposals/{target}/{slug}/` as "two source-side layouts [that] are first-class".

**Evidence:**

`rg "\.cdd/iterations/proposals|\.cdd/proposals/\{target\}" src/packages/cnos.cdd/skills/cdd/`:
- `CDD.md` lines 216, 221 — first-class layouts (positive references)
- `cross-repo/SKILL.md` rules 2.2.a + 3.2 — deprecation references (negative examples)

The new skill ❌-marks paths that the canonical source ✅-marks. Same-cycle contradiction.

**Required α R2 fix:**

1. Update `CDD.md` §"Cross-repo proposal lifecycle" to name the **canonical** source-side path as `.cdd/iterations/cross-repo/{counterpart-repo}/{slug}/` (matching the cnos-side mirror path).
2. Reconcile the two legacy paths: either name them as deprecation references with a one-line note, or remove them. Recommend: remove the two legacy-path blocks; the cross-repo skill carries the canonical path. CDD.md retains the high-level pointer to the cross-repo skill for the full doctrine.
3. Update CDD.md to reference `cdd/cross-repo/SKILL.md` as the canonical surface for the protocol, following the `cdd/SKILL.md` doctrine-vs-skill rule: CDD.md keeps the cross-repo coordination fact (yes, this exists; here's the high-level vocabulary; here's the canonical surface for details); the skill carries the protocol mechanics.

### Non-binding findings

#### N-1: Master/sub `landed` rule is well-defined but the master-close event format lacks a `<release-version>` field

The per-sub `landed` event includes `<release-version>` (e.g. `landed gamma@cnos cnos#379 a3bf7892 3.78.0`). The master-close event has only `master-close` as the marker. Suggest: master-close events also carry a release-version when applicable (e.g. `master-close 3.85.0` if the master closes within a specific release; `master-close` alone for cross-release closures). Non-binding because the absence is not contradicted by any anchor.

#### N-2: `## Source Proposal` block placeholders convention not surfaced in this skill

The protocol observation from the issue body — "Source `## Source Proposal` block with placeholders (Source commit: filled at filing; Disposition: pending) is cleaner than inserting a fresh block target-side — adopt as canonical convention" — was forwarded to cnos#377 by `cn-sigma/agent-activate-skill` LINEAGE. The new skill does not surface this convention. The block lives in `cdd/issue/SKILL.md` (correctly out-of-scope for editing); but the new skill should **cross-reference** the issue skill's block as the target-side integration form, with a note that pre-populating placeholders before filing reduces target-side cycle work. Non-binding because AC6 explicitly says the issue block stays as-is and the convention is incremental polish.

#### N-3: Bundle-state phase `drafted (operator-pending)` is a synonym for `open` for case (d)

The README mapping (`.cdd/iterations/cross-repo/README.md`) and the new skill §2.4 both treat `drafted (operator-pending)` as a recognized phase-synonym for `open`. This is semantically correct but the new skill should make the synonym explicit in §2.4 rather than only naming it in the case-d branch. Non-binding because the README handles it.

### Cross-cycle observations

- **No file conflicts with #375 or #378.** This cycle touches `gamma/SKILL.md §2.1 + §2.7` (different sections than #375's §2.5), `post-release/SKILL.md §5.6b`, `CDD.md §"Cross-repo proposal lifecycle"` (R2 update), `.cdd/iterations/cross-repo/README.md`, and the new skill file. No overlap with #378's `review/`, `alpha/`, `operator/` files. Wave manifest §"Cross-cycle coordination" risk is unchanged.
- **CDD.md update is within wave invariant 1.** CDD.md is at `src/packages/cnos.cdd/skills/cdd/CDD.md`, under `src/packages/cnos.cdd/skills/cdd/` — explicitly named in wave invariant 1's allowed paths.

### Round 1 verdict

**REQUEST CHANGES** — α R2 required to address B-1 (STATUS vocabulary 5 → 8 events per CDD.md canonical source) and B-2 (legacy paths still first-class in CDD.md).

α R2 should NOT touch the case-classification, bundle-file-set, LINEAGE schema, feedback-patch format, archival rule, hat-collapse attribution, or known-edge-cases sections — those passed AC3, AC5, AC6.

Non-binding findings N-1, N-2, N-3 may be addressed in R2 at α's discretion or carried as known debt.

---

## Round 2

(Pending α R2.)
