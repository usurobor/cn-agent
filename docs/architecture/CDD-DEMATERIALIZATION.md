# CDD Dematerialization — `main` holds "what is," history holds "how we got here"

**Status:** Design (L7) · **Tracking issue:** [#682](https://github.com/usurobor/cnos/issues/682) · **Depends on:** #681 (first principle in doctrine — merged)
**Related:** #683 (open-items ledger), #684 (channel plane), #680 (repo-self-coherence methodology)

---

## Summary

`main`'s working tree must contain only artifacts that answer **"what is"** — the current code, specs, config, and a present-tense index. The CDD cell receipts under `.cdd/releases/` and closed `.cdd/unreleased/` answer **"how we got here."** They are already permanent in Git history (the commits that created them); materializing them *again* as 10.4 MB / 1,022 files in every checkout is noise that violates the principle.

**Target:** the closed-cell receipts leave the working tree and live only in `main`'s **ancestry** (reachable commits), retrievable on demand. `main` keeps a single present-tense finding-aid — `.cdd/INDEX.jsonl` — that relates to history exactly as `CHANGELOG.md` does: a curated current view that *points into* the past without *being* it.

This is **event sourcing applied to coherence receipts**: the sealed commits are the append-only log; the index is a rebuildable read-model projection.

---

## Governing principle

> `main` contains only what is needed to answer **"what is."** How the project got here is warranted by history — reachable from `main`, retrievable, auditable — but **not materialized** in the working tree.

Encoded in `DOCUMENTATION-SYSTEM.md §5` (the first principle, landed in #681) and `KERNEL.md §2.1` (no silent drops — history is never dropped, only dematerialized).

---

## 1. Current state (concrete)

`git ls-tree -r main -- .cdd/` today: **1,022 files, 10.4 MB**, materialized in every clone and checkout.

```
.cdd/
├── CADENCE  CDD-VERSION  DISPATCH  MCAs  OPERATORS       ← config: "what is"  (stays)
├── exceptions.yml  legacy-exceptions.yml                 ← config: "what is"  (stays)
├── skills/                                               ← vendored bundle: "what is"  (stays)
├── proposals/  iterations/                               ← queue/backlog: mostly "what is"
├── releases/          596 files · 5.2 MB                 ← closed receipts: "how we got here"  (MOVE)
│   └── {X.Y.Z}/{alpha,beta,gamma}/{issue}.md
├── unreleased/        375 files · 5.0 MB                 ← cells, mostly closed: "how we got here"  (MOVE)
│   └── {N}/{self-coherence,beta-review,gamma-closeout,…}.md
└── waves/             10 files                           ← wave matter, closed: "how we got here"  (MOVE)
```

**The split within `.cdd/`.** Not all of it is history. ~971 files / ~10.2 MB — `releases/` + closed `unreleased/` + `waves/` — are closed-cell receipts. The rest (config + vendored skills) is current-state and stays. The design is about the receipts, not the config.

**What pins the receipts in HEAD today:** tooling reads them from the *working tree*, not from history — `cdd-verify/ledger.go`, `cn-cdd-status`, and `scripts/release.sh` read `.cdd/releases/`. This runtime coupling is the actual reason the directory cannot simply be deleted.

## 2. What is wrong

| # | Problem | Evidence |
|---|---|---|
| P1 | **Principle violation** | 10.2 MB of "how we got here" sits in a tree that should hold only "what is." |
| P2 | **Attention cost** | Every newcomer, agent context load, tree walk, and code search wades through 971 closed-receipt files to find the ~current ones. The problem is *attention*, not storage — the bytes are already Git blobs. |
| P3 | **Unbounded working-set growth** | The receipt set grows with every cell, forever, in the *working tree* — not just in history where unbounded growth is normal. |
| P4 | **Runtime coupling** | Tools read `.cdd/releases/` from the tree, so the directory can't leave HEAD until the coupling is cut — the coupling *is* the blocker. |

## 3. Target state (concrete)

`git ls-tree -r main -- .cdd/` after migration:

```
.cdd/
├── CADENCE  CDD-VERSION  DISPATCH  MCAs  OPERATORS  exceptions.yml  …   ← config (unchanged)
├── skills/                                                             ← vendored bundle (unchanged)
├── INDEX/                              ← NEW · the finding-aid (like CHANGELOG.md)
│   ├── 2026.jsonl                      ← one line per closed cell, sharded by year
│   └── …
└── CURRENT.json                        ← NEW · current-state view (current release, open cells)
```

**Gone from the working tree:** `.cdd/releases/**`, closed `.cdd/unreleased/{N}/**`, `.cdd/waves/**` — every closed-cell payload. They remain **reachable in `main`'s ancestry** (the commit that sealed each cell is an ancestor of `main`) and are retrieved with `cn cdd materialize {N}` / `git show <seal>:.cdd/…`.

**`.cdd/INDEX/{year}.jsonl`** — the present-tense catalog. One append-only line per closed cell:

```jsonc
{"cell":"671","seal":"25bca3ad…","outcome":"accepted","kind":"planning",
 "contract":"cn.wave.v1","closed":"2026-07-25","artifacts":[".cdd/unreleased/671/self-coherence.md", …]}
```

It **is "what is"** — the current index of the archive, exactly like a library catalog or `CHANGELOG.md`: present-tense, points into history, rebuildable from it. Sharded by year to avoid append-contention on a single hot file.

**`.cdd/CURRENT.json`** — the current-state view tools read instead of walking `.cdd/releases/`: current release pointer, open commitments, live-cell references. Replaces the `.cdd/releases/`-as-database read pattern.

**In-flight cells are not here at all.** A cell being worked lives on its own branch `cell/{N}` — in *that* working tree, never in `main`'s. `main` only ever sees a cell as (a) a line in `INDEX` after it closes, or (b) not at all while open.

## 4. Lifecycle — `S ≺ D ≺ P` (commit-level)

A cell's receipts persist **on the commit that created them**, not in the HEAD tree. Three ordered events:

```
cell/671:  ●──α──●──β──●  S  (seal: γ writes MANIFEST.json + closeout; commit contains .cdd/unreleased/671/**)
                          ╲
main:      …──● base ──────● D  (δ merges cell/671 → main; S is now an ANCESTOR of main;
                            │      history-preserving merge — no squash that orphans S)
main:                       ● P  (a DESCENDANT of D: `git rm -r .cdd/unreleased/671/`
                                   + append INDEX/2026.jsonl line + update CURRENT.json)

result:  HEAD tree has NO .cdd/unreleased/671/**   ·   `git show S:.cdd/unreleased/671/self-coherence.md` works
```

- **S — Seal.** γ writes the closeout + `MANIFEST.json` (cell id, contract digest, artifact paths + blob digests, role commits, outcome, residuals). Commit `S` contains the complete payload. The manifest cannot name `S` itself.
- **D — Decide.** `V` validates; δ records the boundary decision by merging (history-preserving). `S` becomes reachable from `main`.
- **P — Project/Prune.** A descendant of `D` removes the payload from the tree and updates `INDEX`/`CURRENT`. The payload leaves HEAD; it stays in `S`.

**The one hard invariant (seal-before-prune):** a payload MUST NOT leave HEAD until its sealing commit is **reachable from `main`.** This forbids the data-loss path *seal → delete-on-branch → squash to empty tree → no `main`-reachable payload*. Reachable-in-ancestry is the real requirement; materialized-in-tree never was.

**Publication by outcome** (the failure-persistence rule): a *rejected* cell's product matter does **not** merge, but the episode still happened — publish a receipt-only seal into ancestry (`base → receipt-only seal → δ-rejection/index commit → prune`). The receipt lands in ancestry regardless of accept/reject; only the *product change* is gated on acceptance.

## 5. The index is a projection, not a source of truth

`.cdd/INDEX/**` and `CURRENT.json` are **derived** and **rebuildable**:

```
cn cdd index --rebuild    # walk main's history, reproduce INDEX/** + CURRENT.json exactly
cn cdd index --check      # verify: every seal SHA resolves & is reachable; every manifest
                          #   blob digest matches; one seal per closed cell; one index line per seal
```

Because it rebuilds exactly from history, it is never a second source of truth — if the file and history disagree, history wins and `--rebuild` corrects it. **Ever-growing is fine**: linear in cell count, one small line each, and compactible (shard/snapshot), never deleted. This is the standard log-structured shape.

**Precedent** (this is a well-worn pattern, not a novel bet):
- **Git itself** — working tree = "what is"; commit DAG = "how we got here." We apply Git's own model recursively to `.cdd`.
- **Event sourcing / CQRS** — append-only log + rebuildable read-model projections. Sealed commits = log; `INDEX`/`CURRENT` = projections.
- **Blockchain UTXO set** — append-only chain + a derived current-state set rebuilt by replay. Ancestry ↔ chain; `CURRENT.json` ↔ UTXO set.
- **Datomic / log-structured stores (Kafka, LSM)** — immutable history + derived indexes, managed by compaction.

## 6. Impact graph (what must change)

| Consumer | Today | After |
|---|---|---|
| `cdd-verify/ledger.go` | reads `.cdd/releases/` tree | reads `CURRENT.json` + `cn cdd` reader over history |
| `cn-cdd-status` | walks `.cdd/` tree | reads `CURRENT.json` / `INDEX` |
| `scripts/release.sh` | reads `.cdd/releases/{X.Y.Z}/` | reads `CURRENT.json` release pointer |
| `cn cdd` (new) | — | `list · show {N} · materialize {N} · index --rebuild/--check · verify` |
| CI | — | (a) no closed payload in HEAD; (b) every INDEX seal reachable; (c) `index --check` passes |
| Merge policy | ad hoc | history-preserving for CHAIN cells; CONTENT cells may squash iff the squash commit carries the complete sealed payload before prune |

## 7. Migration plan (phased; each phase independently shippable and reversible)

- **Phase 0 — Doctrine (done).** First principle in `DOCUMENTATION-SYSTEM.md §5` + `KERNEL.md §2.1` (#681, merged). No file moves.
- **Phase 1 — Reader + schema (no removal).** Build `cn cdd` reader, `INDEX`/`CURRENT`/`MANIFEST` schemas, and `index --rebuild/--check`. Prove `cn cdd materialize {N}` reconstructs a cell **from the seal SHA alone**, in a fresh clone, on a **pruned-HEAD fixture**, with filesystem fallback disabled and tested negatively. Removes nothing.
- **Phase 2 — Break the runtime coupling.** Repoint `ledger.go` / `cn-cdd-status` / `release.sh` from the `.cdd/releases/` tree to `CURRENT.json` + the reader. Still removes nothing — but now nothing *reads* the payload from the tree.
- **Phase 3 — One-time historical prune.** A single mechanical migration commit: `git rm -r` all closed-cell payloads (`releases/`, closed `unreleased/`, `waves/`), build the full `INDEX/**` from history, write `CURRENT.json`. Large diff, purely mechanical, reviewed once. Payloads remain in ancestry (verified by `index --check` reachability).
- **Phase 4 — Lifecycle enforcement.** New cells follow `S ≺ D ≺ P`. CI enforces: no closed payload in HEAD, every INDEX seal reachable, `index --rebuild` reproduces exactly. The steady state.

## 8. Acceptance criteria

- **AC1** — `cn cdd materialize {N}` reconstructs the complete typed cell from the **seal SHA only**, in a fresh full clone, no network, on a **pruned-HEAD fixture** (`git ls-tree HEAD …/{N}` → absent). Filesystem fallback prohibited and tested negatively.
- **AC2** — `cn cdd index --rebuild` reproduces `INDEX/**` + `CURRENT.json` **byte-for-byte** from history; `--check` fails if any seal SHA is unreachable, any manifest blob digest mismatches, or the seal↔index bijection breaks.
- **AC3** — a **rejected** cell's product diff is absent from `main`, yet its receipt seal is `main`-reachable and reconstructable (failure-persistence).
- **AC4** — CHAIN custody survives a real merge (original role commits + exact parents reachable; history-preserving merge); CONTENT survives a squash **iff** the squash commit carries the complete sealed payload before prune. Prune-gate = custody-reachable.
- **AC5** — a shallow clone lacking the seal returns `HISTORY_INCOMPLETE {required_sha, remediation}` — never "cell not found"; `INDEX`/`CURRENT` stay visible.
- **AC6** — post-migration, `git ls-tree -r main -- .cdd/` contains **no** closed-cell payload; only config + `skills/` + `INDEX/**` + `CURRENT.json`.

## 9. Alternatives considered

| Option | Verdict |
|---|---|
| **Orphan ref for `.cdd`** (like the channel plane, #684) | **Rejected.** A cell is *causal lineage* — a decision points back to the cell that warranted it. An orphan ref severs the ancestry walk. Correct for channels (parallel streams); wrong for receipts (main's own past). |
| **Payloads in commit-message bodies** (Move B) | **Rejected as the general form.** Poor fit for review tables / finding lists / cross-file links; loses Markdown + file validation; awkward amendments. Fine only for small boundary events. |
| **Keep a rich projection in HEAD** (early #682) | **Superseded.** The richer in-HEAD projection was a middle step; the destination is index-only — `.cdd` in HEAD is a `CHANGELOG`-scale finding-aid, nothing more. |
| **File-native blobs + commit-native index** (Move A) | **Chosen.** Payloads stay file-native in sealed reachable commits; the index is the derived current view. |

## 10. Risks / negative leverage

| Risk | Mitigation |
|---|---|
| Squash-merge orphans a seal → **data loss** | `S ≺ D ≺ P` invariant + CI reachability check (AC2); CHAIN cells forbid post-review rebase. |
| Shallow clone can't materialize | `HISTORY_INCOMPLETE {required_sha, remediation}`, never "not found" (AC5). |
| Index drift from history | Rebuildable + `--check` in CI (AC2). |
| Large one-time prune diff (Phase 3) | Mechanical, single reviewed commit; content preserved in ancestry, verified reachable. |
| Append-contention on the index | Shard by year (`INDEX/{YYYY}.jsonl`) and/or rebuild rather than hand-append. |

## References

- **Issues:** #682 (this), #681 (first principle — merged), #683 (open-items ledger — a `CURRENT.json` consumer), #684 (channel plane), #680 (repo-self-coherence).
- **Doctrine:** `DOCUMENTATION-SYSTEM.md §5`, `KERNEL.md §2.1`, `docs/architecture/CELL-RUNTIME.md`.
- **Prior art:** Git object model; event sourcing / CQRS; Datomic; blockchain UTXO set; Kafka / LSM log compaction.
