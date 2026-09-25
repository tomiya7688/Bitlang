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

// ParseMixedPreprocessedSource parses declarations using the scope already
// attached to the Preprocessed source.
func ParseMixedPreprocessedSource(specs SpecificationSet, source PreprocessedSource) (PreprocessedSource, error) {
	return ParseMixedPreprocessedSourceInScope(specs, source.Scope, source)
}

// ParseMixedPreprocessedSourceInContext preserves the string-based bootstrap
// entrypoint while attaching a typed scope to the resulting source.
func ParseMixedPreprocessedSourceInContext(specs SpecificationSet, context string, source PreprocessedSource) (PreprocessedSource, error) {
	if context == "" {
		return ParseMixedPreprocessedSourceInScope(specs, nil, source)
	}
	scope, err := NewPreprocessedScope(context, source.Scope)
	if err != nil {
		return PreprocessedSource{}, err
	}
	return ParseMixedPreprocessedSourceInScope(specs, &scope, source)
}

// ParseMixedPreprocessedSourceInScope parses declarations using an explicit
// semantic scope and stores that scope on the resulting Preprocessed source.
func ParseMixedPreprocessedSourceInScope(specs SpecificationSet, scope *PreprocessedScope, source PreprocessedSource) (PreprocessedSource, error) {
	declarations, err := ParseMixedPreprocessedDeclarationsInScope(specs, scope, source.Tokens)
	if err != nil {
		return PreprocessedSource{}, err
	}
	source.Scope = scope
	source.Declarations = declarations
	return source, nil
}
