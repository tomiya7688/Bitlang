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

// ParseMixedPreprocessedSource parses declarations without a context filter.
func ParseMixedPreprocessedSource(specs SpecificationSet, source PreprocessedSource) (PreprocessedSource, error) {
	return ParseMixedPreprocessedSourceInContext(specs, "", source)
}

// ParseMixedPreprocessedSourceInContext parses declarations after filtering
// candidate declaration kinds by the data-defined context.
func ParseMixedPreprocessedSourceInContext(specs SpecificationSet, context string, source PreprocessedSource) (PreprocessedSource, error) {
	declarations, err := ParseMixedPreprocessedDeclarationsInContext(specs, context, source.Tokens)
	if err != nil {
		return PreprocessedSource{}, err
	}
	source.Declarations = declarations
	return source, nil
}
