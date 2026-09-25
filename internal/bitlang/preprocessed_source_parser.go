package bitlang

// ParsePreprocessedSource parses declarations from an existing strict token
// stream using one explicitly selected declaration kind.
func ParsePreprocessedSource(specs SpecificationSet, kindName string, source PreprocessedSource) (PreprocessedSource, error) {
	declarations, err := ParsePreprocessedDeclarations(specs, kindName, source.Tokens)
	if err != nil {
		return PreprocessedSource{}, err
	}
	source.Declarations = declarations
	return source, nil
}

// ParseMixedPreprocessedSource parses declarations whose kinds are determined
// from the machine-readable declaration and property specifications.
func ParseMixedPreprocessedSource(specs SpecificationSet, source PreprocessedSource) (PreprocessedSource, error) {
	declarations, err := ParseMixedPreprocessedDeclarations(specs, source.Tokens)
	if err != nil {
		return PreprocessedSource{}, err
	}
	source.Declarations = declarations
	return source, nil
}
