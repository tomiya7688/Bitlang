# Bitlang Preprocessed Language Specification

This document defines the Bitlang preprocessed language specification.

Bitlang preprocessed is the fully expanded and normalized form produced after preprocessing Bitlang source.

## Position in the language family

Bitlang preprocessed is the canonical form of Bitlang and should be treated as the closest representation to the true core language semantics.

Human-written Bitlang may contain conveniences, abbreviations, multiple equivalent notations, generated constructs, and preprocessor functions. Bitlang-family languages may also introduce very different syntax and programming styles. These differences must be removed before or during preprocessing so that equivalent meaning converges toward the same canonical representation.

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

## Canonicalization requirements

Where practical, semantically equivalent Bitlang programs should normalize to the same Bitlang preprocessed structure.

The preprocessed language should therefore minimize alternate spellings and alternate semantic forms. Convenience syntax belongs in Bitlang or in Bitlang-family front ends, not in Bitlang preprocessed unless the distinction is semantically necessary.

Examples of normalization targets include:

- compound assignment into explicit assignment and arithmetic
- alternate function-value syntax into one function-reference representation
- alternate lambda syntax into one lambda/closure representation
- alternate pattern-matching syntax into one canonical match representation
- shorthand declarations into fully explicit declarations
- inferred/default attributes into explicit attributes where the canonical form requires them
- family-language-specific constructs into their equivalent Bitlang representation

## Assignment normalization

Compound assignment operators do not exist in Bitlang preprocessed.

Forms such as:

```text
a += b
a -= b
a *= b
a /= b
```

must be normalized into explicit assignment forms equivalent to:

```text
a = a + b
a = a - b
a = a * b
a = a / b
```

Increment and decrement operators are not part of Bitlang preprocessed. Prefix forms such as `++a` and `--a` are unsupported in Bitlang. Postfix forms such as `a++` and `a--` must also be eliminated before the preprocessed form is produced, if any source-facing language accepts them at all.

Mutation should therefore be represented explicitly, for example:

```text
a = a + 1
a = a - 1
```

Bitlang preprocessed must not contain value-producing increment/decrement expressions.

## Modifier design principle

A Bitlang modifier should represent one primary semantic property.

Independent properties should be expressed by combining independent modifiers rather than assigning multiple unrelated meanings to one modifier.

This principle is especially important for lifetime, storage, ownership, accessibility, and instance requirements because these properties may vary independently.

## `static`

`static` means **static retention**.

A `static` target exists in a statically retained form and keeps the information or state associated with that target for the lifetime of its static retention region.

The fundamental meaning of `static` is the same regardless of the kind of target to which it is applied. Where applicable, this includes variables, functions, and type members.

`static` does not by itself mean that an instance is unnecessary.

In particular, Bitlang must not overload `static` with the conventional secondary meaning of "callable or accessible without an instance".

The exact start and end of each static retention region, initialization timing, destruction behavior, and lower-level representation are defined separately by the relevant lifetime and compilation rules.

## Instance-free access modifier

Bitlang provides a separate modifier whose sole semantic purpose is to state that access to a target does **not require an instance**.

The final source-level keyword or symbol for this modifier is not yet fixed. Until it is named, this specification refers to it as the **instance-free modifier**.

The instance-free modifier:

- removes the requirement to provide or construct an instance in order to access or call the target
- does not imply static retention
- does not change the target's lifetime by itself
- does not imply state retention
- is semantically independent from `static`

`static` and the instance-free modifier may therefore be combined.

Conceptually, the two properties form independent axes:

| Static retention | Instance-free | Meaning |
|---|---|---|
| no | no | ordinary instance-dependent target |
| yes | no | statically retained target that still requires an instance for access |
| no | yes | instance-free target without static retention |
| yes | yes | statically retained and instance-free target |

No additional compound modifier is required to represent these combinations. Their meaning is obtained from the composition of the two simple modifiers.

Conceptually:

```text
static target
```

means only "retain this target statically", while:

```text
<instance-free> target
```

means only "this target does not require an instance".

The combination:

```text
static <instance-free> target
```

means both properties simultaneously.

Whether every combination is legal for every target kind is defined by the specification for that target kind. The semantic properties themselves remain independent.

## Functional-programming normalization

Functional-programming concepts may remain semantically visible in Bitlang preprocessed when they are part of the core meaning, but each concept should have one canonical representation.

For example:

- a lambda should use one canonical lambda representation
- a captured environment should be represented consistently
- first-class function values should use one function-value model
- pattern matching should use one canonical case/match model
- partial application and currying should normalize into a defined canonical representation before procedural lowering

Bitlang preprocessed is not required to be pleasant to write manually. Human readability is useful, but determinism, explicit semantics, and normalization take priority.
