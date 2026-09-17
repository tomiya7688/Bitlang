package bitlang

import "testing"

func TestLexerBasicSequence(t *testing.T) {
	source := NewSourceText("test.bit", "Int HP = 100\nText Name = \"Alice\"")
	tokens, err := NewLexer(source).Lex()
	if err != nil {
		t.Fatal(err)
	}
	expectedKinds := []TokenKind{TokenIdentifier, TokenIdentifier, TokenSymbol, TokenNumber, TokenIdentifier, TokenIdentifier, TokenSymbol, TokenString, TokenEOF}
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

func TestLexerTracksLines(t *testing.T) {
	tokens, err := NewLexer(NewSourceText("test.bit", "First\nSecond")).Lex()
	if err != nil {
		t.Fatal(err)
	}
	if tokens[1].Line != 2 || tokens[1].Column != 1 {
		t.Fatalf("second token location = %d:%d, want 2:1", tokens[1].Line, tokens[1].Column)
	}
}

func TestLexerRejectsUnterminatedString(t *testing.T) {
	if _, err := NewLexer(NewSourceText("test.bit", "\"unfinished")).Lex(); err == nil {
		t.Fatal("expected unterminated literal error")
	}
}

func TestLexerSkipsComments(t *testing.T) {
	source := NewSourceText("test.bit", "First // line\n/* block\ncomment */ Second")
	tokens, err := NewLexer(source).Lex()
	if err != nil {
		t.Fatal(err)
	}
	if len(tokens) != 3 || tokens[0].Lexeme != "First" || tokens[1].Lexeme != "Second" {
		t.Fatalf("unexpected tokens after comments: %#v", tokens)
	}
	if tokens[1].Line != 3 {
		t.Fatalf("second token line = %d, want 3", tokens[1].Line)
	}
}

func TestLexerPreservesCommentMarkersInString(t *testing.T) {
	tokens, err := NewLexer(NewSourceText("test.bit", "\"// not comment\" \"/* not comment */\"")).Lex()
	if err != nil {
		t.Fatal(err)
	}
	if tokens[0].Kind != TokenString || tokens[1].Kind != TokenString {
		t.Fatal("comment markers inside strings must remain literals")
	}
}

func TestLexerRejectsUnterminatedBlockComment(t *testing.T) {
	if _, err := NewLexer(NewSourceText("test.bit", "/* unfinished")).Lex(); err == nil {
		t.Fatal("expected unterminated block comment error")
	}
}
