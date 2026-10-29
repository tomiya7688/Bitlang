# Bitlang Coding Rules

These rules apply to the bootstrap implementation and to future implementations in other languages.

## 1. Portability is a design requirement

The Go implementation is a bootstrap implementation, not the definition of Bitlang.

Core language logic MUST be written so that it can be reimplemented mechanically in another conventional systems language and, eventually, in Bitlang itself.

When two implementations are possible, prefer the one whose control flow, data model, and error behavior can be expressed similarly in C, Go, Rust, Zig, Java, C#, or Bitlang.

## 2. One file, one responsibility

Every hand-maintained source file MUST have exactly one primary responsibility.

A file MUST NOT become a container for several merely related responsibilities. Package/module membership alone is not sufficient reason to place code in the same file.

The responsibility of a file should be describable in one short sentence. If that sentence requires "and" to join independent behaviors, the file should normally be split.

Examples:

```text
Good:
  token.go              -> defines token representation
  lexer.go              -> converts source text into tokens
  symbol.go             -> represents symbols
  symbol_table.go       -> stores and resolves symbols
  artifact.go           -> represents pipeline artifacts
  stage.go              -> represents and executes one pipeline stage
  pipeline.go           -> orders and executes stages

Avoid:
  lexer.go              -> token definitions + scanning + diagnostics + source loading
  compiler.go           -> parsing + analysis + lowering + output writing
```

This rule is stronger than a numeric line limit. A 90-line file with two independent responsibilities MUST be split, while a longer declarative table may be acceptable if it still represents exactly one responsibility.

The intended long-term mapping is:

```text
one conceptual responsibility
        ~=
one source file in the bootstrap implementation
        ~=
one class/object responsibility in an object-oriented implementation
        ~=
one corresponding component in a future Bitlang implementation
```

When the project is rewritten in an object-oriented language, a source file should therefore be structurally close to one class or one class-sized responsibility. This does not require forcing every language into class syntax; it requires keeping responsibilities narrow enough that such a translation is natural.

A file MAY contain small supporting declarations that exist only to serve its single responsibility, but reusable or independently meaningful concepts MUST be moved to their own file.

## 3. One function, one operation

Every function or method MUST perform one clearly identifiable operation at one abstraction level.

The operation should be describable with one verb phrase, for example:

```text
canonicalize an identifier
resolve a symbol
scan one token
parse one expression
validate one transition
lower one statement
format one diagnostic
```

A function that performs several sequential responsibilities MUST be decomposed, even when the total line count is small.

Avoid functions of the form:

```text
read source -> tokenize -> parse -> analyze -> lower -> write output
```

Prefer orchestration that calls narrowly focused operations:

```text
source = loadSource(...)
tokens = lexSource(source)
ast = parseTokens(tokens)
result = analyzeAst(ast)
```

An orchestration function itself has one operation: coordinating a defined workflow. It MUST NOT also contain the detailed implementation of the operations it coordinates.

A helper function MUST represent a real operation or invariant. Do not create meaningless wrappers solely to satisfy this rule.

## 4. The specification outranks the host language

Bitlang semantics MUST NOT accidentally depend on Go behavior.

Host-language behavior such as integer overflow, map iteration order, Unicode handling, filesystem path rules, goroutine scheduling, panic/recover, interface dispatch, or implicit zero values MUST NOT become Bitlang behavior unless the Bitlang specification explicitly says so.

If host-specific behavior is required, isolate it behind a small adapter boundary and document it.

## 5. Prefer explicit data and explicit control flow

Core code SHOULD use:

- plain named structs / records
- explicit enums or tagged values
- explicit loops and branches
- explicit input and output values
- explicit error returns
- small deterministic functions

Core code SHOULD avoid unless there is a strong reason:

- reflection
- code generation required to understand semantics
- metaprogramming
- host-language generics in semantic data structures
- hidden global state
- implicit initialization with semantic meaning
- exceptions/panics for normal control flow
- concurrency inside deterministic compiler passes
- clever iterator/functional chains that obscure execution order

The goal is not to write Go as if it were C. The goal is to keep the semantic algorithm obvious enough that another implementation can follow the same steps.

## 6. Keep stages as real boundaries

The canonical lowering pipeline is:

```text
Bitlang source
  -> preprocessor
  -> Bitlang preprocessed
  -> static analysis / advisor
  -> Bitlang compiled
  -> tree object
  -> Bitlang VM assembly
```

A stage MUST have an explicit input representation and output representation.

A later optimization MUST NOT erase these conceptual boundaries. Internal fast paths are allowed only when behavior remains equivalent to executing the defined stages.

