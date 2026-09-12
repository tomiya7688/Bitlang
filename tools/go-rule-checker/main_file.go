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
		if !ok {
			continue
		}
		if function.Name.Name == "main" {
			continue
		}
		findings = append(findings, finding{
			level:   "WARN",
			path:    path,
			line:    declarationLine(fileSet, function.Name),
			message: "MAIN extra function: " + function.Name.Name,
		})
	}
	return findings
}
