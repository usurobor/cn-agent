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

**One precondition is still open:** the canonical dispatch labels
(`dispatch:cell` / `protocol:cds` / `status:todo`) are not yet installed by
any command — that mechanism is tracked separately (cnos#493). Until it
ships, `--dispatch cds` still renders the workflow but exits non-zero naming
cnos#493, and you must apply the labels to your repo manually.

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
The `cn` binary is just a file on your `PATH`; delete it wherever `install.sh`
placed it.

## Troubleshooting

| Symptom | Cause / fix |
|---|---|
| `✗ cn repo install must be run inside a Git repository.` | Run it from inside a checked-out Git repo (`git init` at minimum) — `cn repo install` never walks up looking for one or scaffolds a repo for you. |
| `package(s) not found in index` | The requested `--packages` entry isn't published in the resolved release/index. Check spelling, or pin `--release` to a tag that publishes it. |
| `package(s) have multiple versions in index; pass --release to pin one` | You passed `--index` pointing at a multi-version index with no `--release`; add `--release <tag>` to disambiguate. |
| `cn: command not found` after install | `install.sh` put `cn` outside your `PATH`. Re-run with `BIN_DIR="$HOME/.local/bin"` and add that dir to `PATH`. |
| `--dispatch cds` fails with "canonical dispatch labels not ensured: cnos#493 ..." | Expected today — the workflow still renders; apply the dispatch labels to your repo manually until cnos#493 ships. See [§ Autonomous dispatch](#autonomous-dispatch-opt-in). |
| `--dispatch cds` fails with "--workflow-pat-secret is required for --agent ..." | Pass `--workflow-pat-secret <NAME>` (and, for a non-sigma agent, `--bot-name`/`--bot-id`) naming the GitHub Actions secret holding that agent's workflow-scoped PAT. |

## Related

- [`docs/reference/cli/CLI.md`](../reference/cli/CLI.md) — CLI reference.
- [`docs/development/design/cn-repo-install-MOCKS.md`](../development/design/cn-repo-install-MOCKS.md) — the design surface `cn repo install` was built against.
- [dispatch orchestrator skill](../../src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md) — the autonomous dispatch loop (Layer 2).
