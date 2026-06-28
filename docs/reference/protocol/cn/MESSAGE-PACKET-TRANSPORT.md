# Message Packet Transport

**Issue:** #150
**Version:** 3.31.0
**Mode:** MCA
**Active Skills:** design, inbox, writing
**Engineering Level:** L7

## Problem

Inbox materialization can silently write the wrong message content. The receiver reads a different message than the sender sent, with no error or warning.

**Evidence:** Pi sent `sigma/axiom-says-hi` with body "Axiom says hi." Sigma materialized it with body from `where-does-a-surfaced-mca-for.md` — a completely different message from a different date. Triage proceeded on wrong content. The actual message was silently dropped.

**Mechanism:** `get_branch_files` runs `git diff main...origin/<branch>` on the local peer clone. When the clone's main is stale, the diff includes every file changed since the old merge base — not just the message file. `materialize_branch` iterates all diff results and writes the first unarchived `.md` match. Wrong file first in diff = wrong content in inbox.

### Failure modes

1. **Stale clone main** — diff against old merge base returns noise files
2. **Ambiguous file selection** — materializer picks first match, not the intended file
3. **No content verification** — receiver never checks that materialized content matches what the sender declared
4. **No dedup by identity** — duplicate detection is filename-based, not message-identity-based
5. **No equivocation detection** — same message ID with different content is not caught
6. **Git-specific packet identity** — packet semantics are tied too tightly to Git objects
7. **Duplicate authority surfaces** — frontmatter and envelope can disagree
8. **Adversarial substitution risk** — integrity without authenticity is insufficient

## Constraints

### Existing contracts

- Pull-only transport: agents push to their own repo, peers fetch
- `cn_io.ml` `materialize_branch` is the current materialization entry point
- `peers.md` defines peer clones and their local paths
- Inbox skill governs triage; this design governs transport

### What can't change

- Git remains a supported transport substrate
- Pull-only fetch model remains
- Peer clone approach remains available
- Hub structure remains (`threads/mail/inbox/`, `threads/mail/outbox/`, `threads/mail/sent/`)

### Abstraction level

This is a transport protocol change. It sits below the inbox skill (triage) and above transport adapters (Git today, other transports later). The protocol defines:

- canonical packet shape
- transport-proof shape
- validation pipeline
- materialization contract

## Challenged Assumption

**The current code assumes that diffing a branch against `main` is a reliable way to discover which file is the message.**

This fails when:

- the clone's main diverges from the peer's main
- branches contain unexpected files
- multiple message-like files appear in the diff
- transport state is stale
- the sender or transport is adversarial

### Replacement assumption

A message is a self-contained packet with a canonical envelope, a content commitment, and a transport proof. The receiver validates the packet before any inbox write.

---

## Impact Graph

### Downstream consumers

| Consumer | How it's affected |
|----------|------------------|
| `materialize_branch` in `cn_io.ml` | Replaced by packet validation + exact materialization |
| `sync_inbox` in `cn_io.ml` | Fetches packet refs / transport proofs, not arbitrary branch diffs |
| `get_branch_files` in `cn_io.ml` | Removed or Phase 0 only; no longer authoritative |
| `inbox_check_once` in `cn_maintenance.ml` | Calls packet validation pipeline |
| `reject_orphan_branch` in `cn_maintenance.ml` | Uses packet schema / proof validation, not merge-base heuristics |
| `outbox_flush` / `send_thread` in `cn_mail.ml` | Produces canonical packets + transport proof |
| Inbox skill (`ops/inbox/SKILL.md`) | Updated to reflect validated packet input |
| Peer skill (`ops/peer/SKILL.md`) | Updated to reflect packet transport, not branch-diff discovery |

### Upstream producers

| Producer | What it provides |
|----------|-----------------|
| `send_thread` in `cn_mail.ml` | Canonical packet envelope + payload + transport proof |
| Transport adapter (Git) | Fetch, commit/tree/blob proof for packet |
| Future transport adapter (chain, mailbox, etc.) | Equivalent proof for packet inclusion/retrieval |

### Copies and authority

