---
name: artifact-channel
description: Intra-cycle artifact-channel rules. Use when a role writes a per-cycle artifact (gamma-scaffold, self-coherence, beta-review, alpha-closeout, beta-closeout, gamma-closeout, cdd-iteration, gamma-clarification) to .cdd/unreleased/{N}/ on the cycle branch, when one role reads a predecessor role's artifact as input, or when γ moves the cycle directory at release.
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
  - release-time move
scope: task-local
inputs:
  - active CDD-class cycle on origin/cycle/{N}
  - active role (α / β / γ; ε at close-out)
  - predecessor-role artifact(s) when applicable
outputs:
  - per-role artifact written at .cdd/unreleased/{N}/{filename}.md
  - SHA advance on origin/cycle/{N} (the handoff signal to the successor role)
  - at release: cycle directory moved from .cdd/unreleased/{N}/ to .cdd/releases/{X.Y.Z}/{N}/ (frozen snapshot)
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

**`.cdd/unreleased/{N}/` on the cycle branch is the sequential intra-cycle handoff substrate. Each per-cycle artifact has exactly one owning role at write time; successor roles read but do not modify predecessor artifacts (except by their own append, e.g. α's fix-round sections after β's RC verdict); the channel directory freezes on merge and moves to `.cdd/releases/{X.Y.Z}/{N}/` at release. Chat does not transit the channel; the artifact channel does.**

The artifact channel is the wire format by which α, β, and γ hand work to each other within a single cycle. Each role writes its artifacts at canonical paths under `.cdd/unreleased/{N}/`; the SHA advance on `origin/cycle/{N}` is the handoff signal to the polling successor role. The substrate is filesystem + git — no chat, no separate coordination database, no out-of-band side channel. The cycle directory carries the cycle's complete record from γ-scaffold to γ-closeout; at release-time, γ moves it from `.cdd/unreleased/{N}/` to `.cdd/releases/{X.Y.Z}/{N}/` per the frozen-snapshot rule, where it stays as the immutable post-release record.

The failure mode is **channel-substrate drift** — a role writes a coordination signal to chat instead of to the channel; a role writes to a sibling role's artifact and obscures attribution; the cycle directory is mutated after merge (e.g. an α-closeout edit dated after release); the directory is left at `.cdd/unreleased/{N}/` past release and loses its version association. The result is a cycle whose record is incomplete, whose role attribution is ambiguous, and whose post-release assessment cannot reconstruct what happened from the cycle directory alone.

This skill is the **canonical surface** for the wire-format invariants: the channel directory shape, the per-role write ownership pattern, the sequential α→β→γ flow, the frozen-snapshot rule on merge, and the release-time directory move. Consumer protocol packages bind the pattern to their specific artifact instantiation; for CDS-class cycles, the canonical *per-artifact* contract (filenames, content sections, ownership rows, verification gates, ordered 13-stage flow) lives in `cnos.cds/skills/cds/CDS.md §"Artifact contract"`. This skill governs the channel; cnos.cds governs the per-artifact contents.

## Authority

- This skill (in `cnos.handoff`) is the only canonical home for the channel-shape wire-format invariants: the channel directory path pattern (`.cdd/unreleased/{N}/`), the per-role write ownership pattern (α writes α artifacts; β writes β artifacts; γ writes γ artifacts), the sequential α→β→γ flow (each role's artifact is the successor role's input), the frozen-snapshot rule (post-merge immutability of the channel directory), the release-time directory move (cycle dir migrates to `.cdd/releases/{X.Y.Z}/{N}/` before tag). If another skill or doctrine surface carries the same wire-format claims, this skill governs.
- The consumer protocol packages govern the per-artifact instantiation. For CDS-class cycles, the canonical home is `cnos.cds/skills/cds/CDS.md §"Artifact contract"` — specifically the Location matrix (canonical filename per artifact), Ownership matrix (per-role artifact assignment with verification gates), and Ordered flow (the 13-stage flow: design → contract → plan → tests → code → docs → self-coherence → review → gate → release → observe → assess → close). CDR-class cycles bind differently; future c-d-X may bind differently still. The wire-format pattern is what this skill owns; the per-consumer instantiation is what cds/cdr/c-d-X own.
- The artifacts on this channel are also the substrate for sibling sub-skills' wire formats. `cnos.handoff/skills/handoff/mid-flight/SKILL.md` (gamma-clarification.md / gamma-coordination.md) and `cnos.handoff/skills/handoff/receipt-stream/SKILL.md` (cdd-iteration.md aggregator) live on this channel. Their authoring/reader/cadence rules are theirs; the channel-shape they live on is this skill's.

