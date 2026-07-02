// Package issuesmap implements `cn issues map`: it reads the repository's
// open issues (live via the GitHub API, or from a fixture/stdin for offline
// use), derives a taxonomy-aware board record for each, and renders the
// self-contained interactive issue-board views into an output directory.
//
// The Go command owns the data plane (fetch, parse, effort computation,
// emit); the browser owns rendering. The HTML templates embed D3 and are
// spliced with the board JSON at generate time, so the output is a
// self-contained page with no runtime network dependency.
//
// Design authority: cnos#545. The dispatch wrapper lives in
// internal/cli/cmd_issues_map.go and holds no domain logic, per the
// dispatch boundary (INVARIANTS.md T-002, eng/go §2.18).
package issuesmap

import (
	"context"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

//go:embed templates/index.template.html templates/pivot.template.html
var templates embed.FS

// effortWeight maps the effort/* label suffix to the dashboard area weight
// (gentle doubling — see docs/development/issues/TAXONOMY.md §Effort labels).
var effortWeight = map[string]int{"S": 1, "M": 2, "L": 4, "XL": 8}

// kindOrder is the canonical primary-kind ordering (TAXONOMY.md §Primary
// kind), plus "none" and "" (unlabeled) as trailing buckets so the board's
// kind grouping is deterministic.
var kindOrder = []string{
	"bugfix", "cleanup", "process", "feature", "tooling", "doctrine",
	"audit", "tracking", "research", "skill", "spike", "chore", "none", "",
}

const defaultRepo = "usurobor/cnos"
const defaultOut = "docs/development/board"

// ghLabel and ghIssue mirror the subset of the GitHub issue schema the
// board needs. The same shape is used for the live API response and for
// fixtures, so one parser serves both paths.
type ghLabel struct {
	Name string `json:"name"`
}

type ghIssue struct {
	Number      int              `json:"number"`
	Title       string           `json:"title"`
	HTMLURL     string           `json:"html_url"`
	State       string           `json:"state"`
	Labels      []ghLabel        `json:"labels"`
	PullRequest *json.RawMessage `json:"pull_request,omitempty"`
}

// Record is one board cell: an issue projected onto the taxonomy. Field
// names are the short keys the browser templates already read; effort_weight
// / unestimated / tracking back the effort-aware encoding (cnos#552).
type Record struct {
	N            int      `json:"n"`
	T            string   `json:"t"`
	Kind         string   `json:"kind"`
	Areas        []string `json:"areas"`
	Pri          string   `json:"pri"`
	Status       string   `json:"status"`
	Dispatch     bool     `json:"dispatch"`
	Protocol     string   `json:"protocol"`
	Effort       string   `json:"effort"`
	EffortWeight int      `json:"effort_weight"`
	Unestimated  bool     `json:"unestimated"`
	Tracking     bool     `json:"tracking"`
	URL          string   `json:"url"`
}

// toRecord projects a GitHub issue onto a board Record using the label
// taxonomy. Effort is derived from the effort/* label; a missing effort
// label yields unestimated (nominal weight 1), never effort/S.
func toRecord(is ghIssue, repo string) Record {
	r := Record{N: is.Number, T: is.Title, Areas: []string{}, Pri: "none", URL: is.HTMLURL}
	if r.URL == "" {
		if repo == "" {
			repo = defaultRepo
		}
		r.URL = fmt.Sprintf("https://github.com/%s/issues/%d", repo, is.Number)
	}
	for _, l := range is.Labels {
		n := l.Name
		switch {
		case strings.HasPrefix(n, "kind/"):
			if r.Kind == "" {
				r.Kind = n[len("kind/"):]
			}
		case strings.HasPrefix(n, "area/"):
			r.Areas = append(r.Areas, n[len("area/"):])
		case n == "P0" || n == "P1" || n == "P2" || n == "P3":
			r.Pri = n
		case strings.HasPrefix(n, "status:"):
			r.Status = n[len("status:"):]
		case strings.HasPrefix(n, "protocol:"):
			r.Protocol = n[len("protocol:"):]
		case n == "dispatch:cell":
			r.Dispatch = true
		case strings.HasPrefix(n, "effort/"):
			r.Effort = n[len("effort/"):]
		}
	}
	sort.Strings(r.Areas)
	if w, ok := effortWeight[r.Effort]; ok {
		r.EffortWeight = w
	} else {
		// Unestimated: nominal weight 1, flagged — not treated as small.
		r.Effort = ""
		r.EffortWeight = 1
		r.Unestimated = true
	}
	r.Tracking = r.Kind == "tracking"
	return r
}

// Run is the entry point for `cn issues map`. It returns a non-nil error on
// failure; taxonomy gaps are reported to stderr but are non-fatal.
func Run(ctx context.Context, args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("issues map", flag.ContinueOnError)
	fs.SetOutput(stderr)
	repo := fs.String("repo", "", "target repository as owner/name (default: $GITHUB_REPOSITORY, then "+defaultRepo+")")
	out := fs.String("out", defaultOut, "output directory for the generated board")
	fixture := fs.String("fixture", "", "read issues JSON from this file instead of the GitHub API (offline)")
	useStdin := fs.Bool("stdin", false, "read issues JSON from stdin instead of the GitHub API (offline)")
	token := fs.String("token", "", "GitHub token (default: $GITHUB_TOKEN, then $GH_TOKEN)")
	fs.Usage = func() {
		fmt.Fprintln(stderr, "Usage: cn issues map [--repo owner/name] [--out dir] [--fixture path | --stdin]")
		fmt.Fprintln(stderr, "\nGenerate the interactive issue-board (Voronoi + Pivot views) from the")
		fmt.Fprintln(stderr, "repository's open issues, colored by priority and taxonomy-aware.")
		fmt.Fprintln(stderr)
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *repo == "" {
		*repo = os.Getenv("GITHUB_REPOSITORY")
	}
	if *repo == "" {
		*repo = defaultRepo
	}

	// Load issues from the chosen source.
	var issues []ghIssue
	switch {
	case *fixture != "":
		data, err := os.ReadFile(*fixture)
		if err != nil {
			return fmt.Errorf("read fixture: %w", err)
		}
		if err := json.Unmarshal(data, &issues); err != nil {
			return fmt.Errorf("parse fixture JSON: %w", err)
		}
	case *useStdin:
		data, err := io.ReadAll(stdin)
		if err != nil {
			return fmt.Errorf("read stdin: %w", err)
		}
		if err := json.Unmarshal(data, &issues); err != nil {
			return fmt.Errorf("parse stdin JSON: %w", err)
		}
	default:
		tk := *token
		if tk == "" {
			tk = os.Getenv("GITHUB_TOKEN")
		}
		if tk == "" {
			tk = os.Getenv("GH_TOKEN")
		}
		var err error
		issues, err = fetchIssues(ctx, *repo, tk)
		if err != nil {
			return err
		}
	}

	// Project to records, skipping pull requests.
	recs := make([]Record, 0, len(issues))
	for _, is := range issues {
		if is.PullRequest != nil {
			continue
		}
		recs = append(recs, toRecord(is, *repo))
	}
	sort.Slice(recs, func(i, j int) bool { return recs[i].N < recs[j].N })

	reportGaps(recs, stderr)

	index, pivot, err := render(recs)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(*out, 0o755); err != nil {
		return fmt.Errorf("create out dir: %w", err)
	}
	data, err := json.MarshalIndent(recs, "", "  ")
	if err != nil {
		return err
	}
	files := map[string]string{
		"index.html":      index,
		"pivot.html":      pivot,
		"board-data.json": string(data) + "\n",
		"README.md":       readme(*repo),
	}
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(*out, name), []byte(body), 0o644); err != nil {
			return fmt.Errorf("write %s: %w", name, err)
		}
	}

	fmt.Fprintf(stdout, "cn issues map: wrote %d open issues to %s (index.html, pivot.html, board-data.json, README.md)\n", len(recs), *out)
	return nil
}

