// Package cellbound bounds output whose size nothing can predict.
//
// Every place a cell rents a child process — a provider adapter, a checker
// step, a git read — faces the same problem: the child may emit more than
// memory, and a limit applied to a fully buffered result is not a limit. This
// package owns that one mechanism, so a fix to it is a fix everywhere rather
// than a fix in one of three copies.
//
// The only real decision a caller makes is WHICH bytes it needs when the bound
// bites — the first ones, because it will parse the stream from the start or
// reject it outright, or the last ones, because a failing build says what went
// wrong at the end. That decision is the Policy; nothing else varies.
package cellbound

import "strings"

// Marker is the one sentence a bounded value uses to say bytes were dropped.
// One text rather than one per call site: a reader who has met it once
// recognises it in a step tail, an error, and a receipt alike.
const Marker = "...[truncated]...\n"

// Policy names which bytes survive the bound.
type Policy int

const (
	// KeepHead keeps the FIRST max bytes. For callers whose value is content:
	// they either consume it from the start or refuse it for being over bound,
	// and Truncated is what they act on.
	KeepHead Policy = iota
	// KeepTail keeps the LAST max bytes. For callers whose value is evidence a
	// human reads, where the end is the part that explains the outcome.
	KeepTail
)

// Writer keeps at most max bytes of what is written to it under its policy and
// remembers that it had to drop the rest. The bound applies AS THE CHILD
// WRITES, which is the whole point: buffering everything and trimming after
// the process exits bounds the value, not the cost.
//
// It never reports a short write. A bound must fail the operation in the
// caller's own terms — an error, a refused answer, a failed step — rather than
// kill the child mid-stream with a broken pipe and turn a loud process into an
// unavailable one.
type Writer struct {
	max       int
	policy    Policy
	buf       []byte
	truncated bool
}

// New returns a Writer bounded to max bytes under policy.
func New(policy Policy, max int) *Writer {
	return &Writer{max: max, policy: policy}
}

func (w *Writer) Write(p []byte) (int, error) {
	switch w.policy {
	case KeepHead:
		switch room := w.max - len(w.buf); {
		case room >= len(p):
			w.buf = append(w.buf, p...)
		case room > 0:
			w.buf = append(w.buf, p[:room]...)
			w.truncated = true
		default:
			w.truncated = true
		}
	case KeepTail:
		w.buf = append(w.buf, p...)
		if len(w.buf) > w.max {
			w.truncated = true
			w.buf = w.buf[len(w.buf)-w.max:]
		}
	default:
		// Unreachable for the two declared policies; a third value is a
		// programmer error, and guessing head or tail would silently keep the
		// wrong half of someone's evidence.
		panic("cellbound: unknown policy")
	}
	return len(p), nil
}

// String returns the bytes that survived, and nothing else. It does NOT splice
// in Marker: a caller whose value is content — a diff to reapply, a stream to
// parse — must not find a sentence of ours inside it. Announcing the loss is
// the caller's decision, made with Truncated and Marker or by calling Tail.
func (w *Writer) String() string { return string(w.buf) }

// Truncated reports whether the bound dropped anything.
func (w *Writer) Truncated() bool { return w.truncated }

// Tail returns at most n trailing bytes of s, announced with Marker when it had
// to drop any. Trailing, because the end of a stalled or failing process is
// what says where it went wrong.
//
// It trims surrounding space and keeps the result valid UTF-8, so a cut made
// mid-rune cannot corrupt the error message, receipt, or JSON value the tail is
// embedded in.
func Tail(s string, n int) string {
	s = strings.TrimSpace(s)
	if len(s) <= n {
		return s
	}
	return Marker + strings.ToValidUTF8(s[len(s)-n:], "")
}
