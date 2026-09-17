package bitlang

// PreprocessedSource is the strict token-level input produced for later
// Bitlang semantic analysis. Source spelling and canonical identifier forms
// are both retained so later stages do not depend on host-language casing.
type PreprocessedSource struct {
	Path   string
	Tokens []PreprocessedToken
}
