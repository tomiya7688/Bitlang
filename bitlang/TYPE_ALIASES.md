# Bitlang Type Aliases

## Char

`Char` is not a fundamentally separate text type.

It is a source-level alias for a string constrained to exactly one character.

Canonical representation:

```text
Str1x1
```

Therefore:

```text
Char c
```

normalizes during preprocessing to the equivalent of:

```text
Str1x1 c
```

For this form, the final `1` represents the one-character limit.

The first `1` is currently a nominal/reserved component of the text-type notation. It is intentionally kept as `1` for now and does not represent a numeric radix in the same sense as numeric types such as `Int10x32`.

This keeps `Char` inside the `Str` family instead of introducing a separate core character type.
