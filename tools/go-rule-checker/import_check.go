package gorulechecker

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"
)

const bitlangModule = "github.com/tomiya7688/Bitlang"

func checkImports(fileSet *token.FileSet, file *ast.File, path string) []finding {
	layer := sourceLayer(path)
	var findings []finding
	for _, spec := range file.Imports {
		importPath, err := strconv.Unquote(spec.Path.Value)
		if err != nil || !strings.HasPrefix(importPath, bitlangModule+"/") {
			continue
		}
		target := sourceLayer(strings.TrimPrefix(importPath, bitlangModule+"/"))
		if layer == "ui" && target == "data" {
			findings = append(findings, finding{level: "E", rule: "ARCH", path: path, line: declarationLine(fileSet, spec), message: "UI->Data direct import"})
		}
	}
	return findings
}

func sourceLayer(path string) string {
	clean := strings.ToLower(strings.ReplaceAll(path, "\\", "/"))
	parts := strings.Split(clean, "/")
	for _, part := range parts {
		switch part {
		case "ui":
			return "ui"
		case "process":
			return "process"
		case "data":
			return "data"
		}
	}
	return ""
}
