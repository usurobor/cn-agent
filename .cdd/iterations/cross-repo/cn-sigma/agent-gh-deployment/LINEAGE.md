# Cross-Repo Lineage: cn-sigma → cnos (agent-gh-deployment master tracker) — cnos-side mirror

Case (a) operator-directed cross-repo proposal. cn-sigma operator issued the directive; cn-sigma Sigma drafted the master/tracker issue using the CDD design + write + issue skills; a cnos-scoped body filed it. This file mirrors the source-side LINEAGE at `usurobor/cn-sigma:.cdd/iterations/cross-repo/cnos/agent-gh-deployment/LINEAGE.md`.

## Source

- **Repo:** `usurobor/cn-sigma`
- **Branch:** `main`
- **Authoring commit:** `6c95b53`
- **Path:** `.cdd/iterations/cross-repo/cnos/agent-gh-deployment/`
- **Source posture:** fourth cross-repo bundle from cn-sigma (after agent-activate-skill [shipped 3.78.0 → cnos#379], activate-foreign-body [drafted], release-3.82.0-hygiene [merged cycle/422]).

## Upstream — operator directive

Operator (usurobor / Axiom), 2026-05-24: "Create a master issue using cdd design and write skills to capture the end goal: an agent can be activated on GH via either templates repo or by running cn activate in the repo or if there's a way to deploy a GH action my way too."

The directive followed the five-essay read + the GitHub-Actions-as-compute realization. The master tracker is the scoping artifact for the field-application direction the v3.82.0 pause named.

## Target

- **Repo:** `usurobor/cnos`
- **Issue:** [cnos#431](https://github.com/usurobor/cnos/issues/431) — master/tracker
- **Kind:** tracking (umbrella; coordinates subs; does not implement)
- **Title:** "Agent activation on GitHub — three deployment paths to a running hub (master tracker)"
- **Labels:** `tracking`, `P2`, `core`
- **Disposition:** filed (accepted as tracker); open pending operator-gated sub-dispatch
- **Canonical path on cnos main:** `.cdd/iterations/cross-repo/cn-sigma/agent-gh-deployment/`

## Proposed sub-issues (envelope, not binding)

- Sub A — `cn-hub-template` repository
- Sub B — `cn activate` deploy-mode
- Sub C — deployable GH Action / reusable workflow (6-field receipt by construction)
- Sub D — `cn sync` / `cn agent` Go runtime surface (or documented stub) — the load-bearing dependency
- Sub E — end-to-end field test (mechanical success oracle from the operator realization)

## Pause interaction

Filed under the v3.82.0 protocol-evolution pause. The tracker scopes the field-application direction the pause named; it does NOT lift the pause. Sub-dispatch is operator-gated and incremental. Deployment infra is field-enabling (not theory-expanding), so the operator may authorize subs as field needs surface — but no sub dispatches without explicit authorization.

## Bilateral trace

1. **File:** cnos files the master tracker — done, cnos#431; ISSUE.md is the body. ✓
2. **Mirror:** cnos creates this `cn-sigma/agent-gh-deployment/` mirror (LINEAGE + STATUS + ISSUE.md audit copy). ✓
3. **Sub envelope:** the five proposed subs are filed incrementally per the wave dispatch shape, each on operator authorization. (pending)
4. **Close:** master closes when AC1-AC5 pass or all live subs close with deferred subs named as tracked debt. (pending)

## Notes

- This is a tracker — it does not implement. The subs carry the implementation; they dispatch under the pause's operator-gate.
- Sub D (the `cn sync`/`cn agent` runtime surface) is the load-bearing dependency. Subs A and C can ship stub-mode (wake mechanics, no agent step) pre-v4; full-fidelity deployment waits on the runtime port (v4.0.0 Phase 5).
- The 6-field remote-runner receipt (`BOX-AND-THE-RUNNER`) is a hard constraint on every workflow artifact the subs ship — built in by construction, not retrofitted.
- Fourth cn-sigma → cnos cross-repo proposal; same convention as the prior three.
