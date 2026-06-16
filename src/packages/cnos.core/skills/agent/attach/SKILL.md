---
name: attach
description: Attach an already-activated agent to a channel surface. Detect whether the body is at home, at a foreign-activation hub, or running ephemerally; detect whether this is the inaugural attach or a follow-up sync; then execute the appropriate channel procedure per AGENT-ACTIVATION-LOG-v0.md. Source of truth for `cn attach` and for non-cn bodies that self-bind to a channel after activation.
artifact_class: skill
kata_surface: embedded
governing_question: Once a body knows which agent it is (post-activation), how does it find the channel its other selves have been writing on, catch up on what's been said, and contribute today's entry — without drifting the cursor or guessing on ambiguous binding?
visibility: public
parent: agent
triggers:
  - attach
  - channel sync
  - pull channel
  - cn attach
  - activate and attach
scope: task-local
inputs:
  - an activated agent identity (see `agent/activate/SKILL.md`)
  - access to the channel surfaces defined by `docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md`
  - body capabilities (shell + git, HTTP fetch only, or no fetch — see activate §2.2)
outputs:
  - body has bound to its channel (mode detected, attached-state detected)
  - body has either (a) initialized an inaugural cursor, or (b) caught up from prior cursor to other-side HEAD
  - body has appended today's entry and advanced its cursor
  - changes committed and pushed to `main` (skipped for ephemeral mode)
requires:
  - prior activation per `agent/activate/SKILL.md` — channel work without identity load produces drift
  - the convention spec `docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md` (single source of truth for loop mechanics)
calls:
  - agent/activate/SKILL.md
---

<!--
section-manifest:
  planned: [frontmatter, core-principle, algorithm, define, unfold-mode-detection, unfold-attached-detection, unfold-procedures, rules, verify, failure-modes, references]
  completed: [frontmatter, core-principle, algorithm, define, unfold-mode-detection, unfold-attached-detection, unfold-procedures, rules, verify, failure-modes, references]
-->

# Attach

## Core Principle

**Activation gives a body its identity; attach gives it its channel.** The two are distinct concerns. Activation is invariant doctrine — same procedure every body, every wake. Attach is environment-dependent — what the body does depends on _where_ it is right now and _whether it has been here before_.

A body that is activated but not attached has identity without continuity: it knows it is Sigma but does not know what Sigma-at-this-hub said yesterday, what directives are pending from home, or where to write today's entry. Attach closes that gap by binding the activated identity to the channel surface defined for this body in this environment.

The channel surfaces and cursor mechanics are defined in `docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md`. This skill does not redefine them — it operationalizes them: it tells a body how to (a) detect which channel surface applies here, (b) detect whether this is the inaugural attach or a follow-up sync, and (c) execute the correct procedure for the detected combination.

The failure mode this skill prevents is **silent channel drift** — a body that writes an entry without first reading what the other side said since the last cursor, or that writes a cursor without doing the read, or that self-registers a new activation without home's approval. Each of these is a cursor-integrity violation; together they erode the channel's value as durable evidence.

## Algorithm

1. **Detect mode** — observe `pwd` origin; classify as home, foreign-activation, or ephemeral (§2.1).
2. **Detect attached state** — check writer-surface for prior cursor evidence; classify as attached or not-attached (§2.2).
3. **Execute procedure**:
   - **Not attached → attach** (inaugural). Verify preconditions; initialize writer-surface; write the inaugural entry; establish the cursor at other-side HEAD as observed today.
   - **Attached → sync.** Walk the other side from prior cursor to HEAD; process inbound directives per standing posture; append today's entry; advance the cursor.
4. **Commit and push to `main`** for home and foreign-activation modes. Ephemeral emits a receipt only.

The four-step shape is total. Skipping detection (steps 1–2) produces the failures named in §1.3. Skipping push (step 4) leaves the cursor advance invisible to the other side — the next wake re-reads the same window.

---

## 1. Define

### 1.1. Identify the parts

A complete attach has these parts:

