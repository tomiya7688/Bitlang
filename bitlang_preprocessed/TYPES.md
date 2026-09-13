# Bitlang Preprocessed Type System

Bitlang preprocessed uses explicit canonical type names rather than source-level shorthand wherever defaults are known.

## Canonical capitalization

Bitlang identifiers and type names are case-insensitive in meaning. Source code may therefore use any capitalization for equivalent names.

When the preprocessor emits canonical Bitlang preprocessed type names and type constructors, it must normalize their spelling to start with an uppercase letter.

Examples:

```text
int
INT
Int
```

all refer to the same type meaning and normalize to the same canonical type.

Likewise, source spellings equivalent to array, pointer, or reference type constructors normalize to canonical forms beginning with an uppercase letter.

## Canonical numeric naming

For types that have both a radix/base and a bit width, use:

```text
<TypeName><Radix>x<BitWidth>
```

Example:

```text
Int10x32
```

A source declaration such as:

```text
int value
```

must therefore normalize to:

```text
Int10x32 value
```

when `Int10x32` is the default representation of `int`.

`Int` is signed and `Uint` is unsigned.

## Overflow

Overflow is an error by default. Canonical Bitlang arithmetic must not silently wrap overflowing values.

## Radix compatibility

Radix is part of the canonical numeric type.

Operands with different radix values are incompatible for direct arithmetic or bitwise operations. Source code must convert them to a common radix through an appropriate Bitlang pulse function or preprocessor pulse function before the operation reaches canonical Bitlang preprocessed form.

For example:

```text
Int2x32
Int10x32
```

cannot be direct operands of the same arithmetic or bitwise operation.

## Literal normalization

Source literals may use ordinary type inference, but inferred types must be resolved to concrete canonical Bitlang types before or during preprocessing.

## Floating-point types

Floating-point types also carry a radix component. When source code omits the radix, radix 10 is used by default and must be explicit in canonical form.

## Types without radix

Types for which radix/base is not meaningful omit the radix component and encode only the relevant size/count information.

## Strings

For `Str`, the numeric suffix represents the maximum number of characters, not a bit width.

A plain `Str` has no maximum character limit.

A bounded string form records its explicit maximum character count in the canonical type name.

`Char` is not a separate canonical core type. It normalizes to:

```text
Str1x1
```

which represents the one-character string form. The first `1` is currently reserved for possible future string-representation semantics.

## Canonical container and indirection types

Bitlang preprocessed uses the following canonical type constructors:

```text
Array<T>
Ptr<T>
Ref<T>
```

Their meanings are:

- `Array<T>`: array whose element type is `T`
- `Ptr<T>`: pointer to `T`
- `Ref<T>`: reference to `T`

Any source-level notation with the same semantics must normalize to these forms.

Examples:

```text
int[]
Array<int>
```

normalize to:

```text
Array<Int10x32>
```

and equivalent pointer or reference notations must likewise normalize to `Ptr<T>` or `Ref<T>` with the contained type itself fully canonicalized.

Nested forms are represented by composition, for example:

```text
Array<Ptr<Int10x32>>
Ref<Array<Int10x32>>
```

## Canonicalization requirement

Equivalent source-level type spellings, inferred types, omitted default representation details, alternate capitalization, and alternate syntactic forms must normalize to the same concrete Bitlang preprocessed type whenever they have the same semantics.
