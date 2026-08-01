package labeldoctor

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// withFakeGitHub points githubAPIBase at an httptest server for the
// duration of the test and restores the real default on cleanup, so no
// test in this package ever reaches the network. Mirrors
// cnos.issues/commands/issues-fsm/issuesfsm_test.go's withFakeGitHub
// idiom exactly.
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

func TestGhListLabels_Paginates(t *testing.T) {
	// 150 labels across 2 pages (100 + 50) exercises the "short page
	// ends pagination" loop condition.
	var all []ghLabel
	for i := 0; i < 150; i++ {
		all = append(all, ghLabel{Name: fmt.Sprintf("label-%03d", i), Color: "ededed"})
	}

	var pagesServed []string
	withFakeGitHub(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/repos/usurobor/cnos/labels" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		page := r.URL.Query().Get("page")
		pagesServed = append(pagesServed, page)
		var batch []ghLabel
		switch page {
		case "1":
			batch = all[:100]
		case "2":
			batch = all[100:]
		default:
			batch = nil
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(batch)
	})

	got, err := ghListLabels(context.Background(), "usurobor/cnos", "tok")
	if err != nil {
		t.Fatalf("ghListLabels: %v", err)
	}
	if len(got) != 150 {
		t.Fatalf("len(got) = %d, want 150", len(got))
	}
	if len(pagesServed) != 2 {
		t.Errorf("expected exactly 2 pages fetched, got %d: %v", len(pagesServed), pagesServed)
	}
}

func TestGhCreateLabel_Tolerates422AlreadyExists(t *testing.T) {
	withFakeGitHub(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(`{"message":"Validation Failed","errors":[{"code":"already_exists"}]}`))
	})

	err := ghCreateLabel(context.Background(), "usurobor/cnos", "tok", ghLabel{Name: "status:blocked", Color: "b60205"})
	if err != nil {
		t.Fatalf("ghCreateLabel should tolerate 422 as a no-op, got: %v", err)
	}
}

// TestGhCreateLabel_422WithoutAlreadyExistsCode_IsHardError is a
// regression test for a real bug found live against usurobor/cnos:
// ghCreateLabel originally tolerated ANY 422 as "already exists" — but
// GitHub also returns 422 for other validation failures (e.g. a
// description over 100 characters), which must NOT be silently
// swallowed as a false success.
func TestGhCreateLabel_422WithoutAlreadyExistsCode_IsHardError(t *testing.T) {
	withFakeGitHub(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(`{"message":"Validation Failed","errors":[{"resource":"Label","code":"custom","field":"description","message":"description is too long (maximum is 100 characters)"}]}`))
	})

	err := ghCreateLabel(context.Background(), "usurobor/cnos", "tok", ghLabel{Name: "dispatch:cell", Color: "1d76db", Description: "way too long"})
	if err == nil {
		t.Fatal("expected a 422 with a non-already_exists code to be a hard error, not a tolerated no-op")
	}
	if !strings.Contains(err.Error(), "too long") {
		t.Errorf("error should surface the underlying GitHub validation message, got: %v", err)
	}
}

func TestGhCreateLabel_HardErrorOnOtherStatus(t *testing.T) {
	withFakeGitHub(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"message":"insufficient scope"}`))
	})

	if err := ghCreateLabel(context.Background(), "usurobor/cnos", "tok", ghLabel{Name: "status:blocked", Color: "b60205"}); err == nil {
		t.Fatal("expected an error for a 403 response")
	}
}

func TestGhUpdateLabel_PatchesColorAndDescription(t *testing.T) {
	var gotBody map[string]string
	withFakeGitHub(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		// r.URL.Path is the decoded form (net/http decodes %XX into
		// .Path); the %3A-escaped colon only shows up on the wire /
		// in .RequestURI(), not here.
		if r.URL.Path != "/repos/usurobor/cnos/labels/status:review" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		w.WriteHeader(http.StatusOK)
	})

	if err := ghUpdateLabel(context.Background(), "usurobor/cnos", "tok", "status:review", "5319e7", "Cell complete."); err != nil {
		t.Fatalf("ghUpdateLabel: %v", err)
	}
	if gotBody["color"] != "5319e7" || gotBody["description"] != "Cell complete." {
		t.Errorf("PATCH body = %+v, want color=5319e7 description=%q", gotBody, "Cell complete.")
	}
}

func TestGhUpdateLabel_HardErrorOnNon200(t *testing.T) {
	withFakeGitHub(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	if err := ghUpdateLabel(context.Background(), "usurobor/cnos", "tok", "status:review", "5319e7", "x"); err == nil {
		t.Fatal("expected an error for a 404 response")
	}
}
