---
name: dispatch
description: Dispatch-time wire-format contract between γ (orchestrator) and α/β (workers). Use when γ authors a dispatch prompt for α or β, when δ enriches an under-specified prompt at routing time, or when α reads the implementation-contract block at intake.
artifact_class: skill
kata_surface: embedded
governing_question: How does γ hand a cycle to α/β without leakage — naming the issue, the branch, the Tier 3 skills, and the 7-axis implementation contract — and what δ does when the contract is under-pinned at routing time?
visibility: public
parent: handoff
triggers:
  - dispatch
  - dispatch-prompt
  - implementation-contract
  - inward-membrane
  - prompt template
  - prompt rules
  - alpha prompt
  - beta prompt
  - gamma prompt
  - 7 axes
scope: task-local
inputs:
  - γ-authored dispatch prompt at scaffold time
  - cycle issue body
  - repo conventions (language, package paths, integration target)
  - the active sub-issue's protocol exemption (if any)
outputs:
  - α-ready / β-ready dispatch prompt with populated implementation-contract section
  - δ-enriched dispatch prompt (when γ leaves a row unpopulated)
  - mid-flight `gamma-clarification.md` entry (when a contract axis re-pins mid-cycle)
  - escalation to operator-as-human (when a contract row is genuinely undecidable)
requires:
  - active CDD-class cycle exists with a quality-gated issue, a cycle branch on origin, and `gamma-scaffold.md` filed
calls:
  - HANDOFF.md
  - cnos.cdd/skills/cdd/gamma/SKILL.md
  - cnos.cdd/skills/cdd/alpha/SKILL.md
  - cnos.cdd/skills/cdd/beta/SKILL.md
  - cnos.cdd/skills/cdd/delta/SKILL.md
  - cnos.cdd/skills/cdd/operator/SKILL.md
---

# Dispatch — γ → α/β prompt + δ implementation-contract enrichment

## Core principle

**A dispatch prompt is a wire-format contract between γ (orchestrator) and α/β (workers). It carries exactly the artifact-level facts the worker needs to start (issue, branch, Tier 3 skills, implementation contract) and excludes the protocol-internal reasoning γ used to assemble them. δ stands at the inward face of the membrane and enriches any under-pinned row before α is routed; α MUST NOT improvise.**

The dispatch prompt is one of the wire formats `cnos.handoff` owns. It transports a unit of work from γ-the-orchestrator (or operator-as-γ) to a per-role worker session (α, β, or a fresh γ session) without the worker having to reconstruct what γ already decided. The prompt's envelope carries the issue reference, the cycle branch, the Tier 3 skill list, and the 7-axis implementation contract; the prompt's body never restates the algorithm and never smuggles missing constraints into chat prose.

The failure mode is **implementation-contract drift** — γ leaves an architectural axis empty or "TBD"; δ does not catch the omission at routing time; α improvises; β APPROVE-s on behavior-only AC oracles without checking the diff against the contract. The result is a cycle that satisfies its ACs but ships in the wrong language, at the wrong package path, or as a separate binary instead of a `cn` subcommand. cnos#389 (Python-not-Go) and cnos#391 (wrong package scoping + separate binary) are the canonical anchors for this failure mode.

This skill is the **canonical surface** for the dispatch wire-format: prompt template, prompt rules, the 7-axis implementation-contract schema, the δ-inward-membrane enrichment doctrine, and the four-surface mesh (γ template ↔ δ enrichment ↔ α constraint ↔ β verification) that distributes the contract across the cycle's role boundary. Consumer role-skills in `cnos.cdd` (and any future c-d-X protocol package) cite this skill as canonical and carry only role-local procedure.

## Authority

- This skill (in `cnos.handoff`) is the only canonical home for the dispatch-prompt template, the 7-axis implementation-contract schema, the δ-as-inward-membrane doctrine, the four-surface mesh declaration, and the cnos#389/#391/#392/#393 empirical anchors that motivate them. If another skill or doctrine surface carries the same doctrine, this skill governs.
- The consumer protocol packages govern role-local procedure into which the dispatch contract feeds. γ's coordination loop (issue selection, branch pre-flight, scaffold gate, polling, unblock without leakage) lives in `cnos.cdd/skills/cdd/gamma/SKILL.md`. α's intake and obligation to surface unpinned rows lives in `cnos.cdd/skills/cdd/alpha/SKILL.md §3.6`. β's diff-vs-contract verification lives in `cnos.cdd/skills/cdd/beta/SKILL.md §Role Rules Rule 7`. δ's outward membrane (boundary decision on receipt + V verdict) and override authority live in `cnos.cdd/skills/cdd/delta/SKILL.md` (§1, §3).
- The `## Implementation contract (pinned by δ; α MUST NOT improvise)` block γ writes into the α prompt is the wire-format rendering this skill governs; γ's role-skill cites this skill for the schema and the injection rule.

