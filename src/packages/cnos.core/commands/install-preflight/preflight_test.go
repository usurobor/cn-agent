package installpreflight

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// newFakeGitHub returns an httptest.Server answering both endpoints
// Check calls: GET .../actions/secrets (reports exactly presentSecrets)
// and GET /repos/{owner}/{repo} (reports permissions.push == pushAccess).
func newFakeGitHub(t *testing.T, presentSecrets []string, pushAccess bool) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if strings.HasSuffix(r.URL.Path, "/actions/secrets") {
			var b []byte
			b = append(b, []byte(`{"total_count":`)...)
			b = append(b, []byte(fmt.Sprintf("%d", len(presentSecrets)))...)
			b = append(b, []byte(`,"secrets":[`)...)
			for i, name := range presentSecrets {
				if i > 0 {
					b = append(b, ',')
				}
				b = append(b, []byte(fmt.Sprintf(`{"name":%q,"created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-01T00:00:00Z"}`, name))...)
			}
			b = append(b, []byte(`]}`)...)
			w.Write(b)
			return
		}
		fmt.Fprintf(w, `{"permissions":{"push":%v}}`, pushAccess)
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestCheck_AllPresent_Ready(t *testing.T) {
	srv := newFakeGitHub(t, []string{"CLAUDE_CODE_OAUTH_TOKEN", "CN_DISPATCH_PAT"}, true)

	res, err := Check(context.Background(), Options{
		Repo:        "acme/widgets",
		Token:       "tok",
		APIBaseURL:  srv.URL,
		SecretNames: []string{"CLAUDE_CODE_OAUTH_TOKEN", "CN_DISPATCH_PAT"},
	})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if !res.Ready() {
		t.Errorf("Ready() = false, want true: %+v", res)
	}
	if len(res.Missing) != 0 {
		t.Errorf("Missing = %v, want empty", res.Missing)
	}
	if len(res.Present) != 2 {
		t.Errorf("Present = %v, want both names", res.Present)
	}
}

func TestCheck_MissingSecret_NotReady(t *testing.T) {
	srv := newFakeGitHub(t, []string{"CLAUDE_CODE_OAUTH_TOKEN"}, true)

	res, err := Check(context.Background(), Options{
		Repo:        "acme/widgets",
		Token:       "tok",
		APIBaseURL:  srv.URL,
		SecretNames: []string{"CLAUDE_CODE_OAUTH_TOKEN", "CN_DISPATCH_PAT"},
	})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if res.Ready() {
		t.Error("Ready() = true, want false (CN_DISPATCH_PAT missing)")
	}
	if len(res.Missing) != 1 || res.Missing[0] != "CN_DISPATCH_PAT" {
		t.Errorf("Missing = %v, want [CN_DISPATCH_PAT]", res.Missing)
	}
	if len(res.Present) != 1 || res.Present[0] != "CLAUDE_CODE_OAUTH_TOKEN" {
		t.Errorf("Present = %v, want [CLAUDE_CODE_OAUTH_TOKEN]", res.Present)
	}
}

func TestCheck_NoPushAccess_NotReady(t *testing.T) {
	srv := newFakeGitHub(t, []string{"CLAUDE_CODE_OAUTH_TOKEN", "CN_DISPATCH_PAT"}, false)

	res, err := Check(context.Background(), Options{
		Repo:        "acme/widgets",
		Token:       "tok",
		APIBaseURL:  srv.URL,
		SecretNames: []string{"CLAUDE_CODE_OAUTH_TOKEN", "CN_DISPATCH_PAT"},
	})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if res.Ready() {
		t.Error("Ready() = true, want false (no push access)")
	}
	if res.PushAccess {
		t.Error("PushAccess = true, want false")
	}
	if len(res.Missing) != 0 {
		t.Errorf("Missing = %v, want empty (both secrets present)", res.Missing)
	}
}

func TestCheck_EmptySecretNames_OnlyChecksPushAccess(t *testing.T) {
	// The engine tier (cnos#613/#706) needs no repo secret at all — Check
	// must still report Ready() based solely on push access when
	// SecretNames is empty (repoinstall.go's requiredSecretNames returns
	// nil for the engine tier).
	srv := newFakeGitHub(t, nil, true)

	res, err := Check(context.Background(), Options{
		Repo:       "acme/widgets",
		Token:      "tok",
		APIBaseURL: srv.URL,
	})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if !res.Ready() {
		t.Errorf("Ready() = false, want true (no secrets required, push access present): %+v", res)
	}
}

func TestCheck_RepoOverride_SkipsGitRemoteResolution(t *testing.T) {
	// RepoRoot is deliberately empty/unusable (no git repo at all) — Check
	// must not attempt git-remote resolution when Repo is already set.
	srv := newFakeGitHub(t, []string{"X"}, true)

	res, err := Check(context.Background(), Options{
		RepoRoot:    t.TempDir(), // not a git repository
		Repo:        "acme/widgets",
		Token:       "tok",
		APIBaseURL:  srv.URL,
		SecretNames: []string{"X"},
	})
	if err != nil {
		t.Fatalf("Check: %v (should not have attempted git-remote resolution)", err)
	}
	if res.Repo != "acme/widgets" {
		t.Errorf("Repo = %q, want acme/widgets", res.Repo)
	}
}

func TestCheck_NoGitRemote_SurfacesNamedError(t *testing.T) {
	// No Repo override and RepoRoot has no git remote at all: Check must
	// surface a named, actionable error (mirroring label-doctor's own
	// resolveRepoFromGitRemote failure mode) rather than attempting any
	// network call.
	_, err := Check(context.Background(), Options{
		RepoRoot: t.TempDir(),
	})
	if err == nil {
		t.Fatal("expected an error when no git remote is resolvable")
	}
	if got := err.Error(); got == "" {
		t.Error("expected a non-empty error message")
	}
}
