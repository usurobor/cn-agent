---
schema: cn.operator-review.v1
issue: 500
pr: 502
verdict: iterate
reviewer: human-operator
captured_by: kappa (HI)
captured_at: 2026-06-26 (UTC)
worker_pr_head_at_review: 235d2b99cd6e458931c14be4178e03f79407c2bd
findings_count: 5
ci_note: 1
---

# Operator-review — PR #502 (cnos#500 cycle/500 R1)

This artifact is **κ's translation of the human operator's iterate-narrowly verdict on PR #502** into typed durable input for the cycle's R[N+1] re-entry. κ authored this file; α/β/γ act on it.

## Verdict

`iterate` (narrowly). The review-return concept is right and the stale cnos#500 recovery empirically validates the need. Architecture converges (κ / `cn cell return` / `cn cell resume` / δ runtime / α/β/γ split). 5 mechanical blockers + 1 CI note before merge.

## Self-application note (load-bearing — affects R1 routing)

**PR #502 IS the implementation of `cn cell return` + `cn cell resume`.** The primitive the operator wants used to route this iterate back into cycle/500 is the primitive PR #502 builds. This is a bootstrap chicken-and-egg case.

Until PR #502 merges + the renderer re-emits + the live wake reads the new primitive in its prompt, the canonical mechanical flow (κ writes operator-review → `cn cell return` transitions label → wake resumes → δ §9.10 routes R[N+1] → α/β/γ R1) cannot complete itself. **κ flags this for operator decision before initiating any R[N+1] dispatch sequence** (see "Next steps" below).

## Findings

### F1 — `cn cell return` does not validate the artifact against the command flags

**Surface:** `src/go/internal/cell/return.go` (or analogous; the Go implementation of `Return()`).

**Problem:** Current `Return()` checks only `schema: cn.operator-review.v1`. It trusts the CLI flags for `--issue` and `--verdict` and does not parse the artifact frontmatter to confirm the artifact actually matches.

Example failure mode: an HI runs `cn cell return --issue 500 --verdict iterate --review .cdd/unreleased/497/operator-review.md` and the command accepts the artifact even though it says `issue: 497` and `verdict: converge`. This breaks the point of typed operator-review.

**Expected change:**

```text
cn cell return must parse at least:
  issue
  verdict
from the frontmatter and verify:
  artifact.issue == --issue
  artifact.verdict == --verdict
If not, fail before any label mutation.
```

The artifact is the authority. CLI flags select and confirm it. They must not override silently.

**Class:** load-bearing mechanical invariant (artifact-as-authority); the typed-input doctrine fails if `cn cell return` accepts mismatched artifacts.

### F2 — `cn cell return` does not preflight the issue state

**Surface:** `src/go/internal/cell/return.go` (same file as F1).

**Problem:** The command applies `status:review → status:changes` directly without first verifying the issue's current label state. If the issue is in any state other than `status:review`, this would corrupt the lifecycle (e.g., applying `status:changes` to a cell at `status:in-progress` or `status:blocked`).

**Expected change:**

```text
Before applying labels, inspect the issue labels and require:
  exactly one status:* label
  status:review present
  status:changes absent
  issue open
If the state is wrong, fail with a clear message such as:
  review_return_state_invalid
This should be tested.
```

**Class:** mechanical lifecycle invariant; this is what makes `cn cell return` a primitive (not a wrapper around naive `gh issue edit`).

### F3 — Label transition is unsafe if the second command fails

**Surface:** `src/go/internal/cell/return.go` (label transition implementation).

**Problem:** Current implementation removes `status:review` first, then adds `status:changes`. If remove succeeds and add fails, the issue is left with NO lifecycle status — worse than staying at `status:review`.

**Expected change:**

```text
Use one of these patterns:
  Preferred:
    one gh issue edit call that removes status:review and adds status:changes together,
    if gh supports both flags in one call.
  Acceptable:
    add status:changes first,
    then remove status:review,
    and if cleanup fails, report explicit label drift:
      review_return_label_drift
The command must not silently leave the issue statusless.
```

**Class:** mechanical atomicity invariant; partial failure must leave the issue in a recoverable state.

### F4 — `cn cell resume` does not guarantee it is editing the cycle branch

**Surface:** `src/go/internal/cell/resume.go` (or analogous; the Go implementation of `Resume()`).

**Problem:** Current `Resume()` verifies `origin/cycle/{N}` exists, then edits `.cdd/unreleased/{N}/self-coherence.md` in the current checkout. This can edit the wrong branch if the caller is on `main` or some other branch. Does not commit or push the R[N+1] marker. The help text says it "re-arms an existing cycle" and "preserves branch + artifacts," but the code does not ensure the local working tree is actually `cycle/{N}`.

**Expected change:** choose one of two honest designs.

**Option A — mechanical resume command (preferred):**

```text
cn cell resume should:
  fetch origin/cycle/{N}
  checkout or create local cycle/{N}
  verify artifact dir
  append §R[N+1]
  commit
  push
```

