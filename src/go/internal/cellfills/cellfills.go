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

	"github.com/usurobor/cnos/src/go/internal/cdsadmit"
	"github.com/usurobor/cnos/src/go/internal/cdsassess"
	"github.com/usurobor/cnos/src/go/internal/cdspatch"
	"github.com/usurobor/cnos/src/go/internal/cellfill"
	"github.com/usurobor/cnos/src/go/internal/cellskill"
)

// InstalledPackages is the canonical skill authority under a hub:
// `<hub>/.cn/vendor/packages/<package>/skills/<path>/SKILL.md`.
func InstalledPackages(hubPath string) string {
	return filepath.Join(hubPath, ".cn", "vendor", "packages")
}

// Assemble builds the fill registry this binary ships: the generic cdd fills,
// the CDS patch constructor closed over the installed package root, and the CDS
// admission door.
func Assemble(hubPath string) cellfill.Registry {
	return With(cellskill.Tree{Root: InstalledPackages(hubPath)})
}

// With is Assemble over an explicit skill resolver, for tests that supply
// their own installed tree.
//
// The door is wired here for exactly the reason the fills are: it is the one
// place that may know which profile this binary ships. `cn cell run` is the
// generic runner — it dispatches whatever door it is handed and names no CDS
// package, so a second profile is another line in this function rather than an
// import and a branch inside the runner.
func With(skills cellskill.Resolver) cellfill.Registry {
	reg := cellfill.CddFills()
	// The resolver goes on the REGISTRY, not into a fill. The cell declares one
	// methodology bundle and the loader loads it once; a fill closed over its
	// own resolver would be a second place skills could enter a run, which is
	// what `cds.patch`'s own `skills` list was.
	reg.Skills = skills
	// NeedsSubject is declared here, at the registration, because that is where
	// this binary states what `cds.patch` is: a patch alpha measures a change
	// against the repository and base the contract's pinned subject names, and
	// there is nothing for it to act on without one. Declaring it lets the spec
	// loader refuse a subjectless run before the constructor builds a provider
	// adapter.
	reg.Alpha[cdspatch.Fill] = cellfill.AlphaFill{
		Construct:    cdspatch.Factory(),
		NeedsSubject: true,
	}
	// The CDS assessing seat. It needs the pinned subject too — it reconstructs
	// the candidate from it — but the requirement is not declared here, because
	// the beta side of the registry carries no requirement field: nothing has
	// needed one yet in a way the loader could act on. A cds.assess cell pairs
	// with a cds.patch alpha, whose declared requirement already refuses a
	// subjectless run before either constructor is reached; a cell pairing this
	// beta with a subjectless alpha would instead be refused by the seat at
	// review time, which is later and noisier. Stated rather than left as an
	// asymmetry a reader has to notice.
	reg.Beta[cdsassess.Fill] = cdsassess.Factory()
	reg.Door = cdsadmit.Door
	return reg
}
