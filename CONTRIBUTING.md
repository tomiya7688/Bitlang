# Contributing to Bitlang

English | [日本語](CONTRIBUTING.ja.md)

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

- `README.md` / `README.ja.md` — project overview and basic build instructions
- `CURRENT_STATE.md` / `CURRENT_STATE.ja.md` — compact implementation status
- `SECURITY.md` / `SECURITY.ja.md` — vulnerability reporting and security requirements
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

Run the full local CI gate with:

```sh
go run ./tools/bitlang-ci/cmd/bitlang-ci
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

## Security expectations

Read [SECURITY.md](SECURITY.md) before making security-sensitive changes or reporting a vulnerability.

- never commit credentials, tokens, private keys, production secrets, or real personal data;
- new dependencies require a clear reason and should be kept to the minimum necessary;
- GitHub Actions must use immutable commit SHAs, with the readable release version kept in a comment;
- do not broaden workflow token permissions without an explicit reason;
- checkout credentials should not be persisted unless a job intentionally performs authenticated Git writes;
- do not use `pull_request_target` to execute code from an untrusted pull request;
- security-sensitive fixes should add regression coverage where practical;
- do not disable, weaken, or bypass security checks merely to make CI pass.

Potential vulnerabilities must not be disclosed with exploit details in a public Issue. Follow the private reporting process described in `SECURITY.md`.

## Documentation language policy

Contributor-facing documentation is maintained in English and Japanese.

- the base `.md` file is the English document;
- the `.ja.md` file is the Japanese document;
- when a registered document pair is changed, both files must be updated in the same pull request;
- CI verifies that registered pairs exist, link to each other, and are changed together;
- a translation must preserve meaning rather than introduce a separate specification.

The current CI-enforced pairs are listed in `tools/doc-pair-checker`. Technical specification translations will be added to that list as they are reviewed; an unreviewed translation must not silently become a second source of truth.

## Specifications and design changes

Small implementation fixes can normally be submitted directly against an existing issue.

Changes that alter Bitlang semantics, intermediate representations, or pipeline contracts should be discussed first. The specification is the source of truth for language behavior; an implementation should not silently redefine it.

When you find a specification ambiguity, opening an issue that explains the conflicting interpretations is a valid contribution by itself.

## Pull requests

A useful pull request should generally:

1. solve one clearly scoped problem;
2. reference the relevant issue when one exists;
3. include or update tests when behavior changes;
4. update both registered documentation languages when the public contract changes;
5. pass the repository's strict and security CI checks.

Large unrelated refactors should be separated from functional changes where practical.

## License

By contributing to this repository, you agree that your contribution is provided under the repository's MIT License.
