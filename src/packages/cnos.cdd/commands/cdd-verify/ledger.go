package cddverify

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// runLedgerMode handles the legacy artifact-presence checks
// (--unreleased / --all / --version / --pr / --cycle / --triadic).
//
// Mirrors the bash predecessor's output format and exit-code contract:
//   - exit 0 when no FAILs (PASS, possibly with warnings)
//   - exit 1 when any FAIL emitted
//
// The output is operator-facing prose (BC.b oracle expects "Summary:"
// literal — preserved).
func runLedgerMode(ctx context.Context, args Args, stdout, stderr io.Writer) ExitCode {
	// Resolve repo root.
	repoRoot := args.RepoRootOver
	if repoRoot == "" {
		repoRoot = gitRepoRoot(ctx)
	}
	if repoRoot == "" {
		fmt.Fprintln(stderr, "Error: not in a git repository")
		return ExitFAIL
	}

	// Validate flag combinations.
	// Note: exit 1 (FAIL) matches the bash predecessor's usage-error
	// convention. ExitVError (2) is reserved for V-internal errors per
	// schemas/cdd/validation_verdict.schema.json semantics.
	if args.Cycle != "" && args.Version == "" {
		fmt.Fprintln(stderr, "Error: --cycle requires --version")
		Usage(stderr)
		return ExitFAIL
	}
	if args.All && (args.Version != "" || args.PR != "" || args.Cycle != "" || args.Triadic) {
		fmt.Fprintln(stderr, "Error: --all cannot be combined with other modes")
		Usage(stderr)
		return ExitFAIL
	}
	if args.Unreleased && (args.Version != "" || args.PR != "" || args.Cycle != "" || args.Triadic) {
		fmt.Fprintln(stderr, "Error: --unreleased cannot be combined with other modes")
		Usage(stderr)
		return ExitFAIL
	}
	if !args.All && !args.Unreleased && args.Version == "" && args.PR == "" {
		fmt.Fprintln(stderr, "Error: must specify --version, --pr, --all, or --unreleased")
		Usage(stderr)
		return ExitFAIL
	}

	ledger := newLedgerRun(repoRoot, args.Bundle, stdout)

	// Exceptions.
	if args.Exceptions != "" {
		ledger.loadExceptions(args.Exceptions)
	}

	fmt.Fprintln(stdout, "## Core artifacts")

	// PR check (best-effort).
	prNum := args.PR
	if prNum == "" && args.Version != "" {
		// Try `gh pr list` to find PR number for the version.
		prNum = lookupPRForVersion(ctx, args.Version)
	}
	if prNum != "" {
		state := lookupPRState(ctx, prNum)
		switch state {
		case "MERGED":
			ledger.check(fmt.Sprintf("PR #%s merged", prNum), checkPass, "")
		case "UNKNOWN":
			ledger.check(fmt.Sprintf("PR #%s state", prNum), checkWarn, "could not query (gh auth?)")
		default:
			ledger.check(fmt.Sprintf("PR #%s merged", prNum), checkFail, fmt.Sprintf("state: %s", state))
		}
	} else {
		ledger.check("PR identification", checkWarn, "could not determine PR number")
	}

	// Release-scoped checks.
	if args.Version != "" {
		fmt.Fprintln(stdout)
		fmt.Fprintf(stdout, "## Release artifacts (%s)\n", args.Version)

		// CHANGELOG: bare version required.
		clPath := filepath.Join(repoRoot, "CHANGELOG.md")
		if data, err := os.ReadFile(clPath); err == nil {
			if strings.Contains(string(data), args.Version) {
				ledger.check(fmt.Sprintf("CHANGELOG entry (bare %s)", args.Version), checkPass, "")
			} else {
				ledger.check(fmt.Sprintf("CHANGELOG entry (bare %s)", args.Version), checkFail, fmt.Sprintf("version %s not found in CHANGELOG.md", args.Version))
			}
		} else {
			ledger.check("CHANGELOG entry", checkWarn, "no CHANGELOG.md found")
		}

		// Git tag.
		if gitTagExists(ctx, repoRoot, args.Version) {
			ledger.check(fmt.Sprintf("Git tag (bare %s)", args.Version), checkPass, "")
		} else {
			ledger.check(fmt.Sprintf("Git tag (bare %s)", args.Version), checkFail, fmt.Sprintf("canonical bare tag '%s' not found", args.Version))
		}
		if gitTagExists(ctx, repoRoot, "v"+args.Version) {
			ledger.check("Legacy v-prefixed tag", checkWarn, fmt.Sprintf("v%s exists; v-prefixed tags are legacy/noncanonical", args.Version))
		}

		// Canonical PRA.
		canonPRA := filepath.Join(repoRoot, args.Bundle, args.Version, "POST-RELEASE-ASSESSMENT.md")
		if fileExists(canonPRA) {
			ledger.check("Canonical PRA", checkPass, fmt.Sprintf("%s/%s/POST-RELEASE-ASSESSMENT.md", args.Bundle, args.Version))
		} else {
			ledger.check("Canonical PRA", checkFail, fmt.Sprintf("expected %s/%s/POST-RELEASE-ASSESSMENT.md", args.Bundle, args.Version))
		}

		// Legacy PRA paths warn-only.
		legacyPaths := []string{
			filepath.Join(repoRoot, ".cdd/releases", args.Version, "beta/POST-RELEASE-ASSESSMENT.md"),
			filepath.Join(repoRoot, ".cdd/releases", args.Version, "beta/ASSESSMENT.md"),
		}
		for _, lp := range legacyPaths {
			if fileExists(lp) {
				rel, _ := filepath.Rel(repoRoot, lp)
				ledger.check("Legacy PRA path", checkWarn, fmt.Sprintf("%s present; canonical path is %s/%s/POST-RELEASE-ASSESSMENT.md", rel, args.Bundle, args.Version))
			}
		}

		// RELEASE.md.
		if fileExists(filepath.Join(repoRoot, "RELEASE.md")) {
			ledger.check("RELEASE.md present", checkPass, "")
		} else {
			ledger.check("RELEASE.md present", checkWarn, "absent (may have been overwritten by next release)")
		}
	}

	// Triadic close-outs.
	if args.Triadic && args.Version != "" {
		fmt.Fprintln(stdout)
		fmt.Fprintf(stdout, "## Triadic close-outs (.cdd/releases/%s/)\n", args.Version)
		cddDir := filepath.Join(repoRoot, ".cdd/releases", args.Version)
		for _, role := range []string{"alpha", "beta", "gamma"} {
			canon := filepath.Join(cddDir, role, "CLOSE-OUT.md")
			if fileExists(canon) {
				ledger.check(fmt.Sprintf("%s/CLOSE-OUT.md", role), checkPass, "")
			} else {
				ledger.check(fmt.Sprintf("%s/CLOSE-OUT.md", role), checkFail, fmt.Sprintf("expected .cdd/releases/%s/%s/CLOSE-OUT.md", args.Version, role))
			}
		}
		if fileExists(filepath.Join(cddDir, "gamma", "KATA-VERDICT.md")) {
			ledger.check("γ KATA-VERDICT.md", checkPass, "")
		} else {
			ledger.check("γ KATA-VERDICT.md", checkWarn, "not found (kata may not have run)")
		}
	}

	// Cycle-scoped checks.
	if args.Cycle != "" && args.Version != "" {
		fmt.Fprintln(stdout)
		fmt.Fprintf(stdout, "## Cycle-scoped artifacts (.cdd/releases/%s/%s/)\n", args.Version, args.Cycle)
		cycleDir := filepath.Join(repoRoot, ".cdd/releases", args.Version, args.Cycle)
		hardgate := []string{
			"self-coherence.md",
			"beta-review.md",
			"alpha-closeout.md",
			"beta-closeout.md",
			"gamma-closeout.md",
		}
		for _, artifact := range hardgate {
			ap := filepath.Join(cycleDir, artifact)
			if fileExists(ap) {
				ledger.check(artifact, checkPass, "")
				ledger.validateSections(ap, artifact, "", false)
			} else {
				ledger.check(artifact, checkFail, fmt.Sprintf("expected .cdd/releases/%s/%s/%s", args.Version, args.Cycle, artifact))
			}
		}
		if fileExists(filepath.Join(cycleDir, "gamma-clarification.md")) {
			ledger.check("gamma-clarification.md (optional)", checkPass, "")
		}
		// Orphaned unreleased dirs (only on main branch).
		if branch := gitCurrentBranch(ctx, repoRoot); branch == "main" {
			unrel := filepath.Join(repoRoot, ".cdd/unreleased")
			if entries, err := os.ReadDir(unrel); err == nil {
				var orphans []string
				for _, e := range entries {
					if e.IsDir() {
						orphans = append(orphans, e.Name())
					}
				}
				if len(orphans) > 0 {
					ledger.check("Orphaned unreleased directories", checkWarn, fmt.Sprintf("found %d directories: %s", len(orphans), strings.Join(orphans, " ")))
				} else {
					ledger.check("No orphaned unreleased directories", checkPass, "")
				}
			} else {
				ledger.check("No orphaned unreleased directories", checkPass, "no .cdd/unreleased directory")
			}
		}
	}

	// Repository-wide scans.
	if args.All || args.Unreleased {
		fmt.Fprintln(stdout)
		fmt.Fprintln(stdout, "## Unreleased cycles (.cdd/unreleased/)")
		unrelDir := filepath.Join(repoRoot, ".cdd/unreleased")
		entries, err := os.ReadDir(unrelDir)
		if err != nil {
			ledger.check("Unreleased cycles", checkPass, "no .cdd/unreleased/ directory")
		} else {
			var dirs []string
			for _, e := range entries {
				if e.IsDir() {
					dirs = append(dirs, e.Name())
				}
			}
			sort.Strings(dirs)
			if len(dirs) == 0 {
				ledger.check("Unreleased cycles", checkPass, "no unreleased cycles (.cdd/unreleased/ is empty)")
			} else {
				issueCount := 0
				for _, name := range dirs {
					if isNumeric(name) {
						ledger.checkUnreleasedCycle(name)
						issueCount++
					} else {
						ledger.check("Invalid unreleased directory", checkWarn, fmt.Sprintf("non-numeric directory name: %s", name))
					}
				}
				if issueCount == 0 {
					ledger.check("Unreleased cycles", checkWarn, "no valid issue directories found in .cdd/unreleased/")
				} else {
					ledger.check("Unreleased cycles scanned", checkPass, fmt.Sprintf("%d issue directories checked", issueCount))
				}
			}
		}
	}

	if args.All {
		fmt.Fprintln(stdout)
		fmt.Fprintln(stdout, "## Released cycles (.cdd/releases/)")
		relDir := filepath.Join(repoRoot, ".cdd/releases")
		versionEntries, err := os.ReadDir(relDir)
		if err != nil {
			ledger.check("Released cycles", checkPass, "no .cdd/releases/ directory")
		} else {
			var versions []string
			for _, ve := range versionEntries {
				if ve.IsDir() {
					versions = append(versions, ve.Name())
				}
			}
			sort.Strings(versions)
			if len(versions) == 0 {
				ledger.check("Released cycles", checkPass, "no released versions (.cdd/releases/ is empty)")
			} else {
				versionCount := 0
				totalCycles := 0
				for _, version := range versions {
					fmt.Fprintf(stdout, "  Checking version %s\n", version)
					versionDir := filepath.Join(relDir, version)
					cycleEntries, err := os.ReadDir(versionDir)
					if err != nil {
						ledger.check(fmt.Sprintf("Version %s cycles", version), checkWarn, "no cycle directories found")
						continue
					}
					var cycles []string
					for _, ce := range cycleEntries {
						if ce.IsDir() {
							cycles = append(cycles, ce.Name())
						}
					}
					sort.Strings(cycles)
					if len(cycles) == 0 {
						ledger.check(fmt.Sprintf("Version %s cycles", version), checkWarn, "no cycle directories found")
						continue
					}
					cycleCount := 0
					for _, issueNum := range cycles {
						if isNumeric(issueNum) {
							fmt.Fprintf(stdout, "    Issue #%s\n", issueNum)
							ledger.checkTriadicArtifacts(filepath.Join(versionDir, issueNum), fmt.Sprintf("%s (v%s)", issueNum, version), version)
							cycleCount++
							totalCycles++
						} else {
							ledger.check("Invalid cycle directory", checkWarn, fmt.Sprintf("non-numeric directory in %s: %s", version, issueNum))
						}
					}
					ledger.check(fmt.Sprintf("Version %s cycles", version), checkPass, fmt.Sprintf("%d cycles found", cycleCount))
					versionCount++
				}
				ledger.check("Released versions scanned", checkPass, fmt.Sprintf("%d versions, %d total cycles", versionCount, totalCycles))
			}
		}
	}

	// Summary.
	fmt.Fprintln(stdout)
	total := ledger.pass + ledger.fail + ledger.warn
	fmt.Fprintf(stdout, "## Summary: %d passed, %d failed, %d warnings (%d total)\n", ledger.pass, ledger.fail, ledger.warn, total)
	fmt.Fprintln(stdout)
	if ledger.fail > 0 {
		fmt.Fprintln(stdout, "❌ Cycle artifact verification FAILED")
		return ExitFAIL
	}
	if ledger.warn > 0 {
		fmt.Fprintln(stdout, "⚠️  Cycle artifact verification PASSED with warnings")
		return ExitPASS
	}
	fmt.Fprintln(stdout, "✅ Cycle artifact verification PASSED")
	return ExitPASS
}

