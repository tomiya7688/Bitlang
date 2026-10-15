# Bitlang Current State

English | [日本語](CURRENT_STATE.ja.md)

Compact implementation snapshot for AI-assisted and contributor work.

## Implemented
- Go bootstrap module and CUI entrypoint
- explicit staged pipeline skeleton
- source text and lexical token representations
- lexer and one-pass lexical preprocessing with source locations and identifier canonical forms
- loss-preserving C#-style class/member parser for modifiers, generic syntax, accessors, and opaque initializer/body tokens
- property-axis specifications with explicit state resolution and validation
- spec-driven declaration parsing and source-level property expansion hooks
- case-insensitive identifier canonicalization baseline
- symbol table with duplicate-name detection
- pipeline responsibilities split into artifact / stage / transition / orchestration files
- name responsibilities split into canonical name / symbol / duplicate error / symbol table files
- Go implementation coding rules and responsibility registry
- `go-rule-checker` with compact findings, file/function size checks, and ignore support
- project-owned `bitlang-ci` validation gate
- AI context entrypoint, change routing, and validation routing

## Current Pipeline
- Source: source text, lexer, and initial class/member syntax parsing are implemented
- Explicit: strict declaration-only input can be parsed into a typed `ExplicitArtifact`; required properties are validated and original source text is retained. The data-driven source-default helper covers `retention` (`Dynamic`), `retention_domain` (`Process_retention`), `instance_access` (`Instance_required`), and `finalization_trigger` when `lifetime` is resolved; remaining defaults, legacy migration, and class/member normalization are not yet integrated
- Compiled: stage slot only
- TreeObject: stage slot only
- VM Assembly: stage slot only
- VM / translators: planned

## Known Temporary Limits
- identifier canonicalization currently uses Go lowercase behavior as a bootstrap baseline
- lexer identifier character rules are temporary ASCII rules
- `LexicalPreprocess` performs lexical processing only and retains the original text; `ApplySourcePropertyDefaults` resolves configured unconditional and chained conditional defaults by dependency, leaving unresolved conditions untouched, but is not yet connected to preprocessing. Declaration-location-dependent defaults, source-sugar expansion, legacy migration, full property completion, and whole-pipeline iterative normalization remain unimplemented
- class bodies retain balanced tokens, but expressions and statements do not yet have semantic parsers
- generic constraints and additional declaration forms remain outside the initial source grammar
- grammar-specific multi-character operators are not yet finalized

## Near-Term Implementation
1. define remaining source defaults with their declaration-context inputs, then connect data-driven completion to normalization before deterministic legacy migration
2. implement iterative normalization with deterministic convergence and cycle checks
3. connect parsed source declarations to that normalization boundary without guessing adapter-specific modifier meanings
4. expose complete stage results through reusable compiler-core APIs before expanding the CUI

## Source of Truth Reminder
This file is an index only. Language semantics belong in specification/design sources and compiler tests.
