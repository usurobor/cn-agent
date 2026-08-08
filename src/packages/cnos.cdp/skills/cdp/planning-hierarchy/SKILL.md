---
name: cdp/planning-hierarchy
description: Structure a body of work so every unit rolls up to an owned workstream through a track, with one canonical source of truth per level. Use when organizing a backlog, defining workstreams/tracks, assigning ownership, or making a program legible at a glance.
governing_question: How do we structure a body of work so it reads with extreme clarity — every unit rolling up to an owned workstream through a track, with one canonical source of truth per level?
triggers:
  - organizing a backlog into initiatives
  - defining workstreams or tracks
  - assigning single-threaded ownership
  - deciding where the source of truth for an initiative lives
  - a program is hard to reason about at a glance
scope: task-local
artifact_class: skill
kata_surface: none
inputs:
  - a body of work (backlog of issues or tasks)
  - current ownership and sources of truth
outputs:
  - a `workstream ⊃ track ⊃ unit` structure with single-threaded owners
  - one canonical doc per level
---

# Planning Hierarchy

## Core Principle

A body of work is legible when every unit rolls up to one owned workstream, and every level has one canonical source of truth.

`workstream ⊃ track ⊃ unit`. A workstream is a top-level area of work with a single-threaded owner and its own canonical doc. A track is a sub-workstream inside exactly one workstream, also owned. A unit (issue, task) belongs to exactly one track. No orphan work; no diffuse ownership; no scattered truth.

This adapts Naomi Gleit's practice: workstreams with single-threaded owners, each with its own canonical doc linked from an overall canonical doc, recursively — the discipline she calls Canonical Everything, in service of Extreme Clarity. Source: Naomi Gleit, "Canonical Everything" — https://naomi.com/canonical-everything-c85441a84e70; interview context in Lenny's Newsletter — https://www.lennysnewsletter.com/p/metas-head-of-product-naomi-gleit.

## Algorithm

1. Define — name the workstreams, tracks, units, owners, and canonical docs, and the failure mode.
2. Unfold — build top-down: workstreams, then tracks, then assign every unit; give each level an owner and a canonical doc.
3. Rules — one workstream per initiative, one owner per level, one canonical doc per level, no orphan unit.

---

## 1. Define

### 1.1. Identify the parts

A coherent program has these parts:

- **Workstream** — a top-level discrete area of work. Owned. Has its own canonical doc.
- **Track** — a sub-workstream inside exactly one workstream. Owned.
- **Unit** — one issue or task, inside exactly one track.
- **Single-threaded owner** — the one person or agent accountable for a level.
- **Canonical doc** — the one source of truth for a level.
- **Nomenclature** — the shared, defined vocabulary the program is described in.

- ❌ "These issues are roughly about the platform"
- ✅ "`workstream/platform`, owned by cn-sigma, canonical doc #452, tracks: wake-runtime, scheduler, trace"

### 1.2. Articulate how they fit

The workstream sets the area. Each track answers a sub-area. Each unit does one piece. Each level names one owner and one canonical doc. Reading up from any unit reaches exactly one track, one workstream, one owner, one source of truth.

- ❌ A unit that belongs to two workstreams "because it touches both"
- ✅ A unit in one track; the cross-cutting relationship is a link, not a second parent

### 1.3. Name the governing question

State the program's shape in one sentence:

> This workstream is the discrete area of `<X>`, owned by `<owner>`, canonical at `<doc>`.

If a workstream needs two owners or two canonical docs, it is two workstreams.

- ❌ "Substrate-and-platform-and-tooling, owned by the team"
- ✅ "`workstream/threads`, owned by cn-sigma, canonical #698"

### 1.4. Name the canonical doc

The canonical doc is the one place a level's truth lives; everywhere else points to it.

Per level:

- the overall program has one canonical doc that lists the workstreams
- each workstream has its own canonical doc, linked from the overall one
- Canonical Everything recurses: a track large enough to need its own truth gets its own canonical doc, linked from the workstream's

- ❌ The plan restated in three issues that drift apart
- ✅ One canonical doc per level; the rest link to it

### 1.5. Name the failure mode

Planning fails through four losses:

- **Orphan work** — a unit belongs to no track, so no one owns it
- **Diffuse ownership** — a level has zero or many owners, so decisions stall
- **Scattered truth** — the same plan lives in many places and drifts
- **Vocabulary drift** — the same word means different things to different people

- ❌ "Someone will pick it up"
- ✅ "Every unit has a track; every track has an owner"

---

## 2. Unfold

### 2.1. Name the workstreams (top-down, MECE)

Start at the top. Name the discrete areas of work so they are mutually exclusive and collectively exhaustive — every unit will fit exactly one, and none is left over.

