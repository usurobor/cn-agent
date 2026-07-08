package cli

import (
	"context"

	"github.com/usurobor/cnos/src/go/internal/hubstatus"
)

// StatusCmd implements the "status" command — shows hub info, installed
// packages, and command registry. Thin wrapper per eng/go §2.18 (T-002).
type StatusCmd struct {
	Version  string
	Registry *Registry // set by main.go for command registry display
}

func (c *StatusCmd) Spec() CommandSpec {
	return CommandSpec{
		Name:     "status",
		Summary:  "Show hub status",
		Source:   SourceKernel,
		Tier:     TierKernel,
		NeedsHub: true,
	}
}

func (c *StatusCmd) Run(_ context.Context, inv Invocation) error {
	return hubstatus.Run(inv.HubPath, c.Version, registryToCommandInfo(c.Registry), inv.Stdout)
}

// registryToCommandInfo converts the cli Registry into the pure data
// type hubstatus expects. This keeps hubstatus free of cli/ imports.
func registryToCommandInfo(reg *Registry) []hubstatus.CommandInfo {
	if reg == nil {
		return nil
	}
	tierLabel := map[CommandTier]string{
		TierKernel:    "kernel",
		TierRepoLocal: "repo-local",
		TierPackage:   "package",
	}
	var cmds []hubstatus.CommandInfo
	for _, cmd := range reg.All() {
		spec := cmd.Spec()
		cmds = append(cmds, hubstatus.CommandInfo{
			// InvocationName (cnos#612 AC4): same fix as cmd_help.go — a
			// hyphenated name shadowed by a noun group (e.g. "issues-fsm")
			// is only reachable via the space form ("issues fsm"); status
			// must not display the non-invocable hyphenated key either.
			Name:    InvocationName(reg, spec.Name),
			Summary: spec.Summary,
			Tier:    tierLabel[spec.Tier],
			Package: spec.Package,
		})
	}
	return cmds
}
