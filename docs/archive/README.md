# archive/

**Frozen snapshots and historical material.**

> **Pass 1 overlay.** Version-stamped snapshot directories currently live
> in-place under their feature bundles (for example
> `alpha/runtime-extensions/1.0.6/`, `beta/architecture/3.14.4/`,
> `gamma/cdd/3.x.x/`, `gamma/essays/3.14.5/`). They are **frozen by policy**
> and are not rewritten — corrections ship as new versions or superseding
> notes.

## Policy

Per [`beta/governance/DOCUMENTATION-SYSTEM.md`](../beta/governance/DOCUMENTATION-SYSTEM.md):
version directories are immutable. Link/path repairs in active reader paths
are allowed; semantic edits to a frozen snapshot are not. Release evidence
under `.cdd/releases/` (outside `docs/`) is likewise frozen.

A later pass will gather the in-place snapshots under this directory (or
index them from here) without altering their contents.
