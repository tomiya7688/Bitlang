package bitlang

import "testing"

func TestParsePreprocessedDeclarations(t *testing.T) {
	spec := declarationParserTestProperties()
	source, err := NewPreprocessor().Process(NewSourceText(
		"test.bit",
		"Private unnullable Int First; Public nullable Text Second;",
	))
	if err != nil {
		t.Fatal(err)
	}
	declarations, err := ParsePreprocessedDeclarations(spec, declarationParserTestKind(), source.Tokens)
	if err != nil {
		t.Fatal(err)
	}
	if len(declarations) != 2 {
		t.Fatalf("declaration count = %d, want 2", len(declarations))
	}
	if declarations[0].Name.Canonical != "first" || declarations[1].Name.Canonical != "second" {
		t.Fatalf("unexpected declarations: %#v", declarations)
	}
}

func TestParsePreprocessedDeclarationsRejectsMissingTerminator(t *testing.T) {
	source, err := NewPreprocessor().Process(NewSourceText("test.bit", "Private unnullable Int First"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ParsePreprocessedDeclarations(declarationParserTestProperties(), declarationParserTestKind(), source.Tokens); err == nil {
		t.Fatal("expected missing terminator error")
	}
}
