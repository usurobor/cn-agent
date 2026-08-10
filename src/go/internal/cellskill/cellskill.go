// Package cellskill resolves canonical skill references and loads their
// bodies. A reference names an installed skill — `<package>:<path>`, e.g.
// `cnos.eng:eng/go` — and loading means the actual SKILL.md content, digested,
// so a closure can state exactly which skill text was injected into a seat.
// Printing a skill's name into a prompt is not loading it.
//
// Roots resolve like `$PATH`: an ordered list, first match wins, supplied by
// the composition root from a hub or repository anchor. Nothing here is
// relative to the process working directory, and no filesystem path ever
// reaches a receipt — the content digest is the identity.
package cellskill

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Skill is one resolved, loaded skill: the canonical reference, the exact
// body that was injected, and the body's digest for the record.
type Skill struct {
	Ref    string
	Body   string
	SHA256 string
}

// Resolver loads skills. The consumer-owned port keeps fill constructors
// testable without a package tree on disk.
type Resolver interface {
	Load(ref string) (Skill, error)
}

// Tree resolves refs against an ordered list of package roots:
// `<pkg>:<path>` → `<root>/<pkg>/skills/<path>/SKILL.md`, first root that
// has the file wins. Roots must be absolute.
type Tree struct {
	Roots []string
}

func (t Tree) Load(ref string) (Skill, error) {
	pkg, path, ok := strings.Cut(ref, ":")
	if !ok || pkg == "" || path == "" {
		return Skill{}, fmt.Errorf("cellskill: ref %q is not <package>:<path>", ref)
	}
	if strings.Contains(pkg, "..") || strings.Contains(path, "..") ||
		strings.ContainsRune(pkg, filepath.Separator) || filepath.IsAbs(path) {
		return Skill{}, fmt.Errorf("cellskill: ref %q escapes the package tree", ref)
	}
	if len(t.Roots) == 0 {
		return Skill{}, fmt.Errorf("cellskill: no package roots configured")
	}

	tried := make([]string, 0, len(t.Roots))
	for _, root := range t.Roots {
		file := filepath.Join(root, pkg, "skills", filepath.FromSlash(path), "SKILL.md")
		data, err := os.ReadFile(file)
		if err != nil {
			tried = append(tried, file)
			continue
		}
		sum := sha256.Sum256(data)
		return Skill{Ref: ref, Body: string(data), SHA256: hex.EncodeToString(sum[:])}, nil
	}
	return Skill{}, fmt.Errorf("cellskill: skill %q is not installed (looked in: %s)", ref, strings.Join(tried, ", "))
}

// LoadAll loads refs in order; order is meaning (later skills refine earlier
// ones) and the record preserves it.
func LoadAll(r Resolver, refs []string) ([]Skill, error) {
	out := make([]Skill, 0, len(refs))
	for _, ref := range refs {
		s, err := r.Load(ref)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, nil
}

// Roots is the resolution order for one environment anchor: installed
// packages first (the canonical identity), then the source tree so a
// development checkout that has not vendored anything still resolves. The
// anchor is a hub or repository root — never the process working directory.
func Roots(anchor string) []string {
	return []string{
		filepath.Join(anchor, ".cn", "vendor", "packages"),
		filepath.Join(anchor, "src", "packages"),
	}
}
