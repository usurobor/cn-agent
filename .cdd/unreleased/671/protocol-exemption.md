# Protocol-exemption binding (envelope) — #671

> **Purpose.** Bind the mutable GitHub issue-body `## Protocol exemption`
> section into the receipt head so constitution cannot silently depend on prose
> that changes after review. This is the R17 repair of external-β finding
> [REQUIRED] "The revision-bound gate does not bind the constitution input"
> (comment 5076629728).
>
> The **verbatim exemption payload** is the separate file
> `protocol-exemption-source.md` (this envelope carries only binding metadata).

## Bound exemption identity

| Field | Value |
|---|---|
| Source | GitHub issue [#671](https://github.com/usurobor/cnos/issues/671) body, section `## Protocol exemption` (the last section of the body) |
| Observed `updated_at` | **`2026-07-25T03:28:34Z`** |
| Snapshot payload file | `.cdd/unreleased/671/protocol-exemption-source.md` |
| **Snapshot SHA-256** | **`dccba69c668163b09e00ef79a77f7e6236e39cc048eca96c50fca343b507d473`** |
| Snapshot length | 3,862 bytes (LF line endings, no trailing newline) |

## Normalization (for the live-vs-snapshot comparison)

The GitHub REST representation HTML-entity-encodes some characters in the body
string (`&#39;` for `'`, `&#34;` for `"`, `&gt;` for `>`, `&amp;` for `&`). The
snapshot payload stores the **decoded markdown source** (literal `'`, `"`, `>`,
`&`) as authored. The gate comparison is therefore defined as:

1. Fetch issue #671 body.
2. Extract the `## Protocol exemption` section (from the `## Protocol exemption`
   header through end-of-body).
3. HTML-entity-decode (`&#39;`→`'`, `&#34;`→`"`, `&gt;`→`>`, `&lt;`→`<`,
   `&amp;`→`&`).
4. Strip any trailing newline.
5. Require the result to equal `protocol-exemption-source.md` byte-for-byte
   (equivalently: SHA-256 == `dccba69c…`).

## Gate rule (pre-authorization freshness — binding)

The pre-authorization gate **FAILS STALE** if, at authorization time, either:

- the live issue #671 `## Protocol exemption` section (normalized as above) does
  **not** match the snapshot SHA-256 `dccba69c…`; or
- the live issue `updated_at` is **later** than `2026-07-25T03:28:34Z` **and**
  the normalized section no longer matches the snapshot.

On stale, authorization is refused and the exemption must be re-reviewed and
re-snapshotted (a new R-round), exactly as a changed matter tree re-enters the
gate. This makes the exemption a **revision-bound authorization input**: an
operator can prove the exemption authorized is the one reviewed, and mutating
the issue body after review cannot slip past the fixed authorization SHA.

**Verification (matches the reviewer's stated check).** Mutating the issue
exemption while leaving the Git head unchanged makes step 5 fail → gate fails
stale. The PR body, this envelope, and the γ closeout name the same
`matter_sha` (`614829a4`), `receipt_head`, exemption snapshot hash (`dccba69c…`)
+ observed `updated_at`, and next stage.
