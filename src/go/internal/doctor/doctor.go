// Package doctor implements hub health validation checks.
//
// Each check returns a CheckResult with an explicit Status: Pass,
// Info, or Fail. The caller (cli/cmd_doctor.go) renders a glyph per
// status and treats only Fail as a non-zero exit. Info is reserved
// for legitimate "not yet applicable" hub state (e.g. the lockfile
// is absent before the first `cn deps lock`) — it is visible to the
// operator but does not gate merges or wake.
//
// The three-way status separation keeps doctor's contract truthful
// (eng/design-principles §2.4): ✓ means validated, ○ means pending,
// ✗ means broken. Collapsing "pending" into "validated" would hide
// lifecycle state from status/help/doctor parity.
//
// This separation also keeps cli/ as a thin dispatch layer per
// eng/go §2.18 (INVARIANTS.md T-002).
package doctor

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/usurobor/cnos/src/go/internal/activation"
)

// Status classifies a CheckResult.
type Status int

const (
	// StatusPass — the check validated: artifact present, tool installed,
	// config sane. Rendered as ✓.
	StatusPass Status = iota
	// StatusInfo — the artifact is legitimately absent at this lifecycle
	// stage, or the feature is optional by policy. Not a failure, but
	// visible to the operator. Rendered as ○.
	StatusInfo
	// StatusFail — the hub is unhealthy or a required prerequisite is
	// missing. Doctor exits non-zero if any result has this status.
	// Rendered as ✗.
	StatusFail
)

// CheckResult records the outcome of one health check.
type CheckResult struct {
	Name   string
	Status Status
	Value  string
}

// Passed reports whether the result does not represent a failure.
// Both StatusPass and StatusInfo count as "not failing"; only
// StatusFail is a failure. Provided as an ergonomic helper for
// tests and renderers that care about the pass/fail split rather
// than the three-way status.
func (r CheckResult) Passed() bool { return r.Status != StatusFail }

// CommandIssue is a command integrity problem found during validation.
// Imported from discover package for doctor reporting.
type CommandIssue struct {
	CommandName string
	Package     string
	Problem     string
}

// RunAll executes all doctor checks and returns the results.
// commandIssues is an optional list of command integrity problems
// found during command discovery (pass nil if no discovery was done).
func RunAll(ctx context.Context, hubPath, version string, commandIssues []CommandIssue) []CheckResult {
	var checks []CheckResult

	// Prerequisites — environmental tooling a hub genuinely needs.
	// git identity is kept fatal: commits without it fail at the git
	// layer, which is invisible to doctor's caller. curl is fatal
	// because `cn update` and `cn deps restore` (remote indexes) need
	// it. These are one-time user setup, not per-hub state.
	checks = append(checks, checkTool(ctx, "git", "git", "--version"))
	checks = append(checks, checkGitConfig(ctx, "git user.name", "user.name"))
	checks = append(checks, checkGitConfig(ctx, "git user.email", "user.email"))
	checks = append(checks, checkTool(ctx, "curl", "curl", "--version"))

	// Hub structure — must hold for any hub to function.
	checks = append(checks, checkExists(hubPath, "hub directory"))
	checks = append(checks, checkConfig(hubPath))
	checks = append(checks, checkFile(hubPath, "spec/SOUL.md", "missing (optional)", StatusInfo))
	checks = append(checks, checkPeers(hubPath))
	checks = append(checks, checkExists(filepath.Join(hubPath, "agent"), "agent/"))

	// Package system.
	//
	// deps.json — required after `cn setup`; missing = user skipped setup.
	// deps.lock.json — only after `cn deps lock`; legitimately pending on a fresh hub.
	// packages (vendor) — only after `cn deps restore`; legitimately pending.
	checks = append(checks, checkFile(hubPath, ".cn/deps.json",
		"missing (run 'cn setup')", StatusFail))
	checks = append(checks, checkFile(hubPath, ".cn/deps.lock.json",
		"pending (run 'cn deps lock')", StatusInfo))
	checks = append(checks, checkPackages(hubPath, version))

	// Command integrity — broken vendor commands are fatal.
	checks = append(checks, checkCommandIntegrity(commandIssues))

	// Skill activation — only unreadable/malformed SKILL.md files are
	// fatal (genuine structural breakage). Empty triggers and trigger
	// overlaps between public skills are reported as info: cnos skill
	// keywords are many-to-many hints, and overlapping keywords are a
	// legitimate authoring pattern rather than a hub breakage.
	checks = append(checks, checkSkillActivation(hubPath))

	// Runtime contract — generated at wake; legitimately pending before
	// the first wake. A present-but-malformed contract is fatal.
	checks = append(checks, checkRuntimeContract(hubPath))

	// Git remote — optional by policy: a hub may be local-only, offline,
	// or have its remote added later. Missing = informational, not fatal.
	checks = append(checks, checkGitRemote(ctx, hubPath))

	return checks
}

