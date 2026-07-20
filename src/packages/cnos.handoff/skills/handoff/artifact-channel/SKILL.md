---
name: artifact-channel
description: Intra-cycle artifact-channel rules. Use when roles exchange per-cycle artifacts under .cdd/unreleased/{N}, when γ marks post-merge closeout release-ready, or when γ archives the directory after δ disconnects.
artifact_class: skill
kata_surface: embedded
governing_question: How do α, β, and γ hand work to each other inside a single cycle — using a sequential per-role-owned write channel over .cdd/unreleased/{N}/ — without inventing the channel rules per cycle and without leaking matter back into the channel after merge?
visibility: public
parent: handoff
triggers:
  - artifact-channel
  - unreleased
  - per-role artifact
  - sequential handoff
  - cycle directory
  - frozen snapshot
  - post-disconnect archive
scope: task-local
inputs:
  - active CDD-class cycle on origin/cycle/{N}
  - active role (α / β / γ; ε at close-out)
  - predecessor-role artifact(s) when applicable
outputs:
  - per-role artifact written at .cdd/unreleased/{N}/{filename}.md
  - SHA advance on origin/cycle/{N} (the handoff signal to the successor role)
  - "after disconnect: cycle directory moved from .cdd/unreleased/{N}/ to .cdd/releases/{X.Y.Z}/{N}/, then terminally sealed"
requires:
  - cycle/{N} exists on origin
  - the role's predecessor (if any) has signaled completion via the channel (see §2.3 Sequential rule)
calls:
  - HANDOFF.md
  - mid-flight/SKILL.md
  - dispatch/SKILL.md
  - cross-repo/SKILL.md
---

# Artifact Channel — intra-cycle α → β → γ sequential handoff

## Core principle

**`.cdd/unreleased/{N}/` is the sequential handoff substrate. Predecessor artifacts freeze as their roles exit; post-merge role closeouts and the nonterminal γ marker accumulate there; δ validates and disconnects while it remains there; γ moves it only afterward and then seals terminal closure.**

Before merge, SHA advance on `origin/cycle/{N}` is the handoff signal. After merge, the same directory lives on main while role closeouts and γ's release-readiness marker land. δ's release gate reads that exact unreleased path. After tag/disconnect and green CI, γ moves it to the versioned path, then appends terminal evidence in a separate commit.

The failure mode is **channel-substrate drift** — a role writes a coordination signal to chat instead of to the channel; a role writes to a sibling role's artifact and obscures attribution; the cycle directory is mutated after merge (e.g. an α-closeout edit dated after release); the directory is left at `.cdd/unreleased/{N}/` past release and loses its version association. The result is a cycle whose record is incomplete, whose role attribution is ambiguous, and whose post-release assessment cannot reconstruct what happened from the cycle directory alone.

This skill is the canonical surface for channel shape, role ownership, sequential flow, predecessor-artifact freezing, and the post-disconnect archive. CDS governs the per-artifact contents and lifecycle steps.

## Authority

- This skill is the only canonical home for the channel-shape invariants: `.cdd/unreleased/{N}/`, per-role ownership, sequential α→β→γ flow, predecessor-artifact freezing, retention under `unreleased/` through δ disconnect, and γ's post-disconnect move followed by terminal sealing.
- The consumer protocol packages govern the per-artifact instantiation. For CDS-class cycles, the canonical home is `cnos.cds/skills/cds/CDS.md §"Artifact contract"` — specifically the Location matrix (canonical filename per artifact), Ownership matrix (per-role artifact assignment with verification gates), and Ordered flow (the 13-stage flow: design → contract → plan → tests → code → docs → self-coherence → review → gate → release → observe → assess → close). CDR-class cycles bind differently; future c-d-X may bind differently still. The wire-format pattern is what this skill owns; the per-consumer instantiation is what cds/cdr/c-d-X own.
- The artifacts on this channel are also the substrate for sibling sub-skills' wire formats. `cnos.handoff/skills/handoff/mid-flight/SKILL.md` (gamma-clarification.md / gamma-coordination.md) and `cnos.handoff/skills/handoff/receipt-stream/SKILL.md` (cdd-iteration.md aggregator) live on this channel. Their authoring/reader/cadence rules are theirs; the channel-shape they live on is this skill's.

## Scope

In scope:

