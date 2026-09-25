# Bitlang

Bitlang is the human-writable form of the Bitlang language.

This directory owns source-facing syntax and normalization rules: what a programmer may omit, which conveniences and preprocessing-only constructs may be used, and how ordinary Bitlang is normalized into fully explicit Bitlang Preprocessed.

## Stage boundary

```text
Bitlang source
    -> Bitlang preprocessor
    -> Bitlang Preprocessed
    -> static analysis / compiler
    -> Bitlang Compiled
```

Bitlang source and Bitlang Preprocessed are **not separate semantic languages**. They are two normalization states of Bitlang.

- Bitlang source is the authoring form. It may omit applicable properties, use sugar, and contain preprocessing-only constructs.
- Bitlang source may explicitly write every canonical property that Bitlang Preprocessed can contain.
- Bitlang Preprocessed is the fully resolved form of the same language: applicable properties are explicit and source/preprocessing-only shorthand has been removed.
- The repositories are separate only to keep source-facing normalization rules and the canonical fully explicit representation independently maintainable.
- `tomiya7688/Bitlang_compiled` remains the separate lower-level compiled language.

The Bitlang repository should therefore explain **what authors may write or omit and how Bitlang is normalized**, while `tomiya7688/Bitlang_preprocessed` is the canonical home for the requirements of the fully explicit form.

## Main specification documents

- [LANGUAGE_SPEC.md](LANGUAGE_SPEC.md) — overall Bitlang source language semantics
- [STATIC_AND_DIRECT.md](STATIC_AND_DIRECT.md) ([日本語](STATIC_AND_DIRECT.ja.md)) — `static` (static retention) and `Direct` (instance-free access)
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

This repository separation is an ownership boundary for documentation and implementation, not a statement that source and Preprocessed are different semantic languages.
