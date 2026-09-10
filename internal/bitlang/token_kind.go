package bitlang

// TokenKind identifies the lexical category of one Bitlang token.
//
// These categories are intentionally language-neutral so another compiler
// implementation can reproduce the same token stream without copying Go
// implementation details.
type TokenKind string

const (
	TokenIdentifier TokenKind = "identifier"
	TokenNumber     TokenKind = "number"
	TokenString     TokenKind = "string"
	TokenCharacter  TokenKind = "character"
	TokenSymbol     TokenKind = "symbol"
	TokenEOF        TokenKind = "eof"
)
