package activation

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"

	pkgtypes "github.com/usurobor/cnos/src/go/internal/pkg"
)

// Entry is one row in the public skill activation index. The caller
// (runtime contract, status, help) consumes these as the authoritative
// list of public skills the hub has installed.
type Entry struct {
	SkillID  string
	Package  string
	Summary  string
	Triggers []string
}

// Skill is a discovered skill: the originating package, the skill ID
// (path relative to <pkg>/skills/ with forward slashes), the absolute
// path to its SKILL.md, and the parsed frontmatter. Reported by
// Discover for both index construction and validation.
type Skill struct {
	Package     string
	SkillID     string
	Path        string
	Frontmatter Frontmatter
	// ReadErr is non-nil when the SKILL.md was found on disk but could
	// not be read or parsed. Discover still surfaces the entry so the
	// validator can flag it.
	ReadErr error
}

// IssueKind classifies an activation validation issue. Mirrors the
// three-category shape of cn_activation.issue_kind in OCaml.
type IssueKind int

const (
	// IssueMissingSkill — a SKILL.md was found on disk but could not
	// be read or parsed (unreadable / malformed). Filesystem discovery
	// makes "declared but absent" unrepresentable: if the file is
	// absent, the skill is not a skill.
	IssueMissingSkill IssueKind = iota
	// IssueEmptyTriggers — SKILL.md parses but declares no triggers.
	IssueEmptyTriggers
	// IssueTriggerConflict — the same trigger keyword is claimed by
	// two or more distinct skill IDs.
	IssueTriggerConflict
)

// IssueKindLabel returns the short name used in doctor output.
func IssueKindLabel(k IssueKind) string {
	switch k {
	case IssueMissingSkill:
		return "missing"
	case IssueEmptyTriggers:
		return "empty"
	case IssueTriggerConflict:
		return "conflict"
	default:
		return "unknown"
	}
}

// Issue is one problem found by Validate.
type Issue struct {
	Kind    IssueKind
	Message string
}

// vendorPackagesDir returns the directory where installed packages live.
// Kept private so the activation package does not duplicate pkg.VendorPath.
func vendorPackagesDir(hubPath string) string {
	return filepath.Join(hubPath, ".cn", "vendor", "packages")
}

// ReadFrontmatter reads a SKILL.md file from disk and parses its
// frontmatter. The IO wrapper around ParseFrontmatter (eng/go §2.17).
// Returns an error only if the file cannot be read; malformed content
// yields a zero/partial Frontmatter with nil error.
func ReadFrontmatter(path string) (Frontmatter, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Frontmatter{}, fmt.Errorf("read skill frontmatter %s: %w", path, err)
	}
	return ParseFrontmatter(data), nil
}

// Discover walks every installed package under <hub>/.cn/vendor/packages
// and returns one Skill record per SKILL.md file found on disk.
//
// Filesystem presence is authoritative — the `skills` field in
// cn.package.json is not consulted (DESIGN-CONSTRAINTS.md §1,
// PACKAGE-SYSTEM.md §3, issue #261).
//
// The result is sorted by (Package, SkillID) for deterministic
// downstream output (eng/go §2.13).
func Discover(hubPath string) []Skill {
	vendor := vendorPackagesDir(hubPath)
	entries, err := os.ReadDir(vendor)
	if err != nil {
		slog.Debug("activation: cannot read vendor dir",
			slog.String("path", vendor),
			slog.String("err", err.Error()))
		return nil
	}

	var skills []Skill
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		pkgName := entry.Name()
		pkgDir := filepath.Join(vendor, pkgName)
		skills = append(skills, discoverPackageSkills(pkgName, pkgDir)...)
	}

	slices.SortFunc(skills, func(a, b Skill) int {
		if a.Package != b.Package {
			return strings.Compare(a.Package, b.Package)
		}
		return strings.Compare(a.SkillID, b.SkillID)
	})
	return skills
}

// discoverPackageSkills walks a single installed package's skills/
// subtree and returns one Skill per SKILL.md file found.
func discoverPackageSkills(pkgName, pkgDir string) []Skill {
	skillsRoot := filepath.Join(pkgDir, pkgtypes.ClassSkills)
	info, err := os.Stat(skillsRoot)
	if err != nil || !info.IsDir() {
		return nil
	}

	var out []Skill
	_ = filepath.WalkDir(skillsRoot, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			slog.Debug("activation: walk error",
				slog.String("package", pkgName),
				slog.String("path", path),
				slog.String("err", walkErr.Error()))
			return nil
		}
		if d.IsDir() || d.Name() != "SKILL.md" {
			return nil
		}
		skillDir := filepath.Dir(path)
		rel, err := filepath.Rel(skillsRoot, skillDir)
		if err != nil {
			slog.Debug("activation: cannot compute skill id",
				slog.String("package", pkgName),
				slog.String("path", path),
				slog.String("err", err.Error()))
			return nil
		}
		skillID := filepath.ToSlash(rel)
		if skillID == "." || skillID == "" {
			// Contract: a skill's ID is the path of its SKILL.md's
			// containing directory relative to skills/. A SKILL.md
			// directly under skills/ has no containing subdirectory
			// and therefore no derivable ID, so it is intentionally
			// not a discoverable skill. Authors must place each skill
			// under a named subdirectory (skills/<name>/SKILL.md or
			// deeper, e.g. skills/cdd/alpha/SKILL.md → id "cdd/alpha").
			slog.Debug("activation: SKILL.md directly under skills/ has no id; place under a named subdirectory",
				slog.String("package", pkgName),
				slog.String("path", path))
			return nil
		}
		fm, rerr := ReadFrontmatter(path)
		out = append(out, Skill{
			Package:     pkgName,
			SkillID:     skillID,
			Path:        path,
			Frontmatter: fm,
			ReadErr:     rerr,
		})
		return nil
	})
	return out
}

