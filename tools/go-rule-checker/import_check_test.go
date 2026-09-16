package gorulechecker

import "testing"

func TestSourceLayer(t *testing.T) {
	cases := map[string]string{
		"internal/ui/build.go":      "ui",
		"internal/process/build.go": "process",
		"internal/data/source.go":   "data",
		"internal/bitlang/lexer.go": "",
	}
	for path, expected := range cases {
		if actual := sourceLayer(path); actual != expected {
			t.Fatalf("sourceLayer(%q) = %q, want %q", path, actual, expected)
		}
	}
}
