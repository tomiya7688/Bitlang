# Bitlang Compiler Specification

This document defines the Bitlang compiler behavior, transformation stages, interfaces, and implementation requirements.

## Front-end transformation model

Bitlang-family languages are front-end languages that transform into ordinary Bitlang before Bitlang preprocessing begins.

The compiler toolchain should therefore treat transformation and preprocessing as separate stages:

```text
Bitlang family source
    -> family-language transformer
    -> Bitlang
    -> Bitlang preprocessor
    -> Bitlang preprocessed
    -> compiler
    -> Bitlang compiled
```

A transformer is responsible for translating a family language's syntax and language-specific conveniences into equivalent Bitlang meaning. It should not bypass the Bitlang layer unless a future specification explicitly permits it.

## Preprocessor functions

The Bitlang preprocessing system supports preprocessor functions. These are compile-time functions used for definitions, code generation, and source transformation.

They may provide behavior comparable to:

- Go-style pre-build generation
- C-style macro/definition expansion
- general source-text or structure transformation

Preprocessor functions execute before the normal Bitlang-to-Bitlang-preprocessed compilation stage is completed and do not become runtime functions.

The preprocessing implementation should preserve enough source mapping information to report diagnostics against the original source where practical and should allow developers to inspect the generated Bitlang preprocessed output.

Detailed compiler rules will be added as the compiler design is finalized.
