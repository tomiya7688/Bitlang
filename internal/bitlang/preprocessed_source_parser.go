package bitlang

// ParsePreprocessedSource parses declarations from an existing strict token
// stream and returns the same source with semantic declarations attached.
// One declaration kind is accepted per call until mixed-kind source grammar is
// explicitly defined.
func ParsePreprocessedSource(specs SpecificationSet, kindName string, source PreprocessedSource) (PreprocessedSource, error) {
	declarations, err := ParsePreprocessedDeclarations(specs, kindName, source.Tokens)
	if err != nil {
		return PreprocessedSource{}, err
	}
	source.Declarations = declarations
	return source, nil
}