// HasFailures returns true if any check has Status == StatusFail.
// StatusInfo results are non-fatal and do not count.
func HasFailures(checks []CheckResult) bool {
	for _, ch := range checks {
		if ch.Status == StatusFail {
			return true
		}
	}
	return false
}

// --- Check implementations ---

func checkTool(ctx context.Context, name string, args ...string) CheckResult {
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	out, err := cmd.Output()
	if err != nil {
		return CheckResult{Name: name, Status: StatusFail, Value: "not installed"}
	}
	line := strings.SplitN(strings.TrimSpace(string(out)), "\n", 2)[0]
	if strings.HasPrefix(line, "git version ") {
		line = strings.TrimPrefix(line, "git version ")
	}
	if name == "curl" {
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			line = parts[1]
		}
	}
	return CheckResult{Name: name, Status: StatusPass, Value: line}
}

func checkGitConfig(ctx context.Context, label, key string) CheckResult {
	cmd := exec.CommandContext(ctx, "git", "config", key)
	out, err := cmd.Output()
	if err != nil {
		return CheckResult{Name: label, Status: StatusFail, Value: "not set"}
	}
	return CheckResult{Name: label, Status: StatusPass, Value: strings.TrimSpace(string(out))}
}

func checkExists(path, label string) CheckResult {
	if _, err := os.Stat(path); err == nil {
		return CheckResult{Name: label, Status: StatusPass, Value: "exists"}
	}
	return CheckResult{Name: label, Status: StatusFail, Value: "not found"}
}

func checkConfig(hubPath string) CheckResult {
	jsonPath := filepath.Join(hubPath, ".cn", "config.json")
	yamlPath := filepath.Join(hubPath, ".cn", "config.yaml")
	if _, err := os.Stat(jsonPath); err == nil {
		return CheckResult{Name: ".cn/config", Status: StatusPass, Value: "config.json"}
	}
	if _, err := os.Stat(yamlPath); err == nil {
		return CheckResult{Name: ".cn/config", Status: StatusPass, Value: "config.yaml"}
	}
	return CheckResult{Name: ".cn/config", Status: StatusFail, Value: "missing"}
}

// checkFile reports whether a hub-relative file is present. If the
// file is absent, missingStatus controls how that absence is recorded:
//   - StatusFail: the file is required (e.g. deps.json after setup)
//   - StatusInfo: the file is pending or optional by policy
func checkFile(hubPath, rel, missingMsg string, missingStatus Status) CheckResult {
	path := filepath.Join(hubPath, rel)
	if _, err := os.Stat(path); err == nil {
		return CheckResult{Name: rel, Status: StatusPass, Value: "present"}
	}
	return CheckResult{Name: rel, Status: missingStatus, Value: missingMsg}
}

func checkPeers(hubPath string) CheckResult {
	path := filepath.Join(hubPath, "state", "peers.md")
	data, err := os.ReadFile(path)
	if err != nil {
		return CheckResult{Name: "state/peers.md", Status: StatusFail, Value: "missing"}
	}
	count := 0
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "- name:") {
			count++
		}
	}
	return CheckResult{Name: "state/peers.md", Status: StatusPass,
		Value: fmt.Sprintf("%d peer(s)", count)}
}

