# Bitlang Release and Disposal Properties

This document defines how **Bitlang source** expresses and preprocesses release/disposal semantics.

The canonical fully explicit release-property model belongs to Bitlang Preprocessed:

- https://github.com/tomiya7688/Bitlang_preprocessed/blob/main/PROPERTIES.ja.md

## Source properties

Release behavior and release capability are independent semantic axes.

```text
Auto_release
Manual_release
Releasable
Unreleasable
```

- `Auto_release`: source semantics allow required release/finalization to be performed automatically when the resolved lifetime ends, according to the applicable storage model.
- `Manual_release`: automatic release must not be assumed; an explicit valid release operation is required when the resource requires release.
- `Releasable`: release/finalization may be performed through that declaration or access path.
- `Unreleasable`: that declaration or access path may not release or finalize the resource.

These properties are independent from ownership, lifetime, pointer/reference type, copyability, and movability, although those properties constrain legal combinations.

A borrowed declaration does not own release responsibility and therefore normally resolves to `Unreleasable`.

Examples:

```text
Owned Manual_release Releasable Ptr<My_type>
Borrowed Unreleasable Ref<My_type>
```

## Omission and preprocessing

Source-facing Bitlang may state these properties explicitly or omit them. When omitted, preprocessing resolves them from ownership, type, storage model, lifetime, context, and explicit preprocessing rules.

Invalid property combinations must be rejected rather than silently changing ownership or release semantics.

The current release state itself (`Unreleased / Released`) is also resolved before canonical output, but the exact required Preprocessed representation and its cross-property consistency rules are defined in the separate Bitlang Preprocessed specification.


## Preprocessing-time automatic release generation

Bitlang does not require a tracing runtime garbage collector as the normal memory-management model.

For owned resources that require release, the preprocessor may complete missing cleanup by generating an explicit release operation (for example, a `free`-equivalent operation for an applicable allocation model) at a provably valid lifetime boundary.

Conceptually:

```text
source omits required release
    -> preprocessing analyzes ownership and lifetime
    -> preprocessing inserts the required explicit release
    -> Bitlang Preprocessed contains the resolved cleanup behavior
```

This is GC-like convenience implemented as source preprocessing and explicit generated cleanup, not an implicit runtime collector.

Automatic insertion must obey ownership, borrow state, release capability, move state, lifetime, and other applicable semantic rules. The preprocessor must not insert a release at a point where doing so would be invalid.

## Suppressing generated release

Source preprocessing may explicitly disable automatic release generation for a declaration or preprocessing range.

The target may be a specific variable/declaration, in which case only that target is affected. The operation may also be activated for an entire preprocessing range; while that range is active, automatic release generation is disabled for every applicable declaration encountered within the range.

The source-facing preprocessing operation is conceptually `disable_auto_release`. Its exact invocation form may be used directly on a variable/declaration or through the normal preprocessor-function activation-range mechanism.

Range activation follows the ordinary preprocessor start/end rules. An explicit end marker stops the suppression there; when the end marker is omitted, suppression continues to the end of the enclosing declaration field.

Applying this operation means that the preprocessor must not synthesize an automatic `free`/release for the affected target merely because its lifetime ends without an explicit release.

Semantically, the affected release policy resolves to `Manual_release` unless another explicit and valid source rule provides an equivalent manual-management state.

This operation suppresses **generated** release only. It does not prohibit the programmer from writing an explicit valid release operation.

Because preprocessing controls are consumed before canonical output, `disable_auto_release` itself does not remain in Bitlang Preprocessed. The resulting explicit release policy and any user-written/generated runtime operations are what remain.

A manual-release resource that appears to reach the end of its lifetime unreleased may still produce a warning. Explicit suppression of automatic release indicates deliberate manual management; it does not disable safety analysis or make an otherwise invalid ownership/lifetime state valid.

Where the programmer deliberately transfers ownership, relies on an external lifetime, or intentionally leaves a resource for process termination, that intent should be expressible explicitly so diagnostics can distinguish it from an accidental forgotten release.


## Mandatory destruction safety

Release, destruction, destructor/finalizer execution, and owner teardown are subject to mandatory Bitlang safety validation.

A destruction operation that is proven to violate ownership, borrow, lifetime, release-state, initialization/finalization state, or another applicable safety rule is an **error**. It must not be downgraded to a warning merely because the operation was generated automatically.

This rule applies equally to:

- programmer-written Bitlang source;
- Bitlang-family language conversion/adaptor output;
- generated code produced by macros or preprocessor functions;
- automatic release/finalization inserted by the preprocessor;
- compiler-generated lowering and cleanup paths.

Examples of provably invalid destruction include:

- releasing or finalizing a resource that is already `Released`;
- destroying an owner while a live borrow/reference would thereby become invalid;
- releasing through an `Unreleasable` path;
- destroying through a non-owning path when ownership has not been validly transferred;
- finalizing or releasing before the resolved lifetime/finalization constraints permit it;
- generating two cleanup paths that can both destroy the same live resource;
- converting code in a way that introduces one of the above states.

A language converter, preprocessor, static analyzer, or compiler that can prove such invalidity must stop that stage with an error. Generated/transformed provenance is never an exemption.

Warnings remain appropriate only for suspicious states whose invalidity cannot yet be proven. If a particular Bitlang operation requires a proof of safety and that proof cannot be established, the operation is rejected according to that operation's safety contract rather than accepted merely because it originated from generated code.

Explicit low-level or unsafe operations do not grant permission to violate Bitlang's ownership/release/finalization invariants unless a separate specification explicitly defines a different destruction model. In particular, an unsafe pointer operation is not a general-purpose escape hatch for double destruction or invalid owner teardown.
