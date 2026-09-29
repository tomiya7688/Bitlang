# Bitlang Language Versioning and Compiler Compatibility

## Core policy

Bitlang itself does not guarantee backward compatibility between language-specification versions.

A newer Bitlang language specification is free to remove, replace, or redefine source-language features when doing so improves the language.

Backward compatibility is a responsibility of the Bitlang compiler/toolchain, not of the current language specification.

Conceptually:

```text
Bitlang language
    -> current specification only

Bitlang compiler
    -> may accept source written for supported historical Bitlang versions
```

## Selecting a language version

Compilation may select the Bitlang language version under which source code is interpreted.

The exact CLI/project-file spelling is defined separately, but the semantic operation is:

```text
compile(source, language_version)
```

The selected language version controls the compatibility frontend used for parsing and interpreting source-language behavior.

A project should be able to pin this version for reproducible builds. A command-line selection may override or supply the project selection where the build system permits it.

If no version is specified, the compiler may use its current default language version. Reproducible projects should pin the intended version explicitly rather than depend on that default.

## Compatibility frontend

Historical language behavior is not kept as permanent branches inside the current core language semantics.

Instead, a compiler that supports an older Bitlang version uses a version-specific compatibility frontend.

Conceptually:

```text
Bitlang version N source
    -> version N compatibility parser/adapter
    -> version N semantic compatibility transformation
    -> current canonical Bitlang semantics
    -> normal preprocessing
    -> Bitlang Explicit
    -> Bitlang Lowerer
    -> Bitlang Low
```

The compatibility frontend may rewrite old syntax, defaults, properties, library-facing constructs, or other historical behavior into the current semantic model.

If an old behavior cannot be represented safely and deterministically in the current model, the compiler must reject that source for that compatibility target rather than silently change its meaning.

## Compiler support window

A compiler version advertises the Bitlang language versions it supports.

Supporting an old language version means the compiler can correctly interpret and translate that historical specification. It does not mean the current Bitlang language specification still contains the old syntax or behavior.

A compiler is not required to support every historical Bitlang version forever. Removal of a compatibility frontend is a compiler/toolchain support-policy change rather than a change to the current language semantics.

## Safety

Backward-compatibility mode does not weaken current safety validation.

Historical source is translated into current canonical semantics and then passes the same mandatory safety checks as ordinary current-version Bitlang.

A compatibility frontend must not preserve a historical behavior by emitting a program that violates current mandatory safety invariants.

If faithfully preserving historical behavior would require an unsafe or undefined result that current Bitlang does not permit, compilation fails unless the behavior is represented by an explicitly defined safe/unsafe contract allowed by the current language.

## Standard library compatibility

Language-version compatibility may require a matching standard-library compatibility profile.

Old source-level library names or APIs may be translated or provided through compatibility modules where their behavior can be preserved.

The current standard library itself does not need to retain every historical API indefinitely merely to make old source compile.

Compatibility shims are compiler/toolchain compatibility assets and should disappear from the final program when they are compile-time-only.

## Version identity

Language-version identity and compiler-version identity are separate.

For example:

```text
compiler version: X
language version: Y
```

A newer compiler may compile an older supported language version.

The selected language version must be available to diagnostics, build metadata, and preprocessing so version-specific compatibility behavior is deterministic and inspectable.
