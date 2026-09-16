package gorulechecker

import (
	"fmt"
	"io"
	"sort"
)

// Run coordinates one rule-checker execution and writes all findings.
func Run(paths []string, out io.Writer) (int, error) {
	if len(paths) == 0 {
		paths = []string{"."}
	}
	config, err := loadIgnoreConfig()
	if err != nil {
		return 0, err
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
	responsibilityFindings, err := checkResponsibilities(files)
	if err != nil {
		return 0, err
	}
	findings = append(findings, responsibilityFindings...)

	filtered := findings[:0]
	for _, item := range findings {
		if !config.ignores(item) {
			filtered = append(filtered, item)
		}
	}
	findings = filtered
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
		fmt.Fprintln(out, "OK go-rules")
		return 0, nil
	}
	fmt.Fprintf(out, "NG go-rules: %d\n", len(findings))
	return 1, nil
}
