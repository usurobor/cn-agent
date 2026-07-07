cell_kind: docs

# γ scaffold — cnos#611: bootstrap workflow delegates to `cn repo install` + docs flip

## Source-of-truth pin

- Base SHA (`main`): `ca52fc5081027df185e3a8a5191e621d115cfb66` (verified against the working tree at cycle-branch creation — `git rev-parse HEAD` on `cycle/611` matches this SHA at R0).
- Design surface: `docs/development/design/cn-repo-install-MOCKS.md` §Mock D (GitHub UI / no-terminal path), invariants D1–D4.
- Current guide: `docs/guides/INSTALL-CDS.md` (already names `cn repo install` / `cn repo install --dispatch cds` as canonical CLI path; already sigma-generalized per #610 — no sigma caveat remains to remove).
- Current nav: `docs/guides/README.md` (does not yet link `INSTALL-CDS.md` — confirmed by reading the file; this is #608's debt item 4, explicitly deferred to #611).
- CLI reference: `src/go/internal/cli/cmd_repo_install.go` (`cn repo install` flags: `--release`, `--index`, `--packages`, `--dispatch {none,cds}`, `--agent`, `--workflow-pat-secret`, `--bot-name`, `--bot-id`, `--dry-run`; never writes `.github/workflows/` in base mode; never pushes to `main`; exits nonzero naming cnos#493 for `--dispatch cds` until canonical labels ship).
- Scope-extension comment (operator, cnos#606 C6): tenant secrets must be named per tier in `INSTALL-CDS.md`, linking the runbook Sub 6 (#613) produces.

## Friction note — issue body vs actual repo state (resolved before dispatch)

The issue body's "Exists" framing ("the bootstrap `docs/guides/templates/cnos-install.yml` reimplements the install steps in shell") is **not accurate against `main`**: this file does not exist anywhere on `main` (confirmed: `find . -iname cnos-install.yml` → no hits; also independently confirmed by #608's own `self-coherence.md` §Docs-reconciliation, which found the same and explicitly deferred creating it to this issue). The real gap is **absence**, not duplication — #611's job is to **create** the template (Mock D), not rewrite an existing reimplementation. This scaffold pins the corrected framing; α implements against "create," not "rewrite."

## Friction note — cnos#613 dependency (scope-extension comment)

The operator's scope-extension comment demands `INSTALL-CDS.md` link "the tenant secrets runbook produced by Sub 6 (#613)." #613 (`cds-install Sub 6: PAT-free mechanical FSM-engine wake + tenant secrets runbook`) is **still `status:open`, not yet dispatched** — no runbook file exists in the repo to link. Linking to a nonexistent path would be a false claim (the same discipline #608's self-coherence applied when it dropped a stale Track-B doc reference). Resolution pinned for α: **write the per-tier secrets table inline** in `INSTALL-CDS.md` now (the content is fully specifiable today: `GITHUB_TOKEN` for the mechanical/label-driven tier, a `workflow`-scope PAT + `CLAUDE_CODE_OAUTH_TOKEN` for the autonomous-dispatch tier — all already established facts from the existing `cds-dispatch` orchestrator and `cmd_repo_install.go`), and name #613 as the tracker that will supply the fuller canonical runbook — not a broken link to it. This is disclosed as debt (§Debt in `self-coherence.md`), not silently patched over.

## Implementation contract (δ-pinned; α MUST NOT improvise)

| Axis | Pin |
|---|---|
| Language | YAML (GitHub Actions workflow) + Markdown (docs). No new Go code. |
| CLI integration target | N/A — this cell consumes the existing `cn repo install` CLI surface (cnos#608/#609/#610); it does not modify CLI source. |
| Package scoping | `docs/guides/templates/cnos-install.yml` (new file), `docs/guides/INSTALL-CDS.md` (edit), `docs/guides/README.md` (edit). No `src/go/` changes. |
| Existing-binary disposition | N/A — no binary changes. |
| Runtime dependencies | The rendered template's own runtime deps: `actions/checkout@v4`, `peter-evans/create-pull-request@v6` (standard, widely-used PR-opening action; no existing precedent for PR-opening in this repo's `.github/workflows/`, but `cn repo install` itself never opens a PR — something must, and this is the standard idiom); `curl` + `install.sh` (already the canonical `cn` acquisition path per `INSTALL-CDS.md` and Mock E's tenant-portable precedent). |
| JSON/wire contract | None — this cell does not touch `.cn/deps.json` / `.cn/deps.lock.json` schemas. |
| Backward compat | `docs/guides/INSTALL-CDS.md`'s existing CLI-path content (Layer 1 / Layer 2 sections, troubleshooting table) MUST NOT regress; this cell adds a GitHub UI wrapper section and a tenant-secrets table, it does not rewrite the CLI-path sections. |

## Per-AC oracle list (mirrors the issue body's AC1–AC4 verbatim)

- **AC1 — delegation (Mock D3).** Oracle: `docs/guides/templates/cnos-install.yml`'s job body invokes `cn repo install …`; `grep -n "deps.json\|deps.lock\|restore" docs/guides/templates/cnos-install.yml` shows no inline reimplementation of install logic (only doc comments referencing those filenames are permitted, not shell logic that writes them).
- **AC2 — base PR (Mock D1).** Oracle: with `install_dispatch: false` (default), the job step that opens the PR runs with `${{ secrets.GITHUB_TOKEN }}` (not a named custom secret) — grep-verifiable in the rendered file's `token:` binding for the default-path branch.
- **AC3 — dispatch gating (Mock D2).** Oracle: with `install_dispatch: true`, either (a) the named `workflow_pat_secret` resolves to a non-empty value and the job proceeds using that token for checkout + PR-open (never `GITHUB_TOKEN` in this branch), or (b) the secret is unset/empty and the job fails with an explicit `::error::` naming the missing secret, before any file write — never a half-applied `.github/workflows/cnos-cds-dispatch.yml`.
- **AC4 — docs flip.** Oracle: `docs/guides/INSTALL-CDS.md` names `cn repo install` as the canonical **local** path (already true — preserved); a new section documents the GitHub UI workflow as a **thin wrapper** over the same command (not reimplemented logic); the PAT caveat for dispatch mode is stated; the tenant-secrets-per-tier table is present and does not link a nonexistent #613 artifact (names #613 as the tracker, per the friction note above); `docs/guides/README.md` links `INSTALL-CDS.md`.

## Parity requirement

Closeout MUST carry `mock_parity` rows for **D1–D4** with `missed: 0`, per the issue body's Parity requirement section.

## Non-goals (from the issue body, preserved)

No push to `main`; no reimplemented install logic in YAML; no autonomous-write default; no #613 runbook file fabricated ahead of its own cycle; no CLI (`src/go/`) changes.

## α prompt (self-dispatch under wake-invoked δ; single-session δ plays α next)

Implement per the Implementation contract and AC oracle list above:
1. Create `docs/guides/templates/cnos-install.yml` — a `workflow_dispatch`-triggered template (copied by tenants into their own `.github/workflows/`), with inputs `release` (default `latest`), `install_dispatch` (boolean, default `false`), `workflow_pat_secret` (string, default `CNOS_WORKFLOW_PAT`), `agent` (string, default `cnos`). The job installs `cn` via `install.sh` (tenant-portable, no `go build`), calls `cn repo install` with the resolved flags, and opens a PR with `peter-evans/create-pull-request@v6`, using `GITHUB_TOKEN` in base mode and the named secret (dynamic lookup via `secrets[inputs.workflow_pat_secret]`) in dispatch mode — failing clearly (before any checkout/write) if that secret is unset when `install_dispatch: true`.
2. Update `docs/guides/INSTALL-CDS.md`: add a "GitHub UI (no-terminal) install" section under Layer 1/Layer 2 describing the template as a thin wrapper (Track B), add the tenant-secrets-per-tier table (per the #613 friction note resolution), keep all existing CLI-path content intact.
3. Update `docs/guides/README.md`: add an `INSTALL-CDS.md` link under the Operator section.
4. Write `self-coherence.md` with per-AC evidence, the two friction notes surfaced above, and a `Debt` section naming: (a) the dynamic-secret-lookup GH Actions syntax is lint-checked but not exercised against a live `workflow_dispatch` run in this cycle (no live-fire harness available); (b) the tenant-secrets table duplicates content that #613's runbook will eventually own — reconcile-or-link when #613 lands, don't duplicate forever.

## β prompt

Independently walk AC1–AC4 against the diff on `cycle/611`. Re-derive (do not trust α's self-report): confirm the template file's job body calls `cn repo install` with no inline reimplementation; confirm the base-mode token binding is `secrets.GITHUB_TOKEN`; confirm the dispatch-mode branch fails clearly on an unset secret and never partially writes; confirm the docs diff preserves existing CLI-path content and adds the required sections without a broken #613 link. Render the YAML through `yamllint` (or equivalent) as an independent syntax check. Verdict: converge or iterate with findings.

## Scope guardrails

In scope: the three files named above. Out of scope: `src/go/` CLI changes, `.cn/` schema changes, #613's actual runbook content, any `.github/workflows/cnos-cds-dispatch.yml` rendering changes (that's #609/#610's landed scope).
