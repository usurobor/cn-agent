---
name: harness
description: Harness / platform-driver substrate for CDD. Codifies dispatch mechanics, observability contract, worktree management, identity discipline, polling, and timeout recovery. Invoked by δ; referenced by γ and β.
artifact_class: skill
kata_surface: embedded
governing_question: How does the harness substrate execute dispatch, surface observability, isolate worktrees, and enforce identity — coherently — across role sessions?
visibility: internal
parent: cdd
triggers:
  - harness
  - dispatch-mechanics
  - worktree
  - observability
  - identity
scope: global
inputs:
  - δ's dispatch decisions (which role, which branch, which prompt)
  - γ's α/β prompts (verbatim)
  - active cycle branch state on origin
outputs:
  - dispatched role sessions (one `claude -p` / `cn dispatch` invocation at a time)
  - per-cycle JSONL observability streams
  - per-cycle worktree directories (when parallel α applies)
  - per-actor git identity (worktree-local when `extensions.worktreeConfig=true`)
requires:
  - active CDD cycle
  - cycle branch exists on origin
calls:
  - operator/SKILL.md
---

# Harness (CDD substrate)

## Core Principle

**The harness executes what δ decides — dispatch, observability, worktree, identity — without inventing policy.**

The harness is below the role layer: α/β/γ/δ are the role surfaces; harness is the substrate they all run on. Policy (when to dispatch, when to gate, when to override) lives in `operator/SKILL.md` (δ doctrine), `gamma/SKILL.md` (γ doctrine), and the α/β skills. Mechanics (how the dispatch shell is invoked, how identity is written, how worktrees are pre-created, how an exited session is recovered) live here.

The failure mode is **policy creep into mechanics** — δ writes a shell snippet that quietly carries δ-role policy, or β re-derives a worktree rule because the harness contract is implicit. The harness is the single named home for these mechanics; cross-references from `operator/SKILL.md`, `gamma/SKILL.md`, and `beta/SKILL.md` point here so the role skills stay focused on their own discipline.

## Load Order

When acting as δ (or as γ/β consulting harness contracts):

