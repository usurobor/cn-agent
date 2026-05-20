# α design notes — cycle #377

**Issue:** cnos#377
**Branch:** `cycle/377`
**Phase:** design (precedes build)
**Role:** α (with β-collapsed-on-α self-review at end)

This file is the design half of the design-and-build mode. It converges:
1. Directional cases (the discrete shapes a cross-repo coordination can take).
2. STATUS state machine (event vocabulary + transition graph + emitters + master/sub rule).
3. Bundle file set per directional case.
4. LINEAGE.md schema per directional case.
5. Bundle-state phase mapping (`open | converging | closed`) onto STATUS events.
6. Feedback-patch format.
7. Bundle archival rule (per side).
8. Hat-collapse attribution form.
9. Path-collision resolution.

Each subsection ends with retroactive validation against the 8 empirical anchors.

## 1. Directional cases

A cross-repo coordination unit has a **direction** (who originates the substantive work) and a **fan-out shape** (1:1 vs master/sub vs patch-series).

Three primary directional cases + one degenerate edge case:

- **(a) Inbound proposal** — source repo emits a proposal; cnos accepts/modifies/rejects and tracks `landed`. The substantive target work happens in cnos. Sub-shapes:
  - **(a.1) 1:1** — one source proposal → one cnos issue → one cnos PR/cycle. (`cn-sigma/agent-activate-skill`)
  - **(a.2) master+sub** — one source proposal → one cnos master issue → N cnos sub-issues landing independently. (`cph/bootstrap-cdr` aka `gait-support-paths/bootstrap-cdr`)

- **(b) Outbound iteration trace** — cnos cycle produces patches that land in another repo. The substantive target work happens elsewhere. Sub-shapes:
  - **(b.1) patch-series** — cnos authors N patches against one target repo, all landed before the cycle closes. (`tsc/cdd-supercycle`, 6 patches)

- **(c) Bilateral iteration** — same session/wave produces patches for both repos via hat-collapse (γ in repo A + ε in repo B). The substantive target work spans both repos. (`cph/coherence-drift-sweep-followup-2026-05-18`)
  - Sub-shape carries dual purpose: iteration content (for the source repo's ε to land) + patch trace (for the target repo's record). Bundle is asymmetric — the side that cannot write directly carries the "feedback" surface.

- **(d) Operator-pending bundle** — substantive content drafted for a target that does not yet exist (or for a target that cnos session cannot reach in this cycle). The bundle is content + LINEAGE only; STATUS is `drafted` (operator-pending) until the operator scaffolds the target or applies the patch. Three sub-shapes observed:
  - **(d.1) target-repo-doesn't-exist** — full file drafts staged in cnos for an operator-scaffolded target repo. (`cn-rho/bootstrap-2026-05-19`)
  - **(d.2) target-repo-exists-but-unreachable** — feedback patch for an existing repo that the current cnos session cannot write to (MCP scope or credential limits). (`cn-sigma/discipline-section-2026-05-19`)
  - **(d.3) proposal-as-issue-comment** — proposal shaped as a GitHub issue comment rather than a file patch. (`cph/issue-32-tightening-2026-05-19`)

Case (d) is **structurally degenerate** relative to (a)–(c) — it carries no STATUS state machine (no `submitted` from a counterpart-side γ; no `accepted` from cnos because cnos *is* the source side). Its lifecycle is **operator-applied → archive**. The protocol covers it as a named case with a simpler shape, not by forcing it into (a)/(b)/(c).

### Why these and not others

