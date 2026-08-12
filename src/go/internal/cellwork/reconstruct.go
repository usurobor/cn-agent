package cellwork

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"
)

// maxViewBytes bounds the file CONTENT one reconstructed view carries. It sits
// with the other bounds in this package and mirrors the kernel's matter bound,
// because a view and a matter are the two things that travel into a prompt: a
// view that could not be exceeded by any real change would not need reporting,
// and one that could grow without limit would be an unbounded value crossing a
// bounded boundary.
const maxViewBytes = 1 << 20

// FileStatus is what the matter did to one path, as git reports it. It is a
// closed set: an unrecognized status is an error rather than a default, so a
// future git output shape cannot be silently read as "modified".
type FileStatus string

const (
	FileAdded    FileStatus = "added"
	FileModified FileStatus = "modified"
	FileDeleted  FileStatus = "deleted"
	FileRenamed  FileStatus = "renamed"
)

// FileState is one path after the matter is applied.
//
// Content is the post-application text. It is empty in three different
// situations, and a reader must be able to tell them apart, so emptiness is
// never the signal: a DELETED path says so in Status, an omitted one says so in
// Omitted (with Binary saying why, when that is why), and an empty file is the
// only case where Content is empty and both flags are false.
type FileState struct {
	Path string
	// From is the previous path, and is set only for a rename.
	From    string
	Status  FileStatus
	Content string
	// Omitted says the file has content this view does not carry: it is binary,
	// or the view's byte bound was reached before this path.
	Omitted bool
	Binary  bool
	// Symlink says Content is a LINK TARGET PATH rather than file bytes. It is
	// its own field because a reader must not mistake the two: git stores the
	// target as the blob, and nothing here follows the link.
	Symlink bool
}

// View is the assessing seat's reconstructed evidence: the post-application
// state of exactly the paths the matter touches, ordered by path. It is a
// VALUE, not a directory handle — the worktree it was read from is gone before
// this is returned, so there is nothing for a later reader to reach through.
type View struct {
	Files []FileState
	// Truncated says the byte bound was hit and the view is incomplete. It is
	// reported rather than silently dropped because a reader deciding an
	// obligation against a partial view must be able to know that it is partial.
	Truncated bool
}

// Reconstruct is the third adapter operation: it derives the candidate state
// from `(subject, matter)` alone. Nothing produced it but the pinned base and
// the patch — no producer session, no workspace handle, no claim.
//
// A fresh disposable worktree is cut at baseSHA, the matter is applied to it,
// the touched paths are read back, and the worktree is released before this
// returns. The whole point is that only the value escapes.
func Reconstruct(ctx context.Context, repo, baseSHA, matter string) (View, error) {
	wt, release, err := Materialize(ctx, repo, baseSHA)
	if err != nil {
		return View{}, err
	}
	defer release()

	// A patch that does not apply is its own outcome, distinct from a bad base
	// or a bad repository, and it is the one a reviewer most needs named: it
	// means the matter it was handed does not describe a change to this tree.
	if _, err := gitInput(ctx, wt.Dir, matter, maxRefBytes, "apply", "--whitespace=nowarn", "-"); err != nil {
		return View{}, fmt.Errorf("cellwork: the matter does not apply to %s: %w", wt.BaseSHA, err)
	}
	// Staging is what makes created files visible to `diff`, exactly as in the
	// measurement direction; the index belongs to this disposable worktree.
	if _, err := git(ctx, wt.Dir, maxRefBytes, "add", "-A"); err != nil {
		return View{}, fmt.Errorf("cellwork: stage the reconstruction: %w", err)
	}
	// Against the PINNED BASE and NUL-separated: `-z` is not a nicety, it is
	// what stops git from quoting paths containing spaces or non-ASCII, which
	// would make the reported path differ from the path on disk.
	out, err := git(ctx, wt.Dir, maxDiffBytes, "diff", "--cached", "--no-color", "-z", "--name-status", wt.BaseSHA)
	if err != nil {
		return View{}, fmt.Errorf("cellwork: list the reconstructed paths: %w", err)
	}
	changes, err := parseNameStatus(out)
	if err != nil {
		return View{}, err
	}
	return readStates(wt.Dir, changes)
}

// change is one entry of `git diff --name-status -z`.
type change struct {
	status FileStatus
	from   string // rename source only
	path   string
}

