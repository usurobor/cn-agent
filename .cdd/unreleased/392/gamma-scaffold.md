# γ scaffold — cycle/392

## Issue
[cnos#392](https://github.com/usurobor/cnos/issues/392) — Phase 3 remediation v2: port V to Go as `cn cdd-verify` subcommand in cnos.cdd package (#391 superseded).

## Parent / Predecessors
- Parent: cnos#366 (Phase 3 remediation; gates Phase 4)
- Supersedes: cnos#391 (closed as rescoped after operator caught two architectural under-specifications)
- Source of design: cnos#389 (Python implementation of V; merged at 993d7f93)

## Mode
design-and-build, **γ+α+β-collapsed-on-δ**. β-α-collapse acknowledged. δ-as-architect pinned the 7-axis implementation contract at dispatch (see issue body).

## Surfaces
1. **Go source under** `src/packages/cnos.cdd/commands/cdd-verify/` — V implementation (parser, dispatcher, validator, counterfeit rules).
2. **Go workspace** `go.work` at repo root — brings the new module into the `src/go` build graph so the cn binary can import the V package.
3. **`cn cdd-verify` kernel command** in `src/go/internal/cli/` — thin wrapper that delegates to the V package (per eng/go §2.18 "cli/ owns dispatch only").
4. **`cn` main.go** registration — register the new kernel command.
5. **Python `cn-cdd-validate-receipt`** — removed.
6. **Bash `cn-cdd-verify`** — removed (or thin shim deferring to `cn cdd-verify`). Operator preference: prefer full removal. We'll fully remove and update the cnos.cdd `cn.package.json` accordingly. The kernel-tier `cdd-verify` command supersedes the package-vendored one.
7. **`cn.package.json` for cnos.cdd** — remove the `cdd-verify` entry (now kernel; remains for `cdd-status`).
8. **`tests/cdd/test_cn_cdd_validate_receipt.sh`** — keep invocation through `cn cdd-verify` (set `V` to `cn` binary + `cdd-verify`).
9. **Go unit tests** — table-driven tests for the V package alongside the source.
10. **CHANGELOG** — line in 3.81.0 (or unreleased section).

## AC oracle approach
- **AC1 — Go at cnos.cdd path; Python+bash removed.** Mechanical: `ls src/packages/cnos.cdd/commands/cdd-verify/*.go` returns ≥1; `test ! -f .../cn-cdd-validate-receipt`; `test ! -f .../cn-cdd-verify`. Verified post-build.
- **AC2 — `cn cdd-verify` routes through Go.** Build `cn`; run `cn cdd-verify --help`; confirm output comes from Go (validator_identity present in JSON output for `--receipt` path).
- **AC3 — 37/37 oracle suite passes.** Run `bash tests/cdd/test_cn_cdd_validate_receipt.sh` after build; expect 37 passed, 0 failed (matching Python baseline).
- **AC4 — 5/5 counterfeit fixtures rejected with named diagnostics.** Test suite already covers C1 (AC8.a-b), C2 (AC8.c-d), C3 (AC8.e-f), C5 (AC2.d-e), and C4 (AC6.b-c). Mechanical via test suite. Verify each named diagnostic appears.
- **AC5 — JSON schema unchanged.** `git diff main..cycle/392 -- schemas/cdd/validation_verdict.schema.json` empty. Verify Go emits validating JSON: `cn cdd-verify --receipt ... --json | jq .` parseable; required keys present.
- **AC6 — eng/go SKILL compliance.** Walk every section of eng/go and cite file:line evidence in self-coherence.md.
- **AC7 — Backward compat.** Run `cn cdd-verify --unreleased`, `--all`, `--version`, `--pr`, `--receipt` — verify each behaves equivalent to Python predecessor.
- **AC8 — cnos build works.** `cd src/go && go build ./...` succeeds; `go test ./...` passes; `cn cdd-verify --help` works.

## Expected diff scope
- **New** `go.work` (workspace file).
- **New** `src/packages/cnos.cdd/commands/cdd-verify/go.mod` + `*.go` source files (validator, dispatch table, counterfeit rules, JSON emit, tests).
- **New** `src/go/internal/cli/cmd_cdd_verify.go` — thin kernel-command wrapper.
- **Edit** `src/go/cmd/cn/main.go` — register the new command.
- **Edit** `src/go/go.mod` — add replace directive or workspace include (via go.work).
- **Delete** `src/packages/cnos.cdd/commands/cdd-verify/cn-cdd-validate-receipt` (Python).
- **Delete** `src/packages/cnos.cdd/commands/cdd-verify/cn-cdd-verify` (bash) — operator preference.
- **Edit** `src/packages/cnos.cdd/cn.package.json` — remove `cdd-verify` entry (kernel-tier now owns it).
- **Edit** `src/packages/cnos.cdd/commands/cdd-verify/README.md` — update to reflect Go implementation.
- **Edit** `CHANGELOG.md` — line under unreleased.
- **Cycle artifacts**: γ scaffold (this), design-notes, self-coherence, β-review, α/β/γ close-outs, cdd-iteration, INDEX row.

## CDD trace
- cnos#366 (Phase 3 parent) → cnos#389 (V Python) → cnos#391 (Go port — rescoped) → **cnos#392** (Go port v2 with pinned contract).
- δ-as-architect role exercised at dispatch (issue body's implementation-contract table).
- δ-as-agent (this cycle) executes against the pinned contract.
- ε will fold F1–F4 into a coordinated cdd patch (alpha + beta + gamma skill updates + delta/Phase 4 scope refinement).

## Skills
Tier 3:
- `cnos.eng/skills/eng/go/SKILL.md` — **MANDATORY, BINDING**
- `cnos.eng/skills/eng/code/SKILL.md`
- `cnos.eng/skills/eng/test/SKILL.md`
- `cnos.cdd/skills/cdd/design/SKILL.md`
- `cnos.cdd/skills/cdd/issue/proof/SKILL.md`

## Branch
- Created: `cycle/392` from `origin/main` (HEAD `5ec589a1`)
- Pushed to origin: yes
