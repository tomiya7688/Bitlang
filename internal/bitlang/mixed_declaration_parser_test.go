package bitlang

import "testing"

func TestParseMixedPreprocessedDeclarationsDetectsKinds(t *testing.T) {
	specs := mixedDeclarationTestSpecifications()
	source, err := NewPreprocessor().Process(NewSourceText(
		"test.bit",
		"Private Int Count; Private Protected Text Name;",
	))
	if err != nil {
		t.Fatal(err)
	}

	declarations, err := ParseMixedPreprocessedDeclarations(specs, source.Tokens)
	if err != nil {
		t.Fatal(err)
	}
	if len(declarations) != 2 {
		t.Fatalf("declaration count = %d, want 2", len(declarations))
	}
	if declarations[0].Kind != "variable" || declarations[1].Kind != "field" {
		t.Fatalf("kinds = %q, %q", declarations[0].Kind, declarations[1].Kind)
	}
}

func TestParseMixedPreprocessedDeclarationsRejectsAmbiguousKind(t *testing.T) {
	specs := mixedDeclarationTestSpecifications()
	specs.Declarations.Kinds = append(specs.Declarations.Kinds, DeclarationKindSpec{
		Name: "alias_variable", PropertyTarget: "variable", Terminator: ";",
		Layout: []string{"properties", "type", "name"},
	})
	source, err := NewPreprocessor().Process(NewSourceText("test.bit", "Private Int Count;"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ParseMixedPreprocessedDeclarations(specs, source.Tokens); err == nil {
		t.Fatal("expected ambiguous declaration error")
	}
}

func TestParseMixedPreprocessedDeclarationsRejectsUnknownKind(t *testing.T) {
	specs := mixedDeclarationTestSpecifications()
	source, err := NewPreprocessor().Process(NewSourceText("test.bit", "Protected Int Count;"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ParseMixedPreprocessedDeclarations(specs, source.Tokens); err == nil {
		t.Fatal("expected unknown declaration error")
	}
}

func mixedDeclarationTestSpecifications() SpecificationSet {
	return SpecificationSet{
		Properties: PropertySpecification{Version: 1, Axes: []PropertyAxisSpec{
			{
				Name: "visibility", States: []string{"Public", "Private"},
				Exclusive: true, Required: true, AppliesTo: []string{"variable", "field"},
			},
			{
				Name: "inheritance", States: []string{"Protected", "Unprotected"},
				Exclusive: true, Required: true, AppliesTo: []string{"field"},
			},
		}},
		Declarations: DeclarationSpecification{Version: 1, Kinds: []DeclarationKindSpec{
			{
				Name: "variable", PropertyTarget: "variable", Terminator: ";",
				Layout: []string{"properties", "type", "name"},
			},
			{
				Name: "field", PropertyTarget: "field", Terminator: ";",
				Layout: []string{"properties", "type", "name"},
			},
		}},
	}
}


func TestParseMixedPreprocessedDeclarationsUsesContext(t *testing.T) {
	specs := SpecificationSet{
		Properties: PropertySpecification{Version: 1, Axes: []PropertyAxisSpec{{
			Name: "visibility", States: []string{"Public", "Private"},
			Exclusive: true, Required: true, AppliesTo: []string{"variable", "field"},
		}}},
		Declarations: DeclarationSpecification{Version: 1, Kinds: []DeclarationKindSpec{
			{
				Name: "variable", PropertyTarget: "variable", Terminator: ";",
				Layout: []string{"properties", "type", "name"}, Contexts: []string{"outer"},
			},
			{
				Name: "field", PropertyTarget: "field", Terminator: ";",
				Layout: []string{"properties", "type", "name"}, Contexts: []string{"member"},
			},
		}},
	}
	source, err := NewPreprocessor().Process(NewSourceText("test.bit", "Private Int Count;"))
	if err != nil {
		t.Fatal(err)
	}
	declarations, err := ParseMixedPreprocessedDeclarationsInContext(specs, "MeMbEr", source.Tokens)
	if err != nil {
		t.Fatal(err)
	}
	if len(declarations) != 1 || declarations[0].Kind != "field" {
		t.Fatalf("unexpected declarations: %#v", declarations)
	}
}

func TestParseMixedPreprocessedDeclarationsRequiresConfiguredContext(t *testing.T) {
	specs := SpecificationSet{
		Properties: PropertySpecification{Version: 1, Axes: []PropertyAxisSpec{{
			Name: "visibility", States: []string{"Public", "Private"},
			Exclusive: true, Required: true, AppliesTo: []string{"variable"},
		}}},
		Declarations: DeclarationSpecification{Version: 1, Kinds: []DeclarationKindSpec{{
			Name: "variable", PropertyTarget: "variable", Terminator: ";",
			Layout: []string{"properties", "type", "name"}, Contexts: []string{"outer"},
		}}},
	}
	source, err := NewPreprocessor().Process(NewSourceText("test.bit", "Private Int Count;"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ParseMixedPreprocessedDeclarations(specs, source.Tokens); err == nil {
		t.Fatal("expected missing context error")
	}
}
