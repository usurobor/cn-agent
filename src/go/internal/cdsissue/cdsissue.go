// Package cdsissue owns the CDS typed issue: its shape, its admission
// predicate, and its renderer. Nothing else.
//
// The kernel and cellspec carry the issue as opaque bytes on `contract.issue` —
// neither learns what a CDS issue is. This package is where the meaning lives,
// and it is the only place it lives.
//
// Both seats are told the same thing because there is ONE Render, not because
// two formatters are kept in step: the producing fill and the assessing fill
// call it on the same frozen bytes. Admit has one caller, the admission door,
// which is what makes the issue those bytes rather than a document either seat
// interpreted for itself.
//
// The issue enters as JSON, already structured. There is deliberately no
// Markdown parser: prose that has to be parsed is prose two authorities can
// read differently. And no `$param` hole is substituted anywhere inside it —
// cellspec splices seat declarations only, so an issue is authored literally
// and frozen literally into the record and its digest.
//
// Pure: no IO, no paths. Admit takes bytes.
package cdsissue

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/usurobor/cnos/src/go/internal/cellfill"
)

// nonBlankPattern matches one rune that is not whitespace, and it is the ONE
// blankness predicate this schema has. `!= ""` is not it: "   " satisfies that
// while reading as an authored sentence to every human and to neither
// authority.
//
// It is a REGEXP rather than `strings.TrimSpace(v) != ""` so that CUE's
// #NonBlank can be the same RE2 program rather than a second hand-written
// predicate that merely looks equivalent. They are not equivalent: Go's
// `\s` (and CUE's) is the ASCII class `[\t\n\f\r ]`, while `TrimSpace` strips
// everything `unicode.IsSpace` covers — so `=~"\\S"` beside `TrimSpace`
// accepted a NO-BREAK-SPACE-only field in CUE and rejected it in Go. The class
// below enumerates `unicode.IsSpace` exactly, and the corpus fixture
// issue-blank-unicode-whitespace.json holds the whole set in one field, so an
// omission here is a gate failure rather than a divergence found in
// production.
var nonBlankPattern = `[^\t\n\v\f\r \x{0085}\x{00A0}\x{1680}\x{2000}-\x{200A}\x{2028}\x{2029}\x{202F}\x{205F}\x{3000}]`

var nonBlank = regexp.MustCompile(nonBlankPattern)

// ReservedIDPrefix is the acceptance-id namespace an issue may not use, because
// the runtime mints unit ids under it. Declared here rather than in the
// assessing fill: the door is what admits the document, so the door is what has
// to know which ids are not the author's to choose.
const ReservedIDPrefix = "check:"

// Kind is the pinned issue schema tag. An issue must declare it exactly; the
// CUE definition #CDSIssue pins the same string.
const Kind = "cnos.cds.issue.v0"

// Issue is the complete structured task a CDS cell carries. Every field is
// required, because each one answers a question β would otherwise have to
// invent an answer to.
type Issue struct {
	Kind       string      `json:"kind"`
	ID         string      `json:"id"`
	Problem    Problem     `json:"problem"`
	Sources    []Source    `json:"sources"`
	Scope      Scope       `json:"scope"`
	Acceptance []Criterion `json:"acceptance"`
}

// Problem is the incoherence in three lines: what exists, what is expected,
// where they diverge.
type Problem struct {
	Exists   string `json:"exists"`
	Expected string `json:"expected"`
	Diverges string `json:"diverges"`
}

// Source binds one load-bearing claim to one canonical path.
type Source struct {
	Claim string `json:"claim"`
	Path  string `json:"path"`
}

// Scope is the execution boundary. Out may be empty, but it must be PRESENT:
// non-goals are load-bearing, so the author has to have considered them.
type Scope struct {
	In  []string `json:"in"`
	Out []string `json:"out"`
}

// Criterion is one acceptance criterion. Verification is required because it
// is the "proof or rejection mechanism" cdd/issue names, and it is what makes
// β's job well defined: an issue whose criteria state no verification route
// leaves β judging plausibility instead of deciding against a criterion.
type Criterion struct {
	ID           string `json:"id"`
	Statement    string `json:"statement"`
	Verification string `json:"verification"`
}

