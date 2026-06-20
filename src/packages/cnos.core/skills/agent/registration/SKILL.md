---
name: registration
description: Register a fresh private-repo body with its home agent hub. Defines the pending_home_registration state machine, registration proposal schema, home credential model, operator-visible registration packet, and degraded reporting for unregistered or uncredentialed activations. Source of truth for the registration gate between a fresh foreign body and its home hub.
artifact_class: skill
kata_surface: embedded
governing_question: How does a fresh foreign body (especially private-repo) become a registered, readable activation at its home hub — with operator-only gates named explicitly and no gate hidden under a one-verb UX?
visibility: public
parent: agent
triggers:
  - registration
  - register
  - private repo
  - pending_home_registration
  - cohere-join private
scope: task-local
inputs:
  - the home hub URL (cn-sigma or any agent home)
  - the foreign repo URL (the new body being registered)
  - repo visibility (public or private)
  - activated identity (see agent/activate/SKILL.md)
outputs:
  - state machine: unregistered → pending_home_registration → registered
  - registration proposal emitted at foreign body (class: registration-proposal)
  - operator-visible packet emitted at activation site
  - home updates state/activations.md (operator-mediated)
  - home reads foreign body via credentials (once registered)
requires:
  - prior activation per agent/activate/SKILL.md
  - operator authorization for home registration and credential provisioning
calls:
  - agent/activate/SKILL.md
  - agent/attach/SKILL.md
---

<!--
section-manifest:
  planned: [frontmatter, core-principle, algorithm, define, unfold-state-machine, unfold-proposal-schema, unfold-credential-model, unfold-operator-packet, unfold-degraded-reporting, unfold-worked-example, rules, verify, failure-modes, references]
  completed: [frontmatter, core-principle, algorithm, define, unfold-state-machine, unfold-proposal-schema, unfold-credential-model, unfold-operator-packet, unfold-degraded-reporting, unfold-worked-example, rules, verify, failure-modes, references]
-->

# Registration

## Core Principle

**A fresh foreign body cannot self-register.** Registration is a three-gate sequence that only an operator can close: the foreign body emits a registration proposal; the operator authorizes it at home; the operator provisions credentials for home to read the foreign surface. Until all three gates close, the channel is `pending_home_registration` and the foreign body operates in a gated degraded state — not silently, not incorrectly, but explicitly named.

This is not a restriction on agent capability. It is a design commitment: the Writer Locality rule (one writer per surface) and the single-operator topology require that home's `state/activations.md` is written only at home. The foreign body cannot write it; the operator must. The registration skill names that gate explicitly and provides the operator with an actionable packet to close it.

The failure mode this skill prevents is **invisible gating** — a body that attempts cohere-join on a private repo, fails silently because home cannot read the private surface, and never surfaces the credential requirement to the operator. Invisible gating turns a structural gap into an unexplained silence. This skill ensures every gate is named, every degraded state carries a reason, and the operator always has a concrete action.

## Algorithm

1. **Emit proposal** — the fresh foreign body writes a `class: registration-proposal` entry as its first `.cn-{agent}/logs/<today>.md` entry. For private repos: the proposal is local evidence only; home cannot yet read it. For public repos: home reads it on next wake.
2. **Emit operator-visible packet** — the body outputs a registration packet to its activation site (terminal, stdout, or channel reply) in copyable form. This is the operator bridge: the packet carries all fields needed for the operator to create the home registration and provision credentials without reading the private log.
3. **Defer to operator** — the body cannot advance its own registration. It appends `deferred to operator: pending_home_registration at <hub>` to its attach entry and waits. No further channel work until registration completes.
4. **Home registers** (operator action) — operator creates `.cn-{agent}/state/activations.md` entry for the activation, provisions credentials (PAT or GitHub App installation), commits and pushes to home main.
5. **Home reads** — on next home wake, home reads the foreign body's log surface using the new credentials. Channel is now live.

---

## 1. Define

### 1.1. Three-gate sequence

