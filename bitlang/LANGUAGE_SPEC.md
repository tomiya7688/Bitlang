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
    -> Bitlang Explicit
    -> Bitlang Lowerer
    -> Bitlang Low
```

This allows family languages to focus on their own syntax and usability while Bitlang remains the shared language for expressing the meaning of the program.

## Source and Explicit are the same language

Bitlang source and Bitlang Explicit use the same underlying semantic property system.

Every canonical property that may appear in Bitlang Explicit may also be written explicitly in ordinary Bitlang source when that property applies to the declaration.

The difference is normalization strictness:

```text
Bitlang source
    -> properties may be explicit or omitted
    -> source sugar and preprocessing-only constructs may exist

Bitlang Explicit
    -> the same applicable properties are all resolved and explicit
    -> source sugar and preprocessing-only constructs have been consumed
```

Therefore Bitlang Explicit is the fully explicit normalized form/profile of Bitlang, not a separate semantic language with a property system unavailable to source authors.

A programmer is allowed to write source that is already highly explicit. If all applicable properties are supplied and no source-only constructs remain, preprocessing may have very little semantic information left to add.

The repositories are separated to keep stage-specific specifications and implementations manageable; repository separation does not imply language-semantic separation.

## Modules and naming hierarchy

Bitlang does not provide a separate `namespace` construct.

`module` is the single top-level organizational concept used for namespacing and may also carry broader responsibilities such as import/export boundaries, dependency management, visibility, compilation grouping, initialization boundaries, or future versioning rules.

A module may contain nested modules or named declarations so that fully qualified names can be expressed hierarchically.

Conceptually:

```text
Game.Combat.Player.attack
```

may identify `attack` as a member reached through the `Game` and `Combat` module hierarchy and the `Player` type.

Bitlang Explicit should use fully qualified references wherever practical so that the referenced declaration is explicit and unambiguous. Source-level aliases or shortened forms may be provided through preprocessing rules, but they must normalize to the canonical fully qualified representation.

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

Case differences do not create distinct identifiers, but underscores are significant characters and are preserved as part of identifier identity.

Conceptually:

```text
currentScope == CurrentScope
current_scope == CURRENT_SCOPE
currentScope != current_scope
```

Canonical case formatting may normalize capitalization, but generic identifier normalization must not insert, remove, or ignore `_`.

Bitlang-defined standard-library names and built-in instructions formed from two or more words use `snake_case`. User-defined identifiers remain free to use other naming styles, but a differently placed underscore still denotes a different identifier.

## Property declaration and preprocessing

Bitlang source may explicitly declare **any canonical Bitlang property** that can appear in Bitlang Explicit when that property applies to the declaration. Source authors are never required to rely on inference merely because a property is normally filled by preprocessing.

When an applicable property is omitted in source code, the preprocessor determines and fills in the appropriate property from the declaration kind, lexical context, module configuration, defaults, static analysis, or explicit preprocessor rules.

An explicitly written source property takes precedence over an ordinary inferred/default value unless another explicit language rule makes that combination invalid.

This means source-facing Bitlang may remain comparatively compact while still producing a fully explicit canonical representation.

Conceptually:

```text
int a = 4
```

may omit visibility, readability, writability, reassignment, initialization, nullability, optionality, lifetime, and other applicable properties. Preprocessing resolves those omitted properties before Bitlang Explicit is produced.

The same source declaration may instead explicitly specify one or more of those properties when desired. Preprocessor functions may also inspect, add, remove, or change properties before canonical output is finalized.

Bitlang Explicit must contain the resolved final property set and must not require later compiler stages to reconstruct omitted property semantics.

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

These features may have multiple convenient source-level notations in Bitlang or in Bitlang-family languages, but they must be normalized before reaching Bitlang Explicit.

The intended functional pipeline is:

```text
Bit Function lang
    -> transform
    -> Bitlang functional-semantic representation
    -> preprocess + normalize
    -> Bitlang Explicit canonical representation
    -> Bitlang Lowerer
    -> Bitlang Low procedural representation
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

Postfix increment and decrement forms such as `a++` and `a--` are also not part of the canonical language. If convenience syntax of this kind is ever accepted by a source-facing Bitlang mode or family language, it must be restricted to a standalone mutation statement and normalized before Bitlang Explicit. It must not be usable as a value-producing expression.

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

These qualifiers describe individual capabilities. Bitlang should prefer combinations of these explicit properties rather than broad source-language-style categories such as `mutable` or `flexible` when those categories can be represented more precisely by the independent qualifiers. `Dynamic` is reserved for the explicit opposite of `Static` on the retention axis and does not mean general dynamism or mutability.

