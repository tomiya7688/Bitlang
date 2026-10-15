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

## Interactive compiler shell

The normal interactive entrypoint of the Bitlang compiler is a persistent compiler shell rather than only a one-shot command invocation.

Starting the compiler without an explicit one-shot command should enter the Bitlang compiler shell when an interactive terminal is available.

The shell exposes the compilation pipeline instead of hiding it. Users can inspect:

- current project and selected target/backend;
- available source and intermediate artifacts;
- artifact hashes/fingerprints;
- which stage produced each artifact;
- whether an artifact is current, stale, invalid, or reusable;
- cache hits/misses;
- validation state;
- dependency changes that caused rebuilds;
- stage timing and produced artifact size.

Conceptually:

```text
$ bitlang
Bitlang> status
Bitlang> build
Bitlang> inspect btlex
Bitlang> resume btllw
Bitlang> backend c
Bitlang> build
```

Exact command spelling may evolve, but the shell model is canonical.

One-shot commands remain supported for scripts, CI, IDE/build-system integration, and reproducible automation.

### Visible stages as a build-speed control

Because stage boundaries and intermediate artifacts are explicit, a user may deliberately reuse a valid existing artifact and resume compilation from a later stage.

Examples include:

- reuse an unchanged `.btlex` instead of preprocessing source again;
- reuse a validated `.btllw` when only the selected backend changed;
- regenerate `.btlvm` or generated C without repeating unrelated earlier work;
- rebuild only downstream stages invalidated by an input/configuration change;
- compare two intermediate artifacts before deciding which stage to rerun.

This is a supported optimization mechanism, not an accidental implementation shortcut.

However, reuse must be validated. An artifact is reusable only when the compiler can establish that its relevant inputs and semantic environment still match. The artifact metadata/fingerprint must cover at least the producing stage, compiler/specification version, relevant source/input hashes, preprocessing configuration, language-version selection, target-independent semantic configuration, and any other dependency that can change the artifact's meaning.

A stale or incompatible artifact must not be silently accepted merely because a file with the expected extension exists.

### Safety and explicit trust

Build acceleration may avoid recomputing a stage when an equivalent validated artifact already exists; it must not silently skip safety obligations.

If a user explicitly asks the shell to trust/import an externally produced intermediate artifact, the toolchain must distinguish that from ordinary validated reuse and either revalidate the artifact at the appropriate boundary or mark the build as using an explicit trust override. Ordinary builds remain fail-closed.

### Shell state

The compiler shell may keep in-memory project/build state between commands, including parsed manifests, dependency graphs, artifact fingerprints, and safe cache metadata. This avoids repeatedly rediscovering unchanged build state during an interactive development session.

Persistent shell state is an optimization only. A clean non-interactive build with the same declared inputs must remain reproducible.

## JSON compiler configuration

Bitlang compiler/build settings may be loaded from a selected JSON configuration file.

The project-root default filename is `bitlang.json`. The compiler shell and one-shot CLI may explicitly select another JSON file. One session/build has one active configuration file; version 1 does not implicitly merge multiple configuration files.

Selection precedence is:

```text
explicit shell/CLI config selection
> project-root bitlang.json
> compiler built-in defaults
```

For ordinary build settings, explicit shell/CLI session overrides take precedence over the active JSON configuration. For language/preprocessor semantics, configuration contributes at the project/default level and does not override stronger explicit source rules.

Configuration parsing is strict. Unknown fields, type mismatches, unsupported schema versions, and duplicate object keys are errors. Configuration is declarative JSON and does not execute shell commands or arbitrary code.

Switching the active configuration invalidates only the pipeline stages affected by changed resolved settings. Backend-only changes should preserve reusable semantic artifacts such as `.btlex` and `.btllw` when their own dependencies remain valid.

The active configuration and the resolved values relevant to each artifact participate in artifact fingerprints so configuration changes cannot silently reuse semantically incompatible intermediates.

Compiler-shell operations include the conceptual equivalents of:

```text
Bitlang> config show
Bitlang> config use bitlang.debug.json
Bitlang> config reload
Bitlang> config none
```

The detailed format and initial schema are defined in `CONFIGURATION.md` and `../../../spec/compiler/config.schema.json`.

### GUI frontend parity

The GUI uses the same compiler core, active JSON configuration, resolved-setting precedence, artifact fingerprints, stage invalidation, cache/reuse rules, and build semantics as the interactive shell and one-shot CLI.

The GUI is not a separate compiler implementation. A build configured identically through the GUI, shell, or CLI must be semantically equivalent.

GUI-only presentation preferences may remain frontend-local when they do not affect compilation.

Detailed compiler rules will be added as the compiler design is finalized.
