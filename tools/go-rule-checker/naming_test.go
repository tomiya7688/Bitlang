package gorulechecker

import (
	"go/parser"
	"go/token"
	"testing"
)

func TestCheckNamesFindsGenericName(t *testing.T) {
	const source = `package sample

type DataManager struct{}
`

	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, "sample.go", source, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	findings := checkNames(fileSet, file, "sample.go")
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d: %v", len(findings), findings)
	}
}

func TestSplitNameHandlesSnakeAndCamelCase(t *testing.T) {
	parts := splitName("project_file_manager")
	if len(parts) != 3 || parts[2] != "manager" {
		t.Fatalf("unexpected snake-case parts: %v", parts)
	}

	parts = splitName("ProjectFileManager")
	if len(parts) != 3 || parts[2] != "Manager" {
		t.Fatalf("unexpected camel-case parts: %v", parts)
	}
}
