package bitlang

import "testing"

func TestCanonicalizeIdentifier(t *testing.T) {
	got, err := CanonicalizeIdentifier("PlayerHP")
	if err != nil {
		t.Fatal(err)
	}
	if got != "playerhp" {
		t.Fatalf("got %q, want %q", got, "playerhp")
	}
}
