// cnos#671 R11 — WAVE-BOUNDARY pre-authorization validator (Go + CUE repo; NO Python).
//
// Materializes the whole-wave oracle-ownership / classification BIJECTION named by
// oracle-registry.yaml `wave_predicates.wave_oracle_ownership_bijection_enforced`
// (deferred_owner: "wave", wave_authorization_gated: true).
//
// A PRE-AUTHORIZATION gate runs BEFORE any Working Cell executes, so it cannot be
// deferred to a WC — it must be a real, runnable validator in the plan matter. This
// program is that validator: credential-free, single-file `package main`, standard
// library ONLY (no external deps, no go.mod required, `go run` works anywhere).
//
// It proves, over the six child contracts + the authoritative registry:
//
//	(1) BIJECTION — union(child acceptance.predicates) over (owner, predicate) ⇄ the
//	    registry `assurance:` entries over (owner, predicate), EXACTLY: no missing
//	    (a child predicate with no registry entry), no phantom (a registry entry with
//	    no child predicate), no duplicate on either side. (78 ⇄ 78 for the real wave.)
//	(2) OWNERSHIP — every `mechanically-verifiable` predicate binds EXACTLY ONE owner:
//	    exactly one concrete checker|schema path, and its (owner, predicate) is unique
//	    (a predicate "owned twice" is rejected).
//
// Prints a clear result; EXIT 0 IFF bijective, non-zero otherwise.
//
// Usage:
//
//	go run oracle_ownership_bijection.go <path>
//	  <path> = a WAVE DIRECTORY  (contains contracts/*.cn-cell-contract-v1.yaml +
//	           oracle-registry.yaml) — the real wave; OR
//	  <path> = a self-contained FIXTURE FILE (top-level `child_predicates:` +
//	           `assurance:` blocks) — a minimal registry+predicate set.
//
// Examples:
//
//	go run oracle_ownership_bijection.go ..                                  # real wave (from wave-validators/)
//	go run oracle_ownership_bijection.go fixtures/oracle-ownership.one-checker-each.positive.yaml
//	go run oracle_ownership_bijection.go fixtures/oracle-ownership.double-owned.negative.yaml
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// A (owner, predicate) key — the bijection is over these pairs (the same predicate
// STRING may legitimately appear under different owners, so ownership is per-pair).
type pair struct {
	owner     string
	predicate string
}

func (p pair) String() string { return p.owner + " :: " + p.predicate }

// An assurance/registry entry (only the fields the bijection needs).
type entry struct {
	owner          string
	predicate      string
	classification string
	checker        string
	schema         string
}

type model struct {
	mode      string // "wave-dir" | "fixture-file"
	child     []pair
	assurance []entry
}

func indentOf(line string) int {
	n := 0
	for _, c := range line {
		if c == ' ' {
			n++
		} else {
			break
		}
	}
	return n
}

func isBlankOrComment(line string) bool {
	t := strings.TrimSpace(line)
	return t == "" || strings.HasPrefix(t, "#")
}

// scalarValue extracts a YAML scalar from the text after a "key:" — the double-quoted
// content if quoted, otherwise the bare token with any trailing " # comment" stripped.
func scalarValue(raw string) string {
	s := strings.TrimSpace(raw)
	if strings.HasPrefix(s, "\"") {
		if j := strings.IndexByte(s[1:], '"'); j >= 0 {
			return s[1 : 1+j]
		}
		return s[1:]
	}
	if s == "" || strings.HasPrefix(s, "#") {
		return ""
	}
	if i := strings.Index(s, " #"); i >= 0 {
		s = s[:i]
	}
	return strings.TrimSpace(s)
}

// addField parses a "key: value" line into m (unknown keys are kept but harmless).
func addField(m map[string]string, kv string) {
	idx := strings.Index(kv, ":")
	if idx < 0 {
		return
	}
	key := strings.TrimSpace(kv[:idx])
	if key == "" {
		return
	}
	m[key] = scalarValue(kv[idx+1:])
}

