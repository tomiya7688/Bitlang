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

## 11. Go bootstrap conventions

For the current Go implementation specifically:

- keep compiler core below `internal/bitlang` independent from `cmd/`
- prefer ordinary structs over generic semantic containers
- use `any` only at temporary representation boundaries; replace it with concrete Bitlang IR types as those types become defined
- do not use reflection in compiler semantics
- do not use goroutines/channels inside compiler passes without a deterministic outer abstraction
- return errors instead of panicking for malformed user input
- keep OS-specific operations out of semantic packages
- run `gofmt` on Go source
- keep `go test ./...` and `go build ./cmd/bitlang` working

These restrictions may be revised when Bitlang itself has enough defined semantics to become the reference implementation.