func checkPackages(hubPath, _ string) CheckResult {
	vendorDir := filepath.Join(hubPath, ".cn", "vendor", "packages")
	entries, err := os.ReadDir(vendorDir)
	if err != nil {
		// No vendor directory is pending state before the first restore.
		// It becomes a failure only when a lockfile claims packages
		// should be present.
		lockPath := filepath.Join(hubPath, ".cn", "deps.lock.json")
		if _, lockErr := os.Stat(lockPath); lockErr != nil {
			return CheckResult{Name: "packages", Status: StatusInfo,
				Value: "none installed (pending 'cn deps restore')"}
		}
		return CheckResult{Name: "packages", Status: StatusFail,
			Value: "no vendor directory (run 'cn deps restore')"}
	}
	total := 0
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		total++
	}
	if total == 0 {
		return CheckResult{Name: "packages", Status: StatusPass, Value: "none installed"}
	}
	// Version drift detection: compare installed packages against the
	// lockfile. The installed path no longer carries the version
	// (BUILD-AND-DIST.md), so we read the lockfile as the expected set.
	lockPath := filepath.Join(hubPath, ".cn", "deps.lock.json")
	lockData, err := os.ReadFile(lockPath)
	if err != nil {
		return CheckResult{Name: "packages", Status: StatusPass,
			Value: fmt.Sprintf("%d installed (no lockfile for drift check)", total)}
	}
	var lf struct {
		Packages []struct {
			Name string `json:"name"`
		} `json:"packages"`
	}
	if err := json.Unmarshal(lockData, &lf); err != nil {
		return CheckResult{Name: "packages", Status: StatusPass,
			Value: fmt.Sprintf("%d installed (lockfile parse error)", total)}
	}
	missing := 0
	for _, dep := range lf.Packages {
		pkgDir := filepath.Join(vendorDir, dep.Name)
		if _, err := os.Stat(pkgDir); err != nil {
			missing++
		}
	}
	if missing > 0 {
		return CheckResult{Name: "packages", Status: StatusFail,
			Value: fmt.Sprintf("%d installed, %d missing from lockfile (run 'cn deps restore')", total, missing)}
	}
	return CheckResult{Name: "packages", Status: StatusPass,
		Value: fmt.Sprintf("%d installed, all present", total)}
}

func checkRuntimeContract(hubPath string) CheckResult {
	path := filepath.Join(hubPath, "state", "runtime-contract.json")
	data, err := os.ReadFile(path)
	if err != nil {
		// Generated at wake — legitimately pending on a fresh hub.
		return CheckResult{Name: "runtime contract", Status: StatusInfo,
			Value: "pending (generated at wake)"}
	}
	var doc map[string]any
	if err := json.Unmarshal(data, &doc); err != nil {
		return CheckResult{Name: "runtime contract", Status: StatusFail, Value: "invalid JSON"}
	}
	layers := []string{"identity", "cognition", "body", "medium"}
	var missing []string
	for _, l := range layers {
		if _, ok := doc[l]; !ok {
			missing = append(missing, l)
		}
	}
	if len(missing) > 0 {
		return CheckResult{Name: "runtime contract", Status: StatusFail,
			Value: fmt.Sprintf("incomplete (missing: %s)", strings.Join(missing, ", "))}
	}
	return CheckResult{Name: "runtime contract", Status: StatusPass,
		Value: "valid (identity + cognition + body + medium)"}
}

func checkGitRemote(ctx context.Context, hubPath string) CheckResult {
	cmd := exec.CommandContext(ctx, "git", "remote", "get-url", "origin")
	cmd.Dir = hubPath
	if _, err := cmd.Output(); err != nil {
		// Optional by policy: a hub may be local-only or add a remote later.
		return CheckResult{Name: "origin remote", Status: StatusInfo, Value: "not configured (optional)"}
	}
	return CheckResult{Name: "origin remote", Status: StatusPass, Value: "configured"}
}

