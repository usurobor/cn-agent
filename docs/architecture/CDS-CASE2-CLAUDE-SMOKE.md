# Case-2 rented-Claude smoke receipt

One bounded, disposable episode run against an exact head, recorded because
the corpus cannot carry it: the shared gate rents `fake`, and a CI job that
rents real cognition would be a provider service this project deliberately
does not build. This file is the receipt, not a harness — nothing runs it.

**Why it exists (Pi #57 D1).** `ClaudeArgv`'s exact-argv oracle proves the
recipe; it cannot prove the runtime. Before `--permission-mode acceptEdits`
was sealed, the seat's tools were *available* but not *approved*, so earlier
real episodes were resting on ambient host permission authority rather than
on the cell's declaration. That made the runtime claim unproven at exactly
the boundary the fix changed, so the fix needed a real episode behind it.

## Invocation

| | |
|---|---|
| head | `b867ca02` + this round's changes |
| CLI | `claude` 2.1.226 (Claude Code) |
| requested model selector | `claude-opus-5` |
| provider | `claude-cli` |
| skills | `cnos.eng:eng/code`, `cnos.eng:eng/test`, `cnos.eng:eng/go`, `cnos.eng:eng/write-functional` — loaded from a hermetic hub vendored from `repoinstall.DefaultPackages` |
| spec | `schemas/cds/fixtures/code-cell-spec.json` (the committed cell, unmodified) |
| workspace | throwaway git repo, `base_sha` pinned at construction to `a2b7500a198f1a7c11321fc3f17aad5f4541289e` |

"Requested model selector" is exact wording. The CLI remaps some selectors and
says so on stderr; nothing in the runtime observes what actually served the
request, so the closure records what was asked for and claims nothing more.

## Result

| | |
|---|---|
| episode | `ep-76b472f4a8660726eb1847e6060d85f1` |
| exit | `1` — `needs_repair` |
| `execution_mode` | **`cognitive`** — a provider held the seat, so the episode is not reproducible |
| measured diff | **4173 bytes**, nonempty, touching `CONTRIBUTING.md` (new, 87 lines) and `README.md` |
| β review | `pass: false` — "mechanical seat cannot judge the goal; matter preserved for independent review" |
| `VerifyClosure` | passed |
| `#EpisodeClosure` | vets |
| `#CDSPatchAlphaResolved` | vets |
| source repo after the run | unchanged at `a2b7500`, clean tree |

`VerifyClosure` is proven structurally rather than by a separate check:
`cellrun` self-verifies the emitted closure against the contract and metadata
*this invocation* built, and exits 2 with empty stdout if that fails. A
complete closure on stdout with exit 1 is therefore a closure that verified.

## What this establishes

1. **The declared authority is sufficient.** The seat edited files with only
   `--tools Read,Write,Edit,Glob,Grep` and `--permission-mode acceptEdits` —
   no Bash, no bypass, nothing inherited from the host.
2. **The runtime measures rather than believes.** The diff was computed from
   the disposable worktree, and the source repo was untouched afterwards.
3. **Case-2 honesty holds under real cognition.** A genuine 4 KB change still
   closes `needs_repair`, because the mechanical β will not pass what it
   cannot judge. The matter is preserved for Case 3's independent reviewer.
4. **Mode follows the provider.** `cognitive`, not `mechanical` — the episode
   is honestly marked irreproducible.

## What it does not establish

- Nothing about `codex-cli`, which remains held.
- Nothing about which model actually served the request; only which selector
  was asked for.
- No claim of OS confinement. The authority is the offered tool surface, the
  declared permission mode, and the runtime-measured worktree.
- `--safe-mode` closes the user/project customization layer only; higher
  authority managed substrate policy can still apply above it.
