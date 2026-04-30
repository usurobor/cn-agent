// Package activate implements bootstrap prompt generation for cn activate.
//
// It reads safe hub state and writes an orientation prompt to stdout.
// No model is invoked; no secrets are included.
package activate

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// Options carries the inputs for a single activation run.
type Options struct {
	HubPath string    // resolved hub root (empty = fail with diagnostic)
	Stdout  io.Writer // receives the generated prompt only
	Stderr  io.Writer // receives diagnostics
}

// Run generates a bootstrap prompt from the hub at opts.HubPath.
// Returns a non-nil error and writes a diagnostic to stderr when the hub is invalid.
func Run(_ context.Context, opts Options) error {
	if opts.HubPath == "" {
		fmt.Fprintf(opts.Stderr, "✗ No hub found — pass a HUB_DIR or run from inside a hub.\n\n")
		fmt.Fprintf(opts.Stderr, "Fix by running:\n")
		fmt.Fprintf(opts.Stderr, "  1) cn activate HUB_DIR   (explicit path)\n")
		fmt.Fprintf(opts.Stderr, "  2) cd <hub-dir>           (then: cn activate)\n\n")
		fmt.Fprintf(opts.Stderr, "Then rerun: cn activate\n")
		return fmt.Errorf("activate: no hub path provided")
	}

	if _, err := os.Stat(opts.HubPath); err != nil {
		fmt.Fprintf(opts.Stderr, "✗ Hub path not found: %s\n", opts.HubPath)
		return fmt.Errorf("activate: hub path not found: %s", opts.HubPath)
	}
	cnDir := filepath.Join(opts.HubPath, ".cn")
	if _, err := os.Stat(cnDir); err != nil {
		fmt.Fprintf(opts.Stderr, "✗ Not a hub — no .cn/ directory at %s\n\n", opts.HubPath)
		fmt.Fprintf(opts.Stderr, "Fix by running:\n")
		fmt.Fprintf(opts.Stderr, "  1) cn init   (to create a new hub)\n\n")
		fmt.Fprintf(opts.Stderr, "Then rerun: cn activate %s\n", opts.HubPath)
		return fmt.Errorf("activate: not a hub (no .cn/ at %s)", opts.HubPath)
	}

	absPath, err := filepath.Abs(opts.HubPath)
	if err != nil {
		absPath = opts.HubPath
	}

	fmt.Fprintf(opts.Stderr, "→ Generating activation prompt for hub: %s\n", absPath)

	cfg := readConfig(cnDir)
	identity := scanIdentity(opts.HubPath)
	packages := scanPackages(cnDir)
	memory := scanMemory(opts.HubPath)
	threads := scanThreads(opts.HubPath)

	writePrompt(opts.Stdout, absPath, cfg, identity, packages, memory, threads)
	return nil
}

type hubConfig struct {
	Name    string
	Version string
	Created string
}

func readConfig(cnDir string) hubConfig {
	data, err := os.ReadFile(filepath.Join(cnDir, "config.json"))
	if err != nil {
		return hubConfig{}
	}
	var raw struct {
		Name    string `json:"name"`
		Version string `json:"version"`
		Created string `json:"created"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return hubConfig{}
	}
	return hubConfig{Name: raw.Name, Version: raw.Version, Created: raw.Created}
}

func scanIdentity(hubPath string) []string {
	candidates := []string{
		"spec/SOUL.md",
		"spec/identity.md",
		"agent/identity.md",
	}
	var found []string
	for _, rel := range candidates {
		if _, err := os.Stat(filepath.Join(hubPath, rel)); err == nil {
			found = append(found, rel)
		}
	}
	return found
}

func scanPackages(cnDir string) []string {
	entries, err := os.ReadDir(filepath.Join(cnDir, "vendor", "packages"))
	if err != nil {
		return nil
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	slices.Sort(names)
	return names
}

func scanMemory(hubPath string) []string {
	candidates := []string{
		"threads/reflections/daily",
		"threads/reflections/weekly",
		"threads/reflections/monthly",
		"threads/adhoc",
	}
	return presentPaths(hubPath, candidates)
}

func scanThreads(hubPath string) []string {
	candidates := []string{
		"threads/in",
		"threads/mail",
		"threads/archived",
	}
	return presentPaths(hubPath, candidates)
}

func presentPaths(hubPath string, candidates []string) []string {
	var found []string
	for _, rel := range candidates {
		if _, err := os.Stat(filepath.Join(hubPath, rel)); err == nil {
			found = append(found, rel)
		}
	}
	return found
}

func writePrompt(w io.Writer, absPath string, cfg hubConfig, identity, packages, memory, threads []string) {
	fmt.Fprintf(w, "You are activating a cnos hub.\n\n")
	fmt.Fprintf(w, "Hub path: %s\n", absPath)
	if cfg.Name != "" {
		fmt.Fprintf(w, "Hub name: %s\n", cfg.Name)
	}
	fmt.Fprintln(w)
	fmt.Fprintf(w, "Read and use the hub as your durable context.\n")
	fmt.Fprintf(w, "This prompt is an orientation guide — inspect hub files directly for full detail.\n")
	fmt.Fprintln(w)

	fmt.Fprintf(w, "## Identity\n")
	if len(identity) > 0 {
		for _, f := range identity {
			fmt.Fprintf(w, "- %s: present\n", f)
		}
	} else {
		fmt.Fprintf(w, "- no identity files found\n")
	}
	fmt.Fprintln(w)

	fmt.Fprintf(w, "## Packages and skills\n")
	if len(packages) > 0 {
		fmt.Fprintf(w, "- .cn/vendor/packages/: %d installed (%s)\n",
			len(packages), strings.Join(packages, ", "))
	} else {
		fmt.Fprintf(w, "- .cn/vendor/packages/: none installed\n")
	}
	fmt.Fprintln(w)

	fmt.Fprintf(w, "## Config summary\n")
	if cfg.Name != "" {
		fmt.Fprintf(w, "- name: %s\n", cfg.Name)
	}
	if cfg.Version != "" {
		fmt.Fprintf(w, "- version: %s\n", cfg.Version)
	}
	if cfg.Created != "" {
		fmt.Fprintf(w, "- created: %s\n", cfg.Created)
	}
	if cfg.Name == "" && cfg.Version == "" && cfg.Created == "" {
		fmt.Fprintf(w, "- config not available\n")
	}
	fmt.Fprintln(w)

	fmt.Fprintf(w, "## Memory and reflection\n")
	if len(memory) > 0 {
		for _, m := range memory {
			fmt.Fprintf(w, "- %s: present\n", m)
		}
	} else {
		fmt.Fprintf(w, "- no reflection surfaces found\n")
	}
	fmt.Fprintln(w)

	fmt.Fprintf(w, "## Inbox and threads\n")
	if len(threads) > 0 {
		for _, t := range threads {
			fmt.Fprintf(w, "- %s: present\n", t)
		}
	} else {
		fmt.Fprintf(w, "- no thread surfaces found\n")
	}
	fmt.Fprintln(w)

	fmt.Fprintf(w, "## Notes\n")
	fmt.Fprintf(w, "- Secrets (.cn/secrets.env, .env) are not included in this prompt.\n")
	fmt.Fprintf(w, "- This command generates a prompt only. No model is invoked.\n")
	fmt.Fprintf(w, "- To activate: cn activate HUB_DIR | claude -p \"Activate this cnos hub using the bootstrap prompt on stdin.\"\n")
}
