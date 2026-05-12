# Cross-Repo Trace Bundles

Each cross-repo coordination unit lives at:

  .cdd/iterations/cross-repo/{target}/{slug}/

Where:
- {target}   — the target repository name (e.g., tsc, acme-api)
- {slug}     — a short, lowercase, hyphen-separated descriptor of the coordination
               (e.g., supercycle-v0.7.0, schema-migration-v2)

## When to create a bundle

Create a bundle when:
- This repo's cycle produces deliverables that land in {target}
- A cycle in this repo coordinates sequentially or in parallel with a cycle in {target}
- A design decision in this repo has binding downstream effect on {target}

## Bundle contents

Each bundle directory contains:
- README.md  — bundle purpose, originating and target cycles, current status
- STATUS     — one of: open | converging | closed (updated as cycles complete)

A bundle is closed when both the originating-repo cycle and the target-repo cycle
have shipped close-outs and the STATUS file reads 'closed'.

## Existing bundles

- `tsc/` — cnos-tsc supercycle coordination (see `tsc/` subdirectory)
