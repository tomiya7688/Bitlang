package bitlang

import "testing"

func TestNewExplicitArtifactAppliesSourceDefaults(t *testing.T) {
	specs := declarationParserTestSpecifications()
	specs.Properties.Axes[0].SourceDefault = "Private"
	specs.Properties.Axes[1].SourceDefault = "unnullable"

	artifact, err := NewExplicitArtifact(specs, NewSourceText("sample.bit", "Int count;"), nil)
	if err != nil {
		t.Fatal(err)
	}
	declarations := artifact.Declarations()
	if len(declarations) != 1 {
		t.Fatalf("declarations = %#v, want one declaration", declarations)
	}
	want := []PreprocessedProperty{"Private", "unnullable"}
	if len(declarations[0].Properties) != len(want) {
		t.Fatalf("properties = %#v, want %#v", declarations[0].Properties, want)
	}
	for index := range want {
		if declarations[0].Properties[index] != want[index] {
			t.Fatalf("properties = %#v, want %#v", declarations[0].Properties, want)
		}
	}
}

func TestNewExplicitArtifactPreservesExplicitPropertiesAndRejectsUnresolvedAxes(t *testing.T) {
	specs := declarationParserTestSpecifications()
	specs.Properties.Axes[1].SourceDefault = "unnullable"

	artifact, err := NewExplicitArtifact(specs, NewSourceText("sample.bit", "pUbLiC Int count;"), nil)
	if err != nil {
		t.Fatal(err)
	}
	properties := artifact.Declarations()[0].Properties
	if len(properties) != 2 || properties[0] != "Public" || properties[1] != "unnullable" {
		t.Fatalf("properties = %#v", properties)
	}

	if _, err := NewExplicitArtifact(declarationParserTestSpecifications(), NewSourceText("sample.bit", "Int count;"), nil); err == nil {
		t.Fatal("expected unresolved required property axes to fail")
	}
}
