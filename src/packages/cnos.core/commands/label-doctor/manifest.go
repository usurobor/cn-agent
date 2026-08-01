package labeldoctor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Label is one canonical label entry as recorded in labels.json (schema
// cn.labels.v1). Owner/Group are doctrine metadata this package does not
// diff (they are not part of GitHub's label wire shape: name/color/
// description) but LoadManifest preserves them anyway, so a future
// consumer that does need them is not blocked by this package's parse.
type Label struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Owner       string `json:"owner"`
	Group       string `json:"group"`
}

// Manifest is the top-level labels.json shape (schema cn.labels.v1, per
// src/packages/cnos.core/skills/agent/label-doctrine/SKILL.md §5).
type Manifest struct {
	Schema   string  `json:"schema"`
	Owner    string  `json:"owner"`
	Doctrine string  `json:"doctrine"`
	Labels   []Label `json:"labels"`
}

// LoadManifest reads and parses a labels.json file. It does not validate
// against a JSON schema beyond basic shape (schema cn.labels.v1 is
// forward-compatible per the doctrine skill — implementers MAY add
// fields); it only requires at least one label entry, since an empty
// manifest is never a legitimate audit target.
func LoadManifest(path string) (Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Manifest{}, fmt.Errorf("label-doctor: read manifest %s: %w", path, err)
	}
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return Manifest{}, fmt.Errorf("label-doctor: parse manifest %s: %w", path, err)
	}
	if len(m.Labels) == 0 {
		return Manifest{}, fmt.Errorf("label-doctor: manifest %s has no labels", path)
	}
	return m, nil
}

// defaultManifestRelPaths are the two locations label-doctor looks for
// the canonical labels.json, in order:
//
//  1. src/packages/cnos.core/labels.json — the source-tree path, used
//     when running inside a cnos checkout itself (bootstrap/dogfood;
//     cnos#493's primary use case — auditing the cnos repo's own
//     labels).
//  2. .cn/vendor/packages/cnos.core/labels.json — the installed-package
//     path, used when running inside an arbitrary tenant repo that has
//     already run `cn repo install` (cnos.core vendored as a dependency;
//     no source tree present). This is also the path
//     repoinstall.ensureCanonicalDispatchLabels relies on: applyInstall
//     restores cnos.core into exactly this location before
//     runDispatchCds calls into this package.
var defaultManifestRelPaths = []string{
	filepath.Join("src", "packages", "cnos.core", "labels.json"),
	filepath.Join(".cn", "vendor", "packages", "cnos.core", "labels.json"),
}

// resolveDefaultManifestPath walks up from startDir looking for either
// candidate in defaultManifestRelPaths at each level, mirroring
// cnos.issues/commands/issues-fsm's resolveDefaultTablePath idiom
// (issuesfsm.go) so `cn label-doctor` (and repoinstall's in-process
// call) work without requiring an explicit manifest path.
func resolveDefaultManifestPath(startDir string) (string, error) {
	dir := startDir
	for {
		for _, rel := range defaultManifestRelPaths {
			candidate := filepath.Join(dir, rel)
			if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
				return candidate, nil
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("label-doctor: could not find %s (source-tree) or %s (vendored) above %s", defaultManifestRelPaths[0], defaultManifestRelPaths[1], startDir)
		}
		dir = parent
	}
}