type checkResult int

const (
	checkPass checkResult = iota
	checkFail
	checkWarn
)

// ledgerRun bundles state for one ledger-mode invocation: counters,
// loaded exceptions, the output sink, the bundle-relative PRA dir, and
// the repo root. eng/go §1.3 — avoid globals; collect state in a struct.
type ledgerRun struct {
	repoRoot   string
	bundle     string
	stdout     io.Writer
	pass       int
	fail       int
	warn       int
	exceptions map[string]bool
}

func newLedgerRun(repoRoot, bundle string, stdout io.Writer) *ledgerRun {
	return &ledgerRun{
		repoRoot:   repoRoot,
		bundle:     bundle,
		stdout:     stdout,
		exceptions: map[string]bool{},
	}
}

func (l *ledgerRun) check(label string, result checkResult, detail string) {
	switch result {
	case checkPass:
		l.pass++
		fmt.Fprintf(l.stdout, "  ✅ %s", label)
	case checkFail:
		l.fail++
		fmt.Fprintf(l.stdout, "  ❌ %s", label)
	case checkWarn:
		l.warn++
		fmt.Fprintf(l.stdout, "  ⚠️  %s", label)
	}
	if detail != "" {
		fmt.Fprintf(l.stdout, " — %s", detail)
	}
	fmt.Fprintln(l.stdout)
}

