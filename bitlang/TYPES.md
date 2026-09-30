# Bitlang Type System

## Type-name representation

Bitlang may encode representation information directly in the type name.

For types that have both a radix/base and a bit width, the canonical naming pattern is:

```text
<TypeName><Radix>x<BitWidth>
```

Example:

```text
Int10x32
```

means an `Int` represented in radix 10 with a 32-bit width.

A source-level short type name uses that type's default representation. Therefore:

```text
int value
```

is equivalent, after preprocessing, to:

```text
Int10x32 value
```

## Source byte-width notation

Bitlang source may specify a semantic width in **bytes** by using `xx` instead of `x`.

For types where the canonical `x<N>` suffix denotes a semantic bit width:

```text
<TypeName><Radix>xx<ByteWidth>
```

is source-level shorthand for:

```text
<TypeName><Radix>x<ByteWidth * 8>
```

A Bitlang byte in this notation is exactly 8 semantic bits. It does not depend on a C target's `CHAR_BIT`, ABI, carrier type, or physical storage unit.

Examples:

```text
Int2xx4      -> Int2x32
Int10xx4     -> Int10x32
Uint16xx8    -> Uint16x64
Float10xx4   -> Float10x32
```

`xx` does not create a distinct canonical type identity.

Therefore, after preprocessing:

```text
Int2xx4 == Int2x32
```

in the type system.

The byte-width form is source sugar only. Bitlang Explicit always uses the canonical `x<BitWidth>` representation for bit-width-bearing types, so `xx` must not survive the source-to-Explicit normalization boundary.

The ordinary `x` form remains available for widths that are not whole-byte multiples.

The `xx` shorthand is especially useful when source code is naturally specified in byte capacities, such as binary buffers, packet fields, fixed storage blocks, and backing storage for byte-limited string implementations.

## Signed and unsigned integers

`Int` is signed.

`Uint` is unsigned.

Examples:

```text
Int10x32
Uint10x32
```

## Overflow

Numeric overflow is an error by default.

Bitlang must not silently wrap an overflowing value unless a separate operation explicitly requests different overflow behavior.

## Radix-aware numeric forms

Numeric types may use other radix values while preserving the same bit width.

Examples:

```text
Int2x32
Int8x32
Int10x32
Int16x32
```

Radix is part of the numeric type representation, not merely display formatting.

Values with different radix types cannot be used together directly in arithmetic or bitwise operations. One side must first be converted to the same radix through an appropriate Bitlang pulse function or preprocessor pulse function.

For example, an `Int2x32` and an `Int10x32` are not directly compatible operands even when they represent the same numeric value and bit width.

Bitwise operations are performed over the fixed-width bit representation after operand types have been made compatible. Radix-2 exposes individual bits directly, radix-8 groups bits in sets of three, and radix-16 provides a compact bit-oriented form.

## Shift operations

Bitlang distinguishes **bit-representation shifting** from **radix/digit shifting**.

These are different semantic operations and must not be silently substituted for one another.

### Bit shift

Canonical operation family:

```text
bit_shift_left(value, count)
bit_shift_right_zero(value, count)
bit_shift_right_sign(value, count)
```

A bit shift moves the fixed-width bit representation itself.

Conceptually:

```text
00010
bit_shift_left by 1
-> 00100
```

Rules:

- the result keeps the same numeric type and semantic bit width;
- `bit_shift_left` shifts zero bits in from the right;
- `bit_shift_right_zero` shifts zero bits in from the left;
- `bit_shift_right_sign` preserves the sign bit and is valid only for signed integer types whose canonical signed-bit representation is defined;
- information loss is not implicit: a checked form rejects loss of significant bits, while a discard form explicitly permits bits leaving the semantic width to be discarded;
- `count` is a non-negative integer bit count and is not required to have the same radix as `value`;
- ordinary bit-shift count must satisfy `0 <= count < bit_width`;
- an invalid constant count is a compile-time error;
- an invalid count known only at runtime uses the ordinary Bitlang runtime error/trap path.

Bit-shift semantics are defined by Bitlang's canonical semantic bit representation, not by C's implementation-specific signed-shift behavior.

### Radix shift

Canonical operation family:

```text
radix_shift_left(value, count)
radix_shift_right(value, count)
```

A radix shift moves the value by digits of the value's own radix rather than by individual bits.

For a type with radix `R`:

