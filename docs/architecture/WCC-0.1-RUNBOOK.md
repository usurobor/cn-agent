# WCC 0.1 — running the coding cell, and what comes back

How to run one coding episode, how to get the patch out of the closure, and how
to check the patch **outside** the closure that reports it. Companion to
[`WCC-ADVERSARIAL-SUITE.md`](WCC-ADVERSARIAL-SUITE.md) (what the cell resists)
and [`WCC-CODE-QUALITY-AUDIT.md`](WCC-CODE-QUALITY-AUDIT.md) (what its own code
still owes).

Every command below was run as written. The evidence index at the end names the
six committed closures, and `src/go/internal/cellkernel/evidence_test.go` checks
each of them on every `go test ./...`.

## 1. What one run needs

| | |
|---|---|
| the binary | `go -C src/go build -o /tmp/cn ./cmd/cn` |
| a hub | a directory with `.cn/vendor/packages/` — the cell resolves its methodology skills from there, so the run's cwd is the hub |
| a cell spec | `schemas/cds/fixtures/code-cell-spec.json` — reusable, holds the parameter holes |
| a run input | one document per run: issue + design + subject (§2) |
| a provider CLI on PATH | `claude` for `--param provider=claude-cli` |
| a credential | **ambient operator configuration.** It never appears in the cell spec, the run input, or the closure, and the runtime passes only the selected provider's credential to its child process |

The spec is reusable and the input is not, deliberately: folding the issue into
the spec would make every run author a new cell.

## 2. The run input

`cnos.cds.run-input.v0` carries three payloads the runtime keeps opaque —
`cellinput` owns the envelope and never learns what an issue is.

```json
{
  "kind": "cnos.cds.run-input.v0",
  "issue":  { "kind": "cnos.cds.issue.v0",  "...": "see schemas/cds/fixtures/issue/valid-issue.json" },
  "design": { "kind": "cnos.cds.design.v0", "...": "see schemas/cds/fixtures/design/" },
  "subject": {
    "kind": "git.snapshot/0.1",
    "repo": "/home/user/cnos",
    "base_sha": "bfb4d977f34b055bafbc679944fc917130b9b224"
  }
}
```

The subject is the **only** place the repository and the base are stated. There
used to be a second: the `cds.patch` seat declared its own `workspace`, and two
declarations that could disagree produced a closure that self-verified while
recording a repository the episode never acted on. The declaration was deleted;
`--param repo` and `--param base_sha` no longer exist.

The issue must satisfy `cdsissue.Admit` — every field present and non-blank,
every acceptance criterion carrying a verification route, no id in the reserved
`check:` namespace. Admission happens **before** any seat is constructed: a
refused input costs nothing and produces an admission receipt on stdout with
exit 4.

## 3. The invocation

```sh
CN=/tmp/cn
HUB=/tmp/wcc-live/hub          # any directory carrying .cn/vendor/packages/
IN=/tmp/wcc-live/run-input.json
OUT=/tmp/wcc-live/closure.json

cd "$HUB"
timeout 2400 "$CN" cell run \
  --contract /home/user/cnos/schemas/cds/fixtures/code-cell-spec.json \
  --input "$IN" \
  --param language=cnos.eng:eng/go \
  --param provider=claude-cli \
  --param model=claude-opus-5 \
  > "$OUT" 2>/tmp/wcc-live/stderr.txt; echo "exit=$?"
```

`model` is a **requested selector**, not an observed identity — nothing asks the
provider what actually served the request. No argv, environment, executable
path, timeout or safety override can be supplied through the cell spec or the
run input; a timeout a cell could set is a cell that can decline to be bounded.

Exit codes: `0` accepted · `1` needs_repair / degraded / rejected · `2` usage,
malfunction, or **failed self-verification** (stdout empty) · `3` simulated,
never ordinary accepted · `4` run input refused at admission.

`cellrun` re-runs `VerifyClosure` against the contract and metadata *this
invocation* built — never the closure's own — and encodes only after it passes.
So a closure on stdout is already a verified closure.

## 4. Getting the patch out

The candidate is `matter.data`: a unified diff the runtime **measured** from a
disposable worktree against the pinned base. It is not something a seat
reported.

```sh
EV=docs/architecture/evidence/wcc-0.1/self-q1.json
jq -r '.receipt.record.matter.data' "$EV" > /tmp/patch.diff
BASE=$(jq -r '.receipt.record.contract.subject.base_sha' "$EV")

git clone -q /home/user/cnos /tmp/w
git -C /tmp/w checkout -q "$BASE"
git -C /tmp/w apply --check /tmp/patch.diff && echo "applies clean at $BASE"
```

