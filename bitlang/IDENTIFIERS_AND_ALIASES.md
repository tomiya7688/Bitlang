# Bitlang Unicode Identifiers and Declaration Aliases

Bitlang source supports Unicode identifiers, including Japanese variable, function, type/class, module, field, and parameter names.

Examples:

```bitlang
Int 回数;
function 実行(...)
class 利用者 { ... }
module ログ { ... }
```

Identifiers are source names, not string literals. String/character literal contents are never identifier-canonicalized.

## 1. Source encoding

Bitlang source is UTF-8.

Invalid UTF-8 in source text is a lexical error. The lexer must operate on Unicode scalar values/code points rather than treating non-ASCII bytes as independent characters.

## 2. Ordinary Unicode identifiers

Ordinary identifiers use the Unicode identifier profile based on UAX #31 XID_Start / XID_Continue, with `_` explicitly allowed as an identifier start/continue character.

This includes ordinary Japanese writing systems such as Kanji, Hiragana, and Katakana.

Conceptually:

```text
identifier-start    = '_' or XID_Start
identifier-continue = '_' or XID_Continue
```

An identifier may not begin with a decimal digit unless a future escaped-identifier form explicitly defines otherwise.

Bitlang language versions pin the Unicode data version used by identifier classification/canonicalization. Upgrading Unicode tables is a language-versioned semantic change rather than an implementation-dependent host-library change.

## 3. Canonical identifier comparison

Bitlang identifiers remain case-insensitive.

For comparison/name resolution, an identifier is canonicalized by:

```text
valid UTF-8 / Unicode scalar sequence
-> Unicode NFC normalization
-> Unicode default case folding
-> canonical identifier
```

`_` remains significant and is never inserted, removed, or ignored.

Japanese identifiers are unaffected by Latin capitalization but still participate in the same normalization pipeline.

Examples:

```text
Execute == execute
実行 == 実行
current_scope != currentscope
```

Canonically equivalent Unicode spellings resolve to the same identifier and therefore collide in one scope.

## 4. Confusable characters

Visually similar Unicode identifiers are not automatically merged merely because they look alike.

The compiler/tooling may emit a warning for suspicious Unicode confusables, mixed-script names, or other security-sensitive spelling, but such warnings do not change symbol identity unless a language-versioned rule explicitly says so.

## 5. Declaration aliases

Every named declaration may have zero or more aliases.

At minimum this includes:

```text
module
variable
field
parameter
function/method/constructor-like callable
type/class
```

Future named declaration kinds inherit the same model unless their specification says otherwise.

A declaration has:

```text
primary_name
aliases[]
one semantic symbol identity
```

Example concept:

```text
primary name: 実行
aliases: [execute, run]
```

References using `実行`, `execute`, or `run` all resolve to the same function symbol.

Aliases do not create wrapper functions, copied variables, duplicate types, or additional runtime objects.

Exact source spelling for attaching aliases to declarations is defined by the source grammar/parser; the semantic model is primary-name plus alias set.

## 6. Alias canonicalization and collisions

Every alias is canonicalized using exactly the same Unicode identifier rules as a primary name.

Within one name-resolution scope:

- an alias may not collide with another symbol's primary name;
- an alias may not collide with another symbol's alias;
- a symbol may not declare an alias canonically equal to its own primary name;
- duplicate aliases on the same symbol are errors after canonicalization;
- the same spelling may exist normally in an unrelated scope.

Collision checks use canonical forms while diagnostics preserve original spelling.

## 7. Primary identity versus lookup names

The primary canonical qualified name is the stable semantic/linkage identity of the declaration.

Aliases are additional lookup names for that same identity. They do not create additional backend symbols merely because multiple source names exist.

Conceptually:

```text
Game.実行       -- primary qualified name
Game.execute    -- alias lookup path

both -> same symbol identity
```

If a module itself has aliases, qualified lookup may resolve a module segment through one of its aliases without materializing every possible alias combination as separate declarations.

## 8. Visibility, export, import, and policy

Aliases inherit all restrictions of the declaration.

Using an alias does not bypass:

- scope visibility;
- Exported / Unexported;
- import capabilities;
- call-site policy;
- caller-class restrictions;
- ownership/borrow/access properties;
- any other semantic restriction.

If a declaration is not visible/exported/callable from a context, none of its aliases make it visible/exported/callable.

## 9. Name-based call-site policy

Name-based semantic policies must not depend on which alias spelling happened to be used at the call site.

For predicates such as:

```text
caller_function_name_exact
caller_function_name_prefix
caller_function_name_suffix
caller_function_name_glob
```

the matched name is the caller declaration's **primary canonical name** by default.

Adding an alias such as `log_fake` to a function whose primary name is `main` therefore does not make it satisfy a `log_*` caller policy.

If a future policy intentionally needs to inspect aliases, it must use an explicitly alias-aware predicate rather than changing the default meaning.

## 10. Bitlang Explicit boundary

Bitlang Explicit preserves:

- the primary canonical name;
- original spelling required for diagnostics/source mapping where available;
- the canonical alias set for named declarations when name resolution across Explicit units/interfaces still requires it;
- one shared symbol identity.

References should resolve to canonical symbol identity rather than retain whichever alias spelling was used in source.

After all cross-unit name resolution/linking that requires aliases is complete, aliases may be erased before Bitlang Low unless reflection/debug metadata explicitly requires them.