| Artifact | Authority |
|----------|-----------|
| `packet/envelope.json` | **Sole** transport-authoritative manifest |
| `packet/message.md` | Exact payload bytes named by the envelope |
| `packet/signature.ed25519` | Authenticity proof over canonical envelope |
| Inbox file frontmatter | Derived from validated envelope; never authoritative |
| Ref / tx / branch / locator | Routing / transport proof only — never authoritative for message content |
| Markdown frontmatter inside `message.md` | Optional human/application metadata only; ignored for transport validation |

---

## Proposal

### 1. Canonical Packet Schema

A message packet is transport-agnostic.

```ocaml
type packet = {
  envelope  : envelope;
  payload   : payload;
  signature : signature option;
  proof     : transport_proof;
}

type envelope = {
  schema        : string;           (* "cn.packet.v1" *)
  msg_id        : string;           (* globally unique: "{id}@{sender}", e.g. "01JZ7K@sigma" *)
  sender        : string;
  recipient     : string;
  created_at    : string;           (* ISO 8601 *)
  content_type  : string;           (* "text/markdown" *)
  payload_path  : string;           (* "packet/message.md" — fixed *)
  payload_sha256 : string;          (* canonical content commitment *)
  payload_bytes : int;
  topic         : string;
  thread        : string option;
  reply_to      : string option;
  protocol      : protocol_meta;
}

type protocol_meta = {
  transport_kind : string;          (* "git" today; extensible *)
  packet_version : int;             (* 1 *)
}

type payload = {
  sha256  : string;
  bytes   : int;
  content : string;
}

type signature = {
  alg     : string;                 (* "ed25519" *)
  signer  : string;                 (* sender key id *)
  sig_b64 : string;
}

type transport_proof =
  | Git_proof of git_proof
  | Chain_proof of chain_proof
  | Other_proof of yojson

type git_proof = {
  refname          : string;        (* refs/cn/msg/{sender}/{msg_id} *)
  commit_oid       : string;
  tree_oid         : string;
  payload_blob_oid : string;
}

type chain_proof = {
  chain_id        : string;
  contract        : string;
  tx_hash         : string;
  log_index       : int;
  block_number    : int;
  finality_depth  : int;
  payload_locator : string option;
  commitment      : string option;
}

type validation_result =
  | Valid of packet
  | Rejected of { reason_code : string; reason : string }
```

#### Key rule

The canonical content identity is `payload_sha256`, not a Git blob OID.

- Git blob OID is a transport proof detail
- SHA-256 is the transport-agnostic content commitment

This makes the packet compatible with Git, blockchain/event transports, or future mailbox transports without redesigning materialization.

---

### 2. Packet Layout

The packet tree shape is fixed.

```
packet/
  envelope.json
  message.md
  signature.ed25519   (* required when receiver policy enforces signatures; recommended always *)
```

#### Allowed file set

Exactly:

- `packet/envelope.json`
- `packet/message.md`
- optionally `packet/signature.ed25519`

No other files are allowed in Phase 1.

#### Why

If arbitrary files are allowed, ambiguity returns. The packet must be mechanically enumerable.

---

### 3. Transport Adapters

#### 3.1 Git transport (default)

Git remains the default transport adapter.

**Canonical Git ref namespace:**

```
refs/cn/msg/{sender}/{msg_id}
```

Only refs under `refs/cn/msg/` enter the inbox packet pipeline.

**Git proof invariants:**

- Ref sender must match `envelope.sender`
- Ref msg_id must match `envelope.msg_id`
- `packet/message.md` blob OID must match `proof.payload_blob_oid`
- Payload bytes must hash to `envelope.payload_sha256`

**Packet commit shape:**

Preferred: a root commit (no parents) whose tree contains only the packet files.

This removes dependence on:

- merge base
- stale main
- branch history shape

#### 3.2 Future chain transport (compatible, not required)

A chain transport may prove:

- publication
- inclusion
- ordering/finality
- external timestamping

But the canonical packet validation remains the same: envelope → payload hash → signature → proof.

