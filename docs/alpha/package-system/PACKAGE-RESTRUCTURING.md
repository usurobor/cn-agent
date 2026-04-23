# Package Restructuring — Target Structure

**Version:** 0.1.0
**Status:** Draft
**Issue:** #186
**Mode:** MCA
**Active Skills:** cdd/design, eng/evolve

## Problem

The current package boundary is uneven:

- **cnos.core** is carrying almost everything: doctrine, 9 mindsets, 17 agent/* skills, 7 cdd/* skills, 4 ops/* skills, one extension (cnos.net.http), templates, one orchestrator (daily-review), and 3 commands (daily, weekly, save). That is too much for something called "core."
- **cnos.eng** is already coherent: it contains the engineering bundle, including evolve, code, write-functional, ocaml, go, typescript, test, writing, and related engineering skills.
- **cnos.pm** is thin and overlaps the stronger CDD bundle: it only contains pm/follow-up, pm/issue, and pm/ship.
- The package system itself is now explicit that packages can carry skills, extensions, templates, commands, and orchestrators, all through the same source → build → install pipeline.

The best structure makes package boundaries follow install/use cohesion, not history.

## Target Structure

| Package | What it should contain | Why |
|---------|----------------------|-----|
| **cnos.kernel** | Embedded built-in bootstrap commands only: help, init, setup, deps, build, doctor, status, update | The trusted minimal kernel should stay tiny and always available. |
| **cnos.core** | Doctrine, mindsets, templates, and only the minimal bootstrap cognition skills | This becomes the lean substrate every hub needs. |
| **cnos.agent** | The broader agent/* bundle, plus default agent commands/orchestrators like daily, weekly, save, daily-review | This is the missing package today: reusable agent behavior should not bloat core. |
| **cnos.cdd** | All cdd/* skills and future CDD-specific commands/orchestrators | CDD is already a coherent method bundle and should be its own package. |
| **cnos.eng** | All eng/* skills | Already coherent; keep it. |
| **cnos.hub** | All ops/* skills and future peer/inbox/A2A workflow commands/orchestrators | Hub-facing operational practice is distinct from both core cognition and engineering. |
| **cnos.net.http** | The HTTP extension package | Extensions should not live in core if they are optional capability providers. |
| **cnos.pm** | Retire or redefine later | In its current shape it is too small and overlaps CDD. |

## Package Rationale

### cnos.kernel

This is not a normal installed package. It is the embedded bootstrap kernel. It should stay very small so the platform can always:

- install/restore packages
- discover commands
- show help/status/doctor
- update itself

That matches the kernel-command design direction (GO-KERNEL-COMMANDS.md, #216).

### cnos.core

Core should mean:

- always-on principles
- baseline identity/config scaffolding
- enough cognition to wake safely

**Keep in cnos.core:**

- doctrine
- mindsets
- templates
- only the truly bootstrap agent skills

**Move out:**

- cdd/*
- ops/*
- the HTTP extension
- the default rituals/orchestrator
- most of the heavier agent/* surface

That is the biggest immediate coherence gain, because the current cnos.core manifest is simply overloaded.

### cnos.agent

This is the missing package. Today the system already distinguishes commands and orchestrators as first-class content classes, and current cnos.core is shipping daily, weekly, save, and daily-review together with the broader agent skill bundle. Those belong naturally together as "default agent behavior," not as part of the core substrate.

**cnos.agent should own:**

- the non-bootstrap agent/* skills
- daily
- weekly
- save
- daily-review

### cnos.cdd

The CDD bundle is already a real system with:

- canonical docs
- its own lifecycle
- templates
- traceability
- sub-skills for design/review/release/post-release

So it deserves its own package rather than living inside core.

### cnos.eng

Leave this alone. It is already the cleanest package boundary.

### cnos.hub

Current ops/* is really about:

- inbox
- peer
- adhoc-thread
- star-sync

That is closer to hub / peer / communication operations than generic "ops." Named `cnos.hub` for best semantics (alternatively `cnos.ops` for lower churn).

### cnos.net.http

The package system and runtime-extensions docs already treat extensions as their own surface. The HTTP provider should not stay inside cnos.core if "core must be lean" is the rule.

### cnos.pm

Right now it is not coherent enough to keep as-is. Either:

- retire it and move its remaining useful skills elsewhere, or
- redefine it later around a genuinely distinct PM workflow bundle

Today it is mostly overlap with CDD.

## Migration Order

Lowest-churn path to the target:

1. **Extract cnos.cdd** from cnos.core
2. **Extract cnos.hub** from cnos.core
3. **Extract cnos.net.http** from cnos.core
4. **Retire or redefine cnos.pm**
5. **Extract cnos.agent** from what remains of cnos.core

Why this order:

- it removes the clearest non-core bundles first
- it lets you see what "lean core" really means before deciding the exact agent/core split
- it keeps cnos.eng stable

## Summary

The simplest final picture:

```
cnos.kernel    # embedded bootstrap kernel
cnos.core      # doctrine + mindsets + templates + minimal bootstrap cognition
cnos.agent     # default reusable agent behavior
cnos.cdd       # coherence-driven development bundle
cnos.eng       # engineering bundle
cnos.hub       # inbox/peer/adhoc/star-sync bundle
cnos.net.http  # HTTP capability provider
```

And:

```
cnos.pm        # retired or redefined later
```

That is the most coherent package structure for what cnos has today and where it is clearly heading.