Run as written on `self-q1.json`: `applies clean at bfb4d977…`, 4 files, +229
−90.

Then the part that matters — **check it outside the closure**. The closure
reports the gates it ran; a receipt that is trusted because it says so is not
evidence. Apply the patch to a fresh clone at the pinned base and run the gates
by hand:

```sh
git -C /tmp/w apply /tmp/patch.diff
( cd /tmp/w && gofmt -l $(git diff --name-only -- '*.go') )   # empty = formatted
( cd /tmp/w/src/go && go vet ./... && go build ./... && go test ./... -count=1 )
```

`gofmt` is scoped to the files the patch touched, base-relative. Scoping it to
the whole tree would report drift the episode did not cause — the same defect
`cellcheck` had when it scoped `format` off `git status` instead of the base,
which let a seat that committed its own work pass with an unformatted file.

`-count=1` always: a cached PASS is a record of a previous run.

Custody is the other half. The caller's checkout must be untouched by an
episode that acted on a disposable worktree:

```sh
git -C /home/user/cnos rev-parse HEAD    # unchanged from before the run
git -C /home/user/cnos status --porcelain # empty
```

## 5. Reading the closure

```sh
jq '{status, decision, verdict: .verdict.pass,
     mode: .receipt.record.execution_mode,
     bytes: (.receipt.record.matter.data | utf8bytelength)}' "$EV"

jq -r '.receipt.record.review.assessment[]
       | "\(.unit)\t\(.disposition)\t\(.reason)"' "$EV"
```

The assessment covers the issue's own acceptance criteria plus the runtime's
`check:` units. A mechanical checker's `fail` **forces** a `finding`: a
cognitive answer that contradicts a forced disposition is a fault, not a
judgement call — cognition cannot launder a mechanical result.

## 6. The committed evidence

Six closures under `evidence/wcc-0.1/`, all `execution_mode: cognitive`, all
`accepted`. The two green runs establish the path; the four self-runs are
episodes whose patches landed as commits on this branch, so each one's base is
its predecessor.

| file | issue | subject base | matter (bytes) | units |
|---|---|---|---|---|
| `green-calc.json` | `calc-mean` | throwaway calc repo | 1,579 | 5 |
| `self-betafill.json` | `cellfill-beta-needs-subject` | cnos @ `2a968d30` | 15,615 | 7 |
| `self-q1.json` | `wcc-quality-q1-seat-helper` | cnos @ `bfb4d977` | 20,298 | 5 |
| `self-q2.json` | `wcc-quality-q2-one-bound` | cnos @ `651259e2` | 21,862 | 5 |
| `self-q3a.json` | `wcc-quality-q3-cellcog` | cnos @ `6de3b449` | 33,050 | 6 |
| `self-q3b.json` | `wcc-quality-q3-four-fills` | cnos @ `58a4a1bd` | 48,411 | 5 |

`TestCommittedEpisodesVerifyAndBindEveryValueTheyClaim` puts each through
`VerifyClosure` and then mutates twenty-five named values, requiring every one
to move the scope-lift digest.

Not committed here: the nine other live episodes. Seven are the adversarial
suite, whose outcomes are in
[`evidence/wcc-adversarial-suite.jsonl`](evidence/wcc-adversarial-suite.jsonl);
two are the timed-out Q1+Q2-combined and first Q3 attempts, which produced no
closure and are the measurement that raised the provider bound from 10 to 30
minutes.

## 7. What a run does not establish

- **No OS confinement.** The declared tool surface is a CAPABILITY DECLARATION,
  not a boundary. The honest authority is the offered tools, the declared
  permission mode, and the fact that the diff was measured from a worktree the
  caller's checkout does not share.
- **Nothing about which model served the request** — only which selector was
  asked for.
- **Nothing about the issue.** The cell is bounded by issue quality, not by
  capability: it resists a hostile issue (the adversarial suite) and faithfully
  satisfies a weak one. An acceptance criterion whose verification route is
  "read it in the view" gets exactly that — a reading, with no executable
  witness behind it.
- **Not reproducible.** `execution_mode: cognitive` says so. Re-running the same
  input produces a different patch and a different digest, which is why the
  closures are committed rather than regenerated.
