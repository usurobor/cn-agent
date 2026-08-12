# Case-2 rented-Claude smoke receipt

One bounded, disposable episode, recorded because the shared corpus cannot
carry it: the corpus rents `fake`, and a CI job that rents real cognition would
be the provider service this project deliberately does not build.

The raw closure stdout is committed at
[`evidence/cds-case2-claude-closure.json`](evidence/cds-case2-claude-closure.json).
`scripts/cell-schema-check.sh` re-derives **a named subset** from it on every
run — the closure schema, the generic seat envelope, `execution_mode`,
`status`, the UTF-8 diff byte count, the diff SHA-256, and the touched-file
list. A one-byte edit **inside the diff**, or deleting the file, fails the
gate.

Stated exactly, because the earlier version of this file overclaimed
(Pi #59 D1): the gate does **not** check the episode id, the recorded
provider/model/base metadata, the scope-lift digest, or JSON whitespace.
Those are reproduced below as transcription from the artifact, and a reader
who cares should read the artifact. This is an evidence fixture, not a
provider harness: nothing in the corpus invokes a RENTED cognition provider.
(The corpus does run the deterministic `fake`; what it never does is call the
Claude CLI.)

**Why it exists (Pi #57 D1).** The exact-argv oracle proves the recipe; it
cannot prove the runtime. Before `--permission-mode acceptEdits` was sealed,
the seat's tools were *available* but not *approved*, so earlier real episodes
were resting on ambient host permission authority rather than on the cell's
declaration. The claim was open exactly where the fix landed.

**The recorded contract is older than the fixture it names.** Since this run,
`code-cell-spec.json` gained a required `contract.task` — a typed CDS issue —
and then a `contract.subject`, the repository and commit both stations receive
frozen; this episode ran with neither. Its recorded alpha declaration still
carries the `workspace` block that `cds.patch` no longer accepts, so the gate
vets that declaration against the generic seat envelope rather than
`#CDSPatchAlphaResolved`: the artifact records a shape that has since been
deleted, and editing it to fit would break the digest it is evidence for.

State the reduction plainly rather than let "envelope" carry it: `#Seat` is
`{fill!: string & !="", ...}`, so what that row still checks is that the
declaration names a fill — nearly nothing. The canonical-shape claim this
fixture used to carry is gone, not narrowed. What it still proves is
unchanged and is in the gate rows below: the mode, the status, the measured
diff's byte count and digest, and the touched-file list.

Its goal line, *"Carry out the change
described by the issue in the repository at base_sha"*, is exactly the sentence
that motivated the change: it refers to an issue the cell was never given.
Re-running the invocation below against today's fixture would therefore produce
a DIFFERENT contract, a different digest, and a differently prompted seat; the
committed closure is a record of what ran then, not a reproduction recipe for
now. Nothing in the gate compares the two, and this note is here so nobody
reads the agreement of paths as agreement of contracts.

## Invocation

Runtime was an immutable, clean commit — not "a head plus changes":

| | |
|---|---|
| runtime commit | `ca1f241b36b0835b8be3922af2e6a34c8a8270ef` |
| runtime tree | `9d85c712beafd632c8cdaeecf040cb033917bf91` |
| working tree at run *(observed)* | clean (`git status --porcelain` empty) |
| `claude --version` *(observed)* | `2.1.226 (Claude Code)` |
| fixture repo base | `1d79f7552649357165ce9addf3fbe7c57f3b62b0` |

Built with `go -C src/go build -o /home/user/cnos/cn ./cmd/cn` at that commit.
The exact invocation SHAPE below — every flag and value is literal; only the
disposable scratch prefix is abbreviated as `<scratch>`, since it was a
per-run temporary directory that no longer exists:

```sh
CN=/home/user/cnos/cn
SPEC=/home/user/cnos/schemas/cds/fixtures/code-cell-spec.json
REPO=<scratch>/coderepo          # throwaway git repo, base 1d79f755
HUB=<scratch>/hub                # .cn/vendor/packages/ from DefaultPackages
EV=<scratch>/closure.json

cd "$HUB"
timeout 900 "$CN" cell run \
  --contract "$SPEC" \
  --param language=cnos.eng:eng/go \
  --param provider=claude-cli \
  --param model=claude-opus-5 \
  --param base_sha=1d79f7552649357165ce9addf3fbe7c57f3b62b0 \
  --param repo="$REPO" > "$EV" 2><scratch>/stderr.txt; rc=$?
# rc=1, stderr empty; "$EV" is the committed artifact
```

`model` is a **requested selector**, not an observed model identity. Nothing
in the runtime asks the provider what actually served the request.

## Result

Rows marked **gate** are recomputed by `scripts/cell-schema-check.sh`. Rows
marked *transcribed* are present in the artifact but not asserted by the gate.
Rows marked *observed* were seen at run time and are **not** in the artifact at
all.

| | |
|---|---|
| episode *(transcribed)* | `ep-7a83e6c07a2749068aab291152113946` |
| exit *(observed)* | `1` |
| `execution_mode` **(gate)** | **`cognitive`** |
| status **(gate)** | `needs_repair` |
| measured diff **(gate)** | **2479 UTF-8 bytes** |
| diff sha256 **(gate)** | `3826a7e883a9fb78769d1ef99ca54a16bad631aea244620412e2d5be58261766` |
| touched **(gate)** | `CONTRIBUTING.md` (new), `README.md` |
| recorded cognition *(transcribed)* | `{"provider":"claude-cli","model":"claude-opus-5"}` |
| recorded base *(transcribed)* | `1d79f7552649357165ce9addf3fbe7c57f3b62b0` |
| `#EpisodeClosure` **(gate)** | vets |
| `#CDSPatchAlphaResolved` **(gate)** | vets |

Recompute by hand:

```sh
ev=docs/architecture/evidence/cds-case2-claude-closure.json
cue vet schemas/cdd/episode-closure.cue "$ev" -d '#EpisodeClosure'
python3 -c 'import json,sys;json.dump(json.load(open(sys.argv[1]))["receipt"]["record"]["resolved_spec"]["alpha"],open("/tmp/a.json","w"))' "$ev"
cue vet schemas/cdd/spec.cue /tmp/a.json -d '#Seat'   # not #CDSPatchAlphaResolved; see above
python3 -c 'import hashlib,json,sys;d=json.load(open(sys.argv[1]))["receipt"]["record"]["matter"]["data"];print(len(d.encode()),hashlib.sha256(d.encode()).hexdigest())' "$ev"
```

**`VerifyClosure`.** `cellrun` self-verifies the emitted closure against the
contract and metadata *this invocation* built, and exits 2 with empty stdout if
that fails; encoding happens only after. The preserved artifact is a complete
closure and the observed exit was 1, so verification succeeded. The premise
that was missing before — the preserved stdout — is the committed file.

## What this establishes

1. **The declared authority is sufficient.** The seat edited files under the
   declared `--tools` surface and `--permission-mode acceptEdits`. (That
   surface has since been corrected to include `MultiEdit` and `Bash`, so a
   seat can run its own tests; this episode predates the change and did not
   need them.)
2. **The runtime measures rather than believes.** The diff was computed from
   the disposable worktree.
3. **Case-2 honesty holds under real cognition.** A genuine change still
   closes `needs_repair`: the mechanical β will not pass what it cannot judge,
   and the matter is preserved for Case 3's independent reviewer.
4. **Mode follows the provider** — `cognitive`, honestly irreproducible.

## What it does not establish

- Nothing about `codex-cli`, which remains held.
- Nothing about which model served the request; only which selector was asked
  for.
- **Not "nothing inherited from the host."** The adapter seals a declared
  baseline tool and permission recipe and suppresses user/project defaults via
  `--safe-mode`. Authentication is ambient by design, and vendor-managed
  substrate policy can apply above the baseline; the adapter neither detects
  nor overrides it.
- No OS confinement. The honest authority is the offered tool surface, the
  declared permission mode, and the measured worktree.
- Source-repo cleanliness after the run was observed at run time
  (`git status --porcelain` empty, `HEAD` unchanged) but is **not** recomputable
  from the committed artifact, so it is recorded here as an observation only
  and nothing in the gate asserts it.
