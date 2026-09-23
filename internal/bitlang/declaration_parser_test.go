package bitlang

import "testing"

func TestParsePreprocessedDeclaration(t *testing.T) {
	spec := PropertySpecification{Version: 1, Axes: []PropertyAxisSpec{
		{Name: "visibility", States: []string{"Public", "Private"}, Exclusive: true, Required: true, AppliesTo: []string{"variable"}},
		{Name: "nullability", States: []string{"nullable", "unnullable"}, Exclusive: true, Required: true, AppliesTo: []string{"variable"}},
	}}
	source, err := NewPreprocessor().Process(NewSourceText("test.bit", "Private unnullable Int PlayerHP;"))
	if err != nil {
		t.Fatal(err)
	}
	decl, err := ParsePreprocessedDeclaration(spec, "variable", source.Tokens)
	if err != nil {
		t.Fatal(err)
	}
	if decl.Name.Spelling != "PlayerHP" || decl.Name.Canonical != "playerhp" {
		t.Fatalf("unexpected declaration name: %#v", decl.Name)
	}
	if decl.Type.Canonical != "int" {
		t.Fatalf("unexpected type: %#v", decl.Type)
	}
}

func TestParsePreprocessedDeclarationRejectsIncompleteProperties(t *testing.T) {
	spec := PropertySpecification{Version: 1, Axes: []PropertyAxisSpec{
		{Name: "nullability", States: []string{"nullable", "unnullable"}, Exclusive: true, Required: true, AppliesTo: []string{"variable"}},
	}}
	source, err := NewPreprocessor().Process(NewSourceText("test.bit", "Int PlayerHP;"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ParsePreprocessedDeclaration(spec, "variable", source.Tokens); err == nil {
		t.Fatal("expected missing property error")
	}
}
