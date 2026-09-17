package bitlang

import "testing"

func TestLoadPropertySpecification(t *testing.T) {
	data := []byte(`{"version":1,"axes":[{"name":"nullability","states":["nullable","unnullable"],"exclusive":true}]}`)
	spec, err := LoadPropertySpecification(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(spec.Axes) != 1 || spec.Axes[0].States[1] != "unnullable" {
		t.Fatalf("unexpected specification: %#v", spec)
	}
}

func TestLoadPropertySpecificationRejectsIncompleteAxis(t *testing.T) {
	data := []byte(`{"version":1,"axes":[{"name":"nullability","states":[]}]}`)
	if _, err := LoadPropertySpecification(data); err == nil {
		t.Fatal("expected invalid axis error")
	}
}
