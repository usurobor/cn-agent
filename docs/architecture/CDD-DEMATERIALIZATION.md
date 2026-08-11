# CDD Dematerialization — a typed, retention-governed seal-event protocol

**Status:** **Proposed** (design, L7) — revised R1 after external-β `REQUEST CHANGES`; **independent ratification pending**, prune **not** authorized, implementation issues **not** to be derived from this note yet.
**Tracking issue:** [#682](https://github.com/usurobor/cnos/issues/682) · **Depends on:** #681 (first principle — merged)
**Related:** #683 (open-items radar), #684 (channel plane), #642 (ε recurrence), #664 (actor identity), #627/#662 (cell-runtime doctrine)

> **The missing abstraction is not "Git history."** It is a **typed, discoverable, retention-governed CDD seal-event protocol** with separate authorities for *closed evidence*, *live cell state*, *open commitments*, and *independent channels*. Once that separation is explicit, dematerialization is a consequence, not a repository cleanup. This revision makes it explicit.

---

## R1 review-response map (external-β `REQUEST CHANGES`)

| Finding | Where addressed |
|---|---|
| **D1** CURRENT not rebuildable-from-history | §5 Authorities (INDEX = history-derived; CURRENT = declared union; radar = #683) |
| **D2** no seal-discovery protocol | §6 Seal-event protocol (trailers + locator ref + traversal/order/dedup) |
| **D3** stale `kind`/`cell/{N}` vocabulary | §7 Event identity (`cell_class`/`matter_domain`/`protocol_id`/`episode_id`; `cycle/{N}`) |
| **D4** `S≺D≺P` vs CONTENT squash | §8 Custody DAGs (CHAIN vs CONTENT; common invariant) |
| **D5** rejected-episode Git algorithm | §9 Receipt-only publication transaction |
| **D6** Phase-3 can delete active state | §11 Prune-eligibility predicate + dry-run manifest |
| **D7** no retention/trust governance | §12 Retention & trust contract |
| **D8** shallow-runner incompatible | §13 Acquisition contract |
| **D9** projection concurrency hot spots | §10 Projection ownership (per-episode immutable fragments) |
| **D10** CURRENT vs #683 authority overlap | §5 Authorities (explicit split) |
| **D11** missing dual-write/parity prototype | §14 Migration (Phase 1A–1C parity gate) |
| **D12** human retrieval UX | §13.3 Operator retrieval surface |
| **C1** lighter alternatives not compared | §4 Alternatives & leverage (4 structurally distinct moves) |
| **C2** "reversible" overstated | §14 (shippable / rollback-defined / recoverable-from-seals) |
| **C3** impact graph incomplete | §9-impact: §13.4 consumer graph (ε/waves/repair/PRA/installer/TSC/operator) |
| **C4** manifest under-typed | §7.2 `cn.cdd.seal.v1` schema |

---

## 0. Current situation (evidence-bound)

cnos has crossed from prompt-driven orchestration into a partially mechanical runtime, which changes what this migration touches:

- `cn repo install [--dispatch cds]`, package-owned admin + CDS-dispatch wakes, and Go/FSM-owned issue-claim/lifecycle are **live**; `cn cell return|resume|finalize` exist; recovery scanning/finalization run mechanically (`src/go/internal/cell/cell.go`).
- Cells run on **`cycle/{N}`** branches (verified: `cell.go` → `cycle/{Issue}`, `cycle/{N}`), **not** `cell/{N}`.
- The verifier/runtime is **tree-coupled**: `cell.go` resolves `.cdd/unreleased/{N}/` directly from the checkout; release validation resolves `.cdd/releases/{version}/{cycle}`.
- The **cds-dispatch workflow** (`.github/workflows/cnos-cds-dispatch.yml`) checks out with `actions/checkout@v4` + `sparse-checkout` (excludes `.cn-sigma`), with **no full-history setting** — shallow by default, and it excludes only `.cn-sigma`, not historical `.cdd` payloads.
- WC/PC/CC are ratified as output-telos classes of one CCNF kernel; the shipped ontology is **`cell_class` / `matter_domain` / `protocol_id`** (verified present in schemas/runtime); the old "cell kind" vocabulary is retired.
- **TSC is an external cnos tenant** installed via `cn repo install`. Any storage migration must be **versioned, dual-readable, and installer-compatible** — it cannot rely on cnos-repo-specific Git habits.

So #682 is not merely moving files: it changes the repository's **evidence database, merge protocol, release reader, recovery substrate, clone contract, and CI trust model.** That scope is why this note must be a contract, not a cleanup.

---

## 1. Governing principle

> `main`'s working tree contains only what answers **"what is."** How the project got there is warranted by **ancestry** — reachable from `main`, retrievable, auditable — but **not materialized** in the tree.

Encoded in `DOCUMENTATION-SYSTEM.md §5` + `KERNEL.md §2.1`. The distinction is *causal*, not stylistic — which is exactly why CDD receipts stay in ancestry (they warrant current product) and do **not** move to an orphan ref like the channel plane (#684), whose content is *not* product ancestry.

---

## 2. Current state (concrete)

`git ls-tree -r main -- .cdd/` = **1,022 files, 10.4 MB**, materialized in every clone.

```
.cdd/
├── CADENCE CDD-VERSION DISPATCH MCAs OPERATORS exceptions.yml …   ← config: "what is"  (STAYS)
├── skills/                                                        ← vendored bundle    (STAYS)
├── proposals/ iterations/            33 files                     ← queue/backlog: mostly current
├── releases/{X.Y.Z}/{α,β,γ}/{issue}.md   596 files · 5.2 MB       ← closed receipts    (MOVE)
├── unreleased/{N}/{self-coherence,beta-review,…}   375 · 5.0 MB   ← cells (mostly closed) (MOVE closed only)
└── waves/{slug}/…                        10 files                 ← wave matter: NOT all historical
```

**The split is not "`.cdd/` vs not."** Config + `skills/` are current-state and stay. `releases/` + *closed* `unreleased/` are the ~10.2 MB of "how we got here." `waves/` is **not universally historical** — active wave manifests are runtime authority (dispatch context, standing permissions, γ-substitutes) and must be classified, not blanket-pruned (see D6/§11).

## 3. What is wrong

| # | Problem | Evidence |
|---|---|---|
| P1 | principle violation | 10.2 MB of "how we got here" in a tree that should hold only "what is" |
| P2 | **attention, not storage** | every newcomer/agent-context/tree-walk/search wades through 971 closed files; the bytes are already Git blobs — the cost is attention |
| P3 | working-set growth | receipts grow forever *in the working tree*, not just in history |
| P4 | runtime coupling | `cell.go` + release validation read `.cdd/` from the tree; the coupling *is* the blocker |

## 4. Alternatives & leverage (L7 requires ≥3 structurally distinct moves)

| Move | What it does | Positive leverage | Negative leverage / cost | Verdict |
|---|---|---|---|---|
| **M1 · sparse-checkout + search exclusion** | dispatch/CI/search exclude closed `.cdd` payloads (extends the existing `.cn-sigma` exclusion) | removes **most agent-attention cost immediately**; ~zero new machinery; fully reversible | payloads still in HEAD tree for a normal clone; does not satisfy "main = what is"; per-consumer opt-in drifts | **Ship now** as interim relief while M4 is proven |
| **M2 · retention window** | keep only the last *N* cells materialized; older pruned | bounded working set; simpler than full ledger | arbitrary boundary; "how we got here" still partly in HEAD; still needs the seal/prune machinery for the pruned tail | Partial; subsumed by M4 |
| **M3 · release-level packing** | pack each release's receipts into one archive blob + index | large file-count reduction; simple | opaque archive breaks per-artifact review/links/validation; still in tree | Rejected as endpoint |
| **M4 · full ancestry dematerialization** | this note: seals in ancestry, HEAD keeps derived index only | satisfies the principle; unbounded receipts never enter HEAD; typed, auditable evidence DB | ongoing cost: seal+prune events, merge restrictions, full-history fetch, projector, reader/recovery code, operator-UX | **Chosen endpoint** — *only if* it earns the cost below |

**Why M4 earns its ongoing cost over M1/M2:** M1 hides receipts *per consumer* (each must opt in; new consumers re-inherit the noise); M4 makes "receipts are not present-state" a **repository invariant** enforced once, in CI, for every consumer including future tenants. The receipt-bearing model's whole value — a decision provable from evidence — degrades if evidence discovery is ad-hoc per tool. M4 buys a single typed evidence protocol; M1 buys quiet checkouts. **Recommended sequencing: ship M1 immediately; prove and ratify M4 in parallel; M1 becomes redundant once M4's steady state lands.**

## 5. Authorities (the separation that makes this safe) — resolves D1, D10

Four distinct authorities, each with a named source. No file claims another's truth.

| Authority | Owns | Source of truth | Rebuildable from main history? |
|---|---|---|---|
| **Ancestry seal events** | canonical closed-episode evidence | commits reachable from `main` carrying a seal trailer + manifest | — (it *is* the source) |
| **`.cdd/INDEX/**`** | closed-episode **finding aid** | *derived* from ancestry seal events | **Yes, exactly** (`cn cdd index --rebuild`) |
| **`.cdd/CURRENT.json`** | narrow **current CDD projection** | *declared union* — see below | **No** (depends on live inputs) |
| **radar / #683** | cross-system **open commitments** | issues + Known-Debt + deferred clauses + decisions | No (owned by #683) |
| **issue + FSM + `cycle/*` refs** | **live cell state** | GitHub/FSM/Git refs | No (live) |

**D1 fix — CURRENT is not history-derivable, so it is not claimed to be.** `INDEX` rebuilds byte-exactly from history. `CURRENT.json` is an *explicitly declared union* of `{INDEX (closed), current Git refs, issue/FSM state, release state}` — its live fields (open/in-flight cells, active runs) come from off-main sources a history walk **cannot** reconstruct. A `cycle/700` branch with `status:in-progress` and no commit in `main` is invisible to `index --rebuild` **by design** — it belongs to live-state + the #683 radar, not to a history projection.

**D10 fix — open commitments belong to the radar (#683), not CURRENT.** `CURRENT.json` carries only CDD-projection fields whose authority is the CDD ledger (current release pointer, latest sealed episode per lineage). Cross-system open items are the radar's. `cn cdd status` *joins* these sources for operator display; no single file double-claims "current."

## 6. Seal-event protocol (deterministic discovery) — resolves D2

A seal is a **typed event**, discoverable without already trusting the index.

**6.1 Commit trailers** (on the seal commit — the C4/D2 event-typing surface #682 named):

```
CN-Event: cdd.seal.v1
CN-Cell: 671
CN-Episode: 671/1
CN-Seal-Seq: 1
CN-Protocol-ID: cnos.cdd.cds.receipt.v1
CN-Manifest-Blob: <blob-oid of MANIFEST.json>
CN-Custody: CHAIN            # or CONTENT
CN-Outcome: accepted         # accepted|degraded|rejected|invalid|aborted
```

**6.2 Discovery algorithm** (canonical, deterministic):
- **Traversal:** walk **first-parent** history of `main` (`git log --first-parent`), the merge/boundary spine — not every reachable tree. Seal events live on the spine (D) or are bound by an anchor ref (below).
- **Selection:** a commit is a canonical seal iff it carries `CN-Event: cdd.seal.v1` **and** its `CN-Manifest-Blob` resolves. `(CN-Episode, CN-Seal-Seq)` is the primary key.
- **Ordering:** by first-parent topological order, ties broken by `(CN-Episode, CN-Seal-Seq)`.
- **Duplicates:** a repeated `(episode, seq)` is an error surfaced by `index --check`, never silently coalesced.
- **Amendments:** a later event with `CN-Amends: 671/1` supersedes-with-pointer; the original stays immutable (§7.3).

**6.3 Anchor refs (locator, not an evidence plane).** Optionally, a protected ref `refs/cn/cdd/seals/{episode}` points at a seal commit **already reachable from `main`.** This is a durable *locator into ancestry* to make discovery O(1) and survive first-parent edge cases — **not** an orphan evidence plane (that would sever causality, the exact thing #684 is for and this is not). If present, refs are authoritative for discovery and `--check` verifies each still points into `main` ancestry.

**Rebuild proof:** delete `INDEX/**`, retain the commit graph + anchor refs, run `--rebuild`. If it must grep arbitrary trees or read the deleted index to find seals, the projection is **not** genuinely rebuildable — that is the D2 negative test and a CI gate.

## 7. Event identity & manifest — resolves D3, C4

**7.1 Identity.** An issue number is **not** an immutable event. Identity is:

```
cell_id      671         # the cell/issue lineage
episode_id   671/1       # a specific closed episode of that lineage
seal_seq     1           # monotonic per lineage
amends / supersedes      # pointers; never mutate a sealed episode
```

Reopening or post-seal correction creates a **new episode / amendment**, never a mutation of `671/1`.

**7.2 `cn.cdd.seal.v1` manifest** (typed; reconciled with shipped `cell_class`/`matter_domain`/`protocol_id`):

```yaml
schema: cn.cdd.seal.v1
repository_id: usurobor/cnos
cell_id: "671"
episode_id: "671/1"
seal_seq: 1
protocol_id: cnos.cdd.cds.receipt.v1
cell_class: PC                       # WC | PC | CC  (output-telos; NOT "kind")
matter_domain: doctrine
contract_schema: cn.wave.v1
contract_digest: sha256:…
custody_mode: CONTENT | CHAIN
outcome: accepted|degraded|rejected|invalid|aborted
v_verdict: PASS|FAIL
delta_action: merge|reject|repair-dispatch
role_commits: {alpha: …, beta: …, gamma: …}   # required for CHAIN
artifact_blobs: [{path: …, oid: sha256:…}, …]
evidence_refs: […]
amends: null            # e.g. "671/0"
supersedes: null
residuals: […]
```

**7.3** A closed episode is immutable. The index points to `episode + amendments`.

**D3 negative test (CI):** a future reader must **not** be able to derive "planning is a cell kind" from this schema — the retired `kind` vocabulary appears nowhere; `cell_class` uses the ratified WC/PC/CC values.

## 8. Custody DAGs — resolves D4

Common invariant: **`canonical_seal ≤ boundary_decision < prune`.** The canonical seal differs by custody:

```
CHAIN  (role provenance is itself evidence):
   S(role commits + seal) ──< D(history-preserving merge) ──< P(prune)
   canonical_seal = S ;  original α/β/γ commits + exact parents remain main-reachable ;  no post-review rebase

CONTENT  (final content + receipt suffice):
   D(squash/boundary commit carrying the complete sealed payload) ──< P(prune)
   canonical_seal = D ;  the pre-squash branch S is NOT relied upon and may be GC'd
```

`index --check` verifies the *selected* mode: CHAIN → original commits reachable; CONTENT → `D` carries the full manifest+artifacts before `P`. **D4 negative:** indexing a pre-squash `S` for a CONTENT episode is rejected — the canonical seal for CONTENT is `D`.

## 9. Rejected-episode publication — resolves D5

A rejected episode's *product* must not enter `main` ancestry; its *receipt* must persist. Exact transaction:

```
1. base := current main boundary commit
2. tree := base.tree  +  only the terminal receipt + manifest   (NO product files)
3. R := commit(tree, parent = base)          # parent is main's boundary, NOT the product branch
4. R.manifest binds the rejected work by digest: rejected_ref, rejected_diff_digest
5. R records V verdict + δ reject decision (CN-Outcome: rejected)
6. index episode as rejected (seal = R)   ;   R ──< P (prune the receipt from HEAD tree)
```

The product commits stay **outside** `main` ancestry (they were never parents of `R`). **D5 negative:** an "R" that is the rejected branch tip with product files deleted is invalid — its parent chain imports every rejected product commit into ancestry.

## 10. Projection ownership & concurrency — resolves D9

`INDEX/{year}.jsonl` and `CURRENT.json` are **hot spots** — every close/release/reopen writes them; concurrent closes conflict, duplicate, or produce non-byte-identical rebuilds. Two safe models; this note selects **A**:

- **A · per-episode immutable fragments (selected).** Each seal writes its own file `INDEX/{year}/{episode_id}.json` (write-once, never appended by another cell). A **serialized `cdd-projector`** (single writer on `main`, triggered on boundary events) folds fragments into the consolidated `INDEX/{year}.jsonl` view + `CURRENT.json`. Fragments never conflict; the consolidated view is regenerated, not hand-merged.
- **B · serialized projector only** (no fragments) — acceptable but loses the immutable per-episode audit unit.

**Canonical normalization** (required for byte-exact rebuild): UTF-8, LF, sorted keys, RFC-3339-UTC timestamps sourced from the **seal commit's committer date** (not wall-clock), amendments ordered by `seal_seq`, duplicates → `--check` error.

## 11. Prune-eligibility & the one-time migration — resolves D6

No blanket `git rm -r`. A path is prune-eligible **iff** a mechanical predicate holds:

```
prune_eligible(path) :=
     terminal outcome recorded           (accepted|degraded|rejected|invalid|aborted)
  ∧  V verdict recorded
  ∧  δ boundary decision recorded
  ∧  canonical seal reachable from main   (per custody mode)
  ∧  custody check passes
  ∧  NO active repair / review-return lineage referencing the episode
  ∧  (waves only) wave terminal marker present
  ∧  episode NOT retained by a current runtime contract (active dispatch context / standing permission)
```

Phase 3 emits a **dry-run prune manifest** (every path, episode, and proof) for human/independent review **before** any deletion. Active wave manifests, in-flight cells, and repair/review-return state are protected by the predicate — never removed because they merely sit in the directory.

## 12. Retention & trust contract — resolves D7

If ancestry is the evidence database, **history rewriting is data deletion.** Pin governance:

- `main` is **non-rewriteable**; force-push prohibited (branch protection).
- Seal commits + anchor refs (`refs/cn/cdd/seals/*`) protected by repository policy.
- **≥1 independent full-history mirror**; release/export backups include seal history.
- `index --check` records the **last trusted projection root** (a checkpoint SHA) so a silent co-rewrite of history *and* index is detectable (the checkpoint won't match).
- High-assurance seals MAY require **signed commits**; actor attribution intersects **#664** (Git-hosting identity is not yet fully structural — named dependency, not solved here).

**D7 negative:** if `main` + index are force-rewritten to omit an episode and `--check` still reports success, the checkpoint-root defense has failed — `--check` must diff against the last trusted root, not only self-consistency.

## 13. Acquisition, tenants & operator UX — resolves D8, D12, C3

**13.1 Acquisition contract.** After dematerialization, ordinary shallow checkouts lack seal ancestry. Define:
- history-consuming jobs (verify, ε, release) → **full history** (`fetch-depth: 0`) **or** a **targeted fetch** of the indexed reachable seal (`cn cdd fetch {episode}`);
- offline reconstruction → full clone (this is AC1, a *test*, not the normal path);
- shallow + no fetch → typed `HISTORY_INCOMPLETE {required_sha, remediation}`, never "cell not found";
- **installer-generated workflows** (`cn repo install`) carry the correct checkout/fetch policy.

**13.2 Tenants.** TSC (and every `cn repo install` tenant) gets a **versioned, dual-readable** migration: the reader supports both tree-resident (legacy) and ancestry+index (new) layouts, gated by `CDD-VERSION`, so tenants migrate on their own cadence.

**13.3 Operator retrieval surface (D12).** `cn cdd show {episode}` and each `INDEX` entry expose: episode summary, **seal-commit permalink**, manifest, **per-artifact GitHub links at the historical commit**, amendments, PR/issue links, and the `materialize` command — so a GitHub-native operator moves from index → receipt **without** hand-composing `git show`. Operator-final-read is an authority boundary; it must stay one click deep.

**13.4 Consumer impact graph (C3).** Dematerialization changes: `cell.go`/verifier, `cn-cdd-status`, `release.sh` **and** — the ε recurrence detector (#642, scans `.cdd/releases/**`), wave-status/wave-manifest readers, repair/review-return, release PRA generation, installer-generated workflows, tenant repos (TSC), and operator-final-read surfaces. Each is a required migration target, not an afterthought.

## 14. Migration (parity-gated; corrected reversibility) — resolves D11, C2

1. **Phase 0 — Doctrine (done):** first principle in `DOCUMENTATION-SYSTEM.md §5` + `KERNEL.md §2.1` (#681).
2. **Phase 1 — Reader + schema, no deletion:** `cn cdd` reader, `cn.cdd.seal.v1`, trailer/anchor discovery, `index --rebuild/--check`.
3. **Phase 1A — Dual-write:** new cells write seal events (trailers + manifest + fragment) while the tree stays canonical.
4. **Phase 1B — Parity:** dual-read tree vs history; compare artifact set + blob digests; must match.
5. **Phase 1C — Fixtures:** accepted-CONTENT, accepted-CHAIN, rejected-receipt-only, amendment, shallow-history — each reconstructed from ancestry under the **real merge policy**.
6. **Phase 2 — Switch readers** (only after 1B parity stays clean): repoint `cell.go`/verifier/`cn-cdd-status`/`release.sh` + ε + installer workflows to `CURRENT.json` + reader. Removes nothing.
7. **Phase 3 — Dry-run prune manifest** (§11 predicate), independently reviewed.
8. **Phase 4 — Prune** only mechanically-proven terminal episodes; enable the serialized projector.
9. **Rollback material retained**; full-history mirrors verified before prune.

**C2 language:** Phases 1–2 are *independently shippable and reversible*. Phase 3+ is **not** "reversible" — it is **recoverable from verified seals** with a **pre-defined rollback** (re-materialize from ancestry). Re-materializing ~1,000 files + restoring consumers is itself a migration; call it that.

## 15. Acceptance criteria

- **AC1** `cn cdd materialize {episode}` reconstructs the complete typed cell from the **seal SHA only**, fresh full clone, no network, on a **pruned-HEAD fixture**; filesystem fallback prohibited + tested negatively.
- **AC2** `index --rebuild` reproduces `INDEX/**` byte-for-byte from ancestry (fragments + trailers/anchors only); `--check` fails on unreachable seal, digest mismatch, seal↔index non-bijection, or **checkpoint-root divergence** (D7).
- **AC3** rejected episode: product diff absent from `main` ancestry; receipt seal reachable + reconstructable; parent chain imports no product commits (§9 negative).
- **AC4** CHAIN survives history-preserving merge (original commits + parents reachable); CONTENT survives squash iff `D` carries the full sealed payload before `P`; indexing pre-squash `S` for CONTENT is rejected (§8).
- **AC5** shallow clone lacking the seal → `HISTORY_INCOMPLETE {required_sha, remediation}`; `INDEX`/`CURRENT` stay visible; targeted `cn cdd fetch` resolves it.
- **AC6** post-migration `git ls-tree -r main -- .cdd/` has **no** prune-eligible payload; only config + `skills/` + `INDEX/**` + `CURRENT.json`; **no** active wave/repair/live-cell path removed (§11 predicate honored).
- **AC7** parity (Phase 1B): tree-read and history-read produce identical artifact sets + digests across all §14 fixtures before any reader switch.
- **AC8** deterministic discovery (§6): delete `INDEX/**`, rebuild from graph + anchors only, reproduce exactly, without grepping arbitrary trees.

## 16. Risks / negative leverage

| Risk | Mitigation |
|---|---|
| squash orphans a CHAIN seal → data loss | §8 CHAIN reachability + `--check`; forbid post-review rebase |
| history rewrite deletes evidence | §12 non-rewriteable main, mirrors, checkpoint-root `--check` |
| shallow tenant misclassifies old cells | §13 acquisition contract + targeted fetch + dual-read version gate |
| projection concurrency | §10 per-episode immutable fragments + serialized projector |
| Phase-3 deletes active wave/repair | §11 prune predicate + dry-run manifest |
| CURRENT drifts from live state | §5 CURRENT = declared union, radar owns open items (#683) |
| operator can't review a pruned receipt | §13.3 index → seal permalink → artifact links |

## 17. Ratification status

This PR changes one design document; #682 remains design-first and **not dispatched**; there is no `.cdd/unreleased/682/` receipt set. That is correct for planning matter — but it is **not** a ratified architecture contract. After R1 findings clear, choose:

- **A** — keep **Proposed**, merge under an explicit human design-gate, keep #682 open for the prototype + ratification; or
- **B (recommended)** — dispatch an **independent CDS ratification cell** (independent β, no unresolved findings), then merge as canonical architecture.

Given L7 scope, **B**. Do **not** derive reader/prune implementation issues, and do **not** begin any prune, from this note until ratification.

## References

Issues: #682, #681, #683, #684, #642, #664, #627/#662. Doctrine: `DOCUMENTATION-SYSTEM.md §5`, `KERNEL.md §2.1`, `docs/architecture/CELL-RUNTIME.md`. Prior art: Git object model; event sourcing / CQRS (log + rebuildable read-models); Datomic; blockchain UTXO set; Kafka / LSM log compaction.
