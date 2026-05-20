# Tightening proposal — retarget #32 from "author reusable CDR skills in cph" to "author cph Stream-B binding contract"

From cnos γ via cross-repo coordination (`cnos.cdd/skills/cdd/gamma/SKILL.md §"Cross-repo proposal close-out"`).

This comment proposes a scope retarget for #32 in light of three things that landed on the cnos side after #32 was filed:

1. **`ROLES.md §4a` (cnos, 2026-05-19)** — declares the five-layer enforcement chain: persona, operator contract, protocol overlay, project binding, receipt/validator. The chain explicitly separates protocol doctrine (lives in `cnos.cdr`) from project binding (lives in `cph:.cdr/`) from persona identity (lives in `cn-rho/spec/`).
2. **cnos#376 (modified 2026-05-19)** — adds AC6 to cdr v0.1: persona/protocol/project boundary declared. The reusable α/β role skills for research belong in `cnos.cdr/skills/cdr/{alpha,beta}/SKILL.md`, not in `cph/cdr/`.
3. **cn-rho persona hub drafts (cnos, 2026-05-19)** — `usurobor/cn-rho` is the planned research persona hub; persona-level identity (Rho's discipline profile, refusal conditions, voice) belongs there. Drafts staged at `cnos:.cdd/iterations/cross-repo/cn-rho/bootstrap-2026-05-19/`.

Under the new layering, the substance #32 wants to ship splits cleanly across three homes:

| #32's current content | Belongs in | Rationale |
|---|---|---|
| Reusable α/β role doctrine (Stream-B α runs scripts against real data; Stream-B β independently re-runs and audits) | `cnos.cdr/skills/cdr/{alpha,beta}/SKILL.md` (cnos#376 Sub 3) | This is **protocol overlay** — same shape across any research project. |
| Epistemic conservatism, refuse-to-fabricate, no-data-no-claim posture | `cn-rho/spec/PERSONA.md` `## Discipline` section | This is **persona identity** — same shape across any cph-like project a research persona attaches to. |
| Citing each numerical claim to a producing script/data/SHA | `cnos.cdr` receipt schema (claim_refs, data_refs, method_refs, result_refs, claim_status) | This is **receipt schema** — protocol layer 5. |
| The specific cph data root, manifest checksum, script paths, field-report convention, `/opt/gait-data/` mount expectations | `cph:.cdr/STREAM-B.md` (or `cph:cdr/README.md`) | This is **project binding** — cph-specific and stays in cph. |

## Preferred retarget

Rename #32 from "author reusable CDR skills in cph" to **"author cph Stream-B binding contract."**

Produce a single file in cph:

```
cph/.cdr/STREAM-B.md   (or cph/cdr/README.md — pick by cph convention)
```

This file says:

- Stream-B science work on `cph` uses Rho (when `usurobor/cn-rho` lands) + `cnos.cdr` (when v0.1 lands per cnos#376).
- Until those upstream packages land, Stream-B work on cph follows the discipline drafted at `cnos:.cdd/iterations/cross-repo/cn-rho/bootstrap-2026-05-19/PERSONA.md` and the protocol AC6 stub at cnos#376.
- The cph data root is `/opt/gait-data/` (or current truth); a `data-mounted` gate check runs before any compute. The check verifies the mount path exists and the manifest checksum matches `cph/.cdr/data-manifest.sha256` (or current truth).
- Field reports cite producing scripts and SHAs per the receipt schema named in `cnos.cdr` (when shipped).
- The reusable α/β role-skill content is **not** authored in cph; cph imports/loads it from `cnos.cdr/skills/cdr/{alpha,beta}/SKILL.md` once that exists. Until then, cph operates under cdd-with-overrides (as it does today) and cites the cnos#376 stub.

This file is cph's **project binding** — the WHERE in the five-layer chain. It records what cph carries that any persona + protocol pair needs to attach to.

## Acceptable short-term fallback

If retargeting #32 fully is undesirable while cnos.cdr v0.1 hasn't landed, fold the persona/protocol/project boundary into #32 by marking any cph-local CDR role-skill drafts as **temporary prototypes**, with a required banner:

```markdown
> **Temporary cph-local prototype.**
> Protocol content here promotes to `cnos.cdr` (cnos#376) when v0.1 ships.
> Persona content here promotes to `cn-rho` when that hub is scaffolded.
> Do not treat this directory as the durable source of truth.
> See `cnos:ROLES.md §4a` for the persona/protocol/project boundary.
```

This prevents the half-migration trap: cph carries the doctrine temporarily but its provenance and durable home are named.

## Why this matters for #32 specifically

cph#32's empirical anchor (`/opt/gait-data/` not mounted; numbers cited as if produced by execution that did not occur) is exactly the failure that the persona/protocol/project chain is designed to prevent:

- The persona layer (Rho's `## Discipline`) refuses dispatch when the data mount is absent.
- The protocol layer (`cnos.cdr`'s β oracle) requires reproduction; β re-runs the cited script and discovers the mount was absent.
- The project binding layer (`cph/.cdr/STREAM-B.md`) names the specific mount path the gate checks against.
- The receipt validator (V, deferred to cnos #366 Phase 3) rejects receipts that lack `data_refs` or have `reproduction: false`.

The four layers each catch the failure at a different point. The same content authored only in `cph/cdr/` provides one layer of catching, in one project.

## Reach check

cnos session cannot write to `usurobor/cph` (MCP scoped to `usurobor/cnos`). Operator action options:

- Paste this comment on cph#32 directly.
- Or summarise it to cph γ for in-session retargeting.
- Or apply the comment as a markdown file under `cph:.cdr/notes/` for cph γ to action.

Whichever is convenient; the substance is the retarget proposal, not the channel.

## Related cnos artifacts (for context when retargeting)

- `cnos:ROLES.md §4a` (persona / operator / protocol / project / receipt enforcement chain — 2026-05-19)
- `cnos:ROLES.md §4a.1` (discipline-profile required section)
- `cnos:ROLES.md §4a.2` (engineering vs research loss-function distinction)
- `cnos:ROLES.md §4a.3` (receipt schema enforces discipline mechanically)
- cnos#376 (modified 2026-05-19; carries AC6 for the boundary, AC7 for the architecture-choice question)
- `cnos:.cdd/iterations/cross-repo/cn-rho/bootstrap-2026-05-19/` (Rho persona hub drafts; pre-scaffold)

Filed by cnos γ on 2026-05-19.
