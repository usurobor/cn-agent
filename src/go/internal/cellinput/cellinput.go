// Package cellinput decodes the per-run input document: the issue, the design and
// the subject one invocation acts on. It is a SEPARATE document from the cell
// spec, because a spec is reusable while an issue, a design and a revision belong
// to one run — folding them in would make every run author a new cell. The three
// payloads stay OPAQUE: this package owns the envelope and its closed key
// language, and never learns what an issue is. No IO, no paths (eng/go §2.17).
package cellinput

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/usurobor/cnos/src/go/internal/cellfill"
)

// Kind is the pinned run-input schema tag: a document must declare it exactly, and
// the CUE definition #CDSRunInput pins the same string. The spelling follows the
// repo's `cnos.<domain>.<name>.vN` convention (CDS-CELL-MIGRATION.md).
const Kind = "cnos.cds.run-input.v0"

// RunInput is the decoded envelope. The three members are raw because their
// meaning is not this package's: the bytes admitted are the bytes frozen.
type RunInput struct {
	Kind    string          `json:"kind"`
	Issue   json.RawMessage `json:"issue"`
	Design  json.RawMessage `json:"design"`
	Subject json.RawMessage `json:"subject"`
}

// Digest is the identity of one untrusted document: the lowercase hex SHA-256 of
// the EXACT bytes handed in, never a re-serialization — that would digest what
// this package understood, and it understands none of the payloads. Separate from
// Decode, because a document that FAILS to decode still has an identity, and that
// is exactly the document the receipt for a wrong `kind` has to name.
func Digest(raw []byte) string {
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
}

// Decode decodes one run-input document. Presence of the three payloads is
// deliberately NOT checked here: an absent issue and a malformed issue are
// different admission outcomes, and deciding which is admission's job.
func Decode(raw []byte) (RunInput, error) {
	if len(raw) == 0 {
		return RunInput{}, fmt.Errorf("cell run input: document is empty")
	}
	// Duplicate keys and nulls are closed here for the reason they are closed on the
	// cell spec: encoding/json takes the last of a repeated key while CUE rejects the
	// document, so the two authorities reading these exact bytes would disagree.
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
