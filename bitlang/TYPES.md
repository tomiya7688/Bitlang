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

## Radix-aware numeric forms

Numeric types may use other radix values while preserving the same bit width.

Examples:

```text
Int2x32
Int8x32
Int10x32
Int16x32
```

These forms describe the same width of integer storage with different radix representations.

Radix conversion should preserve the numeric value and bit width unless an explicit narrowing, widening, signedness change, or overflow rule says otherwise.

This makes binary and octal forms directly usable for low-level and bit-oriented operations. In particular, radix-2 representation exposes individual bits directly, while radix-8 representation groups bits in sets of three. Radix-16 may likewise be used as a compact bit-oriented representation.

Bitwise operations are defined over the underlying fixed-width bit pattern rather than over the human-readable spelling of the number. The radix component therefore provides an explicit representation while the bit width defines the available bit positions.

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

## Preprocessing rule

Source-level shorthand is for convenience only. Bitlang preprocessing resolves shorthand and defaults into the concrete canonical type representation used by Bitlang preprocessed.
