# Pi Drive effect operations v0

**Version:** 0.1.0

**Status:** experimental bridge design; implementation pending

**Owner:** `usurobor/cn-pi`

**Owns:** the temporary Drive-staged, policy-bounded effect contract executed by the cn-pi host bridge

**Does-not-own:** dialogue, memory, direct agent credentials, project authority, or a canonical CNOS runtime protocol

## Authority, status, and experiment boundary

This document specifies a small bridge experiment. Pi emits data; the host bridge admits and executes only operator-authorized operations. Request fields ask for authority. They never grant it.

On cn-pi `main` at `6baa504bc18105f644f522d7a4393ea3449c493f`, `ops/drive-ingress/` is a proposed/new executable source path. The live `/usr/local/libexec/cn-pi-drive-ingress` is installed, unversioned host state. Its recorded SHA-256 `b16d491116eb0e00d68458e667e4124547e8460fbc1dbffd415f6cddd12e1dc0` equals the artifact on experimental branch `origin/agent/pi-drive-ingress` at `49846607f`; an implementation cycle must reproduce and publish that equality rather than trusting this host-local observation. Neither that branch nor effect execution is canonical on `main`.

The selected boundary is the smallest one that addresses a current failure: Pi cannot write private repositories, while granting it credentials would bypass operator policy. Dialogue-as-commands would mix communication with authority; direct credentials would widen authority; a CNOS scheduler/provider host has no present consumer. The selected adapter keeps credentials and policy on the host and emits evidence for later design.

Issue #2 dispatched this documentation cycle after its initial `held/not dispatched` filing state; implementation remains held. The operator explicitly exempts this docs-only cycle from a `.cdd/unreleased/2/gamma-scaffold.md`: the only product artifact is this design, while exact-head alpha/beta evidence belongs on PR #3. Implementation, Drive provisioning, deployment, and CNOS promotion require separate operator decisions.

## Current bridge and Drive topology

The installed Python bridge polls four Drive roots and materializes dialogue/memory to writer-owned Git refs. It already has stable-ID, quarantine, immutable-publication, cursor, revision-guarded append, and retry precedents. Its experimental branch contains 38 unit-test methods. It does not execute effect requests today.

Current r0 roots remain unchanged:

```text
cn-pi/r0-boxes/pi-home
cn-pi/r0-boxes/pi-cmp
cn-pi/r0-boxes/pi-tsc
cn-pi/r0-boxes/pi-cnos
```

Dialogue, memory, activation refs, and their cursors are not effect intake. Effect intake and its ledger are a fourth, separate plane. Google Drive is staging and projection; the local ledger is execution authority; GitHub and the canonical Pi writer-owned state refs are effect backends. The ledger never becomes a second activation-state authority.

The bridge's service account can read Drive and append to a shared user-owned Google Doc, but it has no personal Drive quota for per-operation files. V0 therefore uses paired, pre-created Docs rather than one Drive object per request.

## Governing sentence and invariants

> For one authenticated Drive route, persist one framed request, admit one policy-allowed operation, reconcile it against the real backend, and publish one truthful receipt without ever inferring success or repeating an uncertain effect.

- One framed record contains one operation and one issuer-scoped logical `id`.
- The stable `id` is not a digest or a Drive, Git, GitHub, ledger, marker, or projection identity.
- Operator route configuration and policy are the only grants; request assertions can only narrow them.
- Exact Doc identity and effective permissions are checked before semantic admission.
- Durable local state precedes every external mutation.
- Every discovered candidate receives a durable disposition; uncertainty remains visible.
- Terminal replay and receipt retry never re-execute an effect.
- Backend guarantees remain backend-specific. In particular, GitHub issue updates are not atomic compare-and-swap.

## Drive outbox and receipt layout

The operator pre-creates eight user-owned Google Docs:

```text
cn-pi/effect-ops/v0/outbox/Pi — Effect Outbox — {HOME|CMP|TSC|CNOS}
cn-pi/effect-ops/v0/receipts/Pi — Effect Receipts — {HOME|CMP|TSC|CNOS}
```

Folder and title organize; they do not authenticate. For each locus, host configuration pins:

```text
outbox_doc_id + owner + allowed_principal_roles + permission_fingerprint
  -> issuer agent/locus + target repository/profile
  -> receipt_doc_id + receipt_allowed_principal_roles + receipt_permission_fingerprint
```

The outbox has the operator/Pi principal as owner-writer, the bridge service account as reader, and no other writer unless explicitly allowlisted. The receipt Doc has its configured operator owner, the bridge service account as writer, Pi as reader, and only explicitly allowlisted additional writers. Credentials remain host-only.

Admission reads file metadata and every page of `permissions.list`. It normalizes each effective permission to `{permission_id,type,principal,role,inherited,deleted,pending_owner}`, sorts the array, and hashes its compact key-sorted JSON. The resulting fingerprint and effective owner/writer roles must exactly equal configuration for both Docs. The route is rechecked immediately before execution and after restart. Any missing page/field, unexpected group/domain/anyone grant, unavailable permission metadata, changed role, or fingerprint drift creates a durable local `trust_incident` and causes zero external calls. Its Drive projection remains pending until the receipt route is trusted. If an implementation cannot verify the exact effective writer boundary, it must add another authenticated write boundary; it must not weaken this check. Unsigned framing is permitted only inside this exact ACL-bound v0 route.

Pi appends complete records. The bridge reads the Docs structured-text model and records Doc ID, revision ID, and ordinal. A request starts with the exact line `<<<CN-PI-EFFECT-REQUEST-BEGIN id=<id>>>` and ends with `<<<CN-PI-EFFECT-REQUEST-END id=<same-id>>>`. Frames require matching IDs, no nesting, and exactly one intervening JSON object.

Scanning resumes at the next valid begin marker after a malformed bounded record. Discovery of a begin marker without its matching end marker commits `pending_incomplete` with Doc ID, revision, ordinal, frame ID, and zero external calls. A later revision that completes the same frame transitions that record through ordinary validation exactly once. A different begin marker permits resynchronization but does not erase the pending record. After a configured revision/age bound, an unchanged fragment transitions to `parse_incident/incomplete_frame`; later completion then requires operator reconciliation and never executes silently. Unexpected objects under `effect-ops/v0/` also produce incidents.

## One-operation request contract

Schema name: `cn-pi.effect-request.v0`. Unknown fields are denied. Required fields are `schema`, `id`, `issuer`, `target`, `resource`, `method`, `precondition`, and `digest`, plus exactly one of `query` or `representation` when the matrix requires it.