// render splices the board records and kind ordering into the embedded
// templates. json.Marshal HTML-escapes <, >, & so the data is safe to embed
// directly inside a <script> element.
func render(recs []Record) (index, pivot string, err error) {
	board, err := json.Marshal(recs)
	if err != nil {
		return "", "", err
	}
	kinds, err := json.Marshal(kindOrder)
	if err != nil {
		return "", "", err
	}
	idxT, err := templates.ReadFile("templates/index.template.html")
	if err != nil {
		return "", "", err
	}
	pvT, err := templates.ReadFile("templates/pivot.template.html")
	if err != nil {
		return "", "", err
	}
	return splice(string(idxT), board, kinds), splice(string(pvT), board, kinds), nil
}

func splice(tmpl string, board, kinds []byte) string {
	tmpl = strings.Replace(tmpl, "__BOARD__", string(board), 1)
	tmpl = strings.Replace(tmpl, "__KINDS__", string(kinds), 1)
	return tmpl
}

// reportGaps prints taxonomy gaps (issues missing a primary kind or any
// area) to stderr. It is advisory — the board still renders these issues.
func reportGaps(recs []Record, stderr io.Writer) {
	var noKind, noArea []int
	for _, r := range recs {
		if r.Kind == "" {
			noKind = append(noKind, r.N)
		}
		if len(r.Areas) == 0 {
			noArea = append(noArea, r.N)
		}
	}
	if len(noKind) == 0 && len(noArea) == 0 {
		return
	}
	fmt.Fprintln(stderr, "cn issues map: taxonomy gaps (advisory):")
	if len(noKind) > 0 {
		fmt.Fprintf(stderr, "  missing kind/*: %s\n", joinInts(noKind))
	}
	if len(noArea) > 0 {
		fmt.Fprintf(stderr, "  missing area/*: %s\n", joinInts(noArea))
	}
}