// parseNameStatus decodes `--name-status -z` output. Pure — it is the half of
// reconstruction that has no worktree in it, and the half whose edge cases
// (renames carrying two paths, a trailing NUL, an empty result) are worth
// testing without cutting a repository first.
func parseNameStatus(out string) ([]change, error) {
	fields := strings.Split(strings.TrimSuffix(out, "\x00"), "\x00")
	var changes []change
	for i := 0; i < len(fields); i++ {
		code := fields[i]
		if code == "" {
			continue
		}
		// A rename or copy carries the score in the same field and TWO paths.
		twoPaths := code[0] == 'R' || code[0] == 'C'
		want := 1
		if twoPaths {
			want = 2
		}
		if i+want >= len(fields) {
			return nil, fmt.Errorf("cellwork: truncated name-status record %q", code)
		}
		c := change{path: fields[i+want]}
		if twoPaths {
			c.from = fields[i+1]
		}
		i += want
		switch code[0] {
		case 'A':
			c.status = FileAdded
		case 'M', 'T': // T is a type change; its post-application content is what a reader needs.
			c.status = FileModified
		case 'D':
			c.status = FileDeleted
		case 'R':
			c.status = FileRenamed
		case 'C':
			// A copy's destination is new, so `added` is what it is after the
			// matter; the source is untouched and therefore not in this view.
			c.status = FileAdded
		default:
			return nil, fmt.Errorf("cellwork: unhandled git status %q for %q", code, c.path)
		}
		changes = append(changes, c)
	}
	return changes, nil
}

// readStates reads the post-application content of each changed path, in path
// order. The order is imposed here rather than trusted from git, because the
// byte bound is spent in that order — an unstable order would make the same
// (subject, matter) yield different views.
func readStates(dir string, changes []change) (View, error) {
	sort.Slice(changes, func(i, j int) bool { return changes[i].path < changes[j].path })

	v := View{Files: make([]FileState, 0, len(changes))}
	total := 0
	for _, c := range changes {
		st := FileState{Path: c.path, From: c.from, Status: c.status}
		if c.status == FileDeleted {
			v.Files = append(v.Files, st)
			continue
		}
		abs := filepath.Join(dir, filepath.FromSlash(c.path))
		// Lstat, never Stat, and the distinction is load-bearing rather than
		// pedantic. A symlink's content in git is its TARGET PATH (mode
		// 120000), so reading through it reports something the matter never
		// contained: a producing seat that adds one symlink could otherwise
		// put the bytes of any host file the runtime can read straight into
		// the reviewing seat's prompt, under a header promising the view came
		// from (subject, matter) alone. That would falsify the purity claim
		// this whole reconstruction rests on, so the link is reported AS a
		// link and nothing is followed.
		info, err := os.Lstat(abs)
		if err != nil {
			// A path git says exists but the filesystem does not is a
			// reconstruction fault, not a degraded file.
			return View{}, fmt.Errorf("cellwork: stat reconstructed %s: %w", c.path, err)
		}
		if info.Mode()&os.ModeSymlink != 0 {
			target, err := os.Readlink(abs)
			if err != nil {
				return View{}, fmt.Errorf("cellwork: read link %s: %w", c.path, err)
			}
			// The target path IS the content, exactly as the blob stores it.
			// A dangling link is therefore an ordinary readable state here,
			// not the episode-ending error that following it produced: a
			// relative shim pointing outside the change is a legal patch.
			st.Symlink, st.Content = true, target
			total += len(target)
			v.Files = append(v.Files, st)
			continue
		}
		// Size before content. The bound has to be decided from metadata or it
		// is not a bound: a one-line patch to a file already large in the base
		// would otherwise pull the whole file into memory before anything
		// checked it — the failure this package already names elsewhere, that
		// "a limit checked on a fully buffered result is not a limit".
		if !info.Mode().IsRegular() || total+int(info.Size()) > maxViewBytes {
			if !info.Mode().IsRegular() {
				st.Omitted, st.Binary = true, true // device, socket, fifo: no text to carry
			} else {
				st.Omitted, v.Truncated = true, true
			}
			v.Files = append(v.Files, st)
			continue
		}
		data, err := os.ReadFile(abs)
		if err != nil {
			return View{}, fmt.Errorf("cellwork: read reconstructed %s: %w", c.path, err)
		}
		if !utf8.Valid(data) {
			// Named, not carried. The view travels as text into a prompt and the
			// kernel's artifact contract is UTF-8, so "not valid UTF-8" is the
			// operative property here — not git's own NUL-byte heuristic, which
			// would let a differently-broken encoding through as content.
			st.Binary, st.Omitted = true, true
		} else {
			st.Content = string(data)
			total += len(data)
		}
		v.Files = append(v.Files, st)
	}
	return v, nil
}
