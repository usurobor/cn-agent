# Package Content Classes

## What the build system distributes via packages

**Version:** 2.0.0
**Status:** Draft
**Addresses:** #119 (template distribution), build-system content model
**Related:**
- docs/beta/architecture/PACKAGE-SYSTEM.md (retired — redirect stub)
- docs/alpha/cognitive-substrate/COGNITIVE-SUBSTRATE.md (asset taxonomy)
- docs/alpha/runtime-extensions/RUNTIME-EXTENSIONS.md (extension model)

---

## 0. Purpose

This document defines the set of content classes that `cn build` knows how
to assemble from `src/packages/` into distributable packages. It covers only
the build-system content pipeline -- not resolution, registry,
or other package-system layers (see §2.4 below for the distribution model).
It answers:

1. What kinds of content can a package distribute?
2. How does each content class flow from source to installed package?
3. Why is the set explicit rather than generic?
4. What is the cost model for adding a new content class?

---

## 1. Content Classes

A package content class is a named category of assets that lives inside
a package's source directory at `src/packages/<name>/<class>/`.

### 1.1 Content classes

| Class | Copy mode | Source location | Runtime role |
|-------|-----------|-----------------|--------------|
| doctrine | wildcard or named | `src/packages/<name>/doctrine/` | Always-on core principles |
| mindsets | named files | `src/packages/<name>/mindsets/` | Always-on behavioral frames |
| skills | directory trees | `src/packages/<name>/skills/` | Keyword-scored, bounded |
| extensions | directory trees | `src/packages/<name>/extensions/` | Runtime capability providers |
| templates | named files | `src/packages/<name>/templates/` | Identity/config scaffolding for new hubs |
| commands | directory trees | `src/packages/<name>/commands/` | Operator-facing CLI commands |
| orchestrators | directory trees | `src/packages/<name>/orchestrators/` | Mechanical workflows (`cn.orchestrator.v1`) |
| katas | directory trees | `src/packages/<name>/katas/` | Executable verification scenarios bundled with the artefact they prove |
| (metadata) | implicit | `src/packages/<name>/` | `cn.package.json` — package identity, version, engine constraint |

All eight content classes are co-located with their package manifest.
`cn build` assembles each package from `src/packages/<name>/` into
a tarball in `dist/packages/`. `cn deps restore` installs from dist
into `.cn/vendor/packages/<name>/` on a hub. Content is authored
in-place — no copy step, no drift possible.

The `commands` class: each command id has a directory at
`src/packages/<name>/commands/<id>/` containing an entrypoint script
`cn-<id>` (executable bit preserved). Per-command metadata (entrypoint
path and summary) lives at the top level of `cn.package.json` under
a `commands` object.

The `orchestrators` class: each orchestrator id has a directory at
`src/packages/<name>/orchestrators/<id>/orchestrator.json`. The JSON
is a `cn.orchestrator.v1` manifest with steps, permissions, and a
trigger — see `docs/alpha/agent-runtime/ORCHESTRATORS.md` for the
full schema.

Metadata (`cn.package.json`) is not a content class — it is always
present at the package root.

### 1.2 Copy modes

The build system uses two copy strategies:

- **Individual file copy**: Used by doctrine (when named), mindsets, and templates.
  Each declared entry is a single file name relative to the package's class directory.

- **Directory tree copy**: Used by skills, extensions, commands, orchestrators, and katas.
  Each declared entry is a directory path. The entire subtree is copied recursively.

Doctrine has a special case: `"*"` copies all `.md` files (wildcard mode).

---

## 2. Content Class Pipeline

Each content class flows through the same four-stage pipeline.
Angle-bracket segments are substitution variables, not literal paths.

```
src/packages/{name}/{class}/{entry}
    |
    | cn build         (assemble into tarball)
    v
dist/packages/{name}-{version}.tar.gz
    |
    | cn deps restore  (fetch + materialize into hub)
    v
.cn/vendor/packages/{name}/{class}/{entry}
    |
    | runtime read     (loaded at wake or on-demand)
    v
agent context / hub state
```

Concrete example for the templates content class:

```
src/packages/cnos.core/templates/SOUL.md
    |  cn build
    v
dist/packages/cnos.core-1.0.0.tar.gz  (contains templates/SOUL.md)
    |  cn deps restore
    v
.cn/vendor/packages/cnos.core/templates/SOUL.md
    |  read_template
    v
spec/SOUL.md (in hub)
```

### 2.1 Build (cn build)

- Reads each `src/packages/<name>/cn.package.json`
- Assembles content into a tarball in `dist/packages/`
- Generates `dist/packages/index.json` and `dist/packages/checksums.txt`

### 2.2 Check (cn build --check)

- Validates package structure and manifest completeness
- Verifies all declared content is present and buildable

### 2.3 Clean (cn build clean)

- Removes `dist/packages/` build output
- Leaves `src/packages/` untouched

### 2.4 Install (cn deps restore)

First-party packages are distributed as versioned tarball artifacts.
GitHub Releases is the current publishing/hosting surface, not the
long-term consumer endpoint — the package index abstracts the URL,
so hosting can move without changing any consumer. The restore flow is:

1. Read the lockfile (`.cn/deps.lock.json`) — name + version + sha256 per package
2. Look the name+version up in the package index (`packages/index.json`) for a URL
3. Download `<name>-<version>.tar.gz` over HTTPS
4. Verify the SHA-256 against the lockfile entry
5. Extract into `.cn/vendor/packages/<name>@<version>/`
6. Validate the extracted `cn.package.json`

The package index is the resolution authority (name+version → URL);
the lockfile is the integrity authority (name+version → sha256).
Hosting can move by republishing a new index without touching any
hub's lockfile.

