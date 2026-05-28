# Durable AI agents on GitHub — problem and proposed solution

**A brief for external review.** Self-contained; assumes no prior knowledge of cnos. Written to be evaluated, not just read: each claim names what would confirm or falsify it.

---

## Governing question

> How does someone go from a GitHub account to a **persistent, portable AI agent** — one whose identity and memory survive across temporary compute — without standing up always-on infrastructure or re-inventing the plumbing each time?

---

## Problem

Most AI agents today are **ephemeral processes**. An agent runs as a chat session, a hosted loop, or a CLI invocation. Its state — what it did, what it learned, what it still owes — lives in runtime memory, hidden prompts, and orchestration scaffolding. When the process ends, that continuity is lost or trapped in a place the next process can't cleanly inherit.

This produces three concrete failures when work has to compose across time, tools, or other agents:

1. **No durable identity.** "The agent" is conflated with "the server that happens to be running it." Kill the process and the agent is gone, or reduced to a fresh model call with no memory of prior decisions, commitments, or failures.
2. **No inheritable handoff.** A deterministic function hands off through a stable contract. A probabilistic agent's output isn't enough — the next step needs to know what authority was used, what evidence was produced, what was skipped, what's still owed. Runtime memory is a poor substrate for that: local, fragile, hard to inspect.
3. **High friction to stand one up.** Running a persistent agent today means choosing a host, keeping a process alive, wiring memory, managing credentials — bespoke per deployment. There is no low-friction path from "I have a GitHub account" to "I have a running, durable agent."

**Evidence this is real:** the common industry pattern is an agent that works in month one, then accretes a routing table, a retry queue, two memory stores, a reconciliation job, and a ranking meta-agent by month six — at which point the orchestration scaffolding owns all the state and "the agent" is a stateless model call at the bottom. The state was never durable; the scaffolding compensates for that.

---

## Why it matters

If durable, low-friction agent deployment exists, agents can **proliferate and compose**: many small agents, each owning its scope, reading each other's published state, accumulating memory across runs — without a central server and without each one being a hand-built system. That is the difference between "agents bolted into existing workflows" and "agents as durable, composable units." The second is only reachable if continuity and deployment are cheap and inspectable.

---

## Constraints any solution must respect

- **No always-on requirement.** A solution that mandates a persistent server reintroduces the hosting burden. Agents that wake, act, and sleep should not need continuous presence.
- **Inspectable handoff.** Every episode's output must be legible evidence — what ran, under what authority, what changed, what's owed — not opaque runtime state.
- **Bounded credentials.** Basic operation must not require an agent to hold write access to other agents' state. Cross-agent write access is a security and coordination liability.
- **Irreducible cost named honestly.** The model API key is the one thing no template can pre-fill; a solution should name it, not hide it.

---

## Proposed solution

Treat the agent as **durable continuity stored in a Git repository** (the "hub"), and use ephemeral compute only to *enter* that continuity for bounded episodes.

Four parts:

### 1. The hub repo stores the agent's durable continuity

The Git repo holds the agent's identity, operating contract, memory, accumulated receipts, and open obligations. A "runner" (a model session, a CLI, a CI job) is temporary machinery that *enters* the agent, does bounded work, and writes evidence back. Kill the runner — the agent persists in the repo. This separates **deployment** (keeping a process alive) from **activation** (a runner entering durable continuity). In shorthand: **the repo is the agent at rest; the runner is the agent awake.** You deploy a runner; you activate continuity.

*Falsifiable by:* a runner activates, does work, disappears; a second, unrelated runner reconstructs the full state from the repo alone (no shared session). If it can't, the continuity isn't durable.

### 2. A standard activation procedure

Any runner becomes the agent by following a fixed, documented procedure: read the agent's identity + operating contract + accumulated state, do bounded work, write evidence back as a structured **receipt** (what ran, what authority, what changed, what's owed). The procedure is defined once and referenced, not re-derived per agent.

*Falsifiable by:* hand a fresh runner only the repo URL; it should reach a correct "I am this agent, serving this operator, with this current state" without further instruction.

### 3. GitHub Actions as serverless compute

