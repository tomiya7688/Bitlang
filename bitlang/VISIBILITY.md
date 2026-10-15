# Bitlang Visibility Model

Status: current design

Bitlang visibility is designed to absorb access models from Bitlang-family and foreign languages without forcing them into only `public/private`.

## 1. Canonical scope-visibility states

Bitlang Explicit uses the following standard canonical states for the `scope_visibility` axis:

```text
Owner_visible
File_visible
Module_visible
Module_tree_visible
Project_visible
Global_visible
Custom_scope_visibility
```

These states describe the ordinary visibility boundary. They remain independent from inheritance visibility (`Protected / Unprotected`) and declaration export (`Exported / Unexported`).

### Owner_visible

Visible only from the declaration's ordinary owner scope.

Typical uses include class/type-private members and other declarations whose natural owner is a lexical/type scope.

### File_visible

Visible only from declarations originating in the same physical source file.

This is the canonical state produced by `Private module`.

### Module_visible

Visible from the same canonical module only.

Nested child modules are not automatically included unless the selected policy explicitly includes them.

### Module_tree_visible

Visible within the declaring module and its descendant module tree.

This is useful for APIs intended for one module hierarchy without exposing them to unrelated modules.

### Project_visible

Visible from the same Bitlang project.

This can represent common `internal` / crate-like visibility models when one project is the corresponding compilation/access unit.

### Global_visible

May participate in cross-project lookup/import, subject to export, dependency, import-capability, and other independent restrictions.

This is the canonical state normally produced by source `Public`.

### Custom_scope_visibility

Uses an explicit compile-time visibility expression when the standard boundaries cannot represent the source-language rule exactly.

A custom visibility expression is data, not executable code.

## 2. Source-facing shorthands

Bitlang source keeps familiar shorthands.

```text
Public
Private
```

They normalize before Bitlang Explicit.

`Public` normally becomes:

```text
Global_visible
```

`Private` is declaration-kind/context aware because the natural private boundary differs by construct.

At minimum:

```text
Private module
    -> File_visible

Private class/type member
    -> Owner_visible
```

Other declaration kinds must define their deterministic private expansion before the Explicit boundary.

The canonical Explicit representation must not require later stages to guess what `Private` meant.

Additional source conveniences may map directly to canonical boundaries, for example:

```text
File_visible
Module_visible
Module_tree_visible
Project_visible
```

A language adapter may also provide source-language spelling such as `internal`, `package`, `crate`, `assembly`, `friend`, or `pub(super)`, provided it lowers deterministically to a standard state or a custom expression.

## 3. Custom visibility expressions

`Custom_scope_visibility` carries a pure compile-time relation expression.

The expression may use boolean composition:

```text
any(...)
all(...)
not(...)
```

and scope-relation atoms such as:

```text
same_owner
same_file
same_module
within_module(<canonical-module>)
module_descendant_of(<canonical-module>)
same_project
within_project(<project-id>)
project_descendant_of(<project-id>)
subtype_of(<canonical-type>)
named_scope(<scope-id>)
named_symbol(<canonical-symbol>)
adapter_domain(<adapter-id>, <domain-id>)
```

The exact serialization is machine-readable and versioned.

Custom expressions must not execute arbitrary preprocessing/compiler code during access checks. They evaluate only already-resolved identities and relationships.

This allows adapters to preserve models such as:

- Java package-private through an adapter-defined package domain;
- C# `internal` through `Project_visible`;
- C# `protected internal` through a union of project visibility and subtype visibility;
- C# `private protected` through an intersection of project/owner and subtype relations;
- Rust `pub(crate)` through project visibility;
- Rust `pub(super)` / `pub(in path)` through named/ancestor module relations;
- C++-style friend access through named symbols/scopes;
- Swift `fileprivate` through `File_visible`;
- language-specific package/assembly/unit visibility through `adapter_domain`.

## 4. Visibility is not export

Visibility determines which source contexts are allowed to resolve/use a declaration.

`Exported / Unexported` determines whether an applicable declaration is exported from its module boundary.

Both may restrict an external access path.

For example:

```text
Global_visible + Unexported
    -> declaration is globally visible in principle,
       but not exported through that module boundary

Project_visible + Exported
    -> export does not widen visibility beyond the project
```

No axis grants a capability forbidden by another applicable axis.

## 5. Visibility is not inheritance access

`Protected / Unprotected` remains a separate inheritance-related axis for the existing model.

When a foreign language has visibility whose semantics are a union or intersection of lexical/project and subtype relations, the adapter may use `Custom_scope_visibility` so the exact relation is preserved instead of forcing an incorrect approximation.

A future normalization may reduce some `Protected` combinations into the general visibility-expression engine, but the axes are not silently merged.

## 6. Lookup hiding versus access denial

Visibility answers whether access is semantically permitted. Some source languages also distinguish whether an inaccessible declaration is hidden from name lookup or found and then rejected.

Bitlang must preserve that distinction when required by the source semantics.

The current native rule for a `Private module` is:

```text
Private module
    -> File_visible
    -> excluded from cross-file module discovery
    -> external import observes "module not found"
```

The general discoverability model is tracked separately from the visibility boundary so adapters can represent either hidden or discoverable-but-inaccessible behavior without weakening access control.

## 7. Compile-time only

Visibility is a compile-time semantic restriction.

After name resolution, import validation, static analysis, and required diagnostics have completed, backends do not need to carry source visibility into Bitlang Low, Assam Core, native assembly, or Bitlang VM runtime unless an explicit reflection/runtime metadata feature separately requires it.
