package bitlang

// PreprocessedSource is the strict token-level input produced for later
// Bitlang semantic analysis. Source spelling is retained in each token so
// diagnostics can still refer to the original program.
type PreprocessedSource struct {
	Path   string
	Tokens []Token
}
