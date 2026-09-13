# Bitlang Preprocessed Type System

Bitlang preprocessed uses explicit canonical type names rather than source-level shorthand wherever defaults are known.

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

## Types without radix

Types for which radix/base is not meaningful omit the radix component and encode only the relevant size/count information.

## Strings

For `Str`, the numeric suffix represents the maximum number of characters, not a bit width.

A plain `Str` has no maximum character limit.

A bounded string form records its explicit maximum character count in the canonical type name.

## Canonicalization requirement

Equivalent source-level type spellings and omitted default representation details should normalize to the same concrete Bitlang preprocessed type whenever possible.
