package gorulechecker

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"
	"unicode"
)

var discouragedNameParts = map[string]bool{
	"manager": true,
	"helper":  true,
	"helpers": true,
	"util":    true,
	"utils":   true,
	"common":  true,
}

func checkFilename(path string) []finding {
	base := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	parts := splitName(base)
	for _, part := range parts {
		if discouragedNameParts[strings.ToLower(part)] {
			return []finding{{
				level:   "WARN",
				path:    path,
				message: "generic filename part \"" + part + "\" hides responsibility",
			}}
		}
	}
	return nil
}

func checkNames(fileSet *token.FileSet, file *ast.File, path string) []finding {
	var findings []finding

	for _, declaration := range file.Decls {
		switch node := declaration.(type) {
		case *ast.FuncDecl:
			findings = append(findings, checkName(fileSet, node.Name, path)...)
		case *ast.GenDecl:
			for _, spec := range node.Specs {
				switch item := spec.(type) {
				case *ast.TypeSpec:
					findings = append(findings, checkName(fileSet, item.Name, path)...)
				case *ast.ValueSpec:
					for _, name := range item.Names {
						findings = append(findings, checkName(fileSet, name, path)...)
					}
				}
			}
		}
	}

	return findings
}

func checkName(fileSet *token.FileSet, name *ast.Ident, path string) []finding {
	for _, part := range splitName(name.Name) {
		if discouragedNameParts[strings.ToLower(part)] {
			return []finding{{
				level:   "WARN",
				path:    path,
				line:    declarationLine(fileSet, name),
				message: "generic identifier part \"" + part + "\" should be replaced by a responsibility-specific name",
			}}
		}
	}
	return nil
}

func splitName(name string) []string {
	name = strings.ReplaceAll(name, "-", "_")
	var parts []string
	for _, chunk := range strings.Split(name, "_") {
		if chunk == "" {
			continue
		}
		start := 0
		runes := []rune(chunk)
		for index := 1; index < len(runes); index++ {
			if unicode.IsUpper(runes[index]) && (unicode.IsLower(runes[index-1]) || unicode.IsDigit(runes[index-1])) {
				parts = append(parts, string(runes[start:index]))
				start = index
			}
		}
		parts = append(parts, string(runes[start:]))
	}
	return parts
}
