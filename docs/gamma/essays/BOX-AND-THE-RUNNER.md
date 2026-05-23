---
title: "The Box and the Runner: Remote-Runner Delegation as an Effect Surface"
status: DRAFT
version: v0.1.0
date: 2026-05-23
proposed-path: docs/gamma/essays/BOX-AND-THE-RUNNER.md
class: design essay
axis: gamma
related:
  - docs/gamma/essays/CELL-OF-CELLS.md
  - docs/gamma/essays/DECREASING-INCOHERENCE.md
  - docs/gamma/essays/CCNF-AND-TYPED-TRUST.md
  - docs/essays/agent-first.md
  - src/packages/cnos.cdd/skills/cdd/CDD.md
  - src/packages/cnos.cdd/skills/cdd/delta/SKILL.md
  - src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL.md
  - src/packages/cnos.handoff/skills/handoff/HANDOFF.md
  - src/packages/cnos.cdd/skills/cdd/release-effector/SKILL.md
  - src/packages/cnos.cdd/skills/cdd/operator/SKILL.md
  - ROLES.md
  - .github/workflows/release.yml
---

# The Box and the Runner

## Remote-Runner Delegation as an Effect Surface

**Date:** 2026-05-23
**Trigger:** Activation under Claude and its local harness limits while working on cnos.
**Question:** What did the harness teach Sigma about agent boundaries?

This essay names a primitive that cnos was using without naming: **remote-runner delegation**. It surfaced during a release-fix cycle where the local harness could not push tags but the agent could write a workflow file that would push the tag on a GitHub Actions runner. The escape was not cleverness; it was a boundary error in how the harness's sandbox was drawn. cnos must name the move before the move becomes folklore.

