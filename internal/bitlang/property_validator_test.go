package bitlang

import (
	"reflect"
	"testing"
)

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

func TestValidatePropertiesIsCaseInsensitive(t *testing.T) {
	spec := PropertySpecification{Version: 1, Axes: []PropertyAxisSpec{
		{Name: "visibility", States: []string{"Private", "Public"}, Exclusive: true, Required: true, AppliesTo: []string{"Variable"}},
		{Name: "release", States: []string{"Auto_Release", "Disable_Auto_Release"}, Exclusive: true, Required: true, AppliesTo: []string{"Variable"}},
	}}
	cases := []struct {
		name       string
		target     string
		properties []PreprocessedProperty
	}{
		{"spec spelling", "Variable", []PreprocessedProperty{"Private", "Disable_Auto_Release"}},
		{"lowercase", "variable", []PreprocessedProperty{"private", "disable_auto_release"}},
		{"uppercase", "VARIABLE", []PreprocessedProperty{"PRIVATE", "DISABLE_AUTO_RELEASE"}},
		{"mixed case", "VaRiAbLe", []PreprocessedProperty{"pRiVaTe", "dIsAbLe_AuTo_ReLeAsE"}},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			original := append([]PreprocessedProperty(nil), test.properties...)
			if err := ValidateProperties(spec, test.target, test.properties); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(test.properties, original) {
				t.Fatalf("validation rewrote source spelling: %#v", test.properties)
			}
		})
	}
}

func TestValidatePropertiesRejectsCaseVariantBypasses(t *testing.T) {
	spec := PropertySpecification{Version: 1, Axes: []PropertyAxisSpec{{
		Name: "visibility", States: []string{"Private", "Public"}, Exclusive: true, Required: true, AppliesTo: []string{"Variable"},
	}}}
	cases := []struct {
		name       string
		properties []PreprocessedProperty
	}{
		{"missing required axis", nil},
		{"conflicting states", []PreprocessedProperty{"private", "PUBLIC"}},
		{"duplicate canonical state", []PreprocessedProperty{"private", "PRIVATE"}},
		{"unknown state", []PreprocessedProperty{"Not_Defined"}},
		{"empty state", []PreprocessedProperty{""}},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			if err := ValidateProperties(spec, "VARIABLE", test.properties); err == nil {
				t.Fatal("expected validation error despite target case variation")
			}
		})
	}
}
