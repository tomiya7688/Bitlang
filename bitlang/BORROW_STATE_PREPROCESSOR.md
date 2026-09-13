# Borrow-state preprocessing

Borrow state is represented by the properties:

```text
Unborrowed
Shared_borrowed
Exclusive_borrowed
```

Preprocessor functions may inspect and explicitly rewrite these states.

A rewrite that weakens an active restriction, such as `Exclusive_borrowed -> Unborrowed` or `Shared_borrowed -> Unborrowed`, is a semantic override rather than ordinary inference and must come from an explicit preprocessing rule, configuration, or source-directed transformation.

The preprocessor may warn when it cannot prove that such an override is safe. If the resulting state is provably invalid, preprocessing must emit an error.

Typical diagnostics include:

- releasing or moving a resource while an active borrow appears to remain
- creating a conflicting borrow while the resource is `Exclusive_borrowed`
- allowing a borrow to outlive its source
- forcing a borrowed state to `Unborrowed` when the active borrow cannot be proven to have ended

The resolved final borrow state must be emitted explicitly in Bitlang preprocessed.