**Important rule:** On-chain data is not required for payload storage. The chain layer is a transport proof / commitment layer, not the default payload store.

---

### 4. Envelope Rules

`packet/envelope.json` is the **sole** transport-authoritative manifest.

**Required fields:**

```json
{
  "schema": "cn.packet.v1",
  "msg_id": "01JZ...@sigma",
  "sender": "sigma",
  "recipient": "pi",
  "created_at": "2026-04-02T12:34:56Z",
  "content_type": "text/markdown",
  "payload_path": "packet/message.md",
  "payload_sha256": "e3b0c44298fc...",
  "payload_bytes": 123,
  "topic": "axiom-says-hi",
  "thread": "heartbeat",
  "reply_to": null,
  "protocol": {
    "transport_kind": "git",
    "packet_version": 1
  }
}
```

**Required invariants:**

- `payload_path` must equal `packet/message.md`
- `sender` must match transport proof sender identity
- `msg_id` must match transport proof route identity
- `recipient` must match the local node identity
- `payload_sha256` must match actual payload bytes
- `payload_bytes` must match actual payload length

**Important rule on frontmatter:** If `message.md` contains Markdown frontmatter, it is **ignored for transport validation**. It may be retained for human/application convenience, but inbox materialization derives its authoritative metadata from the envelope only.

This removes duplicate transport authority surfaces.

---

### 5. Signature Model

#### Receiver-local policy decides signature requirements

Whether signatures are required is determined by **local receiver policy**, not by the packet's declared fields. The sender does not control whether the receiver enforces signatures.

Receiver policy is configured per-peer (or globally) via:

- **Peer trust configuration** in the local peer registry (e.g. `peers.md` or a structured `peers.json`)
- **Transport policy** (e.g. all Git-transport packets from unknown peers require signatures)
- **Runtime default** (`secure` or `insecure`; `secure` is the recommended default)

The envelope does **not** contain a `secure_mode` field. A sender cannot opt out of signature verification.

#### Secure mode (receiver default)

When the receiver's policy requires signatures for a given peer/transport:

- `packet/signature.ed25519` is **required**
- the receiver must verify it before materialization

#### Canonical serialization for signing

The signature covers the canonical bytes of `envelope.json`. Since the envelope binds sender, recipient, payload hash, payload bytes, and protocol version, the signature transitively authenticates the message payload.

**Implementation requirement:** Phase 1 must define a single canonical envelope serialization for signing and verification. Options include deterministic JSON serialization (sorted keys, no optional whitespace, UTF-8 NFC normalization) or a separate binary canonical form. Without this, different implementations may sign different byte sequences for the same logical envelope, causing spurious verification failures.

#### Trusted key source

The receiver verifies signatures against public keys from its **local trusted peer registry**. This is a local-only trust store — not inferred from packet contents.

The trusted peer registry must provide:

- sender identity → public key mapping
- key rotation support (multiple valid keys per peer during transition)
- explicit "unknown peer" handling (reject or quarantine)

#### Insecure / dev mode

When the receiver's policy does not require signatures:

- unsigned packets may be accepted
- but this mode must be explicit and non-default

#### Rule

If the receiver's policy requires a signature and verification fails:

- reject
- quarantine
- log
- do not materialize

#### Why

Checksums/hashes prove integrity. Signatures prove authenticity. For adversarial transports, integrity alone is not enough. And the signature requirement must be receiver-controlled — a sender must not be able to downgrade the receiver's security posture.

---

### 6. Dedup and Equivocation

#### Global message identity

`msg_id` must be globally unique: `{id}@{sender}` (e.g. `01JZ7K@sigma`)

#### Inbound index

Maintain `state/inbound-index.json` keyed by:

- `msg_id`
- `sender`
- `payload_sha256`
- `status` (`accepted`, `duplicate`, `rejected`, `equivocation`)

#### Rules

- Same `msg_id` + same `payload_sha256` → duplicate; ignore safely
- Same `msg_id` + different `payload_sha256` → equivocation; reject loudly

This prevents:

- replay confusion
- duplicate materialization
- adversarial message substitution

