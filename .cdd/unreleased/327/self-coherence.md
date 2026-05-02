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
