# Build and Distribution Model

**Version:** 0.1.0
**Status:** Draft
**Issue:** #219, #186

## Core Design Decisions

### Independent package versioning

- Each package has its own semantic version
- The cnos binary has its own version
- Compatibility is expressed explicitly via `engines.cnos` in `cn.package.json`
- `cnos.core` can be 1.4.0, `cnos.eng` can be 0.9.0, `cnos` binary can be 3.47.0

### Artifact format

Plain tarballs: `<name>-<version>.tar.gz`

- `cnos.core-1.4.0.tar.gz`
- `cnos.eng-0.9.0.tar.gz`
- `cnos.transport.git-0.1.0.tar.gz`

Why:

- every machine already knows how to fetch and unpack them
- no custom archive reader needed
- easy to inspect manually
- checksums fit naturally
- the package manifest inside the tarball carries cnos-specific identity

A custom extension like `.cnar` can come later if the artifact format ever becomes meaningfully different. Right now it would add a concept without buying much. Keep the format standard, put the cnos semantics in the manifest and index.

## Directory Structure

### Source (authored)

```text
src/
  packages/
    cnos.core/
      cn.package.json
      doctrine/
      mindsets/
      skills/
      templates/
    cnos.eng/
      cn.package.json
      skills/
    cnos.cdd/
      cn.package.json
      skills/
      commands/
      orchestrators/
    cnos.agent/
      cn.package.json
      skills/
      commands/
      orchestrators/
    cnos.hub/
      cn.package.json
      skills/
      commands/
      orchestrators/
    cnos.net.http/
      cn.package.json
      extensions/
        cnos.net.http/
          cn.extension.json
          schemas/
          docs/
    cnos.transport.git/
      cn.package.json
      commands/
        git-cn/
      extensions/
        cnos.transport.git/
          cn.extension.json
          schemas/
          docs/
      providers/
        git-cn/
          Cargo.toml
          src/
          tests/
```

### Distribution (publishable artifacts)

```text
dist/
  packages/
    index.json
    checksums.txt
    cnos.core-1.4.0.tar.gz
    cnos.eng-0.9.0.tar.gz
    cnos.cdd-0.4.0.tar.gz
    cnos.agent-0.3.0.tar.gz
    cnos.hub-0.2.0.tar.gz
    cnos.net.http-0.1.0.tar.gz
    cnos.transport.git-0.1.0.tar.gz
  bin/
    cn-linux-x64
    cn-linux-arm64
    cn-macos-x64
    cn-macos-arm64
```

### Installed (active hub state)

```text
.cn/
  vendor/
    packages/
      cnos.core/
      cnos.eng/
      cnos.cdd/
      cnos.agent/
      cnos.hub/
      cnos.net.http/
      cnos.transport.git/
```

### Key points

- `src/packages/` is the only authored source of truth for packages
- `dist/packages/` is the publishable artifact output
- `.cn/vendor/packages/` is the installed active state on a hub
- `.cn/` is not part of the source tree design — it only exists in a working hub, including when cnos itself is being developed by agents
- cnos-as-a-project is treated the same way as any other project

## What `cn build` does

### `cn build` (default)

For each package in `src/packages/`:

1. Read `cn.package.json`
2. Validate:
   - package name
   - package version
   - content classes present
   - `engines.cnos` compatibility field shape
3. Collect declared package contents from that package's source tree
4. If the package contains implementation payloads that need building, build them first (e.g. `cnos.transport.git/providers/git-cn/`)
5. Stage the package into a temporary build directory
6. Create the final tarball in `dist/packages/`
7. Compute SHA-256
8. Update `dist/packages/index.json`
9. Update `dist/packages/checksums.txt`

`cn build` does two jobs: **validation** and **artifact production**.

### `cn build --check`

Same traversal, but instead of publishing artifacts it verifies that:

- the package is structurally valid
- everything required can be staged
- implementation payloads build if needed
- the resulting artifact would be valid

This is the fast "is the package buildable and coherent?" path. Powers the I1 CI check.

### `cn build clean`

Remove build artifacts from `dist/packages/`, leave source untouched.

## What `cn deps restore` does

On a hub:

1. Read lock/config
2. Read `dist/packages/index.json` or the configured package index
3. Resolve name + version
4. Download tarball
5. Verify checksum
6. Extract into `.cn/vendor/packages/<name>/`
7. Validate embedded `cn.package.json`
8. Rebuild runtime registries as needed

Important: the artifact filename carries the version (`cnos.core-1.4.0.tar.gz`), the installed path does not (`.cn/vendor/packages/cnos.core/`). This keeps the active runtime simple while making artifacts and indexes versioned and explicit.

## Flow

```
src/packages/<name>/              → authored package source
dist/packages/<name>-<version>.tar.gz  → distributable artifact
.cn/vendor/packages/<name>/       → installed active package
```

## Version model

```
package version ≠ cnos binary version
```

The package declares compatibility. The runtime enforces it.

## Migration from current layout

### Current state (on main today)

- `src/agent/<class>/` — authored skill/doctrine/template content
- `packages/<name>/cn.package.json` — manifests (source) + copied content (derived)
- `packages/` is both source (manifests) and build output (copied content from `src/agent/`)
- Manual sync required: edit in `src/agent/`, copy to `packages/`, update manifest if new
- I1 CI check (`cn build --check`) catches drift after the fact
- `scripts/build-packages.sh` handles tarball creation for releases
- Package versions locked to binary version

### Why this breaks

The manual sync is a reliability leak. Every skill edit requires copying to the right package directory. Humans forget. I1 catches it, but after the push. The `packages/` directory mixes authored manifests with derived content — you can't tell what's source and what's build output by looking at it.

### Target state

- `src/packages/<name>/cn.package.json` — manifest (source)
- `src/packages/<name>/<class>/` — content (source, authored here, not copied)
- `dist/packages/<name>-<version>.tar.gz` — tarball (build output)
- `dist/packages/index.json` — package index (build output)
- `dist/packages/checksums.txt` — checksums (build output)
- No `packages/` directory at all
- No manual sync step
- No drift possible — content is authored in `src/packages/`, not copied from elsewhere
- Each package has its own version
- `engines.cnos` declares compatibility range

### What disappears

- The `packages/` directory
- The `src/agent/` → `packages/` copy step
- The entire manual sync problem (and the I1 failure class that catches it)
- `scripts/build-packages.sh` (replaced by `cn build`)
- Version lockstep between packages and binary

### Migration steps

1. **Merge manifests and content into `src/packages/`**
   - For each package, move `packages/<name>/cn.package.json` to `src/packages/<name>/cn.package.json`
   - Move content from `src/agent/<class>/<id>/` into `src/packages/<name>/<class>/<id>/` per the manifest's `sources` declarations
   - Remove the `sources` field from manifests (content is now co-located, not referenced from elsewhere)
   - Delete `packages/` directory
   - Delete `src/agent/` (content now lives in `src/packages/`)

2. **Add `dist/` to `.gitignore`** *(accomplished: #266, DESIGN-266)*
   - `dist/` is build output, never committed
   - Removing committed dist/ eliminates the main-side edit vs. PR rebase race described in #266. See `DESIGN-266-dist-out-of-git.md` for the rationale and alternatives rejected.

3. **Update `cn build` to target new layout**
   - Read from `src/packages/<name>/cn.package.json`
   - Package content from `src/packages/<name>/<class>/`
   - Produce tarballs in `dist/packages/`
   - Generate `dist/packages/index.json` and `dist/packages/checksums.txt`

4. **Update CI**
   - I1: `cn build --check` validates `src/packages/` structure and manifest completeness
   - Release: `cn build` produces `dist/packages/` tarballs, uploads to GitHub release

5. **Update `cn deps restore`**
   - Reads package index (from dist or remote URL)
   - Downloads tarballs
   - Extracts to `.cn/vendor/packages/<name>/`

6. **Decouple package versions from binary version**
   - Each `cn.package.json` gets its own `version`
   - `engines.cnos` declares compatible kernel range
   - `scripts/stamp-versions.sh` no longer forces all packages to the binary version

### Sequencing

Steps 1–2 can land as one PR (the big layout move). Steps 3–4 depend on the Go `cn build` being updated. Steps 5–6 follow naturally. The whole migration is #186.