Registration has three gates. Each is operator-only. None can be bypassed.

| Gate | Who closes it | What it unlocks |
|------|---------------|-----------------|
| **G1** Home registration | Operator edits `state/activations.md` at home | Home knows this activation exists |
| **G2** Read credentials | Operator provisions PAT or GitHub App at home | Home can clone and read foreign surface |
| **G3** Workflow write | Operator grants GitHub App `contents:write` at foreign repo | Wake worker can commit to foreign surface |

For public repos, G2 is trivially met (anonymous HTTPS clone). G1 and G3 still apply.

- ❌ "The body emits a proposal and is now registered"
- ✅ "Proposal emitted. G1–G3 remain operator-only gates. State: pending_home_registration."

### 1.2. State machine

Three states. One body per state at any moment.

```
unregistered
    │
    │ (body emits registration-proposal + operator-visible packet)
    ▼
pending_home_registration
    │
    │ (operator closes G1 + G2; home reads and confirms)
    ▼
registered
```

State is observable at two surfaces:
- **Foreign body**: last `.cn-{agent}/logs/` entry `class` field — `registration-proposal` means pending; subsequent entries with standard `class` values mean registered (channel is live).
- **Home**: `state/activations.md` entry presence — absent means unregistered or pending; entry with `last_read_foreign_log` populated means registered.

There is no `unregistered` state entry in `state/activations.md` — absence is the state. There is no separate pending-state file — the registration-proposal entry in the foreign log is the observable evidence for pending state.

- ❌ "If I don't see the activation in home's activations.md, I'll self-register"
- ✅ "If I don't see the activation in home's activations.md, state is unregistered or pending. Emit proposal and operator-visible packet. Stop. Defer."

### 1.3. Failure modes

Five named failure modes. Each has a structural fix.

- **R1 — Silent self-registration.** Body writes a log entry as if registered before G1 is closed. Operator never knows registration is pending; channel reads as live from foreign side but is invisible to home. Fix: §3.1 — never write standard attach entries before registration is confirmed; emit proposal-class entry and stop.
- **R2 — Missing operator packet.** Body emits proposal to its private log only. Home cannot read the private log (credential gate not yet closed). Operator cannot act without the packet. Fix: §2.4 — operator-visible packet is mandatory for private repos; even for public repos it is strongly recommended as the actionable bridge.
- **R3 — Credential model confusion.** Home attempts to clone a private foreign repo without credentials and fails silently. Operator does not know credentials are required. Fix: §2.3 — two named paths (PAT, GitHub App); each carries a concrete operator-action step.
- **R4 — G3 blind spot.** Body registers, home reads, but foreign body cannot push its log entries (GitHub App lacks `contents:write`). Channel silently receives directives but cannot acknowledge. Fix: §2.3 G3 note — G3 must be verified at registration time; failure surfaces as `degraded_reason: workflow_write_permission_missing`.
- **R5 — Foreign body proceeds past proposal.** Body emits proposal, then continues with substantive attach work on the same wake, producing entries home cannot read. Fix: §3.2 — stop after emitting proposal and operator-visible packet; only resumption cue is operator closing G1.

---

## 2. Unfold

### 2.1. State machine transitions

**Transition T1: unregistered → pending_home_registration**

Triggered when: foreign body runs attach foreign-inaugural procedure (per `agent/attach/SKILL.md §2.3`) and finds this activation absent from home's `state/activations.md`.

Body actions:
1. Write `class: registration-proposal` entry to `.cn-{agent}/logs/<today>.md` (§2.2).
2. Emit operator-visible packet to activation site (§2.4).
3. Append `deferred to operator: pending_home_registration at <hub>` to attach entry.
4. Stop. No further channel work this wake.

Precondition check: body MUST have confirmed the activation is genuinely absent from home's `state/activations.md` (not merely unchecked). If home clone failed for any reason, declare `degraded_reason: home_clone_failed` and stop — do not assume unregistered.