// parseEntries reads a top-level (indent-0) YAML block `blockKey:` whose members are
// list entries `  - key: val` with `    key: val` continuation fields. It stops at the
// next indent-0 non-comment line, so it never bleeds into a following block.
func parseEntries(lines []string, blockKey string) []map[string]string {
	var out []map[string]string
	var cur map[string]string
	inBlock := false
	flush := func() {
		if cur != nil {
			out = append(out, cur)
			cur = nil
		}
	}
	for _, line := range lines {
		if !inBlock {
			if indentOf(line) == 0 && strings.TrimRight(line, " \t") == blockKey+":" {
				inBlock = true
			}
			continue
		}
		if isBlankOrComment(line) {
			continue
		}
		if indentOf(line) == 0 { // a new top-level key ends the block
			break
		}
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "- ") {
			flush()
			cur = map[string]string{}
			addField(cur, strings.TrimPrefix(trimmed, "- "))
		} else if cur != nil {
			addField(cur, trimmed)
		}
	}
	flush()
	return out
}

// parseAcceptancePredicates extracts the acceptance.predicates string list from a
// contract file's lines (`acceptance:` at indent 0 → `predicates:` at indent 2 →
// `- "..."` items at indent 4; comments between items skipped).
func parseAcceptancePredicates(lines []string) []string {
	var out []string
	inAccept := false
	inPreds := false
	for _, line := range lines {
		if !inAccept {
			if indentOf(line) == 0 && strings.TrimRight(line, " \t") == "acceptance:" {
				inAccept = true
			}
			continue
		}
		if indentOf(line) == 0 && !isBlankOrComment(line) {
			break // acceptance block ended
		}
		if isBlankOrComment(line) {
			continue
		}
		if !inPreds {
			if strings.TrimSpace(line) == "predicates:" {
				inPreds = true
			}
			continue
		}
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "- ") {
			out = append(out, scalarValue(strings.TrimPrefix(trimmed, "- ")))
		} else if indentOf(line) <= 2 {
			break // dedented out of the predicates list
		}
	}
	return out
}

func readLines(path string) ([]string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(b), "\n"), nil
}

func toEntries(maps []map[string]string) []entry {
	var out []entry
	for _, m := range maps {
		out = append(out, entry{
			owner:          m["owner"],
			predicate:      m["predicate"],
			classification: m["classification"],
			checker:        m["checker"],
			schema:         m["schema"],
		})
	}
	return out
}

func loadWaveDir(dir string) (*model, error) {
	glob := filepath.Join(dir, "contracts", "wc-*.cn-cell-contract-v1.yaml")
	files, err := filepath.Glob(glob)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no child contracts matched %s", glob)
	}
	sort.Strings(files)
	m := &model{mode: "wave-dir"}
	for _, f := range files {
		base := filepath.Base(f)
		owner := strings.SplitN(base, ".", 2)[0] // "wc-1.cn-cell-contract-v1.yaml" -> "wc-1"
		lines, err := readLines(f)
		if err != nil {
			return nil, err
		}
		for _, p := range parseAcceptancePredicates(lines) {
			m.child = append(m.child, pair{owner: owner, predicate: p})
		}
	}
	regPath := filepath.Join(dir, "oracle-registry.yaml")
	regLines, err := readLines(regPath)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", regPath, err)
	}
	m.assurance = toEntries(parseEntries(regLines, "assurance"))
	if len(m.assurance) == 0 {
		return nil, fmt.Errorf("no assurance entries parsed from %s", regPath)
	}
	return m, nil
}

func loadFixture(path string) (*model, error) {
	lines, err := readLines(path)
	if err != nil {
		return nil, err
	}
	m := &model{mode: "fixture-file"}
	for _, cm := range parseEntries(lines, "child_predicates") {
		m.child = append(m.child, pair{owner: cm["owner"], predicate: cm["predicate"]})
	}
	m.assurance = toEntries(parseEntries(lines, "assurance"))
	if len(m.child) == 0 && len(m.assurance) == 0 {
		return nil, fmt.Errorf("fixture %s has neither child_predicates nor assurance blocks", path)
	}
	return m, nil
}

func load(path string) (*model, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		return loadWaveDir(path)
	}
	return loadFixture(path)
}

