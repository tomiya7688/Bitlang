package bitlang

import "testing"

func TestLoadDeclarationSpecification(t *testing.T) {
	data := []byte(`{"version":1,"kinds":[{"name":"variable","property_target":"variable","terminator":";","layout":["properties","type","name"]}]}`)
	spec, err := LoadDeclarationSpecification(data)
	if err != nil {
		t.Fatal(err)
	}
	kind, err := spec.DeclarationKind("variable")
	if err != nil {
		t.Fatal(err)
	}
	if kind.Terminator != ";" || kind.PropertyTarget != "variable" {
		t.Fatalf("unexpected declaration kind: %#v", kind)
	}
}

func TestLoadDeclarationSpecificationAcceptsReorderedLayout(t *testing.T) {
	data := []byte(`{"version":1,"kinds":[{"name":"variable","property_target":"variable","terminator":";","layout":["type","name","properties"]}]}`)
	if _, err := LoadDeclarationSpecification(data); err != nil {
		t.Fatal(err)
	}
}

func TestLoadDeclarationSpecificationRejectsUnknownLayout(t *testing.T) {
	data := []byte(`{"version":1,"kinds":[{"name":"variable","property_target":"variable","terminator":";","layout":["name","type"]}]}`)
	if _, err := LoadDeclarationSpecification(data); err == nil {
		t.Fatal("expected unsupported layout error")
	}
}

func TestDeclarationKindRejectsUnknownName(t *testing.T) {
	spec := DeclarationSpecification{Version: 1, Kinds: []DeclarationKindSpec{{Name: "variable"}}}
	if _, err := spec.DeclarationKind("field"); err == nil {
		t.Fatal("expected unknown declaration kind error")
	}
}


func TestLoadDeclarationSpecificationAcceptsContexts(t *testing.T) {
	data := []byte(`{"version":1,"kinds":[{"name":"field","property_target":"field","terminator":";","layout":["properties","type","name"],"contexts":["member"]}]}`)
	spec, err := LoadDeclarationSpecification(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(spec.Kinds[0].Contexts) != 1 || spec.Kinds[0].Contexts[0] != "member" {
		t.Fatalf("unexpected contexts: %#v", spec.Kinds[0].Contexts)
	}
}

func TestLoadDeclarationSpecificationRejectsDuplicateContexts(t *testing.T) {
	data := []byte(`{"version":1,"kinds":[{"name":"field","property_target":"field","terminator":";","layout":["properties","type","name"],"contexts":["member","MeMbEr"]}]}`)
	if _, err := LoadDeclarationSpecification(data); err == nil {
		t.Fatal("expected duplicate context error")
	}
}
