# Nullable and Optional

`Nullable<T>` and `Optional<T>` are distinct Bitlang types.

- `Nullable<T>`: a value of type `T` may also be `null`.
- `Optional<T>`: explicit presence or absence of a value is represented by the type itself.

`T`, `Nullable<T>`, and `Optional<T>` are different types.

Values of different types cannot participate directly in arithmetic, comparison, assignment requiring exact compatibility, or other type-sensitive operations. They must first be converted, unwrapped, or otherwise made type-compatible.

For example:

```text
Int10x32
Nullable<Int10x32>
Optional<Int10x32>
```

are three different types.

`Nullable<T>` and `Optional<T>` are not implicitly interchangeable.

Nested forms such as:

```text
Optional<Nullable<T>>
```

are allowed when that distinction is semantically intended.
