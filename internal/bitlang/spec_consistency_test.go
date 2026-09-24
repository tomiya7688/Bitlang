package bitlang

import "testing"

func TestValidateSpecificationConsistency(t *testing.T) {
	properties := PropertySpecification{Version: 1, Axes: []PropertyAxisSpec{{
		Name: "nullability", States: []string{"nullable", "unnullable"},
		AppliesTo: []string{"variable"},
	}}}
	declarations := DeclarationSpecification{Version: 1, Kinds: []DeclarationKindSpec{{
		Name: "variable", PropertyTarget: "variable",
	}}}
	if err := ValidateSpecificationConsistency(properties, declarations); err != nil {
		t.Fatal(err)
	}
}

func TestValidateSpecificationConsistencyRejectsUnknownTarget(t *testing.T) {
	properties := PropertySpecification{Version: 1, Axes: []PropertyAxisSpec{{
		Name: "nullability", States: []string{"nullable", "unnullable"},
		AppliesTo: []string{"variable"},
	}}}
	declarations := DeclarationSpecification{Version: 1, Kinds: []DeclarationKindSpec{{
		Name: "field", PropertyTarget: "field",
	}}}
	if err := ValidateSpecificationConsistency(properties, declarations); err == nil {
		t.Fatal("expected unknown property target error")
	}
}

func TestValidateSpecificationConsistencyRejectsSharedStateName(t *testing.T) {
	properties := PropertySpecification{Version: 1, Axes: []PropertyAxisSpec{
		{Name: "first", States: []string{"Enabled"}, AppliesTo: []string{"variable"}},
		{Name: "second", States: []string{"Enabled"}, AppliesTo: []string{"variable"}},
	}}
	declarations := DeclarationSpecification{Version: 1, Kinds: []DeclarationKindSpec{{
		Name: "variable", PropertyTarget: "variable",
	}}}
	if err := ValidateSpecificationConsistency(properties, declarations); err == nil {
		t.Fatal("expected shared property state error")
	}
}
