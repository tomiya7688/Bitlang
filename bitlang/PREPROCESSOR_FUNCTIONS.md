# Bitlang Preprocessor Functions

Bitlang preprocessor functions are compile-time functions used to reduce repetitive source code and remove transformation or code-generation bottlenecks without weakening the strict Bitlang language model.

They combine ideas similar to code generation, macro definitions, and source transformation tools.

## Core language

The preprocessor-function language should support at least:

- variables
- user-defined functions
- `if` / `else`
- `for`
- `switch`
- arithmetic operators: `+`, `-`, `*`, `/`, `%`
- comparison operators
- boolean/logical operators
- string concatenation and basic string operations

## General-purpose preprocessing computation

The Bitlang preprocessing environment is not intentionally restricted to a non-general macro subset. It supports ordinary syntax decomposition and manipulation together with variables, functions, branching, looping, operators, and other shared Bitlang facilities.

As a result, the preprocessing language is intended to have general-purpose computational capability and, under the ordinary theoretical assumption of sufficient memory and execution resources, is Turing-complete.

The language specification does not forbid programmers from using `@preprocesser` code to implement computations that are unrelated to source generation or transformation. Such use may be unusual and may be inefficient compared with ordinary compiler-domain execution, but it is still valid preprocessing-domain code as long as it obeys the applicable preprocessing rules.

This does not merge the preprocessing and compiler execution domains. General computational capability in the preprocessor does not permit a compiler-domain function body or compiler-only instruction to execute during preprocessing. The distinction between the two domains remains based on execution phase and available operations, not on the computational expressiveness of the syntax.

## Execution domains and shared syntax

Bitlang source has two execution domains:

```text
@compiler
@preprocesser
```

Ordinary declarations belong to the `@compiler` domain by default. `@compiler` may be written explicitly, but is normally unnecessary.

A declaration marked with `@preprocesser` belongs to the preprocessing domain and is executed by the Bitlang preprocessor rather than becoming an ordinary runtime function.

Conceptually:

```bitlang
Public void Runtime_function()
{
}

@compiler
Public void Explicit_runtime_function()
{
}

@preprocesser
Public static void Generate_code()
{
}
```

The two domains share the ordinary Bitlang syntax wherever the operation has meaning in both environments. Variable declarations and ordinary value syntax are intentionally common rather than defining a separate variable language for preprocessing.

At minimum, the following language facilities are shared between compiler code and preprocessor code:

- arithmetic and ordinary operators;
- conditional control flow such as `if` / `else`;
- loop control such as `for`;
- type information and type-oriented expressions available to the corresponding environment;
- `console.log`;
- ordinary variable declaration and value syntax.

A standard function or built-in may still be domain-specific when its operation only makes sense in one environment. A compiler-only instruction does not become executable merely because its name is referenced from a preprocessor function, and likewise a preprocessing-only operation must not survive as a runtime compiler instruction.

### Calling and referencing compiler functions from preprocessing

Preprocessor code uses ordinary function-call syntax. There is no special call operator solely for invoking a preprocessor function.

A compiler-domain function may also be referenced from preprocessing code and may be treated as a function value. It may therefore be stored in a variable, passed to another preprocessing operation, inspected as program structure, selected as a transformation target, or used while generating compiler-domain Bitlang code.

However, the body of a compiler-domain function is not executed by the preprocessor.

For example, a preprocessor rule may select a compiler function and arrange for a call to that function to exist in generated Bitlang code, or may add a compiler function to files matching some structural condition. This is manipulation of the compiler program, not execution of that runtime function during preprocessing.

The distinction is therefore:

```text
reference / store / pass / inspect / generate compiler function
    -> allowed during preprocessing

execute compiler function body during preprocessing
    -> not allowed
```

Function identity and function-value handling do not imply cross-domain execution capability.

## Built-in function groups

Built-ins should cover these groups:

### Definition

- define values, aliases, or reusable definitions
- remove definitions
- test whether a definition exists

### Text processing

- replace text
- split and join text
- substring operations
- formatting

### Bitlang generation

- emit Bitlang source text
- emit variables
- emit functions
- emit structs
- emit classes
- emit types

Structured generation is preferred when practical so generated output remains valid Bitlang.

### Source and structure inspection

- enumerate functions
- enumerate variables
- enumerate types
- inspect properties and attributes
- access structured source information
- obtain the set of all variables declared in the current structural scope

The current-scope variable-set operation is a target-selection mechanism, not a global mode. A preprocessor function may use that set as an argument or target and apply one transformation or rule to every variable in the current scope.

The meaning of current scope follows the normal structural scope in which the preprocessor function is evaluated, such as a function, class, struct, file, or lexical block.

Nested scopes remain distinct unless the preprocessor function explicitly traverses into them. This prevents a rule intended for one scope from silently affecting variables owned by child scopes.

### Environment and diagnostics

- current source path and source directory
- include or read source-related text where permitted by the preprocessing environment
- compile-time errors
- compile-time warnings



### Release-generation control

The preprocessing environment may control automatic cleanup generation.

At minimum, it provides a source-facing operation conceptually named `disable_auto_release` that prevents the preprocessor from synthesizing an automatic release/`free` for the affected declaration or activation range.

It may target a single variable/declaration directly. It may also be activated as a range rule, in which case every applicable variable/declaration encountered while the rule is active has automatic release generation suppressed.

This operation does not make release impossible. Explicit release operations remain valid when allowed by the target's properties.

The command is preprocessing-only and disappears before Bitlang Preprocessed. Its semantic result is represented through the resolved release policy, normally `Manual_release`, together with any explicit cleanup that remains in the normalized program.

