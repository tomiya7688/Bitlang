# Bitlang Validation Routing

Run the smallest validation that proves the changed contract, then broaden only when shared behavior changed.

| Change scope | Minimum validation | Broader fallback |
| --- | --- | --- |
| one lexer/token rule | matching lexer test | `go test ./...` if token contract changed |
| names/symbol behavior | matching name/symbol tests | full package tests |
| one pipeline rule | pipeline tests | `go test ./...` if stage contract changed |
| Go rule checker | checker package tests + checker self-run | `go test ./...` |
| CUI only | CUI build / targeted test | `go test ./...` |
| shared compiler API | matching tests | `go test ./...` + `go build ./cmd/bitlang` |
| docs/routing only | link/path consistency review | no full test required |

## Standard Go Completion
When implementation code changes and the environment permits:

```text
go run ./tools/go-rule-checker/cmd/go-rule-checker .
go test ./...
go build ./cmd/bitlang
```

## Rules
- Do not run unrelated visual/runtime validation for pure compiler-core changes.
- Do not treat a successful build as a replacement for semantic tests.
- For shared stage representations, validate consumers of that representation.
- If a required check cannot be run, report it as Unverified instead of expanding unrelated investigation.
- Keep failure logs bounded to the relevant error region.
