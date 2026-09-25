package bitlang

// PreprocessedDeclaration represents one declaration after preprocessing.
// Every property axis applicable to the declaration must be present explicitly;
// later compiler stages must not reconstruct omitted defaults.
type PreprocessedDeclaration struct {
	Kind       string
	Name       CanonicalName
	Type       CanonicalName
	Properties []PreprocessedProperty
	Line       int
	Column     int
}
