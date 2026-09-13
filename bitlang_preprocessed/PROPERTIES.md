# Bitlang Preprocessed Properties

Bitlang preprocessed must emit the resolved canonical state of every property axis that is meaningful for a declaration. Later compiler stages should not need to infer omitted defaults.

## Visibility properties

The canonical visibility properties are separated into independent axes.

### Scope visibility

```text
Public
Private
```

### Inheritance visibility

```text
Protected
Unprotected
```

### Module export visibility

```text
Exported
Unexported
```

These axes are independent. For example, a declaration may be `Public` inside its module while also being `Unexported` from the module boundary.

When an axis is meaningful for the declaration, Bitlang preprocessed must state one resolved side of that axis explicitly rather than relying on a default.

Visibility properties compose with all other applicable canonical properties, including readability, writability, reassignment, ownership, initialization, nullability, and optionality.

A typical local initialized integer may therefore have a canonical declaration conceptually similar to:

```text
Private Unprotected Unexported Readable Writeable Reassignable Initialized Nonnullable Required Int10x32 a = 4
```

Properties that are genuinely inapplicable to a declaration kind are not emitted merely as placeholders.