When `cn deps restore` runs inside a cnos checkout, first-party
packages can be installed from the local `dist/packages/` after
`cn build`. Consumers always go through the index.

---

## 3. Manifest Schema

The package manifest declares identity, version, and command metadata.
Content classes are discovered by directory presence — content is
co-located with the manifest, so no `sources` map is needed.

```json
{
  "schema": "cn.package.v1",
  "name": "cnos.core",
  "version": "1.0.0",
  "kind": "package",
  "engines": { "cnos": ">=3.50.0" },
  "commands": {
    "daily": {
      "entrypoint": "commands/daily/cn-daily",
      "summary": "Create or show today's daily reflection thread"
    },
    "weekly": {
      "entrypoint": "commands/weekly/cn-weekly",
      "summary": "Create or show the current weekly reflection thread"
    },
    "save": {
      "entrypoint": "commands/save/cn-save",
      "summary": "Commit + push"
    }
  }
}
```

`cn build` discovers content classes by scanning for known directory
names (`doctrine/`, `mindsets/`, `skills/`, `templates/`, `commands/`,
`orchestrators/`, `extensions/`, `katas/`) inside each package's
source directory. A package includes only the classes it has
directories for. The canonical list is defined once in
`src/go/internal/pkg/pkg.go` (`pkg.ContentClasses`) and shared by
`cn build --check` and `cn status` via `pkgbuild.FindContentClasses`.

---

## 4. Design Rationale

### 4.1 Why explicit classes, not a generic tree model

The content classes are an explicit, closed set rather than a generic
extensible model.

**Why explicit is preferred:**

- **Semantics**: Different classes have different copy modes (file vs tree vs wildcard). These are not interchangeable.
- **Readability**: Each class is a known directory name. No configuration needed.
- **Debugging**: A mismatch in one class is easy to identify.
- **Current scale**: 8 classes. The explicit model is not a maintenance burden.

**When to reconsider:**

If adding a new content class feels like boilerplate rather than design,
refactor to a generic model.

### 4.2 Why templates are a content class (not special-cased)

Templates became a content class (rather than being special-cased) because:
- They follow the same source -> build -> install -> read path
- They benefit from the same drift detection (`cn build --check`)
- They can be overridden by future third-party packages
- No special cases in the build, install, or restore pipelines

### 4.3 Why the set is closed (for now)

The current 8 content classes cover all shipped cognitive assets plus
identity templates and bundled verification katas. The set is intentionally closed:

- New classes should be added only when an existing asset type cannot
  be served by any current class
- Each new class should follow the existing pattern (add field, add
  copy logic, add check/clean support)
- The design doc (this file) should be updated when the set changes

---

## 5. Tradeoffs

### Explicit classes vs generic model

| Dimension | Explicit (current) | Generic |
|-----------|--------------------|---------|
| Adding a class | ~5 code locations | 0 code, 1 config |
| Typo protection | Compile-time | Runtime/none |
| Copy mode per class | Type-safe branches | String dispatch |
| Readability | High (named fields) | Lower (dynamic) |
| Maintenance at N=6 | Low | Unnecessary |
| Maintenance at N=12 | Noticeable | Justified |

### Templates as content class vs special case

| Dimension | Content class (current) | Special case |
|-----------|------------------------|--------------|
| Drift detection | Automatic (cn build --check) | Must build separately |
| Install path | Same as all other content | Custom install logic |
| Third-party override | Possible | Not without custom code |
| Code complexity | Zero new abstractions | New special-case path |

---

## 6. History

| Version | Change |
|---------|--------|
| v3.4.0 | Package system introduced (doctrine, mindsets, skills) |
| v3.18.0 | Extensions added as 4th content class |
| v3.24.0 | Templates added as 5th content class (#119) |
| v3.37.0 | Commands migrated from built-in to package content class |
| v3.51.0 | Content co-located with manifests in `src/packages/`; `sources` field removed |
| v3.55.0 | Katas added as 8th content class (kata framework split; `cnos.kata`, `cnos.cdd.kata`) |
| v3.56.1 | Single source of truth in `pkg.ContentClasses`; `cn status` and `cn build --check` agree on membership via filesystem presence (#253) |

---

## 7. Command Discovery

CLI commands resolved by `cn <name>` follow a fixed precedence order. The first match wins; later layers are not consulted.

1. **Built-in commands** — compiled into the `cn` binary (`cn status`, `cn doctor`, `cn deps`, ...). Authoritative.
2. **Repo-local commands** — executable files at `.cn/commands/cn-<name>` inside the current hub. Discovered at dispatch time. Used for hub-specific automation that should not be packaged.
3. **Vendored package commands** — declared under the `commands` content class of any installed package in `.cn/vendor/packages/<name>@<version>/`. Resolved by walking each installed package's `cn.package.json` and matching `commands.<name>`.

There is no PATH fallback. If two external commands at the **same precedence level** claim the same name, that is a `cn doctor` error (the package author or hub operator must rename). Built-ins always shadow external commands silently — built-ins are part of the system surface.

External commands are dispatched via `exec`: argv (minus `cn`) is passed through, the working directory is the hub root, and a small set of environment variables (`CN_HUB_PATH`, `CN_PACKAGE_ROOT` for package commands) is exported so the script can locate hub state without re-discovering it.

`cn help` lists all three layers with source and summary. `cn doctor` validates package command integrity: duplicate names within a layer, missing entrypoints, non-executable command files, and malformed metadata are all errors.

Runtime capability extensions (see `RUNTIME-EXTENSIONS.md`) are a **separate** system and unaffected by command discovery. Extensions provide typed capability providers for the agent runtime; commands are operator-facing CLI entry points. Both can ship in the same package without overlap.
