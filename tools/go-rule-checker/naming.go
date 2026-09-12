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
	for _, part := range splitName(base) {
		if discouragedNameParts[strings.ToLower(part)] {
			return []finding{{level: "W", rule: "NAME", path: path, message: "vague word: " + part}}
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
			return []finding{{level: "W", rule: "NAME", path: path, line: declarationLine(fileSet, name), message: "vague word: " + part}}
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