// loadExceptions parses a simple YAML exceptions file matching the bash
// predecessor's expected shape:
//
//	exceptions:
//	  - path: ".cdd/unreleased/286/alpha-closeout.md"
//	    missing_artifacts: [...]
//
// Only the `path:` line is needed; the rest is ignored.
func (l *ledgerRun) loadExceptions(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	fmt.Fprintf(l.stdout, "Loading legacy exceptions from: %s\n", path)
	rePath := regexp.MustCompile(`^\s*-?\s*path:\s*"?([^"\s]+)"?\s*$`)
	count := 0
	for _, line := range strings.Split(string(data), "\n") {
		m := rePath.FindStringSubmatch(line)
		if len(m) == 2 {
			l.exceptions[m[1]] = true
			count++
		}
	}
	if count > 0 {
		fmt.Fprintf(l.stdout, "Loaded %d exception paths\n", count)
	} else {
		fmt.Fprintln(l.stdout, "No exception paths found in file")
	}
}

// isExceptionBacked reports whether path is in the loaded exception set.
func (l *ledgerRun) isExceptionBacked(path string) bool {
	return l.exceptions[path]
}

// isLegacyVersion returns true for versions pre-dating current artifact
// requirements (cutoff: v3.77.0).
func isLegacyVersion(version string) bool {
	re := regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)$`)
	m := re.FindStringSubmatch(version)
	if m == nil {
		return false
	}
	major, _ := strconv.Atoi(m[1])
	minor, _ := strconv.Atoi(m[2])
	if major < 3 {
		return true
	}
	if major == 3 && minor < 77 {
		return true
	}
	return false
}

// classifyCycleType returns "triadic", "small-change", or "unknown".
func classifyCycleType(cycleDir string) string {
	required := []string{
		"self-coherence.md", "beta-review.md", "alpha-closeout.md",
		"beta-closeout.md", "gamma-closeout.md",
	}
	count := 0
	for _, a := range required {
		if fileExists(filepath.Join(cycleDir, a)) {
			count++
		}
	}
	if count == 5 {
		return "triadic"
	}
	if fileExists(filepath.Join(cycleDir, "self-coherence.md")) || fileExists(filepath.Join(cycleDir, "alpha-closeout.md")) {
		if fileExists(filepath.Join(cycleDir, "beta-review.md")) {
			return "triadic"
		}
		return "small-change"
	}
	return "unknown"
}

func (l *ledgerRun) checkUnreleasedCycle(issueNum string) {
	cycleDir := filepath.Join(l.repoRoot, ".cdd/unreleased", issueNum)
	if _, err := os.Stat(cycleDir); err != nil {
		l.check(fmt.Sprintf("Issue #%s directory", issueNum), checkFail, fmt.Sprintf(".cdd/unreleased/%s not found", issueNum))
		return
	}
	switch classifyCycleType(cycleDir) {
	case "triadic":
		fmt.Fprintf(l.stdout, "  Checking triadic cycle #%s\n", issueNum)
		l.checkUnreleasedTriadicArtifacts(cycleDir, issueNum)
	case "small-change":
		fmt.Fprintf(l.stdout, "  Checking small-change cycle #%s\n", issueNum)
		l.checkSmallChangeArtifacts(cycleDir, issueNum)
	default:
		l.check(fmt.Sprintf("Issue #%s classification", issueNum), checkWarn, "cannot determine cycle type - no recognizable artifacts")
	}
}

func (l *ledgerRun) checkUnreleasedTriadicArtifacts(cycleDir, issueNum string) {
	core := []string{"self-coherence.md", "beta-review.md"}
	closeouts := []string{"alpha-closeout.md", "beta-closeout.md", "gamma-closeout.md"}
	closeoutResult, closeoutDetail := l.unreleasedCloseoutState(cycleDir, issueNum)
	for _, artifact := range core {
		ap := filepath.Join(cycleDir, artifact)
		if fileExists(ap) {
			l.check(fmt.Sprintf("%s (issue #%s)", artifact, issueNum), checkPass, "")
			l.validateSections(ap, artifact, "", true)
		} else {
			exPath := fmt.Sprintf(".cdd/unreleased/%s/%s", issueNum, artifact)
			if l.isExceptionBacked(exPath) {
				l.check(fmt.Sprintf("%s (issue #%s)", artifact, issueNum), checkWarn, fmt.Sprintf("missing but exception-backed: %s", exPath))
			} else {
				l.check(fmt.Sprintf("%s (issue #%s)", artifact, issueNum), checkFail, fmt.Sprintf("expected .cdd/unreleased/%s/%s", issueNum, artifact))
			}
		}
	}
	for _, artifact := range closeouts {
		ap := filepath.Join(cycleDir, artifact)
		if fileExists(ap) {
			l.check(fmt.Sprintf("%s (issue #%s)", artifact, issueNum), closeoutResult, closeoutDetail)
			l.validateSections(ap, artifact, "", false)
		} else {
			l.check(fmt.Sprintf("%s (issue #%s)", artifact, issueNum), checkWarn, "missing in unreleased cycle (expected before that role's post-merge closeout)")
		}
	}
}

func (l *ledgerRun) unreleasedCloseoutState(cycleDir, issueNum string) (checkResult, string) {
	gammaPath := filepath.Join(cycleDir, "gamma-closeout.md")
	data, err := os.ReadFile(gammaPath)
	if err != nil {
		return checkWarn, "present before γ establishes release-pending marker and batch"
	}
	marker := false
	batch := ""
	batchCount := 0
	batchRE := regexp.MustCompile(`^CDD-Release-Batch: ([0-9]+\.[0-9]+\.[0-9]+|docs/[0-9]{4}-[0-9]{2}-[0-9]{2})$`)
	for _, line := range strings.Split(string(data), "\n") {
		if line == "CDD-Post-Merge-Closeout: complete" {
			marker = true
		}
		if m := batchRE.FindStringSubmatch(line); m != nil {
			batch = m[1]
			batchCount++
		}
	}
	if !marker || batchCount != 1 {
		return checkWarn, "present but release-pending is unproven (requires exact marker and one canonical batch assignment)"
	}

	archivePath := filepath.Join(l.repoRoot, ".cdd", "releases", filepath.FromSlash(batch), issueNum)
	archiveInfo, archiveErr := os.Stat(archivePath)
	archiveExists := archiveErr == nil && archiveInfo.IsDir()
	tagExists := false
	if !strings.HasPrefix(batch, "docs/") {
		cmd := exec.Command("git", "rev-parse", "--verify", "refs/tags/"+batch)
		cmd.Dir = l.repoRoot
		tagExists = cmd.Run() == nil
	}
	if tagExists || archiveExists {
		return checkWarn, fmt.Sprintf("stale under unreleased: assigned batch %s has disconnect/archive evidence", batch)
	}
	return checkPass, fmt.Sprintf("release-pending for assigned batch %s; disconnect absent", batch)
}

func (l *ledgerRun) checkSmallChangeArtifacts(cycleDir, issueNum string) {
	if fileExists(filepath.Join(cycleDir, "self-coherence.md")) {
		l.check(fmt.Sprintf("self-coherence.md (small-change #%s)", issueNum), checkPass, "")
		l.validateSections(filepath.Join(cycleDir, "self-coherence.md"), "self-coherence.md", "", false)
	} else {
		l.check(fmt.Sprintf("self-coherence.md (small-change #%s)", issueNum), checkWarn, "optional for small-change cycles")
	}
	if fileExists(filepath.Join(cycleDir, "alpha-closeout.md")) {
		l.check(fmt.Sprintf("alpha-closeout.md (small-change #%s)", issueNum), checkPass, "")
		l.validateSections(filepath.Join(cycleDir, "alpha-closeout.md"), "alpha-closeout.md", "", false)
	} else {
		l.check(fmt.Sprintf("alpha-closeout.md (small-change #%s)", issueNum), checkWarn, "optional for small-change cycles")
	}
	for _, artifact := range []string{"beta-review.md", "beta-closeout.md", "gamma-closeout.md"} {
		if fileExists(filepath.Join(cycleDir, artifact)) {
			l.check(fmt.Sprintf("%s (small-change #%s)", artifact, issueNum), checkWarn, "unexpected in small-change cycle - may be misclassified")
		}
	}
}

func (l *ledgerRun) checkTriadicArtifacts(cycleDir, issueNum, version string) {
	required := []string{
		"self-coherence.md", "beta-review.md", "alpha-closeout.md",
		"beta-closeout.md", "gamma-closeout.md",
	}
	isLegacy := version != "" && isLegacyVersion(version)
	for _, artifact := range required {
		ap := filepath.Join(cycleDir, artifact)
		if fileExists(ap) {
			l.check(fmt.Sprintf("%s (issue #%s)", artifact, issueNum), checkPass, "")
			l.validateSections(ap, artifact, version, false)
			continue
		}
		var exPath string
		if version != "" {
			exPath = fmt.Sprintf(".cdd/releases/%s/%s/%s", version, filepath.Base(cycleDir), artifact)
		} else {
			only := strings.SplitN(issueNum, " ", 2)[0]
			exPath = fmt.Sprintf(".cdd/unreleased/%s/%s", only, artifact)
		}
		switch {
		case l.isExceptionBacked(exPath):
			l.check(fmt.Sprintf("%s (issue #%s)", artifact, issueNum), checkWarn, fmt.Sprintf("missing but exception-backed: %s", exPath))
		case isLegacy:
			l.check(fmt.Sprintf("%s (issue #%s)", artifact, issueNum), checkWarn, fmt.Sprintf("missing in legacy cycle (pre-v3.77.0) — expected .cdd/unreleased/%s/%s", issueNum, artifact))
		default:
			l.check(fmt.Sprintf("%s (issue #%s)", artifact, issueNum), checkFail, fmt.Sprintf("expected .cdd/unreleased/%s/%s", issueNum, artifact))
		}
	}
}

// validateSections runs the basic section presence check for an artifact.
// More lenient when forUnreleased=true (in-progress cycle).
func (l *ledgerRun) validateSections(artifactPath, artifactName, version string, forUnreleased bool) {
	data, err := os.ReadFile(artifactPath)
	if err != nil {
		return
	}
	text := string(data)
	isLegacy := version != "" && isLegacyVersion(version)
	switch artifactName {
	case "self-coherence.md":
		var missing []string
		if !sectionPresent(text, "## Gap") {
			missing = append(missing, "Gap")
		}
		if !sectionPresent(text, "## Skills") && !sectionPresent(text, "## Mode") {
			missing = append(missing, "Skills/Mode")
		}
		if !sectionPresent(text, "## ACs") && !sectionPresent(text, "## AC Coverage") {
			missing = append(missing, "ACs/AC Coverage")
		}
		if !sectionPresent(text, "## CDD Trace") {
			missing = append(missing, "CDD Trace")
		}
		if !sectionPresentIgnoreCase(text, "self-check") && !sectionPresentIgnoreCase(text, "debt") {
			missing = append(missing, "Self-check or Debt")
		}
		if len(missing) > 0 {
			label := fmt.Sprintf("%s sections", artifactName)
			detail := strings.Join(missing, " ")
			if l.isExceptionBacked(artifactPath) {
				l.check(label, checkWarn, fmt.Sprintf("missing required sections but exception-backed: %s", detail))
			} else if isLegacy || forUnreleased {
				if isLegacy {
					l.check(label, checkWarn, fmt.Sprintf("missing required sections in legacy cycle (pre-v3.77.0): %s", detail))
				} else {
					l.check(label, checkWarn, fmt.Sprintf("missing sections in unreleased cycle (may be in-progress): %s", detail))
				}
			} else {
				l.check(label, checkFail, fmt.Sprintf("missing required sections: %s", detail))
			}
			return
		}
	case "beta-review.md":
		if !containsAnyIgnoreCase(text, []string{"verdict", "review", "findings", "rc", "approved"}) {
			l.check(fmt.Sprintf("%s sections", artifactName), checkWarn, "may be missing review verdict/findings (basic check)")
			return
		}
	case "alpha-closeout.md", "beta-closeout.md", "gamma-closeout.md":
		if !containsAnyIgnoreCase(text, []string{"summary", "findings", "cycle", "close-out", "closeout", "close out"}) {
			l.check(fmt.Sprintf("%s sections", artifactName), checkWarn, "may be missing close-out structure (basic check)")
			return
		}
	default:
		return
	}
	l.check(fmt.Sprintf("%s sections", artifactName), checkPass, "basic section validation passed")
}

// --- Small file/git helpers ------------------------------------------------

func fileExists(p string) bool {
	info, err := os.Stat(p)
	return err == nil && !info.IsDir()
}

func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func sectionPresent(text, header string) bool {
	for _, line := range strings.Split(text, "\n") {
		if line == header || strings.HasPrefix(line, header+" ") {
			return true
		}
	}
	return false
}

func sectionPresentIgnoreCase(text, needle string) bool {
	low := strings.ToLower(text)
	low2 := strings.ToLower(needle)
	for _, line := range strings.Split(low, "\n") {
		if strings.HasPrefix(line, "##") && strings.Contains(line, low2) {
			return true
		}
	}
	return false
}

func containsAnyIgnoreCase(text string, needles []string) bool {
	low := strings.ToLower(text)
	for _, n := range needles {
		if strings.Contains(low, strings.ToLower(n)) {
			return true
		}
	}
	return false
}

func gitTagExists(ctx context.Context, repoRoot, tag string) bool {
	cmd := exec.CommandContext(ctx, "git", "-C", repoRoot, "tag", "-l", tag)
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == tag {
			return true
		}
	}
	return false
}

func gitCurrentBranch(ctx context.Context, repoRoot string) string {
	cmd := exec.CommandContext(ctx, "git", "-C", repoRoot, "branch", "--show-current")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func lookupPRForVersion(ctx context.Context, version string) string {
	if _, err := exec.LookPath("gh"); err != nil {
		return ""
	}
	cmd := exec.CommandContext(ctx, "gh", "pr", "list", "--state", "merged", "--search", version, "--json", "number", "-q", ".[0].number")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func lookupPRState(ctx context.Context, prNum string) string {
	if _, err := exec.LookPath("gh"); err != nil {
		return "UNKNOWN"
	}
	cmd := exec.CommandContext(ctx, "gh", "pr", "view", prNum, "--json", "state", "-q", ".state")
	out, err := cmd.Output()
	if err != nil {
		return "UNKNOWN"
	}
	s := strings.TrimSpace(string(out))
	if s == "" {
		return "UNKNOWN"
	}
	return s
}
