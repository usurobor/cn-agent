package cellfill

import (
	"encoding/json"
	"strings"
	"testing"
)

// StrictDecode's contract is that a declaration ends where it says it ends.
// The outer contract parser rejects malformed input before production ever
// reaches this helper, so without a direct test the EOF branch is unreached
// and the claim is unverified (Pi #56 C2). These call the helper itself.
func TestStrictDecodeRequiresRealEOF(t *testing.T) {
	type decl struct {
		Fill string `json:"fill"`
	}

	cases := []struct {
		name string
		raw  string
		ok   bool
		want string // substring of the expected error
	}{
		{name: "clean", raw: `{"fill":"x"}`, ok: true},
		{name: "trailing value", raw: `{"fill":"x"} {"fill":"y"}`, want: "trailing data"},
		{name: "trailing garbage", raw: `{"fill":"x"} <not json>`, want: "malformed data"},
		{name: "truncated trailer", raw: `{"fill":"x"} {"fill":`, want: "malformed data"},
		{name: "unknown field", raw: `{"fill":"x","extra":1}`, want: "extra"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var into decl
			err := StrictDecode(json.RawMessage(tc.raw), &into)
			if tc.ok {
				if err != nil {
					t.Fatalf("want clean decode, got %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("want an error containing %q, got nil", tc.want)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error %q does not mention %q", err, tc.want)
			}
		})
	}
}
