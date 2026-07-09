# Install CDS into your repo

**Audience:** a human who wants to bring the cnos **CDS** (Coherence-Driven
Software) process into an existing software repository.

Installing CDS has two layers, and they are separate trust decisions:

- **Layer 1 — Base package install.** Pins the `cn` toolchain reference and
  the cnos packages (`cnos.core`, `cnos.cdd`, `cnos.cds`) in your repo. This is
  the safe default — enough for you, or a Claude attached to the repo, to run
  the CDS method by hand. **This is what this guide covers.**
- **Layer 2 — Autonomous dispatch (opt-in).** A scheduled workflow that wakes
  an agent, claims issues, and opens PRs on a cron. This needs extra secrets
  (a `workflow`-scoped PAT, a model OAuth token) and standing write access —
  see [§ Autonomous dispatch](#autonomous-dispatch-opt-in) before enabling it.

The canonical way to install Layer 1 is one command:

```sh
curl -fsSL https://raw.githubusercontent.com/usurobor/cnos/main/install.sh | sh
cn repo install
```

`cn repo install` is a kernel command — it runs before any hub or package
state exists, so it works in a plain `git init`-only checkout with nothing
more than the `cn` binary on `PATH`.

---

## What "installing CDS" actually means

CDS is a **method**, not a service. Installed into your repo (Layer 1), the
base layer is just three tracked artifacts:

| File | Purpose |
|---|---|
| `.cn/deps.json` | Declares which cnos packages your repo depends on (exact versions). |
| `.cn/deps.lock.json` | The resolved, SHA-256-pinned lockfile (schema `cn.lock.v2`). |
| `.gitignore` (`+ .cn/vendor/`) | Vendored packages rehydrate from the lock; they are not committed. |

`cn repo install` restores the packages into `.cn/vendor/packages/<name>/`
(name-based, not version-suffixed), which includes the CDS skills
(`skills/cds/CDS.md`, the lifecycle and selection overlays) and the
`cn install-wake` renderer.

The `cn` binary is **the body**; a model (Claude, etc.) is **the brain**. `cn`
senses and executes; the model reads the skills and drives the CDS loop. See
[`docs/reference/cli/CLI.md`](../reference/cli/CLI.md).

`cn repo install` never writes anything under `.github/workflows/` and never
requires a workflow or agent secret — that is Layer 2, a separate opt-in (see
below).

---

## Prerequisites

- A GitHub repository, checked out locally, that you can push to (or open PRs
  against).
- A Unix-like shell (Linux, macOS, WSL) with `curl`, to install the `cn`
  binary.

`cn repo install` itself needs nothing beyond `cn` and network access to
GitHub Releases (or a local/offline package index — see `--index` below).

---

## Install (Layer 1 — base)

### 1. Install the `cn` binary

```sh
curl -fsSL https://raw.githubusercontent.com/usurobor/cnos/main/install.sh | sh
```

This downloads the pre-built binary for your platform from the latest cnos
release, verifies its SHA-256, and installs it to `/usr/local/bin/cn`. To
install without root:

```sh
curl -fsSL https://raw.githubusercontent.com/usurobor/cnos/main/install.sh \
  | BIN_DIR="$HOME/.local/bin" sh
# ensure $HOME/.local/bin is on your PATH
```

Verify:

```sh
cn --version
```

### 2. Install the base package set

From your repo root:

```sh
cn repo install
```

This resolves the latest cnos release, writes `.cn/deps.json` +
`.cn/deps.lock.json`, restores `cnos.core` / `cnos.cdd` / `cnos.cds` under
`.cn/vendor/packages/`, and adds `.cn/vendor/` to `.gitignore`. It prints the
resolved release tag and a summary of what it wrote/restored.

Preview what would happen without writing anything:

```sh
cn repo install --dry-run
```

Useful flags (all optional — `cn repo install` alone covers the default
case):

```sh
cn repo install --release <tag>              # pin a specific cnos release instead of latest
cn repo install --packages cnos.core,cnos.cdd,cnos.cds
cn repo install --index ./dist/packages/index.json   # local/offline package index
```

### 3. Commit

```sh
git add .cn/deps.json .cn/deps.lock.json .gitignore
git commit -m "chore(cnos): install CDS packages"
```

That's the base install. You (or an attached Claude) can now open
`.cn/vendor/packages/cnos.cds/skills/cds/CDS.md` and run the CDS lifecycle.

`cn repo install` is idempotent: running it again with the same inputs
produces no further diff.

---

## Autonomous dispatch (opt-in, Layer 2)

The base install gives you the CDS *method*. The **dispatch loop** is the
automation from the cnos repo itself: a scheduled workflow
(`cnos-cds-dispatch.yml`) that wakes an agent, claims `dispatch:cell +
protocol:cds + status:todo` issues, runs each cell through the δ role
contract, and opens PRs with receipts. See the
[dispatch orchestrator skill](../../src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md)
for the full protocol.

`cn repo install --dispatch cds` (cnos#610) is the entry point for installing
this layer:

```sh
cn repo install --dispatch cds \
  --agent acme --workflow-pat-secret ACME_WORKFLOW_PAT \
  --bot-name acme-bot --bot-id 12345678
```

This runs the base install, then renders `.github/workflows/cnos-cds-dispatch.yml`
for the given agent identity — no sigma binding required; any non-sigma
`--agent` requires `--workflow-pat-secret` (fails early, before any render,
if it is missing). The installing token needs `workflow` scope to write
`.github/workflows/`; the command itself never pushes to `main` (PR-only).

**The canonical dispatch labels are ensured automatically:** after
rendering the workflow, `--dispatch cds` audits the installing repo's
GitHub labels against `cnos.core`'s `labels.json` and creates/repairs any
missing or drifted one — the `label-doctor` mechanism (cnos#493), invoked
in-process (no separate command required). This needs the installing
repo to have a resolvable `origin` git remote and a GitHub token
(`$GITHUB_TOKEN`/`$GH_TOKEN`) with permission to manage labels. If
either is unavailable, `--dispatch cds` still renders the workflow but
exits non-zero naming the label-doctor failure, and you can apply the
labels yourself with `cn label doctor` (or manually).

### `--engine`: the PAT-free mechanical tier

`cn repo install --dispatch cds` renders **two separable wake tiers** from
the same command; you pick one with the `--engine` flag:

| Tier | Flag | Work phase | Runtime credential |
|---|---|---|---|
| **Agent** (default) | *(no flag)* | `anthropics/claude-code-action@v1` — an agent claims cells and writes code. | `workflow`-scope PAT **and** `CLAUDE_CODE_OAUTH_TOKEN`. |
| **Engine** | `--engine` | Mechanical CDS issue-state FSM (`cn issues fsm scan/evaluate --apply`) — advances label state deterministically, no agent in the loop. | The default `GITHUB_TOKEN` only. |

```sh
cn repo install --dispatch cds --engine
```

The engine tier (cnos#613) needs **no** `--workflow-pat-secret`,
`--bot-name`, `--bot-id`, or `CLAUDE_CODE_OAUTH_TOKEN`: every token binding
in the rendered workflow resolves to `${{ secrets.GITHUB_TOKEN }}`, and the
workflow declares an explicit **least-privilege** permissions block — not
the agent tier's broader one:

```yaml
permissions:
  contents: read        # checkout only
  issues: write         # the only writes — issue-label FSM transitions
  pull-requests: read
```

That is all `cn issues fsm ... --apply` requires. The engine tier grants no
`contents: write`, no `pull-requests: write`, and no `id-token`: it advances
label state only (through the guarded FSM, via the GitHub API), produces no
local commits, and opens no pull request. `--agent` is optional and only
names the workflow's concurrency group. `--engine` is only valid with
`--dispatch cds` (a rendering variant of the cds workflow); on a base or
`--dispatch none` install it fails early with a clear error.

Note the split between *install* time and *run* time: writing the workflow
file itself still touches `.github/workflows/`, so the **installing** token
(the human/CI running `cn repo install`, or the GitHub UI path's
`workflow_pat_secret`) still needs `workflow` scope — but the rendered
workflow's own **runtime** is PAT-free and least-privilege. (The agent tier,
by contrast, both writes code and opens PRs from its work phase — enable
*Settings → Actions → General → Allow GitHub Actions to create and approve
pull requests* for that tier.)

---

## GitHub UI (no-terminal) install

Both layers above assume a terminal. If you'd rather install from the GitHub
web UI — no local shell, no `cn` binary on your machine — copy
[`docs/guides/templates/cnos-install.yml`](templates/cnos-install.yml) into
your repo at `.github/workflows/cnos-install.yml`, then trigger it from the
**Actions** tab ("Run workflow").

> **Heads up — a more native install surface is coming.** This Action path is
> the **zero-backend fallback**: it works with nothing but GitHub, but you
> must first commit a workflow file before GitHub shows its "Run workflow"
> button. The polished, App-first surface is a **CNOS Installer GitHub App**
> you install onto a repo, which opens the install PR for you — no workflow to
> bootstrap. Its architecture and per-tier permission model are specified in
> [`docs/development/design/cnos-installer-github-app.md`](../development/design/cnos-installer-github-app.md);
> this Action remains the supported fallback once the App ships.

This is a **thin wrapper**, not a second install path: the workflow's job
body installs the `cn` binary and calls the exact same `cn repo install`
command described above — install logic is not duplicated in YAML. It never
runs automatically (`workflow_dispatch` only) and never pushes to `main`;
it opens a pull request with the result, same as a human running the CLI by
hand and opening a PR themselves.

Inputs:

| Input | Default | Meaning |
|---|---|---|
| `release` | `latest` | Same as `cn repo install --release`. |
| `install_dispatch` | `false` | `true` installs Layer 2 (`--dispatch cds`) on top of the base install. |
| `workflow_pat_secret` | `CNOS_WORKFLOW_PAT` | Name of the repo secret holding a workflow-scoped PAT. Required (and must already be set) when `install_dispatch: true` — see below. |
| `agent` | `cnos` | Dispatch agent identity; only used when `install_dispatch: true`. |

**Base run** (`install_dispatch: false`) needs nothing beyond the default
`GITHUB_TOKEN` GitHub Actions already provides — the diff never touches
`.github/workflows/`.

**Dispatch run** (`install_dispatch: true`) writes a *new*
`.github/workflows/cnos-cds-dispatch.yml`, which the default `GITHUB_TOKEN`
cannot push (GitHub blocks `workflow`-scope writes from the default token).
You must set the secret named by `workflow_pat_secret` to a `workflow`-scoped
PAT **before** running the workflow. If it's unset or empty, the run fails
immediately with a clear error naming the missing secret — it never opens a
half-applied PR.

## Tenant secrets, by tier

What you provision depends on how far up the automation ladder you go:

| Tier | What it is | Install command | Runtime secrets |
|---|---|---|---|
| **Tier 1 — base install** | `cn repo install` (this guide's Layer 1). Just the CDS method, no automation. | `cn repo install` | None beyond what GitHub Actions already provides (`GITHUB_TOKEN`), and only if you use the GitHub UI path above — the plain CLI path needs no secrets at all. |
| **Tier 2 — mechanical FSM engine** | A PAT-free mechanical engine (the CDS issue-state FSM, `cn issues fsm scan/evaluate --apply`) that reconciles label state without an agent in the loop. | `cn repo install --dispatch cds --engine` | `GITHUB_TOKEN` only (the default Actions token — no secret to provision). |
| **Tier 3 — autonomous dispatch** | This guide's Layer 2 (`--dispatch cds`): a scheduled agent that claims cells and opens PRs. | `cn repo install --dispatch cds` | A `workflow`-scope PAT (named by `--workflow-pat-secret` / `workflow_pat_secret`) **and** `CLAUDE_CODE_OAUTH_TOKEN` for the agent runtime. |

Both Tier 2 and Tier 3 write `.github/workflows/cnos-cds-dispatch.yml`, so
the **installing** token (or the GitHub UI path's `workflow_pat_secret`)
needs `workflow` scope *at install time* regardless of tier — that is a
one-time cost of committing the workflow file, distinct from the *runtime*
secrets above.

### Tier 2 runbook (mechanical FSM engine — `GITHUB_TOKEN` only)

The engine tier provisions **no runtime secret of its own**. The one setup
choice is whether GitHub Actions may open the recovery/finalizer PRs.

1. Ensure the install target has a resolvable `origin` git remote and that
   your local `$GITHUB_TOKEN`/`$GH_TOKEN` can manage labels (the
   `label-doctor` precondition, cnos#493).
2. Run `cn repo install --dispatch cds --engine`. This renders
   `.github/workflows/cnos-cds-dispatch.yml` with every token binding set to
   `${{ secrets.GITHUB_TOKEN }}` and ensures the canonical dispatch labels.
3. In the repo, enable **Settings → Actions → General → Allow GitHub Actions
   to create and approve pull requests** so the workflow's built-in
   `GITHUB_TOKEN` can open the mechanical recovery/finalizer PRs.
4. Review and merge the install PR. No repo secret needs to be set — the
   scheduled workflow runs on `GITHUB_TOKEN` alone.

### Tier 3 runbook (autonomous dispatch — PAT + OAuth token)

The agent tier runs `claude-code-action`, which needs two credentials:

1. Create a **`workflow`-scope PAT** and store it as a repo Actions secret.
   Name it whatever you pass to `--workflow-pat-secret` (e.g.
   `ACME_WORKFLOW_PAT`); for the default `sigma` agent the renderer expects
   `SIGMA_WORKFLOW_PAT`.
2. Create a **`CLAUDE_CODE_OAUTH_TOKEN`** repo Actions secret holding the
   agent runtime's OAuth token.
3. Run `cn repo install --dispatch cds --agent <name> --workflow-pat-secret
   <SECRET_NAME> --bot-name <name> --bot-id <id>` (a non-sigma agent
   requires all four; a bare `--dispatch cds` uses the sigma defaults).
4. Review and merge the install PR. Rotate both secrets on your normal
   credential-rotation cadence; the PAT must retain `workflow` scope or the
   scheduled wake cannot check out and push.

---

## Verifying the install

```sh
cn repo install          # idempotent; re-running produces no further diff
ls .cn/vendor/packages    # cnos.core, cnos.cdd, cnos.cds
cat .cn/vendor/packages/cnos.cds/skills/cds/SKILL.md   # the CDS loader skill
```

## Uninstalling

Remove `.cn/deps.json`, `.cn/deps.lock.json`, the `.cn/vendor/` gitignore
line, and (if you enabled Layer 2) `.github/workflows/cnos-cds-dispatch.yml`.
If you adopted the GitHub UI path, also remove
`.github/workflows/cnos-install.yml`. The `cn` binary is just a file on your
`PATH`; delete it wherever `install.sh` placed it.

## Troubleshooting

| Symptom | Cause / fix |
|---|---|
| `✗ cn repo install must be run inside a Git repository.` | Run it from inside a checked-out Git repo (`git init` at minimum) — `cn repo install` never walks up looking for one or scaffolds a repo for you. |
| `package(s) not found in index` | The requested `--packages` entry isn't published in the resolved release/index. Check spelling, or pin `--release` to a tag that publishes it. |
| `package(s) have multiple versions in index; pass --release to pin one` | You passed `--index` pointing at a multi-version index with no `--release`; add `--release <tag>` to disambiguate. |
| `cn: command not found` after install | `install.sh` put `cn` outside your `PATH`. Re-run with `BIN_DIR="$HOME/.local/bin"` and add that dir to `PATH`. |
| `--dispatch cds` fails with "canonical dispatch labels not ensured: ..." | The workflow still rendered; label-doctor could not resolve the repo's `origin` git remote, a GitHub token, or reach the GitHub API. Run `cn label doctor` yourself (or apply the labels manually) once the underlying issue (missing remote/token/scope) is fixed. See [§ Autonomous dispatch](#autonomous-dispatch-opt-in). |
| `--dispatch cds` fails with "--workflow-pat-secret is required for --agent ..." | Pass `--workflow-pat-secret <NAME>` (and, for a non-sigma agent, `--bot-name`/`--bot-id`) naming the GitHub Actions secret holding that agent's workflow-scoped PAT. Or, if you want the mechanical (no-agent) tier, add `--engine` — it needs no PAT. |
| `--engine is only valid with --dispatch cds` | `--engine` selects the mechanical tier of the *cds dispatch* render; it is meaningless on a base or `--dispatch none` install. Add `--dispatch cds`, or drop `--engine`. |
| Engine-tier workflow runs but opens no recovery/finalizer PR | The built-in `GITHUB_TOKEN` cannot open PRs until you enable **Settings → Actions → General → Allow GitHub Actions to create and approve pull requests**. |
| GitHub UI install run fails with "`install_dispatch=true` requires a workflow-scoped PAT ..." | Set the repo secret named by the `workflow_pat_secret` input (default `CNOS_WORKFLOW_PAT`) to a `workflow`-scoped PAT before re-running "Run workflow" — see [§ GitHub UI (no-terminal) install](#github-ui-no-terminal-install). |

## Related

- [`docs/reference/cli/CLI.md`](../reference/cli/CLI.md) — CLI reference.
- [`docs/development/design/cn-repo-install-MOCKS.md`](../development/design/cn-repo-install-MOCKS.md) — the design surface `cn repo install` was built against.
- [dispatch orchestrator skill](../../src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md) — the autonomous dispatch loop (Layer 2).
- [`docs/guides/templates/cnos-install.yml`](templates/cnos-install.yml) — the GitHub UI install workflow template.
