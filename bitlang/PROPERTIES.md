# Bitlang Properties

Bitlang source exposes semantic properties so programmers and preprocessor functions can control program meaning explicitly when needed, while still allowing source code to omit or abbreviate properties that the preprocessor can resolve safely.

The **canonical fully explicit property model is owned by `tomiya7688/Bitlang_preprocessed`**. This document defines the Bitlang source-facing role of those properties and how they participate in preprocessing.

Canonical property specification:

- https://github.com/tomiya7688/Bitlang_preprocessed/blob/main/PROPERTIES.ja.md
- borrow-state details: https://github.com/tomiya7688/Bitlang_preprocessed/blob/main/BORROW_STATE.ja.md

## Source-facing principle

Bitlang is human-writable. A declaration may explicitly state semantic properties, may omit properties when the preprocessor can resolve them, and may use source-facing shorthand that expands into one or more canonical property axes.

Resolution may use the type, declaration kind, lexical context, ownership, lifetime, module/project configuration, defaults, control flow, explicit preprocessing rules, or another source construct whose semantics are defined to imply a property.

Before Bitlang Preprocessed is emitted, every applicable semantic axis must be resolved to its final explicit state.

Therefore the source layer and Preprocessed layer have different responsibilities:

```text
Bitlang source
    explicit properties
    + omitted/inferred properties
    + shorthand / aliases / property bundles
    + type/declaration implications
    + preprocessing rules
        -> preprocess / normalize
Bitlang Preprocessed
    fully resolved explicit property state
```

## Loose source notation and canonical normalization

Bitlang source is intentionally allowed to be less verbose than Bitlang Preprocessed.

The following are valid source-language design mechanisms:

- omit an applicable property axis when preprocessing can resolve it mechanically;
- write only the property axes that matter to the programmer and let preprocessing fill the rest;
- use an alias for a property or a property bundle;
- use a compact source qualifier that expands into multiple independent canonical properties;
- let a declaration kind or type imply properties defined by that construct;
- define module/project preprocessing rules that provide defaults or reusable property bundles;
- use preprocessor functions to inspect and transform the source property set.

These conveniences exist only on the Bitlang/source side. They must not survive as unresolved shorthand in Bitlang Preprocessed.

For example, a project or language rule may define a concise source qualifier conceptually as:

```text
Mutable_value
    -> Readable Writeable Reassignable
```

or a compact resource profile as:

```text
Owned_resource
    -> Owned Unborrowed Movable Unmoved Releasable Unreleased
```

The names above illustrate the normalization mechanism; concrete built-in aliases are specified separately when adopted. A user-defined or module-defined alias may provide the same kind of expansion through preprocessing.

A bundle does not erase the independence of the canonical axes. Any axis not fixed by the bundle is still resolved independently before Preprocessed output.

This means a concise Bitlang declaration can normalize into a deliberately verbose canonical declaration without changing its semantics.

Conceptually:

```text
int a = 4
```

may normalize into a form containing explicit visibility, retention, instance-access requirement, access, reassignment, ownership, borrow state, copy/move capability, move state, release state, lifetime, initialization state, initialization trigger, nullability, optionality, const state, and any other applicable canonical properties.

The exact resulting property set depends on the declaration and applicable preprocessing rules.

## Conflict and precedence rules

Source convenience must not create silent ambiguity.

Resolution follows these principles:

1. Directly stated canonical properties are explicit semantic requirements.
2. Source shorthand and property bundles are expanded before final normalization.
3. Inferred/default values fill only still-unresolved axes.
4. If two explicit source requirements produce contradictory states on the same axis, preprocessing must diagnose the conflict instead of silently choosing one.
5. A deliberate semantic override may change an already resolved state only through an explicit preprocessing rule, configuration, or source-directed transformation.
6. The final Bitlang Preprocessed output must contain one valid resolved state for every applicable axis.

Thus source syntax may be permissive in spelling and verbosity while the stage boundary remains strict.

## Property vocabulary available to Bitlang source

Bitlang source may use the same semantic vocabulary that is normalized into Bitlang Preprocessed, including the following independent axes.

### Visibility

```text
Public / Private
Protected / Unprotected
Exported / Unexported
```

Scope visibility, inheritance visibility, and module export visibility are separate concerns.

### Retention and instance access

```text
Static / Dynamic
Instance_required / Instance_unrequired
```

`Static / Dynamic` is the static-retention axis.

- `Static`: the target is retained statically rather than following ordinary non-static retention.
- `Dynamic`: the target is not statically retained and follows the applicable ordinary lifetime/storage relationship.

