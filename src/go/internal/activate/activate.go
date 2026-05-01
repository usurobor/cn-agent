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
	kernel := resolveKernel(cnDir)
	persona := scanPersona(opts.HubPath)
	operator := scanOperator(opts.HubPath)
	deps := scanDeps(cnDir)
	latest := latestReflection(opts.HubPath)
	memory := scanMemory(opts.HubPath)
	threads := scanThreads(opts.HubPath)

	writePrompt(opts.Stdout, absPath, cfg, kernel, persona, operator, deps, latest, memory, threads)
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

// kernelState is the result of resolveKernel.
type kernelState struct {
	state   string // "vendored", "manifest-only", "none"
	path    string // vendored path when state == "vendored"
	version string // kernel version when state == "vendored"
}

// resolveKernel reads the vendored cnos.core package to determine kernel state.
// Three states:
//
//	"vendored"      — doctrine/KERNEL.md is present; path + version returned
//	"manifest-only" — .cn/deps.json declares cnos.core but no vendor dir
//	"none"          — no reference to cnos.core at all
func resolveKernel(cnDir string) kernelState {
	vendoredKernel := filepath.Join(cnDir, "vendor", "packages", "cnos.core", "doctrine", "KERNEL.md")
	if _, err := os.Stat(vendoredKernel); err == nil {
		version := readPackageVersion(filepath.Join(cnDir, "vendor", "packages", "cnos.core", "cn.package.json"))
		return kernelState{state: "vendored", path: vendoredKernel, version: version}
	}
	// Check whether deps.json declares cnos.core.
	if manifestDeclaresCnosCore(filepath.Join(cnDir, "deps.json")) {
		return kernelState{state: "manifest-only"}
	}
	return kernelState{state: "none"}
}

func readPackageVersion(manifestPath string) string {
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return ""
	}
	var raw struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return ""
	}
	return raw.Version
}

