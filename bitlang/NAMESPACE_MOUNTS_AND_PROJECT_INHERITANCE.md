# Namespace Mount and Project Inheritance

This document defines preprocessing-time namespace assignment and parent/child project inheritance for Bitlang source.

These facilities exist before Bitlang Explicit. They are preprocessing configuration, not additional runtime language constructs.

## Namespace Mount

A **Namespace Mount** assigns an outer naming prefix to files based on their project-relative path.

Bitlang itself continues to use `module` as its canonical naming hierarchy. Every file also resolves to an outer file module: an explicit top-of-file module declaration wins, otherwise the source file base name (without extension) is used. A namespace mount is preprocessing metadata placed outside that file module and is normalized into the corresponding Bitlang module/name hierarchy before Bitlang Explicit is emitted.

The standard preprocessing operations are:

```text
mount_namespace_file(file_path, namespace_path)
mount_namespace_tree(directory_path, namespace_path)
```

These names are canonical source-facing preprocessing built-in names.

### File mount

`mount_namespace_file` assigns the specified namespace prefix to one exact source file.

Conceptually:

```text
mount_namespace_file("src/main.bit", "Game.Client")
```

with no explicit file-level module declaration resolves the file under:

```text
Game.Client.main
```

because `main.bit` contributes the implicit file module `main`.

If the file starts with `module App;`, the effective outer module becomes `Game.Client.App` instead. Explicit scoped modules inside the file are nested beneath that.

### Tree mount

`mount_namespace_tree` assigns a namespace root to a directory tree.

Relative child directory names automatically extend the namespace hierarchy.

For example:

```text
mount_namespace_tree("src", "Game")
```

with:

```text
src/player.bit
src/combat/damage.bit
src/combat/effects/fire.bit
```

produces the effective outer hierarchy when those files do not declare their own file module:

```text
src/player.bit               -> Game.player
src/combat/damage.bit        -> Game.Combat.damage
src/combat/effects/fire.bit  -> Game.Combat.Effects.fire
```

The source filename base is the implicit file module. It is not added when the file explicitly chooses its outer module with a top-of-file module declaration.

For example, if `src/combat/damage.bit` begins with `module DamageEngine;`, its effective outer module is `Game.Combat.DamageEngine`.

Explicit scoped modules inside the file are nested below the effective file module.

### Deterministic path rules

- Paths are resolved relative to the active project root unless the API explicitly supplies another supported root.
- Directory segments used as automatic namespace segments must already be valid Bitlang identifier segments. The preprocessor must not silently sanitize an invalid directory name; an explicit mapping is required instead.
- An exact file mount has higher precedence than a directory-tree mount for that file.
- Among directory-tree mounts, the most specific matching directory wins as the namespace root.
- A nested directory mount replaces the inherited mount root for its subtree.
- If equally specific applicable mounts assign incompatible namespace roots and no explicit override/order relation resolves them, preprocessing reports an error.
- A single source file must resolve to one deterministic mount prefix and one deterministic file module.
- The implicit file module is the file base name with the extension removed; invalid identifier text is not silently sanitized.
- An explicit file-level module declaration replaces only the implicit filename-derived file-module segment, not the namespace-mount prefix.

Namespace mounts and implicit filename-module selection are consumed during preprocessing. The preprocessor normalizes the effective file module by enclosing the file body in an equivalent module scope. Bitlang Explicit contains only canonical module scopes/name ownership rather than mount metadata or an implicit filename rule.

## Parent/Child Projects

Bitlang projects support a deterministic single-parent inheritance chain.

The preprocessing operation is:

```text
set_parent_project(parent_project)
```

A project may have at most one direct parent project.

This is deliberately different from ordinary project dependencies:

- **parent project** = preprocessing/configuration inheritance;
- **dependency** = access to another project's code/artifacts.

Declaring a parent does not automatically compile or import all source files from the parent project.

### Inherited project configuration

A child project inherits applicable preprocessing configuration from its parent, including:

- project-level property defaults;
- namespace mounts;
- language-adapter defaults/configuration;
- reusable preprocessing configuration explicitly marked inheritable;
- future project-wide preprocessing settings that declare project inheritance semantics.

Inheritance proceeds from the root parent toward the leaf child.

A child may override an inherited setting where that setting is defined as overridable. The nearest child definition wins.

For project defaults, precedence therefore becomes:

```text
current project default
> nearest parent project default
> next ancestor project default
> ...
> language-adapter default
> Bitlang built-in default
```

The ordinary higher-precedence scopes still remain above project defaults:

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

### Namespace mounts inherited by child projects

Project-relative namespace-mount rules are inherited as configuration.

When inherited by a child project, project-relative path selectors are evaluated relative to the child project's root. This allows a base project to define a standard source layout such as:

```text
src/ -> Company.Product
```

and child projects using the same layout to inherit the namespace policy automatically.

A child project's own mount is more specific in configuration precedence than an inherited parent-project mount. Ordinary file-vs-tree and most-specific-path rules are then applied.

### Parent relation validation

Project inheritance must be deterministic.

- Parent-project cycles are errors.
- A project cannot have more than one direct parent.
- Missing or unresolved parents are errors when the parent relationship is required.
- Conflicting inherited settings that remain unresolved after normal precedence rules are errors.
- The resolved parent chain must be stable before source normalization begins.

## Safety

Namespace mounting and project inheritance do not weaken Bitlang safety.

They may affect name ownership and omitted-property defaults, but the final resolved Bitlang program is validated exactly as though every resulting module relationship and canonical property had been written explicitly.

Unsafe or invalid final semantics remain errors regardless of whether they originated from a local source declaration, an inherited project rule, a namespace mount, a language adapter, or generated preprocessing code.

## Bitlang Explicit boundary

By the Bitlang Explicit boundary:

- namespace mounts have been converted to canonical module/name hierarchy;
- project inheritance has been fully resolved;
- project/default configuration has been applied;
- no parent-project lookup is required by later compiler stages;
- all applicable canonical properties are explicit;
- unresolved namespace/project conflicts are errors.
