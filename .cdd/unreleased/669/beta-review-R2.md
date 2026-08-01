# External β Re-review R2 — cnos#669 recursive-cell CM

**Activation boundary.** Fresh functional β activation over exact PR head `9cc9c3bbe97b11a94dbab6e47bdc95a2fb0a418c`. I authored none of the original matter or either repair round. I re-read the R1 finding and α R2 receipt, inspected the exact one-file diff and complete I5 job, and checked both cited workflow logs and final workflow states. This activation remains hosting-identity-limited: β and α share OpenAI Codex hosting and the `usurobor` GitHub identity, so this is functional role separation and does not create off-house or hosting-independent standing.

**Exact-SHA verdict: CONVERGE — `9cc9c3bbe97b11a94dbab6e47bdc95a2fb0a418c`.**

## Findings

The R1 blocker is closed without weakening the proof:

- `.github/workflows/build.yml` remains depth-1 for the primary checkout and adds one I5-only fetch of exact object `a0d39293a27cfe57b49dacff696345b1ee2cdb40` using `--no-tags --no-recurse-submodules --depth=1`.
- The next command verifies `${revision}^{commit}` before CUE setup or validation. A missing/unreachable object therefore still fails closed; the validator's own `cat-file`/historical-tree gate is unchanged.
- Push I5 run [29687234052](https://github.com/usurobor/cnos/actions/runs/29687234052) shows `fetch-depth: 1`, the exact pinned fetch to `FETCH_HEAD`, the complete positive/negative self-test, and `100 SKILL.md validated; no findings`.
- Pull-request I5 run [29687235912](https://github.com/usurobor/cnos/actions/runs/29687235912) shows the same pinned fetch, fixture result, and 100-skill result.
- Both containing workflows completed successfully, including Go, package/source drift, schema sync, links, I5/I6, dispatch guards, workflow parse, binary verification, and package verification.
- Independent execution of the unchanged historical target validator at this SHA resolves `cc662-system/l0/l1/l2/l3/l4` to `30/4/4/8/12/11` files.

The repair diff is exactly nine added lines in `.github/workflows/build.yml`. It does not change workflow permissions, action pins, validator semantics, methodology/schema bytes, package content, source/package path policy, or any unrelated job. `git diff --check` is clean.

## Original four-finding closure

1. **Frozen composite semantic instrument and manual State-A route:** closed. One digest-checked TSC-v3.2.4-plus-cnos instruction is supplied to the CLI's single instruction input; arbitrary declaration loading remains honestly unshipped.
2. **Registry/manifest/preflight validation and negative proof:** closed. Exact historical resolution is fail-closed and now runs successfully in both required shallow CI event paths.
3. **Replayable self-target, provenance, N=3 evidence, and non-circularity:** closed. Outputs are excluded from the target; retained hashes, prompt, engine/config provenance, and byte-identical mechanical runs remain bound; semantic N=1 remains zero-standing evidence only.
4. **Source/package path-base behavior:** closed. Source-checkout-only authority is explicit; the installed layout refuses activation through the supplied resolver and claims no deployed arbitrary-CM behavior.

## Deliberate non-blocking residuals

- TSC arbitrary-CM declaration loading and distinct authority-root/target-root support are unshipped.
- Installed-package activation remains deliberately refused.
- Semantic evidence remains one same-author sample and carries no independent, consistency, admissibility, held-out, or off-house standing.

These are correctly bounded State-A limitations. This review neither ratifies the CM nor promotes its standing, and #662 remains a design/calibration source rather than a held-out anchor.

**Recommended next CDD action:** γ may perform exact-SHA closeout for cnos#669, binding the two repair rounds, this R2 β result, and the green push/PR workflows. Any readiness, merge, ratification, standing promotion, or subsequent activation decision remains with its separately authorized downstream gate.
