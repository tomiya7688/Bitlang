# Bitlang File Responsibility Table

This file is the canonical responsibility map for source files and major planned files in the Bitlang repository.

The project follows these structural rules:

- one file = one responsibility
- one function = one operation
- a responsibility should map naturally to one class when reimplemented in an object-oriented language
- if one responsibility record becomes too large to describe cleanly, the responsibility is too broad and the file/class SHOULD be split
- adding a new source file requires adding or updating its responsibility record here
- moving behavior between files requires updating this table in the same change

The responsibility text should remain short. It should describe *what the file owns*, not enumerate every helper function it contains.

## Current files

| File | Responsibility |
| --- | --- |
| `cmd/bitlang/main.go` | Start the Bitlang CUI application and hand control to the command-line frontend. |
| `internal/bitlang/names.go` | Represent and canonicalize Bitlang names and provide name-based symbol storage. **Split candidate:** canonical names and symbol storage are separate conceptual responsibilities and should be separated before this area grows. |
| `internal/bitlang/names_test.go` | Verify Bitlang name canonicalization and symbol-name behavior. **Split together with `names.go` when its responsibilities are separated.** |
| `internal/bitlang/pipeline.go` | Represent and execute the ordered Bitlang conversion pipeline. **Split candidate:** artifact representation, stage representation, and pipeline orchestration should become separate files/classes as implementation grows. |
| `internal/bitlang/pipeline_test.go` | Verify pipeline stage ordering, transition validation, and execution behavior. |
| `go.mod` | Define the Go bootstrap module and minimum Go language version. |
| `README.md` | Introduce Bitlang, its pipeline, build procedure, and project-level direction. |
| `CODING_RULES.md` | Define implementation, portability, documentation, and source-structure rules. |
| `FILE_RESPONSIBILITIES.md` | Maintain the canonical mapping from files to their single responsibilities. |

## Planned compiler/runtime responsibilities

These are responsibility slots, not fixed filenames. Names may change when implementation begins, but each responsibility should remain isolated.

| Planned file/class | Responsibility |
| --- | --- |
| `Application` | Coordinate top-level application startup independent of a specific UI. |
| `CuiFrontend` | Accept CUI compiler commands and present compiler results in terminal form. |
| `CompileRequest` | Represent one requested compilation/conversion operation and its options. |
| `CompilerPipeline` | Coordinate the complete multi-stage Bitlang conversion flow. |
| `SourceArtifact` | Represent original Bitlang source input. |
| `PreprocessedArtifact` | Represent Bitlang Preprocessed output. |
| `CompiledArtifact` | Represent Bitlang Compiled output. |
| `TreeObjectArtifact` | Represent Bitlang TreeObject output. |
| `VmAssemblyArtifact` | Represent Bitlang VM Assembly output. |
| `Preprocessor` | Convert Bitlang source into Bitlang Preprocessed. |
| `StaticAnalyzer` | Perform static analysis over Bitlang Preprocessed. |
| `Advisor` | Produce advisory diagnostics that do not change compilation semantics. |
| `Compiler` | Convert analyzed Bitlang Preprocessed into Bitlang Compiled. |
| `TreeLowerer` | Convert Bitlang Compiled into Bitlang TreeObject. |
| `VmAssemblyLowerer` | Convert Bitlang TreeObject into Bitlang VM Assembly. |
| `Diagnostic` | Represent one compiler diagnostic independent of UI rendering. |
| `DiagnosticFormatter` | Convert diagnostics into textual presentation. |

## Planned GUI IDE responsibilities

The GUI IDE is not a replacement for the compiler core. It is another frontend over the same stage APIs used by the CUI compiler and tests.

Because Bitlang has multiple explicit conversion stages, the IDE should expose those stages rather than presenting compilation as a black box.

| Planned file/class | Responsibility |
| --- | --- |
| `IdeApplication` | Start and coordinate the desktop IDE application. |
| `IdeWindow` | Own the primary IDE window layout and top-level UI composition. |
| `SourceEditor` | Edit the currently selected Bitlang source document. |
| `StageNavigator` | Select which conversion stage or representation is being inspected. |
| `StageViewer` | Display one immutable stage artifact without owning conversion logic. |
| `PipelineController` | Request stage conversions from the compiler core and distribute their results to the IDE. |
| `DiagnosticPanel` | Display compiler errors, warnings, and advisor messages. |
| `ArtifactDiffViewer` | Compare adjacent or selected stage artifacts to show what each conversion changed. |
| `ProjectExplorer` | Display and select project source files and related artifacts. |
| `BuildPanel` | Configure and invoke complete builds from the IDE. |
| `VmPanel` | Launch or control Bitlang VM execution using generated VM Assembly. |
| `TranslatorPanel` | Select and invoke VM Assembly translators such as RISC-V or ARM targets. |

## IDE pipeline concept

The intended IDE view is approximately:

```text
Source
  -> Preprocessed
  -> Compiled
  -> TreeObject
  -> VM Assembly
  -> VM / Translator output
```

Each node should be independently inspectable. Ideally the user can select a stage and see:

- the artifact produced at that stage
- diagnostics generated at or before that stage
- the difference from the previous stage
- whether later stages are stale after an edit

The compiler core MUST remain usable without the IDE. The IDE MUST call the same conversion interfaces used by the CUI compiler rather than duplicate compiler logic.

## Split rule for this table

A responsibility record is intentionally expected to fit in one short sentence.

If a record repeatedly needs multiple clauses such as "A and B and C", detailed sub-bullets, or exceptions to explain what the file owns, that is evidence that the corresponding file/class has accumulated more than one responsibility and SHOULD be split.

When a split occurs:

```text
OldFile: A and B
```

should become something like:

```text
FileA: A
FileB: B
```

The responsibility table is therefore not only documentation; it is an architectural size check.