func manifestDeclaresCnosCore(depsPath string) bool {
	data, err := os.ReadFile(depsPath)
	if err != nil {
		return false
	}
	var raw struct {
		Packages []struct {
			Name string `json:"name"`
		} `json:"packages"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return false
	}
	for _, p := range raw.Packages {
		if p.Name == "cnos.core" {
			return true
		}
	}
	return false
}

// scanPersona returns "present" path if spec/PERSONA.md exists, else "".
func scanPersona(hubPath string) string {
	p := filepath.Join(hubPath, "spec", "PERSONA.md")
	if _, err := os.Stat(p); err == nil {
		return "spec/PERSONA.md"
	}
	return ""
}

// scanOperator returns "present" path if spec/OPERATOR.md exists, else "".
func scanOperator(hubPath string) string {
	p := filepath.Join(hubPath, "spec", "OPERATOR.md")
	if _, err := os.Stat(p); err == nil {
		return "spec/OPERATOR.md"
	}
	return ""
}

// depsState is the result of scanDeps.
type depsState struct {
	state    string   // "restored", "manifest-only", "none"
	packages []string // installed package names when state == "restored"
}

// scanDeps reads .cn/deps.json and the vendor directory.
// Three states:
//
//	"restored"      — deps.json present and vendor/packages/ has entries
//	"manifest-only" — deps.json present but vendor/packages/ is empty/absent
//	"none"          — no deps.json
func scanDeps(cnDir string) depsState {
	if _, err := os.Stat(filepath.Join(cnDir, "deps.json")); err != nil {
		return depsState{state: "none"}
	}
	entries, err := os.ReadDir(filepath.Join(cnDir, "vendor", "packages"))
	if err != nil || len(entries) == 0 {
		return depsState{state: "manifest-only"}
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	slices.Sort(names)
	return depsState{state: "restored", packages: names}
}

// latestReflection returns the path of the most recently modified file
// under threads/reflections/daily/ or "" if none exist.
func latestReflection(hubPath string) string {
	dir := filepath.Join(hubPath, "threads", "reflections", "daily")
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) == 0 {
		return ""
	}
	// ReadDir returns entries sorted by name; pick the last (lexically latest).
	// Daily reflection files are named YYYY-MM-DD.md, so lexical == chronological.
	for i := len(entries) - 1; i >= 0; i-- {
		e := entries[i]
		if !e.IsDir() {
			return filepath.Join("threads", "reflections", "daily", e.Name())
		}
	}
	return ""
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

func writePrompt(w io.Writer, absPath string, cfg hubConfig,
	kernel kernelState, persona, operator string,
	deps depsState, latest string,
	memory, threads []string,
) {
	fmt.Fprintf(w, "You are activating a cnos hub.\n\n")
	fmt.Fprintf(w, "Hub path: %s\n", absPath)
	if cfg.Name != "" {
		fmt.Fprintf(w, "Hub name: %s\n", cfg.Name)
	}
	fmt.Fprintln(w)

	// ## Read first — ordered reading path
	fmt.Fprintf(w, "## Read first\n")
	if persona != "" {
		fmt.Fprintf(w, "1. %s\n", persona)
	} else {
		fmt.Fprintf(w, "1. spec/PERSONA.md — not found; consider creating one\n")
	}
	if operator != "" {
		fmt.Fprintf(w, "2. %s\n", operator)
	} else {
		fmt.Fprintf(w, "2. spec/OPERATOR.md — not found; consider creating one\n")
	}
	switch kernel.state {
	case "vendored":
		if kernel.version != "" {
			fmt.Fprintf(w, "3. %s (kernel @%s)\n", kernel.path, kernel.version)
		} else {
			fmt.Fprintf(w, "3. %s\n", kernel.path)
		}
	case "manifest-only":
		fmt.Fprintf(w, "3. kernel — dependency manifest declares cnos.core; not restored — run cn deps restore\n")
	default:
		fmt.Fprintf(w, "3. kernel — no kernel reference\n")
	}
	switch deps.state {
	case "restored":
		fmt.Fprintf(w, "4. .cn/deps.json (manifest; %d packages restored)\n", len(deps.packages))
	case "manifest-only":
		fmt.Fprintf(w, "4. .cn/deps.json (manifest present; packages not restored — run cn deps restore)\n")
	default:
		fmt.Fprintf(w, "4. deps manifest — no dependency manifest\n")
	}
	if latest != "" {
		fmt.Fprintf(w, "5. %s (latest reflection)\n", latest)
	}
	fmt.Fprintln(w)

	// ## Kernel
	fmt.Fprintf(w, "## Kernel\n")
	switch kernel.state {
	case "vendored":
		if kernel.version != "" {
			fmt.Fprintf(w, "vendored at %s@%s\n", kernel.path, kernel.version)
		} else {
			fmt.Fprintf(w, "vendored at %s\n", kernel.path)
		}
	case "manifest-only":
		fmt.Fprintf(w, "dependency manifest declares cnos.core; not restored — run cn deps restore\n")
	default:
		fmt.Fprintf(w, "no kernel reference\n")
	}
	fmt.Fprintln(w)

	// ## Persona
	fmt.Fprintf(w, "## Persona\n")
	if persona != "" {
		fmt.Fprintf(w, "%s: present\n", persona)
	} else {
		fmt.Fprintf(w, "no spec/PERSONA.md found — consider creating one\n")
	}
	fmt.Fprintln(w)

	// ## Operator
	fmt.Fprintf(w, "## Operator\n")
	if operator != "" {
		fmt.Fprintf(w, "%s: present\n", operator)
	} else {
		fmt.Fprintf(w, "no spec/OPERATOR.md found — consider creating one\n")
	}
	fmt.Fprintln(w)

	// ## Dependencies
	fmt.Fprintf(w, "## Dependencies\n")
	switch deps.state {
	case "restored":
		fmt.Fprintf(w, "restored: %s\n", strings.Join(deps.packages, ", "))
	case "manifest-only":
		fmt.Fprintf(w, "dependency manifest present; packages not restored — run cn deps restore\n")
	default:
		fmt.Fprintf(w, "no dependency manifest\n")
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
