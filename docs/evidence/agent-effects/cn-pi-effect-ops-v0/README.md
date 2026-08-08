# cn-pi effect-operations v0 research evidence

Status: **superseded bridge design; retained as non-authoritative CNOS research**

This directory freezes the detailed effect-operation contract developed during
`usurobor/cn-pi#2` and PR #3 at reviewed head
`52a4cb1a813b17fadc4483f18078ca9166081ab2`.

The design was intentionally removed from the temporary bridge's shipping
contract after Pi's exact-head review showed that it was becoming a general
execution protocol. The cn-pi bridge instead ships a small REST-shaped adapter
for issues, issue comments, and textual documents.

The retained document is useful input for the deferred CNOS runtime research in
`usurobor/cnos#714`, especially:

- plan identity versus transport/provider identity;
- operator policy as a non-widening authority boundary;
- ledger-before-effect and exact replay;
- backend-specific preconditions and honest uncertainty;
- typed observations, receipts, and incident projection;
- eventual promotion from a temporary adapter into a governed executor.

It is not a current CNOS specification, implementation issue, compatibility
promise, or dispatch authority. Future work must compare it with empirical
bridge receipts and CNOS's then-current runtime architecture, promoting only
invariants that still earn their cost.

## Provenance

- cn-pi issue: <https://github.com/usurobor/cn-pi/issues/2>
- cn-pi PR: <https://github.com/usurobor/cn-pi/pull/3>
- independent beta R4: <https://github.com/usurobor/cn-pi/pull/3#pullrequestreview-4885734769>
- Pi exact-head re-review: `msg-cn-pi-home-effect-ops-v0-rereview-05` on
  `refs/heads/cn-pi/home/dialogue`
- CNOS deferred research issue: <https://github.com/usurobor/cnos/issues/714>

The adjacent `EFFECT-OPS-v0.md` bytes are copied from the reviewed cn-pi head.
