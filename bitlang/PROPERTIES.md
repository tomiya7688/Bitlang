# Bitlang Properties

Bitlang expresses complex declaration semantics through independent properties wherever practical. These properties may be inspected and transformed by preprocessor functions before canonical Bitlang preprocessed output is produced.

## Visibility properties

Visibility is separated into independent axes rather than represented by one broad modifier.

### Scope visibility

```text
Public
Private
```

`Public` means the declaration is visible from the enclosing visibility domain defined for that declaration kind.

`Private` means access is restricted to its defining scope or owning declaration according to the applicable scope rules.

### Inheritance visibility

```text
Protected
Unprotected
```

`Protected` grants the declaration the inheritance-related access defined by the language.

`Unprotected` means no such inheritance-specific access is granted.

This axis is independent from `Public` / `Private`.

### Module export visibility

```text
Exported
Unexported
```

`Exported` means the declaration is exposed outside its defining module through the module export boundary.

`Unexported` means it is not exposed outside that module.

Module export is independent from ordinary scope visibility. A declaration can therefore be public within a module while remaining unexported from that module.

## Composition

Visibility properties compose with the other independent Bitlang properties, including:

```text
Readable / Unreadable
Writeable / Unwriteable
Reassignable / Unreassignable
Owned / Borrowed
Initialized / Uninitialized
Nullable / Nonnullable
Optional / Required
```

For example, a local declaration may preprocess into a form conceptually similar to:

```text
Private Unprotected Unexported Readable Writeable Reassignable Initialized Nonnullable Required Int10x32 a = 4
```

The exact set of properties emitted depends on which property axes are meaningful for that declaration kind.
