package bitlang

import "testing"

func TestParsePreprocessedDeclaration(t *testing.T) {
	spec := declarationParserTestProperties()
	source, err := NewPreprocessor().Process(NewSourceText("test.bit", "Private unnullable Int PlayerHP;"))
	if err != nil {
		t.Fatal(err)
	}
	decl, err := ParsePreprocessedDeclaration(spec, declarationParserTestKind(), source.Tokens)
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
	source, err := NewPreprocessor().Process(NewSourceText("test.bit", "Int PlayerHP;"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ParsePreprocessedDeclaration(declarationParserTestProperties(), declarationParserTestKind(), source.Tokens); err == nil {
		t.Fatal("expected missing property error")
	}
}

func declarationParserTestKind() DeclarationKindSpec {
	return DeclarationKindSpec{Name: "variable", PropertyTarget: "variable", Terminator: ";", Layout: []string{"properties", "type", "name"}}
}

func declarationParserTestProperties() PropertySpecification {
	return PropertySpecification{Version: 1, Axes: []PropertyAxisSpec{
		{Name: "visibility", States: []string{"Public", "Private"}, Exclusive: true, Required: true, AppliesTo: []string{"variable"}},
		{Name: "nullability", States: []string{"nullable", "unnullable"}, Exclusive: true, Required: true, AppliesTo: []string{"variable"}},
	}}
}