`Dynamic` in this property axis does **not** mean dynamic typing, late binding, or general mutability.

`Instance_required / Instance_unrequired` is a separate access axis.

- `Instance_required`: access to the applicable member requires an instance of its containing type.
- `Instance_unrequired`: the applicable member can be accessed without an instance.

These two axes are independent. A target may therefore be static while still requiring an instance, or may be non-static while not requiring an instance, when that combination is meaningful for the declaration kind.

The canonical applicability is:

- `Static / Dynamic`: variables, fields, and functions;
- `Instance_required / Instance_unrequired`: fields and functions.

Bitlang source provides `Direct` as a source-facing shorthand for:

```text
Direct
    -> Instance_unrequired
```

`Direct` does not imply `Static`, does not alter lifetime, and does not change initialization or ownership semantics.

Static retention is also distinct from the lifetime axis. `Static` and `Static_lifetime` describe related but separate semantic concerns and must satisfy the applicable consistency rules rather than being treated as the same property.

### Access and reassignment

```text
Readable / Unreadable
Writeable / Unwriteable
Reassignable / Unreassignable
```

The canonical model does not collapse these into a single broad `mutable / immutable` category. A source shorthand may group them for convenience, but preprocessing must expand the shorthand back into the independent axes.

A binding may, for example, be unreassignable while the referenced object remains writeable.

### Ownership and borrow

```text
Owned / Borrowed
Unborrowed / Shared_borrowed / Exclusive_borrowed
```

Ownership responsibility and current borrow state are separate axes.

### Copy and move

```text
Copyable / Uncopyable
Movable / Unmovable
Unmoved / Moved
```

Capability and current move state are separate.

### Release / disposal

```text
Auto_release / Manual_release
Releasable / Unreleasable
Unreleased / Released
```

Release policy, release capability, and current release state are separate.

### Initialization and nullability

Initialization state and initialization trigger are separate axes.

```text
Initialized / Uninitialized

Declaration_initialization
Owner_initialization
First_reach_initialization
First_use_initialization
Manual_initialization

nullable / unnullable
Optional / Required
```

The initialization-state axis describes whether a valid value currently exists.

The initialization-trigger axis describes when automatic initialization occurs:

- `Declaration_initialization`: initialize when the declaration's storage instance reaches its normal declaration-initialization point.
- `Owner_initialization`: initialize when the declaration's owning object, type, module, or corresponding owner is initialized.
- `First_reach_initialization`: initialize once when execution first reaches the declaration for that storage instance.
- `First_use_initialization`: initialize once on the first valid use of the declaration for that storage instance.
- `Manual_initialization`: no automatic initialization trigger; initialization must be performed explicitly.

The canonical initialization-trigger axis applies to variables and fields.

`Static / Dynamic` does not itself select an initialization trigger. A language adapter, source declaration, project rule, or preprocessing rule may select the trigger required to preserve the source language's semantics.

For example, a static local may use `First_reach_initialization`, while a static field may use `Owner_initialization` or `First_use_initialization` depending on the intended source-language semantics.

Nullability is not represented by an omitted default at the Preprocessed boundary. Source syntax may be compact, but preprocessing resolves the final state explicitly.

### Lifetime

```text
Local_lifetime
Function_lifetime
Object_lifetime
Module_lifetime
Static_lifetime
```

The source may state lifetime directly or leave it for preprocessing where the applicable lifetime is mechanically resolvable.

### Const

```text
Const / Unconst
```

`Const` is stronger than only making one access path unwritable or unreassignable; it represents a strongly fixed semantic value.

A source shorthand named or described as immutable must not be assumed to mean `Const` unless its expansion rule explicitly says so.

## Preprocessor behavior

Preprocessor functions may inspect, add, remove, expand, or rewrite source properties before canonical output is finalized.

An ordinary inferred value may fill an omitted source property. A shorthand expansion may fill multiple axes at once. A rewrite that deliberately weakens or changes already-resolved semantics is a semantic override and may require diagnostics.

Examples include changing borrow state, move state, release state, ownership-related properties, or other restrictions. A dangerous state that cannot be proven invalid may produce a warning; a provably invalid final state must be rejected.

## Stage boundary

This repository defines **how Bitlang source expresses, abbreviates, or omits properties and how preprocessing is allowed to resolve them**.

The exact canonical property set, required explicitness, final state vocabulary, and cross-property consistency rules belong to the separate Bitlang Preprocessed specification. Later compiler stages must consume that resolved representation rather than reconstructing source omissions or source-only shorthand.
