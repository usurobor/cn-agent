# Install CDS into your repo

**Audience:** a human who wants to bring the cnos **CDS** (Coherence-Driven
Software) process into an existing software repository.

This guide has two tracks:

- **Track A — Terminal.** Run the `cn` command in the root of your repo. Fastest
  if you already work locally.
- **Track B — GitHub only, no terminal.** Commit one installer workflow, click
  **Run workflow**, review the pull request it opens. Nothing runs on your
  machine.

Both tracks install the same thing: the `cn` toolchain reference plus the cnos
packages (`cnos.core`, `cnos.cdd`, `cnos.cds`) pinned in your repo. That is the
**base layer** — enough for you, or a Claude attached to the repo, to run the
CDS method by hand.

Turning on the **autonomous dispatch loop** (an agent that wakes on a schedule,
claims issues, and runs cells) is a separate opt-in with extra prerequisites —
see [§ Autonomous dispatch](#autonomous-dispatch-opt-in) and read its caveats
before enabling it.

---

## What "installing CDS" actually means

CDS is a **method**, not a service. Installed into your repo, the base layer is
just three tracked artifacts:

| File | Purpose |
|---|---|
| `.cn/deps.json` | Declares which cnos packages your repo depends on. |
| `cn.lock` | The resolved, SHA-256-pinned lockfile. |
| `.gitignore` (`+ .cn/vendor/`) | Vendored packages rehydrate from the lock; they are not committed. |

Anyone who then runs `cn deps restore` gets the packages materialized under
`.cn/vendor/packages/`, which includes the CDS skills (`skills/cds/CDS.md`, the
lifecycle and selection overlays) and the `cn install-wake` renderer.

The `cn` binary is **the body**; a model (Claude, etc.) is **the brain**. `cn`
senses and executes; the model reads the skills and drives the CDS loop. See
[`docs/reference/cli/CLI.md`](../reference/cli/CLI.md).

---

## Prerequisites

- A GitHub repository you can push to (or open PRs against).
- **Track A** additionally needs: a Unix-like shell (Linux, macOS, WSL), `curl`,
  `jq`, and `python3` with `pyyaml` (the last two only if you later render a wake
  workflow with `cn install-wake`).
- **Track B** needs nothing local — the GitHub-hosted runner provides all of it.

---

## Track A — Terminal (`cn` in your repo root)

### 1. Install the `cn` binary

```sh
curl -fsSL https://raw.githubusercontent.com/usurobor/cnos/main/install.sh | sh
```

This downloads the pre-built binary for your platform from the latest cnos
release, verifies its SHA-256, and installs it to `/usr/local/bin/cn`. To install
without root:

```sh
curl -fsSL https://raw.githubusercontent.com/usurobor/cnos/main/install.sh \
  | BIN_DIR="$HOME/.local/bin" sh
# ensure $HOME/.local/bin is on your PATH
```

Verify:

```sh
cn --version
```

### 2. Get the package index

`cn` resolves packages from a package **index** (`dist/packages/index.json`).
Each cnos release publishes that index plus the package tarballs as release
assets. Pull them into your repo root:

```sh
# From your repo root. Replace TAG with the release you want (e.g. v3.82.0),
# or resolve the latest tag first.
TAG="$(gh release view -R usurobor/cnos --json tagName -q .tagName)"

mkdir -p dist/packages
gh release download -R usurobor/cnos "$TAG" \
  --dir dist/packages --pattern 'index.json' --pattern '*.tar.gz'
```

> No `gh` CLI? Download `index.json` and the `*.tar.gz` assets from
> `https://github.com/usurobor/cnos/releases` by hand into `dist/packages/`.

`dist/packages/` holds build inputs only — you can delete it after restore; it
does not belong in your repo.

### 3. Declare the CDS packages

Create `.cn/deps.json` in your repo root. Pin the versions to the release you
downloaded (they must match the index):

```json
{
  "schema": "cn.deps.v1",
  "profile": "cds",
  "packages": [
    { "name": "cnos.core", "version": "3.82.0" },
    { "name": "cnos.cdd",  "version": "3.82.0" },
    { "name": "cnos.cds",  "version": "3.82.0" }
  ]
}
```

### 4. Lock and restore

```sh
cn deps lock       # writes cn.lock from .cn/deps.json + dist/packages/index.json
cn deps restore    # installs packages into .cn/vendor/packages/
```

Keep vendored packages out of git:

```sh
grep -qxF '.cn/vendor/' .gitignore || echo '.cn/vendor/' >> .gitignore
rm -rf dist/packages          # build inputs, no longer needed
```

### 5. Commit

```sh
git add .cn/deps.json cn.lock .gitignore
git commit -m "chore(cnos): install CDS packages"
```

That's the base install. You (or an attached Claude) can now open
`.cn/vendor/packages/cnos.cds/skills/cds/CDS.md` and run the CDS lifecycle.

---

## Track B — GitHub only (no terminal)

The "fully mechanical" GitHub path uses **`workflow_dispatch`** — a manual **Run
workflow** button in the Actions tab. You commit one installer workflow once; it
runs everything from Track A on a GitHub runner and opens a pull request with the
result.

### 1. Add the installer workflow

Copy [`templates/cnos-install.yml`](templates/cnos-install.yml) into your repo at:

```
.github/workflows/cnos-install.yml
```

Commit it to your default branch (this file is an ordinary workflow, so the
normal push works).

### 2. Run it

1. Go to your repo's **Actions** tab.
2. Select **cnos install (CDS)** in the left sidebar.
3. Click **Run workflow**. Optionally set the release tag or the package list;
   the defaults install `cnos.core cnos.cdd cnos.cds` from the latest release.
4. Click the green **Run workflow** button.

### 3. Review the PR

The workflow installs `cn`, fetches the package index, writes `.cn/deps.json`,
locks + restores, and opens a pull request titled **"Install cnos CDS (…)"**
containing `.cn/deps.json` + `cn.lock` + the `.gitignore` entry. Review and
merge it — no terminal touched.

The default `GITHUB_TOKEN` is enough here because the base install never writes a
file under `.github/workflows/`.

---

## Autonomous dispatch (opt-in)

The base install gives you the CDS *method*. The **dispatch loop** is the
automation from the cnos repo itself: a scheduled workflow
(`cnos-cds-dispatch.yml`) that wakes an agent, claims `dispatch:cell +
protocol:cds + status:todo` issues, runs each cell through the δ role contract,
and opens PRs with receipts. See [`AUTOMATION.md`](AUTOMATION.md) and the
[dispatch orchestrator skill](../../src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md).

Once packages are restored, the renderer for that workflow is available as a
package command:

```sh
cn install-wake cds-dispatch --out .github/workflows/cnos-cds-dispatch.yml
```

**Before you enable this, know the current limits:**

1. **The renderer is bound to the `sigma` agent.** `cn install-wake` hardcodes
   the substrate secret name `SIGMA_WORKFLOW_PAT` and sigma's bot identity, and
   refuses any other agent. Rendered into your repo, the workflow references
   *sigma's* secret/identity, not yours. Generalizing this per-install binding
   (secret name, bot name/id, agent identity) is required before the dispatch
   layer is turnkey for a third-party repo. Track this as the productization gap
   it is — do not silently ship a sigma-wired workflow into an unrelated repo.
2. **It needs secrets.** The workflow expects an OAuth token for the model
   (`CLAUDE_CODE_OAUTH_TOKEN`) and a workflow PAT (`SIGMA_WORKFLOW_PAT` today).
   The PAT must carry the **`workflow` scope** — the default `GITHUB_TOKEN`
   cannot create or update files under `.github/workflows/`, so the Track B
   installer cannot write the dispatch workflow for you.
3. **It grants standing write access.** A scheduled agent with a write-scoped PAT
   will act on your issues and open PRs on a cron. Enable it deliberately, on a
   repo where that is intended.

Because of (1)–(3), this guide installs the base layer as the default and treats
autonomous dispatch as a manual, eyes-open step rather than part of the one-click
install. When the sigma binding is generalized, the Track B installer can grow a
`workflow_dispatch` input to render and commit the dispatch workflow (via a
`workflow`-scoped PAT) in the same PR.

---

## Verifying the install

```sh
cn deps restore        # idempotent; re-materializes .cn/vendor/packages/
ls .cn/vendor/packages # cnos.core, cnos.cdd, cnos.cds, ...
cat .cn/vendor/packages/cnos.cds/skills/cds/SKILL.md   # the CDS loader skill
```

## Uninstalling

Remove `.cn/deps.json`, `cn.lock`, the `.cn/vendor/` gitignore line, and (if you
added it) `.github/workflows/cnos-install.yml`. The `cn` binary is just a file on
your PATH; delete it wherever `install.sh` placed it.

## Troubleshooting

| Symptom | Cause / fix |
|---|---|
| `package '…' not present in … index.json` | The version in `.cn/deps.json` doesn't match the fetched release. Re-pin to the tag you downloaded, or download the tag matching your pins. |
| `cn deps restore` → "No packages to restore" | No `cn.lock`. Run `cn deps lock` first (needs `dist/packages/index.json` present). |
| `cn: command not found` after install | `install.sh` put `cn` outside your PATH. Re-run with `BIN_DIR="$HOME/.local/bin"` and add that dir to PATH. |
| Track B PR is empty | CDS is already installed (deps.json + cn.lock unchanged). Nothing to do. |
| `cn install-wake` → "unknown agent" | Expected today — the renderer only knows `sigma`. See [§ Autonomous dispatch](#autonomous-dispatch-opt-in). |

## Related

- [`README.md` → Try it](../../README.md#try-it) — the upstream install one-liner.
- [`docs/reference/cli/CLI.md`](../reference/cli/CLI.md) — full `cn` command reference.
- [`docs/reference/cli/SETUP-INSTALLER.md`](../reference/cli/SETUP-INSTALLER.md) — installer design doc.
- [`AUTOMATION.md`](AUTOMATION.md) — the wake/dispatch automation model.
- [`cnos.cds/README.md`](../../src/packages/cnos.cds/README.md) — what the CDS package contains.