## Scope

In scope:

- Dispatch-prompt template format for γ, α, and β prompts (the markdown envelope δ routes via `cn dispatch` or equivalent identity-rotation primitive).
- Prompt rules (issue reference shape; explicit `Branch:` line; `--json title,body,state,comments` convention; do-not-restate-the-algorithm rule; do-not-smuggle-constraints rule).
- The 7-axis implementation-contract schema (Language; CLI integration target; Package scoping; Existing-binary disposition; Runtime dependencies; JSON/wire contract preservation; Backward-compat invariant) and the `## Implementation contract` markdown template γ injects into the α prompt.
- δ as inward membrane: the review-before-routing duty, the fill-or-escalate paths, the cycle-channel logging of enrichments.
- The four-surface mesh (γ template / δ enrichment / α constraint / β verification) declaration.
- Mid-flight contract re-pinning via `gamma-clarification.md` (the cycle-channel hand-off into the mid-flight rescue mechanism whose canonical home is [`cnos.handoff/skills/handoff/mid-flight/SKILL.md`](../mid-flight/SKILL.md), landed under Sub 4 of cnos#404 / cnos#418).

Out of scope:

- γ-owned branch pre-flight mechanics (creating `cycle/{N}` from `origin/main`; `gamma-scaffold.md` presence gate; polling). These are cycle-lifecycle, not wire-format; canonical at `cnos.cdd/skills/cdd/gamma/SKILL.md §2.5` + `cnos.cds/skills/cds/CDS.md §"Development lifecycle"`.
- The mid-flight `gamma-clarification.md` mechanism itself (the cache-bust semantics; the in-flight α/β polling protocol). Cited from this skill; canonical home at [`cnos.handoff/skills/handoff/mid-flight/SKILL.md`](../mid-flight/SKILL.md) (Sub 4 of cnos#404 / cnos#418).
- The `.cdd/unreleased/{N}/` artifact-channel rules (per-role write ownership; frozen-snapshot rule on merge). Cited from this skill; canonical home at [`cnos.handoff/skills/handoff/artifact-channel/SKILL.md`](../artifact-channel/SKILL.md) (Sub 4 of cnos#404 / cnos#418).
- The harness substrate (dispatch invocation mechanics; observability flags; identity rotation; timeout recovery). Lives in `cnos.cdd/skills/cdd/harness/SKILL.md` per Phase 4b of cnos#366.
- δ's outward membrane (BoundaryDecision on receipt + V verdict) and override authority. Lives in `cnos.cdd/skills/cdd/delta/SKILL.md §1, §3`; out of scope for the dispatch wire-format.
- CCNF-O orchestration grammar (how cells compose into cycles/waves/roadmaps). Track A of cnos#405.

---

## 1. Define

### 1.1. What a dispatch prompt is

A dispatch prompt is a **wire-format envelope** that γ writes once per role per cycle, and δ routes via the identity-rotation primitive (`cn dispatch`, `claude -p`, or equivalent). It carries:

- the worker role (γ / α / β),
- the project pin,
- the role-skill file to load,
- the cycle issue reference (in the form the worker can fetch directly via `gh issue view <N> --json title,body,state,comments` or MCP equivalent),
- the cycle branch (an explicit `Branch: cycle/<N>` line so the worker never globs or invents),
- for α: the Tier 3 skill list,
- for α: the `## Implementation contract` section enumerating the 7 architectural axes that pin the cycle's shape.

The prompt's job is to make the next role's session **executable without consulting γ's reasoning transcript**. Everything the role needs to start is on the prompt envelope; everything γ used to assemble the prompt (selection arguments, peer enumeration, gap framing) is in the issue body and the artifact channel.

### 1.2. What a dispatch prompt is not

A dispatch prompt is **not**:

- a paraphrase of the issue (the prompt points at the issue; the worker reads the issue),
- a place to smuggle constraints that did not make it into the issue body (if a constraint is load-bearing, it goes in the issue, not in chat prose),
- a place to restate the algorithm or load order (those are in the role-skill the prompt names),
- a substitute for the artifact channel (per-role artifacts live in `.cdd/unreleased/{N}/`, not in chat).

### 1.3. Failure mode

The dispatch wire-format fails through **implementation-contract drift**:

- γ writes the protocol-level contract (gap, ACs, oracle, evidence) but leaves the architectural axes (language, package path, integration target, etc.) unspecified.
- δ routes the α prompt without reviewing the implementation-contract section.
- α improvises the unpinned axes (picks a language; picks a package path; picks a CLI integration shape).
- β verifies behavior-only ACs (the parser accepts the new shape; V validates the receipt) and APPROVE-s.
- The cycle merges in the wrong architectural shape.

cnos#389 (Python-not-Go) and cnos#391 (wrong package scoping + separate binary) are the canonical instances. Both shipped with passing behavior-only ACs and wrong-shape implementations. cnos#392 was the first cycle where δ pinned the contract at dispatch as an ad-hoc operator action; the cycle succeeded specifically because of it. cnos#393 made the inward-membrane function doctrine — the dispatch-time enrichment that pins implementation-level axes γ's protocol-level contract leaves unpinned.

---

## 2. Unfold

### 2.1. Prompt formats

γ produces the prompts; δ routes them — one role at a time, sequentially, per `cnos.cdd/skills/cdd/operator/SKILL.md` Algorithm + `cnos.cdd/skills/cdd/harness/SKILL.md §1`. γ does not execute dispatch; γ produces prompts.

#### γ prompt (δ dispatches γ first)

```text
You are γ. Project: <project>.
Load src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md and follow its load order.
Issue: gh issue view <N> --json title,body,state,comments
```

#### α prompt (γ produces, δ dispatches)

```text
You are α. Project: <project>.
Load src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md and follow its load order.
Issue: gh issue view <N> --json title,body,state,comments
Branch: cycle/<N>
Tier 3 skills: src/packages/cnos.core/skills/write/SKILL.md, <additional issue-specific skills>
```

`write/SKILL.md` is always included — every α output is a written artifact.

#### β prompt (γ produces, δ dispatches)

```text
You are β. Project: <project>.
Load src/packages/cnos.cdd/skills/cdd/beta/SKILL.md and follow its load order.
Issue: gh issue view <N> --json title,body,state,comments
Branch: cycle/<N>
```

### 2.2. Prompt rules

- Point both roles at the issue, **not a paraphrase** of the issue.
- Include the explicit `Branch: cycle/<N>` line so α and β never have to invent or glob-discover the branch.
- Always use `--json title,body,state,comments` with `gh issue view` — bare `gh issue view` hits deprecated GraphQL fields.
- **Do not restate the algorithm in the prompt.** The role-skill named on the `Load …` line governs.
- **Do not smuggle missing constraints into chat prose.** If a constraint is load-bearing, fix the issue body instead.
- β begins reviewing the cycle branch once α signals review-readiness; the β role-skill handles that protocol — the prompt does not re-state it.
- β receives artifact surfaces (diff, `self-coherence.md`, the issue) — **not α's hidden implementation rationale.** Reasoning transcripts do not transit the dispatch channel.

Tooling and observability flags (`--allowedTools`, `--output-format stream-json --verbose`, `--permission-mode acceptEdits`) live in `cnos.cdd/skills/cdd/harness/SKILL.md §1–§2`; the dispatch wire-format does not restate them.

### 2.3. The `## Implementation contract` section (required for α prompt)

γ **MUST** include a `## Implementation contract` section in the α prompt (and reference it from the β prompt) enumerating the 7 architectural axes that pin the cycle's shape. The template:

````markdown
## Implementation contract (pinned by δ; α MUST NOT improvise)

| Axis | Pinned value |
|---|---|
| Language | <Go ; Python ; TypeScript ; Markdown ; ...> |
| CLI integration target | <`cn` subcommand ; standalone binary ; library only ; N/A> |
| Package scoping | <e.g. `src/go/internal/cdd/...` ; `src/packages/<pkg>/...` ; N/A> |
| Existing-binary disposition | <preserve ; replace ; deprecate ; N/A> |
| Runtime dependencies | <list or "None"> |
| JSON/wire contract preservation | <preserve as-is ; explicit shape change with migration ; N/A> |
| Backward-compat invariant | <e.g. "existing rules preserved; new content additive" ; "breaking change documented in §Migration"> |
````

The 7 axes are the **implementation contract** — the architectural shape the cycle ships, distinct from the behavioral ACs the cycle satisfies. Behavioral ACs test what the implementation *does*; the implementation contract pins what the implementation *is* (language, location, integration shape). The two are independent — a behaviorally-correct implementation can still violate the contract by shipping in the wrong language, at the wrong package path, or as a separate binary instead of a `cn` subcommand.

**γ's obligation.** γ writes the contract values per repo conventions or escalates to δ. **γ MUST NOT dispatch with empty / "TBD" rows.** If γ doesn't know a value, γ asks δ — the inward-membrane function (§3 below) exists for this enrichment. δ fills the row per conventions or, if genuinely undecidable, blocks dispatch and escalates to operator-as-human. Mid-cycle contract changes are logged in `.cdd/unreleased/{N}/gamma-clarification.md` (the mid-flight rescue channel; canonical home at [`cnos.handoff/skills/handoff/mid-flight/SKILL.md`](../mid-flight/SKILL.md)).

- ❌ Dispatch with rows empty / "TBD" / implicit conventions ("everyone knows we use Go").
- ❌ γ silently re-pins a row mid-cycle without logging in `gamma-clarification.md`.
- ✅ Every row populated explicitly; δ enrichment recorded; mid-cycle changes logged.

### 2.4. Spec-staleness propagation (γ → in-flight α/β)

When γ (or δ) pushes spec changes to `origin/main` while α/β sessions are in-flight, loaded skills become stale. γ's responsibility — the wire-format rule:

- **Identity-rotation mode (`cn dispatch` / `claude -p`):** Not applicable — each role invocation loads skills fresh from the filesystem; the next dispatch picks them up.
- **Long-lived polling sessions (legacy):** γ writes a coordination note to `.cdd/unreleased/{N}/gamma-coordination.md` on `cycle/{N}` naming the spec change and affected skill path; or surfaces via the dispatch channel.

The list of *which* skill files trigger propagation (CDD.md, role skill files, release/SKILL.md, review/SKILL.md when they land on `main` while a cycle is in-flight) is consumer-side — it is named in `cnos.cdd/skills/cdd/gamma/SKILL.md §2.5 → "Spec-staleness propagation"` and is CDD-specific. The wire-format invariant ("γ writes a coordination note when spec changes mid-flight") lives in [`cnos.handoff/skills/handoff/mid-flight/SKILL.md §2.6`](../mid-flight/SKILL.md) (Sub 4 of cnos#404 / cnos#418); the consumer-specific file list stays in cdd-gamma.

Derives from cnos#301 §9.1: γ proposed an out-of-spec merge because σ's `4a0f678` ("merge is β authority") landed mid-cycle outside γ's loaded snapshot. The reactive fix is the staleness check in γ's Load Order; this is the proactive side.

---

## 3. δ as inward membrane

The **inward face** of the cell's membrane runs at dispatch time, before α is routed. γ writes the protocol-level contract (gap, ACs, oracle, evidence). γ also writes the `## Implementation contract` section enumerating the 7 architectural axes (see §2.3 above):

1. Language
2. CLI integration target
3. Package scoping
4. Existing-binary disposition
5. Runtime dependencies
6. JSON/wire contract preservation
7. Backward-compat invariant

γ populates the rows from repo conventions and the issue body. **At dispatch time, δ reviews γ's dispatch prompt before routing it to α and ensures every row is populated.** If γ left a row unpopulated or marked "TBD", δ has two paths:

- **(a) Fill the row** per repo conventions (e.g. "this repo is Go-native, row 1 = Go"; "this repo uses the `cn` subcommand pattern for protocol-level commands, row 2 = `cn` subcommand"). Log the enrichment in the cycle's artifact channel (`.cdd/unreleased/{N}/gamma-clarification.md`) so the contract trail is auditable.
- **(b) Block dispatch and escalate to operator-as-human** if the row is genuinely undecidable — typically because the choice is part of the cycle's design question, not its execution shape, or because the row's value would commit the repo to a convention that hasn't been settled.

**Why this is δ's surface, not γ's alone.** γ writes what the cycle is *for* (gap, ACs, oracle, evidence). δ writes what the cycle's output is *shaped like* (language, package path, integration target, dependency footprint). The two contracts are distinct: γ's is protocol-level; δ's is implementation-level. γ knows the protocol; δ knows the repo's standing conventions. Mixing them produced cnos#389 (α improvised language because γ's prompt didn't name one and δ didn't catch the omission) and cnos#391 (α improvised package scoping and binary disposition for the same reason). cnos#392 was the first cycle where δ pinned the implementation contract at dispatch; the cycle succeeded specifically because of it.

**The load-bearing claim.** Implementation-contract decisions belong to δ; α MUST NOT improvise. This is the doctrine of cnos#393. If α encounters an unpinned row at intake, α surfaces it to γ/δ before coding — α does not pick a value. If γ wants to change a pinned row mid-cycle, γ logs the re-pin in `gamma-clarification.md`; γ does not silently mutate the contract. β's verification (Rule 7) catches any diff that drifts from the pinned contract regardless of whether the behavioral ACs pass.

### 3.1. The four-surface mesh

The contract distributes across four surfaces, discoverable from any one:

- **γ template:** §2.3 above — the 7-axis `## Implementation contract` block γ injects into the α prompt; γ MUST NOT dispatch with empty rows. Consumer-side role-skill: `cnos.cdd/skills/cdd/gamma/SKILL.md §2.5`.
- **δ enrichment:** §3 above — inward-membrane function; δ reviews γ's prompt; fills or escalates. Consumer-side role-skill: `cnos.cdd/skills/cdd/delta/SKILL.md §2` (now a pointer to this skill).
- **α constraint:** `cnos.cdd/skills/cdd/alpha/SKILL.md §3.6` — "Implementation contract is δ's, not α's"; α MUST NOT improvise; α surfaces unpinned rows to γ/δ before coding.
- **β verification:** `cnos.cdd/skills/cdd/beta/SKILL.md §Role Rules Rule 7` — "Implementation-contract coherence"; β verifies the diff conforms to every pinned axis before APPROVE; non-conformance → REQUEST CHANGES, severity D, classification `implementation-contract`.

Each surface is locally self-justifying via the empirical anchors below; the mesh is for **discoverability** (a future role session loading any one finds the others), not for circular justification. The α and β surfaces are *consumer-side* — they stay in `cnos.cdd` as cycle-lifecycle role-skill obligations; this skill governs the wire-format and the δ enrichment doctrine.

### 3.2. The δ-inward rules

- ❌ δ routes γ's α prompt with rows blank — "α can figure it out".
- ❌ δ fills rows by guessing — no consultation with γ on intent, no anchor in repo conventions.
- ❌ δ enriches but does not log the change, leaving the contract trail invisible.
- ✅ δ reviews γ's `## Implementation contract` section row-by-row; enriches per repo conventions; logs in `gamma-clarification.md`; escalates the row to operator-as-human if undecidable.
- ✅ δ blocks dispatch (does not route the α prompt) until every row is populated or explicitly escalated.

### 3.3. Mid-flight contract re-pinning

When γ needs to change a pinned axis after α has started coding, γ writes an entry to `.cdd/unreleased/{N}/gamma-clarification.md` on the cycle branch *before* signaling the change. The cycle-branch transition is the signal; γ does not need to chat-relay. The clarification names:

- date,
- edit summary (the axis and the new value),
- which ACs / non-goals / constraints / artifacts the re-pin affects,
- the rationale (typically: a repo-convention emerged; a downstream consumer surfaced; an operator override).

α and β polling sessions detect the `gamma-clarification.md` write as a wake-up event and re-read the issue body via `gh issue view` (cache-bust) before continuing. β's Rule 7 verification treats a `gamma-clarification.md`-logged re-pin as authoritative; an un-logged silent re-pin is a contract violation regardless of the new value.

The mid-flight rescue mechanism — file path; authoring role; reader role; trigger conditions; cache-bust semantics — has its canonical home at [`cnos.handoff/skills/handoff/mid-flight/SKILL.md`](../mid-flight/SKILL.md) (Sub 4 of cnos#404 / cnos#418). This skill cites it; mid-flight owns it.

---

## 4. Rules

### 4.1. Author the dispatch prompt at the wire-format shape

- ❌ Paraphrase the issue body into the prompt; omit the `Branch:` line; use bare `gh issue view`; restate the role-skill's algorithm in the prompt.
- ✅ Point the worker at the issue (`gh issue view <N> --json title,body,state,comments`); include `Branch: cycle/<N>`; name the role-skill on the `Load …` line; for α, name Tier 3 skills.

### 4.2. Inject the `## Implementation contract` section in every α prompt

- ❌ Dispatch α without the implementation-contract section because "the issue body covers it" or "α can read repo conventions".
- ✅ Include the 7-axis table; populate every row; do not dispatch with empty / "TBD" rows.

### 4.3. δ reviews the implementation contract before routing

- ❌ Route γ's α prompt without checking the `## Implementation contract` section row-by-row.
- ✅ δ reviews every row; fills unpopulated rows per repo conventions; escalates undecidable rows to operator-as-human; logs enrichments in `gamma-clarification.md`.

### 4.4. α surfaces unpinned rows; α does not improvise

- ❌ α encounters an unpinned row and picks a value ("Go is the convention; I'll use Go").
- ✅ α surfaces the row to γ/δ via the artifact channel; α does not begin coding until the row is pinned.

### 4.5. β verifies the diff against every pinned axis before APPROVE

- ❌ β APPROVE-s on behavior-only AC oracles ("ACs pass; merge").
- ✅ β confirms each pinned axis against the diff; non-conformance → REQUEST CHANGES, severity D, classification `implementation-contract`.

### 4.6. Mid-cycle re-pin requires `gamma-clarification.md`

- ❌ γ silently mutates a pinned row mid-cycle.
- ✅ γ writes a `gamma-clarification.md` entry naming the axis, new value, date, and affected ACs/non-goals/constraints; the cycle-branch transition is the signal.

### 4.7. Cross-reference, do not duplicate

Consumer role-skills cite this skill for the dispatch-prompt template, the 7-axis schema, and the δ-inward-membrane doctrine. If you find dispatch-prompt template content, 7-axis schema content, or δ-inward-membrane doctrine in `cnos.cdd/skills/cdd/gamma/`, `cnos.cdd/skills/cdd/delta/`, `cnos.cdd/skills/cdd/operator/`, or elsewhere that contradicts this skill, this skill governs and the other surface is patched to reference this one.

---

## 5. Empirical anchors

The dispatch + implementation-contract doctrine has empirical anchors across the cnos#388–#412 wave:

- **[cnos#389](https://github.com/usurobor/cnos/issues/389)** — Python-not-Go. α implemented V in Python despite cnos being Go-native. γ's prompt did not name "language" as a pinned axis; δ did not catch the omission at routing; α improvised; β APPROVE-d on behavior-only AC oracles. The diff drifted from the (implicit) architectural shape; the cycle merged in the wrong language.
- **[cnos#391](https://github.com/usurobor/cnos/issues/391)** — wrong package scoping + separate binary. α placed the Go port in a separate binary at the wrong package path. γ's prompt under-specified package scoping and CLI integration target; δ did not catch the omission; α improvised on both axes. cnos#391 is also the empirical anchor for the **mid-flight rescue mechanism** (the `gamma-clarification.md` mechanism crystallized here as a γ→in-flight-α channel; canonical home at [`cnos.handoff/skills/handoff/mid-flight/SKILL.md`](../mid-flight/SKILL.md), landed under Sub 4 of cnos#404 / cnos#418).
- **[cnos#392](https://github.com/usurobor/cnos/issues/392)** — recovery cycle. First cycle where δ pinned the 7-axis implementation contract at dispatch as an ad-hoc operator action. The cycle succeeded specifically because of it; the `cdd-iteration.md` F1–F4 forecast the four patches cnos#393 ships.
- **[cnos#393](https://github.com/usurobor/cnos/issues/393)** — codification. Made the `## Implementation contract` block a first-class γ obligation in the dispatch prompt; made δ-as-inward-membrane doctrine; declared the four-surface mesh (γ template / δ enrichment / α constraint / β verification). This skill is the cnos.handoff-resident extraction of that doctrine.
- **[cnos#397](https://github.com/usurobor/cnos/issues/397)** — Phase 4a of cnos#366. Relocated the δ-inward-membrane substance from `operator/SKILL.md §3a` to `delta/SKILL.md §2`; unified the two-sided membrane framing (outward §1 + inward §2) and the override semantics in one role-skill home.
- **#406–#412 wave** — the post-codification cycles that exercised the doctrine in practice. Every cycle in the wave was dispatched with a populated 7-axis implementation contract; β's Rule 7 verification fired in production; no post-codification cycle in the wave shipped the wrong-shape failure mode that anchored cnos#389/#391.

---

## 6. Verify

### 6.1. Prompt envelope check (per cycle, at dispatch time)

- Does the α prompt name the role, project, role-skill file (`Load …`), issue (`gh issue view <N> --json title,body,state,comments`), branch (`Branch: cycle/<N>`), and Tier 3 skills?
- Does the α prompt include the `## Implementation contract` section with all 7 rows populated (no "TBD")?
- Does the β prompt name the role, project, role-skill file, issue, branch?
- Does the prompt avoid paraphrasing the issue or restating the algorithm?

### 6.2. δ-inward review check (per α dispatch)

- Did δ review the `## Implementation contract` section row-by-row before routing?
- For any row δ filled, is the enrichment logged in `.cdd/unreleased/{N}/gamma-clarification.md`?
- If δ escalated an undecidable row to operator-as-human, is the escalation noted (in `gamma-clarification.md` or the issue)?

### 6.3. Mid-cycle re-pin check (per re-pin event)

- Did γ write a `gamma-clarification.md` entry naming the axis, new value, date, and affected ACs/non-goals/constraints?
- Did γ commit and push to `cycle/{N}` before signaling the re-pin?

### 6.4. β verification check (at APPROVE time)

- Did β confirm each pinned axis against the diff?
- Are non-conformances flagged with severity D, classification `implementation-contract`?

---

## 7. Related documents

- `cnos.handoff/skills/handoff/HANDOFF.md` — package contract; this skill is one sub-surface.
- `cnos.handoff/skills/handoff/cross-repo/SKILL.md` — Sub 2 / cnos#416 (landed). Cross-repo coordination wire-format; structural precedent for this skill's frontmatter shape.
- `cnos.handoff/skills/handoff/mid-flight/SKILL.md` — Sub 4 of cnos#404 (landed under cnos#418). The `gamma-clarification.md` rescue mechanism; this skill cites it for mid-cycle re-pinning.
- `cnos.handoff/skills/handoff/artifact-channel/SKILL.md` — Sub 4 of cnos#404 (landed under cnos#418). The `.cdd/unreleased/{N}/` per-role write ownership; this skill cites it as the channel `gamma-clarification.md` lives on.
- `cnos.handoff/skills/handoff/receipt-stream/SKILL.md` — Sub 5 of cnos#404 (forthcoming). The cross-cycle `cdd-iteration.md` aggregator; this skill cites it for empirical-anchor receipt-stream context.
- `cnos.cdd/skills/cdd/gamma/SKILL.md` — γ role-skill; γ's coordination loop (selection, branch pre-flight, scaffold gate, polling, unblock) lives here; cites this skill for the dispatch-prompt template and the implementation-contract injection rule.
- `cnos.cdd/skills/cdd/alpha/SKILL.md §3.6` — α constraint; cites this skill for the schema and authority.
- `cnos.cdd/skills/cdd/beta/SKILL.md §Role Rules Rule 7` — β verification; cites this skill for the schema and authority.
- `cnos.cdd/skills/cdd/delta/SKILL.md §2` — δ-inward role-side pointer; cites this skill as canonical for the inward-membrane doctrine.
- `cnos.cdd/skills/cdd/operator/SKILL.md §3a` — operator-side pointer; cites this skill as canonical for the inward-membrane doctrine.
- `cnos.cdd/skills/cdd/harness/SKILL.md §1–§2` — dispatch invocation mechanics, observability flags, identity rotation. This skill does not restate them; the dispatch wire-format and the harness substrate are distinct surfaces.
- `cnos.cds/skills/cds/CDS.md §"Coordination surfaces"` — software-class binding; cites this skill for the dispatch and implementation-contract surfaces.
- [cnos#404](https://github.com/usurobor/cnos/issues/404) — parent tracker (handoff/coordination extraction wave).
- [cnos#393](https://github.com/usurobor/cnos/issues/393) — the codification cycle for this surface.

---

## 8. Non-goals

- This skill does NOT own role-local procedure. γ's coordination loop (selection, branch pre-flight, scaffold gate, polling, unblock), α's intake, β's review protocol, δ's outward membrane (BoundaryDecision on receipt + V verdict), and δ's override authority all live in the consumer role-skills under `cnos.cdd/skills/cdd/`.
- This skill does NOT own the mid-flight rescue mechanism. The `gamma-clarification.md` cache-bust semantics, the in-flight α/β polling protocol, and the file's reader/authoring rules live at [`cnos.handoff/skills/handoff/mid-flight/SKILL.md`](../mid-flight/SKILL.md) (Sub 4 of cnos#404 / cnos#418). This skill cites the mechanism for contract re-pinning; mid-flight owns the mechanism.
- This skill does NOT own the artifact-channel rules. The `.cdd/unreleased/{N}/` per-role write ownership, the sequential α→β→γ flow, and the frozen-snapshot rule on merge live at [`cnos.handoff/skills/handoff/artifact-channel/SKILL.md`](../artifact-channel/SKILL.md) (Sub 4 of cnos#404 / cnos#418). This skill cites the channel as the place `gamma-clarification.md` lives.
- This skill does NOT own CCNF-O orchestration. How cycles compose into waves, roadmaps, gates, and joins is Track A of cnos#405; the dispatch wire-format transports work *within* a cycle (γ → α/β), not *across* cycles.
- This skill does NOT own runtime execution. The dispatch invocation primitive (`cn dispatch`, `claude -p`), observability flags, identity rotation, and timeout recovery live in `cnos.cdd/skills/cdd/harness/SKILL.md`. This skill is the prompt; harness is the runtime that routes the prompt.
- This skill does NOT type the 7-axis implementation contract into a CUE schema. Whether to lift the Markdown table into `schemas/handoff/implementation-contract.cue` is deferred; if a future cycle pressures the type-lift, the Markdown template references the schema by `$ref` then.

---

## 9. Kata

### Scenario

γ is scaffolding cycle/{N} against a freshly-quality-gated issue. The issue body declares the gap, ACs, oracle, and evidence. γ has filed `gamma-scaffold.md` to `cycle/{N}`. γ is about to author the α dispatch prompt; δ is the routing actor.

### Task

1. Produce the α prompt at the wire-format shape.
2. Populate the `## Implementation contract` section.
3. Hand off to δ for inward-membrane review.
4. Detect and recover from one unpinned row.

### Expected reasoning

1. **Prompt envelope.** γ writes the α prompt with: role declaration (`You are α`); project pin; role-skill load line (`Load src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md and follow its load order`); issue reference (`Issue: gh issue view <N> --json title,body,state,comments`); branch line (`Branch: cycle/<N>`); Tier 3 skill list (`Tier 3 skills: src/packages/cnos.core/skills/write/SKILL.md, …`).

2. **Implementation contract section.** γ injects the 7-axis Markdown table. γ populates: Language = `Go` (repo convention); CLI integration target = `cn` subcommand (repo convention); Package scoping = `src/go/internal/cdd/<feature>/` (consistent with the cycle's gap); Existing-binary disposition = `preserve` (no existing binary changes); Runtime dependencies = `None` (Go-native; no external deps); JSON/wire contract preservation = `preserve as-is`; Backward-compat invariant = `existing rules preserved; new content additive`.

3. **δ-inward review.** γ hands the prompt to δ. δ reviews the contract row-by-row. δ finds every row populated; δ routes the α prompt via `cn dispatch`.

4. **Unpinned row + recovery.** Suppose γ instead left "JSON/wire contract preservation" empty because the cycle's effect on the wire contract was unclear. δ catches the empty row at routing time. δ has two paths:
   - **(a)** δ fills the row per repo conventions: the cycle does not touch the wire format → row = "preserve as-is". δ logs the enrichment in `.cdd/unreleased/{N}/gamma-clarification.md` and routes the α prompt.
   - **(b)** If the cycle's effect on the wire format is genuinely undecidable (e.g. the cycle's design question *is* whether to change the wire shape), δ blocks dispatch and escalates to operator-as-human. δ does not route until the row is pinned.

### Common failures

- Dispatching α with the implementation-contract section omitted entirely ("the issue body covers it").
- Dispatching with rows populated as "TBD" / "see issue" / "α decides".
- δ routing without reviewing the contract section.
- γ silently mutating a pinned row mid-cycle without writing `gamma-clarification.md`.
- α picking a value for an unpinned row instead of surfacing to γ/δ.
- β APPROVE-ing on behavior-only ACs without checking the diff against the contract.

### Verification

- The α prompt envelope carries every wire-format element (role, project, load line, issue, branch, Tier 3 skills, contract section).
- Every contract row is populated; no "TBD".
- If δ filled a row, the enrichment is logged in `gamma-clarification.md`.
- If a mid-cycle re-pin happens, `gamma-clarification.md` records it before the cycle-branch transition that signals the re-pin.

### Reflection

A second α reading this kata's output should be able to begin coding without consulting γ's reasoning transcript or chat backlog. Every load-bearing decision is on the prompt envelope, in the issue body, or in the artifact channel — never in chat.

---

**End of skill.**
