# Bitlang Source Lexicon

Bitlang source is intentionally writable in English and Japanese, and the source lexicon is extensible to additional languages.

Japanese support is not limited to Unicode user-defined identifiers. Reserved words, built-in source names, source-facing properties, literals, and other ordinary language vocabulary may have English and Japanese source spellings that normalize to the same canonical semantic token.

## 1. Equal source-language status

English and Japanese spellings are both first-class Bitlang source syntax. Version 1 requires `en` and `ja`; additional registered language aliases may be added under versioned language codes without changing the canonical Explicit vocabulary.

Language codes other than `en` and `ja` are optional extensions of the source lexicon. Every entry still has exactly one canonical English spelling.

At the Bitlang source boundary, Japanese is not implemented as a separate source-to-source language adapter and English is not privileged for source acceptance. Both forms are first-class source spellings.

However, Bitlang Explicit has one canonical **English language vocabulary**. Both source forms normalize directly into the same parser/preprocessor concepts, and preprocessing emits the English canonical keyword/property/built-in spelling when materializing Bitlang Explicit.

Conceptually:

```text
class      -> keyword.class
クラス     -> keyword.class

function   -> keyword.function
関数       -> keyword.function

module     -> keyword.module
モジュール -> keyword.module
```

After lexical normalization, later semantic stages do not care which source-language spelling was used. When serialized as Bitlang Explicit, the canonical language token spelling is English.

## 2. Mixed-language source

English and Japanese source vocabulary may be mixed in one file.

Example:

```bitlang
module Game {
    クラス 利用者 {
        Str 名前;

        関数 execute() {
            // ...
        }
    }
}
```

Mixing languages is a style choice, not a semantic distinction.

## 3. Japanese completeness rule

A source-facing canonical language token is not considered fully specified until its Japanese source spelling is defined, unless the token is language-neutral punctuation/symbol syntax or explicitly documented as intentionally untranslated.

This requirement applies as the grammar grows.

At minimum the source lexicon covers:

- declaration keywords;
- control-flow keywords;
- module/import vocabulary;
- type/class/struct/interface vocabulary;
- built-in primitive/container type names where source aliases are appropriate;
- source-facing property states and common shorthands;
- boolean/null-like literals;
- ownership/borrow/move source vocabulary;
- preprocessing/source directives intended for ordinary source authors.

CI/spec consistency tooling should detect when a newly introduced canonical source token lacks the required Japanese spelling.

## 4. Initial canonical aliases

The initial lexicon includes obvious stable declarations:

```text
module    <-> モジュール
class     <-> クラス
struct    <-> 構造体
interface <-> インターフェース
function  <-> 関数
return    <-> 返す
import    <-> インポート
```

Visibility/source-property conveniences include:

```text
Public  <-> 公開
Private <-> 非公開
Static  <-> 静的
Dynamic <-> 動的
Readable   <-> 読取可能
Unreadable <-> 読取不可
Writeable   <-> 書込可能
Unwriteable <-> 書込不可
Copyable   <-> コピー可能
Uncopyable <-> コピー不可
Movable   <-> 移動可能
Unmovable <-> 移動不可
```

Built-in source type aliases may include:

```text
Bool <-> 真偽
Str  <-> 文字列
Char <-> 文字
```

Numeric canonical type families such as `Int` / `Uint` / `Float` retain their canonical technical names, but source lexicon aliases may additionally provide Japanese spellings when they can be defined without ambiguity. Exact numeric-type Japanese spellings are versioned in the machine-readable lexicon rather than inferred from natural-language translation.

## 5. No semantic translation by natural language

The compiler does not dynamically translate arbitrary Japanese words into English keywords.

Only explicitly registered, versioned lexicon aliases are reserved language vocabulary.

For example, a synonym that humans may consider equivalent is still an ordinary identifier unless it appears in the lexicon.

This keeps parsing deterministic and prevents natural-language ambiguity.

## 6. User identifiers remain multilingual

User-defined identifiers may independently be English, Japanese, or another permitted Unicode identifier.

```bitlang
関数 実行()
function execute()
クラス 利用者
class User
```

Declaration aliases can additionally make one symbol available under multiple names, for example `実行` and `execute`, as specified by `IDENTIFIERS_AND_ALIASES.md`.

Keyword aliases and declaration aliases are different mechanisms:

- source lexicon alias: two reserved spellings mean the same language token;
- declaration alias: multiple user-visible names resolve to the same semantic symbol.

## 7. Canonicalization boundary

Lexicon normalization happens before ordinary parsing/semantic normalization needs to distinguish source-language spelling.

Diagnostics and source maps preserve the original spelling and location.

Bitlang Explicit uses English-only canonical language vocabulary. For example:

```text
Bitlang source                  Bitlang Explicit
クラス                           class
関数                             function
モジュール                       module
公開                             canonical English visibility/property spelling
文字列                           Str
```

Therefore Bitlang Explicit never needs to preserve whether the source language token was written as `class` or `クラス`; preprocessing emits the English canonical language token.

This English-only rule applies to **language vocabulary**, not arbitrary user-defined symbol names. A user identifier such as `利用者` or `実行` remains a valid Unicode identifier in Explicit unless ordinary declaration-alias/name-normalization rules explicitly choose another symbol name. The compiler must not machine-translate arbitrary user identifiers into English.

## 8. Two normalization outputs

The same multilingual source lexicon resolver supports two distinct outputs.

### A. Multilingual Bitlang -> English Bitlang source

This is **source-to-source lexicon normalization only**.

It replaces registered non-English language-vocabulary spellings with each entry's canonical English spelling while preserving ordinary Bitlang source semantics and source-level convenience.

It does **not**:

- fill omitted properties;
- expand source sugar merely because Explicit would require expansion;
- resolve all defaults;
- force fully qualified references;
- turn the program into Bitlang Explicit.

For example:

```text
公開 クラス 利用者 {
    関数 実行() {
        文字列 名前;
    }
}

->

Public class 利用者 {
    function 実行() {
        Str 名前;
    }
}
```

User-defined identifiers, comments, string/character literals, and other non-lexicon text are not translated into English.

The result is still ordinary Bitlang source and may still contain omitted properties, source sugar, preprocessing directives, and other source-only constructs.

### B. Multilingual Bitlang -> Bitlang Explicit

This is the full preprocessing/normalization path.

It first applies the same source-lexicon resolution semantics, then performs all required Bitlang normalization: property/default resolution, source-sugar expansion, module/name normalization, policy/state materialization, and every other transformation required by the Explicit boundary.

The implementation may perform this directly in memory and does not need to serialize an intermediate English Bitlang source file.

The following must be semantically equivalent:

```text
multilingual Bitlang
    -> English Bitlang source
    -> full preprocess
    -> Bitlang Explicit

multilingual Bitlang
    -> direct full preprocess
    -> Bitlang Explicit
```

The English-source conversion and Explicit conversion are therefore separate user-visible operations built on one lexicon resolver.

## 9. Formatter and generated-source preference

The compiler/formatter/tooling may have a source lexicon preference such as:

```text
en
ja
preserve
```

This preference controls emitted/generated/rewritten source spelling only. It does not restrict what the parser accepts.

`preserve` keeps original source spelling where possible.

## 10. Versioning

The lexicon is versioned together with the Bitlang language specification.

Adding/removing/reassigning a reserved spelling can change whether text parses as a keyword or identifier and therefore follows Bitlang's language-version compatibility rules.