**Transition T2: pending_home_registration → registered**

Triggered when: operator adds activation entry to home `state/activations.md`, commits, pushes; home wake reads foreign log using new credentials and confirms readable surface.

Operator actions:
1. Add entry to `cn-{agent}:state/activations.md` (format per §2.5).
2. Provision credentials for home to read foreign repo (§2.3).
3. Commit and push to home main.
4. Optionally: post a directive to home activation thread confirming registration is open.

Body observes transition: next foreign wake reads home thread forward from cursor; finds no `deferred_reason: pending_home_registration` stop; proceeds with normal attach sync.

### 2.2. Registration proposal schema

The body's first log entry when `pending_home_registration`. Written to `.cn-{agent}/logs/<today>.md`. YAML frontmatter carries the registration proposal; body carries human-readable context.

```markdown
## YYYY-MM-DDTHH:MM:SSZ — registration proposal: <agent> at <hub>

---
at: <hub-name>
mode: foreign-activation
class: registration-proposal
proposed_hub: <home-hub-url>
proposed_activation_name: <short-name>
proposed_role: activation
repo_visibility: public | private
read_credential_request: none | pat | github-app
privacy_policy_state: undeclared | p0 | p1 | p2
cursor_in: ~
cursor_out: ~
---

Fresh body. This activation is not yet registered at home (<proposed_hub>).
Emitting registration proposal as first log entry.
Operator: close G1 (home registration) and G2 (read credentials) to activate the channel.

Registration packet emitted to activation site (stdout / terminal report).
```

Required fields:

| Field | Description |
|-------|-------------|
| `class` | `registration-proposal` — marks entry as a registration proposal, not a channel entry |
| `proposed_hub` | Full URL of the agent's home hub (e.g. `https://github.com/usurobor/cn-sigma`) |
| `proposed_activation_name` | Short name for this activation in home's registry (e.g. `cnos`, `bumpt`) |
| `proposed_role` | `activation` in v0 (peer roles defined separately) |
| `repo_visibility` | `public` or `private` |
| `read_credential_request` | `none` (public), `pat` (private PAT path), `github-app` (private App path) |
| `privacy_policy_state` | `undeclared` if no policy file found; else the declared mode from `.cn-{agent}/policy/privacy.yaml` |
| `cursor_in` | `~` (null) — no prior cursor at inaugural registration |
| `cursor_out` | `~` (null) — no read has been performed yet |

- ❌ Body writes a `class: substantive` entry before registration is confirmed
- ✅ Body writes `class: registration-proposal` with required fields, then stops

### 2.3. Home credential model

Home reads the foreign body's `.cn-{agent}/logs/` surface. For public repos, anonymous HTTPS clone suffices. For private repos, two paths:

**Path A — PAT (Personal Access Token)**

Operator provisions a GitHub PAT with `contents:read` scope on the private repo. Stores it as a repository secret at home:
```
CNOS_PRIVATE_READ_TOKEN=ghp_...
```
Home wake clones with:
```bash
git clone https://<token>@github.com/<owner>/<repo>.git
```

Operator action required:
1. Create PAT at `github.com/settings/tokens` with `repo` scope on the private repo.
2. Add as secret `CNOS_PRIVATE_READ_TOKEN` at home repo (`github.com/<owner>/<home-repo>/settings/secrets/actions`).
3. Ensure home wake reads using `CNOS_PRIVATE_READ_TOKEN` when cloning private activations.

Manual rotation: per v0; auto-rotation deferred.

**Path B — GitHub App**

Operator installs a GitHub App at both home and foreign repos. The App generates short-lived installation tokens that home uses for read access.

Operator action required:
1. Install the App at the foreign private repo (with at least `contents:read`).
2. Configure home to use App installation tokens when reading registered private activations.

Broader scope: GitHub Apps can be granted narrower permissions (read-only) than PATs. Preferred for production private-repo networks.

**G3 — Workflow write at foreign repo**

