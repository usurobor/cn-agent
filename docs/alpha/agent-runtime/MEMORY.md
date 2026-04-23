# Memory in cnos — lean triadic model

**Issue:** #100  
**Version:** 0.2.0  
**Mode:** MCA  
**Active Skills:** cdd/design, eng/evolve, eng/writing  
**Engineering Level:** L7

---

## Problem

cnos already has memory in practice, but the current model is still narrower than the way the system actually works.

Today:

- The runtime contract classifies `threads/reflections/` and `state/conversation.json` as memory
- `state/` otherwise remains `private_body`
- The context packer loads reflections into the dynamic system block
- The context packer loads conversation history as recent message turns
- The adhoc-thread skill already treats standalone topic threads as the place where important things should be preserved instead of being lost in chat

So the system already distinguishes:

- Durable reflective memory
- Session continuity
- Runtime internals

What is missing is a simple architectural statement of how these fit together. The current incoherence is:

1. `threads/adhoc/` is functionally memory, but is not yet named as such in the runtime contract
2. Conversation history is useful, but it is too easy to mistake it for canonical memory
3. There is pressure to add a memory index or retrieval layer before the underlying memory model is explicit
4. Packages, commands, and orchestrators could add memory tooling, but core memory semantics are not yet clearly separated from optional memory services

---

## Decision

Use a lean triadic memory model.

### Core memory classes

- **α — episodic memory:** factual retained record of what happened and what was decided. Surface: `threads/adhoc/`
- **β — reflective memory:** interpretation and synthesis over episodic memory. Surface: `threads/reflections/`
- **γ — working continuity:** recent session/conversation continuity across wake boundaries. Surface: `state/conversation.json`

### Canonical rule

Canonical memory is:

- `threads/adhoc/`
- `threads/reflections/`

`state/conversation.json` is:

- Useful
- Retained
- Runtime-visible
- But not canonical memory

It is a working continuity surface.

### KISS / YAGNI rule

Do not add:

- `threads/memory/INDEX.md`
- Vector stores
- Graph stores
- Retrieval indexes
- Memory-specific package roles
- A dedicated memory package

...until the current three-surface model proves insufficient.

---

## Constraints

- Keep canonical memory Git-native and inspectable
- Do not add opaque mandatory storage
- Do not require a retrieval backend
- Preserve current runtime behavior where possible
- Keep the distinction between constitutive self, memory, and private runtime state explicit
- Avoid adding a third canonical memory surface unless there is real operational need

## Challenged Assumption

The challenged assumption is:

> Memory needs an additional index or retrieval layer in core before the underlying memory model is explicit.

That assumption is unnecessary right now. cnos already has enough surfaces to define a coherent memory model without introducing another canonical layer.

---

## Current Evidence

### Runtime contract

The runtime contract already classifies:

- `threads/reflections/`
- `state/conversation.json`

as memory surfaces, distinct from:

- `constitutive_self`
- `private_body`
- `work_medium`
- `projection_surface`

### Context packer

The context packer already loads:

- Daily and weekly reflections into the dynamic system block
- Conversation history into recent message turns

It does not treat all state equally.

### Skills

`adhoc-thread` already defines a durable thread as the place for a proposal, learning, question, or decision that would otherwise be lost if left inline.

`reflect` already defines reflection as: evidence → interpretation → conclusion.

So the practical model already exists. It is just not yet stated cleanly.

---

## Proposal

### 1. Core memory semantics

Core should own only the memory semantics that must be true for every hub:

#### 1.1 Canonical episodic memory — `threads/adhoc/`

What belongs here:

- Decisions
- Investigations
- Proposals
- Durable questions
- Corrections worth preserving
- Significant conversation outcomes once promoted out of chat

This is the factual retained record.

#### 1.2 Canonical reflective memory — `threads/reflections/`

What belongs here:

- Daily synthesis
- Weekly synthesis
- Higher-order lessons
- Explicit MCA/MCI outcomes
- Reflections over patterns in the raw record

This is the interpretive retained record.

#### 1.3 Working continuity — `state/conversation.json`

What belongs here:

- Recent turns
- Session continuity
- Immediate conversational context
- Resumable short-term history

This is useful retained state, but not the canonical record of memory.

---

### 2. Triadic fit

#### α — episodic memory

The α side preserves the pattern of what happened. Its failure mode is:

- Facts stay in chat and disappear
- Decisions are remembered socially but not recorded
- Investigations leave no durable trace

#### β — reflective memory

The β side links the parts into meaning. Its failure mode is:

- Conclusions detached from evidence
- Lessons that cannot be traced back to what actually happened
- Reflections that substitute for the record instead of interpreting it

#### γ — working continuity

The γ side carries the current line of interaction across a wake boundary. Its failure mode is:

- The agent loses the immediate thread of a conversation
- Stale short-term history is mistaken for durable truth
- Runtime projection is treated as if it were canonical memory

#### Memory is coherent when

- α preserves the facts
- β interprets the facts
- γ carries the current thread without claiming to be the source of truth

---

### 3. Ownership split

#### cnos core owns

- The distinction between episodic, reflective, and working memory
- Runtime-contract visibility of those surfaces
- Doctor/status validation of memory surfaces
- The rule that `state/conversation.json` is not canonical memory
- The rule that constitutive self is not memory

#### Agent layer owns

The existing skills already provide most of the needed memory discipline:

- `adhoc-thread` — when to retain a fact or decision durably
- `reflect` — how to synthesize evidence into a durable lesson

