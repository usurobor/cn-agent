# Agent activation on GitHub — three deployment paths to a running hub (master tracker)

Labels: `tracking`, `P2`, `core`

Priority: P2 — Enabling capability for the post-v3.82.0 field-application phase. Not blocking (the pause gates sub-dispatch), but it is the scoping envelope for "deploy a cnos agent to GitHub," which is the named next direction.

Status: The activation *procedure* exists (`agent/activate/SKILL.md`, cnos 3.78.0+). The deployment *paths* that get a fresh GitHub repo to a running, activated hub do not exist as canonical artifacts. The pieces are scattered across a skill, a README-router convention, and three design essays; no template, no published action, no deploy-mode `cn activate`.

This is a **master/tracker issue**. It captures the end goal and the sub-issue envelope. It does not itself implement. Per `cdd/issue/SKILL.md` master+subs pattern, subs are proposed here (not binding until filed) and dispatch under operator authorization (the v3.82.0 pause gates dispatch — see §"Pause posture").

---

## Governing question

> How does someone go from a GitHub account to a running cnos agent hub — via a template repo, via `cn activate` in a cloned repo, or via a deployable GitHub Action — without re-deriving the deployment per instance?

---

## Problem

What exists:
- **`cn activate`** (cnos 3.78.0+, `src/go/internal/activate/`) renders the bootstrap prompt from hub state. Assumes a local clone + the `cn` binary present.
- **`agent/activate/SKILL.md`** is the single source of truth for the activation procedure (six-step load order; three-tier capability matrix). README router template (§2.3) routes a body that has *already landed* on a hub.
- **Three design essays** name the architecture: `GitHub Actions as CN compute` (the operator realization at `cn-sigma:threads/adhoc/20260524-github-actions-as-cn-compute.md`), `BOX-AND-THE-RUNNER.md` (remote-runner delegation as an effect surface + the 6-field receipt), `FIDONET-AND-CNOS.md` (store-and-forward wake model).

What is expected:
- A new GitHub user can stand up a running agent hub through one of three documented, canonical paths, each routing into the same `agent/activate/SKILL.md` procedure.
- The deployment carries the effect-surface governance (`BOX-AND-THE-RUNNER` 6-field receipt) by construction, not by re-derivation.

Where they diverge:
- No template repo exists. A new operator hand-assembles `.cn/`, `spec/`, the README router, and a workflow.
- `cn activate` has no deploy-mode (it renders a prompt; it does not scaffold a hub or wire a wake workflow).
- No published GitHub Action or reusable workflow exists; each hub that wants Actions-driven wake re-derives `cn-wake.yml` — the "per-cycle doctrine re-invention" failure `cnos.handoff` names, applied to deployment.
- The wake loop's runtime dependency (`cn sync` / `cn agent`) is not yet in the Go `cn` (3.82.0 has `activate` + `deps` only; the full loop is the v4.0.0 Phase 5 agent-runtime port).

---

## Impact

- Affects: anyone who wants to run a cnos agent without bespoke setup — the entire "open coherent agency" adoption surface (`COHERENCE-MUST-BE-FREE.md`). Today, standing up a hub requires assembling scattered pieces and re-deriving the wake/deploy wiring.
- Enables: the field-application phase named by the v3.82.0 pause directive — real agents on real repos producing field evidence about intra/inter-agent communication and memory-return.
- Blocks (if unsolved): the social-network vision (`FIDONET-AND-CNOS.md`) — without a low-friction deploy path, the network can't grow past hand-built hubs.

---

## Status truth

| Surface | What ships | What is planned |
|---|---|---|
| `agent/activate/SKILL.md` | shipped 3.78.0+; six-step procedure + 3-tier inbound matrix + README router | the deployment paths route into it; possible §2.5 materialization addition (per the activate-foreign-body bundle, still drafted) |
| `cn activate` Go command | shipped; renders prompt from local hub | deploy-mode (scaffold + wire) does not exist |
| README router (§2.3) | shipped; adopted by cn-sigma | routes a landed body; not a deploy mechanism |
| `cn-hub-template` repo | does not exist | Sub A |
| Published GH Action / reusable workflow | does not exist | Sub C |
| `cn sync` / `cn agent` Go surface | not in 3.82.0 (OCaml-era ops; v4.0.0 Phase 5) | Sub D dependency |
| 6-field remote-runner receipt | doctrine shipped (`BOX-AND-THE-RUNNER.md` + `delta/SKILL.md §8`) | baked into the template workflow by Sub A/C |

---

## Source of truth

