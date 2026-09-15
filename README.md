# Bitlang

English | [日本語](README.ja.md)

Bitlang is a strict language and translation pipeline designed around explicit, inspectable lowering stages.

The initial implementation is written in **Go** so development builds can be turned into standalone executables immediately and cross-compiled without introducing a runtime dependency.

Go is only the bootstrap implementation language. The compiler core is intentionally written in a language-portable style so equivalent implementations can be created in other languages and, ultimately, Bitlang can implement Bitlang itself.

See [CODING_RULES.md](CODING_RULES.md) for the portability and self-hosting rules that apply to compiler code.

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

Intermediate stage boundaries are intended to become cross-implementation contracts. A future compiler written in Bitlang should be able to consume and produce equivalent artifacts without inheriting Go-specific implementation details.

## Bootstrap and self-hosting target

The intended evolution is:

```text
Go bootstrap compiler
  -> implement stable Bitlang semantics and stage formats
  -> implement compiler components in Bitlang
  -> build the Bitlang compiler using the bootstrap compiler
  -> use Bitlang-built Bitlang to build Bitlang again
  -> compare stage outputs / behavior for equivalence
```

Self-hosting is a long-term architectural target, not a requirement to prematurely rewrite unstable compiler components.

## Current implementation

- Go module and CLI entrypoint
- explicit artifact/stage model
- checked pipeline transitions
- case-insensitive identifier canonicalization
- collision-aware symbol table preserving original spelling for diagnostics
- dependency-free unit tests
- portability-first coding rules for future implementations

String and character literal contents are not canonicalized. Identifier canonicalization is a symbol-level operation; source lexing will be implemented separately once the concrete grammar is fixed.

Unicode identifier normalization/case-folding is intentionally not part of the language contract yet. The current implementation uses Go's standard-library lowercase mapping as a temporary deterministic baseline.

## Contributing

Contributors are welcome. You do not need to understand the entire compiler to participate: implementation fixes, tests, documentation, tooling, specification review, and bug reports are all useful.

Start with [CONTRIBUTING.md](CONTRIBUTING.md) ([日本語](CONTRIBUTING.ja.md)), or browse issues labeled [good first issue](https://github.com/tomiya7688/Bitlang/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22).

Contributor-facing documentation is maintained in English and Japanese. Registered language pairs are checked by CI so that a pull request does not update only one side.

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

For the full project validation gate:

```sh
go run ./tools/bitlang-ci/cmd/bitlang-ci
```

## Current CLI

```sh
./bitlang canonicalize MyVariable MYVARIABLE myvariable
```

## License

Bitlang and the Bitlang-related assets maintained in this repository are licensed under the [MIT License](LICENSE), unless explicitly stated otherwise.
