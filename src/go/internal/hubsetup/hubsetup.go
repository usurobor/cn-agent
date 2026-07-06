// Package hubsetup implements `cn setup` — make a hub wake-ready.
//
// Setup materializes the minimum state a hub needs before package
// restore is meaningful:
//
//   - .cn/ and agent/ directories exist
//   - .cn/deps.json is present (default profile if missing)
//   - .gitignore excludes .cn/vendor/
//
// It does NOT run `cn deps restore`. The user (or a downstream
// orchestrator) invokes restore explicitly — that keeps setup
// scoped to "prepare the hub", with no network or package-index
// dependency. This matches the kernel-minimal posture
// (INVARIANTS.md T-002) and keeps the command testable without
// package infrastructure.
//
// This package is cli/-boundary compliant per eng/go §2.18: all
// logic lives here, and cli/cmd_setup.go is a thin wrapper.
package hubsetup

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/usurobor/cnos/src/go/internal/pkg"
)

// Options carries the runtime context for a setup run.
type Options struct {
	// HubPath is the hub to configure (required). An empty string is
	// an error even though cn setup is listed as needs_hub=false in
	// the kernel manifest — the manifest makes setup "always available"
	// but the operation itself still needs a hub to write into. The
	// dispatch layer lets help print setup in the pre-hub command list;
	// this function rejects the empty path cleanly.
	HubPath string

	// Version is the cnos binary version — used to populate the
	// default deps.json when one does not already exist.
	Version string

	Stdout io.Writer
	Stderr io.Writer
}

// Run configures the hub at opts.HubPath.
//
// Behavior is idempotent:
//   - existing deps.json is left untouched
//   - existing .gitignore is amended only when the .cn/vendor/ line
//     is not already present (exact substring match)
//   - directories that already exist are not recreated
func Run(ctx context.Context, opts Options) error {
	_ = ctx // kept for symmetry with other domain Run funcs

	if opts.HubPath == "" {
		return fmt.Errorf("setup: hub path is empty — run this command from inside a hub")
	}

	fmt.Fprintf(opts.Stdout, "→ Setting up hub: %s\n", opts.HubPath)

	// 1. Required directories.
	for _, d := range []string{".cn", "agent"} {
		full := filepath.Join(opts.HubPath, d)
		if err := os.MkdirAll(full, 0755); err != nil {
			return fmt.Errorf("setup: create %s: %w", d, err)
		}
	}

	// 2. Default deps.json (only when missing).
	depsPath := filepath.Join(opts.HubPath, ".cn", "deps.json")
	if err := ensureDefaultDeps(depsPath, opts.Version, opts.Stdout); err != nil {
		return err
	}

	// 3. .gitignore entry for .cn/vendor/.
	if err := EnsureGitignoreEntry(opts.HubPath, opts.Stdout); err != nil {
		return err
	}

	fmt.Fprintf(opts.Stdout, "\n✓ Hub setup complete.\n")
	fmt.Fprintf(opts.Stdout, "  Next: run 'cn deps restore' to install packages.\n")
	return nil
}

// ensureDefaultDeps writes a default deps.json at path if none exists.
// The default manifest pins cnos.core and cnos.eng to the given version
// under the "engineer" profile. Matches
// Cn_deps.default_manifest_for_profile in src/ocaml/cmd/cn_deps.ml.
func ensureDefaultDeps(path, version string, stdout io.Writer) error {
	if _, err := os.Stat(path); err == nil {
		// Present — preserve.
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("setup: stat %s: %w", path, err)
	}

	m := pkg.Manifest{
		Schema:  "cn.deps.v1",
		Profile: "engineer",
		Packages: []pkg.ManifestDep{
			{Name: "cnos.core", Version: version},
			{Name: "cnos.eng", Version: version},
		},
	}
	data, err := json.MarshalIndent(&m, "", "  ")
	if err != nil {
		return fmt.Errorf("setup: marshal deps.json: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("setup: write %s: %w", path, err)
	}
	fmt.Fprintf(stdout, "✓ Created .cn/deps.json (profile: engineer)\n")
	return nil
}

// EnsureGitignoreEntry adds ".cn/vendor/" to .gitignore if it is not
// already present. Creates .gitignore if it does not exist. Idempotent.
//
// Exported (cnos#608) so cn repo install (internal/repoinstall) can reuse
// this exact logic rather than duplicating it — cn setup and cn repo
// install both need the same ".cn/vendor/" gitignore invariant, and this
// is the single canonical writer.
func EnsureGitignoreEntry(hubPath string, stdout io.Writer) error {
	giPath := filepath.Join(hubPath, ".gitignore")
	entry := ".cn/vendor/"

	existing, err := os.ReadFile(giPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("setup: read .gitignore: %w", err)
	}

	// Already present — nothing to do.
	if containsLine(existing, entry) {
		return nil
	}

	var next []byte
	switch {
	case len(existing) == 0:
		next = []byte(entry + "\n")
	case existing[len(existing)-1] == '\n':
		next = append(existing, []byte(entry+"\n")...)
	default:
		next = append(existing, []byte("\n"+entry+"\n")...)
	}
	if err := os.WriteFile(giPath, next, 0644); err != nil {
		return fmt.Errorf("setup: write .gitignore: %w", err)
	}
	fmt.Fprintf(stdout, "✓ Added .cn/vendor/ to .gitignore\n")
	return nil
}

// containsLine returns true if `data` contains `want` as a line,
// ignoring leading/trailing whitespace and matching substring
// occurrences that appear on their own. Prefix matches such as
// ".cn/vendor/subdir" do not count.
func containsLine(data []byte, want string) bool {
	for _, line := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(line) == want {
			return true
		}
	}
	return false
}
