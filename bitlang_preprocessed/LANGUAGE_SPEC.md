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

Semantically equivalent source forms must normalize to the same Bitlang preprocessed representation whenever they express the same meaning.

This is a language-wide rule, not a rule limited to particular syntax categories or types.

Different source spellings, aliases, shorthand, convenience syntax, family-language syntax, or alternate notations may exist before preprocessing, but if their semantics are identical they must converge to one canonical preprocessed form.

The preprocessed language should therefore minimize alternate spellings and alternate semantic forms. Convenience syntax belongs in Bitlang or in Bitlang-family front ends, not in Bitlang preprocessed unless the distinction is semantically necessary.

Examples of normalization targets include:

- compound assignment into explicit assignment and arithmetic
- alternate function-value syntax into one function-reference representation
- alternate lambda syntax into one lambda/closure representation
- alternate pattern-matching syntax into one canonical match representation
- shorthand declarations into fully explicit declarations
- inferred/default attributes into explicit attributes where the canonical form requires them
- family-language-specific constructs into their equivalent Bitlang representation
- alternate array syntax into one canonical array representation
- alternate pointer syntax into one canonical pointer representation
- alternate reference syntax into one canonical reference representation
- aliases and shortened names into fully resolved canonical references

If two source forms are not semantically identical, they must not be merged merely because their surface syntax is similar.

## Explicit property declaration

Bitlang preprocessed should explicitly declare all semantic properties that apply to a declaration or value wherever those properties are represented by the language.

A property may have been written explicitly in Bitlang source or may have been inferred, defaulted, generated, or transformed by preprocessing. That distinction does not survive as ambiguity in Bitlang preprocessed: the resolved final property itself must be emitted.

Defaults, source shorthand, omitted attributes, or preprocessor-only assumptions must therefore be resolved before Bitlang preprocessed is produced.

This includes properties such as:

- `Readable` / `Unreadable`
- `Writeable` / `Unwriteable`
- `Reassignable` / `Unreassignable`
- `Owned` / `Borrowed`
- `Copyable` / `Uncopyable`
- `Movable` / `Unmovable`
- `Moved` / `Unmoved`
- `Initialized` / `Uninitialized`
- nullability and optionality properties
- visibility/export properties
- explicit lifetime properties such as `Local_lifetime`, `Function_lifetime`, `Object_lifetime`, `Module_lifetime`, or `Static_lifetime` where applicable
- instance requirements and other independent semantic properties defined by the language

The general rule is that Bitlang preprocessed should not require later stages to guess which side of a defined property axis applies. If a property is meaningful and applicable to the target, preprocessing should emit its resolved canonical state explicitly.

Properties that are genuinely not applicable to a target do not need meaningless placeholder declarations.

This explicit-property rule exists so that Bitlang preprocessed can act as a deterministic semantic contract for static analysis and later compilation stages.

## Copy and move property declaration

When copying or moving is meaningful for a declaration or value, Bitlang preprocessed must explicitly state the resolved capability using:

```text
Copyable
Uncopyable
Movable
Unmovable
```

`Copyable` and `Uncopyable` form the copy-capability axis. `Movable` and `Unmovable` form a separate move-capability axis.

The properties are independent from ownership. For example:

```text
Owned Uncopyable Movable Ptr<My_type>
```

may describe an owned resource that cannot be duplicated but may have its ownership/value transferred.

An ordinary scalar value may instead resolve to properties such as:

```text
Copyable Movable Int10x32
```

The preprocessed form must not rely on later stages to infer copyability or movability from the type or ownership state. If source Bitlang omits these properties, preprocessing resolves and emits them.

## Move-state declaration

When move state is meaningful, Bitlang preprocessed must explicitly state one of:

```text
Unmoved
Moved
```

`Unmoved` means the declaration currently retains a usable value or ownership state.

`Moved` means that value or ownership has been transferred away from the declaration. Unless another explicit semantic rule has restored the source, ordinary reading, reuse, release, or another move through the moved declaration is invalid.

A normal valid move transitions the source declaration from `Unmoved` to `Moved`.

A valid reinitialization with a new value may transition a declaration back to `Unmoved` where such reinitialization is permitted.

Preprocessor transformations are allowed to change move-state properties, including deliberately rewriting `Moved` to `Unmoved`. Such an override is not ordinary inference: it must come from an explicit preprocessing rule, configuration, or source-directed transformation because it can alter runtime meaning and potentially re-enable access to a resource that was previously moved.

If preprocessing cannot prove that an explicit `Moved -> Unmoved` override is safe, it may emit a warning. If the resulting state is provably invalid, it must produce an error.

The final resolved move-state property is emitted in Bitlang preprocessed so later stages do not need to reconstruct move history merely to determine the declaration's current state.

## Lifetime-property declaration

Bitlang preprocessed must explicitly state the resolved lifetime category of a declaration or resource when lifetime is meaningful for that target.

Canonical lifetime properties include:

```text
Local_lifetime
Function_lifetime
Object_lifetime
Module_lifetime
Static_lifetime
```

The lifetime property is independent from ownership. `Owned` / `Borrowed` answers who controls lifetime responsibility, while the lifetime property states the region for which the value remains valid.

A source declaration may specify its lifetime directly, or it may omit it and allow preprocessing to infer the appropriate lifetime from context. In either case, the Bitlang preprocessed result must contain the same explicit resolved lifetime property when the resulting semantics are the same.

A local declaration should therefore not rely on the compiler inferring that it is local merely from syntax. Its resolved lifetime property is emitted explicitly.

Reference and borrow checks may use these properties directly. A dependent or borrowed value must not have a lifetime that can exceed the value or resource on which it depends.

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

## Initialization-state declaration

Bitlang preprocessed must explicitly state the initialization state of every variable declaration.

The canonical initialization-state properties are:

```text
Initialized
Uninitialized
```

A declaration that already contains a valid value must be marked `Initialized`.

A declaration that exists without a valid value must be marked `Uninitialized`.

The preprocessed form must not depend on inference or source shorthand to determine whether a variable is initialized. Even when the source syntax makes the state obvious, preprocessing must emit the explicit canonical initialization property.

Reading an `Uninitialized` variable as a value is invalid. A valid initialization operation changes its state to `Initialized` before subsequent reads are permitted.

Initialization state is an independent semantic property and may be combined with other explicit properties such as readability, writability, reassignment, ownership, storage, and lifetime properties.

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