```text
radix_shift_left(value, n)
    = value * R^n

radix_shift_right(value, n)
    = value / R^n
```

For integer types, right radix shift uses the language's ordinary integer-division result and therefore truncates toward zero.

Examples:

```text
Int10x32: 10 radix_shift_left 1 -> 100
Int16x32: 0x10 radix_shift_left 1 -> 0x100
Int8x32:  010 radix_shift_left 1 -> 0100
```

Rules:

- the result keeps the same type and radix;
- `count` is a non-negative integer digit count;
- overflow behavior is selected explicitly by the standard-library operation; checked and wrapping forms are distinct;
- right radix shift does not use C implementation-defined signed shift behavior;
- radix shift does not silently choose between checked, wrapping, or other overflow behavior;
- a language adapter may map source-language digit/scale operations to radix shift when their semantics match.

For radix 2, a left radix shift and a left bit shift may produce the same value while no significant bit is discarded, but they remain distinct operations: bit shift explicitly manipulates the fixed-width representation, while radix shift performs checked numeric scaling.

Source-facing shift operations are provided through the standard library rather than by attaching additional semantic properties to shift operators.

Representative APIs are:

```text
bit.shift_left_checked(value, count)
bit.shift_left_discard(value, count)
bit.shift_right_zero_checked(value, count)
bit.shift_right_zero_discard(value, count)
bit.shift_right_sign_checked(value, count)
bit.shift_right_sign_discard(value, count)

radix.shift_left_checked(value, count)
radix.shift_left_wrapping(value, count)
radix.shift_right(value, count)
```

The selected loss/overflow behavior must be explicit in the function used. Compiler lowering may convert these calls to canonical shift intrinsics. Language adapters may map foreign shift operators directly to the matching intrinsic semantics.

## Literals and type inference

Literals use ordinary type inference unless an explicit type is supplied.

The inferred type must ultimately resolve to a concrete Bitlang type before canonical preprocessing is complete.

## Floating-point types

Floating-point types also carry a radix component.

If the radix is omitted, radix 10 is the default.

The same general naming principle applies:

```text
<TypeName><Radix>x<BitWidth>
```

## Types without radix

If radix/base is not meaningful for a type, the type omits the radix component and uses only the relevant size/count information.

## Strings

`Str` uses its numeric suffix as a maximum character count rather than as a bit width.

A plain:

```text
Str
```

has no maximum character limit by default.

A bounded string form specifies a maximum number of characters explicitly.

`Char` is source-level shorthand for a one-character string representation and normalizes to `Str1x1`.

The first `1` in `Str1x1` is currently reserved and does not yet have a finalized semantic meaning. It may be assigned a useful string-representation property later.

The current `Str` maximum-character-count suffix is a separate semantic convention from numeric bit-width notation. The new `xx` byte-width shorthand therefore does not automatically redefine `Str` suffix meaning. Byte-limited strings may use fixed-byte backing storage built from width-bearing types; a direct `Str...xxN` surface form requires a separate explicit string-type decision.

## Generics

Bitlang generics are compile-time type parameterization.

A generic type parameter represents a type, not a runtime value. Supplying a generic argument is not a cast and does not change the type of an existing value.

Generic parameters are resolved during preprocessing. Bitlang Explicit must not contain unresolved generic type parameters.

The same generic declaration instantiated with the same type arguments denotes the same specialization and may be reused rather than regenerated independently.

Generic parameters may be declared on:

- classes;
- structs;
- interfaces;
- functions.

Conceptually:

```text
class Box<T> { ... }
struct Pair<T> { ... }
interface Comparable<T> { ... }
T identity<T>(T value) { ... }
```

A generic function is therefore a compile-time family of functions parameterized by type. After preprocessing, each used specialization has concrete types.

Generic language differences are represented through properties wherever they affect semantics. A language adapter may therefore select generic-parameter properties that preserve the source language's behavior instead of redefining a separate generic system.

Only semantic differences belong in these properties. Pure implementation choices that do not change observable program meaning, such as whether a specialization is emitted as duplicated machine code or shared internally, are compiler/lowering decisions rather than source semantic properties.

The generic declaration is fully resolved during preprocessing. Any generic-specific source/preprocessing properties that still matter must be materialized into the resulting concrete declarations and types before Bitlang Explicit.

The exact property axes, constraints, and advanced generic syntax are defined separately in [GENERICS.md](GENERICS.md).

