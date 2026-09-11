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
| `internal/bitlang/source_text.go` | Represent one untouched Bitlang source input together with diagnostic path metadata. |
| `internal/bitlang/token_kind.go` | Define the language-neutral lexical categories used by Bitlang tokens. |
| `internal/bitlang/token.go` | Represent one lexical token and its original source location. |
| `internal/bitlang/lexer.go` | Convert Bitlang source text into an ordered token stream without applying semantic interpretation. |
| `internal/bitlang/lexer_test.go` | Verify lexical tokenization, spelling preservation, source locations, and malformed literal handling. |
| `internal/bitlang/names.go` | Represent and canonicalize Bitlang names and provide name-based symbol storage. **Split candidate:** canonical names and symbol storage are separate conceptual responsibilities and should be separated before this area grows. |
| `internal/bitlang/names_test.go` | Verify Bitlang name canonicalization and symbol-name behavior. **Split together with `names.go` when its responsibilities are separated.** |
| `internal/bitlang/pipeline.go` | Represent and execute the ordered Bitlang conversion pipeline. **Split candidate:** artifact representation, stage representation, and pipeline orchestration should become separate files/classes as implementation grows. |
| `internal/bitlang/pipeline_test.go` | Verify pipeline stage ordering, transition validation, and execution behavior. |
| `go.mod` | Define the Go bootstrap module and minimum Go language version. |
| `README.md` | Introduce Bitlang, its pipeline, build procedure, and project-level direction. |
| `CODING_RULES.md` | Define language-implementation portability, documentation, and source-structure rules shared across implementations. |
| `GO_CODING_RULES.md` | Define Go-bootstrap-specific application architecture and UPD Commander / Messenger / Processing rules. |
| `FILE_RESPONSIBILITIES.md` | Maintain the canonical mapping from files to their single responsibilities. |

## Planned compiler/runtime responsibilities

These are responsibility slots, not fixed filenames. Names may change when implementation begins, but each responsibility should remain isolated.

| Planned file/class | Responsibility |
| --- | --- |
| `Application` | Coordinate top-level application startup independent of a specific UI. |
| `CuiFrontend` | Accept CUI compiler commands and present official compiler results in terminal form. |
| `CompileRequest` | Represent one requested compilation/conversion operation and its options. |
| `CompilerPipeline` | Coordinate the complete multi-stage official Bitlang conversion flow. |
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

The GUI IDE uses this repository's official Bitlang compiler as its backend.

The IDE MUST NOT contain a second compiler implementation, shadow compiler, simplified parser, or GUI-only semantic pipeline. The official compiler is the single source of truth for preprocessing, parsing, analysis, lowering, diagnostics, and generated artifacts.

The CUI compiler and GUI IDE are therefore sibling frontends over the same official compiler core:

```text
                 +------------------+
                 | Official Bitlang |
                 | compiler core    |
                 +------------------+
                   ^              ^
                   |              |
             +-----------+   +-----------+
             | CUI       |   | GUI IDE   |
             | frontend  |   | frontend  |
             +-----------+   +-----------+
```

Because Bitlang has multiple explicit conversion stages, the IDE should expose those official compiler stage results rather than presenting compilation as a black box.

| Planned file/class | Responsibility |
| --- | --- |
| `IdeApplication` | Start and coordinate the desktop IDE application. |
| `IdeWindow` | Own the primary IDE window layout and top-level UI composition. |
| `SourceEditor` | Edit the currently selected Bitlang source document. |
| `StageNavigator` | Select which official compiler stage or representation is being inspected. |
| `StageViewer` | Display one official compiler artifact without owning conversion logic. |
| `PipelineController` | Request conversions from the official compiler backend and distribute returned artifacts to the IDE. |
| `DiagnosticPanel` | Display diagnostics returned by the official compiler backend. |
| `ArtifactDiffViewer` | Compare official compiler artifacts from selected stages. |
| `ProjectExplorer` | Display and select project source files and related artifacts. |
| `BuildPanel` | Configure and invoke official compiler builds from the IDE. |
| `VmPanel` | Launch or control Bitlang VM execution using VM Assembly produced by the official compiler. |
| `TranslatorPanel` | Select and invoke translators using official compiler VM Assembly output. |

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

Every intermediate representation shown by the IDE MUST come from the official compiler backend.

Each node should be independently inspectable. Ideally the user can select a stage and see:

- the artifact produced by the official compiler at that stage
- diagnostics generated at or before that stage
- the difference from the previous official stage artifact
- whether later stage results are stale after an edit

The compiler core MUST remain usable without the IDE. The IDE MUST call the same official conversion interfaces used by the CUI compiler and tests rather than duplicate compiler logic.

This also means that adding a new compiler stage or changing a stage representation should normally require changing the official compiler first. The IDE should then adapt to the compiler's public stage interface rather than invent its own semantic representation.

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
