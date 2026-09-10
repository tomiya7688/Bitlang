package bitlang

// Token represents one lexical unit produced from Bitlang source.
//
// Lexeme preserves the original source spelling. Identifier canonicalization is
// deliberately separate so string and character literal contents are never
// modified by case-insensitive name handling.
type Token struct {
	Kind   TokenKind
	Lexeme string
	Line   int
	Column int
}

// NewToken creates one token with a 1-based source location.
func NewToken(kind TokenKind, lexeme string, line int, column int) Token {
	return Token{Kind: kind, Lexeme: lexeme, Line: line, Column: column}
}
