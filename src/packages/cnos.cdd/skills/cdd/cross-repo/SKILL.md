---
name: cross-repo
description: Cross-repo CDD coordination. Use when a cnos cycle interacts with another repo's CDD lifecycle — accepting an inbound proposal, emitting an outbound iteration trace, running a bilateral iteration across hats, or staging operator-pending bundles for a counterpart repo.
artifact_class: skill
kata_surface: embedded
governing_question: How does γ resolve a cross-repo event into a single canonical bundle with one path, one STATUS vocabulary, one schema, and one archival rule — without rediscovering the protocol each cycle?
visibility: internal
parent: cdd
triggers:
  - cross-repo
  - proposal intake
  - outbound iteration
  - bilateral iteration
  - feedback patch
scope: task-local
inputs:
  - cross-repo event (inbound proposal / outbound patches / bilateral iteration / operator-pending draft)
  - counterpart-repo identity
  - source or target cycle state
outputs:
  - cross-repo bundle at the canonical path
  - LINEAGE.md per directional case
  - STATUS ledger when proposal lifecycle applies
  - feedback patch when direct write to counterpart is not authorized
  - bundle archival decision when terminal events fire
requires:
  - active role is γ (or ε per `epsilon/SKILL.md §1`)
  - canonical `CDD.md` loaded
  - `gamma/SKILL.md` loaded for intake/close-out steps that invoke this skill
calls:
  - gamma/SKILL.md
  - post-release/SKILL.md
  - issue/SKILL.md
---

# Cross-Repo Coordination

## Core Principle

**Coherent cross-repo coordination resolves a single directional case at a single canonical path, with a known event vocabulary terminating at `landed` (or `rejected`), and a bundle file set + LINEAGE schema selected by the case.**

Cross-repo coordination is the protocol for events that span two CDD-activated repos (or one CDD-activated repo and one agent-hub or unscaffolded counterpart). It covers same-org bilateral coordination only. Multi-hop (A → B → C) and cross-org are not in scope.

The failure mode is **directional ambiguity** — γ invents the bundle shape per cycle because the protocol does not name the directional case. The result is bundles that drift in shape, STATUS ledgers that go stale at `accepted`, and doctrine that lives inside individual bundles rather than in the protocol.

This skill is the **canonical surface** for the protocol. Other CDD skills (`gamma/`, `post-release/`, `issue/`) reference it; they do not carry inline cross-repo doctrine.

## Authority

- `CDD.md` governs the cycle lifecycle and selection rules; this skill governs what happens inside a single cycle's cross-repo interaction.
- This skill is the only canonical home for the cross-repo state machine, bundle file set, LINEAGE schema, feedback-patch format, archival rule, and hat-collapse attribution. If another skill carries the same doctrine, this skill governs.
- The `## Source Proposal` block in `cdd/issue/SKILL.md` minimal output pattern is the target-side issue-body integration; that fragment stays in `issue/`.

## Scope

In scope:
- Bilateral, same-org cross-repo coordination between cnos and one counterpart repo.
- Four directional cases: inbound proposal (1:1 or master/sub), outbound iteration trace, bilateral iteration, operator-pending bundle.
- STATUS state machine, bundle file set, LINEAGE.md schema, feedback-patch format, bundle archival rule, hat-collapse attribution.

Out of scope:
- Multi-hop coordination (A → B → C). When the next cycle hits multi-hop, this skill extends; until then, the protocol covers bilateral only.
- Cross-org coordination (different GitHub orgs).
- Tooling that auto-applies feedback patches, auto-emits `landed` events, or auto-archives bundles. Doctrine first; tooling deferred.

---

## 1. Define

### 1.1. The parts of a cross-repo coordination

A coherent coordination unit has these parts:

- **Direction** — who originates the substantive work (source) and who lands it (target).
- **Fan-out shape** — 1:1, master/sub, or patch-series.
- **Bundle** — the on-disk artifacts that carry the coordination record.
- **STATUS ledger** — the event sequence (proposal lifecycle only; cases a, c when proposal-shaped).
- **LINEAGE.md** — the bilateral trace; the durable record of who did what when.
- **Feedback patch** — when the cnos session cannot write to the counterpart side, the patch carries STATUS events or substantive content for the counterpart γ to apply on receipt.
- **Bundle archival rule** — when each side may delete its copy.

### 1.2. How they fit

Direction selects the case. Case selects the bundle file set. STATUS (when present) tracks proposal lifecycle. LINEAGE records the bilateral trace. Feedback patch is the cross-side write channel when direct write is unavailable. Archival rule names when each side may delete.

A bundle without LINEAGE is uninterpretable. A proposal-lifecycle case without STATUS is unanchored. A case without a named directional reading produces shape drift.

### 1.3. Failure mode

Cross-repo coordination fails through **directional ambiguity**:

