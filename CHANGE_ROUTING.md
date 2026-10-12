# Bitlang Change Routing

Use this map to choose the smallest source/test/doc working set before reading broadly.

| Change | Read / edit first | Matching validation |
| --- | --- | --- |
| source grammar | `bitlang/BASIC_SYNTAX.md`, source parser | `internal/bitlang/source_class_parser_test.go` |
| source modifier preservation | source class representations + source parser | `internal/bitlang/source_modifier_parser_test.go` |
| lexical rules | `internal/bitlang/lexer.go`, token files | `internal/bitlang/lexer_test.go` |
| identifier semantics | canonical-name source, symbol source | name/symbol tests |
| property name matching | property match / resolver / validator + specification consistency | property and declaration parser tests |
| symbol storage | symbol table source | symbol table tests |
| pipeline stages/order | artifact, stage, pipeline, transition source | pipeline tests |
| family pipeline contract | `spec/family-pipeline.json`, neighboring repository contracts | local manifest/stage consistency; cross-repo fixture execution and remote CI inspection remain unimplemented |
| Source -> Explicit | preprocessor + Explicit artifact | lexical preprocessor tests + pipeline tests |
| Go architecture/rules | `GO_CODING_RULES.md`, `FILE_RESPONSIBILITIES.md` | `tools/go-rule-checker/` |
| rule checker | `tools/go-rule-checker/` | checker package tests |
| CUI behavior | `cmd/bitlang/` + compiler entry API | targeted tests + build |
| GUI IDE planning | official compiler API + GUI responsibility section | compiler tests first |
| VM assembly/lowering | TreeObject / VM Assembly components | stage-specific tests |
| VM / translator | VM or translator component only | targeted execution/smoke |

## Routing Rules
- Search first, then open only matching files.
- Prefer matching tests before unrelated implementation files.
- Read detailed coding rules only when a policy question is involved.
- Do not inspect later compiler stages unless the changed stage contract affects them.
- Shared representation or public contract changes require broader pipeline validation.