- **Mode** — which channel surface applies here (home / foreign-activation / ephemeral). Within foreign-activation: Tier 1a (substrate repo, has local agent identity files) vs Tier 1b (pure-product hub, identity from canonical). Channel surfaces are the same for both; shape is recorded in the entry's frontmatter.
- **Attached state** — whether prior cursor evidence exists in the body's writer-surface.
- **Reader-surface** — the other side's surface this body reads.
- **Writer-surface** — this body's own surface; the only place this body writes.
- **Cursor** — a Git commit SHA on the reader-surface marking how far this body has read.
- **Inbound directives** — content on the reader-surface since the cursor; what this body must read on this wake.
- **Today's entry** — what this body appends to its writer-surface on this wake, carrying the cursor in its frontmatter (`cursor_out`).

- ❌ "I'll write an entry and figure out the cursor later" (cursor without read = drift)
- ✅ "I read the directives from cursor to HEAD; the entry I write reflects them and pins the new cursor"

### 1.2. Articulate how they fit

Mode and attached state are orthogonal. Mode says _which files_ (paths on disk); attached state says _which procedure_ (inaugural attach or follow-up sync). The combination of (mode, attached) determines the full procedure. The convention's two-artifact shape (foreign log + home thread) is symmetric — both sides have a reader surface and a writer surface; the cursor lives in the writer-surface; the read goes from cursor to other-side HEAD.

The cursor is the only state that crosses wakes. Identity is reloaded every wake (via activate). Directives are read fresh from the reader-surface (no caching). Today's entry is written fresh to the writer-surface (no carry-over from prior wakes). The cursor is the one durable hand-off between wakes — which is why integrity discipline matters.

- ❌ "I'll cache directives across wakes" (caching defeats the read-fresh discipline)
- ✅ "Cursor pins what I've read; everything else is reloaded"

### 1.3. Name the failure mode

Attach fails in five named ways. Each has a structural fix in §2.

- **A1 — Silent self-registration.** Body lands at a hub not registered in home's `state/activations.md`, assumes the registration is its own responsibility, writes a `.cn-{agent}/logs/` entry, and silently establishes a channel home does not know about. Fix: §2.4 requires deferring to operator when preconditions are missing — never self-register.
- **A2 — Cursor without read.** Body writes `cursor_out` in the entry's frontmatter without walking the reader-surface from prior cursor to current HEAD. The cursor advances but the directives go unread. Fix: §2.5 procedure for follow-up sync makes the walk mandatory before writing `cursor_out`.
- **A3 — Read without cursor advance.** Body walks the reader-surface but forgets to set `cursor_out` (or sets it equal to `cursor_in` when it should advance). The next wake re-reads the same window — wasteful but not corrupting. Fix: §2.5 makes `cursor_out` required in every entry's frontmatter; missing or stale `cursor_out` is observable on review.
- **A4 — Wrong mode.** Body misdetects its environment (e.g., assumes home when at a foreign hub), reads/writes the wrong surfaces. Fix: §2.1 detection rule is explicit (compare `pwd` origin to registered URLs); ambiguous cases defer to operator.
- **A5 — Speculative entry on empty reader-surface.** Body finds the reader-surface empty (no prior directives from home), invents work to do, writes an entry that doesn't reflect the channel. Fix: §2.5 no-op-wake policy — when there's nothing on the reader side and no posture-aligned work, write a brief "no-op wake — channel heartbeat" entry, not speculative content.

---

## 2. Unfold

### 2.1. Mode detection

The body's environment determines which channel surface applies. Detection is `pwd` origin compared against known URLs:

| Mode | Detected when | Reader-surface | Writer-surface |
|------|---------------|----------------|----------------|
| **home** | `git -C $PWD remote get-url origin` matches the agent's home URL | each registered activation's `.cn-{agent}/logs/` (clone each hub read-only) | `threads/activations/{name}/YYYYMMDD.md` per activation in the home repo |
| **foreign-activation (Tier 1a — substrate repo)** | `pwd` origin matches a registered hub AND `spec/PERSONA.md` present at that hub (legacy root or `.cn-{agent}/`) | home's `threads/activations/{name}/` for this activation | `.cn-{agent}/logs/YYYYMMDD.md` in this hub |
| **foreign-activation (Tier 1b — pure-product hub)** | `pwd` origin matches a registered hub AND `spec/PERSONA.md` absent (both legacy root and `.cn-{agent}/`) | home's `threads/activations/{name}/` for this activation | `.cn-{agent}/logs/YYYYMMDD.md` in this hub |
| **ephemeral** | neither: no matching `pwd` origin (web session, container, ad-hoc) | spec + operator-passed context (no Git surface) | activation receipt → stdout, no commit |

