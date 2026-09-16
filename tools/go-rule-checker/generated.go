package gorulechecker

import (
	"bufio"
	"os"
	"strings"
)

func isGeneratedFile(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for line := 0; line < 20 && scanner.Scan(); line++ {
		text := scanner.Text()
		if strings.HasPrefix(text, "// Code generated ") && strings.HasSuffix(text, " DO NOT EDIT.") {
			return true, nil
		}
	}
	return false, scanner.Err()
}