Whenever practical, intermediate representations SHOULD be serializable and independently testable. This allows implementations written in different languages to compare results at stage boundaries.

## 7. Determinism first

Given the same source, options, and declared environment, compiler stages SHOULD produce the same semantic result regardless of host OS, CPU, locale, process timing, or map/hash ordering.

Do not rely on unordered container traversal when output order is observable. Sort or otherwise define the order explicitly.

## 8. Keep the core independent from the CLI and OS

Language semantics, preprocessing, parsing, analysis, lowering, VM assembly generation, and translation logic belong in reusable core packages/modules.

Command-line parsing, terminal output, filesystem access, environment variables, and process exit codes belong at the outer boundary.

The core SHOULD be callable by a CLI, test runner, VM, editor integration, or another compiler without pretending to be a command-line program.

## 9. Minimize dependencies

The bootstrap compiler SHOULD prefer the host standard library.

A third-party dependency in semantic/compiler core code requires a clear reason. Dependencies that make a future Bitlang implementation substantially harder to reproduce SHOULD be avoided.

## 10. Tests are cross-implementation contracts

Tests SHOULD describe semantic inputs and outputs rather than Go-specific implementation details.

Important behavior SHOULD eventually have language-neutral test vectors that every implementation can run.

For stage-based behavior, prefer tests of the form:

```text
input artifact + options -> expected artifact / diagnostics
```

This will allow the Go bootstrap implementation and the future Bitlang implementation to be checked against the same corpus.

## 11. Naming and representation

Names SHOULD describe Bitlang concepts, not Go mechanisms. For example, prefer `Artifact`, `Stage`, `Symbol`, and `Diagnostic` over names tied to interfaces, goroutines, readers, or other host abstractions.

Bitlang identifiers are case-insensitive, but their original spelling SHOULD be retained for diagnostics.

String and character contents MUST NOT be modified by identifier canonicalization.

## 12. Self-hosting is the long-term target

The project SHOULD be designed toward this bootstrap chain:

```text
Stage 0: Go bootstrap implementation
        |
        v
Stage 1: enough Bitlang exists to implement parts of Bitlang in Bitlang
        |
        v
Stage 2: Bitlang compiler written primarily in Bitlang
        |
        v
Stage 3: Bitlang builds Bitlang and reproduces equivalent stage outputs
```

This does not require prematurely implementing every feature in Bitlang. It does require avoiding architecture that can only exist comfortably in Go.

A new core feature SHOULD therefore be reviewed with one additional question:

> Could we implement this same algorithm, data model, and stage boundary in Bitlang without redesigning the compiler?

If the answer is no, either simplify the design or explicitly document why the host-specific mechanism is temporary.

## 13. Comments are part of the implementation

Comments are required where they preserve design intent, semantic constraints, portability assumptions, or non-obvious behavior.

The goal is not maximum comment density. The goal is that a future contributor, another language implementation, or a future Bitlang self-hosted compiler author can understand why the code exists and what behavior must be preserved.

The following SHOULD have comments:

- exported/public types, functions, constants, and important fields
- compiler stages and stage boundaries
- intermediate representations and their invariants
- non-obvious algorithms
- portability-sensitive behavior
- deterministic ordering requirements
- temporary bootstrap limitations
- behavior that deliberately differs from the host language
- error cases whose reason is not immediately obvious
- code where an apparently simpler implementation would be semantically wrong

Comments SHOULD explain one or more of:

- why this code exists
- what invariant must remain true
- what input/output contract is being enforced
- why a particular implementation was chosen
- what must remain compatible across implementations
- what is intentionally temporary during bootstrap

Comments SHOULD NOT merely translate syntax into prose.

Bad:

```text
increment i by one
```

Good:

```text
Advance exactly one source byte here. UTF-8 decoding happens in the lexer so
this scanner must not reinterpret the byte sequence.
```

A TODO comment MUST state enough context to be actionable. Prefer:

```text
TODO: Replace this host lowercase mapping once Bitlang identifier Unicode
normalization is specified. All implementations must then use the same table.
```

over:

```text
TODO: fix later
```

Large logical sections MAY use short header comments, but comments MUST NOT be used as a substitute for splitting an oversized file or function.

When behavior changes, nearby comments MUST be updated in the same change. Stale comments are considered bugs.

## 14. File size is a secondary safety limit

The primary rule is one file, one responsibility. Line count exists only as an additional warning against accidental growth.

For ordinary hand-maintained core source files:

- 0-200 lines: preferred range
- 200-300 lines: acceptable only while the file clearly retains one responsibility
- 300-400 lines: MUST trigger a split review
- more than 400 lines: MUST normally be split
- more than 600 lines: prohibited except for generated data, declarative tables, or another explicitly documented exceptional case