## Arrays

Bitlang source may describe array dimensionality and fixed lengths as attributes inside the `Array<...>` form rather than requiring the programmer to write repeated nested array constructors manually.

Conceptually:

```text
Array<T, Dimension=3>
```

means a three-dimensional array of `T`.

Fixed lengths may also be supplied in the array attributes. When more than one dimension is present, a fixed length may be supplied for each dimension.

Conceptually:

```text
Array<T, Dimension=2, Length=10,20>
```

represents a two-dimensional array whose dimensions have fixed lengths 10 and 20.

These source-level dimension and length attributes are convenience information. The preprocessor expands dimensionality into nested canonical `Array` types. Fixed-length information is attached to the corresponding canonical array level rather than retained as a separate multidimensional-array abstraction.

Therefore programmers may write dimensionality compactly, while Bitlang Explicit uses ordinary nested arrays internally.

## Pointer and reference types

`Ptr<T>` and `Ref<T>` are distinct types with different safety guarantees.

### Ptr<T>

`Ptr<T>` is the low-level raw-pointer type.

It may:

- contain `null`
- be reassigned
- participate in pointer arithmetic
- expose and manipulate raw addresses where the target platform permits it

Because `Ptr<T>` is intentionally low-level, code using it is responsible for avoiding invalid addresses, dangling pointers, and other unsafe memory access.

### Ref<T>

`Ref<T>` is the safe reference type.

It:

- cannot contain `null`
- does not permit pointer arithmetic
- does not expose arbitrary raw-address manipulation as part of normal reference operations
- must not outlive the value or object it references
- cannot be rebound after its initial binding

The compiler and static-analysis stages should reject references whose lifetime is known to exceed the referenced value's lifetime.

Assignment through a `Ref<T>` may modify the referenced value when that value is writable, but it does not change which value the reference is bound to.

For example, after binding a reference to `a`, assigning a new reference target such as `Ref(b)` is invalid. Ordinary assignment through that reference may still write to `a` when permitted by the referenced value's mutability rules.

## Nullability and presence properties

Nullability and presence are separate semantic properties.

The nullability axis is:

```text
nullable
unnullable
```

`nullable` means the value may be `null`. `unnullable` means `null` is not a valid value. These are the canonical property names; the source wrapper `Nullable<T>` is separate source syntax.

The presence axis is:

```text
Optional
Required
```

`Optional` means the value or declaration may be absent. `Required` means it must be present.

These axes are independent. For example, a value may be `Required nullable`, meaning it must exist but may contain `null`, or `Optional unnullable`, meaning it may be absent but, when present, may not be `null`.

Bitlang Explicit must make both properties explicit whenever they apply so that absence and nullability are never inferred from omission.

## Explicit type conversion

Bitlang has no ordinary implicit type conversion. Different types remain incompatible until the program explicitly performs an appropriate conversion.

### Cast

A cast changes the value's type.

It is used when the semantic type itself changes, for example when converting between integer and floating-point types or between other distinct type families.

A cast may fail when the requested target type cannot represent the value. Overflow remains an error.

### Pulse

A pulse produces a converted representation while preserving the original source variable.

Pulse is intended for representation changes such as radix conversion where the source variable must remain unchanged.

For example, an `Int10x32` may be pulsed to an `Int2x32` value for an operation that requires radix 2, while the original `Int10x32` variable remains `Int10x32` and retains its value.

Both ordinary Bitlang code and preprocessor functions may perform cast and pulse operations.

### No automatic conversion

Bitlang itself does not automatically widen, narrow, change radix, change signedness, unwrap wrappers, or otherwise convert a value merely to make an operation type-compatible.

If two operands have different types, they must first be made compatible explicitly.

The preprocessing system may provide explicit opt-in automation rules. For example, a preprocessor function or attribute may mark a declaration as automatically widenable or otherwise permit a specific safe conversion. Such behavior is generated preprocessing logic, not a built-in implicit-conversion rule of Bitlang.

Any automatically generated conversion must be explicit in Bitlang Explicit output so that canonical semantics remain unambiguous.

## Preprocessing rule

Source-level shorthand is for convenience only. Bitlang preprocessing resolves shorthand, inferred types, defaults, compact array-dimension notation, and explicitly configured conversion automation into the concrete canonical type representation and explicit operations used by Bitlang Explicit.
