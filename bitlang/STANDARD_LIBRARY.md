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

## Residual garbage collector

Bitlang does not use garbage collection as its primary memory-management model.

The normal rule is:

```text
ownership/lifetime analysis
    -> automatic release insertion where applicable
    -> explicit release behavior
    -> Bitlang Explicit / Bitlang Low
```

The standard library may provide an optional residual garbage collector as a final safety net, primarily for native C-backend builds.

Its job is limited to memory that remains allocated after the ordinary Bitlang release model has done its work. It is not the normal owner of ordinary Bitlang objects and does not replace deterministic release.

Conceptually:

```text
ordinary Bitlang release model
    -> frees everything that can be deterministically resolved
    -> residual allocations, if any
    -> optional residual collector
```

Bitlang Low and VM-oriented lowering are expected to preserve complete explicit release behavior. The residual collector exists mainly to protect against backend/integration cases where heap allocations may still remain after translation to C or interaction with external/native facilities.

A conforming C optimizer is not assumed to be allowed to remove semantically required release behavior incorrectly. The collector therefore protects against residual allocation states, not against compiler miscompilation.

The collector:

- must not change ordinary ownership semantics;
- must not make unreleased programmer-owned resources silently valid;
- must not replace deterministic finalization of files, sockets, locks, handles, or similar resources;
- must not be required when analysis/backend lowering can prove that no residual managed heap allocation remains;
- follows the standard-library reachability rule and is omitted entirely when unused/unneeded.

The exact residual-tracking mechanism and C-backend integration are defined separately.

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
