# Bitlang Source Normalization

Bitlang source is intentionally more permissive and concise than Bitlang Preprocessed, but they share the same underlying Bitlang semantic property system.

Bitlang source is allowed to optimize for human readability by omitting information or using conveniences. The preprocessor converts those forms into the fully explicit normalized Bitlang form defined by `tomiya7688/Bitlang_preprocessed`.

## Core rule

```text
human-friendly Bitlang
    -> preprocessing / normalization
    -> fully explicit Bitlang Preprocessed
```

A Bitlang source construct is acceptable when its semantics can be resolved deterministically before the Preprocessed boundary.

Bitlang Preprocessed must not inherit unresolved source shorthand, missing semantic axes, or context-dependent aliases.

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

All such forms are resolved before Bitlang Preprocessed is produced.

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

A compact source form therefore does not weaken the Preprocessed model. It only reduces what the programmer must write manually.

## Explicit detail is always allowed

Bitlang does not force users to use the relaxed form.

**Every canonical property available in Bitlang Preprocessed is also available to Bitlang source when applicable.** Preprocessing does not own a hidden property vocabulary that source authors are forbidden to write.

A programmer, generator, language adapter, or debugging tool may therefore write a fully explicit Bitlang declaration containing every applicable property. The preprocessor preserves valid explicit requirements and fills only what remains unresolved.

Consequently, fully explicit Bitlang source may already look almost identical to Bitlang Preprocessed. The difference is that source is allowed to omit or abbreviate information, while the Preprocessed boundary is not.

This makes Bitlang usable both as a human-facing language and as a convenient transformation target for other languages.

## Conflict handling

Relaxed syntax must not mean ambiguous semantics.

- compatible shorthand and explicit properties are merged;
- defaults fill only unresolved axes;
- contradictory explicit requirements are errors unless an explicit preprocessing override rule defines the change;
- dangerous but not provably invalid semantic overrides may generate warnings;
- provably invalid final states are errors;
- normalization must produce a deterministic final result.

## Canonical boundary

At the Bitlang source -> fully explicit Bitlang Preprocessed boundary:

- all applicable semantic property axes are resolved;
- all source aliases are expanded;
- all source-only property bundles are decomposed;
- all required names and references are canonicalized;
- all source sugar relevant to later semantic analysis is expanded;
- no later compiler stage is required to guess what omitted source notation meant.

The canonical result is governed by the Bitlang Preprocessed repository, especially:

- `LANGUAGE_SPEC.ja.md`
- `PROPERTIES.ja.md`
- `BORROW_STATE.ja.md`



## Legacy Preprocessed-shaped input

Older Bitlang Preprocessed representations may be accepted by preprocessing as migration input when their meaning can still be reconstructed deterministically.

A common version change is that a semantic property which was previously implicit becomes public and mandatory in the newer Bitlang Preprocessed specification. In that case, the older representation is no longer considered canonical current Preprocessed, because one or more semantic axes are still implicit.

Conceptually, such input is treated as **Bitlang source written in a Preprocessed-shaped syntax**:

```text
old Preprocessed-shaped input
    -> preprocess / infer formerly implicit properties
    -> materialize newly explicit properties
    -> current fully explicit Bitlang Preprocessed
```

The fact that the text resembles an older Preprocessed format does not grant it direct access to the compiler. The current Preprocessed boundary is defined by the current canonical specification, not by historical syntax.

This keeps backward migration in the preprocessor and allows the Bitlang compiler to remain strict: it only needs to accept the current fully explicit form.

## Design intent

Bitlang should be strict about **meaning**, not unnecessarily strict about **how much repetitive semantic information a human must type**.

The strictness of mandatory explicitness belongs at the normalized boundary. The source form may remain comparatively comfortable, but it retains access to the complete Bitlang property vocabulary at all times.