// Blank reports whether s carries nothing but whitespace. Exported because the
// parity between this predicate and CUE's #NonBlank is itself under test.
func Blank(s string) bool { return !nonBlank.MatchString(s) }

// Admit decodes and admits an issue document. Every rule is mechanical: a
// human judgement about whether an issue is "good enough" is exactly what this
// cycle exists to stop relying on.
//
// The closed key language is stated with cellfill.OnlyKeys at every level, not
// left to the decoder: encoding/json matches field names case-insensitively
// even with DisallowUnknownFields, so `Kind` or a nested `Path` would decode
// in Go while the closed #CDSIssue rejects it — the two authorities would
// admit different documents.
func Admit(raw []byte) (Issue, error) {
	if len(raw) == 0 {
		return Issue{}, fmt.Errorf("cdsissue: no issue bytes to admit")
	}
	if err := exactShape(raw); err != nil {
		return Issue{}, fmt.Errorf("cdsissue: %w", err)
	}

	var iss Issue
	if err := cellfill.StrictDecode(raw, &iss); err != nil {
		return Issue{}, fmt.Errorf("cdsissue: %w", err)
	}

	if iss.Kind != Kind {
		return Issue{}, fmt.Errorf("cdsissue: kind must be %q, got %q", Kind, iss.Kind)
	}
	// Blank, not empty, and the SAME rule everywhere. An earlier revision
	// trimmed the problem triple and tested `!= ""` on the rest, which made a
	// whitespace-only acceptance criterion legal — precisely the ill-defined
	// issue this gate exists to reject. One predicate is also one thing for
	// #NonBlank to mirror.
	if Blank(iss.ID) {
		return Issue{}, fmt.Errorf("cdsissue: id is required")
	}
	for _, f := range []struct{ name, value string }{
		{"exists", iss.Problem.Exists},
		{"expected", iss.Problem.Expected},
		{"diverges", iss.Problem.Diverges},
	} {
		if Blank(f.value) {
			return Issue{}, fmt.Errorf("cdsissue: problem.%s is required", f.name)
		}
	}
	if len(iss.Sources) == 0 {
		return Issue{}, fmt.Errorf("cdsissue: sources is required (one canonical path per load-bearing claim)")
	}
	for i, s := range iss.Sources {
		if Blank(s.Claim) || Blank(s.Path) {
			return Issue{}, fmt.Errorf("cdsissue: sources[%d] needs a claim and a path", i)
		}
	}
	if len(iss.Scope.In) == 0 {
		return Issue{}, fmt.Errorf("cdsissue: scope.in is required")
	}
	for _, s := range [][]string{iss.Scope.In, iss.Scope.Out} {
		for i, item := range s {
			if Blank(item) {
				return Issue{}, fmt.Errorf("cdsissue: scope entry %d is blank", i)
			}
		}
	}
	if len(iss.Acceptance) == 0 {
		return Issue{}, fmt.Errorf("cdsissue: acceptance is required")
	}
	seen := make(map[string]bool, len(iss.Acceptance))
	for i, c := range iss.Acceptance {
		switch {
		case Blank(c.ID):
			return Issue{}, fmt.Errorf("cdsissue: acceptance[%d] has no id", i)
		case Blank(c.Statement):
			return Issue{}, fmt.Errorf("cdsissue: acceptance %q has no statement", c.ID)
		case Blank(c.Verification):
			return Issue{}, fmt.Errorf("cdsissue: acceptance %q states no verification route", c.ID)
		case seen[c.ID]:
			return Issue{}, fmt.Errorf("cdsissue: duplicate acceptance id %q", c.ID)
		case strings.HasPrefix(c.ID, ReservedIDPrefix):
			// A criterion id in the reserved namespace collides with a unit the
			// runtime adds to every assessment catalogue. The collision cannot
			// be resolved later: the catalogue would then carry one id twice,
			// every possible answer would trip the coverage rule, and each
			// episode would die blaming the reviewing seat for a catalogue the
			// runtime built. Refusing the ISSUE is the only place it can be
			// refused honestly, because this is the document that is wrong.
			return Issue{}, fmt.Errorf("cdsissue: acceptance id %q uses the reserved %q prefix, which names runtime-measured units",
				c.ID, ReservedIDPrefix)
		}
		seen[c.ID] = true
	}
	return iss, nil
}

