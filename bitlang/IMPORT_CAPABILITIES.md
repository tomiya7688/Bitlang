# Bitlang Import Capabilities

## Core rule

An import is not only name visibility. It may restrict what the importing code is allowed to do through the imported binding.

Import capabilities do not rewrite the original declaration. They apply to the access path created by the import.

An import may make access more restrictive, but it must never grant a capability that the exported declaration itself does not permit.

Conceptually:

```text
original declaration capability
    AND export boundary capability
    AND import binding capability
    -> effective imported capability
```

The most restrictive applicable rule wins.

## Type-only imports

Bitlang supports imports that expose a type declaration without importing a runtime value binding.

Conceptually:

```text
Type_only import
```

A type-only import may expose the source-visible type structure needed for compilation, such as:

- the type name;
- generic parameters after applicable preprocessing;
- visible field/member names;
- visible field/member types;
- visible nested type relationships;
- visible function/method signatures needed as type information.

It does not grant access to a runtime instance, variable storage, mutable state, or function invocation merely because those declarations are described by the imported type.

Physical ABI/layout information such as exact offsets, packing, padding, or backend storage representation is not implied by ordinary type-only import. Such information requires an explicit layout/ABI contract where applicable.

This is useful for structures and interfaces whose shape must be known while their values must remain inaccessible.

## Struct-oriented examples

A structure may be imported with different capability levels depending on the intended use.

```text
Type_only
    -> know the structure/type contract only

Readable + Unwriteable
    -> read permitted values through the imported binding
    -> writing through that binding is forbidden

Readable + Writeable
    -> read and write where the underlying declaration/export also permits it
```

Field-level visibility and field-level capabilities still apply. Importing a structure as `Readable` does not bypass a field that is itself unreadable or inaccessible.

Likewise, `Writeable` on the import path cannot grant writing to a field that is `Unwriteable` at its declaration or export boundary.

## Variable write capability

Variable imports use the existing Bitlang write-capability axis:

```text
Writeable
Unwriteable
```

Applied to an import binding:

- `Writeable`: the import does not itself forbid writing through that imported path.
- `Unwriteable`: writing through that imported path is forbidden even when the original variable is writeable.

A `Writeable` import cannot override an `Unwriteable` source/export declaration.

This import capability controls the ordinary Bitlang write axis. Reassignment remains a separate `Reassignable / Unreassignable` semantic axis and may receive its own import restriction where required.

## Function call capability

Function imports use a call-capability axis:

```text
Callable
Uncallable
```

Applied to an import binding:

- `Callable`: the import does not itself forbid invocation of the imported function.
- `Uncallable`: the imported function may not be invoked through that import capability.

`Uncallable` is not bypassed by first copying or passing the imported function as a function value. Any derived callable reference/value must preserve the call restriction unless an explicit, valid capability transformation is defined by the language.

This allows a function symbol to be imported for compile-time inspection, signature/type reference, metadata use, or other non-call purposes without granting runtime invocation authority.

## Independence

Variable write capability and function call capability are independent.

An import mechanism may therefore expose only the capabilities relevant to the imported declaration kind.

Other import capabilities such as read access, reassignment, ownership transfer, or future execution capabilities remain separate axes rather than being implicitly bundled into `Writeable` or `Callable`.

## Preprocessing and Bitlang Explicit

Source-facing import syntax may use shorthand, grouped imports, module defaults, or language-adapter rules.

Before Bitlang Explicit is produced, preprocessing must resolve the effective imported capabilities deterministically.

The final normalized program must not rely on an implicit permission escalation.

If an import configuration requests a capability that the source/export declaration does not allow, that request is invalid and preprocessing reports an error rather than silently granting it.

Exact source syntax and source-level defaults are defined separately.
