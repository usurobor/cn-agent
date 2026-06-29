# archive/

**Frozen snapshots and historical material.**

Frozen history is still stored at its original snapshot paths during the
migration grace period — version-stamped directories live in-place under their
legacy bundles (for example `alpha/runtime-extensions/1.0.6/`,
`beta/architecture/3.14.4/`, `gamma/cdd/3.x.x/`, `gamma/essays/3.14.5/`). They
are **frozen by policy** and are not rewritten — corrections ship as new
versions or superseding notes.

For now, historical material may remain under legacy paths so older references
do not break. This directory will become the canonical archive index later; no
snapshot migration happens until after Demo 0.

## Policy

Per [`beta/governance/DOCUMENTATION-SYSTEM.md`](../beta/governance/DOCUMENTATION-SYSTEM.md):
version directories are immutable. Link/path repairs in active reader paths are
allowed; semantic edits to a frozen snapshot are not. Release evidence under
`.cdd/releases/` (outside `docs/`) is likewise frozen.
