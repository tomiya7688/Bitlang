package gorulechecker

import (
	"go/parser"
	"go/token"
	"testing"
)

func TestCheckDeclarationsFindsGenericNameAndMissingComment(t *testing.T) {
	const source = `package sample

type DataManager struct{}
`

	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, "sample.go", source, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	findings := checkDeclarations(fileSet, file, "sample.go")
	if len(findings) != 2 {
		t.Fatalf("expected 2 findings, got %d: %v", len(findings), findings)
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
