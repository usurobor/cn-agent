package cellfill

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/usurobor/cnos/src/go/internal/cellcog"
	"github.com/usurobor/cnos/src/go/internal/cellmethod"
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

// canonical() is the sibling JSON boundary to StrictDecode, and it must make
// the same demand: a fill's declaration ends where it says it ends. Decoding
// one value proves a value arrived, not that only one did — without an EOF
// check the trailing bytes are silently dropped from the record (Pi #57 B1).
func TestCanonicalRequiresRealEOF(t *testing.T) {
	cases := []struct {
		name string
		raw  string
		want string // "" means the call must succeed
	}{
		{name: "clean object", raw: `{"fill":"x","n":1}`},
		{name: "trailing value", raw: `{"fill":"x"} {"fill":"y"}`, want: "trailing data"},
		{name: "trailing garbage", raw: `{"fill":"x"} garbage`, want: "malformed data"},
		{name: "not an object", raw: `["x"]`, want: "not an object"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := canonical("t.fill", json.RawMessage(tc.raw))
			if tc.want == "" {
				if err != nil {
					t.Fatalf("want success, got %v", err)
				}
				if len(got) == 0 {
					t.Fatal("want canonical bytes, got none")
				}
				return
			}
			if err == nil {
				t.Fatalf("want an error containing %q, got %s", tc.want, got)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error %q does not mention %q", err, tc.want)
			}
		})
	}
}

// AdmitSeatDecl is the ONE decoder both CDS seats now go through, so the rules
// that used to be copied per fill are proven once, here, against the helper
// itself. The two fills' own tests still assert their own refusals end to end;
// what this covers is the shape rules no fill states any more.
//
// The projection is built directly rather than loaded: Empty() is a digest
// test, so a bundle on disk would add IO to a decode test that needs none.
func heldBy(role cellmethod.Role) cellmethod.View {
	return cellmethod.View{Role: role, Text: "(projection)", SHA256: strings.Repeat("a", 64)}
}

var testRefusal = SeatRefusal{
	NoMethodology: "a test seat needs the cell's methodology, and this cell declares none",
	WrongRole:     "a test seat takes the constructive projection",
}

// The key language is CLOSED, at both object shapes and case-sensitively: a
// workspace, a skills list, a second cognition argument, or any of them spelled
// with a capital is refused rather than silently decoded. Each of these is a
// second source of something the cell declares once, which is why the negative
// space is asserted key by key.
func TestAdmitSeatDeclClosesItsKeyLanguage(t *testing.T) {
	const cog = `"cognition":{"provider":"fake","model":"m"}`
	for name, raw := range map[string]string{
		"unknown key":             `{"fill":"t.fill",` + cog + `,"extra":1}`,
		"fill, mixed case":        `{"Fill":"t.fill",` + cog + `}`,
		"cognition, mixed case":   `{"fill":"t.fill","Cognition":{"provider":"fake","model":"m"}}`,
		"a workspace":             `{"fill":"t.fill",` + cog + `,"workspace":{"repo":"."}}`,
		"skills of its own":       `{"fill":"t.fill",` + cog + `,"skills":["cnos.eng:eng/go"]}`,
		"an argv inside":          `{"fill":"t.fill","cognition":{"provider":"fake","model":"m","argv":["--x"]}}`,
		"provider, mixed case":    `{"fill":"t.fill","cognition":{"Provider":"fake","model":"m"}}`,
		"cognition not an object": `{"fill":"t.fill","cognition":["fake"]}`,
		"trailing declaration":    `{"fill":"t.fill",` + cog + `} {"fill":"t.fill"}`,
	} {
		t.Run(name, func(t *testing.T) {
			_, err := AdmitSeatDecl(json.RawMessage(raw), "t.fill",
				cellmethod.RoleConstructive, heldBy(cellmethod.RoleConstructive), testRefusal)
			if err == nil {
				t.Fatalf("the closed declaration admitted %s", raw)
			}
			// Every refusal names the fill it came from, because a fill's caller
			// returns this error unchanged.
			if !strings.Contains(err.Error(), `fill "t.fill"`) {
				t.Fatalf("refusal does not name the fill: %v", err)
			}
		})
	}
}

// The admitted shape decodes to exactly what the fill then rents cognition
// with, and a declaration with no cognition at all is a decode that leaves the
// zero config for the port to refuse — not an error this helper invents.
func TestAdmitSeatDeclDecodesTheAdmittedShape(t *testing.T) {
	view := heldBy(cellmethod.RoleConstructive)
	d, err := AdmitSeatDecl(json.RawMessage(`{"fill":"t.fill","cognition":{"provider":"fake","model":"m"}}`),
		"t.fill", cellmethod.RoleConstructive, view, testRefusal)
	if err != nil {
		t.Fatalf("the admitted shape was refused: %v", err)
	}
	if d.Fill != "t.fill" || d.Cognition.Provider != "fake" || d.Cognition.Model != "m" {
		t.Fatalf("decoded declaration is %+v", d)
	}
	bare, err := AdmitSeatDecl(json.RawMessage(`{"fill":"t.fill"}`), "t.fill",
		cellmethod.RoleConstructive, view, testRefusal)
	if err != nil {
		t.Fatalf("a declaration without cognition must decode, not fail here: %v", err)
	}
	if bare.Cognition != (cellcog.Config{}) {
		t.Fatalf("an absent cognition decoded to %+v", bare.Cognition)
	}
}

// The projection check, in both directions and in the fill's own words: no
// methodology at all and the OTHER role's projection are distinct refusals, and
// each carries the sentence the calling fill supplied plus, for the role, what
// it was actually handed.
func TestAdmitSeatDeclChecksTheProjectionRole(t *testing.T) {
	const decl = `{"fill":"t.fill","cognition":{"provider":"fake","model":"m"}}`
	for name, tc := range map[string]struct {
		want cellmethod.Role
		view cellmethod.View
		says string
	}{
		"no methodology at all": {cellmethod.RoleConstructive, cellmethod.View{}, testRefusal.NoMethodology},
		"the adversarial projection to a producing seat": {
			cellmethod.RoleConstructive, heldBy(cellmethod.RoleAdversarial),
			testRefusal.WrongRole + `, got "adversarial"`},
		"the constructive projection to an assessing seat": {
			cellmethod.RoleAdversarial, heldBy(cellmethod.RoleConstructive),
			testRefusal.WrongRole + `, got "constructive"`},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := AdmitSeatDecl(json.RawMessage(decl), "t.fill", tc.want, tc.view, testRefusal)
			if err == nil {
				t.Fatal("the wrong projection constructed a seat")
			}
			if err.Error() != `fill "t.fill": `+tc.says {
				t.Fatalf("refusal is %q, want %q", err, `fill "t.fill": `+tc.says)
			}
		})
	}
	// ...and the projection it does take is admitted, on both sides.
	for _, role := range []cellmethod.Role{cellmethod.RoleConstructive, cellmethod.RoleAdversarial} {
		if _, err := AdmitSeatDecl(json.RawMessage(decl), "t.fill", role, heldBy(role), testRefusal); err != nil {
			t.Fatalf("the %s projection was refused by a seat that takes it: %v", role, err)
		}
	}
}
