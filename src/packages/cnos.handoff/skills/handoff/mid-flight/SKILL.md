---
name: mid-flight
description: Mid-flight γ→in-flight-α rescue channel. Use when γ needs to change a load-bearing fact (pinned axis, scope, AC, blocker) after α has begun coding, when γ edits the cycle's GitHub issue body mid-cycle, or when a spec change lands on origin/main while α/β polling sessions are in-flight.
artifact_class: skill
kata_surface: embedded
governing_question: How does γ change a load-bearing fact mid-cycle without abandoning the branch — naming the change in the cycle's artifact channel so in-flight α/β cache-bust and re-converge, rather than chat-relaying or silently mutating the contract?
visibility: public
parent: handoff
triggers:
  - mid-flight
  - gamma-clarification
  - issue-edit
  - cache-bust
  - rescue
  - re-pin
  - spec-staleness
scope: task-local
inputs:
  - active CDD-class cycle with α and/or β in-flight on origin/cycle/{N}
  - the load-bearing change γ needs to signal (axis re-pin, scope clarification, AC clarification, blocker resolution, issue-body edit, spec change on origin/main)
outputs:
  - .cdd/unreleased/{N}/gamma-clarification.md entry on origin/cycle/{N} (substantive content of the change)
  - cycle-branch SHA transition (the wake-up signal)
  - .cdd/unreleased/{N}/gamma-coordination.md entry (legacy long-lived polling sessions only; spec-staleness propagation)
requires:
  - cycle/{N} exists on origin with at least one in-flight role
  - the load-bearing change has been (or is about to be) committed in some other surface (issue body, origin/main spec file) that the in-flight role's cached read would miss
calls:
  - HANDOFF.md
  - artifact-channel/SKILL.md
  - dispatch/SKILL.md
  - cross-repo/SKILL.md
---

# Mid-Flight Rescue — γ → in-flight α/β re-pinning channel

## Core principle

**A mid-flight change to a load-bearing fact transits the cycle's artifact channel, not chat. γ writes the substance to `.cdd/unreleased/{N}/gamma-clarification.md` on `cycle/{N}` *before* the change becomes visible elsewhere; the cycle-branch SHA transition is the wake-up signal that bursts in-flight α/β caches. The mechanism is bounded recovery — a structural escape valve when γ under-specified at dispatch time — not a substitute for clean dispatch.**

The mid-flight rescue is the wire format γ uses when something the cycle relied on changes after dispatch. The thing that changed may be (1) a pinned implementation-contract axis γ realizes was wrong, (2) an AC γ realizes was ambiguous, (3) a non-goal γ realizes was load-bearing, (4) the issue body itself (γ edited it because α surfaced an ambiguity or δ refined an axis), or (5) a spec file on `origin/main` whose change ripples into the cycle's loaded skill snapshot. In all cases, an in-flight α or β session may be operating on a cached read of the prior state — chat will not reach them; only artifact-channel writes will.