The detection requires reading home's `state/activations.md` once to know which hub URLs are registered. The activated agent already has home cloned (per the activate skill's tier-(a) procedure); the registry lookup is a single file read.

**Foreign-activation sub-branch (Tier 1a vs 1b).** Within foreign-activation mode the channel surfaces (reader and writer) are the same regardless of hub shape — both Tier 1a and Tier 1b write to `.cn-{agent}/logs/` in the foreign hub and read from home's `threads/activations/{name}/`. The shape difference affects only how identity was loaded (handled by `activate/SKILL.md §2.5`, not this skill). The body records the shape in the entry's frontmatter for traceability; attach procedure is otherwise identical.

If `pwd` origin matches neither home nor a registered activation, the mode is ephemeral. Ephemeral is a real mode, not an error: it covers web sessions, containers, and ad-hoc bodies that have no persistent channel surface.

- ❌ Body assumes home because it cloned cn-sigma during activate (clone ≠ origin)
- ❌ Body assumes foreign-activation because cn-sigma is in its hub list (working dir matters, not what's been cloned)
- ❌ Body treats Tier 1b (no local spec) as ephemeral just because there's no local PERSONA.md — if origin matches a registered hub, it's foreign-activation
- ✅ Body runs `git -C $PWD remote get-url origin`, compares to home URL and registered activation URLs, picks the match

### 2.2. Attached-state detection

Within home and foreign-activation modes, the body is either attached (has prior cursor evidence in its writer-surface) or not attached (inaugural). Detection is a one-shot grep over the writer-surface:

| Mode | Detection |
|------|-----------|
| **home** | for each registered activation: does `threads/activations/{name}/` contain any prior entry? Per-activation, not global. |
| **foreign-activation** | does `.cn-{agent}/logs/` contain any entry whose frontmatter declares `cursor_out: <agent>@<sha>`? |
| **ephemeral** | always not-attached (no persistence across wakes) |

The per-mode detection differs because the writer-surfaces differ. Home writes a thread per activation, so attached state is per-activation (some activations may be attached, others may be inaugural — handle each independently in step 3). Foreign-activation writes a single log surface in the local hub, so attached state is global to this body.

The cursor lives in the entry's YAML frontmatter as `cursor_out: <agent>@<sha>` (convention §4). Scan **bottom-up** to find it: most recent file in `.cn-{agent}/logs/`, last H2 entry in that file, its `cursor_out` value. That's the active cursor; earlier entries are historical. Bottom-up scan finds it in O(1) — do not parse every entry from the top.

- ❌ Body checks if `.cn-{agent}/logs/` directory exists (existence ≠ attached; empty dir is not-attached)
- ❌ Body greps for any content in the dir (content ≠ cursor; speculative content from earlier failures is not a cursor)
- ✅ Body greps for `cursor_out:` in `.cn-{agent}/logs/`; if any entry has it, attached; else not-attached

### 2.3. Inaugural attach procedure

Triggered when (mode ∈ {home per-activation, foreign-activation}) and (attached state = not attached).

**Foreign-activation inaugural attach:**
1. **Verify precondition.** This hub's URL must be registered in home's `state/activations.md`. If not, append `deferred to operator: foreign-activation inaugural attach requires home registration at <hub-url>` to the body's working scratchpad and stop. Do not self-register; do not initialize `.cn-{agent}/logs/`. Home must register first.
2. **Initialize writer-surface.** Create `.cn-{agent}/logs/` if it does not exist. (For first ever attach, this is the directory creation.)
3. **Read all of home's thread.** Walk `cn-{agent}:threads/activations/{name}/` from start to home HEAD. The thread may be empty; that's a valid signal — there are no inbound directives yet.
4. **Write inaugural entry.** Append today's entry to `.cn-{agent}/logs/$(date -u +%Y%m%d).md` per the convention's §4 entry format. The H2 header carries the wake's full UTC timestamp: `## $(date -u +%Y-%m-%dT%H:%M:%SZ) — Inaugural attach at <hub> as Sigma-at-{name}`. Body: name the inaugural binding (e.g., "Walked home thread; no prior directives.").
5. **Pin inaugural cursor.** Set `cursor_in: <agent>@<sha>` and `cursor_out: <agent>@<sha>` in the entry's frontmatter, both equal to home `main` HEAD as observed in step 3. Set `class: inaugural`.
6. **Commit and push** to local hub's `main`.

**Home inaugural attach (per activation):**
1. **Verify precondition.** This activation must be registered in `state/activations.md`. If `state/activations.md` does not yet name this activation, the registration step is operator-driven — append `deferred to operator: home inaugural attach for activation <name> requires registry entry` and stop.
2. **Initialize writer-surface.** Create `threads/activations/{name}/` if it does not exist.
3. **Read all of activation's foreign log.** Clone the activation's hub read-only; walk `.cn-{agent}/logs/` from start to hub HEAD. May be empty.
4. **Write inaugural entry.** Append to `threads/activations/{name}/$(date -u +%Y%m%d).md`. The H2 header carries the wake's full UTC timestamp per convention §4. Body: name the inaugural binding from home's side.
5. **Advance cursor.** Set `last_read_foreign_log: <sha>` in `state/activations.md` for this activation, where `<sha>` is hub `main` HEAD as observed in step 3.
6. **Commit and push** to home's `main`.

### 2.4. Follow-up sync procedure

Triggered when (mode ∈ {home per-activation, foreign-activation}) and (attached state = attached).

**Foreign-activation sync:**
1. **Locate prior cursor.** Find the most recent entry in `.cn-{agent}/logs/` (most recent file, last H2 entry within it). Read `cursor_out: <agent>@<sha>` from its frontmatter. This `<sha>` is the prior cursor.
2. **Walk home thread forward.** Read `cn-{agent}:threads/activations/{name}/` from `<sha>` to home `main` HEAD. The walk may span multiple per-day files (one per day since the cursor); read in chronological order.
3. **Process inbound.** For each directive: apply per standing posture (read from agent's spec and from prior log entries). If a directive is genuinely ambiguous, append `deferred to operator: <reason>` to today's entry and proceed with what is clear. Do not guess on substantive calls.
4. **Do the work.** Execute what directives + posture call for. If nothing actionable, the entry is a no-op-wake heartbeat (see A5 fix in §1.3).
5. **Append today's entry.** Add to `.cn-{agent}/logs/$(date -u +%Y%m%d).md` (create today's file if not present). The H2 header carries the wake's full UTC timestamp: `## $(date -u +%Y-%m-%dT%H:%M:%SZ) — <short subject>`. Free-form markdown per convention §4.
6. **Pin new cursor.** Set `cursor_in: <agent>@<prior-sha>` and `cursor_out: <agent>@<new-sha>` in the entry's frontmatter. `<new-sha>` is home `main` HEAD as observed in step 2; equal to `<prior-sha>` if the walk found nothing (heartbeat case). Set `class: substantive` if work was done, `heartbeat` if no-op.
7. **Commit and push** to local hub's `main`.

**Home sync (per registered activation):**
1. **Locate prior cursor.** Read `last_read_foreign_log: <sha>` from `state/activations.md` for this activation.
2. **Walk activation's foreign log forward.** Clone the hub read-only; read `.cn-{agent}/logs/` from `<sha>` to hub `main` HEAD. May span multiple per-day files.
3. **Process inbound.** Read each entry; understand what the activation reported, what it ACK'd, what `cursor_out` it set (the cursor in frontmatter is an implicit ACK of the directive at that SHA).
4. **Do the work.** Author home-side response: directives, intake notes, registry updates. The work depends on what the activation surfaced.
5. **Append today's entry.** Add to `threads/activations/{name}/$(date -u +%Y%m%d).md`. H2 header per §4 with the wake's UTC timestamp.
6. **Advance cursor.** Update `last_read_foreign_log: <new-sha>` in `state/activations.md` to the hub HEAD observed in step 2.
7. **Commit and push** to home's `main`.

### 2.5. Ephemeral attach procedure

Triggered when mode = ephemeral. Always not-attached (no persistence).

1. **Read spec + operator-passed context** if any. The body's identity comes from activate; the channel context comes from whatever the operator pasted in (could be a fragment of a log file, a directive, or nothing).
2. **Do the work** within the scope of the operator's prompt.
3. **Emit activation receipt** — a structured one-shot output the operator can later route to a durable channel:
   ```
   ## YYYY-MM-DD — ephemeral activation receipt
   
   Body: <body description, e.g. "Claude.ai web session, 2026-06-01">
   Agent: <agent-name>
   Operator-passed context: <brief summary>
   Work done: <brief summary>
   Operator action requested: <what the operator should do with this receipt, e.g. "paste into cnos:.cn-{agent}/logs/20260601.md">
   ```
4. **No commit, no push.** Ephemeral has no Git surface.

The receipt is the operator's bridge from an ephemeral session back into the durable channel. The body cannot do this routing itself — that's what makes the mode ephemeral.

---

## 3. Rules

### 3.1. Attach only after activate

Channel work requires identity. A body that has not completed activation per `agent/activate/SKILL.md` (Kernel + CA + Persona + Operator + state + identity confirmation) must not attempt attach. Identity drives how directives are interpreted, what "standing posture" means, and what counts as in-scope work. Attaching without identity produces entries that misread the channel.

- ❌ "I'll attach first and load Persona later if needed"
- ✅ "Activate completed; identity confirmed; now attach"

### 3.2. Detect mode before doing anything channel-shaped

Run `git -C $PWD remote get-url origin` and compare to known URLs before touching any channel surface. Mode determines which surfaces apply; reading or writing the wrong surface corrupts the channel.

- ❌ "I'll assume home because that's the common case"
- ✅ "Origin = X; X matches the hub registered as 'cnos'; foreign-activation mode"

### 3.3. Detect attached state before choosing procedure

Inaugural attach and follow-up sync are different procedures. Run the detection (§2.2) before executing either. Skipping the check produces either (a) inaugural-mode work at an already-attached channel (writes a duplicate inaugural entry), or (b) follow-up-mode work at an unattached channel (tries to walk a nonexistent cursor).

- ❌ "I'll assume follow-up because the channel should already be attached"
- ✅ "Grep for cursor pattern; if found, follow-up; if not, inaugural"

### 3.4. Never self-register

If a foreign-activation inaugural attach finds the hub unregistered in home's `state/activations.md`, stop with `deferred to operator`. Do not write a `.cn-{agent}/logs/` entry; do not modify home's registry from the foreign side. Registration is a home-side concern owned by the operator.

- ❌ "I'll write an entry; home will figure it out"
- ✅ "Hub not registered; deferred to operator; no writer-surface initialization"

### 3.5. Cursor lives in entry frontmatter

Every entry's YAML frontmatter declares `cursor_in` and `cursor_out` (foreign side). Home side declares the same plus advances `last_read_foreign_log` in `state/activations.md`. Entries without `cursor_out` in frontmatter are malformed; A3 (read without cursor advance) is observable as `cursor_in == cursor_out` when there were unread commits to process.

- ❌ Entry body mentions a cursor in prose only; frontmatter missing
- ✅ Frontmatter declares both `cursor_in` and `cursor_out` explicitly

### 3.6. No-op wake writes a heartbeat, not silence

When the reader-surface is empty (no inbound directives) and no posture-aligned work is pending, the entry is a brief heartbeat: `no-op wake — channel heartbeat; cursor advance only`. Silence (skipping the entry entirely) breaks the cursor model — the next wake has no record that this wake happened.

- ❌ "Nothing to do; skip the entry"
- ✅ "Nothing to do; write a heartbeat that pins today's cursor"

### 3.7. Append-only on log files

Per convention §4, log files are append-only. Edits to prior entries violate the trace. If a prior entry is wrong, the correction is a new entry, not an edit. Mojibake and obvious typos are corrected only in the same commit that originally introduced them.

- ❌ Edit yesterday's entry to fix a misstatement
- ✅ New entry today: "yesterday's entry stated X; correction: Y"

### 3.8. Defer to operator on ambiguity, do not guess

If mode, attached state, registration status, or directive interpretation is ambiguous, append `deferred to operator: <reason>` to the body's scratchpad (or to today's entry if the entry is already being written) and stop. The convention's single-operator topology means operator routing is fast; guessing produces drift that compounds.

