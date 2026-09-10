package bitlang

import (
	"errors"
	"testing"
)

func TestCanonicalizeIdentifier(t *testing.T) {
	got, err := CanonicalizeIdentifier("PlayerHP")
	if err != nil {
		t.Fatal(err)
	}
	if got != "playerhp" {
		t.Fatalf("got %q, want %q", got, "playerhp")
	}
}

func TestSymbolTableIsCaseInsensitive(t *testing.T) {
	table := NewSymbolTable()
	if _, err := table.Define("PlayerHP", 100); err != nil {
		t.Fatal(err)
	}

	got, ok := table.Get("PLAYERhp")
	if !ok {
		t.Fatal("symbol not found through case-insensitive lookup")
	}
	if got != 100 {
		t.Fatalf("got %v, want 100", got)
	}
}

func TestSymbolTableRejectsCaseCollision(t *testing.T) {
	table := NewSymbolTable()
	if _, err := table.Define("PlayerHP", 100); err != nil {
		t.Fatal(err)
	}

	_, err := table.Define("playerhp", 200)
	if err == nil {
		t.Fatal("expected duplicate symbol error")
	}
	var duplicate DuplicateSymbolError
	if !errors.As(err, &duplicate) {
		t.Fatalf("unexpected error type: %T", err)
	}
}