**Option B — local-only helper (acceptable for v0 if smaller scope wanted):**

```text
cn cell resume prepares local artifact state only.
Caller must already be on cycle/{N}.
Caller must commit/push after running.
```

Then add a hard preflight: `current branch == cycle/{N}`. If not, refuse.

**Operator's stated preference:** Option A eventually; Option B acceptable for this PR if stated clearly and tested.

**Class:** mechanical correctness invariant; the command must not claim it re-arms the cycle while editing an unspecified checkout.

### F5 — `gamma-interface (HI)` should not become canonical

**Surface:** `src/packages/cnos.cdd/skills/cdd/operator-review/SKILL.md` (operator-review schema's `captured_by` examples).

**Problem:** The current schema lists `gamma-interface (HI)` as a valid canonical example of `captured_by`. This is historically true (cycle/497's `operator-review.md` used `gamma-interface (HI)` because κ had not yet been named). But it is not the name we want to canonize. The whole point of cnos#501 is to prevent the human interface from being taxonomically folded into CDD cell roles. Canonical examples should reflect the corrected naming.

**Expected change:**

```text
Use canonical examples like:
  kappa (HI)
  sigma (HI)
  human-operator-direct
If you want to preserve the history, add:
  legacy/historical witness:
    gamma-interface (HI)
but do not make it the first-class example.
```

**Class:** doctrinal vocabulary consistency; sibling of cycle/497 O4 (CDS closure vs boundary acceptance) and cycle/497 O5 (actor-collapse declaration). Failing to fix this canonizes the very confusion cnos#501 corrects.

## CI note (non-blocking; needs clean statement)

PR #502's checks page shows SKILL frontmatter validation failing. The PR body says these are inherited main-state failures (true; pre-existing infra failures cited at PR creation). But because this PR adds `src/packages/cnos.cdd/skills/cdd/operator-review/SKILL.md`, the final review should either rerun the validator after the F-fixes OR explicitly inventory that the new skill is not contributing to I5.

**Expected change:** clean statement in PR body or β-review §R1 that `operator-review/SKILL.md` passes frontmatter validation independently of the inherited failures.

Not blocking; operator does not block merge on I4/I5/I6 inherited caps; the statement is for cleanliness.

## Recovery-path context (κ provenance)

This artifact is filed BEFORE any role-pass-on-top dispatch. κ has NOT:
- transitioned labels (`status:review → status:changes` would normally be `cn cell return`'s job; the primitive is being iterated, so manual flip would be doing-the-thing-while-installing-it)
- spawned α/β/γ R1 sub-sessions

κ has:
- posted operator's verbatim verdict as PR review comment ([#4813597929](https://github.com/usurobor/cnos/pull/502#issuecomment-4813597929))
- filed this `operator-review.md` as typed durable input artifact

## Next steps (κ flags for operator decision)

Because PR #502 IS the implementation of `cn cell return` / `cn cell resume`, the standard route-back path is bootstrapping itself. Two options:

**Option (a) — manual role-pass-on-top sequence** (cycle/497 precedent; declared `degraded_recovery: human_interface_used_substitute_dispatch_during_self_install`):

1. κ flips label `status:review → status:changes` manually (operator-authorized lifecycle action; documents the manual transition as degraded recovery for the duration of self-install)
2. κ spawns α R1 Agent sub-session to inspect F1–F5 + CI note + take ownership of code changes
3. κ spawns β R1 Agent sub-session to independently review
4. κ spawns γ R1 Agent sub-session to update closeouts + record `degraded_recovery: self_install_bootstrap`
5. cycle/500 reaches `status:review` again; operator-final-read on the corrected PR

**Option (b) — operator-driven sequence** (operator picks up the iterate directly):

1. Operator transitions label `status:review → status:changes` via UI
2. Operator pushes fixes directly OR spawns α R1 themselves
3. Cycle/500 reaches `status:review`; operator-final-read

**Option (c) — wait for the runtime side**: hold cycle/500 until cnos#504's runtime side ships separately. Not practical — cnos#504 depends on cnos#503 which depends on... still requires PR #502 to merge.

**κ's recommendation, presented for operator verdict:** Option (a). The bootstrap-exception path is well-established now (cycle/497 cycle, gamma-closeout §5 `degraded_recovery` declaration). Sub-sessions for narrow mechanical correctness fixes (F1–F5 are all small, scoped code changes) is the cycle/497 pattern. This time the `degraded_recovery` reason changes: `self_install_bootstrap` instead of `missing_review_return_primitive`. Once PR #502 merges, future iterates use `cn cell return` and the bootstrap exception goes away.

**κ does NOT proceed without explicit operator authorization on this choice.** Per κ doctrine: κ may translate intent + record verdict + apply operator-authorized actions; κ must not unilaterally initiate role-pass-on-top sequences. The cycle/497 precedent was operator-authorized after the boundary violation; here κ is asking first.

— κ@cnos.core (herald; human interface; cycle/500 R1)
