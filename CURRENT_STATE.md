# Bitlang Current State

Compact implementation snapshot for AI-assisted work.

## Implemented
- Go bootstrap module and CUI entrypoint
- explicit staged pipeline skeleton
- source text and lexical token representations
- lexer with source locations
- case-insensitive identifier canonicalization baseline
- symbol table with duplicate-name detection
- pipeline responsibilities split into artifact / stage / transition / orchestration files
- name responsibilities split into canonical name / symbol / duplicate error / symbol table files
- Go implementation coding rules and responsibility registry
- `go-rule-checker` with compact findings, file/function size checks, and ignore support
- AI context entrypoint, change routing, and validation routing

## Current Pipeline
- Source: partial
- Preprocessed: representation/processing not yet implemented
- Compiled: stage slot only
- TreeObject: stage slot only
- VM Assembly: stage slot only
- VM / translators: planned

## Known Temporary Limits
- identifier canonicalization currently uses Go lowercase behavior as a bootstrap baseline
- lexer identifier character rules are temporary ASCII rules
- comments and grammar-specific multi-character operators are not yet finalized

## Near-Term Implementation
1. introduce concrete Preprocessed artifact representation
2. add a minimal preprocessor boundary without inventing unspecified syntax
3. expose stage results through reusable compiler-core APIs
4. expand CUI only after compiler-core stage APIs are stable

## Source of Truth Reminder
This file is an index only. Language semantics belong in specification/design sources and compiler tests.
