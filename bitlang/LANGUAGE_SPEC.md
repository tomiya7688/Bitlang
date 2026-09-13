# Bitlang Language Specification

This document defines the Bitlang language specification.

## Role of Bitlang in the Bitlang family

Bitlang is the common semantic target of the Bitlang language family.

Languages in the Bitlang family may use different syntax, convenience features, or domain-specific notation, but they are expected to transform into ordinary Bitlang before the Bitlang preprocessing stage.

Conceptually:

```text
Bitlang family language
    -> transform
    -> Bitlang
    -> preprocess + normalize
    -> Bitlang preprocessed
    -> compile
    -> Bitlang compiled
```

This allows family languages to focus on their own syntax and usability while Bitlang remains the shared language for expressing the meaning of the program.

## Modules and naming hierarchy

Bitlang does not provide a separate `namespace` construct.

`module` is the single top-level organizational concept used for namespacing and may also carry broader responsibilities such as import/export boundaries, dependency management, visibility, compilation grouping, initialization boundaries, or future versioning rules.

A module may contain nested modules or named declarations so that fully qualified names can be expressed hierarchically.

Conceptually:

```text
Game.Combat.Player.attack
```

may identify `attack` as a member reached through the `Game` and `Combat` module hierarchy and the `Player` type.

Bitlang preprocessed should use fully qualified references wherever practical so that the referenced declaration is explicit and unambiguous. Source-level aliases or shortened forms may be provided through preprocessing rules, but they must normalize to the canonical fully qualified representation.

## Naming style

Bitlang does not force source-code authors to use snake_case, camelCase, PascalCase, or another naming style. Source-facing code may use whichever naming style is convenient, subject to the language's normal identifier rules.

Bitlang's own internal and canonical naming follows a stricter convention. Internal functions, class names, and property names formed by joining two or more words use `snake_case`.

Examples of canonical internal-style names include:

```text
memory_manager
current_scope
read_access
```

This convention is intended for Bitlang's own internal functions, generated/internal class names, and multi-word property names. It is not a mandatory coding-style rule imposed on user-written source.

The reason for preferring `snake_case` internally is that Bitlang identifiers are semantically case-insensitive. A naming scheme that relies on capitalization alone to mark word boundaries would lose information during canonical case normalization. Underscores preserve those word boundaries explicitly and therefore remain easy to read after normalization.

Preprocessor-generated canonical names are normally normalized so that the first character is uppercase and the remaining letters of each word are lowercase. Multi-word canonical names preserve `_` as the word separator.

Conceptually:

```text
currentScope
CurrentScope
CURRENT_SCOPE
```

may all resolve to the same semantic name and normalize to a canonical spelling such as:

```text
Current_scope
```

where the initial capitalization is canonical formatting rather than semantic identity.

Because Bitlang names are semantically case-insensitive, differences in capitalization do not create distinct identifiers. Naming-style normalization may therefore be performed by preprocessing where a canonical internal name is required.

## Property declaration and preprocessing

Bitlang source may explicitly declare semantic properties when the programmer wants direct control over them.

When an applicable property is omitted in source code, the preprocessor determines and fills in the appropriate property from the declaration kind, lexical context, module configuration, defaults, static analysis, or explicit preprocessor rules.

An explicitly written source property takes precedence over an ordinary inferred/default value unless another explicit language rule makes that combination invalid.

This means source-facing Bitlang may remain comparatively compact while still producing a fully explicit canonical representation.

Conceptually:

```text
int a = 4
```

may omit visibility, readability, writability, reassignment, initialization, nullability, optionality, lifetime, and other applicable properties. Preprocessing resolves those omitted properties before Bitlang preprocessed is produced.

The same source declaration may instead explicitly specify one or more of those properties when desired. Preprocessor functions may also inspect, add, remove, or change properties before canonical output is finalized.

Bitlang preprocessed must contain the resolved final property set and must not require later compiler stages to reconstruct omitted property semantics.

## Functional-language support

Bitlang must be able to represent functional-programming semantics even though Bitlang itself does not need to be primarily written as a functional language.

The Bitlang family may provide a dedicated **Bit Function lang** front end. Bit Function lang is intended to be written naturally in a functional style and then transform into ordinary Bitlang.

Bitlang should therefore be able to represent concepts such as:

- first-class functions
- functions stored in variables
- functions passed as arguments
- functions returned as values
- lambdas
- closures
- higher-order functions
- immutable values
- recursion
- pattern matching
- sum/variant-like data representations
- Option/Result-like values
- partial application and currying semantics where required by a family language

These features may have multiple convenient source-level notations in Bitlang or in Bitlang-family languages, but they must be normalized before reaching Bitlang preprocessed.

The intended functional pipeline is:

```text
Bit Function lang
    -> transform
    -> Bitlang functional-semantic representation
    -> preprocess + normalize
    -> Bitlang preprocessed canonical representation
    -> compile
    -> Bitlang compiled procedural representation
```

## Assignment and increment syntax

