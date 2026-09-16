# Bitlang Source Normalization

Bitlang source is intentionally more permissive and concise than Bitlang Preprocessed.

The source language is allowed to optimize for human readability. The preprocessor is responsible for converting those convenient forms into the strict, explicit canonical representation defined by `tomiya7688/Bitlang_preprocessed`.

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

## Explicit detail is still allowed

Bitlang does not force users to use the relaxed form.

A programmer, generator, language adapter, or debugging tool may write a highly explicit Bitlang declaration containing many semantic properties. The preprocessor preserves valid explicit requirements while filling only what remains unresolved.

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

At the Bitlang -> Bitlang Preprocessed boundary:

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

## Design intent

Bitlang should be strict about **meaning**, not unnecessarily strict about **how much repetitive semantic information a human must type**.

The strictness belongs at the normalized boundary. The source layer may remain comparatively comfortable as long as preprocessing can transform it into one unambiguous canonical form.
