# Bitlang Header Files

Bitlang supports optional header files as source-level preprocessing support files.

A header file is not a required companion to a source file and is not required to share the source file's name. Its relationship to a source declaration, class, module, or other target is determined by the header information itself rather than by filename matching.

## Allowed contents

A Bitlang header file may contain:

- ordinary header information;
- `@preprocesser` functions;
- preprocessor macros.

Header files therefore belong to the source/preprocessing side of the Bitlang toolchain.

The syntax of preprocessor functions is the same Bitlang syntax documented in `PREPROCESSOR_FUNCTIONS.md`. A function declared with `@preprocesser` is executed by the preprocessor rather than as ordinary compiler/runtime code.

Preprocessor macros contained in a header file are likewise preprocessing constructs and are expanded or otherwise resolved during preprocessing.

## Header information and target association

Header information may contain the ordinary descriptive or structural information expected of a header-like source artifact.

A header file may explicitly state which declaration or structure it describes. For example, it may state that the contained header information belongs to a particular class. The physical filename does not need to match the target class name.

Conceptually:

```text
some_header_file
    -> header information says: target = Some_class
    -> treated as header information for Some_class
```

The exact syntax used to express the target relationship is specified separately, but filename equality is not part of the semantic requirement.

Header information is deliberately weaker than an interface contract. It may describe, annotate, expose preprocessing information for, or otherwise identify a class or other target without requiring that target to satisfy a full interface-style semantic contract.

Header files are optional. A Bitlang source file, class, or module does not require a corresponding header file unless a separate project/module rule explicitly requires one.

## Import from source

Ordinary Bitlang source may import a header file.

Importing a header makes its applicable header information, preprocessor functions, and preprocessor macros available to preprocessing according to the normal visibility and preprocessing rules.

The exact import syntax and conflict-resolution rules are specified separately.

## Preprocessed boundary

Header files are preprocessing-only artifacts.

They do **not** become mandatory explicit structures in Bitlang Preprocessed. The opposite rule applies: header-only information, preprocessor functions, preprocessor macros, and header association metadata are consumed during preprocessing and disappear before Bitlang Preprocessed is emitted.

Only their resolved effects on the program may remain.

Conceptually:

```text
Bitlang source
+ optional imported header files
        -> preprocessing
        -> header metadata/macros/preprocesser functions are consumed
        -> resulting normalized program
        -> Bitlang Preprocessed
```

Bitlang Preprocessed must therefore not depend on the continued existence of a Bitlang header file or on hidden header-only preprocessing state.

If a header causes code generation, declaration modification, property changes, aliases to be resolved, or other semantic transformations, the resulting explicit program state is what proceeds to Bitlang Preprocessed; the header construct itself does not.

## Not yet fixed

This document does not yet define:

- the filename extension for Bitlang header files;
- the exact syntax and vocabulary of header information;
- the exact syntax for associating a header with a class/module/declaration;
- the exact source import syntax;
- duplicate/conflicting header handling;
- ordering rules when several header files apply.

Those points remain separate source-language decisions.
