# Go Rule Checker

`go-rule-checker` is a small static-analysis helper for the Go bootstrap implementation of Bitlang.

It checks mechanically enforceable parts of `GO_CODING_RULES.md`. It does not define Bitlang language semantics.

## Current checks

- discourages responsibility-obscuring names such as `Manager`, `Helper`, `Utils`, and `Common`
- warns when exported declarations have no documentation comment
- warns when a function exceeds 80 lines
- reports an error when a function exceeds 120 lines
- warns when `main.go` contains functions other than `main`
- scans directories recursively while skipping `.git`, `vendor`, `.idea`, and `.vscode`

## Usage

From the repository root:

```sh
go run ./tools/go-rule-checker/cmd/go-rule-checker .
```

Build a standalone checker:

```sh
go build -o go-rule-checker ./tools/go-rule-checker/cmd/go-rule-checker
./go-rule-checker .
```

Exit status:

- `0`: no findings
- `1`: one or more rule findings
- `2`: the checker itself could not complete

## Scope

This tool intentionally checks only rules that can be determined reliably from Go syntax and file structure.

Rules such as `one file = one responsibility` and UPD Commander/ Messenger/ Processing boundaries still require architectural review. Future versions may use `FILE_RESPONSIBILITIES.md` as additional machine-readable evidence, but should prefer warnings over pretending semantic responsibility can be proven automatically.
