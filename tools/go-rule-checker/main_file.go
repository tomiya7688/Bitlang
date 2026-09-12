package gorulechecker

import (
	"go/ast"
	"go/token"
	"path/filepath"
)

func checkMainFile(fileSet *token.FileSet, file *ast.File, path string) []finding {
	if filepath.Base(path) != "main.go" {
		return nil
	}
	var findings []finding
	for _, declaration := range file.Decls {
		function, ok := declaration.(*ast.FuncDecl)
		if !ok || function.Name.Name == "main" {
			continue
		}
		findings = append(findings, finding{level: "W", rule: "MAIN", path: path, line: declarationLine(fileSet, function.Name), message: "extra function: " + function.Name.Name})
	}
	return findings
}
