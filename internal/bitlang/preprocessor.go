package bitlang

// Preprocessor converts human-written Bitlang source into the strict
// preprocessed representation. Macro and preprocessor-function expansion will
// be added here as those source semantics are implemented.
type Preprocessor struct{}

// NewPreprocessor creates the bootstrap preprocessor.
func NewPreprocessor() Preprocessor {
	return Preprocessor{}
}

// Process lexes source and makes identifier comparison semantics explicit.
func (Preprocessor) Process(source SourceText) (PreprocessedSource, error) {
	tokens, err := NewLexer(source).Lex()
	if err != nil {
		return PreprocessedSource{}, err
	}
	strict := make([]PreprocessedToken, 0, len(tokens))
	for _, token := range tokens {
		converted, err := newPreprocessedToken(token)
		if err != nil {
			return PreprocessedSource{}, err
		}
		strict = append(strict, converted)
	}
	return PreprocessedSource{Path: source.Path, Tokens: strict}, nil
}
