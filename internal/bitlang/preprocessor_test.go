package bitlang

import "testing"

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
	name := preprocessed.Tokens[1]
	if name.Lexeme != "PlayerHP" || name.Canonical != "playerhp" {
		t.Fatalf("identifier forms = %q/%q", name.Lexeme, name.Canonical)
	}
	if preprocessed.Tokens[4].Kind != TokenEOF {
		t.Fatalf("last token = %q, want eof", preprocessed.Tokens[4].Kind)
	}
}

func TestPreprocessorDoesNotFoldLiteralContents(t *testing.T) {
	preprocessed, err := NewPreprocessor().Process(NewSourceText("sample.bit", "Name = \"PlayerHP\"; Ch = 'A'"))
	if err != nil {
		t.Fatal(err)
	}
	for _, token := range preprocessed.Tokens {
		if token.Kind == TokenString && (token.Lexeme != "\"PlayerHP\"" || token.Canonical != "") {
			t.Fatalf("string was canonicalized: %#v", token)
		}
		if token.Kind == TokenCharacter && (token.Lexeme != "'A'" || token.Canonical != "") {
			t.Fatalf("character was canonicalized: %#v", token)
		}
	}
}

func TestPreprocessorIdentifiersAreCaseInsensitive(t *testing.T) {
	preprocessed, err := NewPreprocessor().Process(NewSourceText("sample.bit", "PlayerHP playerhp PLAYERHP"))
	if err != nil {
		t.Fatal(err)
	}
	for index := 0; index < 3; index++ {
		if preprocessed.Tokens[index].Canonical != "playerhp" {
			t.Fatalf("token %d canonical = %q", index, preprocessed.Tokens[index].Canonical)
		}
	}
}
