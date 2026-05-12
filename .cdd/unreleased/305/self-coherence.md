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
- **Status:** ⏳ Pending implementation

**AC2: Artifact set completeness is checked**  
- **Invariant:** Every cycle directory must contain 5 hard-gate artifacts
- **Oracle:** Command lists present/missing artifacts with diagnostic naming missing file
- **Evidence:** Validates presence of `self-coherence.md`, `beta-review.md`, `alpha-closeout.md`, `beta-closeout.md`, `gamma-closeout.md`
- **Status:** ⏳ Pending implementation

**AC3: Required sections within artifacts are checked**
- **Invariant:** Each hard-gate artifact has required sections per role SKILL.md
- **Oracle:** Command names artifact and missing section 
- **Evidence:** Section validation logic per CDD role specifications
- **Status:** ⏳ Pending implementation

**AC4: CHANGELOG row validation**
- **Invariant:** Released cycle has CHANGELOG row for its version
- **Oracle:** Command checks CHANGELOG.md for row matching version
- **Evidence:** CHANGELOG parsing and version lookup
- **Status:** ⏳ Pending implementation

**AC5: Orphan detection**
- **Invariant:** `.cdd/unreleased/{N}/` should be empty on main after all cycles released
- **Oracle:** Command warns on any directories in `.cdd/unreleased/`
- **Evidence:** Directory enumeration with warning output
- **Status:** ⏳ Pending implementation

**AC6: Negative fixtures exist**
- **Invariant:** Fixture set includes invalid example per hard-gate check
- **Oracle:** Running invalid fixtures produces non-zero exit and names violation
- **Evidence:** Extended test suite with `--cycle` mode test cases
- **Status:** ⏳ Pending implementation

**AC7: Documentation updated**
- **Invariant:** Command help and comments document `--cycle` mode vs `--triadic` legacy
- **Evidence:** Help output and header comments updated
- **Status:** ⏳ Pending implementation