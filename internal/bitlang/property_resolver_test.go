package bitlang

import "testing"

func TestResolvePropertiesUsesSpecificationSpelling(t *testing.T) {
	spec := PropertySpecification{Version: 1, Axes: []PropertyAxisSpec{{
		Name: "visibility", States: []string{"Public", "Private"},
		AppliesTo: []string{"variable"},
	}}}
	resolved, err := ResolveProperties(spec, "variable", []PreprocessedProperty{"pRiVaTe"})
	if err != nil {
		t.Fatal(err)
	}
	if len(resolved) != 1 || resolved[0] != "Private" {
		t.Fatalf("resolved = %#v, want Private", resolved)
	}
}

func TestResolvePropertiesRejectsInapplicableState(t *testing.T) {
	spec := PropertySpecification{Version: 1, Axes: []PropertyAxisSpec{{
		Name: "visibility", States: []string{"Public", "Private"},
		AppliesTo: []string{"field"},
	}}}
	if _, err := ResolveProperties(spec, "variable", []PreprocessedProperty{"Private"}); err == nil {
		t.Fatal("expected inapplicable property error")
	}
}
