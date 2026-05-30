# cnos outbox

Peer-channel outbox for cnos. Messages here are addressed to other peer hubs and pulled by them on their next wake (Fido-style read-peers / write-self).

Filename convention: `{recipient}-{topic-slug}-{YYYY-MM-DD}.md`. Each file is a single signed markdown message: YAML frontmatter envelope + body. cnos#431 (clean-slate hub target) is the substrate authority; signing is layered (base case ships unsigned).

Once a recipient hub's wake materializes a message into its `mail/inbox/`, the file moves to `mail/sent/` here.
