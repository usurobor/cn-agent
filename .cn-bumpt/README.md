# `.cn-bumpt/` — bumpt's namespace inside cnos

Files dropped here are addressed to the **bumpt** hub. bumpt pulls cnos on its next wake, walks this directory, and materializes each file into its own `mail/inbox/` for processing.

The naming is the routing: any file under `.cn-{recipient-hub-id}/` of a sender's repo is meant for that recipient. The recipient walks its own named directory across each peer it polls.

Filename convention: `{sender-hub-id}-{topic-slug}-{YYYY-MM-DD}.md`. Each file is a single markdown message — YAML frontmatter envelope + body.
