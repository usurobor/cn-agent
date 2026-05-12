# Self-Coherence Report — Issue #328

## Gap

**Problem statement:** CDD artifacts cannot be missing from active or released cycle directories without CI reporting the gap.

**Current incoherence:** The existing `cdd-verify` command validates individual release versions but not the repository-wide CDD artifact ledger. CI runs build/test checks but lacks systematic validation of CDD lifecycle artifacts across both `.cdd/unreleased/{N}/` and `.cdd/releases/{X.Y.Z}/{N}/` directories.

**Specific gaps:**
1. **Stale validator contract** - Current `cdd-verify` checks legacy paths like `.cdd/releases/{version}/alpha/CLOSE-OUT.md` instead of current issue-scoped paths `.cdd/releases/{version}/{issue}/alpha-closeout.md`
2. **Missing CI integration** - No continuous validation of CDD artifact completeness across the repository
3. **Incomplete coverage** - Pre-tag release gate only checks `.cdd/unreleased/` but not historical `.cdd/releases/` ledger
4. **No small-change detection** - Cannot distinguish between substantial triadic cycles (requiring full artifact set) and small-change cycles (with explicit collapse rules)
5. **Silent historical gaps** - Missing historical artifacts go undetected without systematic checking

**Impact:** CDD can produce "witness theater" where artifacts appear to exist but validation is incomplete, leading to:
- Incomplete cycle directories surviving after release
- Missing close-outs not detected until manual inspection
- RELEASE.md present but potentially incorrect for the version
- CDD lifecycle rules remaining prose without mechanical enforcement

**Expected behavior:** CI fails on unexcepted missing required artifacts, with clear diagnostics identifying the missing path and expected contract.