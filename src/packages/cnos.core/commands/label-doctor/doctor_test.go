package labeldoctor

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"path/filepath"
	"sync"
	"testing"
)

func TestAudit_ClassifiesMatchDriftedMissing(t *testing.T) {
	manifest := Manifest{
		Schema: "cn.labels.v1",
		Labels: []Label{
			{Name: "status:backlog", Color: "ededed", Description: "d1"},
			{Name: "status:review", Color: "5319e7", Description: "d2"},
			{Name: "status:blocked", Color: "b60205", Description: "d3"},
		},
	}
	live := map[string]ghLabel{
		"status:backlog": {Name: "status:backlog", Color: "ededed", Description: "d1"},          // match
		"status:review":  {Name: "status:review", Color: "EDEDED", Description: "d2-different"}, // drifted (color AND description)
		// status:blocked absent -> missing
	}

	findings := Audit(manifest, live)
	if len(findings) != 3 {
		t.Fatalf("len(findings) = %d, want 3", len(findings))
	}
	byName := map[string]Finding{}
	for _, f := range findings {
		byName[f.Name] = f
	}
	if byName["status:backlog"].Status != StatusMatch {
		t.Errorf("status:backlog = %v, want match", byName["status:backlog"].Status)
	}
	if byName["status:review"].Status != StatusDrifted {
		t.Errorf("status:review = %v, want drifted", byName["status:review"].Status)
	}
	if byName["status:blocked"].Status != StatusMissing {
		t.Errorf("status:blocked = %v, want missing", byName["status:blocked"].Status)
	}
}

func TestAudit_ColorCaseInsensitive_DescriptionExact(t *testing.T) {
	manifest := Manifest{Labels: []Label{{Name: "x", Color: "AbCdEf", Description: "same"}}}
	live := map[string]ghLabel{"x": {Name: "x", Color: "abcdef", Description: "same"}}
	findings := Audit(manifest, live)
	if findings[0].Status != StatusMatch {
		t.Errorf("expected case-insensitive color match, got %v", findings[0].Status)
	}
}

// fakeLabelStore is a minimal in-memory GitHub label API used to drive
// both the dry-run pass/fail tests and the apply + idempotence test
// against a real HTTP round trip (via withFakeGitHub), without ever
// touching the network.
type fakeLabelStore struct {
	mu      sync.Mutex
	labels  map[string]ghLabel
	posts   int
	patches int
}

func newFakeLabelStore(seed []ghLabel) *fakeLabelStore {
	s := &fakeLabelStore{labels: map[string]ghLabel{}}
	for _, l := range seed {
		s.labels[l.Name] = l
	}
	return s
}

func (s *fakeLabelStore) handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/repos/usurobor/cnos/labels":
			page := r.URL.Query().Get("page")
			var out []ghLabel
			if page == "" || page == "1" {
				for _, l := range s.labels {
					out = append(out, l)
				}
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(out)
		case r.Method == http.MethodPost && r.URL.Path == "/repos/usurobor/cnos/labels":
			s.posts++
			var l ghLabel
			json.NewDecoder(r.Body).Decode(&l)
			if _, exists := s.labels[l.Name]; exists {
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			s.labels[l.Name] = l
			w.WriteHeader(http.StatusCreated)
		case r.Method == http.MethodPatch:
			s.patches++
			name, _ := url.PathUnescape(r.URL.Path[len("/repos/usurobor/cnos/labels/"):])
			var body struct {
				Color       string `json:"color"`
				Description string `json:"description"`
			}
			json.NewDecoder(r.Body).Decode(&body)
			l := s.labels[name]
			l.Name = name
			l.Color = body.Color
			l.Description = body.Description
			s.labels[name] = l
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func (s *fakeLabelStore) mutatingCalls() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.posts + s.patches
}

func writeManifestFixture(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "labels.json")
	writeFile(t, path, fixtureManifestJSON)
	return path
}

// TestDoctor_DryRun_FailsOnDrift is AC5's negative direction: --dry-run
// against a repo with injected drift (status:review at the GitHub
// default ededed instead of canonical 5319e7, status:blocked entirely
// absent) must return ErrDrift and must make zero mutating calls.
func TestDoctor_DryRun_FailsOnDrift(t *testing.T) {
	store := newFakeLabelStore([]ghLabel{
		{Name: "status:backlog", Color: "ededed", Description: "Well-formed scope but not yet refined."},
		{Name: "status:review", Color: "ededed", Description: "Cell complete; awaiting review."}, // drifted color
		// status:blocked missing entirely
	})
	withFakeGitHub(t, store.handler())

	res, err := Doctor(context.Background(), Options{
		Repo:       "usurobor/cnos",
		Token:      "tok",
		LabelsPath: writeManifestFixture(t),
		DryRun:     true,
	})
	if !errors.Is(err, ErrDrift) {
		t.Fatalf("expected ErrDrift, got: %v", err)
	}
	if res == nil || len(res.Findings) != 3 {
		t.Fatalf("expected 3 findings, got: %+v", res)
	}
	if store.mutatingCalls() != 0 {
		t.Errorf("dry-run must make zero mutating calls, got %d", store.mutatingCalls())
	}
}

