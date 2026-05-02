## Gap

**Issue:** #327 — bug(cdd): release gate must validate required cycle artifacts before tag
**Version/mode:** MCA — bug fix to `scripts/release.sh` (pre-tag validation gate)
**Priority:** P1

`release.sh` moves `.cdd/unreleased/{N}/` to `.cdd/releases/{X.Y.Z}/{N}/` but does not
validate whether the directory contains the required lifecycle artifacts before tagging.
3.73.0 shipped with `.cdd/releases/3.73.0/325/` missing `alpha-closeout.md` and
`gamma-closeout.md` — both required by the CDD lifecycle that 3.73.0 itself introduced.

The fix: add `scripts/validate-release-gate.sh` that blocks the release if required
artifacts are missing, and call it from `release.sh` before stamp/move/tag.

Design artifact: not required — single-concern tooling addition (2 new scripts +
3-line change to release.sh); no impact graph beyond the 3 scripts it touches.

Plan artifact: not required — implementation sequencing is trivial (test → implement →
wire in).

---

## Skills

**Tier 1a (CDD authority):**
- `src/packages/cnos.cdd/skills/cdd/CDD.md` — canonical lifecycle and role contract
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — α role surface

**Tier 2 (general engineering):**
- Tooling bundle applies: all outputs are shell scripts in `scripts/`

**Tier 3 (issue-specific):**
- `src/packages/cnos.core/skills/design/SKILL.md` — design principles applied to
  validation contract boundary (one reason to change, truthful interface)
- `src/packages/cnos.core/skills/write/SKILL.md` — script header comments and
  self-coherence prose
- `src/packages/cnos.eng/skills/eng/test/SKILL.md` — test from invariants, not just
  examples; negative space mandatory; oracle explicit
- `src/packages/cnos.eng/skills/eng/tool/SKILL.md` — shell script standards:
  set -euo pipefail, fail-fast, idempotent, zero runtime deps, machine-readable output
