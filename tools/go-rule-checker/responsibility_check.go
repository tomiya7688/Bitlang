package gorulechecker

import (
	"os"
	"path/filepath"
	"strings"
)

func checkResponsibilities(files []string, fullRepository bool) ([]finding, error) {
	records, err := loadResponsibilityTable("FILE_RESPONSIBILITIES.md")
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	actual := make(map[string]bool)
	for _, file := range files {
		actual[normalizeResponsibilityPath(file)] = true
	}
	registered := make(map[string]bool)
	var findings []finding
	for _, record := range records {
		registered[record.path] = true
		if len(record.responsibility) > 180 {
			findings = append(findings, finding{level: "W", rule: "RESP", path: record.path, message: "broad responsibility"})
		}
		if fullRepository && !actual[record.path] {
			if _, statErr := os.Stat(record.path); os.IsNotExist(statErr) {
				findings = append(findings, finding{level: "W", rule: "RESP", path: record.path, message: "stale entry"})
			}
		}
	}
	for path := range actual {
		if !registered[path] {
			findings = append(findings, finding{level: "W", rule: "RESP", path: path, message: "missing entry"})
		}
	}
	return findings, nil
}

func normalizeResponsibilityPath(file string) string {
	return filepath.ToSlash(strings.TrimPrefix(filepath.Clean(file), "./"))
}

func isFullRepositoryCheck(paths []string) bool {
	return len(paths) == 1 && filepath.Clean(paths[0]) == "."
}
