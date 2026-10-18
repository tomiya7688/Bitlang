# Bitlang Value Semantics

Bitlang uses value semantics by default for ordinary typed values.

## Core rule

Writing a type `T` means the value itself.

```text
T       -> value semantics
Ref<T>  -> explicit safe reference semantics
Ptr<T>  -> explicit raw-pointer semantics
```

This rule applies uniformly to structs, classes/types, arrays, numeric values, strings, and other ordinary value types unless a specific type constructor explicitly denotes indirection.

Bitlang does not use a language rule such as `class = reference type` and `struct = value type`.

## Parameters

A parameter declared as `T` receives a value.

```text
function f(T value)
```

The callee does not receive an implicit reference to the caller's storage.

To pass access to the caller's storage, the source must say so explicitly:

```text
function f(Ref<T> value)
function f(Ptr<T> value)
```

`Borrowed` ownership/property state does not silently turn plain `T` into a reference to the caller variable. If storage aliasing is required, the type itself must contain explicit indirection such as `Ref<T>` or `Ptr<T>`.

## Assignment and initialization

Ordinary assignment/initialization of plain `T` has value semantics.

For a copyable source value, assigning or passing an lvalue copies the value.

If the type/value is `Uncopyable`, an operation that would require copying is a compile-time error. The compiler must not fall back to reference semantics.

Ownership transfer of a movable value is a distinct move operation. A move still transfers a `T` value; it does not change the type into a reference.

Conceptually:

```text
a = b          -> value copy, requires Copyable where a copy is required
a = move(b)    -> value move/ownership transfer, requires Movable
a = Ref(b)     -> explicit reference value, type Ref<T>
```

Exact move surface syntax is defined separately; the semantic distinction is mandatory.

## Return values

A function returning `T` returns a value.

```text
function make() -> T
```

A function returning `Ref<T>` or `Ptr<T>` explicitly returns indirection and must satisfy the corresponding lifetime/safety rules.

## Classes and structs

Classes and structs follow the same default transfer rule:

```text
struct S  + parameter T -> value
class  C  + parameter T -> value
```

A class may be `Uncopyable`. In that case passing an existing lvalue class object to a by-value parameter is rejected unless an explicit valid move is performed. It is never silently converted into `Ref<C>`.

If value-copying a class creates a distinct observable object identity, the copied value has its own identity according to that class/type's copy semantics. A type that cannot define a valid independent copy should be `Uncopyable`.

Fields inside a copied value retain their own declared semantics. For example, a `Ref<U>` field remains a reference value and may still refer to the same `U`; value semantics of the outer object do not deep-copy referents unless the type's explicit copy semantics require that.

## Arrays and containers

`Array<T>` contains `T` values. An array of references is written explicitly as `Array<Ref<T>>`; an array of raw pointers is `Array<Ptr<T>>`.

No implicit element-reference model is inferred from whether `T` is a class.

## Backend optimization

Semantic value passing is independent from the physical ABI.

A backend may use hidden addresses, registers, stack slots, copy elision, return-value optimization, or another target-specific mechanism to transport large values efficiently, provided the observable result is exactly the same as value semantics.

In particular, a backend-internal pointer must not allow the callee to mutate caller storage when the source contract is by-value.

## Foreign-language adapters

Foreign reference semantics are made explicit when translated into Bitlang.

Examples:

```text
C struct by value        -> T
C pointer                -> Ptr<T>
C++ reference            -> Ref<T> / appropriate Borrowed reference
Rust &T                  -> Borrowed Ref<T>
Rust owned value         -> T with copy/move semantics
Java/C# class reference  -> explicit reference/handle representation, not plain T
```

This preserves source-language semantics without weakening Bitlang's rule that plain `T` means a value.

## No silent indirection

The compiler, preprocessor, optimizer, language adapter, and backend must never change a plain-value semantic contract into observable shared-reference semantics merely for convenience or performance.