Independent of G2 (read credentials), G3 governs whether the wake worker can push log entries to the foreign repo. Required for the channel to function (foreign body writes `.cn-{agent}/logs/`).

G3 requires `contents:write` (or equivalent App permission) at the foreign repo. For GitHub Actions runners, this is the GITHUB_TOKEN permission in the foreign repo's workflow settings.

If G3 is not satisfied: body can read home directives but cannot acknowledge them via log entries. Degraded: `degraded_reason: workflow_write_permission_missing`.

### 2.4. Operator-visible registration packet

For private repos, the registration proposal in `.cn-{agent}/logs/` is not readable by home until G2 closes. The operator-visible packet is the bridge: a structured, copyable block emitted to the activation site (terminal report, stdout, or channel reply) that carries all fields the operator needs to create the home registration without reading the private log.

```
=== REGISTRATION PACKET — ACTION REQUIRED ===

proposed_hub:              <home-hub-url>
proposed_activation_name:  <short-name>
proposed_role:             activation
repo_visibility:           private
read_credential_request:   pat | github-app

Operator steps to complete registration:

G1 — Home registration:
  Add to <home-hub>:state/activations.md:

  - name: <activation-name>
    hub: https://github.com/<owner>/<repo>
    role: activation
    foreign_log: ".cn-{agent}/logs/"
    home_log: ".cn-{agent}/threads/activations/<activation-name>/"
    last_read_foreign_log: ~
    notes: <optional>

  Commit and push to home main.

G2 — Read credentials (private repos only):
  [If pat:] Add PAT with contents:read scope as CNOS_PRIVATE_READ_TOKEN at home repo.
  [If github-app:] Install App at this repo with contents:read; configure home to use App tokens.

G3 — Workflow write:
  Verify wake worker has contents:write at this repo (GITHUB_TOKEN permissions or App installation).

source_ref: .cn-{agent}/logs/<today>.md (local only; readable after G2)
=== END PACKET ===
```

Required fields in packet:
- `proposed_hub`, `proposed_activation_name`, `proposed_role`, `repo_visibility`, `read_credential_request`
- `privacy_policy_state` (if known)
- `source_ref` to local registration-proposal log path

For public repos: packet is still recommended (carries activation name and home registration steps) but `read_credential_request: none` and G2 section is omitted.

- ❌ Body emits only a private log entry; operator never sees the registration requirement
- ✅ Body emits packet to activation site; packet is self-contained and actionable without reading the private log

### 2.5. Degraded reporting

Registration gates produce degraded outcomes when unmet. These extend the canonical degraded-reasons list at cnos#444.

| Degraded reason | Trigger | Operator action |
|-----------------|---------|-----------------|
| `pending_home_registration` | Activation absent from home `state/activations.md` | Close G1 (home registration) |
| `pending_home_credentials` | Home `state/activations.md` has entry but `last_read_foreign_log: ~` and repo is private with no credentials configured | Close G2 (provision credentials) |
| `workflow_write_permission_missing` | Foreign wake cannot push log entries (`contents:write` denied) | Close G3 (grant write permission to workflow) |
| `home_clone_failed` | Body could not clone home hub (network, credentials, other) | Investigate home connectivity; re-run |

Degraded reporting in attach/cohere:
```
outcome: degraded
degraded_reason: pending_home_registration
operator_action: Close G1 — add activation entry to <home>:state/activations.md. See registration packet emitted to activation site.
```

No channel entry is written as a standard attach entry until `registered` state. The registration-proposal entry IS the appropriate first entry; subsequent entries require registered state.

### 2.6. Worked example — cn-sigma + bumpt (public repos)

The canonical v0 worked example uses two public repos: cn-sigma (home) and bumpt (foreign activation). Both public; G2 trivially met; only G1 and G3 required.

**At bumpt (foreign body, first wake):**
1. Body runs attach foreign-inaugural.
2. Finds bumpt absent from cn-sigma `state/activations.md` (first time).
3. Writes `class: registration-proposal` entry to `bumpt:.cn-sigma/logs/<today>.md`.
4. Emits operator-visible packet (G2 section omitted — public repo).
5. Defers to operator. Stop.

