# Security Policy

English | [日本語](SECURITY.ja.md)

Bitlang treats compiler correctness, build integrity, and contributor supply-chain safety as security concerns.

## Supported versions

Bitlang is currently under active development. Until stable releases are published, security fixes are made against the current `main` branch. Older commits and development snapshots are not independently supported.

## Reporting a vulnerability

Please do **not** publish exploit details or sensitive reproduction steps in a public Issue.

Use GitHub's private vulnerability reporting / Security Advisories flow for this repository when it is available. If private reporting is not available, open only a minimal public Issue asking the maintainer for a private reporting channel; do not include vulnerability details, secrets, exploit code, or affected credentials in that Issue.

Useful reports include:

- the affected component and revision;
- impact and realistic attack preconditions;
- minimal reproduction information that can be shared privately;
- whether the issue can affect the compiler host, generated program, VM, build environment, CI, or repository credentials;
- suggested mitigations if known.

Do not include real credentials or unrelated private data in a report.

## Security scope

Security-sensitive areas include, but are not limited to:

- parser, lexer, preprocessor, compiler, and translator handling of hostile input;
- command injection, path traversal, unsafe temporary-file handling, and unintended file access;
- panics, memory/resource exhaustion, or denial of service caused by crafted source input;
- compiler transformations that silently produce unsafe or materially different behavior;
- VM isolation and host interaction as those components are implemented;
- dependency vulnerabilities and dependency-confusion/supply-chain risks;
- GitHub Actions workflow injection, excessive token permissions, and unpinned actions;
- secret leakage in source, Git history, logs, test fixtures, or build artifacts;
- release/build provenance and tampering risks.

A normal correctness bug is not automatically a security vulnerability, but compiler bugs that can cross a trust boundary or create exploitable output should be treated as security-sensitive.

## Repository security gates

The repository uses separate strict and security CI gates. Security validation includes:

- CodeQL using the `security-extended` query suite;
- `govulncheck` for reachable Go vulnerabilities;
- `gosec` static security analysis;
- Gitleaks secret scanning;
- dependency review for pull-request dependency changes;
- `zizmor` checks for GitHub Actions security and supply-chain hazards;
- `actionlint` workflow validation;
- SHA-pinned third-party GitHub Actions;
- Dependabot monitoring for Go modules and GitHub Actions.

Security checks are intended to fail closed. A scanner failure is not treated as a successful scan.

## Contributor security requirements

Contributions should follow these rules:

- do not add a dependency without a clear need and reviewable justification;
- do not broaden GitHub Actions permissions without explaining why;
- GitHub Actions must be pinned to an immutable commit SHA, with the human-readable version kept in a comment;
- checkout steps should not persist credentials unless a job explicitly requires Git writes;
- never commit credentials, tokens, private keys, production secrets, or real personal data;
- security-sensitive bug fixes should include regression tests where practical;
- avoid `pull_request_target` for code execution from untrusted pull requests;
- do not weaken or bypass a security check solely to make CI green. Fix the cause or document and review a narrowly scoped exception.

## Disclosure and remediation

Security reports should be validated privately first. Fixes should minimize unnecessary disclosure before users can update. Public disclosure can follow after a fix or mitigation is available, with credit to the reporter when requested and appropriate.

This project does not currently promise a fixed response-time SLA, but credible security reports are treated as high priority.
