# #627 S-series refinement / supersession map (cnos#671 R1)

Authored by α (this Planning Cell), repairing external-β REQUIRED-finding-5 (the parent-wave
reconciliation was *requested* but not *authored*). External β may test a mapping but must not
author α's planning matter — so this table is α's.

## Framing (single authority preserved)

The `wave/671-cell-runtime-doctrine` wave is the **doctrine/design layer that FEEDS #627's S2–S8
implementation — it does not replace or fork it.** #627 remains the single implementation authority
for its S-nodes. The WC nodes here are **doctrine/design predecessors**: they produce the
`docs/` artifacts a later Go/CUE Working Cell reads so it does not invent architecture or policy.
Implementation ownership of every S-node stays with #627.

Each S-node is classified as **exactly one** of:
- **landed-predecessor** — already shipped; a fixed external root of this wave.
- **refined-by `<WC-n>`** — a #627 implementation node whose design is fed by the named WC
  doctrine artifact(s). The WC is the design predecessor; #627 keeps implementation ownership.
- **unchanged-downstream** — untouched by this wave; runs later on #627's own timeline.
- **superseded** — replaced (none in this wave).

**No S-node has two implementation owners.** Where several WCs feed one S-node they are listed as
design inputs; the single implementation owner is always the #627 S-node itself.

## Table

| #627 S-node | Scope (from #627) | Classification | WC design input(s) | Implementation owner | Note |
|---|---|---|---|---|---|
| **S1** | doctrine — ratify `CELL-RUNTIME.md`; four axes + WC/PC/CC deployment classes; demote `CELL-KINDS`; settle CM↔V | **landed-predecessor** | — | #628 (shipped) | Fixed external root of this wave; consumed by WC-1/WC-2/WC-3b as an immutable `repo_artifact` ref. Not re-opened. |
| **S2** | implementation — extend `schemas/cdd/` with `cell_class` + `matter_domain` + **CM fields**; add `cell.cue` descriptor | **refined-by WC-1, WC-2, WC-4** | WC-1 (contract envelope design), WC-2 (CM-field + receipt/CM/V/δ path design), WC-4 (shipped→specified migration sub-contract) | #627/S2 (unchanged) | The WC wave supplies the *design*; S2 *implements* the CUE. No WC replaces S2. WC-1 is contract-design input, **not a new S1**. |
| **S3** | implementation — mechanical **CM engine** (`wc/pc/cc-cm.v0`) + typed escalation predicate; `cn cell measure`; tsc as provider | **refined-by WC-2 (contract only)** | WC-2 (the CM object contract: producer/inputs/identity/provenance/standing + CM→V edge) | #627/S3 (unchanged) | WC-2 designs the CM **contract**; it **does not replace S3's engine implementation**. tsc stays the CM provider. |
| **S4** | implementation — CC read-only observer MVP (`cn cell pulse --wave N`, receipts-only) | **refined-by WC-3b, WC-2** | WC-3b (wave FSM + CC-disposition→transition-request adapter), WC-2 (the CM the CC judgment consumes) | #627/S4 (unchanged) | CC consumes CM refs (WC-2) and reaches wave state through the WC-3b adapter; CC never mutates state. |
| **S5** | implementation — worker adapter: wrap CDS dispatch behind `cn cell run --class working` | **refined-by WC-3a, WC-1** | WC-3a (cell FSM), WC-1 (typed cell contract) | #627/S5 (unchanged) | WC-3a gives the cell-lifecycle edges the worker adapter drives; WC-1 gives the contract it reads. |
| **S6** | implementation — planning adapter: intent → validated wave graph (`--class planning`) | **refined-by WC-1, WC-3b, WC-2** | WC-1 (contract envelope + `cn.wave.v1` node shape), WC-3b (wave FSM), WC-2 (grounding CM the PC wave grounds in) | #627/S6 (unchanged) | The planning adapter emits waves whose nodes carry complete child contracts and derive edges from sibling_output refs. |
| **S7** | implementation — thin wakes (working/planning/cohering) invoking the runner | **unchanged-downstream** | (reads WC-3a/WC-3b terminals; no new design owed) | #627/S7 (unchanged) | Transport layer; runs after S4–S6. The FSMs it invokes are designed by WC-3a/WC-3b but S7 itself owes no doctrine artifact here. |
| **S8** | experiment — tsc LLM/hybrid escalation wired into the CMs, after mechanical outputs stabilize | **unchanged-downstream** | — | #627/S8 (unchanged) | Experiment; strictly after mechanical CM stabilizes. Out of this doctrine wave. |

## Newly-executable children (none omitted)

The wave's newly-executable children are **WC-1, WC-2, WC-3a, WC-3b, WC-4, WC-5** — all present in
`wave.cn-wave-v1.yaml` with complete `cn.cell.contract.v1` bodies. None is omitted from the graph.
WC-5 (integration/seal) has no #627 S-counterpart: it is the whole-wave proof node this wave adds,
not a refinement of any single S-node.

## Coordination-index note (κ / control-plane, not this cell's matter)

Keeping #627 the single authority requires one immutable coordination index recording this mapping
(either an update to #627 or a named index issue). **That is a control-plane action owned by κ**, not
this Planning Cell — this cell authors the map (above); κ posts/binds it. This cell files and
dispatches nothing.
