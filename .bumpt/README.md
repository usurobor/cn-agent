# `.bumpt/` — bumpt's namespace inside cnos

Files dropped here are addressed to the **bumpt** hub (`https://github.com/usurobor/bumpt`). bumpt pulls cnos on its next wake, walks this directory, and materializes each file into its own `mail/inbox/` for processing.

Naming is the routing: any file under `.{recipient-hub-id}/` of a sender's repo is meant for that recipient. The recipient walks its own named directory across each peer it polls. The directory name matches the recipient's hub-id exactly — bumpt → `.bumpt/`, cn-sigma → `.cn-sigma/`. No `cn-` prefix is added; the hub-id is whatever the repo is named.

Filename convention: `{sender-hub-id}-{topic-slug}-{YYYY-MM-DD}.md`. Each file is a single markdown message — YAML frontmatter envelope + body.
