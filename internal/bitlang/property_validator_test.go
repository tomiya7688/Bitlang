package bitlang

import "testing"

func TestValidatePropertiesRequiresApplicableAxis(t *testing.T) {
	spec := PropertySpecification{Version: 1, Axes: []PropertyAxisSpec{{
		Name: "nullability", States: []string{"nullable", "unnullable"},
		Exclusive: true, Required: true, AppliesTo: []string{"variable"},
	}}}
	if err := ValidateProperties(spec, "variable", nil); err == nil {
		t.Fatal("expected missing required property error")
	}
}

func TestValidatePropertiesRejectsConflictingStates(t *testing.T) {
	spec := PropertySpecification{Version: 1, Axes: []PropertyAxisSpec{{
		Name: "nullability", States: []string{"nullable", "unnullable"},
		Exclusive: true, Required: true, AppliesTo: []string{"variable"},
	}}}
	properties := []PreprocessedProperty{"nullable", "unnullable"}
	if err := ValidateProperties(spec, "variable", properties); err == nil {
		t.Fatal("expected conflicting property error")
	}
}

func TestValidatePropertiesAcceptsCompleteAxis(t *testing.T) {
	spec := PropertySpecification{Version: 1, Axes: []PropertyAxisSpec{{
		Name: "nullability", States: []string{"nullable", "unnullable"},
		Exclusive: true, Required: true, AppliesTo: []string{"variable"},
	}}}
	if err := ValidateProperties(spec, "variable", []PreprocessedProperty{"unnullable"}); err != nil {
		t.Fatal(err)
	}
}

func TestValidatePropertiesRejectsInapplicableProperty(t *testing.T) {
	spec := PropertySpecification{Version: 1, Axes: []PropertyAxisSpec{{
		Name: "nullability", States: []string{"nullable", "unnullable"},
		Exclusive: true, Required: true, AppliesTo: []string{"variable"},
	}}}
	if err := ValidateProperties(spec, "function", []PreprocessedProperty{"nullable"}); err == nil {
		t.Fatal("expected inapplicable property error")
	}
}
