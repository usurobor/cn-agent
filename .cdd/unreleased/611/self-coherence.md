# self-coherence — cnos#611

## §R0

### Mode

design-and-build (per the issue body's own `Mode:` field), scoped to three files: `docs/guides/templates/cnos-install.yml` (new), `docs/guides/INSTALL-CDS.md` (edit), `docs/guides/README.md` (edit).

### Skills

Active skills for this cycle: `cdd/gamma` (scaffold authorship, see `gamma-scaffold.md`), `cdd/alpha` (this file), `writing` (docs-flip prose). No CLI/Go skills — no `src/go/` change in scope.

### AC evidence

- **AC1 — delegation (Mock D3).** `docs/guides/templates/cnos-install.yml`'s "cn repo install" step body is exactly two branches, each a single `cn repo install …` invocation; no `deps.json`/`deps.lock`/vendor-restore logic is written by this file. Verified: `grep -n "deps.json\|deps.lock\|restore" docs/guides/templates/cnos-install.yml` returns only the two doc-comment lines at the top of the file (prose referencing the filenames, not shell that writes them) — zero lines of reimplemented install logic.
- **AC2 — base PR (Mock D1).** The "Resolve installing token" step's `else` branch (taken when `install_dispatch` is `false`, the default) sets `token=${{ secrets.GITHUB_TOKEN }}`; that token is threaded into both `actions/checkout@v4` and `peter-evans/create-pull-request@v6`. Verified: `grep -n "secrets.GITHUB_TOKEN" docs/guides/templates/cnos-install.yml` → line 67, the only token source in the base-mode branch.
- **AC3 — dispatch gating (Mock D2).** The `if` branch (taken when `install_dispatch` is `true`) reads `secrets[inputs.workflow_pat_secret]` (dynamic secret-name lookup — GitHub Actions supports indexing the `secrets` context by an expression, the same idiom used for per-environment/per-matrix secret selection) and, if empty, emits `::error::` naming the unset secret and exits 1 **before** the `actions/checkout@v4` step runs — no partial write is possible because nothing has been checked out yet. Verified: `grep -n "::error::\|secrets\[inputs" docs/guides/templates/cnos-install.yml` → lines 60, 62; the `exit 1` is inside the same `if` block, before the job's only checkout step.
- **AC4 — docs flip.** `docs/guides/INSTALL-CDS.md` retains its existing "canonical local path = `cn repo install`" framing unchanged (the Layer 1/Layer 2 sections are untouched); the new `## GitHub UI (no-terminal) install` section names the template a "thin wrapper" explicitly and states the PAT precondition; the new `## Tenant secrets, by tier` table names Tier 1/2/3 secrets per the operator's scope-extension comment, naming cnos#613 as the (not-yet-landed) tracker rather than linking a nonexistent runbook path. `docs/guides/README.md` now links `INSTALL-CDS.md` under `## Operator` (previously absent — confirmed by reading the file before editing).

### Friction notes (carried from `gamma-scaffold.md`, confirmed still accurate at implementation time)

1. **Issue-body "Exists" framing was stale.** The issue's own text describes `docs/guides/templates/cnos-install.yml` as an existing file that "reimplements the install steps in shell." It does not exist on `main` (confirmed again at implementation time: `git show main:docs/guides/templates/cnos-install.yml` → `fatal: path ... does not exist in 'main'`). This cycle **creates** the file; it does not rewrite an existing one. The corrected framing was pinned in the γ scaffold before any file was written, so α did not improvise against a false premise.
2. **cnos#613 dependency gap.** The operator's scope-extension comment asks this cycle to link a tenant-secrets runbook that Sub 6 (#613) is supposed to produce; #613 is still `status:open` (confirmed via `gh issue view 613`), so no such file exists to link. Resolution (pinned in the γ scaffold, executed here): write the per-tier secrets content **inline** in `INSTALL-CDS.md` now, naming #613 as the tracker for the fuller runbook, rather than link a path that does not exist. This is disclosed as debt below, not silently patched over.

### Self-check

- Internal consistency: the new `## GitHub UI (no-terminal) install` section's claims (thin wrapper, no reimplementation, PAT precondition, fail-clear) all trace to lines in `docs/guides/templates/cnos-install.yml` cited above — no claim in the docs outruns what the YAML actually does.
- External alignment: `docs/guides/README.md`'s new link target (`INSTALL-CDS.md`) exists at the expected relative path; `docs/guides/INSTALL-CDS.md`'s new link target (`templates/cnos-install.yml`) exists at the expected relative path (both manually verified against the filesystem, not assumed).
- No existing content in `INSTALL-CDS.md`'s Layer 1 / Layer 2 / Verifying / Uninstalling / Troubleshooting / Related sections was removed or contradicted; only additive sections + two targeted edits (Uninstalling's Layer-2-cleanup line, Troubleshooting's new row, Related's new link) were made.
- `yamllint` (relaxed profile, line-length disabled) run against `docs/guides/templates/cnos-install.yml`: zero findings. `python3 -c "import yaml; yaml.safe_load(open(...))"` parses the file without error (the `on:` key's YAML-1.1 boolean-coercion quirk under PyYAML is a linter artifact, not a real GitHub Actions parsing issue — GitHub's own workflow parser treats `on:` as the literal trigger key).

### Debt

Disclosed explicitly, not silently patched over:

1. **The dynamic `secrets[inputs.workflow_pat_secret]` lookup is lint-checked (YAML syntax, grep-verified structure) but not exercised against a live `workflow_dispatch` run in this cycle.** No harness in this environment can fire a real GitHub Actions `workflow_dispatch` event and observe the runner's actual secret-masking/lookup behavior end-to-end. This is the same idiom documented and used elsewhere for per-environment/per-matrix dynamic secret selection, but this specific file's live behavior is unverified beyond static analysis. Flagging rather than silently claiming full behavioral verification.
2. **The Tenant-secrets-by-tier table duplicates content that cnos#613's own runbook will eventually own.** When #613 lands, this table should be reconciled — either replaced with a link to the canonical runbook, or kept as a short summary that links out. Not fixed here because #613 has not landed; fixing it now would mean either blocking this cycle on a sibling sub or fabricating a link to a file that doesn't exist. Named as forward-looking debt per the γ scaffold's friction note 2.
3. **No live PR was opened by this cycle's own template to prove Mock D1–D4 end-to-end** (that would require a second repo or a sandboxed Actions run this environment cannot provide). Verification is static: YAML lint, the delegation/token/fail-clear greps above, and doc cross-reference checks — matching the issue body's own "Proof plan" oracle ("workflow lint + a doc check ... the delegation grep"), not a live-fire requirement.
4. **`docs/guides/README.md`'s new entry is the only nav change made.** No other README/index in the repo was audited for a missing `INSTALL-CDS.md` link; if one exists elsewhere it is out of this cycle's named scope (`docs/guides/README.md` per the issue's own debt item 4 citation from #608).

### Review-readiness signal

R0 implementation complete; all four ACs have grep/read-verifiable evidence above; two friction notes resolved and disclosed; four debt items disclosed. Ready for β review.