// CommandDescriptor is a plain data record describing one registered
// command. Passed from cli/ to doctor/ to avoid circular imports.
type CommandDescriptor struct {
	Name       string
	Source     string // "kernel", "repo-local", "package"
	Tier       int    // 0=kernel, 1=repo-local, 2=package
	Package    string // package name (empty for kernel/repo-local)
	Entrypoint string // absolute path to entrypoint script (empty for kernel)
}

// ValidateCommands checks command integrity from a list of descriptors:
//   - missing entrypoints
//   - non-executable entrypoints
//   - duplicate command names within a single tier
func ValidateCommands(descs []CommandDescriptor) []CommandIssue {
	var issues []CommandIssue

	tierNames := make(map[int]map[string]int)
	for _, d := range descs {
		if tierNames[d.Tier] == nil {
			tierNames[d.Tier] = make(map[string]int)
		}
		tierNames[d.Tier][d.Name]++

		if d.Entrypoint == "" {
			continue
		}

		info, err := os.Stat(d.Entrypoint)
		if err != nil {
			issues = append(issues, CommandIssue{
				CommandName: d.Name,
				Package:     d.Package,
				Problem:     "missing entrypoint: " + d.Entrypoint,
			})
			continue
		}
		if info.Mode()&0111 == 0 {
			issues = append(issues, CommandIssue{
				CommandName: d.Name,
				Package:     d.Package,
				Problem:     "entrypoint not executable: " + d.Entrypoint,
			})
		}
	}

	for tier, names := range tierNames {
		for name, count := range names {
			if count > 1 {
				tierLabel := "unknown"
				switch tier {
				case 0:
					tierLabel = "kernel"
				case 1:
					tierLabel = "repo-local"
				case 2:
					tierLabel = "package"
				}
				issues = append(issues, CommandIssue{
					CommandName: name,
					Problem:     fmt.Sprintf("duplicate in %s tier (%d definitions)", tierLabel, count),
				})
			}
		}
	}

	return issues
}

func checkSkillActivation(hubPath string) CheckResult {
	issues := activation.Validate(hubPath)
	if len(issues) == 0 {
		return CheckResult{
			Name:   "skill activation",
			Status: StatusPass,
			Value:  "all skills valid",
		}
	}
	// Severity split: an unreadable/malformed SKILL.md is the only
	// structural break — the frontmatter parser refuses to produce a
	// usable record. Empty triggers and public-skill overlaps are
	// authoring hints that surface to the operator without failing
	// the hub.
	var fatal, warn []string
	for _, iss := range issues {
		line := fmt.Sprintf("[%s] %s", activation.IssueKindLabel(iss.Kind), iss.Message)
		if iss.Kind == activation.IssueMissingSkill {
			fatal = append(fatal, line)
		} else {
			warn = append(warn, line)
		}
	}
	if len(fatal) > 0 {
		parts := append([]string{}, fatal...)
		parts = append(parts, warn...)
		return CheckResult{
			Name:   "skill activation",
			Status: StatusFail,
			Value:  fmt.Sprintf("%d issue(s): %s", len(issues), strings.Join(parts, "; ")),
		}
	}
	return CheckResult{
		Name:   "skill activation",
		Status: StatusInfo,
		Value:  fmt.Sprintf("%d warning(s): %s", len(warn), strings.Join(warn, "; ")),
	}
}

func checkCommandIntegrity(issues []CommandIssue) CheckResult {
	if len(issues) == 0 {
		return CheckResult{Name: "command integrity", Status: StatusPass, Value: "all commands valid"}
	}
	msgs := make([]string, len(issues))
	for i, iss := range issues {
		if iss.Package != "" {
			msgs[i] = fmt.Sprintf("%s (%s): %s", iss.CommandName, iss.Package, iss.Problem)
		} else {
			msgs[i] = fmt.Sprintf("%s: %s", iss.CommandName, iss.Problem)
		}
	}
	return CheckResult{
		Name:   "command integrity",
		Status: StatusFail,
		Value:  fmt.Sprintf("%d issue(s): %s", len(issues), strings.Join(msgs, "; ")),
	}
}
