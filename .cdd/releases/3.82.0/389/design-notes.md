<!-- sections: [Impact, Implementation choice, CUE invocation strategy, ValidationVerdict schema, Counterfeit-receipt rules, Backward compatibility, Alternatives considered, Risks] -->

# Design Notes — Cycle #389

L7 design surface for the V validator (`cn-cdd-verify` rewrite). See `gamma-scaffold.md` for issue framing.

## §Impact

V is the executable spine of the recursive coherence-cell kernel. Before V is mechanical, every "validated" claim is prose-trusted. After V is mechanical:

- δ (Phase 4) can mechanically gate boundary actions on `ValidationVerdict.result`
- ε (Phase 6) can mechanically aggregate `protocol_gap_refs` across receipt streams without re-parsing prose
- CDR (cnos#376) can ship a research overlay validated by the same V (only the schema dispatch differs)
- `CDD.md` rewrite (Phase 7) becomes safe — V works, so the doctrine can quote V as binding runtime law

The L7 frame: changing this surface makes ~five downstream surfaces cheaper or unnecessary (Phase 4 δ split has a typed verdict to read; Phase 5 γ shrink has a derived `transmissibility` to point at; Phase 6 ε has `protocol_gap_refs` to consume; Phase 7 CDD rewrite has V to cite; cnos#376 CDR can extend cleanly).

## §Implementation choice

**Decision: hybrid bash-wrapper + Python core.**

Three options considered:

### Option A — extend existing bash (`cn-cdd-verify`) in pure bash

- **Pro:** zero new dependencies; existing operator-facing CLI unchanged
- **Con:** parsing YAML frontmatter with bash is fragile; JSON emission needs `jq`-as-templating (poor signal/noise); evidence-ref dereference + counterfeit-receipt graph rules are awkward in bash arrays
- **Rejected because:** counterfeit-receipt rules (AC8) compare git commit timestamps and per-file authors; this is fluent in Python and grim in bash.

### Option B — Go rewrite (cnos v3 → v4 transition direction)

- **Pro:** typed, native CUE Go API available, future-aligned with v4 direction
- **Con:** introduces a new build dependency and shippable binary; requires CI to build/ship binaries before V is usable in cycles; the cnos v3 codebase has no Go infrastructure currently
- **Rejected because:** Phase 3 must ship V *now* (gates Phase 4); a Go binary is a bootstrapping problem that compounds the cycle's risk surface. v4 can re-implement V in Go later once its build infra exists.

### Option C — hybrid (chosen)

- Existing `cn-cdd-verify` bash unchanged for current operator modes (`--version`, `--pr`, `--cycle`, `--all`, `--unreleased`, `--triadic`); new flag `--receipt <path>` (optionally with `--contract <path>` and `--json`) dispatches to a sibling Python helper `cn-cdd-validate-receipt`
- Python helper implements V — YAML parsing, `protocol_id` dispatch, `cue vet` subprocess invocation, evidence-ref dereference, counterfeit-receipt rule evaluation, JSON emission against the JSON Schema
- **Pro:** preserves backward compat with existing operator usage exactly; introduces only a Python script (stdlib `yaml` ships via PyYAML which is already a project dep — see `requirements.txt` / system Python); keeps the CUE invocation a clean subprocess boundary (the schemas remain authoritative; V wraps them)
- **Pro:** Python is already used in the repo (`tools/`); no new build dep
- **Con:** two-language surface (bash + Python). Mitigation: the bash wrapper is a thin dispatch layer; all V logic lives in Python.

PyYAML availability check before commit: `python3 -c "import yaml" 2>&1 || echo MISSING`. Fallback: use stdlib JSON parsing on a YAML→JSON pre-conversion via `cue` itself (CUE can parse YAML and emit JSON).

## §CUE invocation strategy

Per essay §3 ("Let V wrap CUE and evidence checks"):

- **CUE validates structural shape** (`#Receipt` / `#CDSReceipt` / `#CDRReceipt` / `#ValidationVerdict` / `#BoundaryDecision` / `#Transmissibility` / `#Override` / `#ProtocolGapRef`)
- **V wraps CUE** and additionally:
  - Dereferences `evidence_refs.*` (reads the files at the named paths; asserts existence; reads contents for counterfeit-receipt rules)
  - Applies counterfeit-receipt rules (§Counterfeit-receipt rules below)
  - Emits a typed `ValidationVerdict` regardless of CUE pass/fail (CUE failure → V FAIL with failed_predicates = CUE diagnostic lines)
- **δ does not read raw evidence** — δ reads only `ValidationVerdict` + the receipt. V is the only layer that touches `evidence_refs.*` contents.

Dispatch table (read `protocol_id` from receipt frontmatter):

```text
cnos.cdd.generic.receipt.v1 → cue vet -c -d '#Receipt'      schemas/cdd/{contract,boundary_decision,receipt}.cue <receipt>
cnos.cdd.cds.receipt.v1     → cue vet -c -d '#CDSReceipt'   schemas/cds/receipt.cue <receipt>
cnos.cdd.cdr.receipt.v1     → cue vet -c -d '#CDRReceipt'   schemas/cdr/receipt.cue <receipt>
<other>                     → V FAIL with failed_predicates ["unknown protocol_id: <id>"]
```

CUE invocation is via `subprocess.run(["cue", "vet", ...])` with `check=False`; non-zero exit → CUE diagnostic captured in `failed_predicates`. Schema package directory is anchored from `git rev-parse --show-toplevel`.

## §ValidationVerdict schema

`schemas/cdd/validation_verdict.schema.json` (JSON Schema draft-07), aligned with `schemas/cdd/boundary_decision.cue` `#ValidationVerdict` plus the additional `failed_predicates` structure RECEIPT-VALIDATION.md §Validation Interface specifies (predicate refs with `evidence_ref` and `detail`).

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "$id": "cnos.cdd.validation_verdict.v1",
  "type": "object",
  "required": ["result", "failed_predicates", "warnings", "provenance"],
  "properties": {
    "result": { "enum": ["PASS", "FAIL"] },
    "failed_predicates": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["predicate", "diagnostic"],
        "properties": {
          "predicate": { "type": "string" },
          "evidence_ref": { "type": "string" },
          "diagnostic": { "type": "string" }
        }
      }
    },
    "warnings": { "type": "array", "items": { "type": "string" } },
    "provenance": {
      "type": "object",
      "required": ["validator_identity", "validator_version", "checked_at", "input_refs"],
      "properties": {
        "validator_identity": { "type": "string" },
        "validator_version": { "type": "string" },
        "checked_at": { "type": "string" },
        "input_refs": {
          "type": "object",
          "required": ["contract_ref", "receipt_ref"],
          "properties": {
            "contract_ref": { "type": "string" },
            "receipt_ref": { "type": "string" },
            "evidence_root_ref": { "type": "string" }
          }
        }
      }
    }
  }
}
```

Two name choices made:

1. **`result`** (V's emitted JSON), not `verdict` (the CUE-typed field name). Rationale: `verdict` is the CUE structural enum inside `validation.verdict` of the receipt; `result` distinguishes V's JSON emission from the receipt's stored field. The JSON Schema is V's *output* shape, not the receipt's *stored* shape.

2. **`failed_predicates` as object array** (not string list per `boundary_decision.cue` `#ValidationVerdict.failed_predicates: [...string]`). Rationale: RECEIPT-VALIDATION.md §Validation Interface specifies the richer shape (`predicate_ref`, `evidence_ref`, `detail`). The CUE schema accepts strings; the JSON emission carries the structured form. V's JSON-to-receipt write path serializes the structured form into receipt as a string-by-stringify-detail (the receipt's `failed_predicates: [...string]` survives a future schema-bump cycle).

