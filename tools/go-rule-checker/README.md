# Go Rule Checker

`go-rule-checker` is a small static-analysis helper for the Go bootstrap implementation of Bitlang.

It checks mechanically enforceable parts of `GO_CODING_RULES.md`. It does not define Bitlang language semantics.

## Current checks

- `NAME`: vague names such as `Manager`, `Helper`, `Utils`, and `Common`
- `DOC`: exported declarations without documentation comments
- `SIZE`: functions over 80/120 lines
- `FILESIZE`: files over 400/600 lines
- `MAIN`: functions other than `main` inside `main.go`

## Usage

```sh
go run ./tools/go-rule-checker/cmd/go-rule-checker .
```

Output is intentionally compact:

```text
W internal/foo.go:18 NAME vague word: Manager
E internal/bar.go:10 SIZE >120: Compile
NG go-rules: 2
```

No findings:

```text
OK go-rules
```

## Ignore

If `.go-rule-checker-ignore` exists in the current directory, it is loaded automatically.

```text
path generated/**
rule DOC third_party/**
rule NAME legacy.go
```

- `path <pattern>` ignores every rule for a path.
- `rule <RULE> <pattern>` ignores one rule for a path.
- `#` starts a comment.

See `IGNORE_FORMAT.md` and `.go-rule-checker-ignore.example`.

## Exit status

- `0`: no findings
- `1`: one or more findings
- `2`: checker execution failure

## Scope

Rules such as `one file = one responsibility` and UPD Commander / Messenger / Processing boundaries still require architectural review. Future checks should prefer warnings when responsibility cannot be proven mechanically.
