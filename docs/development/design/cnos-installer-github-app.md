# CNOS Installer — GitHub App (native install surface)

> **Operator ruling (this doc's mandate).** For installing CNOS/CDS into an
> **existing** repository, the polished, GitHub-native install surface is a
> **GitHub App that opens a pull request**, not a bootstrap GitHub Action.
> The Action path (`docs/guides/templates/cnos-install.yml`) remains as the
> **zero-backend fallback**. The App is App-first; the Action is the
> escape hatch.

**Status:** Design surface (pre-implementation for the App backend). The
deterministic install *engine* (`cn repo install`) already ships and is
proven end-to-end (release 3.82.3). This doc locks the App architecture,
the per-tier permission model, and the exact division of labor between the
App, `cn`, and Actions.

**Audience:** whoever builds/registers the CNOS Installer App, and reviewers
of that work.

---

## Why an App (and why not the alternatives)

| Surface | What it's good for | Why it's not the polished install UX |
|---|---|---|
| **GitHub App** (this doc) | Native install into an **existing** repo: repo selection + permission grant are first-class; the App opens a reviewable PR; can sync labels; can write workflows/secrets only with explicit, separately-granted scopes. | — (this is the recommendation) |
| **Bootstrap Action** (`cnos-install.yml`) | Zero backend, repo-owned, no external trust. | **Bootstrap problem:** the user must first commit a `workflow_dispatch` workflow to the **default branch** before GitHub will show the "Run workflow" button. "Install the thing that installs CNOS" before CNOS. Keep as fallback. |
| **Repository template** | *Starting a new* CNOS repo. | Templates create **new** repos with unrelated history. They do not retrofit an existing repo, open a PR into existing history, or sync labels/secrets/workflows. Wrong problem. |
| **Marketplace Action** | A reusable building block (`uses: usurobor/cnos-install-action@v1`). | Still an Action — the repo still needs a workflow file that calls it. Same bootstrap problem. Fine as a component, not as the entry UX. |

## Clean architecture

```
GitHub App    =  native installation surface   (repo selection, permissions, PR)
cn repo install =  deterministic installer engine (already shipped, 3.82.3)
GitHub Action  =  optional execution substrate / fallback
```

The App never reimplements install logic. It **invokes the same engine**
(`cn repo install`) and turns the result into a PR. `cn` stays the single
source of installer truth; the App is a delivery surface around it.

## Install flow (base)

1. User clicks **"Install CNOS"** → GitHub opens `/apps/cnos/installations/new`.
2. User selects account/org and **selected repositories**.
3. GitHub shows the App's requested permissions; user reviews and installs.
4. The App receives the `installation` (and `installation_repositories`) webhook.
5. The App runs the base-install engine against each selected repo and
   **opens a PR** — it does not push to the default branch.
6. User reviews the PR (it carries a receipt) and merges.

Resulting PR:

```
branch:  cnos/install
title:   Install CNOS/CDS (base)
files:   .cn/deps.json
         .cn/deps.lock.json
         .gitignore  (+ .cn/vendor/)
```

> **Vendor tree.** `.cn/vendor/packages/**` is **not committed** — it is
> gitignored and rehydrated from `.cn/deps.lock.json` by `cn repo install`.
> The PR commits the two manifests + the `.gitignore` line, matching exactly
> what a human running `cn repo install` then `git add` would commit
> (`docs/guides/INSTALL-CDS.md`). This keeps the App's diff identical to the
> CLI's, so the two install paths are byte-for-byte reconcilable.

### PR body = receipt

The PR body is a receipt, mirroring the CLI's stdout summary:

- **release used** (resolved tag, e.g. `3.82.3`)
- **packages restored** (`name@version`, resolved from the release index —
  not the release tag; see the 3.82.3 resolver fix)
- **files written**
- **explicitly NOT installed** for base tier: no `.github/workflows/` change,
  no secrets, no labels, no autonomous write loop.

## Engine mechanism: Option A (service-side render)

Two ways the App can produce the diff:

- **Option A — service-side render (chosen for base install).** The App runs
  the same repo-install logic server-side, builds the file tree, commits
  through the GitHub Contents/Git APIs (or via a controlled Actions run using
  an installation token), and opens the PR. **Cleaner for base install** —
  it avoids putting execution logic into the tenant repo *just to install
  files*.
- **Option B — app-seeded workflow.** The App writes a temporary bootstrap
  workflow, dispatches it, and lets Actions run `cn repo install` in the
  tenant repo. Useful when the install needs repo-local execution or when we
  want to reuse runner behavior exactly. **Not** the base-install path.

### Concrete Option-A substrate (no separate always-on engine host)

The service-side render is implemented as a controlled Actions run **in
`usurobor/cnos`**, keeping Actions as the execution transport and avoiding a
bespoke build host:

- Trigger: `repository_dispatch` (type `cnos-app-install`) — see
  [`.github/workflows/cnos-app-install.yml`](../../../.github/workflows/cnos-app-install.yml).
- Payload: `{ target_repo, installation_id, release, tier }`.
- The workflow mints a **short-lived installation token** scoped to the
  *target* repo (`actions/create-github-app-token`), checks out the target,
  installs `cn` via `install.sh`, runs `cn repo install`, and opens the PR
  with `peter-evans/create-pull-request` using that installation token.

The only piece that still needs a tiny always-on endpoint is the **webhook
forwarder**: it receives the `installation` webhook and fires the
`repository_dispatch`. That is a thin, stateless function (any serverless
host) — it holds no install logic. **This is the one infra dependency the
repo cannot satisfy by itself** (see [Open infra](#open-infra-what-only-the-operator-can-do)).

## Permission model — minimal, per tier

**Do not request Workflows/Secrets permissions just to do a base install.**
Each tier adds exactly the scopes it needs and no more.

| Tier | What it does | Added GitHub App permissions |
|---|---|---|
| **1 — Install base CNOS** | Opens the base PR (`.cn/deps.json` + lock + `.gitignore`). No workflow, no secrets, no autonomous loop. | `Contents: write`, `Pull requests: write`, `Metadata: read` (implicit) |
| **2 — Install CNOS + labels** | Base, plus canonical label sync (label doctor). Files still PR-based; labels created directly via API. | **+ `Issues: write`** |
| **3 — Install autonomous dispatch** | Base, plus renders `.github/workflows/cnos-cds-dispatch.yml`; optionally writes repo secrets. Requires explicit operator approval. | **+ `Workflows: write`**, and **+ `Secrets: write`** *only if* the App sets repo secrets |

Grounding for each scope (GitHub docs behavior):

- **Workflow files.** Modifying `.github/workflows/**` via the Contents API
  needs `Contents: write` **plus `Workflows: write`** for a fine-grained/App
  token (classic PATs need the `workflow` scope). This is why Tier 3 is a
  separate grant — the base token deliberately cannot touch workflows.
- **Labels.** Label endpoints live under Issues; creating repo labels needs
  `Issues: write` (or `Pull requests: write`). **Label descriptions must be
  ≤ 100 characters** — the canonical `dispatch:cell` description was shortened
  to satisfy this (cycle/493); the label-doctor the App runs must enforce the
  same bound.
- **Secrets.** The repository Actions secrets API accepts App installation
  tokens and requires `Secrets: write`. Relevant to Tier 3 only — never base.

## Division of responsibility

| Concern | Owner |
|---|---|
| Native install surface (repo pick, permission grant, PR open) | **GitHub App** |
| Deterministic install (resolve release → deps/lock → vendor restore) | **`cn repo install`** |
| Execution transport (runners, the Option-A render job, fallback) | **GitHub Actions** |
| Human decision (review + merge the PR; approve Tier-3 upgrade) | **Repo owner** |

## MVP → next → long-term

**MVP (Tier 1):**
1. User installs App on selected repo.
2. App opens "Install CNOS base" PR (deps/lock/gitignore).
3. PR body carries the receipt (release, packages, files, "no workflow/secrets").
4. User reviews and merges.

**Next (Tiers 2–3):**
5. App offers "Enable dispatch" upgrade; requests `Workflows: write` if not
   already granted.
6. App runs label doctor (Tier 2).
7. App opens a PR adding `cnos-cds-dispatch.yml` (Tier 3).
8. App writes required repo secrets if granted `Secrets: write`, else prints
   exact UI instructions.

**Long-term:** the App becomes the native CNOS install/coherence surface;
Actions remain execution transport; `cn` remains the deterministic local
engine.

## Open infra (what only the operator can do)

The repo-side artifacts (this doc, the App manifest, the Option-A render
workflow) are landable and reviewable now. Standing the **live** App up needs
actions outside this repo/sandbox:

1. **Register the App** from the manifest
   ([`templates/cnos-app-manifest.json`](templates/cnos-app-manifest.json)) —
   org-admin action; sets name `cnos`, permissions, webhook URL, and mints the
   App id + private key.
2. **Host the thin webhook forwarder** (installation webhook →
   `repository_dispatch`). Stateless; any serverless host.
3. **Provision `usurobor/cnos` secrets** the render workflow reads:
   `CNOS_APP_ID` and `CNOS_APP_PRIVATE_KEY`. Until these exist the render
   workflow is inert (it only fires on `repository_dispatch`, which nothing
   sends yet).

None of the three block the base CLI install, which is already shipped.

## Related

- [`docs/guides/INSTALL-CDS.md`](../../guides/INSTALL-CDS.md) — human install guide (CLI + Action fallback).
- [`docs/development/design/cn-repo-install-MOCKS.md`](cn-repo-install-MOCKS.md) — the engine's design surface + receipt-parity contract.
- [`docs/guides/templates/cnos-install.yml`](../../guides/templates/cnos-install.yml) — the Action fallback (bootstrap workflow).
- [`templates/cnos-app-manifest.json`](templates/cnos-app-manifest.json) — the App-registration manifest (base tier).
- [`.github/workflows/cnos-app-install.yml`](../../../.github/workflows/cnos-app-install.yml) — the Option-A service-side render (repository_dispatch engine).
