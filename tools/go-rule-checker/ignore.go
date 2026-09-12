package gorulechecker

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ignoreRule struct {
	rule    string
	pattern string
}

type ignoreConfig struct {
	paths []string
	rules []ignoreRule
}

func loadIgnoreConfig(path string) (ignoreConfig, error) {
	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return ignoreConfig{}, nil
	}
	if err != nil {
		return ignoreConfig{}, err
	}
	defer file.Close()

	var config ignoreConfig
	scanner := bufio.NewScanner(file)
	line := 0
	for scanner.Scan() {
		line++
		fields := strings.Fields(strings.TrimSpace(scanner.Text()))
		if len(fields) == 0 || strings.HasPrefix(fields[0], "#") {
			continue
		}
		switch {
		case fields[0] == "path" && len(fields) == 2:
			config.paths = append(config.paths, fields[1])
		case fields[0] == "rule" && len(fields) == 3:
			config.rules = append(config.rules, ignoreRule{rule: fields[1], pattern: fields[2]})
		default:
			return ignoreConfig{}, fmt.Errorf("%s:%d invalid ignore record", path, line)
		}
	}
	return config, scanner.Err()
}

func (c ignoreConfig) ignores(item finding) bool {
	for _, pattern := range c.paths {
		if matchIgnorePath(pattern, item.path) {
			return true
		}
	}
	for _, rule := range c.rules {
		if rule.rule == item.rule && matchIgnorePath(rule.pattern, item.path) {
			return true
		}
	}
	return false
}

func matchIgnorePath(pattern string, path string) bool {
	cleanPattern := filepath.ToSlash(filepath.Clean(pattern))
	cleanPath := filepath.ToSlash(filepath.Clean(path))
	if strings.HasSuffix(cleanPattern, "/**") {
		prefix := strings.TrimSuffix(cleanPattern, "**")
		return strings.HasPrefix(cleanPath, prefix)
	}
	matched, _ := filepath.Match(cleanPattern, cleanPath)
	return matched
}