| Claim / surface | Canonical source | Status |
|---|---|---|
| Activation procedure | `cnos:src/packages/cnos.core/skills/agent/activate/SKILL.md` | Shipped 3.78.0+; the procedure all three paths route into |
| GitHub-Actions-as-compute realization | `cn-sigma:threads/adhoc/20260524-github-actions-as-cn-compute.md` | Operator realization; the compute-projection face |
| Remote-runner governance + 6-field receipt | `cnos:docs/gamma/essays/BOX-AND-THE-RUNNER.md` + `cnos:src/packages/cnos.cdd/skills/cdd/delta/SKILL.md §8` | DRAFT v0.1.0 essay; delta skill rule shipped |
| Store-and-forward wake model | `cnos:docs/gamma/essays/FIDONET-AND-CNOS.md` | v1.0.0 |
| Hub-as-cell + substrate | `cnos:docs/gamma/essays/CELL-OF-CELLS.md` §16.6, §18.5 | DRAFT v0.2.0 |
| Runtime gap (cn sync/agent pending) | `cn-sigma:threads/adhoc/v4-roadmap.md` Phase 5 | v4.0.0 pending |
| Materialization (where clones land when cwd ≠ hub) | `cn-sigma:.cdd/iterations/cross-repo/cnos/activate-foreign-body/ISSUE.md` AC1 | drafted bundle; gated by pause |

---

## Constraints

