# Bitlang Compiler Specification

This document defines the Bitlang compiler behavior, transformation stages, interfaces, and implementation requirements.

## Front-end transformation model

Bitlang-family languages are front-end languages that transform into ordinary Bitlang before Bitlang preprocessing begins.

The compiler toolchain should therefore treat transformation and preprocessing as separate stages:

```text
Bitlang family source
    -> family-language transformer
    -> Bitlang
    -> Bitlang preprocessor
    -> Bitlang preprocessed
    -> compiler
    -> Bitlang compiled
```

A transformer is responsible for translating a family language's syntax and language-specific conveniences into equivalent Bitlang meaning. It should not bypass the Bitlang layer unless a future specification explicitly permits it.

## Preprocessor functions

The Bitlang preprocessing system supports preprocessor functions. These are compile-time functions used for definitions, code generation, source transformation, and attribute automation.

The core preprocessor-function language and built-in function groups are specified in `../../PREPROCESSOR_FUNCTIONS.md`.

Preprocessor functions execute before the normal Bitlang-to-Bitlang-preprocessed compilation stage is completed and do not become runtime functions.

The preprocessing implementation should preserve enough source mapping information to report diagnostics against the original source where practical and should allow developers to inspect the generated Bitlang preprocessed output.


## Build artifact retention

Bitlang's canonical pipeline may materialize large intermediate artifacts such as Bitlang Explicit (`.btlex`), Bitlang Low (`.btllw`), Bitlang VM assembly/output (`.btlvm`), or generated C source when the C backend is selected.

These intermediate artifacts are useful for audit, debugging, reproducibility, regression analysis, and inspecting what semantic expansion/lowering actually occurred. Therefore the **normal/default build keeps generated intermediate artifacts**.

Conceptually:

```text
normal build
    source
    + .btlex
    + .btllw
    + selected backend intermediate (.btlvm or generated .c, etc.)
    + final artifact
```

A separate cleanup/final-only build mode may remove disposable intermediate artifacts after the final artifact has been produced successfully.

Conceptually:

```text
final-only / cleanup mode
    build full pipeline
    -> verify final artifact was produced successfully
    -> remove disposable intermediates
    -> keep final artifact
```

The cleanup mode changes artifact retention only. It must not skip semantic stages, validations, safety checks, or lowering steps merely because their serialized artifacts will later be deleted.


### Why the default build keeps intermediates

This default is also part of the user experience.

Bitlang intentionally spends significant compile time on preprocessing, semantic expansion, explicit-property materialization, safety validation, lowering, and backend preparation. Most users will not read the compiler specification in detail. If a normal build deleted every intermediate artifact automatically, that work would be invisible and Bitlang could appear to be only a slow compiler.

Keeping `.btlex`, `.btllw`, and the selected backend intermediate by default gives the user direct evidence of what the compiler produced and why compilation took time. The intermediate files make the transformation cost inspectable rather than mysterious.

This does not make large intermediate output or slow compilation a goal. They are known side effects of Bitlang's explicit, auditable pipeline. Default retention keeps those costs observable instead of hiding them.

`final-only` is therefore an explicit convenience/space-saving mode rather than the default.

### Failure behavior

Automatic intermediate cleanup occurs only after a successful build by default.

If the build fails, any intermediate artifacts already produced should be retained so they can be inspected for diagnostics and debugging. A separate explicit purge operation may remove them when desired.

### Selective retention

Tooling may also support selectively retaining particular stages, for example keeping `.btlex` while deleting backend-generated C, but the baseline behavior is:

```text
default      -> keep intermediates
cleanup mode -> keep final artifact, delete disposable intermediates after success
failure      -> keep produced intermediates
```

Cache storage is a separate concern from user-visible intermediate artifacts. A cleanup mode may remove materialized intermediates while retaining valid internal build-cache entries according to the build-cache policy.

Detailed compiler rules will be added as the compiler design is finalized.