The operation may be applied using the normal preprocessing activation-range mechanism, so a project may disable automatic release generation for one declaration, a local range, a class/function field, or another supported preprocessing scope. Range activation follows the ordinary start/end semantics; omitting the end marker keeps the suppression active until the end of the enclosing declaration field.

### Module configuration

Preprocessor functions may inspect and modify the configuration of the active Bitlang module.

The module configuration API should be able to work with module-level concerns such as:

- module name and hierarchical path
- imports and dependencies
- exports and public visibility rules
- aliases and name-resolution shortcuts
- initialization settings or initialization requirements
- compile-related module settings
- module metadata and future version-related settings

Preprocessor functions may therefore generate or alter module configuration before Bitlang preprocessed output is produced.

Any source-facing shorthand introduced through module configuration must be resolved during preprocessing. Bitlang preprocessed should retain the normalized, explicit result rather than depending on source-only aliases or preprocessor state.

Module configuration changes are compile-time operations and must not silently become runtime mutation of module state.

## Property-driven semantics

Bitlang should represent complicated declaration behavior as explicit properties wherever practical rather than hiding several independent meanings behind one broad keyword.

Examples include properties such as:

- `Readable` / `Unreadable`
- `Writeable` / `Unwriteable`
- `Reassignable` / `Unreassignable`
- `Owned` / `Borrowed`
- `Moved` / `Unmoved`

These properties are independent semantic axes unless a specific language rule states otherwise.

Preprocessor functions may inspect, add, remove, replace, or otherwise modify these properties during preprocessing. This allows source code to use concise declarations or project-level rules while still producing a strict and explicit Bitlang preprocessed result.

Property changes are compile-time transformations. They must be fully resolved before preprocessing finishes; Bitlang preprocessed must contain the resulting explicit properties and must not depend on hidden mutable preprocessor state.

A preprocessor function may change properties on a single declaration, a selected declaration set, or a current-scope variable set.

Semantic conflicts between properties must be diagnosed rather than silently resolved unless an explicit preprocessing rule defines how to resolve them.

`Const` remains a stronger dedicated semantic concept rather than merely another combination of access properties.

## Attribute and property automation

Preprocessor functions may automatically add, remove, inspect, or provide default attributes or properties for Bitlang declarations.

The system should support operations equivalent to:

- add a property or attribute to a target
- remove a property or attribute from a target
- replace one property with another
- test whether a property or attribute is present
- set default properties or attributes for a declaration category
- infer properties or attributes through static analysis when the result is provable

Automation exists to reduce repetitive annotation while keeping Bitlang's semantic model explicit.

Properties or attributes that can be proven without changing program behavior may be attached automatically. Changes that could alter runtime behavior or program meaning must require an explicit rule, configuration, or declaration rather than heuristic inference alone.

Move state follows the same rule. A preprocessor function may explicitly rewrite `Moved` to `Unmoved` or vice versa, but a `Moved -> Unmoved` rewrite is a semantic override, not an ordinary inferred default. It must therefore come from an explicit preprocessing rule, configuration, or source-directed transformation.

This capability exists so transformers and project-specific preprocessing can deliberately restore or reinterpret state when they have additional semantic knowledge that the ordinary move rules do not capture.

When a move-state override cannot be proven safe, preprocessing may emit a warning. If the resulting state is provably invalid, preprocessing must emit an error.

## Safety diagnostics during preprocessing

The preprocessor may emit warnings when the resolved property set or control flow strongly suggests a dangerous resource or lifetime state even if the program is not yet provably invalid.

Examples include:

- an `Owned Manual_release Releasable` resource reaching the end of its lifetime without any visible release path
- a borrowed or referenced value whose lifetime appears likely to exceed that of its source
- ownership states that are technically representable but leave release responsibility ambiguous
- a moved or otherwise transferred resource that still appears to be used through a stale access path
- an explicit `Moved -> Unmoved` override for which validity cannot be proven

These situations should normally produce warnings when risk is detected but correctness cannot be proven either way.

If the preprocessor or static analysis can prove that the resulting program violates a Bitlang semantic rule, the diagnostic should be an error rather than only a warning.

The intended distinction is:

- **warning**: suspicious or dangerous state, but not conclusively invalid
- **error**: semantic invalidity can be proven

Warnings must not silently rewrite runtime semantics merely to make the warning disappear. Any automatic correction that changes program meaning requires an explicit preprocessing rule or configuration.

## Activation scope

A preprocessor macro or preprocessor function may define an explicit activation start point and may optionally define an explicit end point.

When an end point is omitted, the active range ends automatically at the end of the nearest enclosing structural scope.

Default scope behavior:

- declared inside a function: active until the end of that function
- declared at file level: active until the end of that file
- declared inside a class: active in that class and, when inheritance propagation applies, in child classes
- declared for a variable: active only for that variable/declaration scope
- declared inside a struct: active until the end of that struct definition
- declared inside another lexical block: active until the end of that block

An explicit end marker may be used to terminate the active range before the natural end of the structural scope.

Class-scope propagation should be representable explicitly so a rule can either remain local to the class or propagate to inheriting child classes.

Preprocessor activation markers are preprocessing-only constructs. They do not remain in Bitlang preprocessed output; only their expanded and normalized effects remain.

## Role in the Bitlang family

Bitlang-family languages generally transform into ordinary Bitlang before normal Bitlang preprocessing.

```text
Bitlang family source
    -> transform
    -> Bitlang
    -> preprocessor functions
    -> Bitlang preprocessed
    -> compiler
    -> Bitlang compiled
```

Preprocessor functions may also be reused as part of family-language transformation infrastructure when appropriate.
