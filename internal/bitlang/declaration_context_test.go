package bitlang

import "testing"

func TestDeclarationKindAppliesToContext(t *testing.T) {
	kind := DeclarationKindSpec{Contexts: []string{"member"}}
	if !declarationKindAppliesToContext(kind, "MeMbEr") {
		t.Fatal("case-insensitive context did not match")
	}
	if declarationKindAppliesToContext(kind, "outer") {
		t.Fatal("unexpected context match")
	}
}

func TestDeclarationKindWithoutContextsIsUnrestricted(t *testing.T) {
	if !declarationKindAppliesToContext(DeclarationKindSpec{}, "") {
		t.Fatal("kind without contexts should be unrestricted")
	}
}
