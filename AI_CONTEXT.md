# Bitlang AI Context

AI-assisted development entrypoint. Keep this file small.

## Source of Truth
- README: `README.md`
- Responsibilities: `FILE_RESPONSIBILITIES.md`
- Shared implementation rules: `CODING_RULES.md`
- Go bootstrap rules: `GO_CODING_RULES.md`
- Current state: `CURRENT_STATE.md`
- Change routing: `CHANGE_ROUTING.md`
- Validation routing: `VALIDATION_ROUTING.md`
- Tasks: GitHub Issues

## Read First
1. Current task or Issue
2. Relevant route in `CHANGE_ROUTING.md`
3. Target source and matching tests
4. Detailed docs only if needed

Stop exploring when Goal, Required, Acceptance, and the working set are clear.

## Constraints
- Go behavior does not define Bitlang semantics.
- Go-side architecture rules are implementation rules only.
- CUI and GUI use the same official compiler backend.
- Preserve explicit compiler stage boundaries.
- Prefer deterministic behavior.

## Ignore Normally
- unrelated Issues and history
- generated/build outputs
- large logs
- unrelated compiler stages

## Policy Checks
- Responsibility map: `FILE_RESPONSIBILITIES.md`
- Go checker: `tools/go-rule-checker/`

## Validation
Use the smallest sufficient checks from `VALIDATION_ROUTING.md`.
Report unverified areas explicitly.
