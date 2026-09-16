# Nullable and Optional in Bitlang source

This document defines **source-facing Bitlang** nullable and optional forms.

The canonical final nullability/optionality properties emitted after preprocessing are owned by Bitlang Preprocessed:

- https://github.com/tomiya7688/Bitlang_preprocessed/blob/main/PROPERTIES.ja.md

## Source forms

Bitlang source distinguishes:

```text
T
Nullable<T>
Optional<T>
```

- `Nullable<T>`: source notation expressing that a value of `T` may also be `null`.
- `Optional<T>`: source notation expressing explicit presence or absence of a value.

They are not implicitly interchangeable in source semantics. Operations requiring exact compatibility must first convert, unwrap, or otherwise establish compatible semantics.

Nested source forms such as:

```text
Optional<Nullable<T>>
```

are allowed when both distinctions are semantically intended.

## Preprocessing boundary

Source syntax is not the authority for the final canonical property serialization.

Before Bitlang Preprocessed is emitted, preprocessing resolves the applicable state explicitly. In particular, nullability becomes the explicit canonical property pair:

```text
nullable
unnullable
```

and cannot remain unspecified at the Preprocessed boundary when nullability is applicable.

Optionality is likewise resolved into the canonical representation defined by Bitlang Preprocessed. The source notation remains documented here because it is part of what a Bitlang programmer may write; the fully normalized representation is documented in the separate Preprocessed repository.
