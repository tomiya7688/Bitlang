# Bitlang

Bitlang is the **human-writable source language** in the Bitlang toolchain.

This directory owns source-facing syntax and semantics: what a programmer may write, what may be omitted, which properties may be stated explicitly, what preprocessing rules may transform, and how source conveniences are normalized before the strict intermediate language is produced.

## Stage boundary

```text
Bitlang source
    -> Bitlang preprocessor
    -> Bitlang Preprocessed
    -> static analysis / compiler
    -> Bitlang Compiled
```

Bitlang source and Bitlang Preprocessed are intentionally separate languages/projects.

- This repository defines the human-facing language and preprocessing inputs.
- `tomiya7688/Bitlang_preprocessed` defines the fully resolved canonical language emitted by preprocessing.
- `tomiya7688/Bitlang_compiled` defines the lower-level compiled language.

The Bitlang repository should therefore explain **what source code means and how it is normalized**, without duplicating the complete canonical Preprocessed representation.

## Main specification documents

- [LANGUAGE_SPEC.md](LANGUAGE_SPEC.md) — overall Bitlang source language semantics
- [SOURCE_NORMALIZATION.md](SOURCE_NORMALIZATION.md) — relaxed source notation and canonical normalization rules
- [PROPERTIES.md](PROPERTIES.md) — source-facing semantic properties and preprocessing behavior
- [TYPES.md](TYPES.md) — source type system
- [TYPE_ALIASES.md](TYPE_ALIASES.md) — type alias behavior
- [NULLABLE_OPTIONAL.md](NULLABLE_OPTIONAL.md) — source-facing nullable/optional forms
- [RELEASE_PROPERTIES.md](RELEASE_PROPERTIES.md) — source-facing release/disposal controls
- [BORROW_STATE_PREPROCESSOR.md](BORROW_STATE_PREPROCESSOR.md) — preprocessing of borrow state
- [PREPROCESSOR_FUNCTIONS.md](PREPROCESSOR_FUNCTIONS.md) — preprocessor-function model
- [PREPROCESSOR_EXECUTION.md](PREPROCESSOR_EXECUTION.md) — preprocessing execution order, declaration timing, reuse, and start/end activation ranges
- [HEADER_FILES.md](HEADER_FILES.md) — header information, preprocessor functions, and preprocessor macros in header files

## Canonical Preprocessed specifications

The final resolved representation is not specified here. Its canonical home is:

- https://github.com/tomiya7688/Bitlang_preprocessed/blob/main/LANGUAGE_SPEC.ja.md
- https://github.com/tomiya7688/Bitlang_preprocessed/blob/main/PROPERTIES.ja.md
- https://github.com/tomiya7688/Bitlang_preprocessed/blob/main/BORROW_STATE.ja.md

This separation keeps Bitlang documentation substantial and useful for language authors while preventing two repositories from independently defining the same canonical intermediate-language rules.