func sortedPairs(m map[pair]bool) []pair {
	var ps []pair
	for p := range m {
		ps = append(ps, p)
	}
	sort.Slice(ps, func(i, j int) bool {
		if ps[i].owner != ps[j].owner {
			return ps[i].owner < ps[j].owner
		}
		return ps[i].predicate < ps[j].predicate
	})
	return ps
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run oracle_ownership_bijection.go <wave-dir|fixture-file>")
		os.Exit(2)
	}
	path := os.Args[1]
	m, err := load(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load error: %v\n", err)
		os.Exit(2)
	}

	// Child side: (owner, predicate) multiset.
	childCount := map[pair]int{}
	for _, p := range m.child {
		childCount[p]++
	}
	// Registry side: (owner, predicate) multiset + classification/checker lookup.
	assrCount := map[pair]int{}
	mechCount := map[pair]int{}    // mechanically-verifiable entries per pair
	mechConcrete := map[pair]int{} // # of concrete checker|schema bound per mech pair
	for _, e := range m.assurance {
		k := pair{owner: e.owner, predicate: e.predicate}
		assrCount[k]++
		if e.classification == "mechanically-verifiable" {
			mechCount[k]++
			n := 0
			if strings.TrimSpace(e.checker) != "" {
				n++
			}
			if strings.TrimSpace(e.schema) != "" {
				n++
			}
			mechConcrete[k] += n
		}
	}

	childSet := map[pair]bool{}
	for p := range childCount {
		childSet[p] = true
	}
	assrSet := map[pair]bool{}
	for p := range assrCount {
		assrSet[p] = true
	}

	// Bijection defects.
	missing := map[pair]bool{}     // child -> no registry entry
	phantom := map[pair]bool{}     // registry -> no child predicate
	childDup := map[pair]bool{}    // duplicate child (owner,predicate)
	doubleOwned := map[pair]bool{} // duplicate registry entry — a predicate "owned twice"
	for p, c := range childCount {
		if c > 1 {
			childDup[p] = true
		}
		if !assrSet[p] {
			missing[p] = true
		}
	}
	for p, c := range assrCount {
		if c > 1 {
			doubleOwned[p] = true
		}
		if !childSet[p] {
			phantom[p] = true
		}
	}
	// Ownership: every mechanically-verifiable predicate binds EXACTLY ONE checker|schema.
	mechBadOwner := map[pair]bool{}
	for p, c := range mechCount {
		// c>1 already flagged by doubleOwned; the concrete count must be exactly one PER entry.
		if mechConcrete[p] != c {
			mechBadOwner[p] = true
		}
	}

	bijective := len(missing) == 0 && len(phantom) == 0 && len(childDup) == 0 &&
		len(doubleOwned) == 0 && len(mechBadOwner) == 0

	// ---- report ----
	fmt.Println("wave-oracle-ownership-bijection validator (cnos#671 R11, wave-boundary pre-authorization)")
	fmt.Printf("  input:                        %s  (mode: %s)\n", path, m.mode)
	fmt.Printf("  child acceptance predicates:  %d\n", len(m.child))
	fmt.Printf("  registry assurance entries:   %d\n", len(m.assurance))
	fmt.Printf("  mechanically-verifiable:      %d\n", len(mechCount))
	fmt.Printf("  missing   (child -> registry): %d\n", len(missing))
	fmt.Printf("  phantom   (registry -> child): %d\n", len(phantom))
	fmt.Printf("  child duplicates:              %d\n", len(childDup))
	fmt.Printf("  double-owned (registry dup):   %d\n", len(doubleOwned))
	fmt.Printf("  mech missing/≠1 checker owner: %d\n", len(mechBadOwner))

	report := func(label string, s map[pair]bool) {
		if len(s) == 0 {
			return
		}
		fmt.Printf("  --- %s ---\n", label)
		for _, p := range sortedPairs(s) {
			fmt.Printf("      %s\n", p)
		}
	}
	report("MISSING (child acceptance predicate with no registry entry)", missing)
	report("PHANTOM (registry entry with no child acceptance predicate)", phantom)
	report("DUPLICATE child (owner,predicate)", childDup)
	report("DOUBLE-OWNED registry (owner,predicate) — a predicate owned twice", doubleOwned)
	report("mechanically-verifiable NOT bound to exactly one checker|schema", mechBadOwner)

	fmt.Printf("  bijective: %v\n", bijective)
	if bijective {
		fmt.Println("RESULT: PASS — oracle-ownership / classification bijection holds (exit 0)")
		os.Exit(0)
	}
	fmt.Println("RESULT: FAIL — oracle-ownership / classification bijection broken (exit 1)")
	os.Exit(1)
}
