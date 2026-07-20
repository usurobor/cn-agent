package cddverify

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestUnreleasedCloseoutsAreCanonicalPostMergeState(t *testing.T) {
	root := t.TempDir()
	cycleDir := filepath.Join(root, ".cdd", "unreleased", "669")
	if err := os.MkdirAll(cycleDir, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"self-coherence.md", "beta-review.md", "alpha-closeout.md", "beta-closeout.md", "gamma-closeout.md"} {
		if err := os.WriteFile(filepath.Join(cycleDir, name), []byte("# fixture\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	var out bytes.Buffer
	l := newLedgerRun(root, "", &out)
	l.checkUnreleasedTriadicArtifacts(cycleDir, "669")
	got := out.String()
	if strings.Contains(got, "may indicate stale cycle") {
		t.Fatalf("post-merge closeouts under unreleased were misclassified as stale:\n%s", got)
	}
	for _, name := range []string{"alpha-closeout.md", "beta-closeout.md", "gamma-closeout.md"} {
		if !strings.Contains(got, "✅ "+name+" (issue #669) — present in canonical post-merge/release-pending location") {
			t.Errorf("%s was not accepted in the post-merge location:\n%s", name, got)
		}
	}
}

// Table-driven tests for the helpers and rules. eng/go §2.10 requires
// table-driven + subtests; happy + degraded paths.

func TestIsFilesystemPath(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want bool
	}{
		{"empty", "", false},
		{"absolute", "/tmp/foo.yaml", true},
		{"relative dot", "./foo/bar.yaml", true},
		{"relative slash with ext", "schemas/cds/fixtures/r.yaml", true},
		{"bare name", "foo", false},
		{"git_diff placeholder", "git_diff(origin/main..HEAD)", false},
		{"sha256 placeholder", "sha256:abc123", false},
		{"issue placeholder", "issue:#389", false},
		{"path with at-ref", "schemas/cds/fixtures/r.yaml@sha256:abc", false}, // @sha256: matches
		{"path with merge_sha placeholder", "foo/bar@merge_sha_placeholder", false},
		{"path with draft", "foo/bar@draft", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isFilesystemPath(tc.in); got != tc.want {
				t.Errorf("isFilesystemPath(%q) = %v, want %v", tc.in, got, tc.want)
			}
		})
	}
}

func TestStripAtRef(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"foo/bar.yaml", "foo/bar.yaml"},
		{"foo/bar.yaml@sha256:abc", "foo/bar.yaml"},
		{"@only-ref", ""},
		{"", ""},
	}
	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			if got := stripAtRef(tc.in); got != tc.want {
				t.Errorf("stripAtRef(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestParseISO(t *testing.T) {
	cases := []struct {
		in       string
		wantOK   bool
		approxOK string // optional verbose check
	}{
		{"2026-05-21T10:00:00Z", true, ""},
		{"2026-05-21T10:00:00+00:00", true, ""},
		{"2026-05-21T10:00:00.123Z", true, ""},
		{"not a date", false, ""},
		{"", false, ""},
	}
	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			ts, ok := parseISO(tc.in)
			if ok != tc.wantOK {
				t.Errorf("parseISO(%q) ok = %v, want %v", tc.in, ok, tc.wantOK)
			}
			if ok && ts.IsZero() {
				t.Errorf("parseISO(%q) returned ok=true but zero time", tc.in)
			}
		})
	}
}

func TestKnownProtocolsSorted(t *testing.T) {
	got := KnownProtocols()
	want := []string{
		"cnos.cdd.cdr.receipt.v1",
		"cnos.cdd.cds.receipt.v1",
		"cnos.cdd.generic.receipt.v1",
	}
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d (%v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("[%d] got %q, want %q", i, got[i], want[i])
		}
	}
}

func TestParseReceiptJSON_HappyPath(t *testing.T) {
	data := []byte(`{
        "protocol_id": "cnos.cdd.generic.receipt.v1",
        "validation": {"verdict": "PASS"},
        "evidence_refs": {"evidence_root": "x/"}
    }`)
	r, err := ParseReceiptJSON(data)
	if err != nil {
		t.Fatalf("ParseReceiptJSON: %v", err)
	}
	if protocolID(r) != "cnos.cdd.generic.receipt.v1" {
		t.Errorf("protocol_id = %q", protocolID(r))
	}
	er, ok := evidenceRefs(r)
	if !ok {
		t.Fatal("evidence_refs missing")
	}
	if s, _ := asString(er, "evidence_root"); s != "x/" {
		t.Errorf("evidence_root = %q", s)
	}
}