Blank lines and comments count toward these limits because documentation also affects navigability.

Generated files MAY exceed them but MUST be clearly marked as generated and SHOULD NOT contain hand-maintained semantic logic.

Never keep multiple responsibilities together merely because the file is below the numeric limit.

## 15. Function size is a secondary safety limit

The primary rule is one function, one operation. Line count exists only as an additional warning against mixing operations or abstraction levels.

As a guideline:

- under 40 lines is preferred for ordinary functions
- 40-80 lines MUST be reviewed for hidden multiple operations
- over 80 lines SHOULD normally be decomposed
- over 120 lines requires explicit justification in a comment or review

Parser/lexer state machines may occasionally justify longer functions, but a long state machine still MUST perform one operation. Independent state transitions or semantic actions should be extracted where doing so preserves readable control flow.

## 16. File and function growth must be considered during every change

Before adding code to an existing file, ask:

1. Is this exactly the same responsibility as the file already has?
2. Would this responsibility map naturally to the same class in an object-oriented implementation?
3. Would a future Bitlang implementation naturally keep these operations together?
4. Is an independently meaningful concept being introduced that deserves its own file?

Before adding logic to an existing function, ask:

1. Is this exactly the same operation the function already performs?
2. Is it at the same abstraction level?
3. Can the function still be named accurately with one verb phrase?
4. Is the function beginning to coordinate and implement details at the same time?

If these checks fail, split before adding the new implementation.

## 17. Go bootstrap conventions

For the current Go implementation specifically:

- keep compiler core below `internal/bitlang` independent from `cmd/`
- treat each Go source file as one class-sized responsibility even though Go has no classes
- do not group several types in one file merely because Go permits it
- prefer one primary semantic type/component per file when that type has independent behavior
- place independently reusable support types in their own files
- keep methods associated with the responsibility represented by that file
- prefer ordinary structs over generic semantic containers
- use `any` only at temporary representation boundaries; replace it with concrete Bitlang IR types as those types become defined
- do not use reflection in compiler semantics
- do not use goroutines/channels inside compiler passes without a deterministic outer abstraction
- return errors instead of panicking for malformed user input
- keep OS-specific operations out of semantic packages
- add Go documentation comments to exported identifiers
- comment semantic invariants and portability-sensitive decisions even for unexported code when they are non-obvious
- do not use comments to excuse oversized files; split by responsibility instead
- run `gofmt` on Go source
- keep `go test ./...` and `go build ./cmd/bitlang` working

These restrictions may be revised when Bitlang itself has enough defined semantics to become the reference implementation.

## 継承の共通規約（必須・レビュー基準）

### 禁止：具象クラス起点の連続継承

- **具象クラス A を継承するクラス B を、クラス C がさらに継承してはならない。** すなわち `具象 A → B → C` を禁止する。中間の B が具象か抽象かは問わない。
- **`抽象 A → 具象 B → C` は、この禁止規則の違反ではない。** ただし設計・可読性のレビューは必要とする。
- この規則単独では `具象 A → B` の一段の継承を一律には禁止しない。ただし、既存の「具象クラスの継承そのものを禁止」など、より厳しいプロジェクト規則を常に優先する。
- 言語に暗黙で存在する `Object` / `System.Object` / `java.lang.Object` 等の最上位基底型、インターフェースの継承・実装はこの違反判定の対象外とする。

### 継承を選ぶ前の設計原則

- **具体的な処理を持つクラスの継承は原則として避ける。** 抽象クラスであっても具体的な実装を持つなら同じ。継承関係そのものに意味がなく、処理の再利用だけが目的なら委譲・合成を優先する。
- 基底型と派生型に共通する意味のある形・契約が実際に存在する、または同じ形の型が今後必要になると合理的に予想できる場合に限って継承を検討する。抽象化・継承階層を「いつか使うかもしれない」というだけで作らない。
- 形式上許可される継承でも、階層が読みづらい、責務の理解が難しい、変更の波及が過大、置換可能性が成立しない等の問題があれば**レビューで不合格**とする。「禁止規則に抵触しない」は設計の承認を意味しない。

### 適用範囲と例外

- 当該プロジェクトが意図的にその継承構造を採用し、設計正本に理由・適用範囲を明記している箇所は例外とする。単に既存コードに同じ形があることを例外の根拠にしない。
- コンパイラ・解析器が入力として扱う言語の継承仕様、サンプル、検証用fixtureと、プロジェクト自身の実装規約を区別する。
- 新規・変更するプロジェクト所有コードへ適用し、無関係な既存コードの一斉書き換えを要求しない。
