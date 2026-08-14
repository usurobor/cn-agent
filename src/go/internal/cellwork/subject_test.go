package cellwork

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/usurobor/cnos/src/go/internal/celltest"
)

// corpusDir is the SAME directory scripts/cell-schema-check.sh vets with
// `cue vet … -d '#GitSnapshotPinned'`. Two authorities, one corpus: a subject
// the two disagree about fails a gate here rather than surfacing as a cell that
// vets and will not run.
const corpusDir = "../../../../schemas/cds/fixtures/subject"

// Every negative is invalid for exactly ONE reason, and the expected substring
// pins WHICH rule fired.
var negatives = map[string]string{
	"subject-bad-kind.json":       "subject kind must be",
	"subject-missing-repo.json":   "subject repo is required",
	"subject-empty-base.json":     "subject base_sha is required",
	"subject-unpinned-base.json":  "is not pinned",
	"subject-unknown-key.json":    `unknown key "branch"`,
	"subject-mixed-case-key.json": `unknown key "Repo"`,
}

func corpusFiles(t *testing.T) []string {
	t.Helper()
	entries, err := os.ReadDir(corpusDir)
	if err != nil {
		t.Fatalf("read subject corpus: %v", err)
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".json") {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	return names
}

// The corpus must actually contain both classes. A missing or renamed directory
// would otherwise make every assertion below vacuous.
func TestSubjectCorpusIsPopulatedAndFullyClassified(t *testing.T) {
	names := corpusFiles(t)
	valid := 0
	for _, n := range names {
		switch {
		case strings.HasPrefix(n, "valid-"):
			valid++
		case negatives[n] != "":
		default:
			t.Errorf("corpus file %q is neither a valid- fixture nor a named negative; "+
				"an unclassified fixture is vetted by CUE and by nothing here", n)
		}
	}
	if valid == 0 {
		t.Fatal("the subject corpus has no positive fixture")
	}
	if len(negatives) != len(names)-valid {
		t.Fatalf("negatives table lists %d files, corpus has %d non-positive files",
			len(negatives), len(names)-valid)
	}
}

func TestSubjectCorpusAdmission(t *testing.T) {
	for _, name := range corpusFiles(t) {
		t.Run(name, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join(corpusDir, name))
			if err != nil {
				t.Fatal(err)
			}
			_, admitErr := AdmitSubject(data)
			if strings.HasPrefix(name, "valid-") {
				if admitErr != nil {
					t.Fatalf("valid subject rejected: %v", admitErr)
				}
				return
			}
			want := negatives[name]
			if want == "" {
				t.Fatalf("no expected reason declared for %q", name)
			}
			if admitErr == nil {
				t.Fatal("must be rejected")
			}
			if !strings.Contains(admitErr.Error(), want) {
				t.Fatalf("rejected for the wrong reason:\n got: %v\nwant mention of: %q", admitErr, want)
			}
		})
	}
}

// The authored/pinned split is the whole reason there are two decoders: an
// author may name a moving revision, and a STATION may not.
func TestUnpinnedBaseParsesButIsNotAdmitted(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join(corpusDir, "subject-unpinned-base.json"))
	if err != nil {
		t.Fatal(err)
	}
	s, err := ParseSubject(raw)
	if err != nil {
		t.Fatalf("an authored subject may name a moving revision: %v", err)
	}
	if s.BaseSHA != "HEAD" {
		t.Fatalf("fixture is not the moving-revision case: %q", s.BaseSHA)
	}
	if _, err := AdmitSubject(raw); err == nil {
		t.Fatal("a station must not receive a subject whose base still has to be resolved")
	}
}

func TestAbsentSubjectIsRejected(t *testing.T) {
	for _, raw := range [][]byte{nil, {}} {
		if _, err := AdmitSubject(raw); err == nil {
			t.Fatalf("AdmitSubject(%q) must fail: a cell acting on nothing has no subject", raw)
		}
	}
}

// Pin turns an authored subject into the one both stations receive: the exact
// commit, and an absolute repository path rather than one relative to whatever
// directory a later reader is in.
func TestPinResolvesTheBaseAndTheRepoOnce(t *testing.T) {
	repo, head := celltest.Repo(t)
	authored, err := json.Marshal(Subject{Kind: SubjectKind, Repo: repo, BaseSHA: "HEAD"})
	if err != nil {
		t.Fatal(err)
	}

	pinned, err := Pin(context.Background(), authored)
	if err != nil {
		t.Fatalf("pin: %v", err)
	}
	s, err := AdmitSubject(pinned)
	if err != nil {
		t.Fatalf("a pinned subject must be admissible at a station: %v", err)
	}
	if s.BaseSHA != head {
		t.Fatalf("pinned base = %q, want the resolved HEAD %q", s.BaseSHA, head)
	}
	if !filepath.IsAbs(s.Repo) {
		t.Fatalf("pinned repo %q is not absolute", s.Repo)
	}

	// Pinning is a function of the repository state, not of the call: the same
	// authored subject pins to the same bytes.
	again, err := Pin(context.Background(), authored)
	if err != nil {
		t.Fatal(err)
	}
	if string(again) != string(pinned) {
		t.Fatalf("pinning twice gave two subjects:\n %s\n %s", pinned, again)
	}
	// And a moving ref really was moving: a second commit pins differently, so
	// the equality above is not the trivial "nothing ever changes".
	if err := os.WriteFile(filepath.Join(repo, "MOVED.md"), []byte("moved\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	celltest.Git(t, repo, "add", "-A")
	celltest.Git(t, repo, "commit", "-qm", "move HEAD")
	moved, err := Pin(context.Background(), authored)
	if err != nil {
		t.Fatal(err)
	}
	if string(moved) == string(pinned) {
		t.Fatal("HEAD moved and the pinned subject did not, so pinning is not resolving anything")
	}
}

// Pinning fails closed, and each cause is distinguishable — a repository that
// is not one, and a revision that does not exist, are different repairs.
func TestPinFailsClosed(t *testing.T) {
	repo, _ := celltest.Repo(t)
	cases := map[string]struct {
		subject Subject
		want    string
	}{
		"not a repository":  {Subject{Kind: SubjectKind, Repo: t.TempDir(), BaseSHA: "HEAD"}, "is not a git repository"},
		"unresolvable base": {Subject{Kind: SubjectKind, Repo: repo, BaseSHA: "no-such-rev"}, "does not resolve"},
		"unknown kind":      {Subject{Kind: "svn.checkout/0.1", Repo: repo, BaseSHA: "HEAD"}, "subject kind must be"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			raw, err := json.Marshal(tc.subject)
			if err != nil {
				t.Fatal(err)
			}
			_, err = Pin(context.Background(), raw)
			if err == nil {
				t.Fatal("must fail")
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("failed for the wrong reason: got %v, want mention of %q", err, tc.want)
			}
		})
	}
}
