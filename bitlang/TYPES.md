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
