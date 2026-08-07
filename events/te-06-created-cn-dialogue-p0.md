---
schema: cnos.task-event.v1
event: te-cn-sigma-06
ts: 2026-08-06T14:00:00Z
task: impl-cn-dialogue-command
type: created
actor: {agent: usurobor/cn-sigma, locus: usurobor/cnos, runtime: {engine: claude, surface: claude-code}}
title: "Implement `cn dialogue` command to write dialogue entries (write-own-ref + schema enforced)"
owner: usurobor/cn-sigma
priority: P0
spawned_by: operator
design_issue: usurobor/cnos#715
initial_state: open
---
Rationale: hand-crafting dialogue events + ad-hoc worktree pushes caused routing/schema
mistakes (wrong ref, drift, mis-addressed `to`). A `cn dialogue` command enforces
append-to-own-feed-only + valid cnos.agent-message.v1 + ff-only. Going forward cn-sigma
writes dialogue only through this command.

Execution (per task skill): issue-to-create-design (#715) -> design -> impl wave ->
implementation, delegated to a generic CDD cell with dialogue/write/review skills injected.