// BuildIndex returns the public activation index: one Entry per
// discovered skill whose visibility is not "internal". Skills with
// unreadable SKILL.md files are excluded from the index (the validator
// surfaces them as IssueMissingSkill).
//
// The slice is sorted by (Package, SkillID).
func BuildIndex(hubPath string) []Entry {
	skills := Discover(hubPath)
	out := make([]Entry, 0, len(skills))
	for _, s := range skills {
		if s.ReadErr != nil {
			continue
		}
		if !s.Frontmatter.IsPublic() {
			continue
		}
		out = append(out, Entry{
			SkillID:  s.SkillID,
			Package:  s.Package,
			Summary:  s.Frontmatter.Description,
			Triggers: s.Frontmatter.Triggers,
		})
	}
	return out
}

// Validate is the IO wrapper: discovers skills under hubPath and
// returns issues via ValidateSkills. Pure logic lives in ValidateSkills
// so the same categorisation can be unit-tested without a filesystem.
func Validate(hubPath string) []Issue {
	return ValidateSkills(Discover(hubPath))
}

// ValidateSkills is the pure validator over an already-discovered
// skill list. Returns one Issue per distinct problem:
//   - SKILL.md present but unreadable → IssueMissingSkill
//   - SKILL.md parses but declares no triggers → IssueEmptyTriggers
//   - same trigger claimed by two or more distinct PUBLIC skill IDs →
//     IssueTriggerConflict
//
// Internal skills (visibility: internal) are intentionally excluded
// from conflict detection. They never appear in the public activation
// index (see BuildIndex), so an internal sub-skill sharing a trigger
// keyword with its parent orchestrator is not a user-visible
// ambiguity — the runtime activates the public orchestrator, which
// then delegates internally. This mirrors OCaml-era semantics where
// activation.validate iterated only manifest-exposed skill IDs.
//
// Output is deterministic: missing/empty issues follow input order
// (Discover sorts by package/skill); conflict issues are sorted by
// the offending trigger keyword.
func ValidateSkills(skills []Skill) []Issue {
	var issues []Issue

	// triggerClaimants records which PUBLIC skill IDs claim a trigger,
	// de-duplicated, insertion-ordered. Internal skills are skipped.
	type triggerClaimants struct {
		ids  []string
		seen map[string]struct{}
	}
	byTrigger := map[string]*triggerClaimants{}

	for _, s := range skills {
		if s.ReadErr != nil {
			issues = append(issues, Issue{
				Kind: IssueMissingSkill,
				Message: fmt.Sprintf("package %s: skill %s has unreadable SKILL.md: %v",
					s.Package, s.SkillID, s.ReadErr),
			})
			continue
		}
		if len(s.Frontmatter.Triggers) == 0 {
			issues = append(issues, Issue{
				Kind: IssueEmptyTriggers,
				Message: fmt.Sprintf("package %s: skill %s has no triggers",
					s.Package, s.SkillID),
			})
			continue
		}
		// Internal skills do not participate in public conflict
		// detection — they are not addressable via the activation
		// index, so they cannot ambiguate activation.
		if !s.Frontmatter.IsPublic() {
			continue
		}
		for _, t := range s.Frontmatter.Triggers {
			claim, ok := byTrigger[t]
			if !ok {
				claim = &triggerClaimants{seen: map[string]struct{}{}}
				byTrigger[t] = claim
			}
			if _, dup := claim.seen[s.SkillID]; dup {
				continue
			}
			claim.seen[s.SkillID] = struct{}{}
			claim.ids = append(claim.ids, s.SkillID)
		}
	}

	conflictingTriggers := make([]string, 0)
	for t, claim := range byTrigger {
		if len(claim.ids) > 1 {
			conflictingTriggers = append(conflictingTriggers, t)
		}
	}
	slices.Sort(conflictingTriggers)
	for _, t := range conflictingTriggers {
		ids := slices.Clone(byTrigger[t].ids)
		slices.Sort(ids)
		issues = append(issues, Issue{
			Kind:    IssueTriggerConflict,
			Message: fmt.Sprintf("trigger %s claimed by: %s", t, strings.Join(ids, ", ")),
		})
	}

	return issues
}
