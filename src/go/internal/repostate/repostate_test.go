package repostate

import (
	"encoding/json"
	"strings"
	"testing"
)

func exampleState() *RepoState {
	return &RepoState{
		Schema:  Schema,
		Profile: "cds",
		Source: Source{
			Channel: "stable",
			Release: "3.82.6",
			Index:   "github:usurobor/cnos/releases/3.82.6/index.json",
		},
		ManagedFiles: []ManagedFile{
			{Path: ".cn/deps.lock.json", Kind: KindLockfile, ID: "deps-lockfile", SHA256: strings.Repeat("b", 64)},
			{Path: ".cn/deps.json", Kind: KindManifest, ID: "deps-manifest", SHA256: strings.Repeat("a", 64)},
			{
				Path: ".github/workflows/cnos-cds-dispatch.yml", Kind: KindWorkflow,
				ID: "cnos-cds-dispatch", SHA256: strings.Repeat("d", 64),
				Tier: TierAgent, Renderer: "cnos.core/install-wake",
				RendererPackage: "cnos.core", RendererVersionSource: "lock",
				Agent: "sigma", WorkflowPatSecret: "SIGMA_WORKFLOW_PAT",
				BotName: "sigma@cnos.cn-sigma.cnos", BotID: "41898282",
			},
		},
		ManagedDirs: []ManagedDir{
			{Path: ".cn/vendor/packages/cnos.cds", Package: "cnos.cds"},
			{Path: ".cn/vendor/packages/cnos.core", Package: "cnos.core"},
		},
		ExternalExpectations: ExternalExpectations{
			Labels: LabelExpectations{Mode: LabelModeEnsure, Source: "cnos.core/labels.json", DeleteOnUninstall: false},
		},
	}
}

func TestMarshal_SortsManagedFilesAndDirsByPath(t *testing.T) {
	s := exampleState()
	data, err := s.Marshal()
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("re-parse: %v", err)
	}
	files := got["managed_files"].([]any)
	var paths []string
	for _, f := range files {
		paths = append(paths, f.(map[string]any)["path"].(string))
	}
	want := []string{".cn/deps.json", ".cn/deps.lock.json", ".github/workflows/cnos-cds-dispatch.yml"}
	for i, w := range want {
		if paths[i] != w {
			t.Errorf("managed_files[%d].path = %q, want %q (full order: %v)", i, paths[i], w, paths)
		}
	}

	dirs := got["managed_dirs"].([]any)
	firstDirPath := dirs[0].(map[string]any)["path"].(string)
	if firstDirPath != ".cn/vendor/packages/cnos.cds" {
		t.Errorf("managed_dirs[0].path = %q, want sorted-first %q", firstDirPath, ".cn/vendor/packages/cnos.cds")
	}
}

func TestMarshal_DeterministicAcrossBuildOrder(t *testing.T) {
	a := exampleState()
	b := exampleState()
	// Reverse b's slices before marshaling — Sort() inside Marshal must
	// make the output byte-identical regardless of build/insertion order,
	// which is what makes a same-input `cn repo install` rerun produce no
	// diff (design doc's idempotence contract).
	b.ManagedFiles[0], b.ManagedFiles[1] = b.ManagedFiles[1], b.ManagedFiles[0]
	b.ManagedDirs[0], b.ManagedDirs[1] = b.ManagedDirs[1], b.ManagedDirs[0]

	dataA, err := a.Marshal()
	if err != nil {
		t.Fatalf("Marshal a: %v", err)
	}
	dataB, err := b.Marshal()
	if err != nil {
		t.Fatalf("Marshal b: %v", err)
	}
	if string(dataA) != string(dataB) {
		t.Errorf("Marshal not deterministic across build order:\na=%s\nb=%s", dataA, dataB)
	}
}

func TestMarshal_TrailingNewline(t *testing.T) {
	data, err := exampleState().Marshal()
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if len(data) == 0 || data[len(data)-1] != '\n' {
		t.Errorf("Marshal output does not end with a newline")
	}
}

func TestParse_RoundTrip(t *testing.T) {
	original := exampleState()
	data, err := original.Marshal()
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	parsed, err := Parse(data)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	data2, err := parsed.Marshal()
	if err != nil {
		t.Fatalf("Marshal (round-trip): %v", err)
	}
	if string(data) != string(data2) {
		t.Errorf("round-trip not stable:\noriginal=%s\nround-tripped=%s", data, data2)
	}
}

func TestParse_MalformedJSON(t *testing.T) {
	if _, err := Parse([]byte("{not json")); err == nil {
		t.Error("Parse of malformed JSON: got nil error, want non-nil")
	}
}

func TestFindManagedFile(t *testing.T) {
	s := exampleState()
	if f := s.FindManagedFile(".cn/deps.json"); f == nil || f.Kind != KindManifest {
		t.Errorf("FindManagedFile(.cn/deps.json) = %+v, want manifest entry", f)
	}
	if f := s.FindManagedFile("nonexistent"); f != nil {
		t.Errorf("FindManagedFile(nonexistent) = %+v, want nil", f)
	}
}

func TestFindManagedDir(t *testing.T) {
	s := exampleState()
	if d := s.FindManagedDir("cnos.core"); d == nil || d.Path != ".cn/vendor/packages/cnos.core" {
		t.Errorf("FindManagedDir(cnos.core) = %+v, want vendor path entry", d)
	}
	if d := s.FindManagedDir("nonexistent"); d != nil {
		t.Errorf("FindManagedDir(nonexistent) = %+v, want nil", d)
	}
}

func TestManagedFile_NonWorkflowOmitsRenderContractFields(t *testing.T) {
	s := exampleState()
	data, err := s.Marshal()
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("re-parse: %v", err)
	}
	for _, f := range got["managed_files"].([]any) {
		m := f.(map[string]any)
		if m["kind"] != KindWorkflow {
			for _, renderField := range []string{"tier", "renderer", "renderer_package", "renderer_version_source", "agent", "workflow_pat_secret", "bot_name", "bot_id"} {
				if _, present := m[renderField]; present {
					t.Errorf("non-workflow managed_files entry %q carries render-contract field %q (should be omitted)", m["path"], renderField)
				}
			}
		}
	}
}
