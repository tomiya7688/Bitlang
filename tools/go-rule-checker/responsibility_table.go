package gorulechecker

import (
	"bufio"
	"os"
	"strings"
)

type responsibilityRecord struct {
	path           string
	responsibility string
}

func loadResponsibilityTable(path string) ([]responsibilityRecord, error) {
	// #nosec G304 -- path is the repository-owned FILE_RESPONSIBILITIES.md selected by the checker.
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var records []responsibilityRecord
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "| `") {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) < 4 {
			continue
		}
		filePath := strings.Trim(strings.TrimSpace(parts[1]), "`")
		responsibility := strings.TrimSpace(parts[2])
		if strings.HasSuffix(filePath, ".go") {
			records = append(records, responsibilityRecord{path: filePath, responsibility: responsibility})
		}
	}
	return records, scanner.Err()
}