Bitlang source may provide compound-assignment convenience syntax such as `+=`, `-=`, `*=`, and `/=`. These forms are source-level sugar only and must be expanded during preprocessing into explicit assignment plus the corresponding operation.

For example:

```text
a += b
```

normalizes to the equivalent of:

```text
a = a + b
```

Prefix increment and decrement forms such as `++a` and `--a` are not part of Bitlang. They are intentionally unsupported because they combine mutation and expression evaluation in a way that is easy to misuse.

Postfix increment and decrement forms such as `a++` and `a--` are also not part of the canonical language. If convenience syntax of this kind is ever accepted by a source-facing Bitlang mode or family language, it must be restricted to a standalone mutation statement and normalized before Bitlang preprocessed. It must not be usable as a value-producing expression.

The preferred explicit forms are:

```text
a = a + 1
a = a - 1
```

This keeps mutation and value evaluation separate and removes increment/decrement side-effect semantics from the core language.

## Access and reassignment qualifiers

Bitlang separates read access, write access, and reassignment into independent semantic properties.

The access qualifiers are:

```text
Readable
Unreadable
Writeable
Unwriteable
```

`Readable` means the value may be read through the declaration or reference.

`Unreadable` means reading the value through that declaration or reference is prohibited.

`Writeable` means the value or reachable mutable state may be modified through the declaration or reference.

`Unwriteable` means modification through that declaration or reference is prohibited.

Readability and writability are independent axes.

Reassignment is expressed separately:

```text
Reassignable
Unreassignable
```

`Reassignable` means the declaration may later be assigned another value or binding of the same compatible type.

`Unreassignable` means the declaration cannot be rebound or assigned a replacement value after initialization. This does not by itself make the contained or referenced object immutable.

For example, an `Unreassignable Writeable` reference may remain bound to the same object while still allowing that object's writable state to be changed.

These qualifiers describe individual capabilities. Bitlang should prefer combinations of these explicit properties rather than broad source-language-style categories such as `mutable`, `dynamic`, or `flexible` when those categories can be represented more precisely by the independent qualifiers.

## Ownership qualifiers

Ownership is represented independently from pointer/reference type and independently from read, write, and reassignment capability.

The ownership qualifiers are:

```text
Owned
Borrowed
```

`Owned` means the declaration owns the lifetime responsibility for the represented resource or object. The owning declaration is responsible for ensuring that the resource is released or otherwise finalized according to the applicable storage model.

`Borrowed` means the declaration does not own the resource. It may use the resource only while the actual owner keeps it valid, and it must not independently release or finalize that resource.

Ownership qualifiers are attributes, not type wrappers. Therefore forms such as:

```text
Owned Ptr<MyType>
Borrowed Ref<MyType>
```

express ownership separately from `Ptr<T>` and `Ref<T>` themselves.

The compiler and static-analysis stages should use these qualifiers when checking lifetime and release responsibility, while preserving the distinction between ownership and access capability.

## Borrow state

Borrow state is represented independently from ownership.

The canonical borrow-state properties are:

```text
Unborrowed
Shared_borrowed
Exclusive_borrowed
```

`Unborrowed` means there is no active borrow that restricts ordinary access through the owning declaration.

`Shared_borrowed` means one or more shared borrows are active. Shared borrows may coexist when the applicable access rules allow it, but operations that would invalidate those borrows are restricted.

`Exclusive_borrowed` means an exclusive borrow is active. While it remains active, conflicting borrows or direct operations through other access paths are prohibited.

Borrow state is distinct from the `Borrowed` ownership qualifier. `Borrowed` answers whether a declaration owns a resource, while borrow state describes the current borrowing condition of a resource or declaration.

Source-facing Bitlang may write a borrow-state property explicitly or omit it. When omitted, preprocessing resolves the state from context and emits the final state in Bitlang preprocessed.

Preprocessor rules may deliberately change borrow state. Because forcing a state such as `Exclusive_borrowed -> Unborrowed` can re-enable access while a real borrow may still exist, such a rewrite is a semantic override and must be explicit. Unsafe overrides may produce warnings or errors.

Typical invalid or suspicious states include releasing a resource while it is borrowed, creating a conflicting borrow during `Exclusive_borrowed`, or allowing a borrow to outlive its source.

## Copy and move properties

Copyability and movability are represented as independent semantic properties.

The canonical property pairs are:

```text
Copyable
Uncopyable
Movable
Unmovable
```

`Copyable` means a value may be duplicated as another value with equivalent state according to the type's copy semantics.

`Uncopyable` means duplication is not permitted.

`Movable` means a value or owned resource may be transferred to another declaration according to the applicable move semantics.

`Unmovable` means such transfer is prohibited.

Copyability and movability are independent from ownership. A common owned-resource form may therefore be `Owned Uncopyable Movable`, while ordinary scalar values may be `Copyable Movable`.

