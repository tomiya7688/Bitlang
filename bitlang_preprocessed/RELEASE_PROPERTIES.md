# Bitlang Preprocessed Release Properties

When release semantics apply, Bitlang preprocessed must explicitly state both release behavior and release capability.

## Canonical properties

```text
Auto_release
Manual_release
Releasable
Unreleasable
```

`Auto_release` / `Manual_release` form the release-behavior axis.

`Releasable` / `Unreleasable` form the release-capability axis.

These axes are independent from ownership and lifetime, but legal combinations are constrained by them.

Examples:

```text
Owned Manual_release Releasable Ptr<My_type>
Borrowed Unreleasable Ref<My_type>
```

A borrowed declaration does not own release responsibility and normally resolves to `Unreleasable`.

If source Bitlang omits these properties, preprocessing must resolve and emit them explicitly before later compilation stages.

Later stages must not infer release behavior or release permission from omission.
