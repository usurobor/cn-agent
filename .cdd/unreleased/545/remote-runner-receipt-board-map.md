# Remote-runner effect-surface receipt — `.github/workflows/board-map.yml`

Per `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` §8: any artifact that causes
another body to execute records this 6-field receipt, authored in the same
commit-set as the workflow file.

Honesty note: this workflow was authored as **operator-directed work** (cnos#545,
A2 wave autonomy — the operator explicitly authorized execution and merge), not
as an autonomous δ dispatch cycle. The receipt describes the workflow's actual
effect surface; it does not claim a γ/α/β triad that was not run.

1. **Who asked.** The operator (usurobor) via cnos#545 —
   "board: `cn issues map` — Go generator + refresh Action for the interactive
   issue-board", A2 execution directive.

2. **What runs.** On each trigger, the `regenerate` job in
   `.github/workflows/board-map.yml` executes, in order:
   - `actions/checkout@v4` (checks out the repo at the default branch)
   - `actions/setup-go@v5` (`go-version-file: src/go/go.mod`)
   - `go build -o cn ./cmd/cn` (builds the kernel binary from `src/go`)
   - `./cn issues map --repo "$GITHUB_REPOSITORY" --out docs/development/board`
     — reads the repo's open issues via the GitHub REST API (read-only) using
     the workflow `GITHUB_TOKEN`, and rewrites `docs/development/board/`
     (`index.html`, `pivot.html`, `board-data.json`, `README.md`) in place.
   - `git status --porcelain docs/development/board` (diff check only, no mutation)
   - if (and only if) that diff is non-empty:
     - `git config` a `board-map-bot` identity
     - `git add docs/development/board`
     - `git commit -m "board-map: regenerate docs/development/board from live open issues"`
     - `git pull --rebase` then `git push` to the checked-out ref

3. **Where it runs.** A GitHub-hosted `ubuntu-latest` runner, in the
   `usurobor/cnos` repository's own Actions environment. No self-hosted
   runners, no external compute.

4. **What authority.** The workflow-scoped default `GITHUB_TOKEN`, restricted by
   `permissions: {contents: write}` (the only permission block in the file — no
   `issues`, no `pull-requests`, no `id-token`). No other secrets are read; the
   issue read is a read-only GitHub API call scoped to the same token. The
   `git push` pushes only to this repository, only to the ref the workflow
   already checked out, and only touches paths under `docs/development/board/`.

5. **Evidence.** Post-run-fillable: the first successful `board-map` workflow run
   URL, observed on the repository's Actions tab after the first live trigger
   post-merge (`https://github.com/usurobor/cnos/actions/workflows/board-map.yml`).

6. **Who may accept.** The operator (usurobor). Acceptance criterion: the first
   automated commit to `docs/development/board/` authored by `board-map-bot`,
   following an issue change, is observed on the Actions tab (a `board-map` run
   with a non-skipped "Commit and push refreshed board" step).
