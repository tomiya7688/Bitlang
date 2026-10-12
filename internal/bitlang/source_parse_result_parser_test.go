package bitlang

import "testing"

func TestParseDeclarationsInSourceResultStoresDeclarations(t *testing.T) {
	source, err := LexicalPreprocess(NewSourceText(
		"test.bit",
		"Private unnullable Int First; Public nullable Text Second;",
	))
	if err != nil {
		t.Fatal(err)
	}

	parsed, err := ParseDeclarationsInSourceResult(
		declarationParserTestSpecifications(), "variable", NewSourceParseResult(source),
	)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Source.Path != "test.bit" {
		t.Fatalf("path = %q, want test.bit", parsed.Source.Path)
	}
	if len(parsed.Source.Tokens) != len(source.Tokens) {
		t.Fatalf("token count changed: %d -> %d", len(source.Tokens), len(parsed.Source.Tokens))
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

func TestParseMixedDeclarationsInSourceResultStoresDetectedKinds(t *testing.T) {
	source, err := LexicalPreprocess(NewSourceText(
		"test.bit",
		"Private Int Count; Private Protected Text Name;",
	))
	if err != nil {
		t.Fatal(err)
	}

	parsed, err := ParseMixedDeclarationsInSourceResult(mixedDeclarationTestSpecifications(), NewSourceParseResult(source))
	if err != nil {
		t.Fatal(err)
	}
	if len(parsed.Declarations) != 2 {
		t.Fatalf("declaration count = %d, want 2", len(parsed.Declarations))
	}
	if parsed.Declarations[0].Kind != "variable" || parsed.Declarations[1].Kind != "field" {
		t.Fatalf("unexpected declaration kinds: %#v", parsed.Declarations)
	}
}

func TestParseMixedDeclarationsInSourceResultUsesAttachedScope(t *testing.T) {
	specs := contextSourceTestSpecifications()
	scope, err := NewPreprocessedScope("MeMbEr", nil)
	if err != nil {
		t.Fatal(err)
	}
	source, err := LexicalPreprocess(NewSourceText("test.bit", "Private Int Count;"))
	if err != nil {
		t.Fatal(err)
	}
	parseResult := NewSourceParseResult(source)
	parseResult.Scope = &scope

	parsed, err := ParseMixedDeclarationsInSourceResult(specs, parseResult)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Scope == nil || parsed.Scope.Context.Canonical != "member" {
		t.Fatalf("unexpected scope: %#v", parsed.Scope)
	}
	if len(parsed.Declarations) != 1 || parsed.Declarations[0].Kind != "field" {
		t.Fatalf("unexpected declarations: %#v", parsed.Declarations)
	}
}

func contextSourceTestSpecifications() SpecificationSet {
	return SpecificationSet{
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
}