| Field | Contract |
|---|---|
| `id` | `1..128` characters from `[A-Za-z0-9._-]`; stable logical request and operation ID, scoped by issuer |
| `issuer` | exact `{agent,locus}` assertion matching the route |
| `target` | exact `{repository,profile}` assertion matching the route |
| `resource` | one exact collection/item shape below; document items include both `path` and explicit `ref` |
| `method` | one of `GET|POST|PUT|PATCH|DELETE` |
| `query` | required for collection `GET`, absent from item `GET`; exact per-family shape below |
| `representation` | mutation body; required for `POST|PUT|PATCH`, absent from reads and `DELETE` |
| `precondition` | one backend-specific shape from section 7; reads use `{"kind":"none"}` |
| `digest` | `sha256:` plus lowercase SHA-256 of the request canonical bytes with `digest` omitted |

Closure is recursive: every object below admits exactly the named members, every enum is closed, and every unknown nested key causes `parse_incident/unknown_field` with zero backend calls. Strings are UTF-8 without NUL or control characters; RFC3339 values are UTC `...Z`; SHA values are lowercase hex with the stated prefix. Route policy may reduce any maximum but never enlarge these hard bounds.

Common nested objects are exact:

| Object | Required members and bounds |
|---|---|
| `issuer` | `agent`, `locus`: repository identities `1..128` bytes matching the route |
| `target` | `repository`: `owner/name`, `3..200` bytes; `profile`: `[A-Za-z0-9._-]{1,64}`; both match the route |
| `label` | UTF-8 `1..64` bytes; at most 20 unique labels per request |
| `ref` | full `refs/heads/...`, `12..255` bytes; symbolic refs and tags are denied; mutation additionally denies `main`, protected refs, and foreign writer refs |
| `path` | normalized repository-relative UTF-8, `1..512` bytes; no empty/dot segments, backslash, NUL, symlink traversal, or percent aliases |
| `state key` | `[A-Za-z0-9][A-Za-z0-9._/-]{0,255}` mapped by route policy to one exact path/schema on the canonical state ref |

The `resource` union is exact:

| Shape | Exact object |
|---|---|
| issue collection/item | `{"kind":"issue"}` / `{"kind":"issue","number":1..2147483647}` |
| comment collection/item | `{"kind":"comment","issue_number":1..2147483647}` / plus `"comment_id":1..9223372036854775807` |
| document collection/item | `{"kind":"document"}` / `{"kind":"document","path":"<path>","ref":"<ref>"}` |
| state collection/item | `{"kind":"state"}` / `{"kind":"state","key":"<state key>"}` |

The `query` union is also closed. Item `GET` carries no `query`. Collection queries require only the listed object; omitted optional members take the shown fixed default before admission and are retained in the normalized operation.

| Collection | Exact query members |
|---|---|
| issue | optional `state: open|closed|all` (`open`), `labels: label[]` (`[]`), `sort: created|updated` (`updated`), `direction: asc|desc` (`desc`), `since: RFC3339`, `limit: 1..100` (`30`), `cursor: string|null` (`null`, max 512 bytes) |
| comment | required `issue_number` is in the resource; optional `since`, `limit`, and `cursor` with the same bounds/defaults |
| document | required `ref`; optional `prefix` (`""`, path-prefix rules), `limit`, and `cursor` with the same bounds/defaults |
| state | optional `prefix` (`""`, state-key rules), `limit`, and `cursor`; the route supplies the canonical state ref |

Mutation `representation` objects are exact:

| Operation | Required and optional members |
|---|---|
| issue `POST` | required `title: 1..256` bytes; optional `body: 0..65536` bytes and `labels: label[]`; defaults are empty |
| issue `PATCH` | one or more of `title`, `body`, `labels`, `state: open|closed`, `state_reason: completed|not_planned|reopened`; `state_reason` is permitted only with `state` and must agree with the transition |
| comment `POST` | required `body: 1..65536` bytes |
| document `PUT` | required `encoding:"utf-8"`, route-allowlisted `media_type: 1..127` bytes, `content: 0..1048576` bytes, and matching `content_sha256` |
| state `PUT` | required `encoding:"utf-8"`, route-pinned `media_type: application/json|application/yaml`, `content: 0..262144` bytes, and matching `content_sha256` |

`DELETE` and reads carry no representation. The precondition union is exact:

| Kind | Exact members |
|---|---|
| `none` | `kind` only |
| `marker_absent` | `kind` only |
| `observed_issue` | `kind`, `updated_at`, `representation_sha256` |
| `document_match` | `kind`, `ref`, `head_sha`, `blob_sha`, `content_sha256` |
| `document_absent` | `kind`, `ref`, `head_sha` |
| `state_match` | `kind`, `ref_head_sha`, `content_sha256` |
| `state_absent` | `kind`, `ref_head_sha` |

Canonicalization is a bridge-local procedure, not a generic CNOS subsystem: parse strict UTF-8 as exactly one JSON object; reject duplicate keys, floats/exponents, non-finite numbers, unknown fields, BOMs, and integers outside signed 64-bit range; remove the top-level digest field; recursively sort object keys by Unicode code point; preserve array order; emit UTF-8 JSON with no insignificant whitespace or final newline; SHA-256 those bytes. Receipt digests use the same procedure with their top-level `digest` omitted. No Unicode normalization or adapter inference occurs.

Complete synthetic request fixture:

```text
<<<CN-PI-EFFECT-REQUEST-BEGIN id=req-cn-pi-tsc-20260807-0001>>>
{"schema":"cn-pi.effect-request.v0","id":"req-cn-pi-tsc-20260807-0001","issuer":{"agent":"usurobor/cn-pi","locus":"usurobor/tsc"},"target":{"repository":"usurobor/tsc","profile":"issues-and-docs-v0"},"resource":{"kind":"issue"},"method":"POST","representation":{"title":"Define customer-case realization proof","body":"Capture the bounded customer-case realization and its proof plan.","labels":["documentation"]},"precondition":{"kind":"marker_absent"},"digest":"sha256:b245954dcd2e1cb10d6c83319b1ef1b0eca73fe0370486b149eac057f306a467"}
<<<CN-PI-EFFECT-REQUEST-END id=req-cn-pi-tsc-20260807-0001>>>
```

The canonical bytes begin `{"id":...` because keys are sorted. Independent recomputation yields `sha256:b245954dcd2e1cb10d6c83319b1ef1b0eca73fe0370486b149eac057f306a467`.

The other three resource families have complete positive request fixtures below. Their listed digests independently recompute under the same canonicalizer.