The failure mode is **silent re-pin** — γ edits the issue body, or pushes a spec change to main, or mutates a contract row, without writing `gamma-clarification.md`. In-flight α/β keep working from cached state, finish on the old contract, and ship the wrong shape. β's Rule 7 verification cannot adjudicate "which contract is canonical" when γ's mutation left no trail. The result is a cycle that satisfies the *new* contract by accident (or doesn't), with no audit trail for which decision happened when.

This skill is the **canonical surface** for the rescue mechanism: file path, authoring role, reader role, trigger conditions, cache-bust semantics, and the empirical anchors. Consumer protocol packages (cnos.cdd, cnos.cds, future c-d-X) cite this skill; the **operational realization** (the `gh issue edit` + `git commit` + `git push` shell sequence γ runs at clarification time) stays in `cnos.cds/skills/cds/CDS.md §"Coordination surfaces" → "Mid-flight clarification"` as the v0.1 overlay. The wire-format invariant lives here; the v0.1 software-cycle realization stays at cds.

## Authority

- This skill (in `cnos.handoff`) is the only canonical home for the mid-flight rescue mechanism's wire format: the file path, the authoring/reader roles, the trigger conditions, the cache-bust function, the spec-staleness propagation invariant. If another skill or doctrine surface carries the same wire-format claims, this skill governs.
- The consumer protocol packages govern the operational realization: `cnos.cds/skills/cds/CDS.md §"Coordination surfaces" → "Mid-flight clarification"` carries the post-#411 software-cycle realization (the shell sequence and the polling mapping); `cnos.cdd/skills/cdd/gamma/SKILL.md §2.5 → "Spec-staleness propagation"` carries the consumer-specific file list (which cdd skill files trigger propagation). Both cite this skill for the wire-format invariant.
- The `gamma-clarification.md` artifact lives on the channel substrate; channel rules (per-role write ownership, frozen-snapshot rule on merge, sequential α→β→γ flow) are at `cnos.handoff/skills/handoff/artifact-channel/SKILL.md`. This skill cites the channel; the channel does not redefine the rescue protocol.

## Scope

In scope:

- The `gamma-clarification.md` file format and authoring contract (file path, authoring role, content sections, commit + push sequence).
- The wake-up / cache-bust semantics (cycle-branch SHA transition is the signal; not the GitHub issue mtime, which is invisible to MCP-cached reads).
- Reader role obligations (in-flight α/β polling sessions detect the write; re-read the issue body live via `gh issue view`; re-load any other source γ named in the clarification).
- The resumption protocol (in-flight α/β acknowledge in their own artifact, re-scope, escalate, or stop).
- The spec-staleness propagation invariant (γ writes a coordination note when spec changes land on origin/main mid-cycle).
- Empirical anchors (cnos#391 rescue; cnos#392 recovery).

Out of scope:

- The cycle-lifecycle role-skill obligations of α (intake; pre-implementation surface; review-readiness signal) and β (review CLP; verdict cadence). Those live in `cnos.cdd/skills/cdd/{alpha,beta}/SKILL.md` and are consumers of this wire format.
- The per-artifact channel rules (write ownership, sequential flow, frozen-snapshot). Those live in `cnos.handoff/skills/handoff/artifact-channel/SKILL.md`. This skill cites the channel as the substrate `gamma-clarification.md` lives on.
- The dispatch-time implementation-contract enrichment (δ-as-inward-membrane review before routing). That lives in `cnos.handoff/skills/handoff/dispatch/SKILL.md §3`. This skill picks up where dispatch leaves off — once α is routed, dispatch's enforcement window closes; mid-flight is the post-dispatch channel.
- The operational shell sequence γ runs (`gh issue edit` + `git commit` + `git push`). That stays at `cnos.cds/skills/cds/CDS.md §"Coordination surfaces" → "Mid-flight clarification"` as the v0.1 software-cycle realization overlay.
- The harness substrate (polling primitives, Monitor wrapping, reachability re-probe). That lives in `cnos.cdd/skills/cdd/harness/SKILL.md` per Phase 4b of cnos#366.
- The cross-repo proposal lifecycle. Cited from this skill (issue-body edits propagating to in-flight intake sessions are the same mechanism); the lifecycle itself is at `cnos.handoff/skills/handoff/cross-repo/SKILL.md`.

---

## 1. Define

### 1.1. What a mid-flight rescue is

A mid-flight rescue is **one γ-authored file** (`gamma-clarification.md`) written to the cycle's artifact-channel directory (`.cdd/unreleased/{N}/`) on the cycle branch (`origin/cycle/{N}`), committed and pushed *before* the change becomes visible in any other surface the in-flight role might read. It carries:

- the change (axis re-pin / scope clarification / AC clarification / blocker resolution / issue-body edit summary / spec-change pointer),
- the date,
- the impact (which ACs / non-goals / constraints / artifacts the change affects),
- the rationale (typically: a repo-convention emerged; a downstream consumer surfaced; an operator override; α surfaced an ambiguity; β found a wrong-shape implementation; a spec change landed on main).

The cycle-branch SHA transition is the **only** signal in-flight roles need. They are already polling the branch per `cnos.cds/skills/cds/CDS.md §"Coordination surfaces" → "Polling primitives"`. The SHA advance wakes them; their next iteration re-reads the issue body (and any other named source) live, not from cache.

### 1.2. What a mid-flight rescue is not

A mid-flight rescue is **not**:

- a substitute for clean dispatch (the dispatch-time implementation-contract enrichment at `dispatch/SKILL.md §3` is the *first* line of defense; the rescue is the escape valve when the first line failed),
- a chat-relayed instruction (chat does not reach polling roles; the artifact channel does),
- a silent contract mutation (γ MUST NOT edit the issue body or push a spec change without first writing `gamma-clarification.md` — silent re-pin is the canonical failure mode),
- an authority override of α's intake constraints (α still surfaces unpinned rows to γ/δ; the rescue is γ's tool, not α's permission to improvise),
- a routine cadence (a cycle with > 1–2 `gamma-clarification.md` entries is a signal that dispatch under-specified; ε reads the pattern and codifies).

### 1.3. Failure mode

The mid-flight channel fails through **silent re-pin**:

- γ edits the issue body to clarify an ambiguity (or pin a missing axis, or fold in an operator override).
- γ does not write `gamma-clarification.md` to the cycle branch.
- The in-flight α or β session continues on its cached `gh issue view` read (the MCP cache; the local snapshot; whatever the polling primitive serves).
- α finishes on the old contract; β APPROVE-s on the old contract; the cycle merges in the wrong shape.
- No audit trail explains which decision happened when; the post-release assessment cannot reconstruct the contract history.

The canonical anchor is cnos#391: γ recovered the cycle by editing the issue body to pin missing axes (`Package scoping`, `CLI integration target`) and committing `gamma-clarification.md` to the cycle branch. α detected the SHA transition via cycle-branch polling, re-read the issue body live, and re-shaped the implementation. The rescue worked specifically because the clarification was on-channel; the issue-edit alone would have been invisible.

---

## 2. Unfold

### 2.1. File format

The clarification file lives at the canonical channel path:

```text
.cdd/unreleased/{N}/gamma-clarification.md
```

on the cycle branch (`origin/cycle/{N}`). When multiple clarifications fire in one cycle, γ **appends** to the existing file (one `## Clarification — {ISO-date}` section per event); γ does not author parallel files.

**Required sections per event:**

```markdown
## Clarification — {ISO-date}

**Change:** {one-line summary of the load-bearing change}

**Trigger:** {axis re-pin | scope clarification | AC clarification | blocker resolution | issue-body edit | spec change on main | operator override | α-surfaced ambiguity | β-surfaced wrong-shape}

**Impact:**
- ACs affected: {list or "none"}
- Non-goals affected: {list or "none"}
- Constraints affected: {list or "none"}
- Artifacts affected: {list or "none — γ-authored clarification only"}

**Rationale:** {1–3 sentences naming why the change is needed; cite the surface that surfaced the gap when applicable}

**Resumption guidance:** {what in-flight α/β should do next — re-scope / acknowledge / continue / stop and escalate}
```

The format is small by design: the rescue is bounded recovery, not a re-write of the cycle's substance. If a clarification needs more than ~30 lines, γ is probably mid-redesigning the cycle's intent and should consider stopping and re-scaffolding instead.

### 2.2. Authoring role

`gamma-clarification.md` is **γ-authored only**. The cycle's γ is the writer; the operator-as-γ may write when γ is collapsed onto the operator (per `ROLES.md §4` role-collapse). α and β MUST NOT write to this file — α surfaces gaps via `self-coherence.md` and β surfaces gaps via `beta-review.md` per the per-role write ownership rules at `cnos.handoff/skills/handoff/artifact-channel/SKILL.md §"Per-role write ownership"`.

The write sequence is:

1. γ identifies the load-bearing change (axis re-pin, scope clarification, AC clarification, blocker resolution, issue-body edit, spec change on origin/main).
2. γ writes (or appends to) `.cdd/unreleased/{N}/gamma-clarification.md` on `cycle/{N}` per the §2.1 format.
3. γ commits and pushes to `origin/cycle/{N}` — this is the wake-up signal.
4. *Only then* γ effects the other-surface change (`gh issue edit`, `git push` to main for spec change, etc.). Order matters: the clarification on the cycle branch is the signal that the other surface is about to change.

### 2.3. Reader role

In-flight α and β are the **readers**. Both roles poll `origin/cycle/{N}` per the polling primitives at `cnos.cds/skills/cds/CDS.md §"Coordination surfaces" → "Polling primitives"`. When their polling sees a `gamma-clarification.md` add or update (the SHA advance carries a touch to this file), the reader role's next intake iteration:

1. Reads the clarification file (`git show origin/cycle/{N}:.cdd/unreleased/{N}/gamma-clarification.md` or `git pull` + local read).
2. Re-reads the issue body **live** via `gh issue view {N} --json title,body,state,comments` (the MCP-cached read is stale; the live fetch is the cache-bust).
3. Re-reads any other source γ named in the clarification (e.g. a spec file on `origin/main` per the spec-staleness rule §2.5 below).
4. Decides per §2.4 resumption protocol.

The reader does not need γ to chat-relay; the file write *is* the signal.

### 2.4. Resumption protocol

After reading a clarification, the in-flight role has four paths:

1. **Acknowledge and continue.** The change does not affect work in flight (e.g. γ clarified an AC the role had already interpreted correctly). The role notes the clarification in its own artifact (α in `self-coherence.md`; β in `beta-review.md`) — a one-line entry naming the clarification date and the impact assessment — and continues.
2. **Re-scope and continue.** The change requires a localized adjustment (e.g. a pinned axis re-pinned to a new value; a non-goal added). The role re-shapes work in flight, notes the re-scope in its own artifact, and continues from the next uncompleted step.
3. **Stop and escalate.** The change requires a rewrite (e.g. the language axis re-pinned from Python to Go after substantial implementation in Python). The role writes a stop note in its own artifact, signals via the artifact channel that the cycle needs operator attention, and waits. γ (or operator-as-γ) decides whether to (a) abandon work in flight and resume on the new contract, (b) close the cycle and re-scaffold a new one, or (c) accept the partial work and file a follow-on.
4. **Continue under protest.** Reserved for cases where the role believes the clarification itself is incoherent. The role logs the protest in its own artifact (α or β), continues on the *prior* contract, and surfaces the conflict to γ. γ adjudicates by either revising the clarification or routing through δ for boundary review (`cnos.cdd/skills/cdd/delta/SKILL.md §1 — outward membrane`). The protest path is rare; the resumption guidance field in the clarification (§2.1) should pre-empt it.

The role MUST NOT silently merge the clarification into work in flight without acknowledging it in its own artifact. Acknowledgment is the audit trail that lets β's Rule 7 verification distinguish "α applied the clarification" from "α happened to ship the new shape".

### 2.5. Issue-edit cache-bust

The most frequent rescue trigger is an **issue-body edit**. γ edits the cycle's GitHub issue body to:

- pin a missing implementation-contract axis (post-#391; the dispatch-time gate at `dispatch/SKILL.md §3` did not catch the omission),
- clarify an ambiguous AC,
- add or remove a non-goal,
- fold in a post-filing refinement (per `cnos.handoff/skills/handoff/cross-repo/SKILL.md §2.3.1` `modified` event — the same mechanism propagates a Delta to in-flight readers),
- record an operator override.

**The GitHub issue mtime is not a reliable signal.** MCP `gh issue view` calls may return cached state for an unspecified window; even authenticated `gh` CLI calls may hit caching layers. The reader cannot trust "did the issue update?" as a wake-up condition.

**The cycle-branch SHA transition is the canonical signal.** γ commits `gamma-clarification.md` to `cycle/{N}` *before* `gh issue edit` runs. The polling role sees the SHA advance, knows from the touched file that an issue edit is incoming, and re-fetches the issue body live in the next intake iteration. The clarification on the cycle branch is the explicit cache-bust signal that an issue-edit-only signal would lack.

This is also the cross-repo cache-bust mechanism — when γ records a STATUS event change or refines a `## Source Proposal` block on the cycle's target issue (per `cnos.handoff/skills/handoff/cross-repo/SKILL.md §2.3.4 — Emitters per event`), the same `gamma-clarification.md` write is the signal that propagates the refinement to any in-flight role.

### 2.6. Spec-staleness propagation

When γ (or δ) pushes spec changes to `origin/main` while α/β sessions are in-flight, loaded skills become stale. The wire-format invariant:

- **Identity-rotation mode (`cn dispatch` / `claude -p`):** Not applicable — each role invocation loads skills fresh from the filesystem; the next dispatch picks them up. No `gamma-clarification.md` needed; the next-iteration's fresh load is the cache-bust.
- **Long-lived polling sessions (legacy):** γ writes a coordination note to `.cdd/unreleased/{N}/gamma-coordination.md` on `cycle/{N}` (the gamma-coordination cousin of gamma-clarification — same channel, same authoring role, distinguished filename so the polling role can route by file name). The note names the spec change and the affected skill path; the in-flight role re-loads the affected skill from the new origin/main snapshot.

The list of *which* skill files trigger propagation in a given consumer protocol is **consumer-side** doctrine. For CDD-class cycles, the file list (CDD.md; role skill files; release/SKILL.md; review/SKILL.md when they land on `main` while a cycle is in-flight) lives in `cnos.cdd/skills/cdd/gamma/SKILL.md §2.5 → "Spec-staleness propagation"`. For CDS-class cycles, the canonical realization is at `cnos.cds/skills/cds/CDS.md §"Coordination surfaces"`. The wire-format invariant ("γ writes a coordination note when spec changes mid-flight, on the cycle branch, before the spec change becomes visible") is what this skill owns; the per-consumer trigger list is what cdd/cds own.

Derives from cnos#301 §9.1: γ proposed an out-of-spec merge because σ's `4a0f678` ("merge is β authority") landed mid-cycle outside γ's loaded snapshot. The reactive fix is the staleness check in γ's Load Order; this is the proactive side — γ writes the coordination note *before* the spec change makes the loaded snapshot stale.

---

## 3. Rules

### 3.1. Write the clarification before the other-surface edit

- ❌ Edit the issue body, then write `gamma-clarification.md`. The in-flight role may sample the issue body in between and apply the new contract with no audit trail.
- ✅ Commit and push `gamma-clarification.md` to `cycle/{N}` first; *then* run `gh issue edit` (or `git push` to main for spec change). The branch SHA leads the signal.

### 3.2. The cycle-branch SHA transition is the canonical signal

- ❌ Rely on the GitHub issue mtime to wake polling readers — MCP caches may make it invisible.
- ✅ Commit to `origin/cycle/{N}`; the SHA advance is the wake-up event the polling primitive emits.

### 3.3. γ is the sole author; α and β are the readers

- ❌ α writes to `gamma-clarification.md` to flag a gap it found; β writes to `gamma-clarification.md` to record a finding.
- ✅ α surfaces gaps via `self-coherence.md`; β surfaces gaps via `beta-review.md`; γ adjudicates and (if a clarification is the right tool) writes `gamma-clarification.md` itself.

### 3.4. Acknowledge the clarification in the reader's own artifact

- ❌ Silently merge a clarification into work in flight; β's Rule 7 verification cannot tell whether α applied the new contract or shipped the new shape by accident.
- ✅ α appends a one-line entry to `self-coherence.md` naming the clarification date and impact assessment; β appends a one-line entry to `beta-review.md`. The acknowledgment is the audit trail.

### 3.5. Do not redesign mid-cycle through clarifications

- ❌ Use `gamma-clarification.md` to add ACs, change the cycle's gap, or rewrite the contract substantively (> ~30 lines of clarification per event).
- ✅ Use `gamma-clarification.md` for bounded recovery (pin a missing axis; clarify an ambiguous AC; record an operator override). For substantive scope changes, close the cycle and re-scaffold a new one; the cycle boundary is the right unit of change for redesigns.

### 3.6. The rescue is a structural escape valve, not a routine cadence

- ❌ Treat repeated mid-flight clarifications as normal cadence; ship cycles with 3–4 clarification events as "the contract was iterative".
- ✅ A cycle with > 1–2 clarifications is a dispatch-time signal. ε reads the pattern (per `cnos.cds/skills/cds/CDS.md §"Cycle iteration triggers"`) and files a finding against the dispatch protocol (`dispatch/SKILL.md`) — under-specification at scaffold time is the root cause; codification, not more rescues, is the fix.

### 3.7. Cross-reference, do not duplicate

Consumer protocol skills cite this skill for the mid-flight rescue mechanism. If you find rescue-mechanism doctrine in `cnos.cdd/skills/cdd/gamma/`, `cnos.cdd/skills/cdd/{alpha,beta}/`, `cnos.cds/skills/cds/CDS.md`, or elsewhere that contradicts this skill, this skill governs and the other surface is patched to reference this one. The operational realization (the shell sequence) stays at cds as the v0.1 overlay; the wire-format invariant lives here.

---

## 4. Empirical anchors

The mid-flight rescue mechanism has empirical anchors across the cnos#388–#412 wave:

- **[cnos#391](https://github.com/usurobor/cnos/issues/391)** — the canonical anchor. Wrong-shape implementation (wrong package scoping + a separate-binary axis γ did not pin at dispatch). γ recovered by editing the issue body to pin the missing axes (`Package scoping`, `CLI integration target`, `Existing-binary disposition`) and committing `gamma-clarification.md` to the cycle branch. α picked up the cache-bust via cycle-branch polling and re-shaped the implementation. The cycle was rescued mid-flight without abandoning the branch. The empirical-anchor doctrine cited at `dispatch/SKILL.md §3` (δ-as-architect / cnos#393) names this rescue path as the structural escape valve when γ under-specifies and α improvises.
- **[cnos#392](https://github.com/usurobor/cnos/issues/392)** — recovery cycle. First cycle where δ pinned the 7-axis implementation contract at dispatch (the *upstream* fix); cnos#391's rescue path (the *downstream* fix) is what enabled cnos#392 to absorb the prior cycle's scope without abandonment. The two cycles compose: the dispatch-time enrichment at `dispatch/SKILL.md §3` prevents most mid-flight rescues; the mechanism here remains the escape valve when prevention fails.
- **[cnos#393](https://github.com/usurobor/cnos/issues/393)** — codification cycle. Made the dispatch-time implementation-contract first-class; named the mid-flight rescue as the structural escape valve in `dispatch/SKILL.md §3.3`. The mechanism's *first* doctrinal mention lives in dispatch; the mechanism's *canonical home* is this skill.
- **#406–#412 wave** — the post-codification cycles that exercised both the dispatch-time gate and the rescue channel in practice. Most cycles in the wave shipped with zero `gamma-clarification.md` events — the dispatch-time gate was sufficient. The rescue channel remains available; its rate of fire is the inverse signal of dispatch-time pinning quality.

---

## 5. Verify

### 5.1. Clarification authoring check (per rescue event)

- Does `.cdd/unreleased/{N}/gamma-clarification.md` exist on `origin/cycle/{N}`?
- Was it committed *before* the other-surface change (`gh issue edit`, spec push to main)?
- Does the entry name change / trigger / impact / rationale / resumption guidance per §2.1?
- Is γ the author (commit attribution)?

### 5.2. Reader-side acknowledgment check (at review time)

- Did the in-flight α append a one-line acknowledgment to `self-coherence.md` per §3.4?
- Did the in-flight β append a one-line acknowledgment to `beta-review.md` per §3.4?
- If the role chose Stop and escalate (§2.4 path 3), is the escalation visible in the cycle's record?

### 5.3. Cycle-cadence check (at close-out)

- How many `gamma-clarification.md` events fired in this cycle? If > 1–2, does the close-out (or `cdd-iteration.md`) name the dispatch-protocol gap that motivated them?

### 5.4. Spec-staleness check (when spec changes land mid-cycle)

- For spec changes pushed to `origin/main` during the cycle, did γ write `gamma-coordination.md` on `cycle/{N}` before the push?
- (Identity-rotation mode exempt — next-iteration's fresh load is the cache-bust.)

---

## 6. Related documents

- `cnos.handoff/skills/handoff/HANDOFF.md` — package contract; this skill is one sub-surface.
- `cnos.handoff/skills/handoff/dispatch/SKILL.md` — the dispatch wire format; precondition for this skill (this skill picks up after α is routed and dispatch's enforcement window closes). The dispatch skill cites this one for mid-cycle re-pinning at §3.3.
- `cnos.handoff/skills/handoff/artifact-channel/SKILL.md` — the substrate `gamma-clarification.md` lives on (per-role write ownership; sequential α→β→γ flow; frozen-snapshot rule on merge). Companion sub-skill landed alongside this one under Sub 4 of cnos#404.
- `cnos.handoff/skills/handoff/cross-repo/SKILL.md` — STATUS state machine; `modified` event post-filing refinement uses the same `gamma-clarification.md` mechanism to propagate a Delta to in-flight readers.
- `cnos.handoff/skills/handoff/receipt-stream/SKILL.md` — Sub 5 of cnos#404 (forthcoming). The cross-cycle aggregator; mid-flight rescue events feed into the cycle's `cdd-iteration.md` when the rescue indicates a dispatch-protocol gap.
- `cnos.cdd/skills/cdd/gamma/SKILL.md §2.5` — γ role-skill; γ's coordination loop names the issue-edit cache-bust procedure and the spec-staleness propagation (CDD-specific file list) as role-local procedure, citing this skill for the wire-format invariant.
- `cnos.cdd/skills/cdd/alpha/SKILL.md` — α intake / polling; cites this skill for the issue-edit cache-bust trigger.
- `cnos.cdd/skills/cdd/beta/SKILL.md` — β intake / polling; cites this skill for the issue-edit cache-bust trigger.
- `cnos.cds/skills/cds/CDS.md §"Coordination surfaces" → "Mid-flight clarification"` — v0.1 software-cycle operational realization (the `gh issue edit` + `git commit` + `git push` shell sequence γ runs at clarification time). Cites this skill for the wire-format invariant; this skill cites it for the v0.1 realization overlay.
- [cnos#391](https://github.com/usurobor/cnos/issues/391) — canonical empirical anchor (the rescue cycle).
- [cnos#392](https://github.com/usurobor/cnos/issues/392) — recovery cycle (the downstream cycle that absorbed #391's scope).
- [cnos#404](https://github.com/usurobor/cnos/issues/404) — parent tracker (handoff/coordination extraction wave).
- [cnos#418](https://github.com/usurobor/cnos/issues/418) — the codification cycle for this surface (Sub 4 of cnos#404).

---

## 7. Non-goals

- This skill does NOT own the dispatch-time implementation-contract enrichment. The 7-axis schema, the δ-as-inward-membrane review, the four-surface mesh (γ template / δ enrichment / α constraint / β verification) live in `cnos.handoff/skills/handoff/dispatch/SKILL.md`. This skill picks up where dispatch leaves off — once α is routed, dispatch's enforcement window closes; mid-flight is the post-dispatch channel.
- This skill does NOT own the artifact-channel rules. The per-role write ownership (α writes self-coherence + alpha-closeout; β writes beta-review + beta-closeout; γ writes gamma-scaffold + gamma-closeout + cdd-iteration + gamma-clarification + gamma-coordination), the sequential α→β→γ flow, and the frozen-snapshot rule on merge belong to `cnos.handoff/skills/handoff/artifact-channel/SKILL.md`. This skill cites the channel as the substrate.
- This skill does NOT own the operational realization. The `gh issue edit` + `git commit` + `git push` shell sequence γ runs at clarification time, the polling primitive that detects the SHA transition, the v0.1 software-cycle overlay all stay at `cnos.cds/skills/cds/CDS.md §"Coordination surfaces"` and `cnos.cdd/skills/cdd/harness/SKILL.md §5.4`. This skill names the wire-format invariant; cds/cdd carry the v0.1 realization.
- This skill does NOT redefine the cycle-state observation surfaces or polling primitives. Those live in `cnos.cds/skills/cds/CDS.md §"Coordination surfaces" → "Polling primitives"`. This skill assumes a polling primitive exists that emits on cycle-branch SHA transitions.
- This skill does NOT change the rescue mechanism's behavior. Sub 4 of cnos#404 is a wholesale extraction of an empirically-anchored mechanism (cnos#391); the file path, authoring role, reader role, trigger conditions, and cache-bust function are preserved verbatim from the empirical practice.
- This skill does NOT lift the clarification format into a typed schema. Whether to type `gamma-clarification.md` into `schemas/handoff/clarification.cue` is deferred; if a future cycle pressures the type-lift, the Markdown format references the schema by `$ref` then.

---

## 8. Kata

### Scenario

α is mid-implementation on `cycle/N`, polling `origin/cycle/N`. γ realizes (from an operator note) that implementation-contract row 7 (Backward-compat invariant) is wrong: the operator wants existing rules preserved + new content additive; γ had pinned "breaking change documented in §Migration". γ needs to re-pin without abandoning α's work.

### Expected reasoning

1. **Tool identification.** Single axis re-pin; α has not committed dependent implementation yet. Bounded recovery applies. Mid-flight rescue is the right tool (§2.4 path 2 — re-scope and continue). If the re-pin required a full language switch, γ would stop (§2.4 path 3) and re-scaffold instead.

2. **Clarification authoring.** γ writes (or appends to) `.cdd/unreleased/N/gamma-clarification.md` per the §2.1 format: change = row 7 re-pin; trigger = operator override; impact = AC4 (now requires additive-only diff); rationale = operator clarified the cycle is part of a migration wave; resumption guidance = α re-scope, revert in-flight breaking-change content, acknowledge in self-coherence.md before review-readiness.

3. **Operation sequence.** γ commits and pushes `gamma-clarification.md` to `cycle/N` *first*. Only then does γ effect any other surface change (e.g. `gh issue edit` if the issue body needs updating; in this case the implementation contract lives in the dispatch prompt on the cycle branch, so no issue edit is needed).

4. **α cache-bust.** α's polling primitive detects the SHA advance within the polling cadence. α's next intake iteration reads the clarification, re-reads the issue body live via `gh issue view N --json title,body,state,comments`, re-scopes per the resumption guidance, appends a one-line acknowledgment to `self-coherence.md`, and continues.

### Common failures

- γ runs `gh issue edit` first, then writes `gamma-clarification.md` — α may sample the issue body between the two operations and apply the new contract with no audit trail.
- γ chat-relays the change instead of writing the clarification — chat does not reach polling sessions.
- α reads the clarification but does not acknowledge in `self-coherence.md` — β's Rule 7 verification cannot distinguish "α applied the new contract" from "α happened to ship the new shape".
- γ writes a 100-line clarification rewriting the cycle's scope — at that volume the cycle should be closed and re-scaffolded.
- γ silently mutates the dispatch prompt on the cycle branch — same failure mode as silent issue-body edit.

### Verification

- `.cdd/unreleased/N/gamma-clarification.md` exists on `origin/cycle/N`; commit precedes any other-surface change.
- α's `self-coherence.md` carries a one-line acknowledgment naming the clarification date.
- β at review time confirms the diff matches the *post-clarification* contract; Rule 7 cites the clarification date in the verification record.

### Reflection

A second γ reading this kata's output should be able to reconstruct what changed, when, how it propagated to α, and whether α acknowledged before continuing — without consulting chat history. The artifact channel carries the complete record.

---

**End of skill.**