Workstreams are at **initiative grain**, not theme grain: `generic-cell`, `threads`, `memory` — not one giant `substrate`. Expect roughly a dozen, not three.

- ❌ Five broad buckets that force unrelated work together
- ✅ ~12 initiatives, each a real area someone can own

### 2.2. Give each workstream one owner and one canonical doc

Before subdividing, assign the single-threaded owner and the canonical doc. A workstream without an owner is a wish; without a canonical doc it will drift.

- ❌ "`workstream/memory` — we'll sort ownership later"
- ✅ "`workstream/memory`, owner cn-sigma, canonical #690"

### 2.3. Subdivide into tracks

Inside each workstream, name the sub-workstreams. A track answers one sub-area of its workstream. Give each track an owner too when the workstream is large enough that one owner cannot hold it all.

- ❌ A workstream with forty loose units and no internal structure
- ✅ `workstream/generic-cell` → tracks: cell-fsm, cell-classes, planning-cell, claim-dispatch

### 2.4. Assign every unit to exactly one track

Walk the backlog. Each unit gets one track — and therefore one workstream and one owner. A unit that seems to fit two tracks picks its primary; the other relationship is a link.

- ❌ Leave twenty units untriaged "for now"
- ✅ Every unit carries one track; zero orphans

### 2.5. Maintain the overall canonical doc

Keep one overall canonical doc that lists the workstreams, their owners, and their canonical docs. It is the entry point; it recurses into each workstream's doc.

- ❌ Ask three people where the plan is and get three answers
- ✅ One overall canonical doc; it points to every workstream's canonical doc

### 2.6. Hold the nomenclature fixed

Extreme Clarity needs a shared vocabulary. Define the terms once — what a workstream is, what a track is, what each named workstream means — and use them the same way everywhere.

- ❌ "Stream" here, "track" there, "area" elsewhere, all for the same thing
- ✅ One glossary of terms; every doc uses them identically

---

## 3. Rules

### 3.1. One workstream per initiative

A workstream is one discrete area with one owner and one canonical doc. If it needs two of either, split it.

- ❌ "`workstream/platform-and-product`"
- ✅ "`workstream/platform`" and "`workstream/product`"

### 3.2. Exactly one track per unit

Every unit belongs to one track, and therefore one workstream. Cross-cutting is a link, not a second parent.

- ❌ A unit filed under two tracks
- ✅ One primary track; related work referenced by link

### 3.3. Every level has a single-threaded owner

One accountable owner per workstream, and per track when needed. Not zero. Not a committee.

- ❌ "The team owns it"
- ✅ "cn-sigma owns `workstream/threads`"

### 3.4. Every level has one canonical doc

One source of truth per level; everywhere else points to it. Canonical Everything recurses down.

- ❌ The plan duplicated across issues that drift
- ✅ One canonical doc per level; the rest link

### 3.5. No orphan units

A unit with no track has no owner. Triage it to a track or close it.

- ❌ A backlog of unassigned issues
- ✅ Every open unit carries a track

### 3.6. Workstreams are MECE

Mutually exclusive, collectively exhaustive: every unit fits exactly one, none is left over, none spans two.

- ❌ Overlapping workstreams that argue over the same units
- ✅ A partition — each unit lands in exactly one

### 3.7. Initiative grain, not theme grain

Prefer a dozen real initiatives over a handful of broad themes. A workstream someone can actually own beats a bucket no one can.

- ❌ Three themes: substrate, platform, everything-else
- ✅ generic-cell, threads, memory, dematerialization, wake-runtime, …

### 3.8. Keep the vocabulary fixed

Define the terms once and use them identically. Rename deliberately, everywhere at once, never casually.

- ❌ "stream" and "track" used interchangeably
- ✅ `workstream ⊃ track`, stated once, held everywhere

---

## 4. Realization

This doctrine realizes as issue labels and the board, owned by `cnos.issues` — cite, do not fork:

- `workstream/*` and `track/*` labels are defined in [`docs/development/issues/TAXONOMY.md`](../../../../../../docs/development/issues/TAXONOMY.md).
- The overall canonical doc is `docs/development/issues/WORKSTREAMS.md`; each workstream's canonical doc is its **master issue**.
- `cn issues map` renders `workstream → track → issue` from the labels and registry.

Extend `TAXONOMY.md` first, then let the board consume it. If this doctrine and the taxonomy disagree, change the taxonomy to match, not the reverse.

---

## 5. Final Test

A program is coherent when:

- every unit rolls up to exactly one track and one workstream
- every workstream and needed track names one owner
- every level has one canonical doc, and the overall doc reaches them all
- the workstreams partition the work — MECE, initiative-grained
- one vocabulary describes it all

If a unit has no track, triage it. If a level has no owner, assign one. If the truth lives in two places, pick one home and point to it.
