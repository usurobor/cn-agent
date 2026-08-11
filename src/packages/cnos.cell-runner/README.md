# cnos cell-runner (walking-skeleton spike)

**Exploratory spike — proof-of-life, NOT canonical doctrine or runtime.** This
is a runnable skeleton of the three-cell agentic loop (CC → PC → WC → CC). It
exists to learn from *running* the loop, ahead of the authored R12 cell-runtime
doctrine (#672) and its implementation wave (#627 S2–S8). When that doctrine is
authorized, this spike is reconciled to it or discarded. It invents no policy;
where the design was underspecified it takes the minimal mechanical option and
records the choice in `DESIGN.md` §Build notes.

## What it proves

One turn of the loop, mechanically, with no human in the machinery: measured
incoherence goes **1 → 0** across a `CC.measure → PC.plan → WC.execute →
CC.re-measure` turn, and a chain of `cue vet`-green receipts is produced. The
three cells are distinguished ONLY by output telos (CC → judgment, PC →
relation_graph, WC → artifact); they run the **same** kernel FSM.

- **Kernel FSM** is a declarative transition table (`kernel/transitions.json`)
  evaluated by a generic engine (`kernel/engine.go`, `kernel/drive.go`) that
  never hardcodes a state name in a switch/if-chain — the
  `cnos.issues/issues-fsm` idiom.
- **Actors** (`actor/stub/`) are the "assisted" seam. v0 ships deterministic
  stubs so the loop runs with zero external calls. Agent-backed actors
  implement the same `kernel.Actor` interface next (v1) with no kernel change.
- **Receipts** (`schema/receipt.cue`) are **computed** at γ from transition
  evidence, never authored. The verdict × action → transmissibility table is
  enforced structurally in CUE (fail closed).

## Run it

```sh
export PATH="$PATH:$(go env GOPATH)/bin"   # so `cue` is on PATH
cd src/packages/cnos.cell-runner

# happy path: exits 0, prints the 1 → 0 transcript, writes objects under $WS
go run ./cmd/cell-runner --workspace "$(mktemp -d)" --task testdata/toy

# fail-injection: drives the kernel to `held`, writes a stop receipt, exits non-zero
go run ./cmd/cell-runner --workspace "$(mktemp -d)" --task testdata/toy --fail-injection=beta-reject
go run ./cmd/cell-runner --workspace "$(mktemp -d)" --task testdata/toy --fail-injection=v-fail
```

Emitted under `<workspace>/`: `cm/*.cm.json`, `wave/pc.wave.json`,
`receipts/*.receipt.json`, `judgment/*.judgment.json`, and the working copy of
the target under `target/`.

## Validate & test

```sh
cue vet ./schema/ <workspace>/cm/cc-1.cm.json         -d '#CM'
cue vet ./schema/ <workspace>/wave/pc.wave.json       -d '#Wave'
cue vet ./schema/ <workspace>/receipts/wc.receipt.json -d '#Receipt'

gofmt -l .        # prints nothing
go vet ./...      # clean
go test ./...     # green (the cue-vet integration test self-skips if cue is absent)
```

## Layout

```
cmd/cell-runner/   CLI entry (flags, exit codes)
kernel/            transition table (JSON) + generic evaluator + Drive loop + Actor seam
actor/stub/        deterministic CC/PC/WC actors (the v0 "assisted" boundary)
cell/              kernel-under-telos + the CC→PC→WC→CC loop driver
model/             Go structs for the emitted objects (JSON field names match schema/)
schema/            CUE definitions (cm, wave, contract, judgment, receipt)
testdata/toy/      the toy target (markdown missing `## Coherence`)
```

## Explicitly deferred (not in v0)

Agent-backed actors; multi-WC / parallel waves; the wave FSM beyond one node;
real CM scoring beyond defect-count; the operator-authorization gate;
self-hosting on a real cnos change; tsc adoption. Each is a named next
increment (see `DESIGN.md`), not a silent gap.
