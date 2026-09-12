package gorulechecker

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

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
