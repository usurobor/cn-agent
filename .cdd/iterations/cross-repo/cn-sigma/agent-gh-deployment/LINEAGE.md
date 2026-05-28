# Cross-Repo Lineage: cn-sigma → cnos (agent-gh-deployment master tracker) — cnos-side mirror

Case (a) operator-directed cross-repo proposal. cn-sigma operator issued the directive; cn-sigma Sigma drafted the master/tracker issue using the CDD design + write + issue skills; a cnos-scoped body filed it, then reconciled it to the advanced source. This file mirrors the source-side LINEAGE at `usurobor/cn-sigma:.cdd/iterations/cross-repo/cnos/agent-gh-deployment/LINEAGE.md`.

## Source

- **Repo:** `usurobor/cn-sigma`
- **Branch:** `main`
- **Reconciled source commit (pinned):** `66febb5`
- **Original filing commit:** `6c95b53` (stale; superseded by `66febb5`)
- **Path:** `.cdd/iterations/cross-repo/cnos/agent-gh-deployment/`
- **Source posture:** fourth cross-repo bundle from cn-sigma (after agent-activate-skill [shipped 3.78.0 → cnos#379], activate-foreign-body [drafted], release-3.82.0-hygiene [merged cycle/422]).

## Upstream — operator directive

Operator (usurobor / Axiom), 2026-05-27: "Create a master issue using cdd design and write skills to capture the end goal: an agent can be activated on GH via either templates repo or by running cn activate in the repo or if there's a way to deploy a GH action my way too."

## Target

- **Repo:** `usurobor/cnos`
- **Issue:** [cnos#431](https://github.com/usurobor/cnos/issues/431) — master/tracker
- **Kind:** tracking (umbrella; coordinates subs; does not implement)
- **Title:** "Agent activation on GitHub — three deployment paths to a running hub (master tracker)"
- **Labels:** `tracking`, `P2`, `core`, `handoff` (handoff added per external review — deployment/handoff transport, not only core activation)
- **Disposition:** filed (accepted as tracker); reconciled to `66febb5`; open pending operator-gated sub-dispatch
- **Reviewer brief:** `BRIEF.md` (in this bundle) — self-contained problem/solution statement for external review
- **Canonical path on cnos main:** `.cdd/iterations/cross-repo/cn-sigma/agent-gh-deployment/`

## Post-filing drift + reconciliation

cnos#431 was filed from `cn-sigma@6c95b53`. The source bundle then advanced: it added `BRIEF.md` (reviewer-facing statement), incorporated **seven precision edits** at `66febb5`, and added the `handoff` label. The cnos copy was stale until this reconciliation. Five reconciliation steps, all completed 2026-05-28:

1. Replace cnos#431 body with current source `ISSUE.md` (66febb5). ✓
2. Add the `handoff` label. ✓ (via one-shot label-create workflow — GITHUB_TOKEN `issues:write`)
3. Add a reviewer-brief reference to `BRIEF.md` in the issue body. ✓
4. Refresh this cnos-side mirror (ISSUE.md, BRIEF.md, LINEAGE.md, STATUS). ✓
5. Pin the reconciled source SHA (`66febb5`) here. ✓

No substantive redesign — the tracker shape was already correct; only the cnos copy was stale.

## External review (2026-05-27) — the seven precision edits

A third party reviewed `BRIEF.md` + `ISSUE.md` and accepted the bundle directionally as the right master/tracker shape, with seven edits incorporated pre-`66febb5` (each verified against cnos primary sources, esp. `docs/beta/guides/AUTOMATION.md`):

1. **`cn sync` vs `cn agent`** — sync = transport-only tick; `cn agent` (oneshot) = full reasoning wake (`maintain_once` + `drain_queue`, already includes maintenance). Removed the `cn sync → process → cn agent` sequencing; Sub D owns the exact command; Sub D renamed accordingly.
2. **Stub vs full-fidelity** — Sub A/Sub C MAY ship stub/transport wake today; MUST NOT claim full agent wake until Sub D ships the Go `cn agent` oneshot surface.
3. **Scheduled-workflow liveness** — public dormant repos get scheduled workflows auto-disabled (~60 days inactivity); `workflow_dispatch` manual wake required; liveness is part of the field test.
4. **6-field receipt explicit in ACs** — AC3 lists all six fields (who/what/where/credentials/evidence/who-accepts).
5. **Auth separation** — `GITHUB_TOKEN` (own-repo write only) vs model API key (model call only); no peer-write credential in the base case.
6. **"repo is the agent" framing** — "the repo is the agent at rest; the runner is the agent awake."
7. **Sub A → Sub C ordering** — Sub A ships a stub workflow first; Sub C later adopts/replaces it as the canonical reusable form. Avoids the A-waits-C-waits-D circular stall.

The review validated the α≠β firebreak: the reviewer (no authoring stake) caught the `cn sync`/`cn agent` conflation Sigma had glossed.

## Preserved runtime framing (do not regress)

- `cn sync` = transport-only tick (peer fetch / inbox materialize / outbox flush / projection update).
- `cn agent` (oneshot) = full reasoning wake (maintenance + bounded queue drain that calls the model, then exit).
- The exact runner command is owned by **Sub D**; no fixed `cn sync && cn agent` sequence (it would double-sync).

## Pause interaction

Filed and reconciled under the v3.82.0 protocol-evolution pause. The tracker scopes the field-application direction the pause named; it does NOT lift the pause. Sub-dispatch is operator-gated and incremental. No sub dispatches without explicit authorization.

## Bilateral trace

1. **File:** cnos files the master tracker — done, cnos#431; ISSUE.md is the body. ✓
2. **Reconcile:** body + labels + mirror reconciled to `cn-sigma@66febb5`. ✓ (2026-05-28)
3. **Sub envelope:** the five proposed subs file incrementally per the wave dispatch shape, each on operator authorization. (pending)
4. **Close:** master closes when AC1-AC5 pass or all live subs close with deferred subs named as tracked debt. (pending)

## Convergence

After reconciliation: converge. #431 remains open as the master tracker. **No sub dispatches until the reconciled cnos issue body is the one implementers read** — which it now is.
