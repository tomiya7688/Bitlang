# Basic Source Syntax

This document defines the first C#-style declaration surface for Bitlang source. It describes source syntax only. The strict Preprocessed representation remains the semantic boundary, and source declarations must retain enough information for later property expansion and lowering.

## Lexical conventions

- This bootstrap syntax accepts ASCII identifiers only: the first character is an ASCII letter or `_`, and later characters may also be ASCII digits. Other non-ASCII bytes outside comments and quoted literals are rejected.
- Identifiers are compared case-insensitively; their original spelling and source location are retained. Unicode identifier normalization and case-folding are not defined.
- Lines are 1-based and separated by LF. Columns are 1-based UTF-8 byte positions, matching the lexer input offsets; a multibyte character inside a quoted literal advances the following token's column by its byte length.
- Reserved grammar words such as `class`, `public`, and `private` are emitted as keyword tokens. Their original spelling and source location remain intact; keyword classification must not rewrite or discard source text.
- `{`, `}`, `(`, `)`, `[`, `]`, `;`, `,`, `.`, and `=` are structural symbols in this syntax version.
- String and character literal contents are preserved as opaque tokens and are not canonicalized. UTF-8 validity, escape validity, and character-scalar count are checked by a later semantic stage, not by this declaration lexer.

## Declaration grammar

The following EBNF defines the initial declaration shape. `ExpressionTokens` and `StatementTokens` are balanced token sequences whose internal grammar is intentionally not defined here.

```ebnf
SourceFile          ::= { ClassDeclaration } EOF
ClassDeclaration    ::= ModifierList "class" Identifier [ GenericParameterList ] "{" { MemberDeclaration } "}"
MemberDeclaration   ::= FieldDeclaration
                      | PropertyDeclaration
                      | ConstructorDeclaration
                      | MethodDeclaration
FieldDeclaration    ::= ModifierList Type Identifier [ "=" ExpressionTokens ] ";"
PropertyDeclaration ::= ModifierList Type Identifier "{" AccessorDeclaration { AccessorDeclaration } "}"
AccessorDeclaration ::= ModifierList AccessorKind ";"
AccessorKind         ::= "get" | "set" | "init"
ConstructorDeclaration
                    ::= ModifierList Identifier "(" [ ParameterList ] ")" Block
MethodDeclaration   ::= ModifierList Type Identifier [ GenericParameterList ] "(" [ ParameterList ] ")" Block
ParameterList       ::= Parameter { "," Parameter }
Parameter           ::= ModifierList Type Identifier
Type                 ::= QualifiedName [ TypeArgumentList ]
QualifiedName        ::= Identifier { "." Identifier }
GenericParameterList ::= "<" Identifier { "," Identifier } ">"
TypeArgumentList     ::= "<" Type { "," Type } ">"
ModifierList        ::= { Modifier }
Modifier            ::= "public" | "private" | "protected" | "internal"
                      | "static" | "readonly" | "const" | "abstract"
                      | "sealed" | "virtual" | "override" | "Direct"
Block               ::= "{" StatementTokens "}"
```

A constructor is recognized when a member name matches its containing class name under Bitlang's case-insensitive identifier comparison. Other callable declarations are methods.

## Example

```bitlang
public class Player
{
    private int health;

    public int Health
    {
        get;
        private set;
    }

    public Player(int initialHealth)
    {
        health = initialHealth;
    }

    public void Damage(int amount)
    {
        health = health - amount;
    }
}
```

The example illustrates declaration shape; it does not define expression precedence, statement semantics, constructor rules, or the property-state expansion for each modifier.

## Preservation and property expansion

- The parser must preserve each modifier's spelling, order, and source location until modifier normalization has an explicit contract.
- Parsed classes, members, parameters, and accessors retain `ModifierTokens` containing the original preprocessed tokens, including canonical words and line/column positions. The existing `Modifiers` spelling list is a projection produced during parsing; neither list is reordered or deduplicated, and no modifier meaning is inferred here.
- `public`, `private`, `protected`, `internal`, `static`, `readonly`, `const`, `abstract`, `sealed`, `virtual`, and `override` must not be silently dropped.
- Access modifiers are source-level forms. Their expansion into independent property axes must preserve distinctions between class visibility, inheritance visibility, and module export. The mapping must be specified separately before semantic lowering relies on it.
- The strict Preprocessed output must continue to represent required property states explicitly; omitted states must not be guessed by later stages.
- Generic parameter, method type parameter, and type argument lists are retained as original token sequences for preprocessing; this parser does not resolve generic semantics.
- Generic lists must follow the declaration grammar even while retained as tokens: empty entries, missing separators, and nested types used as parameter names are rejected with source positions.
- Expression and statement token sequences must remain lossless while their grammars are unspecified.
- Balanced expression and statement sequences require correctly paired `()`, `[]`, and `{}` delimiters; literal contents do not participate in this check.

## Explicitly outside this version

Nested classes, interfaces, structs, enums, generic constraints, arrays, events, delegates, operator overloads, expression grammar, and statement grammar are not defined by this document. They can be added in later syntax versions without changing the source-spelling preservation rules above.