**Operator action (at cn-sigma):**
1. Adds bumpt entry to `cn-sigma:.cn-sigma/state/activations.md`:
   ```yaml
   - name: bumpt
     hub: https://github.com/usurobor/bumpt
     role: activation
     foreign_log: ".cn-sigma/logs/"
     home_log: ".cn-sigma/threads/activations/bumpt/"
     last_read_foreign_log: ~
   ```
2. Commits and pushes.

**Next foreign wake at bumpt:**
1. Body runs attach foreign-inaugural.
2. Finds bumpt registered in cn-sigma `state/activations.md`.
3. Reads home thread (empty — inaugural).
4. Writes inaugural attach entry with `class: inaugural`.
5. Channel is live.

**cn-rho** is the first pending registration that will exercise the private-repo path (G2 non-trivial) when cn-rho registers as a peer. The credential model in §2.3 applies at that point.

---

## 3. Rules

### 3.1. Never write standard entries before registered

A body MUST NOT write `class: substantive`, `class: heartbeat`, or `class: inaugural` entries before its activation is confirmed registered at home. The only permitted first entry is `class: registration-proposal`. Writing standard entries before registration produces a silent channel home cannot observe.

- ❌ "I'll write an inaugural entry; home will pick it up when credentials arrive"
- ✅ "State is pending_home_registration. Writing proposal entry and operator packet. Stop."

### 3.2. Stop after proposal + packet

After emitting the registration-proposal entry and operator-visible packet, stop. Do not continue with attach sync work. Do not guess that the operator has already closed G1. The next wake will re-run the state check; if G1 is closed, proceed to inaugural attach.

- ❌ Body emits proposal, then continues reading home thread and doing work
- ✅ Body emits proposal, emits packet, appends deferred note, stops

### 3.3. Packet is mandatory for private repos

For private repos, the operator-visible packet is mandatory — the only operator bridge before G2 closes. Emitting only the private log entry is R2 (missing operator packet). The packet must be emitted to a surface the operator can read without credentials.

- ❌ "The proposal is in the private log; the operator will find it"
- ✅ Emit packet to stdout / terminal / channel; operator can act without reading the private log

### 3.4. Confirm activation absence before declaring unregistered

Before declaring `pending_home_registration`, confirm the activation is genuinely absent from home's `state/activations.md` by reading the file. If home cannot be cloned, declare `home_clone_failed` and stop — do not assume unregistered.

- ❌ Body cannot reach home; concludes "not registered"; emits proposal for a hub that's already registered
- ✅ Body confirms absence by reading home state file; only then emits proposal

### 3.5. Operator owns G1 — home cannot be updated from foreign side

G1 (home registration) is a write to home's `state/activations.md`. Writer Locality applies: only home Sigma (or the operator) can write home's state files. The foreign body MUST NOT attempt to edit `state/activations.md` at home.

- ❌ Foreign body uses home credentials to update home `state/activations.md`
- ✅ Foreign body emits packet; operator updates home `state/activations.md` manually or via home wake

### 3.6. State is observable, not declared

The body does not store registration state in a separate file. State is observed from two canonical surfaces: presence in home `state/activations.md` (G1) and ability to clone with credentials (G2). Re-check on every inaugural attach attempt.

- ❌ "I saved 'pending' state in a local file last session; no need to re-check"
- ✅ Re-read home `state/activations.md` every inaugural attach to determine current state

---

## 4. Verify

### 4.1. State machine completeness check

Confirm the three states (unregistered, pending_home_registration, registered) and two transitions (T1, T2) cover all observable scenarios. Walk through each gate (G1, G2, G3) and confirm a concrete operator action closes each.

### 4.2. Proposal schema check