- The channel directory path pattern (`.cdd/unreleased/{N}/`); the cycle-branch substrate (`origin/cycle/{N}`); the SHA-advance handoff signal.
- The per-role write ownership pattern (which role writes which class of artifact; what successor roles may and may not do to predecessor artifacts).
- The sequential rule (α before β; β before γ; γ-scaffold exception at cycle-start).
- The freeze/archive rule (predecessor artifacts freeze at role exit; the directory remains under `unreleased/` through disconnect; γ moves it afterward and terminally seals it).
- The cross-cycle aggregation pointer (per-cycle `cdd-iteration.md` → `.cdd/iterations/INDEX.md`; aggregator rules canonicalize in Sub 5 / receipt-stream/SKILL.md).
- Empirical anchors (every cycle since the channel was established in pre-#364 form).

Out of scope:

- The per-artifact filenames and content sections. Those live in `cnos.cds/skills/cds/CDS.md §"Artifact contract" → "Location matrix"` (CDS-specific filename set: `self-coherence.md`, `alpha-closeout.md`, `beta-review.md`, `beta-closeout.md`, `gamma-closeout.md`, `gamma-scaffold.md`, `cdd-iteration.md`, `RELEASE.md`, version-snapshot directory, PRA path, cross-repo trace dir). CDR and future c-d-X bind differently.
- The per-artifact role ownership table (who writes which file, when, what the verification gate is, what missing means). Lives in `cnos.cds/skills/cds/CDS.md §"Artifact contract" → "Ownership matrix"`.
- The CDS-specific 13-stage Ordered flow. Lives in `cnos.cds/skills/cds/CDS.md §"Artifact contract" → "Ordered flow"`. CDR has a different flow; future c-d-X may have other flows.
- The mid-flight rescue mechanism (`gamma-clarification.md` write/read protocol; cache-bust semantics). Lives in `cnos.handoff/skills/handoff/mid-flight/SKILL.md`. This skill owns the channel; mid-flight owns the rescue protocol that lives on the channel.
- The cross-cycle aggregator content (per-finding shape; INDEX.md row format; cross-repo trace bundle). Lives in `cnos.handoff/skills/handoff/receipt-stream/SKILL.md` (Sub 5 of cnos#404; landed via cnos#419). This skill owns the per-cycle channel; receipt-stream owns the cross-cycle aggregation.
- The harness substrate (polling primitives, Monitor wrapping, reachability re-probe). Lives in `cnos.cdd/skills/cdd/harness/SKILL.md`. This skill assumes a polling primitive emits on cycle-branch SHA transitions.
- Tooling that auto-verifies channel discipline (`cn cdd verify` channel-integrity checks). Deferred; doctrine first, tooling later.

---

## 1. Define

### 1.1. What the channel is

The channel is **one directory per cycle**: `.cdd/unreleased/{N}/`, first on `origin/cycle/{N}` and then on main after merge. It carries the record from γ scaffold through the marked nonterminal post-merge closeout. δ validates and disconnects without moving it. γ moves it to `.cdd/releases/{X.Y.Z}/{N}/` only after green release CI, then appends terminal evidence.

The substrate is filesystem + git. The handoff signal between roles is the SHA advance on `origin/cycle/{N}` — the polling successor role's polling primitive (per `cnos.cds/skills/cds/CDS.md §"Coordination surfaces" → "Polling primitives"`) emits on the SHA transition and the successor's next intake iteration reads the predecessor's artifact live. No chat, no out-of-band coordination database, no separate inbox.

### 1.2. What the channel is not

The channel is **not**:

- a chat-coordination substitute. Chat does not transit the channel; if a load-bearing fact is communicated only in chat, it is invisible to polling roles and to the post-release record.
- a shared-write directory. Each artifact has one owning role at write time. Sibling roles read but do not modify (except by their own append, e.g. α's fix-round sections after β's RC verdict — α appends to α's own `self-coherence.md`, never to β's `beta-review.md`).
- a multi-cycle channel. One directory = one cycle. Cross-cycle aggregation runs through `.cdd/iterations/INDEX.md` per the receipt-stream sub-skill (Sub 5 of cnos#404; landed via cnos#419).
- freely mutable post-merge. Predecessor artifacts are frozen, but their owning roles add the prescribed closeouts and γ later adds post-release/terminal evidence; no role rewrites a predecessor's matter.
- the canonical post-release snapshot. The frozen post-release snapshot lives under `docs/{tier}/{bundle}/{X.Y.Z}/` per CDS Location matrix; the channel directory at `.cdd/releases/{X.Y.Z}/{N}/` is the cycle's role-local close-out evidence storage.

### 1.3. Failure mode

The channel fails through **channel-substrate drift**:

- A role writes a coordination signal to chat instead of to the channel — the signal is invisible to polling sessions and to the post-release assessment.
- A role writes to a sibling role's artifact — attribution becomes ambiguous; β's review verdict and α's self-coherence are no longer distinguishable.
- A predecessor artifact is rewritten after its role exits — attribution and evidence custody are lost.
- The directory is moved before δ validates/tags, or left under `unreleased/` after green release CI and γ's archive opportunity — either ordering breaks the boundary record.
- A coordination note lives only in the GitHub issue comments or in a chat backlog — the artifact channel is incomplete and the post-release reconstruction is impossible.

The protocol closes each failure mode by naming the channel-shape invariants explicitly and binding role-write ownership, sequential flow, and frozen-snapshot rules to the channel directory.

---

## 2. Unfold

### 2.1. Channel directory

The canonical channel path is:

```text
.cdd/unreleased/{N}/
```

where `{N}` is the cycle's GitHub issue number. The directory exists from γ scaffold through δ disconnect and moves only during γ's post-disconnect archive. The path is consumer-agnostic; only filenames vary.

Before merge the substrate is `origin/cycle/{N}`. After merge the directory is on main under `unreleased/`; post-merge roles and δ read it there through disconnect. After γ moves it, readers use the versioned release path.

The handoff signal between roles is the **SHA advance on `origin/cycle/{N}`**. Each commit to the cycle branch advances the SHA; polling primitives detect the transition and wake successor roles. The signal is git-native — no separate notification channel, no webhook, no inbox. The SHA advance is the only event the successor role needs.

### 2.2. Per-role write ownership

Every per-cycle artifact has exactly one owning role at write time. The pattern (consumer-protocol-agnostic):

| Owner role | Artifact class | Pattern |
|---|---|---|
| γ | scaffold | Cycle-start artifact; carries γ's selection decision and the dispatch surface set. One per cycle. |
| α | role-local primary branch artifact | The role's review-readiness signal; gap / mode / ACs / trace / fix-rounds. One per cycle. |
| α | role close-out | The role's post-merge findings record. One per cycle. |
| β | role-local review record | The review CLP rounds + verdicts. One per cycle; one section per round. |
| β | role close-out | The role's post-review/merge close-out evidence. One per cycle; no δ release evidence. |
| γ | role close-out | Preliminary triage + nonterminal release-readiness marker before disconnect; terminal evidence appended after archive. One evolving γ-owned file per cycle. |
| γ | cross-cycle receipt | The per-cycle iteration artifact (when `protocol_gap_count > 0`; courtesy stub permitted when count is 0). Feeds the cross-cycle aggregator. |
| γ | mid-flight rescue | The clarification artifact (per `cnos.handoff/skills/handoff/mid-flight/SKILL.md`). One file per cycle (appended to per event). |

**Read access is universal; write access is role-local.** α reads γ's scaffold (it is α's intake brief); β reads α's primary branch artifact (it is β's review subject); γ reads both close-outs (they are γ's triage inputs). The reader does not modify the predecessor's artifact — the reader makes its observations and writes its own role-local artifact.

The one structured exception is **α's fix-round append after β's RC verdict.** When β returns RC, α appends a fix-round section to its *own* primary branch artifact (`self-coherence.md` for CDS-class) naming each finding addressed, the commit SHA that addressed it, and any reasoning β needs to re-verify. α does not edit β's review record; β iterates a new round in `beta-review.md`. The fix-round pattern keeps role attribution intact.

The consumer-specific *filename* set lives in `cnos.cds/skills/cds/CDS.md §"Artifact contract" → "Location matrix"` (CDS-class: `gamma-scaffold.md`, `self-coherence.md`, `alpha-closeout.md`, `beta-review.md`, `beta-closeout.md`, `gamma-closeout.md`, `cdd-iteration.md`, plus the cycle-orthogonal `RELEASE.md` and PRA at canonical version-directory paths). CDR has its own filename set; future c-d-X protocols bind their own. The wire-format pattern (per-role ownership, sequential flow) is what this skill owns; the per-protocol filename set is consumer-side.

### 2.3. Sequential rule

The intra-cycle flow is **sequential α → β → γ** with one cycle-start exception:

1. **γ-scaffold (cycle-start exception)** — γ writes the scaffold artifact *before* α is dispatched. This is the only cycle-start γ write; it is the binding gate per `cnos.cdd/skills/cdd/gamma/SKILL.md §2.5 → "Pre-dispatch γ scaffold check"`. The scaffold is α's input.
2. **α phase** — α writes its primary branch artifact incrementally on the cycle branch; the review-readiness section is the last section appended. The α phase ends when α signals review-readiness via the review-readiness section's append commit. α's primary artifact is β's input.
3. **β phase** — β reads α's primary artifact + the diff + the issue body; runs the review CLP; writes the review record incrementally (one section per round). On A verdict, β merges. On RC, α appends a fix-round; β re-reviews. The β phase ends at merge (for non-RC paths) or at A after RC rounds.
4. **Re-dispatched α (close-out)** — after β merge, α is re-dispatched (per `cnos.cdd/skills/cdd/operator/SKILL.md` re-dispatch protocol) and writes its close-out artifact. α exits.
5. **β close-out** — β writes its review/merge close-out findings in the same β session as the merge (per CDS Field 6 collapse). β exits before release; δ's later release evidence is not β-owned.
6. **γ post-merge phase** — γ verifies both closeouts, writes preliminary triage and the exact nonterminal marker under `unreleased/`; γ does not move the directory or declare terminal closure.
7. **δ disconnect** — δ validates the marked unreleased receipt, tags/releases, and reports green CI (or explicit override).
8. **γ terminal phase** — γ observes/assesses, commits the post-disconnect archive move, then appends and commits terminal evidence separately.

The successor role's intake **requires** the predecessor role's signaled completion. β does not begin review until α has signaled review-readiness (in α's primary artifact). γ does not begin triage until both close-outs are on main. Each role's polling primitive emits when the predecessor's relevant artifact reaches the expected state; the role's intake iteration reads the artifact live.

The two exceptions to the strictly-sequential pattern:

- **The γ-scaffold cycle-start write** (§2.3 step 1) — γ writes before α exists in the cycle. Not strictly sequential because α has not yet written anything; γ-scaffold is the cycle's initial state.
- **Mid-flight γ writes** (`gamma-clarification.md` / `gamma-coordination.md`) — γ writes mid-cycle while α and/or β are in-flight. This is the rescue channel (per `cnos.handoff/skills/handoff/mid-flight/SKILL.md`); the channel is the substrate but the sequence is not strict α→β→γ. The clarification is bounded recovery — not routine cadence — and its writes do not violate the sequential rule because they are sibling-channel writes, not predecessor-overwrites.

### 2.4. Freeze/archive rule (post-disconnect move)

At merge, review-time matter freezes; prescribed post-merge role closeouts may
still be added by their owners. The boundary has three ordered phases:

#### 2.4.1. Post-merge closeout and release validation

The directory is on main at `.cdd/unreleased/{N}/`. α and β add their own
closeouts; γ adds preliminary triage, `CDD-Post-Merge-Closeout: complete`, and
exactly one `CDD-Release-Batch: X.Y.Z` assignment.
No actor rewrites predecessor matter. δ validates this exact path and performs
the tagged disconnect while it remains there.

#### 2.4.2. Post-disconnect archive and terminal seal

Only after δ reports the tag and green release CI (or explicit override), γ moves:

```text
.cdd/unreleased/{N}/  →  .cdd/releases/{X.Y.Z}/{N}/
```

The move is a distinct post-disconnect commit on main. After that commit, γ
appends `CDD-Release-Tag: X.Y.Z`, `CDD-Archive-Commit: <sha>`, and
`CDD-Terminal-Closure: complete` to the archived `gamma-closeout.md` and
commits it separately. After terminal seal:

- `.cdd/unreleased/{N}/` no longer exists on main (the directory has moved).
- `.cdd/releases/{X.Y.Z}/{N}/` carries the cycle's complete record, version-anchored.
- Multiple cycle directories may live under one release directory when several issues ship in the same release.

The archived directory is then immutable. Corrections require a new cycle.
Docs-only disconnects use `.cdd/releases/docs/{ISO-date}/{N}/` after δ
explicitly acknowledges the already-merged main SHA as the no-tag disconnect.

The post-disconnect move followed by terminal seal is the **wire-format
invariant**; consumer-side skills define the concrete release and validation
mechanics.

### 2.5. Cross-cycle aggregation

The channel is per-cycle; cross-cycle aggregation runs through one persistent aggregator file at `.cdd/iterations/INDEX.md` at the repo root. Each cycle that produces a `cdd-iteration.md` (per the cadence rule — required when `protocol_gap_count > 0`; courtesy stub permitted when count is 0 per cnos#401) updates INDEX.md with one row naming the cycle / issue / date / findings / patches / MCAs / no-patch / path.

The aggregator's row format, the per-finding shape inside `cdd-iteration.md`, and the cross-repo trace bundle (when a finding's disposition is `patch-landed` with `Cross-repo` target) are owned by **Sub 5 of cnos#404** — the receipt-stream sub-skill ([`cnos.handoff/skills/handoff/receipt-stream/SKILL.md`](../receipt-stream/SKILL.md); landed via cnos#419). This skill names the aggregator's *existence* and *channel-relationship* (per-cycle artifact on this channel feeds the cross-cycle aggregator); the aggregation rules themselves are receipt-stream's.

---

## 3. Rules

### 3.1. Write to the channel, not to chat

- ❌ Communicate a load-bearing cycle decision in chat alone — invisible to polling roles and to the post-release record.
- ✅ Every load-bearing decision lives in the role's own channel artifact; chat is communication, not the cycle's record.

### 3.2. Per-role write ownership; no sibling overwrites

- ❌ β edits α's primary branch artifact to "clarify" what α should have written.
- ❌ α edits β's review record to address findings inline.
- ✅ Each role writes its own role-local artifacts; α's fix-round append after β's RC verdict goes in α's *own* primary artifact, not β's review record.

### 3.3. Predecessor's signaled completion is the successor's intake gate

- ❌ β begins review before α has signaled review-readiness in α's primary artifact.
- ❌ γ begins triage before both close-outs are on main.
- ✅ Each role's intake polls for the predecessor's signaled completion (review-readiness section append; close-out commit on main); the successor's intake iteration reads the predecessor's artifact live.

### 3.4. Role-owned predecessor matter freezes; prescribed closeouts remain writable

- ❌ Edit a predecessor's artifact post-merge to fit the post-merge narrative.
- ❌ Add a "clarification" to α's self-coherence.md after β has approved and merged.
- ✅ α and β add only their own prescribed closeouts; γ adds only γ-owned closeout/terminal artifacts. Review-time predecessor matter is immutable.

### 3.5. Post-disconnect archive is after the release commit

- ❌ Move before δ validates the marked receipt or before the tag/green-CI result.
- ❌ Leave the directory under `unreleased/` after γ's post-disconnect archive opportunity.
- ✅ δ disconnects while the directory is under `unreleased/`; γ then commits the move and commits terminal evidence separately.

### 3.6. After terminal seal, the directory is immutable

- ❌ Edit a file inside `.cdd/releases/{X.Y.Z}/{N}/` to correct a typo discovered later; it is the cycle's frozen post-release record.
- ✅ Corrections go in the post-release assessment, in a new cycle's close-out, or in a doc-only release with its own version-snapshot under `.cdd/releases/docs/{ISO-date}/{N}/`.

### 3.7. Cross-reference, do not duplicate

Consumer protocol skills cite this skill for the channel wire-format invariants. If you find channel-shape doctrine in `cnos.cdd/skills/cdd/`, `cnos.cds/skills/cds/CDS.md`, or elsewhere that contradicts this skill, this skill governs and the other surface is patched to reference this one. The per-artifact contract (which files exist; what each file contains; who owns each row) stays at the consumer protocol's home (cds for software cycles; cdr for research cycles); the channel-shape lives here.

---

## 4. Empirical anchors

Every cycle since the channel was established (pre-#364) uses this channel. The channel is the default substrate — no cycle ships without it. Representative anchors:

- **Historical move timing.** Pre-#669 releases placed the directory move in the release commit. Cycle #669 identified that this made exact-cycle validation and terminal semantics contradictory; current authority is §2.4's post-disconnect move.
- **Per-role write ownership.** Crystallized through cnos#287 (γ-owned branch pre-flight) and cnos#369 (the gamma-scaffold binding gate; β's role-side enforcement via review/SKILL.md rule 3.11b). The pattern is empirically stable: across the cnos#388–#412 wave, no cycle showed cross-role artifact mutation.
- **Sequential α→β→γ.** The pattern is empirically stable: across the cnos#388–#412 wave, every cycle followed the sequential pattern (with the γ-scaffold cycle-start exception). cnos#369 was the empirical anchor for the binding gate that makes the pattern protocol-enforced.
- **Mid-flight on-channel writes.** The cnos#391 rescue cycle (anchored in `cnos.handoff/skills/handoff/mid-flight/SKILL.md`) is the canonical demonstration that the channel handles non-sequential γ writes via sibling-channel files (`gamma-clarification.md`); the sequential rule is not violated.
- **#406–#412 wave** — historical evidence for directory shape and ownership, not authority for the superseded move timing.

---

## 5. Verify

### 5.1. Channel directory check (per cycle, at scaffold time)

- Does `.cdd/unreleased/{N}/` exist on `origin/cycle/{N}` after γ-scaffold commit?
- Is `.cdd/unreleased/{N}/gamma-scaffold.md` present and γ-authored?

### 5.2. Per-role write ownership check (at review and close-out times)

- Are α-owned artifacts (e.g. `self-coherence.md`, `alpha-closeout.md` for CDS-class) authored only by α-attributed commits?
- Are β-owned artifacts (`beta-review.md`, `beta-closeout.md`) authored only by β-attributed commits?
- Are γ-owned artifacts (`gamma-scaffold.md`, `gamma-closeout.md`, `cdd-iteration.md`, mid-flight files) authored only by γ-attributed commits?
- Are α's fix-round appends (after β RC) in α's *own* primary artifact, not in β's review record?

### 5.3. Sequential rule check (at review time)

- Did β begin review only after α's review-readiness signal commit?
- Did γ begin close-out triage only after both close-outs were on main?

### 5.4. Boundary-ordering check (at release and terminal close)

- Does `.cdd/unreleased/{N}/gamma-closeout.md` carry the exact nonterminal
  `CDD-Post-Merge-Closeout: complete` marker before release validation?
- Does that same receipt carry exactly one canonical `CDD-Release-Batch:` line
  naming the selected version/date batch?
- Does δ create the release tag and observe green tag CI while the cycle record
  still lives at `.cdd/unreleased/{N}/`?
- After δ reports successful disconnect, does γ move the directory to
  `.cdd/releases/{X.Y.Z}/{N}/` in a dedicated archive commit?
- Does γ append the terminal declaration only after the archive commit exists,
  binding both the release tag and archive commit?
- Are predecessor artifacts append-only under their owning role throughout?

### 5.5. Cross-cycle aggregation check (at close-out)

- If `protocol_gap_count > 0`, does `.cdd/iterations/INDEX.md` carry a row for cycle N?
- If `protocol_gap_count == 0`, is the courtesy-stub `cdd-iteration.md` present (or explicitly omitted per cnos#401 cadence rule)?

---

## 6. Related documents

- `cnos.handoff/skills/handoff/HANDOFF.md` — package contract; this skill is one sub-surface.
- `cnos.handoff/skills/handoff/mid-flight/SKILL.md` — the rescue protocol that lives on this channel. Companion sub-skill landed alongside this one under Sub 4 of cnos#404. mid-flight owns the `gamma-clarification.md` write/read rules; this skill owns the channel substrate.
- `cnos.handoff/skills/handoff/dispatch/SKILL.md` — the dispatch wire format; the dispatch prompt initiates the cycle that uses this channel. dispatch cites this skill for the `.cdd/unreleased/{N}/` per-role write ownership.
- [`cnos.handoff/skills/handoff/receipt-stream/SKILL.md`](../receipt-stream/SKILL.md) — Sub 5 of cnos#404 (landed via cnos#419). The cross-cycle aggregator (`.cdd/iterations/INDEX.md`; per-finding shape; cross-repo trace bundle). Receipt-stream consumes the per-cycle `cdd-iteration.md` artifact this channel produces.
- `cnos.handoff/skills/handoff/cross-repo/SKILL.md` — STATUS state machine; the cross-repo proposal lifecycle is orthogonal to this channel but uses a parallel directory shape (`.cdd/iterations/cross-repo/{counterpart-repo}/{slug}/`) for bilateral coordination.
- `cnos.cds/skills/cds/CDS.md §"Artifact contract"` — the CDS-class per-artifact contract: Location matrix (canonical filenames); Ownership matrix (per-role artifact assignment with verification gates); Ordered flow (13-stage flow); Frozen snapshot rule (the cds-side realization of this skill's §2.4); Trace format. CDS binds and consumes the channel pattern declared here; it does not own the channel.
- `cnos.cds/skills/cds/CDS.md §"Coordination surfaces" → "Polling primitives"` — the polling primitives that detect the cycle-branch SHA advance (the handoff signal this skill names).
- `cnos.cdd/skills/cdd/gamma/SKILL.md §2.5 → "Pre-dispatch γ scaffold check"` — the binding gate at cycle-start; γ MUST NOT dispatch α until `gamma-scaffold.md` exists on `origin/cycle/{N}`.
- `cnos.cdd/skills/cdd/alpha/SKILL.md` — α's role-local channel obligations (write `self-coherence.md` incrementally; append review-readiness as the final section; append fix-round sections after β RC).
- `cnos.cdd/skills/cdd/beta/SKILL.md` — β's role-local channel obligations (read α's primary artifact; write review record per round; write merge/review close-out before handoff to γ/δ).
- `cnos.cdd/skills/cdd/release/SKILL.md §2.5a` — post-disconnect directory archive mechanics.
- [cnos#364](https://github.com/usurobor/cnos/issues/364) — per-role artifact filename conventions empirical anchor.
- [cnos#369](https://github.com/usurobor/cnos/issues/369) — the gamma-scaffold binding gate empirical anchor.
- [cnos#401](https://github.com/usurobor/cnos/issues/401) — the courtesy-stub cdd-iteration.md cadence rule.
- [cnos#404](https://github.com/usurobor/cnos/issues/404) — parent tracker (handoff/coordination extraction wave).
- [cnos#418](https://github.com/usurobor/cnos/issues/418) — the codification cycle for this surface (Sub 4 of cnos#404).

---

## 7. Non-goals

- This skill does NOT own the per-artifact filenames, content sections, or 13-stage Ordered flow. Those are CDS-specific (and CDR-different; future c-d-X may differ further) and live in `cnos.cds/skills/cds/CDS.md §"Artifact contract"`. This skill names the channel-shape; cds names the file set.
- This skill does NOT own the mid-flight rescue mechanism. `gamma-clarification.md`'s write/read protocol, cache-bust semantics, and resumption protocol are at `cnos.handoff/skills/handoff/mid-flight/SKILL.md`. This skill names the channel that hosts the rescue file; mid-flight names the rescue's behavior.
- This skill does NOT own the cross-cycle aggregator. The per-finding shape, the INDEX.md row format, and the cross-repo trace bundle are at [`cnos.handoff/skills/handoff/receipt-stream/SKILL.md`](../receipt-stream/SKILL.md) (Sub 5 of cnos#404; landed via cnos#419). This skill names the per-cycle `cdd-iteration.md` artifact's existence and channel-relationship; receipt-stream owns the aggregation rules.
- This skill does NOT own the dispatch wire format or the implementation-contract schema. Those live in `cnos.handoff/skills/handoff/dispatch/SKILL.md`. This skill names the channel α/β/γ use after dispatch; dispatch names the prompt that initiates the cycle.
- This skill does NOT own the cross-repo proposal lifecycle (STATUS state machine; bundle file sets; LINEAGE.md schemas). Those live in `cnos.handoff/skills/handoff/cross-repo/SKILL.md`. The cross-repo bundle uses a parallel directory shape (`.cdd/iterations/cross-repo/{counterpart-repo}/{slug}/`); the two channels are orthogonal.
- This skill does NOT own the polling primitives or harness substrate. Those live in `cnos.cds/skills/cds/CDS.md §"Coordination surfaces" → "Polling primitives"` and `cnos.cdd/skills/cdd/harness/SKILL.md`. This skill assumes a polling primitive emits on cycle-branch SHA transitions; it does not specify how.
- This skill does NOT change the rescue mechanism's behavior or the artifact-file names. The current directory-move discipline is the canonical correction established by cycle #669; earlier release-commit timing remains historical evidence, not current authority.
- This skill does NOT lift the channel pattern into a typed schema. Whether to type the channel-directory contents into `schemas/handoff/artifact-channel.cue` is deferred; if a future cycle pressures the type-lift, the Markdown invariants reference the schema by `$ref` then.

---

## 8. Kata

### Scenario

γ is about to dispatch α on `cycle/N`. The cycle is a substantial CDS-class cycle (mode = substantial); γ has already drafted the dispatch prompt with a populated implementation contract per `dispatch/SKILL.md §2.3`. The channel directory `.cdd/unreleased/N/` does not yet exist on `origin/cycle/N`.

### Expected reasoning

1. **Scaffold write (cycle-start exception, §2.3 step 1).** γ writes `.cdd/unreleased/N/gamma-scaffold.md` per the per-role write ownership pattern (γ owns gamma-scaffold). γ commits and pushes to `origin/cycle/N`. The binding gate at `cnos.cdd/skills/cdd/gamma/SKILL.md §2.5` requires the scaffold's presence before α dispatch; γ has satisfied it.

2. **α phase, intake.** γ dispatches α. α's intake polls `origin/cycle/N`; the SHA advance from the scaffold commit wakes α. α reads `gamma-scaffold.md` (the scaffold is α's input per the sequential rule). α begins implementation per the dispatch prompt + scaffold.

3. **α phase, primary artifact.** α writes `self-coherence.md` incrementally on the cycle branch (one section per commit per the CDS Large-file authoring rule). The final section append is the review-readiness signal — α appends `## Review-readiness | round 1 | ...` and pushes. This is α's signaled completion; the SHA advance wakes β.

4. **β phase.** β reads α's primary artifact + the diff + the issue body (β does not modify α's primary artifact). β writes `beta-review.md` incrementally per round. On A verdict, β merges. (If RC, α appends a fix-round section to its own `self-coherence.md` — never to β's `beta-review.md` — and β re-reviews in a new round.)

5. **Post-merge close-out phase.** α is re-dispatched after merge; α writes `alpha-closeout.md`. β writes `beta-closeout.md` (in the same β session as the merge per CDS Field 6 collapse). γ verifies both close-outs are on main, records preliminary triage, and appends `CDD-Post-Merge-Closeout: complete` plus one `CDD-Release-Batch: X.Y.Z` assignment to `gamma-closeout.md`. The pair means release-ready and explicitly does not mean terminally closed.

6. **Release disconnect while unreleased.** γ prepares `RELEASE.md` and runs the default release gate while `.cdd/unreleased/N/` remains in place. δ creates the release tag and observes green tag CI. A successful δ disconnect is the authority for γ to begin archival; δ does not move the cycle directory.

7. **Post-disconnect archive and terminal close (§2.4.2, §3.6).** γ completes PRA/final triage and any `cdd-iteration.md` / `.cdd/iterations/INDEX.md` update, moves `.cdd/unreleased/N/` → `.cdd/releases/{X.Y.Z}/N/` in a dedicated archive commit, then appends the exact release-tag/archive-commit bindings and `CDD-Terminal-Closure: complete` in a later γ commit. Only then is the cycle terminal; the archived directory is immutable thereafter.

### Common failures

- γ dispatches α before writing `gamma-scaffold.md` — the binding gate at γ §2.5 catches this at γ-side; β's rule 3.11b catches it at β-side if it survives.
- β edits α's `self-coherence.md` to "clarify" something — sibling overwrite; rule 3.2 violation.
- α appends a fix-round section to β's `beta-review.md` instead of its own `self-coherence.md` — sibling overwrite; rule 3.2 violation.
- γ moves `.cdd/unreleased/N/` before δ reports a successful release disconnect — rule 3.5 violation; the exact-cycle release gate loses its input before it can authorize the tag.
- γ leaves `.cdd/unreleased/N/` in place after the post-disconnect archive opportunity — the cycle remains nonterminal and can collide with a later cycle on issue N.
- γ edits `alpha-closeout.md` post-release "to add a missing finding" — rule 3.6 violation; corrections go in a follow-on cycle.

### Verification

- The cycle directory exists on `origin/cycle/N` after γ-scaffold commit, remains under `.cdd/unreleased/N/` through merge, release validation, and δ disconnect, then moves to `.cdd/releases/{X.Y.Z}/N/` in a post-disconnect γ archive commit.
- The nonterminal marker precedes release validation; the terminal declaration follows the archive commit and binds both the tag and archive SHA.
- All artifacts attribute commits to their owning role (`α-N:` / `β-N:` / `γ-N:` prefix convention).
- No sibling-overwrite commits (a role committing to another role's owned artifact).
- Post-release, `.cdd/releases/{X.Y.Z}/N/` is unmodified.

### Reflection

A second γ reading the cycle's record (from the frozen `.cdd/releases/{X.Y.Z}/N/` directory alone) should be able to reconstruct the cycle's full history — who wrote what when, what the review verdicts were, how findings were addressed, what the close-out triage decided — without consulting chat history, operator notes, or any out-of-band source. The channel carries the complete record.

---

**End of skill.**
