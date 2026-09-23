# Bitlang Preprocessor Execution and Activation

This document defines the source-facing execution order, declaration timing, reuse rules, and activation ranges of Bitlang preprocessor functions.

It is the authoritative specification for these subjects. If a shorter summary elsewhere differs from this document, this document takes precedence.

## General model

Bitlang preprocessing proceeds while the preprocessor reads and analyzes Bitlang source structures.

Preprocessor functions are ordinary reusable functions in the preprocessing execution domain. They are not restricted to one-shot macro expansion.

A preprocessor function may be declared in an applicable declaration field/scope and then reused wherever that declaration is visible according to the normal preprocessing visibility rules.

The preprocessor normally executes applicable preprocessing behavior as the corresponding source structure is read.



## Multi-pass preprocessing and convergence

Bitlang preprocessing is not a single-pass transformation.

The normal model performs at least two preprocessing passes over the source. The first pass may discover declarations, execute preprocessing behavior, and modify the source representation. A later pass then reads the modified result again so newly generated or changed structures can themselves be analyzed and processed.

More precisely, preprocessing is iterative rather than being limited to exactly two passes.

After each pass, the preprocessor evaluates whether the current program has reached the `preprocess all ok` state. If that state has not been reached, preprocessing continues by applying the required fixes/transformations and running another pass over the resulting source representation.

Conceptually:

```text
initial source
    -> preprocess pass 1
    -> modified source
    -> preprocess pass 2
    -> check: preprocess all ok?
         yes -> finish preprocessing
         no  -> modify/fix
                -> preprocess again
                -> repeat check
```

Therefore "the preprocessor runs twice" describes the minimum normal multi-pass behavior, not a hard maximum. Processing continues until the source reaches a stable valid preprocessing result or preprocessing terminates with an error.

A primary reason for this convergence requirement is compatibility cleanliness at the boundary to Bitlang Preprocessed. The preprocessor must not leave obsolete, legacy, or otherwise outdated Preprocessed-era representations mixed into the final normalized output. If an older Bitlang Preprocessed form remains after an earlier transformation pass, a later preprocessing pass must detect and rewrite or reject it before the result is handed to the Bitlang compiler.

The Bitlang compiler is therefore allowed to assume that input crossing the preprocessing boundary has already been normalized to the currently accepted Bitlang Preprocessed form. Backward-cleanup of stale Preprocessed syntax is a preprocessing responsibility rather than something the compiler must guess around.

In many Bitlang Preprocessed version transitions, the change is primarily the addition of newly explicit semantic properties rather than a change in the underlying program meaning. A property that older Preprocessed versions treated implicitly may become a required explicit property in the current version. Such older input is therefore no longer fully canonical Preprocessed input even if its surface syntax resembles it closely.

For migration purposes, stale Preprocessed input of this kind is treated conceptually like Bitlang source that happens to wear a Preprocessed-like syntax: omitted semantics are reconstructed by preprocessing, the newly required properties are materialized, and only the upgraded fully explicit form may cross into the compiler.

Each additional pass observes the result of the previous pass. Generated declarations, changed properties, inserted cleanup, macro expansion, and other preprocessing transformations may therefore participate in subsequent preprocessing passes.

The default source-order execution rules still apply within each pass unless an explicit preprocessing step, ordering rule, or execution constraint changes them.

A preprocessing implementation must detect a non-converging transformation cycle or another condition that prevents reaching `preprocess all ok`; such a state must not result in an unbounded silent loop.

## Activation ranges

A preprocessor function may be associated with an explicit start point and an explicit end point by referring to the preprocessor function together with start/end activation markers.

The start marker establishes the point from which that preprocessor function applies.

The end marker establishes the point after which that preprocessor function no longer applies.

Conceptually:

```text
Preprocess_function + start
    ... source affected by Preprocess_function ...
Preprocess_function + end
```

The exact surface spelling of the start/end markers is specified separately. The semantic requirement is that the markers identify the preprocessor function and delimit its active source range.

If an explicit end marker is omitted, the activation range ends automatically at the end of the declaration field that contains the start point. The activation does not implicitly escape into an enclosing or following field.

Examples of the default boundary are:

- a start point in a function field applies through the remainder of that function field;
- a start point in a class field applies through the remainder of that class field;
- a start point at file-level applies through the remainder of that file-level field;
- a start point in another lexical/declaration field applies through the end of that field.

An explicit end marker may terminate the activation earlier than that default field boundary.

Because the target is a function, the same declared preprocessor function may be activated more than once in different declaration fields or source ranges where it is visible.

## Default execution order

Unless an explicit preprocessing step, execution-order rule, condition, constraint, or other preprocessing control changes the order, preprocessing follows analyzed source order.

The default model is:

```text
analyze source structure
    -> encounter applicable preprocessing declarations/rules
    -> process them in the order the main source is read
```

Therefore ordinary textual/structural read order is the default sequencing rule, not an unspecified scheduler.

Explicit preprocessing-step definitions or execution constraints may override this default. Such overrides must be represented explicitly rather than silently changing the order.

## Declaration position and read timing

The declaration position determines when a preprocessor declaration becomes visible to the preprocessing walk.

### Outside a class declaration

A preprocessor declaration outside a class is read when that file is read.

Its declaration therefore enters preprocessing at the file-reading stage corresponding to its source position, subject to any explicit preprocessing-step or ordering rule.

### Inside a class declaration

A preprocessor declaration inside a class is read when that class is read.

Class-local preprocessing declarations are read at the beginning of preprocessing that class so that they are available while the class body is subsequently preprocessed.

This does not make them runtime class members. They remain preprocessing-domain declarations.

### Inside a function declaration

A preprocessor declaration inside a function is read when that function is read/analyzed.

However, its execution follows the preprocessing walk through that function body. The function body is preprocessed from top to bottom, and preprocessing behavior executes when the walk reaches the corresponding applicable point.

Conceptually:

```text
read function declaration
    -> discover function-local preprocessing declarations
    -> preprocess function body from top to bottom
    -> execute applicable preprocessing behavior as its point is reached
```

Thus "the declaration is known when the function is read" and "its preprocessing effect executes at its source position while the function body is walked" are distinct rules.

## Reuse

Preprocessor functions are functions and may be reused within declaration fields/scopes where they are visible.

Reusing a preprocessor function does not duplicate its declaration. Multiple activation ranges, calls, or preprocessing applications may refer to the same function declaration.

Normal function arguments and preprocessing-visible values may be used to parameterize reusable preprocessing behavior where the function signature permits it.

## Relationship to ordinary calls

Preprocessor functions use ordinary Bitlang function-call syntax when called directly.

Activation markers are different from an ordinary call: they establish an application range for preprocessing behavior rather than merely performing a single immediate call expression.

A directly called preprocessor function executes as a preprocessing-domain function at that point in the preprocessing walk.

A range-activated preprocessor function applies according to its defined preprocessing behavior while its activation range remains active.

## Explicit ordering and conditions

Bitlang may provide separate mechanisms for:

- named preprocessing steps;
- explicit execution ordering;
- execution conditions;
- activation constraints;
- dependency constraints between preprocessing operations.

When none of those mechanisms are specified, source read order remains the default.

These controls do not change the source meaning of `@preprocesser`; they only control when an otherwise valid preprocessing operation is applied.

## Preprocessed boundary

Activation markers, preprocessing execution-order controls, preprocessing function declarations that exist only for preprocessing, and other preprocessing-only scheduling information are consumed before Bitlang Preprocessed is emitted.

Only the resulting transformed and normalized Bitlang program proceeds to Bitlang Preprocessed.
