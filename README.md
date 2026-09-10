# Bitlang

Bitlang is a strict language and translation pipeline designed around explicit, inspectable lowering stages.

The initial implementation is written in **Go** so development builds can be turned into standalone executables immediately and cross-compiled without introducing a runtime dependency.

## Planned pipeline

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

Each stage is represented explicitly. A stage consumes one artifact kind and produces another; implicit skipping or reordering is rejected by the current core.

## Current implementation

- Go module and CLI entrypoint
- explicit artifact/stage model
- checked pipeline transitions
- case-insensitive identifier canonicalization
- collision-aware symbol table preserving original spelling for diagnostics
- dependency-free unit tests

String and character literal contents are not canonicalized. Identifier canonicalization is a symbol-level operation; source lexing will be implemented separately once the concrete grammar is fixed.

Unicode identifier normalization/case-folding is intentionally not part of the language contract yet. The current implementation uses Go's standard-library lowercase mapping as a temporary deterministic baseline.

## Build

Native executable:

```sh
go build -o bitlang ./cmd/bitlang
```

Windows x86-64 executable from another supported host:

```sh
GOOS=windows GOARCH=amd64 go build -o bitlang.exe ./cmd/bitlang
```

Windows ARM64:

```sh
GOOS=windows GOARCH=arm64 go build -o bitlang-arm64.exe ./cmd/bitlang
```

## Test

```sh
go test ./...
```

## Current CLI

```sh
./bitlang canonicalize MyVariable MYVARIABLE myvariable
```
