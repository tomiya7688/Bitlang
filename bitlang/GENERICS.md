# Bitlang Generics

## Core meaning

Bitlang generics are compile-time type parameterization.

A generic parameter represents a type. It does not represent a runtime value and is not a cast.

Generic parameters may be declared on:

- class
- struct
- interface
- function

Examples:

```text
class Box<T> { ... }
struct Pair<T> { ... }
interface Comparable<T> { ... }
T identity<T>(T value) { ... }
```

Generic arguments are resolved by preprocessing. Bitlang Explicit must not contain unresolved generic type parameters.

The same generic declaration with the same type arguments denotes the same specialization.

## Property-driven language compatibility

Different source languages define generics differently.

Bitlang does not model each source language as a separate generic system. Instead, semantic differences that matter to program behavior are represented as canonical generic properties and are selected by source code, preprocessing rules, or language adapters.

Conceptually:

```text
foreign generic semantics
    -> language adapter
    -> generic declaration + canonical generic properties
    -> preprocessing
    -> concrete Bitlang declarations/types
    -> Bitlang Explicit
```

A property belongs to the generic system only when changing it changes observable typing or program semantics.

Examples of differences that may require generic properties include:

- variance / assignment compatibility between different generic arguments;
- runtime visibility of the original generic type argument when source-language behavior depends on it;
- rules governing what kinds of types satisfy a generic parameter;
- source-language-specific generic compatibility rules that cannot be represented by an ordinary type/interface constraint alone.

Implementation-only choices are not source semantic properties.

For example, whether the compiler emits one machine-code body per concrete specialization or shares an implementation internally is a lowering/optimization decision when observable behavior is unchanged.

## Preprocessing boundary

Generic-specific properties are preprocessing semantics.

Before Bitlang Explicit is produced:

1. generic arguments are resolved to concrete types;
2. generic constraints and compatibility rules are validated;
3. source-language generic properties are applied;
4. the resulting concrete declarations/types receive all ordinary Bitlang properties required by their semantics;
5. unresolved generic parameters and preprocessing-only generic metadata are removed.

If generic resolution is ambiguous, contradictory, or unsafe, preprocessing fails rather than selecting a behavior silently.

## Decisions

Generic semantic differences will be finalized one property axis at a time.

The first axis to define is variance: whether a generic specialization may be substituted for another specialization when their type arguments have an inheritance/compatibility relationship.
