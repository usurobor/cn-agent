# Package Authoring Guide

**Version:** 0.1.0
**Status:** Draft

How to author, structure, test, and ship a cnos package.

---

## 1. What is a package

A package is the unit of distribution in cnos. It contains authored content — skills, commands, orchestrators, templates, doctrine — that ships as a tarball and installs into a hub.

Packages are authored in `src/packages/<name>/`, built into `dist/packages/`, and installed into `.cn/vendor/packages/<name>/`.

Every package has a manifest (`cn.package.json`) that declares what it contains, what it exposes, and what it requires.

## 2. Manifest

```json
{
  "schema": "cn.package.v1",
  "name": "cnos.example",
  "version": "0.1.0",
  "kind": "package",
  "engines": {
    "cnos": ">=3.50.0"
  },
  "commands": {
    "example-cmd": {
      "entrypoint": "commands/example-cmd/cn-example-cmd",
      "summary": "One-line description"
    }
  }
}
```

### Required fields

- `schema` — always `cn.package.v1`
- `name` — package name (e.g. `cnos.core`, `cnos.cdd.kata`)
- `version` — semver
- `kind` — always `package`

### Optional fields

- `engines.cnos` — minimum cnos version constraint
- `commands` — commands the package exposes to the CLI

## 3. When to create a package

A package boundary should follow a **reason to change**, not just a topic label.

### The question

> "If this content changes, does the rest of the system need to change too?"

If yes → same package. If no → different package.

### Core vs package

Policy belongs in the kernel. Implementations belong in packages.

| Belongs in kernel/core | Belongs in a package |
|---|---|
| Command precedence rules | Command implementations |
| Discovery/dispatch logic | Discovered content |
| Doctor/status/help rendering | What doctor checks, status shows |
| Package format and schema | Package content |

A package should **never** widen its own authority or rewrite kernel policy. The kernel decides whether and how the package's content is used.

### Package splitting

Split when two parts of a package change for different reasons:

- ❌ `cnos.core` contains CDD skills — CDD process changes force a core release
- ✅ `cnos.cdd` separated — CDD evolves independently from core ops

Split when install profiles differ:

- ❌ Kata framework bundled with production skills — every hub installs test tooling
- ✅ `cnos.cdd.kata` separate — install only when you want to run katas

### What makes a good package boundary (L7)

A system-shaping package boundary eliminates a class of future work:

- Future CDD method changes don't touch core → `cnos.cdd`
- Future engineering skills don't touch CDD → `cnos.eng`
- Future kata scenarios don't touch any production package → `cnos.cdd.kata`

If the boundary only moves code around without changing what future work looks like, it's L5 convenience, not L7 architecture.

## 4. Content classes

A package may contain any combination of these content class directories:

| Directory | Content class | Description |
|-----------|--------------|-------------|
| `skills/` | Skills | Agent capabilities (SKILL.md + supporting files) |
| `commands/` | Commands | CLI commands (entrypoint scripts) |
| `orchestrators/` | Orchestrators | Multi-step execution patterns |
| `templates/` | Templates | File templates for generation |
| `doctrine/` | Doctrine | Constitutive documents |
| `mindsets/` | Mindsets | Agent behavioral configurations |
| `extensions/` | Extensions | Runtime extensions |

At least one content class directory must be present and non-empty.

### Choosing the right content class

Each runtime surface has a different contract. Do not collapse them:

| If your content... | Use | Not |
|---|---|---|
| Teaches an agent judgment or procedure | Skill | Command |
| Dispatches an action for an operator | Command | Skill |
| Orchestrates multi-step execution | Orchestrator | Command |
| Provides a runtime capability (transport, storage) | Extension | Command |
| Defines agent identity or behavior | Doctrine/Mindset | Skill |
| Generates files from a template | Template | Command |

- ❌ Trigger commands by keyword like skills (smears dispatch into cognition)
- ❌ Treat providers as hidden commands (different contracts)
- ✅ Each surface has its own registry, contract, and policy boundary

There is **no generic `docs/` or `data/` content class.** If you need to ship data with a package, bundle it inside a command directory tree or a skill directory. The content class model is finite and explicit.

## 4. Directory structure

```
src/packages/cnos.example/
  cn.package.json
  skills/
    skill-a/
      SKILL.md
    skill-b/
      SKILL.md
  commands/
    example-cmd/
      cn-example-cmd      # executable entrypoint
      helpers/             # supporting files, bundled with the command
```

### Rules

- One `cn.package.json` at the package root
- Content class directories at the package root (not nested)
- Command entrypoints must be executable files
- Command directory trees are copied intact — bundle any data the command needs inside its directory
- Skill directories must contain a `SKILL.md`

## 5. Commands

### Entrypoint contract

A command entrypoint is an executable file (usually a shell script) that receives:

| Env var | Description |
|---------|-------------|
| `CN_HUB_PATH` | Path to the active hub |
| `CN_PACKAGE_ROOT` | Path to the installed package in vendor |
| `CN_COMMAND_NAME` | Name of the command being dispatched |

Arguments are passed as positional parameters.

### Naming

- Entrypoint file: `cn-<command-name>` (e.g. `cn-kata-list`)
- Directory: `commands/<command-name>/`
- Manifest key: the command name without `cn-` prefix

### Bundling data with commands

Commands that need data (fixtures, rubrics, prompts, corpora) bundle it inside their command directory tree. The entire tree ships in the tarball and installs to vendor.

```
commands/
  kata-run/
    cn-kata-run           # entrypoint
    lib.sh                # shared helpers
    katas/                # bundled kata corpus
      method/
        M1-design/
          kata.md
          rubric.json
```

This is the correct way to distribute non-skill, non-command content today. Do not invent a new content class for data.

