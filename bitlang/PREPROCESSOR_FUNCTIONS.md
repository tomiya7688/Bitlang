# Bitlang Preprocessor Functions

Bitlang preprocessor functions are compile-time functions used to reduce repetitive source code and remove transformation or code-generation bottlenecks without weakening the strict Bitlang language model.

They combine ideas similar to code generation, macro definitions, and source transformation tools.

## Core language

The preprocessor-function language should support at least:

- variables
- user-defined functions
- `if` / `else`
- `for`
- `switch`
- arithmetic operators: `+`, `-`, `*`, `/`, `%`
- comparison operators
- boolean/logical operators
- string concatenation and basic string operations

## Built-in function groups

Built-ins should cover these groups:

### Definition

- define values, aliases, or reusable definitions
- remove definitions
- test whether a definition exists

### Text processing

- replace text
- split and join text
- substring operations
- formatting

### Bitlang generation

- emit Bitlang source text
- emit variables
- emit functions
- emit structs
- emit classes
- emit types

Structured generation is preferred when practical so generated output remains valid Bitlang.

### Source and structure inspection

- enumerate functions
- enumerate variables
- enumerate types
- inspect attributes
- access structured source information

### Environment and diagnostics

- current source path and source directory
- include or read source-related text where permitted by the preprocessing environment
- compile-time errors
- compile-time warnings

## Attribute automation

Preprocessor functions may automatically add, remove, inspect, or provide default attributes for Bitlang declarations.

The system should support operations equivalent to:

- add an attribute to a target
- remove an attribute from a target
- test whether an attribute is present
- set default attributes for a declaration category
- infer attributes through static analysis when the result is provable

Attribute automation exists to reduce repetitive annotation while keeping Bitlang's semantic model explicit.

Attributes that can be proven without changing program behavior may be attached automatically. Attributes that could change runtime behavior or program meaning must require an explicit rule, configuration, or declaration rather than heuristic inference alone.

## Role in the Bitlang family

Bitlang-family languages generally transform into ordinary Bitlang before normal Bitlang preprocessing.

```text
Bitlang family source
    -> transform
    -> Bitlang
    -> preprocessor functions
    -> Bitlang preprocessed
    -> compiler
    -> Bitlang compiled
```

Preprocessor functions may also be reused as part of family-language transformation infrastructure when appropriate.
