package bitlang

import "testing"

func TestParsePreprocessedDeclarations(t *testing.T) {
	spec := declarationSequenceTestSpec()
	source, err := NewPreprocessor().Process(NewSourceText(
		"test.bit",
		"Private unnullable Int First; Public nullable Text Second;",
	))
	if err != nil {
		t.Fatal(err)
	}
	declarations, err := ParsePreprocessedDeclarations(spec, "variable", source.Tokens)
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

func TestParsePreprocessedDeclarationsRejectsMissingSemicolon(t *testing.T) {
	spec := declarationSequenceTestSpec()
	source, err := NewPreprocessor().Process(NewSourceText("test.bit", "Private unnullable Int First"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ParsePreprocessedDeclarations(spec, "variable", source.Tokens); err == nil {
		t.Fatal("expected missing semicolon error")
	}
}

func declarationSequenceTestSpec() PropertySpecification {
	return PropertySpecification{Version: 1, Axes: []PropertyAxisSpec{
		{Name: "visibility", States: []string{"Public", "Private"}, Exclusive: true, Required: true, AppliesTo: []string{"variable"}},
		{Name: "nullability", States: []string{"nullable", "unnullable"}, Exclusive: true, Required: true, AppliesTo: []string{"variable"}},
	}}
}
