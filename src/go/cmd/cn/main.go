// Command cn is the cnos kernel binary.
//
// It discovers the hub, registers kernel commands, and dispatches
// based on the first argument. See GO-KERNEL-COMMANDS.md for the
// full architecture.
package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/usurobor/cnos/src/go/internal/cli"
	"github.com/usurobor/cnos/src/go/internal/discover"
)

// version and commit are set at build time via -ldflags. They default
// to placeholders so `go build ./cmd/cn` works without extra flags.
var (
	version = "dev"
	commit  = ""
)

func main() {
	ctx := context.Background()

	// Build the kernel command registry.
	reg := cli.NewRegistry()
	helpCmd := &cli.HelpCmd{Registry: reg}
	reg.Register(helpCmd)
	reg.Register(&cli.InitCmd{})
	reg.Register(&cli.SetupCmd{Version: version})
	reg.Register(&cli.DepsCmd{})
	statusCmd := &cli.StatusCmd{Version: version, Registry: reg}
	reg.Register(statusCmd)
	doctorCmd := &cli.DoctorCmd{Version: version, Registry: reg}
	reg.Register(doctorCmd)
	reg.Register(&cli.BuildCmd{})
	reg.Register(&cli.UpdateCmd{Version: version, Commit: commit})
	reg.Register(&cli.ActivateCmd{})
	reg.Register(&cli.DispatchCmd{})

	// Discover hub: walk up from cwd to find .cn/.
	hubPath := discoverHub()

	// Discover repo-local and package commands (GO-KERNEL-COMMANDS.md §2.6).
	// Kernel commands are already registered above — they always win by tier.
	if hubPath != "" {
		// Step 4: repo-local commands (.cn/commands/cn-<name>).
		for _, cmd := range discover.ScanRepoLocalCommands(hubPath) {
			reg.Register(cmd)
		}
		// Step 5: vendor package commands (.cn/vendor/packages/*/cn.package.json).
		for _, cmd := range discover.ScanPackageCommands(hubPath) {
			reg.Register(cmd)
		}
	}

	// No args → help.
	if len(os.Args) < 2 {
		inv := makeInvocation(hubPath, nil)
		helpCmd.Run(ctx, inv)
		os.Exit(0)
	}

	// Resolve the command — supports noun-verb form (cn kata run),
	// flat form (cn kata-run), and noun-group listing (cn kata).
	// See docs/alpha/DESIGN-CONSTRAINTS.md §3.1.
	res := cli.ResolveCommand(reg, os.Args[1:])

	// Group case: args[0] is a noun prefix with no matching verb.
	if res.Group != "" {
		unknownVerb := len(res.Remaining) > 0 && !strings.HasPrefix(res.Remaining[0], "-")
		if unknownVerb {
			fmt.Fprintf(os.Stderr, "✗ Unknown subcommand: %s %s\n\n", res.Group, res.Remaining[0])
			cli.PrintGroup(os.Stderr, reg, res.Group)
			os.Exit(1)
		}
		cli.PrintGroup(os.Stdout, reg, res.Group)
		os.Exit(0)
	}

	// No command, no group — unknown.
	if res.Command == nil {
		fmt.Fprintf(os.Stderr, "✗ Unknown command: %s\n\n", os.Args[1])
		fmt.Fprintf(os.Stderr, "Run 'cn help' to see available commands.\n")
		os.Exit(1)
	}

	cmd := res.Command
	args := res.Remaining

	// Check hub requirement.
	spec := cmd.Spec()
	if spec.NeedsHub && hubPath == "" {
		fmt.Fprintf(os.Stderr, "✗ Command '%s' requires a hub, but no hub found.\n\n", spec.Name)
		fmt.Fprintf(os.Stderr, "Fix by running:\n")
		fmt.Fprintf(os.Stderr, "  1) cn init    (to create a new hub)\n")
		fmt.Fprintf(os.Stderr, "  2) cd <hub>   (to enter an existing hub)\n\n")
		fmt.Fprintf(os.Stderr, "Then rerun: cn %s\n", strings.Join(os.Args[1:], " "))
		os.Exit(1)
	}

	inv := makeInvocation(hubPath, args)
	if err := cmd.Run(ctx, inv); err != nil {
		// The command already printed user-facing output to stderr.
		os.Exit(1)
	}
}

// discoverHub walks up from cwd looking for a .cn/ directory.
// Returns "" if no hub is found (not an error — some commands
// work without a hub).
func discoverHub() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	for {
		if info, err := os.Stat(filepath.Join(dir, ".cn")); err == nil && info.IsDir() {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

func makeInvocation(hubPath string, args []string) cli.Invocation {
	return cli.Invocation{
		HubPath: hubPath,
		Args:    args,
		Stdin:   os.Stdin,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
		// Env left nil — commands use os.Getenv directly when needed.
		// Copying the entire env on every invocation was wasteful (F2).
	}
}
