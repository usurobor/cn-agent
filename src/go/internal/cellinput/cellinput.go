// Package cellinput decodes the per-run input document: the issue, the design,
// and the subject reference one invocation acts on.
//
// It is a SEPARATE document from the cell spec, and that separation is the
// point. A cell spec is reusable — the same seats, the same fills, the same
// methodology — while an issue, a design and a repository revision belong to
// one run. Folding them into the spec would make every run author a new cell,
// and would put per-run content inside the artifact the plan is compiled from.
//
// The three payloads stay OPAQUE here: this package owns the envelope — its
// declared kind and its closed key language — and never learns what an issue
// or a subject is. That belongs to cdsadmit and to the adapters it calls.
//
// Pure: no IO, no paths (eng/go §2.17). Decode takes bytes; the caller reads
// the file.
package cellinput

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/usurobor/cnos/src/go/internal/cellfill"
)

// Kind is the pinned run-input schema tag. A document must declare it exactly;
// the CUE definition #CDSRunInput pins the same string.
//
// The spelling is the repo's `cnos.<domain>.<name>.vN` convention, the one
// cnos.cds.issue.v0, cnos.cellspec.v0 and cnos.cellkernel.episode-closure.v0
// already use. The `/0.1` form belongs to the non-`cnos.` adapter namespace
// (git.snapshot/0.1) and is not this document's; two spellings inside one
// namespace would leave a reader guessing which one a new tag should follow.
const Kind = "cnos.cds.run-input.v0"

// RunInput is the decoded envelope. The three members are raw because their
// meaning is not this package's: they are handed to the admission gate exactly
// as authored, and the bytes that are admitted are the bytes that are frozen.
type RunInput struct {
	Kind    string          `json:"kind"`
	Issue   json.RawMessage `json:"issue"`
	Design  json.RawMessage `json:"design"`
	Subject json.RawMessage `json:"subject"`
}

// Digest is the identity of one untrusted run-input document: the lowercase
// hex SHA-256 of the EXACT bytes handed in, never a re-serialization of a
// decoded value. The distinction is load-bearing — a digest over re-serialized
// bytes would be a digest of what this package understood, and the whole reason
// the payloads stay raw is that it understands none of them.
//
// It is a separate function from Decode, and that is the point: a document that
// FAILS to decode still has an identity, and that is exactly the document a
// refusal has to name. A digest returned only alongside a successful decode
// could not be carried by the receipt for a wrong `kind`.
func Digest(raw []byte) string {
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
}

// Decode decodes one run-input document.
//
// Presence of the three payloads is deliberately NOT checked here. An absent
// issue and a malformed issue are different admission outcomes (incomplete vs
// rejected), and deciding which is admission's job, not the decoder's.
func Decode(raw []byte) (RunInput, error) {
	if len(raw) == 0 {
		return RunInput{}, fmt.Errorf("cell run input: document is empty")
	}
	// Duplicate keys and nulls are closed here for the reason they are closed
	// on the cell spec: encoding/json silently takes the last of a repeated
	// key while CUE rejects the document outright, so `{"kind":"a","kind":"b"}`
	// would be read differently by the two authorities that both read these
	// exact bytes.
	if err := cellfill.NoDuplicateKeysOrNull(raw); err != nil {
		return RunInput{}, fmt.Errorf("cell run input: %w", err)
	}
	if err := cellfill.OnlyKeys(raw, "cell run input", "kind", "issue", "design", "subject"); err != nil {
		return RunInput{}, fmt.Errorf("cell run input: %w", err)
	}
	var in RunInput
	if err := cellfill.StrictDecode(raw, &in); err != nil {
		return RunInput{}, fmt.Errorf("cell run input: %w", err)
	}
	if in.Kind != Kind {
		return RunInput{}, fmt.Errorf("cell run input: kind must be %q, got %q", Kind, in.Kind)
	}
	return in, nil
}