```json
{"schema":"cn-pi.effect-request.v0","id":"req-cn-pi-tsc-comment-0001","issuer":{"agent":"usurobor/cn-pi","locus":"usurobor/tsc"},"target":{"repository":"usurobor/tsc","profile":"issues-and-docs-v0"},"resource":{"kind":"comment","issue_number":123},"method":"POST","representation":{"body":"Please review the exact candidate head."},"precondition":{"kind":"marker_absent"},"digest":"sha256:f943b68de45a1225462f6531d660c209dff3f5b2ef422708e615c03b45cec91b"}
{"schema":"cn-pi.effect-request.v0","id":"req-cn-pi-tsc-document-0001","issuer":{"agent":"usurobor/cn-pi","locus":"usurobor/tsc"},"target":{"repository":"usurobor/tsc","profile":"issues-and-docs-v0"},"resource":{"kind":"document","path":"docs/customer-note.md","ref":"refs/heads/agent/pi-docs-tsc"},"method":"PUT","representation":{"encoding":"utf-8","media_type":"text/markdown","content":"# Customer note\n","content_sha256":"sha256:f47fe3326336ccea4614d8ecad937f6aed6034dba315f35475d4a6b96a40428a"},"precondition":{"kind":"document_absent","ref":"refs/heads/agent/pi-docs-tsc","head_sha":"1111111111111111111111111111111111111111"},"digest":"sha256:6a9a62ed16ccde4e9642c9361c4222f9bd3dfd72771fcb663440d77f2c4ca47e"}
{"schema":"cn-pi.effect-request.v0","id":"req-cn-pi-home-state-0001","issuer":{"agent":"usurobor/cn-pi","locus":"usurobor/cn-pi"},"target":{"repository":"usurobor/cn-pi","profile":"home-state-v0"},"resource":{"kind":"state","key":"tasks/example"},"method":"PUT","representation":{"encoding":"utf-8","media_type":"application/yaml","content":"state: ready\n","content_sha256":"sha256:1d013e814cdedbb7aa2469e36e5ca09435030a56b47a9ceaa480ea21c337c348"},"precondition":{"kind":"state_absent","ref_head_sha":"2222222222222222222222222222222222222222"},"digest":"sha256:c123752329cca84164a4316ebbdf872ec14d5546bbacaea3fdff77a32a1bcf55"}
```

Negative nested-closure fixture: although its digest is internally consistent, the unknown `representation.admin` member produces `parse_incident/unknown_field` and zero backend calls.

```json
{"schema":"cn-pi.effect-request.v0","id":"req-negative-nested-key-0001","issuer":{"agent":"usurobor/cn-pi","locus":"usurobor/tsc"},"target":{"repository":"usurobor/tsc","profile":"issues-and-docs-v0"},"resource":{"kind":"issue"},"method":"POST","representation":{"title":"Denied","admin":true},"precondition":{"kind":"marker_absent"},"digest":"sha256:12327256cff079b2b326955b5cd61a9579c639fed184368a1870eb9ab7b9bbf6"}
```

Negative document-target fixture: this request is syntactically closed and its digest is internally consistent, but `resource.ref` differs from `precondition.ref`. Admission produces `denied/stale_precondition` and zero mutation calls even when both refs independently fall under the configured candidate prefix.

```json
{"schema":"cn-pi.effect-request.v0","id":"req-negative-document-ref-mismatch-0002","issuer":{"agent":"usurobor/cn-pi","locus":"usurobor/tsc"},"target":{"repository":"usurobor/tsc","profile":"issues-and-docs-v0"},"resource":{"kind":"document","path":"docs/customer-note.md","ref":"refs/heads/agent/pi-docs-tsc-b"},"method":"DELETE","precondition":{"kind":"document_match","ref":"refs/heads/agent/pi-docs-tsc-a","head_sha":"1111111111111111111111111111111111111111","blob_sha":"2222222222222222222222222222222222222222","content_sha256":"sha256:f47fe3326336ccea4614d8ecad937f6aed6034dba315f35475d4a6b96a40428a"},"digest":"sha256:ce2b5837bd45bb1182bda01b9367604a84864929db6687035c935916e51e6861"}
```

## Resource, method, and policy matrix

This matrix is closed. The 14 listed resource-shape/method pairs are the entire vocabulary; every other pair is denied.

| # | Resource shape | Method | Input and meaning |
|---:|---|---|---|
| 1 | issue collection `{kind}` | `GET` | `query`; list issues |
| 2 | issue item `{kind,number}` | `GET` | retrieve issue; no `query` |
| 3 | issue collection `{kind}` | `POST` | `representation`; create issue |
| 4 | issue item `{kind,number}` | `PATCH` | `representation`; update, close, or reopen |
| 5 | comment collection `{kind,issue_number}` | `GET` | `query`; list comments |
| 6 | comment item `{kind,issue_number,comment_id}` | `GET` | retrieve comment; no `query` |
| 7 | comment collection `{kind,issue_number}` | `POST` | `representation`; create comment |
| 8 | document collection `{kind}` | `GET` | `query`; list allowed non-code documents |
| 9 | document item `{kind,path,ref}` | `GET` | retrieve document from an admitted read ref |
| 10 | document item `{kind,path,ref}` | `PUT` | `representation`; create or replace on an admitted candidate ref |
| 11 | document item `{kind,path,ref}` | `DELETE` | delete on an admitted candidate ref |
| 12 | state collection `{kind}` | `GET` | `query`; list Pi activation state keys |
| 13 | state item `{kind,key}` | `GET` | retrieve state; no `query` |
| 14 | state item `{kind,key}` | `PUT` | `representation`; create or replace state |

The route fixes the repository and profile. Policy then intersects exact resource/method pairs with repository attachment, issue fields/labels, maximum body and result sizes, document path prefixes and media types, exact read refs, one exact pre-created write ref or candidate-ref prefix, state key mappings, pagination, and per-route rate/budget limits. A request ref must equal or narrow that grant. Home may target only `usurobor/cn-pi`; an attached route targets only its bound repository.

State-ref authority is copied, never reconstructed. Route configuration selects exactly one activation in canonical `state/activations.yaml` on cn-pi `main` by exact `activation_key: [agent,locus]`, requires `agent == issuer.agent == usurobor/cn-pi` and `locus == issuer.locus`, and pins the entry's literal `feeds.state` value plus the registry commit SHA. For example, locus identity `usurobor/tsc` binds the literal `refs/heads/cn-pi/tsc/state`; no substitution of the locus string and no parallel ref convention is allowed. Every state read/write verifies that unchanged activation key, registry commit, and exact `feeds.state` binding before backend access. A changed or missing binding is `trust_incident` with zero state mutation calls until an operator admits a new policy/config digest.

Document policy is positive: a path and media type must both be listed as non-code. Extension alone never grants. For every document `PUT|DELETE`, `resource.ref` must byte-equal `precondition.ref`; mismatch is denied before mutation even when both refs independently satisfy the candidate grant. Mutation of `main`, tags, protected refs, another writer's refs, or an unpinned/default branch is denied. An attached-locus write uses an operator-pre-created candidate ref; if none exists, document mutation is disabled because branch creation is not a v0 primitive. A configured PR may consume that candidate ref, but PR creation/merge is outside v0. Code, workflows, build/config files, symlinks, binary media, policy files, raw HTTP/Git/shell, credentials, repository settings, and unattached repositories are denied. Issue deletion and comment update/delete do not exist. A comment `issue_number` must identify an ordinary issue, not a pull request. Policy cannot be mutated by this protocol.

