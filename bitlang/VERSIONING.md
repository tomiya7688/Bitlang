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

## Language-version number format

Bitlang language semantics are versioned at the **major.minor** level.

The source-facing language-version forms are:

```text
ver1.0
ver1.1
ver2.0
```

A patch release does not change the Bitlang language specification or the meaning of valid source code.

Therefore, for example:

```text
1.0.0 -> 1.0.1
```

must not introduce a language syntax or semantic change. Patch versions may contain compiler/tooling bug fixes, implementation fixes, diagnostics, documentation corrections, or other changes that preserve the same language contract.

For this reason, source language-version declarations do not normally include the patch component. `ver1.0` identifies the language specification shared by compatible 1.0.x toolchain releases.

Compiler/toolchain version identity remains separate and may still use a full version such as `1.0.1`.

## Major-only floating selection

A source scope may specify only the major language version:

```text
ver1
```

This means:

> use the newest supported Bitlang language specification in the 1.x series that this compiler/toolchain provides.

For example, if a compiler supports:

```text
ver1.0
ver1.1
ver1.3
ver2.0
```

then:

```text
ver1 -> ver1.3
```

The resolution is constrained to the requested major version. `ver1` never silently selects `ver2.x`.

The resolved major.minor version must be exposed in diagnostics/build metadata so the actual interpretation is inspectable.

A major-only declaration is intentionally floating. A project that requires reproducible language-version selection should pin the major.minor form such as `ver1.1`.

## Selecting a language version

Compilation may select the Bitlang language version under which source code is interpreted.

The exact CLI/project-file spelling is defined separately, but the semantic operation is:

```text
compile(source, language_version)
```

The selected language version controls the compatibility frontend used for parsing and interpreting source-language behavior. This selection may be inherited from project/file scope or overridden by a smaller class/type/function scope as defined below.

A project should be able to pin this version for reproducible builds. A command-line selection may override or supply the project selection where the build system permits it.

If no version is specified, the compiler may use its current default language version. Reproducible projects should pin the intended major.minor version explicitly rather than depend on that default or on a floating major-only form.

## Scoped language-version declarations

Bitlang language versions may be selected at structural scope boundaries rather than only once for an entire project.

A source file, class/type, or function may declare the language version used to interpret that scope.

Conceptually:

```text
file version V3

class Legacy_parser version V2
{
    function parse_old_format version V1
    {
        ...
    }

    function parse_current
    {
        ...
    }
}
```

The version declaration uses ordinary Bitlang notation and belongs to the `@preprocesser` domain. Bitlang does not introduce a separate DOCTYPE-like or foreign metadata syntax solely for language versions.

The exact built-in/function/declaration spelling is defined separately, but semantically it is a preprocessing declaration/directive placed at the beginning of the target scope, before the body whose syntax and semantics depend on that version.

It is not an ordinary compiler/runtime declaration and never becomes runtime state.

The effective language version uses innermost-scope precedence:

```text
function version
    > class/type version
    > file version
    > file-targeted header version
    > project version
    > compiler/toolchain default
```

A nested scope without its own version inherits the effective version of its nearest enclosing scope.

This allows a project to migrate incrementally. A file may use the newest language version while one old class or function remains pinned to an older supported version.

A language-version declaration is semantic configuration, not ordinary runtime data, and disappears before Bitlang Explicit after the selected compatibility frontend has normalized the scope.

### Preprocessor ownership of version declarations

Language-version selection belongs to preprocessing.

Conceptually:

```text
Bitlang source notation
    -> @preprocesser version declaration/directive
    -> select compatibility frontend for target scope
    -> parse/normalize target scope
    -> consume version declaration
    -> Bitlang Explicit
```

The declaration may appear at file, class/type, or function scope according to the normal scoped-version rules.

A version declaration must not survive into Bitlang Explicit as runtime/compiler state. Its only semantic effect is deciding how the corresponding source scope is interpreted and normalized.

Because version selection happens before ordinary version-specific parsing, the minimal Bitlang syntax used to express this preprocessing declaration is part of the stable version-neutral envelope grammar.

### Same-scope conflicts

At the same structural scope, a version written directly in the source scope takes precedence over an inherited/header/project default.

Two incompatible direct version requirements for the same scope are an error unless an explicit preprocessing rule selects one deterministically.

A header may provide the version for a target file or declaration when that target does not provide a more local/direct version itself.

## Version envelope parsing

Because the selected language version may change the grammar used to parse a scope, version metadata must be discoverable before the body of that scope is parsed under a version-specific grammar.

Bitlang therefore treats scope-leading `@preprocesser` version information as part of a small version-neutral source envelope.

This envelope is still Bitlang syntax. It is not a separate header language. The compiler only gives the minimal Bitlang preprocessing syntax required for version selection a cross-version-stable interpretation before dispatching the rest of the scope to its selected compatibility frontend.

Conceptually:

```text
version-neutral envelope reader
    -> discovers file/class/function version metadata
    -> selects supported compatibility frontend for that scope
    -> parses and normalizes the scope body
```

The envelope syntax itself must remain stable enough for supported compilers to discover the version declaration without first knowing the version of the body.

Nested versioned scopes may therefore be dispatched to different supported compatibility frontends and then normalized into the same current canonical Bitlang semantics.

This mechanism is intentionally similar in spirit to a document header declaring how the following content should be interpreted.

## Automatic compatible-version selection

The standard preprocessing/tooling environment may provide a helper that selects the newest supported Bitlang language version under which a target scope compiles successfully.

Conceptually:

```text
select_latest_compatible_version(scope)
```

or, when an explicit search range is desired:

```text
select_latest_compatible_version(scope, newest, oldest)
```

The exact function name and surface syntax are defined separately.

The selection algorithm is conceptually:

```text
try newest supported candidate
    -> if the complete target scope parses, normalizes, and validates:
           select that version
    -> otherwise try the next older supported candidate
    -> continue until success or the allowed range is exhausted
```

A failed candidate does not weaken diagnostics or safety. A candidate is acceptable only if the scope is valid under that language version and can still normalize into current safe Bitlang semantics.

The chosen version becomes the effective explicit version binding for that preprocessing result. Tooling should expose the selected version in diagnostics/build metadata and may offer to write the resulting version declaration into the source or an associated header so future builds are pinned reproducibly.

Automatic probing is a migration convenience, not permission to reinterpret one compilation unit nondeterministically. Reproducible builds should materialize or otherwise pin the selected result rather than depending indefinitely on whatever versions a future compiler happens to support.

This mechanism makes large-scale upgrades practical:

```text
try current version
    -> compatible: keep current
    -> incompatible: pin only the smallest affected scope to the newest older version that works
```

The goal is to avoid rewriting an entire project merely because a small class or function still depends on historical syntax or semantics.

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
