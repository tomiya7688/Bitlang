package gorulechecker

import (
	"go/ast"
	"go/token"
	"strings"
)

func checkDocumentation(fileSet *token.FileSet, file *ast.File, path string) []finding {
	var findings []finding

	for _, declaration := range file.Decls {
		switch node := declaration.(type) {
		case *ast.FuncDecl:
			if node.Name.IsExported() && node.Doc == nil && !isTestEntrypoint(path, node.Name.Name) {
				findings = append(findings, missingComment(fileSet, node.Name, path))
			}
		case *ast.GenDecl:
			for _, spec := range node.Specs {
				switch item := spec.(type) {
				case *ast.TypeSpec:
					if item.Name.IsExported() && node.Doc == nil && item.Doc == nil {
						findings = append(findings, missingComment(fileSet, item.Name, path))
					}
				case *ast.ValueSpec:
					for _, name := range item.Names {
						if name.IsExported() && node.Doc == nil && item.Doc == nil {
							findings = append(findings, missingComment(fileSet, name, path))
						}
					}
				}
			}
		}
	}

	return findings
}

func missingComment(fileSet *token.FileSet, name *ast.Ident, path string) finding {
	return finding{
		level:   "WARN",
		path:    path,
		line:    declarationLine(fileSet, name),
		message: "exported identifier \"" + name.Name + "\" has no documentation comment",
	}
}

func isTestEntrypoint(path string, name string) bool {
	if !strings.HasSuffix(path, "_test.go") {
		return false
	}
	return strings.HasPrefix(name, "Test") || strings.HasPrefix(name, "Benchmark") || strings.HasPrefix(name, "Fuzz") || name == "Example"
}
