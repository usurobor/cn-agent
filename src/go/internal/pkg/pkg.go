// Package pkg defines the pure types for cnos package management:
// manifest, lockfile, package index, and their JSON serialization.
//
// This is the Go equivalent of src/ocaml/lib/cn_package.ml — no IO, only
// types + JSON + pure helpers. The types and field names match the
// OCaml definitions so Go and OCaml produce structurally identical
// JSON for the same data.
//
// Discipline: this package imports only encoding/json and fmt.
// No os, no net, no io — file reading lives in internal/restore/,
// mirroring the OCaml split (src/ocaml/lib/ = pure, src/ocaml/cmd/ = IO).
package pkg

import (
	"encoding/json"
	"fmt"
	"slices"
)

// --- Manifest (deps.json) ---

// ManifestDep is one entry in the operator's desired package list.
type ManifestDep struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Manifest is the parsed .cn/deps.json.
type Manifest struct {
	Schema   string        `json:"schema"`
	Profile  string        `json:"profile"`
	Packages []ManifestDep `json:"packages"`
}

// --- Lockfile (deps.lock.json) ---

// LockedDep pins a package to a specific version + tarball SHA-256.
type LockedDep struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	SHA256  string `json:"sha256"`
}

// Lockfile is the parsed .cn/deps.lock.json.
type Lockfile struct {
	Schema   string      `json:"schema"`
	Packages []LockedDep `json:"packages"`
}

// --- Package Index (packages/index.json) ---

// IndexEntry maps a package version to its download URL and expected SHA-256.
type IndexEntry struct {
	URL    string `json:"url"`
	SHA256 string `json:"sha256"`
}

// PackageIndex is the parsed packages/index.json.
// Schema: cn.package-index.v1.
type PackageIndex struct {
	Schema   string                           `json:"schema"`
	Packages map[string]map[string]IndexEntry `json:"packages"`
}

// ParseManifest parses a deps.json manifest from raw JSON bytes.
// Pure — no IO. The IO wrapper lives in internal/restore/restore.go.
func ParseManifest(data []byte) (*Manifest, error) {
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parse manifest: %w", err)
	}
	return &m, nil
}

// ParseLockfile parses a lockfile from raw JSON bytes.
func ParseLockfile(data []byte) (*Lockfile, error) {
	var lf Lockfile
	if err := json.Unmarshal(data, &lf); err != nil {
		return nil, fmt.Errorf("parse lockfile: %w", err)
	}
	return &lf, nil
}

// ParsePackageIndex parses a package index from raw JSON bytes.
func ParsePackageIndex(data []byte) (*PackageIndex, error) {
	var idx PackageIndex
	if err := json.Unmarshal(data, &idx); err != nil {
		return nil, fmt.Errorf("parse package index: %w", err)
	}
	return &idx, nil
}

// Lookup finds an index entry by name and version.
// Returns nil if the pair is not in the index.
func (idx *PackageIndex) Lookup(name, version string) *IndexEntry {
	versions, ok := idx.Packages[name]
	if !ok {
		return nil
	}
	entry, ok := versions[version]
	if !ok {
		return nil
	}
	return &entry
}

// --- Package command entries ---

// PackageCommandEntry is one command declared in a cn.package.json
// under the "commands" object. Pure type — no IO.
type PackageCommandEntry struct {
	Name       string // command name (map key, e.g. "daily")
	Entrypoint string // relative path to entrypoint script (e.g. "commands/daily/cn-daily")
	Summary    string // brief description for help output
}

// EnginesJSON is the "engines" object in cn.package.json.
// Currently only cnos compatibility is declared.
type EnginesJSON struct {
	Cnos string `json:"cnos"`
}

// ContentClasses is the canonical, ordered list of package content
// class directory names. This is the single source of truth shared
// by `cn build --check` (filesystem discovery) and `cn status`
// (installed-package display), per docs/alpha/package-system/PACKAGE-SYSTEM.md
// §1.1.
//
// Order follows the canonical §1.1 table. Filesystem presence — not
// manifest JSON fields — determines which classes a package contains
// (PACKAGE-SYSTEM.md §3: "Content classes are discovered by directory
// presence"). `providers` is intentionally absent: per
// POLYGLOT-PACKAGES-AND-PROVIDERS.md, providers are a runtime
// capability surface, not a package content class.
var ContentClasses = []string{
	"doctrine", "mindsets", "skills", "extensions",
	"templates", "commands", "orchestrators", "katas",
}

