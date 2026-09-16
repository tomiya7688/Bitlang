# Bitlang Properties

Bitlang source exposes semantic properties so programmers and preprocessor functions can control program meaning explicitly when needed, while still allowing source code to omit properties that the preprocessor can resolve safely.

The **canonical fully explicit property model is owned by `tomiya7688/Bitlang_preprocessed`**. This document defines the Bitlang source-facing role of those properties and how they participate in preprocessing.

Canonical property specification:

- https://github.com/tomiya7688/Bitlang_preprocessed/blob/main/PROPERTIES.ja.md
- borrow-state details: https://github.com/tomiya7688/Bitlang_preprocessed/blob/main/BORROW_STATE.ja.md

## Source-facing principle

Bitlang is human-writable. A declaration may explicitly state semantic properties, or may omit properties when the preprocessor can resolve them from type, declaration kind, lexical context, ownership, lifetime, configuration, defaults, control flow, or explicit preprocessing rules.

Before Bitlang Preprocessed is emitted, every applicable semantic axis must be resolved to its final explicit state.

Therefore the source layer and Preprocessed layer have different responsibilities:

```text
Bitlang source
    explicit properties + omitted/inferred properties + preprocessing rules
        -> preprocess / normalize
Bitlang Preprocessed
    fully resolved explicit property state
```

## Property vocabulary available to Bitlang source

Bitlang source may use the same semantic vocabulary that is normalized into Bitlang Preprocessed, including the following independent axes.

### Visibility

```text
Public / Private
Protected / Unprotected
Exported / Unexported
```

Scope visibility, inheritance visibility, and module export visibility are separate concerns.

### Access and reassignment

```text
Readable / Unreadable
Writeable / Unwriteable
Reassignable / Unreassignable
```

Bitlang does not collapse these into a single broad `mutable / immutable` category. A binding may, for example, be unreassignable while the referenced object remains writeable.

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

```text
Initialized / Uninitialized
nullable / unnullable
Optional / Required
```

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

## Preprocessor behavior

Preprocessor functions may inspect, add, remove, or rewrite source properties before canonical output is finalized.

An ordinary inferred value may fill an omitted source property. A rewrite that deliberately weakens or changes already-resolved semantics is a semantic override and may require diagnostics.

Examples include changing borrow state, move state, release state, ownership-related properties, or other restrictions. A dangerous state that cannot be proven invalid may produce a warning; a provably invalid final state must be rejected.

## Stage boundary

This repository defines **how Bitlang source expresses or omits these properties and how preprocessing is allowed to resolve them**.

The exact canonical property set, required explicitness, final state vocabulary, and cross-property consistency rules belong to the separate Bitlang Preprocessed specification. Later compiler stages must consume that resolved representation rather than reconstructing source omissions.
