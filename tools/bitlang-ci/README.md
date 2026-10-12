# Bitlang CI

English | [日本語](README.ja.md)

`bitlang-ci` is the project-owned strict validation gate for the Bitlang compiler implementation.

It is intentionally stricter than a normal application CI because a silent compiler or transformation defect can change program meaning while still producing apparently valid output.

## Usage

```sh
go run ./tools/bitlang-ci/cmd/bitlang-ci
```

## Current gate

The command runs all checks and fails if any check fails:

- `gofmt` cleanliness
- strict machine-readable language specifications and family pipeline contract manifest
- Bitlang `go-rule-checker` self-check
- bilingual documentation pair existence/link check
- `go mod tidy -diff`
- `go vet ./...`
- `go test ./...`
- shuffled/repeated tests
- `go build ./...`
- Bitlang CUI build into a temporary output directory
- `git diff --check`

GitHub Actions additionally runs:

- the full gate on Linux, Windows, and macOS
- the race detector
- coverage smoke tests
- PR-only documentation pair change synchronization
- cross-build checks for Linux, Windows, and macOS on amd64 and arm64

## Failure policy

Bitlang CI is fail-closed. A missing check target, checker execution failure, malformed checker configuration, panic, suspicious nondeterminism, stale registered translation pair, or other unverified state must not be treated as success.

The family pipeline manifest is checked against the Bitlang artifact identities and transition graph. Until cross-repository acceptance fixtures are implemented, their manifest status remains `not_implemented` and the strict gate fails instead of silently skipping those boundaries. Even after fixture files are added, the gate remains failed until an execution runner is wired in; sibling-repository CI inspection is also not implemented yet.

Future project-owned checks should include semantic invariants, fuzzing, property tests, golden conversion tests, differential tests, round-trip tests where meaningful, and cross-stage consistency checks.
