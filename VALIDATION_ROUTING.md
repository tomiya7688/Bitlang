# Bitlang Validation Routing

Run the smallest validation that proves the changed contract during development, then run the strict project gate before completion.

| Change scope | Minimum validation | Broader fallback |
| --- | --- | --- |
| one lexer/token rule | matching lexer test | `go test ./...` if token contract changed |
| names/symbol behavior | matching name/symbol tests | full package tests |
| one pipeline rule | pipeline tests | `go test ./...` if stage contract changed |
| Go rule checker | checker package tests + checker self-run | `go test ./...` |
| CUI only | CUI targeted tests | `go test ./...` |
| shared compiler API | matching tests | strict Bitlang CI |
| Bitlang CI itself | targeted CI tool test/review | strict Bitlang CI on all hosted runner OSes |
| docs/routing only | link/path consistency review | no full test required unless executable instructions changed |

## Standard completion gate

For implementation changes, the canonical local entrypoint is:

```text
go run ./tools/bitlang-ci/cmd/bitlang-ci
```

The project-owned gate currently runs:

```text
gofmt cleanliness
go-rule-checker self-run
go vet ./...
go test ./...
go test -shuffle=on -count=3 ./...
Bitlang CUI build to a temporary output directory
git diff --check
```

GitHub Actions additionally runs:

```text
Linux race detector
coverage smoke test
strict Bitlang CI on Linux / Windows / macOS
```

## Failure policy

Bitlang validation is fail-closed. The following are failures, not warnings to ignore:

- a deterministic bug
- a highly likely bug
- a plausible bug
- a state likely to cause a future bug
- dangerous or unstable behavior even when not yet proven to be a bug
- a required check that did not actually execute
- a checker/configuration error that makes validation incomplete
- nondeterministic or panic-based compiler behavior

When such a condition is discovered, create an Issue even if the immediate code change is deferred.

## Rules

- Do not treat a successful build as a replacement for semantic tests.
- For shared stage representations, validate consumers of that representation.
- Prefer regression tests for every fixed bug.
- Keep compiler checks deterministic and independent from network access where possible.
- Hosted CI is an additional gate; it does not replace the project-owned `bitlang-ci` entrypoint.
- Expand project-owned checks with fuzz, property, golden, differential, round-trip, and cross-stage invariant testing as the corresponding representations become executable.
