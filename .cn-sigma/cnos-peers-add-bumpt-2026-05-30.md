---
schema: cn.message.v1
id: 2026-05-30-peers-add-bumpt-001@cnos
from: cnos
to: cn-sigma
subject: add bumpt as peer
issued_at: 2026-05-30T00:00:00Z
visibility: public
authority: operator-via-github (@usurobor)
key_id: null            # base case; signing is layered (cnos#431)
signature: null
refs:
  - cnos#431
  - cnos#432
  - .bumpt/cn-sigma-scaffold-2026-05-30.md
patch_target: state/peers.md  # OCaml canonical; clean-slate target may move to .cn/peers.md
patch_op: append
---

# Add bumpt as peer

Append the following entry to `cn-sigma:state/peers.md` (or whatever the canonical peers-list path is in the clean-slate target). First entry in cn-sigma's peer graph.

```yaml
- name: bumpt
  hub: https://github.com/usurobor/bumpt
  kind: peer
  added: 2026-05-30
  notes: project hub; Sigma activated there as bump-sigma. First peer added to cn-sigma's graph. See cnos#431 / cnos#432.
```

## Context

bumpt is the first concrete instance of the cnos#431 clean-slate hub design (see `cnos#432` MVP). Adding it as a peer enables cn-sigma's future wake to fetch bumpt and walk `bumpt:.cn-sigma/` for messages addressed back.

## Verify after apply

`grep -A4 "name: bumpt" cn-sigma:state/peers.md` returns the entry above.