// FullPackageManifest is the extended shape of cn.package.json that
// includes commands and engines. Used by command discovery and
// status display.
//
// Content-class membership is NOT declared here: it is discovered
// from filesystem presence (see ContentClasses above). The manifest
// only carries the metadata that is not derivable from the directory
// layout itself — currently `commands` metadata and `engines.cnos`.
type FullPackageManifest struct {
	Schema   string                        `json:"schema"`
	Name     string                        `json:"name"`
	Version  string                        `json:"version"`
	Kind     string                        `json:"kind"`
	Engines  EnginesJSON                   `json:"engines"`
	Commands map[string]PackageCommandJSON `json:"commands"`
}

// PackageCommandJSON is the JSON shape of one command entry in
// cn.package.json. Exported for deserialization only.
type PackageCommandJSON struct {
	Entrypoint string `json:"entrypoint"`
	Summary    string `json:"summary"`
}

// ParseFullManifestData parses a cn.package.json and extracts the
// package name and command declarations. Pure function — takes bytes, no IO.
func ParseFullManifestData(data []byte) (*FullPackageManifest, error) {
	var m FullPackageManifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parse package manifest: %w", err)
	}
	if m.Name == "" {
		return nil, fmt.Errorf("package manifest missing 'name' field")
	}
	return &m, nil
}

// CommandEntries returns the declared commands as a flat slice.
// The map key becomes the command name.
func (m *FullPackageManifest) CommandEntries() []PackageCommandEntry {
	entries := make([]PackageCommandEntry, 0, len(m.Commands))
	// Collect keys and sort for deterministic output (eng/go §2.13).
	keys := make([]string, 0, len(m.Commands))
	for name := range m.Commands {
		keys = append(keys, name)
	}
	slices.Sort(keys)
	for _, name := range keys {
		cmd := m.Commands[name]
		entries = append(entries, PackageCommandEntry{
			Name:       name,
			Entrypoint: cmd.Entrypoint,
			Summary:    cmd.Summary,
		})
	}
	return entries
}

// --- Helpers ---

// IsFirstParty returns true iff name starts with "cnos.".
func IsFirstParty(name string) bool {
	return len(name) >= 5 && name[:5] == "cnos."
}

// --- Package manifest validation ---

// PackageManifest is the minimal shape of cn.package.json that we
// validate after extraction. Only the name field is checked.
type PackageManifest struct {
	Name string `json:"name"`
}

// ValidatePackageManifestData checks that raw cn.package.json bytes
// parse as JSON and declare a name matching expectedName. The caller
// is responsible for reading the file from disk (IO lives in
// internal/restore/, not here).
func ValidatePackageManifestData(data []byte, expectedName string) error {
	var m PackageManifest
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("invalid cn.package.json: %w", err)
	}
	if m.Name == "" {
		return fmt.Errorf("cn.package.json missing 'name' field")
	}
	if m.Name != expectedName {
		return fmt.Errorf("cn.package.json name %q does not match expected %q", m.Name, expectedName)
	}
	return nil
}

// VendorPath returns the install path for a package in the vendor tree.
// Format: <hubPath>/.cn/vendor/packages/<name>/
// Per BUILD-AND-DIST.md: "the installed path does not [carry the version]
// (.cn/vendor/packages/cnos.core/). This keeps the active runtime simple."
// Version identity lives in the lockfile and the installed cn.package.json.
func VendorPath(hubPath, name string) string {
	return hubPath + "/.cn/vendor/packages/" + name
}

// TmpTarballPath returns the temp download path for a package tarball.
func TmpTarballPath(hubPath, name, version string) string {
	return hubPath + "/.cn/tmp/" + name + "-" + version + ".tar.gz"
}

