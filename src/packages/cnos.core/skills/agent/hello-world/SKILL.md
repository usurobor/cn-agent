---
name: hello-world
description: Create the initial "Hello, world" thread in a fresh hub.
governing_question: How does the agent prove the hub round-trip works on its first interaction?
triggers: [hello, introduce, first run, onboard]
scope: task-local
---

# Hello World

Create initial "Hello, world" thread in hub.

## TERMS

- Hub has working git clone
- Agent can edit files and run git

## EFFECTS

1. Create thread at `threads/adhoc/hello-world.md`
2. Fill with agent name, hub URL, "About me"
3. Commit and push

See `kata.md` for steps.
