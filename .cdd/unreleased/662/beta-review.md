verdict: converge   # R1 (R0 was iterate — see §R0 findings and §R1 repair attestation)

# β review — cnos#662 (`docs/architecture/CELL-RUNTIME-CLASSES.md`)

## Review provenance (bootstrap concurrency disclosure — read first)

This cell's β ran under bootstrap concurrency and this file is the **reconciled** β record over two independent β passes on the same matter:

- **β-pass-A (separate Agent activation, spawned by the δ driver in this session).** A distinct Agent activation reviewed the spec, re-verified every State-A claim directly against repo source, formed its view before reading `self-coherence.md`, and returned `verdict: converge`. Its full six-dimension attestation is retained below (§R0.1–§R0.6). **On AC5/F4 it under-weighted the CDR/CDW protocol-package State-A gap**, filing it only as a non-blocking "CDR under-claim" (its §R0.7 item 3) rather than as the AC5 miss it is.
- **β-pass-B (concurrent bootstrap firing, preserved in git history at `d66d761b`).** An earlier/concurrent β pass on the same alpha-R0 base independently walked the same oracle and returned **`verdict: iterate`**, correctly flagging **AC5 = FAIL**: the note contained zero occurrences of `CDR`/`CDW`/`#376`/`#403` or the draft's §4.4 Protocol-Package distinction that AC5 (F4) explicitly requires as State-A truth, plus a minor §9 citation nit.

The two passes disagreed on exactly one AC (AC5). **The stricter reading (β-pass-B) is correct** — AC5's oracle text ("State A reflects shipped reality … CDR shipped #376; CDW illustrative") is a positive requirement, and the spec did not meet it. The honest synthesized **R0 verdict is therefore `iterate`**, with the two findings below; R1 records the α repair and the convergence attestation. This is recorded rather than papered over so a later reader sees the AC5 gap was caught, not missed. Both β passes are separate activations but share this session's account lineage — the #664 hosting-layer α≠β limitation (see §R0.9).

---

## §R0 — findings (verdict at R0: iterate)

### F1 (AC5 / F4) — Protocol-Package State-A truth absent — **BLOCKING at R0**

