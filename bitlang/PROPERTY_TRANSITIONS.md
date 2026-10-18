# Bitlang Property Transitions

Bitlang variable/declaration properties may be flow-sensitive. A property's state at declaration time is not necessarily its state at every later program point.

## 1. Core model

A declaration has an initial resolved property state. Program operations may then perform valid property transitions.

Conceptually:

```text
declaration
    -> initial property state
    -> control flow
    -> property transition
    -> new effective property state
```

Static analysis evaluates the effective property state at the program point where an operation occurs.

## 2. Unreadable is useful

`Unreadable` does not mean that a variable is permanently useless.

It means that the value may not be read through that declaration/access path **at the current program point**.

For example, conceptually:

```text
Unreadable Writeable Uninitialized Int value;

value = 10;                 // write is permitted
change_property(value, Readable);
print(value);               // now readable
```

The exact source syntax for `change_property` is defined separately. The semantic transition is the important part.

`Unreadable` is useful for:

- output-only/write-only access during one phase;
- values that exist before they are allowed to be observed;
- temporary sealing while an invariant is being rebuilt;
- generated/adapted code that preserves another language's staged access rules;
- explicit capability hand-off between phases of a computation.

Changing `Unreadable -> Readable` does not initialize a value, resurrect a moved/released resource, or bypass another restriction. All other applicable properties must also permit the read.

## 3. Transitionable and structural axes

Not every property axis has the same transition rules.

Each axis defines whether it is:

```text
flow_transitionable
structural
operation_derived
```

- `flow_transitionable`: source/normalized operations may explicitly change the state at a program point, subject to validation.
- `structural`: fixed for the declaration/value identity during ordinary execution; changing it requires a different declaration/normalization contract rather than a runtime-flow transition.
- `operation_derived`: changes as the semantic result of an operation such as move, borrow, initialization, release, or reinitialization.

An axis may be both explicitly transitionable and operation-derived where its specification allows that.

## 4. Initial transition classification

The initial model is:

```text
readability       Readable / Unreadable             flow_transitionable
writeability      Writeable / Unwriteable           flow_transitionable
reassignment      Reassignable / Unreassignable     flow_transitionable

borrow_state      Unborrowed / Shared / Exclusive   operation_derived
move_state        Unmoved / Moved                   operation_derived
release_state     Unreleased / Released             operation_derived
initialization    Initialized / Uninitialized       operation_derived

scope_visibility                                  structural
inheritance_visibility                            structural
module_export                                     structural
retention / retention_domain                      structural
instance_access                                   structural
lifetime                                          structural
nullability                                       structural
optionality                                       structural
caller-class limit / call-site policy             structural
```

Copy/move capability, ownership, release policy/capability, constness, and initialization/finalization policy require axis-specific transition rules. They must not be assumed freely mutable merely because some variable properties are flow-sensitive.

Future language versions may add valid transitions explicitly.

## 5. Capability strengthening is validated

A transition that grants a capability is not automatically valid.

Examples:

```text
Unreadable  -> Readable
Unwriteable -> Writeable
Unreassignable -> Reassignable
```

are permitted only when no stronger invariant forbids the resulting operation.

For example:

- `Readable` does not permit reading an `Uninitialized`, `Moved`, or `Released` value;
- `Writeable` does not override `Const`, a conflicting borrow, or a read-only reference path;
- `Reassignable` does not make a non-rebindable `Ref<T>` rebindable;
- a transition cannot extend a borrow beyond its referent lifetime;
- a transition cannot manufacture ownership or a live resource.

Bitlang remains fail-closed: if the compiler cannot prove a required transition is valid, it rejects it.

## 6. Control flow

Property state is tracked per control-flow path.

Conceptually:

```text
if condition {
    change_property(x, Readable);
} else {
    // x remains Unreadable
}

// x cannot be assumed Readable here
```

At a control-flow join, the analyzer must not silently choose the more permissive incoming state.

An operation after a join is valid only when its required property state is proven on every reachable incoming path, or when control-flow refinement proves the necessary state on the active path.

This rule avoids introducing ad-hoc `Maybe_readable`, `Maybe_released`, etc. source properties merely for analysis joins.

## 7. Bitlang Explicit representation

Bitlang Explicit must represent both:

- the fully resolved initial property state of each declaration; and
- canonical property-transition operations that occur later.

Therefore the phrase "all properties are explicit" means explicit **at each relevant semantic point**, not that a declaration has one immutable property set for its entire lifetime.

Source shorthands and Japanese/English source lexicon spellings are normalized before these transition operations are emitted.

Explicit static analysis consumes the transition operations and validates accesses against the effective flow state.

## 8. Preprocessor changes versus program-flow transitions

These are different mechanisms:

```text
preprocessor override
    -> changes/generated the normalized program before Explicit

program-flow property transition
    -> is part of the normalized program semantics and appears in Explicit
```

A preprocessor may generate a property transition, but it must not erase the semantic distinction between a transition that happens at a program point and a declaration default chosen before execution.

## 9. Backend boundary

Property transitions are compile-time semantic information unless a specific feature requires runtime reflection of them.

After control-flow/property validation and lowering, a backend may erase the property metadata and transition markers when the resulting low-level operations already enforce the semantics.

No general runtime "property object" is required.
