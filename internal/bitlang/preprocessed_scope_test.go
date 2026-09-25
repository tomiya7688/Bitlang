package bitlang

import "testing"

func TestNewPreprocessedScopeCanonicalizesContext(t *testing.T) {
	scope, err := NewPreprocessedScope("MeMbEr", nil)
	if err != nil {
		t.Fatal(err)
	}
	if scope.Context.Spelling != "MeMbEr" || scope.Context.Canonical != "member" {
		t.Fatalf("unexpected context: %#v", scope.Context)
	}
}

func TestNewPreprocessedScopeKeepsParent(t *testing.T) {
	parent, err := NewPreprocessedScope("outer", nil)
	if err != nil {
		t.Fatal(err)
	}
	child, err := NewPreprocessedScope("member", &parent)
	if err != nil {
		t.Fatal(err)
	}
	if child.Parent == nil || child.Parent.Context.Canonical != "outer" {
		t.Fatalf("unexpected parent scope: %#v", child.Parent)
	}
}
