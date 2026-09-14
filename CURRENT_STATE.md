# Bitlang Current State

Compact implementation snapshot for AI-assisted work. Update this file when capabilities materially change.

## Implemented
- Go bootstrap module and CUI entrypoint
- explicit staged pipeline skeleton
- source text representation
- lexer with identifier / number / string / character / symbol / EOF tokens
- source location tracking
- case-insensitive identifier canonicalization baseline
- symbol table with duplicate-name detection
- Go implementation coding rules
- file responsibility registry
- `go-rule-checker` with compact findings, size checks, and ignore support

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
- `names.go` and `pipeline.go` still contain multiple conceptual responsibilities and should be split

## Near-Term Implementation
1. split `names.go` by canonical name / symbol / symbol table responsibilities
2. split `pipeline.go` by artifact / stage / pipeline / transition responsibilities
3. introduce concrete Preprocessed artifact representation
4. add a minimal preprocessor boundary without inventing unspecified syntax
5. expose stage results through reusable compiler-core APIs before expanding CUI/GUI behavior

## Source of Truth Reminder
This file is only a current-state index. Language semantics belong in specification/design sources and compiler tests, not here.