No new agent/memory skill is required in this v1 design.

#### Commands / orchestrators may later own

Optional tooling only:

- Conversation import
- Memory compaction
- Recall helpers
- Search/indexing helpers
- Reflection promotion workflows

Those are layered behaviors, not memory ontology.

---

### 4. Runtime contract changes

Update the memory zone so it reflects the actual lean model.

#### Current

- `threads/reflections/`
- `state/conversation.json`

#### Proposed

Memory should distinguish canonical vs working surfaces more clearly. Suggested rendering:

```json
"medium": {
  "zones": [
    { "path": "threads/adhoc/", "zone": "memory_episodic" },
    { "path": "threads/reflections/", "zone": "memory_reflective" },
    { "path": "state/conversation.json", "zone": "memory_working" },
    { "path": "spec/SOUL.md", "zone": "constitutive_self" },
    { "path": ".cn/", "zone": "private_body" }
  ]
}
```

If introducing three new zone names feels too heavy, then at minimum:

- Add `threads/adhoc/` to memory
- Document in prose that `state/conversation.json` is working continuity, not canonical memory

That lighter version is acceptable for v1.

---

### 5. Context-packing rule

Do not blindly load all adhoc threads into context.

Current behavior is already mostly right:

- Reflections are directly useful as compact memory
- Conversation history is directly useful as recent turns

For `threads/adhoc/`, use:

- Explicit reference from the current task
- Or later, an optional retrieval/indexing layer

This avoids flooding the prompt while still keeping adhoc threads canonical.

---

### 6. What not to add yet

#### No `threads/memory/INDEX.md`

Not yet. Reason:

- Runtime does not use it today
- Current memory surfaces are enough
- A third canonical memory file would add ceremony before necessity

If restore cost becomes too high later, add it as a derived restore map, not a third source of truth.

#### No vector/graph store in core

Not needed. Those can be optional later.

#### No dedicated memory package

Not needed yet. The existing reflect and adhoc-thread split is already enough.

---

## Leverage

This design:

- Makes the current memory surfaces explicit
- Keeps core small
- Keeps memory inspectable and Git-native
- Avoids introducing an unnecessary index layer
- Gives commands/orchestrators/extensions a clear non-core role
- Aligns the runtime contract, context packer, and existing skills

### Negative Leverage

This adds:

- One more explicit distinction to maintain in docs and runtime contract
- One more rule for what is canonical and what is projection
- Some migration work for docs that currently speak about memory more loosely

---

## Alternatives Considered

| Option | Pros | Cons | Decision |
|--------|------|------|----------|
| Keep current implicit model | No work | adhoc stays semantically invisible; confusion remains | Rejected |
| Add `threads/memory/INDEX.md` now | Explicit restore surface | Extra canonical layer too early | Rejected |
| Add vector/index backend in core | Better retrieval | Makes memory backend-defined and opaque too early | Rejected |
| New dedicated agent/memory skill now | One obvious home | Adds a layer before the current split proves insufficient | Rejected |
| **Lean triadic model: adhoc + reflections + conversation** | Small, truthful, aligned with current system | Requires explicit documentation of boundaries | **Chosen** |

---

## File Changes

### Edit

- `docs/alpha/agent-runtime/RUNTIME-CONTRACT-v2.md` — add `threads/adhoc/` to the memory model, clarify `state/conversation.json` as working continuity
- `src/cmd/cn_runtime_contract.ml` — reflect the same distinction in emitted contract
- `docs/alpha/agent-runtime/AGENT-RUNTIME.md` — align memory language with the lean model if needed

### No new files in v1

Do not add:

- `threads/memory/INDEX.md`
- New memory package
- New retrieval extension
- New memory-specific command package

---

## Acceptance Criteria

- [ ] `threads/adhoc/` is explicitly recognized as canonical memory
- [ ] `threads/reflections/` remains canonical reflective memory
- [ ] `state/conversation.json` is explicitly treated as working continuity, not source of truth
- [ ] Runtime contract reflects the distinction clearly
- [ ] Docs no longer imply that memory requires an index layer in core
- [ ] No new canonical memory surfaces are introduced in v1
- [ ] Existing reflect and adhoc-thread skills remain the practice layer for memory

---

## Known Debt

- If wake/restore cost grows, a derived restore map may become useful later
- Optional recall/search/import tooling may still be worth packaging later
- The exact runtime-contract zone names may need one tightening pass for clarity

---

## CDD Trace

| Step | Artifact | Skills loaded | Decision |
|------|----------|---------------|----------|
| 0 Observe | — | — | Current cnos already distinguishes reflections and conversation history, and adhoc-thread already functions as retained memory in practice |
| 1 Select | — | — | Selected gap: memory architecture is overcomplicating itself before the current surfaces are fully named |
| 4 Gap | this artifact | — | Named incoherence: memory semantics are implicit and invite unnecessary new layers |
| 5 Mode | this artifact | cdd/design, eng/evolve, eng/writing | L7 MCA; lean architectural clarification |
| 6 Artifacts | this artifact | — | Design drafted; no implementation plan included yet |

---

## Related

- #100 — Memory as first-class retention faculty (this design supersedes the original spec direction)
- AGENT-NETWORK.md — agents carry memory when deployed to new workspaces
- HUB-PLACEMENT-MODELS.md — hub is memory, workspace is workbench
- #156 — Attached hubs (AC9: agent memory stays in hub, tagged by workspace)
