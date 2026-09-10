package bitlang

// SourceText stores one Bitlang source unit before lexical analysis.
//
// The source path is diagnostic metadata only. Text contains the source exactly
// as received and must not be normalized before tokenization because later
// stages may need exact source locations and literal contents.
type SourceText struct {
	Path string
	Text string
}

// NewSourceText creates an immutable-by-convention Bitlang source value.
func NewSourceText(path string, text string) SourceText {
	return SourceText{Path: path, Text: text}
}