func joinInts(ns []int) string {
	parts := make([]string, len(ns))
	for i, n := range ns {
		parts[i] = "#" + fmt.Sprint(n)
	}
	return strings.Join(parts, " ")
}

func readme(repo string) string {
	base := "https://raw.githack.com/" + repo + "/main/" + defaultOut
	return "# Issue board maps\n\n" +
		"Interactive, self-contained visualizations of the `" + repo + "` open-issue\n" +
		"board, generated by **`cn issues map`** (Go) and labeled per\n" +
		"[the issue taxonomy](../issues/TAXONOMY.md).\n\n" +
		"> **Generated file — do not edit by hand.** Regenerate with\n" +
		"> `cn issues map --repo " + repo + " --out " + defaultOut + "`, or let the\n" +
		"> `board-map` GitHub Action refresh it on issue changes.\n\n" +
		"- [`index.html`](index.html) — **Voronoi tessellation**: organic cells for\n" +
		"  `kind → area → issue`, colored by priority, with facet filters (kind / area /\n" +
		"  priority / status / dispatchable / unestimated / search). Click a kind label to\n" +
		"  zoom, a cell to open the issue; light/dark.\n" +
		"- [`pivot.html`](pivot.html) — **Pivot-feel deep-zoom**: a continuous card wall,\n" +
		"  whole board → single issue card, driven by `d3.zoom()`; the same facet filters\n" +
		"  rearrange the wall.\n" +
		"- `board-data.json` — the embedded board data as a standalone file (the same\n" +
		"  records spliced into the HTML).\n\n" +
		"Each HTML file is fully self-contained (D3 inlined) — open it directly, no server\n" +
		"needed. For a live preview URL without hosting setup:\n\n" +
		"- Voronoi: " + base + "/index.html\n" +
		"- Pivot: " + base + "/pivot.html\n\n" +
		"## Encoding\n\n" +
		"- **Color = priority** (P0 → P3; `none` neutral).\n" +
		"- **Cell area = issue count** in v1; the effort-weighted area encoding\n" +
		"  (`area = effort`, per the `effort/*` labels) is tracked in the follow-up.\n" +
		"- **Data is effort-aware now**: every record carries `effort`, `effort_weight`,\n" +
		"  and `unestimated`. Unestimated issues render with a **dashed cell** (Voronoi)\n" +
		"  / an `unestimated` tag (Pivot) and are selectable via the **unestimated only**\n" +
		"  filter — a visible gap, never treated as small.\n" +
		"- **`kind/tracking`** issues are containers and are excluded from the effort\n" +
		"  rollup (the `effort Σ` stat).\n\n" +
		"## Source of truth\n\n" +
		"The generator is the Go command `cn issues map` (`src/go/internal/issuesmap/`).\n" +
		"There is no standalone script generator on the production path; the browser only\n" +
		"renders the embedded data.\n"
}
