# Bitlang

English | [日本語](README.ja.md)

Bitlang is a strict language and translation pipeline designed around explicit, inspectable lowering stages.

> **好きなコードを 好きな書き方で**  
> *Write the code you want, the way you want.*

This slogan is a design requirement, not only branding. Bitlang itself is intended to be reasonably writable and very safe, but the ecosystem is also deliberately designed so other languages and source styles can be translated into Bitlang without discarding their semantics. Preprocessing, language adapters, scoped defaults, namespace mounting, and explicit normalization exist in part to make that possible while keeping the final safety rules strict.

The initial implementation is written in **Go** so development builds can be turned into standalone executables immediately and cross-compiled without introducing a runtime dependency.

Go is only the bootstrap implementation language. The compiler core is intentionally written in a language-portable style so equivalent implementations can be created in other languages and, ultimately, Bitlang can implement Bitlang itself.

See [CODING_RULES.md](CODING_RULES.md) for the portability and self-hosting rules that apply to compiler code.

## Planned pipeline

```text
Bitlang source
  -> preprocessor
  -> Bitlang Explicit
  -> static analysis / advisor
  -> Bitlang Low
  -> Bitlang VM backend (may use TreeObject internally)
  -> Bitlang VM assembly
  -> VM or architecture translator
```

The public stage contract is Source -> Explicit -> Low -> VM Assembly, identified by the canonical artifact values `source`, `explicit`, `low`, and `vm_assembly`. The Go identifiers `ArtifactPreprocessed` and `ArtifactCompiled` remain source-compatibility aliases for `ArtifactExplicit` and `ArtifactLow` until the first major API version (`v1`), when they will be removed. The old literal values `preprocessed` and `compiled` are rejected, not silently migrated. This repository has no serialized-artifact schema version yet, so persisted values using those old literals must be migrated before entering the pipeline. `ArtifactTreeObject` is retained for VM-backend internal use and is not a public pipeline stage. Skipped or reordered transitions are rejected.

Intermediate stage boundaries are intended to become cross-implementation contracts. A future compiler written in Bitlang should be able to consume and produce equivalent artifacts without inheriting Go-specific implementation details.

## Compile-time and runtime performance

Bitlang assumes that compilation may be expensive by design.

Preprocessing, cross-language translation, full property resolution, repeated normalization passes, static analysis, ownership/borrow/lifetime/release/finalization validation, safety checks, and low-level lowering are intentionally performed before runtime. Compile speed is therefore not the highest-priority optimization target.

The design priority is:

```text
correct semantic translation
    -> safety validation
    -> high-quality lowering
    -> fast native output
    -> compile time
```

Bitlang should not remove required normalization or safety analysis merely to make compilation faster.

For native output, high-level source facilities and preprocessing work should not survive as avoidable runtime overhead. Bitlang Low is lowered to C, target assembly, or another native backend representation so the resulting application can target strong native execution performance.

The project therefore values **a fast and safety-validated finished application more than a fast compiler**.

Bitlang VM has a separate performance target because it also prioritizes portability, inspection, testing, and debugging.

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
- initial C#-style class/member parser preserving modifiers, generic syntax, and opaque body tokens
- case-insensitive identifier canonicalization
- collision-aware symbol table preserving original spelling for diagnostics
- dependency-free unit tests
- portability-first coding rules for future implementations

The initial C#-style class/member declaration shape is specified in [bitlang/BASIC_SYNTAX.md](bitlang/BASIC_SYNTAX.md). The bootstrap lexer classifies its configured keywords while retaining their source spelling; the class parser preserves modifiers and opaque bodies; expression and statement grammar remain future work. String and character literal contents are not canonicalized.

The bootstrap syntax currently accepts ASCII identifiers only; Unicode identifiers are rejected until the language contract defines them. Unicode identifier normalization/case-folding is therefore not part of this syntax version. The identifier canonicalization helper uses Go's standard-library lowercase mapping as a temporary deterministic baseline for names it receives.

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

## Security

Security-sensitive reports should follow [SECURITY.md](SECURITY.md) ([日本語](SECURITY.ja.md)). Do not publish exploit details, credentials, or sensitive reproduction data in a public Issue.

The repository runs a separate Security CI gate for CodeQL, reachable Go vulnerabilities, Go security static analysis, secret leaks, dependency changes, and GitHub Actions supply-chain risks.

## License

Bitlang and the Bitlang-related assets maintained in this repository are licensed under the [MIT License](LICENSE), unless explicitly stated otherwise.