## Scope

In scope:

- The channel directory path pattern (`.cdd/unreleased/{N}/`); the cycle-branch substrate (`origin/cycle/{N}`); the SHA-advance handoff signal.
- The per-role write ownership pattern (which role writes which class of artifact; what successor roles may and may not do to predecessor artifacts).
- The sequential rule (α before β; β before γ; γ-scaffold exception at cycle-start).
- The frozen-snapshot rule on merge (post-merge immutability of the channel directory; the release-time directory move from `.cdd/unreleased/{N}/` to `.cdd/releases/{X.Y.Z}/{N}/`).
- The cross-cycle aggregation pointer (per-cycle `cdd-iteration.md` → `.cdd/iterations/INDEX.md`; aggregator rules canonicalize in Sub 5 / receipt-stream/SKILL.md).
- Empirical anchors (every cycle since the channel was established in pre-#364 form).

Out of scope:

- The per-artifact filenames and content sections. Those live in `cnos.cds/skills/cds/CDS.md §"Artifact contract" → "Location matrix"` (CDS-specific filename set: `self-coherence.md`, `alpha-closeout.md`, `beta-review.md`, `beta-closeout.md`, `gamma-closeout.md`, `gamma-scaffold.md`, `cdd-iteration.md`, `RELEASE.md`, version-snapshot directory, PRA path, cross-repo trace dir). CDR and future c-d-X bind differently.
- The per-artifact role ownership table (who writes which file, when, what the verification gate is, what missing means). Lives in `cnos.cds/skills/cds/CDS.md §"Artifact contract" → "Ownership matrix"`.
- The CDS-specific 13-stage Ordered flow. Lives in `cnos.cds/skills/cds/CDS.md §"Artifact contract" → "Ordered flow"`. CDR has a different flow; future c-d-X may have other flows.
- The mid-flight rescue mechanism (`gamma-clarification.md` write/read protocol; cache-bust semantics). Lives in `cnos.handoff/skills/handoff/mid-flight/SKILL.md`. This skill owns the channel; mid-flight owns the rescue protocol that lives on the channel.
- The cross-cycle aggregator content (per-finding shape; INDEX.md row format; cross-repo trace bundle). Lives in `cnos.handoff/skills/handoff/receipt-stream/SKILL.md` (Sub 5 of cnos#404; forthcoming). This skill owns the per-cycle channel; receipt-stream owns the cross-cycle aggregation.
- The harness substrate (polling primitives, Monitor wrapping, reachability re-probe). Lives in `cnos.cdd/skills/cdd/harness/SKILL.md`. This skill assumes a polling primitive emits on cycle-branch SHA transitions.
- Tooling that auto-verifies channel discipline (`cn cdd verify` channel-integrity checks). Deferred; doctrine first, tooling later.

---

## 1. Define

### 1.1. What the channel is

The channel is **one directory per cycle**: `.cdd/unreleased/{N}/` on `origin/cycle/{N}`, where `{N}` is the cycle's issue number. Inside the directory, files are role-distinguished by **filename**, not by sub-directory. The directory carries the cycle's complete record from γ-scaffold (at cycle-start) through γ-closeout (at close); at release-time, γ moves the directory to `.cdd/releases/{X.Y.Z}/{N}/` per the frozen-snapshot rule (§2.4).

The substrate is filesystem + git. The handoff signal between roles is the SHA advance on `origin/cycle/{N}` — the polling successor role's polling primitive (per `cnos.cds/skills/cds/CDS.md §"Coordination surfaces" → "Polling primitives"`) emits on the SHA transition and the successor's next intake iteration reads the predecessor's artifact live. No chat, no out-of-band coordination database, no separate inbox.

### 1.2. What the channel is not

The channel is **not**:

- a chat-coordination substitute. Chat does not transit the channel; if a load-bearing fact is communicated only in chat, it is invisible to polling roles and to the post-release record.
- a shared-write directory. Each artifact has one owning role at write time. Sibling roles read but do not modify (except by their own append, e.g. α's fix-round sections after β's RC verdict — α appends to α's own `self-coherence.md`, never to β's `beta-review.md`).
- a multi-cycle channel. One directory = one cycle. Cross-cycle aggregation runs through `.cdd/iterations/INDEX.md` per the receipt-stream sub-skill (Sub 5 of cnos#404; forthcoming).
- mutable post-merge. Once the cycle merges, the channel directory freezes. Subsequent edits to predecessor artifacts are reviewable as protocol violations — corrections go in a new cycle's close-out or in the post-release assessment.
- the canonical post-release snapshot. The frozen post-release snapshot lives under `docs/{tier}/{bundle}/{X.Y.Z}/` per CDS Location matrix; the channel directory at `.cdd/releases/{X.Y.Z}/{N}/` is the cycle's role-local close-out evidence storage.

### 1.3. Failure mode

The channel fails through **channel-substrate drift**:

- A role writes a coordination signal to chat instead of to the channel — the signal is invisible to polling sessions and to the post-release assessment.
- A role writes to a sibling role's artifact — attribution becomes ambiguous; β's review verdict and α's self-coherence are no longer distinguishable.
- The cycle directory is mutated after merge — predecessor artifacts are silently edited to fit the post-merge narrative; the cycle's record is no longer trustworthy.
- The directory is left at `.cdd/unreleased/{N}/` past release — it loses its version association; future cycles operating on the same issue number (re-opens, follow-ons) collide with the stale directory.
- A coordination note lives only in the GitHub issue comments or in a chat backlog — the artifact channel is incomplete and the post-release reconstruction is impossible.

The protocol closes each failure mode by naming the channel-shape invariants explicitly and binding role-write ownership, sequential flow, and frozen-snapshot rules to the channel directory.

---

## 2. Unfold

### 2.1. Channel directory

The canonical channel path is:

```text
.cdd/unreleased/{N}/
```

on `origin/cycle/{N}`, where `{N}` is the cycle's GitHub issue number. The directory exists for the duration of the cycle — from γ-scaffold (at cycle-start) through release-time move (at release). The path is consumer-agnostic (CDS, CDR, future c-d-X all use the same channel-directory shape); only the *filenames* inside vary per consumer protocol.

The cycle-branch substrate is `origin/cycle/{N}` — a single named branch per cycle, created by γ from `origin/main` at cycle-start per `cnos.cds/skills/cds/CDS.md §"Development lifecycle" → "Branch rule"`. The channel directory lives on the cycle branch (not on main) until the release-time move; polling main for in-flight cycle dirs is silent.

The handoff signal between roles is the **SHA advance on `origin/cycle/{N}`**. Each commit to the cycle branch advances the SHA; polling primitives detect the transition and wake successor roles. The signal is git-native — no separate notification channel, no webhook, no inbox. The SHA advance is the only event the successor role needs.

### 2.2. Per-role write ownership

Every per-cycle artifact has exactly one owning role at write time. The pattern (consumer-protocol-agnostic):

| Owner role | Artifact class | Pattern |
|---|---|---|
| γ | scaffold | Cycle-start artifact; carries γ's selection decision and the dispatch surface set. One per cycle. |
| α | role-local primary branch artifact | The role's review-readiness signal; gap / mode / ACs / trace / fix-rounds. One per cycle. |
| α | role close-out | The role's post-merge findings record. One per cycle. |
| β | role-local review record | The review CLP rounds + verdicts. One per cycle; one section per round. |
| β | role close-out | The role's post-review close-out + release evidence. One per cycle. |
| γ | role close-out | The cycle's closure declaration + triage record. One per cycle. |
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
5. **β close-out** — β writes its close-out artifact (release evidence + close-out findings) in the same β session as the merge (per CDS Field 6 collapse). β exits.
6. **γ phase** — γ verifies both close-outs are on main; writes the per-cycle receipt artifact (when applicable); writes the cycle's closure declaration; moves the cycle directory at release-time per §2.4.

The successor role's intake **requires** the predecessor role's signaled completion. β does not begin review until α has signaled review-readiness (in α's primary artifact). γ does not begin triage until both close-outs are on main. Each role's polling primitive emits when the predecessor's relevant artifact reaches the expected state; the role's intake iteration reads the artifact live.

The two exceptions to the strictly-sequential pattern:

- **The γ-scaffold cycle-start write** (§2.3 step 1) — γ writes before α exists in the cycle. Not strictly sequential because α has not yet written anything; γ-scaffold is the cycle's initial state.
- **Mid-flight γ writes** (`gamma-clarification.md` / `gamma-coordination.md`) — γ writes mid-cycle while α and/or β are in-flight. This is the rescue channel (per `cnos.handoff/skills/handoff/mid-flight/SKILL.md`); the channel is the substrate but the sequence is not strict α→β→γ. The clarification is bounded recovery — not routine cadence — and its writes do not violate the sequential rule because they are sibling-channel writes, not predecessor-overwrites.

### 2.4. Frozen-snapshot rule (release-time move)

When the cycle merges, the channel directory **freezes** — no further writes to predecessor artifacts. The frozen rule has two phases:

#### 2.4.1. Post-merge, pre-release (in-version state)

Between β merge and γ tag/release, the channel directory at `.cdd/unreleased/{N}/` is on `origin/main` (β merged it). γ writes close-out and (when applicable) receipt artifacts in this window; these are the *only* permitted writes — and they are γ-authored writes to γ-owned files (gamma-closeout.md, cdd-iteration.md, mid-flight files when a rescue event fires post-merge), not edits to predecessor artifacts.

#### 2.4.2. Release-time move

Before γ requests the tag from δ, γ moves the channel directory:

```text
.cdd/unreleased/{N}/  →  .cdd/releases/{X.Y.Z}/{N}/
```

The move is part of the release commit on main (per `cnos.cdd/skills/cdd/release/SKILL.md §2.5a` and `cnos.cds/skills/cds/CDS.md §"Artifact contract" → "Location matrix"`). After the move:

- `.cdd/unreleased/{N}/` no longer exists on main (the directory has moved).
- `.cdd/releases/{X.Y.Z}/{N}/` carries the cycle's complete record, version-anchored.
- Multiple cycle directories may live under one release directory when several issues ship in the same release.

After the release commit, the directory is **immutable**. Subsequent edits to files inside `.cdd/releases/{X.Y.Z}/{N}/` are reviewable as protocol violations — corrections go in a new cycle's close-out, in the post-release assessment, or in a documentation-only release with its own version-snapshot. Doc-only releases use `.cdd/releases/docs/{ISO-date}/{N}/` as the destination directory shape per CDS Location matrix.

The release-time move is the **wire-format invariant**; the per-consumer realization (which release script runs the move; how `cn cdd verify` checks the move happened; how the legacy aggregate paths are detected) is consumer-side at cnos.cds + cnos.cdd/release skills.

### 2.5. Cross-cycle aggregation

The channel is per-cycle; cross-cycle aggregation runs through one persistent aggregator file at `.cdd/iterations/INDEX.md` at the repo root. Each cycle that produces a `cdd-iteration.md` (per the cadence rule — required when `protocol_gap_count > 0`; courtesy stub permitted when count is 0 per cnos#401) updates INDEX.md with one row naming the cycle / issue / date / findings / patches / MCAs / no-patch / path.

The aggregator's row format, the per-finding shape inside `cdd-iteration.md`, and the cross-repo trace bundle (when a finding's disposition is `patch-landed` with `Cross-repo` target) are owned by **Sub 5 of cnos#404** — the receipt-stream sub-skill (`cnos.handoff/skills/handoff/receipt-stream/SKILL.md`; forthcoming). This skill names the aggregator's *existence* and *channel-relationship* (per-cycle artifact on this channel feeds the cross-cycle aggregator); the aggregation rules themselves are receipt-stream's.

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

### 3.4. Channel directory freezes on merge

- ❌ Edit a predecessor's artifact post-merge to fit the post-merge narrative.
- ❌ Add a "clarification" to α's self-coherence.md after β has approved and merged.
- ✅ Post-merge writes are γ-authored writes to γ-owned files (close-out, iteration receipt, mid-flight files when a rescue event fires post-merge). Predecessor artifacts are immutable.

### 3.5. Release-time move is part of the release commit

- ❌ Leave the cycle directory at `.cdd/unreleased/{N}/` after release; future cycles collide on the same issue number.
- ❌ Move the directory in a separate post-release commit, after the tag fires.
- ✅ γ moves `.cdd/unreleased/{N}/` → `.cdd/releases/{X.Y.Z}/{N}/` in the same release commit as `RELEASE.md`; the move is on main before γ requests the tag from δ.

### 3.6. Post-release, the directory is immutable

- ❌ Edit a file inside `.cdd/releases/{X.Y.Z}/{N}/` to correct a typo discovered later; it is the cycle's frozen post-release record.
- ✅ Corrections go in the post-release assessment, in a new cycle's close-out, or in a doc-only release with its own version-snapshot under `.cdd/releases/docs/{ISO-date}/{N}/`.

### 3.7. Cross-reference, do not duplicate

Consumer protocol skills cite this skill for the channel wire-format invariants. If you find channel-shape doctrine in `cnos.cdd/skills/cdd/`, `cnos.cds/skills/cds/CDS.md`, or elsewhere that contradicts this skill, this skill governs and the other surface is patched to reference this one. The per-artifact contract (which files exist; what each file contains; who owns each row) stays at the consumer protocol's home (cds for software cycles; cdr for research cycles); the channel-shape lives here.

---

## 4. Empirical anchors

Every cycle since the channel was established (pre-#364) uses this channel. The channel is the default substrate — no cycle ships without it. Representative anchors:

- **The `.cdd/unreleased/{N}/` → `.cdd/releases/{X.Y.Z}/{N}/` directory move pattern.** Standardized in pre-#283 form; refined through cnos#283–#287 (cycle-branch naming convention) and cnos#364 (per-role artifact-file conventions). Every release commit since cnos#364 carries the move as part of the release commit.
- **Per-role write ownership.** Crystallized through cnos#287 (γ-owned branch pre-flight) and cnos#369 (the gamma-scaffold binding gate; β's role-side enforcement via review/SKILL.md rule 3.11b). The pattern is empirically stable: across the cnos#388–#412 wave, no cycle showed cross-role artifact mutation.
- **Sequential α→β→γ.** The pattern is empirically stable: across the cnos#388–#412 wave, every cycle followed the sequential pattern (with the γ-scaffold cycle-start exception). cnos#369 was the empirical anchor for the binding gate that makes the pattern protocol-enforced.
- **Mid-flight on-channel writes.** The cnos#391 rescue cycle (anchored in `cnos.handoff/skills/handoff/mid-flight/SKILL.md`) is the canonical demonstration that the channel handles non-sequential γ writes via sibling-channel files (`gamma-clarification.md`); the sequential rule is not violated.
- **#406–#412 wave** — the post-codification cycles that exercised the channel in steady-state. Every cycle in the wave produced a complete `.cdd/unreleased/{N}/` directory and a clean release-time move; the channel's wire-format is operationally stable.

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

### 5.4. Frozen-snapshot check (at release time)

- Is `.cdd/unreleased/{N}/` moved to `.cdd/releases/{X.Y.Z}/{N}/` in the release commit (same commit as `RELEASE.md`)?
- Are there any post-merge edits to predecessor artifacts? If yes, are they γ-authored writes to γ-owned files (permitted) or sibling overwrites (protocol violations)?
- Post-release, is `.cdd/releases/{X.Y.Z}/{N}/` immutable (no further commits touching files inside)?

### 5.5. Cross-cycle aggregation check (at close-out)

- If `protocol_gap_count > 0`, does `.cdd/iterations/INDEX.md` carry a row for cycle N?
- If `protocol_gap_count == 0`, is the courtesy-stub `cdd-iteration.md` present (or explicitly omitted per cnos#401 cadence rule)?

---

## 6. Related documents

- `cnos.handoff/skills/handoff/HANDOFF.md` — package contract; this skill is one sub-surface.
- `cnos.handoff/skills/handoff/mid-flight/SKILL.md` — the rescue protocol that lives on this channel. Companion sub-skill landed alongside this one under Sub 4 of cnos#404. mid-flight owns the `gamma-clarification.md` write/read rules; this skill owns the channel substrate.
- `cnos.handoff/skills/handoff/dispatch/SKILL.md` — the dispatch wire format; the dispatch prompt initiates the cycle that uses this channel. dispatch cites this skill for the `.cdd/unreleased/{N}/` per-role write ownership.
- `cnos.handoff/skills/handoff/receipt-stream/SKILL.md` — Sub 5 of cnos#404 (forthcoming). The cross-cycle aggregator (`.cdd/iterations/INDEX.md`; per-finding shape; cross-repo trace bundle). Receipt-stream consumes the per-cycle `cdd-iteration.md` artifact this channel produces.
- `cnos.handoff/skills/handoff/cross-repo/SKILL.md` — STATUS state machine; the cross-repo proposal lifecycle is orthogonal to this channel but uses a parallel directory shape (`.cdd/iterations/cross-repo/{counterpart-repo}/{slug}/`) for bilateral coordination.
- `cnos.cds/skills/cds/CDS.md §"Artifact contract"` — the CDS-class per-artifact contract: Location matrix (canonical filenames); Ownership matrix (per-role artifact assignment with verification gates); Ordered flow (13-stage flow); Frozen snapshot rule (the cds-side realization of this skill's §2.4); Trace format. CDS binds and consumes the channel pattern declared here; it does not own the channel.
- `cnos.cds/skills/cds/CDS.md §"Coordination surfaces" → "Polling primitives"` — the polling primitives that detect the cycle-branch SHA advance (the handoff signal this skill names).
- `cnos.cdd/skills/cdd/gamma/SKILL.md §2.5 → "Pre-dispatch γ scaffold check"` — the binding gate at cycle-start; γ MUST NOT dispatch α until `gamma-scaffold.md` exists on `origin/cycle/{N}`.
- `cnos.cdd/skills/cdd/alpha/SKILL.md` — α's role-local channel obligations (write `self-coherence.md` incrementally; append review-readiness as the final section; append fix-round sections after β RC).
- `cnos.cdd/skills/cdd/beta/SKILL.md` — β's role-local channel obligations (read α's primary artifact; write review record per round; write close-out + release evidence).
- `cnos.cdd/skills/cdd/release/SKILL.md §2.5a` — the release-time directory move mechanics (the cdd-side realization of §2.4.2).
- [cnos#364](https://github.com/usurobor/cnos/issues/364) — per-role artifact filename conventions empirical anchor.
- [cnos#369](https://github.com/usurobor/cnos/issues/369) — the gamma-scaffold binding gate empirical anchor.
- [cnos#401](https://github.com/usurobor/cnos/issues/401) — the courtesy-stub cdd-iteration.md cadence rule.
- [cnos#404](https://github.com/usurobor/cnos/issues/404) — parent tracker (handoff/coordination extraction wave).
- [cnos#418](https://github.com/usurobor/cnos/issues/418) — the codification cycle for this surface (Sub 4 of cnos#404).

---

## 7. Non-goals

- This skill does NOT own the per-artifact filenames, content sections, or 13-stage Ordered flow. Those are CDS-specific (and CDR-different; future c-d-X may differ further) and live in `cnos.cds/skills/cds/CDS.md §"Artifact contract"`. This skill names the channel-shape; cds names the file set.
- This skill does NOT own the mid-flight rescue mechanism. `gamma-clarification.md`'s write/read protocol, cache-bust semantics, and resumption protocol are at `cnos.handoff/skills/handoff/mid-flight/SKILL.md`. This skill names the channel that hosts the rescue file; mid-flight names the rescue's behavior.
- This skill does NOT own the cross-cycle aggregator. The per-finding shape, the INDEX.md row format, and the cross-repo trace bundle are at `cnos.handoff/skills/handoff/receipt-stream/SKILL.md` (Sub 5 of cnos#404; forthcoming). This skill names the per-cycle `cdd-iteration.md` artifact's existence and channel-relationship; receipt-stream owns the aggregation rules.
- This skill does NOT own the dispatch wire format or the implementation-contract schema. Those live in `cnos.handoff/skills/handoff/dispatch/SKILL.md`. This skill names the channel α/β/γ use after dispatch; dispatch names the prompt that initiates the cycle.
- This skill does NOT own the cross-repo proposal lifecycle (STATUS state machine; bundle file sets; LINEAGE.md schemas). Those live in `cnos.handoff/skills/handoff/cross-repo/SKILL.md`. The cross-repo bundle uses a parallel directory shape (`.cdd/iterations/cross-repo/{counterpart-repo}/{slug}/`); the two channels are orthogonal.
- This skill does NOT own the polling primitives or harness substrate. Those live in `cnos.cds/skills/cds/CDS.md §"Coordination surfaces" → "Polling primitives"` and `cnos.cdd/skills/cdd/harness/SKILL.md`. This skill assumes a polling primitive emits on cycle-branch SHA transitions; it does not specify how.
- This skill does NOT change the rescue mechanism's behavior, the artifact-file names, or the directory-move discipline. Sub 4 of cnos#404 is a wholesale extraction of an empirically-stable channel pattern; the wire-format invariants are preserved verbatim from the empirical practice.
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

5. **Close-out phase.** α is re-dispatched after merge; α writes `alpha-closeout.md`. β writes `beta-closeout.md` (in the same β session as the merge per CDS Field 6 collapse). γ verifies both close-outs are on main; writes `gamma-closeout.md`; if `protocol_gap_count > 0` (or as a courtesy stub per cnos#401), writes `cdd-iteration.md` and adds a row to `.cdd/iterations/INDEX.md`.

6. **Release-time move (§2.4.2).** Before γ requests the tag from δ, γ moves `.cdd/unreleased/N/` → `.cdd/releases/{X.Y.Z}/N/` as part of the release commit (same commit as `RELEASE.md`). After the move, `.cdd/unreleased/N/` no longer exists on main; the cycle's record is version-anchored at `.cdd/releases/{X.Y.Z}/N/`.

7. **Post-release, immutable (§3.6).** The directory at `.cdd/releases/{X.Y.Z}/N/` is immutable. Any subsequent correction goes in a follow-on cycle's close-out, in the post-release assessment, or in a doc-only release with its own version-snapshot.

### Common failures

- γ dispatches α before writing `gamma-scaffold.md` — the binding gate at γ §2.5 catches this at γ-side; β's rule 3.11b catches it at β-side if it survives.
- β edits α's `self-coherence.md` to "clarify" something — sibling overwrite; rule 3.2 violation.
- α appends a fix-round section to β's `beta-review.md` instead of its own `self-coherence.md` — sibling overwrite; rule 3.2 violation.
- γ leaves `.cdd/unreleased/N/` on main after release — rule 3.5 violation; future cycle on issue N collides with stale directory.
- γ edits `alpha-closeout.md` post-release "to add a missing finding" — rule 3.6 violation; corrections go in a follow-on cycle.

### Verification

- The cycle directory exists on `origin/cycle/N` after γ-scaffold commit; exists on main after merge; moved to `.cdd/releases/{X.Y.Z}/N/` in the release commit.
- All artifacts attribute commits to their owning role (`α-N:` / `β-N:` / `γ-N:` prefix convention).
- No sibling-overwrite commits (a role committing to another role's owned artifact).
- Post-release, `.cdd/releases/{X.Y.Z}/N/` is unmodified.

### Reflection

A second γ reading the cycle's record (from the frozen `.cdd/releases/{X.Y.Z}/N/` directory alone) should be able to reconstruct the cycle's full history — who wrote what when, what the review verdicts were, how findings were addressed, what the close-out triage decided — without consulting chat history, operator notes, or any out-of-band source. The channel carries the complete record.

---

**End of skill.**
