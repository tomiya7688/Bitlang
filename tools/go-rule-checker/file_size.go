package gorulechecker

import (
	"bufio"
	"os"
)

func checkFileSize(path string) ([]finding, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lines := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines++
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if lines > 600 {
		return []finding{{level: "E", rule: "FILESIZE", path: path, message: ">600 lines"}}, nil
	}
	if lines > 400 {
		return []finding{{level: "W", rule: "FILESIZE", path: path, message: ">400 lines"}}, nil
	}
	return nil, nil
}
