# Borrow-state preprocessing

This document defines how **Bitlang source and preprocessor rules** handle borrow state.

The canonical final borrow-state representation belongs to Bitlang Preprocessed:

- https://github.com/tomiya7688/Bitlang_preprocessed/blob/main/BORROW_STATE.ja.md

## Source-facing states

Bitlang source may explicitly state borrow condition using:

```text
Unborrowed
Shared_borrowed
Exclusive_borrowed
```

or omit it when preprocessing can resolve the applicable state from declarations, references, control flow, lifetime information, defaults, and explicit preprocessing rules.

Ownership and borrow state remain separate. `Owned / Borrowed` describes ownership responsibility; borrow state describes the current borrowing condition.

## Preprocessor rewrite

Preprocessor functions may inspect and explicitly rewrite borrow state.

A rewrite that weakens an active restriction, such as:

```text
Exclusive_borrowed -> Unborrowed
Shared_borrowed -> Unborrowed
```

is a semantic override rather than ordinary inference. It must come from an explicit preprocessing rule, configuration, or source-directed transformation.

The preprocessor may warn when it cannot prove that such an override is safe. If the resulting state is provably invalid, preprocessing must emit an error.

Typical diagnostics include:

- releasing or moving a resource while an active borrow appears to remain
- creating a conflicting borrow while the resource is `Exclusive_borrowed`
- allowing a borrow to outlive its source
- forcing a borrowed state to `Unborrowed` when the active borrow cannot be proven to have ended

## Output contract

Bitlang source does not own the canonical serialized representation. Before preprocessing completes, the final borrow state must be resolved; the exact required representation and consistency rules are defined by the Bitlang Preprocessed specification.
