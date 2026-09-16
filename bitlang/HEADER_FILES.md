# Bitlang Header Files

Bitlang supports header files as a source-level container for information and preprocessing definitions that should be available before ordinary compilation proceeds.

## Allowed contents

A Bitlang header file may contain:

- header information;
- `@preprocesser` functions;
- preprocessor macros.

Header files therefore participate primarily in the preprocessing side of the Bitlang toolchain.

The syntax of preprocessor functions is the same Bitlang syntax documented in `PREPROCESSOR_FUNCTIONS.md`. A function declared with `@preprocesser` is executed by the preprocessor rather than as ordinary compiler/runtime code.

Preprocessor macros contained in a header file are likewise preprocessing constructs and are expanded or otherwise resolved before Bitlang Preprocessed output is finalized.

## Stage relationship

Conceptually:

```text
Bitlang header file
    header information
    + preprocessor functions
    + preprocessor macros
        -> Bitlang preprocessor
        -> normalized Bitlang source / Bitlang Preprocessed
```

Header-file information that affects program semantics must be fully resolved before the Bitlang Preprocessed boundary. Later compiler stages must not depend on unresolved preprocessor macros or hidden header-only preprocessing state.

## Not yet fixed

This document does not yet define:

- the filename extension for Bitlang header files;
- the exact syntax and vocabulary of header information;
- inclusion/import syntax for header files;
- duplicate/conflicting header handling;
- ordering rules when several header files apply.

Those points remain separate source-language decisions.
