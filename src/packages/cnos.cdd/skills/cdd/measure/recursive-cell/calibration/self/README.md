# Candidate CM self-measurement

Status: same-author diagnostic only; no independent, consistency, admissibility,
or held-out standing.

- TSC source: `26aab5023f03dc7d0abf82e5fdba20134fc6adad`
- Engine: `coh 0.12.0 (016c511)`
- Standard instruction SHA-256:
  `ca26d7a1a4dc6bd73e0afff558ed0342d42daeee193d4169de2c51b762759391`
- Frozen 15-file prompt SHA-256:
  `50f8900b0a50480cfbf240f5d816ba364a91b1324e32ab76a91472e015ad03ac`
- Mechanical repetitions: 3, exact; alpha=1.000, beta=1.000,
  gamma=0.7675, C-sigma=0.9155726158, gamma bottleneck.
- Semantic samples: 1; alpha=0.94, beta=0.84, gamma=0.60,
  C-sigma=0.7795658333, gamma bottleneck.

The frozen target contained the CM authority, instruction, README, #662
calibration record, registry, six manifests, fail-closed target preflight,
cnos base/methodology schemas, and cnos skill validator. It intentionally
excluded this directory's response/report: a measurement output cannot become
part of the bytes it claims to have measured after the fact.

`semantic-response.json` is the validated TSC v3.2.4 witness.
`report.json` is the report produced by the engine's hybrid external-response
route.
