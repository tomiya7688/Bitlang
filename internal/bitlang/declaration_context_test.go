package bitlang

import "testing"

func TestDeclarationKindAppliesToScope(t *testing.T) {
	scope, err := NewPreprocessedScope("MeMbEr", nil)
	if err != nil {
		t.Fatal(err)
	}
	kind := DeclarationKindSpec{Contexts: []string{"member"}}
	if !declarationKindAppliesToScope(kind, &scope) {
		t.Fatal("case-insensitive scope context did not match")
	}

	outer, err := NewPreprocessedScope("outer", nil)
	if err != nil {
		t.Fatal(err)
	}
	if declarationKindAppliesToScope(kind, &outer) {
		t.Fatal("unexpected scope context match")
	}
}

func TestDeclarationKindWithoutContextsIsUnrestricted(t *testing.T) {
	if !declarationKindAppliesToScope(DeclarationKindSpec{}, nil) {
		t.Fatal("kind without contexts should be unrestricted")
	}
}
