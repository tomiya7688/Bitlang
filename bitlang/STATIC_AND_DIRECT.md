# static and Direct

This document defines the source-facing Bitlang rules for the `static` and `Direct` modifiers.

These rules close language-design issues:

- [#1 `static` means static retention only](https://github.com/tomiya7688/Bitlang/issues/1)
- [#2 dedicated modifier for instance-free access](https://github.com/tomiya7688/Bitlang/issues/2)

## Definitions

### `static`

> `static` marks a target that exists in a static retention region independent of ordinary local lifetime or instance lifetime, and retains its information or state while that region remains active.

- `static` does **not** mean "instance-free".
- The same basic meaning applies to every applicable target kind (variables, functions, type members, and other targets that admit static retention).

### `Direct`

> `Direct` marks a target that does not require an instance for access or invocation. It does not affect lifetime, retention, or initialization policy by itself.

- The official keyword is `Direct`.
- `Direct` is fully independent of `static`.

## Four orthogonal states

All four combinations are **syntactically permitted** where both modifiers are applicable:

| static | Direct | Meaning |
| --- | --- | --- |
| no | no | Ordinary instance-dependent target |
| yes | no | Statically retained, but still requires an instance for access |
| no | yes | Instance-free, but not statically retained |
| yes | yes | Statically retained and instance-free |

### Policy

> All four states are syntactically allowed. The preprocessor and later static analysis must verify that a given combination does not produce semantic contradictions or illegal access paths when combined with target kind, context, lifetime, ownership, and borrow state. Failures must be diagnosed (error or warning). Silent acceptance is forbidden.

## Target-specific applicability

### Functions

| Kind | static | Direct |
| --- | --- | --- |
| Type-member function | allowed | allowed (all four states) |
| Free function | allowed | **not allowed** (`Direct` is an error) |
| Local function | allowed | **not allowed** (`Direct` is an error) |

Meaning for type-member functions:

| static | Direct | Meaning |
| --- | --- | --- |
| no | no | Ordinary instance method; receiver required |
| yes | no | Has static retention, but still requires an instance to call |
| no | yes | Callable without an instance; state is not statically retained across calls |
| yes | yes | Callable without an instance and has static retention |

### Variables

| Kind | static | Direct |
| --- | --- | --- |
| Type-member variable | allowed | allowed (all four states) |
| Local variable | allowed | **not allowed** (`Direct` is an error) |
| Parameter | **not allowed** | **not allowed** (both are errors) |

Meaning for type-member variables:

| static | Direct | Meaning |
| --- | --- | --- |
| no | no | Ordinary instance field |
| yes | no | Statically retained field that still requires an instance for access |
| no | yes | Accessible without an instance, but not statically retained |
| yes | yes | Accessible without an instance and statically retained |

### Module-level declarations

| Target | static | Direct |
| --- | --- | --- |
| Module-level function | allowed | generally **not meaningful** (warning by default; may escalate to error by configuration) |
| Module-level variable | allowed | generally **not meaningful** (same as above) |

Reason: a module does not have an instance, so `Direct` adds no meaning. `static` remains valid when module-scoped static retention is required.

### Inheritance, override, and interfaces

- Presence or absence of `static` is part of the contract.
- Presence or absence of `Direct` is part of the contract.
- An override must not change either modifier relative to the base declaration.
- An interface implementation must satisfy the modifier contract required by the interface.

## Preprocessor verification obligations

Verification priority:

1. **Inapplicable target** → error  
   (for example `Direct` on locals, or `static`/`Direct` on parameters)

2. **Contract violation** → error  
   (override or interface mismatch of `static` / `Direct`)

3. **Access-path contradiction** → error  
   (for example using a non-`Direct` target through an instance-free path)

4. **Conflict with lifetime, ownership, or borrow state** → error or warning  
   (error when provably invalid; warning when dangerous but not fully provable)

5. **Redundant specification** → warning  
   (for example module-level `Direct`)

**Forbidden:** accepting contradictions or dangerous combinations silently.

Representative checks:

- A non-`Direct` function or member must not be invoked or accessed without an instance.
- A `Direct` function must not depend on an implicit receiver / `this`.
- A non-`static` `Direct` function must not attempt to retain state across calls as if it were statically retained.
- Override and interface pairs must keep matching `static` and `Direct` contracts.
- Access paths must remain consistent with ownership, borrow state, move state, and release state.
- Module-level `Direct` must be diagnosed as redundant unless a future module-instance model explicitly gives it meaning.

## Summary for implementers

```text
static and Direct are independent modifiers.
static means static retention only.
Direct means instance-free access only.

For type-member functions and type-member variables,
all four combinations of static and Direct are syntactically allowed.

Direct is not allowed on free functions, local functions, or local variables.
Neither static nor Direct is allowed on parameters.

At module level, Direct is generally redundant and must be diagnosed.

For overrides and interface implementations, the presence or absence of
static and Direct is part of the contract and must not change.

The preprocessor must verify that combinations do not contradict lifetime,
ownership, borrow state, or access paths.
Provable invalidity is an error; unprovable danger is a warning.
Silent acceptance is forbidden.
```

## Intentionally deferred

The following remain open and are not decided by this document:

- start and end timing of the static retention region
- retention scope per thread / task
- interaction with destructor-like finalization
- concrete Bitlang Preprocessed / Compiled encodings of `static` and `Direct`
- whether a future module-instance model should give module-level `Direct` real meaning
- detailed initialization-order rules for static targets
