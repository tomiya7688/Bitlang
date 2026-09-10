# Bitlang Coding Rules

These rules apply to the bootstrap implementation and to future implementations in other languages.

## 1. Portability is a design requirement

The Go implementation is a bootstrap implementation, not the definition of Bitlang.

Core language logic MUST be written so that it can be reimplemented mechanically in another conventional systems language and, eventually, in Bitlang itself.

When two implementations are possible, prefer the one whose control flow, data model, and error behavior can be expressed similarly in C, Go, Rust, Zig, Java, C#, or Bitlang.

## 2. The specification outranks the host language

Bitlang semantics MUST NOT accidentally depend on Go behavior.

Host-language behavior such as integer overflow, map iteration order, Unicode handling, filesystem path rules, goroutine scheduling, panic/recover, interface dispatch, or implicit zero values MUST NOT become Bitlang behavior unless the Bitlang specification explicitly says so.

If host-specific behavior is required, isolate it behind a small adapter boundary and document it.

## 3. Prefer explicit data and explicit control flow

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

## 4. Keep stages as real boundaries

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

## 5. Determinism first

Given the same source, options, and declared environment, compiler stages SHOULD produce the same semantic result regardless of host OS, CPU, locale, process timing, or map/hash ordering.

Do not rely on unordered container traversal when output order is observable. Sort or otherwise define the order explicitly.

## 6. Keep the core independent from the CLI and OS

Language semantics, preprocessing, parsing, analysis, lowering, VM assembly generation, and translation logic belong in reusable core packages/modules.

Command-line parsing, terminal output, filesystem access, environment variables, and process exit codes belong at the outer boundary.

The core SHOULD be callable by a CLI, test runner, VM, editor integration, or another compiler without pretending to be a command-line program.

## 7. Minimize dependencies

The bootstrap compiler SHOULD prefer the host standard library.

A third-party dependency in semantic/compiler core code requires a clear reason. Dependencies that make a future Bitlang implementation substantially harder to reproduce SHOULD be avoided.

## 8. Tests are cross-implementation contracts

Tests SHOULD describe semantic inputs and outputs rather than Go-specific implementation details.

Important behavior SHOULD eventually have language-neutral test vectors that every implementation can run.

For stage-based behavior, prefer tests of the form:

```text
input artifact + options -> expected artifact / diagnostics
```

This will allow the Go bootstrap implementation and the future Bitlang implementation to be checked against the same corpus.

## 9. Naming and representation

Names SHOULD describe Bitlang concepts, not Go mechanisms. For example, prefer `Artifact`, `Stage`, `Symbol`, and `Diagnostic` over names tied to interfaces, goroutines, readers, or other host abstractions.

Bitlang identifiers are case-insensitive, but their original spelling SHOULD be retained for diagnostics.

String and character contents MUST NOT be modified by identifier canonicalization.

## 10. Self-hosting is the long-term target

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

## 11. Comments are part of the implementation

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

## 12. Keep files small and single-purpose

A source file SHOULD represent one clear responsibility or one tightly related group of data structures and operations.

Do not grow a file merely because new code belongs to the same package/module. Package membership is not sufficient justification for sharing a file.

For ordinary core source files, use these size guidelines:

- 0-200 lines: preferred range
- 200-300 lines: acceptable when the file still has one clear responsibility
- 300-400 lines: SHOULD trigger a split review
- more than 400 lines: MUST normally be split
- more than 600 lines: prohibited except for generated data, declarative tables, or another explicitly documented exceptional case

Blank lines and comments count toward these limits because documentation also affects navigability. Generated files MAY exceed them but MUST be clearly marked as generated and SHOULD NOT contain hand-maintained semantic logic.

Line count is only a warning signal. A file MUST be split earlier when it contains multiple independent responsibilities.

Examples of responsibilities that SHOULD normally be separate files include:

- token definitions vs lexer implementation
- AST node definitions vs parser control flow
- diagnostics data vs diagnostic rendering
- symbol representation vs scope resolution
- IR definitions vs lowering logic
- VM instruction definitions vs VM execution
- architecture mapping tables vs translation algorithms
- CLI argument parsing vs compiler invocation

Prefer names that describe responsibility directly, for example:

```text
lexer_token.go
lexer_scan.go
parser_expr.go
parser_stmt.go
diagnostic.go
diagnostic_format.go
ir_types.go
lower_expr.go
lower_stmt.go
```

Do not create meaningless fragments such as `utils1.go`, `helpers2.go`, or files split only to satisfy a numeric line limit.

A split SHOULD follow semantic boundaries that can also be reproduced in another implementation language.

## 13. Keep functions small enough to understand locally

A function SHOULD normally perform one operation at one abstraction level.

As a guideline:

- under 40 lines is preferred for ordinary functions
- 40-80 lines SHOULD be reviewed for extraction opportunities
- over 80 lines SHOULD normally be split
- over 120 lines requires explicit justification in a comment or review

Complex parser/lexer state machines may occasionally justify longer functions, but even there state transitions SHOULD be separated when doing so improves clarity without hiding control flow.

Do not split functions into tiny wrappers merely to meet line limits. The purpose is local comprehensibility and portability, not numeric compliance.

## 14. File growth must be considered during every change

Before adding a substantial feature to an existing file, check:

1. Is this the same responsibility as the existing file?
2. Will the file remain easy to understand in isolation?
3. Would a future C/Rust/Bitlang implementation naturally put this logic in the same module/file?
4. Is the file approaching the split-review threshold?

If the answer suggests a new responsibility, create a new file before adding the implementation.

Refactoring an oversized file is part of feature work when the feature would otherwise make the problem worse. Do not defer all structural cleanup indefinitely.

## 15. Go bootstrap conventions

For the current Go implementation specifically:

- keep compiler core below `internal/bitlang` independent from `cmd/`
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
