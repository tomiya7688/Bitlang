# Bitlang File Responsibility Table

Canonical map for implementation files. Keep each record to one short responsibility.

## Core source

| File | Responsibility |
| --- | --- |
| `cmd/bitlang/main.go` | Start the CUI and hand control to the frontend. |
| `internal/cui/run.go` | Route one CUI invocation and format terminal output. |
| `internal/cui/run_test.go` | Verify CUI routing and exit behavior. |
| `internal/bitlang/source_text.go` | Represent untouched source text and diagnostic path metadata. |
| `internal/bitlang/token_kind.go` | Define lexical token categories. |
| `internal/bitlang/token.go` | Represent one token and source location. |
| `internal/bitlang/lexer.go` | Convert source text into tokens. |
| `internal/bitlang/lexer_comment.go` | Recognize and skip bootstrap source comments. |
| `internal/bitlang/lexer_test.go` | Verify lexer behavior. |
| `internal/bitlang/preprocessed_token.go` | Represent strict tokens with explicit canonical identifier forms. |
| `internal/bitlang/preprocessed_property.go` | Represent one property state without hardcoded state names. |
| `internal/bitlang/property_spec.go` | Represent machine-readable property axes. |
| `internal/bitlang/property_spec_loader.go` | Decode and validate property specification data. |
| `internal/bitlang/property_spec_loader_test.go` | Verify property specification loading. |
| `internal/bitlang/property_validator.go` | Validate explicit properties against specification axes. |
| `internal/bitlang/property_validator_test.go` | Verify specification-driven property validation. |
| `internal/bitlang/spec_consistency.go` | Validate invariants across machine-readable specifications. |
| `internal/bitlang/spec_consistency_test.go` | Verify cross-specification consistency checks. |
| `internal/bitlang/declaration_spec.go` | Represent machine-readable declaration kinds and grammar. |
| `internal/bitlang/declaration_spec_loader.go` | Decode and validate declaration grammar data. |
| `internal/bitlang/declaration_spec_loader_test.go` | Verify declaration grammar loading and lookup. |
| `internal/bitlang/declaration_parser.go` | Parse one strict declaration and validate its explicit properties. |
| `internal/bitlang/declaration_parser_test.go` | Verify strict declaration parsing and property completeness. |
| `internal/bitlang/declaration_sequence_parser.go` | Split and parse multiple strict declarations of one target kind. |
| `internal/bitlang/declaration_sequence_parser_test.go` | Verify multi-declaration parsing and terminator errors. |
| `internal/bitlang/preprocessed_declaration.go` | Represent one strict declaration and its explicit properties. |
| `internal/bitlang/preprocessed_source.go` | Represent strict token-level preprocessor output. |
| `internal/bitlang/preprocessor.go` | Convert source into the initial preprocessed representation. |
| `internal/bitlang/preprocessor_test.go` | Verify the Source to Preprocessed boundary. |
| `internal/bitlang/canonical_name.go` | Represent and canonicalize Bitlang identifiers. |
| `internal/bitlang/canonical_name_test.go` | Verify identifier canonicalization. |
| `internal/bitlang/symbol.go` | Represent one named semantic value. |
| `internal/bitlang/duplicate_symbol_error.go` | Represent canonical-name collision errors. |
| `internal/bitlang/symbol_table.go` | Store and resolve symbols by canonical name. |
| `internal/bitlang/symbol_table_test.go` | Verify symbol-table behavior. |
| `internal/bitlang/artifact_kind.go` | Define canonical pipeline artifact kinds. |
| `internal/bitlang/artifact.go` | Represent one pipeline artifact. |
| `internal/bitlang/stage.go` | Represent and execute one pipeline stage. |
| `internal/bitlang/stage_transition.go` | Define valid canonical stage transitions. |
| `internal/bitlang/pipeline.go` | Order and execute stages. |
| `internal/bitlang/pipeline_test.go` | Verify pipeline behavior. |

## Machine-readable specifications

| File | Responsibility |
| --- | --- |
| `spec/preprocessed/properties.json` | Define property axes, states, applicability, and completeness requirements. |
| `spec/preprocessed/declarations.json` | Define strict declaration kinds, targets, terminators, and layout. |

## Go rule checker