1. Load `CDD.md` as the canonical lifecycle.
2. Load `operator/SKILL.md` as the δ-role surface (the harness serves δ's gate authority and route discipline).
3. Load this file when executing dispatch mechanics, identity writes, worktree management, polling, or timeout recovery.
4. β consults §Git identity for the pre-merge gate Row 1 identity-truth row.
5. γ consults §Dispatch observability contract when authoring α/β prompts that name the observable-output flag.

This file does **not** redefine δ's gate authority, γ's coordination scope, or β's review independence. It states what the substrate guarantees so the roles can refer to it instead of restating it.

## Phase 4 lineage

This skill consolidates harness substrate content surfaced and pinned across cycles:

- **cnos#371** — dispatch observability contract; codified here as §Dispatch observability contract.
- **cnos#373** — preventive `--worktree` identity write when `extensions.worktreeConfig=true`; codified here as §Git identity for role actors (with explicit worktree subsection).
- **cnos#384** — parallel α dispatch requires per-cycle pre-created worktrees; codified here as §Parallel dispatch precondition.
- **cnos#398** — this cycle; extracts the harness substrate from `operator/SKILL.md` (Phase 4b of cnos#366).

Phase 4a (cycle/397) carves the δ-role boundary into `delta/SKILL.md`. Phase 4c (cycle/399) carves the release-effector surface. The three sibling cycles together complete the operator/SKILL.md split that #366 Phase 4 requires.

---

## 1. Dispatch invocation

### 1.1. One role at a time

δ dispatches γ, then α, then β, sequentially — one `claude -p` (or `cn dispatch`) process at a time. Sequential dispatch keeps memory pressure low, gives δ direct visibility into each session, and isolates failures: if α dies, δ retries α without losing γ or β state.

The role-execution policy (γ first, α before β, etc.) lives in `operator/SKILL.md` §Algorithm. This subsection documents the *mechanics* of the invocation: the shell command, the required flags, and the prompt-delivery shape.

### 1.2. Invocation shape

```bash
# γ — reads issue, creates branch, produces α/β prompts
cat /tmp/gamma-prompt.md | claude -p \
  --allowedTools "Read,Write,Bash" \
  --output-format stream-json --verbose \
  --model <model>

# α — implements (needs Bash for git operations)
cat /tmp/alpha-prompt.md | claude -p \
  --allowedTools "Read,Write,Bash" \
  --output-format stream-json --verbose \
  --permission-mode acceptEdits \
  --model <model>

# β — reviews, merges (needs Bash for git/gh read + git merge)
cat /tmp/beta-prompt.md | claude -p \
  --allowedTools "Read,Write,Bash" \
  --output-format stream-json --verbose \
  --permission-mode acceptEdits \
  --model <model>
```

**Required flags (mechanics, not policy):**

- `--output-format stream-json --verbose` — observability contract. See §2 below. Black-box dispatch violates δ's gate authority; `--verbose` is mandatory when combining `--output-format stream-json` with `-p` (Claude Code exits with an error otherwise).
- `--permission-mode acceptEdits` — required for fresh `claude -p` sessions to bypass the trust dialog and allow file writes.
- `--allowedTools "Read,Write,Bash"` — β gets Bash because it needs `git diff`, `gh issue view`, and `git merge`; role boundaries are enforced by the role skills, not by tool scoping.

### 1.3. `cn dispatch` — the identity-rotation primitive

The Go subcommand `cn dispatch` (under `src/go/internal/dispatch/`) wraps the `claude -p` invocation pattern above and adds identity-rotation: fresh Claude CLI session per role, branch/artifact continuity preserved, dispatch backend pluggable (claude / stub / print for testing). The default backend is `claude -p` with stream-json observability per §2.

```bash
cn dispatch --role α --branch cycle/<N>
cn dispatch --role β --branch cycle/<N>
```

`cn dispatch` is the canonical δ-side dispatch entry-point. The `claude -p` shell forms in §1.2 are the interim form used when `cn dispatch` is unavailable or for explicit shell-level control. Both honor the same observability and identity contracts.

The `cn dispatch` Go implementation is **codified, not reimplemented** by this skill. Behavior changes go through the Go package; this skill documents the contract the Go code obeys.

### 1.4. Prompt-delivery rules

- Deliver the prompt verbatim to the `claude -p` session — do not rewrite it to add constraints δ thinks are missing.
- Run one role at a time; inspect artifacts between dispatches.
- Name the agent-to-role mapping before delivering prompts (which terminal / which subagent is which role).

These are mechanics of the prompt-handoff; the doctrine that γ owns prompt quality and δ does not rewrite prompts lives in `operator/SKILL.md` §6 "What the operator does NOT do."

---

## 2. Dispatch observability contract

### 2.1. The contract

**Cycle dispatch via `claude -p` (or `cn dispatch`) MUST use `--output-format stream-json` and capture stdout to a per-cycle JSONL log path. The rule is normative (MUST), not advisory.**

Black-box dispatch — running `claude -p` with default text output, which emits only the final assistant message on process exit — violates δ's gate authority. δ holds the boundary on the cycle; the cycle is observable to δ in real time, or δ cannot exercise the gate at SIGTERM, at cost overruns, or at mid-cycle drift.

The contract:

| Surface | Requirement |
|---|---|
| Flag | `--output-format stream-json --verbose` on every `claude -p` invocation (γ, α, β, re-dispatches). |
| Log path | Per-cycle JSONL log at `/tmp/cycle-{N}.jsonl` (or `.cdd/unreleased/{N}/dispatch.jsonl` if δ chooses to retain in-cycle). |
| Consumer | δ (the parent session) reads the stream; for §5.2 single-session δ-as-γ, the same parent session reads it. |
| Wake-up | Each JSONL line is a `task-notification`-class event; δ wakes on tool calls, subagent dispatches, commits, SIGTERMs. |
| Cycle hygiene | δ runs `touch /tmp/session-start-sentinel` before each dispatch so §6 timeout recovery can identify session-newer artifacts. |

### 2.2. JSONL line shape

Each line in the stream is a JSON object with at minimum:

- `type` — `"system"` / `"assistant"` / `"user"` / `"result"` (terminal).
- `subtype` (when `type == "system"`) — initialization metadata, model, session id.
- `message.content[]` (when `type == "assistant"`) — array of text/tool_use blocks the model emitted in one turn.
- `message.content[]` (when `type == "user"`) — tool_result blocks returned to the model.
- `usage` (in `result`-type) — token + cost telemetry for cost/time gating.

The schema is owned by Claude Code; this skill names the shape so future role sessions reading the stream know what they're parsing. Refer to `claude -p --help` and the Claude Code SDK docs for the authoritative reference.

### 2.3. Tail and gate

```bash
# Tail the stream during dispatch
tail -f /tmp/cycle-<N>.jsonl

# Wake on tool calls, commits, SIGTERMs — operator can intervene early
jq -c 'select(.type == "assistant" and (.message.content[]? | .type == "tool_use"))' /tmp/cycle-<N>.jsonl
```

The stream enables:

- α SIGTERM detectable in seconds; δ-recovery (see §6) runs before further work is wasted.
- Cost/time gating mid-cycle — δ can stop a runaway dispatch.
- Mid-cycle intervention — δ surfaces an override before damage compounds.
- Close-token gap (#368) early-detection — γ's missing `gh issue close` surfaces in the stream, not only after process exit.

### 2.4. Falsification anchor

This contract exists because the alternative was empirically tested and failed.

**Cycle #369 + #370 parallel dispatch, 2026-05-17.** Both cycles dispatched in default text mode produced zero in-process visibility — only filesystem-side `git log` polling on the cycle branches gave proxy progress signal. If either α had SIGTERMed mid-cycle (as in #365 and #367 precedents), δ would not have known until the entire γ subprocess exited. The pattern normalized the black box; δ stopped trying to inspect.

After #371 freezes the contract: every dispatched cycle is observable by skill rule; black-box dispatch is no longer doctrine-compliant.

### 2.5. Mirror in γ

γ's α/β dispatch (when γ executes dispatch directly, e.g. under §5.2 single-session δ-as-γ) carries the same MUST. γ does not get a weaker observability contract because γ's subprocess dispatch is the same primitive as δ's; see `gamma/SKILL.md` §2.5 dispatch prompt format and the `Identity-rotation primitive` line.

- ❌ "Dispatch with default text mode; check the stream after exit."
- ❌ "Skip `--verbose` because the schema looks busy."
- ✅ Every `claude -p` invocation passes `--output-format stream-json --verbose`; the JSONL log path is named in the dispatch record.

---

## 3. Git identity for role actors

### 3.1. Canonical form

Every CDD role actor configures a git identity in the form `{role}@{project}.cdd.cnos` before making any commits on the cycle branch. DNS domains read broad-to-narrow right-to-left: `cnos` is the origin repository where the cdd protocol is defined and versioned, `cdd` is the protocol namespace inside cnos, and `{project}` is the tenant project running the protocol. The role name is the local part. This form makes the protocol's origin repo visible in every commit trailer and leaves namespace room for sibling protocols (`cnav`, `cnobs`) under the same cnos root.

**Special case — cnos itself.** When the project running the cycle is the cnos repo, the literal form would be `{role}@cnos.cdd.cnos` (redundant `cnos`). The canonical elision is `{role}@cdd.cnos`, which reads as "the cdd protocol at cnos." Existing cnos commit trailers already use this form; the redundancy adds no information.

| Role | Project | Canonical identity | Notes |
|------|---------|-------------------|-------|
| alpha | tsc | `alpha@tsc.cdd.cnos` | tsc project actor |
| beta | cnos | `beta@cdd.cnos` | cnos actor — elision form (see above) |
| gamma | acme | `gamma@acme.cdd.cnos` | hypothetical third project |
| beta | * | `beta@cdd.{project}` | **(deprecated)** — cycle #287 form; cycle #343 cutover |

### 3.2. Worktree-aware identity write (preventive `--worktree` discipline)

**Binding rule.** At session start, before the first identity write, the role actor runs:

```bash
git config --get extensions.worktreeConfig
```

If this returns `true`, **every subsequent identity write MUST use the `--worktree` flag from the first command, not from a recovery commit.** A plain `git config user.email X` without `--worktree`, when `extensions.worktreeConfig=true`, writes to the *shared* `.git/config`. The immediate `git config --get user.email` returns X — but any subsequent process (a sibling worktree's command, a hook, an unrelated tool) that writes to the shared layer overwrites the value. The next commit-time read returns the overwriting role's identity. Role-identity-is-git-observable (`operator/SKILL.md` §"Git identity for role actors", `review/SKILL.md` §Review identity, `beta/SKILL.md` Pre-merge gate Row 1) is then silently violated; the audit trail records a different actor than the role contract names.

```bash
# Step 1 — check the layered-config flag at session start
if [ "$(git config --get extensions.worktreeConfig)" = "true" ]; then
  # Worktree-aware path: write to the worktree-local config layer.
  git config --worktree user.name "{role}"
  git config --worktree user.email "{role}@{project}.cdd.cnos"
else
  # Shared-config path: write to the standard .git/config layer.
  git config user.name "{role}"
  git config user.email "{role}@{project}.cdd.cnos"
fi

# Step 2 — verify
git config --get user.name
git config --get user.email
```

For the cnos project (elision form), the email is `{role}@cdd.cnos`; for tenant projects, `{role}@{project}.cdd.cnos`.

### 3.3. Failure modes the rule catches

| Surfacing | Cycle | What happened |
|---|---|---|
| First | #301 O8 | Merge-test worktree's identity write leaked to shared `.git/config`; subsequent commits authored as `beta-merge-test`. |
| Second | #370 β R1 §2.1 row 1 | Merge-test worktree's `extensions.worktreeConfig=true` toggle flipped the shared layer; main-worktree's `user.email` reverted to a stale shared value after merge-test teardown. β re-asserted with `--worktree`. |
| Third | #370 α F4 | α's session-start `git config user.email "alpha@cdd.cnos"` (no `--worktree`) wrote to the shared layer; γ's subsequent identity write overwrote it; α's first commit carried the wrong author and required `--amend --reset-author` to fix. |

The fix in #301 was local to β's merge-test recipe (worked example) — it did not generalize to every role's session-start identity write. Cycle #373 surfaced this gap; this cycle (cnos#398) absorbs #373 and writes the preventive rule into the harness substrate so every role-skill's identity-write step inherits it by reference.

### 3.4. Cross-references

- **`beta/SKILL.md` Pre-merge gate Row 1** ("Identity truth") — the β-side enforcement. Row 1 verifies the canonical email at merge time and re-asserts under `--worktree` if leaked. The doctrine here in §3.2 is what Row 1 references for the preventive rule.
- **`alpha/SKILL.md` §2.6 row 14** — α's pre-review-gate identity check. The corrective path (a) (`--amend --reset-author && force-push`) and path (b) (configure correctly going forward) both depend on the rule named here.
- **`release/SKILL.md` §2.1** — merge-test recipe worked example. The recipe's `git worktree add` step pairs with §3.2 here for the per-worktree config write under `--worktree`.

### 3.5. Recovery — `--amend --reset-author`

If a commit lands under the wrong identity (the failure mode in §3.3), the recovery is:

```bash
# Re-assert correct identity (use --worktree if §3.2 applies)
git config --worktree user.email "{role}@{project}.cdd.cnos"

# Rewrite recent commits with the correct author
git rebase -i {merge-base} --exec 'git commit --amend --reset-author --no-edit'

# Force-push (only on cycle branches; never on main)
git push --force-with-lease origin cycle/{N}
```

This is reactive; the preventive rule in §3.2 is one flag away and saves the rebase round. See `alpha/SKILL.md` §2.6 row 14 for the α-side path (a) / path (b) decision.

- ❌ `git config user.email "{role}@cdd.cnos"` without checking `extensions.worktreeConfig` first.
- ❌ Skip the `--worktree` flag because "the first read returns the right value."
- ✅ Check `extensions.worktreeConfig`; write with `--worktree` if true; verify with `git config --get user.email` after writing.

---

## 4. Parallel dispatch precondition

### 4.1. Same-repo parallel α requires per-cycle worktrees

**Binding rule.** When δ dispatches two or more α sessions targeting the same repo concurrently, each α MUST receive a pre-created `git worktree add` directory before dispatch. Two α sessions trying to mutate the same `.git` state simultaneously is not a coherent operation; each α needs its own working copy and its own `.git/worktrees/{N}` lockfile.

### 4.2. Failure mode the rule catches

**cph cdr-refactor wave, 2026-05-18 to 2026-05-20 (cph#27 + cph#28).** δ attempted parallel α dispatch for two cycles sharing the same on-disk working copy (`/root/cph`). α-28 attempted `git worktree add /tmp/cph-alpha-28` to isolate its working copy mid-session; the worktree-add setup did not complete cleanly under concurrent α activity and δ killed the first α-28 attempt. Serial dispatch (α-27 → β-27 merge → α-28) was the safe fallback; cph#28 still merged cleanly but ~2h of α work was serialized rather than parallel.

The race is at the harness/filesystem level — the role skills did nothing wrong. The fix lives on the δ surface: pre-create worktrees before dispatching.

### 4.3. The pre-creation contract

```bash
# Before dispatching α for cycle <N> (when parallel α targets the same repo):
git fetch --quiet origin cycle/<N>
git worktree add /path/to/<project>-cycle-<N> origin/cycle/<N>

# Configure per-worktree identity (per §3.2 worktree-aware path)
cd /path/to/<project>-cycle-<N>
git config --worktree user.name "alpha"
git config --worktree user.email "alpha@<project>.cdd.cnos"

# Dispatch α with the worktree path baked into the prompt
# (the α prompt's Branch: cycle/<N> line stays the same; α `cd`s to the worktree path
# named in the dispatch instructions and does NOT cd out of it)
```

The α prompt names the worktree path; α does not `cd` out of it nor re-init a working copy (`alpha/SKILL.md` cross-reference per cnos#384 AC2).

### 4.4. Cleanup contract

After α exits and β merges:

```bash
# Tear down the per-α worktree
git worktree remove /path/to/<project>-cycle-<N>
```

The cycle branch on `origin` remains until β's post-merge delete per `cnos.cds/skills/cds/CDS.md` §"Development lifecycle" → §"Step table" Step 8 (β review + merge) closure step — branch lifecycle is not coupled to worktree lifecycle.

### 4.5. Cross-references

- **`beta/SKILL.md` Pre-merge gate Row 3** ("Non-destructive merge-test") — β's merge-test worktree (`/tmp/cnos-merge-test/wt`) is the same primitive used for a different purpose (testing the merge in isolation). The β rule + this α rule are the same family: worktree-as-isolation, with worktree-local identity write per §3.2.
- **cnos#301 — β worktree-identity isolation discipline** — the first cycle that named per-worktree config writes. cnos#373 generalized; this section absorbs cnos#384 for the α parallel-dispatch case.

### 4.6. Worked example — when to pre-create

```bash
# Wave: cph subs #27 and #28, both targeting /root/cph (same repo)
# Sequential model: dispatch α-27 → wait for β-27 merge → dispatch α-28 (~2h serialized)
# Parallel model with pre-created worktrees:

# Pre-create worktrees BEFORE either α dispatch
cd /root/cph
git worktree add /tmp/cph-cycle-27 origin/cycle/27
git worktree add /tmp/cph-cycle-28 origin/cycle/28
cd /tmp/cph-cycle-27 && git config --worktree user.email "alpha@cph.cdd.cnos"
cd /tmp/cph-cycle-28 && git config --worktree user.email "alpha@cph.cdd.cnos"

# Dispatch α-27 and α-28 in parallel — each runs against its own working copy
# (α-27 in /tmp/cph-cycle-27, α-28 in /tmp/cph-cycle-28); .git is shared but
# .git/worktrees/{N} lockfiles isolate the per-worktree state.

# After both β merges complete:
git worktree remove /tmp/cph-cycle-27
git worktree remove /tmp/cph-cycle-28
```

- ❌ Dispatch two α sessions against the same `/root/cph` working copy concurrently.
- ❌ Let the α session run `git worktree add` after dispatch starts (the race already exists).
- ✅ δ pre-creates worktrees + per-worktree identity *before* α dispatch; cleans up after β merges.

---

## 5. Polling and wake-up

### 5.1. δ's wake-up signals (mechanics)

The δ-role policy (when to act / when to wait, what counts as a gate request, the "do not poll internal work" discipline) lives in `operator/SKILL.md` §2. This subsection documents the *mechanics* — the shell snippets and the wake-up contract — that the policy depends on.

**Mechanism:** transition-only polling under a `Monitor`-equivalent wrapper. Each stdout line from the loop becomes a wake-up event for the parent session.

### 5.2. Issue activity polling

```bash
prev=""
while true; do
  cur="$(cd /path/to/repo && git fetch --quiet origin && \
         gh issue view <N> --json comments --jq '.comments | length')"
  if [ "$cur" != "$prev" ]; then
    echo "issue-activity: comments=$cur"
    prev="$cur"
  fi
  sleep 300
done
```

Run under `Monitor` or equivalent. 5-minute interval is sufficient for δ — γ owns the tight loop on the cycle branch.

### 5.3. Cycle branch polling

Canonical cycle branches are `origin/cycle/{N}` (per `cnos.cds/skills/cds/CDS.md` §"Development lifecycle" → §"Branch rule", since #287). The pre-#287 `'origin/claude/*'` glob is **warn-only / retrospective** — retained for tracking historical cycles whose branches predate the rule, never as a discovery surface for new cycles.

```bash
prev_branches=""; declare -A prev_head
while true; do
  cd /path/to/repo && git fetch --quiet origin
  # Canonical: cycle/{N} branches (γ creates these per cnos.cds/skills/cds/CDS.md §"Development lifecycle" → §"Branch rule" / §"Step table" Step 2).
  cur_branches="$(git branch -r --list 'origin/cycle/*' 2>/dev/null | sed 's| ||g' | sort)"
  comm -13 <(echo "$prev_branches") <(echo "$cur_branches") | sed 's/^/new-branch: /'
  # Per-branch head SHA — cycle artifacts live on cycle branches, not on main.
  for b in $cur_branches; do
    cur_head="$(git rev-parse "$b" 2>/dev/null)"
    [ "$cur_head" != "${prev_head[$b]:-}" ] && [ -n "$cur_head" ] && echo "branch-update: $b → $cur_head"
    prev_head[$b]="$cur_head"
  done
  prev_branches="$cur_branches"
  sleep 300
done
# To track legacy branches retrospectively, swap the glob to 'origin/claude/*'
# (warn-only — pre-#287 cycles only). Do not use for new cycles.
```

### 5.4. Single-named-branch polling (γ-side, under `Monitor`)

γ owns the tight loop on a single named branch (per `gamma/SKILL.md` §2.5 dispatch). The same `Monitor`-wrapped transition-only loop applies, scoped to one branch:

```bash
# Baseline sync — run BEFORE the transition loop.
git fetch --quiet origin cycle/<N>
git rev-parse --verify origin/cycle/<N>
echo "baseline-cycle-dir: $(git ls-tree -r --name-only origin/cycle/<N> .cdd/unreleased/<N>/ 2>/dev/null | tr '\n' ' ')"

# Transition loop — single named branch.
prev_head=""; empty_iters=0
while true; do
  git fetch --quiet origin cycle/<N> || echo "fetch-error: cycle/<N>"
  cur_head="$(git rev-parse origin/cycle/<N> 2>/dev/null)"
  if [ "$cur_head" != "$prev_head" ] && [ -n "$cur_head" ]; then
    echo "branch-update: cycle/<N> → $cur_head"
    prev_head="$cur_head"
    empty_iters=0
  else
    empty_iters=$((empty_iters + 1))
    # cnos.cds/skills/cds/CDS.md §"Coordination surfaces" → §"Polling primitives" git fetch reliability rule: every 10 empty iterations, do a synchronous re-probe.
    if [ "$empty_iters" -ge 10 ]; then
      git fetch --verbose origin cycle/<N> 2>&1 | tee /tmp/cycle-<N>-fetch.log >&2 || \
        echo "reachability-fail: cycle/<N> — surface to operator"
      empty_iters=0
    fi
  fi
  sleep 60
done
```

Each transition line becomes a `task-notification` that wakes the parent session.

### 5.5. Reachability probe

The `git fetch --quiet` form silently swallows transport flake and auth issues. Every N empty iterations (canonical: 10), the loop re-probes synchronously with `git fetch --verbose` and surfaces any error. This prevents the polling-looks-healthy-but-truth-is-stale failure mode (cycle #283 β observation #3).

---

## 6. Timeout recovery

### 6.1. When this section fires

An agent dispatched via `claude -p` (or `cn dispatch`) may be SIGTERM'd by the OS, hit the session timeout, or crash without committing its in-progress work. This section is the executable recovery procedure. The δ-side override declaration (when the operator commits on behalf of the agent) lives in `operator/SKILL.md` §4; this section is the mechanics.

### 6.2. Inspect the worktree

```bash
# 1. What is staged or modified but not committed?
git status --short

# 2. What files were written during the session?
#    Replace /tmp/session-start-sentinel with the sentinel touched before dispatch.
find . -newer /tmp/session-start-sentinel -not -path './.git/*' | sort

# 3. Did the agent stash anything?
git stash list

# 4. What is the diff against main?
git diff origin/main..HEAD
git diff HEAD   # staged + unstaged against last commit
```

To create the sentinel, run `touch /tmp/session-start-sentinel` immediately before dispatch.

### 6.3. Decision tree

| Situation | Action |
|---|---|
| Work is committed — agent wrote one or more commits before SIGTERM | Normal recovery: the agent left durable checkpoints. Read `origin/cycle/{N}`. No further action unless the cycle is incomplete. |
| Work is staged / unstaged but not committed | Commit under agent identity per §3 (`git config user.name "Alpha" && git config user.email "alpha@cdd.cnos"`, with `--worktree` if §3.2 applies), then `git add <files> && git commit -m "<msg>"`. Preserve the agent's intent; do not rewrite scope. |
| Work exists only as new/modified files (never staged) | Same as above — stage and commit under agent identity with a message naming the recovery context (e.g., `"recovery(338): α partial work — committed by operator after SIGTERM"`). |
| Nothing useful recovered — session started but produced no artifact content | Declare a failed dispatch. File the failure in `self-coherence.md §Debt`. Re-dispatch α with a fresh budget per `cnos.cds/skills/cds/CDS.md` §"Field 6: Actor collapse rule" (heuristic constants in `operator/SKILL.md` §5.2). |
| Stash exists | `git stash show -p` to inspect; pop and commit if relevant: `git stash pop && git add <files> && git commit`. |

**Agent-identity vs operator-identity commit.** Prefer committing under the agent's canonical identity (`alpha@cdd.cnos`) to preserve the role-identity-is-git-observable property. If the agent's identity cannot be confirmed, commit under operator identity but declare this as an override in `self-coherence.md §Debt` with the override cross-reference to `operator/SKILL.md` §4.

### 6.4. Override declaration

If the operator commits work on behalf of an agent, this is an implicit override. The declaration shape and doctrine live in `operator/SKILL.md` §4. Mechanics: log the declaration in `self-coherence.md §Debt`:

> Override: operator-identity commit for cycle #N. Reason: α session SIGTERM'd before committing. Committed staged/unstaged work under operator identity at `<SHA>`. Grade implication: per `release/SKILL.md §3.8`, a cycle with operator override has γ < A.

### 6.5. Prevention

The recovery procedure is the failure path. The prevention path is a correctly-sized dispatch budget and commit-checkpoint instruction per `cnos.cds/skills/cds/CDS.md` §"Field 6: Actor collapse rule" (heuristic constants in `operator/SKILL.md` §5.2). If this section is being exercised, the dispatch that spawned the agent did not satisfy the budget contract — record the actual budget and AC count in the PRA telemetry fields (`post-release/SKILL.md` §4: `dispatch_seconds_budget`, `dispatch_seconds_actual`, `commit_count_at_termination`) so the heuristic can be refined.

### 6.6. Close-state assertion parity on recovery closure

If δ's recovery (§6.3) carries a SIGTERM-terminated cycle through merge, δ MUST run the same hard-gate assertion `gamma/SKILL.md §2.10` Phase A requires: `gh issue view {N} --json state --jq .state` MUST return `CLOSED` before the nonterminal post-merge marker is written. If not, δ closes the issue and records pre-state, action, and post-state in `gamma-closeout.md` (or `self-coherence.md §Debt` until that receipt exists). This assertion establishes issue custody only; recovery still requires δ disconnect, green CI, γ archive, and the later terminal declaration. *Derives from cnos#368's close-token gap.*

---

## 7. Branch retries and push restrictions

### 7.1. Harness push restrictions (`claude/*` branches)

Some Claude Code harness environments block updates to existing remote branches (HTTP 403 on push to a previously-written branch). Under this constraint, a fresh-branch chain (`cycle/{N}` → `cycle/{N}-impl` → `cycle/{N}-impl-r2` → `cycle/{N}-merged` → `cycle/{N}-final`) is an acceptable workaround. Each link is a valid cycle branch for that fix-round; the final fast-forward into `main` becomes an external δ action.

The doctrine surface for this constraint (when to invoke the chain, what it means for §5.2 single-session dispatch) lives in `operator/SKILL.md` §5.2 consequence (3). This subsection documents the mechanics:

```bash
# After harness 403 on push to cycle/<N>:
git checkout cycle/<N>
git checkout -b cycle/<N>-impl-r2
git push -u origin cycle/<N>-impl-r2

# Subsequent fix-rounds chain forward; never amend a pushed branch on the harness side.
```

### 7.2. Cycle-branch identity is the cycle branch on origin

Regardless of how many links the chain accumulates locally, the canonical coordination surface remains `origin/cycle/{N}` — γ polls it, β diffs against it, and δ gates on it. The chain is harness-side bookkeeping; the cycle-branch identity is the named coordination surface.

---

## 8. Embedded Kata

### Kata A — Observability gate

#### Scenario

δ is about to dispatch α via `claude -p` for cycle #500. The prompt is at `/tmp/alpha-prompt.md`. δ has the cycle's JSONL log path ready.

#### Task

Write the dispatch shell line.

#### Expected answer

```bash
touch /tmp/session-start-sentinel
cat /tmp/alpha-prompt.md | claude -p \
  --allowedTools "Read,Write,Bash" \
  --output-format stream-json --verbose \
  --permission-mode acceptEdits \
  --model <model> \
  2>&1 | tee /tmp/cycle-500.jsonl
```

#### Common failures

- Omitting `--output-format stream-json --verbose` (black-box dispatch).
- Omitting `--permission-mode acceptEdits` (α can't write files).
- Forgetting the sentinel (timeout recovery §6 can't identify session-newer files).

### Kata B — Worktree identity write

#### Scenario

You are β starting a session in the cnos repo. `git config --get extensions.worktreeConfig` returns `true`.

#### Task

Write the identity-write sequence.

#### Expected answer

```bash
# Confirmed: extensions.worktreeConfig=true
git config --worktree user.name "beta"
git config --worktree user.email "beta@cdd.cnos"
git config --get user.email   # expect: beta@cdd.cnos
```

#### Common failures

- `git config user.email "beta@cdd.cnos"` (no `--worktree`) — writes to the shared layer; the next sibling worktree's identity write overwrites it; β's next commit lands under the wrong author.
- Checking `--get user.email` and trusting the immediate read.

### Kata C — Parallel α worktree pre-creation

#### Scenario

δ has wave-mode dispatch ready for cycles #N1 and #N2 against the same repo at `/root/<project>`. Both cycle branches exist on origin.

#### Task

Pre-create the per-α worktrees with correct identity, then state the dispatch sequence.

#### Expected answer

```bash
cd /root/<project>
git fetch --quiet origin cycle/<N1> cycle/<N2>

git worktree add /tmp/<project>-cycle-<N1> origin/cycle/<N1>
git worktree add /tmp/<project>-cycle-<N2> origin/cycle/<N2>

cd /tmp/<project>-cycle-<N1>
git config --worktree user.name "alpha"
git config --worktree user.email "alpha@<project>.cdd.cnos"

cd /tmp/<project>-cycle-<N2>
git config --worktree user.name "alpha"
git config --worktree user.email "alpha@<project>.cdd.cnos"

# Dispatch α-<N1> and α-<N2> in parallel; each prompt names its worktree path.
# After β merges complete:
cd /root/<project>
git worktree remove /tmp/<project>-cycle-<N1>
git worktree remove /tmp/<project>-cycle-<N2>
```

#### Common failures

- Dispatching α without pre-created worktrees (race on shared `.git`).
- Configuring identity without `--worktree` (per §3.2 leak).
- Forgetting the cleanup (worktree directory accumulates over time).
