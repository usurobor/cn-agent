## Gap

**Issue:** #321 — core(activate): introduce Kernel/Persona/Operator activation triad and retire templates class

**Version/mode:** MCA — system-shaping change across package layout, pkg code, activate command, docs, and kata.

**Incoherence:** `cn activate` 3.70.0 emits a structurally wrong `## Identity` bucket that conflates kernel, persona, and operator into one unstructured list. On `cn-sigma` (the canonical reference hub), it reports `no identity files found`. The kernel ships under `templates/SOUL.md` as a "template" in both docs and code, entrenching a fork-the-kernel pattern. `pkg.ContentClasses` and `PACKAGE-SYSTEM.md` both list `templates` as a content class. Per-hub identity slots (`Name/Role/Operator`) live in the kernel file, contradicting kernel semantics. `scanPackages` reads only the vendor directory, not `deps.json`. No ordered first-read guidance, no latest-reflection pointer.

**Selected gap:** The triad (Kernel / Persona / Operator) is the correct layering. Each layer owns one canonical path. The kernel ships under `doctrine/` (not `templates/`). The `templates` content class is removed from code and docs. `cn activate` produces three layered sections, three kernel states, three deps states, an ordered `## Read first`, and a latest-reflection pointer.