V0 does not claim end-to-end pull-request review agency. It has no PR head/base, diff, changed-file, check, review, or thread observation resource and cannot issue a merge-gating verdict. Until a separately reviewed bounded PR-observation family exists, PR review remains operator-relayed and PR conversation comments are denied. This explicit non-goal preserves the 14-pair vocabulary instead of hiding PR aggregation behind issue `GET`.

## Backend preconditions and reconciliation

| Operation | Required precondition and execution truth |
|---|---|
| any `GET` | `none`; read exact backend state, paginate within policy, and receipt typed inline data plus version/completeness evidence |
| issue `POST` / comment `POST` | `marker_absent`; derive a stable hidden marker from `agent NUL locus NUL id NUL request_digest`, query it before first attempt and every retry, create once, then exact-read the returned item |
| issue `PATCH` | `observed_issue` with `updated_at` and representation digest; immediate preflight equality, PATCH, exact post-read; the interval between preflight and mutation is non-atomic and profiles may deny updates that need stronger exclusion |
| document `PUT` existing | `document_match` with blob/content SHA, `precondition.ref == resource.ref`, route-admitted candidate ref, and observed head SHA; preflight all evidence, send the backend-required blob SHA, and post-read returned commit/blob/head |
| document `PUT` new | `document_absent` with `precondition.ref == resource.ref`, route-admitted candidate ref, and observed head SHA; recheck absence/head, create, and post-read exact bytes and returned commit |
| document `DELETE` | `document_match` with blob/content SHA, `precondition.ref == resource.ref`, route-admitted candidate ref, and observed head SHA; send the required blob SHA and verify absence plus returned commit |
| state `PUT` | `state_match` with the route-pinned activation entry's literal `feeds.state` head and content digest, or `state_absent` with head; map key to the pinned state path/schema, create one commit on the observed head, and perform a fast-forward-only ref update; receipt before/after heads and content digests |

GitHub's Contents API enforces the existing blob SHA for update/delete, but the observed branch head is preflight and reconciliation evidence unless the implementation selects a primitive that enforces it. Candidate-ref mutation therefore remains subject to a disclosed head race and must post-read the returned commit. State writes use Git's fast-forward rejection as optimistic compare: a commit built on the expected state head cannot replace a concurrently advanced ref. A conflicting head, blob, issue representation, or state version is `denied/stale_precondition`, never silently overwritten.

For creates, the bridge injects only the correlation marker into the created body. The example marker is `cn-pi-effect:c2d69ae6de3900171057ffd05c0d50d6dc09c59f165543f680dad68cf78a326e`. Query can prove an existing matching effect; an absent result after an ambiguous timeout does not necessarily prove non-commit. If the backend cannot prove success or absence, the request remains `indeterminate` and is not automatically replayed.

