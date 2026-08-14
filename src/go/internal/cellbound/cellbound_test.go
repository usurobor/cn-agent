package cellbound

import (
	"strings"
	"testing"
	"unicode/utf8"
)

// The invariant the policies exist for is WHICH bytes survive, not that
// something was dropped. Both tests below assert the surviving bytes exactly,
// so substituting one policy for the other fails them — which is the only way
// a caller that picked the wrong half of its evidence gets caught.

func TestKeepHeadKeepsTheFirstBytes(t *testing.T) {
	w := New(KeepHead, 4)
	// Written in pieces, because a child writes in pieces and the bound is
	// applied as it does.
	mustWrite(t, w, "abcde", "fghij")
	if got := w.String(); got != "abcd" {
		t.Fatalf("KeepHead kept %q, want the FIRST 4 bytes %q", got, "abcd")
	}
	if !w.Truncated() {
		t.Fatal("a writer past its bound must report truncation")
	}
}

func TestKeepTailKeepsTheLastBytes(t *testing.T) {
	w := New(KeepTail, 4)
	mustWrite(t, w, "abcde", "fghij")
	if got := w.String(); got != "ghij" {
		t.Fatalf("KeepTail kept %q, want the LAST 4 bytes %q", got, "ghij")
	}
	if !w.Truncated() {
		t.Fatal("a writer past its bound must report truncation")
	}
}

// Exactly at the bound is the boundary both policies must agree on: nothing was
// dropped, so neither may claim it was, and the two must return the same bytes.
func TestAtTheBoundNeitherPolicyDropsAnything(t *testing.T) {
	for name, policy := range map[string]Policy{"KeepHead": KeepHead, "KeepTail": KeepTail} {
		t.Run(name, func(t *testing.T) {
			w := New(policy, 10)
			mustWrite(t, w, "abcde", "fghij")
			if got := w.String(); got != "abcdefghij" {
				t.Fatalf("kept %q, want the whole input", got)
			}
			if w.Truncated() {
				t.Fatal("a writer that dropped nothing must not report truncation")
			}
		})
	}
}

// A short write would make the child die of a broken pipe, turning a loud
// process into an unavailable one. The bound must fail the caller's operation
// instead, so Write always accepts everything it was handed.
func TestWriteNeverReportsAShortWrite(t *testing.T) {
	for name, policy := range map[string]Policy{"KeepHead": KeepHead, "KeepTail": KeepTail} {
		t.Run(name, func(t *testing.T) {
			w := New(policy, 2)
			p := []byte("far past the bound")
			n, err := w.Write(p)
			if n != len(p) || err != nil {
				t.Fatalf("Write = (%d, %v), want (%d, nil)", n, err, len(p))
			}
		})
	}
}

// String is the surviving bytes and nothing else: a caller reapplying a diff or
// parsing a stream must not find our sentence inside its value.
func TestStringDoesNotSpliceInTheMarker(t *testing.T) {
	w := New(KeepTail, 4)
	mustWrite(t, w, "abcdefghij")
	if strings.Contains(w.String(), Marker) {
		t.Fatalf("String must carry no marker, got %q", w.String())
	}
}

func TestTailKeepsTheLastBytesAndAnnouncesTheLoss(t *testing.T) {
	got := Tail("abcdefghij", 4)
	if got != Marker+"ghij" {
		t.Fatalf("Tail = %q, want the announced LAST 4 bytes %q", got, Marker+"ghij")
	}
}

func TestTailReturnsShortInputUntouchedButTrimmed(t *testing.T) {
	if got := Tail("  short  ", 40); got != "short" {
		t.Fatalf("Tail = %q, want %q", got, "short")
	}
	if got := Tail("short", 40); strings.Contains(got, Marker) {
		t.Fatalf("an untruncated tail must not announce a loss, got %q", got)
	}
}

// A tail is embedded in error messages, step tails and JSON receipts. Cutting n
// bytes off the end of a multi-byte rune must not put invalid UTF-8 into any of
// them.
func TestTailStaysValidUTF8AcrossACutRune(t *testing.T) {
	s := strings.Repeat("é", 20) // two bytes per rune
	got := Tail(s, 5)            // an odd bound guarantees a cut mid-rune
	if !utf8.ValidString(got) {
		t.Fatalf("Tail returned invalid UTF-8: %q", got)
	}
	if !strings.HasPrefix(got, Marker) {
		t.Fatalf("Tail = %q, want the loss announced with %q", got, Marker)
	}
	if !strings.HasSuffix(got, "éé") {
		t.Fatalf("Tail = %q, want the END of the input", got)
	}
}

func mustWrite(t *testing.T, w *Writer, chunks ...string) {
	t.Helper()
	for _, c := range chunks {
		if _, err := w.Write([]byte(c)); err != nil {
			t.Fatalf("write %q: %v", c, err)
		}
	}
}
