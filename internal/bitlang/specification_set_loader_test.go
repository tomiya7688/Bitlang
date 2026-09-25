package bitlang

import "testing"

func TestLoadSpecificationSet(t *testing.T) {
	properties := []byte(`{"version":1,"axes":[{"name":"nullability","states":["nullable","unnullable"],"exclusive":true,"required":true,"applies_to":["variable"]}]}`)
	declarations := []byte(`{"version":1,"kinds":[{"name":"variable","property_target":"variable","terminator":";","layout":["properties","type","name"]}]}`)

	specs, err := LoadSpecificationSet(properties, declarations)
	if err != nil {
		t.Fatal(err)
	}
	if len(specs.Properties.Axes) != 1 || len(specs.Declarations.Kinds) != 1 {
		t.Fatalf("unexpected specification set: %#v", specs)
	}
}

func TestLoadSpecificationSetRejectsCrossSpecMismatch(t *testing.T) {
	properties := []byte(`{"version":1,"axes":[{"name":"nullability","states":["nullable","unnullable"],"exclusive":true,"required":true,"applies_to":["variable"]}]}`)
	declarations := []byte(`{"version":1,"kinds":[{"name":"field","property_target":"field","terminator":";","layout":["properties","type","name"]}]}`)

	if _, err := LoadSpecificationSet(properties, declarations); err == nil {
		t.Fatal("expected cross-specification validation error")
	}
}
