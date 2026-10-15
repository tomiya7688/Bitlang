package bitlang

import "testing"

func TestParsePreprocessedDeclaration(t *testing.T) {
	source, err := LexicalPreprocess(NewSourceText("test.bit", "pRiVaTe UnNuLlAbLe Int PlayerHP;"))
	if err != nil {
		t.Fatal(err)
	}
	decl, err := ParsePreprocessedDeclaration(declarationParserTestSpecifications(), "variable", source.Tokens)
	if err != nil {
		t.Fatal(err)
	}
	if decl.Name.Spelling != "PlayerHP" || decl.Name.Canonical != "playerhp" {
		t.Fatalf("unexpected declaration name: %#v", decl.Name)
	}
	if decl.Kind != "variable" {
		t.Fatalf("declaration kind = %q, want variable", decl.Kind)
	}
	if decl.Type.Canonical != "int" {
		t.Fatalf("unexpected type: %#v", decl.Type)
	}
	if len(decl.Properties) != 2 || decl.Properties[0] != "Private" || decl.Properties[1] != "unnullable" {
		t.Fatalf("properties = %#v", decl.Properties)
	}
}

func TestParsePreprocessedDeclarationUsesConfiguredLayout(t *testing.T) {
	specs := declarationParserTestSpecifications()
	specs.Declarations.Kinds[0].Layout = []string{"type", "name", "properties"}
	source, err := LexicalPreprocess(NewSourceText("test.bit", "Int PlayerHP pRiVaTe UnNuLlAbLe;"))
	if err != nil {
		t.Fatal(err)
	}
	decl, err := ParsePreprocessedDeclaration(specs, "variable", source.Tokens)
	if err != nil {
		t.Fatal(err)
	}
	if decl.Type.Canonical != "int" || decl.Name.Canonical != "playerhp" {
		t.Fatalf("unexpected declaration: %#v", decl)
	}
	if len(decl.Properties) != 2 || decl.Properties[0] != "Private" || decl.Properties[1] != "unnullable" {
		t.Fatalf("properties = %#v", decl.Properties)
	}
}

func TestParsePreprocessedDeclarationRejectsIncompleteProperties(t *testing.T) {
	source, err := LexicalPreprocess(NewSourceText("test.bit", "Int PlayerHP;"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ParsePreprocessedDeclaration(declarationParserTestSpecifications(), "variable", source.Tokens); err == nil {
		t.Fatal("expected missing property error")
	}
}

func TestParsePreprocessedDeclarationNormalizesPropertyVariants(t *testing.T) {
	specs := declarationParserTestSpecifications()
	specs.Declarations.Kinds[0].PropertyTarget = "VARIABLE"
	specs.Properties.Axes[0].AppliesTo = []string{"Variable"}
	specs.Properties.Axes[1] = PropertyAxisSpec{
		Name: "release", States: []string{"Disable_Auto_Release"}, Required: true, Exclusive: true, AppliesTo: []string{"Variable"},
	}
	if err := ValidateSpecificationConsistency(specs.Properties, specs.Declarations); err != nil {
		t.Fatal(err)
	}
	for _, text := range []string{
		"Private Disable_Auto_Release Int PlayerHP;",
		"private disable_auto_release Int PlayerHP;",
		"PRIVATE DISABLE_AUTO_RELEASE Int PlayerHP;",
	} {
		t.Run(text, func(t *testing.T) {
			source, err := LexicalPreprocess(NewSourceText("case.bit", text))
			if err != nil {
				t.Fatal(err)
			}
			declaration, err := ParsePreprocessedDeclaration(specs, "variable", source.Tokens)
			if err != nil {
				t.Fatal(err)
			}
			if len(declaration.Properties) != 2 || declaration.Properties[0] != "Private" || declaration.Properties[1] != "Disable_Auto_Release" {
				t.Fatalf("properties = %#v", declaration.Properties)
			}
			if source.Tokens[0].Lexeme != text[:len("Private")] {
				t.Fatal("source property spelling changed")
			}
		})
	}
}

func declarationParserTestSpecifications() SpecificationSet {
	return SpecificationSet{
		Properties: PropertySpecification{Version: 1, Axes: []PropertyAxisSpec{
			{Name: "visibility", States: []string{"Public", "Private"}, Exclusive: true, Required: true, AppliesTo: []string{"variable"}},
			{Name: "nullability", States: []string{"nullable", "unnullable"}, Exclusive: true, Required: true, AppliesTo: []string{"variable"}},
		}},
		Declarations: DeclarationSpecification{Version: 1, Kinds: []DeclarationKindSpec{{
			Name: "variable", PropertyTarget: "variable", Terminator: ";", Layout: []string{"properties", "type", "name"},
		}}},
	}
}