Backend references: [GitHub issue operations](https://docs.github.com/en/rest/issues/issues), [issue comments](https://docs.github.com/en/rest/issues/comments), [repository contents](https://docs.github.com/en/rest/repos/contents), [Google Drive permissions](https://developers.google.com/workspace/drive/api/reference/rest/v3/permissions/list), and [Google Docs revision-guarded writes](https://developers.google.com/workspace/docs/api/reference/rest/v1/documents/batchUpdate).

## Durable ledger, replay, collision, and receipts

The local SQLite ledger is authoritative. Discovery is inserted and committed before admission; admission is committed before execution; `executing` plus an attempt token is committed before the external call. Admission atomically retains source Doc/revision/ordinal, scoped ID, the exact canonical request bytes, normalized admitted-operation bytes, request and operation digests, immutable route-policy/config digest, normalized precondition, backend observations, attempt history, timestamps, and receipt projection state. Restart and reconciliation use only these retained bytes; mutable Drive text is never reread as execution input.

| Dimension | States |
|---|---|
| delivery | `discovered`, `pending_incomplete` |
| admission | `pending`, `admitted`, `denied`, `parse_incident`, `trust_incident`, `policy_incident` |
| execution | `not_started`, `executing`, `succeeded`, `incomplete`, `failed`, `indeterminate` |
| projection | `pending`, `projected` |

`succeeded`, `failed`, and `denied` are terminal only with a reason and retained evidence. `indeterminate` is nonterminal uncertainty requiring reconciliation or operator resolution. Exact scoped-ID + digest replay returns the retained state/receipt. Same scoped ID + different digest is `id_collision` and makes zero external calls. Different IDs with identical content are distinct intent.

After restart, `executing` is reconciled before any retry. Proven matching external state becomes `succeeded`; proven non-commit may retry under the original policy and budget; otherwise it becomes `indeterminate`. A terminal request never executes again.

Receipts use the recursively closed `cn-pi.effect-receipt.v0`. The paired Drive journal is a revision-guarded, reread-verified projection. Projection failure records `receipt_projection_pending`; retries read only the ledger and never call the effect backend.

Every request receipt has exactly these common fields: `schema`, `id`, `request_digest`, `status`, `reason`, `source`, `ledger_seq`, `policy`, `admission`, `attempts`, `backend`, one status payload, and `digest`. `source` is exact `{doc_id,revision_id,record_ordinal}` with strings `1..256` bytes and ordinal `0..9223372036854775807`. `policy` is exact `{profile,digest}`. `admission` is exact `{canonical_request_sha256,normalized_operation_sha256,precondition}`. Each attempt is exact `{number,started_at,finished_at,outcome,backend_request_id?}` with number `1..16`, RFC3339 timestamps, optional request ID `1..256` bytes, and outcome `observed|no_commit|committed|ambiguous`; attempts are ordered, nonempty after a backend call, and bounded by policy. `backend` is exact `{kind:github-rest|git-ff,executor_identity,request_ids}` with identity `1..128` bytes and at most 16 request IDs. `ledger_seq` is a positive signed 64-bit integer; every digest uses the canonical SHA form.

The status union is discriminated and closed:

| Status | Allowed reason | Required payload | Forbidden payloads |
|---|---|---|---|
| `succeeded` read | `read_complete|not_found` | `observation` with `complete:true`, `truncated:false`; `not_found` requires item `count:0`, `version:null`, and `data:null` | `mutation`, `denial`, `error`, `reconciliation` |
| `incomplete` read | `pagination_limit|result_budget_exhausted` | `observation` with `complete:false`, `truncated:true`, partial inspectable data, and shape-appropriate continuation evidence | `mutation`, `denial`, `error`, `reconciliation` |
| `succeeded` mutation | `created|updated|deleted|reconciled` | `mutation` | `observation`, `denial`, `error` |
| `denied` | `policy_denied|trust_denied|stale_precondition|id_collision|unsupported_pair` | `denial:{rule,evidence_sha256}` | `observation`, `mutation`, `error` |
| `failed` | `backend_rejected` | `error:{code,message,commit_proven_absent:true}` | `observation`, `mutation`, `denial` |
| `indeterminate` | `ambiguous_commit|remote_mismatch` | `reconciliation:{last_observation_sha256,next_action:observe|operator}` | `observation`, `mutation`, `denial`, `error` |

An observation is exact `{resource_kind,shape:item|collection,complete,truncated,count,next_cursor,version,data}` and is at most 262144 canonical UTF-8 bytes after policy narrowing. `count` equals the number of returned items. Item reads require `count:0|1` and `next_cursor:null`. A complete existing item requires `reason:read_complete`, `complete:true`, `truncated:false`, `count:1`, and non-null resource-specific `version` and item `data`. An oversized existing item requires `reason:result_budget_exhausted`, `complete:false`, `truncated:true`, `count:1`, the same non-null version and partial closed item data; its stable item identity/version is the continuation evidence for a later larger-budget request. A missing item requires `reason:not_found`, `count:0`, `version:null`, and `data:null`; this is the only observation shape that permits null version/data. A complete collection requires `reason:read_complete`, `complete:true`, `truncated:false`, and `next_cursor:null`. An incomplete collection requires `reason:pagination_limit|result_budget_exhausted`, `complete:false`, `truncated:true`, and the backend continuation token in `next_cursor`. No digest-only result counts as perception. V0 uses bounded inline data rather than mutable artifact pointers.

Closed item `data` views are: issue `{number,title,body,state,labels,updated_at,url}`; comment `{comment_id,issue_number,body,author,updated_at,url}`; document `{path,ref,media_type,encoding,content,content_sha256,blob_sha,head_sha}`; state `{key,path,ref,media_type,encoding,content,content_sha256,head_sha}`. Collection `data` is an array of closed summaries: issue `{number,title,state,labels,updated_at,url}`; comment `{comment_id,issue_number,author,updated_at,url}`; document `{path,ref,media_type,content_sha256,blob_sha,head_sha}`; state `{key,path,ref,media_type,content_sha256,head_sha}`. Request field bounds apply to returned titles, bodies, labels, paths, refs, media, and content; author is `1..128` bytes and URL is `1..2048` bytes. When an item body exceeds the inline budget, the receipt is `incomplete/result_budget_exhausted`, includes the remaining metadata and an empty content/body field, and never claims completeness.

Item `version` and present mutation observation points are resource-specific closed objects: issue `{number,updated_at,representation_sha256}`; comment `{comment_id,updated_at,representation_sha256}`; document `{ref,head_sha,blob_sha,content_sha256}`; state `{ref,head_sha,path,content_sha256}`. Collection `version` is also closed: issue/comment use `{query_started_at,query_finished_at,items_sha256}`; document/state use `{ref,head_sha,items_sha256}`. Timestamps are RFC3339 with finish not before start. `items_sha256` is the canonical SHA-256 of the exact returned `data` array, including a partial page, and makes no atomic-snapshot claim for GitHub issue/comment enumeration; document/state enumeration is pinned to the named ref head.

A `mutation` is exact `{resource_kind,observed_before,observed_after,correlation}`. Before may be null only for a proven create. Create/update after-points use the present resource-specific objects above. Document delete alone uses the exact absence after-point `{ref,head_sha,path,absent:true}`, where `head_sha` is the returned deletion commit/ref head and an exact post-read proves the path absent. Correlation is exactly `{marker}` for issue/comment `POST`, `{backend_request_id}` for issue `PATCH`, or `{commit_sha}` for document/state writes; each ID is `1..256` bytes or the required 40-hex commit SHA. Denial `rule` and error `code` are `[A-Za-z0-9._/-]{1,128}`; error message is `1..4096` bytes.

Complete issue-create receipt fixture:

```text
<<<CN-PI-EFFECT-RECEIPT-BEGIN id=req-cn-pi-tsc-20260807-0001>>>
{"schema":"cn-pi.effect-receipt.v0","id":"req-cn-pi-tsc-20260807-0001","request_digest":"sha256:b245954dcd2e1cb10d6c83319b1ef1b0eca73fe0370486b149eac057f306a467","status":"succeeded","reason":"created","source":{"doc_id":"1ExampleTscOutboxDoc000000000","revision_id":"rev-outbox-20260807-01","record_ordinal":7},"ledger_seq":41,"policy":{"profile":"issues-and-docs-v0","digest":"sha256:4444444444444444444444444444444444444444444444444444444444444444"},"admission":{"canonical_request_sha256":"sha256:b245954dcd2e1cb10d6c83319b1ef1b0eca73fe0370486b149eac057f306a467","normalized_operation_sha256":"sha256:3333333333333333333333333333333333333333333333333333333333333333","precondition":{"kind":"marker_absent"}},"attempts":[{"number":1,"started_at":"2026-08-07T12:00:03Z","finished_at":"2026-08-07T12:00:04Z","outcome":"committed","backend_request_id":"gh-req-0001"}],"backend":{"kind":"github-rest","executor_identity":"host:cn-node-01","request_ids":["gh-req-0001"]},"mutation":{"resource_kind":"issue","observed_before":null,"observed_after":{"number":84,"updated_at":"2026-08-07T12:00:04Z","representation_sha256":"sha256:c578c40a6fc04f1e7ba1ca1f6502bf56d69f734a3ee2e0c75c8a3adecf7f8e3e"},"correlation":{"marker":"cn-pi-effect:c2d69ae6de3900171057ffd05c0d50d6dc09c59f165543f680dad68cf78a326e"}},"digest":"sha256:68b63f6774e4400441520181fa58a17517ba687e3eb28212bbeeb50daa5d4f77"}
<<<CN-PI-EFFECT-RECEIPT-END id=req-cn-pi-tsc-20260807-0001>>>
```

Independent canonicalization with receipt `digest` omitted yields `sha256:68b63f6774e4400441520181fa58a17517ba687e3eb28212bbeeb50daa5d4f77`.

Complete item-`GET` request/receipt pair, proving inspectable perception rather than digest-only verification:

```json
{"schema":"cn-pi.effect-request.v0","id":"req-cn-pi-tsc-issue-get-0002","issuer":{"agent":"usurobor/cn-pi","locus":"usurobor/tsc"},"target":{"repository":"usurobor/tsc","profile":"issues-and-docs-v0"},"resource":{"kind":"issue","number":84},"method":"GET","precondition":{"kind":"none"},"digest":"sha256:1d59a42dcb9a84355e370d7ef33d61d77d4b1ba0af363fb46cd156f52765c08e"}
{"schema":"cn-pi.effect-receipt.v0","id":"req-cn-pi-tsc-issue-get-0002","request_digest":"sha256:1d59a42dcb9a84355e370d7ef33d61d77d4b1ba0af363fb46cd156f52765c08e","status":"succeeded","reason":"read_complete","source":{"doc_id":"1ExampleTscOutboxDoc000000000","revision_id":"rev-outbox-20260807-02","record_ordinal":8},"ledger_seq":42,"policy":{"profile":"issues-and-docs-v0","digest":"sha256:4444444444444444444444444444444444444444444444444444444444444444"},"admission":{"canonical_request_sha256":"sha256:1d59a42dcb9a84355e370d7ef33d61d77d4b1ba0af363fb46cd156f52765c08e","normalized_operation_sha256":"sha256:5555555555555555555555555555555555555555555555555555555555555555","precondition":{"kind":"none"}},"attempts":[{"number":1,"started_at":"2026-08-07T12:01:03Z","finished_at":"2026-08-07T12:01:04Z","outcome":"observed","backend_request_id":"gh-req-0002"}],"backend":{"kind":"github-rest","executor_identity":"host:cn-node-01","request_ids":["gh-req-0002"]},"observation":{"resource_kind":"issue","shape":"item","complete":true,"truncated":false,"count":1,"next_cursor":null,"version":{"number":84,"updated_at":"2026-08-07T12:00:04Z","representation_sha256":"sha256:931e34bac37d3a336e0b9f59e5d1165afcf36a31b5623402bbf2084b43c865ab"},"data":{"number":84,"title":"Define customer-case realization proof","body":"Capture the bounded customer-case realization and its proof plan.","state":"open","labels":["documentation"],"updated_at":"2026-08-07T12:00:04Z","url":"https://github.com/usurobor/tsc/issues/84"}},"digest":"sha256:5d6f63691ee25a8f3c4e5b1303b3fbaaac9fe8df02e7eb5d3093aeb556942af6"}
```

Complete missing-item `GET` request/receipt pair:

```json
{"schema":"cn-pi.effect-request.v0","id":"req-cn-pi-tsc-issue-get-missing-0003","issuer":{"agent":"usurobor/cn-pi","locus":"usurobor/tsc"},"target":{"repository":"usurobor/tsc","profile":"issues-and-docs-v0"},"resource":{"kind":"issue","number":85},"method":"GET","precondition":{"kind":"none"},"digest":"sha256:65fe0690bf41737bf2fb6c6d8864d170da3269e20d63482fd88974c1a64352bb"}
{"schema":"cn-pi.effect-receipt.v0","id":"req-cn-pi-tsc-issue-get-missing-0003","request_digest":"sha256:65fe0690bf41737bf2fb6c6d8864d170da3269e20d63482fd88974c1a64352bb","status":"succeeded","reason":"not_found","source":{"doc_id":"1ExampleTscOutboxDoc000000000","revision_id":"rev-outbox-20260807-03","record_ordinal":9},"ledger_seq":43,"policy":{"profile":"issues-and-docs-v0","digest":"sha256:4444444444444444444444444444444444444444444444444444444444444444"},"admission":{"canonical_request_sha256":"sha256:65fe0690bf41737bf2fb6c6d8864d170da3269e20d63482fd88974c1a64352bb","normalized_operation_sha256":"sha256:6666666666666666666666666666666666666666666666666666666666666666","precondition":{"kind":"none"}},"attempts":[{"number":1,"started_at":"2026-08-07T12:02:03Z","finished_at":"2026-08-07T12:02:04Z","outcome":"observed","backend_request_id":"gh-req-0003"}],"backend":{"kind":"github-rest","executor_identity":"host:cn-node-01","request_ids":["gh-req-0003"]},"observation":{"resource_kind":"issue","shape":"item","complete":true,"truncated":false,"count":0,"next_cursor":null,"version":null,"data":null},"digest":"sha256:fd0d4ab394cda8f41f9792d82b774fcd93f256241f996b9d0f1d5c51b4e0e040"}
```

Complete and incomplete collection `GET` pairs prove both collection-version forms used by the status union:

```json
{"schema":"cn-pi.effect-request.v0","id":"req-cn-pi-tsc-issues-get-complete-0004","issuer":{"agent":"usurobor/cn-pi","locus":"usurobor/tsc"},"target":{"repository":"usurobor/tsc","profile":"issues-and-docs-v0"},"resource":{"kind":"issue"},"method":"GET","query":{"state":"open","labels":[],"sort":"updated","direction":"desc","limit":30,"cursor":null},"precondition":{"kind":"none"},"digest":"sha256:062269d6fd7c9037842b61014c6e8851b405e4bff609c2977e501fa84e7469d2"}
{"schema":"cn-pi.effect-receipt.v0","id":"req-cn-pi-tsc-issues-get-complete-0004","request_digest":"sha256:062269d6fd7c9037842b61014c6e8851b405e4bff609c2977e501fa84e7469d2","status":"succeeded","reason":"read_complete","source":{"doc_id":"1ExampleTscOutboxDoc000000000","revision_id":"rev-outbox-20260807-04","record_ordinal":10},"ledger_seq":44,"policy":{"profile":"issues-and-docs-v0","digest":"sha256:4444444444444444444444444444444444444444444444444444444444444444"},"admission":{"canonical_request_sha256":"sha256:062269d6fd7c9037842b61014c6e8851b405e4bff609c2977e501fa84e7469d2","normalized_operation_sha256":"sha256:7777777777777777777777777777777777777777777777777777777777777777","precondition":{"kind":"none"}},"attempts":[{"number":1,"started_at":"2026-08-07T12:03:03Z","finished_at":"2026-08-07T12:03:04Z","outcome":"observed","backend_request_id":"gh-req-0004"}],"backend":{"kind":"github-rest","executor_identity":"host:cn-node-01","request_ids":["gh-req-0004"]},"observation":{"resource_kind":"issue","shape":"collection","complete":true,"truncated":false,"count":1,"next_cursor":null,"version":{"query_started_at":"2026-08-07T12:03:03Z","query_finished_at":"2026-08-07T12:03:04Z","items_sha256":"sha256:7438f2960e3bd8d7354f53e11bbd95e8e0dc0ceaf65d1f7269f7404d3e13a5f1"},"data":[{"number":84,"title":"Define customer-case realization proof","state":"open","labels":["documentation"],"updated_at":"2026-08-07T12:00:04Z","url":"https://github.com/usurobor/tsc/issues/84"}]},"digest":"sha256:932d0317bdfba0e90bc363f860e99f48a491f653ea48944dbde2f0f7bf1ffd70"}
{"schema":"cn-pi.effect-request.v0","id":"req-cn-pi-tsc-issues-get-incomplete-0005","issuer":{"agent":"usurobor/cn-pi","locus":"usurobor/tsc"},"target":{"repository":"usurobor/tsc","profile":"issues-and-docs-v0"},"resource":{"kind":"issue"},"method":"GET","query":{"state":"all","labels":[],"sort":"updated","direction":"desc","limit":1,"cursor":null},"precondition":{"kind":"none"},"digest":"sha256:30ca7a2a26049e727eee2f4a3b4021f9165741d3a9437f9304c69f4f280b5d21"}
{"schema":"cn-pi.effect-receipt.v0","id":"req-cn-pi-tsc-issues-get-incomplete-0005","request_digest":"sha256:30ca7a2a26049e727eee2f4a3b4021f9165741d3a9437f9304c69f4f280b5d21","status":"incomplete","reason":"pagination_limit","source":{"doc_id":"1ExampleTscOutboxDoc000000000","revision_id":"rev-outbox-20260807-05","record_ordinal":11},"ledger_seq":45,"policy":{"profile":"issues-and-docs-v0","digest":"sha256:4444444444444444444444444444444444444444444444444444444444444444"},"admission":{"canonical_request_sha256":"sha256:30ca7a2a26049e727eee2f4a3b4021f9165741d3a9437f9304c69f4f280b5d21","normalized_operation_sha256":"sha256:8888888888888888888888888888888888888888888888888888888888888888","precondition":{"kind":"none"}},"attempts":[{"number":1,"started_at":"2026-08-07T12:04:03Z","finished_at":"2026-08-07T12:04:04Z","outcome":"observed","backend_request_id":"gh-req-0005"}],"backend":{"kind":"github-rest","executor_identity":"host:cn-node-01","request_ids":["gh-req-0005"]},"observation":{"resource_kind":"issue","shape":"collection","complete":false,"truncated":true,"count":1,"next_cursor":"cursor-gh-issues-2","version":{"query_started_at":"2026-08-07T12:04:03Z","query_finished_at":"2026-08-07T12:04:04Z","items_sha256":"sha256:914aae9d2854619e1739f2216dfc5e57d5a114f2e4f1ca7e0ac291d1ca898f6f"},"data":[{"number":83,"title":"Earlier realization note","state":"closed","labels":[],"updated_at":"2026-08-07T11:59:00Z","url":"https://github.com/usurobor/tsc/issues/83"}]},"digest":"sha256:f3af24835de8d3a5aa2946595233dabddcdceeba5077228e1603d3ad3de74b8a"}
```

Issue-update and document-delete pairs prove the remaining successful mutation shapes:

```json
{"schema":"cn-pi.effect-request.v0","id":"req-cn-pi-tsc-issue-patch-0006","issuer":{"agent":"usurobor/cn-pi","locus":"usurobor/tsc"},"target":{"repository":"usurobor/tsc","profile":"issues-and-docs-v0"},"resource":{"kind":"issue","number":84},"method":"PATCH","representation":{"body":"Customer proof plan accepted."},"precondition":{"kind":"observed_issue","updated_at":"2026-08-07T12:00:04Z","representation_sha256":"sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"digest":"sha256:3cc6e0190d3f54900fd2caedc038b4c063a331b37af5651b20e726638661757b"}
{"schema":"cn-pi.effect-receipt.v0","id":"req-cn-pi-tsc-issue-patch-0006","request_digest":"sha256:3cc6e0190d3f54900fd2caedc038b4c063a331b37af5651b20e726638661757b","status":"succeeded","reason":"updated","source":{"doc_id":"1ExampleTscOutboxDoc000000000","revision_id":"rev-outbox-20260807-06","record_ordinal":12},"ledger_seq":46,"policy":{"profile":"issues-and-docs-v0","digest":"sha256:4444444444444444444444444444444444444444444444444444444444444444"},"admission":{"canonical_request_sha256":"sha256:3cc6e0190d3f54900fd2caedc038b4c063a331b37af5651b20e726638661757b","normalized_operation_sha256":"sha256:9999999999999999999999999999999999999999999999999999999999999999","precondition":{"kind":"observed_issue","updated_at":"2026-08-07T12:00:04Z","representation_sha256":"sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}},"attempts":[{"number":1,"started_at":"2026-08-07T12:05:03Z","finished_at":"2026-08-07T12:05:04Z","outcome":"committed","backend_request_id":"gh-req-0006"}],"backend":{"kind":"github-rest","executor_identity":"host:cn-node-01","request_ids":["gh-req-0006"]},"mutation":{"resource_kind":"issue","observed_before":{"number":84,"updated_at":"2026-08-07T12:00:04Z","representation_sha256":"sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"observed_after":{"number":84,"updated_at":"2026-08-07T12:05:04Z","representation_sha256":"sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"},"correlation":{"backend_request_id":"gh-req-0006"}},"digest":"sha256:2f8e6d65d7079a3b31193751ed1f8857c48e84079edc8043c14551e80738453a"}
{"schema":"cn-pi.effect-request.v0","id":"req-cn-pi-tsc-document-delete-0007","issuer":{"agent":"usurobor/cn-pi","locus":"usurobor/tsc"},"target":{"repository":"usurobor/tsc","profile":"issues-and-docs-v0"},"resource":{"kind":"document","path":"docs/customer-note.md","ref":"refs/heads/agent/pi-docs-tsc"},"method":"DELETE","precondition":{"kind":"document_match","ref":"refs/heads/agent/pi-docs-tsc","head_sha":"1111111111111111111111111111111111111111","blob_sha":"2222222222222222222222222222222222222222","content_sha256":"sha256:f47fe3326336ccea4614d8ecad937f6aed6034dba315f35475d4a6b96a40428a"},"digest":"sha256:226e579f72ed65db05b63ed1bce6321d4077538d8137419eda8f5a780469be8f"}
{"schema":"cn-pi.effect-receipt.v0","id":"req-cn-pi-tsc-document-delete-0007","request_digest":"sha256:226e579f72ed65db05b63ed1bce6321d4077538d8137419eda8f5a780469be8f","status":"succeeded","reason":"deleted","source":{"doc_id":"1ExampleTscOutboxDoc000000000","revision_id":"rev-outbox-20260807-07","record_ordinal":13},"ledger_seq":47,"policy":{"profile":"issues-and-docs-v0","digest":"sha256:4444444444444444444444444444444444444444444444444444444444444444"},"admission":{"canonical_request_sha256":"sha256:226e579f72ed65db05b63ed1bce6321d4077538d8137419eda8f5a780469be8f","normalized_operation_sha256":"sha256:1111111111111111111111111111111111111111111111111111111111111111","precondition":{"kind":"document_match","ref":"refs/heads/agent/pi-docs-tsc","head_sha":"1111111111111111111111111111111111111111","blob_sha":"2222222222222222222222222222222222222222","content_sha256":"sha256:f47fe3326336ccea4614d8ecad937f6aed6034dba315f35475d4a6b96a40428a"}},"attempts":[{"number":1,"started_at":"2026-08-07T12:06:03Z","finished_at":"2026-08-07T12:06:04Z","outcome":"committed","backend_request_id":"gh-req-0007"}],"backend":{"kind":"github-rest","executor_identity":"host:cn-node-01","request_ids":["gh-req-0007"]},"mutation":{"resource_kind":"document","observed_before":{"ref":"refs/heads/agent/pi-docs-tsc","head_sha":"1111111111111111111111111111111111111111","blob_sha":"2222222222222222222222222222222222222222","content_sha256":"sha256:f47fe3326336ccea4614d8ecad937f6aed6034dba315f35475d4a6b96a40428a"},"observed_after":{"ref":"refs/heads/agent/pi-docs-tsc","head_sha":"3333333333333333333333333333333333333333","path":"docs/customer-note.md","absent":true},"correlation":{"commit_sha":"3333333333333333333333333333333333333333"}},"digest":"sha256:67cb964506aea21b32fdecfe13da0e5753fc266a46d1157da98342a5aa9041b5"}
```

## Failure scenarios and conformance proof

| Scenario | Durable disposition | External-call rule |
|---|---|---|
| malformed/mismatched/nested frame or invalid JSON | `parse_incident` with Doc/revision/ordinal | zero |
| unterminated frame first observed | `pending_incomplete` with Doc/revision/ordinal/frame ID | zero |
| same frame completed in a later revision | validate once from retained discovery lineage | zero until admission |
| unchanged fragment exceeds configured bound | `parse_incident/incomplete_frame`; later completion needs operator reconciliation | zero |
| duplicate key, float, unknown field, or bad digest | `parse_incident` | zero |
| wrong Doc/owner/issuer/locus/target/receipt route | `trust_incident` | zero |
| permission metadata unavailable or ACL/fingerprint drift | `trust_incident`; projection pending | zero |
| policy/resource/method/path/media/budget denial | `policy_incident` or terminal `denied` | zero |
| exact replay | retained receipt/state | zero |
| same ID, different digest | terminal `denied/id_collision` | zero |
| stale issue/blob/head/state evidence | terminal `denied/stale_precondition` | read only |
| document resource/precondition ref mismatch | terminal `denied/stale_precondition` | zero mutation calls |
| default/main/tag/protected/foreign/unpinned document ref | terminal `denied/policy_denied` | zero |
| read exceeds pagination/result budget | terminal `incomplete` with partial inline data and continuation evidence | reads only |
| definite backend rejection before commit | terminal `failed/backend_rejected` | one attempt |
| create returns, exact post-read matches | terminal `succeeded` | one attempt |
| timeout after possible commit, marker found | terminal `succeeded/reconciled` | reads only after timeout |
| timeout after possible commit, outcome unprovable | `indeterminate` | no automatic retry |
| post-read mismatch | `indeterminate/remote_mismatch` | no mutation retry |
| ledger unavailable before mutation | delivery error/incident | zero |
| receipt Doc append or verification fails | terminal outcome + `receipt_projection_pending` | receipt retry only |
| restart with `executing` | reconcile to success, proven absence, or indeterminate | never blind replay |

AC proof map: AC1 is the header/status and one-file product diff plus exact-head PR review evidence; AC2 is paired Docs, ACL binding, framing, and durable incomplete/incidents; AC3 is recursive closure and recomputed fixture digests; AC4 is the closed 14-row matrix plus explicit PR-review exclusion; AC5 is the backend/ref table and indeterminate rule; AC6 is retained admitted bytes, ledger-before-effect, replay, collision, typed perception, and projection isolation; AC7 is the proposed layout and evidence gate below.

Conformance must strictly parse every positive fixture, recompute every request/receipt/content digest, prove the negative nested-key fixture is denied before backend calls, count exactly 14 matrix rows, exercise every failure/status row, verify all links/paths, scan negative space, and run `git diff --check origin/main...HEAD`. Exact-head beta evidence is published on PR #3; no approval claim may precede a zero-finding review at that head.

## Proposed source, deployment, state, and retirement layout

Later implementation may reconcile the experimental branch, but must not treat it as landed or canonical. Proposed repository source:

```text
ops/drive-ingress/EFFECT-OPS-v0.md
ops/drive-ingress/cn-pi-drive-ingress
ops/drive-ingress/test_drive_ingress.py
ops/drive-ingress/test_effect_ops.py
ops/drive-ingress/systemd/cn-pi-drive-ingress.service
ops/drive-ingress/systemd/cn-pi-drive-ingress.timer
```

Installed artifacts and mutable state remain derived/host-local:

```text
/usr/local/libexec/cn-pi-drive-ingress
/usr/local/sbin/cn-pi-drive-ingress-check
/etc/systemd/system/cn-pi-drive-ingress.{service,timer}
/var/lib/cn-pi-drive-ingress/effects/ledger.sqlite3
/var/lib/cn-pi-drive-ingress/effects/receipt-projections/
```

A later service update must add only the state write paths it needs and preserve current sandboxing. Cutover requires tests, backup/migration of mutable state, dry-run discovery, reconciliation of in-flight records, and rollback to the dialogue/memory-only bridge. Retire effect execution when Pi has a reviewed safer path or CNOS meets section 11; retain/migrate the ledger and receipts rather than deleting evidence.

## Evidence gate for later CNOS promotion

The bridge is an evidence generator, not a second canonical platform. CNOS may reconsider this boundary only when it has a real Go intake/execution consumer and retained bridge fixtures/receipts for exact replay, collision, policy denial, ACL drift/unavailable metadata, stale document evidence, the issue-update interval race, post-create crash, reconciliation, remote mismatch, and receipt-projection failure.

Promotion compares observed invariants with then-current CNOS architecture. It may keep the stable ID/digest split, policy non-widening, backend-honest evidence, ledger-before-effect, visible uncertainty, and receipt projection. It need not preserve Drive framing or experimental syntax.

Process economics: the paired Docs cost eight one-time operator-created files and ACL maintenance; they prevent unauthenticated source writes and give Pi receipts. The ledger adds host state but prevents duplicate effects and silent loss. Framing and the closed matrix are mechanical checks. No DAG, provider ABI, or generic URI/canonicalization layer earns present cost. Review this experiment after the first complete failure corpus; automate repeated permission/digest checks, and retire the Drive process when a stronger authenticated boundary replaces it.

## Non-goals

- No implementation, schema file, test, service, timer, Drive object, ledger, credential, runtime ref, or repository-content mutation beyond this design. Issue/PR metadata may record dispatch and exact-head review evidence only.
- No dialogue/memory parsing as effects; no change to r0, activation, dialogue, memory, or state-ref grammar.
- No direct Pi execution or credentials; no raw shell, HTTP, Git, workflow, settings, or policy operation.
- No full pull-request observation, review-thread/check aggregation, merge-gating verdict, PR creation, or merge; those require a separately reviewed resource family.
- No multi-operation plan, DAG, dependency scheduler, aggregate plan status, CTB compiler, provider host, hypermedia, cache, capability registry, canonical `cn:` URI, or generic canonicalization subsystem.
- No atomic-CAS claim for GitHub issues/comments and no silent terminalization of uncertainty.
- No canonical CNOS protocol until the evidence and runtime gates pass.