func TestParseReceiptJSON_Malformed(t *testing.T) {
	_, err := ParseReceiptJSON([]byte("not json"))
	if err == nil {
		t.Fatal("expected error for malformed JSON")
	}
}

func TestRuleC3OverrideDoesNotRewrite(t *testing.T) {
	r := Receipt{
		"validation": map[string]any{"verdict": "FAIL"},
		"boundary_decision": map[string]any{
			"action": "override",
			"override": map[string]any{
				"original_validation_verdict": map[string]any{"verdict": "PASS"},
			},
		},
	}
	var fail []FailedPredicate
	ruleC3OverrideDoesNotRewrite(r, &fail)
	if len(fail) != 1 {
		t.Fatalf("expected 1 fail, got %d (%v)", len(fail), fail)
	}
	if fail[0].Predicate != "counterfeit.override_does_not_rewrite" {
		t.Errorf("predicate = %q", fail[0].Predicate)
	}
}

func TestRuleC3_NotOverride_NoFire(t *testing.T) {
	r := Receipt{
		"validation":        map[string]any{"verdict": "FAIL"},
		"boundary_decision": map[string]any{"action": "accept"},
	}
	var fail []FailedPredicate
	ruleC3OverrideDoesNotRewrite(r, &fail)
	if len(fail) != 0 {
		t.Errorf("expected 0 fails for non-override action, got %d", len(fail))
	}
}

func TestRuleC5_CDSMissingKeys(t *testing.T) {
	r := Receipt{
		"protocol_id":   "cnos.cdd.cds.receipt.v1",
		"evidence_refs": map[string]any{"evidence_root": "x/"},
	}
	var fail []FailedPredicate
	ruleC5ProtocolIDMismatch(r, &fail)
	if len(fail) != 1 {
		t.Fatalf("expected 1 fail, got %d (%v)", len(fail), fail)
	}
	if !strings.Contains(fail[0].Diagnostic, "missing required CDS evidence_refs keys") {
		t.Errorf("diagnostic = %q", fail[0].Diagnostic)
	}
}

func TestRuleC5_CDRMissingKeys(t *testing.T) {
	r := Receipt{
		"protocol_id":   "cnos.cdd.cdr.receipt.v1",
		"evidence_refs": map[string]any{"claim": "x"},
	}
	var fail []FailedPredicate
	ruleC5ProtocolIDMismatch(r, &fail)
	if len(fail) != 1 {
		t.Fatalf("expected 1 fail, got %d", len(fail))
	}
	if !strings.Contains(fail[0].Diagnostic, "missing required CDR evidence_refs keys") {
		t.Errorf("diagnostic = %q", fail[0].Diagnostic)
	}
}

func TestRuleC5_GenericProtocol_NoFire(t *testing.T) {
	r := Receipt{
		"protocol_id":   "cnos.cdd.generic.receipt.v1",
		"evidence_refs": map[string]any{},
	}
	var fail []FailedPredicate
	ruleC5ProtocolIDMismatch(r, &fail)
	if len(fail) != 0 {
		t.Errorf("expected no fails for generic protocol")
	}
}

func TestEmitVerdict_JSONShape(t *testing.T) {
	v := emitVerdict(ResultPASS, nil, nil, "/tmp/r.yaml", "", "")
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	// Required schema keys present.
	for _, k := range []string{`"result"`, `"failed_predicates"`, `"warnings"`, `"provenance"`} {
		if !strings.Contains(string(data), k) {
			t.Errorf("missing required key %s in %s", k, string(data))
		}
	}
	if !strings.Contains(string(data), `"validator_identity":"cnos.cdd.validate_receipt"`) {
		t.Errorf("validator_identity not pinned in %s", string(data))
	}
	// PASS verdict has empty failed_predicates array (not null/missing).
	if !strings.Contains(string(data), `"failed_predicates":[]`) {
		t.Errorf("failed_predicates should be [] for PASS, got %s", string(data))
	}
}

