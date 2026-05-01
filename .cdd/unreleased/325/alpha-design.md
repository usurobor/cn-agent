# Design — CDD Lifecycle Audit and Refactor (#325)

**Issue:** #325 — cdd: systematic lifecycle audit and refactor for role, artifact, and release coherence
**Author:** α
**Status:** final

---

## Problem

CDD exists to reduce incoherence across the system as a whole. When CDD itself is ambiguous or internally inconsistent, each cycle invents part of the lifecycle at runtime — the method becomes a source of the incoherence it is supposed to cure.

The specific gaps that motivated this cycle:

1. **Coordination model is not declared.** CDD.md contains both polling-session language and sequential-dispatch language without reconciling them. A reader cannot answer: is α expected to stay alive after signaling review-readiness, or exit?

2. **α close-out timing is broken.** CDD.md §1.4 α algorithm step 10 says "When β approves and merges: write alpha-closeout.md." In bounded sequential dispatch, α has already exited when β approves. No re-entry path is defined.

3. **No role/artifact ownership matrix exists.** The canonical filename table (§Tracking) and artifact location matrix (§5.3a) scatter ownership across two sections without a single mechanical lookup: "which artifact, owned by whom, verified by whom, blocks what if missing."

4. **γ closure gate is incomplete.** γ/SKILL.md §2.10 has 12 rows but does not explicitly make `gamma-closeout.md` itself a precondition for δ tag/release, and the failure consequence of each missing artifact is not stated.

5. **Small-change path is vague.** CDD.md §1.2 says "self-coherence is optional" but does not state which close-outs, release gates, and RELEASE.md requirements apply or collapse.

6. **No failure-mode regression surface.** The known failure patterns (stale unreleased dirs, missing alpha-closeout, missing RELEASE.md, δ tag before γ closure) exist only in issue prose — not in a reviewable checklist embedded in CDD.

---

## Decisions

### D1: Coordination model — sequential bounded role invocations

**Decision:** CDD's current coordination model is **sequential bounded role invocations dispatched one at a time via `claude -p` or equivalent.**

Each role is dispatched, completes its phase, and exits. In-session polling (CDD §Tracking) applies during each role's active session window — it is not a background daemon loop.

The polling instructions in §Tracking describe how each role detects transitions *while its session is running* (not across sessions). A role that has exited is not polling.

**Not current:** persistent polling role sessions where α, β, and γ remain alive indefinitely. That model requires daemon infrastructure (#310) not yet implemented.

**Compatibility:** When dispatch infrastructure provides long-lived session support, the polling instructions are already compatible — they describe in-session behavior and work either way. The declaration changes what "current" means, not the polling mechanics.

---

### D2: α close-out mechanism — re-dispatch after merge

**Decision:** In the sequential bounded dispatch model, α close-out is written through **γ-requested δ re-dispatch of α after β merge.**

Flow:
```
β merge → β writes beta-closeout.md → β exits
γ (in γ close-out session) reads beta-closeout.md
γ requests δ to re-dispatch α for close-out
δ dispatches α with close-out prompt
α reads beta-review.md (approved verdict) and the merged branch state
α writes alpha-closeout.md → commits to main → exits
γ continues closure
```

The re-dispatch prompt format is defined in CDD.md alongside the existing α dispatch prompt format. The re-dispatch uses a narrower `--allowedTools` set (Read, Write, Bash for git commit/push only).

**Acceptable alternative** (when re-dispatch is unavailable): α writes a provisional close-out at review-readiness time — explicitly marked provisional, noting findings known at that point. γ may supplement it with PRA observations. If used, this must be declared in the close-out and the self-coherence debt section.

---

### D3: Role/artifact ownership matrix — new §5.3b in CDD.md

**Decision:** Add a dedicated ownership matrix section to CDD.md §5.3b with the required columns:
`Artifact | Owner | Written when | Verified by | Required before | Missing means`

This supplements (not replaces) §Tracking canonical filename table and §5.3a Artifact Location Matrix, which own different concerns.

---

### D4: Lifecycle state table — new §4.1a in CDD.md

**Decision:** Add a lifecycle state/transition table after §4.1 Lifecycle steps, covering the major states from dispatch through closure. Format:
`State | Owner | Required inputs | Required outputs | Next state | Failure/retry`

---

### D5: Small-change path — explicit artifact collapse in §1.2

**Decision:** Expand CDD.md §1.2 to explicitly state which artifacts are required, optional, or not applicable for small-change cycles:
- alpha-closeout.md: optional (write if there are findings; otherwise one-liner "no findings")
- beta-closeout.md: not applicable (no β in single-author small change)
- gamma-closeout.md: not applicable unless γ coordinated
- RELEASE.md: required before tag/release (no exception)
- `.cdd/releases/{X.Y.Z}/{N}/` move: required at release time
- self-coherence.md: optional (may use commit message instead)
- PRA: required after every release including small-change releases

---

### D6: Failure-mode regression surface — new §8.1 in CDD.md

**Decision:** Add §8.1 "Closure verification checklist" to CDD.md §8 Gate that explicitly names the known failure modes as mechanical checks. This is not a full automated verifier — it is an inspectable checklist that a reviewer can run through per-cycle to catch the known failure modes.

---

### D7: Polling-era text — label as in-session

**Decision:** CDD.md §Tracking polling instructions are already correct for in-session use. Add a one-sentence clarification at the top of §Tracking: "Polling applies during each role's active session. A role that has completed its phase and exited is not polling." No mass rewrite needed.

---

## Scope of changes

Primary file: `CDD.md` — adds §1.6 (coordination model), §4.1a (lifecycle state table), §5.3b (ownership matrix), expands §1.2 (small-change), adds §8.1 (failure-mode checklist), fixes §1.4 α step 10 (close-out mechanism), adds §1.6a (re-dispatch prompt formats).

Role skill changes:
- `alpha/SKILL.md` §2.8: fix close-out timing; add re-dispatch path; add provisional fallback
- `gamma/SKILL.md` §2.10: add α close-out re-dispatch to closure gate; make gamma-closeout.md prerequisite for δ tag explicit; add missing-means column awareness
- `operator/SKILL.md` §1.2 or §3: add α close-out re-dispatch dispatch step and prompt
- `beta/SKILL.md`, `release/SKILL.md`, `post-release/SKILL.md`, `review/SKILL.md`: confirm agreement (minimal or no changes expected)

---

## Peer enumeration

Peer set for the coordination model change (role-skill peers):
- alpha/SKILL.md — owns close-out timing detail
- beta/SKILL.md — owns merge boundary; should agree
- gamma/SKILL.md — owns re-dispatch request; closure gate
- operator/SKILL.md — owns dispatch execution including re-dispatch
- review/SKILL.md — no change expected (review mechanics unchanged)
- release/SKILL.md — no change expected (RELEASE.md ownership already correct)
- post-release/SKILL.md — no change expected (γ-owned PRA timing unchanged)

Lifecycle-skill peers (separate enumeration class):
- review/SKILL.md — confirm no drift from CDD merge authority
- release/SKILL.md — confirm §2.5a movement timing matches new closure gate
- post-release/SKILL.md — confirm PRA timing relative to re-dispatch

All peers will be updated or explicitly exempted with a reason.

---

## What this does not change

- Runtime implementation (cn dispatch, #310, #295)
- Release script behavior
- Git hosting model
- Historical `.cdd/unreleased/` migration (covered by prior immediate-output in #324)
- β/γ scoring axes (those measure artifact integrity, surface agreement, cycle economics — unchanged)
- Role names