**Source/evidence:** `gamma-scaffold.md §4` AC5 requires "State A reflects shipped reality (return/resume/finalize; **CDR shipped; CDW illustrative**; #504 open)"; issue #662 F4 and the draft's §4.4 κ-note require the CDS-shipped/CDR-shipped-v0.1/CDW-illustrative distinction. `grep -nE "CDR|CDW|#376|#403|Protocol Package"` over the alpha-R0 spec returned only one passing `cdr/others` mention in §6 — **zero** coverage of the protocol-package shipped/illustrative truth.
**Why it fails:** AC5 is a positive State-A-honesty requirement, not merely an anti-overclaim; omitting the CDR-shipped / CDW-illustrative truth leaves State A incomplete against the pinned F4 finding.
**What α must change:** add a "Protocol Package state truth (F4)" statement to §11.2 (State A) naming `cnos.cds` shipped (#403), `cnos.cdr` shipped v0.1 (#376), `cnos.cdw` illustrative-only, and that every worked example uses `protocol: cds`.

### F2 (§9 citation) — `doctrine_affecting` mis-cited to §2 — **non-blocking**

**Source/evidence:** §9 stated the doctrine-gate flag is "`matter_domain: doctrine` / `doctrine_affecting: true`, as #662's own contract carries — §2," but §2's worked snippet (quoting `gamma-scaffold.md §1`) shows only `matter_domain: doctrine`; `doctrine_affecting: true` appears in the operator-authorization comment's contract block, not in §2's snippet.
**Why it fails:** the §2 cross-reference is an inaccurate evidentiary pointer for the `doctrine_affecting` half.
**What α must change:** cite the operator-authorization comment's contract block for the full flag pair, or add the field to §2's snippet.

---

## §R1 — repair attestation (verdict: converge)

α (the cell's producer, repairing against the R0 findings — the normal β-iterate → α-repair loop) applied both repairs; δ re-verified each against the enumerated finding:

- **F1 repaired.** §11.2 now carries the "Protocol Package state truth (F4)" paragraph: `cnos.cds` shipped (#403); `cnos.cdr` shipped v0.1 (#376); `cnos.cdw` illustrative-only per `CDD.md` v4.0.0 (`cdo`/`cdh` future bindings); every worked example uses `protocol: cds`. Re-verified: `grep -nE "cnos.cdr|cnos.cdw|#376|#403|Protocol Package"` now matches the new paragraph; the CDR-shipped fact was independently confirmed against issue #376's shipped state, and CDS #403 / CDW-illustrative against `CDD.md`. The paragraph is verbatim-equal to the content the concurrent β-pass-B + finalize already validated as the correct AC5 remedy (preserved at `f84d155d`). **AC5 now fully met.**
- **F2 repaired.** §9 now cites the operator-authorization comment's contract block for the `doctrine_affecting: true` half and notes §2's snippet quotes only the `matter_domain: doctrine` half. Citation is now accurate.

**R1 independence note (bootstrap limitation, #664):** the R1 repair is a verbatim insertion of the exact CDR/CDW State-A truth AC5 demands plus a one-line citation correction — both mechanically verifiable against the enumerated R0 findings. R1 convergence is **δ-attested** (finding-by-finding) rather than re-verified by a freshly re-spawned β activation; this is recorded as a bootstrap-β limitation, not presented as a third independent pass. The R1 delta touches no other AC; every R0 attestation below (A–F, unchanged AC set) still holds against the R1 file.

---

## §R0.1 — Dimension A: Relation-consistency (β-pass-A attestation, retained)

ATTEST — coherent. CCNF §Scope-Lift Projection 3 match (`COHERENCE-CELL-NORMAL-FORM.md:276,:289,:297`); CELL-RUNTIME.md operationalized not restated (`:36-42`); #530 depended-on not redefined (§4:137); #500 shipped / #504 open stated correctly (`transitions.json:231-250`); #644/#654 named as targets, not overclaimed; #583/#584 cited as the mechanism/cognition boundary; κ boundary carried; shipped FSM verified in §R0.2.

## §R0.2 — Dimension B: State-A / State-B honesty (D10)

ATTEST — every State-A claim re-verified against source. `states` array `["ready","todo","in-progress","review","changes"]` exact-match `transitions.json:18`; `blocked` a `target_state` only (`:155-158`); all twelve guards match `transitions.json:19-31`; the four request-marker rules match (`:110-116,:137-142,:153-158,:161-176`) incl. #574 tightening and #368 protection; `changes → todo` repair re-entry gated `repair_contract_present` (#516, `:235-240`); `cn cell return/resume/finalize` (`cmd_cell.go`), `cn issues dispatch` (`cmd_issues_dispatch.go`), `cn issues fsm evaluate`+`scan` all ship; `run/pulse/measure/bundle/act` confirmed absent; CellKind seam observation-only (`TestSeam_CellKindNotEnforced`, `issuesfsm_test.go:810`). With the F1 repair, the CDR/CDW protocol-package State-A truth is now also present (§11.2) — closing the one State-A completeness gap.

## §R0.3 — Dimension C: No policy invented silently

ATTEST. κ≠α stated (§8:190, §12:303) and never violated; CC↔ε carried as a fenced verbatim block (§4:118-135) matching D3/CCNF Projection 3; `cell_class` post-claim with readiness label-gated only (§6:164, §2:67); contract-incompleteness → `status:blocked`+`degraded_reason`. No untraceable normative content.

## §R0.4 — Dimension D: Mechanical expressibility of class-specific guards/V (D8)

ATTEST. PC "no child auto-dispatched" → "PC applied no `status:todo` to any child" (checkable over label events); CC "no implementation surface modified" → "matter paths ⊆ judgment artifacts, no code/product diff" (checkable over artifact paths/branch diff); common floor + WC/PC additions all predicates over receipt fields/artifact paths/labels; `checks_passing` a real shipped guard (`transitions.json:28`). No V predicate is narrative-only.

## §R0.5 — Dimension E: Hard non-goals respected

ATTEST. `git diff --stat` (base `5ca785cd`) = exactly `docs/architecture/CELL-RUNTIME-CLASSES.md` + the cell's own `.cdd/unreleased/662/` artifacts. No Go/schema/FSM code; no other `docs/architecture/*.md`; no `.github/workflows/**`; no child issue filed or dispatched. §15 states the non-goals reflexively.

## §R0.6 — Dimension F: FSM + wave FSM + three-cell contracts + Coherence-Loop nomenclature + State A→B

ATTEST. Cell FSM §11.1 (shipped-grounded); wave FSM §11.4 (specified, unshipped, typed `wave_transition_request`); three-cell contracts §3.1/§3.2 (D0+Wave)/§3.3 each with input-AC/matter/output-AC/result; all five Coherence-Loop terms present and consistent (Coherence Loop §7, Cohering Cell §3.3, Cell Runtime title/§1, Cell Runner §10:210, Wake Provider §10); State A→B three-way partition (§11 shipped/specified/illustrative-future) consistent.

## §R0.7 — β-pass-A non-blocking items (retained; item 3 superseded by F1)

1. Stale cli-wrapper comment `cmd_issues_fsm.go:11` ("exactly one sub-verb") — a defect in repo source outside this cell's forbidden surface; spec is correct against the real package. Not a spec finding.
2. §2's worked-instance YAML nesting differs cosmetically from the envelope template; a faithful labeled quote of the γ-scaffold, and D9 makes field placement non-load-bearing. Cosmetic.
3. ~~"§6 cds/cdr under-claim … not actionable."~~ **Superseded — this is the AC5/F4 gap now recorded as blocking finding F1 and repaired in R1.** β-pass-A's original non-blocking classification was too lenient; the correct disposition is F1.

## §R0.8 — Contradiction-rule check

No TRUE contradiction in the pinned architecture (D1–D10 / settled findings / shipped state). The AC5 gap was an authoring omission repairable by exposition (F1), not an irreconcilable conflict. **No candidate operator-gate hold.**

## §R0.9 — β independence disclosure (#664 structural limitation)

Both β passes are **separate Agent activations** under the **same account/session lineage** as the authoring cell. Therefore α≠β holds at the **protocol layer** (neither β pass authored the matter; both re-verified State-A claims directly against repo ground truth; both formed their view before reading `self-coherence.md`) but is **bootstrap-limited at the hosting layer** (#664). The two-pass concurrency here is itself evidence for #664's thesis: a single lenient pass (β-pass-A converge) was corrected only because a second independent pass (β-pass-B iterate) caught the AC5 gap — protocol-level independence is real but hosting-level identity separation would make the check more robust. R1 convergence is δ-attested (see §R1). The later external CC ratification in the exit sequence supplies additional warrant.

---

**Verdict: converge at R1** (R0 was iterate on F1/AC5 + F2). The spec operationalizes CELL-RUNTIME.md without restating or contradicting it, honestly partitions shipped vs specified vs illustrative-future state (D10) with every State-A claim — now including the CDR/CDW protocol-package truth — verified against source, carries D1–D10 as settled input, keeps every class-specific V predicate mechanically expressible, invents no untraceable policy, and wrote exactly one spec file plus its own cell artifacts. Ready for the exit sequence (separate CC review → operator-final-read → merge → doctrine cell) the note's §16 names.