## 6. Skills

### Structure

```
skills/
  skill-name/
    SKILL.md              # required — defines the skill
    supporting-file.md    # optional supporting content
```

### Visibility

Skill visibility is declared in the skill's own frontmatter, not in the package manifest:

- `visibility: public` (default) — visible to agents for loading
- `visibility: internal` — used by the package's own entrypoints (e.g. CDD sub-skills loaded by the `cdd` entrypoint)

The package manifest does not carry a skill inventory. Skills are discovered from directory presence; visibility is declared per-artifact.

## 7. Katas and testing

### Where katas live

| What you're testing | Where katas go |
|---|---|
| Kernel commands (`cn build`, `cn init`, `cn status`, etc.) | `scripts/kata/` — system-level, not a package |
| A package's own capabilities | Inside the package — as commands that test the package |

### Runtime katas (system-level)

`scripts/kata/` contains executable scripts that prove the cnos package pipeline works: build, install, discover, dispatch, self-describe. These test the kernel, not any specific package.

### Package katas (package-level)

A package that wants to prove its capabilities work ships katas as commands. Example: `cnos.cdd.kata` exposes `kata-run`, `kata-list`, `kata-judge` as commands, with the kata corpus bundled inside the command directory tree.

This is not a special mechanism — it's just commands + bundled data. The package system already supports this.

### Package capability verification

Every package should be provable. The question each package must answer: **"What can the system do with this package installed that it couldn't do without it?"**

Create a companion `.kata` package for your package:

| Package | Kata package |
|---|---|
| `cnos.core` | `cnos.core.kata` |
| `cnos.cdd` | `cnos.cdd.kata` |
| `cnos.eng` | `cnos.eng.kata` |

The kata package ships commands (`kata-list`, `kata-run`, `kata-judge`) with bundled test scenarios, rubrics, and judge prompts. It proves the parent package's capabilities work.

Examples:
- `cnos.core.kata`: `cn daily` dispatches, skills are loadable, `cn weekly` creates a thread
- `cnos.cdd.kata`: CDD review produces evidence-bound findings, release produces coherence delta
- `cnos.eng.kata`: Go skill constraints are followed, design principles are applied

The kata package is separate from the parent so it can be installed independently — you can run katas without the overhead of the parent's dev tooling, and you can skip katas in production installs.

## 8. Source, artifact, and installed state

Three layers. Keep them explicit.

| Layer | Location | Who produces it | Editable? |
|---|---|---|---|
| **Source** | `src/packages/<name>/` | Human author | Yes |
| **Artifact** | `dist/packages/<name>-<version>.tar.gz` | `cn build` | No — derived |
| **Installed** | `.cn/vendor/packages/<name>/` | `cn deps restore` | No — derived |

- ❌ Edit files in `.cn/vendor/` (installed state treated as source — will be overwritten on next restore)
- ❌ Edit tarballs in `dist/` (artifact treated as source — will be overwritten on next build)
- ✅ Edit in `src/packages/`, build, restore. Source → artifact → installed.

The manifest (`cn.package.json`) is a **contract**. It declares what the package provides. The runtime holds it to that contract:
- `cn help` derives command listings from the manifest
- `cn doctor` validates that declared entrypoints and skills exist
- `cn status` derives package inventory from installed manifests
- `cn build --check` validates source structure against manifest claims

If the manifest says it, it must be true. If it's not in the manifest, the runtime doesn't know about it.

## 9. Build and distribution

### Pipeline

```
src/packages/<name>/  →  cn build  →  dist/packages/<name>-<version>.tar.gz
                                       dist/packages/index.json
                                       dist/packages/checksums.txt
```

### Validation

`cn build --check` validates all packages in `src/packages/`:
- Manifest has name + version
- At least one content class directory exists and is non-empty
- Structure matches the `cn.package.v1` schema

### Installation

```
cn deps lock      # generate lockfile from index
cn deps restore   # install from lockfile into .cn/vendor/packages/<name>/
```

Installed packages live at `.cn/vendor/packages/<name>/` (no version in path).

## 9. Naming conventions

| Convention | Example |
|---|---|
| Core packages | `cnos.core`, `cnos.cdd`, `cnos.eng` |
| Extension packages | `cnos.cdd.kata`, `cnos.transport.git` |
| Third-party packages | `<org>.<name>` (future) |

Package names use dots as namespace separators. The name in the manifest must match the directory name under `src/packages/`.

## 10. Checklist for a new package

- [ ] `cn.package.json` with schema, name, version, kind
- [ ] At least one content class directory, non-empty
- [ ] Every command entrypoint exists and is executable
- [ ] Every declared skill directory contains a SKILL.md
- [ ] `cn build --check` passes
- [ ] `cn build && cn deps restore` round-trips successfully
- [ ] Package capability is demonstrable — create a `<name>.kata` companion package

## 11. Anti-patterns

- ❌ Inventing a new content class for data — bundle inside commands
- ❌ Putting runtime katas in a package — kernel tests go in `scripts/kata/`
- ❌ Putting package-level katas in `scripts/` — package tests ship with the package
- ❌ Empty content class directories (CI will catch this)
- ❌ Declaring commands in manifest without creating the entrypoint
- ❌ Non-executable entrypoint files
- ❌ Version in vendor path (`name@version/`) — vendor path is `name/` only
- ❌ Package that widens its own authority or overrides kernel policy
- ❌ Smearing runtime surfaces (skills that dispatch, commands that choose, providers that orchestrate)
- ❌ Editing installed state in `.cn/vendor/` — it's derived, not source
- ❌ Splitting by topic label instead of reason to change
- ❌ Manifest declares content that doesn't exist on disk