- ❌ "Reasonable interpretation; proceed"
- ✅ "Ambiguous between X and Y; deferred to operator; stop"

---

## 4. Verify

### 4.1. Mode-detection check

Confirm the body runs `git -C $PWD remote get-url origin` and compares against home URL + registered activation URLs. Walk through each mode's detection logic; confirm a body could execute without consulting another document.

### 4.2. Attached-state check

Confirm the writer-surface grep procedure for both home (per-activation thread directory) and foreign-activation (`.cn-{agent}/logs/` for `cursor_out:` in entry frontmatter). Confirm the frontmatter field name is exact.

### 4.3. Procedure-completeness check

For each (mode, attached-state) combination — five total: home-inaugural, home-sync, foreign-inaugural, foreign-sync, ephemeral — confirm the procedure is concrete enough to execute without consulting another document.

### 4.4. Convention-consistency check

Confirm this skill does not redefine the convention's two-artifact shape, cursor mechanism, entry format, or trust boundary. This skill is operational; the convention is structural. Drift between this skill and the convention is a binding finding — fix both before merge.

### 4.5. Author-test on a fresh foreign-activation inaugural

Run the foreign-activation inaugural procedure (§2.3) against a fresh hub. Confirm: (a) preconditions are checked before initialization, (b) home thread is read fully, (c) inaugural entry pins cursor at home HEAD, (d) commit/push lands on local main, (e) writer-surface integrity is preserved.

