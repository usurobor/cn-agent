package cellinput

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"strings"
	"testing"
)

const validPath = "../../../../schemas/cds/fixtures/runinput/valid-run-input.json"

// The digest is over the EXACT bytes handed in. This is the assertion that it
// is not a digest of a re-serialization: the fixture is pretty-printed JSON,
// so any normalizing pass would change the bytes and therefore the digest.
func TestDigestIsOverTheExactBytes(t *testing.T) {
	raw, err := os.ReadFile(validPath)
	if err != nil {
		t.Fatal(err)
	}
	digest := Digest(raw)
	sum := sha256.Sum256(raw)
	if want := hex.EncodeToString(sum[:]); digest != want {
		t.Fatalf("digest = %q, want the sha256 of the file %q", digest, want)
	}

	// One byte of insignificant whitespace moves it — so the digest binds the
	// document as authored, not as understood.
	spaced := append(append([]byte(nil), raw...), ' ')
	if _, err := Decode(spaced); err != nil {
		t.Fatalf("a trailing space must still decode: %v", err)
	}
	if Digest(spaced) == digest {
		t.Fatal("a byte the decoder ignores left the digest unmoved")
	}
}

// A document that FAILS to decode still has an identity. This is why Digest is
// not a second return value of Decode: the receipt for a refused document is
// the only record that the document was ever seen, and it has to name which
// document that was.
func TestAnUndecodableDocumentStillHasADigest(t *testing.T) {
	const garbage = `{"kind":`
	if _, err := Decode([]byte(garbage)); err == nil {
		t.Fatal("the fixture must be undecodable, or this proves nothing")
	}
	sum := sha256.Sum256([]byte(garbage))
	if got, want := Digest([]byte(garbage)), hex.EncodeToString(sum[:]); got != want {
		t.Fatalf("digest of an undecodable document = %q, want %q", got, want)
	}
}

// The envelope is closed and the payloads are not decoded here. Both halves
// matter: the first is what keeps Go and CUE reading one language, the second
// is what keeps this package from learning what an issue is.
func TestEnvelopeIsClosedAndPayloadsStayRaw(t *testing.T) {
	raw, err := os.ReadFile(validPath)
	if err != nil {
		t.Fatal(err)
	}
	in, err := Decode(raw)
	if err != nil {
		t.Fatal(err)
	}
	for name, payload := range map[string][]byte{
		"issue": in.Issue, "design": in.Design, "subject": in.Subject,
	} {
		if len(payload) == 0 {
			t.Fatalf("%s payload is empty", name)
		}
		if payload[0] != '{' {
			t.Fatalf("%s payload was normalized rather than carried raw: %s", name, payload)
		}
	}
}

func TestRejects(t *testing.T) {
	valid, err := os.ReadFile(validPath)
	if err != nil {
		t.Fatal(err)
	}
	cases := map[string]struct {
		raw  string
		want string
	}{
		"empty":         {"", "document is empty"},
		"not an object": {`[]`, "not an object"},
		"malformed":     {`{"kind":`, "EOF"},
		"wrong kind":    {`{"kind":"cnos.cell.run-input/9.9"}`, "kind must be"},
		"unknown key":   {`{"kind":"` + Kind + `","task":{}}`, `unknown key "task"`},
		// encoding/json matches field names case-insensitively even with
		// DisallowUnknownFields, so `Kind` would decode here while the closed
		// #CDSRunInput rejects it.
		"mixed case key": {`{"Kind":"` + Kind + `"}`, `unknown key "Kind"`},
		// Last-wins in Go, a hard error in CUE: two authorities reading one
		// document differently is the exact divergence this closes.
		"duplicate key": {`{"kind":"` + Kind + `","kind":"other"}`, `duplicate key "kind"`},
		"null payload":  {`{"kind":"` + Kind + `","issue":null}`, "null is not allowed"},
		"trailing data": {string(valid) + "}", `invalid character '}' after top-level value`},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := Decode([]byte(tc.raw))
			if err == nil {
				t.Fatal("must be rejected")
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("rejected for the wrong reason:\n got: %v\nwant mention of: %q", err, tc.want)
			}
		})
	}
}

// Presence of the payloads is admission's decision, not the decoder's: an
// absent issue and a malformed one are different outcomes, and only the door
// can say which. A decoder that required them would collapse the two.
func TestPayloadPresenceIsNotDecided(t *testing.T) {
	in, err := Decode([]byte(`{"kind":"` + Kind + `"}`))
	if err != nil {
		t.Fatalf("an envelope with no payloads must still decode: %v", err)
	}
	if in.Issue != nil || in.Design != nil || in.Subject != nil {
		t.Fatalf("absent payloads must stay absent: %+v", in)
	}
}