// TestDoctor_DryRun_PassesOnClean is AC5's positive direction: --dry-run
// against a repo whose live labels already match the manifest exactly
// must return nil (not ErrDrift), and still make zero mutating calls.
func TestDoctor_DryRun_PassesOnClean(t *testing.T) {
	store := newFakeLabelStore([]ghLabel{
		{Name: "status:backlog", Color: "ededed", Description: "Well-formed scope but not yet refined."},
		{Name: "status:review", Color: "5319e7", Description: "Cell complete; awaiting review."},
		{Name: "status:blocked", Color: "b60205", Description: "Gated on external input."},
	})
	withFakeGitHub(t, store.handler())

	res, err := Doctor(context.Background(), Options{
		Repo:       "usurobor/cnos",
		Token:      "tok",
		LabelsPath: writeManifestFixture(t),
		DryRun:     true,
	})
	if err != nil {
		t.Fatalf("expected nil error on a clean repo, got: %v", err)
	}
	for _, f := range res.Findings {
		if f.Status != StatusMatch {
			t.Errorf("finding %+v should be match on a clean repo", f)
		}
	}
	if store.mutatingCalls() != 0 {
		t.Errorf("dry-run must make zero mutating calls, got %d", store.mutatingCalls())
	}
}

// TestDoctor_Apply_RepairsDriftAndIsIdempotent is AC4's core functional
// + idempotence oracle: a non-dry-run call against a drifted/missing
// live state creates/repairs everything to canonical, and a SECOND call
// (against the now-repaired live state) performs zero mutating API
// calls — asserted on the fake server's actual request counts, not just
// "exit 0 twice".
func TestDoctor_Apply_RepairsDriftAndIsIdempotent(t *testing.T) {
	store := newFakeLabelStore([]ghLabel{
		{Name: "status:backlog", Color: "ededed", Description: "Well-formed scope but not yet refined."},
		{Name: "status:review", Color: "ededed", Description: "wrong description"}, // drifted
		// status:blocked missing
	})
	withFakeGitHub(t, store.handler())
	labelsPath := writeManifestFixture(t)

	res1, err := Doctor(context.Background(), Options{
		Repo:       "usurobor/cnos",
		Token:      "tok",
		LabelsPath: labelsPath,
	})
	if err != nil {
		t.Fatalf("first Doctor call: %v", err)
	}
	if len(res1.Applied) != 2 {
		t.Fatalf("expected 2 labels applied (1 create + 1 repair), got %v", res1.Applied)
	}
	if store.posts != 1 || store.patches != 1 {
		t.Fatalf("expected exactly 1 POST + 1 PATCH, got posts=%d patches=%d", store.posts, store.patches)
	}

	// Verify the live store actually reflects canonical state now.
	if store.labels["status:review"].Color != "5319e7" {
		t.Errorf("status:review color after repair = %q, want 5319e7", store.labels["status:review"].Color)
	}
	if _, ok := store.labels["status:blocked"]; !ok {
		t.Error("status:blocked should have been created")
	}

	// Second call against the now-canonical live state: idempotence.
	res2, err := Doctor(context.Background(), Options{
		Repo:       "usurobor/cnos",
		Token:      "tok",
		LabelsPath: labelsPath,
	})
	if err != nil {
		t.Fatalf("second Doctor call: %v", err)
	}
	if len(res2.Applied) != 0 {
		t.Errorf("second call should apply nothing, got %v", res2.Applied)
	}
	if store.posts != 1 || store.patches != 1 {
		t.Errorf("second call must make zero additional mutating calls, got posts=%d patches=%d", store.posts, store.patches)
	}
}

func TestDoctor_ResolvesRepoAndManifestFromRepoRoot(t *testing.T) {
	store := newFakeLabelStore([]ghLabel{
		{Name: "status:backlog", Color: "ededed", Description: "Well-formed scope but not yet refined."},
		{Name: "status:review", Color: "5319e7", Description: "Cell complete; awaiting review."},
		{Name: "status:blocked", Color: "b60205", Description: "Gated on external input."},
	})
	withFakeGitHub(t, store.handler())

	root := t.TempDir()
	runGit(t, root, "init", "-q")
	runGit(t, root, "remote", "add", "origin", "https://github.com/usurobor/cnos.git")
	writeFile(t, filepath.Join(root, "src", "packages", "cnos.core", "labels.json"), fixtureManifestJSON)

	res, err := Doctor(context.Background(), Options{RepoRoot: root, Token: "tok", DryRun: true})
	if err != nil {
		t.Fatalf("Doctor with resolved repo/manifest: %v", err)
	}
	if res.Repo != "usurobor/cnos" {
		t.Errorf("Repo = %q, want usurobor/cnos", res.Repo)
	}
}

func TestDoctor_ManifestLoadErrorPropagates(t *testing.T) {
	_, err := Doctor(context.Background(), Options{
		Repo:       "usurobor/cnos",
		Token:      "tok",
		LabelsPath: filepath.Join(t.TempDir(), "nope.json"),
	})
	if err == nil {
		t.Fatal("expected a manifest-load error")
	}
}
