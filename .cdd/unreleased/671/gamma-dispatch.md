# γ dispatch / activation record — #671 (durable provenance)

> **Purpose.** Establish the γ closure as the output of a **γ activation
> distinct from κ**, provable without relying on caller-selected Git author
> metadata. This is the R17 repair of external-β finding [BLOCKER] "A Git author
> label does not establish a distinct non-κ γ actor" (comment 5076629728).
>
> This record is authored by the **κ/δ control plane** (dispatch is a
> control-plane act; *performing* γ is not). It binds: the **exact frozen
> inputs** given to the γ activation, the **activation identity**, and the **γ
> output hash** — so a fresh reviewer can reproduce the input set and match the
> output to `gamma-closeout.md`, and establish that κ only **transported** the
> result.

## Honest scope of the claim (bootstrap limitation, disclosed)

The γ activation is a **distinct in-host activation** — a separate agent
session with its own activation id, dispatched with a frozen input set, which
**independently re-derives** the evidence (re-runs the CUE suite, the Go gate,
the ancestry checks, the hash bindings) and **could return HOLD**. It is **not**
an independent third party of a different model lineage — that independence is
**β's** guarantee (external "Codex"), and it is β's discriminating judgment the
γ closure binds. In a *manually bootstrapped* Planning Cell this distinct,
durably-recorded, re-deriving activation is the **strongest γ-independence
available**; the exemption (#671 body) names the bootstrap collapse honestly.
What this record proves is the narrow, load-bearing thing the firebreak needs:
**the γ closure was produced by an activation separate from the κ control-plane
activation, over fixed inputs, and κ only transported it** — not that κ
relabelled its own closeout.

## Frozen input set (what the γ activation was given)

| Input | Identity |
|---|---|
| **matter_sha** (frozen reviewed matter) | `614829a4682e148d98c70371e600ffdc3fa6386e` (R16); subtree `.cdd/waves/cell-runtime-doctrine/**` |
| **β source identity** | PR #672 comment `5076109763` (external, different-lineage "Codex"), verdict CONVERGE |
| **β raw-body SHA-256** | `75cdb9b67619570c668574ca27efc4107709f8c5a3f52199484033695ccfe3ae` (9,894 bytes) |
| **β payload file** | `.cdd/unreleased/671/beta-review-source-5076109763.md` (byte-exact to the source) |
| **exemption snapshot SHA-256** | `dccba69c668163b09e00ef79a77f7e6236e39cc048eca96c50fca343b507d473` (3,862 bytes) |
| **exemption payload file** | `.cdd/unreleased/671/protocol-exemption-source.md` (issue #671, `updated_at 2026-07-25T03:28:34Z`) |
| **applicable contract / matter** | the six `cn.cell.contract.v1` WC contracts + `wave.cn-wave-v1.yaml` + `schema/**` + `wave-validators/**` + `oracle-registry.yaml`, all at `matter_sha` |
| **pinned base for #628 ancestry** | `6e40d93497589afd96e6c891e94851cdabe2ef3a` |

## Dispatch instruction (verbatim — the activation's task)

> You are **γ (in-cell closer)** for the #671 Planning Cell. You are **not** κ
> and **not** α. Independently re-verify and either CLOSE (CONVERGED) or HOLD —
> reach your **own** judgment; return HOLD if any check fails.
>
> Re-run and record observed values: (1) `git diff 614829a4 -- .cdd/waves/cell-runtime-doctrine/`
> is EMPTY; (2) `make -C .cdd/waves/cell-runtime-doctrine/schema all` exit 0 (31
> negatives reject); (3) `make -C .cdd/waves/cell-runtime-doctrine/wave-validators all`
> exit 0 (78⇄78); (4) both #628 merges (`562e8025…`, `a08c56ad…`) are ancestors
> of `6e40d934…`; (5) the β payload `beta-review-source-5076109763.md` hashes to
> `75cdb9b6…` / 9,894 bytes (verdict CONVERGE, review-target `614829a4`); (6) the
> exemption payload `protocol-exemption-source.md` hashes to `dccba69c…`.
>
> Author `.cdd/unreleased/671/gamma-closeout.md` binding: `matter_sha`
> `614829a4`; the β **raw-body** hash `75cdb9b6…` (name THIS hash — not any κ
> wrapper); the exemption snapshot `dccba69c…`; the assurance PASS you
> reproduced. Distinguish **transport** (κ commits/pushes) from **authorship**
> (this activation produced the bytes). State that Git author metadata is not
> the proof of identity; this record is. Leave the sole OBSERVATION (README
> R13→R14 tail) untouched to preserve the R16 freeze. Do **not** modify the
> frozen matter; do **not** authorize/dispatch any child; do **not** act as κ.
> Sign as `γ (Planning Cell #671), in-cell closer (activation distinct from κ)`.

## Transport receipt (appended by κ after the activation completes)

_Filled in after the γ activation returns and pushes its output; records the
activation identity and the γ output hash, closing dispatch → authorship →
transport._

- **Activation identity:** _(pending — recorded post-activation)_
- **γ output file:** `.cdd/unreleased/671/gamma-closeout.md`
- **γ output SHA-256:** _(pending — recorded post-activation)_
- **Transport (κ):** _(commit that κ used to push the γ output)_
