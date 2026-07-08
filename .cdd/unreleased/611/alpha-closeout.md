# α closeout — cnos#611

## Summary

All four ACs (AC1–AC4) closed with grep/read-verifiable evidence, converged at R0 (no iteration — β's independent re-derivation raised zero blocking findings). Scope: three files (`docs/guides/templates/cnos-install.yml` new, `docs/guides/INSTALL-CDS.md` edit, `docs/guides/README.md` edit). No `src/go/` changes, no `.cn/` schema changes.

```yaml
mock_parity:
  source: docs/development/design/cn-repo-install-MOCKS.md
  source_commit: ca52fc5081027df185e3a8a5191e621d115cfb66
  rows:
    - id: D1
      expectation: "Base run (dispatch off) opens a PR with the Mock-B diff, using the default GITHUB_TOKEN."
      observed: "docs/guides/templates/cnos-install.yml's else branch (install_dispatch != 'true') sets token=${{ secrets.GITHUB_TOKEN }}, threaded into checkout + create-pull-request."
      evidence: "docs/guides/templates/cnos-install.yml lines 55-67, 72, 101; beta-review.md §AC2"
      verdict: match
      how: "Static structural match; no live workflow_dispatch run available in this environment to fire the actual PR (disclosed as debt item 3)."
    - id: D2
      expectation: "Dispatch run either opens a PR containing the Mock-C workflow or fails with a clear message that a workflow-scoped PAT is required -- never a half-applied state."
      observed: "The token-resolution step (first in the job) reads secrets[inputs.workflow_pat_secret]; empty -> ::error:: + exit 1 before actions/checkout@v4 runs, so no partial write is reachable. Non-empty -> proceeds through cn repo install --dispatch cds using that token."
      evidence: "docs/guides/templates/cnos-install.yml lines 55-68; beta-review.md §AC3"
      verdict: match
      how: "Structural proof (fail-clear step precedes the only checkout/write steps in the job); not exercised against a live run (debt item 1)."
    - id: D3
      expectation: "The job body invokes cn repo install ...; install logic is not duplicated in shell."
      observed: "Single 'cn repo install' step, two branches (base / --dispatch cds), no deps.json/deps.lock/vendor-restore logic anywhere in the file."
      evidence: "grep -n 'deps.json\\|deps.lock\\|restore' docs/guides/templates/cnos-install.yml -> only 2 comment lines; beta-review.md §AC1"
      verdict: match
      how: "Direct grep + full-file read, independently re-run by both α and β."
    - id: D4
      expectation: "Never pushes to main."
      observed: "The template has no push-to-main step anywhere; the only write path is peter-evans/create-pull-request@v6 (PR-only by construction), and checkout does not check out main for direct commit."
      evidence: "docs/guides/templates/cnos-install.yml full read (114 lines); beta-review.md scope check"
      verdict: match
      how: "No git push / no direct-commit action referenced in the file; PR-opening action is the only mutation primitive used."
  summary:
    matched: 4
    exceeded: 0
    missed: 0
    exceed_justified: false
```

## Debt disclosed (carried from self-coherence.md, restated for the closeout)

1. Dynamic `secrets[inputs.workflow_pat_secret]` lookup is lint-checked and structurally reasoned about (by both α and β, independently) but not exercised against a live `workflow_dispatch` run in this cycle. No live-fire harness was available in this environment.
2. The Tenant-secrets-by-tier table in `INSTALL-CDS.md` duplicates content cnos#613's own runbook will eventually own; reconcile-or-link once #613 lands, don't duplicate indefinitely.
3. No live PR was opened by the template itself to prove Mock D1–D4 end-to-end; verification is static (lint, grep, doc cross-reference), matching the issue's own stated Proof plan oracle.
4. `docs/guides/README.md`'s new entry is the only nav change made; no repo-wide nav audit was performed beyond the file the issue names.

## What I'd tell a future cycle

The issue body's own "Exists" framing was stale (the template file did not exist on `main`) and the scope-extension comment asked for a link to a sibling cell's (#613) not-yet-existing artifact. Both were caught and resolved *before* writing code — by re-reading `main` directly (`find . -iname cnos-install.yml`, `gh issue view 613`) rather than trusting the issue text at face value. Corroborates the same discipline #608's own self-coherence applied to a similar stale-reference problem: verify the premise against the actual repo state before implementing against it.

## Artifact list

`docs/guides/templates/cnos-install.yml` (new), `docs/guides/INSTALL-CDS.md` (edit, additive), `docs/guides/README.md` (edit, one line), `.cdd/unreleased/611/{gamma-scaffold.md,self-coherence.md,beta-review.md,alpha-closeout.md}` (this cycle's process artifacts).

## Final AC status

AC1–AC4: all PASS, all independently re-derived by β (fresh greps/reads, not read off this artifact or self-coherence.md). Zero iteration rounds — R0 converged directly.
