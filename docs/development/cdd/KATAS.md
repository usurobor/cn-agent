# Agent Kata

Kata are executable proofs. Scripts and packages are the source of truth;
this catalog is a secondary index.

## Architecture

cnos separates the kata **framework** (the mechanism that finds and runs
katas) from **kata content** (the actual proofs). Both ship as packages.

- **`cnos.kata`** — the framework: `cn kata run`, `cn kata list`,
  `cn kata judge`. Discovers katas across any installed package
  (convention: `<package>/katas/<id>/`). Also carries its own Tier 2
  runtime katas as dogfood of the framework.
- **`cnos.cdd.kata`** — content-only: CDD method katas (M0–M5). No
  commands. The framework in `cnos.kata` picks them up on discovery.
- **Future domain packages** (eng, ops, …) may ship their own katas by
  adding a `katas/` directory; no framework changes needed.

Class is declared inside each `kata.md` as `**Class:** <runtime|method>`.
The framework dispatches on that class:

- **runtime** — binary pass/fail. `run.sh` is executed; exit code is the verdict.
- **method** — scored per rubric. Produces a run bundle (`metadata.json`,
  `artifacts/`) for the judge.

## Tier model

| Tier | Proves | Home |
|------|--------|------|
| 1 — bare binary | `cn` binary works end-to-end before any package is installed | `scripts/kata/` |
| 2 — runtime/package | post-install behavior (dispatch, roundtrip, doctor-broken, self-describe) | `cnos.kata/katas/R*/` |
| 3 — method/CDD | CDD adds value over ad hoc execution | `cnos.cdd.kata/katas/M*/` |

Release-bootstrap and network-dependent compatibility checks belong in
`scripts/smoke/`, not here — they are production-facing, not kata.

## Tier 1 — bare binary (`scripts/kata/`)

CI gate: `.github/workflows/ci.yml` job `kata-tier1`. Failure = red build.

| # | Name | Proves | Script |
|---|------|--------|--------|
| 01 | binary  | `cn` runs                            | `scripts/kata/01-binary.sh` |
| 02 | init    | `cn init` creates a hub              | `scripts/kata/02-init.sh`   |
| 03 | status  | `cn status` reads hub state          | `scripts/kata/03-status.sh` |
| 04 | doctor  | `cn doctor` validates clean hub      | `scripts/kata/04-doctor.sh` |
| 05 | build   | `cn build` produces dist/            | `scripts/kata/05-build.sh`  |
| 06 | install | `cn deps restore` installs packages  | `scripts/kata/06-install.sh` |

```bash
scripts/kata/run-all.sh
```

## Tier 2 — runtime/package (`cnos.kata/katas/`)

Proves post-install behavior after at least one package is installed.
CI gate: `.github/workflows/ci.yml` job `kata-tier2` (`needs: kata-tier1`).
Failure = red build.

| ID | Name | Proves |
|----|------|--------|
| R1 | Command         | package command discovery + dispatch (`cn help`, `cn daily`) |
| R2 | Round-trip      | author → `cn build` → `cn deps restore` → dispatch in an isolated workdir |
| R3 | Doctor (broken) | `cn doctor` catches `chmod -x` on an installed command entrypoint |
| R4 | Self-describe   | `cn status` surfaces installed package name, version, and commands |

```bash
cn kata run --class runtime   # run R1..R4, stop on first failure
cn kata run R2-roundtrip      # run one
```

`cnos.kata` must be installed (`cn deps restore`) before `cn kata run`
can be dispatched — the dispatch itself is weak proof that the package
loop works; R1-R4 make the proof explicit.

## Tier 3 — method/CDD (`cnos.cdd.kata/katas/`)

Proves CDD adds value over ad hoc execution. Content-only package; the
commands come from `cnos.kata`.

| ID | Name | Purpose |
|---|------|---------|
| M0 | Gap | frame the right incoherence before design |
| M1 | Design | artifact completeness + traceability |
| M2 | Review | evidence-bound + architecture-aware |
| M3 | Post-release | closure quality + measurement |
| M4 | Full cycle | end-to-end CDD loop vs ad hoc on the same change |
| M5 | Issue authoring | schema validation + README alignment + witness-theater check |

```bash
cn kata list                           # list all discovered katas
cn kata list --class method            # method only
cn kata run M1-design --mode cdd       # scored, produces a run bundle
cn kata judge <run-dir>                # honest stub until judge wiring lands
```

## Adding katas from a new package

1. Create `src/packages/<pkg>/katas/<id>/kata.md` with `**Class:** runtime` or `method`.
2. For runtime, add an executable `run.sh` alongside `kata.md`.
3. For method, add `rubric.json`, `baseline.prompt.md`, `cdd.prompt.md`; `run.sh` is optional (corpus-only stubs are supported).
4. `cn build && cn deps restore` — your kata now appears in `cn kata list`.

No framework change is required. Discovery is convention-driven.

## Post-release usage

Reference which katas were run, on what environment, pass/fail.
