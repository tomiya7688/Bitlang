# Contributing to Bitlang

Thank you for your interest in Bitlang.

Bitlang is still under active design and implementation. Contributions do not need to be large, and you do not need to understand the whole compiler pipeline before participating.

## Good ways to contribute

Contributions are welcome in areas such as:

- Go bootstrap compiler implementation
- lexer, preprocessing, static analysis, and pipeline work
- tests and regression cases
- compiler diagnostics and error handling
- developer tooling
- documentation and examples
- specification review and inconsistency reports
- VM / assembly / translator design and implementation as those stages mature

Non-code contributions are welcome too. Finding an unclear rule, contradictory specification, missing test case, or confusing document is useful work.

## Start here

For a first contribution, check issues labeled [`good first issue`](https://github.com/tomiya7688/Bitlang/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22).

If an issue is unclear, ask on the issue before making a large design assumption.

Bitlang deliberately separates stable language semantics from implementation details. Do not invent new language syntax or semantics merely to complete an implementation task. If a required behavior is unspecified, surface the ambiguity instead.

## Repository orientation

The intended lowering pipeline is:

```text
Bitlang source
  -> preprocessor
  -> Bitlang preprocessed
  -> static analysis / advisor
  -> Bitlang compiled
  -> tree object
  -> Bitlang VM assembly
  -> VM or architecture translator
```

Useful entry points:

- `README.md` — project overview and basic build instructions
- `CURRENT_STATE.md` — compact implementation status
- `bitlang/LANGUAGE_SPEC.md` — language specification
- `CODING_RULES.md` — portability and compiler implementation rules
- `GO_CODING_RULES.md` — Go bootstrap implementation rules
- `FILE_RESPONSIBILITIES.md` — source-file responsibility index

## Development setup

The current bootstrap implementation uses Go.

Run the test suite with:

```sh
go test ./...
```

Build the current CLI with:

```sh
go build -o bitlang ./cmd/bitlang
```

Before submitting a change, make sure relevant tests pass and keep changes focused on the issue being solved.

## Coding expectations

For Go implementation work:

- follow `GO_CODING_RULES.md`
- preserve the repository's responsibility boundaries
- prefer deterministic errors over panic or silent failure
- add regression tests for bug fixes
- avoid adding dependencies unless they are clearly justified
- keep compiler behavior explicit and inspectable

The Go implementation rules are not Bitlang language rules. Do not accidentally turn bootstrap implementation conventions into user-facing language requirements.

## Specifications and design changes

Small implementation fixes can normally be submitted directly against an existing issue.

Changes that alter Bitlang semantics, intermediate representations, or pipeline contracts should be discussed first. The specification is the source of truth for language behavior; an implementation should not silently redefine it.

When you find a specification ambiguity, opening an issue that explains the conflicting interpretations is a valid contribution by itself.

## Pull requests

A useful pull request should generally:

1. solve one clearly scoped problem;
2. reference the relevant issue when one exists;
3. include or update tests when behavior changes;
4. update documentation when the public contract changes;
5. pass the repository's CI checks.

Large unrelated refactors should be separated from functional changes where practical.

## License

By contributing to this repository, you agree that your contribution is provided under the repository's MIT License.
