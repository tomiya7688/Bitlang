# Bitlang Language Specification

This document defines the Bitlang language specification.

## Role of Bitlang in the Bitlang family

Bitlang is the common semantic target of the Bitlang language family.

Languages in the Bitlang family may use different syntax, convenience features, or domain-specific notation, but they are expected to transform into ordinary Bitlang before the Bitlang preprocessing stage.

Conceptually:

```text
Bitlang family language
    -> transform
    -> Bitlang
    -> preprocess
    -> Bitlang preprocessed
    -> compile
    -> Bitlang compiled
```

This allows family languages to focus on their own syntax and usability while Bitlang remains the shared language for expressing the normalized meaning of the program.

## Preprocessor functions

Bitlang provides preprocessor functions as a language-specific mechanism for reducing repetitive or inconvenient source code without weakening the strict Bitlang language model.

A preprocessor function is executed before ordinary compilation and does not remain as a runtime function.

Its role combines ideas similar to:

- Go-style code generation (`go generate`)
- C-style definitions/macros (`#define`)
- source-text transformation tools

Preprocessor functions may be used for three broad purposes:

1. **Define** - declarations, aliases, constants, abbreviations, or reusable definitions.
2. **Generate** - generate Bitlang declarations, functions, types, or other source structures.
3. **Transform** - transform existing source or another supported representation into Bitlang.

The output of preprocessing must be valid **Bitlang preprocessed** input.

Preprocessor functions are intended both to make strict Bitlang somewhat easier to write and to remove transformation/code-generation bottlenecks. They are also part of the infrastructure that can be used by Bitlang-family language transformers.

Detailed syntax and semantic rules will be added as the language design is finalized.
