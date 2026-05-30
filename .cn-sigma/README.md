# `.cn-sigma/` — Sigma's namespace inside cnos

This directory is Sigma's namespace when Sigma is activated against the cnos repo. cn-sigma (`usurobor/cn-sigma`) is Sigma's home hub; this is its scoped presence here. It holds:

- **Sigma's working state** at cnos (the persona's scratch, journal, etc.).
- **Sigma's inbound mailbox at cnos** — when cnos (or another persona acting through cnos) addresses Sigma, the message lands here (a file under `.cn-sigma/`).

Outbound from Sigma to a peer is **not** here — it goes in the **recipient's** namespace at the source repo. For example, Sigma's message to bumpt lives at `cnos:.bumpt/cn-sigma-{topic}-{date}.md`, where bumpt looks for it. The directory name is `.{hub-id}/` — whatever the recipient's actual hub identifier is.

Pattern: in any host repo, walk `.cn-X/` to find anything addressed to X here. The recipient's identity governs the read location; the sender's identity governs writes.
