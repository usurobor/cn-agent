# γ scaffold — cycle #609

**Issue:** #609 — cds-install Sub 2: generalize `cn install wake` identity (agent/PAT/bot) — successor to #549
**Mode:** substantial (multi-file renderer change + golden/test additions; not a small-change or immediate-output cell)
**Base SHA:** `612a96d6a5cac69486a48b50ca8432daf3bdaac2` (origin/main at scaffold time)
**Branch:** `cycle/609`
**Wake-invoked-δ:** claimed by `cds-dispatch` wake firing; run_class `first_pass` (per `.cdd/unreleased/609/CLAIM-REQUEST.yml`, committed to main pre-branch per the `cnos#575` claim-request convention).

## Source of truth

| Surface | Source | Status |
|---|---|---|
| Renderer under change | `src/packages/cnos.core/commands/install-wake/cn-install-wake` | On main |
| Golden fixtures + CI gate | `.github/workflows/install-wake-golden.yml`, `src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml`, `src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml`, live `.github/workflows/cnos-cds-dispatch.yml` | On main |
| **Canonical design surface (binding oracle)** | `docs/development/design/cn-repo-install-MOCKS.md` §Mock C (lines ~104-152) + §Mock E (lines ~179-206) | On main (landed via cycle/608, merged as part of #617) |
| Issue AC1-AC4 (original scope) | issue #609 body | — |
| Scope-extension comments | issue #609 comments (tenant finalizer E2/E4; full-identity-set framing; final operator launch directive) | — |

**Binding-oracle resolution (δ decision, logged per delta/SKILL.md §2 inward-membrane duty):** the issue body's AC1 oracle command (`--agent acme --workflow-pat-secret ACME_WORKFLOW_PAT --bot-name acme-bot --bot-id 12345678`) is *byte-identical* to Mock C's console example. Mock C's invariant table (C1-C6) and Mock E's invariant table (E1-E4) are the more precise, machine-checkable statement of the same requirement the issue's prose comments gesture at. **This cycle implements exactly Mock C + Mock E's renderer-facing invariants (C2, C4, C5, C3 partially — C3 is `cn repo install`'s concern per #610, not this renderer's; C1/C6 are `cn repo install`'s concern, out of scope here) plus E2/E4.** The later issue-comment's `--git-user-name`/`--git-user-email` framing does **not** appear in Mock C or Mock E (the pinned design surface) and is **not** in this cycle's scope — see Friction note 1 below.

## Scope (restated from issue, cross-checked against Mock C/E)

**In:** `--workflow-pat-secret`, `--bot-name`, `--bot-id` flags on `cn-install-wake` (i.e. `cn install-wake`); the agent name already interpolates the concurrency group via existing `{agent}` substitution — this cycle extends the same substitution discipline to the one self-referential literal that was missed (see AC2 friction below); a leak-audit test (C4); tenant-portable finalizer/scanner acquisition (E2/E4); sigma defaults remain available but are not the only path (C5/E4 byte-identical backward compat).

**Out (per issue + confirmed by Mock C/E's own scope split):** `cn repo install` / repo-installer wiring (Mock A/B/C1/C3/C6/D — that is issue #610's job, which *consumes* this cycle's new flags); wake behavior/prompt/schema changes beyond the one mechanical self-reference fix named below; new triggers; `--git-user-name`/`--git-user-email` (not in the pinned Mock C/E design; see Friction note 1).

## AC oracle approach (per-AC, mapped to Mock C/E invariants)

| AC | Oracle | Implementation surface |
|---|---|---|
| AC1 (config identity) | `cn-install-wake cds-dispatch --agent acme --workflow-pat-secret ACME_WORKFLOW_PAT --bot-name acme-bot --bot-id 12345678 --out <tmp>` → output contains `ACME_WORKFLOW_PAT`, `acme-bot`, `12345678`, `cds-dispatch-acme`. Negative (Mock C2): omit `--workflow-pat-secret`/`--bot-name`/`--bot-id` for `--agent acme` → renderer dies **before** writing `--out` (no partial file), non-zero exit, message names the missing flag(s). | `cn-install-wake` arg parsing + identity resolution (§"Resolve substrate bindings") |
| AC2 (Mock C4 — zero sigma leak) | Two-layer test: (a) a **synthetic fixture manifest** (mirrors the existing AC5/AC2-negative fixture pattern already in `install-wake-golden.yml`) rendered with a non-sigma agent + full identity flags → `grep -iE 'sigma\|SIGMA_WORKFLOW_PAT\|41898282'` returns nothing; this isolates the renderer's own leak-freedom from any given package's prompt-body prose. (b) a **scoped** grep against the real `cds-dispatch` acme-render for the renderer-controlled tokens only (`SIGMA_WORKFLOW_PAT`, literal `bot_name: "sigma`, literal `bot_id: "41898282"`, literal `cds-dispatch-sigma` group line) → zero matches, confirming the fix in the next row closes the one real self-reference gap. See Friction note 2 for why an unconditional `grep -i sigma` against the *live* `cds-dispatch` wake's full prompt body is not used as the oracle. | new CI step in `install-wake-golden.yml`; one self-reference fix in `cnos.cds/orchestrators/cds-dispatch/SKILL.md` |
| AC3 (Mock C5 / E4 — sigma byte-identical) | Existing "Verify goldens unchanged" + idempotence steps in `install-wake-golden.yml` already re-render with default (`--agent sigma`, no new flags) and diff against the committed golden. No new step needed — this is the regression oracle: if the sigma default path renders anything different, that step already fails the build. | `cn-install-wake` default-path preservation (new flags only take effect when explicitly passed or when agent ≠ sigma) |
| AC4 (gates green) | `install-wake-golden`, Go (`cd src/go && go build ./... && go test ./...`), Package gates all green before closeout. | CI + local verification |
| E2 (tenant: no `cd src/go && go build`) | Render with `--agent acme` (+ required identity flags) → output contains **no** `cd src/go` / `go build ./cmd/cn`; contains an `install.sh`-based acquisition step instead, invoked via a fixed path. | `cn-install-wake` "Mechanical recovery scanner" + "Mechanical checkpoint + PR finalizer" step bodies, gated on `$agent != sigma` |
| E4 (sigma still self-builds) | Sigma-default render still contains `(cd src/go && go build -o /tmp/cn-scan ./cmd/cn)` / `-o /tmp/cn-finalize` exactly as today (covered by the same byte-identical check as AC3). | same gate as above, `sigma` branch unchanged |

## Expected diff scope

- `src/packages/cnos.core/commands/install-wake/cn-install-wake` — new flags, identity resolution refactor, tenant-portable acquisition branch. No change to argument parsing for `--agent`, `--out`, `--manifest`, `--activation-state-override`.
- `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` — one-line self-reference fix (`cds-dispatch-sigma` → `cds-dispatch-{agent}` in the disallowed-surfaces bullet). This is the wake's own prompt body; γ names this narrowly-scoped exception explicitly because it is a mechanical parameterization-consistency fix (the same value is already `{agent}`-substituted three other times in the same file), not a behavioral/schema change, and because AC2 cannot pass against the live wake without it.
- `.github/workflows/install-wake-golden.yml` — new CI steps for AC1 positive/negative, AC2 (synthetic fixture + scoped live-wake grep), E2/E4.
- No change to `cnos-cds-dispatch.golden.yml`, `cnos-agent-admin.golden.yml`, or the live `.github/workflows/cnos-cds-dispatch.yml` — the sigma-default render must stay byte-identical (verified mechanically by existing steps, not by manual inspection).
- No change to `src/go/**`, `docs/development/design/cn-repo-install-MOCKS.md`, or issue #610's surfaces.

## Friction notes

1. **`--git-user-name`/`--git-user-email` not implemented.** A later issue comment ("PR 2 — enacting the operator launch cut") asks the renderer to also expose "the git commit identity (`--git-user-name`, `--git-user-email`)". The pinned design surface (`cn-repo-install-MOCKS.md` Mock C, which the issue's own "Source of truth" table names as authoritative and which the final "Launch scoping" comment reaffirms as the pinned deliverable) does **not** include these flags in its console example or invariant table (C1-C6), and no existing renderer output emits a `git config user.*` step for sigma today (git commit authorship in this repo's dispatch cycles is configured by the substrate/harness outside this renderer's YAML, not by a step this renderer emits) — there is no existing behavior to parameterize. Adding two new flags plus a wholly new conditionally-rendered workflow step would be a real behavior addition beyond Mock C/E's pinned shape, and the issue's own scope line explicitly lists "wake behavior/prompt/schema changes" as **Out**. δ resolves this ambiguity by following the pinned Mock C/E surface (which is more specific and more recently landed than the prose comment) and **not** adding these flags in this cycle; if the operator wants them, that is a follow-up sub, not a silent scope-creep item here.
2. **AC2's oracle is split into a synthetic-fixture check + a scoped live-wake check, not one unconditional `grep -i sigma` against the live `cds-dispatch` wake.** The live wake's prompt body (package-owned prose per `cnos.cds/orchestrators/cds-dispatch/SKILL.md`) contains two literal, out-of-scope "sigma" substrings: (a) `agent-admin-sigma` — the **admin** wake's own concurrency-group name, a different wake entirely, not a parameter of *this* wake's `--agent` flag; (b) a dated historical citation (`cn-sigma:.cn-sigma/logs/20260624.md`) naming a real past incident. Neither is a renderer-authority leak; both are package-authority prose the issue's scope line places out of bounds ("wake behavior/prompt/schema changes" — Out). An unconditional case-insensitive `grep sigma` against the live wake's full render would therefore never reach zero, regardless of how the renderer's own bindings are parameterized. The two-layer oracle in the AC2 row above tests the renderer's actual leak-freedom (synthetic fixture, fully controlled) and the one genuine self-reference bug this cycle fixes (scoped grep on the renderer-controlled tokens), without being defeated by unrelated legacy prose it is not this cycle's job to rewrite.
3. **Command naming.** A later comment says "the installer is `cn repo install`; this renderer command stays `cn install wake`." The existing renderer is invoked as `cn install-wake <wake-name>` (single hyphenated subcommand via the package-command discovery shim, per the script's own usage header) — there is no two-word `cn install wake` subcommand distinct from it. Read as prose shorthand for the existing `cn install-wake` command (which the comment itself confirms is unaffected — "stays"), not as a rename directive. No CLI renaming performed.

## Scope guardrails (binding on α)

- Do not touch `src/go/**`, `cn repo install`, or any file under `docs/development/design/cn-repo-install-MOCKS.md` (read-only source of truth).
- Do not touch `cnos-cds-dispatch.golden.yml`, `cnos-agent-admin.golden.yml`, or `.github/workflows/cnos-cds-dispatch.yml` — if any of these three files show a diff after your change with `--agent sigma`/default, that is a regression, not a feature.
- The ONLY line to change in `cnos.cds/orchestrators/cds-dispatch/SKILL.md` is the one self-reference named in Friction note 2(a)'s companion fix row (AC2). Do not touch any other line of that file.
- New flags must be optional and must not change output when omitted for `--agent sigma` (default).
- Fail-early discipline (Mock C2): identity-resolution failures for a non-sigma agent must `die` (exit 1, per the script's existing `die()` helper) **before** the render/write phase — no partial `--out` file.

## Empirical anchor

cycle/608 (just-merged, `cn repo install` base installer) is the sibling cycle this one's flags will be consumed by (issue #610). No direct code dependency — #608 does not touch `cn-install-wake` — but the design doc it landed (`cn-repo-install-MOCKS.md`) is this cycle's binding oracle, landed on `main` before this cycle's base SHA.
