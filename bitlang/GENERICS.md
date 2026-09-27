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

## No runtime generic machinery

Bitlang generics are fully resolved during preprocessing.

After preprocessing, the generic mechanism itself does not remain as runtime machinery.

Conceptually:

```text
generic declaration + concrete type arguments
    -> preprocessing
    -> concrete specialized declarations/types
    -> Bitlang Explicit
```

Bitlang Explicit must not require runtime generic dispatch, runtime generic substitution, runtime generic constraint checks, or unresolved generic type parameters.

This rule exists to avoid carrying generic-resolution overhead into the finished application.

If a source language has explicitly observable runtime type/reflection behavior related to generics, the adapter/preprocessor must translate that behavior into ordinary concrete Bitlang metadata or operations that are actually required by the program. Such metadata is not treated as retained generic machinery and must not be emitted when it is unobservable or unnecessary.

The compiler may deduplicate or share generated machine code internally when doing so preserves all observable semantics, but that optimization must not reintroduce runtime generic resolution.

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
- source-language generic reflection requirements, when explicitly observable, must be materialized as ordinary concrete metadata/operations rather than by retaining a runtime generic mechanism;
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

## Variance

Variance is a canonical generic-parameter property axis.

The states are:

```text
Invariant
Covariant
Contravariant
```

Meaning:

- `Invariant`: no substitution compatibility is derived merely from compatibility between the type arguments.
- `Covariant`: if `Dog` is usable as `Animal`, a corresponding generic specialization may be usable in the same direction, subject to all other Bitlang safety and type rules.
- `Contravariant`: the generic specialization may be usable in the opposite direction, again subject to all other Bitlang safety and type rules.

The ordinary Bitlang default is:

```text
Invariant
```

This is the safety-oriented default. A source language, declaration, or language adapter may explicitly select `Covariant` or `Contravariant` when required to preserve source semantics.

Variance applies per generic type parameter, so a multi-parameter generic may assign a different variance state to each parameter.

Variance never authorizes an otherwise unsafe write, ownership transfer, borrow, lifetime extension, or other invalid operation. If a variance declaration would make a concrete specialization unsafe or internally inconsistent, preprocessing must reject it.

## Decisions

Generic semantic differences will be finalized one property axis at a time.

Variance is finalized. A generic constraint is a compile-time requirement on the concrete type supplied for a generic parameter. Constraint checks happen during preprocessing rather than at runtime. The next generic semantic area to define is which constraint categories Bitlang supports.


## Preprocessor-function constraints

A generic constraint may be supplied by a preprocessing-domain function.

This allows reusable type-acceptance rules to be defined once, imported, and attached to generic declarations.

Conceptually:

```text
constraint preprocessor function
    -> inspect candidate type and its canonical properties
    -> prove that the candidate satisfies the required rule
    -> accept or emit compile-time error
```

A generic declaration may use an imported constraint function, or may express an equivalent constraint directly at the declaration site.

The exact source syntax and function signature are defined separately. The semantic requirements are:

- the constraint executes during preprocessing;
- it may inspect type identity, implemented interfaces, inheritance relationships, canonical properties, and other compile-time type metadata exposed by the preprocessing environment;
- it must not defer the decision to runtime;
- a rejected candidate type is a compile-time error;
- when the constraint requires a safety property and preprocessing cannot prove that property, the specialization is rejected with an error;
- imported constraint functions are reusable preprocessing definitions and do not remain in Bitlang Explicit;
- after successful resolution, only the concrete type and the ordinary resolved Bitlang semantics remain.

This mechanism is intended to absorb source-language-specific generic restrictions without creating separate generic systems for each source language.

Inline constraints and imported preprocessing constraints are semantically equivalent ways to express the acceptance rule. Their exact combination/precedence rules are defined separately.


## Constraint composition

When multiple constraints are attached to the same generic parameter without an explicit boolean operator, they are combined with logical AND.

Conceptually:

```text
T requires A
T requires B

==

T requires (A AND B)
```

AND may also be written explicitly.

OR is never inferred from multiple constraint declarations. If a parameter may satisfy either one condition or another, OR must be written explicitly.

Conceptually:

```text
T requires (A OR B)
```

This rule applies equally to inline constraints and imported preprocessor-function constraints.

Parentheses may be used to group compound constraint expressions. Exact surface syntax is defined separately, but the semantic boolean model is:

```text
implicit multiple constraints -> AND
explicit AND                 -> AND
explicit OR                  -> OR
```

Constraint evaluation remains compile-time only. If the final composed constraint cannot be proven satisfied for the supplied concrete type, preprocessing reports an error.


## Explicit generic arguments

Ordinary Bitlang source requires generic type arguments to be written explicitly when a generic declaration is instantiated or called.

Conceptually:

```text
identity<Int10x32>(10)
Box<Str>
```

Bitlang does not perform implicit generic type-argument inference as a core language rule.

The reason is semantic visibility and determinism: the concrete type used to specialize a generic declaration must be visible in source rather than reconstructed from surrounding expressions.

A language adapter or explicit preprocessing rule may implement source-language-specific generic inference before ordinary Bitlang generic resolution. Such inference is preprocessing behavior, not native Bitlang generic semantics. Its resolved output must contain explicit concrete type arguments before the generic specialization is accepted.

If inferred generic information is ambiguous, incomplete, or cannot be proven safe, preprocessing reports an error.

This keeps native Bitlang generic usage deliberately explicit and avoids hidden type specialization.
