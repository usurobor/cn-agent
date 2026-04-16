package cli

import (
	"fmt"
	"io"
	"strings"
)

// Resolution is the outcome of ResolveCommand.
//
// Exactly one of Command or Group is non-empty; both are empty when no
// command matched and args[0] is not a known noun prefix.
type Resolution struct {
	// Command is the resolved command, or nil if no command matched.
	Command Command

	// Remaining is the args to pass to the command after the command
	// name. When Group is set, it holds the args after the noun.
	Remaining []string

	// Group is set when no command matched but args[0] is a noun
	// prefix of one or more "<noun>-<verb>" commands. The caller
	// should list the group's members.
	Group string
}

// ResolveCommand returns the command to dispatch for the given argv
// (without the program name).
//
// Lookup order — see docs/alpha/DESIGN-CONSTRAINTS.md §3.1:
//
//  1. Noun-verb form: args[0]+"-"+args[1] as a registry key
//     (e.g., "kata run" → "kata-run").
//  2. Flat form: args[0] as a registry key
//     (e.g., "kata-run", "doctor"). Backward-compat for the
//     hyphenated form.
//  3. Group: args[0] matches "<noun>-*" in the registry with no
//     matching verb. Caller lists the group.
//
// The user-facing surface is noun-verb. The flat form is kept as an
// undocumented fallback so existing scripts do not break.
func ResolveCommand(reg *Registry, args []string) Resolution {
	if len(args) == 0 {
		return Resolution{}
	}

	// 1. Noun-verb form.
	if len(args) >= 2 {
		if cmd, ok := reg.Lookup(args[0] + "-" + args[1]); ok {
			return Resolution{Command: cmd, Remaining: args[2:]}
		}
	}

	// 2. Flat form.
	if cmd, ok := reg.Lookup(args[0]); ok {
		return Resolution{Command: cmd, Remaining: args[1:]}
	}

	// 3. Group prefix.
	if len(GroupMembers(reg, args[0])) > 0 {
		return Resolution{Group: args[0], Remaining: args[1:]}
	}

	return Resolution{}
}

// GroupMembers returns commands whose name starts with prefix+"-".
// Results are returned in registry order so help output is deterministic.
func GroupMembers(reg *Registry, prefix string) []Command {
	pfx := prefix + "-"
	var cmds []Command
	for _, c := range reg.All() {
		if strings.HasPrefix(c.Spec().Name, pfx) {
			cmds = append(cmds, c)
		}
	}
	return cmds
}

// PrintGroup writes a help listing for a noun prefix to w.
//
// The listing uses the user-facing space-separated form
// ("kata run") rather than the internal hyphenated key ("kata-run").
// Returns false when the prefix has no members.
func PrintGroup(w io.Writer, reg *Registry, prefix string) bool {
	members := GroupMembers(reg, prefix)
	if len(members) == 0 {
		return false
	}

	pfx := prefix + "-"

	// Compute the widest verb for column alignment.
	maxVerb := 0
	for _, c := range members {
		verb := strings.TrimPrefix(c.Spec().Name, pfx)
		if len(verb) > maxVerb {
			maxVerb = len(verb)
		}
	}

	fmt.Fprintf(w, "Usage: cn %s <subcommand>\n\n", prefix)
	fmt.Fprintf(w, "Subcommands:\n")
	for _, c := range members {
		spec := c.Spec()
		verb := strings.TrimPrefix(spec.Name, pfx)
		fmt.Fprintf(w, "  %s %-*s  %s\n", prefix, maxVerb, verb, spec.Summary)
	}
	return true
}