func TestSameAny(t *testing.T) {
	cases := []struct {
		a, b any
		want bool
	}{
		{"PASS", "PASS", true},
		{"PASS", "FAIL", false},
		{nil, nil, true},
		{nil, "x", false},
		{1, 1, true},
		{1, "1", true}, // fmt.Sprint coercion — acceptable for verdict-string compare
	}
	for i, tc := range cases {
		got := sameAny(tc.a, tc.b)
		if got != tc.want {
			t.Errorf("[%d] sameAny(%v, %v) = %v, want %v", i, tc.a, tc.b, got, tc.want)
		}
	}
}

func TestSortedKeys(t *testing.T) {
	m := map[string]any{"b": 1, "a": 1, "c": 1}
	got := sortedKeys(m)
	want := []string{"a", "b", "c"}
	for i, k := range want {
		if got[i] != k {
			t.Errorf("[%d] got %q, want %q", i, got[i], k)
		}
	}
}

func TestParseArgs_Receipt(t *testing.T) {
	a, err := ParseArgs([]string{"--receipt", "/tmp/r.yaml", "--json"})
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if a.Receipt != "/tmp/r.yaml" {
		t.Errorf("Receipt = %q", a.Receipt)
	}
	if !a.JSON {
		t.Error("JSON should be true")
	}
}

func TestParseArgs_LedgerVersion(t *testing.T) {
	a, err := ParseArgs([]string{"--version", "3.81.0", "--cycle", "392"})
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if a.Version != "3.81.0" {
		t.Errorf("Version = %q", a.Version)
	}
	if a.Cycle != "392" {
		t.Errorf("Cycle = %q", a.Cycle)
	}
}

func TestParseArgs_UnknownFlag(t *testing.T) {
	_, err := ParseArgs([]string{"--nonexistent"})
	if err == nil {
		t.Fatal("expected error for unknown flag")
	}
	if !strings.Contains(err.Error(), "unknown option") {
		t.Errorf("error = %q", err.Error())
	}
}

func TestParseArgs_MissingValue(t *testing.T) {
	_, err := ParseArgs([]string{"--receipt"})
	if err == nil {
		t.Fatal("expected error for flag without value")
	}
}

func TestParseArgs_Help(t *testing.T) {
	a, err := ParseArgs([]string{"--help"})
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if !a.Help {
		t.Error("Help should be true")
	}
}

func TestRenderProse_PASS(t *testing.T) {
	v := ValidationVerdict{
		Result: ResultPASS,
		Provenance: Provenance{
			ValidatorIdentity: ValidatorIdentity,
			ValidatorVersion:  ValidatorVersion,
			CheckedAt:         time.Now().UTC().Format(time.RFC3339),
		},
	}
	out := RenderProse(v, "/tmp/r.yaml")
	if !strings.Contains(out, "PASS") {
		t.Errorf("expected PASS in output: %q", out)
	}
	if !strings.Contains(out, "cnos.cdd.validate_receipt") {
		t.Errorf("expected validator identity in output: %q", out)
	}
}

func TestRenderProse_FAIL(t *testing.T) {
	v := ValidationVerdict{
		Result: ResultFAIL,
		FailedPredicates: []FailedPredicate{{
			Predicate:   "counterfeit.actor_separation",
			EvidenceRef: "alpha_closeout / beta_closeout",
			Diagnostic:  "same actor",
		}},
	}
	out := RenderProse(v, "r.yaml")
	if !strings.Contains(out, "FAIL") {
		t.Errorf("expected FAIL in output")
	}
	if !strings.Contains(out, "counterfeit.actor_separation") {
		t.Errorf("expected predicate in output")
	}
	if !strings.Contains(out, "alpha_closeout / beta_closeout") {
		t.Errorf("expected evidence_ref in output")
	}
}

func TestIsLegacyVersion(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"3.76.0", true},
		{"3.77.0", false},
		{"3.81.0", false},
		{"2.0.0", true},
		{"not-a-version", false},
	}
	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			if got := isLegacyVersion(tc.in); got != tc.want {
				t.Errorf("isLegacyVersion(%q) = %v, want %v", tc.in, got, tc.want)
			}
		})
	}
}

func TestSectionPresent(t *testing.T) {
	text := "# Title\n\n## Gap\n\nbody\n\n## ACs\n"
	if !sectionPresent(text, "## Gap") {
		t.Error("Gap not detected")
	}
	if !sectionPresent(text, "## ACs") {
		t.Error("ACs not detected")
	}
	if sectionPresent(text, "## Missing") {
		t.Error("Missing should not be detected")
	}
}