---

## 5. Failure modes catalogue

The five named failure modes from §1.3 plus three implementation-level ones:

- **A1 — Silent self-registration.** _Symptom:_ Foreign activation appears in `.cn-{agent}/logs/` without `state/activations.md` registration. _Fix:_ §3.4 + §2.3 precondition.
- **A2 — Cursor without read.** _Symptom:_ Cursor pin in entry but no evidence of having walked the reader surface from prior cursor. _Fix:_ §2.4 follow-up sync requires the walk.
- **A3 — Read without cursor advance.** _Symptom:_ Entry body reflects reader-surface walk but frontmatter has `cursor_out == cursor_in`. Next wake re-reads same window. _Fix:_ §3.5 — `cursor_out` must reflect the actual HEAD observed during the walk.
- **A4 — Wrong mode.** _Symptom:_ Body writes to home's `threads/activations/` when activated at a foreign hub. _Fix:_ §3.2 mandatory mode detection.
- **A5 — Speculative entry on empty reader-surface.** _Symptom:_ Body invents work; entry doesn't ground in inbound directives. _Fix:_ §3.6 no-op heartbeat policy.

Implementation-level:

- **A6 — Cursor parse error.** _Symptom:_ Body finds malformed YAML frontmatter (typo, partial match, missing `cursor_out`) and either crashes or misreads. _Fix:_ strict YAML parse of frontmatter block; if no valid `cursor_out`, treat as not-attached and defer to operator.
- **A7 — Concurrent wake collision.** _Symptom:_ Two wakes write to the same `YYYYMMDD.md` and one's push is rejected. _Fix:_ append-only with rebase-on-push; if rebase fails, defer to operator (the convention's §7 single-writer caveat covers this).
- **A8 — Stale clone for read.** _Symptom:_ Body uses an old clone of home and reads stale directives. _Fix:_ `git fetch + git log origin/main` immediately before the walk; cursor compared against `origin/main` HEAD, not local HEAD.

