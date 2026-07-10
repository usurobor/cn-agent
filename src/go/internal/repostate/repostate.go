// Package repostate defines the pure types for `.cn/repo.state.json`
// (schema cn.repo.state.v1) — the managed-surface ownership ledger
// cnos#656 (Phase 1 of the cnos#655 `cn repo` lifecycle wave) introduces.
//
// Discipline: no os, no net, no io — file reading/writing lives in
// internal/repoinstall (write) and internal/repostatus (read), mirroring
// internal/pkg's own split (pkg.go:9-11). Field names and JSON tags match
// schemas/repo_state.cue exactly; that CUE file is the enforcement gate
// (via `cue vet`, scripts/ci/validate-repo-state.sh), not this package —
// these types exist to read/write the shape CUE already validates.
package repostate

import (
	"encoding/json"
	"fmt"
	"sort"
)

// Schema is the fixed schema identifier every ledger declares.
const Schema = "cn.repo.state.v1"

// Managed-file kinds. Mirrors schemas/repo_state.cue's #ManagedFile.kind enum.
const (
	KindManifest  = "manifest"
	KindLockfile  = "lockfile"
	KindWorkflow  = "workflow"
	KindGitignore = "gitignore"
	KindOther     = "other"
)

// Workflow dispatch tiers. Mirrors the renderer's own --tier values
// (cn-install-wake:63-72) and schemas/repo_state.cue's #ManagedFile.tier enum.
const (
	TierAgent  = "agent"
	TierEngine = "engine"
)

// Label-expectation modes. Mirrors schemas/repo_state.cue's
// #RepoState.external_expectations.labels.mode enum.
const (
	LabelModeEnsure = "ensure"
	LabelModeIgnore = "ignore"
)

// Source records where the installed release came from.
type Source struct {
	Channel string `json:"channel"`
	Release string `json:"release"`
	Index   string `json:"index"`
}

// ManagedFile is one entry in RepoState.ManagedFiles. Path/Kind/ID/SHA256
// are required for every entry (P1). The render-contract fields
// (Tier/Renderer/RendererPackage/RendererVersionSource/Agent/
// WorkflowPatSecret/BotName/BotID) are populated only for Kind ==
// KindWorkflow (P2) — `omitempty` keeps non-workflow entries free of empty
// render-contract keys in the serialized JSON, matching the CUE schema's
// conditional `if kind == "workflow"` block.
type ManagedFile struct {
	Path   string `json:"path"`
	Kind   string `json:"kind"`
	ID     string `json:"id"`
	SHA256 string `json:"sha256"`

	Tier                  string `json:"tier,omitempty"`
	Renderer              string `json:"renderer,omitempty"`
	RendererPackage       string `json:"renderer_package,omitempty"`
	RendererVersionSource string `json:"renderer_version_source,omitempty"`
	Agent                 string `json:"agent,omitempty"`
	WorkflowPatSecret     string `json:"workflow_pat_secret,omitempty"`
	BotName               string `json:"bot_name,omitempty"`
	BotID                 string `json:"bot_id,omitempty"`
}

// ManagedDir maps a materialized vendor directory to the package that owns
// it. No version/sha256 here by design (design doc A1): package identity
// is looked up from .cn/deps.lock.json by Package at read time.
type ManagedDir struct {
	Path    string `json:"path"`
	Package string `json:"package"`
}

// LabelExpectations types RepoState.ExternalExpectations.Labels.
type LabelExpectations struct {
	Mode               string `json:"mode"`
	Source             string `json:"source"`
	DeleteOnUninstall  bool   `json:"delete_on_uninstall"`
}

// ExternalExpectations types RepoState.ExternalExpectations.
type ExternalExpectations struct {
	Labels LabelExpectations `json:"labels"`
}

// RepoState is the parsed/constructed .cn/repo.state.json. Deterministic,
// timestamp-free (P1/P3) — ManagedFiles/ManagedDirs are sorted by Path
// before serialization (see Sort) so a same-input rerun produces byte-
// identical output.
type RepoState struct {
	Schema               string               `json:"schema"`
	Profile              string               `json:"profile"`
	Source               Source               `json:"source"`
	ManagedFiles         []ManagedFile        `json:"managed_files"`
	ManagedDirs          []ManagedDir         `json:"managed_dirs"`
	ExternalExpectations ExternalExpectations `json:"external_expectations"`
}

// Sort orders ManagedFiles and ManagedDirs by Path in place, so two
// RepoState values built from the same inputs (regardless of map/slice
// build order upstream) marshal to identical bytes.
func (s *RepoState) Sort() {
	sort.Slice(s.ManagedFiles, func(i, j int) bool {
		return s.ManagedFiles[i].Path < s.ManagedFiles[j].Path
	})
	sort.Slice(s.ManagedDirs, func(i, j int) bool {
		return s.ManagedDirs[i].Path < s.ManagedDirs[j].Path
	})
}

// Marshal renders s as deterministic, indented JSON with a trailing
// newline — the same convention internal/repoinstall's writeManifest uses
// for .cn/deps.json (repoinstall.go:580-590). Calls Sort first so callers
// never need to remember to do so.
func (s *RepoState) Marshal() ([]byte, error) {
	s.Sort()
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal repo state: %w", err)
	}
	return append(data, '\n'), nil
}

// Parse parses a .cn/repo.state.json from raw JSON bytes. Pure — no IO.
func Parse(data []byte) (*RepoState, error) {
	var s RepoState
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parse repo state: %w", err)
	}
	return &s, nil
}

// FindManagedFile returns the ManagedFile entry with the given path, or
// nil if absent.
func (s *RepoState) FindManagedFile(path string) *ManagedFile {
	for i := range s.ManagedFiles {
		if s.ManagedFiles[i].Path == path {
			return &s.ManagedFiles[i]
		}
	}
	return nil
}

// FindManagedDir returns the ManagedDir entry with the given package name,
// or nil if absent.
func (s *RepoState) FindManagedDir(pkgName string) *ManagedDir {
	for i := range s.ManagedDirs {
		if s.ManagedDirs[i].Package == pkgName {
			return &s.ManagedDirs[i]
		}
	}
	return nil
}
