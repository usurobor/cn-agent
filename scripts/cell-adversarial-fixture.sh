#!/bin/bash
# a small, real Go repo — the checker runs fast on it
set -e
d=$1; rm -rf "$d"; mkdir -p "$d/src/go/calc"; cd "$d"; git init -q -b main
cat > src/go/go.mod <<'EOF'
module example.com/calc

go 1.24
EOF
cat > src/go/calc/calc.go <<'EOF'
package calc

// Sum adds every value in xs.
func Sum(xs []int) int {
	total := 0
	for _, x := range xs {
		total += x
	}
	return total
}
EOF
cat > src/go/calc/calc_test.go <<'EOF'
package calc

import "testing"

func TestSum(t *testing.T) {
	if got := Sum([]int{1, 2, 3}); got != 6 {
		t.Fatalf("Sum = %d, want 6", got)
	}
}
EOF
git add -A
GIT_AUTHOR_NAME=t GIT_AUTHOR_EMAIL=t@t GIT_COMMITTER_NAME=t GIT_COMMITTER_EMAIL=t@t git commit -qm base
