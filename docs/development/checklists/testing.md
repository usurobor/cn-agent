# Testing Checklist

From `skills/testing/SKILL.md`

| Item | Severity |
|------|----------|
| New functions have tests | C |
| Edge cases covered (empty, missing, malformed) | C |
| `go test ./...` passes (from `src/go/`) | D |
| No `[%expect]` mismatches | D |
| Parsing roundtrips tested | B |
| Critical paths have integration tests | B |
