# Self-Coherence — Issue #305

## Gap

**Issue:** #305 — tools(cdd): extend cdd-verify for current cycle-scoped artifacts

**Version/Mode:** substantial change / release-scoped triadic cycle

**Selected gap:** The existing `cn cdd-verify` command validates the pre-#283 release-scoped CDD artifact layout (`.cdd/releases/{X.Y.Z}/{role}/CLOSE-OUT.md`) but does not validate the current cycle-scoped layout introduced in #283 (`.cdd/releases/{X.Y.Z}/{N}/self-coherence.md`, `beta-review.md`, etc.). No tool currently validates the current cycle-scoped artifact set, required sections within each artifact, or orphaned unreleased directories. This creates a gap where γ can declare a cycle closed with missing close-outs, CHANGELOG rows can be malformed, and CDD's process discipline relies on agent self-assertion rather than structural verification.

**What exists:** `cnos.cdd` ships `cn cdd-verify` with `--triadic` mode that validates the legacy aggregate close-out layout.

**What is expected:** CDD validation should check the current cycle-scoped artifact layout per CDD.md §5.3a Artifact Location Matrix.

**Where they diverge:** The verifier checks `.cdd/releases/{X.Y.Z}/{role}/CLOSE-OUT.md` but current cycles produce `.cdd/releases/{X.Y.Z}/{N}/{artifact-name}.md`.

## Skills

**Tier 1 (Lifecycle):**
- `CDD.md` — canonical lifecycle and role contract
- `alpha/SKILL.md` — α role surface  
- `design/SKILL.md` — not required (single command extension, no design needed)
- `plan/SKILL.md` — not required (straightforward script extension)

**Tier 2 (Engineering):**
- `eng/*` bundles per cnos.eng README

**Tier 3 (Issue-specific):**
- `eng/tool` — shell script standards, command behavior, fail-fast
- `eng/test` — positive/negative proof, fixtures  
- `cnos.core/skills/write` — documentation standards

**Active constraint:** All loaded skills constrain generation. Shell script must follow `eng/tool` standards (set -euo pipefail, prereq checks, NO_COLOR support). Tests must prove invariants per `eng/test` (positive/negative fixtures, oracle-driven). Documentation follows `write` standards (one governing question per file, front-load the point).

## ACs

**AC1: `--cycle` mode added to `cn cdd-verify`**
- **Invariant:** `cn cdd-verify --version X.Y.Z --cycle N` validates `.cdd/releases/{X.Y.Z}/{N}/` against current artifact layout  
- **Oracle:** `cn cdd-verify --version 3.61.0 --cycle 283` exits 0 for real completed cycle
- **Evidence:** Command accepts `--cycle` parameter, validates cycle directory path, produces pass/fail diagnostics
- **Status:** ✅ Complete

**AC2: Artifact set completeness is checked**  
- **Invariant:** Every cycle directory must contain 5 hard-gate artifacts
- **Oracle:** Command lists present/missing artifacts with diagnostic naming missing file
- **Evidence:** Validates presence of `self-coherence.md`, `beta-review.md`, `alpha-closeout.md`, `beta-closeout.md`, `gamma-closeout.md`
- **Status:** ✅ Complete

**AC3: Required sections within artifacts are checked**
- **Invariant:** Each hard-gate artifact has required sections per role SKILL.md
- **Oracle:** Command names artifact and missing section 
- **Evidence:** Section validation logic per CDD role specifications
- **Status:** ✅ Complete

**AC4: CHANGELOG row validation**
- **Invariant:** Released cycle has CHANGELOG row for its version
- **Oracle:** Command checks CHANGELOG.md for row matching version
- **Evidence:** CHANGELOG parsing and version lookup
- **Status:** ✅ Complete

**AC5: Orphan detection**
- **Invariant:** `.cdd/unreleased/{N}/` should be empty on main after all cycles released
- **Oracle:** Command warns on any directories in `.cdd/unreleased/`
- **Evidence:** Directory enumeration with warning output
- **Status:** ✅ Complete

**AC6: Negative fixtures exist**
- **Invariant:** Fixture set includes invalid example per hard-gate check
- **Oracle:** Running invalid fixtures produces non-zero exit and names violation
- **Evidence:** Extended test suite with `--cycle` mode test cases
- **Status:** ✅ Complete

**AC7: Documentation updated**
- **Invariant:** Command help and comments document `--cycle` mode vs `--triadic` legacy
- **Evidence:** Help output and header comments updated
- **Status:** ✅ Complete

## Self-check

**Did α's work push ambiguity onto β?** No. The implementation provides clear validation with specific error messages, comprehensive test coverage, and follows established shell script patterns. β can verify functionality through the test suite.

**Is every claim backed by evidence in the diff?**
- ✅ **AC1:** `--cycle` parameter added, validates cycle directory, see `cn-cdd-verify` lines 31, 73, 84-86
- ✅ **AC2:** Hard-gate artifact validation, see `cn-cdd-verify` lines 283-302  
- ✅ **AC3:** Section validation implemented, see `validate_artifact_sections` function
- ✅ **AC4:** CHANGELOG validation already exists in original code
- ✅ **AC5:** Orphan detection implemented, see lines 320-331
- ✅ **AC6:** Extended test suite with 4 new test cases, see `test-cn-cdd-verify.sh` lines 117-206
- ✅ **AC7:** Documentation updated in usage function and header comments

**Role boundaries clean?** Yes. α implemented the feature completely without requiring β to discover missing work. No authoring work deferred to review phase.

## Debt

**Known debt:**

1. **Section validation is basic** — The `validate_artifact_sections` function uses simple grep patterns rather than structured markdown parsing. It catches obvious missing sections but may not validate section content quality or completeness. This is acceptable for v0 as stated in issue scope ("Out of scope: Validating prose quality or content correctness").

2. **No CI integration yet** — The verifier extension is implemented but not yet integrated into CI workflows. This is explicitly deferred per issue scope ("Deferred: CI integration - document path, implement later").

3. **Cross-artifact consistency not validated** — The verifier doesn't check consistency between artifacts (e.g., γ triage references every β finding). This is explicitly out of scope per issue requirements.

4. **No runtime enforcement** — The tool validates after-the-fact but doesn't prevent invalid state during authoring. This is intentional as the tool is designed for post-completion verification.

**If a loaded skill would have prevented remaining debt:** No applicable skill gaps identified. The debt items are either acceptable limitations for v0 or explicitly deferred scope.