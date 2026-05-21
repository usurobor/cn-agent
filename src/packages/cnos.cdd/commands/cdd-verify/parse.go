package cddverify

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

// Receipt is a parsed receipt at the YAML/JSON boundary.
//
// Per eng/go §2.3, map[string]any is permitted at true JSON boundaries
// (which a receipt is — its schema is governed by CUE, not Go). Counterfeit
// rules access fields through small typed accessor helpers in
// counterfeit.go so the rest of the code is concrete-typed.
type Receipt map[string]any

// ParseReceiptJSON parses receipt JSON bytes into a Receipt.
//
// Pure (eng/go §2.17): no IO, no os imports inside the function body
// beyond encoding/json. Used by tests that supply JSON directly and by
// ReadReceipt after cue-export coerces YAML → JSON.
func ParseReceiptJSON(data []byte) (Receipt, error) {
	var r Receipt
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, fmt.Errorf("decode receipt json: %w", err)
	}
	return r, nil
}

// ReadReceipt loads a receipt from a YAML or JSON file on disk.
//
// IO wrapper (eng/go §2.17) — calls ParseReceiptJSON after coercing the
// file via `cue export --out json` (cue is already a pinned runtime
// dependency; this avoids adding a YAML parser).
//
// JSON files (.json) are read directly; everything else is treated as
// YAML and routed through cue.
func ReadReceipt(ctx context.Context, path string) (Receipt, error) {
	if len(path) > 5 && path[len(path)-5:] == ".json" {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read receipt %q: %w", path, err)
		}
		return ParseReceiptJSON(data)
	}
	// YAML path — coerce to JSON via cue export.
	cmd := exec.CommandContext(ctx, "cue", "export", "--out", "json", path)
	out, err := cmd.Output()
	if err != nil {
		var stderr string
		if ee, ok := err.(*exec.ExitError); ok {
			stderr = string(ee.Stderr)
		}
		return nil, fmt.Errorf("cue export %q: %w (%s)", path, err, stderr)
	}
	return ParseReceiptJSON(out)
}

// asMap returns r[key] coerced to map[string]any, or nil + ok=false.
//
// Defensive accessor — receipts come from untrusted YAML; type assertions
// must never panic. eng/go §1.3 "panic-driven control flow" → never.
func asMap(r Receipt, key string) (map[string]any, bool) {
	v, ok := r[key]
	if !ok || v == nil {
		return nil, false
	}
	m, ok := v.(map[string]any)
	return m, ok
}

// asString returns m[key] coerced to string, or "" + ok=false.
func asString(m map[string]any, key string) (string, bool) {
	v, ok := m[key]
	if !ok || v == nil {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

// asStringOrEmpty returns m[key] as string, "" if absent or wrong type.
func asStringOrEmpty(m map[string]any, key string) string {
	s, _ := asString(m, key)
	return s
}

// asSlice returns m[key] coerced to []any, or nil + ok=false.
func asSlice(m map[string]any, key string) ([]any, bool) {
	v, ok := m[key]
	if !ok || v == nil {
		return nil, false
	}
	s, ok := v.([]any)
	return s, ok
}

// protocolID extracts the receipt's top-level protocol_id.
func protocolID(r Receipt) string {
	return asStringOrEmpty(map[string]any(r), "protocol_id")
}

// evidenceRefs extracts the receipt's evidence_refs map.
func evidenceRefs(r Receipt) (map[string]any, bool) {
	return asMap(r, "evidence_refs")
}

// boundary extracts the receipt's boundary_decision map.
func boundary(r Receipt) (map[string]any, bool) {
	return asMap(r, "boundary_decision")
}

// validationBlock extracts the receipt's validation map.
func validationBlock(r Receipt) (map[string]any, bool) {
	return asMap(r, "validation")
}

// mode returns the receipt's top-level mode field, or "".
func mode(r Receipt) string {
	return asStringOrEmpty(map[string]any(r), "mode")
}