- γ scans the wrong path and misses an active proposal.
- γ invents a LINEAGE schema by adapting a precedent shaped for a different case.
- γ leaves a STATUS at `accepted` indefinitely because the `landed` emitter is unclear.
- γ places protocol doctrine inside the bundle README, where it archives with the bundle.

This skill closes each failure mode by naming the directional case explicitly and binding bundle shape, STATUS rules, and archival to it.

---

## 2. Unfold

### 2.1. Directional cases

A cross-repo coordination unit fits **exactly one** of four directional cases. γ identifies the case before authoring the bundle.

#### Case (a) — Inbound proposal

Counterpart repo emits a proposal; cnos accepts, modifies, or rejects; substantive target work happens in cnos. Sub-shapes:

- **(a.1) 1:1** — one source proposal → one cnos issue → one cnos cycle. Example: `cn-sigma/agent-activate-skill` → cnos#379.
- **(a.2) master+sub** — one source proposal → one cnos master issue → N cnos sub-issues landing independently. Example: `cph/bootstrap-cdr` → cnos#376 + subs.

cnos side is the **target**. The cnos-side bundle is the **mirror**; the source-side bundle is the originating artifact.

#### Case (b) — Outbound iteration trace

A cnos cycle produces patches that land in another repo. The substantive target work happens elsewhere. Sub-shape:

