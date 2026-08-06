# Install CDS into your repo

**Audience:** a human who wants to bring the cnos **CDS** (Coherence-Driven
Software) process into an existing software repository.

## Terms

A few words below are easy to gloss over if you haven't set up a GitHub
Actions workflow before. Defined here, before anything below uses them:

- **PAT (Personal Access Token)** — a credential tied to *your own* GitHub
  account that authorizes API and git actions as you. Distinct from the
  built-in `GITHUB_TOKEN` GitHub Actions provisions automatically for every
  workflow run (which is scoped to that one run and cannot push new
  `.github/workflows/` files or trigger other workflows). This guide uses a
  **fine-grained PAT** — one scoped to a single repository and an explicit,
  minimal set of permissions, rather than a "classic" token with blanket
  account-wide access.
- **Repo secret** — a named value stored at a repository's
  **Settings → Secrets and variables → Actions**. Only that repo's own
  workflow runs can read it at runtime; the value is never shown again
  after creation, and `cn repo install` never reads, logs, or otherwise
  handles a secret's *value* — it only checks that a secret with a given
  *name* exists (see [§ Preflight](#preflight-what-the-operator-provides-before-anything-else-cnos706)).
- **Default branch** — the branch (usually `main`) GitHub treats as a
  repo's canonical line of history. A workflow file under
  `.github/workflows/` only takes effect once it is merged to the default
  branch — sitting on any other branch, it is inert.
- **"Bot"** — **there is no bot account to create for this install.**
  Earlier revisions of the rendered dispatch workflow carried a cosmetic
  `bot_name`/`bot_id` commit-author label naming a fixed, historical
  identity string — this was never a separate GitHub account, only a label
  layered over the operator's own PAT (verified against the live workflow
  and the GitHub API; see cnos#706). That cosmetic default was deleted: a
  fresh install authors commits as whichever account your `CN_DISPATCH_PAT`
  belongs to (your own), with no bot-flavored anything. A **real** dedicated
  bot account — a second GitHub login, added to the repo as a collaborator,
  so dispatch commits/PRs are attributable to a distinct identity — is
  optional future work (cnos#449 / cnos#702), not required and not built
  today.

---

Installing CDS has two layers, and they are separate trust decisions:

- **Layer 1 — Base package install.** Pins the `cn` toolchain reference and
  the cnos packages (`cnos.core`, `cnos.cdd`, `cnos.cds`) in your repo. This is
  the safe default — enough for you, or a Claude attached to the repo, to run
  the CDS method by hand. **This is what this guide covers.**
- **Layer 2 — Autonomous dispatch (opt-in).** A scheduled workflow that wakes
  an agent, claims issues, and opens PRs on a cron. This needs two secrets on
  **your own GitHub account** (a Claude OAuth token, a fine-grained PAT) and a
  manual merge to your default branch — no second account to create. See
  [§ Terms](#terms) for what these words mean and
  [§ Autonomous dispatch](#autonomous-dispatch-opt-in) before enabling it.

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
cn repo install --dispatch cds
```

That's the whole command for the common case — no `--agent`, no
`--workflow-pat-secret`, no bot flags. It runs the base install, then —
once [§ Preflight](#preflight-what-the-operator-provides-before-anything-else-cnos706)
passes — renders `.github/workflows/cnos-cds-dispatch.yml` bound to the
`sigma` identity (just a concurrency-group label, not an account — see
[§ Terms](#terms)) and the `CN_DISPATCH_PAT` / `CLAUDE_CODE_OAUTH_TOKEN`
secrets. Commits it makes are authored by whichever account
`CN_DISPATCH_PAT` belongs to — **your own** — with no bot-flavored label of
any kind.

Installing under a different named identity only changes the
concurrency-group label and the PAT secret's name — still bot-less by
default:

```sh
cn repo install --dispatch cds --agent acme --workflow-pat-secret ACME_DISPATCH_PAT
```

A non-sigma `--agent` requires `--workflow-pat-secret` naming your own
secret (fails early, before any render, if it is missing). `--bot-name` /
`--bot-id` still exist as a **strictly opt-in** cosmetic commit-author
override — you do not need them, and a fresh install never sets them by
default (cnos#706 AC9). The installing token needs `workflow` scope to
write `.github/workflows/`; the command itself never pushes to `main`
(PR-only).

### Preflight: what the operator provides, before anything else (cnos#706)

Before `--dispatch cds` renders a single file, it checks three things only
*you* (the operator) can provide, and refuses to proceed until they are all
in place — so you are never left with a half-deployed, inert workflow:

1. **`CLAUDE_CODE_OAUTH_TOKEN`** repo secret exists (checked by name only —
   its value is never read).
2. **`CN_DISPATCH_PAT`** repo secret exists (checked by name only — see
   [§ Tier 3 runbook](#tier-3-runbook-autonomous-dispatch--two-own-account-secrets)
   for how to create one).
3. **Push access** — the token you're running `cn repo install` with has
   push access to this repo (checked via a single read-only GitHub API
   call; no secret value is ever sent to or read by the CLI).

If any of these is missing, the command exits non-zero and prints, for each
missing item, what it is, why it's needed, and the exact steps to get it —
the same text as the runbook linked above. **Nothing is written** — no
`.cn/`, no `.github/workflows/` — until every prerequisite is satisfied.
Once they are, re-run the exact same command: `cn repo install` is
idempotent, so it resumes and completes cleanly (no special "resume" flag
needed).

A fourth gate is deliberately **not** automated: merging the install PR to
your default branch. An autonomous, PR-opening agent must not be able to
self-activate its own scheduled automation — that step stays a manual
merge, by design, not an oversight.

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
| **Tier 3 — autonomous dispatch** | This guide's Layer 2 (`--dispatch cds`): a scheduled agent that claims cells and opens PRs. | `cn repo install --dispatch cds` | Two secrets on **your own GitHub account** — `CLAUDE_CODE_OAUTH_TOKEN` and `CN_DISPATCH_PAT` (a fine-grained PAT; see the runbook below). No separate bot account. |

Both Tier 2 and Tier 3 write `.github/workflows/cnos-cds-dispatch.yml`, so
the **installing** token (or the GitHub UI path's `workflow_pat_secret`)
needs `workflow` scope *at install time* regardless of tier — that is a
one-time cost of committing the workflow file, distinct from the *runtime*
secrets above.

### Tier 2 runbook (mechanical FSM engine — `GITHUB_TOKEN` only)

The engine tier provisions **no runtime secret of its own** and needs **no
repository setting changes**: it advances issue-label state only, opens no
pull request, and runs entirely on the default `GITHUB_TOKEN`.

1. Ensure the install target has a resolvable `origin` git remote and that
   your local `$GITHUB_TOKEN`/`$GH_TOKEN` can manage labels (the
   `label-doctor` precondition, cnos#493).
2. Run `cn repo install --dispatch cds --engine`. This renders
   `.github/workflows/cnos-cds-dispatch.yml` with every token binding set to
   `${{ secrets.GITHUB_TOKEN }}`, a least-privilege permissions block
   (`contents: read`, `issues: write`, `pull-requests: read`), and ensures
   the canonical dispatch labels.
3. Review and merge the install PR. Nothing else to provision — the
   scheduled workflow runs on `GITHUB_TOKEN` alone and only advances
   issue-label FSM state through the GitHub API. You do **not** enable
   *Allow GitHub Actions to create and approve pull requests* for this tier:
   the engine wake opens no PR (`contents` and `pull-requests` are both
   read-only), so that setting would grant permissions it never uses. (It is
   the Tier 3 agent tier — which writes code and opens PRs — that needs it.)

### Tier 3 runbook (autonomous dispatch — two own-account secrets)

There is no bot account to create (see [§ Terms](#terms)) — just two
secrets on your own GitHub account. This is the exact text the
[§ Preflight](#preflight-what-the-operator-provides-before-anything-else-cnos706)
check prints when either is missing:

1. **`CLAUDE_CODE_OAUTH_TOKEN`** — authorizes the dispatch agent to call
   Claude. Get it by running **`claude setup-token`** locally (requires a
   Claude Pro/Max subscription); paste the printed token as this repo
   secret. **Settings → Secrets and variables → Actions → New repository
   secret.**
2. **`CN_DISPATCH_PAT`** — a **fine-grained Personal Access Token on your
   own GitHub account**, scoped to this repo with **Contents + Issues +
   Pull requests + Workflows = write**. The dispatch workflow uses it to
   check out, move FSM labels via the API, and push the cell branch + open
   the PR. Create it at **Settings → Developer settings → Personal access tokens → Fine-grained tokens.**

   *Why a PAT and not the built-in `GITHUB_TOKEN`?* Installing needs
   workflow-write to commit `.github/workflows/cnos-cds-dispatch.yml`,
   which `GITHUB_TOKEN` cannot do; and at runtime, GitHub blocks
   `GITHUB_TOKEN`-authored pushes from triggering *other* workflows — so a
   `GITHUB_TOKEN`-only setup would silently lose PR CI and the
   `issues: labeled` fast-path (it would still limp along on the cron
   backstop, but that's a caveat, not the default). This has been verified
   against the dispatch workflow's actual triggers, not assumed.
3. Run `cn repo install --dispatch cds` (bare — the sigma default needs no
   `--agent` / `--workflow-pat-secret` / `--bot-name` / `--bot-id`). If
   either secret above is missing, or the installing token lacks push
   access, the command exits non-zero explaining exactly which and how to
   fix it, *before* touching your repo.
4. Review and merge the install PR to your **default branch** — deliberately
   manual; an autonomous PR-opening agent must not be able to self-activate
   its own scheduled automation.
5. Rotate both secrets on your normal credential-rotation cadence; the PAT
   must retain all four scopes or the scheduled wake cannot check out, move
   labels, or push.

A **dedicated bot account** — a second GitHub login, added as a repo
collaborator, so dispatch commits/PRs are attributable to a distinct
identity — is optional future work (cnos#449 / cnos#702), not required
here.

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
| `--dispatch cds` fails with "preflight: operator prerequisites are not yet satisfied ..." | Nothing was written yet (cnos#706 — this is the first thing the command checks). The message names exactly which of `CLAUDE_CODE_OAUTH_TOKEN` / `CN_DISPATCH_PAT` / push access is missing and how to get it — see [§ Tier 3 runbook](#tier-3-runbook-autonomous-dispatch--two-own-account-secrets). Re-run the exact same command once satisfied; it resumes cleanly. |
| `--dispatch cds` fails with "canonical dispatch labels not ensured: ..." | The workflow still rendered (preflight already passed); label-doctor could not resolve the repo's `origin` git remote, a GitHub token, or reach the GitHub API. Run `cn label doctor` yourself (or apply the labels manually) once the underlying issue (missing remote/token/scope) is fixed. See [§ Autonomous dispatch](#autonomous-dispatch-opt-in). |
| `--dispatch cds` fails with "--workflow-pat-secret is required for --agent ..." | Pass `--workflow-pat-secret <NAME>` naming the GitHub Actions secret holding that agent's workflow-scoped PAT. `--bot-name`/`--bot-id` are optional cosmetic overrides, not required. Or, if you want the mechanical (no-agent) tier, add `--engine` — it needs no PAT. |
| `--engine is only valid with --dispatch cds` | `--engine` selects the mechanical tier of the *cds dispatch* render; it is meaningless on a base or `--dispatch none` install. Add `--dispatch cds`, or drop `--engine`. |
| Engine-tier workflow advances labels but opens no PR | Expected — by design. The `--engine` tier is label-FSM only: it renders `contents: read` + `pull-requests: read` and never opens a pull request. If you want a tier that opens PRs, that is the Tier 3 agent tier (`--dispatch cds` without `--engine`), which needs a `workflow`-scope PAT + `CLAUDE_CODE_OAUTH_TOKEN`. |
| GitHub UI install run fails with "`install_dispatch=true` requires a workflow-scoped PAT ..." | Set the repo secret named by the `workflow_pat_secret` input (default `CNOS_WORKFLOW_PAT`) to a `workflow`-scoped PAT before re-running "Run workflow" — see [§ GitHub UI (no-terminal) install](#github-ui-no-terminal-install). |

## Related

- [`docs/reference/cli/CLI.md`](../reference/cli/CLI.md) — CLI reference.
- [dispatch orchestrator skill](../../src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md) — the autonomous dispatch loop (Layer 2).
- [`docs/guides/templates/cnos-install.yml`](templates/cnos-install.yml) — the GitHub UI install workflow template.
