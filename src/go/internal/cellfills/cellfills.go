// Package cellfills is the application composition root: the one place that
// knows which fills this binary ships and what they depend on.
//
// It exists so neither side of the boundary has to know the other. The generic
// runner receives an already-closed registry and only dispatches it; the CLI
// stays a thin wrapper; and `cds.patch` keeps its own dependency on installed
// skills without the runner ever learning that a patch alpha needs any.
//
// Skill authority is the canonical INSTALLED package root under the hub.
// There is no fallback search, discovery, or service locator: if a skill is
// not installed, construction fails and says so. Tests inject an explicit
// tree instead of relying on a search order.
package cellfills

import (
	"path/filepath"

	"github.com/usurobor/cnos/src/go/internal/cdspatch"
	"github.com/usurobor/cnos/src/go/internal/cellfill"
	"github.com/usurobor/cnos/src/go/internal/cellskill"
)

// InstalledPackages is the canonical skill authority under a hub:
// `<hub>/.cn/vendor/packages/<package>/skills/<path>/SKILL.md`.
func InstalledPackages(hubPath string) string {
	return filepath.Join(hubPath, ".cn", "vendor", "packages")
}

// Assemble builds the fill registry this binary ships: the generic cdd fills
// plus the CDS patch constructor, closed over the installed package root.
func Assemble(hubPath string) cellfill.Registry {
	return With(cellskill.Tree{Root: InstalledPackages(hubPath)})
}

// With is Assemble over an explicit skill resolver, for tests that supply
// their own installed tree.
func With(skills cellskill.Resolver) cellfill.Registry {
	reg := cellfill.CddFills()
	reg.Alpha[cdspatch.Fill] = cdspatch.Factory(skills)
	return reg
}
