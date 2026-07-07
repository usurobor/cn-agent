# β close-out — cycle #609

**Issue:** [cnos#609](https://github.com/usurobor/cnos/issues/609) — cds-install Sub 2: generalize `cn install wake` identity (agent/PAT/bot), successor to #549. **Verdict:** converge at R0 — no iteration required.

---

## Review Summary

One round, on `origin/cycle/609`. Base `612a96d6` (per γ's scaffold, confirmed still current — no rebase needed). Reviewed head `01f08008` (α's implementation + self-coherence, 4 commits ahead of base: γ scaffold, two implementation commits, self-coherence). Every AC/invariant α claimed — AC1 positive, AC1 negative (Mock C2), AC2 two-layer (Mock C4), AC3/E4 (Mock C5 sigma byte-identical), E2 (tenant acquisition) — was independently re-executed from a fresh command line, not read off α's self-coherence transcript. No blocking findings. `beta-review.md` (R0) recorded `verdict: converge`, pushed as commit `93ddc21c`.

---

## Implementation Assessment

The renderer change holds to every guardrail γ's scaffold made binding on α:

- **Identity resolution is fail-early and structurally so, not just empirically.** I traced the actual line numbers: `workflow_pat_secret`/`bot_name`/`bot_id` resolution happens at line ~695-724, `tmp_out` (the render buffer) is not created until line 748, and the final `cp` to `out_path` doesn't happen until line 1360. A `die` in the resolution block cannot leave a partial `--out` file — this is a property of the script's control flow, not an artifact of the specific inputs I happened to test.
- **Sigma default path is untouched.** Re-rendering both `cds-dispatch` and `agent-admin` with no flags produced `(unchanged)` against both committed goldens, `git status --porcelain` was clean, and `sha256sum` of the live `.github/workflows/cnos-cds-dispatch.yml` matched the golden exactly (`3dee3d1574...`). This is the single highest-consequence regression this cycle could have introduced (a renderer used by every dispatch wake in the repo), and it held.
- **Scope discipline was total.** The diff touches exactly 4 files (`cn-install-wake`, one line of `cds-dispatch/SKILL.md`, `install-wake-golden.yml`, `self-coherence.md`) — nothing under `src/go/**`, no change to the pinned design doc, no change to either protected golden or the live `cnos-cds-dispatch.yml`, and no `--git-user-name`/`--git-user-email` flags (the issue-comment scope-creep item γ's Friction note 1 explicitly declined).
- **The one in-scope SKILL.md edit was exactly as narrow as promised.** `git diff` on `cnos.cds/orchestrators/cds-dispatch/SKILL.md` shows a single line changed (`cds-dispatch-sigma` → `cds-dispatch-{agent}`); the adjacent `agent-admin-sigma` reference and the `20260624.md` historical citation in the same paragraph are untouched.

---

## Technical Review

**Why I didn't just trust the AC2 split-oracle rationale.** γ's scaffold and α's self-coherence both assert that an *unscoped* `grep -i sigma` against the live `cds-dispatch` wake's rendered prompt body would never reach zero, regardless of how well the renderer's own bindings are parameterized — because two legitimate, out-of-scope strings (`agent-admin-sigma`, a different wake's own identity; and a dated historical log citation) live in package-owned prose this cycle isn't chartered to rewrite. Rather than accept that as asserted, I ran the *unscoped* grep myself against my own acme render and confirmed the only two hits were exactly those two named strings, verbatim. That converts the friction note from "a plausible-sounding excuse for a narrower oracle" into an empirically checked claim — the narrower, scoped oracle (four literal renderer-controlled tokens) is not covering for a real residual leak; it is correctly excluding two things that are not renderer leaks at all.

**Why I re-derived, not re-read, the fail-early guarantee.** α's self-coherence claims the resolution-before-write ordering closes Mock C2 "before any output is written." Rather than accepting that as a description of intent, I located the actual `tmp_out`/`out_path` write sites in the file (`grep -n 'tmp_out\|out_path'`) and confirmed their line numbers sit strictly after the identity-resolution block. This matters because a `die()` call that is merely "usually" reached first (e.g., guarded by a condition that could be bypassed by some flag combination) would not satisfy Mock C2's guarantee — I wanted the ordering to be unconditional in the script's structure, and it is.

**Regression-suite re-verification, not re-trust.** I independently re-ran, byte-for-byte, the shell commands from three pre-existing CI steps that predate this cycle (AC5 declaration-only refusal, AC2-negative malformed-SKILL.md rejection, AC7/AC8 renderer-source leak audit) against the changed renderer, rather than assuming "α says these still pass." All three passed with output matching their pre-existing expected shape (exit 3 + precise stderr; exit-rejection + exact error string; zero role-decision-string leaks in the renderer source).

**CI YAML scrutiny beyond the two cited fixes.** α's self-coherence names two step-name edits made to work around a YAML `: `-in-unquoted-scalar mis-parse. Rather than spot-checking just those two, I wrote a small script to check every `- name:` value in the file for an embedded `": "` — zero hits repo-wide in that file, confirming the fix pattern was applied consistently rather than just at the two cited call sites.

---

## Process Observations

- This cycle's γ scaffold did unusually heavy resolution work up front — the binding-oracle resolution section, two friction notes, and explicit scope guardrails all pre-empted the most likely sources of ambiguity (the issue's later `--git-user-name`/`--git-user-email` comment, the AC2-oracle-scope question, the "cn install wake" vs. "cn repo install" naming confusion). α's implementation and self-coherence report tracked that scaffold precisely — no scope drift, no re-litigation of already-made calls. The result was a single-round converge, which is the expected outcome when the scaffold does this much of the ambiguity-resolution work before implementation starts, not evidence that review was perfunctory (every AC was still independently re-executed, per the Technical Review section above).
- No CI run was polled for this cycle as part of my review (this was a local-verification round); the gates (`install-wake-golden`, Go, Package) should still be confirmed green in CI before merge, per AC4, though nothing in the diff touches Go code and the local re-run of every CI-mirrored shell step passed.

