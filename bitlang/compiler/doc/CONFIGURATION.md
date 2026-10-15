# Bitlang Compiler Configuration

Bitlang compiler/build settings may be loaded from a JSON configuration file.

The compiler shell, GUI frontend, and one-shot CLI can select which configuration file is active. They use the same compiler configuration loader and semantics. This allows projects to keep separate configurations such as development, release, VM, C-backend, or target-specific builds without changing source code.

## 1. File format

Configuration files are JSON.

The default project configuration filename is:

```text
bitlang.json
```

Alternative configuration files may use any filename ending in `.json`, for example:

```text
bitlang.debug.json
bitlang.release.json
configs/vm.json
configs/c-x64.json
```

Version 1 uses a strict machine-readable schema:

```text
spec/compiler/config.schema.json
```

Unknown keys, invalid value types, unsupported schema versions, and duplicate object keys are errors. Bitlang does not silently ignore misspelled configuration keys.

JSON comments and trailing commas are not part of the configuration format.

## 2. Selecting the active configuration

One compiler session/build has one active configuration file.

Selection order:

```text
explicit shell/GUI/CLI config selection
> project-root bitlang.json
> compiler built-in defaults
```

The compiler does not implicitly merge several JSON configuration files in version 1.

Conceptual compiler-shell operations:

```text
Bitlang> config show
Bitlang> config use bitlang.debug.json
Bitlang> config use configs/c-x64.json
Bitlang> config reload
Bitlang> config none
```

`config none` disables automatic use of the project configuration for the current session and uses explicit session overrides plus built-in defaults.

The GUI and one-shot CLI also support an explicit configuration path. Exact UI/CLI spelling may evolve, but selecting a configuration in any frontend must produce equivalent compiler behavior.

Changing the active configuration causes the compiler to recompute which existing artifacts remain reusable.

## 3. Configuration precedence

For ordinary compiler/build settings:

```text
explicit CLI / compiler-shell / GUI session override
> active JSON configuration
> compiler built-in default
```

For language/preprocessor semantics, the JSON configuration supplies project-level configuration/defaults. It does not override a stronger explicit source declaration where the language's existing precedence rules say the source declaration wins.

Thus configuration-file support does not replace Bitlang's existing scoped property/version/preprocessor precedence rules.

## 4. Initial configuration sections

Version 1 defines these top-level sections:

```text
schema_version
build
language
preprocessor
cache
shell
```

### build

Build/toolchain settings include:

```text
backend
target
output_dir
artifact_retention
incremental
```

`artifact_retention` currently supports:

```text
keep
final-only
```

The normal default remains `keep`.

### language

`default_version` supplies the project-level default Bitlang language-version selector when one is not selected by a stronger source/version rule.

### preprocessor

Project-level preprocessing configuration may include the parent project and namespace mounts.

Namespace mount paths retain the existing project-relative semantics defined by `NAMESPACE_MOUNTS_AND_PROJECT_INHERITANCE.md`.

### cache

Controls whether validated compiler/build cache reuse is enabled and where compiler-managed cache data is stored.

### shell

Controls presentation preferences for the interactive compiler shell, such as whether stage timing and artifact sizes are shown automatically.

## 5. Example

```json
{
  "schema_version": 1,
  "build": {
    "backend": "c",
    "target": "x64",
    "output_dir": "build/debug",
    "artifact_retention": "keep",
    "incremental": true
  },
  "language": {
    "default_version": "1.0"
  },
  "preprocessor": {
    "parent_project": null,
    "namespace_mounts": [
      {
        "kind": "tree",
        "path": "src",
        "namespace": "Game"
      }
    ]
  },
  "cache": {
    "enabled": true,
    "directory": ".bitlang-cache"
  },
  "shell": {
    "show_stage_timing": true,
    "show_artifact_sizes": true
  }
}
```

## 6. Configuration and artifact reuse

The compiler must not treat every configuration change as invalidating every stage.

Each setting declares which pipeline stage(s) it can affect. Artifact fingerprints use only the configuration dependencies relevant to that artifact.

Examples:

```text
preprocessor/language semantic change
    -> invalidates .btlex and downstream

lowering-semantic change
    -> invalidates .btllw and downstream

backend/target change
    -> may preserve valid .btlex/.btllw and invalidate backend output only

artifact_retention/output presentation change
    -> does not change semantic artifact identity
```

This stage-aware dependency model is required so the configuration system improves usability without destroying Bitlang's explicit artifact-reuse performance advantage.

Switching configuration files therefore causes targeted invalidation rather than an unconditional clean rebuild.

## 7. Reproducibility

The active configuration identity and the relevant canonicalized setting values are recorded in build metadata/fingerprints.

A clean one-shot build, compiler-shell build, and GUI build using the same source, compiler/specification versions, active configuration, and explicit overrides must have equivalent semantics.

Configuration paths themselves are not sufficient identity: the compiler fingerprints the relevant resolved setting values so changing file contents is detected.

## 8. Security and path handling

Configuration files may select output/cache/project-relative paths, so path use must be validated.

- cleanup must never delete source/input merely because a configuration path points to it;
- build cleanup must not escape the permitted output root;
- namespace-mount project-relative semantics remain unchanged;
- configuration loading must not execute arbitrary code;
- JSON values do not perform implicit shell expansion or command execution.

Configuration is declarative build/preprocessing data, not a scripting language.

## 9. GUI frontend parity

The GUI is a frontend to the same compiler configuration/build engine used by the shell and one-shot CLI.

The GUI must not define a second configuration format, silently translate settings into different semantics, or maintain hidden build-semantic settings that cannot be represented by the shared compiler configuration/session model.

The GUI provides equivalents of the shell configuration operations:

```text
show active configuration
select another JSON configuration
reload the selected file
disable automatic project configuration for the session
show resolved settings and their source/precedence
```

Changing the active configuration in the GUI uses the same stage-aware invalidation and artifact-reuse decisions as the shell/CLI.

GUI-only presentation state such as window layout, panel sizes, editor theme, and other non-compiler UI preferences may be stored separately. Such UI state must not affect compiler semantics or artifact fingerprints unless it corresponds to an explicit shared compiler setting.