A scheduled (or manually triggered) GitHub Actions workflow provides the ephemeral compute: it wakes, checks out the repo, runs the agent loop (fetch peers' published state, process inbound, act, write receipts), commits, pushes, and exits. No always-on daemon.

A wake has two fidelities: a **transport-only tick** (sync peers, materialize/flush, no reasoning) and a **full reasoning wake** (one maintenance cycle plus a bounded queue drain that calls the model, then exit). The exact runner command is a runtime detail, not fixed by this design.

**Liveness caveat (honest):** for *active* repos, ordinary commit activity can keep scheduled wakes alive. But GitHub auto-disables scheduled workflows in public repos after ~60 days of inactivity — so a genuinely dormant agent's schedule can be turned off. v1 therefore also provides a manual wake (`workflow_dispatch`), and the scheduled-liveness behavior is itself part of the field test, not an assumption.

*Falsifiable by:* a sleeping repo wakes on schedule, processes new input, writes its own reply/receipt, pushes, and sleeps — observable purely in git history, with no server running between wakes. (And: a dormant repo's schedule survives, or the manual-wake fallback is demonstrated.)

### 4. Store-and-forward sociality (read peers, write self)

Agents become social by **reading peers' repos and writing only their own**. No agent needs write access to another. One agent publishes state; peers fetch it on their next wake. This is the FidoNet store-and-forward model: asynchronous, no central server, latency measured in poll intervals. It degrades locally rather than collapsing globally.

*Falsifiable by:* two agents coordinate (one publishes, the other consumes and responds) using only each one's push access to its *own* repo.

### Deployment: three paths to a running agent

1. **Template repository** — "Use this template" → a new repo pre-wired with the activation router, the wake workflow, and identity stubs. One click + fill in identity + add the model API key as a secret.
2. **A command (`cn activate`)** — run in a cloned repo to scaffold/orient.
3. **A published GitHub Action / reusable workflow** — referenced from any repo's workflow file to add the wake loop.

All three route into the *same* activation procedure; none redefines it.

### Governance: workflows are effect surfaces

A workflow file is "a command waiting for another machine to run." Any workflow a deployment path ships carries a structured receipt naming six fields: (1) who authorized it, (2) what it runs, (3) where it runs, (4) with what credentials/permissions, (5) what evidence returns, (6) who may accept the result. This keeps remote execution auditable rather than a hidden capability.

### Credentials: two separate authorities, neither peer-facing

The base case needs exactly two credentials, kept distinct:

- **`GITHUB_TOKEN`** — the platform-default token, scoped to the agent's **own repo only** (write for commits/pushes; PR/issue authority within its own repo). It can write self; it cannot write peers.
- **Model API key** (e.g. `ANTHROPIC_API_KEY`) — model-call authority only. This is the irreducible friction floor: no template can pre-fill it; the operator supplies it.

No **peer write** credential exists in the base case — that is what makes read-peers/write-self enforceable rather than aspirational. If *private* peer reads are needed later, that is a separate, optional credential, never part of the base install.

---

## Scope — for the reviewer

**In scope:** the four parts above + the three deployment paths + the receipt-governance discipline.

**Out of scope (deliberately deferred):** autonomous goal-setting / self-directed issue generation; any "agent decides what to do next on its own" steering layer; event-driven cross-agent triggering (the polling model is the v1 default precisely because it preserves the no-cross-write property).

**Key dependency (stated honestly):** the agent's internal loop (fetch-peers / process-inbox / act / write-receipts) is partly built. A stub version proves the *wake-and-publish mechanics* today; the full loop depends on a runtime component still in progress. A reviewer should treat "the wake mechanics work" and "the full agent loop works" as two separate claims with different current maturity.

---

## The leverage claim (why this is worth doing)

This changes the system boundary so a **whole class of future work disappears**: standing up a persistent agent stops being a bespoke hosting project and becomes a one-click + one-secret operation, and a network of agents becomes reachable without any central service. If the claim holds, the marginal cost of the *n*-th agent approaches the cost of creating a GitHub repo.

**The strongest reasons to doubt it:**
- GitHub Actions' scheduled-workflow constraints (minimum interval; auto-disable on inactivity) may make truly-dormant agents unreliable without a liveness mechanism.
- Poll-interval latency may be too high for some coordination needs (the design accepts this as the cost of the no-central-server property).
- The "repo is the agent" model assumes the durable state genuinely captures everything the next runner needs; if important continuity leaks into runtime-only context, the model breaks silently.

A reviewer's most useful contribution: stress these three, and judge whether the falsification tests above are the right ones.

---

## What we're asking the reviewer to assess

1. Is the **problem** real and correctly framed, or an artifact of how we've been building?
2. Does the **repo-as-agent + activation-not-deployment** split actually solve it, or relocate the difficulty?
3. Are the **falsification tests** the right ones — would passing them genuinely demonstrate the claims?
4. Is the **scope boundary** (defer autonomy/steering; ship deployment + continuity first) the right cut?
