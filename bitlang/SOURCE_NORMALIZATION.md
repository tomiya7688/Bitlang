# Bitlang Source Normalization

Bitlang source is intentionally more permissive and concise than Bitlang Explicit, but they share the same underlying Bitlang semantic property system.

Bitlang source is allowed to optimize for human readability by omitting information or using conveniences. The preprocessor converts those forms into the fully explicit normalized Bitlang form defined by `tomiya7688/Bitlang-Explicit`.

## Core rule

```text
human-friendly Bitlang
    -> preprocessing / normalization
    -> fully explicit Bitlang Explicit
```

A Bitlang source construct is acceptable when its semantics can be resolved deterministically before the Explicit boundary.

Bitlang Explicit must not inherit unresolved source shorthand, missing semantic axes, or context-dependent aliases.

## Forms that may be normalized

Bitlang source may use:

- omitted properties whose values can be resolved mechanically;
- partially specified property sets;
- aliases for properties, types, names, or declaration patterns;
- property bundles that expand into several canonical properties;
- declaration forms whose kind implies semantic properties;
- type forms that imply semantic properties;
- module/project defaults;
- preprocessing rules that inspect context and fill omitted information;
- source sugar that expands into more explicit statements or declarations;
- source-only names or short references that resolve into canonical qualified references.

All such forms are resolved before Bitlang Explicit is produced.

## Property normalization

The canonical property axes remain independent even when Bitlang source uses a compact description.

For example, a source-facing convenience may conceptually expand as:

```text
Mutable_value
    -> Readable Writeable Reassignable
```

while another profile may conceptually expand as:

```text
Owned_resource
    -> Owned Unborrowed Movable Unmoved Releasable Unreleased
```

These names are examples of the normalization mechanism, not a declaration that those exact aliases are mandatory built-ins.

After expansion, any still-missing axes are resolved independently.

A compact source form therefore does not weaken the Explicit model. It only reduces what the programmer must write manually.

## Explicit detail is always allowed

Bitlang does not force users to use the relaxed form.

**Every canonical property available in Bitlang Explicit is also available to Bitlang source when applicable.** Preprocessing does not own a hidden property vocabulary that source authors are forbidden to write.

A programmer, generator, language adapter, or debugging tool may therefore write a fully explicit Bitlang declaration containing every applicable property. The preprocessor preserves valid explicit requirements and fills only what remains unresolved.

Consequently, fully explicit Bitlang source may already look almost identical to Bitlang Explicit. The difference is that source is allowed to omit or abbreviate information, while the Explicit boundary is not.

This makes Bitlang usable both as a human-facing language and as a convenient transformation target for other languages.

## Namespace mount normalization

Preprocessing may assign outer namespace/module ownership from project-relative paths.

Canonical preprocessing operations:

```text
mount_namespace_file(file_path, namespace_path)
mount_namespace_tree(directory_path, namespace_path)
set_parent_project(parent_project)
```

An exact file mount applies only to that file. A tree mount applies recursively to a directory and automatically extends the namespace using relative child-directory segments.

Filename text is not automatically added as a namespace segment. Explicit module declarations inside a mounted file are nested under the mounted outer namespace.

Resolution order is exact file mount, then most-specific tree mount, then inherited ancestor tree mount. Ambiguous equally specific mounts are errors.

Project inheritance is single-parent and deterministic. Child projects inherit applicable project defaults, namespace mounts, adapter configuration, and explicitly inheritable preprocessing configuration. Child configuration overrides parent configuration.

Parent-project cycles and multiple direct parents are errors. Parent project source files are not automatically included merely because configuration is inherited.

Namespace/project preprocessing state is fully consumed before Bitlang Explicit.

## Scoped property defaults

The preprocessor may define default property values for declarations by structural scope. This exists especially to make language adapters and generated Bitlang concise while preserving a fully explicit Bitlang Explicit result.

Supported default scopes include:

- file scope;
- class/type scope;
- namespace/module scope;
- project scope;
- language-adapter scope.

