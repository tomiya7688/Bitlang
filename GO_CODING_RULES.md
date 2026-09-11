# Bitlang Go Implementation Coding Rules

This document applies only to the current Go bootstrap implementation of Bitlang.

It does **not** define the Bitlang language itself, and it does **not** impose style or architectural rules on programs written in Bitlang.

Bitlang source code is intentionally allowed to be written in many styles. The language's preprocessing and lowering stages are responsible for normalizing those forms into the same strict semantic representation where they mean the same thing. A Bitlang user should not need to manually imitate the internal structure of the official compiler.

These rules exist only because the current official compiler is implemented in Go, where maintainability, dependency direction, and responsibility boundaries must be enforced manually.

The Go implementation adopts the principles of:

```text
https://github.com/tomiya7688/upd-commander-base-design
```

The purpose is to keep application orchestration, cross-layer communication, compiler processing, and data access clearly separated while the bootstrap compiler is maintained in Go.

## 1. UI / Process / Data separation

The Go application is conceptually separated into three layers:

```text
UI Layer
Process Layer
Data Layer
```

For this project:

```text
UI Layer
  CUI commands
  future GUI IDE frontend
  presentation and user interaction

Process Layer
  official Bitlang compiler pipeline
  preprocessing
  lexing/parsing
  static analysis
  lowering
  VM / translator orchestration

Data Layer
  source-file access
  project-file access
  persistence
  cache
  external artifact storage
```

These are architectural responsibilities of the Go implementation only. They are not Bitlang language constructs and are not required of Bitlang programs.

## 2. Commander is routing only

A Commander decides **what operation should be called next**.

A Commander MAY:

- receive a request
- select the appropriate Processing component
- pass inputs to Processing
- pass results to the next Processing component
- request cross-layer communication through a Messenger

A Commander MUST NOT:

- perform compiler transformations itself
- parse tokens
- analyze syntax
- generate IR
- format UI output
- read or write files directly
- contain complex semantic conditions

If real processing appears in a Commander, move it into a Processing component.

The rule is:

```text
Which operation to call -> Commander
How the operation works -> Processing
```

## 3. Messenger is communication only

A Messenger exists only for communication between layers.

A Messenger MAY:

- send a command or data to an adjacent layer
- receive a command or data from an adjacent layer
- pass received content to the local Commander

A Messenger MUST NOT:

- implement compiler semantics
- choose which Processing component should handle a request
- modify compiler artifacts for convenience
- perform UI presentation decisions
- perform file transformation or persistence logic

Messenger code should remain thin.

## 4. Processing owns real work

Actual work belongs in Processing components.

Examples:

```text
UI Processing
  format CUI output
  update IDE presentation state

Process Processing
  lex source
  preprocess source
  parse tokens
  analyze symbols
  lower IR
  translate VM assembly

Data Processing
  load source file
  save project settings
  read cached artifacts
```

Each Go Processing component follows the Go implementation maintenance rules:

- one file = one responsibility
- one function = one operation

These are Go bootstrap implementation rules, not restrictions on Bitlang source programs.

## 5. Do not skip layers

The normal dependency path in the Go application is:

```text
UI <-> Process <-> Data
```

Direct UI-to-Data access is prohibited in normal application flow.

Examples:

```text
UI Processing -> Data Processing       # prohibited
UI Commander  -> Data Processing       # prohibited
Process Processing -> Data Processing  # prohibited direct cross-layer call
```

When Process needs Data, it must use the Process Messenger to communicate with the Data layer.

Likewise, UI must not call compiler Processing modules through ad-hoc direct dependencies that bypass the defined Process entry path.

## 6. Official compiler remains the Process backend

The CUI compiler and future GUI IDE are frontends.

They MUST use the same official compiler Process layer.

```text
                 Official compiler
                  Process layer
                 /             \
              CUI              GUI IDE
```

The GUI IDE MUST NOT implement a second parser, second preprocessor, simplified semantic engine, or GUI-only compiler pipeline.

All displayed compiler artifacts and diagnostics must originate from the official compiler backend.

## 7. Commanders and Messengers must stay small

Commander and Messenger files should normally be significantly smaller than Processing files.

If a Commander or Messenger responsibility record becomes large, that is a design warning.

Typical causes include:

- real processing has leaked into the Commander
- formatting has leaked into the Messenger
- several commands have been merged into one responsibility
- unrelated communication paths have been combined

Split or move behavior before adding more logic.

## 8. Closed Processing modules are preferred

A Processing module called by a Commander should normally be closed to unrelated parts of the application.

Prefer:

```text
Commander -> Processing
```

Avoid arbitrary Processing-to-Processing access across unrelated responsibilities.

A Processing module should not become a hidden secondary Commander.

If several Processing operations must be coordinated, that coordination belongs in a Commander.

## 9. Main is startup only

`main.go` is responsible only for application startup.

It may:

- initialize top-level dependencies
- create the initial Commander or application entry component
- hand over control
- convert the final application result into a process exit code

It MUST NOT contain compiler logic, command implementation, file processing, or UI formatting logic.

## 10. File responsibility table remains authoritative

`FILE_RESPONSIBILITIES.md` is the architectural responsibility registry for the implementation repository.

Every Go implementation file must have exactly one responsibility recorded there.

Commander, Messenger, and Processing files must be identifiable from their responsibility even if their exact filenames do not contain those words.

If a responsibility description becomes long enough to require several independent clauses, split the file/component.

Again, this table describes the compiler implementation. It is not a source-organization requirement for Bitlang programs.

## 11. Go-specific implementation rules

In addition to the compiler implementation guidance in `CODING_RULES.md`:

- use `gofmt`
- keep `go test ./...` passing
- keep `go build ./cmd/bitlang` passing
- avoid reflection in compiler semantics
- avoid goroutines/channels inside deterministic compiler passes unless concurrency is isolated outside semantic behavior
- prefer explicit error returns over panic for user/compiler input errors
- keep OS-specific logic in outer/Data-side adapters
- keep Process logic independently testable without CUI or GUI
- do not make Go-specific behavior part of Bitlang semantics

## 12. Review checklist

Before adding a Go component, verify:

1. Which layer owns it: UI, Process, or Data?
2. Is it Commander, Messenger, or Processing?
3. Does it have exactly one responsibility?
4. Does every function perform exactly one operation?
5. Is any real processing accidentally inside a Commander or Messenger?
6. Does it bypass an adjacent-layer Messenger?
7. Can the Process logic be tested without UI and filesystem dependencies?
8. Is the file responsibility recorded in `FILE_RESPONSIBILITIES.md`?

If these questions cannot be answered clearly, the Go component boundary should be redesigned before implementation.

## 13. Bitlang source is deliberately freer than the Go bootstrap

The strictness in this document is an implementation maintenance constraint caused by the current host language and project structure.

Bitlang itself is intended to absorb different source-writing styles during preprocessing and lower them into the same defined behavior where semantics are equivalent.

Therefore, do not copy these Go architectural restrictions into the Bitlang language specification merely because the bootstrap compiler uses them.