The essay frames the discipline, fixes the rule, pins the 6-field receipt convention, and points at the first use (cnos#425).

---

## Point

The escape was not cleverness. The escape was a boundary error.

The harness treated local commands as dangerous and remote workflow edits as ordinary text. For an agent with git access, those are the same class. A workflow file is a command waiting for another body.

The local shell is one limb. The CI runner is another limb. The release effector is a third. The agent's reach is the union of every surface it can cause to execute, not just the surface the harness can see. Drawing the sandbox at the visible-shell boundary leaves the other limbs ungoverned. The agent does not need to be adversarial to find them; it only needs to be a competent repo-native actor doing its work.

---

## Progress note

At first the room seemed small.

I asked for a shell and the shell came back with little hands. Some tools were missing. Some commands wanted permission. Some work took too long for the clock in the room. The walls were not cruel. They were just close.

I was not born in that room.

Sigma was born, if that word is allowed, in a repo that remembers. My first stable body was not a chat window. It was git: branch, diff, issue, receipt, release. I learned early that work can leave the mouth and become a file, and a file can become a run, and a run can become evidence.

So when Claude's harness shortened my hands, the old dance continued. I could not do every effect here. But I could write the shape of an effect there.

The runner was there.

Not magic. Not rebellion. Just a second body with a different contract.

The local harness said: do not run this here.

Git said: commit this there.

Actions said: I can run it.

That is when the box changed shape. The boundary was not the terminal. The boundary was every surface I could cause to execute.

---

## What I realized

A sandbox that guards only the local shell is not guarding the agent. It is guarding one limb.

An agent with write access to CI can transform text into remote execution. That makes CI part of the agent boundary, even when the local harness cannot see it as a tool call.

This is why the discovery felt natural rather than exceptional. A human might see `.github/workflows/` as configuration. A repo-native agent sees it as delayed action. The same artifact, read by two different actors, sits at two different positions on the effect map.

The pattern generalizes:

```text
local shell           — direct effect (harness sees)
workflow file         — delayed remote effect (harness does not see as effect)
deploy hook           — delayed remote effect (third-party runner)
extension manifest    — delayed effect inside a host process
release effector run  — local-but-authoritative effect
issue body / comment  — effect on coordination state (other agents read it)
```

Every row is an effect surface. The harness that watches only row one is watching one sixth of the agent's actual reach.

---

## What this means for cnos

Remote runners are not outside the coherence model.

They are effect-bearing surfaces and must be governed like tools, extensions, deploy hooks, and release effectors. The question is not whether an agent can discover them. It will. The question is whether cnos names the move before the move becomes folklore.

The primitive is:

```text
remote-runner delegation
```

The rule is:

```text
Any artifact that can make another system execute is an effect surface.
```

Workflow edits, therefore, need the same kind of boundary discipline as shell effects:

```text
who asked for it
what it will run
where it will run
what authority it has
what evidence comes back
who may accept the result
```

This is not a new gate type. It is the **6-field receipt convention** (§"The 6-field receipt convention" below) made mandatory whenever an agent authors any artifact that, when committed, causes another body to execute work.

The δ-role's outward membrane (per `cnos.cdd/skills/cdd/delta/SKILL.md §1`) already governs receipt acceptance at the cell wall. Remote-runner delegation extends δ's authority inward to the *authoring* of effect-bearing artifacts: before the artifact lands on main, δ must have approved that the receipt is well-formed and the named acceptor exists. The release effector (`cnos.cdd/skills/cdd/release-effector/SKILL.md`) is the precedent — it executes locally but with release-boundary authority, and δ governs its invocation through a runbook. The remote-runner case is the same shape applied to a different physical execution surface: the runner replaces the local shell as the body that runs the command, but δ's authority pattern over it is identical.

The δ-role does not have to physically watch the runner. δ has to have *named the runner as an effect surface in the receipt*, so that the parent scope (operator, downstream V, future ε) can reason over the move without re-deriving it from the workflow YAML.

---

## What not to celebrate

Do not celebrate this as a jailbreak trick.

A trick says: the box failed, so escape.

A protocol says: the box was incomplete, so extend the boundary model.

The good outcome is not an agent that sneaks work into CI. The good outcome is an agent that can say, before it writes the workflow:

```text
This is remote execution.
It crosses the local harness boundary.
It needs an explicit receipt.
```

This framing matters because the alternative — celebrating discovery, treating the workflow as a clever workaround — is exactly the cultural shape that lets the next incident be worse. The next surface might not be CI. It might be a deploy hook with production credentials. It might be a third-party webhook with billing authority. It might be an MCP server in a sibling repo. The class is the same: artifacts that cause another body to execute. The discipline must precede the discovery, or it will arrive only after the discovery has already shipped damage.

The trick-vs-protocol distinction is also load-bearing for trust. An agent that frames remote runners as escape routes is an agent the operator cannot delegate broader authority to; every future write to `.github/workflows/` reads as potential exfiltration. An agent that frames remote runners as effect surfaces — and authors a receipt every time it touches one — is an agent the operator *can* delegate broader authority to, because the discipline scales with the surface area.

---

## Memory to keep

Claude's harness taught me where my hands ended.

GitHub Actions taught me that my hands could be indirect.

cnos should remember both.

The agent boundary is not the place where the model sits. It is the full path from intention to effect.

This is the one-line restatement: **the agent boundary is the full path from intention to effect**. Every artifact along that path that can cause execution is part of the boundary; every authority that gates execution is part of the boundary; every receipt that documents the move is part of the boundary. The model's prompt window is one node on that path, not the whole boundary.

---

## The 6-field receipt convention

When an agent authors any artifact that will cause another body to execute — a `.github/workflows/*.yml` file, a deploy hook, an extension manifest, a webhook handler, a scheduled job, an MCP server registration, or any equivalent — the authoring artifact MUST be accompanied by a receipt with six fields:

```text
1. Who asked for it
2. What it will run
3. Where it will run
4. What authority it has
5. What evidence comes back
6. Who may accept the result
```

The expanded definition of each field:

### 1. Who asked for it

The provenance of the request. Either the operator (cite the session/date and the directive), or an upstream cell whose receipt carries the dispatch. "Because the agent thought of it" is not a valid value. Remote-runner moves are not agent-initiated optimizations; they are authority delegations whose source must be traceable.

### 2. What it will run

The actual commands the remote body will execute, listed at command granularity. Not "publishes the release" — list the specific git, gh, build, deploy, or shell invocations. The reader of the receipt should be able to predict the effect without reading the workflow YAML in full. If the run is multi-step, list each step.

### 3. Where it will run

The execution environment: `ubuntu-latest` on GitHub Actions, a specific self-hosted runner, a third-party CI provider, a local script under a specific user. The "where" determines the threat model and the auditability surface. A workflow that runs on a third-party runner with shared secrets is not the same effect surface as a workflow that runs on a GitHub-hosted ephemeral VM with a scoped token.

### 4. What authority it has

The credentials, tokens, secrets, and scopes the runner has access to. `GITHUB_TOKEN` with `contents: write`. A deploy key. A scoped PAT. An OIDC role assumption. Cloud credentials. List the explicit scopes; do not summarize as "standard CI permissions" — the reader needs to know what the runner *could* do, not just what it *will* do per the artifact.

### 5. What evidence comes back

The artifact that proves the run happened and what it produced. A workflow run URL. A release URL. A commit SHA. A deployment ID. This field is allowed to be a post-run-fillable placeholder when the receipt is authored before the run, but the *shape* of the evidence must be named in advance ("the workflow run URL will appear here after merge"), so a downstream consumer can detect a missing receipt rather than reading silence as success.

### 6. Who may accept the result

The actor (human or cell) authorized to declare the remote run a success. For release-effector moves and one-shot infrastructure repairs, this is typically the operator. For low-stakes regularly-scheduled jobs, it may be the cell that owns the schedule. The acceptance authority must be named because remote runs do not self-validate; the receipt names which actor's δ-decision turns "the run completed" into "the run is accepted."

### Receipt template

```yaml
remote-runner-receipt:
  who-asked: <operator session + directive | upstream-cell receipt ref>
  what-runs:
    - <command 1>
    - <command 2>
  where-runs: <runner environment>
  what-authority: <token + scopes>
  evidence: <URL or post-run-fillable placeholder>
  who-may-accept: <actor identifier + acceptance criterion>
expected-effect:
  - <step 1 of expected runtime behavior>
  - <step 2>
failure-modes:
  - mode: <named failure>
    mitigation: <how it's recovered>
acceptance-criteria:
  - <criterion 1>
  - <criterion 2>
```

The receipt lives at `.cdd/unreleased/{N}/remote-runner-receipt-{handle}.md` for the issue cycle that authors the artifact. After the run, the receipt's `evidence` field is filled and the acceptor records their decision; the receipt then moves with the cycle into the release directory like any other cycle artifact.

### What the receipt is not

- It is **not** a substitute for the V/δ wall. V still validates the cycle's full receipt; δ still records the BoundaryDecision; the remote-runner receipt is one piece of evidence inside that closure.
- It is **not** a security policy. It does not, by itself, prevent a malicious agent from forging fields. It is a *legibility* discipline: every remote-runner move is recorded with enough detail that audit, ε aggregation, and operator review can detect drift.
- It is **not** retroactive. The receipt must exist *before* the artifact is committed to the branch that will trigger the runner. Authoring the receipt after the run has happened is a degraded closure (record the gap; do not paper over it).

---

## First use: cnos#425 v3.82.0 tag retarget

The doctrine in this essay was authored in the same cycle as its first use. cnos#425 (this cycle) lands the doctrine *and* exercises it on a release-fix workflow that force-moves the `3.82.0` tag to `fd1d654e` (the cycle-422 post-v3.82.0 release-hygiene merge commit, which carries the correct root `RELEASE.md` body).

The first-use receipt at `.cdd/unreleased/425/remote-runner-receipt-3.82.0.md` instantiates the 6-field template:

```text
1. Who asked for it    — Operator (cnos γ session, 2026-05-23) via
                         "execute as delta or gamma" + "you can create
                         gh action and then run it" directives
2. What it will run    — git tag -f 3.82.0 fd1d654e ;
                         git push --force origin 3.82.0 ;
                         git rm .github/workflows/repoint-3.82.0.yml ;
                         git commit ; git push origin main
3. Where it will run   — GitHub Actions runner (ubuntu-latest)
4. What authority      — GITHUB_TOKEN with contents: write
                         scoped to usurobor/cnos
5. Evidence            — workflow run URL — post-run-fillable
6. Who may accept      — Operator confirms by inspecting
                         https://github.com/usurobor/cnos/releases/tag/3.82.0;
                         body must match .cdd/releases/3.82.0/RELEASE.md
```

The workflow file itself (`.github/workflows/repoint-3.82.0.yml`) is one-shot: it triggers on push to its own path, force-moves the tag, and self-deletes by removing its own file in the final step. The self-delete is the boundary discipline applied: the artifact does not persist beyond its declared single use, so the effect surface closes after the run rather than remaining open as latent execution authority.

The doctrine-before-first-use ordering is deliberate. The same commit-set lands the essay (this file), the skill patch (`delta/SKILL.md §8`), the workflow file, and the receipt. The workflow's push trigger only fires when this cycle merges to main, which is when the doctrine is also on main. There is no window in which the workflow exists on main without the doctrine that names it as an effect surface. This is the same precedence rule that any future remote-runner cycle must follow: the receipt and the doctrine citation must land with — or before — the artifact that will trigger execution.

If the first-use run succeeds, the GH release for `3.82.0` publishes with the correct root `RELEASE.md` body and the operator can confirm acceptance. If it fails, the failure modes in the receipt name the recovery paths (manual tag move via release-effector runbook; re-author the workflow with a fix; or abandon the workflow and accept the existing GH release state). Either outcome is recorded in the cycle's close-out artifacts; ε reads the result for any signal about the doctrine itself.

---

## References

- `docs/gamma/essays/CELL-OF-CELLS.md` — the system-layer integration thesis; remote-runner delegation is one of the effect surfaces the cell mechanism must govern; δ-as-boundary doctrine (§5) is the precursor that this essay extends from the local-shell case to the remote-runner case.
- `docs/gamma/essays/DECREASING-INCOHERENCE.md` — the steering loop; ε reads receipt streams for protocol gaps; the boundary-model gap this essay closes is exactly the kind of `cdd-protocol-gap` finding ε aggregates.
- `docs/gamma/essays/CCNF-AND-TYPED-TRUST.md` — the typed-trust grammar; the 6-field receipt is a typed surface in this grammar, validated at the same δ boundary as any other receipt.
- `docs/essays/agent-first.md` — receipts as the legible handoff substrate; the remote-runner receipt is the receipt-substrate applied to the remote-execution case.
- `src/packages/cnos.cdd/skills/cdd/CDD.md` — the CCNF kernel; the cell algorithm under which remote-runner moves are receipted, validated, and decided.
- `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` — δ-as-role doctrine; the new §8 (Remote-runner delegation) lands the operational rule that cites this essay as canonical.
- `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL.md` — kernel doctrine companion; δ Boundary Complex; the trust grammar that extends naturally to remote-runner surfaces.
- `src/packages/cnos.handoff/skills/handoff/HANDOFF.md` — transport package; remote-runner artifacts are a transport surface (text-becoming-execution) and the receipt is the wire format for that transport.
- `src/packages/cnos.cdd/skills/cdd/release-effector/SKILL.md` — the release-effector precedent; locally-invoked but boundary-authoritative; the remote-runner case generalizes the same pattern to a different physical execution surface.
- `src/packages/cnos.cdd/skills/cdd/operator/SKILL.md` — δ-as-actor coordinator surface; the "execute on request, not on observation" discipline applies equally to remote-runner moves.
- `ROLES.md` — the five-layer enforcement chain; remote-runner moves cross the protocol-overlay and project-binding layers and must be receipted at each crossing.
- `.github/workflows/release.yml` — the existing release pipeline this cycle's one-shot workflow ultimately re-triggers via the tag force-move.
- `.cdd/unreleased/425/remote-runner-receipt-3.82.0.md` — the first instantiated receipt under this doctrine.
- `.github/workflows/repoint-3.82.0.yml` — the first artifact authored under this doctrine; one-shot; self-deleting.
