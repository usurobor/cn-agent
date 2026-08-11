package installpreflight

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestGhSecret_HasNoValueField(t *testing.T) {
	// AC3's mechanical oracle: the wire type this package decodes GitHub's
	// secrets-list response into must have NO field that could ever hold a
	// secret's value — not "value", not "Value", not any case variant.
	// This is a structural guarantee, not a runtime behavior: even a
	// malicious or buggy server response carrying a "value" key has
	// nowhere in ghSecret to land.
	typ := reflect.TypeOf(ghSecret{})
	for i := 0; i < typ.NumField(); i++ {
		name := strings.ToLower(typ.Field(i).Name)
		if strings.Contains(name, "value") || strings.Contains(name, "secret") {
			t.Errorf("ghSecret field %q could hold a secret value — AC3 violation", typ.Field(i).Name)
		}
	}
	if typ.NumField() != 3 {
		t.Errorf("ghSecret has %d fields, want exactly 3 (name, created_at, updated_at) — an unexpected field is a potential AC3 leak surface", typ.NumField())
	}
}

// TestGhListSecrets_PresenceOnlyResponse_DecodesCorrectly is AC3's
// integration oracle: an httptest fixture returning GitHub's REAL
// response shape for GET .../actions/secrets — {name, created_at,
// updated_at}, no "value" key at all — decodes correctly and drives
// presence detection. Proves the code path never expects a value field
// to be present.
func TestGhListSecrets_PresenceOnlyResponse_DecodesCorrectly(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/acme/widgets/actions/secrets" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		// Exact GitHub response shape (no "value" key — GitHub's API
		// never returns one for this endpoint).
		fmt.Fprint(w, `{
			"total_count": 2,
			"secrets": [
				{"name": "CLAUDE_CODE_OAUTH_TOKEN", "created_at": "2026-01-01T00:00:00Z", "updated_at": "2026-01-01T00:00:00Z"},
				{"name": "CN_DISPATCH_PAT", "created_at": "2026-01-01T00:00:00Z", "updated_at": "2026-01-02T00:00:00Z"}
			]
		}`)
	}))
	defer srv.Close()

	got, err := ghListSecrets(context.Background(), srv.URL, "acme/widgets", "tok")
	if err != nil {
		t.Fatalf("ghListSecrets: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len(got) = %d, want 2", len(got))
	}
	if got[0].Name != "CLAUDE_CODE_OAUTH_TOKEN" || got[1].Name != "CN_DISPATCH_PAT" {
		t.Errorf("unexpected names: %+v", got)
	}
}

func TestGhListSecrets_Paginates(t *testing.T) {
	var pagesServed []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		pagesServed = append(pagesServed, page)
		var batch []ghSecret
		switch page {
		case "1":
			for i := 0; i < 100; i++ {
				batch = append(batch, ghSecret{Name: fmt.Sprintf("SECRET_%03d", i)})
			}
		case "2":
			for i := 100; i < 150; i++ {
				batch = append(batch, ghSecret{Name: fmt.Sprintf("SECRET_%03d", i)})
			}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ghSecretsPage{TotalCount: 150, Secrets: batch})
	}))
	defer srv.Close()

	got, err := ghListSecrets(context.Background(), srv.URL, "acme/widgets", "tok")
	if err != nil {
		t.Fatalf("ghListSecrets: %v", err)
	}
	if len(got) != 150 {
		t.Fatalf("len(got) = %d, want 150", len(got))
	}
	if len(pagesServed) != 2 {
		t.Errorf("expected exactly 2 pages fetched, got %d: %v", len(pagesServed), pagesServed)
	}
}

func TestGhListSecrets_HTTPErrorSurfaced(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, `{"message":"Bad credentials"}`)
	}))
	defer srv.Close()

	_, err := ghListSecrets(context.Background(), srv.URL, "acme/widgets", "bad-tok")
	if err == nil {
		t.Fatal("expected an error for HTTP 401")
	}
	if !strings.Contains(err.Error(), "401") {
		t.Errorf("error should name the HTTP status, got: %v", err)
	}
}

func TestGhCheckPushAccess_TrueAndFalse(t *testing.T) {
	for _, want := range []bool{true, false} {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/repos/acme/widgets" {
				t.Errorf("unexpected path: %s", r.URL.Path)
			}
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"full_name":"acme/widgets","permissions":{"admin":false,"push":%v,"pull":true}}`, want)
		}))
		got, err := ghCheckPushAccess(context.Background(), srv.URL, "acme/widgets", "tok")
		srv.Close()
		if err != nil {
			t.Fatalf("ghCheckPushAccess: %v", err)
		}
		if got != want {
			t.Errorf("ghCheckPushAccess = %v, want %v", got, want)
		}
	}
}

func TestGhRequest_SetsAuthHeaderOnlyWhenTokenPresent(t *testing.T) {
	var gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	if _, err := ghRequest(context.Background(), http.MethodGet, srv.URL, "", nil); err != nil {
		t.Fatal(err)
	}
	if gotAuth != "" {
		t.Errorf("Authorization header set with empty token: %q", gotAuth)
	}

	if _, err := ghRequest(context.Background(), http.MethodGet, srv.URL, "sekret", nil); err != nil {
		t.Fatal(err)
	}
	if gotAuth != "Bearer sekret" {
		t.Errorf("Authorization header = %q, want Bearer sekret", gotAuth)
	}
}
