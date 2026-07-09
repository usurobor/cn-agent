package issuesdispatch

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

// withFakeGitHub points githubAPIBase at an httptest server, mirroring
// issues-fsm's issuesfsm_test.go helper of the same name.
func withFakeGitHub(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(handler)
	orig := githubAPIBase
	githubAPIBase = srv.URL
	t.Cleanup(func() {
		srv.Close()
		githubAPIBase = orig
	})
	return srv
}

// TestLiveDefaults_RouteThroughSharedPrimitives proves the caller-path
// trace for the live default wiring (pre-review gate row 12): Dispatch's
// zero-value Options, once setDefaults runs, actually reaches
// liveGetIssue / ghEditIssueBody / ghRemoveLabel / ghAddLabel via real
// HTTP calls against a fake server -- not just the injected-fake path the
// other tests in this package exercise.
func TestLiveDefaults_RouteThroughSharedPrimitives(t *testing.T) {
	var mu sync.Mutex
	var requests []string
	withFakeGitHub(t, func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requests = append(requests, r.Method+" "+r.URL.Path)
		mu.Unlock()
		switch r.Method {
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"body":` + jsonQuote(contradictoryBody) + `,"labels":[{"name":"dispatch:cell"},{"name":"protocol:cds"},{"name":"status:ready"}]}`))
		case http.MethodPatch:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{}`))
		case http.MethodDelete:
			w.WriteHeader(http.StatusOK)
		case http.MethodPost:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[]`))
		}
	})

	opts := &Options{Repo: "acme/widgets", Issue: 640, Apply: true, Token: "tok"}
	res, err := Dispatch(context.Background(), opts)
	if err != nil {
		t.Fatalf("Dispatch via live defaults: %v", err)
	}
	if !res.Applied || !res.BodyChanged || !res.LabelFlipped {
		t.Fatalf("expected a fully-applied fresh-dispatch result, got %+v", res)
	}

	mu.Lock()
	got := append([]string(nil), requests...)
	mu.Unlock()

	var gets, patches, deletes, posts int
	for _, r := range got {
		switch {
		case strings.HasPrefix(r, "GET "):
			gets++
		case strings.HasPrefix(r, "PATCH "):
			patches++
		case strings.HasPrefix(r, "DELETE "):
			deletes++
		case strings.HasPrefix(r, "POST "):
			posts++
		}
	}
	if gets != 1 {
		t.Errorf("expected exactly 1 GET (fetch issue), got %d: %v", gets, got)
	}
	if patches != 1 {
		t.Errorf("expected exactly 1 PATCH (edit body), got %d: %v", patches, got)
	}
	if deletes != 1 {
		t.Errorf("expected exactly 1 DELETE (remove status:ready), got %d: %v", deletes, got)
	}
	if posts != 1 {
		t.Errorf("expected exactly 1 POST (add status:todo), got %d: %v", posts, got)
	}
}

// jsonQuote is a tiny helper so the fake handler above can embed
// contradictoryBody (which itself contains a backtick and newlines) as a
// valid JSON string literal without pulling in encoding/json just for
// this one test fixture.
func jsonQuote(s string) string {
	var b strings.Builder
	b.WriteByte('"')
	for _, r := range s {
		switch r {
		case '"':
			b.WriteString(`\"`)
		case '\\':
			b.WriteString(`\\`)
		case '\n':
			b.WriteString(`\n`)
		default:
			b.WriteRune(r)
		}
	}
	b.WriteByte('"')
	return b.String()
}