- **The activation skill is the single source of truth.** All three paths route into `agent/activate/SKILL.md`; none redefines the procedure (per its §2.3 adoption rule: hubs are routers, not redefinitions). A template's README is the §2.3 router verbatim; a deploy-mode `cn activate` renders the same six-step order; the GH Action runs the same procedure in the runner.
- **Every wake workflow is a remote-runner delegation** (`BOX-AND-THE-RUNNER`). Any artifact a path ships that causes a runner to execute (`cn-wake.yml`, a published action) MUST carry the 6-field receipt convention by construction. A template that ships an ungoverned workflow is distributing ungoverned effect surfaces — the failure the essay's "what not to celebrate" section names.
- **Store-and-forward, read-peers/write-self** (`FIDONET`). The default `GITHUB_TOKEN` (contents:write, own-repo-scoped) suffices for basic social operation; no peer write access. The deploy paths must not require cross-repo credentials for the base case.
- **Model credentials are the irreducible friction floor.** `cn agent` calls the model; the runner needs the operator's API key in repo secrets. No path can pre-fill this. The template can name the secret; the operator must provide it.
- **Runtime dependency.** Paths 1-2 (template, `cn activate`) partially work today (scaffold + prompt render). Path 3's full wake loop (`cn sync` → process → `cn agent`) is gated on the Go runtime port (v4.0.0 Phase 5). Sub D either ships that surface or documents a stub-mode wake (checkout → activate-render → push, no agent step) that proves the mechanics pre-v4.
- **Pause posture.** This tracker is authored under operator direction during the v3.82.0 pause. The pause permits *scoping* the field-application direction; it gates *dispatch* of the implementation subs. Each sub dispatches only on explicit operator authorization (the pause directive's "field evidence gates next theory expansion" applies; deployment infra is field-enabling, not theory-expanding, so the operator may authorize subs incrementally as field needs surface).

---

## Impact graph

What this wave touches (enumerated before any sub dispatches, per `design/SKILL.md §2.3`):

- `cnos:src/packages/cnos.core/skills/agent/activate/SKILL.md` — the procedure all paths route into; possible §2.5 materialization addition (consumes the activate-foreign-body bundle)
- `cnos:src/go/internal/activate/` — `cn activate` deploy-mode (Sub B)
- **new** `usurobor/cn-hub-template` repo — template (Sub A)
- **new** published action or `cnos:.github/workflows/cn-wake.yml` reusable workflow (Sub C)
- `cnos:src/go/` — `cn sync` / `cn agent` surface (Sub D; the v4 runtime port dependency)
- `cnos:docs/gamma/essays/BOX-AND-THE-RUNNER.md` 6-field receipt — embedded into the template workflow's accompanying receipt
- `cnos:install.sh` — the curl-install path the template/action may use to get `cn` onto the runner
- repo secrets convention (model API key) — documented by Sub A
- `cn-sigma` + `cn-pi` + future hubs — consumers; each adopts a path

Authority: `agent/activate/SKILL.md` governs the procedure; each path is a router/host, not a redefinition. If a path wants to change the load order or capability matrix, that change belongs in the skill, not the path.

---

## Proposed sub-issues (not binding until filed)

Per `cdd/issue/SKILL.md` master+subs and `cn-sigma:spec/OPERATOR.md` "Wave dispatch shape" (sub N filed when N-1 closes, unless independent + merge-safe):

### Sub A — `cn-hub-template` repository
"Use this template" → a new hub carrying: `.cn/` skeleton, `spec/PERSONA.md` + `spec/OPERATOR.md` stubs, the §2.3 README router verbatim, a `.github/workflows/cn-wake.yml` workflow, and the 6-field remote-runner receipt accompanying that workflow. One click + persona/operator fill + model-key secret = activatable hub.
Mode: design-and-build. Depends on: Sub C (the workflow it ships) OR ships a stub workflow first.

### Sub B — `cn activate` deploy-mode
Extend `cn activate` from prompt-render to optional scaffold/wire: `cn activate --init <hub-dir>` materializes the hub skeleton + router + workflow in a cloned/new repo. Consumes the activate-foreign-body §2.5 materialization (where clones land when cwd ≠ hub).
Mode: design-and-build. Depends on: activate-foreign-body bundle landing (currently drafted, pause-gated).

### Sub C — Deployable GitHub Action / reusable workflow
A `cnos:.github/workflows/cn-wake.yml` reusable workflow (`workflow_call`) hubs reference in three lines, OR a published `usurobor/cn-action@v1` composite action. Carries the 6-field receipt by construction. Triggers: `workflow_dispatch` + `schedule` (per FIDONET poll model); later `repository_dispatch` for peer-waking.
Mode: design-and-build. Depends on: Sub D (the `cn sync`/`cn agent` surface it invokes) OR ships stub-mode (wake mechanics only).

### Sub D — `cn sync` / `cn agent` Go runtime surface (or documented stub)
The wake loop's runtime dependency. Either port the OCaml-era `cn sync` (peer fetch / inbox materialize / outbox flush / projection update) + `cn agent` (oneshot scheduler / maintain_once / drain_queue) to Go (v4.0.0 Phase 5 work), OR document a stub-mode wake that proves checkout → activate-render → push without the agent step.
Mode: explore → design-and-build. This is the load-bearing dependency; Subs A and C ship full-fidelity only after this.

### Sub E — End-to-end field test
A real fresh GitHub repo activated all three ways; the mechanical success oracle from the operator realization: "a sleeping repo wakes, pulls peers, processes inbox, writes its own reply or receipt, pushes, and sleeps again." Records field evidence (the pause's gating deliverable).
Mode: explore. Depends on: at least one of Subs A/B/C reaching working state + Sub D's runtime (or stub).

---

## Acceptance criteria (wave-level)

This master closes when:

### AC1 — Three documented paths exist
Each of {template repo, `cn activate` deploy-mode, deployable GH Action} is a canonical, documented artifact that a new operator can use without re-deriving the deployment. Oracle: each path has a shipped artifact + a doc entry; a new operator following any one reaches an activatable hub.

### AC2 — All three route into the single activation procedure
No path redefines the activation procedure; each routes into `agent/activate/SKILL.md`. Oracle: diff each path's activation logic against the skill — zero procedure divergence (only the *mechanism* of reaching the skill differs).

### AC3 — Effect-surface governance is built in
Every workflow artifact a path ships carries the 6-field remote-runner receipt (`BOX-AND-THE-RUNNER`). Oracle: each shipped `*.yml` that causes runner execution has an accompanying receipt; no ungoverned effect surface ships.

### AC4 — Base case needs no peer write access
The default `GITHUB_TOKEN` (own-repo, contents:write) suffices for a hub to wake, process its own inbox, and push its own state. Oracle: a hub activated via any path runs a full wake cycle with no cross-repo credentials.

### AC5 — Field evidence recorded
At least one real hub is activated end-to-end and the wake-cycle success oracle is observed in git history. Oracle: commits from the Actions runner showing peer-fetch + inbox-process + self-publish.

---

## Non-goals

- Does not redefine the activation procedure (lives in `agent/activate/SKILL.md`).
- Does not build CCNF-O, IssueProposal.v1, RiskPolicy.v1, or TSC steering (pause-gated; separate roadmap per cnos#405).
- Does not require peer write access for the base case (read-peers/write-self).
- Does not solve event-driven peer-waking in v1 (schedule-poll is the FIDONET-aligned default; `repository_dispatch` is a later optimization that breaks the no-peer-write property).
- Does not pre-fill model credentials (irreducible operator-provided secret).
- Does not promote any of the five design essays from DRAFT to RATIFIED.
- Does not dispatch any sub during the v3.82.0 pause without explicit operator authorization.

---

## Pause posture

This tracker is authored under operator direction ("we need to do this next") during the v3.82.0 protocol-evolution pause. The tracker is a scoping artifact — it names the field-application envelope the pause directive itself pointed at (intra/inter-agent communication + deployment for field evidence). Sub-issues dispatch only on explicit operator authorization, incrementally, as field needs surface. The tracker does not lift the pause; it scopes the work the pause's "next phase" named.

---

## Leverage

**Positive:** every future hub becomes a one-click + one-secret deploy instead of a hand-assembled setup. The social-network vision (FIDONET) becomes reachable — agents can proliferate. The activation skill gains its deployment surface; the GitHub-Actions-as-compute realization gains its shipped artifacts.

**Negative:** adds a template repo, a published action/reusable workflow, and a deploy-mode to maintain — three new surfaces with their own versioning. The 6-field-receipt-by-construction requirement adds discipline overhead to every workflow the paths ship. Gated hard on the Go `cn sync`/`cn agent` runtime port (v4.0.0 Phase 5); full-fidelity deployment can't ship before that, only stub-mode.

---

## Success / closure condition

This master closes when AC1-AC5 pass, or when all live subs are closed and any deferred sub is named as tracked debt in this master's closure. The wave is the field-application enabler; its closure is the point at which "deploy a cnos agent to GitHub" is a documented, governed, one-of-three-paths operation rather than a hand-assembly.
