# Historical candidate CM self-measurement — v0.2 repair

Status: same-author diagnostic only; no independent, consistency,
admissibility, or held-out standing. R3 supersedes this instrument and sample;
see `supersession-r3.json`.

`target.tsc` is the exact ordered 55-file source bundle and intentionally
excludes every measurement output in this directory. `receipt.json` binds the
manifest, framed bundle, registry, composite instruction, config, prompt,
three raw mechanical outputs, semantic response, hybrid report, commands,
engine revision/version, witness boundary, and UTC run times.

- TSC revision: `26aab5023f03dc7d0abf82e5fdba20134fc6adad`
- Engine: `coh 0.12.0 (26aab50)`
- Composite instruction SHA-256: `0b7e26a2fce364c2f89bbd04a72e8e1068cb7339de9e9b9e9d33e520eff75645`
- Ordered bundle SHA-256: `d6031d7f107f4bcce4f86d36e0c8255148ada28df80949912f0fe2ea57df4001`
- Frozen prompt SHA-256: `4ccef013f935b1733df82400c0b7eda9e05c04ba5f90ea8afb8ebd6df7ebeaa2`
- Mechanical N=3: exact numeric equality; alpha `0.94625`, beta
  `0.44117647058823534`, gamma `0.7432575757575758`, C-sigma
  `0.6769956211903662`, beta bottleneck.
- Semantic N=1: alpha `0.94`, beta `0.64`, gamma `0.44`, C-sigma
  `0.6420765882044609`, gamma bottleneck; zero standing.

Replay requires the historical cnos evidence revision; current R3 bytes use a
different instruction and must not be mixed with this response:

```bash
COH=/path/to/tsc/engine/ocaml/_build/default/bin/main.exe
HISTORICAL_REVISION=c8d680bf822b46b9f931d62b07d68f50ade1b07e
BASE=src/packages/cnos.cdd/skills/cdd/measure/recursive-cell
REGISTRY="$BASE/calibration/self/registry.tsc"
INSTRUCTION="$BASE/INSTRUCTION.md"
PROMPT=/tmp/recursive-cell-self.prompt.md

test "$(git rev-parse HEAD)" = "$HISTORICAL_REVISION"
"$BASE/instruction/assemble-instruction.sh" --check
"$COH" --mode llm --target recursive-cell-self --registry "$REGISTRY" \
  --instruction "$INSTRUCTION" --root . --emit-prompt "$PROMPT"
sha256sum "$PROMPT"
"$COH" --mode mechanical --target recursive-cell-self \
  --registry "$REGISTRY" --instruction "$INSTRUCTION" --root . \
  --output /tmp/recursive-cell-mechanical
"$COH" --mode hybrid --target recursive-cell-self --registry "$REGISTRY" \
  --instruction "$INSTRUCTION" --root . \
  --llm-response "$BASE/calibration/self/semantic-response.json" \
  --output /tmp/recursive-cell-hybrid
```

The semantic response was produced by a same-author OpenAI Codex hosted
activation; the host exposes the GPT-5 family but no immutable model revision.
That limitation is recorded rather than promoted into a provenance claim.
No R3 semantic sample was generated. The historical N=1 response cannot count
toward a later k=3 consistency claim because R3 changed the instrument.
