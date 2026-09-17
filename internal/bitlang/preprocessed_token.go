package bitlang

// PreprocessedToken stores one strict preprocessor output token.
// Lexeme always retains source spelling. Canonical is populated only for
// identifiers, keeping literal contents completely outside name folding.
type PreprocessedToken struct {
	Kind      TokenKind
	Lexeme    string
	Canonical string
	Line      int
	Column    int
}

func newPreprocessedToken(token Token) (PreprocessedToken, error) {
	result := PreprocessedToken{
		Kind: token.Kind, Lexeme: token.Lexeme,
		Line: token.Line, Column: token.Column,
	}
	if token.Kind != TokenIdentifier {
		return result, nil
	}
	canonical, err := CanonicalizeIdentifier(token.Lexeme)
	if err != nil {
		return PreprocessedToken{}, err
	}
	result.Canonical = canonical
	return result, nil
}
