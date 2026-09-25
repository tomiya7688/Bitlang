package bitlang

// PreprocessedScope represents one semantic scope together with the declaration
// context used to interpret declarations directly contained by that scope.
// Context names remain data-defined rather than hardcoded in the Go bootstrap.
type PreprocessedScope struct {
	Context CanonicalName
	Parent  *PreprocessedScope
}

// NewPreprocessedScope creates a scope with a case-insensitive context name.
func NewPreprocessedScope(context string, parent *PreprocessedScope) (PreprocessedScope, error) {
	name, err := NewCanonicalName(context)
	if err != nil {
		return PreprocessedScope{}, err
	}
	return PreprocessedScope{Context: name, Parent: parent}, nil
}