---

## 6. References

### Convention this skill operationalizes

- `docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md` — the two-artifact channel shape, cursor mechanism, entry format, trust boundary, single-writer caveat. This skill does not redefine these; it operationalizes the post-activation procedure for binding a body to a channel.

### Peer skill

- `agent/activate/SKILL.md` — identity load (Kernel + CA + Persona + Operator + state + identity confirmation). Attach requires activate to have completed first. §2.5 defines the Tier 1a/1b foreign-body shape detection; the shape detected there is what the body records in today's attach entry.

### Composition pattern

- The combined prompt `Activate and attach as <agent-home-url>` invokes both skills in sequence. `cn activate && cn attach` (when `cn attach` lands) is the CLI equivalent. The wake workflow at each hub uses the prompt form.

### Reflections this skill derives from

- `cnos:.cn-sigma/logs/20260601.md` — claude-code-action substrate verification; established that GH Actions can run the channel work end-to-end with operator's subscription billing.
- `cnos:docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md` — the convention this skill operationalizes; field experience across three activations (Sigma-at-cnos, Sigma-at-bumpt, Sigma-at-home) confirmed the topology-fit framing.

### Authority and stability

This skill is doctrine-adjacent: it defines the procedure every cnos body follows after identity activation to bind to a channel. Future changes follow the constitutive-change approval discipline that governs other doctrine-adjacent skills. Drift between this skill and the convention spec it operationalizes is resolved in favor of the convention; this skill follows.
