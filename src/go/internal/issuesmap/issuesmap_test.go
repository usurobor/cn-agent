package issuesmap

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// findRec returns the record with number n, or fails.
func findRec(t *testing.T, recs []Record, n int) Record {
	t.Helper()
	for _, r := range recs {
		if r.N == n {
			return r
		}
	}
	t.Fatalf("record #%d not found", n)
	return Record{}
}

func TestToRecord_LabelParsingAndEffort(t *testing.T) {
	cases := []struct {
		name   string
		issue  ghIssue
		expect func(t *testing.T, r Record)
	}{
		{
			name: "estimated dispatch cell",
			issue: ghIssue{Number: 545, Title: "board", Labels: []ghLabel{
				{"kind/tooling"}, {"area/ci"}, {"area/cli"}, {"P1"},
				{"status:changes"}, {"dispatch:cell"}, {"protocol:cds"}, {"effort/XL"},
			}},
			expect: func(t *testing.T, r Record) {
				if r.Kind != "tooling" || r.Pri != "P1" || r.Status != "changes" {
					t.Errorf("bad core fields: %+v", r)
				}
				if !r.Dispatch || r.Protocol != "cds" {
					t.Errorf("dispatch/protocol wrong: %+v", r)
				}
				if r.Effort != "XL" || r.EffortWeight != 8 || r.Unestimated {
					t.Errorf("effort wrong: %+v", r)
				}
				// areas sorted
				if strings.Join(r.Areas, ",") != "ci,cli" {
					t.Errorf("areas not sorted: %v", r.Areas)
				}
			},
		},
		{
			name:  "unestimated is weight 1, flagged, not S",
			issue: ghIssue{Number: 1, Title: "x", Labels: []ghLabel{{"kind/bugfix"}, {"area/ci"}}},
			expect: func(t *testing.T, r Record) {
				if r.Effort != "" || r.EffortWeight != 1 || !r.Unestimated {
					t.Errorf("unestimated wrong: %+v", r)
				}
				if r.Pri != "none" {
					t.Errorf("default priority should be none: %+v", r)
				}
			},
		},
		{
			name:  "tracking flagged",
			issue: ghIssue{Number: 326, Title: "t", Labels: []ghLabel{{"kind/tracking"}, {"area/cdd"}}},
			expect: func(t *testing.T, r Record) {
				if !r.Tracking {
					t.Errorf("tracking not flagged: %+v", r)
				}
			},
		},
		{
			name:  "no labels: gaps, url synthesized",
			issue: ghIssue{Number: 547, Title: "gap"},
			expect: func(t *testing.T, r Record) {
				if r.Kind != "" || len(r.Areas) != 0 {
					t.Errorf("expected empty kind/areas: %+v", r)
				}
				if r.URL != "https://github.com/usurobor/cnos/issues/547" {
					t.Errorf("url not synthesized: %q", r.URL)
				}
				if r.Areas == nil {
					t.Errorf("areas must be non-nil for JSON []")
				}
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, toRecord(tc.issue, "usurobor/cnos"))
		})
	}
}

func TestEffortWeights(t *testing.T) {
	for lbl, want := range map[string]int{"S": 1, "M": 2, "L": 4, "XL": 8} {
		got := toRecord(ghIssue{Number: 1, Labels: []ghLabel{{"effort/" + lbl}}}, "r/r").EffortWeight
		if got != want {
			t.Errorf("effort/%s weight = %d, want %d", lbl, got, want)
		}
	}
}

// TestRun_Fixture drives the full generate path offline from testdata and
// asserts the outputs exist, are self-contained, and are deterministic.
func TestRun_Fixture(t *testing.T) {
	out := t.TempDir()
	var stdout, stderr bytes.Buffer
	err := Run(context.Background(),
		[]string{"--repo", "usurobor/cnos", "--fixture", "testdata/issues.json", "--out", out},
		nil, &stdout, &stderr)
	if err != nil {
		t.Fatalf("Run: %v\nstderr: %s", err, stderr.String())
	}

	// All artifacts written.
	for _, name := range []string{"index.html", "board-data.json", "README.md"} {
		if _, err := os.Stat(filepath.Join(out, name)); err != nil {
			t.Errorf("missing output %s: %v", name, err)
		}
	}

	// The PR (#999) must be excluded; 5 issues remain.
	data, _ := os.ReadFile(filepath.Join(out, "board-data.json"))
	var recs []Record
	if err := json.Unmarshal(data, &recs); err != nil {
		t.Fatalf("board-data.json not valid JSON: %v", err)
	}
	if len(recs) != 5 {
		t.Errorf("want 5 records (PR excluded), got %d", len(recs))
	}
	for _, r := range recs {
		if r.N == 999 {
			t.Errorf("pull request #999 leaked into the board")
		}
	}

	// Effort-aware data present; #545 estimated XL, #547 unestimated.
	if r := findRec(t, recs, 545); r.EffortWeight != 8 {
		t.Errorf("#545 effort_weight = %d, want 8", r.EffortWeight)
	}
	if r := findRec(t, recs, 547); !r.Unestimated {
		t.Errorf("#547 should be unestimated")
	}

	// index.html is self-contained (D3 inlined, data spliced, no placeholder left).
	idx, _ := os.ReadFile(filepath.Join(out, "index.html"))
	sidx := string(idx)
	if strings.Contains(sidx, "__BOARD__") || strings.Contains(sidx, "__KINDS__") {
		t.Errorf("index.html still contains an unspliced placeholder")
	}
	if !strings.Contains(sidx, "const DATA=[") {
		t.Errorf("index.html missing spliced board data")
	}
	// Self-contained = no external *resource loads*. Anchor hrefs to issue
	// URLs are legitimate data links, so only script/style/font loads count.
	for _, bad := range []string{"src=\"http", "<link ", "@import", "url(http"} {
		if strings.Contains(sidx, bad) {
			t.Errorf("index.html loads an external resource (%q) — not self-contained", bad)
		}
	}
	// Taxonomy gap for #547 reported.
	if !strings.Contains(stderr.String(), "#547") {
		t.Errorf("expected taxonomy gap report for #547, stderr: %s", stderr.String())
	}

	// Determinism: a second run into a fresh dir yields byte-identical files.
	out2 := t.TempDir()
	var so2, se2 bytes.Buffer
	if err := Run(context.Background(),
		[]string{"--repo", "usurobor/cnos", "--fixture", "testdata/issues.json", "--out", out2},
		nil, &so2, &se2); err != nil {
		t.Fatalf("second Run: %v", err)
	}
	for _, name := range []string{"index.html", "board-data.json"} {
		a, _ := os.ReadFile(filepath.Join(out, name))
		b, _ := os.ReadFile(filepath.Join(out2, name))
		if !bytes.Equal(a, b) {
			t.Errorf("%s not deterministic across runs", name)
		}
	}
}

func TestRun_Stdin(t *testing.T) {
	out := t.TempDir()
	in := strings.NewReader(`[{"number":10,"title":"s","labels":[{"name":"kind/feature"},{"name":"area/ui"}]}]`)
	var so, se bytes.Buffer
	if err := Run(context.Background(), []string{"--repo", "o/r", "--stdin", "--out", out}, in, &so, &se); err != nil {
		t.Fatalf("stdin run: %v", err)
	}
	data, _ := os.ReadFile(filepath.Join(out, "board-data.json"))
	if !strings.Contains(string(data), "\"n\": 10") {
		t.Errorf("stdin issue not rendered: %s", data)
	}
}
