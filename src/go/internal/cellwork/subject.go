package cellwork

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/usurobor/cnos/src/go/internal/cellfill"
)

// SubjectKind is the pinned tag of the only subject this adapter speaks. It is
// versioned in the tag itself rather than in a sibling field so that a document
// carrying an unknown snapshot shape is rejected by the one comparison every
// reader already makes.
const SubjectKind = "git.snapshot/0.1"

// Subject is what an episode acts on: one repository at one commit. It rides
// the contract as opaque bytes — the kernel and the generic loader never decode
// it — and this package is the ONLY place its language lives, so the producing
// and assessing seats cannot read the same bytes differently.
type Subject struct {
	Kind    string `json:"kind"`
	Repo    string `json:"repo"`
	BaseSHA string `json:"base_sha"`
}

// pinnedSHA is the recorded form of a base: a full commit object name. Short
// and upper-case spellings are excluded because two spellings of one commit are
// two documents, and the record's bytes are digested.
var pinnedSHA = regexp.MustCompile(`^[0-9a-f]{40}$`)

// ParseSubject decodes an AUTHORED subject. Pure: bytes in, value out, no IO
// and no paths touched (eng/go §2.17), which is what lets both the pinning step
// and the two seats share one decoder instead of three.
//
// It deliberately admits a base that is not yet pinned — `HEAD`, a branch, a
// tag — because that is exactly what an author may write and what Pin exists to
// resolve. Emptiness rather than blankness is the rule for the two string
// fields, matching CUE's `!=""` character for character: a second blankness
// predicate written by hand in two languages is what diverged the last time
// (see cdsissue's nonBlankPattern), and a path made of spaces fails at the
// filesystem anyway.
//
// The closed key language is stated with cellfill.OnlyKeys rather than left to
// the decoder: encoding/json matches field names case-insensitively even with
// DisallowUnknownFields, so `Kind` or `Repo` would decode here while the closed
// CUE definition rejects it.
func ParseSubject(raw []byte) (Subject, error) {
	if len(raw) == 0 {
		return Subject{}, fmt.Errorf("cellwork: contract carries no subject")
	}
	if err := cellfill.OnlyKeys(raw, "cellwork subject", "kind", "repo", "base_sha"); err != nil {
		return Subject{}, fmt.Errorf("cellwork: %w", err)
	}
	var s Subject
	if err := cellfill.StrictDecode(raw, &s); err != nil {
		return Subject{}, fmt.Errorf("cellwork: subject: %w", err)
	}
	switch {
	case s.Kind != SubjectKind:
		return Subject{}, fmt.Errorf("cellwork: subject kind must be %q, got %q", SubjectKind, s.Kind)
	case s.Repo == "":
		return Subject{}, fmt.Errorf("cellwork: subject repo is required")
	case s.BaseSHA == "":
		return Subject{}, fmt.Errorf("cellwork: subject base_sha is required")
	}
	return s, nil
}

// AdmitSubject decodes a FROZEN subject: ParseSubject plus the one rule a
// recorded subject must satisfy — its base names an exact commit. Both stations
// call this and neither calls ParseSubject, so an unpinned base cannot reach a
// station: re-resolving a moving name per station is precisely how the two
// could measure against different trees while the record claimed one.
//
// Pure, for the same reason ParseSubject is.
func AdmitSubject(raw []byte) (Subject, error) {
	s, err := ParseSubject(raw)
	if err != nil {
		return Subject{}, err
	}
	if !pinnedSHA.MatchString(s.BaseSHA) {
		return Subject{}, fmt.Errorf("cellwork: subject base_sha %q is not pinned "+
			"(a recorded subject names a full commit, not a revision to resolve)", s.BaseSHA)
	}
	return s, nil
}

// Pin resolves an authored subject to the exact repository and commit the
// episode will use, ONCE, before either station exists. This is the effectful
// half of the subject language and the only one: everything downstream reads
// the pinned bytes.
//
// The repository path is made absolute as well, so a recorded subject is not
// silently relative to whatever directory a later reader happens to be in — the
// same reason the base is resolved rather than carried as `HEAD`.
func Pin(ctx context.Context, raw json.RawMessage) (json.RawMessage, error) {
	s, err := ParseSubject(raw)
	if err != nil {
		return nil, err
	}
	repoAbs, err := repoPath(s.Repo)
	if err != nil {
		return nil, err
	}
	sha, err := ResolveBase(ctx, repoAbs, s.BaseSHA)
	if err != nil {
		return nil, err
	}
	pinned, err := json.Marshal(Subject{Kind: s.Kind, Repo: repoAbs, BaseSHA: sha})
	if err != nil {
		return nil, fmt.Errorf("cellwork: canonicalize subject: %w", err)
	}
	return pinned, nil
}
