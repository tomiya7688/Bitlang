package bitlang

// PreprocessedSource is the strict output produced for later Bitlang semantic
// analysis. Tokens preserve source spelling while Declarations preserve parsed
// semantic structure that later stages must not reconstruct from token order.
type PreprocessedSource struct {
	Path         string
	Tokens       []PreprocessedToken
	Declarations []PreprocessedDeclaration
	Scope        *PreprocessedScope
}