A native Bitlang `module` is the normal Bitlang naming scope. A foreign-language adapter may expose its source-language namespace as a preprocessing namespace scope and map that scope into the corresponding Bitlang module/name hierarchy during normalization.

A scoped default may set one property axis or a reusable set of property axes, and may optionally restrict the declaration kinds it applies to.

Defaults fill unresolved axes only. They do not silently overwrite a directly written canonical property or an already stronger explicit preprocessing requirement.

When multiple defaults apply to the same unresolved axis, precedence is:

```text
explicit declaration property
    > declaration-targeted preprocessing rule
    > innermost class/type default
    > file default
    > innermost namespace/module default
    > current project default
    > parent-project defaults from nearest to farthest
    > language-adapter default
    > Bitlang built-in default
```

Within nested class/type or namespace/module scopes, the innermost matching scope wins.

If two rules at the same precedence level assign contradictory defaults to the same target and axis, preprocessing must report an error unless an explicit ordering/override relation between those rules has been defined.

A scoped default is preprocessing state only. The default declaration itself disappears before Bitlang Explicit; each affected declaration receives the resolved canonical property value.

Conceptually:

```text
file default:       Dynamic Process_retention
namespace default:  Public
class default:      Unreassignable
```

may allow many declarations to omit those properties in source, while Bitlang Explicit still contains each final property explicitly.

### Defaults never weaken safety

Property defaults affect only how omitted semantic axes are resolved. They do not weaken Bitlang's safety requirements.

After all applicable defaults have been resolved, the resulting program is validated exactly as if every resolved property had been written explicitly on each declaration.

If the final resolved property set, operation, control flow, ownership/lifetime relation, release/finalization behavior, reference usage, or another semantic condition is provably unsafe or invalid, preprocessing must report an error.

A file/class/namespace/module/project/language-adapter default cannot suppress, downgrade, or bypass that error merely because it supplied the property value.

A later static-analysis/compiler stage must also reject a provably unsafe final program even if the unsafe state originated entirely from defaults or generated transformation rules.

## Conflict handling

Relaxed syntax must not mean ambiguous semantics.

- compatible shorthand and explicit properties are merged;
- defaults fill only unresolved axes;
- contradictory explicit requirements are errors unless an explicit preprocessing override rule defines the change;
- dangerous but not provably invalid semantic overrides may generate warnings;
- provably invalid final states are errors;
- normalization must produce a deterministic final result.

## Canonical boundary

At the Bitlang source -> fully explicit Bitlang Explicit boundary:

- all applicable semantic property axes are resolved;
- all source aliases are expanded;
- all source-only property bundles are decomposed;
- all required names and references are canonicalized;
- all source sugar relevant to later semantic analysis is expanded;
- no later compiler stage is required to guess what omitted source notation meant.

The canonical result is governed by the Bitlang Explicit repository, especially:

- `LANGUAGE_SPEC.ja.md`
- `PROPERTIES.ja.md`
- `BORROW_STATE.ja.md`



## Legacy Explicit-shaped input

Older Bitlang Explicit representations may be accepted by preprocessing as migration input when their meaning can still be reconstructed deterministically.

A common version change is that a semantic property which was previously implicit becomes public and mandatory in the newer Bitlang Explicit specification. In that case, the older representation is no longer considered canonical current Explicit, because one or more semantic axes are still implicit.

Conceptually, such input is treated as **Bitlang source written in a Explicit-shaped syntax**:

```text
old Explicit-shaped input
    -> preprocess / infer formerly implicit properties
    -> materialize newly explicit properties
    -> current fully explicit Bitlang Explicit
```

The fact that the text resembles an older Explicit format does not grant it direct access to the compiler. The current Explicit boundary is defined by the current canonical specification, not by historical syntax.

This keeps backward migration in the preprocessor and allows the Bitlang compiler to remain strict: it only needs to accept the current fully explicit form.

## Design intent

Bitlang should be strict about **meaning**, not unnecessarily strict about **how much repetitive semantic information a human must type**.

The strictness of mandatory explicitness belongs at the normalized boundary. The source form may remain comparatively comfortable, but it retains access to the complete Bitlang property vocabulary at all times.
