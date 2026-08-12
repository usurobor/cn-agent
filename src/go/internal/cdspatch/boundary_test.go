package cdspatch

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/usurobor/cnos/src/go/internal/cellwork"
)

// D3: Go must accept exactly the keys the closed CUE overlay accepts.
// encoding/json matches field names case-insensitively even with
// DisallowUnknownFields, so without an exact-key check `Fill` or a nested
// `Provider` would execute in Go while CUE rejects it.
func TestMixedCaseKeysRejected(t *testing.T) {
	f := Factory(skillTree(t, testSkills...))
	base := `{"fill":"cds.patch","cognition":{"provider":"fake","model":""},` +
		`"skills":["cnos.eng:eng/go"]}`
	if _, err := f(context.Background(), json.RawMessage(base)); err != nil {
		t.Fatalf("canonical declaration must construct: %v", err)
	}
	mixed := map[string]string{
		"seat tag":      strings.Replace(base, `"fill"`, `"Fill"`, 1),
		"top-level arg": strings.Replace(base, `"cognition"`, `"Cognition"`, 1),
		"nested arg":    strings.Replace(base, `"provider"`, `"Provider"`, 1),
		"skill list":    strings.Replace(base, `"skills"`, `"Skills"`, 1),
	}
	for name, decl := range mixed {
		t.Run(name, func(t *testing.T) {
			if _, err := f(context.Background(), json.RawMessage(decl)); err == nil {
				t.Fatalf("mixed-case %s must be rejected", name)
			}
		})
	}
}

// D4: an oversized diff is refused without ever buffering it whole.
func TestOversizedDiffIsRefused(t *testing.T) {
	repo, _ := testRepo(t)
	wt, release, err := cellwork.Materialize(context.Background(), repo, "HEAD")
	if err != nil {
		t.Fatalf("materialize: %v", err)
	}
	defer release()

	// Comfortably past the 1 MiB matter bound.
	big := strings.Repeat("every line of this file is evidence nobody asked for\n", 40000)
	if err := os.WriteFile(filepath.Join(wt.Dir, "huge.txt"), []byte(big), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := wt.Diff(context.Background()); err == nil {
		t.Fatal("an oversized diff must be refused")
	} else if !strings.Contains(err.Error(), "more than") {
		t.Fatalf("want a bound error, got %v", err)
	}
}
