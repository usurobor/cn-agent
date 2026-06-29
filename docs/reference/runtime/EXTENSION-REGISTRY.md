# Extension Registry

## Publishing, Trust, Installation, and App Bundles in cnos

**Version:** 1.0.0
**Status:** Draft
**Scope:** Ecosystem/distribution layer above Runtime Extensions
**Depends on:** Runtime Extensions (docs/reference/runtime/extensions/RUNTIME-EXTENSIONS.md)

---

## 0. Coherence Contract

### Gap

Runtime Extensions makes capability growth coherent at wake and execution time, but it does not yet define the ecosystem layer that makes extensions practical to distribute and operate. Without a registry/app-layer design, cnos lacks a coherent answer to:

- how packages and extensions are published
- how operators discover them
- how trust and signatures work
- how install/update/remove/rollback works
- how channels (`stable`, `beta`, `local`) are modeled
- how higher-level apps/bundles are represented
- how local installed truth relates to registry-visible available artifacts

### Mode

MCA — change the ecosystem/distribution architecture.

### α / β / γ target

- α PATTERN: one explicit ecosystem model instead of ad hoc package distribution
- β RELATION: align registry, trust, installer, Runtime Contract, doctor, and traceability
- γ EXIT: make future marketplace/app-store surfaces possible without redesigning runtime extension mechanics

### Smallest coherent intervention

Do not redesign Runtime Extensions. Build the layer above them:

- registry
- trust model
- install lifecycle
- bundle/app model

---

## 1. Relationship to Runtime Extensions

Runtime Extensions defines:

- what an extension is
- how it is discovered locally
- how it is loaded
- how it is isolated
- how it is executed
- how it is surfaced in the Runtime Contract

This document defines:

- how extension-bearing packages become available to install
- how they are trusted
- how they are selected through channels
- how they are grouped into bundles/apps
- how they move through install/update/remove/rollback lifecycle

### Core rule

The registry is not a wake-time dependency. Wake-time truth comes from:

- local installed packages
- local installed extension manifests
- local runtime contract
- local config/policy

The registry is only used during:

- discovery
- resolution
- install
- update
- rollback

This preserves the local-first invariant.

---

## 2. Core Decision

cnos should keep a minimal trusted core. Everything else should be installable through:

- packages
- extensions
- bundles/apps
- a registry/index layer

This allows:

- a small, stable core
- additive capability growth
- independent release cadence for non-core features
- richer ecosystems without core accretion

---

## 3. Concepts

### 3.1 Package

The atomic distribution unit. A package may contain:

- doctrine
- mindsets
- skills
- profiles
- extensions
- metadata

### 3.2 Extension

A capability provider inside a package. Extensions are runtime-facing and execute through the Runtime Extensions architecture.

### 3.3 Bundle / App

A higher-level installable composition of:

- packages
- extensions
- optional profiles
- optional configuration presets
- optional operator-facing metadata

A bundle/app is not a new runtime primitive. It is an install-time composition layer.

### 3.4 Registry

A signed index/discovery/trust layer through which packages and bundles become installable.

### 3.5 Channel

A named release track for registry content. Examples:

- stable
- beta
- local

### 3.6 Publisher

A named identity that signs and publishes packages/bundles. Examples:

- first-party cnos publisher
- trusted third-party org
- local developer

### 3.7 Desired state

The operator-declared intent about what should be installed.

### 3.8 Resolved state

The exact locked set of installed artifacts, versions, hashes, trust metadata, and enablement state.

---

## 4. Design Principles

### 4.1 Registry is discovery/trust, not runtime truth

The runtime should never need the network to know what capabilities exist.

### 4.2 Bundles are user-facing, packages are system-facing

Operators should be able to install meaningful units without reasoning about every package individually.

### 4.3 Trust is explicit

Nothing is trusted by accident. Publisher identity, signatures, and policy must be explicit.

### 4.4 Installed state is deterministic

Installed layout and resolved state must be reproducible.

### 4.5 Runtime sees installed truth, not marketplace noise

At wake, the agent should know:

- what is installed
- what is enabled
- what is active

It should not depend on knowing what is merely available remotely.

### 4.6 Removal and rollback are first-class

An ecosystem without clean disable/remove/rollback is not coherent.

---

## 5. Registry Model

The registry is a signed index over available packages and bundles.

### 5.1 Registry contents

A registry may contain:

- publisher records
- package records
- bundle records
- channels
- signatures / hashes
- compatibility metadata

### 5.2 Example registry entry shape

