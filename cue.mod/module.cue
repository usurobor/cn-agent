// cue.mod/module.cue — CUE module declaration for the cnos repository.
//
// Established by cycle/388 (Phase 2.5 of cnos#366) to enable cross-package
// import semantics for the three-package schema split:
//   schemas/cdd/  — generic coherence-cell kernel
//   schemas/cds/  — software protocol overlay (imports cdd)
//   schemas/cdr/  — research protocol overlay (imports cdd)
//
// The module path `cnos.dev/cnos` is a canonical-style placeholder identifier
// (not a real domain). CUE only requires the path to be stable and unique
// within the build; resolution is filesystem-relative under `cue.mod/`. The
// existing `schemas/skill.cue` (used by scripts/ci/validate-skill-frontmatter.sh)
// is unaffected — it is invoked directly by path, not by import.

module: "cnos.dev/cnos"
language: version: "v0.10.0"
