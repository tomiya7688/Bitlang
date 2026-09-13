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

## Explicit declaration properties

Bitlang preprocessed should explicitly state every semantic property that applies to a declaration rather than relying on defaults or omission.

For example, a source declaration inside a function such as:

```text
int a = 4
```

may normalize conceptually to:

```text
Private Readable Writeable Reassignable Initialized Nonnullable Required Int10x32 a = 4
```

The exact set depends on which properties apply to the declaration, but applicable properties should not be left implicit.

## Nullability and presence

Nullability and presence are separate canonical property axes.

Nullability:

```text
Nullable
Nonnullable
```

Presence:

```text
Optional
Required
```

`Nullable` means `null` is a valid value. `Nonnullable` means `null` is invalid.

`Optional` means the value or declaration may be absent. `Required` means it must be present.

All applicable declarations must state both axes explicitly in Bitlang preprocessed. Therefore the absence of `Nullable` must not be interpreted implicitly as non-nullable, and the absence of `Optional` must not be interpreted implicitly as required.

The axes are independent. Examples include:

```text
Required Nullable
Required Nonnullable
Optional Nullable
Optional Nonnullable
```

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
- `Ptr<T>`: raw pointer to `T`
- `Ref<T>`: safe reference to `T`

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

## Ptr semantics

`Ptr<T>` is the canonical raw-pointer type.

It may contain `null`, may be reassigned, may use pointer arithmetic, and may expose or manipulate raw addresses when supported by the target environment.

Invalid-address, dangling-pointer, and equivalent low-level memory errors are the responsibility of code using `Ptr<T>`.

## Ref semantics

`Ref<T>` is the canonical safe-reference type.

It cannot contain `null`, does not permit pointer arithmetic, does not provide normal arbitrary raw-address manipulation, must not outlive its referent, and cannot be rebound after its initial binding.

A reference whose lifetime is provably longer than its referent is invalid and must be rejected before execution.

Assignment through a writable `Ref<T>` changes the referenced value; it does not change which referent the reference is bound to.

## Multidimensional arrays

Bitlang preprocessed does not keep a separate multidimensional-array abstraction.

Source-level dimension-count notation is expanded into nested `Array` constructors.

Conceptually:

```text
Array<T, Dimension=3>
```

normalizes to:

```text
Array<Array<Array<T>>>
```

with `T` itself fully canonicalized.

Fixed-length information is attached to the corresponding nested array level. A source form conceptually equivalent to:

```text
Array<T, Dimension=2, Length=10,20>
```

normalizes to a nested form conceptually equivalent to:

```text
Array<Array<T, 20>, 10>
```

Thus dimensionality is represented structurally by nesting, while each fixed length belongs to the specific array level it constrains.

## Canonicalization requirement

Equivalent source-level type spellings, inferred types, omitted default representation details, alternate capitalization, and alternate syntactic forms must normalize to the same concrete Bitlang preprocessed type whenever they have the same semantics.
