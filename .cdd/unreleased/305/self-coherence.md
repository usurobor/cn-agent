# Self-Coherence — Issue #305

## Gap

**Issue:** #305 — tools(cdd): extend cdd-verify for current cycle-scoped artifacts

**Version/Mode:** substantial change / release-scoped triadic cycle

**Selected gap:** The existing `cn cdd-verify` command validates the pre-#283 release-scoped CDD artifact layout (`.cdd/releases/{X.Y.Z}/{role}/CLOSE-OUT.md`) but does not validate the current cycle-scoped layout introduced in #283 (`.cdd/releases/{X.Y.Z}/{N}/self-coherence.md`, `beta-review.md`, etc.). No tool currently validates the current cycle-scoped artifact set, required sections within each artifact, or orphaned unreleased directories. This creates a gap where γ can declare a cycle closed with missing close-outs, CHANGELOG rows can be malformed, and CDD's process discipline relies on agent self-assertion rather than structural verification.

**What exists:** `cnos.cdd` ships `cn cdd-verify` with `--triadic` mode that validates the legacy aggregate close-out layout.

**What is expected:** CDD validation should check the current cycle-scoped artifact layout per CDD.md §5.3a Artifact Location Matrix.

**Where they diverge:** The verifier checks `.cdd/releases/{X.Y.Z}/{role}/CLOSE-OUT.md` but current cycles produce `.cdd/releases/{X.Y.Z}/{N}/{artifact-name}.md`.