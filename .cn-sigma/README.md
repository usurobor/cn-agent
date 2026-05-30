# `.cn-sigma/` — Sigma's working state inside cnos

This directory holds the **Sigma persona's** working state when Sigma is activated against the cnos repo. cn-sigma (`usurobor/cn-sigma`) is Sigma's home hub; this is its scoped presence here.

Pattern: a persona acting inside any host repo writes its state under `.cn-<persona-hub-id>/`. The host repo (cnos) keeps its own identity; the persona keeps its own state and outbox; the two don't collide.

Contents follow the hub-shape design (cnos#431 clean-slate target):

- `mail/outbox/` — Sigma's outbound messages addressed to peers (Fido-style, write-self).
- (additional zones — `memory/`, `work/`, `mail/inbox/`, `mail/sent/` — added when needed.)

Filenames in `mail/outbox/`: `{recipient}-{topic-slug}-{YYYY-MM-DD}.md`. Each is a single markdown file with a YAML envelope frontmatter + body. Peers pull on their next wake; once materialized into their `mail/inbox/`, the file moves to `mail/sent/` here.
