# cnos recursive-cell CM

Candidate coherence methodology produced from the external CC analysis of
cnos#662. `SKILL.md` is the authority, `INSTRUCTION.md` is the semantic witness
contract, and `calibration/662/README.md` records the first non-independent
calibration.

Publication issue: cnos#669. Parent runtime wave: cnos#627.

The current TSC release does not consume arbitrary CM declarations directly.
To apply this candidate today:

1. check out receipt head `a0d39293...`, overlay this CM package, and run
   `calibration/662/verify-target.sh <checkout>` (it fails on any other HEAD
   or matter bytes);
2. emit frozen prompts using TSC's `coh --emit-prompt`;
3. obtain three validated semantic responses using `INSTRUCTION.md`;
4. ingest each response with `coh --mode hybrid --llm-response`;
5. run mechanical determinism and semantic spread;
6. apply the hard-invariant gate from `SKILL.md` before disposition.

Do not use #662 as a held-out validation target; it designed this CM.

The bundled #662 registry points to its canonical package paths. Its target
corpus is the historical R7 receipt revision, not whatever version of #662 may
eventually land on main. The preflight makes that distinction fail-closed:
publication of the candidate CM does not pretend the unratified architecture
already ships or silently measure later bytes as its calibration source.
