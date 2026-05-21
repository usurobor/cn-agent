# CDD Verify (`cn cdd-verify`)

V (Contract × Receipt → ValidationVerdict) + the legacy CDD artifact-presence ledger, implemented in Go and exposed as the `cn cdd-verify` kernel subcommand.

This directory ports cnos#389's Python design source to Go (cnos#392). The previous `cn-cdd-validate-receipt` (Python) and `cn-cdd-verify` (bash) scripts that lived here are removed; the Go source compiles into the cn binary via the repo-root `go.work` workspace.

## Invocation

Canonical noun-verb form (matches `cn` dispatcher convention):

```bash
# V receipt validation
cn cdd verify --receipt schemas/cds/fixtures/valid-receipt.yaml
cn cdd verify --receipt path/to/receipt.yaml --json
cn cdd verify --receipt path/to/receipt.yaml --structural-only

# Legacy artifact-presence ledger
cn cdd verify --all                          # entire repository CDD ledger
cn cdd verify --unreleased                   # only .cdd/unreleased/ cycles
cn cdd verify --version 3.81.0               # one release
cn cdd verify --version 3.81.0 --cycle 392   # one cycle
cn cdd verify --pr 392                       # PR-scoped cycle
```

The hyphenated form `cn cdd-verify` is intercepted as a group listing by the `cn` dispatcher (per `src/go/internal/cli/dispatch.go` lines 54–69); the canonical form is `cn cdd verify`.

## Modes

| Flag | Effect |
|---|---|
| `--receipt <path>` | V verdict (Contract × Receipt → ValidationVerdict). Exit 0 PASS, 1 FAIL, 2 V-error. |
| `--receipt <path> --json` | Emit ValidationVerdict as JSON conforming to `schemas/cdd/validation_verdict.schema.json`. |
| `--receipt <path> --structural-only` | Skip filesystem-dereference rules (C1/C2/C4). |
| `--receipt <path> --contract <p>` | Optional contract reference recorded in provenance. |
| `--all` | Repository-wide ledger (unreleased + released). |
| `--unreleased` | Only `.cdd/unreleased/` cycles. |
| `--version X` | One release. |
| `--version X --cycle N` | One cycle (current artifact layout). |
| `--version X --triadic` | Include legacy `.cdd/` close-out checks. |
| `--pr N` | PR-scoped cycle. |
| `--bundle <path>` | Override PRA bundle dir (default `docs/gamma/cdd`). |
| `--exceptions <path>` | Load YAML exceptions for missing-artifact warnings. |
| `--repo-root <path>` | Override repo root detection (used in tests). |

## Counterfeit-receipt rules (V receipt mode)

| Rule | Predicate | Scope | Trigger |
|---|---|---|---|
| C1 | `counterfeit.actor_separation` | CDS | α and β closeouts authored by same actor (without `mode: collapsed`). |
| C2 | `counterfeit.verdict_precedes_merge` | CDS | β closeout commit time postdates `boundary_decision.decided_at` (action ∈ {accept, release}, without `mode: collapsed`). |
| C3 | `counterfeit.override_does_not_rewrite` | all | `boundary_decision.override.original_validation_verdict.verdict` ≠ `validation.verdict`. |
| C4 | `counterfeit.evidence_ref_unresolved` | all | `evidence_refs.*` path-like values don't resolve under repo root. |
| C5 | `counterfeit.protocol_id_mismatch` | CDS / CDR | `protocol_id` declares CDS/CDR but required keys missing in `evidence_refs`. |

Plus structural validation via `cue vet -c -d <#Definition> <schemas/...>` for the declared `protocol_id`.

## Source layout

```
cddverify.go    — top-level Validate orchestrator + emit_verdict
verdict.go      — ValidationVerdict struct types (matches JSON schema)
dispatch.go     — protocol_id → schema-package dispatch table
counterfeit.go  — C1–C5 counterfeit rules
parse.go        — Receipt type, ParseReceiptJSON (pure), ReadReceipt (IO wrapper)
cuevet.go       — adapter: `cue vet` subprocess
git.go          — adapter: `git log` / `git rev-parse` subprocess
ledger.go       — legacy artifact-presence ledger (--all/--unreleased/etc.)
run.go          — argv parsing, Run dispatcher (V mode vs ledger mode)
cddverify_test.go — table-driven unit tests
```

## Testing

```bash
go test ./src/packages/cnos.cdd/commands/cdd-verify/...
bash tests/cdd/test_cn_cdd_validate_receipt.sh   # 37-check AC oracle
```

## Exit codes

| Code | Meaning |
|---|---|
| `0` | PASS — all predicates held / all required artifacts present (warnings allowed). |
| `1` | FAIL — at least one predicate failed / at least one required artifact missing. |
| `2` | V itself errored (receipt not found, `cue` missing, internal error). |

## Migrating from the predecessors

| Predecessor | Disposition |
|---|---|
| `cn-cdd-validate-receipt` (Python) | Removed. All semantics ported to Go. |
| `cn-cdd-verify` (bash) | Removed. Operator surface now `cn cdd verify`. |

## References

- `src/packages/cnos.cdd/skills/cdd/CDD.md` §5.3a — Artifact Location Matrix
- `src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md` — V's validation contract
- `src/packages/cnos.eng/skills/eng/go/SKILL.md` — binding doctrine for this Go code
- `schemas/cdd/validation_verdict.schema.json` — stable JSON-schema contract
- `docs/gamma/essays/CCNF-AND-TYPED-TRUST.md` §3 — CUE/V split design source
- cnos#389 — Python predecessor (merged at 993d7f93)
- cnos#391 — closed as rescoped (this issue supersedes)
- cnos#392 — this cycle
