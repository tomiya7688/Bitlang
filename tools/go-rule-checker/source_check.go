package gorulechecker

import (
	"go/ast"
	"go/parser"
	"go/token"
)

func inspectFile(path string) ([]finding, error) {
	fileFindings, err := checkFileSize(path)
	if err != nil {
		return nil, err
	}
	fileSet := token.NewFileSet()
	parsed, err := parser.ParseFile(fileSet, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	findings := fileFindings
	findings = append(findings, checkFilename(path)...)
	findings = append(findings, checkNames(fileSet, parsed, path)...)
	findings = append(findings, checkDocumentation(fileSet, parsed, path)...)
	findings = append(findings, checkFunctionSizes(fileSet, parsed, path)...)
	findings = append(findings, checkMainFile(fileSet, parsed, path)...)
	return findings, nil
}

func declarationLine(fileSet *token.FileSet, node ast.Node) int {
	return fileSet.Position(node.Pos()).Line
}
