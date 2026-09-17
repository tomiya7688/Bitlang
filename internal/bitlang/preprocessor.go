package bitlang

// Preprocessor converts human-written Bitlang source into the strict
// preprocessed representation. Macro and preprocessor-function expansion will
// be added here as those source semantics are implemented.
type Preprocessor struct{}

// NewPreprocessor creates the bootstrap preprocessor.
func NewPreprocessor() Preprocessor {
	return Preprocessor{}
}

// Process lexes source into the first strict PreprocessedSource form.
func (Preprocessor) Process(source SourceText) (PreprocessedSource, error) {
	tokens, err := NewLexer(source).Lex()
	if err != nil {
		return PreprocessedSource{}, err
	}
	return PreprocessedSource{Path: source.Path, Tokens: tokens}, nil
}
