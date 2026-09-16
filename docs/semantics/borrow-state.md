# Borrow state in Bitlang source

This document defines the **source-facing Bitlang** rules for borrow state. Canonical Bitlang Preprocessed representation is specified in the separate `tomiya7688/Bitlang_preprocessed` repository, and low-level lowering is specified in `tomiya7688/Bitlang_compiled`.

## Source properties

Bitlang source may use the following borrow-state properties:

```text
Unborrowed
Shared_borrowed
Exclusive_borrowed
```

Borrow state is independent from ownership. `Owned` / `Borrowed` describes ownership responsibility; the properties above describe the current borrowing condition of a declaration or resource.

Source code may state the borrow state explicitly or omit it. When omitted, the Bitlang preprocessor resolves the applicable state from declarations, references, control flow, lifetime information, defaults, and explicit preprocessing rules.

## Source-level meaning

- `Unborrowed`: no active borrow currently restricts normal access through the owning declaration.
- `Shared_borrowed`: one or more shared borrows are active. Shared borrows may coexist where the access rules allow it, while operations that would invalidate those borrows are restricted.
- `Exclusive_borrowed`: an exclusive borrow is active. Conflicting borrows and conflicting direct access are prohibited while it remains active.

Typical invalid or suspicious source states include releasing or moving a resource while a borrow remains active, creating a conflicting borrow while `Exclusive_borrowed`, and allowing a borrow to outlive its source.

## Preprocessor overrides

Preprocessor functions may inspect and explicitly change borrow-state properties.

A rewrite that weakens restrictions, such as:

```text
Exclusive_borrowed -> Unborrowed
```

is a semantic override, not ordinary inference. It must therefore come from an explicit preprocessing rule, configuration, or source-directed transformation.

If an override appears dangerous but invalidity cannot be proven, preprocessing may emit a warning. If the resulting program is provably invalid, preprocessing must emit an error.

## Stage boundary

Bitlang source does not define the final canonical serialized property set. The resolved state passed to the next stage is governed by the Bitlang Preprocessed specification.
