# Bitlang CI

`bitlang-ci` is the project-owned strict validation gate for the Bitlang compiler implementation.

It is intentionally stricter than a normal application CI because a silent compiler or transformation defect can change program meaning while still producing apparently valid output.

## Usage

```sh
go run ./tools/bitlang-ci/cmd/bitlang-ci
```

## Current gate

The command runs all checks and fails if any check fails:

- `gofmt` cleanliness
- Bitlang `go-rule-checker` self-check
- `go vet ./...`
- `go test ./...`
- shuffled/repeated tests
- Bitlang CUI build into a temporary output directory
- `git diff --check`

GitHub Actions additionally runs the race detector and coverage smoke tests.

## Failure policy

Bitlang CI is fail-closed. A missing check target, checker execution failure, malformed checker configuration, panic, suspicious nondeterminism, or other unverified state must not be treated as success.

Future project-owned checks should include semantic invariants, fuzzing, property tests, golden conversion tests, differential tests, round-trip tests where meaningful, and cross-stage consistency checks.
