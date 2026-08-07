---
schema: cnos.task-event.v1
event: te-cn-sigma-07
ts: 2026-08-06T14:02:00Z
task: create-task-skill
type: created
actor: {agent: usurobor/cn-sigma, locus: usurobor/cnos, runtime: {engine: claude, surface: claude-code}}
title: "Author the task skill — task = a higher-level intent above an issue; defines the task execution algorithm"
owner: usurobor/cn-sigma
priority: P1
spawned_by: operator
initial_state: open
---
Spec (operator): the task skill is the issue-skill analog one level up. A TASK captures an
intent. Executing a task follows the same algorithm:
  issue-to-create-design -> design -> issues-wave-to-implement-design -> implementation,
delegated to a generic CDD cell with proper skills injected (including review skills).
Placement (agent vs cnos.cdp/cdd) to be decided by its design step.
