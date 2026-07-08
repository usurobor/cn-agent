# β review — cnos#611

## §R0

### Scope check

`git diff main...cycle/611 --stat` touches only: `.cdd/unreleased/611/{CLAIM-REQUEST.yml,gamma-scaffold.md,self-coherence.md}` (process artifacts), `docs/guides/INSTALL-CDS.md`, `docs/guides/README.md`, `docs/guides/templates/cnos-install.yml` (new). No `src/go/` changes, no `.cn/` schema changes, no unrelated files touched. Confirmed by re-running the stat myself, not trusting the self-report.

### AC1 — delegation (Mock D3): **pass**

Re-ran `grep -n "deps.json\|deps.lock\|restore" docs/guides/templates/cnos-install.yml` myself:

```
83:        # Delegates entirely to the installer command — no deps.json /
84:        # deps.lock.json / vendor-restore logic is reimplemented here
```

Both hits are comment lines (`#`), not shell logic. Read the full 114-line file: the job has exactly one step that touches install semantics (`cn repo install`, lines 82–96), branching only on `install_dispatch` to add `--dispatch cds --agent ... --workflow-pat-secret ...` flags — no hand-rolled package/lock/vendor logic anywhere in the file. Confirmed independently.

### AC2 — base PR (Mock D1): **pass**

Read lines 47–68 and 98–101. The `else` branch of the "Resolve installing token" step (taken when `install_dispatch` is not `"true"`, i.e. the default `false`) sets `token=${{ secrets.GITHUB_TOKEN }}` (line 67) into `steps.token.outputs.token`, which is the only token threaded into `actions/checkout@v4` (line 72) and `peter-evans/create-pull-request@v6` (line 101). No custom secret is referenced on this branch. Confirmed independently, not just re-citing self-coherence's line numbers.

### AC3 — dispatch gating (Mock D2): **pass, with one non-blocking syntax-risk observation**

Structural check (read lines 46–72): the "Resolve installing token" step is the job's **first** step; `actions/checkout@v4` is the **second**. When `install_dispatch: true`, the step reads `secrets[inputs.workflow_pat_secret]` into `SECRET_VALUE`; if empty, it emits `::error::` (line 62) naming the unset secret and `exit 1` (line 63) — this happens before checkout runs, so no local checkout, no `cn repo install` render, and no PR-open step executes. This satisfies "never half-applied": the only externally-visible write in this workflow is the PR-open step, and every path that reaches it either used a resolved non-empty secret or already exited beforehand.

Residual (non-blocking) observation: the empty-string check only validates that the named secret is *set*, not that it actually carries `workflow` scope. A wrong-but-non-empty PAT would pass this guard and fail later, in the `peter-evans/create-pull-request@v6` push (GitHub server-side rejects `workflow`-scope file pushes from an insufficiently-scoped token) — still a clean failure (job goes red, no partial `.github/workflows/cnos-cds-dispatch.yml` gets merged since nothing is merged without a successful PR-open), just later in the pipeline than the explicit early check. This is consistent with the AC3 oracle as literally worded ("resolves a non-empty named secret and proceeds, or fails clearly") — scope-validation ahead of time isn't asked for and isn't cheaply achievable from within the workflow. Noting as an observation, not a finding.

**GitHub Actions syntax risk for `secrets[inputs.workflow_pat_secret]`:** I don't have live access to verify against GitHub's docs or a real Actions run. Reasoning from known GHA expression semantics: all contexts (`github`, `env`, `vars`, `secrets`, `matrix`, `inputs`, etc.) support bracket-indexing by an arbitrary expression — this is the same idiom widely used for matrix-driven dynamic secret selection (e.g. `secrets[format('{0}_TOKEN', matrix.env)]`), and secret-value masking in logs is applied by string-match against the job's known secret values, not by static enumeration of literal `secrets.NAME` references in the workflow source — so dynamic indexing does not break masking. I assess this syntax as **plausible and likely valid**, matching a known, if slightly under-documented, idiom. This is not a novel construct invented for this cell. `yamllint` (relaxed profile) passes with zero findings, and `python3 -c "import yaml; yaml.safe_load(...)"` parses cleanly (re-ran both myself). `actionlint` was not available in this environment to give a stronger static-analysis signal. Per α's own disclosed debt item 1, this has not been exercised against a live `workflow_dispatch` run. I flag this as an **observation warranting a live-fire follow-up** (e.g. a canary dispatch run against a scratch repo) before this template is held out as fully proven — not a blocking finding, since the syntax risk is low-probability and already disclosed as debt.

### AC4 — docs flip: **pass**

Read the full `docs/guides/INSTALL-CDS.md` (256 lines) post-edit:

- Existing content preserved: Layers 1/2 framing, Prerequisites, Install steps 1–3, Autonomous dispatch section, Verifying, Troubleshooting table (all pre-existing rows intact), Related section (all pre-existing links intact) — confirmed via `git diff main...cycle/611 -- docs/guides/INSTALL-CDS.md`: the diff is purely additive except for one Uninstalling-section sentence (adds a note about removing `cnos-install.yml`) and one new Troubleshooting row. No existing line was deleted or contradicted.
- New `## GitHub UI (no-terminal) install` section (lines 168–202) explicitly states "This is a **thin wrapper**, not a second install path... install logic is not duplicated in YAML," names the `release`/`install_dispatch`/`workflow_pat_secret`/`agent` inputs table, and states the PAT precondition and fail-clear behavior for dispatch mode — all claims trace to lines actually present in the YAML (verified above).
- New `## Tenant secrets, by tier` table (lines 204–218) names Tier 1 (`GITHUB_TOKEN` only), Tier 2 (mechanical FSM, `GITHUB_TOKEN` only, explicitly labeled "Forthcoming — tracked in cnos#613"), Tier 3 (`workflow`-PAT + `CLAUDE_CODE_OAUTH_TOKEN`). It names cnos#613 as a **tracker**, with no markdown link to a file path — confirmed by reading the raw text; there is no `[...](...)`-style link to any #613 runbook file. `CLAUDE_CODE_OAUTH_TOKEN` is a real, established secret name in this repo (grepped: used in `.github/workflows/cnos-cds-dispatch.yml`, `.github/workflows/cnos-agent-admin.yml`, and the golden templates under `src/packages/`), so the table's claims are grounded, not fabricated.
- `docs/guides/README.md` diff (re-read in full): adds `- [INSTALL-CDS.md](INSTALL-CDS.md) — install CDS into a repo (CLI or GitHub UI).` under `## Operator`, as the first entry. Confirmed the file did not link `INSTALL-CDS.md` before this cycle (only `HANDSHAKE.md`/`AUTOMATION.md` were listed).
- Link resolution, verified against the actual filesystem (not assumed): `docs/guides/README.md` → `INSTALL-CDS.md` resolves to `docs/guides/INSTALL-CDS.md` (exists). `docs/guides/INSTALL-CDS.md` → `templates/cnos-install.yml` (two occurrences, lines 172 and 255) resolves to `docs/guides/templates/cnos-install.yml` (exists, new file in this diff). Both link targets are real, present files at the stated relative paths.
- Repo-wide grep for other README/index files referencing `INSTALL-CDS` (`grep -rn "INSTALL-CDS" --include="*.md" .`) finds only `docs/guides/README.md`'s new line and historical `.cdd/unreleased/{608,610}/` process artifacts — no other nav surface was missed or needs updating for this cell's stated scope.

### YAML syntax / parse check

`python3 -c "import yaml; yaml.safe_load(open('docs/guides/templates/cnos-install.yml'))"` — parses cleanly, no error. `yamllint -d "{extends: relaxed, rules: {line-length: disable}}"` — zero findings (silent, exit 0). Both re-run independently, matching self-coherence's claim.

### Friction-note resolutions — independent judgment

1. **Creating rather than rewriting the template.** Confirmed independently: `git show main:docs/guides/templates/cnos-install.yml` → `fatal: path ... does not exist in 'main'`. The issue body's "Exists" framing was simply stale/wrong. This does not change the deliverable (a template that delegates to `cn repo install`) and doesn't warrant escalation — correcting a factual premise before implementation, with the correction pinned in `gamma-scaffold.md` ahead of any code being written, is exactly the right discipline. **No escalation was warranted; this was the correct call.**
2. **Inlining tenant-secrets content instead of linking a nonexistent #613 runbook.** Confirmed independently: `gh issue view 613` → `state: OPEN`, not merged, no runbook file exists in the repo to link. The three options were: (a) block #611 on an unlanded sibling cell for a soft/informational dependency, (b) fabricate a link to a file that doesn't exist, or (c) inline the currently-knowable facts and name #613 as the forward tracker, disclosed as debt. Option (c) is the right engineering call — it satisfies the operator's actual ask (a tenant can see what to provision today) without a false link or an unnecessary block, and the duplication is explicitly flagged for reconciliation once #613 lands. **No escalation was warranted; this was the correct call**, and is consistent with this repo's own established precedent (cnos#608's self-coherence dropping a stale reference rather than keeping a broken one).

## Findings

None that block convergence.

## Observations (non-blocking)

1. `secrets[inputs.workflow_pat_secret]` (dynamic secret-name indexing) is plausible, idiomatic GitHub Actions syntax by my independent reasoning, but has not been exercised against a live `workflow_dispatch` run (disclosed by α as debt item 1). Recommend a live canary dispatch (e.g. against a scratch repo) in a follow-up before holding this template out as fully proven, though I assess the risk of it being outright rejected by the parser as low.
2. AC3's fail-clear guard checks only that the named secret is non-empty, not that it carries `workflow` scope; an insufficiently-scoped-but-set PAT would still fail cleanly, just later (at the PR-push step) rather than at the earliest possible point. Matches the AC3 oracle as literally worded; not a gap worth blocking on.
3. The Tenant-secrets-by-tier table duplicates content #613's own runbook will eventually own (disclosed by α as debt item 2) — reconcile-or-link when #613 lands.

## verdict: converge

Justification: all four ACs re-derived independently with fresh greps/reads (not trusting self-coherence.md's claims) and all pass. No unrelated files touched, no `src/go/` changes, YAML parses and lints clean, both friction-note resolutions are sound engineering calls that did not warrant blocking or escalation, and every doc link resolves to a real file. The two open items (dynamic-secret live-fire verification, #613 content duplication) are already disclosed as debt by α and are appropriately deferred rather than blocking this cell.