## Static retention and instance access

Bitlang separates static retention from whether an instance is required for access.

The retention axis is:

```text
Static
Dynamic
```

`Static` means the target is statically retained. `Dynamic` is its explicit opposite and means the target is not statically retained, instead following the applicable ordinary non-static retention/lifetime relationship.

This use of `Dynamic` does not mean dynamic typing, dynamic dispatch, or general runtime mutability.

The instance-access axis is:

```text
Instance_required
Instance_unrequired
```

`Instance_required` means access to an applicable type member requires an instance. `Instance_unrequired` means the member can be accessed without an instance.

These axes are independent:

```text
Static  + Instance_required
Static  + Instance_unrequired
Dynamic + Instance_required
Dynamic + Instance_unrequired
```

A combination is only valid when the declaration kind can meaningfully support it.

Bitlang source provides the convenience modifier:

```text
Direct
```

which normalizes to:

```text
Instance_unrequired
```

`Direct` does not imply `Static`. It changes only the instance-access requirement.

Retention and lifetime also remain distinct. `Static / Dynamic` is not a replacement for lifetime properties such as `Static_lifetime`; preprocessing and static analysis must ensure that the resolved retention and lifetime states are consistent.

### Retention domain

Static retention and the execution domain in which retained state is shared are separate semantic axes.

The retention-domain axis is:

```text
Process_retention
Thread_retention
Task_retention
```

- `Process_retention`: one retained state is shared across the process/program.
- `Thread_retention`: each thread owns an independent retained state.
- `Task_retention`: each task/coroutine-like execution unit owns an independent retained state.

This axis applies to variables, fields, and functions. For functions it applies to the same function-associated state governed by `Static / Dynamic`.

`Static / Dynamic`, lifetime, and retention domain are independent properties. A source-language adapter may therefore preserve thread-local or task-local static semantics without redefining `Static`.

When omitted in ordinary Bitlang source:

```text
retention domain -> Process_retention
```

Bitlang Explicit must always contain the resolved retention-domain state where the axis applies.

### Static initialization and finalization order

Retained state uses deterministic initialization and finalization ordering.

Initialization order is derived from a dependency graph. The graph includes:

- explicit initialization-order constraints supplied by source, a language adapter, or preprocessing;
- statically known dependencies used by an initializer;
- owner/module initialization dependencies required by the declaration.

Dependencies are initialized before dependents.

When multiple declarations are otherwise unordered, Bitlang uses deterministic tie-breaking:

1. declarations in the same declaration field use lexical declaration order;
2. declarations in unrelated fields/modules use canonical fully qualified declaration-name order.

A language adapter may supply explicit ordering constraints when preserving another language's initialization semantics. Exact source syntax for those constraints is defined separately; the semantic relation itself must survive into Bitlang Explicit when required.

Lazy triggers remain lazy:

- `First_reach_initialization` and `First_use_initialization` are not eagerly initialized merely to establish global order;
- when triggered, any unresolved dependencies required by that state are initialized first using the same dependency rules;
- `Manual_initialization` is excluded from automatic initialization ordering.

Initialization cycles that are provable statically are compile errors. If a lazy/dynamic initialization cycle can only be discovered by runtime re-entry into a state whose initialization is already in progress, the runtime must fail explicitly rather than observe a partially initialized value.

Finalization defaults to the reverse of the **actual successful initialization order** within each retention-domain instance. This rule naturally handles lazy initialization and means a dependency normally remains alive while its dependent is finalized.

Only states whose initialization completed successfully participate in automatic finalization.

If explicit finalization-order constraints are supplied, they may refine the default order but must not contradict required dependency safety. A contradictory order is an error.

For `Thread_retention` and `Task_retention`, each thread/task instance maintains its own actual initialization order and therefore its own reverse finalization order. `Process_retention` uses the process-wide retained-state order.

The finalization trigger determines **when** a state becomes eligible for finalization; this ordering rule determines the relative order among states being finalized. Cross-trigger combinations must still satisfy lifetime and destruction-safety rules.

### Function retention semantics

For a function declaration, `Static / Dynamic` governs function-associated semantic state rather than executable-code lifetime.

Function-associated state may include closure environments, captured storage, first-class function-object state, or other state owned by the function representation.

```text
Static function
    -> function-associated state uses static retention

Dynamic function
    -> function-associated state follows its ordinary owner/lifetime
```