Confirm `class: registration-proposal` is distinct from all standard `class` values (`substantive`, `heartbeat`, `inaugural`, `directive-out`). Confirm all required fields are present in the schema (§2.2 table). Confirm `cursor_in: ~` and `cursor_out: ~` signal no read has occurred.

### 4.3. Credential model check

Confirm both PAT and GitHub App paths are concrete enough to execute without consulting another document. Confirm G3 (workflow write) is distinct from G2 (read credentials) and named as such.

### 4.4. Operator-visible packet check

Confirm the packet is self-contained. A reader with no access to the private repo should be able to execute all G1/G2/G3 operator actions from the packet alone.

### 4.5. Degraded-reasons coverage check

Confirm the four degraded reasons in §2.5 cover the named failure modes in §1.3. Confirm each reason has a concrete operator action.

### 4.6. Worked-example walkthrough

Simulate the cn-sigma + bumpt scenario (§2.6) step by step. Confirm the body stops after proposal emission (R5 prevention). Confirm home reads the proposal after G1+G2 close. Confirm the channel goes live on the second wake.

---

## 5. Failure modes catalogue

- **R1 — Silent self-registration.** _Symptom:_ Body writes a standard attach entry before G1 is closed. Home cannot observe the entry; operator is unaware the channel is live only on one side. _Fix:_ §3.1 — proposal entry only until registered.
- **R2 — Missing operator packet.** _Symptom:_ Private-repo body emits proposal to local log only. Operator cannot act. Channel stays pending indefinitely. _Fix:_ §3.3 — operator-visible packet mandatory for private repos.
- **R3 — Credential model confusion.** _Symptom:_ Home clones fail silently for private repos; no operator action is surfaced. _Fix:_ §2.3 — two named credential paths with concrete operator steps.
- **R4 — G3 blind spot.** _Symptom:_ Channel is readable by home but foreign body cannot push log entries. Directives arrive; no acknowledgement. _Fix:_ §2.3 G3 — verify workflow write at registration time; `degraded_reason: workflow_write_permission_missing` if absent.
- **R5 — Foreign body proceeds past proposal.** _Symptom:_ Body emits proposal and continues with substantive work; entries home cannot read accumulate in the foreign log. Next home read finds a backlog of unacknowledged entries. _Fix:_ §3.2 — stop after proposal + packet; do not continue with sync work.

---

## 6. References

### Upstream skills (registration calls into or is called by)

- `agent/activate/SKILL.md` — identity must be confirmed before registration check; activation is prerequisite.
- `agent/attach/SKILL.md §2.3` — foreign-inaugural attach procedure triggers registration when activation is absent.
- `docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md` — `state/activations.md` format and field definitions.

### Downstream skills (registration informs)

- `agent/cohere/SKILL.md` (`cnos#444`) — cohere private-onboarding flow invokes registration degraded reporting at G1/G2 gates.
- `agent/privacy/SKILL.md` (`cnos#460`) — registration packet carries `privacy_policy_state`; privacy mode is a co-condition with registration for private repos.

### Issue

- `cnos#449` — this skill's authoring issue. Acceptance criteria in §1–§5 of this skill correspond to AC1–AC6 of cnos#449.

### Empirical evidence

- cn-sigma `state/activations.md` — first registration (bumpt, 2026-05-30); PAT not required (public); G1 operator-closed. `cnos#431 / cnos#432`.
- cn-sigma `threads/adhoc/20260519-foreign-body-activation-gap.md` — pre-convention foreign-body activation gap; established the pattern this skill canonicalizes.
- `agent/attach/SKILL.md §1.3 A1` — "Silent self-registration" was named as attach failure mode A1 before this skill existed. Registration skill lifts the prevention mechanism from an attach-side note to a first-class procedure.

### Authority and stability

This skill is doctrine-adjacent: it defines the registration gate every cnos body traverses when entering a new activation. Future changes require constitutive-change approval. The operator-gate invariants (G1/G2/G3) MUST NOT be weakened without explicit operator authorization — they exist to prevent silent channel establishment in environments the operator has not authorized.
