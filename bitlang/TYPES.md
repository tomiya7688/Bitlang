# Bitlang Type System

## Type-name representation

Bitlang may encode representation information directly in the type name.

For types that have both a radix/base and a bit width, the canonical naming pattern is:

```text
<TypeName><Radix>x<BitWidth>
```

Example:

```text
Int10x32
```

means an `Int` represented in radix 10 with a 32-bit width.

A source-level short type name uses that type's default representation. Therefore:

```text
int value
```

is equivalent, after preprocessing, to:

```text
Int10x32 value
```

## Signed and unsigned integers

`Int` is signed.

`Uint` is unsigned.

Examples:

```text
Int10x32
Uint10x32
```

## Overflow

Numeric overflow is an error by default.

Bitlang must not silently wrap an overflowing value unless a separate operation explicitly requests different overflow behavior.

## Radix-aware numeric forms

Numeric types may use other radix values while preserving the same bit width.

Examples:

```text
Int2x32
Int8x32
Int10x32
Int16x32
```

Radix is part of the numeric type representation, not merely display formatting.

Values with different radix types cannot be used together directly in arithmetic or bitwise operations. One side must first be converted to the same radix through an appropriate Bitlang pulse function or preprocessor pulse function.

For example, an `Int2x32` and an `Int10x32` are not directly compatible operands even when they represent the same numeric value and bit width.

Bitwise operations are performed over the fixed-width bit representation after operand types have been made compatible. Radix-2 exposes individual bits directly, radix-8 groups bits in sets of three, and radix-16 provides a compact bit-oriented form.

## Literals and type inference

Literals use ordinary type inference unless an explicit type is supplied.

The inferred type must ultimately resolve to a concrete Bitlang type before canonical preprocessing is complete.

## Floating-point types

Floating-point types also carry a radix component.

If the radix is omitted, radix 10 is the default.

The same general naming principle applies:

```text
<TypeName><Radix>x<BitWidth>
```

## Types without radix

If radix/base is not meaningful for a type, the type omits the radix component and uses only the relevant size/count information.

## Strings

`Str` uses its numeric suffix as a maximum character count rather than as a bit width.

A plain:

```text
Str
```

has no maximum character limit by default.

A bounded string form specifies a maximum number of characters explicitly.

`Char` is source-level shorthand for a one-character string representation and normalizes to `Str1x1`.

The first `1` in `Str1x1` is currently reserved and does not yet have a finalized semantic meaning. It may be assigned a useful string-representation property later.

## Arrays

Bitlang source may describe array dimensionality and fixed lengths as attributes inside the `Array<...>` form rather than requiring the programmer to write repeated nested array constructors manually.

Conceptually:

```text
Array<T, Dimension=3>
```

means a three-dimensional array of `T`.

Fixed lengths may also be supplied in the array attributes. When more than one dimension is present, a fixed length may be supplied for each dimension.

Conceptually:

```text
Array<T, Dimension=2, Length=10,20>
```

represents a two-dimensional array whose dimensions have fixed lengths 10 and 20.

These source-level dimension and length attributes are convenience information. The preprocessor expands dimensionality into nested canonical `Array` types. Fixed-length information is attached to the corresponding canonical array level rather than retained as a separate multidimensional-array abstraction.

Therefore programmers may write dimensionality compactly, while Bitlang preprocessed uses ordinary nested arrays internally.

## Preprocessing rule

Source-level shorthand is for convenience only. Bitlang preprocessing resolves shorthand, inferred types, defaults, and compact array-dimension notation into the concrete canonical type representation used by Bitlang preprocessed.
