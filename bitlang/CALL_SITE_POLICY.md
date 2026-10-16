# Bitlang Call-Site Policy

Status: current design

Bitlang supports compile-time restrictions on **which call sites are allowed to invoke a callable declaration**.

This is independent from visibility and import capability:

- visibility answers whether the declaration may be resolved from a context;
- import capability can reduce whether an imported callable may be invoked;
- call-site policy answers whether this particular caller is permitted to perform the call.

A declaration may therefore be visible and imported as Callable while a specific call is still rejected by its call-site policy.

## 1. Canonical property axis

Callable declarations use:

```text
Call_site_unrestricted
Custom_call_site_policy
```

The canonical target is `function` (including methods and constructor-like callable declarations after normalization).

`Call_site_unrestricted` places no additional caller restriction.

`Custom_call_site_policy` carries a versioned, pure compile-time call-policy expression.

Source default:

```text
call_site_policy -> Call_site_unrestricted
```

## 2. Function-name restrictions

A call policy may restrict the caller by its canonical function name.

Supported deterministic name-match forms include:

```text
caller_function_name_exact(<name>)
caller_function_name_prefix(<prefix>)
caller_function_name_suffix(<suffix>)
caller_function_name_glob(<glob>)
```

Matching is performed after ordinary Bitlang identifier canonicalization:

- identifier case is not significant;
- `_` remains significant;
- the caller declaration's primary canonical function name is matched, not source spelling preserved only for diagnostics.

For example, conceptually:

```text
Caller_function_name_matches("log_*")
function output_log(...)
```

may normalize to a custom policy equivalent to:

```text
caller_function_name_glob("log_*")
```

Then callers such as `log_error` and `log_request` are permitted, while `main` or `handle_request` are rejected.

Aliases do not affect these default name predicates. Adding an alias `log_main` to a function whose primary name is `main` does not make it satisfy a `log_*` policy. Alias-aware matching requires an explicitly separate predicate if one is ever introduced.

The exact source shorthand spelling may remain ergonomic; Bitlang Explicit receives the resolved canonical expression.

## 3. Glob semantics

The canonical glob language is intentionally small and deterministic.

For local function-name matching:

- `*` matches zero or more canonical identifier characters;
- `?` matches exactly one canonical identifier character;
- all other characters match literally after identifier canonicalization;
- no regular-expression execution is embedded in the compiler;
- invalid pattern syntax is a compile-time error.

A language adapter that needs a different source-language pattern mechanism must normalize it to equivalent standard predicates or another explicit call-policy relation before Bitlang Explicit.

## 4. Other caller relations

The expression model may also restrict by already-resolved caller identity/scope:

```text
caller_function(<canonical-symbol>)
caller_owner_type(<canonical-type>)
caller_module(<canonical-module>)
caller_project(<project-id>)
caller_adapter_domain(<adapter-id>, <domain-id>)
```

Boolean composition:

```text
any(...)
all(...)
not(...)
```

Example:

```text
all(
    caller_function_name_prefix("log_"),
    caller_module(Game.Logging)
)
```

This permits log functions with the required naming convention only when they are also inside the selected module.

## 5. Logging use case

A low-level log-output function may deliberately be callable only from functions following the logging naming convention.

Conceptually:

```text
Custom_call_site_policy(
    caller_function_name_glob("log_*")
)
function output_log(...)
```

This can enforce architecture such as:

```text
application code
    -> log_info / log_error / log_request
    -> output_log
```

while rejecting:

```text
main
    -> output_log        ERROR

handle_request
    -> output_log        ERROR
```

The restriction is compile-time only and introduces no runtime caller-name lookup.

## 6. Indirect calls and function values

Call-site policy follows the target callable and cannot be bypassed by copying or passing a function value.

For:

```text
restricted_function -> function value -> indirect call
```

the compiler validates the **ultimate invocation site** against every reachable restricted target.

If static analysis cannot prove that every possible target permits the caller, the indirect call is rejected unless an explicitly specified safe dispatch mechanism has already resolved the policy.

Taking a function value without invoking it is not itself a call-site-policy violation, but capability/visibility restrictions may independently restrict obtaining that value.

## 7. Interaction with wrappers

A wrapper does not inherit permission merely because it was called by an allowed function.

For example, if `output_log` requires caller name `log_*`:

```text
log_error -> generic_emit -> output_log
```

is valid only if `generic_emit` itself satisfies the policy (or the policy explicitly allows it).

The policy checks the direct semantic caller of the invocation. Bitlang does not inspect the runtime stack and does not grant transitive permission from an ancestor caller.

If transitive/trusted-wrapper behavior is required, it must be represented explicitly by an allowed caller symbol, name pattern, module, adapter domain, or another defined policy relation.

## 8. Interaction with class caller limits

This feature is independent from the existing class/type caller-count restriction:

```text
Single_caller_class / Caller_class_unrestricted
```

The class property limits the number of distinct caller classes.

Call-site policy filters whether an individual invocation from a particular caller is permitted.

Both restrictions apply when relevant.

## 9. Interaction with visibility/imports

Permission is intersected:

```text
visibility
AND export
AND import capability
AND call-site policy
AND caller-class/count restrictions
-> valid call
```

No policy can widen a restriction imposed by another applicable axis.

An import marked `Callable` therefore cannot override a call-site policy.

## 10. Backend boundary

Call-site policy is compile-time semantic information.

After all direct/indirect calls have been validated and lowered, the policy may be erased before Bitlang Low unless a separate reflection/runtime metadata feature explicitly requires it.

No Assam Core or Bitlang VM instruction is added for caller-name checking.