The function's executable code is not considered to be created and destroyed on each invocation merely because the function is `Dynamic`.

For an ordinary stateless named function, the retention property may have no observable runtime storage effect. It still remains explicit in Bitlang Explicit because it is part of the declaration's semantic contract.

`Instance_required / Instance_unrequired` remains independent. Therefore all meaningful combinations of function retention and instance access remain representable.

Initialization and finalization of function-associated state are controlled by their own properties/rules; `Static` does not imply a specific initialization or destruction time.

### Source default resolution

When omitted in ordinary Bitlang source:

```text
retention       -> Dynamic
instance access -> Instance_required
```

Initialization-trigger defaults are context-sensitive:

```text
dynamic local variable      -> Declaration_initialization
static local variable       -> First_reach_initialization
dynamic instance field      -> Owner_initialization
static field                -> Owner_initialization
module/file-level variable  -> Owner_initialization
```

`First_use_initialization` and `Manual_initialization` require explicit selection unless a language adapter or explicit preprocessing rule supplies them.

These defaults belong to Bitlang source normalization. Bitlang Explicit never relies on them implicitly; it contains the resolved property explicitly.

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

Source-facing Bitlang may write a borrow-state property explicitly or omit it. When omitted, preprocessing resolves the state from context and emits the final state in Bitlang Explicit.

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

Source-facing Bitlang may specify these properties explicitly or omit them. When omitted, preprocessing resolves them from the type, declaration kind, ownership state, context, and explicit preprocessor rules. The final resolved properties must be emitted in Bitlang Explicit.

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

Source-facing Bitlang may explicitly specify a lifetime property. When omitted, preprocessing resolves the applicable lifetime from context and emits it explicitly in Bitlang Explicit.

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

## Initialization trigger

Bitlang represents the timing of automatic initialization independently from both retention and current initialization state.

The canonical initialization-trigger states are:

```text
Declaration_initialization
Owner_initialization
First_reach_initialization
First_use_initialization
Manual_initialization
```

Their meanings are:

- `Declaration_initialization`: initialize at the ordinary declaration-initialization point for each storage instance.
- `Owner_initialization`: initialize when the owning object, type, module, or corresponding owner is initialized.
- `First_reach_initialization`: initialize once when execution first reaches the declaration for that storage instance.
- `First_use_initialization`: initialize once when that storage instance is first validly used.
- `Manual_initialization`: perform no automatic initialization; a valid explicit initialization operation is required before use.

This axis applies to variables and fields.

The trigger is independent from `Static / Dynamic`. In particular, `Static` does not imply eager initialization, first-use initialization, or first-reach initialization.

This separation exists so Bitlang-family language adapters can preserve different source-language initialization semantics without overloading the meaning of `Static`.

The trigger is also independent from `Initialized / Uninitialized`: the trigger describes the policy/timing, while the initialization-state axis describes the current semantic state.

## Finalization trigger

Bitlang represents destruction/finalization timing independently from lifetime, retention, and release policy.

The canonical finalization-trigger states are:

```text
Scope_end_finalization
Owner_end_finalization
Module_end_finalization
Program_end_finalization
Manual_finalization
```

Their meanings are:

- `Scope_end_finalization`: finalize automatically when the owning lexical/function scope ends.
- `Owner_end_finalization`: finalize automatically when the owning object/type/storage owner ends.
- `Module_end_finalization`: finalize automatically when the owning module is finalized or unloaded.
- `Program_end_finalization`: finalize automatically during program/process termination.
- `Manual_finalization`: no automatic finalization trigger; explicit finalization is required when applicable.

The canonical axis applies to variables, fields, and parameters.

Finalization is distinct from `Auto_release / Manual_release`. Finalization may run destructor/finalizer logic without necessarily performing memory/resource release, and release policy may generate release independently when valid.

When omitted in Bitlang source, the default is derived from resolved lifetime:

```text
Local_lifetime    -> Scope_end_finalization
Function_lifetime -> Scope_end_finalization
Object_lifetime   -> Owner_end_finalization
Module_lifetime   -> Module_end_finalization
Static_lifetime   -> Program_end_finalization
```

`Manual_finalization` requires explicit selection unless supplied by an explicit language-adapter or preprocessing rule.

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

The output of preprocessing must be valid **Bitlang Explicit** input.

Preprocessor functions are intended both to make strict Bitlang somewhat easier to write and to remove transformation/code-generation bottlenecks. They are also part of the infrastructure that can be used by Bitlang-family language transformers.

Detailed syntax and semantic rules will be added as the language design is finalized.
