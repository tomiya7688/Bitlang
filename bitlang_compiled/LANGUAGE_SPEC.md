# Bitlang Compiled Language Specification

This document defines the Bitlang compiled language specification.

Bitlang compiled is a human-readable, C-like common low-level representation designed for simple lowering to C and to the Bitlang VM virtual assembly.

## Relationship to Bitlang preprocessed

Bitlang compiled receives canonical Bitlang preprocessed input and lowers high-level semantic constructs into a simpler procedural representation.

Conceptually:

```text
Bitlang preprocessed
    -> compile / lower
    -> Bitlang compiled
```

## Functional-language lowering

Functional-programming semantics supported by Bitlang must be lowered into ordinary procedural structures before or while producing Bitlang compiled.

Typical lowering strategies include:

- lambdas -> generated normal functions
- closures -> environment structure + generated function
- captured variables -> explicit fields in the closure environment
- first-class functions -> function reference/value representation compatible with procedural calls
- higher-order calls -> ordinary function calls through function values/references
- pattern matching -> tags/conditions/switch-style control flow
- sum/variant data -> explicit tagged data representation
- partial application/currying -> stored arguments + generated callable wrapper/environment

Bitlang compiled should not preserve multiple equivalent functional notations. By this stage, functional source styles must already have converged through Bitlang preprocessed canonicalization and should be represented in one procedural form suitable for C and VM lowering.

The goal is:

```text
functional source style
    -> Bitlang semantic representation
    -> canonical Bitlang preprocessed
    -> procedural Bitlang compiled
```

This keeps Bit Function lang expressive and natural while keeping Bitlang compiled simple, readable, and compatible with C-like execution models.