```json
{
  "schema": "cn.registry.v1",
  "publishers": [
    {
      "name": "cnos",
      "id": "publisher:cnos",
      "public_key": "ed25519:..."
    }
  ],
  "packages": [
    {
      "name": "cnos.net.http",
      "version": "1.0.0",
      "publisher": "publisher:cnos",
      "channels": ["stable"],
      "source": "git+https://example.com/cnos-net-http.git",
      "rev": "abc123...",
      "subdir": "packages/cnos.net.http",
      "integrity": "sha256:...",
      "engines": { "cnos": ">=3.13.0 <4.0.0" },
      "signature": "sig:..."
    }
  ],
  "bundles": [
    {
      "name": "web-observer",
      "version": "1.0.0",
      "publisher": "publisher:cnos",
      "channels": ["stable"],
      "packages": [
        { "name": "cnos.core", "version": "^1.0.0" },
        { "name": "cnos.net.http", "version": "^1.0.0" }
      ],
      "signature": "sig:..."
    }
  ]
}
```

### 5.3 Registry transport

A registry may be distributed through:

- git
- static hosted JSON
- local filesystem index
- future remote API

The architecture should not depend on one transport form.

### 5.4 Registry is cacheable

Registry data may be cached locally. Cached registry data is still distinct from installed runtime truth.

---

## 6. Trust Model

### 6.1 Trust classes

The ecosystem must distinguish at least:

- first-party signed
- third-party signed
- local unsigned/dev
- rejected / untrusted

### 6.2 Trust decisions

Operator/runtime policy decides which trust classes are allowed.

### 6.3 Integrity model

Minimum integrity layers:

- signed metadata
- content hash
- publisher identity
- compatibility constraints

### 6.4 Trust is orthogonal to enablement

A package can be:

- trusted but disabled
- installed but disabled
- trusted and enabled
- rejected and never installed

These distinctions must remain visible.

---

## 7. Channels

### 7.1 Standard channels

- stable
- beta
- local

### 7.2 Meaning

- stable — suitable for general use
- beta — preview / lower confidence
- local — operator/developer-only / no shared trust assumption

### 7.3 Resolution rule

Desired state may request:

- exact version
- semver range
- preferred channel

The resolver chooses a matching artifact according to:

- trust policy
- channel policy
- compatibility
- lock/update strategy

---

## 8. Bundle / App Model

### 8.1 Why bundles/apps exist

Most operators should not have to install raw packages one by one.

### 8.2 Bundle contents

A bundle may declare:

- packages
- extension-bearing packages
- optional profile defaults
- optional runtime defaults
- optional operator-facing metadata

### 8.3 Bundle is not a runtime primitive

The runtime still reasons in:

- installed packages
- installed extensions
- enabled capabilities

A bundle is resolved into those lower-level primitives.

### 8.4 Example bundle manifest

```json
{
  "schema": "cn.bundle.v1",
  "name": "web-observer",
  "version": "1.0.0",
  "packages": [
    { "name": "cnos.core", "version": "^1.0.0" },
    { "name": "cnos.net.http", "version": "^1.0.0" }
  ],
  "defaults": {
    "profile": "engineer"
  }
}
```

---

## 9. Installation Lifecycle

### 9.1 Desired state

The operator declares desired installs:

- packages
- bundles/apps
- channels / version preferences

### 9.2 Resolution

The installer resolves desired state against:

- registry
- trust policy
- channel policy
- compatibility
- dependency/conflict rules

### 9.3 Fetch

Artifacts are fetched from their declared source.

### 9.4 Verify

Before install:

- signatures checked
- integrity hashes checked
- compatibility checked

### 9.5 Materialize

Resolved packages are materialized into:

- `.cn/vendor/packages/<pkg>@<ver>/...`

Bundles may optionally materialize metadata under:

- `.cn/vendor/bundles/<bundle>@<ver>/...`

### 9.6 Enable

Installed artifacts may still be disabled by runtime/operator policy.

### 9.7 Update

Updates re-resolve against the registry and produce a new resolved state.

### 9.8 Remove

Remove the installed artifact only when:

- no desired state references it
- and no other installed artifact depends on it

### 9.9 Rollback

Rollback restores a previous known-good resolved state.

---

## 10. Desired vs Resolved State

### 10.1 Desired state

Human/operator declared intent. Examples:

- `.cn/deps.json`
- future bundle/app install intent

### 10.2 Resolved state

Machine-written exact state. Examples:

- `.cn/deps.lock.json`
- future bundle lock metadata

### 10.3 Rule

The runtime trusts only resolved local state at wake. The registry is not consulted during wake.

---

## 11. Dependency and Conflict Resolution

### 11.1 Dependencies

Packages and bundles may depend on:

- packages
- versions
- engine compatibility

### 11.2 Conflicts

The resolver must reject:

- incompatible versions
- duplicate active op kinds
- incompatible engines
- conflicting bundle/package requirements
- incompatible trust/channel constraints

