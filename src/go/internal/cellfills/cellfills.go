// Package cellfills is the application composition root: the one place that knows
// which fills this binary ships, so the runner receives a closed registry and only
// dispatches it. Skill authority is the canonical INSTALLED package root under the
// hub — no fallback search or service locator, so an uninstalled skill fails
// construction and says so.
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

// Assemble builds the fill registry this binary ships.
func Assemble(hubPath string) cellfill.Registry {
	return With(cellskill.Tree{Root: InstalledPackages(hubPath)})
}

// With is Assemble over an explicit skill resolver, for tests with their own tree.
// The door is wired here for the reason the fills are: the runner names no CDS
// package, so a second profile is another line here, not a branch in the runner.
func With(skills cellskill.Resolver) cellfill.Registry {
	reg := cellfill.CddFills()
	// The resolver goes on the REGISTRY: a fill closed over its own resolver is a
	// second place skills could enter a run, which is what `cds.patch`'s list was.
	reg.Skills = skills
	// NeedsSubject is declared at the registration: a patch alpha measures a change
	// against the pinned subject's repository, so the spec loader refuses without one.
	reg.Alpha[cdspatch.Fill] = cellfill.AlphaFill{
		Construct:    cdspatch.Factory(),
		NeedsSubject: true,
	}
	// The assessing seat needs the subject too (it reconstructs the candidate), so a
	// pairing with a subjectless alpha is refused by the spec, not later in Review.
	reg.Beta[cdsassess.Fill] = cellfill.BetaFill{
		Construct:    cdsassess.Factory(),
		NeedsSubject: true,
	}
	reg.Door = cdsadmit.Door
	return reg
}
