package gorulechecker

import "testing"

func TestFileRole(t *testing.T) {
	cases := map[string]string{
		"internal/ui/build_commander.go":         "commander",
		"internal/process/compiler_messenger.go": "messenger",
		"internal/data/source_processing.go":     "processing",
		"internal/bitlang/lexer.go":              "",
	}
	for path, expected := range cases {
		if actual := fileRole(path); actual != expected {
			t.Fatalf("fileRole(%q) = %q, want %q", path, actual, expected)
		}
	}
}
