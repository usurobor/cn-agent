# Issue board maps

Interactive, self-contained visualizations of the `usurobor/cnos` open-issue board,
labeled per [the issue taxonomy](../issues/TAXONOMY.md).

- [`index.html`](index.html) — **Voronoi tessellation**: organic cells for
  `kind → area → issue`, colored by priority, with facet filters (kind / area /
  priority / status / dispatchable / search). Click a kind label to zoom in, a cell
  to open the issue; light/dark.
- [`heatmap.html`](heatmap.html) — an `area × kind` heat map + zoomable treemap
  companion.

Each file is fully self-contained (libraries inlined) — open it directly, no server
needed. For a live preview URL without hosting setup, serve the raw file through a
static preview (e.g. `raw.githack.com`).

## Status

These are a **manual static snapshot** (2026-07-02). Automation — a GitHub Action
that regenerates these on issue changes — and a continuous **Pivot-feel deep-zoom**
view are tracked in **#545**. A fully client-side "go live" (fetch issues in the
browser) is a later phase.
