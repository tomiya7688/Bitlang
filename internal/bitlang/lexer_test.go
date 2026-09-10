package bitlang

import "testing"

// TestLexerBasicSequence verifies that the bootstrap lexer preserves source
// spelling while separating identifiers, numbers, symbols, and literals.
func TestLexerBasicSequence(t *testing.T) {
	source := NewSourceText("test.bit", "Int HP = 100\nText Name = \"Alice\"")
	lexer := NewLexer(source)

	tokens, err := lexer.Lex()
	if err != nil {
		t.Fatal(err)
	}

	expectedKinds := []TokenKind{
		TokenIdentifier,
		TokenIdentifier,
		TokenSymbol,
		TokenNumber,
		TokenIdentifier,
		TokenIdentifier,
		TokenSymbol,
		TokenString,
		TokenEOF,
	}

	if len(tokens) != len(expectedKinds) {
		t.Fatalf("got %d tokens, want %d", len(tokens), len(expectedKinds))
	}

	for index, expectedKind := range expectedKinds {
		if tokens[index].Kind != expectedKind {
			t.Fatalf("token %d kind = %q, want %q", index, tokens[index].Kind, expectedKind)
		}
	}

	if tokens[0].Lexeme != "Int" {
		t.Fatalf("identifier spelling changed: %q", tokens[0].Lexeme)
	}
	if tokens[7].Lexeme != "\"Alice\"" {
		t.Fatalf("string literal changed: %q", tokens[7].Lexeme)
	}
}

// TestLexerTracksLines verifies diagnostic locations across newlines.
func TestLexerTracksLines(t *testing.T) {
	lexer := NewLexer(NewSourceText("test.bit", "First\nSecond"))
	tokens, err := lexer.Lex()
	if err != nil {
		t.Fatal(err)
	}

	if tokens[1].Line != 2 || tokens[1].Column != 1 {
		t.Fatalf("second token location = %d:%d, want 2:1", tokens[1].Line, tokens[1].Column)
	}
}

// TestLexerRejectsUnterminatedString verifies that malformed quoted input does
// not silently become a partial token stream.
func TestLexerRejectsUnterminatedString(t *testing.T) {
	lexer := NewLexer(NewSourceText("test.bit", "\"unfinished"))
	if _, err := lexer.Lex(); err == nil {
		t.Fatal("expected unterminated literal error")
	}
}