### 11.3 Determinism

Given the same:

- desired state
- registry view
- trust/channel policy

resolution must produce the same resolved state.

### 11.4 No silent shadowing

Conflicts must be explicit and visible.

---

## 12. Install Roots and Caches

### 12.1 Installed roots

- packages: `.cn/vendor/packages/<pkg>@<ver>/`
- bundles (optional metadata): `.cn/vendor/bundles/<bundle>@<ver>/`

### 12.2 Cache roots

A separate local cache may exist for:

- fetched registry metadata
- fetched package artifacts

The cache is not the same thing as the installed roots.

### 12.3 Rule

Runtime discovery works against the installed roots, not cache roots.

---

## 13. Runtime Contract Integration

The Runtime Contract surfaces installed truth, not registry possibilities.

### 13.1 Cognition

Installed packages and installed extensions appear under cognition.

### 13.2 Body

Enabled/active extension-provided capabilities appear under body.

### 13.3 Optional provenance

The Runtime Contract may include:

- publisher
- trust class
- channel

for installed artifacts if useful for operator or agent behavior.

### 13.4 Bundle visibility

Installed bundle names may be surfaced for operator clarity, but actual capability surface is still derived from installed package/extension state.

---

## 14. Doctor Integration

cn doctor should validate:

- desired state syntax
- resolved state syntax
- registry compatibility (when registry or cache is available)
- signature/integrity verification
- installed roots vs lockfile
- enable/disable consistency
- conflict-free effective op set

Doctor should explain:

- what is installed
- what is enabled
- what is rejected
- why

---

## 15. Traceability Integration

The ecosystem layer should emit lifecycle events such as:

- registry.refreshed
- registry.unavailable
- package.resolved
- package.installed
- package.updated
- package.removed
- bundle.installed
- bundle.updated
- bundle.removed
- extension.enabled
- extension.disabled
- artifact.verified
- artifact.rejected

These are distinct from runtime execution events.

---

## 16. Security and Policy

### 16.1 Minimal trusted core

The ecosystem layer must not enlarge the trusted runtime core unnecessarily.

### 16.2 No ambient trust

Artifacts are not trusted merely because they appear in a registry.

### 16.3 Policy controls

The system must support policy over:

- allowed publishers
- allowed channels
- signed vs unsigned artifacts
- allowed extension classes
- network-enabled packages/extensions

### 16.4 Disable is first-class

Disable is not the same as remove. Both must be supported.

---

## 17. Marketplace / App Store Surface

A future marketplace is a product/UI surface over the registry model. The architectural mechanism remains:

- registry
- trust
- install lifecycle
- bundles/apps
- local installed state

The marketplace is therefore a presentation layer over the registry, not a separate ontology.

---

## 18. Alternatives Considered

### 18.1 No registry, direct git URLs only

**Pros:** simple, aligned with cnos/git lineage
**Cons:** weak discovery, poor trust ergonomics, no channel model, poor app-style install experience
**Decision:** insufficient as long-term ecosystem model

### 18.2 Registry without bundles/apps

**Pros:** simpler, package-centric
**Cons:** poor operator ergonomics, user-facing install units remain too low-level
**Decision:** packages alone are not enough

### 18.3 Marketplace-first design

**Pros:** product-facing clarity
**Cons:** risks conflating UI with architecture, obscures the runtime/install model
**Decision:** registry first, marketplace later

---

## 19. Migration Strategy

### 19.1 Phase 1

- define registry metadata model
- define trust/signature hooks
- define desired/resolved state interplay
- support package installs via registry records

### 19.2 Phase 2

- add bundle/app manifests
- add install/update/remove/rollback semantics
- integrate with doctor and traceability

### 19.3 Phase 3

- add marketplace/app-store surface
- support richer publisher and channel workflows

---

## 20. Acceptance Criteria

This is done when:

1. A package/extension can be published and installed without core runtime code changes.
2. A bundle/app can be installed as one unit.
3. Trust and channel state are explicit.
4. Install/update/remove/rollback are coherent and reversible.
5. Installed, enabled, disabled, and rejected states are visible in Runtime Contract, doctor, and traceability.
6. The runtime still wakes from local installed state only.

---

## 21. Summary

Runtime Extensions makes capability growth coherent at runtime. Extension Registry makes that growth usable at ecosystem scale.

The key stack is:

1. **Package** — atomic distribution unit
2. **Extension** — capability provider inside a package
3. **Bundle / App** — higher-level installable composition
4. **Registry** — discovery, trust, channels, and lifecycle

This lets cnos keep:

- a minimal coherent core
- local wake-time truth
- extensible capabilities
- and a future marketplace/app ecosystem built on stable architectural primitives
