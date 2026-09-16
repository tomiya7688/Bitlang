# Bitlang Header Files

Bitlang supports optional header files as source-level preprocessing support files.

A header file is not a required companion to a source file and is not required to share the source file's name. Its relationship to a source declaration, class, module, or other target is determined by the header information itself rather than by filename matching.

## Primary purpose

The primary practical purpose of a Bitlang header file is to provide a reusable place for **preprocessor macros** and related preprocessing definitions.

Header information and `@preprocesser` functions are supported, but a header file should not be treated as the canonical declaration of a class, module, interface, or runtime program structure. In normal use it is closer to a reusable preprocessing-definition file than to a mandatory declaration contract.

This means a header file may exist solely to collect macros and preprocessing helpers for import by ordinary Bitlang source.

## Allowed contents

A Bitlang header file may contain:

- ordinary header information;
- `@preprocesser` functions;
- preprocessor macros;
- preprocessing `export` declarations/directives.

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

## Header export

Bitlang also provides an `export` concept for preprocessing definitions.

This `export` is a **header-file-only source feature**. Ordinary Bitlang source files do not use this form of `export` to publish runtime/compiler declarations.

Its purpose is to make preprocessing definitions from a header available at project scope rather than requiring every source file to import the header individually.

Conceptually:

```text
header file
    -> export preprocessing macro/function/definition
    -> available to preprocessing across the project
```

An exported preprocessing definition may therefore be used by applicable files, classes, functions, or other source structures throughout the project according to the normal preprocessing visibility, activation, ordering, and conflict rules.

`import` and header `export` solve opposite distribution problems:

```text
import
    -> one source location explicitly brings a header into its preprocessing context

header export
    -> a header exposes preprocessing definitions to the project-wide preprocessing context
```

Header export does not turn a preprocessor function into compiler/runtime code and does not create a runtime module export. It only affects preprocessing availability.

The exact surface syntax for exporting one definition, several definitions, or an entire header is specified separately. Likewise, conflict handling and explicit project-level restrictions may further constrain exported definitions.

## Preprocessed boundary

Header files are preprocessing-only artifacts.

They do **not** become mandatory explicit structures in Bitlang Preprocessed. The opposite rule applies: header-only information, preprocessor functions, preprocessor macros, header association metadata, and header `export` controls are consumed during preprocessing and disappear before Bitlang Preprocessed is emitted.

Only their resolved effects on the program may remain.

Conceptually:

```text
Bitlang source
+ optional imported/exported header preprocessing definitions
        -> preprocessing
        -> header metadata/macros/preprocesser functions/export controls are consumed
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
- the exact header export syntax;
- whether export applies to selected definitions or an entire header by shorthand;
- duplicate/conflicting header handling;
- ordering rules when several header files apply.

Those points remain separate source-language decisions.
