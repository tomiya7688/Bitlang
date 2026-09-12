package gorulechecker

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Run scans Go source files below the supplied paths and writes findings.
func Run(paths []string, out io.Writer) (int, error) {
	if len(paths) == 0 {
		paths = []string{"."}
	}

	files, err := collectGoFiles(paths)
	if err != nil {
		return 0, err
	}

	var findings []finding
	for _, path := range files {
		fileFindings, err := inspectFile(path)
		if err != nil {
			return 0, err
		}
		findings = append(findings, fileFindings...)
	}

	sort.Slice(findings, func(i, j int) bool {
		if findings[i].path != findings[j].path {
			return findings[i].path < findings[j].path
		}
		return findings[i].line < findings[j].line
	})

	for _, item := range findings {
		fmt.Fprintln(out, item.String())
	}

	if len(findings) == 0 {
		fmt.Fprintln(out, "go-rule-checker: no findings")
		return 0, nil
	}

	fmt.Fprintf(out, "go-rule-checker: %d finding(s)\n", len(findings))
	return 1, nil
}

func collectGoFiles(paths []string) ([]string, error) {
	seen := map[string]bool{}
	var files []string

	for _, root := range paths {
		info, err := os.Stat(root)
		if err != nil {
			return nil, err
		}

		if !info.IsDir() {
			if strings.HasSuffix(root, ".go") {
				files = appendUnique(files, seen, root)
			}
			continue
		}

		err = filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if entry.IsDir() && shouldSkipDir(entry.Name()) {
				return filepath.SkipDir
			}
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go") {
				files = appendUnique(files, seen, path)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	sort.Strings(files)
	return files, nil
}

func appendUnique(files []string, seen map[string]bool, path string) []string {
	clean := filepath.Clean(path)
	if seen[clean] {
		return files
	}
	seen[clean] = true
	return append(files, clean)
}

func shouldSkipDir(name string) bool {
	switch name {
	case ".git", "vendor", ".idea", ".vscode":
		return true
	default:
		return false
	}
}

func inspectFile(path string) ([]finding, error) {
	fileSet := token.NewFileSet()
	parsed, err := parser.ParseFile(fileSet, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var findings []finding
	findings = append(findings, checkFilename(path)...)
	findings = append(findings, checkDeclarations(fileSet, parsed, path)...)
	findings = append(findings, checkFunctionSizes(fileSet, parsed, path)...)
	findings = append(findings, checkMainFile(fileSet, parsed, path)...)
	return findings, nil
}

func declarationLine(fileSet *token.FileSet, node ast.Node) int {
	return fileSet.Position(node.Pos()).Line
}
