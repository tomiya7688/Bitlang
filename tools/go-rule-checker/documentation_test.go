package gorulechecker

import (
	"go/parser"
	"go/token"
	"testing"
)

func TestCheckDocumentationFindsMissingExportedComment(t *testing.T) {
	const source = `package sample

type PublicType struct{}
`

	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, "sample.go", source, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	findings := checkDocumentation(fileSet, file, "sample.go")
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d: %v", len(findings), findings)
	}
}

func TestCheckDocumentationIgnoresGoTestEntrypoint(t *testing.T) {
	const source = `package sample

func TestFeature() {}
`

	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, "sample_test.go", source, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	findings := checkDocumentation(fileSet, file, "sample_test.go")
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %v", findings)
	}
}
