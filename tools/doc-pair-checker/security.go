package docpairchecker

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func validateBaseCommit(value string) error {
	if len(value) != 40 && len(value) != 64 {
		return fmt.Errorf("base ref must be a full commit object id")
	}
	for _, character := range value {
		if !isHexadecimal(character) {
			return fmt.Errorf("base ref must contain only hexadecimal characters")
		}
	}
	return nil
}

func isHexadecimal(character rune) bool {
	return character >= '0' && character <= '9' ||
		character >= 'a' && character <= 'f' ||
		character >= 'A' && character <= 'F'
}

func readRegisteredDocumentation(root string, relativePath string) ([]byte, error) {
	if filepath.IsAbs(relativePath) {
		return nil, fmt.Errorf("registered documentation path must be relative: %s", relativePath)
	}
	cleanRelative := filepath.Clean(filepath.FromSlash(relativePath))
	if cleanRelative == ".." || strings.HasPrefix(cleanRelative, ".."+string(filepath.Separator)) {
		return nil, fmt.Errorf("registered documentation path escapes root: %s", relativePath)
	}

	realRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		return nil, fmt.Errorf("resolve documentation root: %w", err)
	}
	realRoot, err = filepath.Abs(realRoot)
	if err != nil {
		return nil, fmt.Errorf("make documentation root absolute: %w", err)
	}

	candidate := filepath.Join(realRoot, cleanRelative)
	info, err := os.Lstat(candidate)
	if err != nil {
		return nil, err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("registered documentation must not be a symlink: %s", relativePath)
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("registered documentation must be a regular file: %s", relativePath)
	}

	realPath, err := filepath.EvalSymlinks(candidate)
	if err != nil {
		return nil, fmt.Errorf("resolve documentation path: %w", err)
	}
	inside, err := filepath.Rel(realRoot, realPath)
	if err != nil {
		return nil, fmt.Errorf("compare documentation path: %w", err)
	}
	if inside == ".." || strings.HasPrefix(inside, ".."+string(filepath.Separator)) || filepath.IsAbs(inside) {
		return nil, fmt.Errorf("registered documentation resolves outside root: %s", relativePath)
	}

	// #nosec G304 -- realPath is constrained to a registered regular file beneath the resolved repository root.
	return os.ReadFile(realPath)
}
