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

## Capability views instead of defensive copies

Bitlang access control is designed so code does not need to create a copied structure or a separate reduced-value object merely to express restricted access.

The same underlying declaration/value may be exposed through different capability-restricted access paths.

For example, one module may receive a structure through a path that is:

```text
Readable + Unwriteable
```

while another path may be:

```text
Type_only
```

and the owning module may retain broader capabilities.

These access paths do not imply separate runtime copies. They are semantic capability views enforced by preprocessing/static analysis/compiler validation.

A compiler may erase capability metadata after all relevant accesses are validated and lowered, provided doing so cannot reintroduce an operation that the restricted view forbids.

This mechanism allows fine-grained access control over structure knowledge, value observation, mutation, reassignment, calls, ownership transfer, and other independent capabilities without requiring wrapper objects or defensive copies solely for permission management.


## Module discovery and Private modules

Import capability processing begins only after module discovery has found an eligible module.

A source module marked `Private` normalizes to canonical `File_visible`. When the importing source file is different from the module's defining physical source file:

- the private module is excluded from import lookup;
- the import must not create a binding to it;
- a fully qualified module reference must not bypass the same exclusion;
- import capabilities such as `Callable`, `Type_only`, or write/read restrictions cannot grant access to it.

Conceptually:

```text
cross-file module discovery
    -> remove modules whose canonical visibility does not admit the importer
    -> resolve Public module candidate
    -> apply export/member visibility
    -> apply import capabilities
```

This ordering prevents import-capability syntax from turning a file-private module into an externally visible one.

Within the same physical source file, a Private module remains resolvable under the ordinary module/name rules.

If an external import names only a Private module, the observable import result is "module not found" rather than a successful lookup followed by an access-capability grant.

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

## Type knowledge does not grant value inspection

Knowledge of a type or structure never grants permission to inspect a value of that type.

Type metadata capability and value-read capability are independent.

For example, importing the structure of a type may allow code to know:

- which visible fields exist;
- the declared type of each visible field;
- visible nested type relationships;
- signatures and other permitted type metadata.

However, if a structure/class variable, instance, member, or field is `Unreadable` through the effective access path, code must not read its current value merely because the corresponding type structure is known.

This distinction also applies to validation.

A check that only asks whether a declaration's **declared type** matches a known structural/type contract may be performed from type metadata alone.

A check that asks whether the **current stored value** satisfies some structural, content, range, invariant, serialization, or runtime-value condition requires the corresponding value-read capability.

Therefore:

```text
Type_only + Unreadable value
    -> type/structure inspection allowed
    -> current value inspection forbidden
    -> direct value-based validation forbidden
```

The restriction cannot be bypassed by reflection, serialization helpers, generated code, preprocessor-generated access paths, aliases, references, or another indirect mechanism. Any operation that semantically observes the protected value requires effective `Readable` capability.

A trusted function or module that independently has the necessary read capability may perform such validation and expose only its permitted result. That does not grant the caller direct access to the protected value.

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


## Generalized scope visibility

Module discovery must evaluate the canonical visibility policy from `bitlang/VISIBILITY.md`, not hard-code only Public/Private.

Standard boundaries include owner, file, module, module-tree, project, and global visibility. `Custom_scope_visibility` may carry a pure scope-relation expression for language-adapter semantics such as package, assembly, friend, or path-relative visibility.

Import capabilities remain capability-reducing only and are evaluated after the declaration/module is eligible under its visibility policy.


## Call-site policy and Callable imports

`Callable` on an import path only means that the import capability does not itself prohibit invocation.

It cannot override the target function's canonical call-site policy.

A call through an imported function value is valid only when visibility/export/import capability **and** the target call-site policy all permit the direct caller.

See `bitlang/CALL_SITE_POLICY.md`.
