package gorulechecker

import (
	"go/ast"
	"go/token"
	"strings"
)

func checkRoleResponsibilities(fileSet *token.FileSet, file *ast.File, path string) []finding {
	role := fileRole(path)
	if role == "" {
		return nil
	}
	var findings []finding
	for _, spec := range file.Imports {
		importPath := strings.Trim(spec.Path.Value, "\"")
		if role == "commander" && isDirectIOImport(importPath) {
			findings = append(findings, finding{level: "W", rule: "ARCH", path: path, line: declarationLine(fileSet, spec), message: "Commander imports direct I/O"})
		}
		if role == "messenger" && isTransformImport(importPath) {
			findings = append(findings, finding{level: "W", rule: "ARCH", path: path, line: declarationLine(fileSet, spec), message: "Messenger imports transform package"})
		}
	}
	return findings
}

func fileRole(path string) string {
	name := strings.ToLower(strings.ReplaceAll(path, "\\", "/"))
	switch {
	case strings.Contains(name, "commander"):
		return "commander"
	case strings.Contains(name, "messenger"):
		return "messenger"
	case strings.Contains(name, "processing"):
		return "processing"
	default:
		return ""
	}
}

func isDirectIOImport(path string) bool {
	switch path {
	case "os", "io", "io/fs", "bufio":
		return true
	default:
		return false
	}
}

func isTransformImport(path string) bool {
	return strings.Contains(strings.ToLower(path), "parser") || strings.Contains(strings.ToLower(path), "compiler") || strings.Contains(strings.ToLower(path), "lexer")
}
