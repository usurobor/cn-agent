---
retroactive: true
cycle: 333
issue: "#330"
pr: "#333"
merge_sha: "6ffdf48a"
reconstructed_by: "#335"
reconstructed_date: "2026-05-09"
---

# CDD Iteration — Cycle #333

> **RETROACTIVE RECONSTRUCTION** — This artifact was not produced during the cycle.
> Reconstructed post-merge from git history by cycle #335 (issue #335).
> Original merge: `6ffdf48a` on 2026-05-08. Reconstruction date: 2026-05-09.
>
> Step 5.6b (requiring this file for cycles with `cdd-*-gap` findings) was introduced by cycle #331
> (merged same day). Cycle #333 had immediate access to this requirement and produced no such file.

This file lists the three `cdd-*-gap` findings addressed by cycle #333's patches. All dispositions are `patch-landed` since the patches are already on `cnos:main`.

---

### F1: alpha pre-review gate missing artifact enumeration, caller-path trace, and test assertion count rows

- **Source:** γ-triage / PRA §3 — upstream proposal from `usurobor/tsc`
- **Class:** `cdd-skill-gap`
- **Trigger:** γ process-gap check — α pre-review gate lacked structured checks for three specific verification categories
- **Description:** `cdd/alpha/SKILL.md` pre-review gate did not include: (row 11) artifact enumeration verifying that all artifacts named in the issue body exist on the branch; (row 12) caller-path trace confirming that the chain from issue → skill → implementation → artifact is unbroken; (row 13) test assertion count ensuring that test coverage claims in self-coherence.md match actual test counts in the diff.
- **Root cause:** The pre-review gate grew organically over multiple cycles, adding checks for known failure modes. Three recently-observed failure patterns (artifact claims without corresponding files, broken caller chains, over-counted test coverage) had not yet been codified as gate rows.
- **Disposition:** `patch-landed`
- **Patch:** `efed11e4`
- **Affects:** `cdd/alpha/SKILL.md` pre-review gate rows 11–13

---

### F2: operator skill missing --output-format stream-json requirement for all dispatches

- **Source:** γ-triage / PRA §3 — upstream proposal from `usurobor/tsc`
- **Class:** `cdd-protocol-gap`
- **Trigger:** γ process-gap check — operator dispatch templates lacked the stream-json output format requirement, causing silent JSON parse failures in harness integrations
- **Description:** `cdd/operator/SKILL.md` dispatch templates did not require `--output-format stream-json` when invoking `claude -p`. Without this flag, the output format defaults to text, which cannot be parsed as structured JSON by downstream harness consumers. This caused silent failures in automated dispatch pipelines.
- **Root cause:** The operator skill was authored before the harness integration pattern was established. The stream-json output format requirement emerged from operational experience but was not retroactively applied to the dispatch templates.
- **Disposition:** `patch-landed`
- **Patch:** `c48d5a95`
- **Affects:** `cdd/operator/SKILL.md` dispatch templates

---

### F3: CDD §1.6b re-dispatch prompt complexity note missing

- **Source:** γ-triage / PRA §3 — upstream proposal from `usurobor/tsc`
- **Class:** `cdd-protocol-gap`
- **Trigger:** γ process-gap check — re-dispatch prompts were growing in complexity without guidance on managing that complexity
- **Description:** `CDD.md` §1.6 (re-dispatch mechanism) lacked a note on re-dispatch prompt complexity. As cycles became more complex, re-dispatch prompts compensated for underspecified issues and skills by growing longer. This hid the root problem (weak issue quality or skill gap) behind prompt verbosity.
- **Root cause:** The re-dispatch mechanism was designed for simple scenarios. As it was applied to more complex situations, the natural tendency was to add context to the prompt rather than fix the underlying specification gap. §1.6b formalizes the correct response: fix the issue or skill, not the prompt.
- **Disposition:** `patch-landed`
- **Patch:** `30c02d1c`
- **Affects:** `CDD.md` §1.6b