| File | Responsibility |
| --- | --- |
| `tools/go-rule-checker/cmd/go-rule-checker/main.go` | Start the checker and map results to exit status. |
| `tools/go-rule-checker/run.go` | Coordinate one checker execution. |
| `tools/go-rule-checker/file_discovery.go` | Discover Go source files. |
| `tools/go-rule-checker/source_check.go` | Coordinate checks for one source file. |
| `tools/go-rule-checker/finding.go` | Represent and format one finding. |
| `tools/go-rule-checker/naming.go` | Check responsibility-obscuring names. |
| `tools/go-rule-checker/documentation.go` | Check exported documentation comments. |
| `tools/go-rule-checker/function_size.go` | Check function size limits. |
| `tools/go-rule-checker/file_size.go` | Check file size limits. |
| `tools/go-rule-checker/main_file.go` | Check startup-only main.go structure. |
| `tools/go-rule-checker/ignore.go` | Load and apply checker ignore rules. |
| `tools/go-rule-checker/generated.go` | Detect generated Go source. |
| `tools/go-rule-checker/import_check.go` | Check forbidden architecture-layer imports. |
| `tools/go-rule-checker/role_check.go` | Check Commander and Messenger dependency responsibilities. |
| `tools/go-rule-checker/responsibility_table.go` | Parse Go entries from the file responsibility table. |
| `tools/go-rule-checker/responsibility_check.go` | Compare Go files with registered responsibilities. |
| `tools/go-rule-checker/generated_test.go` | Verify generated-source detection. |
| `tools/go-rule-checker/import_check_test.go` | Verify architecture-layer classification. |
| `tools/go-rule-checker/role_check_test.go` | Verify architecture-role classification. |
| `tools/go-rule-checker/responsibility_check_test.go` | Verify responsibility path and scan-scope handling. |
| `tools/go-rule-checker/README.md` | Document checker usage. |
| `tools/go-rule-checker/IGNORE_FORMAT.md` | Document ignore configuration syntax. |

Other tests under `tools/go-rule-checker/*_test.go` verify the matching checker responsibility.

## Bitlang CI

| File | Responsibility |
| --- | --- |
| `tools/bitlang-ci/cmd/bitlang-ci/main.go` | Start the shared strict CI gate and map its exit status. |
| `tools/bitlang-ci/run.go` | Coordinate strict project-owned validation checks. |
| `tools/bitlang-ci/spec_check.go` | Validate machine-readable specifications during strict CI. |
| `tools/bitlang-ci/spec_check_test.go` | Verify strict CI specification validation. |
| `tools/bitlang-ci/README.md` | Document strict CI behavior and failure policy. |
| `.github/workflows/strict-ci.yml` | Run Bitlang CI and additional hosted checks on supported runner OSes. |
| `.gitattributes` | Keep source and validation text line endings deterministic across operating systems. |

## Project / AI routing documents

| File | Responsibility |
| --- | --- |
| `README.md` | Introduce Bitlang and its overall pipeline. |
| `CODING_RULES.md` | Define cross-language compiler implementation rules. |
| `GO_CODING_RULES.md` | Define Go-bootstrap-specific architecture rules. |
| `FILE_RESPONSIBILITIES.md` | Map files to responsibilities. |
| `AI_CONTEXT.md` | Provide the smallest AI development entrypoint. |
| `CURRENT_STATE.md` | Summarize current implementation capability and limits. |
| `CHANGE_ROUTING.md` | Route change categories to source/tests/docs. |
| `VALIDATION_ROUTING.md` | Route changes to the smallest sufficient validation. |
| `go.mod` | Define the Go bootstrap module and language version. |

## Planned core responsibilities

Keep these separate when implemented:

- `Application`: top-level application coordination
- `CuiFrontend`: terminal request/presentation handling
- `CompileRequest`: one requested conversion and options
- `CompilerPipeline`: complete official conversion orchestration
- `SourceArtifact`, `PreprocessedArtifact`, `CompiledArtifact`, `TreeObjectArtifact`, `VmAssemblyArtifact`: stage-specific representations
- `Preprocessor`: Source to Preprocessed
- `StaticAnalyzer`: static analysis over Preprocessed
- `Advisor`: non-semantic advisory diagnostics
- `Compiler`: Preprocessed to Compiled
- `TreeLowerer`: Compiled to TreeObject
- `VmAssemblyLowerer`: TreeObject to VM Assembly
- `Diagnostic`: UI-independent diagnostic representation
- `DiagnosticFormatter`: textual diagnostic presentation

## GUI rule

The GUI IDE and CUI are sibling frontends over the same official compiler core. GUI code must not duplicate compiler semantics.

Planned GUI responsibilities remain separate: `IdeApplication`, `IdeWindow`, `SourceEditor`, `StageNavigator`, `StageViewer`, `PipelineController`, `DiagnosticPanel`, `ArtifactDiffViewer`, `ProjectExplorer`, `BuildPanel`, `VmPanel`, `TranslatorPanel`.

## Split rule

If a responsibility cannot be described in one short sentence, split the file/component before extending it.
