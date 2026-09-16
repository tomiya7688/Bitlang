# Bitlang Release and Disposal Properties

This document defines how **Bitlang source** expresses and preprocesses release/disposal semantics.

The canonical fully explicit release-property model belongs to Bitlang Preprocessed:

- https://github.com/tomiya7688/Bitlang_preprocessed/blob/main/PROPERTIES.ja.md

## Source properties

Release behavior and release capability are independent semantic axes.

```text
Auto_release
Manual_release
Releasable
Unreleasable
```

- `Auto_release`: source semantics allow required release/finalization to be performed automatically when the resolved lifetime ends, according to the applicable storage model.
- `Manual_release`: automatic release must not be assumed; an explicit valid release operation is required when the resource requires release.
- `Releasable`: release/finalization may be performed through that declaration or access path.
- `Unreleasable`: that declaration or access path may not release or finalize the resource.

These properties are independent from ownership, lifetime, pointer/reference type, copyability, and movability, although those properties constrain legal combinations.

A borrowed declaration does not own release responsibility and therefore normally resolves to `Unreleasable`.

Examples:

```text
Owned Manual_release Releasable Ptr<My_type>
Borrowed Unreleasable Ref<My_type>
```

## Omission and preprocessing

Source-facing Bitlang may state these properties explicitly or omit them. When omitted, preprocessing resolves them from ownership, type, storage model, lifetime, context, and explicit preprocessing rules.

Invalid property combinations must be rejected rather than silently changing ownership or release semantics.

The current release state itself (`Unreleased / Released`) is also resolved before canonical output, but the exact required Preprocessed representation and its cross-property consistency rules are defined in the separate Bitlang Preprocessed specification.