// exactShape states the closed key language at every object level, and checks
// the one thing decoding cannot tell apart afterwards: `scope.out` PRESENT but
// empty (non-goals considered, and there are none) versus ABSENT (never
// considered). Both decode to an empty slice, so presence has to be read off
// the raw document — the same present-vs-absent discipline cellspec.Resolve
// applies to supplied parameters.
func exactShape(raw json.RawMessage) error {
	if err := cellfill.OnlyKeys(raw, "the issue document", "kind", "id", "problem", "sources", "scope", "acceptance"); err != nil {
		return err
	}
	if p, ok := cellfill.Field(raw, "problem"); ok {
		if err := cellfill.OnlyKeys(p, "problem", "exists", "expected", "diverges"); err != nil {
			return err
		}
	}
	if err := eachElement(raw, "sources", "a source", "claim", "path"); err != nil {
		return err
	}
	if err := eachElement(raw, "acceptance", "an acceptance criterion", "id", "statement", "verification"); err != nil {
		return err
	}
	scope, ok := cellfill.Field(raw, "scope")
	if !ok {
		return fmt.Errorf("scope is required")
	}
	if err := cellfill.OnlyKeys(scope, "scope", "in", "out"); err != nil {
		return err
	}
	out, ok := cellfill.Field(scope, "out")
	if !ok || string(out) == "null" {
		return fmt.Errorf("scope.out must be present (an empty list declares that non-goals were considered; " +
			"an absent one declares nothing)")
	}
	return nil
}

// eachElement applies one closed key set to every element of a named array.
// Non-array and non-object values are left to the strict decode to judge —
// this reports on the key language only.
func eachElement(raw json.RawMessage, field, where string, allowed ...string) error {
	member, ok := cellfill.Field(raw, field)
	if !ok {
		return nil
	}
	var items []json.RawMessage
	if err := json.Unmarshal(member, &items); err != nil {
		return nil
	}
	for _, item := range items {
		if len(item) == 0 || item[0] != '{' {
			continue
		}
		if err := cellfill.OnlyKeys(item, where, allowed...); err != nil {
			return err
		}
	}
	return nil
}

// Render is the ONE rendering of an issue into seat-visible text. Both CDS
// seats call it and neither formats an issue field itself, so "α and β were
// told the same thing" is a property of there being one function, not of two
// formatters being kept in step.
func Render(iss Issue) string {
	var b strings.Builder
	fmt.Fprintf(&b, "===== ISSUE %s =====\n", iss.ID)
	b.WriteString("PROBLEM\n")
	fmt.Fprintf(&b, "  what exists:   %s\n", iss.Problem.Exists)
	fmt.Fprintf(&b, "  what is expected: %s\n", iss.Problem.Expected)
	fmt.Fprintf(&b, "  where they diverge: %s\n", iss.Problem.Diverges)
	b.WriteString("\nSOURCE OF TRUTH\n")
	for _, s := range iss.Sources {
		fmt.Fprintf(&b, "  %s: %s\n", s.Claim, s.Path)
	}
	b.WriteString("\nIN SCOPE\n")
	for _, s := range iss.Scope.In {
		fmt.Fprintf(&b, "  - %s\n", s)
	}
	b.WriteString("\nOUT OF SCOPE\n")
	if len(iss.Scope.Out) == 0 {
		b.WriteString("  (the author declared no non-goals)\n")
	}
	for _, s := range iss.Scope.Out {
		fmt.Fprintf(&b, "  - %s\n", s)
	}
	b.WriteString("\nACCEPTANCE CRITERIA\n")
	for _, c := range iss.Acceptance {
		fmt.Fprintf(&b, "  %s: %s\n", c.ID, c.Statement)
		fmt.Fprintf(&b, "    verified by: %s\n", c.Verification)
	}
	return b.String()
}