Source-facing Bitlang may specify these properties explicitly or omit them. When omitted, preprocessing resolves them from the type, declaration kind, ownership state, context, and explicit preprocessor rules. The final resolved properties must be emitted in Bitlang preprocessed.

## Move state

Move state is represented by:

```text
Unmoved
Moved
```

A normal move changes the source declaration from `Unmoved` to `Moved`. A `Moved` declaration cannot ordinarily be read, released again, or moved again until a valid reinitialization or explicit preprocessing rule changes its state.

Preprocessor rules may deliberately change move state, including `Moved -> Unmoved`, but such a semantic override must be explicit and may produce warnings or errors when unsafe.

## Release state

Release state is represented independently from release capability and release policy.

The state properties are:

```text
Unreleased
Released
```

`Unreleased` means the resource represented by the declaration has not yet been released.

`Released` means that resource has already been released or otherwise destroyed and is no longer a valid live resource.

Releasing an `Unreleased` resource transitions it to `Released`.

After `Released`, ordinary access to the released resource, creation of new references to it, or a second release is invalid unless a later explicit operation establishes a new live resource for the declaration.

The declaration or handle itself may remain in scope after the underlying resource is released. This is intentional: the language can retain `Released` as a semantic state so static analysis can detect use-after-release and double-release. At runtime the underlying resource may already be gone, while the variable slot or handle may still exist until normal scope end or optimization removes it.

A valid reallocation or reinitialization may transition a declaration from `Released` to `Unreleased`. Preprocessor rules may also explicitly rewrite this state, but an override that claims a resource is live without establishing a valid resource is unsafe and should be diagnosed.

## Lifetime properties

Lifetime is represented independently from ownership, access capability, and pointer/reference type.

Canonical lifetime properties include:

```text
Local_lifetime
Function_lifetime
Object_lifetime
Module_lifetime
Static_lifetime
```

These properties describe the lifetime region in which a declaration or resource remains valid.

`Local_lifetime` describes a value whose lifetime is limited to its local lexical region.

`Function_lifetime` describes a value valid for the lifetime of the containing function invocation.

`Object_lifetime` describes state whose lifetime follows the lifetime of an owning object or instance.

`Module_lifetime` describes state retained for the lifetime of the active module.

`Static_lifetime` describes state retained for the applicable static-retention lifetime.

Lifetime properties do not by themselves define ownership. For example, a `Borrowed Ref<T>` may have a shorter lifetime than the `Owned` value it refers to.

The compiler and static-analysis stages must reject uses where a borrowed reference or other dependent value can outlive the value on which it depends.

Source-facing Bitlang may explicitly specify a lifetime property. When omitted, preprocessing resolves the applicable lifetime from context and emits it explicitly in Bitlang preprocessed.

## Initialization state

Initialization state is represented as an explicit property of a declaration.

The initialization properties are:

```text
Initialized
Uninitialized
```

`Initialized` means the declaration currently has a valid initialized value.

`Uninitialized` means the declaration exists but does not yet contain a valid initialized value.

An `Uninitialized` declaration must not be read as a value. It may transition to `Initialized` only through a valid initialization operation.

Initialization state is independent from readability, writability, reassignment, ownership, and other declaration properties. Preprocessor functions may inspect and transform this property under the same explicit-property rules used elsewhere in Bitlang.

Bitlang source may omit some initialization-state detail when it is directly inferable from the declaration syntax, but preprocessing must resolve that state explicitly.

## Const

`Const` is separate from `Readable`, `Unreadable`, `Writeable`, `Unwriteable`, `Reassignable`, and `Unreassignable`.

`Const` represents a strongly fixed value.

A `Const` declaration cannot be reassigned, and the value represented by it must remain semantically unchanged after initialization. For aggregate or object values, this strong fixed-value rule also applies to the state that forms part of that value rather than only to the outer variable binding.

Therefore `Const` is stronger than merely combining `Unreassignable` and `Unwriteable` as access restrictions: it states that the value itself is fixed, not only that one particular access path lacks permission to modify it.

Preprocessing and compilation must preserve this strong fixed-value meaning.

## Preprocessor functions

Bitlang provides preprocessor functions as a language-specific mechanism for reducing repetitive or inconvenient source code without weakening the strict Bitlang language model.

A preprocessor function is executed before ordinary compilation and does not remain as a runtime function.

Its role combines ideas similar to:

- Go-style code generation (`go generate`)
- C-style definitions/macros (`#define`)
- source-text transformation tools

Preprocessor functions may be used for three broad purposes:

1. **Define** - declarations, aliases, constants, abbreviations, or reusable definitions.
2. **Generate** - generate Bitlang declarations, functions, types, or other source structures.
3. **Transform** - transform existing source or another supported representation into Bitlang.

The output of preprocessing must be valid **Bitlang preprocessed** input.

Preprocessor functions are intended both to make strict Bitlang somewhat easier to write and to remove transformation/code-generation bottlenecks. They are also part of the infrastructure that can be used by Bitlang-family language transformers.

Detailed syntax and semantic rules will be added as the language design is finalized.