---

## Release Notes

A release note for #609 should say:

- `cn install-wake` (the shell renderer behind `cds-dispatch` / `agent-admin` wakes) no longer hard-binds dispatch-wake identity to `sigma`. Three new optional flags — `--workflow-pat-secret <NAME>`, `--bot-name <name>`, `--bot-id <id>` — let a caller render a fully non-sigma dispatch wake (workflow PAT secret, claude-code-action bot identity, and concurrency-group self-reference all caller-supplied).
- Omitting these flags for a non-sigma `--agent` fails early (before any `--out` file is written) rather than silently defaulting to sigma's identity or partially rendering.
- `--agent sigma` (or omitting `--agent`) is byte-identical to today's output — verified via golden diff and sha256 match on both `cnos-cds-dispatch.golden.yml`/`cnos-agent-admin.golden.yml` and the live `.github/workflows/cnos-cds-dispatch.yml`.
- Non-sigma dispatch renders now acquire the `cn` binary via `install.sh` rather than `cd src/go && go build ./cmd/cn` (which only exists inside the cnos repo itself) — unblocking tenant repos (cnos#606 dogfooding) that vendor cnos packages without a `src/go/` tree. Sigma's self-build path is unchanged.
- One mechanical, non-behavioral fix: `cnos.cds/orchestrators/cds-dispatch/SKILL.md`'s one remaining literal `cds-dispatch-sigma` self-reference is now `{agent}`-substituted, matching the same value's substitution elsewhere in the same file.
- Explicitly not in this cycle's scope (deferred, not forgotten): `cn repo install` wiring that consumes these new flags (issue #610); `--git-user-name`/`--git-user-email` flags (not in the pinned Mock C/E design surface — a candidate follow-up if an operator wants git-commit-identity parameterization too).