- Multi-hop (A → B → C) and cross-org are out of scope per issue body §"Out of scope".
- Inbound iteration trace (someone else's cycle producing cnos patches) is structurally the same as (c) bilateral — the actor is the source repo's ε. When that ε also wears cnos γ's hat, hat-collapse applies; when it doesn't, the cnos side simply receives the patch as if from any contributor and no cross-repo bundle is needed on the cnos side.
- "Outbound proposal" (cnos emits a proposal asking another repo to accept) is structurally case (b) — cnos is the source, the proposal becomes patches the target lands. The proposal-acceptance handshake mirrors (a) from the target's perspective.

## 2. STATUS state machine

### 2.1. Event vocabulary

Five canonical events:

| Event | Meaning |
|---|---|
| `submitted` | Source repo's γ (or another role) emits the proposal as ready for the target's intake. Target sees this in its intake scan. |
| `accepted` | Target's γ accepts the proposal substantially as proposed and files a target issue. |
| `modified` | Target's γ accepts the governing gap but changes scope, split, wording, implementation, proof, or patch application materially. Filed with a `Delta` field in the target issue's `## Source Proposal` block. |
| `rejected` | Target's γ declines the proposal. The target issue is not filed. |
| `landed` | Target work has merged. For 1:1 — one event per merge. For master/sub — one event per sub merge, plus one terminal `landed` event when the master closes; see §2.5. |

### 2.2. `drafted` reconciliation

The `cn-sigma/agent-activate-skill` STATUS opens with `2026-05-19 drafted sigma cn-sigma@8f153c15e` — `drafted` is not in the canonical 5-event vocabulary above. Two options:

- **(i) expand vocabulary to 6 events** — add `drafted` as a pre-`submitted` event meaning "authored, not yet emitted to target's intake".
- **(ii) name `drafted` as a synonym for `submitted` in source-posture cases** — accept that an agent hub (cn-sigma, not a CDD-activated tenant repo) may use `drafted` to mean "filing-ready".

**Decision: (ii)** — keep the vocabulary at five events per the issue's non-goal "Do NOT add a new bundle-state phase or STATUS event. Work with the existing five events + three phases; the work is reconciliation, not extension." `drafted` is named as a recognized synonym for `submitted` in source-posture cases (agent hubs, dormant branches, draft-only artifacts). The target's intake scan treats `drafted` as `submitted`-equivalent for the purpose of triage.

A target γ that finds `drafted` in source STATUS treats it as filing-ready and proceeds to `accepted | modified | rejected`. The target's STATUS mirror **normalizes** `drafted` to `submitted` so the mirror's vocabulary is canonical, but the source STATUS may retain `drafted`.

### 2.3. Transition graph

```
                          ┌────────────┐
                          │  (start)   │
                          └─────┬──────┘
                                │
              source γ emits    │
              (or `drafted`     │
              recognized        │
              equivalent)       ▼
                          ┌────────────┐
                          │  submitted │
                          └─────┬──────┘
                                │
                  target γ      │
                  triage        │
              ┌─────────────────┼─────────────────┐
              │                 │                 │
              ▼                 ▼                 ▼
        ┌──────────┐      ┌──────────┐      ┌──────────┐
        │ accepted │      │ modified │      │ rejected │
        └─────┬────┘      └─────┬────┘      └──────────┘
              │                 │              (terminal)
              │                 │
              │  target work    │
              │  merges         │
              └────────┬────────┘
                       │
                       ▼
                 ┌──────────┐
                 │  landed  │
                 └──────────┘
                  (terminal)
```

Legal transitions only:

- `(start) → submitted`
- `submitted → accepted | modified | rejected`
- `accepted → landed`
- `modified → landed`

Illegal:
- `rejected → *` (terminal)
- `landed → *` (terminal)
- `submitted → landed` (must pass through accepted/modified)
- `accepted → modified` (the disposition is decided once at intake; re-disposition requires a new proposal)

### 2.4. Emitters per event

| Event | Emitter | Notes |
|---|---|---|
| `submitted` (or `drafted` synonym) | source γ (or source role authoring the bundle) | Writes to source-side STATUS at filing time. |
| `accepted` / `modified` / `rejected` | target γ | Writes to source-side STATUS via direct push if cross-repo write is authorized, or via `FEEDBACK.patch` if not. Constraint: after target γ's filing decision, source STATUS must not remain at `submitted`. |
| `landed` (1:1) | target γ | Writes when the cycle implementing the target issue merges to target main. |
| `landed` (master/sub) | target γ per §2.5 below | One event per sub merge + one terminal master-closure event. |

The cnos-side **mirror** STATUS (at `cnos:.cdd/iterations/cross-repo/{source-repo}/{slug}/STATUS`) is target γ's record; it mirrors the source STATUS but may carry mirror-only annotations (e.g. `landed` events when the source-side cannot be written from the cycle session).

### 2.5. Master/sub rule for `landed`

For master/sub-shaped proposals:

1. **Per-sub `landed` event** fires each time a sub merges to target main, with the form `<ISO-date> landed gamma@<target> <target-sub-issue#> <merge-sha> <release-version>`. Source-side STATUS receives one row per sub merge (via direct write or feedback patch).

2. **Terminal master-close `landed` event** fires when the master issue is closed (i.e. all subs landed or deferred subs explicitly named as tracked debt in the master's closure condition). The master-close event reads `<ISO-date> landed gamma@<target> <target-master-issue#> master-close`.

**Rationale:** the per-sub events preserve granular trace for retroactive analysis (which sub landed when); the master-close event preserves the protocol invariant that `accepted/modified` must terminate at `landed` and is the single event the source γ may treat as "the work is fully done". A wave with 4 subs landing across 3 releases produces 5 `landed` rows: 4 per-sub + 1 master-close.

For 1:1 proposals: one `landed` event total. There is no separate "master close" because the issue is the only target unit.

### 2.6. Mapping to bundle-state phases

The `.cdd/iterations/cross-repo/README.md` defines three bundle-state phases: `open | converging | closed`. Mapping to STATUS events:

| Bundle phase | STATUS event(s) | Meaning |
|---|---|---|
| `open` | `submitted` (or `drafted`) | Filed; no target disposition yet. |
| `converging` | `accepted` or `modified` (no `landed` yet, or only some `landed` rows for master/sub) | Target accepted; work in flight. |
| `closed` | terminal `landed` (1:1) **or** master-close `landed` (master/sub) **or** `rejected` | No further work expected. |

A bundle in `converging` may carry partial `landed` events (per-sub) without transitioning to `closed`; transition fires only on terminal `landed` or `rejected`.

### 2.7. Retroactive validation

| Anchor | STATUS observed | Compliance | Notes |
|---|---|---|---|
| `tsc/cdd-supercycle` | (no STATUS file; LINEAGE-only outbound trace) | Case (b) — outbound; STATUS lives in source-side bundle. Compliant by exemption. | Outbound case § 3.5 below: STATUS lives source-side. |
| `cph/bootstrap-cdr` | `submitted gait` → `accepted gamma@cnos cnos#376` | Compliant. Per-sub `landed` rows to be appended as cnos.cdr v0.1 wave subs close. |  |
| `gait-support-paths/bootstrap-cdr` | (historical mirror; STATUS at `cph/bootstrap-cdr` is the canonical post-rename home) | Compliant. Pre-rename mirror preserved as audit artifact only. | Per `cph/bootstrap-cdr/LINEAGE.md` repository rename event. |
| `cn-sigma/agent-activate-skill` | `drafted sigma` → `accepted gamma@cnos cnos#379` → `modified gamma@cnos cnos#379 ac4-...` → `landed gamma@cnos cnos#379 a3bf7892 3.78.0` | **drafted** synonym applied (§2.2 reconciliation); transitions follow §2.3 except shows `accepted → modified → landed`. **Contradicts §2.3 illegal `accepted → modified`.** | See §2.3a below for resolution. |
| `cph/coherence-drift-sweep-followup-2026-05-18` | (no STATUS — case (c) bilateral; iteration content + patch trace, not proposal lifecycle) | Compliant by case-classification. | Case (c) has no STATUS state machine; it has a separate phase machine — see §3.4. |
| `cn-rho/bootstrap-2026-05-19` | (no STATUS; LINEAGE-only with disposition `drafted (operator-pending)`) | Compliant by case-classification. Case (d.1). | Case (d) has no STATUS state machine; LINEAGE carries disposition. |
| `cn-sigma/discipline-section-2026-05-19` | (no STATUS; LINEAGE-only) | Compliant. Case (d.2). |  |
| `cph/issue-32-tightening-2026-05-19` | (no STATUS; LINEAGE-only) | Compliant. Case (d.3). |  |

### 2.3a. Reconciling `accepted → modified` (cn-sigma/agent-activate-skill)

The cn-sigma anchor STATUS shows:
```
2026-05-19 submitted sigma cn-sigma@8f153c15e   (drafted-normalized)
2026-05-19 accepted gamma@cnos cnos#379
2026-05-19 modified gamma@cnos cnos#379 ac4-capability-matrix-foldin@bdda457f5
2026-05-19 landed gamma@cnos cnos#379 a3bf7892 3.78.0
```

Two readings:

**Reading A:** `accepted → modified` is a legal transition recording **post-filing refinement** of the disposition. The proposal was first accepted; then the issue body was edited to fold in upstream refinements; the second event records the delta. The transition graph in §2.3 should permit `accepted → modified`.

**Reading B:** The `modified` event here is documenting a Delta against the originally-accepted body, not changing the disposition. The transition is `accepted → accepted` with a delta annotation; the protocol vocabulary should permit a `delta` event instead of overloading `modified`.

**Decision: Reading A.** Permit `accepted → modified` as a legal transition, with the semantic "disposition refined after initial acceptance; Delta field of the target issue's Source Proposal block is updated to match". This:
- preserves the 5-event vocabulary (no new `delta` event needed; honors the issue's non-goal "Do NOT add a new STATUS event");
- matches the observed practice in cn-sigma/agent-activate-skill;
- does not break the terminal `landed` invariant — the chain still ends at `landed`.

Updated transition graph: `accepted → modified` is legal (and `modified → modified` if further refinement occurs); `modified → accepted` is **not** legal (cannot un-modify once delta is recorded).

The skill's STATUS state machine §2.3 graph is updated accordingly.

## 3. Bundle file set per directional case

### 3.1. Canonical path

**One canonical path:** `.cdd/iterations/cross-repo/{counterpart-repo}/{slug}/` on whichever side carries the bundle.

- For inbound (case a) on the cnos side: `cnos:.cdd/iterations/cross-repo/{source-repo}/{slug}/` (the source repo's name appears in the path; the cnos side carries the mirror).
- For outbound (case b) on the cnos side: `cnos:.cdd/iterations/cross-repo/{target-repo}/{slug}/` (the target repo's name appears in the path; the cnos side carries the source-side bundle until the target-side mirror lands).
- For bilateral (case c) on the cnos side: `cnos:.cdd/iterations/cross-repo/{counterpart-repo}/{slug}/` (whichever repo the hat-collapse counterpart was — cnos session staging surface for cross-repo deliverables).
- For operator-pending (case d) on the cnos side: `cnos:.cdd/iterations/cross-repo/{target-repo}/{slug}/` (the prospective or unreachable target).

The path is **direction-agnostic**: `{counterpart-repo}` is "the repo that is *not* cnos". This resolves the path-collision identified in the issue §"Where they diverge" #1 — intake-scan and outbound-trace now share the path shape; they are distinguished by **event semantics** (presence of source-vs-target events in STATUS) and **LINEAGE schema** (Source vs Target which-side-am-I block).

### 3.2. Bundle file set — case (a) inbound proposal (1:1 or master/sub)

Required files on the cnos side:

| File | Required? | Purpose |
|---|---|---|
| `LINEAGE.md` | yes | Bilateral trace — names source/target, disposition, delta, per-sub or per-patch confirmation. |
| `STATUS` | yes | Event ledger (mirror of source STATUS + cnos-side authoritative `landed` rows). |
| `FEEDBACK.patch` | conditional | Required when the cnos session cannot write to source-side STATUS directly. Carries `accepted/modified/rejected/landed` events as a unified diff against source STATUS. |
| `ISSUE.md` | optional | Verbatim copy of source `ISSUE.md` for audit. Useful when source bundle is dormant/branch-only. Not required because the target cnos issue body carries the substance. |
| `README.md` | **discouraged** | Per AC5: do not place protocol doctrine in bundle README. Bundle-specific narrative may live in `LINEAGE.md` instead. |

### 3.3. Bundle file set — case (b) outbound iteration trace

Required files on the cnos side (mirror, after source-side bundle is archived):

| File | Required? | Purpose |
|---|---|---|
| `LINEAGE.md` | yes | Bilateral trace — names source cycle, target patches, per-patch confirmation. |

Source-side bundle (under counterpart-repo's `.cdd/iterations/cross-repo/cnos/{slug}/`) carries the working files (patches, lineage, STATUS if a proposal lifecycle was attached). After target patches merge **and** target-side mirror confirms receipt, the source-side may be deleted (archival rule §6).

### 3.4. Bundle file set — case (c) bilateral iteration

Required files on the cnos side:

| File | Required? | Purpose |
|---|---|---|
| `LINEAGE.md` | yes | Bilateral trace + hat-collapse acknowledgment + per-patch trace per destination repo. |
| `cdd-iteration.md` | yes | The ε work product (per `post-release/SKILL.md §5.6b` and `epsilon/SKILL.md §1`) when the bilateral iteration carries ε output. |
| `<patch>.patch` | one per outbound patch | Feedback patches for the counterpart repo to apply (unified diff format). Named descriptively (e.g. `cdr-changelog-rule.patch`). |

Case (c) bundles are **dual-purpose**: they carry both the iteration content (for the counterpart repo's γ to apply) and the cnos-side patch trace (for cnos's record of patches that landed in cnos as part of the same iteration). The dual purpose is named explicitly in `LINEAGE.md §"Bilateral trace direction"`.

### 3.5. Bundle file set — case (d) operator-pending

Required files on the cnos side:

| File | Required? | Purpose |
|---|---|---|
| `LINEAGE.md` | yes | Names source trigger, reference doctrine, sibling exemplars, application gate, archival rule. Carries `Status: drafted (operator-pending)`. |
| Substantive content | varies | The drafted artifact(s): file content for case (d.1), `<file>.patch` for case (d.2), `issue-comment.md` for case (d.3). |

No `STATUS` ledger — there is no proposal lifecycle to track. LINEAGE.md `Disposition` field carries the operator-pending state.

### 3.6. Retroactive validation

| Anchor | Case | Observed bundle file set | LINEAGE schema | Compliance |
|---|---|---|---|---|
| `tsc/cdd-supercycle` | (b.1) | `LINEAGE.md` only | Source / Target / Bilateral trace / Per-patch confirmation | Compliant — matches §3.3. |
| `cph/bootstrap-cdr` | (a.2) | `LINEAGE.md`, `STATUS`, `FEEDBACK.patch` | Source / Target / Repository rename event / Bilateral trace / Per-sub confirmation / Related cross-repo wave context | Compliant — matches §3.2. |
| `gait-support-paths/bootstrap-cdr` | (a.2) historical | `LINEAGE.md`, `FEEDBACK.patch` (no STATUS — pre-rename mirror) | Source / Target / Bilateral trace / Per-sub confirmation | Compliant as **archived pre-rename mirror**. Active mirror is `cph/bootstrap-cdr`. |
| `cn-sigma/agent-activate-skill` | (a.1) | `LINEAGE.md`, `STATUS`, `FEEDBACK.patch` | Source / Target / Delta / Bilateral trace / Related cross-repo wave context | Compliant — matches §3.2. |
| `cph/coherence-drift-sweep-followup-2026-05-18` | (c) | `LINEAGE.md`, `cdd-iteration.md`, `cdr-changelog-rule.patch` | Source / ε actor / Target patches / Bundle contents / Bilateral trace direction / Protocol observations | Compliant — matches §3.4. |
| `cn-rho/bootstrap-2026-05-19` | (d.1) | `LINEAGE.md`, `PERSONA.md`, `OPERATOR.md` | Purpose / Source / Bundle contents / Operator-side scaffolding steps / Not in scope / Disposition | Compliant — matches §3.5. |
| `cn-sigma/discipline-section-2026-05-19` | (d.2) | `LINEAGE.md`, `PERSONA-discipline.patch` | Purpose / Source / Bundle contents / Application / Verification / Disposition | Compliant — matches §3.5. |
| `cph/issue-32-tightening-2026-05-19` | (d.3) | `LINEAGE.md`, `issue-comment.md` | Purpose / Source / Bundle contents / Application / Disposition | Compliant — matches §3.5. |

**Zero contradictions.** All 8 anchors validate without revision to the schema.

## 4. LINEAGE.md schema per case

### 4.1. Case (a) inbound proposal — 1:1 or master/sub

Required sections:
- `## Source` (repo, branch, commit, path)
- `## Target` (repo, issue #, mode, disposition, first filing commit/branch, canonical-path on cnos main)
- `## Delta` (only when disposition = `modified`; names what changed)
- `## Bilateral trace` (which side carries what; archival predicate)
- `## Per-sub confirmation` (master/sub only — table of subs, status, target issue numbers, commits) — **omit for 1:1**
- `## Per-patch confirmation` (1:1 with multi-patch) — table of patches, target commits, status
- `## Related cross-repo wave context` (optional — links to wave artifacts, repository rename events)

### 4.2. Case (b) outbound iteration trace

Required sections:
- `## Source` (cnos cycle / branch / commit / path) — cnos is the source here
- `## Target` (counterpart repo, target issue, target PR/commit, patches landed)
- `## Bilateral trace` (mirror status)
- `## Per-patch confirmation` (table — patch number, target commit, status)

### 4.3. Case (c) bilateral iteration

Required sections:
- `## Source` (counterpart repo, branch, wave, master issue, sub issues, wave verdict)
- `## ε actor` (who played ε; authority basis; date) — hat-collapse acknowledged here
- `## Target patches` (per-patch table per destination repo)
- `## Bundle contents` (file-by-file purpose)
- `## Bilateral trace direction` (why this is on this side; what the counterpart side carries)
- `## Protocol observations forwarded` (optional — when the iteration surfaces protocol gaps for codification)

### 4.4. Case (d) operator-pending

Required sections:
- `## Purpose` (one-paragraph what-this-stages)
- `## Source` (trigger, reference doctrine, sibling exemplar, empirical anchor)
- `## Bundle contents` (file-by-file purpose)
- `## Operator-side scaffolding steps` (case d.1) OR `## Application` (case d.2/d.3)
- `## Verification after application` (case d.2 specifically — how operator confirms the patch landed)
- `## What is NOT in scope here` / `## What this bundle does NOT do` (optional but recommended)
- `## Disposition` (always `drafted (operator-pending)` for case d)

### 4.5. Retroactive validation

(Already covered in §3.6 — every anchor's LINEAGE.md carries the sections its case requires. Specifically:
- `tsc/cdd-supercycle` LINEAGE has Source / Target / Bilateral trace / Per-patch confirmation — matches §4.2.
- `cph/bootstrap-cdr` LINEAGE has all §4.1 master/sub sections including the optional repository-rename-event subsection.
- `cn-sigma/agent-activate-skill` LINEAGE has all §4.1 sections including `## Delta` (because disposition was `modified` post-filing).
- `cph/coherence-drift-sweep-followup-2026-05-18` LINEAGE has all §4.3 sections including the optional `## Protocol observations forwarded`.
- `cn-rho/bootstrap-2026-05-19` LINEAGE has all §4.4 sections; matches case (d.1).
- `cn-sigma/discipline-section-2026-05-19` LINEAGE has all §4.4 sections including `## Verification after application`; matches case (d.2).
- `cph/issue-32-tightening-2026-05-19` LINEAGE has all §4.4 sections; matches case (d.3).

Zero contradictions across all 8 anchors.

## 5. Feedback-patch format

A feedback patch is a unified diff against the counterpart-side STATUS file (or other counterpart-side file), wrapped in a header that names:
- `From:` actor identity (cnos role + branch)
- `Date:` ISO date
- `Subject:` one-line summary

The header is followed by free-prose context (apply command, rationale, links to related artifacts), then a `---` separator, then the unified diff.

Apply command convention: `git apply <patch-filename>` from the counterpart repo root.

Patch filename convention: descriptive lowercase-with-hyphens (e.g. `FEEDBACK.patch` for STATUS-event patches, `cdr-changelog-rule.patch` for substantive doctrine patches, `PERSONA-discipline.patch` for persona-file patches). When the patch's purpose is STATUS-event emission only, the canonical name is `FEEDBACK.patch`. When it carries substantive content, the name describes the content.

### 5.1. Retroactive validation

| Anchor | Patch file | Format compliance |
|---|---|---|
| `cph/bootstrap-cdr/FEEDBACK.patch` | header + apply command + diff against source STATUS | Compliant. |
| `gait-support-paths/bootstrap-cdr/FEEDBACK.patch` | identical to above | Compliant. |
| `cn-sigma/agent-activate-skill/FEEDBACK.patch` | header + apply command + diff against source STATUS + protocol-observation note on `drafted` event | Compliant. |
| `cph/coherence-drift-sweep-followup-2026-05-18/cdr-changelog-rule.patch` | header + apply command + diff against `CDR.md` | Compliant. |
| `cn-sigma/discipline-section-2026-05-19/PERSONA-discipline.patch` | header + apply command + diff against `spec/PERSONA.md` | Compliant. |

All 5 observed feedback patches comply. Zero contradictions.

## 6. Bundle archival rule

A bundle's archival rule depends on direction and side:

### 6.1. Source-side may be deleted when

For case (a) inbound: source-side bundle may be deleted by source γ once the cnos-side mirror exists on cnos main **and** the cycle implementing the target work has closed (`landed` event recorded). Per `cph/bootstrap-cdr/LINEAGE.md`: "once this cnos-side mirror exists on cnos main + the cnos.cdr v0.1 wave closes, the source-side bundle may be archived."

For case (b) outbound: source-side (cnos-side until the counterpart's mirror lands) bundle may be deleted by cnos γ once the counterpart-side mirror confirms receipt of all patches. Per existing `post-release/SKILL.md §5.6b` text: "The cross-repo directory persists until the target PR merges; γ may delete it thereafter (the lineage is preserved in the target repo's own `cdd-iteration.md`)."

For case (c) bilateral: counterpart-side bundle (the side that cannot write directly) may be deleted by the cnos γ once the counterpart applies the patches and emits its own iteration trace. Per `cph/coherence-drift-sweep-followup-2026-05-18/LINEAGE.md`: "cph γ may delete the cph-side outbound bundle once cnos main carries the F1+F2 patch and cph#21 closes."

For case (d) operator-pending: bundle may be deleted by cnos γ once the operator effects the application (target repo created with files / patch applied / comment posted).

### 6.2. Target-side may be deleted

Mirror-side bundles (the one on `cnos:.cdd/iterations/cross-repo/{counterpart}/{slug}/` when cnos is the **target**) are **preserved indefinitely** as audit artifacts and retroactive validators. Per issue §"Do NOT delete any existing cross-repo bundle. Bundles persist as audit artifacts and retroactive validators."

Source-side bundles (cnos as **source**, case b/c) may be deleted post-archival per §6.1.

This asymmetry is intentional: the target-side mirror is the durable record; the source-side bundle is a transient working surface.

### 6.3. Retroactive validation

All 8 anchors are currently preserved. The archival rule is invoked only when source γ decides to archive; none have. The rule does not contradict observed practice.

## 7. Hat-collapse attribution form

When a single actor plays γ in repo A and ε (or another role) in repo B during the same session, the attribution form is:

1. **In the bundle LINEAGE.md**, a `## ε actor` (or equivalent role-name) section names:
   - `Played by:` (actor identity + originating session branch)
   - `Authority basis:` (citing `epsilon/SKILL.md §2` "ε and δ are often the same actor" + `ROLES.md §4` role-collapse rules + any operator-specific delegation)
   - `Output date:` ISO date

2. **In the cdd-iteration.md** (when case (c) produces ε output), a top-level `## Hat-collapse acknowledgment` section names:
   - which two roles are collapsed in this session
   - the authority basis citation
   - an explicit statement that the substantive role-A and role-B work are independent
   - the reason hat-collapse does not collapse with α/β/γ responsibilities for any in-flight cycle

This is exactly the shape observed in `cph/coherence-drift-sweep-followup-2026-05-18/LINEAGE.md §"ε actor"` and `cdd-iteration.md §"Hat-collapse acknowledgment"`. The protocol codifies the observed shape.

**Boundary preservation:** the attribution form names role-collapse as factual record, not as authority override. Per `ROLES.md §4a` persona/protocol/project boundary, the **personas** stay distinct (cnos γ identity ≠ cph ε identity); the actor wearing both hats acknowledges the collapse and carries the substantive work for each role independently. The skill's hat-collapse section does **not** break the persona/protocol/project boundary in `ROLES.md §4a` because it documents an actor-level fact (one human/agent ran two sessions of work) without merging the role-level identities or authorities.

### 7.1. Retroactive validation

Only one observed anchor — `cph/coherence-drift-sweep-followup-2026-05-18` — uses hat-collapse. Both LINEAGE.md and cdd-iteration.md carry the required form. Compliant.

## 8. Path-collision resolution

Per §3.1: `.cdd/iterations/cross-repo/{counterpart-repo}/{slug}/` is the **single canonical path** for cross-repo bundles, direction-agnostic.

The gamma §2.1 intake-scan needs updating: instead of scanning `.cdd/iterations/proposals/*/STATUS` or `.cdd/proposals/{target}/*/STATUS`, γ scans source-side counterpart-repos at `.cdd/iterations/cross-repo/cnos/*/STATUS` for active proposal lifecycle files.

Specifically: a source repo X authors a proposal *for cnos*; the source-side bundle lives at `X:.cdd/iterations/cross-repo/cnos/{slug}/`. The cnos-side mirror, once filed, lives at `cnos:.cdd/iterations/cross-repo/X/{slug}/`. The {counterpart} name in the path is "the repo that is not me".

After this protocol lands, the gamma §2.1 doctrine reads:

> Cross-repo proposal intake (see `cdd/cross-repo/SKILL.md §<intake>`). Scan source repos at `.cdd/iterations/cross-repo/cnos/*/STATUS` (the source-side path naming cnos as counterpart) for the last non-comment event = `submitted` (or `drafted`-synonym per `cdd/cross-repo/SKILL.md §STATUS state machine`). For each, decide accepted/modified/rejected per the protocol.

The legacy paths `.cdd/iterations/proposals/*/` and `.cdd/proposals/{target}/*/` are dropped — they did not match observed practice (`cph/bootstrap-cdr` shipped at the canonical path). The legacy paths are removed from gamma §2.1; they are not retained as deprecation references because the issue body §"Out of scope" excludes retroactive cleanup and the legacy paths are not anchored to any in-flight bundle.

## 9. Design summary

The protocol surface to encode in `cdd/cross-repo/SKILL.md`:

1. **Coherence formula:** Coherent cross-repo coordination resolves a single directional case at a single canonical path, with a known event vocabulary terminating at `landed` (or `rejected`), and a bundle file set + LINEAGE schema selected by the case.
2. **Failure mode:** Cross-repo coordination fails through **directional ambiguity** — γ invents the bundle shape per cycle because the protocol doesn't name the directional case.
3. **Domain:** Cross-repo coordination between cnos and a single counterpart repo (same-org, bilateral, no multi-hop).
4. **Governing question:** How does γ resolve a cross-repo event (proposal intake / iteration outbound / bilateral / operator-pending) into a single canonical bundle with one path, one STATUS vocabulary, and one schema?
5. **Parts:** directional case + canonical path + STATUS state machine + bundle file set + LINEAGE schema + feedback-patch format + bundle archival rule + hat-collapse attribution form.

### Section order (the skill's Unfold)

1. Define (parts, formula, failure mode).
2. Directional cases (a–d with sub-shapes).
3. Canonical path.
4. STATUS state machine (events, transitions, emitters, master/sub `landed` rule, `drafted` reconciliation).
5. Bundle-state phases + mapping to STATUS.
6. Bundle file set (per case).
7. LINEAGE.md schema (per case).
8. Feedback-patch format.
9. Bundle archival rule.
10. Hat-collapse attribution.
11. Known protocol edge cases (the observations carried for future codification).
12. Cross-references (to gamma, post-release, issue skills).
13. Kata.

Total expected length: ~600–800 lines (matches issue's projection).

## 10. Build sequence

After this design phase:
1. Author `src/packages/cnos.cdd/skills/cdd/cross-repo/SKILL.md` per the design above.
2. Update `gamma/SKILL.md §2.1 + §2.7` to reference the new skill; remove inline directional doctrine.
3. Update `post-release/SKILL.md §5.6b` to reference the new skill for the outbound case; keep the `cdd-iteration.md` text since that's post-release's own concern.
4. Update `.cdd/iterations/cross-repo/README.md` to cross-reference the new skill and align bundle-state phases with the canonical state machine.
5. Author `self-coherence.md` per `design/SKILL.md §3.1` output format.
6. Author `beta-review.md` (β-collapsed-on-α self-review) verifying AC1–AC6.
7. Author three closeouts (`alpha-closeout.md`, `beta-closeout.md`, `gamma-closeout.md`).
8. Merge to main.

End of design phase.
