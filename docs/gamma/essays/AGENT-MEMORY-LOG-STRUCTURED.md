# Agent Memory Is Coherence-Preserving Log Structure

Agent memory exists to preserve coherence across time.

The log remembers what happened. Reflection remembers what mattered for coherence. State remembers where to resume. Spec and protocol remember what must constrain future work.

## Governing question

How should cnos make memory first-class without adding a memory database before the existing Git-native surfaces are exhausted?

## The law

Use one law:

```text
Log first.
Compact later.
Promote when stable.
Measure coherence.
```

A shorter form:

```text
Memory is coherence-preserving compaction over append-only evidence.
```

This law absorbs three separate threads: activation logs, eventual consistency, and KISS for v0.

Activation logs give the raw stream. Eventual consistency gives the cursor discipline. KISS says not to build a retrieval system before the log and compaction model fail. Coherence explains why anything is kept, compacted, or promoted at all.

## Coherence governs memory

Salience is not the governing principle. Salience is a signal.

A thing becomes memory because retaining it reduces future incoherence. It helps the next activation avoid re-deriving a frame, repeating a failure, violating an obligation, losing a decision, or drifting away from the system it claims to continue.

So the acquisition question is not:

```text
Is this interesting?
```

It is:

```text
Will losing this make future work less coherent?
```

The compaction question is not:

```text
What is the summary?
```

It is:

```text
What relation must survive so the next cycle still describes the same system?
```

The promotion question is not:

```text
Was this repeated?
```

It is:

```text
Has this become a constraint on future coherent action?
```

## TSC is the measurement ground

TSC is the foundation under the word coherence.

TSC asks whether three independent descriptions of the same system still describe one system. It observes three axes: pattern coherence, relational coherence, and process coherence.

Agent memory should satisfy the same triad.

```text
α / pattern:
  memory has stable terms, ranks, paths, and entry shapes

β / relation:
  memories cite their inputs, outputs, cursors, and authority surfaces consistently

γ / process:
  memory can evolve through logs, reflections, and promotions without losing identity
```

That makes memory measurable, not just poetic.

A memory corpus with high α but low β has tidy folders whose files do not actually relate. A memory corpus with high β but low γ can explain the past but cannot guide the next change. A memory corpus with high γ but low α keeps moving, but its own terms drift.

`.tsc/` is not the memory store. It is the coherence measurement output. TSC does not replace logs, reflections, state, or specs. It measures whether those surfaces still cohere.

## What this changes for #100

#100 already names the right gap: cnos has persistence in practice, but memory is not yet a declared retention faculty. The agent has reflections, adhoc threads, hub state, Git history, activation logs, and runtime continuity, but no single answer to what to read first, what to write back, what counts as durable, and what is only runtime projection.

The missing answer is not a new store.

The missing answer is structural:

```text
canonical memory is log-structured and coherence-governed
```

`threads/adhoc/`, activation logs, mail, and other event streams are `r0`: retained evidence. `threads/reflections/` is the recursive compaction tree over that evidence. `state/` is an index and cursor surface. `spec/` and protocol packages are promotion targets. `.tsc/` is the measurement surface that can tell whether the memory corpus still forms one system.

## Rank is not frequency

Rank names derivation depth.

Frequency names schedule.

Those are different axes.

```text
r0  raw retained evidence
r1  first compaction over r0
r2  compaction over r1
r3  compaction over r2
rN  compaction over rN-1
```

Weekly, monthly, quarterly, and yearly are cadence labels. They are not the ontology.

A weekly reflection is usually `r1` if it reads raw logs. A monthly reflection is usually `r2` if it reads weekly reflections. But an intense day can need an `r1` daily reflection, and a missed week can force a monthly reflection to read `r0` directly.

So the rule is:

```text
frequency is operational
rank is structural
```

Do not encode truth in the folder name alone. Encode it in the artifact basis.

## r0 is concise evidence

`r0` refuses to guess what will matter later.

It records events, decisions, corrections, directives, receipts, handoffs, activation notes, and facts that would otherwise vanish in chat.

But `r0` is not a transcript dump. Raw does not mean garbage.

A good `r0` entry is short evidence:

```markdown
## 2026-05-30 — short subject

What happened.

What changed.

What needs attention.

Read/seen cursor if relevant.
```

The current agent activation log shape is already a good `r0` substrate: one writer per stream, date-sharded Markdown, H2 entry boundaries, and Git commit cursors.

The coherence rule for `r0` is:

```text
Record enough to preserve future coherence.
Do not record so much that future coherence is buried in noise.
```

## Compaction is a receipt

A reflection is not a replacement for the log. A reflection is a receipt over the log.

It says what it read, what it retained, what it dropped, and what should be promoted. It also says why that selection preserves coherence.

Use explicit basis metadata:

```yaml
rank: r1
period: weekly
covers:
  from: 2026-05-24
  to: 2026-05-30
inputs:
  - rank: r0
    path: threads/activations/cnos/20260530.md
    at: ec9d7bb3
  - rank: r0
    path: threads/adhoc/20260530-agent-comms-futures-kiss.md
    at: 1da16aa
basis:
  cn-sigma: 1da16aa
  cnos: ec9d7bb3
coherence:
  retained_because:
    - activation logs should be agent-level, not Sigma-only
    - flat logs became ambiguous after multi-activation traffic
  risks:
    - signing layer remains deferred
outputs:
  retained:
    - activation logs should be agent-level, not Sigma-only
  promoted:
    - target: docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md
      claim: activation home logs should be date-sharded directories
  dropped:
    - old flat home-log path
  open:
    - signing layer remains deferred
```

Pin inputs by commit SHA, not just by path.

A path can keep changing. A compaction needs a reproducible basis. The SHA is the cheap version of the signed-feed discipline: it says exactly which file state the reflection read.

## Higher rank does not erase lower rank

Higher-rank memory can be wrong. That is safe only if the lower rank still exists.

`r1` may compress `r0` badly. `r2` may overfit the week. `r3` may promote the wrong doctrine candidate. The repair path is possible because the raw evidence stays available.

So the safety rule is:

```text
higher-rank memory must cite lower-rank memory
```

No compaction deletes the basis it compacted.

This is also a coherence rule. If the reflection and the log disagree, the system can inspect the relation instead of trusting the most recent prose.

## Promotion is not a rank

Do not call promoted memory `r4`.

Promotion is not a deeper reflection. Promotion is a different event.

The reflective tower stays one kind of artifact:

```text
r0 → r1 → r2 → r3 → ...
```

Promotion moves a stable lesson out of the reflective tower into a governing surface:

```text
reflection candidate → explicit edit to state/spec/protocol
```

Examples:

```text
activation path rule       → convention/spec
operator discipline        → spec/OPERATOR.md
persona identity change    → spec/PERSONA.md
routing cursor             → state/activations.md
software lifecycle rule    → cnos.cds
handoff transport rule     → cnos.handoff
project-local constraint   → project repo docs/state
```

The reflection says why the rule emerged. The promoted surface says what the rule is.

## Promotion needs a citation

Promotion must be explicit.

No auto-promotion. No silent crystallization. No rule becoming authoritative merely because it appeared in a reflection.

A reflection may name candidates:

```yaml
outputs:
  promoted:
    - target: spec/OPERATOR.md
      claim: activation logs are date-sharded on both sides
      reason: flat files became ambiguous after multi-activation traffic
```

The actual spec or state edit should carry a back-reference:

```yaml
derived_from:
  - rank: r1
    path: threads/reflections/weekly/2026-W22.md
    at: 1da16aa
    claim: activation logs are date-sharded on both sides
```

The exact syntax can wait. The invariant should not.

A promoted rule must name its reflective basis.

## Activation reflections are identity-attached

Activation logs are per activation.

Activation reflections default to the home identity.

That means Sigma reflecting across cnos, bumpt, and cph should usually write one home-owned reflection:

```text
cn-sigma:threads/reflections/weekly/2026-W22.md
```

not three separate per-activation reflection trees.

The reason is ownership. The reflection belongs to Sigma as an identity, not to the foreign body that happened to expose the lesson.

A per-activation reflection path may be useful later for high-volume or long-lived bodies:

```text
threads/reflections/activations/cnos/weekly/2026-W22.md
```

But v0 should not add that tree by default. One identity-level reflection can synthesize across many activation logs.

## State is an index

State helps the next activation resume.

It should not become the place where meaning lives.

Examples:

```text
state/activations.md     activation registry and cursors
state/conversation.json  working continuity
.cn-*                    local body, outbox, inbox, runtime attachment
.tsc/                    generated coherence measurements
```

State may be retained. State may be loaded. State may be essential for continuity.

But state is not the canonical memory of what happened or why it matters. Canonical memory remains inspectable prose plus Git history. TSC reports remain generated measurements over that memory, not replacements for it.

## This is event sourcing for coherent agent memory

The analogy is old and useful.

Event sourcing stores state changes as a sequence of events, then derives current state by replaying or projecting those events. Kappa-style systems similarly prefer a durable log and derived processing over duplicate batch and speed paths.

Agent memory follows the same shape, but with a different payload:

```text
r0 logs           event stream
r1+ reflections   projections / compactions
state files       indexes / materialized views
spec/protocol     promoted constraints
.tsc reports      coherence measurements
```

The distinctive agent rule is citation-bearing compaction at every level.