- **(b.1) patch-series** — cnos authors N patches against one target repo, all landed before the cycle closes. Example: `tsc/cdd-supercycle` (6 patches → cnos#331/#332).

cnos side is the **source**. The cnos-side bundle is the source-side staging surface; the target-side bundle (when authored) is the mirror.

#### Case (c) — Bilateral iteration

A single session produces patches for both repos via hat-collapse (an actor playing γ in one repo and ε — or another role — in the other during the same session). The substantive target work spans both repos. Example: `cph/coherence-drift-sweep-followup-2026-05-18` (ε work for cph; patches landing in cnos and cph).

Case (c) bundles are **dual-purpose**: they carry both the iteration content (for the counterpart repo's γ to apply) and the cnos-side patch trace (for cnos's record of patches that landed in cnos as part of the same iteration).

#### Case (d) — Operator-pending bundle

Substantive content drafted for a target the current cnos session cannot reach. Three sub-shapes observed:

- **(d.1) target-repo-doesn't-exist** — full file drafts staged in cnos for a repo the operator will scaffold. Example: `cn-rho/bootstrap-2026-05-19` (drafts for `usurobor/cn-rho` before the repo exists).
- **(d.2) target-repo-exists-but-unreachable** — feedback patch for an existing repo that the current session's MCP / credential scope cannot write to. Example: `cn-sigma/discipline-section-2026-05-19`.
- **(d.3) proposal-as-issue-comment** — proposal shaped as a GitHub issue comment rather than a file patch (e.g. retargeting an existing issue). Example: `cph/issue-32-tightening-2026-05-19`.

Case (d) is **structurally degenerate** relative to (a)–(c): it carries no STATUS state machine (there is no source-side γ on the other end emitting `submitted`). Its lifecycle is **operator-applied → archive**. The protocol covers it as a named case with a simpler shape, not by forcing it into the proposal-lifecycle frame.

#### Rule 2.1.a — Identify the case before authoring the bundle

- ❌ Begin authoring a bundle "following the precedent" without naming which case the cycle is in.
- ✅ Name the case (a.1 / a.2 / b / c / d.1 / d.2 / d.3) in the first paragraph of `LINEAGE.md` so the schema selection is auditable.

### 2.2. Canonical path

**One canonical path:** `.cdd/iterations/cross-repo/{counterpart-repo}/{slug}/` on whichever side carries the bundle.

`{counterpart-repo}` is "the repo that is *not* me". On cnos:
- For case (a) inbound: `cnos:.cdd/iterations/cross-repo/{source-repo}/{slug}/` (source name appears).
- For case (b) outbound: `cnos:.cdd/iterations/cross-repo/{target-repo}/{slug}/` (target name appears).
- For case (c) bilateral: `cnos:.cdd/iterations/cross-repo/{counterpart-repo}/{slug}/` (the not-cnos repo).
- For case (d) operator-pending: `cnos:.cdd/iterations/cross-repo/{prospective-or-unreachable-target}/{slug}/`.

The corresponding source-side path on the counterpart repo (when the counterpart is CDD-activated) is `{counterpart}:.cdd/iterations/cross-repo/cnos/{slug}/`.

This path resolves the historical collision between "intake scan" and "outbound trace" — both directions share the path shape; they are distinguished by **event semantics** (presence of source-vs-target events in STATUS) and **LINEAGE schema** (Source vs Target section content).

#### Rule 2.2.a — Use exactly one canonical path

- ❌ Stage an inbound proposal mirror at `.cdd/proposals/{target}/{slug}/` or `.cdd/iterations/proposals/{slug}/`.
- ✅ Stage all cross-repo bundles at `.cdd/iterations/cross-repo/{counterpart-repo}/{slug}/`, direction-agnostic.

### 2.3. STATUS state machine

A STATUS ledger applies to cases (a) and (c) when the case carries proposal-shaped lifecycle. Cases (b) and (d) do not carry STATUS (case b's trace is patch-confirmation in LINEAGE; case d's disposition is recorded in LINEAGE).

The STATUS vocabulary is **canonical in `cnos.cds/skills/cds/CDS.md` §"Coordination surfaces" → §"Cross-repo proposals"**. This section codifies the transition graph, emitter rules, master/sub landed rule, and bundle-state phase mapping that the CDS vocabulary list does not enumerate.

#### 2.3.1. Event vocabulary (per CDS §"Coordination surfaces" → §"Cross-repo proposals")

Eight events:

| Event | Meaning |
|---|---|
| `drafted` | Source has written the proposal but has not requested target action. Pre-intake. |
| `submitted` | Source requests target intake. This is the only event required for target intake. |
| `accepted` | Target γ will act substantially as proposed and has filed a target reference. |
| `modified` | Target γ accepts the governing gap but changes scope, split, wording, implementation, proof, or patch application materially. The target issue's `## Source Proposal` block carries a `Delta` field. May also fire post-`accepted` to record a refinement (e.g. issue body edit folding in upstream changes). |
| `landed` | Target work merged or otherwise became target truth. For 1:1 — one event. For master/sub — one per sub merge + one terminal master-close (see §2.3.4). |
| `rejected` | Target γ declines the proposal. Terminal. |
| `withdrawn` | Source retracts the request. Terminal. |
| `revised` / `corrected` | Optional audit events for post-submission revisions or corrections. Do not rewrite old events after sharing — append a `corrected` event instead. |

#### 2.3.2. Transition graph

```
                ┌──────────────┐
                │   (start)    │
                └──────┬───────┘
                       │ source role
                       │ authors bundle
                       ▼
                ┌──────────────┐  source retracts
                │   drafted    │──────────────────┐
                └──────┬───────┘                  │
                       │                          │
                       │ source γ                 │
                       │ requests intake          │
                       ▼                          │
                ┌──────────────┐                  │
                │  submitted   │──────────────────┤
                └──────┬───────┘                  │
                       │                          │
                       │ target γ                 │
                       │ disposition              │
       ┌───────────────┼──────────────────┐       │
       ▼               ▼                  ▼       │
  ┌─────────┐    ┌──────────┐       ┌──────────┐  │
  │accepted │◀──▶│ modified │       │ rejected │  │
  └────┬────┘    └────┬─────┘       └──────────┘  │
       │              │              (terminal)   │
       │              │                           │
       │ target       │                           ▼
       │ work         │                     ┌──────────┐
       │ merges       │                     │withdrawn │
       └──────┬───────┘                     └──────────┘
              │                              (terminal)
              ▼
         ┌────────┐
         │ landed │
         └────────┘
          (terminal — per §2.3.4 master/sub)

  *  → revised | corrected   (optional audit events; do not change state;
                              may fire from any non-terminal state)
```

Legal transitions:
- `(start) → drafted | submitted` (source authors directly into either state)
- `drafted → submitted` (source γ requests intake)
- `drafted → accepted | modified | rejected` (permitted when source explicitly delegates filing-authority to target without intermediate `submitted` — see §2.3.3 below)
- `drafted → withdrawn`
- `submitted → accepted | modified | rejected | withdrawn`
- `accepted → modified` (post-filing refinement; Delta updated)
- `modified → modified` (further refinement)
- `accepted → landed`
- `modified → landed`
- `* → revised | corrected` (any non-terminal state may emit an audit event without changing lifecycle state)

Illegal transitions:
- `rejected → *` (terminal)
- `landed → *` (terminal — for 1:1; master/sub permits multiple `landed` rows per §2.3.4)
- `withdrawn → *` (terminal)
- `submitted → landed` (must pass through `accepted` or `modified`)
- `modified → accepted` (cannot un-modify once delta is recorded)

#### 2.3.3. `drafted → accepted` (direct acceptance from drafted state)

The cn-sigma/agent-activate-skill anchor shows `drafted sigma → accepted gamma@cnos` without an intermediate `submitted` event. `cnos.cds/skills/cds/CDS.md` §"Coordination surfaces" → §"Cross-repo proposals" says `submitted` is "the only event required for target intake" — implying intake without `submitted` is unusual but not impossible. The protocol permits `drafted → accepted` (or `modified` / `rejected`) when:

1. The source explicitly delegates filing-authority to the target (e.g. an agent hub authoring a proposal whose first reader is expected to act, not approve-then-act).
2. The target γ records the delegation in the cnos-side LINEAGE's `## Source` or `## Target` section.

The mirror STATUS may carry the original `drafted` event; the cnos-side mirror does not retroactively insert a `submitted` event the source did not emit.

#### 2.3.4. Emitters per event

| Event | Emitter | Write target |
|---|---|---|
| `drafted` | source role (γ or other) authoring the bundle | source-side STATUS at authoring time |
| `submitted` | source γ | source-side STATUS when source γ requests target intake |
| `accepted` / `modified` / `rejected` | target γ | source-side STATUS via direct push if cross-repo write is authorized; otherwise via `FEEDBACK.patch` for source γ to apply |
| `landed` (1:1) | target γ | source-side STATUS (and cnos-side mirror STATUS) when the cycle implementing the target issue merges to target main |
| `landed` (master/sub) | target γ per §2.3.5 | source-side STATUS (one row per sub merge + one terminal master-close row) |
| `withdrawn` | source γ | source-side STATUS when source retracts the request |
| `revised` / `corrected` | event-originator (source or target as relevant) | source-side STATUS; optional audit-only — does not change lifecycle state |

Constraint: **after the target γ's filing decision, source STATUS must not remain at `submitted`** (or at `drafted` once target γ has taken action on a delegated `drafted → accepted` path). The decision must be recorded — directly or via feedback patch — within the same target session that made the decision.

#### 2.3.5. Master/sub rule for `landed`

For master/sub-shaped proposals (case a.2):

1. **Per-sub `landed` event** fires each time a sub merges to target main, with the form:
   ```
   <ISO-date> landed gamma@<target> <target-sub-issue#> <merge-sha> <release-version>
   ```
2. **Terminal master-close `landed` event** fires when the master issue is closed (all subs landed, or deferred subs explicitly named as tracked debt in the master's closure condition):
   ```
   <ISO-date> landed gamma@<target> <target-master-issue#> master-close
   ```

A wave with 4 subs landing across 3 releases produces 5 `landed` rows: 4 per-sub + 1 master-close.

For 1:1 proposals (case a.1): one `landed` event total. No separate master-close because the issue is the only target unit.

- ❌ Leave a master/sub proposal at `accepted` indefinitely after subs land, because "the master hasn't closed yet."
- ✅ Emit per-sub `landed` rows as each sub merges; emit the master-close row when the master closes.

### 2.4. Bundle-state phases

`.cdd/iterations/cross-repo/README.md` defines three bundle-state phases: `open | converging | closed`. Mapping to STATUS:

| Bundle phase | STATUS event(s) | Meaning |
|---|---|---|
| `open` | `drafted` or `submitted` (or initial state) | Filed; no target disposition yet. Both pre-intake (`drafted`) and intake-requested (`submitted`) sit here. |
| `converging` | `accepted` or `modified` (no terminal `landed` yet; may carry partial per-sub `landed` rows) | Target accepted; work in flight. |
| `closed` | terminal `landed` (1:1) **or** master-close `landed` (master/sub) **or** `rejected` **or** `withdrawn` | No further work expected. |

Audit events (`revised`, `corrected`) do not change the bundle's phase — they record post-hoc corrections to earlier events without re-opening or re-closing the bundle.

For cases (b), (c), (d) — which carry no STATUS — the bundle-state phase is derived from LINEAGE state:

| Case | Phase derivation |
|---|---|
| (b) outbound | `open` while patches are in flight; `closed` when all patches confirmed landed in the target's per-patch table. |
| (c) bilateral | `open` while ε deliverables (cnos-side and counterpart-side) are pending; `closed` when both sides land their respective patches and the iteration record is committed in each repo's `cdd-iteration.md`. |
| (d) operator-pending | `drafted (operator-pending)` (LINEAGE.md `Disposition` field) until the operator effects the application; then `closed`. This is a phase-name synonym for `open` for inventory purposes — case (d) bundles report as `open` in bundle-state inventories until their application gate fires. |

### 2.5. Bundle file set per case

#### Case (a) — Inbound proposal (1:1 or master/sub)

Required files on the cnos side:

| File | Required? | Purpose |
|---|---|---|
| `LINEAGE.md` | yes | Bilateral trace; per §2.6 schema |
| `STATUS` | yes | Event ledger (mirror of source STATUS + cnos-side authoritative `landed` rows) |
| `FEEDBACK.patch` | conditional | Required when the cnos session cannot write to source-side STATUS directly |
| `ISSUE.md` | optional | Verbatim copy of source `ISSUE.md` for audit (useful when source bundle is dormant) |
| `README.md` | **discouraged** | Do not place protocol doctrine here; bundle-specific narrative belongs in `LINEAGE.md` |

#### Case (b) — Outbound iteration trace

Required files on the cnos side (mirror, after source-side bundle is archived):

| File | Required? | Purpose |
|---|---|---|
| `LINEAGE.md` | yes | Bilateral trace; per-patch confirmation; mirror status |

The source-side bundle (under counterpart-repo's `.cdd/iterations/cross-repo/cnos/{slug}/`) carries the working files (patches, lineage, STATUS if a proposal lifecycle is attached). After target patches merge and target-side mirror confirms receipt, the source-side may be archived (see §2.8).

#### Case (c) — Bilateral iteration

Required files on the cnos side:

| File | Required? | Purpose |
|---|---|---|
| `LINEAGE.md` | yes | Bilateral trace + hat-collapse acknowledgment + per-patch trace per destination repo |
| `cdd-iteration.md` | yes | The ε work product when the bilateral iteration carries ε output (per `post-release/SKILL.md §5.6b` and `epsilon/SKILL.md §1`) |
| `<patch>.patch` | one per outbound patch | Feedback patches for the counterpart repo to apply, named descriptively |

#### Case (d) — Operator-pending

Required files on the cnos side:

| File | Required? | Purpose |
|---|---|---|
| `LINEAGE.md` | yes | Carries `Status: drafted (operator-pending)` and the application gate |
| Substantive content | varies | Drafted artifact(s): file content for (d.1), `<file>.patch` for (d.2), `issue-comment.md` for (d.3) |

No STATUS ledger — there is no proposal lifecycle to track.

#### Rule 2.5.a — Use the bundle file set the case requires, no more, no less

- ❌ Author a case (a) bundle without a STATUS file because "the precedent didn't have one" (the precedent was a different case).
- ❌ Add a README.md carrying protocol doctrine because "future readers will want context" (doctrine archives with the bundle if the bundle is archived).
- ✅ Match the file set to the case in §2.5; place narrative inside `LINEAGE.md`; place protocol doctrine here in this skill.

### 2.6. LINEAGE.md schema per case

#### Case (a) inbound proposal — required sections

- `## Source` (counterpart repo, branch, commit, path)
- `## Target` (cnos repo, issue #, mode, disposition, first filing commit/branch, canonical-path on cnos main)
- `## Delta` — required only when disposition = `modified`; names what changed and why
- `## Bilateral trace` — which side carries what; archival predicate
- `## Per-sub confirmation` — master/sub only; table of subs (proposed shape, cnos issue #, cnos commit, status). **Omit for 1:1.**
- `## Per-patch confirmation` — 1:1 with multi-patch only; table of patches and target commits
- `## Related cross-repo wave context` — optional; links to wave artifacts, repository rename events, sibling bundles

#### Case (b) outbound — required sections

- `## Source` (cnos cycle, branch, commit, source-side path)
- `## Target` (counterpart repo, target issue, target PR/commit, patches landed)
- `## Bilateral trace` (mirror status)
- `## Per-patch confirmation` — table of patch numbers, target commits, status

#### Case (c) bilateral — required sections

- `## Source` (counterpart repo, branch, wave, master issue, sub issues, wave verdict)
- `## ε actor` (or relevant role name) — hat-collapse acknowledgment: played-by, authority basis, output date (per §2.9)
- `## Target patches` — per-patch table, one row per destination repo
- `## Bundle contents` — file-by-file purpose
- `## Bilateral trace direction` — why this is on this side; what the counterpart side carries
- `## Protocol observations forwarded` — optional; when the iteration surfaces protocol gaps for codification in this skill

#### Case (d) operator-pending — required sections

- `## Purpose` — one-paragraph what-this-stages
- `## Source` — trigger, reference doctrine, sibling exemplar, empirical anchor
- `## Bundle contents` — file-by-file purpose
- `## Operator-side scaffolding steps` (case d.1) OR `## Application` (case d.2/d.3)
- `## Verification after application` — case d.2 specifically; how the operator confirms the patch landed
- `## What is NOT in scope here` / `## What this bundle does NOT do` — optional but recommended; prevents scope creep
- `## Disposition` — always `drafted (operator-pending)` for case d

#### Rule 2.6.a — Select the LINEAGE schema by case, not by precedent

- ❌ Copy the LINEAGE shape from the most recent bundle, even if the case differs.
- ✅ Identify the case (§2.1), then use the schema in §2.6 for that case.

### 2.7. Feedback-patch format

A feedback patch is a unified diff against a counterpart-side file, wrapped in a header.

**Header form:**

```
From: <actor identity> (<originating branch>)
Date: <ISO date>
Subject: <one-line summary>

<free-prose context: apply command, rationale, links to related artifacts>

---

<unified diff>
```

**Apply command convention:** `git apply <patch-filename>` from the counterpart repo root.

**Patch filename convention:**

- `FEEDBACK.patch` — when the patch's purpose is STATUS-event emission only. Canonical name.
- `<descriptor>.patch` — when the patch carries substantive content. Lowercase-with-hyphens descriptor (e.g. `cdr-changelog-rule.patch`, `PERSONA-discipline.patch`).

**Content rules:**

- Diff is against the counterpart-side file path (e.g. `a/.cdd/iterations/cross-repo/cnos/{slug}/STATUS`, `a/CDR.md`, `a/spec/PERSONA.md`).
- Header includes the cnos-side cross-reference to where the bilateral trace lives (`cnos:.cdd/iterations/cross-repo/{counterpart}/{slug}/LINEAGE.md`).
- When the patch documents a protocol observation worth forwarding (e.g. cn-sigma's `drafted`-event observation forwarded to cnos#377), include the observation in the header context.

#### Rule 2.7.a — Use the canonical feedback-patch format

- ❌ Emit a feedback "patch" as a YAML script, a replacement-file blob, or a free-prose change request.
- ✅ Emit a unified diff with the header form above; the counterpart γ applies with `git apply`.

### 2.8. Bundle archival rule

A bundle's archival rule depends on direction and side.

#### 2.8.1. Source-side may be deleted when

| Case | Source side may archive when |
|---|---|
| (a) inbound | The cnos-side mirror exists on cnos main **and** the cycle implementing the target work has closed (terminal `landed` event recorded). |
| (b) outbound | The cnos-side (source) bundle may be archived once the counterpart-side mirror confirms receipt of all patches landed. |
| (c) bilateral | The counterpart-side bundle (the side that cannot write directly) may be archived once the counterpart applies the patches and emits its own iteration trace. |
| (d) operator-pending | The bundle may be archived once the operator effects the application (target repo created with files / patch applied / comment posted). |

#### 2.8.2. Target-side (mirror) is preserved indefinitely

Mirror bundles on `cnos:.cdd/iterations/cross-repo/{counterpart}/{slug}/` when cnos is the **target** (case a) are preserved indefinitely as audit artifacts and retroactive validators. They are not subject to archival.

Source-side bundles where cnos is the **source** (case b / c) may be deleted post-archival per §2.8.1.

#### Rule 2.8.a — Asymmetric archival

- ❌ Delete the cnos-side mirror of an inbound proposal once the work lands ("cycle is done; mirror not needed").
- ✅ Preserve target-side mirrors as audit artifacts; archive source-side bundles once the counterpart records receipt.

### 2.9. Hat-collapse attribution

When a single actor plays γ in repo A and ε (or another role) in repo B during the same session, attribution is recorded in **two places**:

1. **In the bundle `LINEAGE.md`**, a `## ε actor` (or equivalent role-name) section names:
   - `Played by:` — actor identity + originating session branch
   - `Authority basis:` — citing `epsilon/SKILL.md §2` ("ε and δ are often the same actor") + `ROLES.md §4` role-collapse rules + any operator-specific delegation
   - `Output date:` — ISO date

2. **In the `cdd-iteration.md`** (when case (c) produces ε output), a top-level `## Hat-collapse acknowledgment` section names:
   - which two roles are collapsed in this session
   - the authority basis citation
   - an explicit statement that the substantive role-A and role-B work are independent
   - the reason hat-collapse does not collapse with α/β/γ responsibilities for any in-flight cycle

**Boundary preservation:** the attribution form names role-collapse as factual record, not as authority override. Per `ROLES.md §4a` persona/protocol/project boundary, the **personas** stay distinct (e.g. cnos γ identity ≠ cph ε identity); the actor wearing both hats acknowledges the collapse and carries the substantive work for each role independently. This skill's hat-collapse section documents an actor-level fact (one human/agent ran two sessions of work) without merging the role-level identities or authorities — the persona/protocol/project boundary is preserved.

#### Rule 2.9.a — Record hat-collapse where it happens

- ❌ Hide hat-collapse by attributing both pieces of work to "the session" without naming the actor or the two roles.
- ✅ Name the collapse in `LINEAGE.md §"<role> actor"` and (for case c) in `cdd-iteration.md §"Hat-collapse acknowledgment"`.

### 2.10. Known protocol edge cases

These are observed shapes that the protocol covers explicitly. New edge cases discovered in future cycles are added here.

- **`drafted → accepted` direct acceptance** — agent-hub source repos may author a bundle in `drafted` state and delegate filing-authority to the target, who acts on the bundle without an intermediate `submitted` event. The transition is legal (per §2.3.3) when the cnos-side LINEAGE records the delegation in `## Source` or `## Target`. Empirical anchor: `cn-sigma/agent-activate-skill` (source `drafted sigma cn-tau@8f153c15e` → target `accepted gamma@cnos cnos#379`).

- **Repository rename mid-coordination** — when a source repo is renamed (e.g. `gait-support-paths` → `cph` on 2026-05-18), the cnos-side mirror's path is updated to the new name; the pre-rename mirror is preserved as an archived audit artifact. The LINEAGE.md `## Repository rename event` section names the rename and its consequences. Empirical anchor: `cph/bootstrap-cdr` (with archived pre-rename mirror at `gait-support-paths/bootstrap-cdr`).

- **Post-filing refinement** — a target γ may transition `accepted → modified` when post-filing upstream changes warrant folding into the target issue body. The `Delta` field of the target issue's `## Source Proposal` block is updated; the source STATUS receives a `modified` event. Empirical anchor: `cn-sigma/agent-activate-skill` (AC4 capability-matrix refinement folded in 7 minutes after filing).

- **Asymmetric source posture** — when the source repo is an agent hub (not a CDD-activated tenant repo), the source-side bundle may live on source main rather than on a dormant branch. The cnos-side mirror remains canonical. Empirical anchor: `cn-sigma` (full bundle on cn-sigma main, asymmetric to `cph/bootstrap-cdr` whose bundle is on a dormant branch).

- **Dual-purpose case-c bundle** — case (c) bundles carry both iteration content (for the counterpart's γ to apply) and the cnos-side patch trace. The dual purpose is named explicitly in `LINEAGE.md §"Bilateral trace direction"`. Empirical anchor: `cph/coherence-drift-sweep-followup-2026-05-18`.

- **Proposal-as-issue-comment (case d.3)** — a proposal that retargets or comments on an existing target issue, rather than emitting a new ISSUE.md or patch. The bundle carries `issue-comment.md` as the deliverable. The operator pastes the comment on the target issue. Empirical anchor: `cph/issue-32-tightening-2026-05-19`.

- **Target-repo-doesn't-exist (case d.1)** — drafts staged for an operator-scaffolded repo. The bundle carries the prospective files (e.g. `PERSONA.md`, `OPERATOR.md`). Empirical anchor: `cn-rho/bootstrap-2026-05-19`.

---

## 3. Rules

### 3.1. Identify the directional case first

- ❌ Author a bundle "following the precedent" without naming which case the cycle is in.
- ✅ Name the case (a.1 / a.2 / b / c / d.1 / d.2 / d.3) in the first paragraph of `LINEAGE.md`; the schema selection becomes auditable.

### 3.2. Use the canonical path

- ❌ `.cdd/proposals/{target}/{slug}/` or `.cdd/iterations/proposals/{slug}/`.
- ✅ `.cdd/iterations/cross-repo/{counterpart-repo}/{slug}/`.

### 3.3. After filing decision, source STATUS must not remain at `submitted` or `drafted`

- ❌ File the target issue, then forget to emit the `accepted` event to source STATUS (directly or via FEEDBACK.patch).
- ✅ Emit the disposition event in the same session that made the decision. (For `drafted → accepted` direct acceptance per §2.3.3, the same constraint applies — source STATUS must not remain at `drafted` after target γ acts.)

### 3.4. Emit `landed` per master/sub rule

- ❌ Treat `landed` as a single event for a master/sub proposal ("the work is done").
- ✅ Per-sub `landed` rows as each sub merges; terminal master-close `landed` row when the master closes (§2.3.5).

### 3.5. Keep protocol doctrine in this skill, not in bundle READMEs

- ❌ Place archival rules, schema definitions, or STATUS doctrine inside a bundle's `README.md`.
- ✅ Doctrine lives here; bundle `LINEAGE.md` describes only the bundle's specific context.

### 3.6. Mirror target-side bundles are preserved indefinitely

- ❌ Delete a cnos-side mirror of an inbound proposal once the work lands.
- ✅ Preserve target-side mirrors as audit artifacts and retroactive validators.

### 3.7. Use the canonical feedback-patch format

- ❌ Emit a feedback "patch" as YAML, a replacement file, or prose.
- ✅ Unified diff with the header form in §2.7; `git apply` from counterpart repo root.

### 3.8. Record hat-collapse where it happens

- ❌ Attribute the work to "the session" without naming the actor or the two roles.
- ✅ `LINEAGE.md §"<role> actor"` + (for case c) `cdd-iteration.md §"Hat-collapse acknowledgment"`.

### 3.9. Identify case (d) explicitly when no STATUS lifecycle applies

- ❌ Force a case-d bundle into the proposal-lifecycle frame with a hollow STATUS.
- ✅ Name `(d.1)` / `(d.2)` / `(d.3)` in LINEAGE; carry `Disposition: drafted (operator-pending)` and the application gate.

### 3.10. Cross-reference, do not duplicate

Other CDD skills cite this skill for cross-repo doctrine. If you find cross-repo doctrine in `gamma/`, `post-release/`, or elsewhere that contradicts this skill, this skill governs and the other surface is patched to reference this one.

---

## 4. Verify

### 4.1. Case-identification check

- Did `LINEAGE.md` name the case (a.1 / a.2 / b / c / d.1 / d.2 / d.3) in its first paragraph?
- Does the bundle file set match §2.5 for that case?
- Does the LINEAGE schema match §2.6 for that case?

### 4.2. STATUS check (cases a, c when proposal-shaped)

- Is the source STATUS at a non-`submitted` event after target γ's filing decision?
- For master/sub: are per-sub `landed` rows present for each merged sub?
- Are illegal transitions absent (e.g. no `submitted → landed` skip)?

### 4.3. Archival check

- For source-side bundles: is the archival predicate (per §2.8.1) named in `LINEAGE.md §"Bilateral trace"`?
- For target-side mirrors: are they preserved (not subject to archival)?

### 4.4. Doctrine-locality check

- Is the bundle's `README.md` (if present) free of protocol doctrine?
- Does protocol doctrine live only in this skill?

---

## 5. Cross-references

- `cdd/gamma/SKILL.md §2.1` — cross-repo proposal intake: references this skill for case (a) directional rules.
- `cdd/gamma/SKILL.md §2.7` — cross-repo proposal close-out: references this skill for `landed` emitter and master/sub rule.
- `cdd/post-release/SKILL.md §5.6b` — cross-repo trace: references this skill for case (b) outbound bundle file set and §2.8 archival.
- `cdd/issue/SKILL.md` minimal output pattern `## Source Proposal` block — target-side issue-body integration; stays in `issue/`.
- `.cdd/iterations/cross-repo/README.md` — bundle convention; references this skill for the canonical state machine.
- `cdd/epsilon/SKILL.md §1` — ε's `cdd-iteration.md` work product; relates to case (c) bilateral bundle file set §2.5.
- `ROLES.md §4` — role-collapse rules; cited from §2.9 hat-collapse attribution.
- `ROLES.md §4a` — persona / operator / protocol / project / receipt boundary; cited from §2.9 to confirm the attribution form preserves the boundary.

---

## 6. Kata

### Scenario

cnos γ is dispatched on a new cycle. While running the intake scan in `gamma/SKILL.md §2.1`, γ finds a `STATUS` ledger in a source repo `usurobor/cn-tau` at `.cdd/iterations/cross-repo/cnos/agent-trace-hook/STATUS` whose last non-comment event reads `drafted sigma cn-tau@abc1234`.

### Task

1. Identify the directional case.
2. Decide the disposition (accepted / modified / rejected).
3. Stage the cnos-side bundle correctly.
4. Emit the STATUS update.

### Expected reasoning

1. **Case identification.** Source is cn-tau; substantive target work would happen in cnos. This is case (a) inbound. Single proposal; sub-shape is (a.1) unless the proposal body declares master/sub.

2. **`drafted` reading.** Per §2.3.1, `drafted` is a distinct first-class event meaning "source has written the proposal but has not requested target action". The cn-tau STATUS at `drafted` means filing-authority has not yet been requested by source γ — γ checks whether the source posture (agent hub) is delegating filing-authority. If yes, γ proceeds via §2.3.3 `drafted → accepted` direct-acceptance and records the delegation in LINEAGE. If no, γ asks the source γ to emit `submitted` first.

3. **Disposition.** γ reads `ISSUE.md` and decides; let's say `accepted`. The cnos-side target issue is filed with a `## Source Proposal` block citing cn-tau path + commit + disposition.

4. **Bundle staging.** Cnos-side bundle at `cnos:.cdd/iterations/cross-repo/cn-tau/agent-trace-hook/`:
   - `LINEAGE.md` — case (a.1) schema: `## Source`, `## Target`, `## Bilateral trace`. Names the case in the first paragraph.
   - `STATUS` — mirror of source STATUS, preserving the `drafted` event (do not retroactively insert `submitted`) and appending the new `accepted gamma@cnos cnos#<new-issue#>` row per §2.3.3 direct-acceptance.
   - `FEEDBACK.patch` — required if cnos session cannot write to cn-tau. Header + apply command + unified diff against source STATUS appending the `accepted` event.

5. **STATUS update.** Source STATUS must not remain at `submitted`/`drafted` post-decision. Either direct push to cn-tau (if authorized) or `FEEDBACK.patch` for cn-tau γ to apply.

### Common failures

- Filing the cnos issue but forgetting to emit the `accepted` event to source STATUS.
- Authoring LINEAGE.md without naming the case in the first paragraph.
- Placing the FEEDBACK.patch under `.cdd/proposals/...` instead of the canonical path.
- Treating `drafted` as unknown and refusing intake.

### Verification

- `LINEAGE.md` first paragraph names case (a.1).
- Source STATUS last event = `accepted gamma@cnos cnos#<N>`.
- Bundle file set matches §2.5 case (a).
- No bundle `README.md` carrying protocol doctrine.

### Reflection

A second γ reading this kata's output should be able to reconstruct the disposition, locate every artifact, and verify all events fired in the right order — without consulting the original cycle session.
