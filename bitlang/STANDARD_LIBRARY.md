# Bitlang Standard Library Lowering Policy

## Core rule

The Bitlang standard library is not automatically carried into lower compilation stages as one monolithic runtime library.

Only standard-library declarations and support code that are actually required by the program may descend through the compilation pipeline.

Conceptually:

```text
standard library
    -> preprocessing / reachability resolution
    -> only reachable required declarations
    -> Bitlang Explicit
    -> Bitlang Low
    -> backend output
```

A declaration is not considered required merely because:

- it exists in the standard library;
- its module/package is installed;
- a broad standard-library module is imported;
- another unused declaration in the same source file/module refers to it.

## Reachability

Required standard-library code is determined transitively.

If user code references standard-library declaration A, and A requires B and C, then A, B, and C are retained.

Unreachable declarations are not emitted to lower stages.

Conceptually:

```text
user code -> A -> B
              -> C

retained: A, B, C
removed: every unrelated standard-library declaration
```

The same principle applies to types, functions, constants, generated helpers, generic specializations, metadata, and other standard-library artifacts.

## Import is not use

Importing or making a standard-library name available does not by itself require the corresponding runtime implementation to survive lowering.

Only semantic use that requires the declaration at a later stage makes it reachable.

Type-only information may be consumed during compilation and disappear before runtime when no runtime representation is required.

Preprocessor-only standard-library utilities execute during preprocessing and do not survive into Bitlang Explicit or runtime output unless they deliberately generate required ordinary Bitlang declarations.

## Generics

Standard-library generics follow the normal Bitlang generic rules.

Only concrete specializations that are actually referenced are generated.

Unused generic declarations and unused specializations do not descend into Bitlang Explicit.

## Optional garbage collector

Bitlang does not require a garbage collector as part of the core runtime.

The standard library may provide one or more explicitly selected garbage-collection facilities for programs that want runtime-managed memory.

Conceptually:

```text
ordinary Bitlang program
    -> no garbage collector runtime

program that uses standard-library GC
    -> selected GC API/runtime support becomes reachable
    -> only required collector code and metadata descend
```

Using the garbage collector must be explicit through the standard-library API or an explicit language-adapter/preprocessing rule that lowers to that API.

Ordinary Bitlang allocations do not silently become GC-managed merely because the GC library is available or imported.

The exact public API, handle/reference type, and collection algorithm are defined separately. The standard library may provide multiple collector implementations as long as their observable contracts are explicit.

GC support follows the normal standard-library reachability rule. If no GC-managed allocation or collector operation is semantically used, collector code, tracing metadata, runtime tables, and initialization must not be emitted to lower stages.

Garbage collection is primarily a memory-management facility. Code must not rely on an unspecified collection moment for deterministic resource finalization. Resources that require deterministic release, such as external handles, should continue to use the ordinary Bitlang ownership/release/finalization model unless a GC API explicitly defines a stronger contract.

The presence of a collector does not bypass Bitlang safety validation. GC-managed references, borrows, finalization hooks, and generated runtime support must remain consistent with the applicable ownership, lifetime, and access rules.

## Required runtime/compiler support

A language or standard-library operation may require helper code even when the helper is not named directly by user source.

Such helper code is considered reachable when lowering of a used operation requires it.

Examples include a required bounds-check helper, trap implementation, numeric helper, allocator operation, or backend support routine.

The compiler must retain only the support actually required by the selected operations and target semantics.

## Initialization and side effects

Unused standard-library declarations must not become reachable merely because they could have initialization side effects.

Standard-library design should avoid hidden global initialization whose only purpose is triggered by library presence rather than semantic use.

If a used standard-library feature requires initialization, that initialization and its transitive dependencies are reachable together with the feature.

## Optimization boundary

Removal of unused standard-library code is a pipeline rule, not merely an optional final-linker optimization.

The compiler should avoid lowering unreachable standard-library declarations into Bitlang Explicit / Bitlang Low in the first place when their unreachability is already known.

Later dead-code elimination and linker garbage collection may still remove additional unreachable backend code, but they are secondary safeguards rather than the primary standard-library inclusion model.

## Design goal

The standard library may be broad without imposing proportional runtime size or execution overhead on every program.

Bitlang therefore follows this principle:

```text
available standard library != emitted runtime library

used semantics -> emitted required implementation
unused semantics -> do not descend
```

This rule supports Bitlang's general goal of accepting expensive compile-time analysis in exchange for small, fast, safety-validated runtime output.
