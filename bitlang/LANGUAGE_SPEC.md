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
    -> preprocess + normalize
    -> Bitlang preprocessed
    -> compile
    -> Bitlang compiled
```

This allows family languages to focus on their own syntax and usability while Bitlang remains the shared language for expressing the meaning of the program.

## Modules and naming hierarchy

Bitlang does not provide a separate `namespace` construct.

`module` is the single top-level organizational concept used for namespacing and may also carry broader responsibilities such as import/export boundaries, dependency management, visibility, compilation grouping, initialization boundaries, or future versioning rules.

A module may contain nested modules or named declarations so that fully qualified names can be expressed hierarchically.

Conceptually:

```text
Game.Combat.Player.attack
```

may identify `attack` as a member reached through the `Game` and `Combat` module hierarchy and the `Player` type.

Bitlang preprocessed should use fully qualified references wherever practical so that the referenced declaration is explicit and unambiguous. Source-level aliases or shortened forms may be provided through preprocessing rules, but they must normalize to the canonical fully qualified representation.

## Functional-language support

Bitlang must be able to represent functional-programming semantics even though Bitlang itself does not need to be primarily written as a functional language.

The Bitlang family may provide a dedicated **Bit Function lang** front end. Bit Function lang is intended to be written naturally in a functional style and then transform into ordinary Bitlang.

Bitlang should therefore be able to represent concepts such as:

- first-class functions
- functions stored in variables
- functions passed as arguments
- functions returned as values
- lambdas
- closures
- higher-order functions
- immutable values
- recursion
- pattern matching
- sum/variant-like data representations
- Option/Result-like values
- partial application and currying semantics where required by a family language

These features may have multiple convenient source-level notations in Bitlang or in Bitlang-family languages, but they must be normalized before reaching Bitlang preprocessed.

The intended functional pipeline is:

```text
Bit Function lang
    -> transform
    -> Bitlang functional-semantic representation
    -> preprocess + normalize
    -> Bitlang preprocessed canonical representation
    -> compile
    -> Bitlang compiled procedural representation
```

## Assignment and increment syntax

Bitlang source may provide compound-assignment convenience syntax such as `+=`, `-=`, `*=`, and `/=`. These forms are source-level sugar only and must be expanded during preprocessing into explicit assignment plus the corresponding operation.

For example:

```text
a += b
```

normalizes to the equivalent of:

```text
a = a + b
```

Prefix increment and decrement forms such as `++a` and `--a` are not part of Bitlang. They are intentionally unsupported because they combine mutation and expression evaluation in a way that is easy to misuse.

Postfix increment and decrement forms such as `a++` and `a--` are also not part of the canonical language. If convenience syntax of this kind is ever accepted by a source-facing Bitlang mode or family language, it must be restricted to a standalone mutation statement and normalized before Bitlang preprocessed. It must not be usable as a value-producing expression.

The preferred explicit forms are:

```text
a = a + 1
a = a - 1
```

This keeps mutation and value evaluation separate and removes increment/decrement side-effect semantics from the core language.

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
