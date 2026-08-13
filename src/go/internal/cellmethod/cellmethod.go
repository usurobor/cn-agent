// Package cellmethod loads the cell's ONE methodology bundle and projects it
// into the views the seats receive.
//
// Before this package a producing seat carried its own `skills` list and any
// assessing seat would have carried a second one. Two lists that are supposed
// to state the same obligations are two lists that drift, and nothing could
// tell that they had: each was internally consistent, and the record showed two
// independently digested collections with no statement that they were meant to
// agree. The bundle is declared once, on the cell, loaded once, digested once,
// and PROJECTED — so a seat cannot be held to obligations the cell never
// declared, and two seats cannot be held to different ones.
//
// PROJECTION IS IDENTITY UNDER A FIXED ROLE WRAPPER, and that is the whole of
// it in this bootstrap. Both views carry the same ordered bodies, byte for
// byte, differing only in a leading role paragraph. This makes "every
// obligation appears in both projections" TRIVIALLY TRUE today: the two texts
// cannot disagree about an obligation because neither selects anything. It
// becomes a real property only when an obligation can be absent from one view
// — that is, when a methodology carries stable property ids and a projection
// can drop or transform one. Nothing here should be read as evidence that
// projection preserves obligations; it is evidence that the bootstrap does not
// yet select.
//
// The digest covers the ordered (ref, body-digest) list, so it changes when a
// skill body changes, when a skill is added or removed, and when two skills
// swap places. Order is meaning — later skills refine earlier ones — so a
// reordering is a different methodology and gets a different digest.
package cellmethod

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/usurobor/cnos/src/go/internal/cellskill"
)

// Kind is the pinned methodology declaration kind. A bundle must declare it
// exactly; there is one shape and no version negotiation.
const Kind = "skills.methodology.v0"

// Ref is one entry of the bundle as the record carries it: the canonical skill
// reference and the digest of the body that was actually loaded. Naming a
// skill is not loading it, so the digest is what makes the reference a claim
// about content rather than about a filename.
type Ref struct {
	Ref    string `json:"ref"`
	SHA256 string `json:"sha256"`
}

// Bundle is the loaded methodology: what the cell declared, what it resolved
// to, and the one digest that identifies it.
type Bundle struct {
	Kind   string
	Skills []Ref
	// SHA256 digests the ordered (ref, body-digest) list — nothing else. Kind
	// is deliberately outside it: exactly one kind is admissible, so including
	// it would mix a constant into every digest and identify nothing.
	SHA256 string
}

// Role names which projection a view is. It is on the view because a seat's
// record says which projection held it, and "constructive" and "adversarial"
// must not be told apart by inspecting prose.
type Role string

const (
	RoleConstructive Role = "constructive"
	RoleAdversarial  Role = "adversarial"
)

// View is what a seat receives: the rendered text, and the digest of the
// bundle it was projected from. The digest travels with the text so a seat's
// record can state which methodology held it without re-deriving anything from
// the prose.
type View struct {
	Role   Role
	Text   string
	SHA256 string
}

// Empty reports a view no methodology was projected into. A cell may declare
// no methodology at all — every cell did before this package existed — so the
// zero View is a legitimate value, and a fill that cannot act without a
// methodology says so itself rather than having the loader decide for everyone.
func (v View) Empty() bool { return v.SHA256 == "" }

// The role wrappers. They are the ONLY difference between the two projections,
// and each states what it is rather than dressing up a selection that does not
// happen: both texts carry the identical ordered bodies below.
const (
	constructiveWrapper = "METHODOLOGY (constructive projection)\n" +
		"You are the PRODUCING seat. The ordered skill bodies below are the methodology\n" +
		"this cell declared, in the order it declared them: later skills refine earlier\n" +
		"ones. Follow them while producing.\n"

	adversarialWrapper = "METHODOLOGY (adversarial projection)\n" +
		"You are the ASSESSING seat. The ordered skill bodies below are the methodology\n" +
		"the candidate is held to, in the order the cell declared them: later skills\n" +
		"refine earlier ones. Try to falsify each obligation against the candidate.\n"
)

