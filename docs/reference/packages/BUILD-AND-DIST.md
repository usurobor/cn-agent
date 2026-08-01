# Build and Distribution Model

**Version:** 0.1.0
**Status:** Implemented
**Implementation:** `src/go/internal/pkgbuild/build.go` (`cn build`), `src/go/internal/cli/cmd_build.go`, `src/go/internal/restore/restore.go` (`cn deps restore`), `src/go/internal/cli/cmd_deps.go`. The `cnos.transport.git` package and its Rust `git-cn` provider payload shown in the source-layout example do not yet exist — that provider build step is aspirational.

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
