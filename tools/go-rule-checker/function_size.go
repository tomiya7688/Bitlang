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
			findings = append(findings, finding{
				level:   "ERROR",
				path:    path,
				line:    start,
				message: "function \"" + function.Name.Name + "\" exceeds 120 lines and requires decomposition or explicit justification",
			})
			continue
		}

		if lines > 80 {
			findings = append(findings, finding{
				level:   "WARN",
				path:    path,
				line:    start,
				message: "function \"" + function.Name.Name + "\" exceeds 80 lines; review for multiple operations",
			})
		}
	}

	return findings
}
