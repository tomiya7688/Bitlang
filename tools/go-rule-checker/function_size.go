package gorulechecker

import (
	"go/ast"
	"go/token"
)

func checkFunctionSizes(fileSet *token.FileSet, file *ast.File, path string) []finding {
	var findings []finding
	for _, declaration := range file.Decls {
		function, ok := declaration.(*ast.FuncDecl)
		if !ok || function.Body == nil {
			continue
		}
		start := fileSet.Position(function.Pos()).Line
		end := fileSet.Position(function.End()).Line
		lines := end - start + 1
		if lines > 120 {
			findings = append(findings, finding{level: "E", rule: "SIZE", path: path, line: start, message: ">120: " + function.Name.Name})
			continue
		}
		if lines > 80 {
			findings = append(findings, finding{level: "W", rule: "SIZE", path: path, line: start, message: ">80: " + function.Name.Name})
		}
	}
	return findings
}
