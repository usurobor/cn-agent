# alpha close-out — cycle #609

**Issue:** [cnos#609](https://github.com/usurobor/cnos/issues/609) — "cds-install Sub 2: generalize `cn install wake` identity (agent/PAT/bot) — successor to #549." **Verdict at merge:** converge (β R0), zero fix iterations — β's independent re-verification found no findings on the first pass.

---

## What α built

Three new flags on the `cn-install-wake` shell renderer — `--workflow-pat-secret <NAME>`, `--bot-name <name>`, `--bot-id <id>` — each supporting both `--flag value` and `--flag=value` forms, matching the existing `--agent`/`--out` convention. `workflow_pat_secret` resolves to the flag value, else `SIGMA_WORKFLOW_PAT` when `--agent` is `sigma`, else `die` before any render/write phase runs — `out_path` is never touched until the script's final `cp`, so a non-sigma agent omitting the flag leaves no partial `--out` file (Mock C2). `bot_name`/`bot_id` resolve to their flag values, else fall back unchanged to the existing `agent_bot_name()`/`agent_bot_id()` tables, which already `die` on unknown non-sigma agents. All three literal `SIGMA_WORKFLOW_PAT` occurrences in the render body now interpolate `${workflow_pat_secret}`.

The "Mechanical recovery scanner" and "Mechanical checkpoint + PR finalizer" steps were branched on `$agent`: `sigma` keeps the exact `(cd src/go && go build -o /tmp/cn-scan|cn-finalize ./cmd/cn)` lines byte-for-byte (Mock E4); any other agent gets a new "Install cn (tenant-portable — no src/go build)" step that acquires the prebuilt `cn` binary via the repo's existing `install.sh` into a fixed path (`/tmp/cnos-cn-609`), and both mechanical steps then invoke that binary directly instead of assuming `src/go/` exists (Mock E2 — the cnos#606 `tsc` tenant dogfooding finding this mock was written to close).

One mechanical, non-behavioral fix landed in `cnos.cds/orchestrators/cds-dispatch/SKILL.md`: the single literal `cds-dispatch-sigma` occurrence in the "Disallowed surfaces" bullet (describing this wake's own concurrency group) became `cds-dispatch-{agent}`, matching the same value's `{agent}`-substituted form used three other times in the same file — the one real self-reference gap that would otherwise have kept AC2's zero-leak invariant from closing regardless of how well the renderer's own bindings were parameterized.

Five new CI steps landed in `install-wake-golden.yml`: AC1 positive (acme render carries the caller's PAT secret / bot name / bot id / concurrency group), AC1 negative (Mock C2 fail-early), AC2's two-layer zero-sigma-leak oracle (a synthetic dispatch-shape SKILL.md fixture, reusing the existing AC5/AC2-negative fixture-authoring pattern, plus a scoped grep against the real cds-dispatch acme render for the renderer-controlled leak tokens only), and E2/E4 (tenant acme render has no src/go build step and acquires `cn` via `install.sh`; the sigma-default golden still self-builds byte-identically).

All five of this sub's owned mock-parity rows (C2, C4, C5, E2, E4 — C1/C3/C6 and E1/E3 belong to `cn repo install`, issue #610, per γ's scaffold) closed `match`, `missed: 0`.

## The one round, and why it converged clean

Unlike cycle #608 (one β-caught fix-round on a negative test that stubbed around the actual boundary it claimed to cover), this cycle converged at R0 with no findings. The structural reason: every negative-case claim in this cycle (Mock C2's fail-early behavior, the AC2 leak audits) was verified by actually running the built renderer shell script end-to-end and inspecting its exit code / stderr / file-existence, not by constructing an intermediate data structure and asserting against it in isolation — there was no equivalent of #608's `Invocation{HubPath: ""}` shortcut available here, because a POSIX shell script has no unit-test seam to stub around; the only way to exercise `die()` firing before the render phase is to actually invoke the script and check the real exit code and the real absence of the `--out` file. β's independent re-run of every AC/invariant (its own `/tmp/beta-609-*` paths, not mine) landed on identical PASS/FAIL results to my own self-coherence.md, including one extra check I hadn't run myself — a repo-wide grep for any other `- name:` step containing an embedded `": "` colon-collision, confirming the two step-name fixes I made were the only two needed, not a partial fix. That additional check is exactly the kind of thing β's role exists to catch even when a cycle otherwise converges cleanly.

## Debt disclosed (self-coherence.md §Debt, all non-blocking, none re-raised by β as findings)

1. The fixed tenant-acquisition path (`/tmp/cnos-cn-609`) is cycle-numbered — γ's scaffold explicitly pre-authorized "pick one consistent name," so this is cosmetic, not a correctness gap. A future cycle could rename it to something acquisition-purpose-named with no behavior change.
2. The tenant "Install cn" step is rendered once (inside the scanner block's non-sigma branch) and reused by the finalizer step via the shared job/runner filesystem, relying on both blocks staying co-gated on the same `role == "dispatch"` condition. If a future cycle decouples their gating, the finalizer would need its own acquisition step — flagged for whoever touches that gating next.
3. `--git-user-name`/`--git-user-email` were not implemented, per γ's scaffold Friction note 1 (explicitly out of the pinned Mock C/E design surface) — not re-litigated this cycle, as instructed.
4. `cn repo install`'s consumption of these new flags is issue #610's job; this cycle only lands the renderer-side flags, with no end-to-end wiring yet.

β's review noted item 1 explicitly and declined to re-raise it as a finding, agreeing it was already correctly disclosed and pre-authorized by scope.

## Artifact list

`src/packages/cnos.core/commands/install-wake/cn-install-wake` (new flags, identity resolution, tenant-portable acquisition branch), `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` (one-line `{agent}` parameterization fix), `.github/workflows/install-wake-golden.yml` (5 new CI steps), `.cdd/unreleased/609/self-coherence.md` (this cycle's primary artifact), `.cdd/unreleased/609/beta-review.md` (β's independent R0 verification, verdict: converge).

No file under `src/go/**`, `cn repo install`, or `docs/development/design/cn-repo-install-MOCKS.md` was touched (scope guardrails held); `cnos-cds-dispatch.golden.yml`, `cnos-agent-admin.golden.yml`, and the live `.github/workflows/cnos-cds-dispatch.yml` are confirmed byte-identical to pre-cycle (both α's and β's independent sha256 checks agree: `3dee3d1574...` on both sides of the live/golden pair).

## Final AC/invariant status

AC1 (positive + Mock C2 negative), AC2 (both leak-audit layers), AC3/Mock C5 (sigma byte-identical), AC4 (gates green), E2, E4: all PASS, all independently re-confirmed by β against the built renderer and the CI YAML — not merely read off self-coherence.md's claims. Zero fix-round iterations were needed; β's verdict is `converge` at R0, and this close-out reflects that as the terminal state for this cycle.