`validator_identity` is pinned to `cnos.cdd.validate_receipt` per RECEIPT-VALIDATION.md illustrative example.

`validator_version` is read from the version constant baked into `cn-cdd-validate-receipt` itself (`V_VERSION = "phase3.0"`).

## §Counterfeit-receipt rules

V applies these rules **after** CUE structural validation passes (or, for AC3/AC4/AC5, V reports CUE failure first; counterfeit rules don't run if CUE already failed). Each rule is a single mechanical predicate.

### Rule C1 — Actor separation (α actor ≠ β actor)

**For:** CDS protocol only (CDR's research-reviewer-rerun discipline is a different rule).

**Predicate:** read `evidence_refs.alpha_closeout` and `evidence_refs.beta_closeout`; parse YAML frontmatter or first non-blank "Actor:" / "Author:" / "α actor:" / "β actor:" lines; assert the α actor name ≠ β actor name.

**Fallback when actor field is absent:** read `git log --format='%ae' -1 -- <file>` for each file; assert author emails differ.

**Failure form:** `failed_predicates += [{predicate: "counterfeit.actor_separation", evidence_ref: "<alpha_closeout> / <beta_closeout>", diagnostic: "α and β closeout authored by same actor: <name>"}]`.

### Rule C2 — β verdict precedes merge (timestamp/ancestry)

**For:** CDS protocol only.

**Predicate:** read `evidence_refs.beta_closeout` git log; record commit timestamp T_β. Read `boundary_decision.decided_at` (timestamp T_δ). Assert T_β ≤ T_δ.

**Fallback when no merge SHA available** (unreleased cycle): assert `beta_closeout` exists and is non-empty before any `boundary_decision.action == accept | release`.

**Failure form:** `failed_predicates += [{predicate: "counterfeit.verdict_precedes_merge", ..., diagnostic: "β closeout commit T_β <ts> postdates boundary decision T_δ <ts>"}]`.

### Rule C3 — Override does not rewrite verdict

**For:** all protocols. (This is partly enforced by CUE; V double-checks at the JSON level.)

**Predicate:** if `boundary_decision.action == "override"`, assert `boundary_decision.override.original_validation_verdict.verdict == validation.verdict`. Override carries the original verdict verbatim; if a counterfeit author rewrote validation.verdict to PASS but kept action=override, original_validation_verdict will not match.

**Failure form:** `failed_predicates += [{predicate: "counterfeit.override_does_not_rewrite", ..., diagnostic: "override claims original verdict <X> but receipt.validation.verdict is <Y>"}]`.

### Rule C4 — Evidence-ref dereferenceability

**For:** all protocols.

**Predicate:** for every entry in `evidence_refs.*` whose value looks like a filesystem path (starts with `.`, `/`, or matches `*.md` / `*.yaml` / `*.json` etc.; explicit exclusion of `git_diff(...)` / `sha256:...` / `issue:#...` synthetic refs), assert the file exists at the receipt's evidence root.

**Failure form:** `failed_predicates += [{predicate: "counterfeit.evidence_ref_unresolved", evidence_ref: "<key>: <path>", diagnostic: "evidence_ref <key> points at non-existent path: <resolved-path>"}]`.

### Rule C5 — protocol_id matches schema-declared id

**For:** all protocols.

**Predicate:** if `protocol_id == cnos.cdd.cds.receipt.v1`, assert all CDS-required `evidence_refs` keys are present (covered by `cue vet -c` — but V re-asserts for a friendlier diagnostic).

**Failure form:** `failed_predicates += [{predicate: "counterfeit.protocol_id_mismatch", ..., diagnostic: "protocol_id declares cds.receipt.v1 but missing required evidence_refs keys: [<list>]"}]`.

### Rule scope

These rules together discharge AC8. They are **mechanical** — every rule reduces to a finite-time predicate over (receipt frontmatter, evidence files on disk, git log). No prose interpretation.

The rule set is the cycle's load-bearing contribution. It is intentionally narrow: V does not attempt to detect *all* counterfeits, only those with structural signal. Detecting "the author lied in the prose of self-coherence.md" is out of scope (that's β's job, not V's). Detecting "the actor metadata says α and β are the same person" is in scope (mechanical).

## §Backward compatibility

The existing `cn-cdd-verify` command's invocations all continue to work unchanged:

- `cn cdd-verify --version 3.74.0` → existing artifact-presence check (no V invocation)
- `cn cdd-verify --all` → existing repo-wide ledger check (no V invocation)
- `cn cdd-verify --unreleased` → existing unreleased check (no V invocation)

New flags only:

- `cn cdd-verify --receipt <path>` → invoke V on the receipt; emit ValidationVerdict (prose to stdout by default)
- `cn cdd-verify --receipt <path> --json` → invoke V; emit JSON to stdout
- `cn cdd-verify --receipt <path> --contract <path>` → invoke V with explicit contract path (default: derived from receipt's `provenance.input_refs.contract_ref` if resolvable, otherwise V reports a `warnings` entry but proceeds — contract presence is not load-bearing for AC1-AC8)

Exit codes:

- `--receipt` mode: exit 0 if `result: PASS`, exit 1 if `result: FAIL`, exit 2 if V itself errored (parse fail, missing schema package, etc.)
- Other modes: unchanged

No existing operator workflow breaks. No new top-level commands.

## §Alternatives considered

### Alt-1: Make V a separate CLI (`cn validate-receipt`)

**Rejected.** #389 non-goals explicitly forbid new commands ("Do NOT add CLI surface for `cn schemas verify` or `cn validate` beyond what `cn-cdd-verify` already exposes"). Extending `cn-cdd-verify` is the constrained path.

### Alt-2: Implement counterfeit-receipt rules in CUE alone

**Rejected.** CUE cannot read git history or dereference filesystem paths; counterfeit rules C1, C2, C4 require runtime side-effects. Per essay §3, V wraps CUE — these checks belong in V.

### Alt-3: V emits YAML instead of JSON

**Rejected.** AC7 requires "mechanically consumable by δ"; JSON has stable schema-validated tooling (`jq`, JSON Schema). YAML for V output adds parser-ambiguity surface. The receipt YAML and the V verdict JSON are *different surfaces* — YAML for human-authored receipts, JSON for machine-emitted verdicts.

## §Risks

### Risk: PyYAML not installed in target environment

**Mitigation:** V helper falls back to using `cue export --out json <receipt.yaml>` to coerce YAML → JSON via CUE itself; CUE is required for the cycle and is verified present at startup. PyYAML is preferred but not required.

### Risk: `cue vet` diagnostic format changes between cue versions

**Mitigation:** V captures the diagnostic verbatim into `failed_predicates[i].diagnostic`; V does not parse the diagnostic for structure. The `predicate` field uses V-internal predicate IDs (`cue.structural_vet`, `counterfeit.actor_separation`, etc.).

### Risk: counterfeit rules misfire on legitimate edge cases (e.g., a sole-maintainer repo where α and β actors must be the same in practice)

**Mitigation:** rules emit FAIL — they do not refuse to emit a verdict. An operator who needs to ship a sole-maintainer cycle uses `boundary_decision.action: override` per Q4; the override carries the failed_predicates, and the receipt closes degraded. This is correct per the design — sole-maintainer cycles *are* degraded relative to two-actor review.

### Risk: AC6 oracle (evidence files must exist) refuses on fixtures that point at placeholder paths (`@merge_sha_placeholder`)

**Mitigation:** counterfeit-fixture files include actual closure-record stub files under `schemas/cds/fixtures/_closure_records/<fixture-name>/`. The valid `schemas/cds/fixtures/valid-receipt.yaml` already points at `.cdd/releases/0.0.0/369/*.md` which exist on main; V on cycle/389 resolves these against the worktree root and they resolve successfully.
