package bitlang

import (
	"reflect"
	"testing"
)

func TestLexicalPreprocessReturnsLexicalSource(t *testing.T) {
	source := NewSourceText("sample.bit", "Int PlayerHP = 100")
	lexical, err := LexicalPreprocess(source)
	if err != nil {
		t.Fatal(err)
	}
	if lexical.Path != "sample.bit" {
		t.Fatalf("path = %q, want sample.bit", lexical.Path)
	}
	if len(lexical.Tokens) != 5 {
		t.Fatalf("token count = %d, want 5", len(lexical.Tokens))
	}
	name := lexical.Tokens[1]
	if name.Lexeme != "PlayerHP" || name.Canonical != "playerhp" {
		t.Fatalf("identifier forms = %q/%q", name.Lexeme, name.Canonical)
	}
	if lexical.Tokens[4].Kind != TokenEOF {
		t.Fatalf("last token = %q, want eof", lexical.Tokens[4].Kind)
	}
}

func TestLexicalPreprocessDoesNotFoldLiteralContents(t *testing.T) {
	lexical, err := LexicalPreprocess(NewSourceText("sample.bit", "Name = \"PlayerHP\"; Ch = 'A'"))
	if err != nil {
		t.Fatal(err)
	}
	for _, token := range lexical.Tokens {
		if token.Kind == TokenString && (token.Lexeme != "\"PlayerHP\"" || token.Canonical != "") {
			t.Fatalf("string was canonicalized: %#v", token)
		}
		if token.Kind == TokenCharacter && (token.Lexeme != "'A'" || token.Canonical != "") {
			t.Fatalf("character was canonicalized: %#v", token)
		}
	}
}

func TestLexicalPreprocessIdentifiersAreCaseInsensitive(t *testing.T) {
	lexical, err := LexicalPreprocess(NewSourceText("sample.bit", "PlayerHP playerhp PLAYERHP"))
	if err != nil {
		t.Fatal(err)
	}
	for index := 0; index < 3; index++ {
		if lexical.Tokens[index].Canonical != "playerhp" {
			t.Fatalf("token %d canonical = %q", index, lexical.Tokens[index].Canonical)
		}
	}
}

func TestLexicalPreprocessLeavesDeclarationSyntaxAsTokens(t *testing.T) {
	lexical, err := LexicalPreprocess(NewSourceText("test.bit", "Private unnullable Int First;"))
	if err != nil {
		t.Fatal(err)
	}
	got := make([]string, len(lexical.Tokens))
	for index, token := range lexical.Tokens {
		got[index] = token.Lexeme
	}
	want := []string{"Private", "unnullable", "Int", "First", ";", ""}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("lexical tokens = %#v, want %#v", got, want)
	}
}
