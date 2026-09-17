package bitlang

import "testing"

// TestPreprocessorProducesStrictTokenInput verifies the first Source to
// Preprocessed boundary without inventing later macro semantics.
func TestPreprocessorProducesStrictTokenInput(t *testing.T) {
	source := NewSourceText("sample.bit", "Int PlayerHP = 100")
	preprocessed, err := NewPreprocessor().Process(source)
	if err != nil {
		t.Fatal(err)
	}
	if preprocessed.Path != "sample.bit" {
		t.Fatalf("path = %q, want sample.bit", preprocessed.Path)
	}
	if len(preprocessed.Tokens) != 5 {
		t.Fatalf("token count = %d, want 5", len(preprocessed.Tokens))
	}
	if preprocessed.Tokens[1].Lexeme != "PlayerHP" {
		t.Fatalf("source spelling changed: %q", preprocessed.Tokens[1].Lexeme)
	}
	if preprocessed.Tokens[4].Kind != TokenEOF {
		t.Fatalf("last token = %q, want eof", preprocessed.Tokens[4].Kind)
	}
}
