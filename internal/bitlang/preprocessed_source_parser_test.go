package bitlang

import "testing"

func TestParsePreprocessedSourceStoresDeclarations(t *testing.T) {
	source, err := NewPreprocessor().Process(NewSourceText(
		"test.bit",
		"Private unnullable Int First; Public nullable Text Second;",
	))
	if err != nil {
		t.Fatal(err)
	}

	parsed, err := ParsePreprocessedSource(declarationParserTestSpecifications(), "variable", source)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Path != "test.bit" {
		t.Fatalf("path = %q, want test.bit", parsed.Path)
	}
	if len(parsed.Tokens) != len(source.Tokens) {
		t.Fatalf("token count changed: %d -> %d", len(source.Tokens), len(parsed.Tokens))
	}
	if len(parsed.Declarations) != 2 {
		t.Fatalf("declaration count = %d, want 2", len(parsed.Declarations))
	}
	if parsed.Declarations[0].Kind != "variable" || parsed.Declarations[1].Kind != "variable" {
		t.Fatalf("declaration kinds = %#v", parsed.Declarations)
	}
	if parsed.Declarations[0].Name.Canonical != "first" || parsed.Declarations[1].Name.Canonical != "second" {
		t.Fatalf("unexpected declarations: %#v", parsed.Declarations)
	}
}

func TestPreprocessorLeavesDeclarationsUnparsed(t *testing.T) {
	source, err := NewPreprocessor().Process(NewSourceText("test.bit", "Private unnullable Int First;"))
	if err != nil {
		t.Fatal(err)
	}
	if len(source.Declarations) != 0 {
		t.Fatalf("preprocessor unexpectedly parsed declarations: %#v", source.Declarations)
	}
}
