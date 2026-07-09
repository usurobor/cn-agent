# γ Scaffold — cnos#640

**Issue:** [#640](https://github.com/usurobor/cnos/issues/640) — dispatch-protocol: design-first hold state duplicated across issue-body prose and the `status:*` label, and the two can drift out of sync.
**Mode:** design-and-build (process doctrine + a likely-small mechanism).
**cell_kind:** `doctrine` (see Friction note 5 — α may redeclare if the shipped diff turns out code-dominant).
**Branch:** `cycle/640` (already created from `origin/main`; see Friction note 6 for the base-SHA note).
**Family:** #614 (first occurrence + RCA), #633 (live recurrence, same session), #630 (mandatory-learning doctrine), #583/#584 (mechanism-over-prose), #501 (κ role), #639 (adjacent — **do NOT dispatch**).
**Authority already granted on the issue:** operator posted "✅ Operator-authorized dispatch... A2/CAP: use MCA for body cleanup, PR metadata, tests, and CI fixes without asking. Escalate only if the fix requires new status semantics, weakens operator hold authority, or changes the dispatch/FSM lifecycle." α does not need to re-ask for routine implementation authority within that envelope.

---

## 1. Per-AC oracle list

### AC1 — dispatching a design-first issue can no longer leave body-hold and label state contradictory

**Invariant:** after this cycle ships, no mechanism exists (or survives) by which an issue's body can assert "held" while its `status:*` label asserts "dispatch now" without either (a) that state being structurally impossible, or (b) it being caught and named before it costs a firing.

**Oracle:**
- A fixture reconstructing the exact #614/#633 shape — issue body ending in the hold phrase (`"Not dispatched — status:ready. ... dispatch on explicit operator authorization."` or equivalent) plus a `status:ready`/`status:todo` label — run through whatever mechanism α ships, with the **positive case** (contradictory fixture) producing a corrected/atomic result and the **negative case** (already-clean issue, no hold phrase in body) producing a safe no-op — both provable by an automated test (unit test with an injectable GitHub-API double, mirroring `issues-fsm`'s `TerminalOptions`/`ScanOptions` pattern; NOT a live-network manual check).
- If a `cn issues dispatch --issue N` (or equivalent) primitive is built: idempotency oracle — running it twice against the same issue produces the same end state both times, and its own test suite proves the label flip and the body edit happen as one caller-visible operation (no half-applied state observable between the two, even if the underlying GitHub calls are still two HTTP requests).

### AC2 — the wake no longer has to decide "body says held, label says go"

**Invariant:** the ambiguity δ/κ hit at #614 and #633 (reading the body, finding a contradiction, deferring) is either structurally absent (the body no longer carries hold-state prose that could contradict the label) or mechanically flagged (a gate names the contradiction before a wake wastes a firing on it) — not resolved by an agent's in-the-moment judgment call.

**Oracle:**
- Read the diff to `src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md` §1.2 "Human dispatch gate" and/or `src/packages/cnos.core/skills/agent/label-doctrine/SKILL.md` and confirm it now states, in one place, which surface is authoritative for dispatch-readiness (labels) and what a body-hold sentence means if one is still found (non-authoritative / to be removed).
- If a contradiction gate is built: a fixture test proving it fires exactly on the contradictory combination (body hold phrase present + `status:todo`/`status:in-progress`) and not on any clean combination (checked against ≥2 negative fixtures: no hold phrase; hold phrase + `status:ready` which is not a contradiction).
- Cross-check `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` §9.2 input #1 ("δ does NOT improvise additions" to the body's named fields) still reads consistently after the change — δ's wake-invoked contract must not retain language that tells δ to treat a stray body sentence as authority-bearing once labels are pinned as sole source of truth.

### AC3 — the observe → capture → mechanize → verify loop is documented

**Invariant:** a doctrine section exists naming the 4-step loop (observe / capture / detect-recurrence / mechanize / verify-non-recurrence — κ's comment names 5 steps, folding "detect-recurrence" in as step 3; see below), who owns each step **today**, and where the current gap is — honestly, including the gap κ's comment names (nobody/nothing runs cross-cell recurrence detection today).

**Oracle:**
- `grep` finds a new, named section (landing candidate: `src/packages/cnos.cdd/skills/cdd/CELL-KINDS.md` extending §"Mandatory terminal learning section", cross-linked from `CDD.md` — see source-of-truth table) that: (1) names all of observe / capture / detect-recurrence / mechanize / verify-non-recurrence explicitly; (2) states today's owner per step — capture = γ (already doctrine, #630), detect-recurrence = **named as an open gap** with a concrete next-MCA (an issue number, or explicit "ε made concrete" commitment) rather than silently dropped, mechanize = promotion-to-`kind/process`-issue-with-an-MCA-bar, verify = "the mechanism ships with a test/gate that would fail on recurrence"; (3) cites #614 → #633 as the worked example and links κ's #640 meta-comment as the seed, per the issue's own instruction ("the meta comment is the seed; land it in doctrine").
- Negative oracle: the section does **not** claim the ownership gap is closed if it isn't — an AC3 doc that asserts "ε now runs this" without a shipped ε pass is a false-closure violation (`issue/SKILL.md` failure mode). Documenting the gap honestly, with a named next step, satisfies AC3; building ε itself is out of scope for this cell's size (see Scope guardrails).

### AC4 — all gates green; no weakening of the operator's genuine ability to hold work design-first

**Oracle:**
- CI green on the merge commit for every gate the touched surfaces trigger (Go build/test if a CLI primitive ships; Package; the doctrine-file lint/link-check gates if any apply to `CELL-KINDS.md`/`dispatch-protocol/SKILL.md`/`label-doctrine/SKILL.md`).
- Negative-case proof: `label-doctrine/SKILL.md` §1.2's invariant — "The `status:ready → status:todo` transition is the operator's authorization event. It is not automated." — is unchanged in substance after the diff. If a `cn issues dispatch` primitive ships, it is an **operator/κ-invoked** command (a human or an explicitly-authorized planner runs it naming the issue), never something a dispatch wake calls on its own initiative to self-authorize a claim. β greps the diff for any wake/scan-side call site of the new primitive and confirms none exists.
- `status:ready` remains a valid, honored hold state — an issue at `status:ready` with no `status:todo` is still not claimable by any wake (unchanged FSM/claim-selector behavior, `src/packages/cnos.cds/skills/cds/fsm/transitions.json` byte-identical in its state-machine semantics, only doc/prose-adjacent surfaces touched).

---

## 2. Source-of-truth table

| Claim / surface | Canonical source | Status / role for α |
|---|---|---|
| The "Not dispatched" body-hold prose convention | **Not defined in any skill or template** — confirmed via `rg -rn "Not dispatched"` and `rg -rn "dispatch on explicit operator"` across the whole repo (0 hits outside issue-body artifacts). The only textual record of the exact phrase is `.cdd/unreleased/614/beta-review.md` (quotes it verbatim while RCA'ing the first occurrence) and #640's own issue body. | Do not go looking for a canonical doc to "fix" — there isn't one. The convention was ad hoc prose γ/κ wrote by hand at filing time. The fix is either "stop writing it" (doctrine) or "reconcile it mechanically" (primitive/gate) — not an edit to a pre-existing definition. |
| `status:ready` / `status:todo` semantics, the human-authorization gate | `src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md` §1.2 "Human dispatch gate" (~lines 121–136) | Primary landing spot for AC1/AC2's "labels are the sole source of truth for dispatch readiness" doctrine statement. **Do not weaken** the "not automated" operator-authorization sentence there (AC4). |
| Generic lifecycle-label table + failure-mode catalogue | `src/packages/cnos.core/skills/agent/label-doctrine/SKILL.md` §1.1 (lifecycle table) — secondary/complementary doctrine home; and `src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md` §1.3 (named failure modes D1–D7, D10 — **D8 and D9 are currently unused/unassigned**, a natural free slot for a new "body/label authorization drift" failure-mode entry, structurally the same shape as the existing D-entries: symptom + `Fix: §ref`). | Candidate slot for a new D-numbered failure mode documenting #614/#633 and pointing at whatever mechanism α ships. Not mandatory to use this exact slot, but a clean fit if α wants a structured failure-mode entry rather than only prose in §1.2. |
| Issue-body template (what γ/κ fill in when filing a design-first issue) | `.github/ISSUE_TEMPLATE/cdd-issue.md` | Currently silent on the body-hold convention (only says `dispatchable: yes → dispatch:cell + protocol:cds/cdd + status:ready|todo`). No change required by the ACs, but if α's doctrine fix explicitly forbids the hold-prose pattern, a one-line comment here reinforcing "labels alone gate dispatch; do not write a body hold sentence" is a cheap, high-leverage addition — optional, not required by any AC. |
| δ's wake-invoked input contract (how δ reads the issue body as spec) | `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` §9.2, input #1 (~lines 401) | Read for consistency (AC2 oracle). δ's contract already says "δ does NOT improvise additions" beyond the named Mode/Gap/Impact/Source-of-truth/Constraints/ACs/Proof/Cross-references fields — confirm nothing here tells δ to treat a stray body sentence as dispatch-authority-bearing once labels are pinned as sole source of truth. This file is **framework-level** (cnos.cdd); if it needs a change, keep it a small clarifying addition, not a rewrite — γ-level lifecycle-transition semantics are explicitly out of scope. |
| FSM transition table for the `ready` state | `src/packages/cnos.cds/skills/cds/fsm/transitions.json` (state `"ready"` → single rule, `action: propose_status_todo`, reason: `"status:ready has no Phase-1 gating guard; release for dispatch (operator flips ready -> todo)"`) | **Read-only reference.** Confirms the FSM already treats the ready→todo flip as the sole state-machine mechanism — this file's semantics are unaffected by the body/label reconciliation fix. **Do not edit this file's transition rules or add/remove states/guards** (explicit out-of-scope: "do NOT ... change FSM lifecycle semantics"). |
| `cn issues fsm` / `cn issues map` — sibling CLI-command precedent | Dispatch wiring: `src/go/internal/cli/cmd_issues_fsm.go`, `src/go/internal/cli/cmd_issues_map.go`. Domain logic: `src/packages/cnos.issues/commands/issues-fsm/*.go` (esp. `fetch.go`'s `ghRequest`/`ghAddLabel`/`ghRemoveLabel`/`ghEnsureLabelExists` HTTP primitives and `terminal.go`'s injectable-dependency `*Options` struct pattern — `ListClosedCandidates`/`RemoveLabel`/`AddLabel` are function fields with live defaults, exactly the shape a testable `cn issues dispatch` would want). Registration: `src/go/cmd/cn/main.go` lines 48–49 (`reg.Register(&cli.IssuesMapCmd{})` / `reg.Register(&cli.IssuesFsmCmd{})`). CLI routing: `src/go/internal/cli/dispatch.go` lines 47–51 — noun-verb resolution joins `args[0]+"-"+args[1]` and looks it up in the registry, so `cn issues dispatch` resolves to a command registered under the name `"issues-dispatch"`, the exact sibling shape of `"issues-fsm"`/`"issues-map"`. | If the atomic-primitive path is taken: new files at `src/go/internal/cli/cmd_issues_dispatch.go` (thin dispatch wiring, `CommandSpec{Name: "issues-dispatch", ...}`) + `src/packages/cnos.issues/commands/issues-dispatch/` (domain logic, own `go.mod` mirroring `issues-fsm`'s), registered in `main.go` beside the two existing siblings. Reuse `ghRequest`/label primitives verbatim where possible. Body-editing has **no existing precedent in this codebase** — `ghEditIssueBody` (PATCH `/repos/{repo}/issues/{issue}` with `{"body": "..."}`) would be the first of its kind; model it on `ghAddLabel`'s shape (marshal payload, `ghRequest`, status-code check). |
| Cell-kind observation seam | `src/packages/cnos.cdd/skills/cdd/CELL-KINDS.md` §"Observation, not enforcement" + `issues-fsm/fetch.go`'s `assembleLive`, which parses a `cell_kind:` line out of `.cdd/unreleased/{N}/gamma-scaffold.md` when present (`FactSnapshot.CellKind`). Observation-only; no transition guard consumes it. | Informational. This scaffold declares `cell_kind: doctrine` in its header per Friction note 5. |
| Landing candidate for AC3's loop doctrine | `src/packages/cnos.cdd/skills/cdd/CELL-KINDS.md` §"Mandatory terminal learning section" (already frames capture as "ε captured by γ"; already links `gamma/SKILL.md`'s γ-binds-learning clause) | Natural extension point — add a named subsection for the observe→capture→detect-recurrence→mechanize→verify loop, cross-linked from `CDD.md` and from this section's own §"Cross-references" list (do not leave the new section orphaned — `issue/SKILL.md`'s "source drift" / discoverability discipline applies to doctrine too). |
| Kernel doctrine the issue cites directly | `src/packages/cnos.core/doctrine/KERNEL.md` §1.2 "Action" ("Change the system when you can (MCA)... 'Won't repeat' without a mechanism is not a fix.") and §2.2 "Engineering invariants" ("One source of truth per fact") | Cite directly in AC3's doc and in the AC1/AC2 doctrine addition; do not re-derive or paraphrase — these are the exact clauses #640's body names as governing. |
| AC3's literal seed content | Issue #640's own comment thread — κ's "Meta: how does process self-improvement actually work here? (κ, 2026-07-09)" comment (`gh issue view 640 --json comments`) | The issue explicitly calls this comment "the seed" for AC3. Read it via `gh issue view 640 --json comments -q '.comments[].body'` and use its proposed 5-step loop and its honest "who owns this" analysis as the basis — do not re-derive from scratch. |
| A second live instance of the same failure family (reference only) | Issue #618's body: `"Filing-only — operator holds dispatch."` while it happened to carry `status:todo` at #640's own claim time — δ's claim comment on #640 explicitly names this and states it was deferred, not claimed. | **Reference only.** Illustrates the pattern generalizes beyond the literal "Not dispatched — status:ready" phrasing. #618 itself is untouched — not named in the issue's out-of-scope list explicitly, but γ treats "do not touch other named issues" as covering it; see Scope guardrails. |

---

## 3. α prompt

```text
You are α. Project: cnos.

Load src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md and follow its load order.
Issue: gh issue view 640 --repo usurobor/cnos --json title,body,state,comments
Branch: cycle/640
Tier 3 skills:
  - src/packages/cnos.core/skills/write/SKILL.md
  - src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md
  - src/packages/cnos.core/skills/agent/label-doctrine/SKILL.md
  - src/packages/cnos.cdd/skills/cdd/CELL-KINDS.md
  - src/packages/cnos.core/doctrine/KERNEL.md (§1.2, §2.2 — cite, do not restate)

Read .cdd/unreleased/640/gamma-scaffold.md in full FIRST — it names the
per-AC oracle list, the concrete source-of-truth surfaces (including which
files do NOT already define the "Not dispatched" convention — do not go
hunting for a canonical doc that doesn't exist), and pre-identified friction.

## What you are fixing

A design-first issue's body can say "Not dispatched — status:ready ...
dispatch on explicit operator authorization" while its status:* label
independently says status:todo. These two can drift (#614, then #633 in
the same session). The wake/δ then has to defer on the contradiction,
costing firings. Root cause: the hold state lives in two places (Kernel
§2.2 "one source of truth per fact" violation), and a captured learning
(#614's RCA) was never mechanized (Kernel §1.2 "MCA over MCI").

## Build to the operator's stated preference

The operator's own issue body states a lean toward candidates (1)+(2) of
the four named in the issue: remove the duplicate hold-state prose from
the body (labels become the sole source of truth for dispatch readiness),
AND ship a one-verb dispatch primitive (documented command, e.g. `cn
issues dispatch --issue N`) so a human/κ never has to hand-edit both
surfaces again. Build both. Candidate (3) — a mechanical contradiction
gate — is optional; add it if it is cheap given what you've already built
(e.g. a fixture-testable check reusing the same detection logic your
primitive needs internally to decide whether a body edit is required).
Candidate (4) — a wake precedence rule that lets the label win WITHOUT
also cleaning the body — is disfavored by the operator's own risk framing
in the issue ("riskier... likely rejected"); only reach for it, with an
explicit justification in self-coherence.md, if (1)+(2) turn out
insufficient to satisfy AC1/AC2.

## Implementation contract (pinned by γ; α MUST NOT improvise)

| Axis | Pinned value |
|---|---|
| Language | Go (the CLI primitive, if built) + Markdown (the doctrine changes). No other language. |
| CLI integration target | `cn` subcommand — `cn issues dispatch`, a sibling of the existing `cn issues fsm` / `cn issues map` (same kernel binary, same registry, same noun-verb routing). Not a standalone binary. |
| Package scoping | `src/go/internal/cli/cmd_issues_dispatch.go` (thin dispatch wiring only, per INVARIANTS.md T-002 / eng/go §2.18) + `src/packages/cnos.issues/commands/issues-dispatch/` (domain logic: label mutation + body edit), registered in `src/go/cmd/cn/main.go` beside `IssuesMapCmd`/`IssuesFsmCmd`. Mirror `issues-fsm`'s structure (own `go.mod`, injectable-dependency `*Options` pattern for testability). |
| Existing-binary disposition | N/A — additive new subcommand inside the existing `cn` kernel binary. Nothing is replaced or deprecated. |
| Runtime dependencies | GitHub REST API only, via the same HTTP primitives shape as `issues-fsm/fetch.go` (`ghRequest`, `ghAddLabel`/`ghRemoveLabel` reused verbatim where possible; a new `ghEditIssueBody` PATCH primitive is the one genuinely new HTTP call — no existing precedent in this repo, model it on `ghAddLabel`'s shape). No other new runtime dependency. |
| JSON/wire contract preservation | N/A / additive. This is a new command; no existing wire contract to preserve. If you add a `--json` output flag, define its shape fresh — no migration needed. |
| Backward-compat invariant | `status:*` label semantics, `src/packages/cnos.cds/skills/cds/fsm/transitions.json`'s transition rules, and the existing claim/dispatch-protocol mechanics are unchanged. The fix is additive (new doctrine sentence + new optional command); no existing dispatchable issue's behavior changes unless it happens to carry the legacy body-hold phrase. |

## Scope guardrails (see gamma-scaffold.md §5 for the full list — restated here for emphasis)

Do NOT: weaken the operator's ability to hold work design-first · add new
status labels unless absolutely necessary · change FSM lifecycle
semantics (transitions.json's rules stay byte-identical) · change #630's
recovery behavior · touch #626 · touch #618 (a second live instance of
this same pattern-family, named only in this scaffold, not in the issue
body — leave it alone) · start Demo 0 · dispatch #639.

## AC3 — the loop doctrine

Read κ's meta-comment on #640 (`gh issue view 640 --repo usurobor/cnos
--json comments`) — it is explicitly named as "the seed." Land the
observe → capture → detect-recurrence → mechanize → verify-non-recurrence
loop in doctrine (landing candidate: CELL-KINDS.md, extending
§"Mandatory terminal learning section"; cross-link from CDD.md). Name
today's owner per step honestly — capture is γ (existing, #630);
detect-recurrence has NO current owner (κ's comment says so explicitly:
"nobody/nothing runs it") — do not paper over that gap; name it and
commit a concrete next step (an issue number, or an explicit "next: ε
made concrete" commitment) rather than asserting it's solved. Building a
full ε pass is out of scope for this cell's size — document the gap, do
not build the pass.

## Before requesting review

Write .cdd/unreleased/640/self-coherence.md per alpha/SKILL.md, mapping
each AC to concrete evidence (diff hunks, test names, doctrine section
names/paths). State your final cell_kind explicitly (γ pinned `doctrine`
tentatively; reclassify with reasoning if the shipped diff is
code-dominant). Name any open design call honestly (e.g. whether you used
the CELL-KINDS.md landing spot for AC3 or a different one, and why).
```

---

## 4. β prompt

```text
You are β. Project: cnos.

Load src/packages/cnos.cdd/skills/cdd/beta/SKILL.md and follow its load order.
Issue: gh issue view 640 --repo usurobor/cnos --json title,body,state,comments
Branch: cycle/640

Read .cdd/unreleased/640/gamma-scaffold.md §1 (per-AC oracle list) and §2
(source-of-truth table) FIRST — walk α's diff against that oracle list
independently. Do not trust α's self-coherence.md framing; re-derive.

## AC-by-AC verification

- AC1: locate the fixture test(s) proving the #614/#633 shape (body-hold
  phrase + status:ready/status:todo) is corrected atomically, AND that a
  clean issue (no hold phrase) is a safe no-op. If no automated fixture
  test exists — only a described manual check or prose claim — this is a
  finding (the issue/SKILL.md "false closure" failure mode: proof without
  a proof mechanism).
- AC2: confirm dispatch-protocol/SKILL.md §1.2 and/or label-doctrine/SKILL.md
  now state, in one place, that labels alone gate dispatch readiness.
  Cross-check delta/SKILL.md §9.2 input #1 for consistency — no
  contradicting clause should remain. If a contradiction gate was built,
  verify its fixture tests (≥1 positive, ≥2 negative per gamma-scaffold.md
  §1 AC2 oracle).
- AC3: confirm the loop doctrine names all 5 steps (observe / capture /
  detect-recurrence / mechanize / verify-non-recurrence), names today's
  owner per step, and is HONEST about the detect-recurrence gap (does not
  claim it's solved if it isn't — cross-check against κ's meta-comment on
  #640). Confirm the new section is cross-linked, not orphaned (check
  CDD.md and CELL-KINDS.md's own cross-reference list for a matching
  entry).
- AC4: confirm CI is green on the branch. Grep the diff for any call site
  where a dispatch wake or scan/reconciler process invokes the new
  primitive (if built) on its own initiative — none should exist; the
  primitive is operator/κ-invoked only. Confirm label-doctrine/SKILL.md
  §1.2's "not automated" operator-authorization sentence for status:ready
  → status:todo is unchanged in substance. Confirm
  src/packages/cnos.cds/skills/cds/fsm/transitions.json is byte-identical
  to its pre-cycle state (no diff at all) — this is a hard check, not a
  judgment call.

## Standard β obligations

- Implementation-contract coherence (Rule 7): verify the diff conforms to
  every pinned axis in gamma-scaffold.md §3's Implementation contract
  table (Go + Markdown only; `cn issues dispatch` as a `cn` subcommand if
  built, not a standalone binary; package paths as pinned).
- Non-goal / out-of-scope check: confirm no diff touches #618, #626,
  #630's recovery-behavior code, adds a new status label, or changes
  transitions.json — per gamma-scaffold.md §5.
- gamma-scaffold.md presence gate (rule 3.11b): already satisfied — this
  file exists on cycle/640 before α's dispatch.
```

---

## 5. Scope guardrails

Restated verbatim from the issue body's "Out of scope / do NOT" list, plus one γ-added item:

- **Do NOT** weaken the operator's genuine ability to hold work design-first.
- **Do NOT** add new status labels unless absolutely necessary.
- **Do NOT** change FSM lifecycle semantics (`src/packages/cnos.cds/skills/cds/fsm/transitions.json`'s transition rules/states/guards stay unchanged).
- **Do NOT** change #630's recovery behavior.
- **Do NOT** touch #626.
- **Do NOT** start Demo 0.
- **Do NOT** dispatch #639.
- **(γ-added, not literal issue text)** Do NOT touch #618 — a second live instance of the same body/label-contradiction pattern, surfaced only via δ's own #640 claim comment (which deferred on it rather than claiming it). It is not named in #640's own out-of-scope list, but "do not touch other named issues" is read to cover it; #618 gets its own cell if/when the operator authorizes it.
- **Escalation boundary** (from the operator's dispatch-authorization comment): use MCA freely for body cleanup, PR metadata, tests, and CI fixes; escalate to operator only if the fix would require new status semantics, weaken operator hold authority, or change dispatch/FSM lifecycle semantics.

---

## 6. Friction notes

1. **No canonical file defines the "Not dispatched" body-hold convention.** A repo-wide `rg -rn "Not dispatched"` / `rg -rn "dispatch on explicit operator"` returns zero hits outside issue-body artifacts (`.cdd/unreleased/614/beta-review.md` quoting it while RCA'ing #614, and #640's own body). This means AC1/AC2 are not "fix an existing doc" — they're "write the doctrine that should have existed." γ confirmed this via grep per gamma/SKILL.md §2.2a's peer-enumeration discipline (do not assert "X does not exist" without grep evidence).
2. **`dispatch-protocol/SKILL.md` §1.3's D-numbered failure-mode catalogue has a gap at D8/D9** (defined entries: D1, D2, D3, D4, D5, D6, D7, D10 — D8/D9 unused). γ flags this as a clean structural slot for a new failure-mode entry but does not mandate α use it — a prose addition to §1.2 "Human dispatch gate" may be sufficient depending on what α actually ships.
3. **AC3's "meta comment"** is κ's issue-#640 comment titled "Meta: how does process self-improvement actually work here? (κ, 2026-07-09)". It proposes a 5-step loop (observe / capture / detect-recurrence / mechanize / verify-non-recurrence) and is explicit that step 3 (detect-recurrence) has **no current owner** — ε is named only as a "role name, not a running process." γ interprets AC3 as satisfied by documenting this honestly (loop + per-step owner + an explicit named gap-and-next-step for detect-recurrence), not by building a working ε pass this cycle — building ε would blow this cell well past "process/doctrine cell with a possibly-small mechanism" sizing per the operator's explicit "clean first dispatch... no Demo 0" framing.
4. **Candidate-mechanism weighting.** The issue names 4 candidate fixes and says the mechanism is "design sub picks; not pre-decided," but the operator's own "Preferred design" section leans explicitly toward (1) atomic primitive + (2) single source of truth. γ's α prompt treats (1)+(2) as the target and (3) contradiction gate as optional/cheap-if-convenient, and flags (4) wake-precedence-without-body-cleanup as disfavored (the issue's own text calls it "riskier... likely rejected") rather than hard-banning it — α retains room to justify (4) explicitly if (1)+(2) prove insufficient, per design-and-build mode's "design converges within this cycle" shape.
5. **cell_kind ambiguity.** The issue itself hedges: "Cell kind: `doctrine` → possibly `implementation` (a guard/helper)." γ pins `doctrine` in this scaffold's header (the dominant matter is a label/dispatch-doctrine fix; the CLI primitive, if built, is a small helper serving that doctrine, not the other way around) but explicitly authorizes α to reclassify in `self-coherence.md` with stated reasoning if the shipped diff turns out code-dominant. This is a non-blocking, low-stakes call.
6. **Branch base-SHA note (informational, not a blocker).** δ's claim comment on #640 recorded head commit `2114d64979f669f191d94ef31a53e65f9dfd7102` as the base for `cycle/640`. `origin/main` has since advanced one commit to `59e7fdb74fd7ea2abef863cd36378f080383e330` (a `board-map: regenerate docs/development/board` commit — no code/doctrine-surface overlap with this cell). `cycle/640`'s current tip (`f5c1204...`, the `CLAIM-REQUEST.yml` commit) branches from `59e7fdb`, i.e. from the newer SHA, not the older one δ's comment cited. This is a harmless one-commit drift; β should still SHA-pin-verify per its own pre-merge gate, but γ does not expect it to surface a real conflict.
7. **`.cdd/unreleased/640/CLAIM-REQUEST.yml` already exists on `cycle/640`** (written by δ's claim sequence per `dispatch-protocol/SKILL.md` §2.2 step 4 / `cnos#575` Phase 3). γ leaves it untouched — it is δ's artifact, not γ's, and no rule in `gamma/SKILL.md` or `delta/SKILL.md` §9.5 asks γ to remove or modify it at scaffold time.