A database projection can be overwritten because the event log can rebuild it. An agent reflection can be challenged because the lower-rank evidence and commit basis remain visible. A promoted rule can be reviewed because it points back to the reflection that earned it.

## The restore path

A fresh activation should not read everything.

It should restore from the highest useful layer downward:

```text
1. load constitutive spec and operator rules
2. load state indexes and cursors
3. load latest high-rank reflection
4. load lower-rank reflections since that basis
5. load r0 logs since the last cursor only when needed
6. follow explicit references into adhoc threads or receipts
7. consult .tsc reports when coherence measurement is relevant
```

Read enough to recover continuity. Do not scan the archive because the archive exists.

## The write path

A working activation should write at the lowest layer that truthfully captures the event.

```text
something happened              → r0 log
something mattered              → r1 reflection
pattern repeated                → r2/r3 synthesis
rule must constrain future      → promote to state/spec/protocol
coherence needs measurement     → run/record TSC report
```

This prevents two failures.

First, the agent does not edit doctrine for a one-off observation.

Second, the agent does not bury behavior-changing lessons in raw logs where no future activation will read them.

## TSC closes the loop

Memory is coherent when it can be measured as one system across time.

For memory targets, TSC should eventually ask:

```text
α: Do the memory artifacts use stable ranks, paths, terms, and entry forms?
β: Do reflections, promotions, cursors, and source citations actually line up?
γ: Can the memory corpus continue through new activations without losing identity?
```

This does not mean every log entry needs a measurement run.

KISS still holds. Use TSC when the memory surface itself changes, when a promotion crosses into spec/state/protocol, or when a reflection wave claims a major coherence delta.

## KISS for v0

Do not add a memory database yet.

Do not add embeddings to core yet.

Do not add `threads/memory/INDEX.md` yet.

Do not add a dedicated memory package yet.

Do not make the memory faculty depend on a retrieval backend.

Do not require TSC measurement for every r0 entry.

The current shape is enough:

```text
r0 streams
recursive reflections
state indexes
explicit promotion
Git history
optional .tsc measurement
```

Add machinery only when a concrete failure appears:

```text
restore cost is too high
privacy boundaries require redacted/public split
many activations need delivery receipts
search becomes necessary for real work
manual compaction is skipped often enough to hurt continuity
coherence drift is not visible without measurement
```

Until then, the implementation is discipline, not infrastructure.

## What #100 should absorb

#100 should not grow into a memory platform.

It should name the memory faculty as a restore and write discipline over existing surfaces:

```text
backend: git+threads+state
selection rule: expected coherence delta
entrypoint: latest promoted spec/state plus latest reflection basis
surfaces: r0 evidence streams, r1+ reflection streams, state indexes, promoted spec/protocol surfaces
measurement: optional TSC target/report when coherence claims need grounding
freshness: latest compaction basis and unread r0 span
scope: restore continuity, preserve decisions, and promote stable rules
```

The memory skill, if added, should be thin.

`adhoc-thread` can still own durable topic capture. `reflect` can still own compaction. Activation logs can still own cross-body continuity. TSC can still own measurement. The memory faculty owns the restore map and the rules that connect them.

## The model

```text
r0 log
  → r1 reflection
    → r2 synthesis
      → r3+ deeper synthesis
        → explicit promotion into state/spec/protocol
          → optional TSC measurement of coherence
```

The log remembers what happened.

The reflection remembers what mattered for coherence.

The spec remembers what must constrain the future.

TSC tells whether the memory surfaces still describe one system.

That is agent memory without a memory store.

## Related analogues

- TSC, [docs/THESIS.md](https://github.com/usurobor/tsc/blob/main/docs/THESIS.md)
- TSC Core, [spec/tsc-core.md](https://github.com/usurobor/tsc/blob/main/spec/tsc-core.md)
- TSC Self-Measure, [runtime/SELF-MEASURE.md](https://github.com/usurobor/tsc/blob/main/runtime/SELF-MEASURE.md)
- Martin Fowler, [Event Sourcing](https://martinfowler.com/eaaDev/EventSourcing.html)
- Microsoft Azure Architecture Center, [Event Sourcing pattern](https://learn.microsoft.com/en-us/azure/architecture/patterns/event-sourcing)
- Jay Kreps, [Questioning the Lambda Architecture](https://www.oreilly.com/radar/questioning-the-lambda-architecture/)
- Agent Activation Log Convention v0, [`docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md`](../conventions/AGENT-ACTIVATION-LOG-v0.md)
- Companion essays: [`AGENT-ACTIVATION-LOGS-AND-EVENTUAL-CONSISTENCY.md`](./AGENT-ACTIVATION-LOGS-AND-EVENTUAL-CONSISTENCY.md), [`AGENT-COMMS-FUTURES-KISS.md`](./AGENT-COMMS-FUTURES-KISS.md)
