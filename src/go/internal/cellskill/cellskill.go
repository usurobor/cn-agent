// Package cellskill resolves canonical skill references and loads their
// bodies. A reference names an installed skill — `<package>:<path>`, e.g.
// `cnos.eng:eng/go` — and loading means the actual SKILL.md content, digested,
// so a closure can state exactly which skill text was injected into a seat.
// Printing a skill's name into a prompt is not loading it.
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
	Ref    string // canonical "<package>:<path>"
	Body   string
	SHA256 string
}

// Resolver loads skills. The consumer-owned port keeps the fill constructors
// testable without a package tree on disk.
type Resolver interface {
	Load(ref string) (Skill, error)
}

// Tree resolves refs against a source/installed package tree:
// `<pkg>:<path>` → `<root>/<pkg>/skills/<path>/SKILL.md`.
type Tree struct {
	Root string // e.g. "src/packages"
}

func (t Tree) Load(ref string) (Skill, error) {
	pkg, path, ok := strings.Cut(ref, ":")
	if !ok || pkg == "" || path == "" {
		return Skill{}, fmt.Errorf("cellskill: ref %q is not <package>:<path>", ref)
	}
	if strings.Contains(pkg, "..") || strings.Contains(path, "..") || filepath.IsAbs(path) {
		return Skill{}, fmt.Errorf("cellskill: ref %q escapes the package tree", ref)
	}
	file := filepath.Join(t.Root, pkg, "skills", filepath.FromSlash(path), "SKILL.md")
	data, err := os.ReadFile(file)
	if err != nil {
		return Skill{}, fmt.Errorf("cellskill: skill %q is not installed (%s): %w", ref, file, err)
	}
	sum := sha256.Sum256(data)
	return Skill{Ref: ref, Body: string(data), SHA256: hex.EncodeToString(sum[:])}, nil
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
