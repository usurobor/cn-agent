# α Close-out — Cycle #394

**Date:** 2026-05-21
**Issue:** [#394](https://github.com/usurobor/cnos/issues/394) — Sub 2 (#376): bootstrap cnos.cdr package skeleton
**Branch:** `cycle/394`
**α identity:** alpha / alpha@cdd.cnos
**Outcome:** R1 APPROVE; 6/6 ACs PASS.

---

## §1 Cycle summary

Authored three package-skeleton files under `src/packages/cnos.cdr/`:

- `cn.package.json` (9 lines) — `cn.package.v1` schema; `name: cnos.cdr`; `version: 0.1.0`; `kind: package`; `engines.cnos: ">=3.81.0"`. No commands declared.
- `README.md` (~55 lines) — What CDR Does, Package Structure (naming CDR.md + SKILL.md), Forthcoming surfaces (naming Sub 3 + Sub 4), Quick Start, Status (v0.1 skeleton), License. Cross-references CDR.md 8 times.
- `skills/cdr/SKILL.md` (~95 lines) — loader skill with standard frontmatter (name, description, artifact_class, governing_question, triggers, scope, plus optional fields); Load order; Rule; Role overlays (forthcoming — Sub 3); Cross-protocol relationship; Conflict rule. `calls:` lists CDR.md + five role overlays (α/β/γ/δ/ε per ROLES.md §1).

Plus cycle evidence under `.cdd/unreleased/394/`: gamma-scaffold, design-notes, self-coherence, beta-review, this file, and forthcoming beta-closeout / gamma-closeout / cdd-iteration.

## §2 Findings (factual only)

### F1 (α) — Schema source-of-truth is `pkg.go`, not a separate JSON schema document

The cn.package.v1 schema is defined in Go (`src/go/internal/pkg/pkg.go` lines 107-167) rather than in a JSON-Schema document. The exemplar source-of-truth for what fields are required is the `FullPackageManifest` struct and the `ValidatePackageManifestData` function. This is observed, not a finding worth filing — the design is consistent with cnos's Go-as-source-of-truth discipline.

### F2 (α) — `cn build --check` is the discovery oracle; no separate `cn doctor` package-list step

The kernel's package-discovery output comes from `cn build --check`. There is no separate `cn doctor` or `cn list` that enumerates packages without validating them. The check + list are unified. This was confirmed by reading `cn help` output: `doctor` exists but requires a hub (none in the worktree); `build --check` works without a hub and lists all six packages including the new cnos.cdr.

### F3 (α) — The cnos.cdr v0.1 skeleton is **loadable** but **not yet dispatchable**

The package is discoverable (`✓ cnos.cdr: valid`), but the role-overlay layer (Sub 3) is not yet authored, so a dispatch matching `cdr` triggers will load CDR.md as the canonical contract but will not find a per-role overlay until Sub 3 ships. This is by design (the issue's `Non-goals` explicitly defers role-overlays to Sub 3). The SKILL.md loader body explicitly names this as "Role overlays (forthcoming — Sub 3 of cnos#376)". A future dispatch error message should distinguish "role overlay missing because Sub 3 unshipped" from "role overlay missing due to misconfiguration". Recorded for Sub 3 dispatch context.

## §3 Debt

- **D1** — `engines.cnos` range vs pin policy is not centrally documented. The convention "first-party packages aligned with kernel pin exact; new sub-packages range" is observable from the package set but not stated in any normative document. If future kernel releases enforce a stricter compatibility policy, the cnos.cdr engine range may need re-evaluation. Not blocking.
- **D2** — The Sub 3 role-overlay paths (`alpha/SKILL.md` ... `epsilon/SKILL.md`) named in the loader's `calls:` do not yet exist. The current `cn build --check` does not enforce `calls:` path existence; if a stricter skill-loader is introduced later, the loader will need a "forthcoming" marker or Sub 3 will need to land in lockstep. Recorded as forward-looking debt.
- **D3** — Per cycle 390 F1, the verdict-enum schema (`schemas/cdr/verdict.cue`) is not yet authored. cnos.cdr's manifest does not currently reference any cdr-specific schema beyond `schemas/cdr/receipt.cue` (cited in README.md). If a verdict.cue lands later, README.md may want to cross-reference it. Carried forward from cycle 390; not a Sub 2 finding.

## §4 What unblocks downstream

- **Sub 3 (cnos#376)** — role overlays can author `skills/cdr/{alpha,beta,gamma,delta,epsilon}/SKILL.md` against the loader's named paths. Sub 3 will not need to re-author the loader; only the role files.
- **Sub 4 (cnos#376)** — empirical-anchor doc can cite cnos.cdr/README.md's "Forthcoming surfaces" section as its dispatch context. Sub 4's content goes into `skills/cdr/EMPIRICAL-ANCHOR.md` or similar; cnos.cdr/README.md may be updated to cross-reference the new file once it ships.
- **First CDR cycle** — once Sub 3 + Sub 4 ship, a CDR cycle can run end-to-end with the typed `#CDRReceipt` (`schemas/cdr/receipt.cue`) backed by α/β/γ/δ/ε role overlays + project-binding from cph.

## §5 Hand-off

β review complete (APPROVE). γ proceeds to closeout + INDEX update + merge.