// Load decodes a methodology declaration, resolves every skill ref against
// `r`, and returns the bundle plus the loaded bodies. The bodies are returned
// rather than kept because a Bundle is what a RECORD carries — refs and
// digests — while the bodies are what a PROJECTION renders, and a value that
// travels into a receipt should not carry a megabyte of prose it never states.
//
// The strict decode is written out here rather than reusing cellfill's
// OnlyKeys/StrictDecode: cellfill imports THIS package for the alpha factory
// signature, so importing it back is an import cycle. Hoisting those two
// helpers into a third package would move them, and two other fills' imports,
// for no gain this increment.
func Load(r cellskill.Resolver, decl []byte) (Bundle, []cellskill.Skill, error) {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(decl, &obj); err != nil {
		return Bundle{}, nil, fmt.Errorf("cellmethod: methodology is not an object: %w", err)
	}
	for k := range obj {
		// Exact and case-sensitive: encoding/json would otherwise decode `Kind`
		// and `Skills` happily while the closed CUE definition rejects them.
		if k != "kind" && k != "skills" {
			return Bundle{}, nil, fmt.Errorf(
				"cellmethod: methodology has unknown key %q (keys are exact and case-sensitive)", k)
		}
	}
	var d struct {
		Kind   string   `json:"kind"`
		Skills []string `json:"skills"`
	}
	dec := json.NewDecoder(bytes.NewReader(decl))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&d); err != nil {
		return Bundle{}, nil, fmt.Errorf("cellmethod: decode methodology: %w", err)
	}
	var extra json.RawMessage
	if err := dec.Decode(&extra); err != io.EOF {
		// One value decoded is not one value supplied: `{...} garbage` must not
		// be read as the leading object with the rest silently discarded.
		if err == nil {
			return Bundle{}, nil, fmt.Errorf("cellmethod: trailing data after the methodology")
		}
		return Bundle{}, nil, fmt.Errorf("cellmethod: malformed data after the methodology: %w", err)
	}
	if d.Kind != Kind {
		return Bundle{}, nil, fmt.Errorf("cellmethod: methodology kind must be %q, got %q", Kind, d.Kind)
	}
	if len(d.Skills) == 0 {
		return Bundle{}, nil, fmt.Errorf("cellmethod: a methodology with no skills states no obligations")
	}
	if r == nil {
		return Bundle{}, nil, fmt.Errorf("cellmethod: no skill resolver configured")
	}
	bodies, err := cellskill.LoadAll(r, d.Skills)
	if err != nil {
		return Bundle{}, nil, fmt.Errorf("cellmethod: %w", err)
	}
	refs := make([]Ref, 0, len(bodies))
	for _, s := range bodies {
		refs = append(refs, Ref{Ref: s.Ref, SHA256: s.SHA256})
	}
	return Bundle{Kind: d.Kind, Skills: refs, SHA256: digest(refs)}, bodies, nil
}

// digest hashes the ordered (ref, body-digest) pairs. Each ref is
// LENGTH-PREFIXED so no ref content can forge a record boundary: without it a
// ref containing the separator could produce the same byte stream as a
// different two-skill list, and two different methodologies would share a
// digest.
func digest(refs []Ref) string {
	h := sha256.New()
	for _, r := range refs {
		fmt.Fprintf(h, "%d:%s %s\n", len(r.Ref), r.Ref, r.SHA256)
	}
	return hex.EncodeToString(h.Sum(nil))
}

// Constructive is the projection the producing seat receives.
func Constructive(b Bundle, bodies []cellskill.Skill) View {
	return project(RoleConstructive, constructiveWrapper, b, bodies)
}

// Adversarial is the projection an assessing seat receives.
//
// IT HAS NO PRODUCTION CALLER. Nothing in the runtime projects it today: the
// assessing seat that will consume it is a later increment, and until that seat
// exists this function is exercised by this package's tests alone.
// TestAdversarialHasNoProductionCaller fails if that stops being true, so the
// sentence above cannot quietly become false.
func Adversarial(b Bundle, bodies []cellskill.Skill) View {
	return project(RoleAdversarial, adversarialWrapper, b, bodies)
}

func project(role Role, wrapper string, b Bundle, bodies []cellskill.Skill) View {
	return View{Role: role, Text: wrapper + render(bodies), SHA256: b.SHA256}
}

// render writes the ordered bodies VERBATIM, each under a header naming its
// canonical ref and the digest of the body that follows. Pure and shared by
// both projections — one renderer, so the two views cannot differ in anything
// but their wrapper by construction rather than by a comment.
func render(bodies []cellskill.Skill) string {
	var b strings.Builder
	for _, s := range bodies {
		fmt.Fprintf(&b, "\n===== SKILL %s (sha256 %s) =====\n", s.Ref, s.SHA256)
		b.WriteString(s.Body)
		if !strings.HasSuffix(s.Body, "\n") {
			b.WriteString("\n")
		}
	}
	return b.String()
}
