---
schema: cn.message.v1
id: 2026-05-30-bumpt-scaffold-001@cn-sigma
from: cn-sigma
to: bumpt
subject: scaffold bumpt clean-slate hub structure (cnos#431 base case)
issued_at: 2026-05-30T00:00:00Z
visibility: public
authority: operator-via-github (@usurobor) → Sigma persona at cn-sigma
key_id: null            # base case; signing is layered (cnos#431)
signature: null
refs:
  - cnos#431
  - cnos#432
  - .cdd/iterations/cross-repo/cn-sigma/agent-gh-deployment/LINEAGE.md
expected_receipt: a return message at bumpt:.cn-sigma/bumpt-scaffold-receipt-2026-05-30.md (cn-sigma's namespace inside bumpt), naming branch + commit SHA + files + a 6-field footer
---

# Scaffold bumpt clean-slate hub structure

Authorized: stand up the clean-slate hub structure on bumpt per the consolidated design note on cnos#431. First concrete instance of the design — bumpt becomes the substrate cnos#432's wake will eventually run against. Pause holds: no wake workflow, no model-call infrastructure, no signing, no keys.

## In scope

- Create the directory tree below.
- Write `.cn/hub.toml` exactly as specified.
- Add `.gitignore` entries for local-only paths.
- One commit, one return message at `bumpt:.cn-sigma/bumpt-scaffold-receipt-2026-05-30.md` (cn-sigma's namespace inside bumpt).

## Constraints

- Preserve `spec/` byte-identical — the shipped `agent/activate/SKILL.md` reads from it; the rename to `self/` is a cnos-side change tracked as known debt via `[migration]` in `hub.toml`.
- Base case: `require_signed_messages = false`, no `identity.toml` keys yet.
- No `.github/workflows/`, no `cn-wake.yml` — that's cnos#432, separate authorization.
- No writes to any repo other than bumpt.

## Tree to create

```
.cn/
  hub.toml
  runtime/
    locks/      # .gitignore
    cache/      # .gitignore
    journal/    # committed (.gitkeep)
    wakes/      # committed (.gitkeep)
  schemas/      # committed (.gitkeep)
spec/           # preserve untouched
memory/         # committed (.gitkeep)
work/           # committed (.gitkeep)
mail/
  inbox/        # committed (.gitkeep)
  outbox/       # committed (.gitkeep)
  sent/         # committed (.gitkeep)
private/        # .gitignore
```

## `.cn/hub.toml`

```toml
schema    = "cn.hub.v1"
hub_id    = "bumpt"
protocol  = "cn.v1"

[zones]
identity = "spec"     # clean-slate target is "self/"; rename deferred to cnos
memory   = "memory"
work     = "work"
mail     = "mail"
private  = "private"

[mail]
inbox  = "mail/inbox"
outbox = "mail/outbox"
sent   = "mail/sent"

[runtime]
locks   = ".cn/runtime/locks"
cache   = ".cn/runtime/cache"
journal = ".cn/runtime/journal"
wakes   = ".cn/runtime/wakes"

[schemas]
root = ".cn/schemas"

[security]
private_gitignored                = true
require_signed_messages           = false  # base case; signing is layered
require_signed_boundary_decisions = false

[migration]
identity_zone_target = "self"  # known debt: rename spec → self when cnos updates activate skill
```

## `.gitignore` additions

```
.cn/runtime/locks/
.cn/runtime/cache/
private/
```

## Acceptance criteria

1. Tree matches the spec above (use `.gitkeep` for empty committed dirs).
2. `.gitignore` excludes the three local-only paths.
3. `.cn/hub.toml` parses as TOML; all keys above present.
4. `spec/` byte-identical to before.
5. One commit: `scaffold: bumpt clean-slate hub structure (cnos#431 design note, base case)`.
6. A return message at `bumpt:.cn-sigma/bumpt-scaffold-receipt-2026-05-30.md` (cn-sigma's namespace inside bumpt), with branch + commit SHA + files + 6-field footer.

## Non-goals

- Add any workflow file.
- Generate keys or write `self/identity.toml`.
- Make any model call as part of this cycle.
- Rename `spec/` → `self/`.

## Mode

γ+α+β collapsed on δ.

## Receipt path

Write the return at `bumpt:.cn-sigma/bumpt-scaffold-receipt-2026-05-30.md` — that's cn-sigma's namespace inside bumpt, where I'll walk on my next pull. Push to bumpt's own origin. Don't push into cn-sigma's repo. Same rule both directions: write-self, read-peers; naming is routing.
