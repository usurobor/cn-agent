# Skill Bundle — cnos origin declaration

cnos is the origin repository for the CDD skill bundle.

The canonical skill location is `src/packages/cnos.cdd/skills/cdd/`. A vendored copy under
`.cdd/skills/cdd/` is not required because this repo IS the source. The skills here are not
vendored from an upstream; they are authored and versioned in-tree.

Tenant repos that adopt CDD vendor their skill bundle from the cnos commit SHA pinned in their
own `.cdd/CDD-VERSION`. They copy `src/packages/cnos.cdd/skills/cdd/` into `.cdd/skills/cdd/`
at that SHA and refresh only when bumping `CDD-VERSION`. See `src/packages/cnos.cdd/skills/cdd/activation/SKILL.md §16`
for the full pull/sync procedure.

For cnos itself: the §24 verification check treats a non-empty `.cdd/skills/README.md` containing
this origin declaration as equivalent to a populated vendored bundle.