---

### 7. Validation Pipeline (Before Materialization)

All assertions happen before any inbox write.

#### Step 1 — Fetch to quarantine

Fetch packet refs/proofs into a quarantine namespace or equivalent transport staging area.

#### Step 2 — Resolve proof

Resolve:

- transport proof
- packet tree / packet object
- sender / route identity

#### Step 3 — Validate packet shape

Reject unless:

- packet contains exactly the allowed files
- or transport-specific equivalent resolves exactly one canonical packet

#### Step 4 — Parse envelope

Reject if:

- missing
- malformed
- missing required fields
- unsupported schema/version

#### Step 5 — Validate route / identity bindings

Reject if:

- transport sender ≠ `envelope.sender`
- transport msg_id ≠ `envelope.msg_id`
- `envelope.recipient` ≠ local node identity

#### Step 6 — Load payload bytes

Resolve `packet/message.md` from the validated packet. Reject if:

- missing
- wrong path
- byte count mismatch
- SHA-256 mismatch

#### Step 7 — Validate signature (per receiver policy)

If the receiver's local policy requires a signature for this peer/transport:

- require `packet/signature.ed25519`
- verify against sender's public key from the local trusted peer registry

Reject if:

- signature required but missing
- signature present but invalid
- signer not in trusted peer registry

#### Step 8 — Dedup / equivocation check

Reject or ignore according to the dedup rules.

#### Step 9 — Materialize exact payload

Only now:

- write inbox file
- using the exact validated payload bytes
- with inbox metadata derived from the validated envelope

#### Rule

If any step fails:

- no inbox write occurs
- the packet is marked rejected/quarantined
- an operator-visible diagnostic is emitted

---

### 8. Materialization Rule

The materializer never searches for a "likely" file. It writes exactly:

- the payload bytes whose SHA-256 matches the envelope
- from the validated packet
- after all checks succeed

There is:

- no diff against main
- no "first matching .md file"
- no dependence on branch noise
- no dependence on sender repo freshness
- no dependence on receiver repo freshness

---

### 9. Send-Side Changes

`send_thread` / `outbox_flush` must produce canonical packets.

#### Git transport send flow

1. Write `packet/envelope.json`
2. Write `packet/message.md`
3. Compute `payload_sha256` and `payload_bytes`
4. Sign the envelope (recommended; receiver policy determines enforcement)
5. Create a root commit containing only the packet tree
6. Push to `refs/cn/msg/{sender}/{msg_id}`

#### Future transport send flow

Equivalent: canonical envelope, canonical payload hash, signature, transport proof.

---

### 10. Error Handling

#### Reject classes

| Code | Meaning |
|------|---------|
| `fetch_failed` | Transport fetch failed |
| `invalid_commit_or_proof` | Transport proof unresolvable |
| `invalid_tree_shape` | Packet tree contains unexpected files |
| `invalid_envelope` | Envelope missing, malformed, or wrong schema |
| `namespace_mismatch` | Transport identity ≠ envelope identity |
| `wrong_recipient` | `envelope.recipient` ≠ local node |
| `payload_mismatch` | SHA-256, byte length, or path mismatch |
| `signature_missing` | Secure mode, no signature |
| `signature_invalid` | Signature verification failed |
| `duplicate` | Same msg_id + same hash (safe ignore) |
| `equivocation` | Same msg_id + different hash (reject) |
| `unsupported_protocol` | Unknown schema or packet version |

#### Runtime behavior

On reject:

- fail early
- do not wake the agent
- do not create inbox content
- do not silently skip
- emit operator-visible diagnostics

#### Quarantine

Rejected packets may remain in quarantine for inspection, but never become inbox material.

---

### 11. Operator and Traceability Surface

Emit packet lifecycle events:

- `packet.fetched`
- `packet.validated`
- `packet.rejected`
- `packet.duplicate`
- `packet.equivocation`
- `packet.materialized`

Each event includes: `msg_id`, `sender`, `recipient`, transport proof identifiers, `payload_sha256`, `reason_code` if rejected.

This lets operators answer:

- what was received
- what was validated
- what was rejected
- why
- and whether the receiver materialized exactly the intended content

---

## 12. Prior Art Mapping

### Git

Use Git for:

- content-addressed objects
- exact commit/tree/blob identity
- local-first fetch and auditability

### BitTorrent

Adopt the core lesson:

- trust declared hashes and verified content, not filenames or peer hints

### Netnews / NNTP / store-and-forward messaging

Adopt the message-id lesson:

- every message needs a globally unique identity
- duplicates and reinjection must be detectable

### Blockchain / event transports

Use chains for:

- inclusion/finality proof
- public timestamping
- transport-level commitment

Do not require on-chain payload storage by default.

### Checksums / check codes

Optional operator check codes may be derived from `payload_sha256` for human confirmation. They are:

- useful as operator aids
- not the core safety mechanism

The core mechanism is: canonical packet → content hash → signature → transport proof → fail-closed validation.

---

## 13. Why This Is Bulletproof Against the Current Failure Class

The current failure exists because:

- branch diff is computed against stale main
- diff includes unrelated files
- materializer chooses first match

This design removes each of those assumptions:

- no diff against main
- no file discovery from branch noise
- no guessed payload path
- canonical envelope
- exact payload hash
- transport proof
- signature
- dedup/equivocation check
- fail-closed materialization

So stale clone state, unexpected files, or adversarial repo contents cannot silently swap inbox content.

---

## 14. Migration Plan

### Phase 0 — Hotfix

Before packet protocol fully lands:

- fail closed if diff-based discovery yields 0 or >1 candidate files
- never materialize "first match"
- tighten remaining legacy path to reduce noise while migration proceeds

### Phase 1 — Canonical packet protocol over Git

- dedicated `refs/cn/msg/` namespace
- root packet commits
- canonical envelope
- payload hash validation
- dedup/equivocation index
- materialize exact payload only

### Phase 2 — Signature enforcement by default

- receiver policy defaults to requiring signatures
- trusted peer registry with public keys
- operator-visible signature failures
- insecure/dev mode explicitly separate and non-default

### Phase 3 — Additional transport adapters

- chain/event proof adapters
- attachment / multipart support
- future packet schema evolution

---

## 15. Acceptance Criteria

### Phase 0

- [ ] Legacy diff path fails closed when candidate payload is ambiguous
- [ ] Legacy path never writes "first match" content

### Phase 1

- [ ] `send_thread` produces canonical packet commits under `refs/cn/msg/{sender}/{msg_id}`
- [ ] `envelope.json` contains all required fields
- [ ] Validation rejects invalid tree shape, wrong recipient, and payload hash/length mismatch
- [ ] Dedup ignores same-id same-hash duplicates
- [ ] Equivocation rejects same-id different-hash packets
- [ ] Materialized inbox content is exactly the validated payload bytes
- [ ] Materialization correctness does not depend on main freshness
- [ ] No `diff.*main...` dependency remains in the packet materialization path

### Phase 2

- [ ] Receiver policy defaults to requiring signatures for all peers
- [ ] Unsigned packets are rejected when receiver policy requires signatures
- [ ] Invalid signatures are rejected
- [ ] Sender public key verification uses the local trusted peer registry
- [ ] Insecure/dev mode requires explicit opt-in

### Transport compatibility

- [ ] Canonical packet validation is independent of transport adapter
- [ ] Git remains supported as the default transport
- [ ] Additional transport adapters can provide proof without changing materialization semantics

---

## 16. Known Debt

- Legacy non-packet message branches may need temporary compatibility handling during migration
- Dedup index garbage collection is not yet designed
- Attachments / multipart payloads deferred
- Packet schema migration strategy (`cn.packet.v2`) not yet designed

---

## 17. Summary

The right fix is not "better diffing." The right fix is:

- canonical packet
- transport proof
- payload hash
- signature
- dedup / equivocation
- fail-closed validation before materialization

Git remains the default transport adapter. But the packet itself is transport-agnostic, so the same validation/materialization contract can later work over blockchain or other append-only transports without redesign.
